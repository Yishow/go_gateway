package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// RecordingStartResponse separates start progress from production SQL delivery.
type RecordingStartResponse struct {
	Success bool                               `json:"success"`
	Data    *workspace.RecordingStartOperation `json:"data"`
}

// StudioV2RecordingStartHandler exposes only the scoped application operation.
type StudioV2RecordingStartHandler struct {
	start  *workspace.RecordingStartService
	groups *workspace.WriteGroupService
}

// NewStudioV2RecordingStartHandler shares the existing configuration services.
func NewStudioV2RecordingStartHandler(start *workspace.RecordingStartService, groups *workspace.WriteGroupService) *StudioV2RecordingStartHandler {
	return &StudioV2RecordingStartHandler{start: start, groups: groups}
}

func strictRecordingStartBody(c *gin.Context, destination any) error {
	if c.Request.Body == nil {
		return workspace.ErrRecordingStartInvalid
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return workspace.ErrRecordingStartInvalid
	}
	return nil
}

// Start persists or resumes one selected intent without schema mutation.
// @Summary Start recording for selected devices and groups
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body workspace.RecordingStartRequest true "Scoped start identity and expected revisions"
// @Success 200 {object} RecordingStartResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/recording-start [post]
func (h *StudioV2RecordingStartHandler) Start(c *gin.Context) {
	if h.start == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	var request workspace.RecordingStartRequest
	if strictRecordingStartBody(c, &request) != nil {
		renderWriteGroupTyped(c, http.StatusBadRequest, "RECORDING_START_INVALID", "recording start request is incomplete", false, "review_selection")
		return
	}
	result, err := h.start.Start(c.Request.Context(), request)
	if err != nil {
		renderRecordingStartError(c, err)
		return
	}
	c.JSON(http.StatusOK, RecordingStartResponse{Success: true, Data: result})
}

// Get reads the durable operation; it never repeats a start action.
// @Summary Read scoped recording start progress
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Operation ID"
// @Success 200 {object} RecordingStartResponse
// @Failure 404,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/recording-start/operations/{id} [get]
func (h *StudioV2RecordingStartHandler) Get(c *gin.Context) {
	if h.start == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.start.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderRecordingStartError(c, err)
		return
	}
	c.JSON(http.StatusOK, RecordingStartResponse{Success: true, Data: result})
}

// EnsureBasic reuses the persisted workspace/device/basic-role identity.
// @Summary Save or reuse one device's basic managed group
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param request body WriteGroupMutationRequest true "Basic managed draft and expected revisions"
// @Success 200 {object} WriteGroupSaveResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/basic/{id} [post]
func (h *StudioV2RecordingStartHandler) EnsureBasic(c *gin.Context) {
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	var request WriteGroupMutationRequest
	if strictRecordingStartBody(c, &request) != nil || request.Group == nil ||
		request.ExpectedWorkspaceRevision == nil || request.ExpectedConnectorRevision == nil || request.ExpectedGroupRevision != nil {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "basic group request is incomplete", false, "review_selection")
		return
	}
	result, err := h.groups.EnsureBasicManaged(c.Request.Context(), c.Param("id"), workspace.WriteGroupMutation{
		WorkspaceID: request.WorkspaceID, ExpectedWorkspaceRevision: *request.ExpectedWorkspaceRevision,
		ExpectedConnectorRevision: *request.ExpectedConnectorRevision, Group: request.Group,
	})
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupSaveResponse{Success: true, Data: result})
}

func renderRecordingStartError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, workspace.ErrRecordingStartBusy):
		renderWriteGroupTyped(c, http.StatusConflict, "RECORDING_START_BUSY", "another recording start is still running; retry this request after it completes", true, "retry_same_request")
	case errors.Is(err, workspace.ErrRecordingStartInvalid):
		renderWriteGroupTyped(c, http.StatusBadRequest, "RECORDING_START_INVALID", "recording start selection is invalid", false, "review_selection")
	case errors.Is(err, recordingplan.ErrSchemaOperationMismatch):
		renderWriteGroupTyped(c, http.StatusConflict, "RECORDING_START_INTENT_CHANGED", "request identity belongs to a different intent", false, "revalidate")
	case errors.Is(err, recordingplan.ErrSchemaOperationNotFound):
		renderWriteGroupTyped(c, http.StatusNotFound, "RECORDING_START_NOT_FOUND", "recording start operation is unavailable", false, "reload")
	default:
		renderWriteGroupError(c, err)
	}
}
