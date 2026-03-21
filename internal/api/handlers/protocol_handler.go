package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

// ProtocolHandler 協議資訊 API Handler
type ProtocolHandler struct{}

// NewProtocolHandler 建立新的協議 Handler
func NewProtocolHandler() *ProtocolHandler {
	return &ProtocolHandler{}
}

// ProtocolInfo 協議資訊
type ProtocolInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ConfigSchema string `json:"config_schema"`
}

// List 列出所有支援的協議
// GET /datalink/protocols
func (h *ProtocolHandler) List(c *gin.Context) {
	protocols := []ProtocolInfo{
		{
			Type:        string(schema.ProtocolModbusTCP),
			Name:        "Modbus TCP",
			Description: "Modbus TCP/IP 協議",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"host": {"type": "string", "description": "主機位址"},
					"port": {"type": "integer", "default": 502, "description": "埠號"},
					"slave_id": {"type": "integer", "minimum": 1, "maximum": 247, "description": "Slave ID"},
					"timeout": {"type": "integer", "default": 3000, "description": "逾時 (ms)"}
				},
				"required": ["host", "slave_id"]
			}`,
		},
		{
			Type:        string(schema.ProtocolModbusRTU),
			Name:        "Modbus RTU",
			Description: "Modbus RTU 串列協議",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"serial_port": {"type": "string", "description": "串列埠"},
					"baud_rate": {"type": "integer", "default": 9600, "description": "鮑率"},
					"data_bits": {"type": "integer", "default": 8, "description": "資料位元"},
					"stop_bits": {"type": "integer", "default": 1, "description": "停止位元"},
					"parity": {"type": "string", "enum": ["none", "odd", "even"], "default": "none", "description": "同位檢查"},
					"slave_id": {"type": "integer", "minimum": 1, "maximum": 247, "description": "Slave ID"},
					"timeout": {"type": "integer", "default": 3000, "description": "逾時 (ms)"}
				},
				"required": ["serial_port", "slave_id"]
			}`,
		},
		{
			Type:        string(schema.ProtocolModbusUDP),
			Name:        "Modbus UDP",
			Description: "Modbus UDP 協議",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"host": {"type": "string", "description": "主機位址"},
					"port": {"type": "integer", "default": 502, "description": "埠號"},
					"slave_id": {"type": "integer", "minimum": 1, "maximum": 247, "description": "Slave ID"},
					"timeout": {"type": "integer", "default": 3000, "description": "逾時 (ms)"}
				},
				"required": ["host", "slave_id"]
			}`,
		},
		{
			Type:        string(schema.ProtocolFatekFBs),
			Name:        "FATEK FBs",
			Description: "FATEK FBs PLC 協議",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"mode": {"type": "string", "enum": ["tcp", "serial"], "description": "連線模式"},
					"host": {"type": "string", "description": "主機位址 (TCP 模式)"},
					"port": {"type": "integer", "default": 500, "description": "埠號 (TCP 模式)"},
					"serial_port": {"type": "string", "description": "串列埠 (Serial 模式)"},
					"baud_rate": {"type": "integer", "default": 9600, "description": "鮑率 (Serial 模式)"},
					"data_bits": {"type": "integer", "enum": [7, 8], "default": 7, "description": "資料位元 (Serial 模式)"},
					"stop_bits": {"type": "integer", "enum": [1, 2], "default": 1, "description": "停止位元 (Serial 模式)"},
					"parity": {"type": "string", "enum": ["none", "odd", "even"], "default": "even", "description": "同位檢查 (Serial 模式)"},
					"station_no": {"type": "integer", "minimum": 0, "maximum": 255, "description": "站號"},
					"timeout": {"type": "integer", "default": 3000, "description": "逾時 (ms)"}
				},
				"required": ["mode", "station_no"]
			}`,
		},
		{
			Type:        string(schema.ProtocolMC3E),
			Name:        "Mitsubishi MC 3E",
			Description: "三菱 MC 3E 協議 (Binary)",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"host": {"type": "string", "description": "主機位址"},
					"port": {"type": "integer", "default": 5000, "description": "埠號"},
					"network_no": {"type": "integer", "default": 0, "description": "網路號"},
					"pc_no": {"type": "integer", "default": 255, "description": "PC 號"},
					"io_no": {"type": "integer", "default": 1023, "description": "IO 號"},
					"station_no": {"type": "integer", "default": 0, "description": "站號"},
					"data_format": {"type": "string", "enum": ["ABCD", "BADC", "CDAB", "DCBA"], "default": "CDAB", "description": "字節序格式"},
					"timeout": {"type": "integer", "default": 3000, "description": "逾時 (ms)"}
				},
				"required": ["host"]
			}`,
		},
		{
			Type:        string(schema.ProtocolMQTT),
			Name:        "MQTT",
			Description: "MQTT 訂閱接收",
			ConfigSchema: `{
				"type": "object",
				"properties": {
					"broker": {"type": "string", "description": "Broker 位址"},
					"port": {"type": "integer", "default": 1883, "description": "埠號"},
					"client_id": {"type": "string", "description": "客戶端 ID"},
					"username": {"type": "string", "description": "使用者名稱"},
					"password": {"type": "string", "description": "密碼"},
					"topic": {"type": "string", "description": "訂閱主題"},
					"qos": {"type": "integer", "enum": [0, 1, 2], "default": 0, "description": "QoS 等級"}
				},
				"required": ["broker", "topic"]
			}`,
		},
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": protocols})
}
