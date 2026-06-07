package handlers

import (
	"net/http"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type RuntimeHandler struct {
	deviceSvc    *device.Service
	pointSvc     *point.Service
	groupSvc     *pollinggroup.Service
	scheduler    *collector.Scheduler
	runtimeSvc   *datalinkruntime.Service
	workspaceSvc *workspace.Service
}

type runtimeStatusResponse struct {
	Running       bool                       `json:"running"`
	UptimeSeconds int64                      `json:"uptime_seconds"`
	Metrics       datalinkruntime.Stats      `json:"metrics"`
	Collectors    []runtimeCollectorResponse `json:"collectors"`
}

type runtimeCollectorResponse struct {
	DeviceID                   string     `json:"device_id"`
	DeviceName                 string     `json:"device_name"`
	Protocol                   string     `json:"protocol"`
	Status                     string     `json:"status"`
	AvailabilityStatus         string     `json:"availability_status"`
	AvailabilityReason         *string    `json:"availability_reason,omitempty"`
	Running                    bool       `json:"running"`
	PointsTotal                int        `json:"points_total"`
	PointsHealthy              int        `json:"points_healthy"`
	PointsStale                int        `json:"points_stale"`
	PointsError                int        `json:"points_error"`
	LastReadAt                 *time.Time `json:"last_read_at"`
	LastError                  *string    `json:"last_error"`
	BreakerState               string     `json:"breaker_state"`
	ProjectionAlignment        string     `json:"projection_alignment,omitempty"`
	RuntimeProjectionVersion   string     `json:"runtime_projection_version,omitempty"`
	WorkspaceProjectionVersion string     `json:"workspace_projection_version,omitempty"`
	ProjectionMessage          string     `json:"projection_message,omitempty"`
}

type runtimeWorkspaceContextResponse struct {
	WorkspaceID     string                               `json:"workspace_id"`
	Devices         []runtimeWorkspaceContextDeviceEntry `json:"devices"`
	DefaultDeviceID *string                              `json:"default_device_id"`
}

type runtimeWorkspaceContextDeviceEntry struct {
	DeviceID                   string  `json:"device_id"`
	Name                       string  `json:"name"`
	Protocol                   string  `json:"protocol"`
	Running                    bool    `json:"running"`
	AvailabilityStatus         string  `json:"availability_status"`
	AvailabilityReason         *string `json:"availability_reason,omitempty"`
	ProjectionAlignment        string  `json:"projection_alignment,omitempty"`
	RuntimeProjectionVersion   string  `json:"runtime_projection_version,omitempty"`
	WorkspaceProjectionVersion string  `json:"workspace_projection_version,omitempty"`
	ProjectionMessage          string  `json:"projection_message,omitempty"`
}

func NewRuntimeHandler(
	deviceSvc *device.Service,
	pointSvc *point.Service,
	groupSvc *pollinggroup.Service,
	scheduler *collector.Scheduler,
	runtimeSvc *datalinkruntime.Service,
	workspaceSvc *workspace.Service,
) *RuntimeHandler {
	return &RuntimeHandler{
		deviceSvc:    deviceSvc,
		pointSvc:     pointSvc,
		groupSvc:     groupSvc,
		scheduler:    scheduler,
		runtimeSvc:   runtimeSvc,
		workspaceSvc: workspaceSvc,
	}
}

func (h *RuntimeHandler) WorkspaceContext(c *gin.Context) {
	if h.workspaceSvc == nil || h.deviceSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error": gin.H{
				"message": "runtime workspace services unavailable",
			},
		})
		return
	}

	record, err := h.workspaceSvc.GetOrCreate(c.Request.Context())
	if err != nil {
		renderStudioV2WorkspaceBootstrapError(c)
		return
	}

	statusByDeviceID := make(map[string]datalinkruntime.DeviceRuntimeStatus)
	if h.runtimeSvc != nil {
		snapshot, err := h.runtimeSvc.RuntimeStatusSnapshot(c.Request.Context(), "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"message": "failed to build workspace runtime context: " + err.Error(),
				},
			})
			return
		}
		for _, collector := range snapshot.Collectors {
			statusByDeviceID[collector.DeviceID] = collector
		}
	}

	response := runtimeWorkspaceContextResponse{
		WorkspaceID: record.ID,
		Devices:     make([]runtimeWorkspaceContextDeviceEntry, 0, len(record.OrderedDeviceIDs)),
	}
	for _, deviceID := range record.OrderedDeviceIDs {
		savedDevice, err := h.deviceSvc.GetByID(c.Request.Context(), deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"message": "failed to load workspace runtime device: " + err.Error(),
				},
			})
			return
		}

		deviceResponse := mapRuntimeWorkspaceContextDevice(savedDevice, statusByDeviceID[deviceID])
		response.Devices = append(response.Devices, deviceResponse)
		if response.DefaultDeviceID == nil && deviceResponse.AvailabilityStatus != device.AvailabilityStatusUnavailable {
			id := deviceResponse.DeviceID
			response.DefaultDeviceID = &id
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
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
			response.Collectors = append(response.Collectors, mapRuntimeCollectorResponse(collector))
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
			DeviceID:           dev.ID,
			DeviceName:         dev.Name,
			Protocol:           string(dev.Protocol),
			Status:             "idle",
			AvailabilityStatus: device.AvailabilityStatusAvailable,
			Running:            false,
			BreakerState:       "closed",
		}
		if availabilityStatus, availabilityReason := device.AvailabilityOf(dev); availabilityStatus == device.AvailabilityStatusUnavailable {
			collectorResp.AvailabilityStatus = availabilityStatus
			collectorResp.AvailabilityReason = &availabilityReason
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
		collectorResp.Running = collectorResp.Status == "running" && collectorResp.AvailabilityStatus == device.AvailabilityStatusAvailable
		if collectorResp.AvailabilityStatus == device.AvailabilityStatusUnavailable {
			collectorResp.Running = false
			collectorResp.Status = "idle"
		}

		response.Collectors = append(response.Collectors, collectorResp)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func mapRuntimeWorkspaceContextDevice(savedDevice *schema.Device, runtimeStatus datalinkruntime.DeviceRuntimeStatus) runtimeWorkspaceContextDeviceEntry {
	response := runtimeWorkspaceContextDeviceEntry{
		DeviceID:           savedDevice.ID,
		Name:               savedDevice.Name,
		Protocol:           string(savedDevice.Protocol),
		Running:            isWorkspaceDeviceRunning(savedDevice),
		AvailabilityStatus: availabilityStatusOf(savedDevice),
	}
	if reason := availabilityReasonOf(savedDevice); reason != "" {
		response.AvailabilityReason = &reason
	}

	if runtimeStatus.DeviceID == "" {
		return response
	}

	response.Running = runtimeStatus.Running
	response.AvailabilityStatus = runtimeStatus.AvailabilityStatus
	response.AvailabilityReason = runtimeStatus.AvailabilityReason
	response.ProjectionAlignment = runtimeStatus.ProjectionAlignment
	response.RuntimeProjectionVersion = runtimeStatus.RuntimeProjectionVersion
	response.WorkspaceProjectionVersion = runtimeStatus.WorkspaceProjectionVersion
	response.ProjectionMessage = runtimeStatus.ProjectionMessage
	return response
}

func mapRuntimeCollectorResponse(collector datalinkruntime.DeviceRuntimeStatus) runtimeCollectorResponse {
	return runtimeCollectorResponse{
		DeviceID:                   collector.DeviceID,
		DeviceName:                 collector.DeviceName,
		Protocol:                   collector.Protocol,
		Status:                     collector.Status,
		AvailabilityStatus:         collector.AvailabilityStatus,
		AvailabilityReason:         collector.AvailabilityReason,
		Running:                    collector.Running,
		PointsTotal:                collector.PointsTotal,
		PointsHealthy:              collector.PointsHealthy,
		PointsStale:                collector.PointsStale,
		PointsError:                collector.PointsError,
		LastReadAt:                 collector.LastReadAt,
		LastError:                  collector.LastError,
		BreakerState:               collector.BreakerState,
		ProjectionAlignment:        collector.ProjectionAlignment,
		RuntimeProjectionVersion:   collector.RuntimeProjectionVersion,
		WorkspaceProjectionVersion: collector.WorkspaceProjectionVersion,
		ProjectionMessage:          collector.ProjectionMessage,
	}
}
