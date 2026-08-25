package handlers

import (
	"errors"
	"net/http"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

func (h *SourceRuleHandler) ApplyDatabaseOutputs(c *gin.Context) {
	var req sourcerule.ApplyOutputCandidatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "validation",
				"message": err.Error(),
			},
		})
		return
	}

	response, err := h.svc.ApplyDatabaseOutputCandidates(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		writeOutputApplyError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

func (h *SourceRuleHandler) ApplyLocalModbusOutputs(c *gin.Context) {
	var req sourcerule.ApplyOutputCandidatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "validation",
				"message": err.Error(),
			},
		})
		return
	}

	response, err := h.svc.ApplyLocalModbusOutputCandidates(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		writeOutputApplyError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": response})
}

func writeOutputApplyError(c *gin.Context, err error) {
	var shareErr *modbusshare.Error
	if errors.As(err, &shareErr) {
		statusCode := http.StatusUnprocessableEntity
		switch shareErr.Code {
		case modbusshare.ErrCodeWorkspaceScope:
			statusCode = http.StatusForbidden
		case modbusshare.ErrCodeRevisionConflict:
			statusCode = http.StatusConflict
		case modbusshare.ErrCodeHydrationRequired:
			statusCode = http.StatusServiceUnavailable
		}
		renderModbusShareAPIError(c, statusCode, shareErr, shareErr.Code, shareErr.Retryable)
		return
	}
	statusCode := http.StatusInternalServerError
	errorCode := "internal"
	switch {
	case errors.Is(err, sourcerule.ErrSourceRuleNotFound):
		statusCode = http.StatusNotFound
		errorCode = "not_found"
	case errors.Is(err, sourcerule.ErrOutputApplyRevisionConflict):
		statusCode = http.StatusConflict
		errorCode = "revision_mismatch"
	case errors.Is(err, sourcerule.ErrInvalidOutputApplyRequest):
		statusCode = http.StatusBadRequest
		errorCode = "validation"
	}

	c.JSON(statusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    errorCode,
			"message": err.Error(),
		},
	})
}
