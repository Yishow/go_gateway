package recordingplan

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// MaintainMySQLPartitions verifies and extends one confirmed native scope under
// the existing operation lease, settling its current projection atomically.
func (s *Service) MaintainMySQLPartitions(ctx context.Context, req MySQLPartitionMaintenanceRequest, target MySQLPartitionMaintenanceTarget) (MySQLPartitionMaintenanceResult, error) {
	if !validateMySQLPartitionMaintenanceScope(req.Scope) || strings.TrimSpace(req.GroupID) == "" || req.ConnectorRevision == "" || req.At.IsZero() || target.Inspect == nil || target.Reorganize == nil {
		return MySQLPartitionMaintenanceResult{}, ErrMySQLPartitionMaintenanceRequest
	}
	repo, ok := s.repo.(mysqlPartitionMaintenanceRepository)
	if !ok {
		return MySQLPartitionMaintenanceResult{}, ErrMySQLPartitionMaintenanceUnsupported
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		if errors.Is(err, ErrStorageLineageNotFound) {
			return maintenanceBlocked(nil, nil, "confirmation_missing"), nil
		}
		return MySQLPartitionMaintenanceResult{}, err
	}
	if lineage.RepairState != "" {
		return maintenanceBlocked(nil, nil, "storage_repair_blocked"), nil
	}
	baseDigest, err := validateMySQLPartitionMaintenanceLineage(lineage, req)
	if err != nil {
		return maintenanceBlocked(nil, nil, "confirmation_mismatch"), err
	}
	state, err := ReadMySQLPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	if state != nil {
		if !stateMatchesLineage(state, lineage, req, baseDigest) {
			return MySQLPartitionMaintenanceResult{}, ErrMySQLPartitionMaintenanceStateMismatch
		}
	}
	scopeKey := mysqlPartitionMaintenanceScopeKey(req.Scope, lineage.Layout)
	if scopeKey == "" {
		return maintenanceBlocked(state, nil, "confirmation_mismatch"), nil
	}
	if active, activeErr := s.repo.FindActiveSchemaOperation(ctx, req.Scope.WorkspaceID, scopeKey); activeErr == nil {
		if active.Action == MySQLPartitionMaintenanceAction && schemaOperationLeaseExpired(active) {
			resolved, err := s.resolveStaleMySQLPartitionMaintenance(ctx, repo, req, lineage, state, active, target, baseDigest)
			return resolved, err
		}
		return maintenanceBlocked(state, withoutOwner(*active), "scope_busy"), nil
	} else if !errors.Is(activeErr, ErrSchemaOperationNotFound) {
		return MySQLPartitionMaintenanceResult{}, fmt.Errorf("find active mysql partition operation: %w", activeErr)
	}

	expected := lineage.Layout.MySQL.NativeMonthly
	if state != nil {
		expected = state.Current
	}
	actual, err := inspectMySQLPartitionMaintenanceCatalog(ctx, repo.storageDB(), req, lineage, target, expected, nil)
	if err != nil {
		return s.persistMaintenanceBlock(ctx, repo, state, nil, maintenanceReason(err))
	}
	base := cloneMySQLPartitionEvidence(lineage.Layout.MySQL.NativeMonthly)
	if state == nil {
		if !sameMySQLPartitionEvidence(actual, base) {
			return maintenanceBlocked(nil, nil, mysqlPartitionReasonCatalogMismatch), nil
		}
		state = &MySQLPartitionMaintenanceState{
			Scope: req.Scope, BaseLayoutDigest: baseDigest, ConnectorRevision: req.ConnectorRevision,
			Generation: lineage.Generation, Current: base, Revision: 1, Status: mysqlPartitionMaintenanceReady,
			UpdatedAt: time.Now().UTC(),
		}
		if err := repo.ensureMySQLPartitionMaintenanceState(ctx, state); err != nil {
			return MySQLPartitionMaintenanceResult{}, err
		}
		state, err = ReadMySQLPartitionMaintenanceState(ctx, repo.storageDB(), req.Scope)
		if err != nil {
			return MySQLPartitionMaintenanceResult{}, err
		}
		if !stateMatchesLineage(state, lineage, req, baseDigest) || !sameMySQLPartitionEvidence(state.Current, base) {
			return MySQLPartitionMaintenanceResult{}, ErrMySQLPartitionMaintenanceStateMismatch
		}
	}
	if !sameMySQLPartitionEvidence(actual, state.Current) {
		return s.persistMaintenanceBlock(ctx, repo, state, nil, mysqlPartitionReasonCatalogMismatch)
	}
	after, err := PlanMySQLNativeMonthlyExtension(state.Current, req.At)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	if sameMySQLPartitionEvidence(after, state.Current) {
		if state.Status != mysqlPartitionMaintenanceReady || state.Reason != "" {
			state, err = persistMaintenanceState(ctx, repo.storageDB(), state, state.Current, state.LastOperationID, mysqlPartitionMaintenanceReady, "")
			if err != nil {
				return MySQLPartitionMaintenanceResult{}, err
			}
		}
		return maintenanceResult(state, nil, "", "", false), nil
	}
	operation, outcome, err := s.claimMySQLPartitionMaintenance(ctx, repo, req, lineage, state, after, baseDigest, scopeKey)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	if outcome != ClaimAcquired {
		return maintenanceBlocked(state, operation, "scope_busy"), nil
	}
	return s.runMySQLPartitionMaintenance(ctx, repo, req, lineage, state, operation, after, target)
}

func validateMySQLPartitionMaintenanceLineage(lineage *StorageLineage, req MySQLPartitionMaintenanceRequest) (string, error) {
	if lineage == nil || !lineage.Confirmed() || lineage.Generation < 1 || lineage.Layout == nil || lineage.Layout.MySQL == nil ||
		lineage.Layout.MySQL.NativeMonthly == nil || lineage.ConnectorRevision != req.ConnectorRevision ||
		lineage.StorageScope != req.Scope || lineage.Layout.TableName != lineage.TableName || lineage.Layout.OwnerColumn != lineage.OwnerColumn ||
		lineage.RepairState != "" {
		return "", ErrMySQLPartitionMaintenanceStateMismatch
	}
	table := "gw_long_" + mysqlStorageIdentityHash(req.Scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(req.Scope, "owner")[:24]
	canonical, err := ManagedMySQLLongGroupLayout(lineage.TableName, lineage.OwnerColumn, *lineage.Layout.MySQL)
	if err != nil || !groupSchemaLayoutsEqual(canonical, *lineage.Layout) || lineage.TableName != table ||
		lineage.OwnerColumn != owner || digestGroupSchemaLayout(*lineage.Layout) != lineage.LayoutDigest ||
		validateMySQLNativeMonthlyPartitionEvidence(lineage.Layout.MySQL.NativeMonthly) != nil {
		return "", ErrMySQLPartitionMaintenanceStateMismatch
	}
	return lineage.LayoutDigest, nil
}

func validateMySQLNativeMonthlyPartitionEvidence(evidence *MySQLPartitionEvidence) error {
	_, err := normalizeMySQLPartitionEvidence(evidence)
	return err
}

func stateMatchesLineage(state *MySQLPartitionMaintenanceState, lineage *StorageLineage, req MySQLPartitionMaintenanceRequest, baseDigest string) bool {
	if validateMySQLPartitionMaintenanceState(state) != nil || state.Scope != req.Scope ||
		state.BaseLayoutDigest != baseDigest || state.ConnectorRevision != req.ConnectorRevision || state.Generation != lineage.Generation {
		return false
	}
	base := lineage.Layout.MySQL.NativeMonthly
	if sameMySQLPartitionEvidence(state.Current, base) {
		return true
	}
	return len(state.Current.Partitions) >= len(base.Partitions) &&
		slices.Equal(state.Current.Partitions[:len(base.Partitions)-1], base.Partitions[:len(base.Partitions)-1])
}

func mysqlPartitionMaintenanceScopeKey(scope StorageScope, layout *GroupSchemaLayout) string {
	if layout == nil {
		return ""
	}
	return (&SchemaPreviewToken{WorkspaceID: scope.WorkspaceID, ConnectorID: scope.ConnectorID,
		Database: scope.Database, Schema: scope.TableSchema, GroupLayout: layout}).ScopeKey()
}

func inspectMySQLPartitionMaintenanceCatalog(ctx context.Context, q DBTX, req MySQLPartitionMaintenanceRequest, lineage *StorageLineage, target MySQLPartitionMaintenanceTarget, expected, pending *MySQLPartitionEvidence) (*MySQLPartitionEvidence, error) {
	group, err := target.Inspect(ctx, lineage.TableName)
	if err != nil {
		return nil, fmt.Errorf("%w: inspect mysql group table: %w", ErrTargetInspectionUnconfirmed, err)
	}
	receipt, err := target.Inspect(ctx, managedEffectReceiptTable)
	if err != nil {
		return nil, fmt.Errorf("%w: inspect mysql receipt table: %w", ErrTargetInspectionUnconfirmed, err)
	}
	if group.Status != inspectionExists || group.MySQL == nil || group.MySQL.Partitions == nil ||
		receipt.Status != inspectionExists || receipt.MySQL == nil || receipt.MySQL.Partitions != nil {
		return nil, fmt.Errorf("%w: mysql partition catalog evidence is incomplete", ErrTargetInspectionUnconfirmed)
	}
	actual := cloneMySQLPartitionEvidence(group.MySQL.Partitions)
	if err := validateCurrentMySQLPartitionEvidence(actual); err != nil {
		return nil, fmt.Errorf("%w: current partition evidence is not canonical: %w", ErrTargetInspectionUnconfirmed, err)
	}
	groupForNormalization := group
	if !sameMySQLPartitionEvidence(actual, expected) {
		if pending == nil || !sameMySQLPartitionEvidence(actual, pending) {
			return nil, fmt.Errorf("%w: actual partitions differ from the authorized before/after evidence", ErrIncompatibleExistingTable)
		}
		groupForNormalization = cloneTargetTableInspection(group)
		groupForNormalization.MySQL.Partitions = cloneMySQLPartitionEvidence(expected)
	}
	normalized, err := NormalizeMySQLPartitionCatalog(ctx, q, req.Scope, *lineage.Layout, groupForNormalization)
	if err != nil {
		return nil, fmt.Errorf("%w: normalize mysql partition catalog: %w", ErrTargetInspectionUnconfirmed, err)
	}
	digest, err := VerifyManagedMySQLCatalog(*lineage.Layout, req.Scope.Database, normalized, receipt)
	if err != nil || digest != lineage.ConfirmedDigest {
		if err != nil {
			return nil, fmt.Errorf("%w: confirmed mysql catalog facts changed: %w", ErrTargetInspectionUnconfirmed, err)
		}
		return nil, fmt.Errorf("%w: confirmed mysql catalog digest changed", ErrTargetInspectionUnconfirmed)
	}
	return actual, nil
}

func maintenanceReason(err error) string {
	switch {
	case errors.Is(err, ErrTargetPermissionDenied):
		return "permission_lost"
	case errors.Is(err, ErrIncompatibleExistingTable):
		return mysqlPartitionReasonCatalogMismatch
	default:
		return mysqlPartitionReasonUnverifiable
	}
}

func maintenanceBlocked(state *MySQLPartitionMaintenanceState, op *SchemaOperation, reason string) MySQLPartitionMaintenanceResult {
	return maintenanceResult(state, op, mysqlPartitionMaintenanceCode, reason, false)
}

func maintenanceResult(state *MySQLPartitionMaintenanceState, op *SchemaOperation, code, reason string, changed bool) MySQLPartitionMaintenanceResult {
	var stateCopy *MySQLPartitionMaintenanceState
	if state != nil {
		stateCopyValue := *state
		stateCopyValue.Current = cloneMySQLPartitionEvidence(state.Current)
		stateCopy = &stateCopyValue
	}
	if op != nil {
		op = withoutOwner(*op)
	}
	return MySQLPartitionMaintenanceResult{State: stateCopy, Operation: op, Code: code, Reason: reason, Changed: changed}
}

func (s *Service) persistMaintenanceBlock(ctx context.Context, repo mysqlPartitionMaintenanceRepository, state *MySQLPartitionMaintenanceState, op *SchemaOperation, reason string) (MySQLPartitionMaintenanceResult, error) {
	if state == nil {
		return maintenanceBlocked(nil, op, reason), nil
	}
	updated, err := persistMaintenanceState(ctx, repo.storageDB(), state, state.Current,
		state.LastOperationID, mysqlPartitionMaintenanceBlocked, reason)
	if err != nil {
		return MySQLPartitionMaintenanceResult{}, err
	}
	return maintenanceBlocked(updated, op, reason), nil
}

func persistMaintenanceState(ctx context.Context, q DBTX, previous *MySQLPartitionMaintenanceState, current *MySQLPartitionEvidence, lastOperationID, status, reason string) (*MySQLPartitionMaintenanceState, error) {
	next := *previous
	next.Current = cloneMySQLPartitionEvidence(current)
	next.Revision++
	next.LastOperationID = lastOperationID
	next.Status, next.Reason, next.UpdatedAt = status, reason, time.Now().UTC()
	if err := updateMySQLPartitionMaintenanceState(ctx, q, previous.Revision, &next); err != nil {
		return nil, err
	}
	return &next, nil
}
