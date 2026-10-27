package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/common"
)

func (s *Service) claimMySQLPartitionMaintenance(ctx context.Context, repo mysqlPartitionMaintenanceRepository, req MySQLPartitionMaintenanceRequest, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, after *MySQLPartitionEvidence, baseDigest, scopeKey string) (*SchemaOperation, ClaimOutcome, error) {
	prefix := fmt.Sprintf("partition:%s:%d:%d:", storageScopeFingerprint(req.Scope), lineage.Generation, state.Revision)
	attempt, err := repo.countMySQLPartitionMaintenanceAttempts(ctx, req.Scope.WorkspaceID, prefix)
	if err != nil {
		return nil, "", err
	}
	token := fmt.Sprintf("%s%d", prefix, attempt)
	sum := sha256.Sum256([]byte(token))
	operationID := "partition-" + hex.EncodeToString(sum[:16])
	owner, err := common.NewUUID()
	if err != nil {
		return nil, "", fmt.Errorf("create mysql partition maintenance owner: %w", err)
	}
	detail, err := json.Marshal(mysqlPartitionMaintenanceDetail{
		Kind: mysqlPartitionMaintenanceKind, Scope: req.Scope, GroupID: req.GroupID, ConnectorRevision: req.ConnectorRevision,
		BaseLayoutDigest: baseDigest, Generation: lineage.Generation, StateRevision: state.Revision, At: req.At.UTC(),
		Before: cloneMySQLPartitionEvidence(state.Current), After: cloneMySQLPartitionEvidence(after),
	})
	if err != nil {
		return nil, "", fmt.Errorf("encode mysql partition maintenance detail: %w", err)
	}
	now := time.Now().UTC()
	return s.repo.ClaimSchemaOperation(ctx, &SchemaOperation{
		OperationID: operationID, Token: token, WorkspaceID: req.Scope.WorkspaceID, ScopeKey: scopeKey,
		Owner: "claim-" + owner, Action: MySQLPartitionMaintenanceAction, Status: SchemaOperationRunning,
		Detail: string(detail), CreatedAt: now, UpdatedAt: now,
	})
}

func (s *Service) runMySQLPartitionMaintenance(ctx context.Context, repo mysqlPartitionMaintenanceRepository, req MySQLPartitionMaintenanceRequest, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, op *SchemaOperation, after *MySQLPartitionEvidence, target MySQLPartitionMaintenanceTarget) (MySQLPartitionMaintenanceResult, error) {
	detail := op.Detail
	if err := s.saveMySQLPartitionMaintenanceProgress(ctx, op, detail); err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	if err := s.requireMySQLPartitionMaintenanceOwner(ctx, op); err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	actual, inspectErr := inspectMySQLPartitionMaintenanceCatalog(ctx, repo.storageDB(), req, lineage, target, state.Current, nil)
	if inspectErr != nil {
		return s.finishMySQLPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationUnknown, Reason: mysqlPartitionReasonUnverifiable, NextAction: mysqlPartitionMaintenanceNext,
		}, state.Current, mysqlPartitionMaintenanceBlocked, mysqlPartitionReasonUnverifiable, false)
	}
	if !sameMySQLPartitionEvidence(actual, state.Current) {
		return s.finishMySQLPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: mysqlPartitionReasonCatalogMismatch, NextAction: mysqlPartitionMaintenanceNext,
		}, state.Current, mysqlPartitionMaintenanceBlocked, mysqlPartitionReasonCatalogMismatch, false)
	}
	if err := s.saveMySQLPartitionMaintenanceProgress(ctx, op, detail); err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	if err := s.requireMySQLPartitionMaintenanceOwner(ctx, op); err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	execErr, heartbeatErr := s.reorganizeMySQLPartitionWithHeartbeat(ctx, op, detail, target, state.Current, after)
	if heartbeatErr != nil {
		return MySQLPartitionMaintenanceResult{}, fmt.Errorf("mysql partition maintenance heartbeat: %w", heartbeatErr)
	}
	if err := s.requireMySQLPartitionMaintenanceOwner(ctx, op); err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	observed, verifyErr := inspectMySQLPartitionMaintenanceCatalog(ctx, repo.storageDB(), req, lineage, target, state.Current, after)
	switch {
	case verifyErr == nil && sameMySQLPartitionEvidence(observed, after):
		return s.finishMySQLPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationSucceeded, ExecutedStatements: 1, VerifiedDigest: mysqlPartitionEvidenceDigest(after),
		}, after, mysqlPartitionMaintenanceReady, "", true)
	case verifyErr == nil && sameMySQLPartitionEvidence(observed, state.Current):
		reason := "statement_failed"
		if errors.Is(execErr, ErrTargetPermissionDenied) {
			reason = "permission_lost"
		}
		return s.finishMySQLPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationFailed, Reason: reason, NextAction: mysqlPartitionMaintenanceNext,
		}, state.Current, mysqlPartitionMaintenanceBlocked, reason, false)
	default:
		reason := mysqlPartitionReasonCatalogMismatch
		if verifyErr != nil {
			reason = mysqlPartitionReasonUnverifiable
		}
		return s.finishMySQLPartitionMaintenance(ctx, repo, req, state, op, SchemaOperationResult{
			Status: SchemaOperationUnknown, ExecutedStatements: 1, Reason: reason, NextAction: mysqlPartitionMaintenanceNext,
		}, state.Current, mysqlPartitionMaintenanceBlocked, reason, false)
	}
}

func (s *Service) reorganizeMySQLPartitionWithHeartbeat(ctx context.Context, op *SchemaOperation, detail string, target MySQLPartitionMaintenanceTarget, before, after *MySQLPartitionEvidence) (execErr, heartbeatErr error) {
	progressRepo, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return nil, ErrSchemaOperationUnacknowledged
	}
	heartbeatCtx, cancel := context.WithCancel(ctx)
	heartbeatDone := make(chan error, 1)
	go func() {
		ticker := time.NewTicker(SchemaOperationLease / 3)
		defer ticker.Stop()
		for {
			select {
			case <-heartbeatCtx.Done():
				heartbeatDone <- nil
				return
			case <-ticker.C:
				if err := progressRepo.SaveSchemaOperationProgress(heartbeatCtx, op.OperationID, op.Owner, detail); err != nil {
					cancel()
					heartbeatDone <- err
					return
				}
			}
		}
	}()
	execErr = target.Reorganize(heartbeatCtx, cloneMySQLPartitionEvidence(before), cloneMySQLPartitionEvidence(after))
	cancel()
	heartbeatErr = <-heartbeatDone
	return execErr, heartbeatErr
}

func (s *Service) saveMySQLPartitionMaintenanceProgress(ctx context.Context, op *SchemaOperation, detail string) error {
	progressRepo, ok := s.repo.(operationExecutionRepository)
	if !ok {
		return fmt.Errorf("%w: operation progress repository unavailable", ErrSchemaOperationUnacknowledged)
	}
	if err := progressRepo.SaveSchemaOperationProgress(ctx, op.OperationID, op.Owner, detail); err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	return nil
}

func (s *Service) requireMySQLPartitionMaintenanceOwner(ctx context.Context, op *SchemaOperation) error {
	current, err := s.repo.GetSchemaOperation(ctx, op.WorkspaceID, op.OperationID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	if current.Action != MySQLPartitionMaintenanceAction || current.ScopeKey != op.ScopeKey || current.Owner != op.Owner || !current.Status.IsActive() {
		return fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, ErrSchemaOperationNotOwned)
	}
	return nil
}

func (s *Service) finishMySQLPartitionMaintenance(ctx context.Context, repo mysqlPartitionMaintenanceRepository, req MySQLPartitionMaintenanceRequest, state *MySQLPartitionMaintenanceState, op *SchemaOperation, result SchemaOperationResult, current *MySQLPartitionEvidence, status, reason string, changed bool) (MySQLPartitionMaintenanceResult, error) {
	result.Detail = op.Detail
	apply := func(ctx context.Context, q DBTX) error {
		next := *state
		next.Current = cloneMySQLPartitionEvidence(current)
		next.Revision++
		next.LastOperationID = op.OperationID
		next.Status, next.Reason, next.UpdatedAt = status, reason, time.Now().UTC()
		return updateMySQLPartitionMaintenanceState(ctx, q, state.Revision, &next)
	}
	finished, err := repo.finishMySQLPartitionMaintenance(ctx, op.OperationID, op.Owner, result, apply)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, fmt.Errorf("%w: %w", ErrSchemaOperationUnacknowledged, err)
	}
	updated, err := ReadMySQLPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	code := ""
	if result.Status != SchemaOperationSucceeded {
		code = mysqlPartitionMaintenanceCode
	}
	return maintenanceResult(updated, finished, code, reason, changed), nil
}

func mysqlPartitionEvidenceDigest(evidence *MySQLPartitionEvidence) string {
	encoded, err := json.Marshal(evidence)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:])
}

func sameMySQLPartitionPrefix(before, after *MySQLPartitionEvidence) bool {
	if before == nil || after == nil || len(after.Partitions) <= len(before.Partitions) {
		return false
	}
	return slices.Equal(before.Partitions[:len(before.Partitions)-1], after.Partitions[:len(before.Partitions)-1])
}
