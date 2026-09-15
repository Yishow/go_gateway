package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ConnectionHandler 處理連線池相關的 API 請求
type ConnectionHandler struct {
	testHandler *TestHandler
}

// NewConnectionHandler 建立新的連線處理器
func NewConnectionHandler(testHandler *TestHandler) *ConnectionHandler {
	return &ConnectionHandler{
		testHandler: testHandler,
	}
}

// List 取得連線池狀態列表
func (h *ConnectionHandler) List(c *gin.Context) {
	h.testHandler.mu.RLock()
	defer h.testHandler.mu.RUnlock()

	connections := make([]*ConnectionState, 0, len(h.testHandler.connections))
	for _, conn := range h.testHandler.connections {
		connections = append(connections, conn)
	}

	c.JSON(http.StatusOK, gin.H{"connections": connections})
}

// Get 取得特定連線詳情
func (h *ConnectionHandler) Get(c *gin.Context) {
	id := c.Param("id")

	h.testHandler.mu.RLock()
	defer h.testHandler.mu.RUnlock()

	if conn, exists := h.testHandler.connections[id]; exists {
		c.JSON(http.StatusOK, conn)
	} else {
		c.JSON(http.StatusNotFound, gin.H{apiResponseErrorKey: connectionNotFoundMessage})
	}
}
