package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func removeSQLServerGroupLongTestRows(ctx context.Context, db *sql.DB, ref GroupLongTestRowRef) (int64, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return 0, fmt.Errorf("begin SQL Server long test cleanup: %w", err)
	}
	rows, err := readSQLServerGroupLongTestRows(ctx, tx, ref, []string{longColRecordID}, true)
	if err != nil {
		return 0, rollbackWith(tx, err)
	}
	receipts, err := readSQLServerLongTestReceipts(ctx, tx, ref, true)
	if err != nil {
		return 0, rollbackWith(tx, err)
	}
	expectedEffects := make(map[string]string, len(ref.RecordIDs))
	for index, recordID := range ref.RecordIDs {
		expectedEffects[recordID] = ref.EffectKeys[index]
	}
	for _, row := range rows {
		recordID, ok := groupLongText(row[0])
		if !ok || receipts[expectedEffects[recordID]] == "" {
			return 0, rollbackWith(tx, fmt.Errorf("%w: SQL Server test row has no matching receipt", ErrOwnedRowRefInvalid))
		}
	}

	deleteSQL, deleteArgs := sqlServerGroupLongTestDelete(ref)
	result, err := tx.ExecContext(ctx, deleteSQL, deleteArgs...)
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("delete SQL Server long test rows: %w", err))
	}
	removed, err := result.RowsAffected()
	if err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("read deleted SQL Server long test rows: %w", err))
	}
	if removed != int64(len(rows)) {
		return 0, rollbackWith(tx, fmt.Errorf("%w: SQL Server long test row count changed", ErrOwnedRowRefInvalid))
	}
	receiptQuery, receiptArgs := sqlServerGroupLongTestReceiptDelete(ref)
	if _, err := tx.ExecContext(ctx, receiptQuery, receiptArgs...); err != nil {
		return 0, rollbackWith(tx, fmt.Errorf("delete SQL Server long test receipts: %w", err))
	}
	remainingRows, err := readSQLServerGroupLongTestRows(ctx, tx, ref, []string{longColRecordID}, true)
	if err != nil {
		return 0, rollbackWith(tx, err)
	}
	if len(remainingRows) != 0 {
		return 0, rollbackWith(tx, fmt.Errorf("%w: SQL Server long test rows remain", ErrOwnedRowRefInvalid))
	}
	remainingReceipts, err := readSQLServerLongTestReceipts(ctx, tx, ref, true)
	if err != nil {
		return 0, rollbackWith(tx, err)
	}
	if len(remainingReceipts) != 0 {
		return 0, rollbackWith(tx, fmt.Errorf("%w: SQL Server long test receipts remain", ErrOwnedRowRefInvalid))
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOwnedCleanupAmbiguous, err)
	}
	return removed, nil
}

func sqlServerGroupLongTestDelete(ref GroupLongTestRowRef) (query string, args []any) {
	args = make([]any, 0, 4+2*len(ref.RecordIDs))
	args = append(args, ref.OwnerValue, ref.WorkspaceID, ref.GroupID, ref.GroupRevision)
	owner, workspaceID, groupID, revision := "@p1", "@p2", "@p3", "@p4"
	ownerPredicate := sqlServerLongExactText(quoteIdentifier(ref.Kind, ref.OwnerColumn), owner)
	workspacePredicate := sqlServerLongExactText(quoteIdentifier(ref.Kind, longColWorkspaceID), workspaceID)
	groupPredicate := sqlServerLongExactText(quoteIdentifier(ref.Kind, longColGroupID), groupID)
	revisionPredicate := sqlServerLongExactText(quoteIdentifier(ref.Kind, longColGroupRev), revision)
	conditions := make([]string, len(ref.RecordIDs))
	for index, recordID := range ref.RecordIDs {
		recordParam := fmt.Sprintf("@p%d", len(args)+1)
		args = append(args, recordID)
		sampleParam := fmt.Sprintf("@p%d", len(args)+1)
		args = append(args, ref.SampleIDs[index])
		provenance := quoteIdentifier(ref.Kind, ref.ProvenanceColumn)
		conditions[index] = "(" + sqlServerLongExactRecord(quoteIdentifier(ref.Kind, longColRecordID), recordParam) +
			" AND ISJSON(" + provenance + ") = 1 AND " +
			sqlServerLongExactText("JSON_VALUE("+provenance+", '$.sample_id')", sampleParam) + ")"
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE %s AND %s AND %s AND %s AND (%s)",
		qualifiedTableName(ref.Kind, ref.SchemaName, ref.TableName), ownerPredicate, workspacePredicate,
		groupPredicate, revisionPredicate, strings.Join(conditions, " OR "))
	return query, args
}

func sqlServerGroupLongTestReceiptDelete(ref GroupLongTestRowRef) (query string, args []any) {
	conditions := make([]string, len(ref.EffectKeys))
	args = make([]any, 0, len(ref.EffectKeys)*2)
	for index, key := range ref.EffectKeys {
		keyParam := fmt.Sprintf("@p%d", len(args)+1)
		args = append(args, key)
		digestParam := fmt.Sprintf("@p%d", len(args)+1)
		args = append(args, ref.EffectDigests[key])
		conditions[index] = "(" + sqlServerLongExactRecord(quoteIdentifier(ref.Kind, "effect_key"), keyParam) +
			" AND " + sqlServerLongExactRecord(quoteIdentifier(ref.Kind, "payload_digest"), digestParam) + ")"
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE %s", qualifiedTableName(ref.Kind, ref.SchemaName, EffectReceiptTable), strings.Join(conditions, " OR "))
	return query, args
}

func readSQLServerLongTestReceipts(ctx context.Context, q groupLongQueryer, ref GroupLongTestRowRef, lock bool) (map[string]string, error) {
	args := make([]any, len(ref.EffectKeys))
	conditions := make([]string, len(ref.EffectKeys))
	for index, key := range ref.EffectKeys {
		args[index] = key
		conditions[index] = sqlServerLongExactRecord(quoteIdentifier(ref.Kind, "effect_key"), fmt.Sprintf("@p%d", index+1))
	}
	lockHint := ""
	if lock {
		lockHint = " WITH (UPDLOCK, HOLDLOCK)"
	}
	query := fmt.Sprintf("SELECT %s, %s FROM %s%s WHERE (%s)",
		quoteIdentifier(ref.Kind, "effect_key"), quoteIdentifier(ref.Kind, "payload_digest"),
		qualifiedTableName(ref.Kind, ref.SchemaName, EffectReceiptTable), lockHint, strings.Join(conditions, " OR "))
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("read SQL Server long test receipts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	result := make(map[string]string, len(ref.EffectKeys))
	for rows.Next() {
		var rawKey, rawDigest any
		if err := rows.Scan(&rawKey, &rawDigest); err != nil {
			return nil, fmt.Errorf("scan SQL Server long test receipt: %w", err)
		}
		key, keyOK := groupLongText(rawKey)
		digest, digestOK := groupLongText(rawDigest)
		want := ref.EffectDigests[key]
		if !keyOK || !digestOK || want == "" || digest != want {
			return nil, fmt.Errorf("%w: SQL Server long test receipt proof differs", ErrOwnedRowRefInvalid)
		}
		if _, duplicate := result[key]; duplicate {
			return nil, fmt.Errorf("%w: duplicate SQL Server long test receipt", ErrOwnedRowRefInvalid)
		}
		result[key] = digest
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read SQL Server long test receipts: %w", err)
	}
	return result, nil
}
