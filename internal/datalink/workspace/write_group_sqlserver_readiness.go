package workspace

import (
	"context"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

type managedSQLServerRuntimeInspector interface {
	InspectSQLServerManagedTableAtRevision(context.Context, string, string, string, string, string, string) (*dbtarget.SQLServerManagedInspection, error)
}

func (s *WriteGroupService) sqlServerReadiness(ctx context.Context, id string, content *WriteGroup, before *writeGroupReadinessSnapshot, result *WriteGroupReadiness) (*WriteGroupReadiness, error) {
	s.populateSQLServerPartitionMaintenanceReadiness(ctx, before.connectorKind, before.validated, result)
	frozen, inspectErr := s.readSQLServerRuntimeLayout(ctx, before)
	after, err := s.readinessSnapshot(ctx, id, content)
	if err != nil {
		return nil, err
	}
	if !sameWriteGroupReadinessSnapshot(before, after) || after.validateErr != nil {
		result = s.newSnapshotReadiness(ctx, after)
		result.Issues = append(result.Issues, writeGroupReadinessIssue("readiness-snapshot-changed", "configuration changed during schema verification", id))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	if inspectErr != nil || frozen == nil {
		result.Issues = append(result.Issues, writeGroupReadinessIssue("managed-schema-unverified", "managed destination ownership and canonical schema could not be verified", id))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	result.SchemaReady, result.SchemaDigest, result.runtimeLayout = true, frozen.SchemaDigest, frozen
	return finalizeWriteGroupReadiness(result), nil
}

func (s *WriteGroupService) readSQLServerRuntimeLayout(ctx context.Context, snapshot *writeGroupReadinessSnapshot) (*WriteGroupRuntimeLayout, error) {
	layout, err := s.managedSQLServerSchemaLayout(ctx, snapshot.validated)
	if err != nil {
		return nil, err
	}
	reader, ok := s.managedTableInspector.(managedSQLServerRuntimeInspector)
	if !ok {
		return nil, ErrWriteGroupReadinessUnavailable
	}
	group := snapshot.validated
	read := func(table string) (*dbtarget.SQLServerManagedInspection, error) {
		return reader.InspectSQLServerManagedTableAtRevision(ctx, group.Destination.ConnectorID, group.Destination.ConnectorRevision,
			group.Destination.Database, group.Destination.TableSchema, layout.SQLServer.CapabilityRevision, table)
	}
	physical, err := read(group.Destination.TableName)
	if err != nil || physical == nil || physical.Inspection == nil || physical.Evidence == nil || physical.Inspection.Status != dbtarget.TableInspectionExists {
		return nil, ErrWriteGroupValidation
	}
	receipt, err := read(dbtarget.EffectReceiptTable)
	if err != nil || receipt == nil || receipt.Inspection == nil || receipt.Evidence == nil || receipt.Inspection.Status != dbtarget.TableInspectionExists {
		return nil, ErrWriteGroupValidation
	}
	groupInspection := recordingPlanSQLServerInspection(physical)
	groupPhysical := dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: physical}).Group
	if group.StorageMode == WriteGroupStorageModeNativeMonthly && layout.SQLServer != nil && layout.SQLServer.NativeMonthly != nil {
		normalized, err := s.NormalizeSQLServerPartitionInspection(ctx, recordingplan.SchemaPreviewScope{
			WorkspaceID: group.WorkspaceID, ConnectorID: group.Destination.ConnectorID, Database: group.Destination.Database,
			Schema: group.Destination.TableSchema, GroupLayout: layout,
		}, groupInspection, recordingPlanSQLServerInspection(receipt))
		if err != nil {
			return nil, err
		}
		groupInspection = normalized
		groupPhysical.Evidence.NativeMonthly = sqlServerReadinessPartition(layout.SQLServer.NativeMonthly)
		groupPhysical.Evidence.PartitionInspected = true
		groupPhysical.Evidence.Plain = false
	}
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(*layout, group.Destination.Database, groupInspection, recordingPlanSQLServerInspection(receipt))
	if err != nil || digest == "" {
		return nil, ErrWriteGroupValidation
	}
	frozen := &WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: physical.Inspection.Columns, TagTypes: snapshot.tagTypes,
		ReceiptTableReady: true, SchemaDigest: digest,
		SQLServer: &dbtarget.SQLServerManagedRowProof{
			ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Group: groupPhysical, Receipt: receipt,
		},
	}
	if err := ValidateSQLServerRuntimeLayout(group, frozen); err != nil {
		return nil, err
	}
	return frozen, nil
}

func sqlServerReadinessPartition(evidence *recordingplan.SQLServerPartitionEvidence) *dbtarget.SQLServerManagedPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]dbtarget.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = dbtarget.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &dbtarget.SQLServerManagedPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: partitions, FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}
