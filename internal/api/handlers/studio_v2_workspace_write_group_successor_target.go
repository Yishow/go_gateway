package handlers

import (
	"context"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// checkedGroupSchemaTarget reuses target adapters without granting schema authority.
// Each caller supplies its own mandatory private authority revalidation.
type checkedGroupSchemaTarget struct {
	groups     *workspace.WriteGroupService
	inspector  *dbtarget.ManagedTableInspector
	revalidate func(context.Context) error
}

func (h *checkedGroupSchemaTarget) check(ctx context.Context) error {
	if h == nil || h.inspector == nil || h.revalidate == nil {
		return recordingplan.ErrPreviewTokenStale
	}
	return h.revalidate(ctx)
}
func (h *StudioV2WorkspaceWriteGroupSchemaHandler) checkedTarget(id string, req WriteGroupSchemaRequest, scope recordingplan.SchemaPreviewScope) *checkedGroupSchemaTarget {
	return &checkedGroupSchemaTarget{groups: h.groups, inspector: h.inspector, revalidate: func(ctx context.Context) error { return h.checkScope(ctx, id, req, scope) }}
}
func (h *checkedGroupSchemaTarget) targetInspector(scope recordingplan.SchemaPreviewScope) recordingplan.TargetInspector {
	if scope.GroupLayout != nil && scope.GroupLayout.MySQL != nil {
		return h.mysqlGroupTargetInspector(scope)
	}
	if scope.GroupLayout != nil && scope.GroupLayout.SQLServer != nil {
		return h.sqlServerGroupTargetInspector(scope)
	}
	return func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		if err := h.check(ctx); err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		result, err := h.inspector.InspectTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}
		if result == nil {
			return recordingplan.TargetTableInspection{}, errRecordingTargetInspectionEmpty
		}
		state := recordingplan.TargetTableInspection{Status: string(result.Status), ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}}
		for _, column := range result.Columns {
			state.Columns = append(state.Columns, column.Name)
			state.ColumnTypes[column.Name] = column.DataType
			state.ColumnNullable[column.Name] = column.Nullable
			state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
		}
		// A layout that declares indexes (the long table) needs real catalog
		// evidence of them; wide layouts declare none and never pay for the query.
		if scope.GroupLayout != nil && len(scope.GroupLayout.Indexes) > 0 && table == scope.GroupLayout.TableName &&
			result.Status == dbtarget.TableInspectionExists {
			indexes, err := h.inspector.InspectTableIndexesAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
			if err != nil {
				return recordingplan.TargetTableInspection{}, err
			}
			state.IndexesInspected = true
			for _, index := range indexes {
				state.Indexes = append(state.Indexes, recordingplan.TargetIndexInspection{Name: index.Name, Columns: index.Columns, Unique: index.Unique})
			}
		}
		if scope.GroupLayout != nil && scope.GroupLayout.Partition != nil && table == scope.GroupLayout.TableName {
			evidence, err := h.timescaleEvidence(ctx, scope, result.Status == dbtarget.TableInspectionExists)
			if err != nil {
				return recordingplan.TargetTableInspection{}, err
			}
			state.Timescale = evidence
		}
		return state, h.check(ctx)
	}
}

func (h *checkedGroupSchemaTarget) timescaleEvidence(ctx context.Context, scope recordingplan.SchemaPreviewScope, exists bool) (*recordingplan.TargetTimescaleEvidence, error) {
	if err := h.check(ctx); err != nil {
		return nil, err
	}
	capability, err := h.inspector.ProbeTimescaleAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema)
	if err != nil {
		return nil, err
	}
	evidence := &recordingplan.TargetTimescaleEvidence{
		Available: capability.Available, Reasons: capability.Reasons, ExtensionSchema: capability.ExtensionSchema, Version: capability.Version,
	}
	if !exists {
		return evidence, nil
	}
	hypertable, err := h.inspector.InspectHypertableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, scope.GroupLayout.TableName)
	if err != nil {
		return nil, err
	}
	if hypertable != nil {
		evidence.Hypertable = &recordingplan.TargetHypertable{
			TimeColumn: hypertable.TimeColumn, ChunkInterval: hypertable.ChunkInterval, Dimensions: hypertable.Dimensions,
			CompressionEnabled: hypertable.CompressionEnabled, Jobs: hypertable.Jobs,
		}
	}
	return evidence, nil
}
