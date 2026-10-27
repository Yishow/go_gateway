package recordingplan

import (
	"context"
	"fmt"
	"strings"
)

// NormalizeMySQLPartitionCatalog projects a ledger-authorized current catalog
// back to its immutable confirmation for runtime and automatic-repair schema
// verification. It never adopts an unrecorded external partition change or
// mutates the supplied inspection.
func NormalizeMySQLPartitionCatalog(ctx context.Context, q DBTX, scope StorageScope, base GroupSchemaLayout, actual TargetTableInspection) (TargetTableInspection, error) {
	return normalizeMySQLPartitionCatalog(ctx, q, scope, base, actual, "", false)
}

// NormalizeMySQLPartitionCatalogForExplicitReconfirmation is the explicit
// preview/apply path. It may bridge a verified capability/server projection
// change to the old confirmed partition shape; it is never used by runtime or
// automatic repair.
func NormalizeMySQLPartitionCatalogForExplicitReconfirmation(ctx context.Context, q DBTX, scope StorageScope, base GroupSchemaLayout, actual TargetTableInspection, currentConnectorRevision string) (TargetTableInspection, error) {
	return normalizeMySQLPartitionCatalog(ctx, q, scope, base, actual, currentConnectorRevision, true)
}

func normalizeMySQLPartitionCatalog(ctx context.Context, q DBTX, scope StorageScope, base GroupSchemaLayout, actual TargetTableInspection, currentConnectorRevision string, explicit bool) (TargetTableInspection, error) {
	if q == nil || !validateMySQLPartitionMaintenanceScope(scope) || base.MySQL == nil || base.MySQL.NativeMonthly == nil ||
		actual.Status != inspectionExists || actual.MySQL == nil {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	canonical, err := ManagedMySQLLongGroupLayout(base.TableName, base.OwnerColumn, *base.MySQL)
	if err != nil || !groupSchemaLayoutsEqual(canonical, base) {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	lineage, err := ReadStorageLineage(ctx, q, scope)
	if err != nil {
		return actual, err
	}
	state, err := ReadMySQLPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return actual, err
	}
	if lineage.Layout == nil {
		// A same-operation partial CREATE may exist before the first storage
		// confirmation. It can only prove the exact initial layout, never an extension.
		table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
		owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]
		if state != nil || lineage.ConfirmedOperationID != "" || lineage.LayoutDigest != UnconfirmedMySQLStorageDigest(scope) ||
			lineage.TableName != table || lineage.OwnerColumn != owner || base.TableName != table || base.OwnerColumn != owner ||
			!sameMySQLPartitionEvidence(actual.MySQL.Partitions, base.MySQL.NativeMonthly) {
			return actual, ErrMySQLPartitionMaintenanceStateMismatch
		}
		return cloneTargetTableInspection(actual), nil
	}
	req := MySQLPartitionMaintenanceRequest{Scope: scope, ConnectorRevision: lineage.ConnectorRevision}
	digest, err := validateMySQLPartitionMaintenanceLineage(lineage, req)
	if err != nil {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	baseDigest := digestGroupSchemaLayout(base)
	if !explicit {
		return normalizeMySQLPartitionAuthorized(base, actual, lineage, state, req, digest)
	}
	if !sameMySQLActualCapability(actual.MySQL, base.MySQL) {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	if baseDigest == digest {
		return normalizeMySQLPartitionExplicit(base, actual, lineage, state, req, digest)
	}
	if strings.TrimSpace(currentConnectorRevision) == "" {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	// This is the only cross-binding path. The old confirmed layout and the new
	// verified layout must describe the same physical schema; only the server
	// version/capability projection may differ within the same family.
	if !sameMySQLPhysicalLayoutExceptProjection(*lineage.Layout, base) {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	return normalizeMySQLPartitionExplicit(base, actual, lineage, state, req, digest)
}

func normalizeMySQLPartitionAuthorized(base GroupSchemaLayout, actual TargetTableInspection, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, req MySQLPartitionMaintenanceRequest, digest string) (TargetTableInspection, error) {
	if digestGroupSchemaLayout(base) != digest {
		return actual, ErrMySQLPartitionMaintenanceStateMismatch
	}
	expected := base.MySQL.NativeMonthly
	if state != nil {
		if !stateMatchesLineage(state, lineage, req, digest) ||
			!sameMySQLPartitionEvidence(expected, state.Current) && !sameMySQLPartitionPrefix(expected, state.Current) {
			return actual, ErrMySQLPartitionMaintenanceStateMismatch
		}
		expected = state.Current
	}
	if !sameMySQLPartitionEvidence(actual.MySQL.Partitions, expected) {
		return actual, fmt.Errorf("%w: actual partitions differ from authorized current state", ErrIncompatibleExistingTable)
	}
	normalized := cloneTargetTableInspection(actual)
	normalized.MySQL.Partitions = cloneMySQLPartitionEvidence(base.MySQL.NativeMonthly)
	return normalized, nil
}

func normalizeMySQLPartitionExplicit(base GroupSchemaLayout, actual TargetTableInspection, lineage *StorageLineage, state *MySQLPartitionMaintenanceState, req MySQLPartitionMaintenanceRequest, digest string) (TargetTableInspection, error) {
	expected := base.MySQL.NativeMonthly
	if state != nil {
		if !stateMatchesLineage(state, lineage, req, digest) ||
			!sameMySQLPartitionEvidence(actual.MySQL.Partitions, state.Current) {
			return actual, fmt.Errorf("%w: actual partitions differ from authorized current state", ErrIncompatibleExistingTable)
		}
	} else if !sameMySQLPartitionEvidence(actual.MySQL.Partitions, expected) {
		return actual, fmt.Errorf("%w: actual partitions differ from confirmed initial state", ErrIncompatibleExistingTable)
	}
	normalized := cloneTargetTableInspection(actual)
	normalized.MySQL.Partitions = cloneMySQLPartitionEvidence(base.MySQL.NativeMonthly)
	return normalized, nil
}

func sameMySQLActualCapability(actual *TargetMySQLTableEvidence, options *MySQLGroupSchemaOptions) bool {
	return actual != nil && options != nil && actual.ServerFamily == options.ServerFamily &&
		actual.ServerVersion == options.ServerVersion && actual.CapabilityRevision == options.CapabilityRevision
}

// sameMySQLPhysicalLayoutExceptProjection permits the explicit preview path to
// retain the old partition authorization while rebinding the new verified
// capability/server projection. Everything that defines the table's physical
// shape and initial partitions must remain identical.
func sameMySQLPhysicalLayoutExceptProjection(old, next GroupSchemaLayout) bool {
	if old.MySQL == nil || next.MySQL == nil || old.MySQL.ServerFamily != next.MySQL.ServerFamily {
		return false
	}
	projected := next
	options := *next.MySQL
	options.ServerVersion = old.MySQL.ServerVersion
	options.CapabilityRevision = old.MySQL.CapabilityRevision
	projected.MySQL = &options
	return groupSchemaLayoutsEqual(old, projected)
}
