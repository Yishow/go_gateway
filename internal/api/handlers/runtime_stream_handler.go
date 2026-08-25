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

type runtimeStreamStateResponse struct {
	DeviceID    string                            `json:"device_id"`
	StreamState datalinkruntime.RuntimeTruthState `json:"stream_state"`
}

func NewRuntimeStreamHandler(source runtimeValueStreamSource) *RuntimeStreamHandler {
	return &RuntimeStreamHandler{
		source:            source,
		heartbeatInterval: 15 * time.Second,
	}
}

// Stream emits runtime value, status, and heartbeat events for a device.
// @Summary Stream runtime updates
// @Description Opens an SSE stream for runtime updates. Failure responses use the typed safe-error envelope.
// @Tags datalink
// @Produce text/event-stream
// @Param device_id query string true "Device ID"
// @Param point_ids query string false "Comma-separated point IDs"
// @Success 200 {string} string "SSE stream"
// @Failure 400 {object} APIErrorResponse "Runtime device is missing"
// @Failure 503 {object} APIErrorResponse "Runtime stream unavailable"
// @Router /datalink/runtime/stream [get]
func (h *RuntimeStreamHandler) Stream(c *gin.Context) {
	if h.source == nil {
		renderTypedAPIErrorWithData(c, http.StatusServiceUnavailable, ErrCodeRuntimeStreamUnavailable, true, runtimeStreamUnavailableResponse("", "runtime stream unavailable"))
		return
	}

	deviceID := strings.TrimSpace(c.Query("device_id"))
	if deviceID == "" {
		renderTypedAPIError(c, http.StatusBadRequest, ErrCodeRuntimeDeviceNotFound, false)
		return
	}

	pointIDs := parseRuntimePointIDs(c.Query("point_ids"))

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		renderTypedAPIError(c, http.StatusServiceUnavailable, ErrCodeRuntimeStreamUnavailable, true)
		return
	}

	valueStream, unsubscribeValues := h.source.SubscribeValueEvents(deviceID, pointIDs)
	statusStream, unsubscribeStatus := h.source.SubscribeStatusEvents(deviceID)
	if valueStream == nil || statusStream == nil {
		if unsubscribeValues != nil {
			unsubscribeValues()
		}
		if unsubscribeStatus != nil {
			unsubscribeStatus()
		}
		renderTypedAPIErrorWithData(c, http.StatusServiceUnavailable, ErrCodeRuntimeStreamUnavailable, true, runtimeStreamUnavailableResponse(deviceID, "runtime stream unavailable"))
		return
	}
	defer unsubscribeValues()
	defer unsubscribeStatus()

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	heartbeatInterval := h.heartbeatInterval
	if heartbeatInterval <= 0 {
		heartbeatInterval = 15 * time.Second
	}
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	if err := writeRuntimeSSE(c.Writer, flusher, "stream_state", datalinkruntime.RuntimeStreamStateEvent{
		DeviceID:    deviceID,
		StreamState: datalinkruntime.RuntimeReadyTruthState(),
		Timestamp:   time.Now().UTC(),
		Code:        "runtime_stream_ready",
		RequestID:   getOrGenerateRequestID(c),
	}); err != nil {
		return
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case evt, ok := <-valueStream:
			if !ok {
				_ = writeRuntimeSSE(c.Writer, flusher, "stream_state", datalinkruntime.RuntimeStreamStateEvent{
					DeviceID:    deviceID,
					StreamState: datalinkruntime.RuntimeUnavailableTruthState("runtime value stream closed"),
					Timestamp:   time.Now().UTC(),
					Code:        ErrCodeRuntimeStreamUnavailable,
					Message:     typedAPIErrorMessage(ErrCodeRuntimeStreamUnavailable),
					Retryable:   true,
					Action:      "retry the runtime stream",
					RequestID:   getOrGenerateRequestID(c),
				})
				return
			}
			if err := writeRuntimeSSE(c.Writer, flusher, "value", evt); err != nil {
				return
			}
		case evt, ok := <-statusStream:
			if !ok {
				_ = writeRuntimeSSE(c.Writer, flusher, "stream_state", datalinkruntime.RuntimeStreamStateEvent{
					DeviceID:    deviceID,
					StreamState: datalinkruntime.RuntimeUnavailableTruthState("runtime status stream closed"),
					Timestamp:   time.Now().UTC(),
					Code:        ErrCodeRuntimeStreamUnavailable,
					Message:     typedAPIErrorMessage(ErrCodeRuntimeStreamUnavailable),
					Retryable:   true,
					Action:      "retry the runtime stream",
					RequestID:   getOrGenerateRequestID(c),
				})
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

func runtimeStreamUnavailableResponse(deviceID string, reason string) runtimeStreamStateResponse {
	return runtimeStreamStateResponse{
		DeviceID:    strings.TrimSpace(deviceID),
		StreamState: datalinkruntime.RuntimeUnavailableTruthState(reason),
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
