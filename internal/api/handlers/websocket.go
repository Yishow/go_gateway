package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允許所有來源（生產環境應限制）
	},
}

// WebSocketHandler 處理 WebSocket 連線
type WebSocketHandler struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// NewWebSocketHandler 建立新的 WebSocket 處理器
func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		clients: make(map[*websocket.Conn]bool),
	}
}

// HandleWebSocket 處理 WebSocket 連線
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}
	defer conn.Close()

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// 處理訊息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("讀取 WebSocket 訊息失敗: %v", err)
			break
		}

		// 回傳訊息
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("寫入 WebSocket 訊息失敗: %v", err)
			break
		}
	}

	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
}

// HandleMonitorStream 處理監控數據流
func (h *WebSocketHandler) HandleMonitorStream(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}
	defer conn.Close()

	// TODO: 實作監控數據流推送
	for {
		// 發送監控數據
		data := map[string]interface{}{
			"timestamp": "2024-01-01T00:00:00Z",
			"values":    []interface{}{},
		}

		if err := conn.WriteJSON(data); err != nil {
			log.Printf("發送監控數據失敗: %v", err)
			break
		}
	}
}
