package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// Lookup reads a recorded start without submitting or resuming it. POST keeps
// the original readiness proof and full request out of the URL/query string.
// @Summary Look up a scoped recording start by its original request
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param request body workspace.RecordingStartRequest true "Original complete scoped start request"
// @Success 200 {object} RecordingStartResponse
// @Failure 400,404,409,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/recording-start/lookup [post]
func (h *StudioV2RecordingStartHandler) Lookup(c *gin.Context) {
	if h.start == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	var request workspace.RecordingStartRequest
	if strictRecordingStartBody(c, &request) != nil {
		renderWriteGroupTyped(c, http.StatusBadRequest, "RECORDING_START_INVALID", "recording start request is incomplete", false, "review_selection")
		return
	}
	result, err := h.start.Lookup(c.Request.Context(), request)
	if err != nil {
		renderRecordingStartError(c, err)
		return
	}
	c.JSON(http.StatusOK, RecordingStartResponse{Success: true, Data: result})
}
