package handlers

import (
	"net/http"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"

	"github.com/gin-gonic/gin"
)

type RuntimeHandler struct {
	deviceSvc  *device.Service
	pointSvc   *point.Service
	groupSvc   *pollinggroup.Service
	scheduler  *collector.Scheduler
	runtimeSvc *datalinkruntime.Service
}

type runtimeStatusResponse struct {
	Running       bool                       `json:"running"`
	UptimeSeconds int64                      `json:"uptime_seconds"`
	Metrics       datalinkruntime.Stats      `json:"metrics"`
	Collectors    []runtimeCollectorResponse `json:"collectors"`
}

type runtimeCollectorResponse struct {
	DeviceID      string     `json:"device_id"`
	DeviceName    string     `json:"device_name"`
	Protocol      string     `json:"protocol"`
	Status        string     `json:"status"`
	PointsTotal   int        `json:"points_total"`
	PointsHealthy int        `json:"points_healthy"`
	PointsStale   int        `json:"points_stale"`
	PointsError   int        `json:"points_error"`
	LastReadAt    *time.Time `json:"last_read_at"`
	LastError     *string    `json:"last_error"`
	BreakerState  string     `json:"breaker_state"`
}

func NewRuntimeHandler(
	deviceSvc *device.Service,
	pointSvc *point.Service,
	groupSvc *pollinggroup.Service,
	scheduler *collector.Scheduler,
	runtimeSvc *datalinkruntime.Service,
) *RuntimeHandler {
	return &RuntimeHandler{
		deviceSvc:  deviceSvc,
		pointSvc:   pointSvc,
		groupSvc:   groupSvc,
		scheduler:  scheduler,
		runtimeSvc: runtimeSvc,
	}
}

func (h *RuntimeHandler) Status(c *gin.Context) {
	if h.runtimeSvc != nil {
		snapshot, err := h.runtimeSvc.RuntimeStatusSnapshot(c.Request.Context(), c.Query("device_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   gin.H{"message": "failed to build runtime status snapshot: " + err.Error()},
			})
			return
		}

		response := runtimeStatusResponse{
			Running:       snapshot.Running,
			UptimeSeconds: snapshot.UptimeSeconds,
			Metrics:       snapshot.Metrics,
			Collectors:    make([]runtimeCollectorResponse, 0, len(snapshot.Collectors)),
		}
		for _, collector := range snapshot.Collectors {
			response.Collectors = append(response.Collectors, runtimeCollectorResponse{
				DeviceID:      collector.DeviceID,
				DeviceName:    collector.DeviceName,
				Protocol:      collector.Protocol,
				Status:        collector.Status,
				PointsTotal:   collector.PointsTotal,
				PointsHealthy: collector.PointsHealthy,
				PointsStale:   collector.PointsStale,
				PointsError:   collector.PointsError,
				LastReadAt:    collector.LastReadAt,
				LastError:     collector.LastError,
				BreakerState:  collector.BreakerState,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    response,
		})
		return
	}

	if h.deviceSvc == nil || h.pointSvc == nil || h.groupSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "runtime services unavailable"},
		})
		return
	}

	ctx := c.Request.Context()

	devices, err := h.deviceSvc.List(ctx, device.ListFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": "failed to load devices: " + err.Error()},
		})
		return
	}

	points, err := h.pointSvc.List(ctx, point.ListFilter{Limit: 100000})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": "failed to load points: " + err.Error()},
		})
		return
	}

	groups, err := h.groupSvc.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": "failed to load polling groups: " + err.Error()},
		})
		return
	}

	deviceIDFilter := c.Query("device_id")
	groupIntervalByID := make(map[string]time.Duration, len(groups))
	for _, group := range groups {
		groupIntervalByID[group.ID] = time.Duration(group.IntervalMs) * time.Millisecond
	}

	pointsByDeviceID := make(map[string][]*point.PointDetail, len(devices))
	for _, pt := range points {
		pointsByDeviceID[pt.DeviceID] = append(pointsByDeviceID[pt.DeviceID], &point.PointDetail{Point: pt})
	}

	now := time.Now()
	running := h.scheduler != nil && h.scheduler.IsRunning()
	uptimeSeconds := int64(0)
	if h.runtimeSvc != nil {
		running = h.runtimeSvc.IsRunning()
		uptimeSeconds = h.runtimeSvc.UptimeSeconds()
	}
	response := runtimeStatusResponse{
		Running:       running,
		UptimeSeconds: uptimeSeconds,
		Metrics:       datalinkruntime.Stats{},
		Collectors:    make([]runtimeCollectorResponse, 0, len(devices)),
	}

	for _, dev := range devices {
		if deviceIDFilter != "" && dev.ID != deviceIDFilter {
			continue
		}

		collectorResp := runtimeCollectorResponse{
			DeviceID:     dev.ID,
			DeviceName:   dev.Name,
			Protocol:     string(dev.Protocol),
			Status:       "idle",
			BreakerState: "closed",
		}

		if h.scheduler != nil {
			if state, exists := h.scheduler.GetDeviceBreakerState(dev.ID); exists {
				collectorResp.BreakerState = string(state)
			}
		}

		devicePoints := pointsByDeviceID[dev.ID]
		collectorResp.PointsTotal = len(devicePoints)

		var lastReadAt *time.Time
		var lastError *string
		for _, pointDetail := range devicePoints {
			if pointDetail.LastReadAt == nil {
				collectorResp.PointsStale++
				continue
			}

			if lastReadAt == nil || pointDetail.LastReadAt.After(*lastReadAt) {
				lastReadAt = pointDetail.LastReadAt
			}

			interval := time.Second
			if pointDetail.PollingGroupID != nil {
				if groupInterval, exists := groupIntervalByID[*pointDetail.PollingGroupID]; exists && groupInterval > 0 {
					interval = groupInterval
				}
			}

			isStale := now.Sub(*pointDetail.LastReadAt) > interval
			if pointDetail.LastError != "" {
				collectorResp.PointsError++
				lastError = &pointDetail.LastError
				continue
			}
			if isStale {
				collectorResp.PointsStale++
				continue
			}
			collectorResp.PointsHealthy++
		}

		collectorResp.LastReadAt = lastReadAt
		collectorResp.LastError = lastError

		switch {
		case collectorResp.PointsError > 0 || collectorResp.BreakerState == "open":
			collectorResp.Status = "error"
		case collectorResp.PointsStale > 0:
			collectorResp.Status = "warning"
		case response.Running:
			collectorResp.Status = "running"
		default:
			collectorResp.Status = "idle"
		}

		response.Collectors = append(response.Collectors, collectorResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
