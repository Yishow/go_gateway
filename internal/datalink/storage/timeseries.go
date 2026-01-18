// Package storage 提供時序資料儲存功能。
//
// 本套件實作時序資料的批次寫入、快取和查詢功能。
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
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

// BatchWriter 批次寫入器
type BatchWriter struct {
	config     BatchWriterConfig
	underlying Writer
	buffer     []TimeSeriesRecord
	mu         sync.Mutex
	flushTimer *time.Timer
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewBatchWriter 建立新的批次寫入器
func NewBatchWriter(underlying Writer, config BatchWriterConfig) *BatchWriter {
	bw := &BatchWriter{
		config:     config,
		underlying: underlying,
		buffer:     make([]TimeSeriesRecord, 0, config.BatchSize),
		stopCh:     make(chan struct{}),
	}

	// 啟動定時刷新
	bw.startFlushTimer()

	return bw
}

// Write 寫入單筆記錄
func (bw *BatchWriter) Write(ctx context.Context, record TimeSeriesRecord) error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.buffer = append(bw.buffer, record)

	// 達到批次大小時刷新
	if len(bw.buffer) >= bw.config.BatchSize {
		return bw.flushLocked(ctx)
	}

	return nil
}

// WriteBatch 批次寫入記錄
func (bw *BatchWriter) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	bw.buffer = append(bw.buffer, records...)

	// 達到批次大小時刷新
	if len(bw.buffer) >= bw.config.BatchSize {
		return bw.flushLocked(ctx)
	}

	return nil
}

// Flush 強制刷新緩衝區
func (bw *BatchWriter) Flush(ctx context.Context) error {
	bw.mu.Lock()
	defer bw.mu.Unlock()

	return bw.flushLocked(ctx)
}

// flushLocked 刷新緩衝區 (必須持有鎖)
func (bw *BatchWriter) flushLocked(ctx context.Context) error {
	if len(bw.buffer) == 0 {
		return nil
	}

	// 複製緩衝區
	toWrite := make([]TimeSeriesRecord, len(bw.buffer))
	copy(toWrite, bw.buffer)

	// 清空緩衝區
	bw.buffer = bw.buffer[:0]

	// 寫入底層
	return bw.underlying.WriteBatch(ctx, toWrite)
}

// startFlushTimer 啟動定時刷新
func (bw *BatchWriter) startFlushTimer() {
	bw.flushTimer = time.NewTimer(bw.config.FlushInterval)

	bw.wg.Add(1)
	go func() {
		defer bw.wg.Done()
		for {
			select {
			case <-bw.stopCh:
				bw.flushTimer.Stop()
				return
			case <-bw.flushTimer.C:
				ctx, cancel := context.WithTimeout(context.Background(), bw.config.WriteTimeout)
				_ = bw.Flush(ctx)
				cancel()
				bw.flushTimer.Reset(bw.config.FlushInterval)
			}
		}
	}()
}

// Close 關閉寫入器
func (bw *BatchWriter) Close() error {
	close(bw.stopCh)
	bw.wg.Wait()

	// 最後刷新
	ctx, cancel := context.WithTimeout(context.Background(), bw.config.WriteTimeout)
	defer cancel()
	_ = bw.Flush(ctx)

	return bw.underlying.Close()
}

// BufferSize 取得目前緩衝區大小
func (bw *BatchWriter) BufferSize() int {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	return len(bw.buffer)
}

// =============================================================================
// 記憶體儲存 (用於測試)
// =============================================================================

// MemoryStorage 記憶體時序儲存
type MemoryStorage struct {
	mu      sync.RWMutex
	records []TimeSeriesRecord
	maxSize int
}

// NewMemoryStorage 建立新的記憶體儲存
func NewMemoryStorage(maxSize int) *MemoryStorage {
	if maxSize <= 0 {
		maxSize = 100000
	}
	return &MemoryStorage{
		records: make([]TimeSeriesRecord, 0),
		maxSize: maxSize,
	}
}

// Write 寫入單筆記錄
func (ms *MemoryStorage) Write(ctx context.Context, record TimeSeriesRecord) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.records = append(ms.records, record)

	// 超過最大大小時移除舊資料
	if len(ms.records) > ms.maxSize {
		ms.records = ms.records[len(ms.records)-ms.maxSize:]
	}

	return nil
}

// WriteBatch 批次寫入記錄
func (ms *MemoryStorage) WriteBatch(ctx context.Context, records []TimeSeriesRecord) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.records = append(ms.records, records...)

	// 超過最大大小時移除舊資料
	if len(ms.records) > ms.maxSize {
		ms.records = ms.records[len(ms.records)-ms.maxSize:]
	}

	return nil
}

// Flush 刷新 (記憶體儲存無需操作)
func (ms *MemoryStorage) Flush(ctx context.Context) error {
	return nil
}

// Close 關閉 (記憶體儲存無需操作)
func (ms *MemoryStorage) Close() error {
	return nil
}

// Query 查詢時序資料
func (ms *MemoryStorage) Query(ctx context.Context, query TimeSeriesQuery) ([]TimeSeriesRecord, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result := make([]TimeSeriesRecord, 0)

	for _, record := range ms.records {
		// 篩選 TagID
		if record.TagID != query.TagID {
			continue
		}

		// 篩選時間範圍
		if query.StartTime != nil && record.Timestamp.Before(*query.StartTime) {
			continue
		}
		if query.EndTime != nil && record.Timestamp.After(*query.EndTime) {
			continue
		}

		// 篩選品質
		if query.Quality != nil && record.Quality != *query.Quality {
			continue
		}

		result = append(result, record)
	}

	// 排序
	if query.Order == "desc" {
		// 反轉切片
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
	}

	// 限制數量
	if query.Limit > 0 && len(result) > query.Limit {
		result = result[:query.Limit]
	}

	return result, nil
}

// GetLatest 取得標籤的最新值
func (ms *MemoryStorage) GetLatest(ctx context.Context, tagID string) (*TimeSeriesRecord, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var latest *TimeSeriesRecord

	for i := len(ms.records) - 1; i >= 0; i-- {
		if ms.records[i].TagID == tagID {
			record := ms.records[i]
			latest = &record
			break
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("找不到標籤 %s 的資料", tagID)
	}

	return latest, nil
}

// GetLatestBatch 批次取得多個標籤的最新值
func (ms *MemoryStorage) GetLatestBatch(ctx context.Context, tagIDs []string) (map[string]*TimeSeriesRecord, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result := make(map[string]*TimeSeriesRecord)
	tagSet := make(map[string]bool)
	for _, id := range tagIDs {
		tagSet[id] = true
	}

	// 從後往前遍歷找最新值
	for i := len(ms.records) - 1; i >= 0 && len(tagSet) > 0; i-- {
		record := ms.records[i]
		if tagSet[record.TagID] {
			r := record
			result[record.TagID] = &r
			delete(tagSet, record.TagID)
		}
	}

	return result, nil
}

// Count 計算記錄數量
func (ms *MemoryStorage) Count() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.records)
}

// Clear 清空所有記錄
func (ms *MemoryStorage) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.records = make([]TimeSeriesRecord, 0)
}

// =============================================================================
// 值轉換輔助
// =============================================================================

// ValueToRecord 將收集值轉換為時序記錄
func ValueToRecord(tagID string, value interface{}, rawValue interface{}, timestamp time.Time, quality schema.QualityFlag, dataType schema.DataType) TimeSeriesRecord {
	record := TimeSeriesRecord{
		TagID:     tagID,
		Timestamp: timestamp,
		Quality:   quality,
	}

	// 設定原始值
	if rawValue != nil {
		record.RawValue = rawValue
	}

	// 根據資料型別設定值
	switch dataType {
	case schema.DataTypeBool:
		if v, ok := value.(bool); ok {
			record.ValueBool = &v
		}
	case schema.DataTypeString:
		if v, ok := value.(string); ok {
			record.ValueText = &v
		} else {
			s := fmt.Sprintf("%v", value)
			record.ValueText = &s
		}
	default:
		// 數值型別
		if v, ok := toFloat64(value); ok {
			record.ValueNum = &v
		}
	}

	return record
}

// toFloat64 將值轉換為 float64
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return f, true
		}
	}
	return 0, false
}
