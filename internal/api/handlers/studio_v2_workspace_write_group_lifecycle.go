package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// WriteGroupLifecycleResponse exposes safe, revision-scoped recording actions.
type WriteGroupLifecycleResponse struct {
	Success bool                           `json:"success"`
	Data    *workspace.WriteGroupLifecycle `json:"data"`
}

// Lifecycle reads action availability without authorizing a later mutation.
// @Summary Read saved recording lifecycle actions
// @Description Local read only: saved status/actions, immutable per-generation history and sanitized successor_seed. Actions and seed are advisory, never mutation authority. Saved configuration, runtime health and delivery outcomes remain distinct; samples, outbox and unknown counts use separate units.
// @Tags Studio V2 Write Groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} WriteGroupLifecycleResponse
// @Failure 404,409,422,500,503 {object} APIErrorResponse
// @Router /v1/datalink/studio-v2/workspace/write-groups/{id}/lifecycle [get]
func (h *StudioV2WorkspaceWriteGroupsHandler) Lifecycle(c *gin.Context) {
	if !h.available(c) {
		return
	}
	result, err := h.groups.Lifecycle(c.Request.Context(), c.Param("id"))
	if err != nil {
		renderWriteGroupError(c, err)
		return
	}
	c.JSON(http.StatusOK, WriteGroupLifecycleResponse{Success: true, Data: result})
}
