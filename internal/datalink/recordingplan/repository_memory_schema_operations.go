package recordingplan

import (
	"context"
	"fmt"
	"time"
)

// ClaimSchemaOperation records op as the owner of its scope unless the same
// operation, its token or another active operation of the scope already exists.
func (r *MemoryRepository) ClaimSchemaOperation(_ context.Context, op *SchemaOperation) (*SchemaOperation, ClaimOutcome, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.operations[op.OperationID]; ok {
		if existing.WorkspaceID != op.WorkspaceID || existing.Token != op.Token {
			return nil, "", ErrSchemaOperationMismatch
		}
		return withoutOwner(existing), claimOutcomeFor(&existing), nil
	}
	for _, existing := range r.operations {
		if existing.Token == op.Token {
			return nil, "", ErrSchemaOperationMismatch
		}
		if existing.WorkspaceID == op.WorkspaceID && existing.ScopeKey == op.ScopeKey && existing.Status.IsActive() {
			return withoutOwner(existing), ClaimScopeBusy, nil
		}
	}
	r.operations[op.OperationID] = *op
	stored := *op
	return &stored, ClaimAcquired, nil
}

// GetSchemaOperation reads one operation of the workspace.
func (r *MemoryRepository) GetSchemaOperation(_ context.Context, workspaceID, operationID string) (*SchemaOperation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	op, ok := r.operations[operationID]
	if !ok || op.WorkspaceID != workspaceID {
		return nil, fmt.Errorf("%w: %s", ErrSchemaOperationNotFound, operationID)
	}
	return &op, nil
}

// FindActiveSchemaOperation returns the pending or running operation of a scope.
func (r *MemoryRepository) FindActiveSchemaOperation(_ context.Context, workspaceID, scopeKey string) (*SchemaOperation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, op := range r.operations {
		if op.WorkspaceID == workspaceID && op.ScopeKey == scopeKey && op.Status.IsActive() {
			return &op, nil
		}
	}
	return nil, fmt.Errorf("%w: active scope", ErrSchemaOperationNotFound)
}

// FinishSchemaOperation records the terminal result of an owned, still active operation.
func (r *MemoryRepository) FinishSchemaOperation(_ context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	op, ok := r.operations[operationID]
	if !ok || owner == "" || op.Owner != owner || !op.Status.IsActive() {
		return nil, fmt.Errorf("%w: %s", ErrSchemaOperationNotOwned, operationID)
	}
	now := time.Now().UTC()
	op.Status, op.ExecutedStatements, op.VerifiedDigest = result.Status, result.ExecutedStatements, result.VerifiedDigest
	op.Reason, op.NextAction, op.UpdatedAt, op.CompletedAt = result.Reason, result.NextAction, now, &now
	r.operations[operationID] = op
	return &op, nil
}

// FinishStaleSchemaOperation records a terminal result for an active operation
// whose claim lease has run out, whoever owns it. A nil operation with no error
// means the owning execution finished it first and its result stands.
func (r *MemoryRepository) FinishStaleSchemaOperation(_ context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time) (*SchemaOperation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	op, ok := r.operations[operationID]
	if !ok || !op.Status.IsActive() || !op.UpdatedAt.Before(staleBefore) {
		return nil, nil
	}
	now := time.Now().UTC()
	op.Status, op.ExecutedStatements, op.VerifiedDigest = result.Status, result.ExecutedStatements, result.VerifiedDigest
	op.Reason, op.NextAction, op.UpdatedAt, op.CompletedAt = result.Reason, result.NextAction, now, &now
	r.operations[operationID] = op
	return &op, nil
}
