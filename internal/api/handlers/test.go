package handlers

import (
"sync"
"time"
)

// TestHandler 處理協議測試相關的 API 請求
type TestHandler struct {
connections    map[string]*ConnectionState
activeMonitors map[string]chan struct{}
wsHandler      *WebSocketHandler // 保留用於其他用途
sseHandler     *SSEHandler       // SSE 處理器
debugHandler   *DebugHandler     // Debug 處理器
mu             sync.RWMutex
}

// ConnectionState 連線狀態
type ConnectionState struct {
ID        string                 `json:"id"`
Protocol  string                 `json:"protocol"`
Config    map[string]interface{} `json:"config"`
Connected bool                   `json:"connected"`
CreatedAt time.Time              `json:"created_at"`
Client    interface{}            `json:"-"` // 協議客戶端實例
}

// NewTestHandler 建立新的測試處理器
func NewTestHandler(wsHandler *WebSocketHandler, sseHandler *SSEHandler, debugHandler *DebugHandler) *TestHandler {
return &TestHandler{
connections:    make(map[string]*ConnectionState),
activeMonitors: make(map[string]chan struct{}),
wsHandler:      wsHandler,
sseHandler:     sseHandler,
debugHandler:   debugHandler,
}
}
