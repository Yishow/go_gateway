package handlers

import (
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type StudioV2WorkspaceDatabaseHandler struct {
	workspaceSvc *workspace.Service
	deviceSvc    *device.Service
	ruleSvc      *sourcerule.Service
	connectorSvc *dbtarget.ConnectorService
	mappingSvc   *dbtarget.MappingService
}

func NewStudioV2WorkspaceDatabaseHandler(workspaceSvc *workspace.Service, deviceSvc *device.Service, ruleSvc *sourcerule.Service, connectorSvc *dbtarget.ConnectorService, mappingSvc *dbtarget.MappingService) *StudioV2WorkspaceDatabaseHandler {
	return &StudioV2WorkspaceDatabaseHandler{
		workspaceSvc: workspaceSvc,
		deviceSvc:    deviceSvc,
		ruleSvc:      ruleSvc,
		connectorSvc: connectorSvc,
		mappingSvc:   mappingSvc,
	}
}

func (h *StudioV2WorkspaceDatabaseHandler) GetConfig(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
		return
	}

	connector, err := h.connectorSvc.GetByID(c.Request.Context(), record.DatabaseConnectorID)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": buildWorkspaceDatabaseConfigResponse(record.ID, connector)})
}

func (h *StudioV2WorkspaceDatabaseHandler) UpdateConfig(c *gin.Context) {
	record, req, ok := h.parseWorkspaceDatabaseConfigRequest(c)
	if !ok {
		return
	}

	savedConnector, err := h.saveWorkspaceConnector(c, record, req)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	payload := buildWorkspaceDatabaseConfigResponse(record.ID, savedConnector)
	payload.RuntimeApplyStatus, payload.RuntimeApplyMessage = resolveStudioV2WorkspaceRuntimeApplyStatus(c.Request.Context(), h.deviceSvc, record.OrderedDeviceIDs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func (h *StudioV2WorkspaceDatabaseHandler) ListTargets(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []studioV2WorkspaceDatabaseTargetResponse{}})
		return
	}

	bindings, err := h.listWorkspacePointBindings(c)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	tagToPoint := make(map[string]string, len(bindings))
	for _, binding := range bindings {
		tagToPoint[binding.TagID] = binding.PointID
	}
	rows, err := h.mappingSvc.List(c.Request.Context(), dbtarget.TargetMappingListFilter{ConnectorID: &record.DatabaseConnectorID})
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	payload := make([]studioV2WorkspaceDatabaseTargetResponse, 0, len(rows))
	for _, row := range rows {
		pointID := tagToPoint[row.TagID]
		if pointID == "" {
			continue
		}
		payload = append(payload, studioV2WorkspaceDatabaseTargetResponse{
			ID:          row.ID,
			WorkspaceID: record.ID,
			PointID:     pointID,
			TagID:       row.TagID,
			ColumnName:  row.ColumnName,
			Enabled:     row.Enabled,
			SaveState:   "saved",
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func (h *StudioV2WorkspaceDatabaseHandler) UpsertTarget(c *gin.Context) {
	record, connector, binding, req, ok := h.resolveWorkspaceDatabaseTarget(c)
	if !ok {
		return
	}

	filter := dbtarget.TargetMappingListFilter{
		ConnectorID: &connector.ID,
		TagID:       &binding.TagID,
	}
	rows, err := h.mappingSvc.List(c.Request.Context(), filter)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	var savedRow *schema.DatabaseTargetMapping
	if len(rows) == 0 {
		savedRow, err = h.mappingSvc.Create(c.Request.Context(), dbtarget.CreateTargetMappingRequest{
			TagID:                binding.TagID,
			ConnectorID:          connector.ID,
			TableSchema:          connectorConfigString(connector, "schema"),
			TableName:            connectorConfigString(connector, "table"),
			ColumnName:           strings.TrimSpace(req.ColumnName),
			WriteMode:            connectorConfigWriteMode(connector),
			TimestampColumn:      workspaceOptionalString(connectorConfigString(connector, "timestamp_column")),
			WriteIntervalSeconds: workspaceOptionalInt(connector.DefaultWriteIntervalSeconds),
			Enabled:              boolPtr(req.Enabled),
		})
	} else {
		savedRow, err = h.mappingSvc.Update(c.Request.Context(), rows[0].ID, dbtarget.UpdateTargetMappingRequest{
			TableSchema:          workspaceOptionalString(connectorConfigString(connector, "schema")),
			TableName:            workspaceOptionalString(connectorConfigString(connector, "table")),
			ColumnName:           workspaceOptionalString(strings.TrimSpace(req.ColumnName)),
			WriteMode:            workspaceOptionalWriteMode(connectorConfigWriteMode(connector)),
			TimestampColumn:      workspaceOptionalString(connectorConfigString(connector, "timestamp_column")),
			WriteIntervalSeconds: workspaceOptionalInt(connector.DefaultWriteIntervalSeconds),
			Enabled:              boolPtr(req.Enabled),
		})
	}
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	payload := studioV2WorkspaceDatabaseTargetResponse{
		ID:          savedRow.ID,
		WorkspaceID: record.ID,
		PointID:     binding.PointID,
		TagID:       savedRow.TagID,
		ColumnName:  savedRow.ColumnName,
		Enabled:     savedRow.Enabled,
		SaveState:   "saved",
		CreatedAt:   savedRow.CreatedAt,
		UpdatedAt:   savedRow.UpdatedAt,
	}
	payload.RuntimeApplyStatus, payload.RuntimeApplyMessage = resolveStudioV2WorkspaceRuntimeApplyStatus(c.Request.Context(), h.deviceSvc, record.OrderedDeviceIDs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}
