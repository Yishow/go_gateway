package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceBootstrapResponse struct {
	ID               string                      `json:"id"`
	Kind             workspace.Kind              `json:"kind"`
	Status           workspace.Status            `json:"status"`
	OrderedDeviceIDs []string                    `json:"ordered_device_ids"`
	CreatedAt        any                         `json:"created_at"`
	UpdatedAt        any                         `json:"updated_at"`
	ReadinessSummary *workspace.ReadinessSummary `json:"readiness_summary"`
}

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
	readinessSummary, err := h.svc.Readiness(c.Request.Context())
	if err != nil {
		readinessSummary = studioV2WorkspaceReadinessUnavailableSummary()
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": studioV2WorkspaceBootstrapResponse{
			ID:               record.ID,
			Kind:             record.Kind,
			Status:           record.Status,
			OrderedDeviceIDs: append([]string{}, record.OrderedDeviceIDs...),
			CreatedAt:        record.CreatedAt,
			UpdatedAt:        record.UpdatedAt,
			ReadinessSummary: readinessSummary,
		},
	})
}

func studioV2WorkspaceReadinessUnavailableSummary() *workspace.ReadinessSummary {
	return &workspace.ReadinessSummary{
		Ready:         false,
		BlockingCount: 1,
		Issues: []workspace.ReadinessIssue{
			{
				Code:     "readiness-unavailable",
				Severity: workspace.ReadinessSeverityBlocking,
				Step:     workspace.ReadinessStep1,
				Scope:    "workspace",
				Message:  "workspace readiness is unavailable",
			},
		},
	}
}
