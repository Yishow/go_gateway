package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

// DatalinkSSEHandler Datalink SSE Handler
type DatalinkSSEHandler struct {
	mappingSvc *mapping.Service
}

// NewDatalinkSSEHandler 建立新的 Datalink SSE Handler
func NewDatalinkSSEHandler() *DatalinkSSEHandler {
	// 使用現有的 Mapping Service（實際應用中應使用依賴注入）
	repo := mapping.NewMemoryRepository()
	svc := mapping.NewService(repo)
	return &DatalinkSSEHandler{mappingSvc: svc}
}

// PreviewEvent SSE 預覽事件
type PreviewEvent struct {
	Type       string                 `json:"type"`         // "preview" | "error" | "heartbeat"
	MappingID  string                 `json:"mapping_id,omitempty"`
	RawValue   interface{}            `json:"raw_value,omitempty"`
	FinalValue interface{}            `json:"final_value,omitempty"`
	Steps      []mapping.StepResult   `json:"steps,omitempty"`
	Quality    int                    `json:"quality,omitempty"`
	Timestamp  string                 `json:"timestamp"`
	Error      string                 `json:"error,omitempty"`
}

// PreviewStream 即時預覽 SSE 串流
// GET /datalink/preview/stream?mapping_id=xxx
func (h *DatalinkSSEHandler) PreviewStream(c *gin.Context) {
	mappingID := c.Query("mapping_id")

	// 設定 SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// 取得 flusher
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": "Streaming not supported"},
		})
		return
	}

	// 發送初始連線確認
	h.sendEvent(c.Writer, flusher, PreviewEvent{
		Type:      "connected",
		MappingID: mappingID,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	// 建立客戶端斷線通知 channel
	clientGone := c.Request.Context().Done()

	// 心跳計時器
	heartbeatTicker := time.NewTicker(15 * time.Second)
	defer heartbeatTicker.Stop()

	// 模擬資料更新計時器（實際應用中應監聽實際資料變更）
	dataTicker := time.NewTicker(2 * time.Second)
	defer dataTicker.Stop()

	// 模擬計數器
	counter := 0

	for {
		select {
		case <-clientGone:
			// 客戶端斷線
			return

		case <-heartbeatTicker.C:
			// 發送心跳
			h.sendEvent(c.Writer, flusher, PreviewEvent{
				Type:      "heartbeat",
				Timestamp: time.Now().Format(time.RFC3339),
			})

		case <-dataTicker.C:
			// 發送模擬預覽資料
			counter++
			rawValue := float64(1000 + counter*10) // 模擬原始值

			// 模擬轉換管線結果
			steps := []mapping.StepResult{
				{
					StepIndex: 0,
					StepType:  string(schema.TransformScale),
					Input:     rawValue,
					Output:    rawValue * 0.1, // scale by 0.1
				},
			}

			finalValue := rawValue * 0.1

			event := PreviewEvent{
				Type:       "preview",
				MappingID:  mappingID,
				RawValue:   rawValue,
				FinalValue: finalValue,
				Steps:      steps,
				Quality:    192, // Good quality
				Timestamp:  time.Now().Format(time.RFC3339),
			}

			h.sendEvent(c.Writer, flusher, event)
		}
	}
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
