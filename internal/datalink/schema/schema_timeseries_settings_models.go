package schema

import "time"

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
	// QualityMissing 資料缺失 (如通訊逾時或無讀值)
	QualityMissing QualityFlag = "missing"
	// QualityStale 資料陳舊 (未更新)
	QualityStale QualityFlag = "stale"
	// QualityInvalid 資料不合法 (如 NaN/Inf 或溢位)
	QualityInvalid QualityFlag = "invalid"
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
