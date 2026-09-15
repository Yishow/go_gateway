package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 連接器工廠函數
// =============================================================================

// ConnectorFactory 連接器工廠函數類型
// 用於建立新的協議連接器實例
type ConnectorFactory func() Protocol //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.

// =============================================================================
// 註冊表實作
// =============================================================================

// Registry 協議連接器註冊表
// 管理所有已註冊的協議連接器工廠
type Registry struct {
	mu        sync.RWMutex
	factories map[schema.ProtocolType]ConnectorFactory
}

// 全域註冊表實例
var globalRegistry = &Registry{
	factories: make(map[schema.ProtocolType]ConnectorFactory),
}

// Register 註冊協議連接器工廠到全域註冊表
//
// 參數:
//   - protocolType: 協議類型
//   - factory: 建立連接器的工廠函數
//
// 範例:
//
//	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol {
//	    return modbus.NewTCPConnector()
//	})
func Register(protocolType schema.ProtocolType, factory ConnectorFactory) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.factories[protocolType] = factory
}

// Unregister 從全域註冊表移除協議連接器
func Unregister(protocolType schema.ProtocolType) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	delete(globalRegistry.factories, protocolType)
}

// Get 從全域註冊表取得並初始化連接器
//
// 參數:
//   - protocolType: 協議類型
//   - config: JSON 格式的連線配置
//
// 返回:
//   - Protocol: 已連線的協議連接器
//   - error: 錯誤 (若有)
func Get(ctx context.Context, protocolType schema.ProtocolType, config string) (Protocol, error) {
	return globalRegistry.Get(ctx, protocolType, config)
}

// GetFactory 取得協議連接器工廠 (不初始化連線)
func GetFactory(protocolType schema.ProtocolType) (ConnectorFactory, bool) {
	return globalRegistry.GetFactory(protocolType)
}

// ListProtocols 列出所有已註冊的協議類型
func ListProtocols() []schema.ProtocolType {
	return globalRegistry.ListProtocols()
}

// IsRegistered 檢查協議是否已註冊
func IsRegistered(protocolType schema.ProtocolType) bool {
	return globalRegistry.IsRegistered(protocolType)
}

// =============================================================================
// Registry 方法實作
// =============================================================================

// Get 從註冊表取得並初始化連接器
func (r *Registry) Get(ctx context.Context, protocolType schema.ProtocolType, config string) (Protocol, error) {
	factory, ok := r.GetFactory(protocolType)
	if !ok {
		return nil, fmt.Errorf("未註冊的協議類型: %s", protocolType)
	}

	// 建立連接器實例
	conn := factory()

	// 連線到設備
	if err := conn.Connect(ctx, config); err != nil {
		return nil, fmt.Errorf("連線失敗 (%s): %w", protocolType, err)
	}

	return conn, nil
}

// GetFactory 取得協議連接器工廠
func (r *Registry) GetFactory(protocolType schema.ProtocolType) (ConnectorFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, ok := r.factories[protocolType]
	return factory, ok
}

// ListProtocols 列出所有已註冊的協議類型
func (r *Registry) ListProtocols() []schema.ProtocolType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	protocols := make([]schema.ProtocolType, 0, len(r.factories))
	for p := range r.factories {
		protocols = append(protocols, p)
	}
	return protocols
}

// IsRegistered 檢查協議是否已註冊
func (r *Registry) IsRegistered(protocolType schema.ProtocolType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.factories[protocolType]
	return ok
}

// =============================================================================
// 協議資訊結構
// =============================================================================

// ProtocolInfo 協議資訊
type ProtocolInfo struct {
	// Type 協議類型
	Type schema.ProtocolType `json:"type"`

	// Name 協議顯示名稱
	Name string `json:"name"`

	// Description 協議描述
	Description string `json:"description"`

	// ConfigSchema 配置欄位 Schema (用於 UI 動態表單)
	ConfigSchema json.RawMessage `json:"config_schema"`
}

// protocolInfos 預定義的協議資訊
var protocolInfos = map[schema.ProtocolType]ProtocolInfo{
	schema.ProtocolModbusTCP: {
		Type:        schema.ProtocolModbusTCP,
		Name:        "Modbus TCP",
		Description: "Modbus TCP/IP 協議，用於乙太網路連接的 PLC 和設備",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"host": {"type": "string", "title": "主機位址", "description": "IP 位址或主機名稱"},
				"port": {"type": "integer", "title": "埠號", "default": 502, "minimum": 1, "maximum": 65535},
				"slave_id": {"type": "integer", "title": "從站 ID", "default": 1, "minimum": 1, "maximum": 247},
				"timeout": {"type": "integer", "title": "逾時 (秒)", "default": 5, "minimum": 1, "maximum": 60},
				"data_format": {"type": "string", "title": "資料格式（字節序）", "enum": ["ABCD", "BADC", "CDAB", "DCBA"], "default": "ABCD", "description": "對應 hsllogic 浮點／多暫存器解碼字節序"}
			},
			"required": ["host", "port", "slave_id"]
		}`),
	},
	schema.ProtocolModbusRTU: {
		Type:        schema.ProtocolModbusRTU,
		Name:        "Modbus RTU",
		Description: "Modbus RTU 協議，用於 RS-232/RS-485 串列通訊",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"serial_port": {"type": "string", "title": "串列埠", "description": "如 COM1 或 /dev/ttyUSB0"},
				"baud_rate": {"type": "integer", "title": "鮑率", "enum": [9600, 19200, 38400, 57600, 115200], "default": 9600},
				"data_bits": {"type": "integer", "title": "資料位元", "enum": [7, 8], "default": 8},
				"stop_bits": {"type": "integer", "title": "停止位元", "enum": [1, 2], "default": 1},
				"parity": {"type": "string", "title": "同位檢查", "enum": ["none", "odd", "even"], "default": "none"},
				"slave_id": {"type": "integer", "title": "從站 ID", "default": 1, "minimum": 1, "maximum": 247},
				"timeout": {"type": "integer", "title": "逾時 (秒)", "default": 5, "minimum": 1, "maximum": 60},
				"data_format": {"type": "string", "title": "資料格式（字節序）", "enum": ["ABCD", "BADC", "CDAB", "DCBA"], "default": "ABCD", "description": "對應 hsllogic 浮點／多暫存器解碼字節序"}
			},
			"required": ["serial_port", "baud_rate", "slave_id"]
		}`),
	},
	schema.ProtocolModbusUDP: {
		Type:        schema.ProtocolModbusUDP,
		Name:        "Modbus UDP",
		Description: "Modbus UDP 協議，用於 UDP 通訊",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"host": {"type": "string", "title": "主機位址"},
				"port": {"type": "integer", "title": "埠號", "default": 502},
				"slave_id": {"type": "integer", "title": "從站 ID", "default": 1},
				"timeout": {"type": "integer", "title": "逾時 (秒)", "default": 5},
				"data_format": {"type": "string", "title": "資料格式（字節序）", "enum": ["ABCD", "BADC", "CDAB", "DCBA"], "default": "ABCD", "description": "對應 hsllogic 浮點／多暫存器解碼字節序"}
			},
			"required": ["host", "port", "slave_id"]
		}`),
	},
	schema.ProtocolFatekFBs: {
		Type:        schema.ProtocolFatekFBs,
		Name:        "FATEK FBs",
		Description: "FATEK FBs 系列 PLC ASCII 協議",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"mode": {"type": "string", "title": "連線模式", "enum": ["tcp", "serial"], "default": "tcp"},
				"host": {"type": "string", "title": "主機位址", "description": "TCP 模式使用"},
				"port": {"type": "integer", "title": "埠號", "default": 500, "description": "TCP 模式使用"},
				"serial_port": {"type": "string", "title": "串列埠", "description": "Serial 模式使用"},
				"baud_rate": {"type": "integer", "title": "鮑率", "default": 9600, "description": "Serial 模式使用"},
				"data_bits": {"type": "integer", "title": "資料位元", "enum": [7, 8], "default": 7, "description": "Serial 模式使用"},
				"stop_bits": {"type": "integer", "title": "停止位元", "enum": [1, 2], "default": 1, "description": "Serial 模式使用"},
				"parity": {"type": "string", "title": "同位檢查", "enum": ["none", "odd", "even"], "default": "even", "description": "Serial 模式使用"},
				"station_no": {"type": "integer", "title": "站號", "default": 1, "minimum": 0, "maximum": 255},
				"timeout": {"type": "integer", "title": "逾時 (秒)", "default": 5}
			},
			"required": ["mode", "station_no"]
		}`),
	},
	schema.ProtocolMC3E: {
		Type:        schema.ProtocolMC3E,
		Name:        "Mitsubishi MC 3E",
		Description: "三菱 MC Protocol 3E Frame (Binary)",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"host": {"type": "string", "title": "主機位址"},
				"port": {"type": "integer", "title": "埠號", "default": 5000},
				"network_no": {"type": "integer", "title": "網路編號", "default": 0},
				"pc_no": {"type": "integer", "title": "PC 編號", "default": 255},
				"io_no": {"type": "integer", "title": "I/O 編號", "default": 1023},
				"station_no": {"type": "integer", "title": "站號", "default": 0},
				"timeout": {"type": "integer", "title": "逾時 (秒)", "default": 5},
				"data_format": {"type": "string", "title": "資料格式（字節序）", "enum": ["ABCD", "BADC", "CDAB", "DCBA"], "default": "CDAB", "description": "對應 hsllogic 浮點／多暫存器解碼字節序"}
			},
			"required": ["host", "port"]
		}`),
	},
	schema.ProtocolMQTT: {
		Type:        schema.ProtocolMQTT,
		Name:        "MQTT",
		Description: "MQTT 訂閱接收，用於接收 IoT 設備資料",
		ConfigSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"broker_url": {"type": "string", "title": "Broker URL", "description": "如 tcp://localhost:1883"},
				"client_id": {"type": "string", "title": "Client ID"},
				"username": {"type": "string", "title": "使用者名稱"},
				"password": {"type": "string", "title": "密碼", "format": "password"},
				"use_tls": {"type": "boolean", "title": "使用 TLS", "default": false},
				"topics": {"type": "array", "title": "訂閱主題", "items": {"type": "string"}},
				"qos": {"type": "integer", "title": "QoS", "enum": [0, 1, 2], "default": 0}
			},
			"required": ["broker_url", "client_id", "topics"]
		}`),
	},
}

// GetProtocolInfo 取得協議資訊
func GetProtocolInfo(protocolType schema.ProtocolType) (ProtocolInfo, bool) {
	info, ok := protocolInfos[protocolType]
	return info, ok
}

// ListProtocolInfos 列出所有協議資訊
func ListProtocolInfos() []ProtocolInfo {
	infos := make([]ProtocolInfo, 0, len(protocolInfos))
	for _, info := range protocolInfos {
		infos = append(infos, info)
	}
	return infos
}
