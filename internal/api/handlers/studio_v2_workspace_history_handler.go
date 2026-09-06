package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go-gateway/internal/datalink/history"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

// StudioV2WorkspaceHistoryHandler 處理歷史數據查詢與 CSV 匯出請求。
type StudioV2WorkspaceHistoryHandler struct {
	workspaceSvc *workspace.Service
	historySvc   *history.Service
}

// NewStudioV2WorkspaceHistoryHandler 建立歷史處理器。
func NewStudioV2WorkspaceHistoryHandler(
	workspaceSvc *workspace.Service,
	historySvc *history.Service,
) *StudioV2WorkspaceHistoryHandler {
	return &StudioV2WorkspaceHistoryHandler{
		workspaceSvc: workspaceSvc,
		historySvc:   historySvc,
	}
}

// Query 查詢遙測歷史資料與統計曲線。
func (h *StudioV2WorkspaceHistoryHandler) Query(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var req history.HistoryQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_QUERY_PAYLOAD",
				"message": err.Error(),
			},
		})
		return
	}

	req.WorkspaceID = record.ID
	report, err := h.historySvc.QueryHistory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HISTORY_QUERY_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
	})
}

// Export 匯出歷史資料為安全 CSV 檔案。
func (h *StudioV2WorkspaceHistoryHandler) Export(c *gin.Context) {
	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	var req history.HistoryQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_EXPORT_PAYLOAD",
				"message": err.Error(),
			},
		})
		return
	}

	req.WorkspaceID = record.ID
	report, err := h.historySvc.QueryHistory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HISTORY_EXPORT_FAILED",
				"message": err.Error(),
			},
		})
		return
	}

	filename := fmt.Sprintf("telemetry_history_%s_%s.csv", req.PlanID, time.Now().UTC().Format("20060102_150405"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	if err := history.ExportCSV(c.Writer, report.Points); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
