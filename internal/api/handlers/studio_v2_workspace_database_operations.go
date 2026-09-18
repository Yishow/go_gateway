package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
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
// @Success 200 {object} map[string]interface{} "success=true; data is the schema operation"
// @Failure 404 {object} APIErrorResponse "success=false; code=RECORDING_OPERATION_NOT_FOUND; unknown or foreign operation"
// @Failure 503 {object} APIErrorResponse "success=false; code=RECORDING_PLAN_UNAVAILABLE; retryable=true"
// @Router /v1/datalink/studio-v2/workspace/database-operations/{operation_id} [get]
func (h *StudioV2WorkspaceDatabaseOperationsHandler) Get(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}
	op, err := h.operations.GetSchemaOperation(c.Request.Context(), record.ID, strings.TrimSpace(c.Param("operation_id")))
	switch {
	case errors.Is(err, recordingplan.ErrSchemaOperationNotFound):
		renderRecordingPlanError(c, http.StatusNotFound, "RECORDING_OPERATION_NOT_FOUND", "schema operation not found", false, "check the operation id")
	case err != nil:
		renderRecordingPlanError(c, http.StatusServiceUnavailable, recordingPlanUnavailableCode, recordingPlanUnavailableMessage, true, "retry reading the operation")
	default:
		c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: op})
	}
}
