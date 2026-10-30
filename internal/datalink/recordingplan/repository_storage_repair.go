package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// storageRepairRepository is what restoring a lost managed table needs from the
// local database beyond the operation ledger: ending an operation together with
// the lineage change it proves, counting earlier attempts, and listing repairs
// whose owner never came back. The SQL repository provides it; a repository
// without it cannot restore storage.
type storageRepairRepository interface {
	storageDB() DBTX
	countStorageRepairAttempts(ctx context.Context, workspaceID, tokenPrefix string) (int, error)
	finishStorageRepair(ctx context.Context, operationID, owner string, result SchemaOperationResult, apply func(context.Context, DBTX) error) (*SchemaOperation, error)
	finishStaleStorageRepair(ctx context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time, apply func(context.Context, DBTX) error) (*SchemaOperation, error)
	listStaleStorageRepairs(ctx context.Context, staleBefore time.Time) ([]SchemaOperation, error)
}

var _ storageRepairRepository = (*SQLRepository)(nil)

func (r *SQLRepository) storageDB() DBTX { return r.db }

func (r *SQLRepository) countStorageRepairAttempts(ctx context.Context, workspaceID, tokenPrefix string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, adaptPlaceholders(`SELECT COUNT(*) FROM managed_schema_operations
		WHERE workspace_id = $1 AND action = $2 AND substr(token, 1, length($3)) = $3`),
		workspaceID, StorageRepairAction, tokenPrefix, tokenPrefix).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count storage repair attempts: %w", err)
	}
	return count, nil
}

// finishStorageRepair ends an owned, still active repair operation and applies
// its lineage change in the same local transaction, so the ledger and the
// scope never disagree about whether the table was rebuilt.
func (r *SQLRepository) finishStorageRepair(ctx context.Context, operationID, owner string, result SchemaOperationResult, apply func(context.Context, DBTX) error) (*SchemaOperation, error) {
	return r.finishRepairInTx(ctx, apply, func(tx *sql.Tx, now time.Time) (bool, error) {
		res, err := tx.ExecContext(ctx, adaptPlaceholders(finishSchemaOperationQuery),
			string(result.Status), result.ExecutedStatements, result.VerifiedDigest, result.Reason, result.NextAction,
			now, now, result.PayloadDigest, result.WriteOutcome, result.CleanupStatus, result.CleanupReason, result.Detail,
			operationID, owner)
		if err != nil {
			return false, fmt.Errorf("finish storage repair: %w", err)
		}
		affected, err := res.RowsAffected()
		return affected == 1, err
	}, operationID, true, result)
}

// finishStaleStorageRepair ends an active repair whose lease ran out, whoever
// owns it. A nil operation without error means its owner finished it first.
func (r *SQLRepository) finishStaleStorageRepair(ctx context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time, apply func(context.Context, DBTX) error) (*SchemaOperation, error) {
	return r.finishRepairInTx(ctx, apply, func(tx *sql.Tx, now time.Time) (bool, error) {
		res, err := tx.ExecContext(ctx, adaptPlaceholders(finishStaleSchemaOperationQuery),
			string(result.Status), result.ExecutedStatements, result.VerifiedDigest, result.Reason, result.NextAction,
			now, now, result.PayloadDigest, result.WriteOutcome, result.CleanupStatus, result.CleanupReason, result.Detail,
			operationID, staleBefore)
		if err != nil {
			return false, fmt.Errorf("finish stale storage repair: %w", err)
		}
		affected, err := res.RowsAffected()
		return affected == 1, err
	}, operationID, false, result)
}

func (r *SQLRepository) finishRepairInTx(ctx context.Context, apply func(context.Context, DBTX) error, finish func(*sql.Tx, time.Time) (bool, error), operationID string, mustOwn bool, result SchemaOperationResult) (*SchemaOperation, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin storage repair finish: %w", err)
	}
	settledAt := time.Now().UTC()
	finished, err := finish(tx, settledAt)
	if err == nil && finished && apply != nil {
		err = apply(ctx, tx)
	}
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback storage repair finish: %w", rollbackErr))
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		if finished {
			recordStorageSettlement(ctx, operationID, mustOwn, result, settledAt, nil, err)
		}
		return nil, fmt.Errorf("commit storage repair finish: %w", err)
	}
	if !finished {
		if mustOwn {
			return nil, fmt.Errorf("%w: %s", ErrSchemaOperationNotOwned, operationID)
		}
		return nil, nil
	}
	op, err := r.querySchemaOperation(ctx, schemaOperationByIDQuery, operationID)
	recordStorageSettlement(ctx, operationID, mustOwn, result, settledAt, op, err)
	return op, err
}

func (r *SQLRepository) listStaleStorageRepairs(ctx context.Context, staleBefore time.Time) ([]SchemaOperation, error) {
	rows, err := r.db.QueryContext(ctx, adaptPlaceholders(schemaOperationSelect+
		`WHERE action = $1 AND status IN ('pending', 'running') AND updated_at < $2 ORDER BY created_at LIMIT 50`),
		StorageRepairAction, staleBefore)
	if err != nil {
		return nil, fmt.Errorf("list stale storage repairs: %w", err)
	}
	defer rows.Close()
	var operations []SchemaOperation
	for rows.Next() {
		var op SchemaOperation
		var status string
		var completed sql.NullTime
		if err := rows.Scan(&op.OperationID, &op.Token, &op.WorkspaceID, &op.ScopeKey, &op.Owner, &op.Action, &status,
			&op.ExecutedStatements, &op.VerifiedDigest, &op.Reason, &op.NextAction, &op.CreatedAt, &op.UpdatedAt, &completed,
			&op.PayloadDigest, &op.WriteOutcome, &op.CleanupStatus, &op.CleanupReason, &op.Detail); err != nil {
			return nil, fmt.Errorf("scan stale storage repair: %w", err)
		}
		op.Status = SchemaOperationStatus(status)
		operations = append(operations, op)
	}
	return operations, rows.Err()
}
