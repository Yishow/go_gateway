package handlers

import (
"fmt"
"net/http"

"github.com/gin-gonic/gin"
)

// WriteRequest 寫入請求
type WriteRequest struct {
	ConnectionID string      `json:"connection_id" binding:"required"`
	Operation    string      `json:"operation" binding:"required"`
	Address      uint16      `json:"address"`
	Values       interface{} `json:"values" binding:"required"`
	Symbol       string      `json:"symbol,omitempty"` // For Fatek
	Device       string      `json:"device,omitempty"` // For MC
	UnitID       *byte       `json:"unit_id,omitempty"` // For Modbus (可選，覆蓋連線配置的站號)
	Station      *int        `json:"station,omitempty"` // For Fatek (可選，覆蓋連線配置的站號)
}

// Write 執行寫入操作
func (h *TestHandler) Write(c *gin.Context) {
	var req WriteRequest
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
		// 創建臨時客戶端配置
		tempConfig := make(map[string]interface{})
		for k, v := range state.Config {
			tempConfig[k] = v
		}
		// 覆蓋站號
		if req.UnitID != nil {
			tempConfig["unitID"] = int(*req.UnitID)
		}
		if req.Station != nil {
			tempConfig["station"] = *req.Station
		}
		// 創建臨時客戶端並連線
		tempClient, err := h.createClientWithDebug(state.Protocol, tempConfig, state.ID+"_temp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("創建臨時客戶端失敗: %v", err)})
			return
		}
		// 連線臨時客戶端
		if err := h.connectClient(tempClient, state.Protocol); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("連線臨時客戶端失敗: %v", err)})
			return
		}
		clientToUse = tempClient
	}

	if err := h.executeWrite(clientToUse, state.Protocol, req); err != nil {
		// 如果使用了臨時客戶端，關閉它
		if tempClient != nil {
			h.closeClient(tempClient, state.Protocol)
		}
		// 記錄錯誤日誌
		if h.debugHandler != nil {
			h.debugHandler.RecordLog("error", fmt.Sprintf("寫入失敗: %s", err.Error()), map[string]interface{}{
				"connection_id": req.ConnectionID,
				"operation":     req.Operation,
				"address":       req.Address,
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": categorizeError(err).Error()})
		return
	}

	// 如果使用了臨時客戶端，關閉它
	if tempClient != nil {
		h.closeClient(tempClient, state.Protocol)
	}

	// 記錄成功日誌
	if h.debugHandler != nil {
		h.debugHandler.RecordLog("info", fmt.Sprintf("寫入成功: %s (地址: %d)", req.Operation, req.Address), map[string]interface{}{
			"connection_id": req.ConnectionID,
			"operation":     req.Operation,
			"address":       req.Address,
		})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// BatchOperation 定義批量操作中的單個項目
type BatchOperation struct {
	Type         string      `json:"type" binding:"required"` // "read" or "write"
	ReadRequest  *ReadRequest `json:"read_request,omitempty"`
	WriteRequest *WriteRequest `json:"write_request,omitempty"`
}

// BatchRequest 批量測試請求
type BatchRequest struct {
	ConnectionID string           `json:"connection_id" binding:"required"`
	Operations   []BatchOperation `json:"operations" binding:"required"`
}

// BatchResult 批量測試結果
type BatchResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Batch 執行批量測試
func (h *TestHandler) Batch(c *gin.Context) {
	var req BatchRequest
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

	results := make([]BatchResult, len(req.Operations))

	for i, op := range req.Operations {
		var err error
		var data interface{}
		var tempClient interface{} = nil

		switch op.Type {
		case "read":
			if op.ReadRequest == nil {
				err = fmt.Errorf("missing read_request")
			} else {
				// Override ConnectionID to match the batch request context
				op.ReadRequest.ConnectionID = req.ConnectionID
				// 如果請求中指定了站號，創建臨時客戶端
				var clientToUse interface{} = state.Client
				if op.ReadRequest.UnitID != nil || op.ReadRequest.Station != nil {
					// 創建臨時客戶端配置
					tempConfig := make(map[string]interface{})
					for k, v := range state.Config {
						tempConfig[k] = v
					}
					// 覆蓋站號
					if op.ReadRequest.UnitID != nil {
						tempConfig["unitID"] = int(*op.ReadRequest.UnitID)
					}
					if op.ReadRequest.Station != nil {
						tempConfig["station"] = *op.ReadRequest.Station
					}
					// 創建臨時客戶端並連線
					tempClient, err = h.createClientWithDebug(state.Protocol, tempConfig, state.ID+"_temp_batch")
					if err != nil {
						err = fmt.Errorf("創建臨時客戶端失敗: %v", err)
					} else {
						// 連線臨時客戶端
						if err = h.connectClient(tempClient, state.Protocol); err != nil {
							err = fmt.Errorf("連線臨時客戶端失敗: %v", err)
						} else {
							clientToUse = tempClient
						}
					}
				}
				if err == nil {
					data, err = h.executeRead(clientToUse, state.Protocol, *op.ReadRequest)
				}
				// 如果使用了臨時客戶端，關閉它
				if tempClient != nil {
					h.closeClient(tempClient, state.Protocol)
				}
			}
		case "write":
			if op.WriteRequest == nil {
				err = fmt.Errorf("missing write_request")
			} else {
				op.WriteRequest.ConnectionID = req.ConnectionID
				// 如果請求中指定了站號，創建臨時客戶端
				var clientToUse interface{} = state.Client
				if op.WriteRequest.UnitID != nil || op.WriteRequest.Station != nil {
					// 創建臨時客戶端配置
					tempConfig := make(map[string]interface{})
					for k, v := range state.Config {
						tempConfig[k] = v
					}
					// 覆蓋站號
					if op.WriteRequest.UnitID != nil {
						tempConfig["unitID"] = int(*op.WriteRequest.UnitID)
					}
					if op.WriteRequest.Station != nil {
						tempConfig["station"] = *op.WriteRequest.Station
					}
					// 創建臨時客戶端並連線
					tempClient, err = h.createClientWithDebug(state.Protocol, tempConfig, state.ID+"_temp_batch")
					if err != nil {
						err = fmt.Errorf("創建臨時客戶端失敗: %v", err)
					} else {
						// 連線臨時客戶端
						if err = h.connectClient(tempClient, state.Protocol); err != nil {
							err = fmt.Errorf("連線臨時客戶端失敗: %v", err)
						} else {
							clientToUse = tempClient
						}
					}
				}
				if err == nil {
					err = h.executeWrite(clientToUse, state.Protocol, *op.WriteRequest)
				}
				// 如果使用了臨時客戶端，關閉它
				if tempClient != nil {
					h.closeClient(tempClient, state.Protocol)
				}
			}
		default:
			err = fmt.Errorf("unknown operation type: %s", op.Type)
		}

		if err != nil {
			results[i] = BatchResult{
				Success: false,
				Error:   categorizeError(err).Error(),
			}
		} else {
			results[i] = BatchResult{
				Success: true,
				Data:    data,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
