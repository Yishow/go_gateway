package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// recordingSavedConnectorResolver resolves the saved target, reads its actual
// table metadata and runs confirmed schema statements on it.
type recordingSavedConnectorResolver interface {
	ResolveSavedTarget(ctx context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error)
	InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*dbtarget.TableInspection, error)
	ExecuteSchemaStatements(ctx context.Context, connectorID, schemaName string, statements []string) (*dbtarget.SchemaStatementExecution, error)
}

type StudioV2WorkspaceRecordingPlansHandler struct {
	workspaceSvc *workspace.Service
	planSvc      *recordingplan.Service
	dbTargetSvc  recordingSavedConnectorResolver
	membership   recordingPlanMembershipServices
}

func NewStudioV2WorkspaceRecordingPlansHandler(
	workspaceSvc *workspace.Service,
	planSvc *recordingplan.Service,
	dbTargetSvc recordingSavedConnectorResolver,
) *StudioV2WorkspaceRecordingPlansHandler {
	return &StudioV2WorkspaceRecordingPlansHandler{
		workspaceSvc: workspaceSvc,
		planSvc:      planSvc,
		dbTargetSvc:  dbTargetSvc,
	}
}

func (h *StudioV2WorkspaceRecordingPlansHandler) List(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	plans, err := h.planSvc.ListByWorkspace(c.Request.Context(), record.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseCodeKey:    "RECORDING_PLAN_LIST_FAILED",
				apiResponseMessageKey: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    plans,
	})
}

func (h *StudioV2WorkspaceRecordingPlansHandler) Create(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var plan recordingplan.RecordingPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	if workspaceID := strings.TrimSpace(plan.WorkspaceID); workspaceID != "" && workspaceID != record.ID {
		renderRecordingPlanTargetError(c, errRecordingPlanDestinationScope)
		return
	}
	plan.WorkspaceID = record.ID
	if err := h.validatePlanDestinations(c.Request.Context(), record, plan.Destinations); err != nil {
		renderRecordingPlanTargetError(c, err)
		return
	}
	if err := h.validateRecordingPlanMembership(c.Request.Context(), record, &plan); err != nil {
		renderRecordingPlanMembershipError(c, err)
		return
	}
	if err := h.planSvc.CreatePlan(c.Request.Context(), &plan); err != nil {
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_PLAN_CREATE_FAILED", "recording plan could not be created", false, "review the recording plan")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    plan,
	})
}

// Get returns a plan only from the current workspace; missing and foreign
// plans share the same safe 404.
func (h *StudioV2WorkspaceRecordingPlansHandler) Get(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	plan, err := h.getPlanForWorkspace(c.Request.Context(), record.ID, strings.TrimSpace(c.Param("id")))
	if err != nil {
		renderRecordingPlanReadError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    plan,
	})
}

// Update changes a plan inside the current workspace only. The plan must exist
// there before destinations or members are resolved.
func (h *StudioV2WorkspaceRecordingPlansHandler) Update(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	var plan recordingplan.RecordingPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	if _, err := h.getPlanForWorkspace(c.Request.Context(), record.ID, id); err != nil {
		renderRecordingPlanReadError(c, err)
		return
	}
	if workspaceID := strings.TrimSpace(plan.WorkspaceID); workspaceID != "" && workspaceID != record.ID {
		renderRecordingPlanTargetError(c, errRecordingPlanDestinationScope)
		return
	}
	plan.ID = id
	plan.WorkspaceID = record.ID

	if err := h.validatePlanDestinations(c.Request.Context(), record, plan.Destinations); err != nil {
		renderRecordingPlanTargetError(c, err)
		return
	}
	if err := h.validateRecordingPlanMembership(c.Request.Context(), record, &plan); err != nil {
		renderRecordingPlanMembershipError(c, err)
		return
	}
	if err := h.planSvc.UpdatePlanByWorkspace(c.Request.Context(), &plan, record.ID); err != nil {
		if errors.Is(err, recordingplan.ErrPlanNotFound) {
			renderRecordingPlanReadError(c, errRecordingPlanNotFound)
			return
		}
		renderRecordingPlanError(c, http.StatusUnprocessableEntity, "RECORDING_PLAN_UPDATE_FAILED", "recording plan could not be updated", false, "review the recording plan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    plan,
	})
}

// Delete removes a plan from the current workspace only.
func (h *StudioV2WorkspaceRecordingPlansHandler) Delete(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	id := strings.TrimSpace(c.Param("id"))
	if err := h.planSvc.DeletePlanByWorkspace(c.Request.Context(), id, record.ID); err != nil {
		if errors.Is(err, recordingplan.ErrPlanNotFound) {
			renderRecordingPlanReadError(c, errRecordingPlanNotFound)
			return
		}
		renderRecordingPlanError(c, http.StatusInternalServerError, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry deleting the recording plan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: gin.H{
			apiResponseDeletedKey: true,
			"id":                  id,
		},
	})
}

// Capabilities returns static capabilities or a read-only saved-connector probe.
// @Summary Read recording connector capabilities
// @Tags datalink
// @Produce json
// @Param kind query string false "Connector kind; required for a saved-connector probe"
// @Param connector_id query string false "Saved connector identity; paired with expected_connector_revision"
// @Param expected_connector_revision query string false "Expected saved connector revision"
// @Param locator_workspace_id query string false "Optional locator workspace identity; paired with locator_workspace_revision"
// @Param locator_workspace_revision query string false "Observed locator workspace revision; explicit empty is valid only for matching legacy-empty setup"
// @Success 200 {object} map[string]interface{} "Success envelope with a capability object or the static capability list"
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Failure 503 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/recording-plans/capabilities [get]
func (h *StudioV2WorkspaceRecordingPlansHandler) Capabilities(c *gin.Context) {
	// A saved-scope SQL Server query names kind=sqlserver explicitly; the
	// predicate gates on that kind itself, so the MySQL saved-scope branch
	// below keeps every MySQL query.
	if sqlserverRecordingScopeRequested(c) {
		h.sqlserverRecordingCapabilities(c)
		return
	}
	if mysqlRecordingScopeRequested(c) {
		h.mysqlRecordingCapabilities(c)
		return
	}
	kind := c.Query("kind")
	if kind != "" {
		capResult := h.planSvc.GetConnectorCapability(kind)
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: capResult})
		return
	}

	kinds := []string{recordingDialectSQLite, recordingDialectPostgres, recordingDialectMySQL, "sqlserver", "oracle"}
	caps := make([]recordingplan.ConnectorCapability, 0, len(kinds))
	for _, k := range kinds {
		caps = append(caps, h.planSvc.GetConnectorCapability(k))
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: caps})
}
