package recordingplan

import (
	"context"
	"database/sql"
	"sort"
)

// readTestWriteRecovery reads both ledger and original bindings from one local
// transaction. No join may silently remove an operation whose binding is lost.
func (r *SQLRepository) readTestWriteRecovery(ctx context.Context, workspaceID string) ([]testWriteRecoveryRecord, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, adaptPlaceholders(schemaOperationSelect+`WHERE workspace_id = $1 AND action = 'test_write' ORDER BY created_at, operation_id`), workspaceID)
	if err != nil {
		return nil, err
	}
	records := make([]testWriteRecoveryRecord, 0)
	for rows.Next() {
		var op SchemaOperation
		var completed sql.NullTime
		err = rows.Scan(&op.OperationID, &op.Token, &op.WorkspaceID, &op.ScopeKey, &op.Owner, &op.Action, &op.Status,
			&op.ExecutedStatements, &op.VerifiedDigest, &op.Reason, &op.NextAction, &op.CreatedAt, &op.UpdatedAt, &completed,
			&op.PayloadDigest, &op.WriteOutcome, &op.CleanupStatus, &op.CleanupReason, &op.Detail)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		if completed.Valid {
			at := completed.Time
			op.CompletedAt = &at
		}
		// Parse timestamps and nullable SQL values before deciding a row is
		// resolved. SQL text comparisons and NULL predicates cannot establish
		// valid operation facts. Resolved history needs no retained token and
		// consumes neither candidate memory nor the unresolved safety bound.
		if recoveryResolved(op) {
			continue
		}
		records = append(records, testWriteRecoveryRecord{operation: op})
		if len(records) > testWriteRecoveryLimit {
			_ = rows.Close()
			return nil, ErrTestWriteRecoveryUnavailable
		}
	}
	err = rows.Err()
	closeErr := rows.Close()
	if err != nil {
		return nil, err
	}
	if closeErr != nil {
		return nil, closeErr
	}
	for i := range records {
		records[i].token, err = readSchemaPreviewToken(ctx, tx, records[i].operation.Token)
		if err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *MemoryRepository) readTestWriteRecovery(_ context.Context, workspaceID string) ([]testWriteRecoveryRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	records := make([]testWriteRecoveryRecord, 0)
	for _, op := range r.operations {
		if op.WorkspaceID != workspaceID || op.Action != TestWriteAction || recoveryResolved(op) {
			continue
		}
		token, exists := r.tokens[op.Token]
		if !exists {
			return nil, ErrTestWriteRecoveryUnavailable
		}
		records = append(records, testWriteRecoveryRecord{operation: op, token: &token})
		if len(records) > testWriteRecoveryLimit {
			return nil, ErrTestWriteRecoveryUnavailable
		}
	}
	sort.Slice(records, func(i, j int) bool {
		a, b := records[i].operation, records[j].operation
		if a.CreatedAt.Equal(b.CreatedAt) {
			return a.OperationID < b.OperationID
		}
		return a.CreatedAt.Before(b.CreatedAt)
	})
	return records, nil
}
