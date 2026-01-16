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

// Broadcast 向所有連線的客戶端廣播訊息
func (h *WebSocketHandler) Broadcast(message interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("廣播訊息失敗, 移除連線: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

// HandleWebSocket 處理一般 WebSocket 連線 (Echo)
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}
	defer conn.Close()

	// 這裡僅作簡單 Echo 測試，不加入廣播池
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, message); err != nil {
			break
		}
	}
}

// HandleMonitorStream 處理監控數據流
func (h *WebSocketHandler) HandleMonitorStream(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升級失敗: %v", err)
		return
	}
	// 不在此處 defer conn.Close()，由連線斷開或錯誤時處理

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// 保持連線開啟，監聽關閉事件
	// 對於監控流，我們主要從伺服器推送數據，客戶端不需要發送太多數據
	// 但我們需要讀取以偵測斷線
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	conn.Close()
}
