package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type schemaOperationReader interface {
	GetSchemaOperation(ctx context.Context, workspaceID, operationID string) (*recordingplan.SchemaOperation, error)
}

// StudioV2WorkspaceDatabaseOperationsHandler reports schema operations of the current workspace.
type StudioV2WorkspaceDatabaseOperationsHandler struct {
	workspaceSvc *workspace.Service
	operations   schemaOperationReader
}

// NewStudioV2WorkspaceDatabaseOperationsHandler creates the operation status handler.
func NewStudioV2WorkspaceDatabaseOperationsHandler(workspaceSvc *workspace.Service, operations schemaOperationReader) *StudioV2WorkspaceDatabaseOperationsHandler {
	return &StudioV2WorkspaceDatabaseOperationsHandler{workspaceSvc: workspaceSvc, operations: operations}
}

// Get reports one schema operation of the current workspace, including running,
// partial and unknown outcomes. Unknown and foreign operations share one safe 404.
// @Summary Read a schema operation
// @Description Returns the durable state of a schema operation issued by a preview of the current workspace.
// @Tags Studio V2 Recording Plans
// @Produce json
// @Param operation_id path string true "Operation ID returned with the schema preview"
// @Param locator_workspace_id query string false "Locator-origin workspace identity; requires revision pair"
// @Param locator_workspace_revision query string false "Exact observed workspace revision; present empty only matches legacy empty"
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Success 200 {object} map[string]interface{} "success=true; data is the schema operation"
// @Failure 404 {object} APIErrorResponse "success=false; code=RECORDING_OPERATION_NOT_FOUND; unknown or foreign operation"
// @Failure 503 {object} APIErrorResponse "success=false; code=RECORDING_PLAN_UNAVAILABLE; retryable=true"
// @Router /v1/datalink/studio-v2/workspace/database-operations/{operation_id} [get]
func (h *StudioV2WorkspaceDatabaseOperationsHandler) Get(c *gin.Context) {
	record, locatorScope, ok := beginStudioLocatorNativeRead(c, h.workspaceSvc, func() { renderStudioV2WorkspaceBootstrapError(c) })
	if !ok {
		return
	}
	var err error
	op, err := h.operations.GetSchemaOperation(c.Request.Context(), record.ID, strings.TrimSpace(c.Param("operation_id")))
	if !finishStudioLocatorNativeRead(c, h.workspaceSvc, locatorScope) {
		return
	}
	if locatorScope != nil && err == nil {
		if op == nil {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonUnavailable))
			return
		}
		if op.WorkspaceID != record.ID {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonForeign))
			return
		}
		if op.OperationID != strings.TrimSpace(c.Param("operation_id")) {
			renderRuntimeEvidenceError(c, runtimeEvidenceError(runtimeview.ReasonStale))
			return
		}
	}
	switch {
	case errors.Is(err, recordingplan.ErrSchemaOperationNotFound):
		renderRecordingPlanError(c, http.StatusNotFound, "RECORDING_OPERATION_NOT_FOUND", "schema operation not found", false, "check the operation id")
	case err != nil:
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry reading the operation")
	default:
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: op})
	}
}

// TestWriteRecoveryResponse carries a complete recovery observation. An error
// never carries a partial list and must not be treated as proof of no write.
type TestWriteRecoveryResponse struct {
	Success bool                                 `json:"success"`
	Data    *recordingplan.TestWriteRecoveryList `json:"data"`
}

type testWriteRecoveryReader interface {
	ListUnresolvedTestWrites(context.Context, string) (*recordingplan.TestWriteRecoveryList, error)
}

// ListTestWrites reports unresolved test writes of the server's current workspace.
// @Summary Recover unresolved test writes
// @Description Read-only ledger recovery independent of browser references. Complete empty results do not authorize replay of an unknown write. Query failures never return a partial list. Tokens and destination details are not exposed.
// @Tags Studio V2 Recording Plans
// @Produce json
// @Param locator_workspace_id query string false "Locator-origin workspace identity; requires revision pair"
// @Param locator_workspace_revision query string false "Exact observed workspace revision; present empty only matches legacy empty"
// @Failure 400 {object} APIErrorResponse
// @Failure 403 {object} APIErrorResponse
// @Failure 409 {object} APIErrorResponse
// @Success 200 {object} TestWriteRecoveryResponse
// @Failure 503 {object} APIErrorResponse "success=false; code=RECORDING_PLAN_UNAVAILABLE; retryable=true; recovery could not be established"
// @Router /v1/datalink/studio-v2/workspace/database-operations/test-writes [get]
func (h *StudioV2WorkspaceDatabaseOperationsHandler) ListTestWrites(c *gin.Context) {
	if !validateStudioLocatorNativeQuery(c) {
		return
	}
	unavailable := func() {
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry reading test write recovery")
	}
	reader, ok := h.operations.(testWriteRecoveryReader)
	if !ok || h.workspaceSvc == nil {
		unavailable()
		return
	}
	record, locatorScope, ready := beginStudioLocatorNativeRead(c, h.workspaceSvc, unavailable)
	if !ready {
		return
	}
	result, err := reader.ListUnresolvedTestWrites(c.Request.Context(), record.ID)
	if !finishStudioLocatorNativeRead(c, h.workspaceSvc, locatorScope) {
		return
	}
	if err != nil || result == nil || !result.Complete || result.WorkspaceID != record.ID {
		unavailable()
		return
	}
	c.JSON(http.StatusOK, TestWriteRecoveryResponse{Success: true, Data: result})
}
