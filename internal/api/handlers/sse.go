package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEHandler 處理 Server-Sent Events 連線
type SSEHandler struct {
	clients map[string]chan []byte // connection_id -> channel
	mu      sync.RWMutex
}

// NewSSEHandler 建立新的 SSE 處理器
func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		clients: make(map[string]chan []byte),
	}
}

// HandleMonitorStream 處理監控數據流的 SSE 連接
func (h *SSEHandler) HandleMonitorStream(c *gin.Context) {
	connectionID := c.Query("connection_id")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection_id is required"})
		return
	}

	// 設置 SSE 響應頭
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用 Nginx 緩衝

	// 創建客戶端通道
	clientChan := make(chan []byte, 10) // 緩衝區大小為 10

	h.mu.Lock()
	// 如果該連接已存在，關閉舊通道
	if oldChan, exists := h.clients[connectionID]; exists {
		close(oldChan)
	}
	h.clients[connectionID] = clientChan
	h.mu.Unlock()

	// 發送初始連接確認
	c.SSEvent("connected", gin.H{
		"connection_id": connectionID,
		"timestamp":     time.Now().Format(time.RFC3339Nano),
	})
	c.Writer.Flush()

	// 監聽客戶端斷線
	ctx := c.Request.Context()
	ticker := time.NewTicker(30 * time.Second) // 每 30 秒發送一次心跳
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// 客戶端斷線
			log.Printf("SSE 客戶端斷線: %s", connectionID)
			h.mu.Lock()
			if ch, exists := h.clients[connectionID]; exists && ch == clientChan {
				close(clientChan)
				delete(h.clients, connectionID)
			}
			h.mu.Unlock()
			return

		case <-ticker.C:
			// 發送心跳保持連接
			c.SSEvent("ping", gin.H{"timestamp": time.Now().Format(time.RFC3339Nano)})
			c.Writer.Flush()

		case data, ok := <-clientChan:
			if !ok {
				// 通道已關閉
				return
			}
			// 發送數據
			c.Writer.Write(data)
			c.Writer.Flush()
		}
	}
}

// Broadcast 向指定連接的客戶端發送訊息
func (h *SSEHandler) Broadcast(connectionID string, message interface{}) {
	h.mu.RLock()
	clientChan, exists := h.clients[connectionID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	// 將消息轉換為 SSE 格式
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("序列化消息失敗: %v", err)
		return
	}

	// SSE 格式: data: {json}\n\n
	sseData := []byte("data: " + string(jsonData) + "\n\n")

	select {
	case clientChan <- sseData:
		// 成功發送
	default:
		// 通道已滿，跳過這條消息
		log.Printf("SSE 通道已滿，跳過消息: %s", connectionID)
	}
}

// RemoveClient 移除客戶端連接
func (h *SSEHandler) RemoveClient(connectionID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, exists := h.clients[connectionID]; exists {
		close(ch)
		delete(h.clients, connectionID)
	}
}

// GetClientCount 獲取當前連接的客戶端數量
func (h *SSEHandler) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
