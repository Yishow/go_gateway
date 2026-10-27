package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
)

func (s *Service) claimSQLServerPartitionMaintenance(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, state *SQLServerPartitionMaintenanceState, after *SQLServerPartitionEvidence, baseDigest, scopeKey string) (*SchemaOperation, ClaimOutcome, error) {
	prefix := fmt.Sprintf("sqlserver-partition:%s:%d:%d:", storageScopeFingerprint(req.Scope), lineage.Generation, state.Revision)
	attempt, err := repo.countSQLServerPartitionMaintenanceAttempts(ctx, req.Scope.WorkspaceID, prefix)
	if err != nil {
		return nil, "", err
	}
	token := fmt.Sprintf("%s%d", prefix, attempt)
	sum := sha256.Sum256([]byte(token))
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create sqlserver partition maintenance owner: %w", err)
	}
	detail, err := json.Marshal(sqlServerPartitionMaintenanceDetail{
		Kind: sqlServerPartitionMaintenanceKind, Scope: req.Scope, GroupID: req.GroupID, ConnectorRevision: req.ConnectorRevision,
		BaseLayoutDigest: baseDigest, Generation: lineage.Generation, StateRevision: state.Revision, At: req.At.UTC(),
		Before: CloneSQLServerPartitionEvidence(state.Current), After: CloneSQLServerPartitionEvidence(after),
	})
	if err != nil {
		return nil, "", fmt.Errorf("encode sqlserver partition maintenance detail: %w", err)
	}
	now := time.Now().UTC()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: "sqlserver-partition-" + hex.EncodeToString(sum[:16]), Token: token, WorkspaceID: req.Scope.WorkspaceID,
		ScopeKey: scopeKey, Owner: "claim-" + owner, Action: SQLServerPartitionMaintenanceAction,
		Status: SchemaOperationRunning, Detail: string(detail), CreatedAt: now, UpdatedAt: now,
	})
}

func (s *Service) runSQLServerPartitionMaintenance(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, state *SQLServerPartitionMaintenanceState, op *SchemaOperation, after *SQLServerPartitionEvidence, target SQLServerPartitionMaintenanceTarget) (SQLServerPartitionMaintenanceResult, error) {
	if err := s.saveSQLServerPartitionProgress(ctx, op); err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	actual, inspectErr := inspectSQLServerPartitionMaintenanceCatalog(ctx, req, lineage, target, state.Current, nil)
	if inspectErr != nil {
		return s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationUnknown, Reason: sqlServerPartitionReasonUnverifiable, NextAction: sqlServerPartitionMaintenanceNext,
		}, state.Current, sqlServerPartitionMaintenanceBlocked, sqlServerPartitionReasonUnverifiable, false)
	}
	if !sameSQLServerPartitionEvidence(actual, state.Current) {
		return s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: sqlServerPartitionReasonCatalogMismatch, NextAction: sqlServerPartitionMaintenanceNext,
		}, state.Current, sqlServerPartitionMaintenanceBlocked, sqlServerPartitionReasonCatalogMismatch, false)
	}
	if err := s.requireSQLServerPartitionOwner(ctx, op); err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	splitErr := target.Split(ctx, CloneSQLServerPartitionEvidence(state.Current), CloneSQLServerPartitionEvidence(after))
	if err := s.requireSQLServerPartitionOwner(ctx, op); err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	observed, verifyErr := inspectSQLServerPartitionMaintenanceCatalog(ctx, req, lineage, target, state.Current, after)
	switch {
	case verifyErr == nil && sameSQLServerPartitionEvidence(observed, after):
		return s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationSucceeded, ExecutedStatements: 1, VerifiedDigest: sqlServerPartitionEvidenceDigest(after),
		}, after, sqlServerPartitionMaintenanceReady, "", true)
	case verifyErr == nil && sameSQLServerPartitionEvidence(observed, state.Current):
		reason := sqlServerPartitionReasonStatementFailed
		if errors.Is(splitErr, ErrTargetPermissionDenied) {
			reason = sqlServerPartitionReasonPermissionLost
		}
		return s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: reason, NextAction: sqlServerPartitionMaintenanceNext,
		}, state.Current, sqlServerPartitionMaintenanceBlocked, reason, false)
	default:
		reason := sqlServerPartitionReasonCatalogMismatch
		if verifyErr != nil || splitErr != nil {
			reason = sqlServerPartitionReasonUnverifiable
		}
		return s.finishSQLServerPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationUnknown, ExecutedStatements: 1, Reason: reason, NextAction: sqlServerPartitionMaintenanceNext,
		}, state.Current, sqlServerPartitionMaintenanceBlocked, reason, false)
	}
}

func (s *Service) saveSQLServerPartitionProgress(ctx context.Context, op *SchemaOperation) error {
	progress, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return fmt.Errorf("%w: operation progress repository unavailable", ErrSchemaOperationUnacknowledged)
	}
	if err := progress.SaveSchemaOperationProgress(ctx, op.OperationID, op.Owner, op.Detail); err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	return nil
}

func (s *Service) requireSQLServerPartitionOwner(ctx context.Context, op *SchemaOperation) error {
	current, err := s.repo.GetSchemaOperation(ctx, op.WorkspaceID, op.OperationID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	if current.Action != SQLServerPartitionMaintenanceAction || current.ScopeKey != op.ScopeKey || current.Owner != op.Owner || !current.Status.IsActive() {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, ErrSchemaOperationNotOwned)
	}
	return nil
}

func (s *Service) finishSQLServerPartitionMaintenance(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, state *SQLServerPartitionMaintenanceState, op *SchemaOperation, result SchemaOperationResult, current *SQLServerPartitionEvidence, status, reason string, changed bool) (SQLServerPartitionMaintenanceResult, error) {
	result.Detail = op.Detail
	apply := func(ctx context.Context, q DBTX) error {
		next := *state
		next.Current = CloneSQLServerPartitionEvidence(current)
		next.Revision++
		next.LastOperationID = op.OperationID
		next.Status, next.Reason, next.UpdatedAt = status, reason, time.Now().UTC()
		return updateSQLServerPartitionMaintenanceState(ctx, q, state.Revision, &next)
	}
	finished, err := repo.finishSQLServerPartitionMaintenance(ctx, op.OperationID, op.Owner, result, apply)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	updated, err := ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	code := ""
	if result.Status != SchemaOperationSucceeded {
		code = sqlServerPartitionMaintenanceCode
	}
	return sqlServerMaintenanceResult(updated, finished, code, reason, changed), nil
}

func (s *Service) persistSQLServerMaintenanceBlock(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, state *SQLServerPartitionMaintenanceState, op *SchemaOperation, reason string) (SQLServerPartitionMaintenanceResult, error) {
	if state == nil {
		return sqlServerMaintenanceBlocked(nil, op, reason), nil
	}
	next, err := persistSQLServerMaintenanceState(ctx, repo.storageDB(), state, state.Current, state.LastOperationID, sqlServerPartitionMaintenanceBlocked, reason)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	return sqlServerMaintenanceBlocked(next, op, reason), nil
}

func persistSQLServerMaintenanceState(ctx context.Context, q DBTX, state *SQLServerPartitionMaintenanceState, current *SQLServerPartitionEvidence, operationID, status, reason string) (*SQLServerPartitionMaintenanceState, error) {
	next := *state
	next.Current = CloneSQLServerPartitionEvidence(current)
	next.Revision++
	next.LastOperationID, next.Status, next.Reason, next.UpdatedAt = operationID, status, reason, time.Now().UTC()
	if err := updateSQLServerPartitionMaintenanceState(ctx, q, state.Revision, &next); err != nil {
		return nil, err
	}
	return ReadSQLServerPartitionMaintenanceState(ctx, q, state.Scope)
}
