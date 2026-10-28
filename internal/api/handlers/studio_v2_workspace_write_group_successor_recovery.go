package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/workspace"
)

type successorOperationStatusRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorOperationStatusRequest
}

// firstStartConfirm starts a genuinely never-applied lineage. It never returns
// a cutover with an invented predecessor end, event, or ownership fence.
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) firstStartConfirm(c *gin.Context, raw []byte) {
	var req successorActivationConfirmRequest
	if decodeSuccessorRequest(raw, &req) != nil || !validActivationEnvelope(req.Envelope) || req.Envelope.ActivationKind != "first_start" || len(req.ReviewDigest) != 64 || len(req.ReadinessDigest) != 64 {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.ConfirmSuccessorFirstStart(c.Request.Context(), c.Param("id"), req.WriteGroupSuccessorActivationConfirmation)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

// operationStatus reads the original root operation in one local transaction.
// It deliberately requires neither a current revision nor live target access.
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) operationStatus(c *gin.Context, raw []byte) {
	var req successorOperationStatusRequest
	if decodeSuccessorRequest(raw, &req) != nil || !successorExactIdentity(req.WorkspaceID) || !successorExactIdentity(req.LineageID) || !successorExactIdentity(req.OperationID) {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.SuccessorOperationStatus(c.Request.Context(), c.Param("id"), req.WriteGroupSuccessorOperationStatusRequest)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func successorExactIdentity(value string) bool {
	return value != "" && strings.TrimSpace(value) == value
}
