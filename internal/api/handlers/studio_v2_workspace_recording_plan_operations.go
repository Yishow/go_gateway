package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type testWriteRequest struct {
	PlanID      string `json:"plan_id" binding:"required"`
	StreamID    string `json:"stream_id"`
	TablePrefix string `json:"table_prefix"`
}

// TestWrite reports that recording test writes are not connected yet.
// @Summary Execute a recording test write
// @Description Recording test writes are unavailable until a verified operation is connected.
// @Tags Studio V2 Recording Plans
// @Accept json
// @Produce json
// @Param request body testWriteRequest true "Recording test write request"
// @Failure 400 {object} APIErrorResponse "Invalid recording test write request"
// @Failure 501 {object} APIErrorResponse "success=false; code=RECORDING_TEST_WRITE_NOT_IMPLEMENTED; message=recording test write is not implemented; retryable=false; action=wait_for_supported_operation; request_id is returned"
// @Router /v1/datalink/studio-v2/workspace/recording-plans/test-write [post]
func (h *StudioV2WorkspaceRecordingPlansHandler) TestWrite(c *gin.Context) {
	var req testWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		renderStudioV2WorkspaceValidationError(c, err)
		return
	}

	renderRecordingOperationNotImplemented(c, "RECORDING_TEST_WRITE_NOT_IMPLEMENTED", "recording test write is not implemented")
}

func renderRecordingOperationNotImplemented(c *gin.Context, code, message string) {
	c.JSON(http.StatusNotImplemented, gin.H{
		apiResponseSuccessKey: false,
		apiResponseErrorKey: TypedAPIErrorEnvelope{
			Code:      code,
			Message:   message,
			Retryable: false,
			RequestID: getOrGenerateRequestID(c),
			Action:    "wait_for_supported_operation",
		},
	})
}
