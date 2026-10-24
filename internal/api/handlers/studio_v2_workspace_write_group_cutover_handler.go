package handlers

import (
	"net/http"
	"strings"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupCutoverReviewResponse is the read-only difference of a cutover.
type WriteGroupCutoverReviewResponse struct {
	Success bool                     `json:"success"`
	Data    *workspace.CutoverReview `json:"data"`
}

// WriteGroupCutoverResponse is the immutable record of an applied cutover.
type WriteGroupCutoverResponse struct {
	Success bool                         `json:"success"`
	Data    *workspace.WriteGroupCutover `json:"data"`
}

// ReviewCutover returns, without changing anything, which sources a forward
// cutover (or its rollback) moves, what ends and starts at the exact boundary,
// and every reason it cannot be applied. The result carries a digest that
// the apply must present.
// @Summary Review a forward cutover or rollback
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Source write group ID"
// @Param request body workspace.CutoverRequest true "Cutover request"
// @Success 200 {object} WriteGroupCutoverReviewResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/cutover-review [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ReviewCutover(c *gin.Context) {
	request, ok := h.cutoverRequest(c)
	if !ok {
		return
	}
	review, err := h.groups.ReviewCutover(c.Request.Context(), request)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupCutoverReviewResponse{Success: true, Data: review})
}

// ApplyCutover applies a reviewed cutover as one compare-and-swap transaction.
// A stale revision, a changed review or any failure leaves the group that held
// the sources exactly as it was.
// @Summary Apply a reviewed forward cutover or rollback
// @Tags Studio V2 Write Groups
// @Accept json
// @Produce json
// @Param id path string true "Source write group ID"
// @Param request body workspace.CutoverRequest true "Reviewed cutover request"
// @Success 200 {object} WriteGroupCutoverResponse
// @Failure 400,404,409,422,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/cutover [post]
func (h *StudioV2WorkspaceWriteGroupsHandler) ApplyCutover(c *gin.Context) {
	request, ok := h.cutoverRequest(c)
	if !ok {
		return
	}
	cutover, err := h.groups.ApplyCutover(c.Request.Context(), request)
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupCutoverResponse{Success: true, Data: cutover})
}

func (h *StudioV2WorkspaceWriteGroupsHandler) cutoverRequest(c *gin.Context) (workspace.CutoverRequest, bool) {
	if !h.available(c) {
		return workspace.CutoverRequest{}, false
	}
	var request workspace.CutoverRequest
	if err := c.ShouldBindJSON(&request); err != nil || strings.TrimSpace(c.Param("id")) == "" {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "write-group request is incomplete", false, "review_request")
		return workspace.CutoverRequest{}, false
	}
	// The path names the source; a body naming another is a client mistake.
	if body := strings.TrimSpace(request.SourceGroupID); body != "" && body != strings.TrimSpace(c.Param("id")) {
		renderWriteGroupTyped(c, http.StatusBadRequest, "WRITE_GROUP_INVALID_REQUEST", "write-group request is incomplete", false, "review_request")
		return workspace.CutoverRequest{}, false
	}
	request.SourceGroupID = c.Param("id")
	return request, true
}

// WriteGroupCutoverListResponse lists the recorded cutovers of a source group.
type WriteGroupCutoverListResponse struct {
	Success bool                           `json:"success"`
	Data    []*workspace.WriteGroupCutover `json:"data"`
}

// GetCutover reads one recorded cutover of a source group, so an operator whose
// answer to an apply was lost can see that it happened and what it recorded.
// @Summary Read a recorded cutover
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Source write group ID"
// @Param cutover_id path string true "Cutover ID"
// @Success 200 {object} WriteGroupCutoverResponse
// @Failure 404,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/cutovers/{cutover_id} [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) GetCutover(c *gin.Context) {
	if !h.available(c) {
		return
	}
	cutover, err := h.groups.GetCutover(c.Request.Context(), c.Param("id"), c.Param("cutover_id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupCutoverResponse{Success: true, Data: cutover})
}

// ListCutovers lists the recorded cutovers and rollbacks of a source group.
// @Summary List recorded cutovers
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Source write group ID"
// @Success 200 {object} WriteGroupCutoverListResponse
// @Failure 404,500,503 {object} APIErrorResponse
// @Router /datalink/studio-v2/workspace/write-groups/{id}/cutovers [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) ListCutovers(c *gin.Context) {
	if !h.available(c) {
		return
	}
	cutovers, err := h.groups.ListCutovers(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupCutoverListResponse{Success: true, Data: cutovers})
}
