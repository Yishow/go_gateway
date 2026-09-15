package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DatalinkHealthHandler Datalink 健康檢查 API Handler
type DatalinkHealthHandler struct{}

// NewDatalinkHealthHandler 建立新的健康檢查 Handler
func NewDatalinkHealthHandler() *DatalinkHealthHandler {
	return &DatalinkHealthHandler{}
}

// HealthStatus 健康狀態
type HealthStatus struct {
	Status    string    `json:"status" example:"healthy"`
	Service   string    `json:"service" example:"datalink"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
	Version   string    `json:"version,omitempty" example:"1.0.0"`
}

// Check 健康檢查
// @Summary 健康檢查
// @Description 檢查 Datalink 服務的健康狀態
// @Tags datalink
// @Accept json
// @Produce json
// @Success 200 {object} HealthStatus "健康狀態"
// @Router /datalink/health [get]
func (h *DatalinkHealthHandler) Check(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Service:   "datalink",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: status})
}
