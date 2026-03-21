package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReadRequest 讀取請求
type ReadRequest struct {
	ConnectionID string `json:"connection_id" binding:"required"`
	Operation    string `json:"operation" binding:"required"`
	Address      uint16 `json:"address"` // For Modbus/MC
	Count        uint16 `json:"count"`
	Symbol       string `json:"symbol,omitempty"`  // For Fatek/MC
	Device       string `json:"device,omitempty"`  // For MC
	UnitID       *byte  `json:"unit_id,omitempty"` // For Modbus (可選，覆蓋連線配置的站號)
	Station      *int   `json:"station,omitempty"` // For Fatek (可選，覆蓋連線配置的站號)
}

// Read 執行讀取操作
func (h *TestHandler) Read(c *gin.Context) {
	// 添加 panic 恢復
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC] Read 操作發生 panic: %v\n", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("內部錯誤: %v", r)})
		}
	}()

	var req ReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.mu.RLock()
	state, exists := h.connections[req.ConnectionID]
	h.mu.RUnlock()

	if !exists || !state.Connected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection not found or not connected"})
		return
	}

	// 如果請求中指定了站號，創建臨時客戶端
	var clientToUse interface{} = state.Client
	var tempClient interface{} = nil
	if req.UnitID != nil || req.Station != nil {
		var err error
		clientToUse, tempClient, err = h.prepareOverrideClient(state, req.UnitID, req.Station, "_temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	result, err := h.executeRead(clientToUse, state.Protocol, req)

	// 如果使用了臨時客戶端，關閉它
	if tempClient != nil {
		h.closeTempClient(tempClient, state.Protocol)
	}
	if err != nil {
		// 記錄錯誤日誌
		if h.debugHandler != nil {
			h.debugHandler.RecordLog("error", fmt.Sprintf("讀取失敗: %s", err.Error()), map[string]interface{}{
				"connection_id": req.ConnectionID,
				"operation":     req.Operation,
				"address":       req.Address,
				"count":         req.Count,
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": categorizeError(err).Error()})
		return
	}

	// 計算結果數量（安全處理）
	count := 0
	if result != nil {
		count = len(resultToString(result))
	}

	// 記錄成功日誌
	if h.debugHandler != nil {
		h.debugHandler.RecordLog("info", fmt.Sprintf("讀取成功: %s (地址: %d, 數量: %d)", req.Operation, req.Address, req.Count), map[string]interface{}{
			"connection_id": req.ConnectionID,
			"operation":     req.Operation,
			"address":       req.Address,
			"count":         req.Count,
			"result_count":  count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"values": result,
		"count":  count,
	})
}

// resultToString 將結果轉換為字符串數組（用於計算長度）
func resultToString(v interface{}) []string {
	if v == nil {
		return []string{}
	}

	// 處理各種類型的數組
	switch arr := v.(type) {
	case []bool:
		res := make([]string, len(arr))
		return res
	case []uint16:
		res := make([]string, len(arr))
		return res
	case []int:
		res := make([]string, len(arr))
		return res
	case []interface{}:
		res := make([]string, len(arr))
		return res
	case []string:
		return arr
	default:
		// 嘗試使用反射獲取長度
		// 如果無法確定，返回空數組
		return []string{}
	}
}
