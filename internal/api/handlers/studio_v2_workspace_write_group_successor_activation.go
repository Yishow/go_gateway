package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/datalink/workspace"
)

type successorActivationReviewRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorActivationRequest
}

type successorActivationReadinessRequest struct {
	Action       string                                          `json:"action"`
	Envelope     workspace.WriteGroupSuccessorActivationEnvelope `json:"envelope"`
	ReviewDigest string                                          `json:"review_digest"`
}

type successorActivationConfirmRequest struct {
	Action string `json:"action"`
	workspace.WriteGroupSuccessorActivationConfirmation
}

// activationReview reads local evidence only. It never probes or activates G2.
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) activationReview(c *gin.Context, raw []byte) {
	var req successorActivationReviewRequest
	if decodeSuccessorRequest(raw, &req) != nil || !validSuccessorStage(req.Stage) || req.Boundary.IsZero() {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.SuccessorActivationReview(c.Request.Context(), c.Param("id"), req.WriteGroupSuccessorActivationRequest)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

// activationReadiness is an explicitly selected target inspection, never DDL.
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) activationReadiness(c *gin.Context, raw []byte) {
	var req successorActivationReadinessRequest
	if decodeSuccessorRequest(raw, &req) != nil || !validActivationEnvelope(req.Envelope) || len(req.ReviewDigest) != 64 {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.SuccessorActivationReadiness(c.Request.Context(), c.Param("id"), workspace.WriteGroupSuccessorActivationConfirmation{Envelope: req.Envelope, ReviewDigest: req.ReviewDigest})
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func validActivationEnvelope(envelope workspace.WriteGroupSuccessorActivationEnvelope) bool {
	return envelope.Version != "" && validSuccessorStage(envelope.Request.Stage) && !envelope.Request.Boundary.IsZero()
}

// activationConfirm is distinct from schema confirmation and rechecks readiness.
func (h *StudioV2WorkspaceWriteGroupSuccessorHandler) activationConfirm(c *gin.Context, raw []byte) {
	var req successorActivationConfirmRequest
	if decodeSuccessorRequest(raw, &req) != nil || !validActivationEnvelope(req.Envelope) || (req.Envelope.ActivationKind != "" && req.Envelope.ActivationKind != "cutover") || len(req.ReviewDigest) != 64 || len(req.ReadinessDigest) != 64 {
		invalidSuccessorRequest(c)
		return
	}
	if h.groups == nil {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	result, err := h.groups.ConfirmSuccessorActivation(c.Request.Context(), c.Param("id"), req.WriteGroupSuccessorActivationConfirmation)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}
