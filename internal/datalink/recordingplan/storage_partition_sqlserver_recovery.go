package recordingplan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *Service) resolveStaleSQLServerPartitionMaintenance(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, state *SQLServerPartitionMaintenanceState, stale *SchemaOperation, target SQLServerPartitionMaintenanceTarget, baseDigest string) (SQLServerPartitionMaintenanceResult, error) {
	var detail sqlServerPartitionMaintenanceDetail
	if json.Unmarshal([]byte(stale.Detail), &detail) != nil ||
		!validSQLServerPartitionMaintenanceDetail(detail, req, lineage, state, baseDigest) {
		return s.finishStaleSQLServerMaintenanceUnknown(ctx, repo, req, state, stale, sqlServerPartitionReasonUnverifiable)
	}
	observed, err := inspectSQLServerPartitionMaintenanceCatalog(ctx, req, lineage, target, detail.Before, detail.After)
	now := time.Now().UTC()
	cutoff := now.Add(-SchemaOperationLease)
	var result SchemaOperationResult
	var current *SQLServerPartitionEvidence
	status, reason, changed := sqlServerPartitionMaintenanceBlocked, sqlServerPartitionReasonUnverifiable, false
	switch {
	case err == nil && sameSQLServerPartitionEvidence(observed, detail.After):
		result = SchemaOperationResult{Status: SchemaOperationSucceeded, ExecutedStatements: 1, VerifiedDigest: sqlServerPartitionEvidenceDigest(detail.After), Reason: "recovered"}
		current, status, reason, changed = detail.After, sqlServerPartitionMaintenanceReady, "", true
	case err == nil && sameSQLServerPartitionEvidence(observed, detail.Before):
		result = SchemaOperationResult{Status: SchemaOperationFailed, Reason: sqlServerPartitionReasonAttemptAbandoned, NextAction: sqlServerPartitionMaintenanceNext}
		current, reason = detail.Before, sqlServerPartitionReasonAttemptAbandoned
	default:
		result = SchemaOperationResult{Status: SchemaOperationUnknown, Reason: sqlServerPartitionReasonUnverifiable, NextAction: sqlServerPartitionMaintenanceNext}
		if state != nil {
			current = state.Current
		}
	}
	result.Detail = stale.Detail
	var apply func(context.Context, DBTX) error
	if state != nil {
		apply = func(ctx context.Context, q DBTX) error {
			next := *state
			next.Current = CloneSQLServerPartitionEvidence(current)
			next.Revision++
			next.LastOperationID = stale.OperationID
			next.Status, next.Reason, next.UpdatedAt = status, reason, now
			return updateSQLServerPartitionMaintenanceState(ctx, q, state.Revision, &next)
		}
	}
	finished, finishErr := repo.finishStaleSQLServerPartitionMaintenance(ctx, stale.OperationID, result, cutoff, apply)
	if finishErr != nil {
		return SQLServerPartitionMaintenanceResult{}, fmt.Errorf("resolve stale sqlserver partition maintenance: %w", finishErr)
	}
	updated, readErr := ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if readErr != nil {
		return SQLServerPartitionMaintenanceResult{}, readErr
	}
	code := ""
	if result.Status != SchemaOperationSucceeded {
		code = sqlServerPartitionMaintenanceCode
	}
	return sqlServerMaintenanceResult(updated, finished, code, reason, changed), nil
}

func (s *Service) finishStaleSQLServerMaintenanceUnknown(ctx context.Context, repo sqlServerPartitionMaintenanceRepository, req SQLServerPartitionMaintenanceRequest, state *SQLServerPartitionMaintenanceState, stale *SchemaOperation, reason string) (SQLServerPartitionMaintenanceResult, error) {
	now := time.Now().UTC()
	result := SchemaOperationResult{Status: SchemaOperationUnknown, Reason: reason, NextAction: sqlServerPartitionMaintenanceNext, Detail: stale.Detail}
	var apply func(context.Context, DBTX) error
	if state != nil {
		apply = func(ctx context.Context, q DBTX) error {
			next := *state
			next.Revision++
			next.LastOperationID = stale.OperationID
			next.Status, next.Reason, next.UpdatedAt = sqlServerPartitionMaintenanceBlocked, reason, now
			return updateSQLServerPartitionMaintenanceState(ctx, q, state.Revision, &next)
		}
	}
	finished, err := repo.finishStaleSQLServerPartitionMaintenance(ctx, stale.OperationID, result, now.Add(-SchemaOperationLease), apply)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	updated, err := ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	return sqlServerMaintenanceResult(updated, finished, sqlServerPartitionMaintenanceCode, reason, false), nil
}

func validSQLServerPartitionMaintenanceDetail(detail sqlServerPartitionMaintenanceDetail, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, state *SQLServerPartitionMaintenanceState, baseDigest string) bool {
	if detail.Kind != sqlServerPartitionMaintenanceKind || detail.Scope != req.Scope || detail.GroupID == "" ||
		detail.ConnectorRevision != req.ConnectorRevision || detail.BaseLayoutDigest != baseDigest ||
		detail.Generation != lineage.Generation || state == nil || detail.StateRevision != state.Revision ||
		!sameSQLServerPartitionEvidence(detail.Before, state.Current) || validateCurrentSQLServerPartitionEvidence(detail.After) != nil {
		return false
	}
	planned, err := PlanSQLServerNativeMonthlyExtension(detail.Before, detail.At)
	return err == nil && sameSQLServerPartitionEvidence(planned, detail.After) && sameSQLServerPartitionPrefix(detail.Before, detail.After)
}
