package recordingplan

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"
)

func advanceSQLServerPartitionStateAfterReconstruction(
	ctx context.Context,
	q DBTX,
	scope StorageScope,
	lineage *StorageLineage,
	fromGeneration int,
	operationID string,
	now time.Time,
) error {
	if lineage == nil || lineage.Layout == nil || lineage.Layout.SQLServer == nil || lineage.Layout.SQLServer.NativeMonthly == nil {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	state, err := ReadSQLServerPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return fmt.Errorf("read sqlserver partition state after storage repair: %w", err)
	}
	if state == nil {
		return nil
	}
	request := SQLServerPartitionMaintenanceRequest{Scope: scope, ConnectorRevision: lineage.ConnectorRevision}
	if lineage.Generation != fromGeneration || state.Generation != fromGeneration ||
		!sqlServerStateMatchesLineage(state, lineage, request, lineage.LayoutDigest) {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	next := *state
	next.Generation = fromGeneration + 1
	next.Current = CloneSQLServerPartitionEvidence(state.Current)
	next.Revision++
	next.LastOperationID = operationID
	next.UpdatedAt = now.UTC()
	if err := validateSQLServerPartitionMaintenanceState(&next); err != nil {
		return err
	}
	encoded, err := json.Marshal(next.Current)
	if err != nil {
		return fmt.Errorf("encode authorized sqlserver partition state: %w", err)
	}
	result, err := q.ExecContext(ctx, numberedPlaceholders(`UPDATE `+sqlServerPartitionMaintenanceStateTable+` SET
		generation = $7, current_evidence = $8, revision = $9, last_operation_id = $10, updated_at = $11
		WHERE `+storageScopeWhere+` AND base_layout_digest = $12 AND connector_revision = $13
		AND generation = $14 AND revision = $15`), append(scopeArgs(scope), next.Generation, string(encoded), next.Revision,
		next.LastOperationID, next.UpdatedAt.Format(time.RFC3339Nano), next.BaseLayoutDigest, next.ConnectorRevision,
		fromGeneration, state.Revision)...)
	if err != nil {
		return fmt.Errorf("advance sqlserver partition state after storage repair: %w", err)
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read sqlserver partition state reconstruction: %w", err)
	}
	if changed != 1 {
		return ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return nil
}

func (s *Service) sqlServerNativeRepairStatements(
	ctx context.Context,
	req StorageRestoreRequest,
	layout GroupSchemaLayout,
	op *SchemaOperation,
	statements []string,
) ([]string, error) {
	if layout.SQLServer == nil || layout.SQLServer.NativeMonthly == nil {
		return slices.Clone(statements), nil
	}
	repo, ok := s.repo.(storageRepairRepository)
	if !ok {
		return nil, ErrStorageRepairUnsupported
	}
	lineage, err := ReadStorageLineage(ctx, repo.storageDB(), req.Scope)
	if err != nil {
		return nil, err
	}
	if refusal := sqlServerStorageLineageRefusal(lineage, req); refusal != nil ||
		!groupSchemaLayoutsEqual(layout, *lineage.Layout) {
		return nil, ErrSchemaOperationMismatch
	}
	var detail storageRepairDetail
	if op == nil || op.Action != StorageRepairAction || json.Unmarshal([]byte(op.Detail), &detail) != nil ||
		detail.Scope != req.Scope || detail.Generation != lineage.Generation || detail.SQLServer == nil {
		return nil, ErrSchemaOperationMismatch
	}
	current, err := authorizedSQLServerNativeCurrent(ctx, repo.storageDB(), req.Scope, lineage)
	if err != nil {
		return nil, err
	}
	partition, err := sqlServerRestorePartitionDDL(layout, current)
	if err != nil {
		return nil, err
	}
	remainder := make([]string, 0, len(statements))
	for _, statement := range statements {
		if strings.Contains(statement, "CREATE PARTITION FUNCTION ") || strings.Contains(statement, "CREATE PARTITION SCHEME ") ||
			strings.Contains(statement, sqlserverPartitionOwnerProperty) {
			continue
		}
		remainder = append(remainder, statement)
	}
	return append(partition, remainder...), nil
}

func sqlServerRestorePartitionDDL(layout GroupSchemaLayout, current *SQLServerPartitionEvidence) ([]string, error) {
	if layout.SQLServer == nil || layout.SQLServer.NativeMonthly == nil || current == nil {
		return nil, ErrInvalidGroupSchemaLayout
	}
	base := layout.SQLServer.NativeMonthly
	if err := validateSQLServerPartitionEvidence(base, layout.TableName); err != nil {
		return nil, err
	}
	if err := validateCurrentSQLServerPartitionEvidence(current); err != nil || current.FunctionName != base.FunctionName ||
		current.SchemeName != base.SchemeName || current.Method != base.Method || current.Column != base.Column || current.Filegroup != base.Filegroup {
		return nil, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	values := make([]string, 0, len(current.Partitions)-1)
	for _, partition := range current.Partitions[:len(current.Partitions)-1] {
		boundary, err := parseSQLServerPartitionBoundary(partition.LessThan)
		if err != nil {
			return nil, fmt.Errorf("%w: sqlserver repair partition boundary: %w", ErrInvalidGroupSchemaLayout, err)
		}
		values = append(values, "N'"+boundary.Format(sqlServerPartitionTimeLayout)+"'")
	}
	function := quoteSQLServerIdentifier(current.FunctionName)
	scheme := quoteSQLServerIdentifier(current.SchemeName)
	owner := quoteSQLServerLiteral(layout.SQLServer.Schema + "." + layout.TableName)
	return []string{
		"CREATE PARTITION FUNCTION " + function + " (" + sqlserverPartitionFunctionType + ") AS " +
			sqlserverPartitionMethod + " FOR VALUES (" + strings.Join(values, ", ") + ");",
		"EXEC sys.sp_addextendedproperty @name=N'" + sqlserverPartitionOwnerProperty + "', @value=N'" + owner +
			"', @level0type=N'PARTITION FUNCTION', @level0name=N'" + current.FunctionName + "';",
		"CREATE PARTITION SCHEME " + scheme + " AS PARTITION " + function + " ALL TO ([" + current.Filegroup + "]);",
		"EXEC sys.sp_addextendedproperty @name=N'" + sqlserverPartitionOwnerProperty + "', @value=N'" + owner +
			"', @level0type=N'PARTITION SCHEME', @level0name=N'" + current.SchemeName + "';",
	}, nil
}

func authorizedSQLServerNativeCurrent(ctx context.Context, q DBTX, scope StorageScope, lineage *StorageLineage) (*SQLServerPartitionEvidence, error) {
	if lineage == nil || !lineage.Confirmed() || lineage.Layout == nil || lineage.Layout.SQLServer == nil ||
		lineage.Layout.SQLServer.NativeMonthly == nil || lineage.StorageScope != scope || lineage.RepairState != "" {
		return nil, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	state, err := ReadSQLServerPartitionMaintenanceState(ctx, q, scope)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return CloneSQLServerPartitionEvidence(lineage.Layout.SQLServer.NativeMonthly), nil
	}
	request := SQLServerPartitionMaintenanceRequest{Scope: scope, ConnectorRevision: lineage.ConnectorRevision}
	if !sqlServerStateMatchesLineage(state, lineage, request, lineage.LayoutDigest) {
		return nil, ErrSQLServerPartitionMaintenanceStateMismatch
	}
	return CloneSQLServerPartitionEvidence(state.Current), nil
}

func normalizeSQLServerNativeRestoreGroup(layout GroupSchemaLayout, group TargetTableInspection) (TargetTableInspection, error) {
	if layout.SQLServer == nil || layout.SQLServer.NativeMonthly == nil {
		return group, nil
	}
	if group.SQLServer == nil || !group.SQLServer.PartitionInspected {
		return TargetTableInspection{}, ErrTargetInspectionUnconfirmed
	}
	actual := group.SQLServer.NativeMonthly
	if actual == nil {
		return group, nil
	}
	if err := validateSQLServerNativeRestoreEvidence(layout, group, layout.SQLServer.NativeMonthly, false); err != nil {
		return TargetTableInspection{}, err
	}
	normalized := group
	evidence := *group.SQLServer
	evidence.NativeMonthly = sqlServerRestoreTargetPartition(layout.SQLServer.NativeMonthly, true, true, true, group.Status == inspectionExists)
	normalized.SQLServer = &evidence
	return normalized, nil
}

func validateSQLServerNativeRestoreEvidence(
	layout GroupSchemaLayout,
	group TargetTableInspection,
	expected *SQLServerPartitionEvidence,
	exact bool,
) error {
	if layout.SQLServer == nil || expected == nil || group.SQLServer == nil || !group.SQLServer.PartitionInspected {
		return ErrTargetInspectionUnconfirmed
	}
	actual := group.SQLServer.NativeMonthly
	functionName, schemeName, err := SQLServerNativeMonthlyObjectNames(layout.TableName)
	if err != nil || actual == nil || !actual.FunctionExists || !actual.Exclusive ||
		actual.FunctionName != functionName || actual.SchemeName != schemeName ||
		(expected.FunctionName != functionName || expected.SchemeName != schemeName) {
		return sqlServerPartitionUnavailable("native partition objects are not exclusively owned by the confirmed table")
	}
	if actual.SchemeExists {
		if actual.Filegroup != expected.Filegroup {
			return sqlServerPartitionUnavailable("native partition scheme does not map to the authorized filegroup")
		}
	} else if group.Status != inspectionMissing {
		return sqlServerPartitionUnavailable("an existing native table is missing its partition scheme")
	}
	if group.Status == inspectionExists && !actual.IndexesAligned {
		return sqlServerPartitionUnavailable("existing native table indexes are not aligned")
	}
	current := &SQLServerPartitionEvidence{
		Method: actual.Method, Column: actual.Column, FunctionName: actual.FunctionName, SchemeName: actual.SchemeName,
		Filegroup: expected.Filegroup, Partitions: slices.Clone(actual.Partitions),
	}
	if err := validateCurrentSQLServerPartitionEvidence(current); err != nil {
		return sqlServerPartitionUnavailable("native partition boundaries are not canonical")
	}
	if exact {
		if !sameSQLServerPartitionEvidence(current, expected) {
			return sqlServerPartitionUnavailable("native partition state differs from the authorized current projection")
		}
		return nil
	}
	if !sameSQLServerPartitionEvidence(current, expected) && !sameSQLServerPartitionPrefix(expected, current) {
		return sqlServerPartitionUnavailable("native partition boundaries do not extend the frozen layout")
	}
	return nil
}

func normalizeSQLServerNativeRepairGroup(
	ctx context.Context,
	q DBTX,
	req StorageRestoreRequest,
	lineage *StorageLineage,
	layout GroupSchemaLayout,
	group TargetTableInspection,
) (TargetTableInspection, error) {
	if layout.SQLServer == nil || layout.SQLServer.NativeMonthly == nil {
		return group, nil
	}
	current, err := authorizedSQLServerNativeCurrent(ctx, q, req.Scope, lineage)
	if err != nil {
		return TargetTableInspection{}, err
	}
	if group.SQLServer != nil && group.SQLServer.NativeMonthly != nil {
		if err := validateSQLServerNativeRestoreEvidence(layout, group, current, true); err != nil {
			return TargetTableInspection{}, err
		}
	}
	return normalizeSQLServerNativeRestoreGroup(layout, group)
}

func inspectSQLServerNativeRestoreGroupTables(
	ctx context.Context,
	layout GroupSchemaLayout,
	database string,
	inspect TargetInspector,
) (groupInspection, error) {
	var found groupInspection
	database = strings.TrimSpace(database)
	if inspect == nil || database == "" || layout.SQLServer == nil || layout.SQLServer.NativeMonthly == nil {
		return found, ErrTargetInspectionUnconfirmed
	}
	if err := validateSQLServerPreviewScope(database, layout.SQLServer.Schema, dialectSQLServer, *layout.SQLServer); err != nil {
		return found, err
	}
	if _, err := GenerateGroupSchemaDDL(dialectSQLServer, layout); err != nil {
		return found, err
	}
	group, err := inspect(ctx, layout.TableName)
	if err != nil {
		return found, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, layout.TableName, err)
	}
	receipt, err := inspect(ctx, managedEffectReceiptTable)
	if err != nil {
		return found, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, managedEffectReceiptTable, err)
	}
	normalized, err := normalizeSQLServerNativeRestoreGroup(layout, group)
	if err != nil {
		return found, err
	}
	return inspectSQLServerGroupSnapshot(layout, database, normalized, receipt)
}

func sqlServerRestoreTargetPartition(
	native *SQLServerPartitionEvidence,
	functionExists, schemeExists, exclusive, indexesAligned bool,
) *TargetSQLServerPartitionEvidence {
	partitions := slices.Clone(native.Partitions)
	return &TargetSQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: partitions, FunctionExists: functionExists, SchemeExists: schemeExists,
		Exclusive: exclusive, IndexesAligned: indexesAligned,
	}
}
