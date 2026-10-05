package handlers

import (
	"context"
	"errors"
	"net/http"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupResolutionResponse contains durable operator-decision evidence.
type WriteGroupResolutionResponse struct {
	Success bool                          `json:"success"`
	Data    *groupdelivery.DecisionResult `json:"data"`
}

// ResolveDelivery records a repair retry or an explicitly confirmed skip.
// @Summary Resolve a blocked or quarantined write-group row
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body groupdelivery.OperatorDecision true "Scoped CAS decision"
// @Success 200 {object} WriteGroupResolutionResponse
// @Failure 400,404,409,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/delivery/resolve [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ResolveDelivery(c *gin.Context) {
	if !h.available(c) {
		return
	}
	group, err := h.groups.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	resolver, ok := h.delivery.(interface {
		ResolveAttention(context.Context, string, string, groupdelivery.OperatorDecision) (*groupdelivery.DecisionResult, error)
	})
	if !ok {
		renderWriteGroupError(c, workspace.ErrWriteGroupServiceUnavailable)
		return
	}
	var request groupdelivery.OperatorDecision
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: "invalid_delivery_decision"})
		return
	}
	result, err := resolver.ResolveAttention(c.Request.Context(), group.Group.WorkspaceID, group.Group.ID, request)
	if err != nil {
		status, code := http.StatusServiceUnavailable, "delivery_resolution_unavailable"
		switch {
		case errors.Is(err, groupdelivery.ErrInvalidDecision):
			status, code = http.StatusBadRequest, "invalid_delivery_decision"
		case errors.Is(err, groupdelivery.ErrOutboxItemNotFound):
			status, code = http.StatusNotFound, "delivery_row_not_found"
		case errors.Is(err, groupdelivery.ErrDecisionConflict):
			status, code = http.StatusConflict, "delivery_decision_conflict"
		}
		c.JSON(status, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: code})
		return
	}
	c.JSON(http.StatusOK, WriteGroupResolutionResponse{Success: true, Data: result})
}
