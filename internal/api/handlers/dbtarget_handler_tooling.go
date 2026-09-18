package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/gin-gonic/gin"
)

func (h *DatabaseTargetHandler) GenerateSchema(c *gin.Context) {
	var req dbtarget.SchemaGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseCodeKey: apiValidationErrorCode, apiResponseMessageKey: err.Error()}})
		return
	}

	if !requireSchemaConfirmation(c, req.DryRun) {
		return
	}

	result, err := h.connectorSvc.GenerateSchema(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		statusCode, errorCode := dbTargetErrorStatusCode(err)
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseCodeKey: errorCode, apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func (h *DatabaseTargetHandler) DryRunMappings(c *gin.Context) {
	var req dbtarget.MappingDryRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseCodeKey: apiValidationErrorCode, apiResponseMessageKey: err.Error()}})
		return
	}

	result, err := h.mappingSvc.DryRun(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		statusCode, errorCode := dbTargetErrorStatusCode(err)
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseCodeKey: errorCode, apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func (h *DatabaseTargetHandler) ListWriteHistory(c *gin.Context) {
	limit := parseHistoryLimit(c.Query("limit"))
	records, err := h.connectorSvc.ListWriteHistory(c.Request.Context(), c.Param("id"), limit)
	if err != nil {
		statusCode, errorCode := dbTargetErrorStatusCode(err)
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseCodeKey: errorCode, apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: gin.H{
		"connector_id": c.Param("id"),
		"records":      records,
	}})
}

func parseHistoryLimit(value string) int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 20
	}
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 20
	}
	if parsed <= 0 {
		return 20
	}
	if parsed > 200 {
		return 200
	}
	return parsed
}

func dbTargetErrorStatusCode(err error) (statusCode int, errorCode string) {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "not found"), strings.Contains(message, "不存在"):
		return http.StatusNotFound, "not_found"
	case strings.Contains(message, "連接"), strings.Contains(message, "unreachable"), strings.Contains(message, "connector"):
		return http.StatusUnprocessableEntity, "connector_unavailable"
	case strings.Contains(message, "schema"), strings.Contains(message, "table"), strings.Contains(message, "column"), strings.Contains(message, "欄位"):
		return http.StatusUnprocessableEntity, "schema_missing"
	default:
		return http.StatusInternalServerError, apiInternalErrorCode
	}
}
