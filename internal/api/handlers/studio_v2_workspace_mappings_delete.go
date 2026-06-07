package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StudioV2WorkspaceMappingsHandler) Delete(c *gin.Context) {
	_, rule, links, link, mappingRecord, _, ok := h.requireWorkspaceMapping(c)
	if !ok {
		return
	}
	if err := h.mappingSvc.Delete(c.Request.Context(), mappingRecord.ID); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	oldTagID := mappingRecord.TagID
	link.MappingID = nil
	link.TagID = nil
	link.UpdatedAt = time.Now()
	if err := h.ruleSvc.ReplaceLinks(c.Request.Context(), rule.ID, links); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}
	if err := h.deleteOrphanRuleManagedTag(c.Request.Context(), oldTagID); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, []string{rule.DeviceID, link.PointID})
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    mapStudioV2RuntimeApplyResponse(applyOutcome),
	})
}
