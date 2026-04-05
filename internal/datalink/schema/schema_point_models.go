package schema

import "time"

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

	// DataFormat Modbus 等多暫存器型別的字節序（ABCD/BADC/CDAB/DCBA）；空字串表示與歷史行為相同（等同 ABCD）。
	DataFormat string `json:"data_format,omitempty" db:"data_format"`

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
