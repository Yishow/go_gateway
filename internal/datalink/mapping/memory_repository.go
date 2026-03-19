package mapping

import (
	"context"
	"fmt"
	"sync"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 記憶體 Repository
// =============================================================================

// MemoryRepository 記憶體內映射儲存庫
type MemoryRepository struct {
	mu       sync.RWMutex
	mappings map[string]*schema.Mapping
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		mappings: make(map[string]*schema.Mapping),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, mapping *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[mapping.ID]; exists {
		return fmt.Errorf("映射 ID 已存在: %s", mapping.ID)
	}

	mappingCopy := *mapping
	r.mappings[mapping.ID] = &mappingCopy
	return nil
}

func (r *MemoryRepository) Update(ctx context.Context, mapping *schema.Mapping) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[mapping.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrMappingNotFound, mapping.ID)
	}

	mappingCopy := *mapping
	r.mappings[mapping.ID] = &mappingCopy
	return nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mappings[id]; !exists {
		return fmt.Errorf("%w: %s", ErrMappingNotFound, id)
	}

	delete(r.mappings, id)
	return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mapping, exists := r.mappings[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrMappingNotFound, id)
	}

	mappingCopy := *mapping
	return &mappingCopy, nil
}

func (r *MemoryRepository) GetByPointID(ctx context.Context, pointID string) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0)
	for _, mapping := range r.mappings {
		if mapping.PointID == pointID {
			mappingCopy := *mapping
			result = append(result, &mappingCopy)
		}
	}
	return result, nil
}

func (r *MemoryRepository) GetByTagID(ctx context.Context, tagID string) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0)
	for _, mapping := range r.mappings {
		if mapping.TagID == tagID {
			mappingCopy := *mapping
			result = append(result, &mappingCopy)
		}
	}
	return result, nil
}

func (r *MemoryRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Mapping, 0, len(r.mappings))

	for _, mapping := range r.mappings {
		if filter.PointID != nil && mapping.PointID != *filter.PointID {
			continue
		}
		if filter.TagID != nil && mapping.TagID != *filter.TagID {
			continue
		}
		if filter.Enabled != nil && mapping.Enabled != *filter.Enabled {
			continue
		}

		mappingCopy := *mapping
		result = append(result, &mappingCopy)
	}

	if filter.Offset > 0 {
		if filter.Offset >= len(result) {
			return []*schema.Mapping{}, nil
		}
		result = result[filter.Offset:]
	}
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// Count 計算映射數量
func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.mappings)
}

// Clear 清空所有映射
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.mappings = make(map[string]*schema.Mapping)
}
