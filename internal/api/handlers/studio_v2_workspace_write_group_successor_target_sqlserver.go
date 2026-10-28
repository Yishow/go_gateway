package handlers

import (
	"context"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

func (h *checkedGroupSchemaTarget) sqlServerGroupTargetInspector(scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		if err := h.check(ctx); err != nil {
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
		return state, h.check(ctx)
	}
}

func (h *checkedGroupSchemaTarget) executeSQLServerGroupSchemaStatements(ctx context.Context, scope recordingplan.SchemaPreviewScope, _ *recordingplan.SchemaPreviewToken, statements []string) (recordingplan.SchemaExecution, error) {
	if err := h.check(ctx); err != nil {
		return recordingplan.SchemaExecution{RolledBack: true}, err
	}
	execution, err := h.inspector.ExecuteSQLServerManagedStatementsAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision,
		scope.Database, scope.Schema, scope.GroupLayout.SQLServer.CapabilityRevision, scope.GroupLayout.SQLServer.NativeMonthly != nil, statements)
	if execution == nil {
		return recordingplan.SchemaExecution{}, sqlServerGroupSchemaError(err)
	}
	return recordingplan.SchemaExecution{Committed: execution.Committed, RolledBack: execution.RolledBack}, sqlServerGroupSchemaError(err)
}
