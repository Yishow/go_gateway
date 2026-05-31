package handlers

import (
	"context"
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type workspaceActivator interface {
	ActivateEligible(ctx context.Context) (*workspace.ActivationResponse, error)
}

type StudioV2WorkspaceActivationHandler struct {
	activator workspaceActivator
}

func NewStudioV2WorkspaceActivationHandler(activator workspaceActivator) *StudioV2WorkspaceActivationHandler {
	return &StudioV2WorkspaceActivationHandler{activator: activator}
}

func (h *StudioV2WorkspaceActivationHandler) Activate(c *gin.Context) {
	response, err := h.activator.ActivateEligible(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
