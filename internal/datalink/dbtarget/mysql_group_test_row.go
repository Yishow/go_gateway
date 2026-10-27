package dbtarget

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// MaxGroupLongTestRows bounds both preview and cleanup identity arrays.
const MaxGroupLongTestRows = 256

var groupLongRecordIDPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

// GroupLongTestRowRef names exact operation-derived rows in the shared long
// table. The owner column remains the immutable table ownership marker; the
// operation owner, workspace/group/revision, sample IDs, record IDs and receipt
// identities together prove the row belongs to this test operation.
type GroupLongTestRowRef struct {
	Kind             schema.DatabaseConnectorKind
	SchemaName       string
	TableName        string
	OwnerColumn      string
	OwnerValue       string
	ProvenanceColumn string
	WorkspaceID      string
	GroupID          string
	GroupRevision    string
	OperationOwner   string
	RecordIDs        []string
	SampleIDs        []string
	EffectKeys       []string
	EffectDigests    map[string]string
}

func (r GroupLongTestRowRef) validate() error {
	if (r.Kind != schema.DatabaseConnectorKindMySQL && r.Kind != schema.DatabaseConnectorKindSQLServer) || strings.TrimSpace(r.SchemaName) == "" || strings.TrimSpace(r.TableName) == "" ||
		strings.TrimSpace(r.OwnerColumn) == "" || r.OwnerValue != "gateway" ||
		strings.TrimSpace(r.ProvenanceColumn) == "" || strings.TrimSpace(r.WorkspaceID) == "" ||
		strings.TrimSpace(r.GroupID) == "" || strings.TrimSpace(r.GroupRevision) == "" ||
		!strings.HasPrefix(r.OperationOwner, TestRowOwnerPrefix) || len(r.OperationOwner) == len(TestRowOwnerPrefix) {
		return fmt.Errorf("%w: long test operation scope", ErrOwnedRowRefInvalid)
	}
	if len(r.RecordIDs) == 0 || len(r.RecordIDs) > MaxGroupLongTestRows ||
		len(r.RecordIDs) != len(r.SampleIDs) || len(r.RecordIDs) != len(r.EffectKeys) ||
		len(r.EffectDigests) != len(r.EffectKeys) {
		return fmt.Errorf("%w: long test row proof", ErrOwnedRowRefInvalid)
	}
	samples := make(map[string]struct{}, len(r.SampleIDs))
	records := make(map[string]struct{}, len(r.RecordIDs))
	for _, id := range r.RecordIDs {
		if !groupLongRecordIDPattern.MatchString(id) {
			return fmt.Errorf("%w: long test record ID", ErrOwnedRowRefInvalid)
		}
		if _, duplicate := records[id]; duplicate {
			return fmt.Errorf("%w: duplicate long test record", ErrOwnedRowRefInvalid)
		}
		records[id] = struct{}{}
	}
	for i, sampleID := range r.SampleIDs {
		if !strings.HasPrefix(sampleID, r.OperationOwner+"-") || strings.TrimSpace(sampleID) == "" {
			return fmt.Errorf("%w: long test sample proof", ErrOwnedRowRefInvalid)
		}
		if _, duplicate := samples[sampleID]; duplicate {
			return fmt.Errorf("%w: duplicate long test sample", ErrOwnedRowRefInvalid)
		}
		samples[sampleID] = struct{}{}
		if strings.TrimSpace(r.EffectKeys[i]) == "" || strings.TrimSpace(r.EffectDigests[r.EffectKeys[i]]) == "" {
			return fmt.Errorf("%w: long test effect proof", ErrOwnedRowRefInvalid)
		}
	}
	for key := range r.EffectDigests {
		found := false
		for _, expected := range r.EffectKeys {
			if key == expected {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%w: unexpected long test effect", ErrOwnedRowRefInvalid)
		}
	}
	if r.Kind == schema.DatabaseConnectorKindSQLServer {
		return validateSQLServerGroupLongTestRowRef(r)
	}
	return nil
}

type groupLongQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// ReadGroupLongTestRows reads exact record IDs only; it never reads by owner
// alone and fails closed when any returned row lacks this operation's proof.
func ReadGroupLongTestRows(ctx context.Context, db *sql.DB, ref GroupLongTestRowRef, columns []string) ([][]any, error) {
	if err := ref.validate(); err != nil {
		return nil, err
	}
	if db == nil || len(columns) == 0 {
		return nil, fmt.Errorf("%w: connection and columns are required", ErrOwnedRowRefInvalid)
	}
	return readGroupLongTestRows(ctx, db, ref, columns)
}

func readGroupLongTestRows(ctx context.Context, q groupLongQueryer, ref GroupLongTestRowRef, columns []string) ([][]any, error) {
	if ref.Kind == schema.DatabaseConnectorKindSQLServer {
		return readSQLServerGroupLongTestRows(ctx, q, ref, columns, false)
	}
	selected := append([]string{}, columns...)
	for _, mandatory := range []string{longColRecordID, longColWorkspaceID, longColGroupID, longColGroupRev, longColProvenance, ref.OwnerColumn} {
		found := false
		for _, column := range selected {
			if strings.EqualFold(column, mandatory) {
				found = true
				break
			}
		}
		if !found {
			selected = append(selected, mandatory)
		}
	}
	quoted := make([]string, len(selected))
	for i, column := range selected {
		quoted[i] = quoteIdentifier(ref.Kind, column)
	}
	args := make([]any, 0, len(ref.RecordIDs)+1)
	args = append(args, ref.OwnerValue)
	for _, id := range ref.RecordIDs {
		args = append(args, id)
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? AND %s IN (%s) LIMIT %d",
		strings.Join(quoted, ", "), qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName),
		quoteIdentifier(ref.Kind, ref.OwnerColumn), quoteIdentifier(ref.Kind, longColRecordID),
		strings.TrimSuffix(strings.Repeat("?,", len(ref.RecordIDs)), ","), MaxGroupLongTestRows)
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("read long test rows: %w", err)
	}
	defer func() { _ = rows.Close() }()
	index := make(map[string]int, len(selected))
	for i, column := range selected {
		index[strings.ToLower(column)] = i
	}
	requested := make([]int, len(columns))
	for i, column := range columns {
		requested[i] = index[strings.ToLower(column)]
	}
	expected := make(map[string]string, len(ref.RecordIDs))
	for i, id := range ref.RecordIDs {
		expected[id] = ref.SampleIDs[i]
	}
	var out [][]any
	for rows.Next() {
		values := make([]any, len(selected))
		dest := make([]any, len(selected))
		for i := range dest {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("scan long test row: %w", err)
		}
		if err := validateGroupLongTestRow(values, index, expected, ref); err != nil {
			return nil, err
		}
		row := make([]any, len(requested))
		for i, at := range requested {
			row[i] = normalizeGroupLongValue(values[at])
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read long test rows: %w", err)
	}
	return out, nil
}

func validateGroupLongTestRow(values []any, index map[string]int, expected map[string]string, ref GroupLongTestRowRef) error {
	recordID, ok := groupLongText(values[index[longColRecordID]])
	if !ok {
		return fmt.Errorf("%w: long test record ID", ErrOwnedRowRefInvalid)
	}
	sampleID, ok := expected[recordID]
	if !ok {
		return fmt.Errorf("%w: unowned long test record", ErrOwnedRowRefInvalid)
	}
	if workspace, ok := groupLongText(values[index[longColWorkspaceID]]); !ok || workspace != ref.WorkspaceID {
		return fmt.Errorf("%w: long test workspace proof", ErrOwnedRowRefInvalid)
	}
	if group, ok := groupLongText(values[index[longColGroupID]]); !ok || group != ref.GroupID {
		return fmt.Errorf("%w: long test group proof", ErrOwnedRowRefInvalid)
	}
	if revision, ok := groupLongText(values[index[longColGroupRev]]); !ok || revision != ref.GroupRevision {
		return fmt.Errorf("%w: long test revision proof", ErrOwnedRowRefInvalid)
	}
	if owner, ok := groupLongText(values[index[strings.ToLower(ref.OwnerColumn)]]); !ok || owner != ref.OwnerValue {
		return fmt.Errorf("%w: long test table owner proof", ErrOwnedRowRefInvalid)
	}
	provenance, ok := groupLongText(values[index[longColProvenance]])
	if !ok {
		return fmt.Errorf("%w: long test provenance", ErrOwnedRowRefInvalid)
	}
	var proof struct {
		SampleID string `json:"sample_id"`
	}
	if json.Unmarshal([]byte(provenance), &proof) != nil || proof.SampleID != sampleID {
		return fmt.Errorf("%w: long test provenance proof", ErrOwnedRowRefInvalid)
	}
	return nil
}

func normalizeGroupLongValue(value any) any {
	if raw, ok := value.([]byte); ok {
		return append([]byte(nil), raw...)
	}
	return value
}

func groupLongText(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return typed, true
	case []byte:
		return string(typed), true
	}
	return "", false
}

// RemoveGroupLongTestRows deletes only the exact operation record IDs and its
// receipts. It validates row ownership proof and receipt digests before any
// deletion; a foreign row or receipt conflicts and aborts the whole cleanup.
func RemoveGroupLongTestRows(ctx context.Context, db *sql.DB, ref GroupLongTestRowRef) (int64, error) {
	if err := ref.validate(); err != nil {
		return 0, err
	}
	if db == nil {
		return 0, fmt.Errorf("%w: connection is required", ErrOwnedRowRefInvalid)
	}
	if ref.Kind == schema.DatabaseConnectorKindSQLServer {
		return removeSQLServerGroupLongTestRows(ctx, db, ref)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin long test row cleanup: %w", err)
	}
	if _, err := readGroupLongTestRows(ctx, tx, ref, []string{longColRecordID}); err != nil {
		return 0, rollbackWith(tx, err)
	}
	if err := verifyLongTestReceipts(ctx, tx, ref); err != nil {
		return 0, rollbackWith(tx, err)
	}
	rowConditions := make([]string, len(ref.RecordIDs))
	rowArgs := []any{ref.OwnerValue, ref.WorkspaceID, ref.GroupID, ref.GroupRevision}
	for i, id := range ref.RecordIDs {
		// The DELETE repeats the validated record ID, scope and sample proof.
		// HEX compares the extracted JSON string bytes exactly on both server families,
		// avoiding MariaDB JSON-versus-parameter comparison coercion. A concurrent
		// replacement cannot pass a different operation sample or scope.
		rowConditions[i] = fmt.Sprintf("(%s = ? AND HEX(JSON_EXTRACT(%s, '$.sample_id')) = HEX(JSON_EXTRACT(?, '$.sample_id')))",
			quoteIdentifier(ref.Kind, longColRecordID), quoteIdentifier(ref.Kind, ref.ProvenanceColumn))
		rowArgs = append(rowArgs, id, fmt.Sprintf(`{"sample_id":%q}`, ref.SampleIDs[i]))
	}
	// #nosec G201 -- Identifiers are quoted; generated predicates contain only bound placeholders.
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ? AND %s = ? AND %s = ? AND (%s)",
		qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName), quoteIdentifier(ref.Kind, ref.OwnerColumn),
		quoteIdentifier(ref.Kind, longColWorkspaceID), quoteIdentifier(ref.Kind, longColGroupID),
		quoteIdentifier(ref.Kind, longColGroupRev), strings.Join(rowConditions, " OR "))
	res, err := tx.ExecContext(ctx, query, rowArgs...)
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("delete long test rows: %w", err))
	}
	removed, err := res.RowsAffected()
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("read deleted long test rows: %w", err))
	}
	receipts := qualifiedTableName(ref.Kind, ref.SchemaName, EffectReceiptTable)
	receiptConditions := make([]string, len(ref.EffectKeys))
	receiptArgs := make([]any, 0, len(ref.EffectKeys)*2)
	for i, key := range ref.EffectKeys {
		receiptConditions[i] = "(effect_key = ? AND payload_digest = ?)"
		receiptArgs = append(receiptArgs, key, ref.EffectDigests[key])
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s", receipts, strings.Join(receiptConditions, " OR ")), receiptArgs...); err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("delete long test receipts: %w", err))
	}
	// A concurrent replacement after the verified read must not be committed
	// as cleaned: check both exact rows and exact receipts inside the transaction.
	if rows, err := readGroupLongTestRows(ctx, tx, ref, []string{longColRecordID}); err != nil {
		return 0, rollbackWith(tx, err)
	} else if len(rows) > 0 {
		return 0, rollbackWith(tx, fmt.Errorf("%w: long test rows remain", ErrOwnedRowRefInvalid))
	}
	var remaining int
	remainingKeys := make([]any, len(ref.EffectKeys))
	for i, key := range ref.EffectKeys {
		remainingKeys[i] = key
	}
	if err := tx.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE effect_key IN (%s)", receipts,
		strings.TrimSuffix(strings.Repeat("?,", len(remainingKeys)), ",")), remainingKeys...).Scan(&remaining); err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("read remaining long test receipts: %w", err))
	}
	if remaining > 0 {
		return 0, rollbackWith(tx, fmt.Errorf("%w: long test receipts remain", ErrOwnedRowRefInvalid))
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOwnedCleanupAmbiguous, err)
	}
	return removed, nil
}

func verifyLongTestReceipts(ctx context.Context, tx *sql.Tx, ref GroupLongTestRowRef) error {
	args := make([]any, len(ref.EffectKeys))
	for i, key := range ref.EffectKeys {
		args[i] = key
	}
	receipts := qualifiedTableName(ref.Kind, ref.SchemaName, EffectReceiptTable)
	rows, err := tx.QueryContext(ctx, fmt.Sprintf("SELECT effect_key, payload_digest FROM %s WHERE effect_key IN (%s)",
		receipts, strings.TrimSuffix(strings.Repeat("?,", len(args)), ",")), args...)
	if err != nil {
		return fmt.Errorf("read long test receipts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var key, digest string
		if err := rows.Scan(&key, &digest); err != nil {
			return fmt.Errorf("scan long test receipt: %w", err)
		}
		if want := ref.EffectDigests[key]; want != "" && want != digest {
			return fmt.Errorf("%w: long test receipt belongs to another effect", ErrOwnedRowRefInvalid)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("read long test receipts: %w", err)
	}
	return nil
}
