package handlers

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
)

// TestHandler 處理協議測試相關的 API 請求
type TestHandler struct {
	connections    map[string]*ConnectionState
	activeMonitors map[string]chan struct{}
	wsHandler      *WebSocketHandler
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
func NewTestHandler(wsHandler *WebSocketHandler) *TestHandler {
	return &TestHandler{
		connections:    make(map[string]*ConnectionState),
		activeMonitors: make(map[string]chan struct{}),
		wsHandler:      wsHandler,
	}
}

// ConnectRequest 連線請求
type ConnectRequest struct {
	Protocol string                 `json:"protocol" binding:"required"`
	Config   map[string]interface{} `json:"config" binding:"required"`
}

// Connect 建立連線
func (h *TestHandler) Connect(c *gin.Context) {
	var req ConnectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.createClient(req.Protocol, req.Config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to create client: %v", err)})
		return
	}

	// 嘗試連線
	if err := h.connectClient(client, req.Protocol); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": categorizeError(err).Error()})
		return
	}

	connID := generateConnectionID()
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

// ReadRequest 讀取請求
type ReadRequest struct {
	ConnectionID string `json:"connection_id" binding:"required"`
	Operation    string `json:"operation" binding:"required"`
	Address      uint16 `json:"address"` // For Modbus/MC
	Count        uint16 `json:"count"`
	Symbol       string `json:"symbol,omitempty"` // For Fatek/MC
	Device       string `json:"device,omitempty"` // For MC
}

// Read 執行讀取操作
func (h *TestHandler) Read(c *gin.Context) {
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

	result, err := h.executeRead(state.Client, state.Protocol, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": categorizeError(err).Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"values": result,
		"count":  len(resultToString(result)),
	})
}

// WriteRequest 寫入請求
type WriteRequest struct {
	ConnectionID string      `json:"connection_id" binding:"required"`
	Operation    string      `json:"operation" binding:"required"`
	Address      uint16      `json:"address"`
	Values       interface{} `json:"values" binding:"required"`
	Symbol       string      `json:"symbol,omitempty"` // For Fatek
	Device       string      `json:"device,omitempty"` // For MC
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

	if err := h.executeWrite(state.Client, state.Protocol, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": categorizeError(err).Error()})
		return
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

		switch op.Type {
		case "read":
			if op.ReadRequest == nil {
				err = fmt.Errorf("missing read_request")
			} else {
				// Override ConnectionID to match the batch request context
				op.ReadRequest.ConnectionID = req.ConnectionID
				data, err = h.executeRead(state.Client, state.Protocol, *op.ReadRequest)
			}
		case "write":
			if op.WriteRequest == nil {
				err = fmt.Errorf("missing write_request")
			} else {
				op.WriteRequest.ConnectionID = req.ConnectionID
				err = h.executeWrite(state.Client, state.Protocol, *op.WriteRequest)
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

			// 推送數據
			msg := map[string]interface{}{
				"type":          "monitor_update",
				"connection_id": connID,
				"timestamp":     time.Now().Format(time.RFC3339Nano),
				"data":          results,
			}
			h.wsHandler.Broadcast(msg)
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

// categorizeError 將底層錯誤轉換為更易讀的錯誤訊息
func categorizeError(err error) error {
	if err == nil {
		return nil
	}
	
	// 處理 net.Error (Timeout, Connection refused)
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return fmt.Errorf("TIMEOUT: %v", err)
		}
		return fmt.Errorf("NETWORK_ERROR: %v", err)
	}

	// 檢查常見的錯誤字串 (因為部分庫可能返回普通 error)
	s := err.Error()
	if s == "EOF" {
		return fmt.Errorf("CONNECTION_CLOSED: Remote host closed connection")
	}
	// TODO: 可以根據具體協議庫的錯誤類型進行更細緻的分類
	// 例如: CRC Checksum Error, Illegal Function, etc.

	return err
}

// ExecuteScript 執行測試腳本
func (h *TestHandler) ExecuteScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ListScripts 取得測試腳本列表
func (h *TestHandler) ListScripts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"scripts": []interface{}{}})
}

// SaveScript 保存測試腳本
func (h *TestHandler) SaveScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DeleteScript 刪除測試腳本
func (h *TestHandler) DeleteScript(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// --- Helper Functions ---

func (h *TestHandler) createClient(protocol string, config map[string]interface{}) (interface{}, error) {
	// Helper to safe cast config values
	getString := func(key string, def string) string {
		if v, ok := config[key].(string); ok {
			return v
		}
		return def
	}
	getInt := func(key string, def int) int {
		if v, ok := config[key].(float64); ok {
			return int(v)
		}
		if v, ok := config[key].(int); ok {
			return v
		}
		return def
	}
	getDuration := func(key string, def time.Duration) time.Duration {
		if v, ok := config[key].(float64); ok {
			return time.Duration(v) * time.Millisecond
		}
		return def
	}

	switch protocol {
	case "modbus_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 502)
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return modbus.CreateTCPClient(host, port, unitID, timeout), nil
	
	case "modbus_rtu":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 8)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "N")
		unitID := byte(getInt("unitID", 1))
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return modbus.CreateRTUClient(port, baudRate, dataBits, stopBits, parity, timeout, unitID), nil

	case "fatek_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 500)
		station := getInt("station", 1)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return fatek.CreateTCPClient(host, port, station, timeout), nil

	case "fatek_serial":
		port := getString("port", "COM1")
		baudRate := getInt("baudRate", 9600)
		dataBits := getInt("dataBits", 7)
		stopBits := getInt("stopBits", 1)
		parity := getString("parity", "E")
		station := getInt("station", 1)
		timeout := getDuration("timeout", 1000*time.Millisecond)
		return fatek.CreateSerialClient(port, station, baudRate, dataBits, stopBits, parity, timeout), nil
	
	case "mc_tcp":
		host := getString("host", "localhost")
		port := getInt("port", 6000)
		timeout := getDuration("timeout", 2000*time.Millisecond)
		return mcprotocol.CreateTCPClient(host, port, timeout), nil

	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

func (h *TestHandler) connectClient(client interface{}, protocol string) error {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		return c.Connect()
	case *fatek.FatekClient:
		return c.Connect()
	case *mcprotocol.MCClient:
		return c.Connect()
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) closeClient(client interface{}, protocol string) error {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		return c.Close()
	case *fatek.FatekClient:
		return c.Close()
	case *mcprotocol.MCClient:
		return c.Close()
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) executeRead(client interface{}, protocol string, req ReadRequest) (interface{}, error) {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		switch req.Operation {
		case "read_coils":
			return c.ReadCoils(req.Address, req.Count)
		case "read_discrete_inputs":
			return c.ReadDiscreteInputs(req.Address, req.Count)
		case "read_holding_registers":
			return c.ReadHoldingRegisters(req.Address, req.Count)
		case "read_input_registers":
			return c.ReadInputRegisters(req.Address, req.Count)
		default:
			return nil, fmt.Errorf("unsupported operation for modbus: %s", req.Operation)
		}

	case *fatek.FatekClient:
		switch req.Operation {
		case "read_status":
			return c.ReadStatus(req.Symbol, int(req.Address), int(req.Count))
		case "read_registers":
			return c.ReadRegisters(req.Symbol, int(req.Address), int(req.Count))
		default:
			return nil, fmt.Errorf("unsupported operation for fatek: %s", req.Operation)
		}

	case *mcprotocol.MCClient:
		switch req.Operation {
		case "batch_read_word":
			return c.BatchReadWord(req.Device, int(req.Address), int(req.Count))
		case "batch_read_bit":
			return c.BatchReadBit(req.Device, int(req.Address), int(req.Count))
		default:
			return nil, fmt.Errorf("unsupported operation for mcprotocol: %s", req.Operation)
		}

	default:
		return nil, fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) executeWrite(client interface{}, protocol string, req WriteRequest) error {
	// Helper to convert interface{} to []uint16 or []bool
	// This is tricky because JSON unmarshaling might give []interface{}
	
	toUint16Slice := func(v interface{}) ([]uint16, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]uint16, len(arr))
		for i, val := range arr {
			if f, ok := val.(float64); ok {
				res[i] = uint16(f)
			} else {
				return nil, fmt.Errorf("invalid value type at index %d", i)
			}
		}
		return res, nil
	}
	
	toIntSlice := func(v interface{}) ([]int, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]int, len(arr))
		for i, val := range arr {
			if f, ok := val.(float64); ok {
				res[i] = int(f)
			} else {
				return nil, fmt.Errorf("invalid value type at index %d", i)
			}
		}
		return res, nil
	}

	toBoolSlice := func(v interface{}) ([]bool, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]bool, len(arr))
		for i, val := range arr {
			if b, ok := val.(bool); ok {
				res[i] = b
			} else {
				return nil, fmt.Errorf("invalid value type at index %d", i)
			}
		}
		return res, nil
	}

	switch c := client.(type) {
	case *modbus.ModbusClient:
		switch req.Operation {
		case "write_single_coil":
			val, ok := req.Values.(bool)
			if !ok { return fmt.Errorf("value must be bool") }
			return c.WriteSingleCoil(req.Address, val)
		case "write_single_register":
			val, ok := req.Values.(float64)
			if !ok { return fmt.Errorf("value must be number") }
			return c.WriteSingleRegister(req.Address, uint16(val))
		case "write_multiple_coils":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.WriteMultipleCoils(req.Address, vals)
		case "write_multiple_registers":
			vals, err := toUint16Slice(req.Values)
			if err != nil { return err }
			return c.WriteMultipleRegisters(req.Address, vals)
		default:
			return fmt.Errorf("unsupported operation for modbus: %s", req.Operation)
		}

	case *fatek.FatekClient:
		switch req.Operation {
		case "write_status":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.WriteStatus(req.Symbol, int(req.Address), vals)
		case "write_registers":
			vals, err := toIntSlice(req.Values)
			if err != nil { return err }
			return c.WriteRegisters(req.Symbol, int(req.Address), vals)
		default:
			return fmt.Errorf("unsupported operation for fatek: %s", req.Operation)
		}

	case *mcprotocol.MCClient:
		switch req.Operation {
		case "batch_write_word":
			vals, err := toIntSlice(req.Values)
			if err != nil { return err }
			return c.BatchWriteWord(req.Device, int(req.Address), vals)
		case "batch_write_bit":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.BatchWriteBit(req.Device, int(req.Address), vals)
		default:
			return fmt.Errorf("unsupported operation for mcprotocol: %s", req.Operation)
		}

	default:
		return fmt.Errorf("unknown client type")
	}
}

func resultToString(v interface{}) []string {
	// Simple helper for counting only
	// Implementation dependent on usage, currently just len
	if arr, ok := v.([]interface{}); ok {
		res := make([]string, len(arr))
		return res
	}
	return []string{}
}


