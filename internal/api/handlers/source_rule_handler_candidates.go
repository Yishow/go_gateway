package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) Candidates(c *gin.Context) {
	view, err := h.svc.GetCandidateView(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": view})
}

func (h *SourceRuleHandler) RecomputeCandidates(c *gin.Context) {
	view, err := h.svc.RecomputeCandidateView(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": view})
}
