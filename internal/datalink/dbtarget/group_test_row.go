package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// TestRowOwnerPrefix namespaces the owner value of operation-owned test rows.
// A row is only ever read or deleted by the exact owner value of one test
// operation, and a value outside this namespace is refused, so test cleanup can
// never reach a production row.
const TestRowOwnerPrefix = "gw-test-"

// maxOwnedRowsRead bounds how many rows a readback fetches; one operation owns
// exactly one row, so more than one is itself a finding.
const maxOwnedRowsRead = 3

var (
	// ErrOwnedRowRefInvalid marks a reference that cannot prove ownership.
	ErrOwnedRowRefInvalid = errors.New("owned test row reference invalid")
	// ErrOwnedCleanupAmbiguous means the cleanup commit failed in a way that
	// does not say whether the rows were removed.
	ErrOwnedCleanupAmbiguous = errors.New("owned test row cleanup outcome unknown")
)

// OwnedRowRef names the rows one test operation owns: those whose OwnerColumn
// holds exactly OwnerValue.
type OwnedRowRef struct {
	Kind        schema.DatabaseConnectorKind
	SchemaName  string
	TableName   string
	OwnerColumn string
	OwnerValue  string
}

func (r OwnedRowRef) validate() error {
	switch {
	case r.Kind != schema.DatabaseConnectorKindSQLite && r.Kind != schema.DatabaseConnectorKindPostgres:
		return fmt.Errorf("%w: destination kind", ErrOwnedRowRefInvalid)
	case strings.TrimSpace(r.TableName) == "" || strings.TrimSpace(r.OwnerColumn) == "":
		return fmt.Errorf("%w: table and owner column are required", ErrOwnedRowRefInvalid)
	case !strings.HasPrefix(r.OwnerValue, TestRowOwnerPrefix) || len(r.OwnerValue) == len(TestRowOwnerPrefix):
		return fmt.Errorf("%w: owner value is outside the test namespace", ErrOwnedRowRefInvalid)
	}
	return nil
}

// ReadOwnedRows returns the raw values of the requested columns for the rows
// the operation owns, in column order. It never writes.
func ReadOwnedRows(ctx context.Context, db *sql.DB, ref OwnedRowRef, columns []string) ([][]any, error) {
	if err := ref.validate(); err != nil {
		return nil, err
	}
	if db == nil || len(columns) == 0 {
		return nil, fmt.Errorf("%w: connection and columns are required", ErrOwnedRowRefInvalid)
	}
	rows, err := db.QueryContext(ctx, ownedRowsSelect(ref, columns), ref.OwnerValue)
	if err != nil {
		return nil, fmt.Errorf("read owned test rows: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out [][]any
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range dest {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("scan owned test row: %w", err)
		}
		for i, value := range values {
			if raw, ok := value.([]byte); ok {
				values[i] = append([]byte(nil), raw...)
			}
		}
		out = append(out, values)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read owned test rows: %w", err)
	}
	return out, nil
}

// RemoveOwnedRows deletes only the rows the operation owns, and the effect
// receipt of its test effect when one was written, in one transaction. It
// returns how many table rows were removed. effectKey is set only when the
// operation wrote a receipt, so a missing receipt table is a real failure.
func RemoveOwnedRows(ctx context.Context, db *sql.DB, ref OwnedRowRef, effectKey string) (int64, error) {
	if err := ref.validate(); err != nil {
		return 0, err
	}
	if db == nil {
		return 0, fmt.Errorf("%w: connection is required", ErrOwnedRowRefInvalid)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin owned test row cleanup: %w", err)
	}
	res, err := tx.ExecContext(ctx, ownedRowsDelete(ref), ref.OwnerValue)
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("delete owned test rows: %w", err))
	}
	removed, err := res.RowsAffected()
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("read deleted test rows: %w", err))
	}
	if strings.TrimSpace(effectKey) != "" {
		receipts := qualifiedTableName(ref.Kind, ref.SchemaName, EffectReceiptTable)
		if _, err := tx.ExecContext(ctx, receiptDeleteStatement(ref.Kind, receipts), effectKey); err != nil {
			return 0, rollbackWith(tx, fmt.Errorf("delete owned test receipt: %w", err))
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOwnedCleanupAmbiguous, err)
	}
	return removed, nil
}

// ReadEffectReceiptDigest reports the payload digest the destination holds for
// an effect key, if any.
func ReadEffectReceiptDigest(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName, effectKey string) (digest string, found bool, err error) {
	if db == nil || strings.TrimSpace(effectKey) == "" {
		return "", false, fmt.Errorf("%w: connection and effect key are required", ErrOwnedRowRefInvalid)
	}
	receipts := qualifiedTableName(kind, schemaName, EffectReceiptTable)
	err = db.QueryRowContext(ctx, receiptLookupStatement(kind, receipts), effectKey).Scan(&digest)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", false, nil
	case err != nil:
		return "", false, fmt.Errorf("read destination receipt: %w", err)
	}
	return digest, true, nil
}

func rollbackWith(tx *sql.Tx, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
		return errors.Join(err, rollbackErr)
	}
	return err
}

func ownedRowsSelect(ref OwnedRowRef, columns []string) string {
	quoted := make([]string, len(columns))
	for i, column := range columns {
		quoted[i] = quoteIdentifier(ref.Kind, column)
	}
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s LIMIT %d", strings.Join(quoted, ", "),
		qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName), quoteIdentifier(ref.Kind, ref.OwnerColumn),
		buildPlaceholders(ref.Kind, 1)[0], maxOwnedRowsRead)
}

func ownedRowsDelete(ref OwnedRowRef) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName), quoteIdentifier(ref.Kind, ref.OwnerColumn),
		buildPlaceholders(ref.Kind, 1)[0])
}

func receiptDeleteStatement(kind schema.DatabaseConnectorKind, receipts string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE effect_key = %s", receipts, buildPlaceholders(kind, 1)[0])
}
