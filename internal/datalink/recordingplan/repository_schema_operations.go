package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	schemaOperationColumns = `operation_id, token, workspace_id, scope_key, owner, action, status,
		executed_statements, verified_digest, reason, next_action, created_at, updated_at, completed_at`
	schemaOperationSelect           = `SELECT ` + schemaOperationColumns + ` FROM managed_schema_operations `
	schemaOperationByIDQuery        = schemaOperationSelect + `WHERE operation_id = $1`
	schemaOperationByWorkspaceQuery = schemaOperationSelect + `WHERE operation_id = $1 AND workspace_id = $2`
	schemaOperationByTokenQuery     = schemaOperationSelect + `WHERE token = $1`
	activeSchemaOperationQuery      = schemaOperationSelect +
		`WHERE workspace_id = $1 AND scope_key = $2 AND status IN ('pending', 'running') ORDER BY created_at LIMIT 1`
	// The primary key, the unique token and the partial unique index on active
	// scopes turn this single conditional insert into the execution claim.
	claimSchemaOperationQuery = `INSERT INTO managed_schema_operations (` + schemaOperationColumns + `)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT DO NOTHING`
	finishSchemaOperationQuery = `UPDATE managed_schema_operations
		SET status = $1, executed_statements = $2, verified_digest = $3, reason = $4, next_action = $5,
			updated_at = $6, completed_at = $7
		WHERE operation_id = $8 AND owner = $9 AND owner <> '' AND status IN ('pending', 'running')`
	finishStaleSchemaOperationQuery = `UPDATE managed_schema_operations
		SET status = $1, executed_statements = $2, verified_digest = $3, reason = $4, next_action = $5,
			updated_at = $6, completed_at = $7
		WHERE operation_id = $8 AND status IN ('pending', 'running') AND updated_at < $9`
)

// ClaimSchemaOperation records op as the owner of its scope unless the same
// operation, its token or another active operation of the scope already
// exists. Separate processes sharing the database cannot both acquire it.
func (r *SQLRepository) ClaimSchemaOperation(ctx context.Context, op *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	for range 2 {
		inserted, err := r.insertSchemaOperation(ctx, op)
		if err != nil {
			return nil, "", err
		}
		if inserted {
			stored := *op
			return &stored, ClaimAcquired, nil
		}
		existing, outcome, err := r.explainClaimConflict(ctx, op)
		if err != nil || existing != nil {
			return existing, outcome, err
		}
		// The blocking operation ended between the insert and the lookup.
	}
	return nil, "", fmt.Errorf("claim schema operation %s did not settle", op.OperationID)
}

// GetSchemaOperation reads one operation of the workspace.
func (r *SQLRepository) GetSchemaOperation(ctx context.Context, workspaceID, operationID string) (*SchemaOperation, error) {
	return r.querySchemaOperation(ctx, schemaOperationByWorkspaceQuery, operationID, workspaceID)
}

// FindActiveSchemaOperation returns the pending or running operation of a scope.
func (r *SQLRepository) FindActiveSchemaOperation(ctx context.Context, workspaceID, scopeKey string) (*SchemaOperation, error) {
	return r.querySchemaOperation(ctx, activeSchemaOperationQuery, workspaceID, scopeKey)
}

// FinishSchemaOperation records the terminal result of an owned, still active operation.
func (r *SQLRepository) FinishSchemaOperation(ctx context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error) {
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, adaptPlaceholders(finishSchemaOperationQuery),
		string(result.Status), result.ExecutedStatements, result.VerifiedDigest, result.Reason, result.NextAction,
		now, now, operationID, owner)
	if err != nil {
		return nil, fmt.Errorf("finish schema operation: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read finished schema operation: %w", err)
	}
	if affected != 1 {
		return nil, fmt.Errorf("%w: %s", ErrSchemaOperationNotOwned, operationID)
	}
	return r.querySchemaOperation(ctx, schemaOperationByIDQuery, operationID)
}

// FinishStaleSchemaOperation records a terminal result for an active operation
// whose claim lease has run out, whoever owns it. A nil operation with no error
// means the owning execution finished it first and its result stands.
func (r *SQLRepository) FinishStaleSchemaOperation(ctx context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time) (*SchemaOperation, error) {
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx, adaptPlaceholders(finishStaleSchemaOperationQuery),
		string(result.Status), result.ExecutedStatements, result.VerifiedDigest, result.Reason, result.NextAction,
		now, now, operationID, staleBefore)
	if err != nil {
		return nil, fmt.Errorf("finish stale schema operation: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read stale finished schema operation: %w", err)
	}
	if affected != 1 {
		return nil, nil
	}
	return r.querySchemaOperation(ctx, schemaOperationByIDQuery, operationID)
}

func (r *SQLRepository) insertSchemaOperation(ctx context.Context, op *SchemaOperation) (bool, error) {
	res, err := r.db.ExecContext(ctx, adaptPlaceholders(claimSchemaOperationQuery),
		op.OperationID, op.Token, op.WorkspaceID, op.ScopeKey, op.Owner, op.Action, string(op.Status),
		op.ExecutedStatements, op.VerifiedDigest, op.Reason, op.NextAction, op.CreatedAt, op.UpdatedAt,
		optionalOperationTime(op.CompletedAt))
	if err != nil {
		return false, fmt.Errorf("claim schema operation: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read claimed schema operation: %w", err)
	}
	return affected == 1, nil
}

// explainClaimConflict reports what blocked an insert: the same operation, the
// token's other operation, or the active operation of the scope. It returns no
// operation and no error when the blocking operation has since ended.
func (r *SQLRepository) explainClaimConflict(ctx context.Context, op *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	existing, err := r.querySchemaOperation(ctx, schemaOperationByIDQuery, op.OperationID)
	switch {
	case err == nil:
		if existing.WorkspaceID != op.WorkspaceID || existing.Token != op.Token {
			return nil, "", ErrSchemaOperationMismatch
		}
		return withoutOwner(*existing), claimOutcomeFor(existing), nil
	case !errors.Is(err, ErrSchemaOperationNotFound):
		return nil, "", err
	}
	_, err = r.querySchemaOperation(ctx, schemaOperationByTokenQuery, op.Token)
	switch {
	case err == nil:
		return nil, "", ErrSchemaOperationMismatch
	case !errors.Is(err, ErrSchemaOperationNotFound):
		return nil, "", err
	}
	busy, err := r.FindActiveSchemaOperation(ctx, op.WorkspaceID, op.ScopeKey)
	switch {
	case err == nil:
		return withoutOwner(*busy), ClaimScopeBusy, nil
	case errors.Is(err, ErrSchemaOperationNotFound):
		return nil, "", nil
	default:
		return nil, "", err
	}
}

func (r *SQLRepository) querySchemaOperation(ctx context.Context, query string, args ...any) (*SchemaOperation, error) {
	var op SchemaOperation
	var status string
	var completed sql.NullTime
	err := r.db.QueryRowContext(ctx, adaptPlaceholders(query), args...).Scan(
		&op.OperationID, &op.Token, &op.WorkspaceID, &op.ScopeKey, &op.Owner, &op.Action, &status,
		&op.ExecutedStatements, &op.VerifiedDigest, &op.Reason, &op.NextAction, &op.CreatedAt, &op.UpdatedAt, &completed,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSchemaOperationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read schema operation: %w", err)
	}
	op.Status = SchemaOperationStatus(status)
	if completed.Valid {
		completedAt := completed.Time
		op.CompletedAt = &completedAt
	}
	return &op, nil
}

func optionalOperationTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return *value
}
