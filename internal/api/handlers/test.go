package handlers

import (
	"fmt"
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
	connections map[string]*ConnectionState
	mu          sync.RWMutex
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
func NewTestHandler() *TestHandler {
	return &TestHandler{
		connections: make(map[string]*ConnectionState),
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to connect: %v", err)})
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
	state, exists := h.connections[connID]
	if !exists {
		h.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "connection not found"})
		return
	}

	// 關閉連線
	if err := h.closeClient(state.Client, state.Protocol); err != nil {
		// Log error but continue to remove connection
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"values": result,
		"count":  len(resultToString(result)), // Helper to count items
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// BatchRequest 批量測試請求
type BatchRequest struct {
	ConnectionID string   `json:"connection_id" binding:"required"`
	Operations   []string `json:"operations" binding:"required"`
}

// Batch 執行批量測試
func (h *TestHandler) Batch(c *gin.Context) {
	// 暫時僅保留框架，具體實作需定義更複雜的結構
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
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

// StartMonitor 啟動監控模式
func (h *TestHandler) StartMonitor(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// StopMonitor 停止監控模式
func (h *TestHandler) StopMonitor(c *gin.Context) {
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


