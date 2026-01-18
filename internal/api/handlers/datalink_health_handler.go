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
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// Check 健康檢查
// GET /datalink/health
func (h *DatalinkHealthHandler) Check(c *gin.Context) {
	status := HealthStatus{
		Status:    "healthy",
		Service:   "datalink",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": status})
}
