package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type studioV2WorkspaceAuditHistoryResponse struct {
	Entries []audit.Entry `json:"entries"`
}

// StudioV2WorkspaceAuditHandler exposes recent workspace audit history.
type StudioV2WorkspaceAuditHandler struct {
	workspaceSvc *workspace.Service
	auditSvc     *audit.Service
}

// NewStudioV2WorkspaceAuditHandler creates a Studio V2 workspace audit handler.
func NewStudioV2WorkspaceAuditHandler(workspaceSvc *workspace.Service, auditSvc *audit.Service) *StudioV2WorkspaceAuditHandler {
	return &StudioV2WorkspaceAuditHandler{
		workspaceSvc: workspaceSvc,
		auditSvc:     auditSvc,
	}
}

// List returns recent audit history for the current Studio V2 workspace.
func (h *StudioV2WorkspaceAuditHandler) List(c *gin.Context) {
	if h == nil || h.workspaceSvc == nil || h.auditSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseMessageKey: "workspace audit history services unavailable",
			},
		})
		return
	}

	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	limit, ok := parseStudioV2WorkspaceAuditLimit(c)
	if !ok {
		return
	}
	entries, err := h.auditSvc.List(c.Request.Context(), audit.ListFilter{
		WorkspaceID: record.ID,
		EventType:   audit.EventType(strings.TrimSpace(c.Query("event_type"))),
		Limit:       limit,
	})
	if err != nil {
		log.Printf("list workspace audit history failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey: gin.H{
				apiResponseMessageKey: "failed to load workspace audit history",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey: studioV2WorkspaceAuditHistoryResponse{
			Entries: entries,
		},
	})
}

func parseStudioV2WorkspaceAuditLimit(c *gin.Context) (int, bool) {
	raw := strings.TrimSpace(c.Query("limit"))
	if raw == "" {
		return 0, true
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 0 {
		renderStudioV2WorkspaceValidationError(c, errors.New("limit must be a non-negative integer"))
		return 0, false
	}
	return limit, true
}
