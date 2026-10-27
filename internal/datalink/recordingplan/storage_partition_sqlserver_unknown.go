package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
)

type sqlServerPartitionUnknownRepository interface {
	reclaimSQLServerPartitionUnknown(context.Context, *SchemaOperation, string, string) (*SchemaOperation, error)
}

// Reclaim the recorded operation under the existing active-scope constraint.
// Only the matching unknown result can move back to running; a concurrent
// recovery or schema operation wins without granting this caller ownership.
func (r *SQLRepository) reclaimSQLServerPartitionUnknown(ctx context.Context, previous *SchemaOperation, detail, owner string) (*SchemaOperation, error) {
	query := adaptPlaceholders(`UPDATE managed_schema_operations
		SET status = 'running', owner = $1, updated_at = $2, completed_at = NULL,
			reason = '', next_action = '', detail = $3
		WHERE operation_id = $4 AND workspace_id = $5 AND token = $6 AND scope_key = $7
			AND action = 'sqlserver_partition_maintenance' AND status = 'unknown' AND detail = $8
			AND NOT EXISTS (SELECT 1 FROM managed_schema_operations active
				WHERE active.workspace_id = $9 AND active.scope_key = $10
				AND active.operation_id <> $11 AND active.status IN ('pending', 'running'))`)
	result, err := r.db.ExecContext(ctx, query, owner, time.Now().UTC(), detail,
		previous.OperationID, previous.WorkspaceID, previous.Token, previous.ScopeKey, previous.Detail,
		previous.WorkspaceID, previous.ScopeKey, previous.OperationID)
	if err != nil {
		return nil, fmt.Errorf("reclaim unknown sqlserver partition operation: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("read reclaimed sqlserver partition operation: %w", err)
	}
	if changed != 1 {
		return nil, nil
	}
	return r.querySchemaOperation(ctx, schemaOperationByIDQuery, previous.OperationID)
}

func (s *Service) recoverUnknownSQLServerPartitionMaintenance(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, state *SQLServerPartitionMaintenanceState, target SQLServerPartitionMaintenanceTarget, baseDigest, scopeKey string) (*SQLServerPartitionMaintenanceResult, error) {
	if state == nil || state.Status != sqlServerPartitionMaintenanceBlocked || state.LastOperationID == "" {
		return nil, nil
	}
	op, err := s.repo.GetSchemaOperation(ctx, req.Scope.WorkspaceID, state.LastOperationID)
	if err != nil {
		return nil, err
	}
	if op.Status != SchemaOperationUnknown || op.Action != SQLServerPartitionMaintenanceAction || op.ScopeKey != scopeKey {
		return nil, nil
	}
	var detail sqlServerPartitionMaintenanceDetail
	if json.Unmarshal([]byte(op.Detail), &detail) != nil || detail.StateRevision >= state.Revision {
		blocked := sqlServerMaintenanceBlocked(state, op, sqlServerPartitionReasonUnverifiable)
		return &blocked, nil //nolint:nilerr // Invalid persisted detail is a visible blocked result, never authority.
	}
	// Finishing an unknown attempt advanced the mutable state revision, while
	// retaining its exact before proof and operation identity. Validate every
	// original binding before rebinding the next claim to the current revision.
	original := *state
	original.Revision = detail.StateRevision
	if !validSQLServerPartitionMaintenanceDetail(detail, req, lineage, &original, baseDigest) {
		blocked := sqlServerMaintenanceBlocked(state, op, sqlServerPartitionReasonUnverifiable)
		return &blocked, nil
	}
	reclaimer, ok := s.repo.(sqlServerPartitionUnknownRepository)
	if !ok {
		return nil, ErrSQLServerPartitionMaintenanceUnsupported
	}
	detail.StateRevision = state.Revision
	encoded, err := json.Marshal(detail)
	if err != nil {
		return nil, fmt.Errorf("encode recovered sqlserver partition detail: %w", err)
	}
	owner, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create sqlserver partition recovery owner: %w", err)
	}
	claimed, err := reclaimer.reclaimSQLServerPartitionUnknown(ctx, op, string(encoded), "claim-"+owner)
	if err != nil {
		return nil, err
	}
	if claimed == nil {
		blocked := sqlServerMaintenanceBlocked(state, op, "scope_busy")
		return &blocked, nil
	}
	observed, inspectErr := inspectSQLServerPartitionMaintenanceCatalog(ctx, req, lineage, target, detail.Before, detail.After)
	status, reason, changed := SchemaOperationUnknown, sqlServerPartitionReasonUnverifiable, false
	current := state.Current
	switch {
	case inspectErr == nil && sameSQLServerPartitionEvidence(observed, detail.After):
		status, reason, changed, current = SchemaOperationSucceeded, "", true, detail.After
	case inspectErr == nil && sameSQLServerPartitionEvidence(observed, detail.Before):
		status, reason = SchemaOperationFailed, sqlServerPartitionReasonAttemptAbandoned
	}
	result := SchemaOperationResult{Status: status, Reason: reason, NextAction: sqlServerPartitionMaintenanceNext}
	stateStatus := sqlServerPartitionMaintenanceBlocked
	if status == SchemaOperationSucceeded {
		result.ExecutedStatements, result.VerifiedDigest = 1, sqlServerPartitionEvidenceDigest(current)
		result.NextAction, stateStatus = "", sqlServerPartitionMaintenanceReady
	}
	recovered, err := s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, claimed, result, current, stateStatus, reason, changed)
	if errors.Is(err, ErrSchemaOperationNotOwned) {
		blocked := sqlServerMaintenanceBlocked(state, nil, "scope_busy")
		return &blocked, nil
	}
	return &recovered, err
}
