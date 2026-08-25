package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// DebugHandler 處理 Debug 相關的 API 請求
type DebugHandler struct {
	packets []PacketRecord
	logs    []LogRecord
	mu      sync.RWMutex
}

// PacketRecord 數據包記錄
type PacketRecord struct {
	ID           string      `json:"id"`
	Timestamp    time.Time   `json:"timestamp"`
	Direction    string      `json:"direction"` // "request" or "response"
	Protocol     string      `json:"protocol"`
	RawData      []byte      `json:"raw_data"`
	HexData      string      `json:"hex_data"`
	ParsedData   interface{} `json:"parsed_data,omitempty"`
	ConnectionID string      `json:"connection_id"`
}

// LogRecord 日誌記錄
type LogRecord struct {
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
}

// NewDebugHandler 建立新的 Debug 處理器
func NewDebugHandler() *DebugHandler {
	return &DebugHandler{
		packets: make([]PacketRecord, 0, 1000),
		logs:    make([]LogRecord, 0, 1000),
	}
}

// RecordPacket 記錄數據包
func (h *DebugHandler) RecordPacket(connectionID, protocol, direction string, data []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 限制數據包數量，保留最近 1000 條
	if len(h.packets) >= 1000 {
		h.packets = h.packets[1:]
	}

	packet := PacketRecord{
		ID:           fmt.Sprintf("pkt_%d_%s", time.Now().UnixNano(), connectionID),
		Timestamp:    time.Now(),
		Direction:    direction,
		Protocol:     protocol,
		RawData:      data,
		HexData:      bytesToHex(data),
		ConnectionID: connectionID,
	}

	h.packets = append(h.packets, packet)
}

// RecordLog 記錄日誌
func (h *DebugHandler) RecordLog(level, message string, details interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 限制日誌數量，保留最近 1000 條
	if len(h.logs) >= 1000 {
		h.logs = h.logs[1:]
	}

	log := LogRecord{
		ID:        fmt.Sprintf("log_%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Details:   details,
	}

	h.logs = append(h.logs, log)
}

// bytesToHex 將字節數組轉換為十六進制字符串
func bytesToHex(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	hexStr := ""
	for _, b := range data {
		hexStr += fmt.Sprintf("%02X ", b)
	}
	return strings.TrimSpace(hexStr)
}

// GetPackets 取得數據包記錄
func (h *DebugHandler) GetPackets(c *gin.Context) {
	limit := 100
	if l := c.Query("limit"); l != "" {
		// TODO: 解析 limit 參數
	}
	connectionID := c.Query("connection_id")

	h.mu.RLock()
	packets := h.packets

	// 按連接 ID 過濾
	if connectionID != "" {
		filtered := make([]PacketRecord, 0)
		for _, p := range packets {
			if p.ConnectionID == connectionID {
				filtered = append(filtered, p)
			}
		}
		packets = filtered
	}

	if len(packets) > limit {
		packets = packets[len(packets)-limit:]
	}
	h.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"packets": packets})
}

// GetLogs 取得日誌記錄
func (h *DebugHandler) GetLogs(c *gin.Context) {
	limit := 100
	if l := c.Query("limit"); l != "" {
		// TODO: 解析 limit 參數
	}
	connectionID := c.Query("connection_id")

	h.mu.RLock()
	logs := h.logs

	// 按連接 ID 過濾
	if connectionID != "" {
		filtered := make([]LogRecord, 0)
		for _, l := range logs {
			if details, ok := l.Details.(map[string]interface{}); ok {
				if cid, ok := details["connection_id"].(string); ok && cid == connectionID {
					filtered = append(filtered, l)
				}
			}
		}
		logs = filtered
	}

	if len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}
	h.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// ClearData 清空數據
func (h *DebugHandler) ClearData(c *gin.Context) {
	connectionID := c.Query("connection_id")

	h.mu.Lock()
	defer h.mu.Unlock()

	beforePacketsCount := len(h.packets)
	beforeLogsCount := len(h.logs)

	if connectionID != "" {
		// 只清空指定連接的數據
		// 使用預分配容量以提高性能
		filteredPackets := make([]PacketRecord, 0, len(h.packets))
		for _, p := range h.packets {
			if p.ConnectionID != connectionID {
				filteredPackets = append(filteredPackets, p)
			}
		}
		h.packets = filteredPackets

		filteredLogs := make([]LogRecord, 0, len(h.logs))
		for _, l := range h.logs {
			if details, ok := l.Details.(map[string]interface{}); ok {
				if cid, ok := details["connection_id"].(string); ok && cid != connectionID {
					filteredLogs = append(filteredLogs, l)
				}
			} else {
				filteredLogs = append(filteredLogs, l)
			}
		}
		h.logs = filteredLogs
	} else {
		// 清空所有數據
		h.packets = make([]PacketRecord, 0, 1000)
		h.logs = make([]LogRecord, 0, 1000)
	}

	afterPacketsCount := len(h.packets)
	afterLogsCount := len(h.logs)

	c.JSON(http.StatusOK, gin.H{
		apiResponseStatusKey: "cleared",
		"cleared_packets":    beforePacketsCount - afterPacketsCount,
		"cleared_logs":       beforeLogsCount - afterLogsCount,
		"remaining_packets":  afterPacketsCount,
		"remaining_logs":     afterLogsCount,
	})
}

// SendRawRequest 發送原始數據包請求
type SendRawRequest struct {
	ConnectionID string `json:"connection_id" binding:"required"`
	HexData      string `json:"hex_data" binding:"required"`
}

// SendRaw 發送原始數據包
func (h *DebugHandler) SendRaw(c *gin.Context) {
	var req SendRawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 實作原始數據包發送
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// AnalyzePacket 分析數據包結構
func (h *DebugHandler) AnalyzePacket(c *gin.Context) {
	packetID := c.Param("packetId")

	h.mu.RLock()
	var packet *PacketRecord
	for i := range h.packets {
		if h.packets[i].ID == packetID {
			packet = &h.packets[i]
			break
		}
	}
	h.mu.RUnlock()

	if packet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "packet not found"})
		return
	}

	// TODO: 實作協議分析邏輯
	c.JSON(http.StatusOK, gin.H{
		"packet": packet,
		"analysis": map[string]interface{}{
			"protocol":  packet.Protocol,
			"structure": "待實作",
		},
	})
}
