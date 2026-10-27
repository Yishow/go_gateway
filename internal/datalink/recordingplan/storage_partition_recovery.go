package recordingplan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *Service) resolveStaleMySQLPartitionMaintenance(ctx context.Context, repo mysqlPartitionMaintenanceRepository, req MySQLPartitionMaintenanceRequest, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, stale *SchemaOperation, target MySQLPartitionMaintenanceTarget, baseDigest string) (MySQLPartitionMaintenanceResult, error) {
	var detail mysqlPartitionMaintenanceDetail
	if json.Unmarshal([]byte(stale.Detail), &detail) != nil || !validMySQLPartitionMaintenanceDetail(detail, req, lineage, state, baseDigest) {
		return s.finishStalePartitionMaintenanceUnknown(ctx, repo, req, state, stale, mysqlPartitionReasonUnverifiable)
	}
	observed, err := inspectMySQLPartitionMaintenanceCatalog(ctx, repo.storageDB(), req, lineage, target, detail.Before, detail.After)
	now := time.Now().UTC()
	cutoff := now.Add(-SchemaOperationLease)
	var result SchemaOperationResult
	var current *MySQLPartitionEvidence
	status, reason, changed := mysqlPartitionMaintenanceBlocked, mysqlPartitionReasonUnverifiable, false
	switch {
	case err == nil && sameMySQLPartitionEvidence(observed, detail.After):
		result = SchemaOperationResult{Status: SchemaOperationSucceeded, ExecutedStatements: 1, VerifiedDigest: mysqlPartitionEvidenceDigest(detail.After), Reason: "recovered"}
		current, status, reason, changed = detail.After, mysqlPartitionMaintenanceReady, "", true
	case err == nil && sameMySQLPartitionEvidence(observed, detail.Before):
		result = SchemaOperationResult{Status: SchemaOperationFailed, Reason: "attempt_abandoned", NextAction: mysqlPartitionMaintenanceNext}
		current, reason = detail.Before, "attempt_abandoned"
	default:
		result = SchemaOperationResult{Status: SchemaOperationUnknown, Reason: mysqlPartitionReasonUnverifiable, NextAction: mysqlPartitionMaintenanceNext}
		if state != nil {
			current = state.Current
		}
	}
	result.Detail = stale.Detail
	var apply func(context.Context, DBTX) error
	if state != nil {
		apply = func(ctx context.Context, q DBTX) error {
			next := *state
			next.Current = cloneMySQLPartitionEvidence(current)
			next.Revision++
			next.LastOperationID = stale.OperationID
			next.Status, next.Reason, next.UpdatedAt = status, reason, now
			return updateMySQLPartitionMaintenanceState(ctx, q, state.Revision, &next)
		}
	}
	finished, finishErr := repo.finishStaleMySQLPartitionMaintenance(ctx, stale.OperationID, result, cutoff, apply)
	if finishErr != nil {
		return MySQLPartitionMaintenanceResult{}, fmt.Errorf("resolve stale mysql partition maintenance: %w", finishErr)
	}
	updated, readErr := ReadMySQLPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if readErr != nil {
		return MySQLPartitionMaintenanceResult{}, readErr
	}
	code := ""
	if result.Status != SchemaOperationSucceeded {
		code = mysqlPartitionMaintenanceCode
	}
	return maintenanceResult(updated, finished, code, reason, changed), nil
}

func (s *Service) finishStalePartitionMaintenanceUnknown(ctx context.Context, repo mysqlPartitionMaintenanceRepository, req MySQLPartitionMaintenanceRequest, state *MySQLPartitionMaintenanceState, stale *SchemaOperation, reason string) (MySQLPartitionMaintenanceResult, error) {
	now := time.Now().UTC()
	result := SchemaOperationResult{Status: SchemaOperationUnknown, Reason: reason, NextAction: mysqlPartitionMaintenanceNext, Detail: stale.Detail}
	var apply func(context.Context, DBTX) error
	if state != nil {
		apply = func(ctx context.Context, q DBTX) error {
			next := *state
			next.Revision++
			next.LastOperationID = stale.OperationID
			next.Status, next.Reason, next.UpdatedAt = mysqlPartitionMaintenanceBlocked, reason, now
			return updateMySQLPartitionMaintenanceState(ctx, q, state.Revision, &next)
		}
	}
	finished, err := repo.finishStaleMySQLPartitionMaintenance(ctx, stale.OperationID, result, now.Add(-SchemaOperationLease), apply)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	updated, err := ReadMySQLPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	return maintenanceResult(updated, finished, mysqlPartitionMaintenanceCode, reason, false), nil
}

func validMySQLPartitionMaintenanceDetail(detail mysqlPartitionMaintenanceDetail, req MySQLPartitionMaintenanceRequest, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, baseDigest string) bool {
	if detail.Kind != mysqlPartitionMaintenanceKind || detail.Scope != req.Scope || detail.GroupID == "" ||
		detail.ConnectorRevision != req.ConnectorRevision || detail.BaseLayoutDigest != baseDigest ||
		detail.Generation != lineage.Generation || state == nil || detail.StateRevision != state.Revision ||
		!sameMySQLPartitionEvidence(detail.Before, state.Current) || validateCurrentMySQLPartitionEvidence(detail.After) != nil {
		return false
	}
	planned, err := PlanMySQLNativeMonthlyExtension(detail.Before, detail.At)
	return err == nil && sameMySQLPartitionEvidence(planned, detail.After) && sameMySQLPartitionPrefix(detail.Before, detail.After)
}
