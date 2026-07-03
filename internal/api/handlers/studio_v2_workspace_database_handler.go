package handlers

import (
	"context"
	"errors"
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
	auditSvc     workspaceAuditRecorder
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

	c.JSON(http.StatusOK, gin.H{"success": true, "data": buildWorkspaceDatabaseConfigResponse(record.ID, connector, record.DatabaseRowGroups)})
}

func (h *StudioV2WorkspaceDatabaseHandler) UpdateConfig(c *gin.Context) {
	record, req, ok := h.parseWorkspaceDatabaseConfigRequest(c)
	if !ok {
		return
	}
	if req.RowGroups != nil {
		if err := workspace.ValidateDatabaseRowGroups(record.DatabaseConnectorID, req.Schema, req.Table, req.RowGroups); err != nil {
			renderStudioV2WorkspaceDatabaseError(c, err)
			return
		}
	}

	savedConnector, err := h.saveWorkspaceConnector(c, record, req)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	if req.RowGroups != nil {
		record, err = h.workspaceSvc.SaveDatabaseRowGroups(
			c.Request.Context(),
			savedConnector.ID,
			connectorConfigString(savedConnector, "schema"),
			connectorConfigString(savedConnector, "table"),
			req.RowGroups,
		)
		if err != nil {
			renderStudioV2WorkspaceDatabaseError(c, err)
			return
		}
	} else {
		record, err = h.workspaceSvc.GetOrCreate(c.Request.Context())
		if err != nil {
			renderStudioV2WorkspaceBootstrapError(c)
			return
		}
	}
	payload := buildWorkspaceDatabaseConfigResponse(record.ID, savedConnector, record.DatabaseRowGroups)
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, record.OrderedDeviceIDs, nil)
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	h.recordDatabaseConfigAudit(c.Request.Context(), record.ID, savedConnector, payload.RuntimeApplyStatus)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func (h *StudioV2WorkspaceDatabaseHandler) GenerateSchema(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	if connectorID == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("workspace database config not found"))
		return
	}

	var req dbtarget.SchemaGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	result, err := h.connectorSvc.GenerateSchema(c.Request.Context(), connectorID, req)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// EnsureWorkspaceSchema 確保 workspace 已綁定的 database connector 的目標資料表存在。
// 若 workspace 尚未設定 database（無 connector），視為無需建表並回 nil。
// 用於 activation 前的保險步驟，實作 workspaceSchemaEnsurer 介面。
func (h *StudioV2WorkspaceDatabaseHandler) EnsureWorkspaceSchema(ctx context.Context) error {
	record, err := h.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return err
	}
	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	if connectorID == "" {
		return nil
	}
	if _, err := h.connectorSvc.GenerateSchema(ctx, connectorID, dbtarget.SchemaGenerateRequest{DryRun: false}); err != nil {
		return err
	}
	return nil
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
	rowGroupByPoint := workspaceDatabaseTargetRefsByPoint(record.DatabaseTargetRefs, record.DatabaseRowGroups)

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
			RowGroupID:  rowGroupByPoint[pointID],
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
	rowGroupID := strings.TrimSpace(req.RowGroupID)
	if rowGroupID != "" {
		if !workspaceDatabaseRowGroupExists(record.DatabaseRowGroups, rowGroupID) {
			renderStudioV2WorkspaceValidationError(c, errors.New("database target row group does not exist"))
			return
		}
		if !workspaceDatabaseRowGroupContainsPoint(record.DatabaseRowGroups, rowGroupID, binding.PointID) {
			renderStudioV2WorkspaceValidationError(c, errors.New("database target row group does not contain point"))
			return
		}
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
			GroupKey:             workspaceOptionalString(rowGroupTargetKey(rowGroupID, binding.PointID)),
			WriteIntervalSeconds: workspaceOptionalInt(connector.DefaultWriteIntervalSeconds),
			Enabled:              boolPtr(req.Enabled),
			// 允許在目標表/欄位尚未建立時先儲存，稍後由建表流程補建。
			AllowMissingTable: true,
		})
	} else {
		savedRow, err = h.mappingSvc.Update(c.Request.Context(), rows[0].ID, dbtarget.UpdateTargetMappingRequest{
			TableSchema:          workspaceOptionalString(connectorConfigString(connector, "schema")),
			TableName:            workspaceOptionalString(connectorConfigString(connector, "table")),
			ColumnName:           workspaceOptionalString(strings.TrimSpace(req.ColumnName)),
			WriteMode:            workspaceOptionalWriteMode(connectorConfigWriteMode(connector)),
			TimestampColumn:      workspaceOptionalString(connectorConfigString(connector, "timestamp_column")),
			GroupKey:             workspaceOptionalString(rowGroupTargetKey(rowGroupID, binding.PointID)),
			WriteIntervalSeconds: workspaceOptionalInt(connector.DefaultWriteIntervalSeconds),
			Enabled:              boolPtr(req.Enabled),
			AllowMissingTable:    true,
		})
	}
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	if _, err := h.workspaceSvc.SaveDatabaseTargetReference(c.Request.Context(), binding.PointID, rowGroupID); err != nil {
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
		RowGroupID:  rowGroupID,
		SaveState:   "saved",
		CreatedAt:   savedRow.CreatedAt,
		UpdatedAt:   savedRow.UpdatedAt,
	}
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, record.OrderedDeviceIDs, []string{binding.PointID, connector.ID})
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	h.recordDatabaseTargetAudit(c.Request.Context(), record.ID, savedRow, binding.PointID, payload.RuntimeApplyStatus)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}
