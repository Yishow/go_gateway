package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExecuteScript 執行測試腳本
func (h *TestHandler) ExecuteScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{apiResponseErrorKey: notImplementedMessage})
}

// ListScripts 取得測試腳本列表
func (h *TestHandler) ListScripts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"scripts": []interface{}{}})
}

// SaveScript 保存測試腳本
func (h *TestHandler) SaveScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{apiResponseErrorKey: notImplementedMessage})
}

// DeleteScript 刪除測試腳本
func (h *TestHandler) DeleteScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{apiResponseErrorKey: notImplementedMessage})
}
