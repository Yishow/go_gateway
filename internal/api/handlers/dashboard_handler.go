package handlers

import (
	"net/http"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

// DashboardHandler 儀表板統計 API 處理器
type DashboardHandler struct {
	deviceService *device.Service
	pointService  *point.Service
	tagService    *tag.Service
}

// NewDashboardHandler 建立新的儀表板處理器
func NewDashboardHandler(
	deviceService *device.Service,
	pointService *point.Service,
	tagService *tag.Service,
) *DashboardHandler {
	return &DashboardHandler{
		deviceService: deviceService,
		pointService:  pointService,
		tagService:    tagService,
	}
}

// DashboardStats 儀表板統計數據結構
type DashboardStats struct {
	// 設備統計
	TotalDevices    int `json:"total_devices"`
	ActiveDevices   int `json:"active_devices"`
	DisabledDevices int `json:"disabled_devices"`
	ErrorDevices    int `json:"error_devices"`

	// 點位統計
	TotalPoints   int `json:"total_points"`
	EnabledPoints int `json:"enabled_points"`

	// 標籤統計
	TotalTags   int `json:"total_tags"`
	ActiveTags  int `json:"active_tags"`
	RetiredTags int `json:"retired_tags"`

	// 錯誤統計（過去24小時）
	ErrorsLast24h int `json:"errors_last_24h"`

	// 吞吐量統計（每秒操作數，基於點位輪詢頻率估算）
	EstimatedThroughput float64 `json:"estimated_throughput"`
}

// GetStats 取得儀表板統計數據
// GET /datalink/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats := DashboardStats{}

	// 獲取所有設備
	devices, err := h.deviceService.List(ctx, device.ListFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "取得設備列表失敗: " + err.Error()},
		})
		return
	}

	stats.TotalDevices = len(devices)
	for _, dev := range devices {
		switch dev.Status {
		case schema.DeviceStatusActive:
			stats.ActiveDevices++
		case schema.DeviceStatusDisabled:
			stats.DisabledDevices++
		}

		// 檢查是否有錯誤（最近測試失敗）
		if dev.LastTestSuccess != nil && !*dev.LastTestSuccess {
			stats.ErrorDevices++
		}
	}

	// 獲取所有點位
	points, err := h.pointService.List(ctx, point.ListFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "取得點位列表失敗: " + err.Error()},
		})
		return
	}

	stats.TotalPoints = len(points)
	enabled := true
	enabledPoints, err := h.pointService.List(ctx, point.ListFilter{Enabled: &enabled})
	if err == nil {
		stats.EnabledPoints = len(enabledPoints)
	}

	// 獲取所有標籤
	tags, err := h.tagService.List(ctx, tag.ListFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "取得標籤列表失敗: " + err.Error()},
		})
		return
	}

	stats.TotalTags = len(tags)
	for _, t := range tags {
		switch t.Status {
		case schema.TagStatusActive:
			stats.ActiveTags++
		case schema.TagStatusRetired:
			stats.RetiredTags++
		}
	}

	// 計算錯誤統計（過去24小時）
	// 目前基於設備測試失敗次數，未來可以擴展到專門的錯誤日誌表
	stats.ErrorsLast24h = 0
	now := time.Now()
	for _, dev := range devices {
		if dev.LastTestAt != nil {
			// 檢查是否在過去24小時內
			if now.Sub(*dev.LastTestAt) <= 24*time.Hour {
				if dev.LastTestSuccess != nil && !*dev.LastTestSuccess {
					stats.ErrorsLast24h++
				}
			}
		}
	}

	// 估算吞吐量（基於啟用點位的輪詢頻率）
	// 假設每個啟用點位每秒輪詢一次（實際頻率取決於輪詢群組配置）
	stats.EstimatedThroughput = float64(stats.EnabledPoints)

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    stats,
	})
}

// DeviceStatus 設備狀態資訊
type DeviceStatus struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Protocol    string     `json:"protocol"`
	Status      string     `json:"status"`
	LastTestAt  *time.Time `json:"last_test_at,omitempty"`
	LastSuccess *bool      `json:"last_success,omitempty"`
	LastError   string     `json:"last_error,omitempty"`
}

// GetDeviceStatuses 取得設備狀態列表
// GET /datalink/dashboard/device-statuses
func (h *DashboardHandler) GetDeviceStatuses(c *gin.Context) {
	ctx := c.Request.Context()

	devices, err := h.deviceService.List(ctx, device.ListFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "取得設備列表失敗: " + err.Error()},
		})
		return
	}

	statuses := make([]DeviceStatus, 0, len(devices))
	for _, dev := range devices {
		// 檢查 LastTestAt 是否為零值或無效時間
		var lastTestAt *time.Time
		if dev.LastTestAt != nil && !dev.LastTestAt.IsZero() {
			lastTestAt = dev.LastTestAt
		}

		statuses = append(statuses, DeviceStatus{
			ID:          dev.ID,
			Name:        dev.Name,
			Protocol:    string(dev.Protocol),
			Status:      string(dev.Status),
			LastTestAt:  lastTestAt,
			LastSuccess: dev.LastTestSuccess,
			LastError:   dev.LastTestError,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		apiResponseSuccessKey: true,
		apiResponseDataKey:    statuses,
	})
}
