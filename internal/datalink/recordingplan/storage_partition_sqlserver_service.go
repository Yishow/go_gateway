package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// MaintainSQLServerPartitions verifies and extends one confirmed SQL Server
// native scope under the existing operation lease and durable projection.
func (s *Service) MaintainSQLServerPartitions(ctx context.Context, req SQLServerPartitionMaintenanceRequest, target SQLServerPartitionMaintenanceTarget) (SQLServerPartitionMaintenanceResult, error) {
	if !validateSQLServerPartitionMaintenanceScope(req.Scope) || req.GroupID == "" || req.ConnectorRevision == "" ||
		req.At.IsZero() || target.Inspect == nil || target.Split == nil {
		return SQLServerPartitionMaintenanceResult{}, ErrSQLServerPartitionMaintenanceRequest
	}
	repo, ok := s.repo.(sqlServerPartitionMaintenanceRepository)
	if !ok {
		return SQLServerPartitionMaintenanceResult{}, ErrSQLServerPartitionMaintenanceUnsupported
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		if errors.Is(err, ErrStorageLineageNotFound) {
			return sqlServerMaintenanceBlocked(nil, nil, "confirmation_missing"), nil
		}
		return SQLServerPartitionMaintenanceResult{}, err
	}
	if lineage.RepairState != "" {
		return sqlServerMaintenanceBlocked(nil, nil, "storage_repair_blocked"), nil
	}
	baseDigest, err := validateSQLServerPartitionMaintenanceLineage(lineage, req)
	if err != nil {
		return sqlServerMaintenanceBlocked(nil, nil, "confirmation_mismatch"), err
	}
	state, err := ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	if state != nil && !sqlServerStateMatchesLineage(state, lineage, req, baseDigest) {
		return SQLServerPartitionMaintenanceResult{}, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	scopeKey := sqlServerPartitionMaintenanceScopeKey(req.Scope, lineage.Layout)
	if scopeKey == "" {
		return sqlServerMaintenanceBlocked(state, nil, "confirmation_mismatch"), nil
	}
	if active, activeErr := s.repo.FindActiveSchemaOperation(ctx, req.Scope.WorkspaceID, scopeKey); activeErr == nil {
		if active.Action == SQLServerPartitionMaintenanceAction && schemaOperationLeaseExpired(active) {
			return s.resolveStaleSQLServerPartitionMaintenance(ctx, repo, req, lineage, state, active, target, baseDigest)
		}
		return sqlServerMaintenanceBlocked(state, withoutOwner(*active), "scope_busy"), nil
	} else if !errors.Is(activeErr, ErrSchemaOperationNotFound) {
		return SQLServerPartitionMaintenanceResult{}, fmt.Errorf("find active sqlserver partition operation: %w", activeErr)
	}
	if recovered, err := s.recoverUnknownSQLServerPartitionMaintenance(ctx, repo, req, lineage, state, target, baseDigest, scopeKey); err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	} else if recovered != nil {
		return *recovered, nil
	}

	base := CloneSQLServerPartitionEvidence(lineage.Layout.SQLServer.NativeMonthly)
	expected := lineage.Layout.SQLServer.NativeMonthly
	if state != nil {
		expected = state.Current
	}
	actual, err := inspectSQLServerPartitionMaintenanceCatalog(ctx, req, lineage, target, expected, nil)
	if err != nil {
		inspectErr := err
		if state == nil {
			state = &SQLServerPartitionMaintenanceState{
				Scope: req.Scope, BaseLayoutDigest: baseDigest, ConnectorRevision: req.ConnectorRevision,
				Generation: lineage.Generation, Current: base, Revision: 1, Status: sqlServerPartitionMaintenanceReady,
				UpdatedAt: time.Now().UTC(),
			}
			if ensureErr := repo.ensureSQLServerPartitionMaintenanceState(ctx, state); ensureErr != nil {
				return SQLServerPartitionMaintenanceResult{}, ensureErr
			}
			state, err = ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
			if err != nil {
				return SQLServerPartitionMaintenanceResult{}, err
			}
		}
		return s.persistSQLServerMaintenanceBlock(ctx, repo, state, nil, sqlServerMaintenanceReason(inspectErr))
	}
	if state == nil {
		if !sameSQLServerPartitionEvidence(actual, base) {
			return sqlServerMaintenanceBlocked(nil, nil, sqlServerPartitionReasonCatalogMismatch), nil
		}
		state = &SQLServerPartitionMaintenanceState{
			Scope: req.Scope, BaseLayoutDigest: baseDigest, ConnectorRevision: req.ConnectorRevision,
			Generation: lineage.Generation, Current: base, Revision: 1, Status: sqlServerPartitionMaintenanceReady,
			UpdatedAt: time.Now().UTC(),
		}
		if err := repo.ensureSQLServerPartitionMaintenanceState(ctx, state); err != nil {
			return SQLServerPartitionMaintenanceResult{}, err
		}
		state, err = ReadSQLServerPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
		if err != nil {
			return SQLServerPartitionMaintenanceResult{}, err
		}
		if !sqlServerStateMatchesLineage(state, lineage, req, baseDigest) || !sameSQLServerPartitionEvidence(state.Current, base) {
			return SQLServerPartitionMaintenanceResult{}, ErrSQLServerPartitionMaintenanceStateMismatch
		}
	}
	if !sameSQLServerPartitionEvidence(actual, state.Current) {
		return s.persistSQLServerMaintenanceBlock(ctx, repo, state, nil, sqlServerPartitionReasonCatalogMismatch)
	}
	after, err := PlanSQLServerNativeMonthlyExtension(state.Current, req.At)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	if sameSQLServerPartitionEvidence(after, state.Current) {
		if state.Status != sqlServerPartitionMaintenanceReady || state.Reason != "" {
			state, err = persistSQLServerMaintenanceState(ctx, repo.storageDB(), state, state.Current, state.LastOperationID, sqlServerPartitionMaintenanceReady, "")
			if err != nil {
				return SQLServerPartitionMaintenanceResult{}, err
			}
		}
		return sqlServerMaintenanceResult(state, nil, "", "", false), nil
	}
	operation, outcome, err := s.claimSQLServerPartitionMaintenance(ctx, repo, req, lineage, state, after, baseDigest, scopeKey)
	if err != nil {
		return SQLServerPartitionMaintenanceResult{}, err
	}
	if outcome != ClaimAcquired {
		return sqlServerMaintenanceBlocked(state, operation, "scope_busy"), nil
	}
	return s.runSQLServerPartitionMaintenance(ctx, repo, req, lineage, state, operation, after, target)
}

func validateSQLServerPartitionMaintenanceLineage(lineage *StorageLineage, req SQLServerPartitionMaintenanceRequest) (string, error) {
	if lineage == nil || !lineage.Confirmed() || lineage.Generation < 1 || lineage.Layout == nil || lineage.Layout.SQLServer == nil ||
		lineage.Layout.SQLServer.NativeMonthly == nil || lineage.ConnectorRevision != req.ConnectorRevision ||
		lineage.StorageScope != req.Scope || lineage.Layout.TableName != lineage.TableName ||
		lineage.Layout.OwnerColumn != lineage.OwnerColumn || lineage.RepairState != "" {
		return "", ErrSQLServerPartitionMaintenanceStateMismatch
	}
	table := "gw_long_" + mysqlStorageIdentityHash(req.Scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(req.Scope, "owner")[:24]
	canonical, err := ManagedSQLServerLongGroupLayout(lineage.TableName, lineage.OwnerColumn, *lineage.Layout.SQLServer)
	if err != nil || !groupSchemaLayoutsEqual(canonical, *lineage.Layout) || lineage.TableName != table ||
		lineage.OwnerColumn != owner || digestGroupSchemaLayout(*lineage.Layout) != lineage.LayoutDigest ||
		validateSQLServerPartitionEvidence(lineage.Layout.SQLServer.NativeMonthly, lineage.TableName) != nil {
		return "", ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return lineage.LayoutDigest, nil
}

func sqlServerStateMatchesLineage(state *SQLServerPartitionMaintenanceState, lineage *StorageLineage, req SQLServerPartitionMaintenanceRequest, baseDigest string) bool {
	if validateSQLServerPartitionMaintenanceState(state) != nil || state.Scope != req.Scope ||
		state.BaseLayoutDigest != baseDigest || state.ConnectorRevision != req.ConnectorRevision || state.Generation != lineage.Generation {
		return false
	}
	base := lineage.Layout.SQLServer.NativeMonthly
	if sameSQLServerPartitionEvidence(state.Current, base) {
		return true
	}
	return sameSQLServerPartitionPrefix(base, state.Current)
}

func sqlServerPartitionMaintenanceScopeKey(scope StorageScope, layout *GroupSchemaLayout) string {
	if layout == nil || layout.SQLServer == nil {
		return ""
	}
	return (&SchemaPreviewToken{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID,
		Database: scope.Database, Schema: scope.TableSchema, GroupLayout: layout}).ScopeKey()
}

func sqlServerMaintenanceReason(err error) string {
	switch {
	case errors.Is(err, ErrTargetPermissionDenied):
		return sqlServerPartitionReasonPermissionLost
	case errors.Is(err, ErrIncompatibleExistingTable):
		return sqlServerPartitionReasonCatalogMismatch
	default:
		return sqlServerPartitionReasonUnverifiable
	}
}

func sqlServerMaintenanceBlocked(state *SQLServerPartitionMaintenanceState, op *SchemaOperation, reason string) SQLServerPartitionMaintenanceResult {
	return sqlServerMaintenanceResult(state, op, sqlServerPartitionMaintenanceCode, reason, false)
}

func sqlServerMaintenanceResult(state *SQLServerPartitionMaintenanceState, op *SchemaOperation, code, reason string, changed bool) SQLServerPartitionMaintenanceResult {
	var stateCopy *SQLServerPartitionMaintenanceState
	if state != nil {
		value := *state
		value.Current = CloneSQLServerPartitionEvidence(state.Current)
		stateCopy = &value
	}
	if op != nil {
		op = withoutOwner(*op)
	}
	return SQLServerPartitionMaintenanceResult{State: stateCopy, Operation: op, Code: code, Reason: reason, Changed: changed}
}
