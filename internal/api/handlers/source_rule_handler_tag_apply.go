package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) ApplyTags(c *gin.Context) {
	var req sourcerule.ApplyTagCandidatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	response, err := h.svc.ApplyTagCandidates(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch {
		case errors.Is(err, sourcerule.ErrSourceRuleNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, sourcerule.ErrTagApplyRevisionConflict):
			statusCode = http.StatusConflict
		case errors.Is(err, sourcerule.ErrInvalidTagApplyRequest):
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}
