package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/collector/health"

	"github.com/gin-gonic/gin"
)

// DeviceHealthHandler 設備健康度處理器
type DeviceHealthHandler struct {
	scheduler *collector.Scheduler
}

// NewDeviceHealthHandler 建立設備健康度處理器
func NewDeviceHealthHandler(scheduler *collector.Scheduler) *DeviceHealthHandler {
	return &DeviceHealthHandler{scheduler: scheduler}
}

// BreakerStateResponse 熔斷器狀態回應
type BreakerStateResponse struct {
	DeviceID     string        `json:"device_id"`
	State        health.State  `json:"state"`
	ErrorRate    float64       `json:"error_rate"`
	IsProbing    bool          `json:"is_probing"`
	TotalRequests int          `json:"total_requests"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
}

// GetBreakerState 取得單一設備的熔斷器狀態
// GET /datalink/devices/:id/health
func (h *DeviceHealthHandler) GetBreakerState(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "排程器未初始化"},
		})
		return
	}

	deviceID := c.Param("id")
	stats, exists := h.scheduler.GetDeviceBreakerStats(deviceID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"message": "設備不存在或無健康追蹤資料"},
		})
		return
	}

	resp := BreakerStateResponse{
		DeviceID:      deviceID,
		State:         stats.State,
		ErrorRate:     stats.ErrorRate,
		IsProbing:     stats.IsProbing,
		TotalRequests: stats.TrackerStats.TotalRequests,
		SuccessCount:  stats.TrackerStats.SuccessCount,
		FailureCount:  stats.TrackerStats.FailureCount,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// ListBreakerStates 取得所有設備的熔斷器狀態
// GET /datalink/devices/health
func (h *DeviceHealthHandler) ListBreakerStates(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "排程器未初始化"},
		})
		return
	}

	states := h.scheduler.GetAllDeviceBreakerStates()
	result := make([]BreakerStateResponse, 0, len(states))

	for deviceID, state := range states {
		stats, _ := h.scheduler.GetDeviceBreakerStats(deviceID)
		result = append(result, BreakerStateResponse{
			DeviceID:      deviceID,
			State:         state,
			ErrorRate:     stats.ErrorRate,
			IsProbing:     stats.IsProbing,
			TotalRequests: stats.TrackerStats.TotalRequests,
			SuccessCount:  stats.TrackerStats.SuccessCount,
			FailureCount:  stats.TrackerStats.FailureCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// ResetBreaker 重置設備熔斷器
// POST /datalink/devices/:id/health/reset
func (h *DeviceHealthHandler) ResetBreaker(c *gin.Context) {
	if h.scheduler == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "排程器未初始化"},
		})
		return
	}

	deviceID := c.Param("id")
	if !h.scheduler.ResetDeviceBreaker(deviceID) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"message": "設備不存在"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gin.H{"message": "熔斷器已重置", "device_id": deviceID},
	})
}
