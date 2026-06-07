package handlers

import (
	"context"
	"errors"
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type workspaceActivator interface {
	ActivateEligible(ctx context.Context) (*workspace.ActivationResponse, error)
}

// workspaceSchemaEnsurer 在啟動採集前確保 workspace database target 對應的
// 資料表已存在（CREATE TABLE IF NOT EXISTS 為冪等操作）。
type workspaceSchemaEnsurer interface {
	EnsureWorkspaceSchema(ctx context.Context) error
}

type StudioV2WorkspaceActivationHandler struct {
	activator     workspaceActivator
	schemaEnsurer workspaceSchemaEnsurer
	auditSvc      workspaceAuditRecorder
}

func NewStudioV2WorkspaceActivationHandler(activator workspaceActivator, schemaEnsurer ...workspaceSchemaEnsurer) *StudioV2WorkspaceActivationHandler {
	handler := &StudioV2WorkspaceActivationHandler{activator: activator}
	if len(schemaEnsurer) > 0 {
		handler.schemaEnsurer = schemaEnsurer[0]
	}
	return handler
}

func (h *StudioV2WorkspaceActivationHandler) Activate(c *gin.Context) {
	// 先確保目標資料表存在，避免採集到值卻因缺表而寫入失敗。
	if h.schemaEnsurer != nil {
		if err := h.schemaEnsurer.EnsureWorkspaceSchema(c.Request.Context()); err != nil {
			renderStudioV2WorkspaceSchemaEnsureError(c, err)
			return
		}
	}

	response, err := h.activator.ActivateEligible(c.Request.Context())
	if err != nil {
		var blockedErr *workspace.ReadinessBlockedError
		if errors.As(err, &blockedErr) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "readiness_blocked",
					"message": blockedErr.Error(),
					"issues":  blockedErr.BlockingIssues(),
					"summary": blockedErr.Summary,
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"message": err.Error(),
			},
		})
		return
	}
	h.recordActivationAudit(c.Request.Context(), response)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
