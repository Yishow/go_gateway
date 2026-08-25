package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ConnectRequest 連線請求
type ConnectRequest struct {
	Protocol string                 `json:"protocol" binding:"required"`
	Config   map[string]interface{} `json:"config" binding:"required"`
}

// Connect 建立連線
func (h *TestHandler) Connect(c *gin.Context) {
	// 添加調試日誌
	fmt.Printf("[DEBUG] Connect 請求收到: Method=%s, Path=%s\n", c.Request.Method, c.Request.URL.Path)
	fmt.Printf("[DEBUG] Request Headers: %+v\n", c.Request.Header)

	var req ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[ERROR] JSON 綁定失敗: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("請求格式錯誤: %v", err)})
		return
	}

	fmt.Printf("[DEBUG] 收到連線請求: Protocol=%s, Config=%+v\n", req.Protocol, req.Config)

	// 驗證協議名稱
	if req.Protocol == "" {
		fmt.Printf("[ERROR] 協議名稱為空\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "協議名稱不能為空"})
		return
	}

	// 驗證配置
	if req.Config == nil {
		fmt.Printf("[ERROR] 配置為空\n")
		c.JSON(http.StatusBadRequest, gin.H{"error": "配置不能為空"})
		return
	}

	fmt.Printf("[DEBUG] 開始創建客戶端: Protocol=%s\n", req.Protocol)

	// 先生成 connectionID，這樣可以在創建客戶端時就記錄數據包
	connID := generateConnectionID()

	client, err := h.createClientWithDebug(req.Protocol, req.Config, connID)
	if err != nil {
		fmt.Printf("[ERROR] 創建客戶端失敗: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("創建客戶端失敗: %v", err)})
		return
	}
	fmt.Printf("[DEBUG] 客戶端創建成功\n")

	// 嘗試連線
	fmt.Printf("[DEBUG] 開始嘗試連線...\n")
	if err := h.connectClient(client, req.Protocol); err != nil {
		// 使用更友好的錯誤訊息
		errMsg := categorizeError(err).Error()
		fmt.Printf("[ERROR] 連線失敗: %v\n", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	fmt.Printf("[DEBUG] 連線成功\n")

	state := &ConnectionState{
		ID:        connID,
		Protocol:  req.Protocol,
		Config:    req.Config,
		Connected: true,
		CreatedAt: time.Now(),
		Client:    client,
	}

	h.mu.Lock()
	h.connections[connID] = state
	h.mu.Unlock()

	// 記錄連線日誌
	if h.debugHandler != nil {
		h.debugHandler.RecordLog("info", fmt.Sprintf("連線建立: %s (%s)", connID, req.Protocol), map[string]interface{}{
			"connection_id": connID,
			"protocol":      req.Protocol,
			"config":        req.Config,
		})
	}

	fmt.Printf("[DEBUG] 連線建立完成: ConnectionID=%s\n", connID)
	c.JSON(http.StatusOK, gin.H{
		"connection_id": connID,
		"status":        "connected",
	})
}

// Disconnect 斷開連線
func (h *TestHandler) Disconnect(c *gin.Context) {
	connID := c.Query("connection_id")
	if connID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection_id is required"})
		return
	}

	h.mu.Lock()
	// 先停止該連線的監控任務（如果有的話）
	if stopChan, ok := h.activeMonitors[connID]; ok {
		close(stopChan)
		delete(h.activeMonitors, connID)
	}

	state, exists := h.connections[connID]
	if !exists {
		h.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "connection not found"})
		return
	}

	// 關閉連線
	if err := h.closeClient(state.Client, state.Protocol); err != nil {
		fmt.Printf("Error closing connection %s: %v\n", connID, err)
	}

	state.Connected = false
	delete(h.connections, connID)
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "disconnected"})
}

// GetStatus 取得連線狀態
func (h *TestHandler) GetStatus(c *gin.Context) {
	connID := c.Query("connection_id")
	if connID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection_id is required"})
		return
	}

	h.mu.RLock()
	state, exists := h.connections[connID]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "connection not found"})
		return
	}

	c.JSON(http.StatusOK, state)
}
