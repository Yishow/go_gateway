package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const reclaimTestWriteCleanupOperationQuery = `UPDATE managed_schema_operations
SET owner = $1, status = 'running', updated_at = $2
WHERE operation_id = $3 AND workspace_id = $4 AND token = $5 AND scope_key = $6 AND action = 'test_write'
  AND write_outcome = 'written_verified' AND cleanup_status IN ('failed', 'unknown')
  AND status NOT IN ('pending', 'running') AND detail = $7
  AND NOT EXISTS (
      SELECT 1 FROM managed_schema_operations active
      WHERE active.workspace_id = $8 AND active.scope_key = $9 AND active.status IN ('pending', 'running') AND active.operation_id <> $10
  )`

// TestWriteCleanupRepository is an optional repository capability for reclaiming
// a test-write operation that verified write and readback but failed or lost cleanup.
type TestWriteCleanupRepository interface {
	ReclaimTestWriteCleanupOperation(ctx context.Context, operationID, workspaceID, token, scopeKey, expectedDetail, newOwner string) (*SchemaOperation, ClaimOutcome, error)
}

// ReclaimTestWriteCleanupOperation atomically re-claims a completed test write
// whose cleanup is failed or unknown, locking the scope and preventing concurrent
// target operations.
func (r *SQLRepository) ReclaimTestWriteCleanupOperation(ctx context.Context, operationID, workspaceID, token, scopeKey, expectedDetail, newOwner string) (*SchemaOperation, ClaimOutcome, error) {
	if operationID == "" || workspaceID == "" || token == "" || scopeKey == "" || expectedDetail == "" || newOwner == "" {
		return nil, "", ErrSchemaOperationMismatch
	}
	res, err := r.db.ExecContext(ctx, adaptPlaceholders(reclaimTestWriteCleanupOperationQuery),
		newOwner, time.Now().UTC(), operationID, workspaceID, token, scopeKey, expectedDetail,
		workspaceID, scopeKey, operationID)
	if err != nil {
		active, activeErr := r.FindActiveSchemaOperation(ctx, workspaceID, scopeKey)
		if activeErr == nil {
			if active.OperationID == operationID {
				return withoutOwner(*active), ClaimInProgress, nil
			}
			return withoutOwner(*active), ClaimScopeBusy, nil
		}
		if !errors.Is(activeErr, ErrSchemaOperationNotFound) {
			return nil, "", fmt.Errorf("inspect active test write operation: %w", activeErr)
		}
		return nil, "", fmt.Errorf("reclaim test write cleanup: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, "", fmt.Errorf("read reclaimed test write cleanup: %w", err)
	}
	if affected == 1 {
		reclaimed, err := r.querySchemaOperation(ctx, schemaOperationByIDQuery, operationID)
		if err != nil {
			return nil, "", err
		}
		return reclaimed, ClaimAcquired, nil
	}
	current, err := r.querySchemaOperation(ctx, schemaOperationByWorkspaceQuery, operationID, workspaceID)
	if err != nil {
		return nil, "", err
	}
	if current.Token != token || current.ScopeKey != scopeKey || current.Action != TestWriteAction {
		return nil, "", ErrSchemaOperationMismatch
	}
	if current.CleanupStatus == cleanupStatusCleaned {
		return withoutOwner(*current), ClaimCompleted, nil
	}
	if current.Status.IsActive() {
		return withoutOwner(*current), ClaimInProgress, nil
	}
	active, activeErr := r.FindActiveSchemaOperation(ctx, workspaceID, scopeKey)
	if activeErr == nil && active.OperationID != operationID {
		return withoutOwner(*active), ClaimScopeBusy, nil
	}
	return withoutOwner(*current), ClaimCompleted, nil
}

// ReclaimTestWriteCleanupOperation in MemoryRepository provides the same atomic CAS for testing.
func (r *MemoryRepository) ReclaimTestWriteCleanupOperation(_ context.Context, operationID, workspaceID, token, scopeKey, expectedDetail, newOwner string) (*SchemaOperation, ClaimOutcome, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	op, ok := r.operations[operationID]
	if !ok || op.WorkspaceID != workspaceID {
		return nil, "", ErrSchemaOperationNotFound
	}
	if op.Token != token || op.ScopeKey != scopeKey || op.Action != TestWriteAction {
		return nil, "", ErrSchemaOperationMismatch
	}
	if op.CleanupStatus == cleanupStatusCleaned {
		return withoutOwner(op), ClaimCompleted, nil
	}
	if op.Status.IsActive() {
		return withoutOwner(op), ClaimInProgress, nil
	}
	if op.WriteOutcome != writeOutcomeVerified || (op.CleanupStatus != "failed" && op.CleanupStatus != "unknown") {
		return withoutOwner(op), ClaimCompleted, nil
	}
	if op.Detail != expectedDetail {
		return withoutOwner(op), ClaimCompleted, nil
	}
	for _, other := range r.operations {
		if other.OperationID != operationID && other.WorkspaceID == workspaceID && other.ScopeKey == scopeKey && other.Status.IsActive() {
			return withoutOwner(other), ClaimScopeBusy, nil
		}
	}
	op.Owner, op.Status, op.UpdatedAt = newOwner, SchemaOperationRunning, time.Now().UTC()
	r.operations[operationID] = op
	return &op, ClaimAcquired, nil
}
