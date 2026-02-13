package handlers

import (
"fmt"
"net/http"
"time"

"github.com/gin-gonic/gin"
)

// MonitorRequest 監控請求
type MonitorRequest struct {
	ConnectionID string        `json:"connection_id" binding:"required"`
	Items        []ReadRequest `json:"items" binding:"required"`
	Interval     int           `json:"interval"` // 毫秒
}

// StartMonitor 啟動監控模式
func (h *TestHandler) StartMonitor(c *gin.Context) {
	var req MonitorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Interval < 100 {
		req.Interval = 100 // 最小 100ms
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// 檢查連線
	state, exists := h.connections[req.ConnectionID]
	if !exists || !state.Connected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection not found or not connected"})
		return
	}

	// 如果已經有監控任務，先停止
	if stopChan, ok := h.activeMonitors[req.ConnectionID]; ok {
		close(stopChan)
		delete(h.activeMonitors, req.ConnectionID)
	}

	// 啟動新的監控任務
	stopChan := make(chan struct{})
	h.activeMonitors[req.ConnectionID] = stopChan

	go h.runMonitorLoop(req.ConnectionID, state, req.Items, req.Interval, stopChan)

	c.JSON(http.StatusOK, gin.H{"status": "monitoring_started"})
}

func (h *TestHandler) runMonitorLoop(connID string, state *ConnectionState, items []ReadRequest, interval int, stopChan chan struct{}) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			// 執行讀取
			results := make(map[string]interface{})
			for i, item := range items {
				res, err := h.executeRead(state.Client, state.Protocol, item)
				key := fmt.Sprintf("item_%d", i) // 或者使用地址/符號作為 key
				if err != nil {
					results[key] = map[string]string{"error": categorizeError(err).Error()}
				} else {
					results[key] = res
				}
			}

			// 推送數據（使用 SSE）
			msg := map[string]interface{}{
				"type":          "monitor_update",
				"connection_id": connID,
				"timestamp":     time.Now().Format(time.RFC3339Nano),
				"data":          results,
			}
			h.sseHandler.Broadcast(connID, msg)
		}
	}
}

// StopMonitor 停止監控模式
func (h *TestHandler) StopMonitor(c *gin.Context) {
	var req struct {
		ConnectionID string `json:"connection_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.mu.Lock()
	if stopChan, ok := h.activeMonitors[req.ConnectionID]; ok {
		close(stopChan)
		delete(h.activeMonitors, req.ConnectionID)
		h.mu.Unlock()
		c.JSON(http.StatusOK, gin.H{"status": "monitoring_stopped"})
	} else {
		h.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "no active monitor for this connection"})
	}
}
