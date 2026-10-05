package handlers

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/mapping"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *StudioV2WorkspaceMappingsHandler) Delete(c *gin.Context) {
	row, ok := h.requireWorkspaceMapping(c)
	if !ok {
		return
	}

	defer row.release()
	rule, links, link, mappingRecord := row.rule, row.links, row.link, row.mapping
	previousLinks := cloneWorkspaceLinks(links)

	oldTagID := mappingRecord.TagID
	link.MappingID = nil
	link.TagID = nil
	link.UpdatedAt = time.Now()
	if err := h.ruleSvc.ReplaceLinks(c.Request.Context(), rule.ID, links); err != nil {
		renderStudioV2WorkspaceMappingError(c, err)
		return
	}

	if deleteErr := h.mappingSvc.Delete(c.Request.Context(), mappingRecord.ID); deleteErr != nil {
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(c.Request.Context()), 5*time.Second)
		defer cancel()
		current, readErr := h.mappingSvc.GetByID(rollbackCtx, mappingRecord.ID)
		switch {
		case readErr == nil && current.PointID == mappingRecord.PointID && current.TagID == mappingRecord.TagID:
			deleteErr = errors.Join(deleteErr, h.ruleSvc.ReplaceLinks(rollbackCtx, rule.ID, previousLinks))
		case readErr != nil:
			deleteErr = errors.Join(deleteErr, readErr)
		default:
			deleteErr = errors.Join(deleteErr, mapping.ErrWorkspaceConfirmationConflict)
		}
		renderStudioV2WorkspaceMappingError(c, deleteErr)
		return
	}
	cleanupStatus, cleanupMessage := "complete", ""
	if err := h.deleteOrphanRuleManagedTag(c.Request.Context(), oldTagID); err != nil {
		cleanupStatus, cleanupMessage = workspaceMappingCleanupFailedStatus, "Mapping deleted; orphan tag cleanup could not be completed"
	}

	row.release()
	applyOutcome := resolveStudioV2ScopedRuntimeApplyOutcome(c.Request.Context(), h.workspaceSvc, h.deviceSvc, []string{rule.DeviceID}, []string{rule.DeviceID, link.PointID})
	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: struct {
			studioV2RuntimeApplyResponse
			CleanupStatus  string `json:"cleanup_status"`
			CleanupMessage string `json:"cleanup_message,omitempty"`
		}{mapStudioV2RuntimeApplyResponse(applyOutcome), cleanupStatus, cleanupMessage},
	})
}
