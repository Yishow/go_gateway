package recordingplan

import (
	"context"
	"fmt"
	"maps"
)

// NormalizeSQLServerPartitionCatalog projects a ledger-authorized current
// SQL Server partition projection back to the immutable confirmed layout. It
// never suppresses arbitrary catalog drift.
func NormalizeSQLServerPartitionCatalog(ctx context.Context, q DBTX, scope StorageScope, base GroupSchemaLayout, actual, receipt TargetTableInspection) (TargetTableInspection, error) {
	if q == nil || base.SQLServer == nil || base.SQLServer.NativeMonthly == nil || actual.Status != inspectionExists ||
		actual.SQLServer == nil {
		return actual, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	lineage, err := ReadStorageLineage(ctx, q, scope)
	if err != nil {
		return actual, err
	}
	if lineage == nil || !lineage.Confirmed() || lineage.Layout == nil || !groupSchemaLayoutsEqual(*lineage.Layout, base) ||
		lineage.ConfirmedDigest == "" {
		return actual, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	current, err := sqlServerPartitionEvidenceFromTarget(actual.SQLServer)
	if err != nil {
		return actual, err
	}
	if !sameSQLServerPartitionEvidence(current, base.SQLServer.NativeMonthly) {
		state, err := ReadSQLServerPartitionMaintenanceState(ctx, q, scope)
		if err != nil {
			return actual, err
		}
		if state == nil || !sqlServerStateMatchesLineage(state, lineage, SQLServerPartitionMaintenanceRequest{
			Scope: scope, ConnectorRevision: lineage.ConnectorRevision,
		}, lineage.LayoutDigest) || !sameSQLServerPartitionEvidence(current, state.Current) {
			return actual, ErrSQLServerPartitionMaintenanceStateMismatch
		}
	}
	normalized, err := normalizeSQLServerPartitionTarget(actual, base.SQLServer.NativeMonthly)
	if err != nil {
		return actual, err
	}
	normalizedReceipt, err := normalizeSQLServerPartitionTarget(receipt, nil)
	if err != nil {
		return actual, err
	}
	digest, err := VerifyManagedSQLServerCatalog(base, scope.Database, normalized, normalizedReceipt)
	if err != nil || digest != lineage.ConfirmedDigest {
		if err != nil {
			return actual, err
		}
		return actual, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return normalized, nil
}

func inspectSQLServerPartitionMaintenanceCatalog(ctx context.Context, req SQLServerPartitionMaintenanceRequest, lineage *StorageLineage, target SQLServerPartitionMaintenanceTarget, expected, pending *SQLServerPartitionEvidence) (*SQLServerPartitionEvidence, error) {
	if target.Inspect == nil || lineage == nil || lineage.Layout == nil || lineage.Layout.SQLServer == nil || lineage.Layout.SQLServer.NativeMonthly == nil {
		return nil, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	group, err := target.Inspect(ctx, lineage.TableName)
	if err != nil {
		return nil, fmt.Errorf("%w: inspect sqlserver group table: %w", ErrTargetInspectionUnconfirmed, err)
	}
	receipt, err := target.Inspect(ctx, managedEffectReceiptTable)
	if err != nil {
		return nil, fmt.Errorf("%w: inspect sqlserver receipt table: %w", ErrTargetInspectionUnconfirmed, err)
	}
	if group.Status != inspectionExists || group.SQLServer == nil {
		return nil, fmt.Errorf("%w: sqlserver group catalog evidence is incomplete", ErrTargetInspectionUnconfirmed)
	}
	if receipt.Status != inspectionExists || receipt.SQLServer == nil || receipt.SQLServer.NativeMonthly != nil {
		return nil, fmt.Errorf("%w: sqlserver receipt table must remain plain", ErrIncompatibleExistingTable)
	}
	actual, err := sqlServerPartitionEvidenceFromTarget(group.SQLServer)
	if err != nil || actual == nil {
		return nil, fmt.Errorf("%w: sqlserver partition evidence is incomplete", ErrTargetInspectionUnconfirmed)
	}
	base := lineage.Layout.SQLServer.NativeMonthly
	if actual.FunctionName != base.FunctionName || actual.SchemeName != base.SchemeName || actual.Filegroup != base.Filegroup {
		return nil, fmt.Errorf("%w: sqlserver partition object scope differs", ErrIncompatibleExistingTable)
	}
	if err := validateCurrentSQLServerPartitionEvidence(actual); err != nil {
		return nil, fmt.Errorf("%w: current sqlserver partition evidence is not canonical: %w", ErrTargetInspectionUnconfirmed, err)
	}
	if !sameSQLServerPartitionEvidence(actual, expected) && (pending == nil || !sameSQLServerPartitionEvidence(actual, pending)) {
		return nil, fmt.Errorf("%w: actual sqlserver partitions differ from authorized before/after evidence", ErrIncompatibleExistingTable)
	}
	normalized, err := normalizeSQLServerPartitionTarget(group, base)
	if err != nil {
		return nil, fmt.Errorf("%w: normalize sqlserver partition catalog: %w", ErrTargetInspectionUnconfirmed, err)
	}
	normalizedReceipt, err := normalizeSQLServerPartitionTarget(receipt, nil)
	if err != nil {
		return nil, err
	}
	digest, err := VerifyManagedSQLServerCatalog(*lineage.Layout, req.Scope.Database, normalized, normalizedReceipt)
	if err != nil || digest != lineage.ConfirmedDigest {
		if err != nil {
			return nil, fmt.Errorf("%w: confirmed sqlserver catalog facts changed: %w", ErrTargetInspectionUnconfirmed, err)
		}
		return nil, fmt.Errorf("%w: confirmed sqlserver catalog digest changed", ErrTargetInspectionUnconfirmed)
	}
	return actual, nil
}

func sqlServerPartitionEvidenceFromTarget(evidence *TargetSQLServerTableEvidence) (*SQLServerPartitionEvidence, error) {
	if evidence == nil || !evidence.PartitionInspected || evidence.NativeMonthly == nil {
		return nil, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	native := evidence.NativeMonthly
	return &SQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName,
		SchemeName: native.SchemeName, Filegroup: native.Filegroup,
		Partitions: append([]SQLServerPartitionBoundary(nil), native.Partitions...),
	}, nil
}

func normalizeSQLServerPartitionTarget(state TargetTableInspection, base *SQLServerPartitionEvidence) (TargetTableInspection, error) {
	if state.SQLServer == nil {
		return state, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	cloned := state
	evidence := *state.SQLServer
	if state.SQLServer.NativeMonthly != nil {
		native := *state.SQLServer.NativeMonthly
		native.Partitions = append([]SQLServerPartitionBoundary(nil), state.SQLServer.NativeMonthly.Partitions...)
		evidence.NativeMonthly = &native
	}
	evidence.ColumnDefaults = cloneStringPtrMap(state.SQLServer.ColumnDefaults)
	evidence.ColumnCollations = maps.Clone(state.SQLServer.ColumnCollations)
	cloned.SQLServer = &evidence
	if base == nil {
		if cloned.SQLServer.NativeMonthly != nil {
			return cloned, ErrIncompatibleExistingTable
		}
		return cloned, nil
	}
	projected := &TargetSQLServerPartitionEvidence{
		Method: base.Method, Column: base.Column, FunctionName: base.FunctionName, SchemeName: base.SchemeName,
		Filegroup: base.Filegroup, Partitions: append([]SQLServerPartitionBoundary(nil), base.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
	if !evidence.PartitionInspected {
		return cloned, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	if evidence.NativeMonthly != nil {
		projected.FunctionExists = evidence.NativeMonthly.FunctionExists
		projected.SchemeExists = evidence.NativeMonthly.SchemeExists
		projected.Exclusive = evidence.NativeMonthly.Exclusive
		projected.IndexesAligned = evidence.NativeMonthly.IndexesAligned
	}
	cloned.SQLServer.NativeMonthly = projected
	return cloned, nil
}

func cloneStringPtrMap(source map[string]*string) map[string]*string {
	if source == nil {
		return nil
	}
	cloned := make(map[string]*string, len(source))
	for key, value := range source {
		if value == nil {
			cloned[key] = nil
			continue
		}
		copyValue := *value
		cloned[key] = &copyValue
	}
	return cloned
}
