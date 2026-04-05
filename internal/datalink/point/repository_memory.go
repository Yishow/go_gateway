package point

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 記憶體 Repository (用於測試和原型)
// =============================================================================

// MemoryRepository 記憶體內點位儲存庫
type MemoryRepository struct {
	mu     sync.RWMutex
	points map[string]*schema.Point
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		points: make(map[string]*schema.Point),
	}
}

// Create 建立新點位
func (r *MemoryRepository) Create(ctx context.Context, point *schema.Point) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.points[point.ID]; exists {
		return fmt.Errorf("點位 ID 已存在: %s", point.ID)
	}

	pointCopy := *point
	r.points[point.ID] = &pointCopy

	return nil
}

// Update 更新點位
func (r *MemoryRepository) Update(ctx context.Context, point *schema.Point) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.points[point.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrPointNotFound, point.ID)
	}

	pointCopy := *point
	r.points[point.ID] = &pointCopy

	return nil
}

// Delete 刪除點位
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.points[id]; !exists {
		return fmt.Errorf("%w: %s", ErrPointNotFound, id)
	}

	delete(r.points, id)
	return nil
}

// GetByID 根據 ID 取得點位
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	point, exists := r.points[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrPointNotFound, id)
	}

	pointCopy := *point
	return &pointCopy, nil
}

// ListByDevice 列出設備的所有點位
func (r *MemoryRepository) ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Point, 0)
	for _, point := range r.points {
		if point.DeviceID == deviceID {
			pointCopy := *point
			result = append(result, &pointCopy)
		}
	}

	return result, nil
}

// ListByPollingGroup 列出輪詢群組的所有點位
func (r *MemoryRepository) ListByPollingGroup(ctx context.Context, groupID string) ([]*schema.Point, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Point, 0)
	for _, point := range r.points {
		if point.PollingGroupID != nil && *point.PollingGroupID == groupID {
			pointCopy := *point
			result = append(result, &pointCopy)
		}
	}

	return result, nil
}

// List 列出所有點位
func (r *MemoryRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Point, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Point, 0, len(r.points))

	for _, point := range r.points {
		// 套用篩選條件
		if filter.DeviceID != nil && point.DeviceID != *filter.DeviceID {
			continue
		}
		if filter.PollingGroupID != nil {
			if point.PollingGroupID == nil || *point.PollingGroupID != *filter.PollingGroupID {
				continue
			}
		}
		if filter.Enabled != nil && point.Enabled != *filter.Enabled {
			continue
		}
		if filter.DataType != nil && point.DataType != *filter.DataType {
			continue
		}

		pointCopy := *point
		result = append(result, &pointCopy)
	}

	// 套用分頁
	if filter.Offset > 0 {
		if filter.Offset >= len(result) {
			return []*schema.Point{}, nil
		}
		result = result[filter.Offset:]
	}
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// UpdateReadResult 更新讀取結果
func (r *MemoryRepository) UpdateReadResult(ctx context.Context, id string, value interface{}, errMsg string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	point, exists := r.points[id]
	if !exists {
		return fmt.Errorf("%w: %s", ErrPointNotFound, id)
	}

	now := time.Now()
	point.LastReadAt = &now

	if errMsg != "" {
		point.LastError = errMsg
	} else {
		point.LastError = ""
		// 序列化值為 JSON
		valueJSON, err := json.Marshal(value)
		if err == nil {
			valueStr := string(valueJSON)
			point.LastValue = &valueStr
		}
	}

	point.UpdatedAt = now

	return nil
}

// BatchUpdateReadResult 批次更新讀取結果
func (r *MemoryRepository) BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	for _, result := range results {
		point, exists := r.points[result.PointID]
		if !exists {
			continue
		}

		point.LastReadAt = &now

		if result.Error != "" {
			point.LastError = result.Error
		} else {
			point.LastError = ""
			valueJSON, err := json.Marshal(result.Value)
			if err == nil {
				valueStr := string(valueJSON)
				point.LastValue = &valueStr
			}
		}

		point.UpdatedAt = now
	}

	return nil
}

// Count 計算點位數量
func (r *MemoryRepository) Count(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.points)
}

// Clear 清空所有點位 (用於測試)
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.points = make(map[string]*schema.Point)
}

// =============================================================================
// 輪詢群組記憶體 Repository
// =============================================================================

// MemoryPollingGroupRepository 記憶體內輪詢群組儲存庫
type MemoryPollingGroupRepository struct {
	mu     sync.RWMutex
	groups map[string]*schema.PollingGroup
}

// NewMemoryPollingGroupRepository 建立新的記憶體儲存庫
func NewMemoryPollingGroupRepository() *MemoryPollingGroupRepository {
	return &MemoryPollingGroupRepository{
		groups: make(map[string]*schema.PollingGroup),
	}
}

// Create 建立新輪詢群組
func (r *MemoryPollingGroupRepository) Create(ctx context.Context, group *schema.PollingGroup) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[group.ID]; exists {
		return fmt.Errorf("輪詢群組 ID 已存在: %s", group.ID)
	}

	// 檢查名稱唯一性
	for _, existing := range r.groups {
		if existing.Name == group.Name {
			return fmt.Errorf("輪詢群組名稱已存在: %s", group.Name)
		}
	}

	groupCopy := *group
	r.groups[group.ID] = &groupCopy

	return nil
}

// Update 更新輪詢群組
func (r *MemoryPollingGroupRepository) Update(ctx context.Context, group *schema.PollingGroup) error {
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
func (r *MemoryPollingGroupRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.groups[id]; !exists {
		return fmt.Errorf("輪詢群組不存在: %s", id)
	}

	delete(r.groups, id)
	return nil
}

// GetByID 根據 ID 取得輪詢群組
func (r *MemoryPollingGroupRepository) GetByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
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
func (r *MemoryPollingGroupRepository) List(ctx context.Context) ([]*schema.PollingGroup, error) {
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
func (r *MemoryPollingGroupRepository) Count(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.groups)
}

// Clear 清空所有輪詢群組 (用於測試)
func (r *MemoryPollingGroupRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups = make(map[string]*schema.PollingGroup)
}
