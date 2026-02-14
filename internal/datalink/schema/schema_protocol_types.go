package schema

// =============================================================================
// 時間精度相關常數
// =============================================================================

// TimePrecision 定義時序資料的時間精度
type TimePrecision string

const (
	// TimePrecisionSecond 秒級精度
	TimePrecisionSecond TimePrecision = "second"
	// TimePrecisionMillisecond 毫秒級精度
	TimePrecisionMillisecond TimePrecision = "millisecond"
)

// PartitionInterval 定義 Postgres 時序表的分區間隔
type PartitionInterval string

const (
	// PartitionDaily 每日分區
	PartitionDaily PartitionInterval = "daily"
	// PartitionWeekly 每週分區
	PartitionWeekly PartitionInterval = "weekly"
	// PartitionMonthly 每月分區 (預設)
	PartitionMonthly PartitionInterval = "monthly"
)

// =============================================================================
// 設備狀態與協議相關常數
// =============================================================================

// DeviceStatus 設備生命週期狀態
type DeviceStatus string

const (
	// DeviceStatusDraft 草稿狀態 (未啟用)
	DeviceStatusDraft DeviceStatus = "draft"
	// DeviceStatusActive 啟用狀態 (參與調度)
	DeviceStatusActive DeviceStatus = "active"
	// DeviceStatusDisabled 停用狀態 (暫停調度)
	DeviceStatusDisabled DeviceStatus = "disabled"
)

// ProtocolType 協議類型
type ProtocolType string

const (
	// ProtocolModbusTCP Modbus TCP 協議
	ProtocolModbusTCP ProtocolType = "modbus_tcp"
	// ProtocolModbusRTU Modbus RTU 協議
	ProtocolModbusRTU ProtocolType = "modbus_rtu"
	// ProtocolModbusUDP Modbus UDP 協議
	ProtocolModbusUDP ProtocolType = "modbus_udp"
	// ProtocolFatekFBs FATEK FBs 協議
	ProtocolFatekFBs ProtocolType = "fatek_fbs"
	// ProtocolMC3E Mitsubishi MC 3E 協議
	ProtocolMC3E ProtocolType = "mc_3e"
	// ProtocolMQTT MQTT 接收
	ProtocolMQTT ProtocolType = "mqtt"
)
