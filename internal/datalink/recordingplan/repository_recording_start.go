package recordingplan

import (
	"context"
	"fmt"
	"time"
)

// ResumeRecordingStartOperation reopens only a matching start action. Schema
// and test-write operations keep their existing terminal-result contracts.
func (r *SQLRepository) ResumeRecordingStartOperation(ctx context.Context, id, owner, digest string) (*SchemaOperation, error) {
	result, err := r.db.ExecContext(ctx, adaptPlaceholders(`UPDATE managed_schema_operations
		SET owner=$1,status='running',reason='',next_action='',completed_at=NULL,updated_at=$2
		WHERE operation_id=$3 AND action=$4 AND payload_digest=$5
		AND status IN ('partial','failed','unknown')`), owner, time.Now().UTC(), id, RecordingStartAction, digest)
	if err != nil {
		return nil, fmt.Errorf("resume recording start operation: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read resumed recording start operation: %w", err)
	}
	if affected != 1 {
		return nil, nil
	}
	return r.querySchemaOperation(ctx, schemaOperationByIDQuery, id)
}

// ResumeRecordingStartOperation provides the same owner fencing for in-memory
// tests as the SQL ledger, without introducing another operation repository.
func (r *MemoryRepository) ResumeRecordingStartOperation(_ context.Context, id, owner, digest string) (*SchemaOperation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	op, ok := r.operations[id]
	if !ok || op.Action != RecordingStartAction || op.PayloadDigest != digest || op.Status.IsActive() || op.Status == SchemaOperationSucceeded {
		return nil, nil
	}
	for _, active := range r.operations {
		if active.Status.IsActive() && active.WorkspaceID == op.WorkspaceID && active.ScopeKey == op.ScopeKey {
			return nil, nil
		}
	}
	op.Owner, op.Status, op.UpdatedAt = owner, SchemaOperationRunning, time.Now().UTC()
	op.Reason, op.NextAction, op.CompletedAt = "", "", nil
	r.operations[id] = op
	return &op, nil
}
