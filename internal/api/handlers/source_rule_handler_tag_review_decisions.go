package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) ListTagReviewDecisions(c *gin.Context) {
	scope := candidateScopeFromQuery(c)
	if !h.requireCandidateScope(c, scope) {
		return
	}
	var decisions []*schema.SourceRuleTagReviewDecision
	var err error
	if h.svc.CandidateScopeConfigured() {
		decisions, err = h.svc.ListTagReviewDecisionsInWorkspace(c.Request.Context(), c.Param("id"), scope)
	} else {
		decisions, err = h.svc.ListTagReviewDecisions(c.Request.Context(), c.Param("id"))
	}
	if err != nil {
		if h.svc.CandidateScopeConfigured() {
			writeCandidateScopeError(c, err)
			return
		}
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: decisions})
}

func (h *SourceRuleHandler) UpsertTagReviewDecision(c *gin.Context) {
	var req sourcerule.UpsertTagReviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	decision, err := h.svc.UpsertTagReviewDecision(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		if h.svc.CandidateScopeConfigured() {
			writeCandidateScopeError(c, err)
			return
		}
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: decision})
}
