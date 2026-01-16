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
func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		testHandler: NewTestHandler(),
	}
}

// List 取得連線池狀態列表
func (h *ConnectionHandler) List(c *gin.Context) {
	// TODO: 從 TestHandler 取得連線列表
	c.JSON(http.StatusOK, gin.H{"connections": []interface{}{}})
}

// Get 取得特定連線詳情
func (h *ConnectionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	// TODO: 從 TestHandler 取得連線詳情
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "not implemented"})
}
