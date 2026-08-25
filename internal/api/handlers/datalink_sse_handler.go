package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

// DatalinkSSEHandler Datalink SSE Handler
type DatalinkSSEHandler struct {
	mappingSvc *mapping.Service
	runtime    interface {
		SubscribeValueEvents(string, []string) (<-chan datalinkruntime.ValueEvent, func())
	}
	scopeValidator func(context.Context, string, *schema.Mapping) (string, error)
}

// NewDatalinkSSEHandler connects preview consumers to persisted mappings and
// the runtime value stream; it never fabricates values.
func NewDatalinkSSEHandler(mappingSvc *mapping.Service, runtimeSource interface {
	SubscribeValueEvents(string, []string) (<-chan datalinkruntime.ValueEvent, func())
}) *DatalinkSSEHandler {
	return &DatalinkSSEHandler{mappingSvc: mappingSvc, runtime: runtimeSource}
}

// WithPreviewScopeValidator requires durable workspace ownership proof before
// opening a runtime preview stream.
func (h *DatalinkSSEHandler) WithPreviewScopeValidator(validator func(context.Context, string, *schema.Mapping) (string, error)) *DatalinkSSEHandler {
	h.scopeValidator = validator
	return h
}

// PreviewEvent SSE 預覽事件
type PreviewEvent struct {
	Type       string                 `json:"type"` // "preview" | "error" | "heartbeat"
	Code       string                 `json:"code,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Retryable  bool                   `json:"retryable"`
	Action     string                 `json:"action,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	MappingID  string                 `json:"mapping_id,omitempty"`
	RawValue   interface{}            `json:"raw_value,omitempty"`
	FinalValue interface{}            `json:"final_value,omitempty"`
	Steps      []mapping.StepResult   `json:"steps,omitempty"`
	Quality    int                    `json:"quality,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	Error      *TypedAPIErrorEnvelope `json:"error,omitempty"`
}

// PreviewStream 即時預覽 SSE 串流
// @Summary Stream a persisted mapping preview
// @Description Opens a text/event-stream backed by the runtime value source for the persisted mapping. Failure responses use the typed safe-error envelope before the stream starts; a runtime stream failure is emitted as an SSE error event.
// @Tags datalink
// @Produce text/event-stream
// @Param mapping_id query string true "Persisted mapping identifier"
// @Param workspace_id query string true "Workspace scope"
// @Success 200 {string} string "SSE stream of connected, preview, heartbeat, and error events"
// @Failure 422 {object} APIErrorResponse "Mapping identifier is missing or invalid"
// @Failure 503 {object} APIErrorResponse "Preview mapping or runtime stream is unavailable"
// @Router /datalink/preview/stream [get]
func (h *DatalinkSSEHandler) PreviewStream(c *gin.Context) {
	mappingID := c.Query("mapping_id")
	if mappingID == "" {
		renderSafeError(c, http.StatusUnprocessableEntity, ErrCodePreviewInvalidRequest, false)
		return
	}
	if h.mappingSvc == nil || h.runtime == nil {
		renderSafeError(c, http.StatusServiceUnavailable, ErrCodePreviewUnavailable, true)
		return
	}
	persistedMapping, err := h.mappingSvc.GetByID(c.Request.Context(), mappingID)
	if err != nil {
		if errors.Is(err, mapping.ErrMappingNotFound) {
			renderSafeError(c, http.StatusUnprocessableEntity, ErrCodePreviewInvalidRequest, false)
			return
		}
		renderSafeError(c, http.StatusServiceUnavailable, ErrCodePreviewUnavailable, true)
		return
	}
	if persistedMapping == nil || persistedMapping.PointID == "" {
		renderSafeError(c, http.StatusUnprocessableEntity, ErrCodePreviewInvalidRequest, false)
		return
	}
	workspaceID := strings.TrimSpace(c.Query("workspace_id"))
	if workspaceID == "" || h.scopeValidator == nil {
		renderSafeError(c, http.StatusUnprocessableEntity, modbusshare.ErrCodeWorkspaceScope, true)
		return
	}
	deviceID, scopeErr := h.scopeValidator(c.Request.Context(), workspaceID, persistedMapping)
	if scopeErr != nil || strings.TrimSpace(deviceID) == "" {
		if scopeErr == nil {
			scopeErr = modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "preview device ownership could not be verified", true)
		}
		renderModbusShareAPIError(c, http.StatusForbidden, scopeErr, modbusshare.ErrCodeWorkspaceScope, false)
		return
	}
	valueStream, unsubscribe := h.runtime.SubscribeValueEvents(deviceID, []string{persistedMapping.PointID})
	if valueStream == nil {
		renderSafeError(c, http.StatusServiceUnavailable, ErrCodePreviewUnavailable, true)
		return
	}
	defer unsubscribe()

	// 設定 SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 取得 flusher
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		renderSafeError(c, http.StatusServiceUnavailable, ErrCodePreviewUnavailable, true)
		return
	}

	requestID := getOrGenerateRequestID(c)
	// 發送初始連線確認
	h.sendEvent(c.Writer, flusher, PreviewEvent{
		Type:      "connected",
		Code:      "preview_connected",
		MappingID: mappingID,
		RequestID: requestID,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	clientGone := c.Request.Context().Done()
	heartbeatTicker := time.NewTicker(15 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-clientGone:
			// 客戶端斷線
			return

		case <-heartbeatTicker.C:
			// 發送心跳
			h.sendEvent(c.Writer, flusher, PreviewEvent{
				Type:      "heartbeat",
				Code:      "preview_heartbeat",
				RequestID: requestID,
				Timestamp: time.Now().Format(time.RFC3339),
			})

		case value, ok := <-valueStream:
			if !ok {
				h.sendEvent(c.Writer, flusher, PreviewEvent{
					Type: apiResponseErrorKey, MappingID: mappingID,
					Error: &TypedAPIErrorEnvelope{
						Code: ErrCodePreviewStreamClosed, Message: typedAPIErrorMessage(ErrCodePreviewStreamClosed),
						Retryable: true, RequestID: requestID, Action: "retry the preview stream",
					}, Timestamp: time.Now().UTC().Format(time.RFC3339),
					Code: ErrCodePreviewStreamClosed, Message: typedAPIErrorMessage(ErrCodePreviewStreamClosed), Retryable: true, RequestID: requestID, Action: "retry the preview stream",
				})
				return
			}
			h.sendEvent(c.Writer, flusher, PreviewEvent{
				Type:       "preview",
				Code:       "preview_value",
				RequestID:  requestID,
				MappingID:  mappingID,
				RawValue:   value.RawValue,
				FinalValue: value.TransformedValue,
				Quality:    previewQuality(value.Quality),
				Timestamp:  value.Timestamp.UTC().Format(time.RFC3339),
			})
		}
	}
}

func previewQuality(quality schema.QualityFlag) int {
	if quality == schema.QualityGood {
		return 192
	}
	return 0
}

// sendEvent 發送 SSE 事件
func (h *DatalinkSSEHandler) sendEvent(w http.ResponseWriter, flusher http.Flusher, event PreviewEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()
}
