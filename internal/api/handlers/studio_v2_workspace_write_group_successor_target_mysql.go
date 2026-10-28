package handlers

import (
	"context"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

func (h *checkedGroupSchemaTarget) executeGroupSchemaStatements(ctx context.Context, scope recordingplan.SchemaPreviewScope, statements []string) (recordingplan.SchemaExecution, error) {
	isMySQL := scope.GroupLayout != nil && scope.GroupLayout.MySQL != nil
	if err := h.check(ctx); err != nil {
		return recordingplan.SchemaExecution{RolledBack: !isMySQL}, err
	}
	if !isMySQL {
		execution, err := h.inspector.ExecuteSchemaStatements(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, statements)
		return recordingSchemaExecution(execution, err)
	}
	if len(statements) != 1 {
		return recordingplan.SchemaExecution{}, recordingplan.ErrInvalidGroupSchemaLayout
	}
	err := h.inspector.ExecuteMySQLManagedStatementAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
		scope.Database, scope.GroupLayout.MySQL.CapabilityRevision, statements[0])
	if err != nil {
		return recordingplan.SchemaExecution{}, mysqlGroupSchemaError(err)
	}
	return recordingplan.SchemaExecution{Committed: 1}, nil
}

func (h *checkedGroupSchemaTarget) mysqlGroupTargetInspector(scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		if err := h.check(ctx); err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		result, err := h.inspector.InspectMySQLManagedTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
			scope.Database, scope.GroupLayout.MySQL.CapabilityRevision, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, mysqlGroupSchemaError(err)
		}
		if result == nil || result.Inspection == nil || result.Evidence == nil {
			return recordingplan.TargetTableInspection{}, recordingplan.ErrTargetInspectionUnconfirmed
		}
		state := recordingplan.TargetTableInspection{
			Status: string(result.Inspection.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{},
			ColumnPrimaryKey: map[string]bool{}, IndexesInspected: result.Inspection.IndexesInspected,
		}
		for _, column := range result.Inspection.Columns {
			state.Columns = append(state.Columns, column.Name)
			state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.DataType, column.Nullable, column.PrimaryKey
		}
		for _, index := range result.Indexes {
			state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
		}
		e := result.Evidence
		state.MySQL = &recordingplan.TargetMySQLTableEvidence{
			ServerFamily: e.ServerFamily, ServerVersion: e.ServerVersion, CapabilityRevision: e.CapabilityRevision, Database: e.Database,
			Engine: e.Engine, Charset: e.Charset, Collation: e.Collation, TableType: e.TableType, PrimaryKey: e.PrimaryKey,
			ColumnCollations: e.ColumnCollations, ColumnDefaults: e.ColumnDefaults, JSONColumns: e.JSONColumns,
			Partitions: mirrorMySQLPartitionEvidence(e.Partitions),
		}
		if table != dbtarget.EffectReceiptTable && h.groups != nil {
			var err error
			state, err = h.groups.NormalizeMySQLPartitionInspection(ctx, scope, state)
			if err != nil {
				return state, mysqlGroupSchemaError(err)
			}
		}
		return state, h.check(ctx)
	}
}
