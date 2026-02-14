// Package storage 提供時序資料儲存功能。
//
// 本套件實作時序資料的批次寫入、快取和查詢功能。
package storage

import (
	"context"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 時序記錄
// =============================================================================

// TimeSeriesRecord 時序資料記錄
type TimeSeriesRecord struct {
	TagID     string             `json:"tag_id"`
	Timestamp time.Time          `json:"ts"`
	ValueNum  *float64           `json:"value_num,omitempty"`
	ValueText *string            `json:"value_text,omitempty"`
	ValueBool *bool              `json:"value_bool,omitempty"`
	RawValue  interface{}        `json:"raw_value,omitempty"`
	Quality   schema.QualityFlag `json:"quality"`
}

// =============================================================================
// Writer 介面定義
// =============================================================================

// Writer 時序資料寫入器介面
type Writer interface {
	// Write 寫入單筆記錄
	Write(ctx context.Context, record TimeSeriesRecord) error

	// WriteBatch 批次寫入記錄
	WriteBatch(ctx context.Context, records []TimeSeriesRecord) error

	// Flush 強制刷新緩衝區
	Flush(ctx context.Context) error

	// Close 關閉寫入器
	Close() error
}

// Reader 時序資料讀取器介面
type Reader interface {
	// Query 查詢時序資料
	Query(ctx context.Context, query TimeSeriesQuery) ([]TimeSeriesRecord, error)

	// GetLatest 取得標籤的最新值
	GetLatest(ctx context.Context, tagID string) (*TimeSeriesRecord, error)

	// GetLatestBatch 批次取得多個標籤的最新值
	GetLatestBatch(ctx context.Context, tagIDs []string) (map[string]*TimeSeriesRecord, error)
}

// TimeSeriesQuery 時序查詢條件
type TimeSeriesQuery struct {
	// TagID 標籤 ID (必填)
	TagID string

	// StartTime 開始時間 (含)
	StartTime *time.Time

	// EndTime 結束時間 (含)
	EndTime *time.Time

	// Quality 品質篩選
	Quality *schema.QualityFlag

	// Limit 限制數量
	Limit int

	// Order 排序 ("asc" 或 "desc")
	Order string
}

// =============================================================================
// 批次寫入器
// =============================================================================

// BatchWriterConfig 批次寫入器配置
type BatchWriterConfig struct {
	// BatchSize 批次大小
	BatchSize int

	// FlushInterval 刷新間隔
	FlushInterval time.Duration

	// WriteTimeout 寫入逾時
	WriteTimeout time.Duration
}

// DefaultBatchWriterConfig 預設批次寫入器配置
func DefaultBatchWriterConfig() BatchWriterConfig {
	return BatchWriterConfig{
		BatchSize:     1000,
		FlushInterval: 5 * time.Second,
		WriteTimeout:  10 * time.Second,
	}
}
