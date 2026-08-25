package handlers

import (
	"errors"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) Candidates(c *gin.Context) {
	scope := candidateScopeFromQuery(c)
	if !h.requireCandidateScope(c, scope) {
		return
	}
	var view *sourcerule.CandidateSnapshotView
	var err error
	if h.svc.CandidateScopeConfigured() {
		view, err = h.svc.GetCandidateViewAtRevision(c.Request.Context(), c.Param("id"), scope.RevisionID, &scope)
	} else {
		view, err = h.svc.GetCandidateView(c.Request.Context(), c.Param("id"))
	}
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	applyCandidateScopeToView(view, scope)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": view})
}

func (h *SourceRuleHandler) RecomputeCandidates(c *gin.Context) {
	scope := candidateScopeFromQuery(c)
	if !h.requireCandidateScope(c, scope) {
		return
	}
	var view *sourcerule.CandidateSnapshotView
	var err error
	if h.svc.CandidateScopeConfigured() {
		view, err = h.svc.RecomputeCandidateViewAtRevision(c.Request.Context(), c.Param("id"), scope.RevisionID, &scope)
	} else {
		view, err = h.svc.RecomputeCandidateView(c.Request.Context(), c.Param("id"))
	}
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	applyCandidateScopeToView(view, scope)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": view})
}

func candidateScopeFromQuery(c *gin.Context) sourcerule.CandidateScopeRequest {
	return sourcerule.CandidateScopeRequest{
		WorkspaceID:               c.Query("workspace_id"),
		ExpectedWorkspaceRevision: c.Query("expected_workspace_revision"),
		RevisionID:                c.Query("revision_id"),
	}
}

func (h *SourceRuleHandler) requireCandidateScope(c *gin.Context, scope sourcerule.CandidateScopeRequest) bool {
	if !h.svc.CandidateScopeConfigured() {
		return true
	}
	if err := h.svc.ValidateCandidateScope(c.Request.Context(), c.Param("id"), scope); err != nil {
		writeCandidateScopeError(c, err)
		return false
	}
	return true
}

func applyCandidateScopeToView(view *sourcerule.CandidateSnapshotView, scope sourcerule.CandidateScopeRequest) {
	if view == nil || strings.TrimSpace(scope.WorkspaceID) == "" {
		return
	}
	view.WorkspaceID = scope.WorkspaceID
	view.WorkspaceRevision = scope.ExpectedWorkspaceRevision
}

func writeCandidateScopeError(c *gin.Context, err error) {
	if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	var typed *modbusshare.Error
	if !errors.As(err, &typed) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	status := http.StatusUnprocessableEntity
	switch typed.Code {
	case modbusshare.ErrCodeWorkspaceScope:
		status = http.StatusForbidden
	case modbusshare.ErrCodeRevisionConflict:
		status = http.StatusConflict
	case modbusshare.ErrCodeHydrationRequired:
		status = http.StatusServiceUnavailable
	}
	renderModbusShareAPIError(c, status, typed, typed.Code, typed.Retryable)
}
