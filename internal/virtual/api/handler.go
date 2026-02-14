package api

import (
	"encoding/hex"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-gateway/internal/virtual/memory"
	modbusserver "go-gateway/internal/virtual/server/modbus"
	"go-gateway/internal/virtual/simulation"
)

// =============================================================================
// 虛擬設備 API Handler
// =============================================================================

// VirtualDeviceHandler 虛擬設備 API 處理器
type VirtualDeviceHandler struct {
	bank       *memory.MemoryBank
	server     *modbusserver.Server
	simulation *simulation.SimulationEngine
}

// NewVirtualDeviceHandler 建立新的虛擬設備處理器
func NewVirtualDeviceHandler(bank *memory.MemoryBank, server *modbusserver.Server, sim *simulation.SimulationEngine) *VirtualDeviceHandler {
	return &VirtualDeviceHandler{
		bank:       bank,
		server:     server,
		simulation: sim,
	}
}

// RegisterRoutes 註冊路由
func (h *VirtualDeviceHandler) RegisterRoutes(r *gin.RouterGroup) {
	v := r.Group("/virtual")
	// 記憶體操作
	v.GET("/memory/dump", h.GetMemoryDump)
	v.GET("/memory/range", h.GetMemoryRange)
	v.POST("/memory/write", h.WriteMemory)
	v.POST("/memory/clear", h.ClearMemory)

	// 伺服器控制
	v.POST("/server/start", h.StartServer)
	v.POST("/server/stop", h.StopServer)
	v.GET("/server/status", h.GetServerStatus)

	// 模擬控制
	v.POST("/simulation/start", h.StartSimulation)
	v.POST("/simulation/stop", h.StopSimulation)
	v.POST("/simulation/rules", h.AddSimulationRule)
}

func toUint16(value int) (uint16, bool) {
	if value < 0 || value > math.MaxUint16 {
		return 0, false
	}
	return uint16(value), true
}

// =============================================================================
// 記憶體 API
// =============================================================================

// MemoryDumpResponse 記憶體快照回應
type MemoryDumpResponse struct {
	Size int    `json:"size"`
	Data string `json:"data"` // Hex 編碼
}

// GetMemoryDump 取得完整記憶體快照
// GET /api/virtual/memory/dump
func (h *VirtualDeviceHandler) GetMemoryDump(c *gin.Context) {
	dump := h.bank.Dump()
	c.JSON(http.StatusOK, MemoryDumpResponse{
		Size: len(dump),
		Data: hex.EncodeToString(dump),
	})
}

// MemoryRangeRequest 範圍讀取請求
type MemoryRangeRequest struct {
	Offset int `form:"offset"`
	Length int `form:"length"`
}

// MemoryRangeResponse 範圍讀取回應
type MemoryRangeResponse struct {
	Offset int      `json:"offset"`
	Length int      `json:"length"`
	Data   string   `json:"data"`  // Hex 編碼
	Words  []uint16 `json:"words"` // 16-bit 值陣列
}

// GetMemoryRange 取得指定範圍記憶體
// GET /api/virtual/memory/range?offset=0&length=10
func (h *VirtualDeviceHandler) GetMemoryRange(c *gin.Context) {
	var req MemoryRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Length <= 0 {
		req.Length = 10
	}

	data, err := h.bank.ReadSlice(req.Offset, req.Length)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換為 16-bit 值
	words := make([]uint16, 0)
	for i := 0; i+1 < len(data); i += 2 {
		word, err := h.bank.ReadWord(req.Offset + i)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		words = append(words, word)
	}

	c.JSON(http.StatusOK, MemoryRangeResponse{
		Offset: req.Offset,
		Length: req.Length,
		Data:   hex.EncodeToString(data),
		Words:  words,
	})
}

// WriteMemoryRequest 寫入記憶體請求
type WriteMemoryRequest struct {
	Offset int    `json:"offset"`
	Data   string `json:"data"`  // Hex 編碼
	Value  *int   `json:"value"` // 或直接寫入 16-bit 值
}

// WriteMemory 寫入記憶體
// POST /api/virtual/memory/write
func (h *VirtualDeviceHandler) WriteMemory(c *gin.Context) {
	var req WriteMemoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Value != nil {
		// 寫入單一 16-bit 值
		word, ok := toUint16(*req.Value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "value 超出 uint16 範圍"})
			return
		}
		err := h.bank.WriteWord(req.Offset, word)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if req.Data != "" {
		// 寫入 Hex 數據
		data, err := hex.DecodeString(req.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的 Hex 數據"})
			return
		}
		err = h.bank.WriteSlice(req.Offset, data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ClearMemory 清空記憶體
// POST /api/virtual/memory/clear
func (h *VirtualDeviceHandler) ClearMemory(c *gin.Context) {
	h.bank.Clear()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// =============================================================================
// 伺服器 API
// =============================================================================

// StartServerRequest 啟動伺服器請求
type StartServerRequest struct {
	Port int `json:"port"`
}

// StartServer 啟動 Modbus 伺服器
// POST /api/virtual/server/start
func (h *VirtualDeviceHandler) StartServer(c *gin.Context) {
	var req StartServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Port = 502 // 預設 Modbus 端口
	}

	if h.server == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "伺服器未初始化"})
		return
	}

	err := h.server.Start(req.Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"port":    h.server.Port(),
		"address": h.server.Address(),
	})
}

// StopServer 停止 Modbus 伺服器
// POST /api/virtual/server/stop
func (h *VirtualDeviceHandler) StopServer(c *gin.Context) {
	if h.server == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "伺服器未初始化"})
		return
	}

	err := h.server.Stop()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetServerStatus 取得伺服器狀態
// GET /api/virtual/server/status
func (h *VirtualDeviceHandler) GetServerStatus(c *gin.Context) {
	if h.server == nil {
		c.JSON(http.StatusOK, gin.H{"running": false})
		return
	}

	port := h.server.Port()
	c.JSON(http.StatusOK, gin.H{
		"running": port > 0,
		"port":    port,
		"address": h.server.Address(),
	})
}

// =============================================================================
// 模擬 API
// =============================================================================

// StartSimulation 啟動模擬
// POST /api/virtual/simulation/start
func (h *VirtualDeviceHandler) StartSimulation(c *gin.Context) {
	if h.simulation == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模擬引擎未初始化"})
		return
	}

	h.simulation.Start()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// StopSimulation 停止模擬
// POST /api/virtual/simulation/stop
func (h *VirtualDeviceHandler) StopSimulation(c *gin.Context) {
	if h.simulation == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模擬引擎未初始化"})
		return
	}

	h.simulation.Stop()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// AddSimulationRuleRequest 新增模擬規則請求
type AddSimulationRuleRequest struct {
	ID         string `json:"id" binding:"required"`
	Type       string `json:"type" binding:"required"` // increment, toggle, random, sine
	Offset     int    `json:"offset"`
	IntervalMs int    `json:"interval_ms"`
	Increment  int    `json:"increment"`
	MinValue   int    `json:"min_value"`
	MaxValue   int    `json:"max_value"`
}

// AddSimulationRule 新增模擬規則
// POST /api/virtual/simulation/rules
func (h *VirtualDeviceHandler) AddSimulationRule(c *gin.Context) {
	var req AddSimulationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.simulation == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "模擬引擎未初始化"})
		return
	}

	var ruleType simulation.RuleType
	switch req.Type {
	case "increment":
		ruleType = simulation.RuleAutoIncrement
	case "toggle":
		ruleType = simulation.RuleToggle
	case "random":
		ruleType = simulation.RuleRandom
	case "sine":
		ruleType = simulation.RuleSineWave
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的規則類型: " + req.Type})
		return
	}

	interval := 1000
	if req.IntervalMs > 0 {
		interval = req.IntervalMs
	}

	rule := &simulation.Rule{
		ID:        req.ID,
		RuleType:  ruleType,
		Offset:    req.Offset,
		Interval:  time.Duration(interval) * time.Millisecond,
		Increment: req.Increment,
		MinValue:  req.MinValue,
		MaxValue:  req.MaxValue,
	}

	h.simulation.AddRule(rule)
	c.JSON(http.StatusOK, gin.H{"success": true, "rule_id": req.ID})
}
