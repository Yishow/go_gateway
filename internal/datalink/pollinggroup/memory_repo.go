package pollinggroup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 記憶體 Repository
// =============================================================================

// MemoryRepository 記憶體內輪詢群組儲存庫
type MemoryRepository struct {
	mu     sync.RWMutex
	groups map[string]*schema.PollingGroup
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		groups: make(map[string]*schema.PollingGroup),
	}
}

// Create 建立新輪詢群組
func (r *MemoryRepository) Create(ctx context.Context, group *schema.PollingGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group.ID]; exists {
		return fmt.Errorf("輪詢群組 ID 已存在: %s", group.ID)
	}

	groupCopy := *group
	r.groups[group.ID] = &groupCopy

	return nil
}

// Update 更新輪詢群組
func (r *MemoryRepository) Update(ctx context.Context, group *schema.PollingGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group.ID]; !exists {
		return fmt.Errorf("輪詢群組不存在: %s", group.ID)
	}

	groupCopy := *group
	r.groups[group.ID] = &groupCopy

	return nil
}

// Delete 刪除輪詢群組
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[id]; !exists {
		return fmt.Errorf("輪詢群組不存在: %s", id)
	}

	delete(r.groups, id)
	return nil
}

// GetByID 根據 ID 取得輪詢群組
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	group, exists := r.groups[id]
	if !exists {
		return nil, fmt.Errorf("輪詢群組不存在: %s", id)
	}

	groupCopy := *group
	return &groupCopy, nil
}

// List 列出所有輪詢群組
func (r *MemoryRepository) List(ctx context.Context) ([]*schema.PollingGroup, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.PollingGroup, 0, len(r.groups))

	for _, group := range r.groups {
		groupCopy := *group
		result = append(result, &groupCopy)
	}

	return result, nil
}

// Count 計算輪詢群組數量
func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.groups)
}

// Clear 清空所有輪詢群組 (用於測試)
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups = make(map[string]*schema.PollingGroup)
}

// Seed 預設資料 (用於開發測試)
func (r *MemoryRepository) Seed() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.groups) > 0 {
		return nil
	}

	now := time.Now()
	defaults := []*schema.PollingGroup{
		{
			ID:          "pg-fast",
			Name:        "Fast Polling",
			Description: "高頻輪詢群組 (1秒)",
			IntervalMs:  1000,
			Priority:    10,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "pg-normal",
			Name:        "Normal Polling",
			Description: "一般輪詢群組 (5秒)",
			IntervalMs:  5000,
			Priority:    50,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "pg-slow",
			Name:        "Slow Polling",
			Description: "低頻輪詢群組 (30秒)",
			IntervalMs:  30000,
			Priority:    100,
			Enabled:     true,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, g := range defaults {
		r.groups[g.ID] = g
	}

	return nil
}
