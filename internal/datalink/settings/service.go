// Package settings 提供系統設定管理功能。
//
// 本套件實作 Datalink 系統設定的讀取與更新。
package settings

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// 設定鍵常數
// =============================================================================

const (
	// KeyWritePrecision 寫入時間精度
	KeyWritePrecision = "write_precision"
	// KeyPartitionInterval 分區間隔
	KeyPartitionInterval = "partition_interval"
	// KeyBatchSize 批次大小
	KeyBatchSize = "batch_size"
	// KeyDefaultRetryCount 預設重試次數
	KeyDefaultRetryCount = "default_retry_count"
	// KeyDefaultRetryDelay 預設重試延遲
	KeyDefaultRetryDelay = "default_retry_delay"
)

// =============================================================================
// 設定項目
// =============================================================================

// SettingItem 設定項目
type SettingItem struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 設定資料存取介面
type Repository interface {
	// Get 取得設定值
	Get(ctx context.Context, key string) (*SettingItem, error)

	// Set 設定值
	Set(ctx context.Context, key string, value interface{}) error

	// List 列出所有設定
	List(ctx context.Context) ([]*SettingItem, error)
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 設定管理服務
type Service struct {
	repo Repository
	mu   sync.RWMutex
}

// NewService 建立新的設定服務
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Get 取得設定值
func (s *Service) Get(ctx context.Context, key string) (*SettingItem, error) {
	return s.repo.Get(ctx, key)
}

// Set 設定值
func (s *Service) Set(ctx context.Context, key string, value interface{}) error {
	return s.repo.Set(ctx, key, value)
}

// List 列出所有設定
func (s *Service) List(ctx context.Context) ([]*SettingItem, error) {
	return s.repo.List(ctx)
}

// =============================================================================
// 記憶體 Repository
// =============================================================================

// MemoryRepository 記憶體內設定儲存庫
type MemoryRepository struct {
	mu       sync.RWMutex
	settings map[string]*SettingItem
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	repo := &MemoryRepository{
		settings: make(map[string]*SettingItem),
	}

	// 初始化預設設定
	now := time.Now()
	defaults := []*SettingItem{
		{Key: KeyWritePrecision, Value: "millisecond", Description: "寫入時間精度 (second/millisecond)", UpdatedAt: now},
		{Key: KeyPartitionInterval, Value: "monthly", Description: "時序表分區間隔 (daily/weekly/monthly)", UpdatedAt: now},
		{Key: KeyBatchSize, Value: 1000, Description: "批次寫入大小", UpdatedAt: now},
		{Key: KeyDefaultRetryCount, Value: 3, Description: "預設重試次數", UpdatedAt: now},
		{Key: KeyDefaultRetryDelay, Value: 1000, Description: "預設重試延遲 (ms)", UpdatedAt: now},
	}

	for _, item := range defaults {
		repo.settings[item.Key] = item
	}

	return repo
}

// Get 取得設定值
func (r *MemoryRepository) Get(ctx context.Context, key string) (*SettingItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.settings[key]
	if !exists {
		return nil, fmt.Errorf("設定不存在: %s", key)
	}

	itemCopy := *item
	return &itemCopy, nil
}

// Set 設定值
func (r *MemoryRepository) Set(ctx context.Context, key string, value interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, exists := r.settings[key]
	if !exists {
		// 新增設定
		r.settings[key] = &SettingItem{
			Key:       key,
			Value:     value,
			UpdatedAt: time.Now(),
		}
	} else {
		// 更新設定
		item.Value = value
		item.UpdatedAt = time.Now()
	}

	return nil
}

// List 列出所有設定
func (r *MemoryRepository) List(ctx context.Context) ([]*SettingItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*SettingItem, 0, len(r.settings))

	for _, item := range r.settings {
		itemCopy := *item
		result = append(result, &itemCopy)
	}

	return result, nil
}
