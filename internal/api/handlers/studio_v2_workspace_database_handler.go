package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
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
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: nil})
		return
	}

	connector, err := h.connectorSvc.GetByID(c.Request.Context(), record.DatabaseConnectorID)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	payload := buildWorkspaceDatabaseConfigResponse(record.ID, connector, record.DatabaseRowGroups)
	payload.SetupRevision = record.DatabaseSetupRevision
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}

// UpdateConfig saves the connector, its workspace binding and row groups in
// one local transaction guarded by the expected setup revision.
func (h *StudioV2WorkspaceDatabaseHandler) UpdateConfig(c *gin.Context) {
	record, req, ok := h.parseWorkspaceDatabaseConfigRequest(c)
	if !ok {
		return
	}
	if req.RowGroups != nil {
		connectorID := record.DatabaseConnectorID
		if req.ConnectorID != "" {
			connectorID = req.ConnectorID
		}
		if err := workspace.ValidateDatabaseRowGroups(connectorID, req.Schema, req.Table, req.RowGroups); err != nil {
			renderStudioV2WorkspaceDatabaseError(c, err)
			return
		}
	}

	prepared, err := h.prepareWorkspaceConnector(c.Request.Context(), record, req)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	savedConnector := prepared.Connector
	record, err = h.commitWorkspaceConnectorSetup(c.Request.Context(), record.DatabaseConnectorID, req, prepared)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}
	payload := buildWorkspaceDatabaseConfigResponse(record.ID, savedConnector, record.DatabaseRowGroups)
	payload.SetupRevision = record.DatabaseSetupRevision
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, record.OrderedDeviceIDs, nil)
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	h.recordDatabaseConfigAudit(c.Request.Context(), record.ID, savedConnector, payload.RuntimeApplyStatus)
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
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
	if !requireSchemaConfirmation(c, req.DryRun) {
		return
	}

	result, err := h.connectorSvc.GenerateSchema(c.Request.Context(), connectorID, req)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

// EnsureWorkspaceSchema 只核對 workspace 已綁定 connector 的目標資料表是否齊備，
// 不建立也不修改任何資料表；缺表時回報需要先準備，由操作者以確認流程建立。
// 判定以 dry-run 語句數加上實際 table inspection 佐證，兩者任一指出缺表即要求準備，
// inspection 無法確認時回安全錯誤而不是誤判就緒。
// 若 workspace 尚未設定 database（無 connector），視為無需資料庫並回 nil。
func (h *StudioV2WorkspaceDatabaseHandler) EnsureWorkspaceSchema(ctx context.Context) error {
	record, err := h.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return err
	}
	connectorID := strings.TrimSpace(record.DatabaseConnectorID)
	if connectorID == "" {
		return nil
	}
	plan, err := h.connectorSvc.GenerateSchema(ctx, connectorID, dbtarget.SchemaGenerateRequest{DryRun: true})
	if err != nil {
		return err
	}
	statementCount := 0
	if plan != nil {
		statementCount = len(plan.Statements)
	}
	inspections, err := h.inspectWorkspaceTargetTables(ctx, connectorID)
	if err != nil {
		return err
	}
	return evaluateWorkspaceSchemaReadiness(statementCount, inspections)
}

// evaluateWorkspaceSchemaReadiness decides activation readiness from the
// dry-run statement count and corroborating table inspections. A missing
// table or any remaining dry-run statement requires preparation; an
// unconfirmable inspection is a safe error instead of a false readiness.
func evaluateWorkspaceSchemaReadiness(statementCount int, inspections []*dbtarget.TableInspection) error {
	for _, inspection := range inspections {
		switch inspection.Status {
		case dbtarget.TableInspectionMissing:
			return errWorkspaceSchemaPreparationRequired
		case dbtarget.TableInspectionForbidden, dbtarget.TableInspectionFailed:
			return fmt.Errorf("workspace schema table %s could not be confirmed: %s", inspection.Table, inspection.Reason)
		}
	}
	if statementCount > 0 {
		return errWorkspaceSchemaPreparationRequired
	}
	return nil
}

// inspectWorkspaceTargetTables inspects every distinct table referenced by the
// connector's enabled target mappings — the same live projection
// GenerateSchema acts on — so readiness is corroborated by the actual database
// instead of the dry-run statement list alone. Disabled mappings are excluded:
// their tables are never written, so they cannot block activation.
func (h *StudioV2WorkspaceDatabaseHandler) inspectWorkspaceTargetTables(ctx context.Context, connectorID string) ([]*dbtarget.TableInspection, error) {
	enabled := true
	mappings, err := h.mappingSvc.List(ctx, dbtarget.TargetMappingListFilter{ConnectorID: &connectorID, Enabled: &enabled})
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	inspections := make([]*dbtarget.TableInspection, 0, len(mappings))
	for _, mapping := range mappings {
		table := strings.TrimSpace(mapping.TableName)
		if table == "" {
			continue
		}
		schemaName := strings.TrimSpace(mapping.TableSchema)
		key := schemaName + "." + table
		if seen[key] {
			continue
		}
		seen[key] = true
		inspection, err := h.connectorSvc.InspectTable(ctx, connectorID, schemaName, table)
		if err != nil {
			return nil, err
		}
		inspections = append(inspections, inspection)
	}
	return inspections, nil
}

func (h *StudioV2WorkspaceDatabaseHandler) ListTargets(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	if strings.TrimSpace(record.DatabaseConnectorID) == "" {
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: []studioV2WorkspaceDatabaseTargetResponse{}})
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
			SaveState:   workspaceSaveStateSaved,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
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

	var prepared *dbtarget.PreparedTargetMapping
	if len(rows) == 0 {
		prepared, err = h.mappingSvc.PrepareCreate(c.Request.Context(), dbtarget.CreateTargetMappingRequest{
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
		prepared, err = h.mappingSvc.PrepareUpdate(c.Request.Context(), rows[0].ID, dbtarget.UpdateTargetMappingRequest{
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
	record, err = h.commitWorkspaceTargetSetup(c.Request.Context(), connector.ID, binding.PointID, rowGroupID, req.ExpectedSetupRevision, prepared)
	if err != nil {
		renderStudioV2WorkspaceDatabaseError(c, err)
		return
	}

	savedRow := prepared.Mapping
	payload := studioV2WorkspaceDatabaseTargetResponse{
		ID:            savedRow.ID,
		WorkspaceID:   record.ID,
		PointID:       binding.PointID,
		TagID:         savedRow.TagID,
		ColumnName:    savedRow.ColumnName,
		Enabled:       savedRow.Enabled,
		RowGroupID:    rowGroupID,
		SaveState:     workspaceSaveStateSaved,
		SetupRevision: record.DatabaseSetupRevision,
		CreatedAt:     savedRow.CreatedAt,
		UpdatedAt:     savedRow.UpdatedAt,
	}
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, record.OrderedDeviceIDs, []string{binding.PointID, connector.ID})
	payload.RuntimeApplyStatus = applyOutcome.Status
	payload.RuntimeApplyMessage = applyOutcome.Message
	payload.RuntimeApplyIssues = applyOutcome.Issues
	h.recordDatabaseTargetAudit(c.Request.Context(), record.ID, savedRow, binding.PointID, payload.RuntimeApplyStatus)
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: payload})
}
