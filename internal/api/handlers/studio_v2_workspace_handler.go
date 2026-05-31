package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type StudioV2WorkspaceHandler struct {
	svc *workspace.Service
}

func NewStudioV2WorkspaceHandler(svc *workspace.Service) *StudioV2WorkspaceHandler {
	return &StudioV2WorkspaceHandler{svc: svc}
}

func (h *StudioV2WorkspaceHandler) Get(c *gin.Context) {
	record, err := h.svc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    record,
	})
}
