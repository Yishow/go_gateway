package storage

import (
	"context"
	"fmt"
	"sync"
)

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
