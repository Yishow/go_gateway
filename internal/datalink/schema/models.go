// Package schema 定義設備資料收集管線的資料庫結構。
//
// 本套件提供完整的資料庫模型定義，支援：
//   - PostgreSQL (生產環境，支援表分區)
//   - SQLite (測試環境，不支援分區)
package schema

import (
	"time"
)

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

// =============================================================================
// 設備模型 (Device)
// =============================================================================

// Device 設備實體，代表可連接的 PLC 或資料來源
type Device struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// Name 設備顯示名稱
	Name string `json:"name" db:"name"`

	// Description 設備描述
	Description string `json:"description,omitempty" db:"description"`

	// Protocol 通訊協議類型
	Protocol ProtocolType `json:"protocol" db:"protocol"`

	// Status 設備狀態
	Status DeviceStatus `json:"status" db:"status"`

	// ConnectionConfig 協議特定的連線配置 (JSON 格式)
	// 不同協議有不同的欄位結構，詳見 ConnectionConfigXXX 類型
	ConnectionConfig string `json:"connection_config" db:"connection_config"`

	// LastTestAt 最後連線測試時間
	LastTestAt *time.Time `json:"last_test_at,omitempty" db:"last_test_at"`

	// LastTestSuccess 最後連線測試是否成功
	LastTestSuccess *bool `json:"last_test_success,omitempty" db:"last_test_success"`

	// LastTestError 最後連線測試錯誤訊息
	LastTestError string `json:"last_test_error,omitempty" db:"last_test_error"`

	// LastCollectedAt 最後資料收集時間 (新增欄位)
	LastCollectedAt *time.Time `json:"last_collected_at,omitempty" db:"last_collected_at"`

	// CollectionCount 總收集次數 (新增欄位)
	CollectionCount int64 `json:"collection_count" db:"collection_count"`

	// ErrorCount 總錯誤次數 (新增欄位)
	ErrorCount int64 `json:"error_count" db:"error_count"`

	// ReadinessStatus 設備就緒狀態 (新增欄位)
	// JSON 儲存: {"status": "ready|warning|error", "details": [...], "missing_steps": [...]}
	ReadinessStatus string `json:"readiness_status,omitempty" db:"readiness_status"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ReadinessCheck 單項檢查結果
type ReadinessCheck struct {
	Name    string `json:"name"`
	Pass    bool   `json:"pass"`
	Message string `json:"message"`
}

// DeviceReadiness 設備就緒狀態詳情
type DeviceReadiness struct {
	DeviceID string           `json:"device_id"`
	Status   string           `json:"status"` // ready, warning, error
	Checks   []ReadinessCheck `json:"checks"`
}

// =============================================================================
// 連線配置結構 (協議特定)
// =============================================================================

// ConnectionConfigModbusTCP Modbus TCP 連線配置
type ConnectionConfigModbusTCP struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port TCP 埠號 (預設 502)
	Port int `json:"port"`
	// SlaveID 從站 ID (1-247)
	SlaveID byte `json:"slave_id"`
	// Timeout 連線逾時秒數
	Timeout int `json:"timeout"`
}

// ConnectionConfigModbusRTU Modbus RTU 連線配置
type ConnectionConfigModbusRTU struct {
	// SerialPort 串列埠名稱 (如 COM1, /dev/ttyUSB0)
	SerialPort string `json:"serial_port"`
	// BaudRate 鮑率 (9600, 19200, 38400, 57600, 115200)
	BaudRate int `json:"baud_rate"`
	// DataBits 資料位元 (7 或 8)
	DataBits int `json:"data_bits"`
	// StopBits 停止位元 (1 或 2)
	StopBits int `json:"stop_bits"`
	// Parity 同位檢查 (none, odd, even)
	Parity string `json:"parity"`
	// SlaveID 從站 ID (1-247)
	SlaveID byte `json:"slave_id"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
}

// ConnectionConfigModbusUDP Modbus UDP 連線配置
type ConnectionConfigModbusUDP struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port UDP 埠號
	Port int `json:"port"`
	// SlaveID 從站 ID
	SlaveID byte `json:"slave_id"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
}

// ConnectionConfigFatekFBs FATEK FBs 連線配置
type ConnectionConfigFatekFBs struct {
	// Mode 連線模式 (tcp 或 serial)
	Mode string `json:"mode"`
	// Host TCP 模式的主機位址
	Host string `json:"host,omitempty"`
	// Port TCP 模式的埠號 (預設 500)
	Port int `json:"port,omitempty"`
	// SerialPort 串列模式的埠名稱
	SerialPort string `json:"serial_port,omitempty"`
	// BaudRate 串列模式的鮑率
	BaudRate int `json:"baud_rate,omitempty"`
	// StationNo 站號 (0-255)
	StationNo byte `json:"station_no"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
}

// ConnectionConfigMC3E Mitsubishi MC 3E 連線配置
type ConnectionConfigMC3E struct {
	// Host 目標主機位址
	Host string `json:"host"`
	// Port TCP 埠號 (預設依 PLC 型號)
	Port int `json:"port"`
	// NetworkNo 網路編號
	NetworkNo byte `json:"network_no"`
	// PCNo PC 編號
	PCNo byte `json:"pc_no"`
	// IONo I/O 編號
	IONo uint16 `json:"io_no"`
	// StationNo 站號
	StationNo byte `json:"station_no"`
	// Timeout 通訊逾時秒數
	Timeout int `json:"timeout"`
	// DataFormat 字節序格式 (ABCD, BADC, CDAB, DCBA)，預設為 CDAB (三菱標準)
	DataFormat string `json:"data_format,omitempty"`
}

// ConnectionConfigMQTT MQTT 接收連線配置
type ConnectionConfigMQTT struct {
	// BrokerURL MQTT Broker 位址 (如 tcp://localhost:1883)
	BrokerURL string `json:"broker_url"`
	// ClientID 客戶端識別碼
	ClientID string `json:"client_id"`
	// Username 使用者名稱 (選用)
	Username string `json:"username,omitempty"`
	// Password 密碼 (選用)
	Password string `json:"password,omitempty"`
	// UseTLS 是否使用 TLS
	UseTLS bool `json:"use_tls"`
	// Topics 訂閱的主題清單
	Topics []string `json:"topics"`
	// QoS 服務品質等級 (0, 1, 2)
	QoS byte `json:"qos"`
}

// =============================================================================
// 點位模型 (Point)
// =============================================================================

// PointMode 點位讀寫模式
type PointMode string

const (
	// PointModeReadOnly 唯讀
	PointModeReadOnly PointMode = "read"
	// PointModeReadWrite 讀寫
	PointModeReadWrite PointMode = "readwrite"
)

// DataType 資料型別
type DataType string

const (
	// DataTypeBool 布林值
	DataTypeBool DataType = "bool"
	// DataTypeInt16 16 位元有號整數
	DataTypeInt16 DataType = "int16"
	// DataTypeUint16 16 位元無號整數
	DataTypeUint16 DataType = "uint16"
	// DataTypeInt32 32 位元有號整數
	DataTypeInt32 DataType = "int32"
	// DataTypeUint32 32 位元無號整數
	DataTypeUint32 DataType = "uint32"
	// DataTypeInt64 64 位元有號整數
	DataTypeInt64 DataType = "int64"
	// DataTypeUint64 64 位元無號整數
	DataTypeUint64 DataType = "uint64"
	// DataTypeFloat32 32 位元浮點數
	DataTypeFloat32 DataType = "float32"
	// DataTypeFloat64 64 位元浮點數
	DataTypeFloat64 DataType = "float64"
	// DataTypeString 字串
	DataTypeString DataType = "string"
)

// RegisterCountForDataType returns word count for 16-bit word protocols (Modbus/FATEK).
func RegisterCountForDataType(dataType DataType) int {
	switch dataType {
	case DataTypeBool, DataTypeInt16, DataTypeUint16:
		return 1
	case DataTypeInt32, DataTypeUint32, DataTypeFloat32:
		return 2
	case DataTypeInt64, DataTypeUint64, DataTypeFloat64:
		return 4
	default:
		return 1
	}
}

// Point 點位實體，代表設備上的可讀取/寫入位址
type Point struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// DeviceID 所屬設備 ID (外鍵)
	DeviceID string `json:"device_id" db:"device_id"`

	// Name 點位顯示名稱
	Name string `json:"name" db:"name"`

	// Description 點位描述
	Description string `json:"description,omitempty" db:"description"`

	// Address 協議特定的位址字串 (如 Modbus: "40001", FATEK: "D0100")
	Address string `json:"address" db:"address"`

	// Function 協議功能碼或存取方式 (如 Modbus FC03)
	Function string `json:"function,omitempty" db:"function"`

	// DataType 資料型別
	DataType DataType `json:"data_type" db:"data_type"`

	// Mode 讀寫模式
	Mode PointMode `json:"mode" db:"mode"`

	// PollingGroupID 所屬輪詢群組 ID (外鍵，可為空)
	PollingGroupID *string `json:"polling_group_id,omitempty" db:"polling_group_id"`

	// LastReadAt 最後讀取時間
	LastReadAt *time.Time `json:"last_read_at,omitempty" db:"last_read_at"`

	// LastValue 最後讀取的值 (JSON 格式儲存)
	LastValue *string `json:"last_value,omitempty" db:"last_value"`

	// LastError 最後讀取錯誤訊息
	LastError string `json:"last_error,omitempty" db:"last_error"`

	// Enabled 是否啟用
	Enabled bool `json:"enabled" db:"enabled"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 輪詢群組模型 (PollingGroup)
// =============================================================================

// PollingGroup 輪詢群組，用於控制點位的採集間隔
type PollingGroup struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// Name 群組名稱
	Name string `json:"name" db:"name"`

	// Description 群組描述
	Description string `json:"description,omitempty" db:"description"`

	// IntervalMs 輪詢間隔 (毫秒)
	IntervalMs int `json:"interval_ms" db:"interval_ms"`

	// Priority 優先級 (數字越小優先級越高)
	Priority int `json:"priority" db:"priority"`

	// Enabled 是否啟用
	Enabled bool `json:"enabled" db:"enabled"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 標籤模型 (Tag)
// =============================================================================

// TagStatus 標籤生命週期狀態
type TagStatus string

const (
	// TagStatusDraft 草稿狀態
	TagStatusDraft TagStatus = "draft"
	// TagStatusActive 啟用狀態
	TagStatusActive TagStatus = "active"
	// TagStatusRetired 退役狀態 (不允許新映射)
	TagStatusRetired TagStatus = "retired"
)

// Tag 全域標籤實體
type Tag struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// Key 標籤鍵 (唯一，大小寫不敏感)
	// 格式: ASCII 字母、數字、底線、連字號、點、斜線
	// 長度: 1-128 字元
	Key string `json:"key" db:"key"`

	// KeyLower 標籤鍵小寫版本 (用於唯一性檢查)
	KeyLower string `json:"-" db:"key_lower"`

	// DisplayName 顯示名稱
	DisplayName string `json:"display_name" db:"display_name"`

	// Description 標籤描述
	Description string `json:"description,omitempty" db:"description"`

	// Unit 單位 (如 ℃, kWh, rpm)
	Unit string `json:"unit,omitempty" db:"unit"`

	// DataType 資料型別
	DataType DataType `json:"data_type" db:"data_type"`

	// Status 標籤狀態
	Status TagStatus `json:"status" db:"status"`

	// Labels 標籤屬性 (JSON 物件，用於分類和篩選)
	Labels string `json:"labels,omitempty" db:"labels"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 映射模型 (Mapping)
// =============================================================================

// Mapping 點位到標籤的映射關係
type Mapping struct {
	// ID 主鍵，UUID 格式
	ID string `json:"id" db:"id"`

	// PointID 來源點位 ID (外鍵)
	PointID string `json:"point_id" db:"point_id"`

	// TagID 目標標籤 ID (外鍵)
	TagID string `json:"tag_id" db:"tag_id"`

	// TransformPipeline 轉換管線步驟 (JSON 陣列)
	// 詳見 TransformStep 類型定義
	TransformPipeline string `json:"transform_pipeline" db:"transform_pipeline"`

	// Enabled 是否啟用
	Enabled bool `json:"enabled" db:"enabled"`

	// CreatedAt 建立時間
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// =============================================================================
// 轉換步驟模型 (TransformStep)
// =============================================================================

// TransformType 轉換步驟類型
type TransformType string

const (
	// TransformDecode 解碼 (如 BCD, Gray code)
	TransformDecode TransformType = "decode"
	// TransformCast 型別轉換
	TransformCast TransformType = "cast"
	// TransformScale 線性縮放 (value * scale + offset)
	TransformScale TransformType = "scale"
	// TransformLookup 查表替換
	TransformLookup TransformType = "lookup"
	// TransformConditional 條件判斷
	TransformConditional TransformType = "conditional"
	// TransformFormula 公式表達式
	TransformFormula TransformType = "formula"
)

// TransformStep 單一轉換步驟的配置
type TransformStep struct {
	// Type 轉換類型
	Type TransformType `json:"type"`

	// Order 執行順序 (從小到大)
	Order int `json:"order"`

	// Params 轉換參數 (依類型不同)
	Params map[string]interface{} `json:"params"`
}

// TransformParamsScale 線性縮放參數
type TransformParamsScale struct {
	// Scale 縮放係數
	Scale float64 `json:"scale"`
	// Offset 偏移量
	Offset float64 `json:"offset"`
}

// TransformParamsCast 型別轉換參數
type TransformParamsCast struct {
	// TargetType 目標型別
	TargetType DataType `json:"target_type"`
}

// TransformParamsLookup 查表參數
type TransformParamsLookup struct {
	// Table 對照表 (key-value 映射)
	Table map[string]interface{} `json:"table"`
	// Default 無匹配時的預設值
	Default interface{} `json:"default,omitempty"`
}

// TransformParamsConditional 條件判斷參數
type TransformParamsConditional struct {
	// Condition 條件表達式
	Condition string `json:"condition"`
	// TrueValue 條件為真時的值或表達式
	TrueValue interface{} `json:"true_value"`
	// FalseValue 條件為假時的值或表達式
	FalseValue interface{} `json:"false_value"`
}

// TransformParamsFormula 公式表達式參數
type TransformParamsFormula struct {
	// Expression 表達式字串
	Expression string `json:"expression"`
}

// TransformParamsDecode 解碼參數
type TransformParamsDecode struct {
	// Format 解碼格式 (bcd, gray, hex_string, etc.)
	Format string `json:"format"`
}

// =============================================================================
// 時序資料模型 (TimeSeriesRecord)
// =============================================================================

// QualityFlag 資料品質標誌
type QualityFlag string

const (
	// QualityGood 資料品質良好
	QualityGood QualityFlag = "good"
	// QualityBad 資料品質不佳 (讀取失敗)
	QualityBad QualityFlag = "bad"
	// QualityUncertain 資料品質不確定 (如逾時後使用舊值)
	QualityUncertain QualityFlag = "uncertain"
)

// TimeSeriesRecord 時序資料記錄
type TimeSeriesRecord struct {
	// ID 主鍵 (自增或 UUID，依資料庫)
	ID int64 `json:"id" db:"id"`

	// TagID 標籤 ID (外鍵)
	TagID string `json:"tag_id" db:"tag_id"`

	// Timestamp 時間戳記
	Timestamp time.Time `json:"timestamp" db:"ts"`

	// ValueNum 數值型態值 (互斥，依 DataType 選擇)
	ValueNum *float64 `json:"value_num,omitempty" db:"value_num"`

	// ValueText 文字型態值
	ValueText *string `json:"value_text,omitempty" db:"value_text"`

	// ValueBool 布林型態值
	ValueBool *bool `json:"value_bool,omitempty" db:"value_bool"`

	// RawValue 原始值 (JSON 格式，用於追溯)
	RawValue string `json:"raw_value,omitempty" db:"raw_value"`

	// Quality 資料品質標誌
	Quality QualityFlag `json:"quality" db:"quality"`
}

// =============================================================================
// 系統設定模型 (SystemSettings)
// =============================================================================

// SystemSettings 系統設定
type SystemSettings struct {
	// ID 主鍵
	ID string `json:"id" db:"id"`

	// Key 設定鍵
	Key string `json:"key" db:"key"`

	// Value 設定值 (JSON 格式)
	Value string `json:"value" db:"value"`

	// Description 設定描述
	Description string `json:"description,omitempty" db:"description"`

	// UpdatedAt 更新時間
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// 系統設定鍵常數
const (
	// SettingWritePrecision 寫入時間精度
	SettingWritePrecision = "write_precision"
	// SettingPartitionInterval 分區間隔
	SettingPartitionInterval = "partition_interval"
	// SettingBatchSize 批次寫入大小
	SettingBatchSize = "batch_size"
	// SettingFlushInterval 批次刷新間隔
	SettingFlushInterval = "flush_interval"
	// SettingDefaultRetryCount 預設重試次數
	SettingDefaultRetryCount = "default_retry_count"
	// SettingDefaultRetryDelay 預設重試延遲
	SettingDefaultRetryDelay = "default_retry_delay"
)
