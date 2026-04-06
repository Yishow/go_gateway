package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) ListTagReviewDecisions(c *gin.Context) {
	decisions, err := h.svc.ListTagReviewDecisions(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": decisions})
}

func (h *SourceRuleHandler) UpsertTagReviewDecision(c *gin.Context) {
	var req sourcerule.UpsertTagReviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	decision, err := h.svc.UpsertTagReviewDecision(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": decision})
}
