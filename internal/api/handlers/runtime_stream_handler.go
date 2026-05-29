package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	datalinkruntime "go-gateway/internal/datalink/runtime"

	"github.com/gin-gonic/gin"
)

type runtimeValueStreamSource interface {
	SubscribeValueEvents(deviceID string, pointIDs []string) (<-chan datalinkruntime.ValueEvent, func())
	SubscribeStatusEvents(deviceID string) (<-chan datalinkruntime.DeviceStatusEvent, func())
}

type RuntimeStreamHandler struct {
	source            runtimeValueStreamSource
	heartbeatInterval time.Duration
}

func NewRuntimeStreamHandler(source runtimeValueStreamSource) *RuntimeStreamHandler {
	return &RuntimeStreamHandler{
		source:            source,
		heartbeatInterval: 15 * time.Second,
	}
}

func (h *RuntimeStreamHandler) Stream(c *gin.Context) {
	if h.source == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "runtime stream unavailable"},
		})
		return
	}

	deviceID := strings.TrimSpace(c.Query("device_id"))
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"message": "device_id is required"},
		})
		return
	}

	pointIDs := parseRuntimePointIDs(c.Query("point_ids"))

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": "streaming not supported"},
		})
		return
	}

	valueStream, unsubscribeValues := h.source.SubscribeValueEvents(deviceID, pointIDs)
	statusStream, unsubscribeStatus := h.source.SubscribeStatusEvents(deviceID)
	if valueStream == nil || statusStream == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error":   gin.H{"message": "runtime stream unavailable"},
		})
		return
	}
	defer unsubscribeValues()
	defer unsubscribeStatus()

	heartbeatInterval := h.heartbeatInterval
	if heartbeatInterval <= 0 {
		heartbeatInterval = 15 * time.Second
	}
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case evt, ok := <-valueStream:
			if !ok {
				return
			}
			if err := writeRuntimeSSE(c.Writer, flusher, "value", evt); err != nil {
				return
			}
		case evt, ok := <-statusStream:
			if !ok {
				return
			}
			if err := writeRuntimeSSE(c.Writer, flusher, "status", evt); err != nil {
				return
			}
		case ts := <-heartbeatTicker.C:
			payload := gin.H{"ts": ts.Format(time.RFC3339)}
			if err := writeRuntimeSSE(c.Writer, flusher, "heartbeat", payload); err != nil {
				return
			}
		}
	}
}

func parseRuntimePointIDs(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	pointIDs := make([]string, 0, len(parts))
	for _, part := range parts {
		pointID := strings.TrimSpace(part)
		if pointID == "" {
			continue
		}
		pointIDs = append(pointIDs, pointID)
	}
	return pointIDs
}

func writeRuntimeSSE(w http.ResponseWriter, flusher http.Flusher, event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}
