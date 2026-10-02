package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/grouptestwrite"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

const (
	testWriteNotFoundCode        = "WRITE_GROUP_TEST_WRITE_PREVIEW_NOT_FOUND"
	testWriteStaleCode           = "WRITE_GROUP_TEST_WRITE_PREVIEW_STALE"
	testWriteExpiredCode         = "WRITE_GROUP_TEST_WRITE_PREVIEW_EXPIRED"
	testWriteKindCode            = "WRITE_GROUP_TEST_WRITE_TOKEN_KIND"
	testWriteMismatchCode        = "WRITE_GROUP_TEST_WRITE_OPERATION_MISMATCH"
	testWriteBusyCode            = "WRITE_GROUP_TEST_WRITE_BUSY"
	testWriteUnsupportedCode     = "WRITE_GROUP_TEST_WRITE_UNSUPPORTED"
	testWriteResultUnknownCode   = "WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN"
	testWritePlanUnresolvedCode  = "RECORDING_TEST_WRITE_PLAN_UNRESOLVED"
	testWriteUnavailableCode     = "WRITE_GROUP_TEST_WRITE_UNAVAILABLE"
	testWritePreviewAgainAction  = "preview the test write again"
	testWriteUnavailableResponse = "test write is not available"
)

// WriteGroupTestWriter previews and confirms explicit test writes of one group.
type WriteGroupTestWriter interface {
	Preview(ctx context.Context, ws grouptestwrite.Workspace, groupID string) (*grouptestwrite.Preview, error)
	Confirm(ctx context.Context, ws grouptestwrite.Workspace, req grouptestwrite.Confirmation) (*grouptestwrite.Outcome, error)
}

// RecordingPlanGroupResolver maps a legacy recording plan to its canonical group.
type RecordingPlanGroupResolver interface {
	ResolveRecordingPlanGroup(ctx context.Context, planID string) (string, error)
}

// StudioV2WorkspaceWriteGroupTestWriteHandler exposes explicit test writes for
// canonical groups and the legacy plan routes that resolve to them.
type StudioV2WorkspaceWriteGroupTestWriteHandler struct {
	workspaceSvc *workspace.Service
	testWriter   WriteGroupTestWriter
	plans        RecordingPlanGroupResolver
}

// NewStudioV2WorkspaceWriteGroupTestWriteHandler builds the handler; a nil
// writer makes every route answer that test writes are unavailable.
func NewStudioV2WorkspaceWriteGroupTestWriteHandler(workspaceSvc *workspace.Service, testWriter WriteGroupTestWriter, plans RecordingPlanGroupResolver) *StudioV2WorkspaceWriteGroupTestWriteHandler {
	return &StudioV2WorkspaceWriteGroupTestWriteHandler{workspaceSvc: workspaceSvc, testWriter: testWriter, plans: plans}
}

// WriteGroupTestWriteConfirmRequest confirms a previewed test write.
type WriteGroupTestWriteConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	OperationID string `json:"operation_id" binding:"required"`
}

// WriteGroupTestWritePreviewResponse is the mutation-free preview.
type WriteGroupTestWritePreviewResponse struct {
	Success bool                    `json:"success"`
	Data    *grouptestwrite.Preview `json:"data"`
}

// WriteGroupTestWriteOperationResponse is a test-write operation, running or finished.
type WriteGroupTestWriteOperationResponse struct {
	Success bool                           `json:"success"`
	Data    *recordingplan.SchemaOperation `json:"data"`
}

// legacyTestWritePreviewRequest is the legacy plan-addressed preview body.
type legacyTestWritePreviewRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}

// legacyTestWriteConfirmRequest is the legacy confirmation body. It must carry
// the token and operation of a preview; a bare plan_id never writes.
type legacyTestWriteConfirmRequest struct {
	PlanID      string `json:"plan_id"`
	Token       string `json:"token" binding:"required"`
	OperationID string `json:"operation_id" binding:"required"`
}

func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) workspaceRef(c *gin.Context) (grouptestwrite.Workspace, bool) {
	if h.testWriter == nil {
		renderTestWriteError(c, http.StatusServiceUnavailable, testWriteUnavailableCode, testWriteUnavailableResponse, true, "retry later", "")
		return grouptestwrite.Workspace{}, false
	}
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return grouptestwrite.Workspace{}, false
	}
	return grouptestwrite.Workspace{ID: record.ID, Revision: record.DatabaseSetupRevision}, true
}

// Preview describes the exact test row and its cleanup for a saved group and
// persists an action-bound token. It reads the group and the real target table
// but never writes to the target.
// @Summary Preview a write-group test write
// @Description Returns the test row, its operation-owned cleanup and an action-bound token without changing the target.
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} WriteGroupTestWritePreviewResponse
// @Failure 404 {object} APIErrorResponse "unknown group"
// @Failure 409 {object} APIErrorResponse "destination revision changed"
// @Failure 422 {object} APIErrorResponse "code=WRITE_GROUP_TEST_WRITE_UNSUPPORTED; the target cannot safely identify and remove an operation-owned test row"
// @Failure 503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/test-write-preview [post]
func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) Preview(c *gin.Context) {
	h.preview(c, c.Param("id"))
}

func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) preview(c *gin.Context, groupID string) {
	ws, ok := h.workspaceRef(c)
	if !ok {
		return
	}
	preview, err := h.testWriter.Preview(c.Request.Context(), ws, groupID)
	if err != nil {
		renderTestWriteFailure(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, WriteGroupTestWritePreviewResponse{Success: true, Data: preview})
}

// Confirm runs a previewed test write once. A running duplicate returns 202, a
// retained result returns 200 without touching the target (even after the
// token expired), and another operation holding the table returns 409.
// @Summary Confirm a write-group test write
// @Description Writes one operation-owned test row through the production row codec, reads it back, then removes only that row. Results keep write verification and cleanup status separate.
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body WriteGroupTestWriteConfirmRequest true "Test write confirmation"
// @Success 200 {object} WriteGroupTestWriteOperationResponse "written_verified, written_unverified, failed or unknown, with cleanup_status"
// @Success 202 {object} WriteGroupTestWriteOperationResponse "still running"
// @Failure 400 {object} APIErrorResponse "missing token or operation_id"
// @Failure 404 {object} APIErrorResponse "unknown or foreign token"
// @Failure 409 {object} APIErrorResponse "stale or expired first claim, or another operation holds the table"
// @Failure 422 {object} APIErrorResponse "token for another action, or unsupported target"
// @Failure 503 {object} APIErrorResponse "code=WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN; the result could not be recorded"
// @Router /datalink/studio-v2/workspace/write-groups/{id}/test-write [post]
func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) Confirm(c *gin.Context) {
	var req WriteGroupTestWriteConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	h.confirm(c, c.Param("id"), req.Token, req.OperationID)
}

func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) confirm(c *gin.Context, groupID, token, operationID string) {
	token, operationID = strings.TrimSpace(token), strings.TrimSpace(operationID)
	if token == "" || operationID == "" {
		renderStudioV2WorkspaceValidationError(c, errors.New("token and operation_id are required"))
		return
	}
	ws, ok := h.workspaceRef(c)
	if !ok {
		return
	}
	outcome, err := h.testWriter.Confirm(c.Request.Context(), ws, grouptestwrite.Confirmation{Token: token, OperationID: operationID, GroupID: groupID})
	if err != nil {
		var op *recordingplan.SchemaOperation
		if outcome != nil {
			op = outcome.Operation
		}
		renderTestWriteFailure(c, err, op)
		return
	}
	switch outcome.Claim {
	case recordingplan.ClaimAcquired, recordingplan.ClaimCompleted:
		c.JSON(http.StatusOK, WriteGroupTestWriteOperationResponse{Success: true, Data: outcome.Operation})
	case recordingplan.ClaimInProgress:
		c.JSON(http.StatusAccepted, WriteGroupTestWriteOperationResponse{Success: true, Data: outcome.Operation})
	case recordingplan.ClaimScopeBusy:
		renderTestWriteError(c, http.StatusConflict, testWriteBusyCode, "another test write is using this table", true,
			"wait for the running test write, then retry", outcome.Operation.OperationID)
	default:
		renderTestWriteError(c, http.StatusServiceUnavailable, testWriteUnavailableCode, testWriteUnavailableResponse, true, "retry later", "")
	}
}

// LegacyPreview resolves a legacy recording plan to its canonical group and
// previews that group. A plan that does not map to exactly one group is refused.
// @Summary Preview a test write through a legacy recording plan
// @Tags Studio V2 Recording Plans
// @Accept json
// @Produce json
// @Param request body legacyTestWritePreviewRequest true "Legacy plan reference"
// @Success 200 {object} WriteGroupTestWritePreviewResponse
// @Failure 400 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse "code=RECORDING_TEST_WRITE_PLAN_UNRESOLVED or WRITE_GROUP_TEST_WRITE_UNSUPPORTED"
// @Router /v1/datalink/studio-v2/workspace/recording-plans/test-write-preview [post]
func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) LegacyPreview(c *gin.Context) {
	var req legacyTestWritePreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	groupID, ok := h.resolvePlan(c, req.PlanID)
	if !ok {
		return
	}
	h.preview(c, groupID)
}

// LegacyConfirm confirms a previewed test write addressed by the legacy plan
// route. It needs the preview's token and operation; a bare plan_id is refused
// before anything is resolved or written.
// @Summary Confirm a test write through a legacy recording plan
// @Tags Studio V2 Recording Plans
// @Accept json
// @Produce json
// @Param request body legacyTestWriteConfirmRequest true "Test write confirmation"
// @Success 200 {object} WriteGroupTestWriteOperationResponse
// @Success 202 {object} WriteGroupTestWriteOperationResponse
// @Failure 400 {object} APIErrorResponse "missing token or operation_id"
// @Failure 409 {object} APIErrorResponse
// @Failure 422 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/recording-plans/test-write [post]
func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) LegacyConfirm(c *gin.Context) {
	var req legacyTestWriteConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}
	groupID := ""
	if strings.TrimSpace(req.PlanID) != "" {
		resolved, ok := h.resolvePlan(c, req.PlanID)
		if !ok {
			return
		}
		groupID = resolved
	}
	h.confirm(c, groupID, req.Token, req.OperationID)
}

func (h *StudioV2WorkspaceWriteGroupTestWriteHandler) resolvePlan(c *gin.Context, planID string) (string, bool) {
	if h.plans == nil {
		renderTestWriteError(c, http.StatusServiceUnavailable, testWriteUnavailableCode, testWriteUnavailableResponse, true, "retry later", "")
		return "", false
	}
	groupID, err := h.plans.ResolveRecordingPlanGroup(c.Request.Context(), planID)
	if err != nil {
		if errors.Is(err, workspace.ErrRecordingPlanGroupUnresolved) {
			renderTestWriteError(c, http.StatusUnprocessableEntity, testWritePlanUnresolvedCode,
				"the recording plan does not map to exactly one write group", false, "use the write group test write instead", "")
			return "", false
		}
		renderWriteGroupError(c, err)
		return "", false
	}
	return groupID, true
}

func renderTestWriteFailure(c *gin.Context, err error, op *recordingplan.SchemaOperation) {
	var unsupported *grouptestwrite.UnsupportedError
	operationID := ""
	if op != nil {
		operationID = op.OperationID
	}
	switch {
	case errors.Is(err, grouptestwrite.ErrResultUnacknowledged):
		renderTestWriteError(c, http.StatusServiceUnavailable, testWriteResultUnknownCode, "the test write result could not be recorded", true,
			"check the operation status before retrying", operationID)
	case errors.As(err, &unsupported):
		renderTestWriteError(c, http.StatusUnprocessableEntity, testWriteUnsupportedCode,
			"test write is not supported for this group: "+unsupported.Reason, false, "choose a table with an entity key column the gateway can use to remove its own test row", "")
	case errors.Is(err, grouptestwrite.ErrGroupNotFound), errors.Is(err, workspace.ErrWriteGroupNotFound):
		renderWriteGroupTyped(c, http.StatusNotFound, "WRITE_GROUP_NOT_FOUND", "write group or resource is not available", false, "reload")
	case errors.Is(err, recordingplan.ErrPreviewTokenNotFound):
		renderTestWriteError(c, http.StatusNotFound, testWriteNotFoundCode, "test write preview not found", false, testWritePreviewAgainAction, "")
	case errors.Is(err, recordingplan.ErrPreviewTokenKind):
		renderTestWriteError(c, http.StatusUnprocessableEntity, testWriteKindCode, "the token was issued for a different action", false, testWritePreviewAgainAction, "")
	case errors.Is(err, recordingplan.ErrSchemaOperationMismatch):
		renderTestWriteError(c, http.StatusConflict, testWriteMismatchCode, "operation does not belong to this test write preview", false, "use the operation id returned with the preview", "")
	case errors.Is(err, recordingplan.ErrPreviewTokenExpired):
		renderTestWriteError(c, http.StatusConflict, testWriteExpiredCode, "test write preview expired", false, testWritePreviewAgainAction, "")
	case errors.Is(err, recordingplan.ErrPreviewTokenStale), errors.Is(err, recordingplan.ErrPreviewTokenLegacy), errors.Is(err, grouptestwrite.ErrDestinationChanged), errors.Is(err, grouptestwrite.ErrGroupMismatch):
		renderTestWriteError(c, http.StatusConflict, testWriteStaleCode, "test write preview is stale", false, testWritePreviewAgainAction, "")
	default:
		renderTestWriteError(c, http.StatusServiceUnavailable, testWriteUnavailableCode, testWriteUnavailableResponse, true, "retry later", operationID)
	}
}

func renderTestWriteError(c *gin.Context, status int, code, message string, retryable bool, action, operationID string) {
	c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: TypedAPIErrorEnvelope{
		Code: code, Message: message, Retryable: retryable, RequestID: getOrGenerateRequestID(c), Action: action, OperationID: operationID,
	}})
}
