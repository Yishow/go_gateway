package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// studioV2WorkspaceDatabaseMetadataResponse binds actual table metadata to the
// saved workspace scope it was inspected for.
type studioV2WorkspaceDatabaseMetadataResponse struct {
	WorkspaceID       string                         `json:"workspace_id"`
	ConnectorID       string                         `json:"connector_id"`
	ConnectorRevision string                         `json:"connector_revision"`
	Database          string                         `json:"database"`
	Schema            string                         `json:"schema"`
	Table             string                         `json:"table"`
	InspectionStatus  dbtarget.TableInspectionStatus `json:"inspection_status"`
	Reason            string                         `json:"reason,omitempty"`
	Columns           []dbtarget.ColumnInfo          `json:"columns"`
}

// GetMetadata reports the actual metadata of the saved workspace or group table.
// Only the persisted connector, database, schema and table are inspected, and
// the caller must present the connector identity revision it is displaying.
// Missing, forbidden and failed outcomes stay distinct and never carry columns.
func (h *StudioV2WorkspaceDatabaseHandler) GetMetadata(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		renderStudioV2WorkspaceDatabaseError(c, dbtarget.ErrConnectorNotFound)
		return
	}
	expectedRevision := strings.TrimSpace(c.Query("expected_connector_revision"))
	if expectedRevision == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("expected_connector_revision is required"))
		return
	}

	connectorID := record.DatabaseConnectorID
	var group *workspace.WriteGroup
	if groupID := strings.TrimSpace(c.Query("group_id")); groupID != "" {
		if strings.TrimSpace(c.Query("expected_group_revision")) == "" {
			renderStudioV2WorkspaceValidationError(c, errors.New("expected_group_revision is required"))
			return
		}
		group, err = h.metadataGroup(c, record.ID, groupID, expectedRevision)
		if err != nil {
			renderWriteGroupError(c, err)
			return
		}
		connectorID = group.Destination.ConnectorID
	} else if c.Query("expected_group_revision") != "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("group_id is required"))
		return
	}
	connector, err := h.connectorSvc.ResolveSavedTarget(c.Request.Context(), connectorID, expectedRevision)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	// The metadata scope must name the target exactly like the schema preview
	// scope does, so the same saved connector resolves to the same database and
	// schema on both paths.
	dialect := canonicalRecordingDialect(string(connector.Kind))
	schemaName, tableName := recordingTargetSchema(connector, dialect), connectorConfigString(connector, "table")
	if group != nil {
		schemaName, tableName = group.Destination.TableSchema, group.Destination.TableName
	}
	inspection, err := h.inspectMetadataTable(c.Request.Context(), connector.ID, expectedRevision, schemaName, tableName, group)
	if err != nil {
		if errors.Is(err, dbtarget.ErrManagedDestinationInternal) {
			renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_INTERNAL_DATABASE", "recording requires a distinct application-owned destination", false, "select_recording_destination")
			return
		}
		if errors.Is(err, dbtarget.ErrManagedDestinationInvalid) {
			renderWriteGroupTyped(c, http.StatusUnprocessableEntity, "WRITE_GROUP_DESTINATION_INVALID", "the saved recording destination cannot be verified", false, "select_recording_destination")
			return
		}
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	if _, err := h.connectorSvc.ResolveSavedTarget(c.Request.Context(), connector.ID, expectedRevision); err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	if group != nil {
		if _, err := h.metadataGroup(c, record.ID, group.ID, expectedRevision); err != nil {
			renderWriteGroupError(c, err)
			return
		}
		current, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
		if err != nil {
			renderStudioV2WorkspaceBootstrapError(c)
			return
		}
		if current.ID != record.ID || current.DatabaseSetupRevision != record.DatabaseSetupRevision {
			renderWriteGroupError(c, workspace.ErrSetupRevisionConflict)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: studioV2WorkspaceDatabaseMetadataResponse{
		WorkspaceID:       record.ID,
		ConnectorID:       connector.ID,
		ConnectorRevision: connector.IdentityRevision,
		Database:          recordingTargetDatabase(connector, dialect),
		Schema:            inspection.Schema,
		Table:             inspection.Table,
		InspectionStatus:  inspection.Status,
		Reason:            inspection.Reason,
		Columns:           inspection.Columns,
	}})
}

func (h *StudioV2WorkspaceDatabaseHandler) inspectMetadataTable(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string, group *workspace.WriteGroup) (*dbtarget.TableInspection, error) {
	if group == nil || group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged {
		return dbtarget.NewReadOnlyTableInspector(h.connectorSvc).InspectTable(ctx, connectorID, schemaName, tableName)
	}
	groups, ok := h.writeGroups.(*workspace.WriteGroupService)
	if !ok || groups == nil {
		return nil, workspace.ErrWriteGroupServiceUnavailable
	}
	inspector := groups.ManagedTableInspector(h.connectorSvc)
	if inspector == nil {
		return nil, workspace.ErrWriteGroupServiceUnavailable
	}
	return inspector.InspectTableAtRevision(ctx, connectorID, expectedRevision, schemaName, tableName)
}

// metadataGroup resolves only persisted current-workspace scope. A stale or
// foreign request fails before any destination inspection.
func (h *StudioV2WorkspaceDatabaseHandler) metadataGroup(c *gin.Context, workspaceID, groupID, connectorRevision string) (*workspace.WriteGroup, error) {
	expected := strings.TrimSpace(c.Query("expected_group_revision"))
	if expected == "" {
		return nil, workspace.ErrWriteGroupValidation
	}
	if h.writeGroups == nil {
		return nil, workspace.ErrWriteGroupServiceUnavailable
	}
	result, err := h.writeGroups.List(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, workspace.ErrWriteGroupServiceUnavailable
	}
	if result.WorkspaceID != workspaceID {
		return nil, workspace.ErrWriteGroupNotFound
	}
	for _, group := range result.Groups {
		if group == nil || group.ID != groupID || group.WorkspaceID != workspaceID || group.Status == workspace.WriteGroupStatusDeleted {
			continue
		}
		if group.Revision != expected || group.Destination.ConnectorRevision != connectorRevision {
			return nil, workspace.ErrWriteGroupRevisionConflict
		}
		return group, nil
	}
	return nil, workspace.ErrWriteGroupNotFound
}
