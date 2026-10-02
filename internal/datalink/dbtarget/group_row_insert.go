package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// Effect strategies a destination can offer for recognizing a repeated row.
const (
	GroupEffectReceipt   = "receipt"
	GroupEffectUniqueKey = "unique_key"
	GroupEffectNone      = "none"

	// EffectReceiptTable is the destination-side receipt table name.
	EffectReceiptTable = "gw_effect_receipts"
)

// Phases of a failed insert. Only the commit phase is ambiguous: the
// destination may have committed even though the caller saw an error.
const (
	InsertPhasePreCommit = "pre_commit"
	InsertPhaseCommit    = "commit"
)

var (
	// ErrReceiptDigestMismatch means the destination already holds a receipt for
	// this effect key with different content; it is never overwritten.
	ErrReceiptDigestMismatch = errors.New("destination receipt digest mismatch")
	// ErrExistingRowDiffers means a row with the same record key already exists
	// with different values.
	ErrExistingRowDiffers = errors.New("destination row with the same record key differs")
	// ErrGroupInsertUnsupported marks an unsupported kind, strategy or request.
	ErrGroupInsertUnsupported = errors.New("group row insert unsupported")
)

// GroupInsertRequest is one row to commit to a destination transactionally.
type GroupInsertRequest struct {
	Kind            schema.DatabaseConnectorKind
	SchemaName      string
	TableName       string
	Row             EncodedRow
	PayloadDigest   string
	Strategy        string
	RecordKeyColumn string
	// CommittedAt is written into the destination receipt.
	CommittedAt string
}

// GroupInsertResult reports what the destination did.
type GroupInsertResult struct {
	// AlreadyCommitted is true when the destination already held this effect.
	AlreadyCommitted bool
}

// GroupInsertError carries the phase so callers can tell a safe retry from an
// ambiguous commit. The cause is kept for callers; Error() stays value-free.
type GroupInsertError struct {
	Phase string
	Cause error
}

func (e *GroupInsertError) Error() string { return "group row insert failed in " + e.Phase }

// Unwrap exposes the underlying cause.
func (e *GroupInsertError) Unwrap() error { return e.Cause }

// Ambiguous reports whether the destination may have committed.
func (e *GroupInsertError) Ambiguous() bool { return e.Phase == InsertPhaseCommit }

// InsertGroupRow commits one row (and, by strategy, its effect evidence) in a
// single destination transaction.
//
//   - receipt: the row and a receipt (effect key and payload digest) commit
//     together; an existing receipt with the same digest means the effect is
//     already committed and nothing is inserted, a different digest is refused.
//   - unique_key: the record key column must be unique; a conflicting row with
//     identical values is the same effect, different values are refused.
//   - none: a plain insert. A failure at commit is ambiguous and the caller
//     must not retry blindly.
func InsertGroupRow(ctx context.Context, db *sql.DB, req GroupInsertRequest) (GroupInsertResult, error) {
	if err := validateGroupInsert(db, req); err != nil {
		return GroupInsertResult{}, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return GroupInsertResult{}, &GroupInsertError{Phase: InsertPhasePreCommit, Cause: err}
	}
	result, err := insertInTx(ctx, tx, req)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, rollbackErr)
		}
		return GroupInsertResult{}, preCommit(err)
	}
	if result.AlreadyCommitted {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			return GroupInsertResult{}, &GroupInsertError{Phase: InsertPhasePreCommit, Cause: err}
		}
		return result, nil
	}
	if err := tx.Commit(); err != nil {
		return GroupInsertResult{}, &GroupInsertError{Phase: InsertPhaseCommit, Cause: err}
	}
	return result, nil
}

// preCommit keeps decisive sentinels intact and wraps driver errors as safe
// pre-commit failures.
func preCommit(err error) error {
	if errors.Is(err, ErrReceiptDigestMismatch) || errors.Is(err, ErrExistingRowDiffers) {
		return err
	}
	return &GroupInsertError{Phase: InsertPhasePreCommit, Cause: err}
}

func validateGroupInsert(db *sql.DB, req GroupInsertRequest) error {
	unsupported := func(reason string) error {
		return fmt.Errorf("%w: %s", ErrGroupInsertUnsupported, reason)
	}
	switch {
	case db == nil:
		return unsupported("no destination connection")
	case req.Kind != schema.DatabaseConnectorKindSQLite && req.Kind != schema.DatabaseConnectorKindPostgres:
		return unsupported("destination kind")
	case strings.TrimSpace(req.TableName) == "":
		return unsupported("table is required")
	case req.Row.EffectKey == "" || len(req.Row.Cells) == 0:
		return unsupported("row identity and cells are required")
	}
	for _, cell := range req.Row.Cells {
		if strings.TrimSpace(cell.Column) == "" {
			return unsupported("column name is required")
		}
	}
	switch req.Strategy {
	case GroupEffectReceipt:
		if req.PayloadDigest == "" {
			return unsupported("receipt strategy needs a payload digest")
		}
	case GroupEffectUniqueKey:
		if strings.TrimSpace(req.RecordKeyColumn) == "" {
			return unsupported("unique key strategy needs the record key column")
		}
	case GroupEffectNone:
	default:
		return unsupported("effect strategy")
	}
	return nil
}

func insertInTx(ctx context.Context, tx *sql.Tx, req GroupInsertRequest) (GroupInsertResult, error) {
	table := qualifiedTableName(req.Kind, req.SchemaName, req.TableName)
	if req.Strategy == GroupEffectReceipt {
		receipts := qualifiedTableName(req.Kind, req.SchemaName, EffectReceiptTable)
		var existing string
		err := tx.QueryRowContext(ctx, receiptLookupStatement(req.Kind, receipts), req.Row.EffectKey).Scan(&existing)
		switch {
		case err == nil && existing == req.PayloadDigest:
			return GroupInsertResult{AlreadyCommitted: true}, nil
		case err == nil:
			return GroupInsertResult{}, ErrReceiptDigestMismatch
		case !errors.Is(err, sql.ErrNoRows):
			return GroupInsertResult{}, fmt.Errorf("read destination receipt: %w", err)
		}
	}

	columns := make([]string, 0, len(req.Row.Cells))
	args := make([]any, 0, len(req.Row.Cells))
	for _, cell := range req.Row.Cells {
		columns = append(columns, quoteIdentifier(req.Kind, cell.Column))
		args = append(args, normalizeDBValue(cell.Value))
	}
	res, err := tx.ExecContext(ctx, groupInsertStatement(req, table, columns), args...)
	if err != nil {
		return GroupInsertResult{}, fmt.Errorf("insert destination row: %w", err)
	}

	switch req.Strategy {
	case GroupEffectUniqueKey:
		affected, err := res.RowsAffected()
		if err != nil {
			return GroupInsertResult{}, fmt.Errorf("read affected rows: %w", err)
		}
		if affected == 0 {
			return existingRowResult(ctx, tx, req, table, columns)
		}
	case GroupEffectReceipt:
		receipts := qualifiedTableName(req.Kind, req.SchemaName, EffectReceiptTable)
		if _, err := tx.ExecContext(ctx, receiptInsertStatement(req.Kind, receipts),
			req.Row.EffectKey, req.PayloadDigest, req.CommittedAt); err != nil {
			return GroupInsertResult{}, fmt.Errorf("insert destination receipt: %w", err)
		}
	}
	return GroupInsertResult{}, nil
}

// Statement builders keep SQL text in one place; identifiers are quoted and
// every value is a bound parameter.
func receiptLookupStatement(kind schema.DatabaseConnectorKind, receipts string) string {
	return fmt.Sprintf("SELECT payload_digest FROM %s WHERE effect_key = %s", receipts, buildPlaceholders(kind, 1)[0])
}

func receiptInsertStatement(kind schema.DatabaseConnectorKind, receipts string) string {
	marks := buildPlaceholders(kind, 3)
	return fmt.Sprintf("INSERT INTO %s (effect_key, payload_digest, committed_at) VALUES (%s, %s, %s)", receipts, marks[0], marks[1], marks[2])
}

func groupInsertStatement(req GroupInsertRequest, table string, columns []string) string {
	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(columns, ", "), strings.Join(buildPlaceholders(req.Kind, len(columns)), ", "))
	if req.Strategy == GroupEffectUniqueKey {
		statement += fmt.Sprintf(" ON CONFLICT (%s) DO NOTHING", quoteIdentifier(req.Kind, req.RecordKeyColumn))
	}
	return statement
}

func existingRowStatement(req GroupInsertRequest, table string, columns []string) string {
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
		strings.Join(columns, ", "), table, quoteIdentifier(req.Kind, req.RecordKeyColumn), buildPlaceholders(req.Kind, 1)[0])
}

// existingRowResult decides whether a row that already holds the record key is
// this same effect (identical values) or a conflicting one.
func existingRowResult(ctx context.Context, tx *sql.Tx, req GroupInsertRequest, table string, columns []string) (GroupInsertResult, error) {
	var keyValue any
	for _, cell := range req.Row.Cells {
		if strings.EqualFold(cell.Column, req.RecordKeyColumn) {
			keyValue = normalizeDBValue(cell.Value)
		}
	}
	if keyValue == nil {
		return GroupInsertResult{}, fmt.Errorf("%w: row does not carry the record key column", ErrGroupInsertUnsupported)
	}
	query := existingRowStatement(req, table, columns)
	got := make([]any, len(columns))
	dest := make([]any, len(columns))
	for i := range dest {
		dest[i] = &got[i]
	}
	if err := tx.QueryRowContext(ctx, query, keyValue).Scan(dest...); err != nil {
		return GroupInsertResult{}, fmt.Errorf("read existing destination row: %w", err)
	}
	for i, cell := range req.Row.Cells {
		if !groupCellEqual(normalizeDBValue(cell.Value), got[i]) {
			return GroupInsertResult{}, ErrExistingRowDiffers
		}
	}
	return GroupInsertResult{AlreadyCommitted: true}, nil
}

// groupCellEqual compares a bound value with what the driver reads back,
// tolerating only representation differences (0/1 booleans, text bytes, text
// timestamps), never value changes.
func groupCellEqual(want, got any) bool {
	switch w := want.(type) {
	case nil:
		return got == nil
	case bool:
		switch g := got.(type) {
		case bool:
			return g == w
		case int64:
			return (g == 1) == w && (g == 0 || g == 1)
		}
	case int64:
		g, ok := got.(int64)
		return ok && g == w
	case float64:
		g, ok := got.(float64)
		return ok && g == w
	case string:
		switch g := got.(type) {
		case string:
			return g == w
		case []byte:
			return string(g) == w
		}
	case time.Time:
		switch g := got.(type) {
		case time.Time:
			return g.Equal(w)
		case string:
			for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05.999999999 -0700 MST"} {
				if parsed, err := time.Parse(layout, g); err == nil {
					return parsed.Equal(w)
				}
			}
		}
	}
	return false
}

// CreateEffectReceiptTable creates the destination receipt table. It is a
// setup helper for managed schema flows and tests; senders never run DDL.
func CreateEffectReceiptTable(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName string) error {
	if kind != schema.DatabaseConnectorKindSQLite && kind != schema.DatabaseConnectorKindPostgres {
		return fmt.Errorf("%w: destination kind", ErrGroupInsertUnsupported)
	}
	_, err := db.ExecContext(ctx, receiptTableStatement(qualifiedTableName(kind, schemaName, EffectReceiptTable)))
	if err != nil {
		return fmt.Errorf("create destination receipt table: %w", err)
	}
	return nil
}

func receiptTableStatement(table string) string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (effect_key TEXT PRIMARY KEY, payload_digest TEXT NOT NULL, committed_at TEXT NOT NULL)", table)
}
