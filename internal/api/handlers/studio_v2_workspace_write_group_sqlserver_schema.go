package handlers

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) sqlServerGroupTargetInspector(id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		if err := h.checkScope(ctx, id, req, scope); err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		physical, err := h.inspector.InspectSQLServerManagedTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
			scope.Database, scope.Schema, scope.GroupLayout.SQLServer.CapabilityRevision, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, sqlServerGroupSchemaError(err)
		}
		state, err := sqlServerTargetInspection(physical)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		if table != dbtarget.EffectReceiptTable && h.groups != nil &&
			scope.GroupLayout.SQLServer.NativeMonthly != nil && state.Status == string(dbtarget.TableInspectionExists) {
			receiptPhysical, err := h.inspector.InspectSQLServerManagedTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
				scope.Database, scope.Schema, scope.GroupLayout.SQLServer.CapabilityRevision, dbtarget.EffectReceiptTable)
			if err != nil {
				return recordingplan.TargetTableInspection{}, sqlServerGroupSchemaError(err)
			}
			receipt, err := sqlServerTargetInspection(receiptPhysical)
			if err != nil {
				return recordingplan.TargetTableInspection{}, err
			}
			state, err = h.groups.NormalizeSQLServerPartitionInspection(ctx, scope, state, receipt)
			if err != nil {
				return recordingplan.TargetTableInspection{}, sqlServerGroupSchemaError(err)
			}
		}
		return state, h.checkScope(ctx, id, req, scope)
	}
}

func sqlServerTargetInspection(physical *dbtarget.SQLServerManagedInspection) (recordingplan.TargetTableInspection, error) {
	if physical == nil || physical.Inspection == nil || physical.Evidence == nil {
		return recordingplan.TargetTableInspection{}, recordingplan.ErrTargetInspectionUnconfirmed
	}
	state := recordingplan.TargetTableInspection{
		Status: string(physical.Inspection.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
		ColumnPrimaryKey: map[string]bool{}, IndexesInspected: physical.Inspection.IndexesInspected,
	}
	for _, column := range physical.Inspection.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
	}
	for _, index := range physical.Indexes {
		state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
	}
	e := physical.Evidence
	state.SQLServer = &recordingplan.TargetSQLServerTableEvidence{
		ServerVersion: e.ServerVersion, ServerEdition: e.ServerEdition, CapabilityRevision: e.CapabilityRevision,
		Database: e.Database, Schema: e.Schema, EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel,
		MetadataVerified: e.MetadataVerified, TableType: e.TableType, ObjectID: e.ObjectID, Plain: e.Plain,
		PrimaryKey: e.PrimaryKey, ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults,
		ProvenanceJSONChecked: e.ProvenanceJSONChecked, RepairOperationID: e.RepairOperationID,
		NativeMonthly: sqlServerTargetPartition(e.NativeMonthly), PartitionInspected: e.PartitionInspected,
	}
	return state, nil
}

func (h *StudioV2WorkspaceWriteGroupSchemaHandler) executeSQLServerGroupSchemaStatements(ctx context.Context, id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope, _ *recordingplan.SchemaPreviewToken, statements []string) (recordingplan.SchemaExecution, error) {
	if err := h.checkScope(ctx, id, req, scope); err != nil {
		return recordingplan.SchemaExecution{RolledBack: true}, err
	}
	execution, err := h.inspector.ExecuteSQLServerManagedStatementsAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
		scope.Database, scope.Schema, scope.GroupLayout.SQLServer.CapabilityRevision, scope.GroupLayout.SQLServer.NativeMonthly != nil, statements)
	if execution == nil {
		return recordingplan.SchemaExecution{}, sqlServerGroupSchemaError(err)
	}
	return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, sqlServerGroupSchemaError(err)
}

func sqlServerTargetPartition(evidence *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.TargetSQLServerPartitionEvidence {
	if evidence == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(evidence.Partitions))
	for i, partition := range evidence.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.TargetSQLServerPartitionEvidence{
		Method: evidence.Method, Column: evidence.Column, FunctionName: evidence.FunctionName, SchemeName: evidence.SchemeName,
		Filegroup: evidence.Filegroup, Partitions: partitions, FunctionExists: evidence.FunctionExists, SchemeExists: evidence.SchemeExists,
		Exclusive: evidence.Exclusive, IndexesAligned: evidence.IndexesAligned,
	}
}

func sqlServerGroupSchemaError(err error) error {
	switch {
	case errors.Is(err, dbtarget.ErrSQLServerManagedScopeChanged):
		return fmt.Errorf("%w: %w", recordingplan.ErrPreviewTokenStale, err)
	case errors.Is(err, dbtarget.ErrSQLServerManagedTableForbidden), errors.Is(err, dbtarget.ErrSchemaExecutionPermissionDenied):
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetPermissionDenied, err)
	case errors.Is(err, dbtarget.ErrSQLServerManagedIndexUnusable), errors.Is(err, dbtarget.ErrSQLServerManagedPhysicalMismatch):
		return fmt.Errorf("%w: %w", recordingplan.ErrIncompatibleExistingTable, err)
	case errors.Is(err, dbtarget.ErrSQLServerManagedTableUnconfirmed), errors.Is(err, recordingplan.ErrSQLServerPartitionMaintenanceStateMismatch):
		return fmt.Errorf("%w: %w", recordingplan.ErrTargetInspectionUnconfirmed, err)
	default:
		return err
	}
}
