package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// MemoryRepository 記憶體內標籤儲存庫
type MemoryRepository struct {
	mu   sync.RWMutex
	tags map[string]*schema.Tag
}

// NewMemoryRepository 建立新的記憶體儲存庫
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tags: make(map[string]*schema.Tag),
	}
}

// Create 建立新標籤
func (r *MemoryRepository) Create(ctx context.Context, tag *schema.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[tag.ID]; exists {
		return fmt.Errorf("標籤 ID 已存在: %s", tag.ID)
	}

	// 檢查鍵唯一性
	for _, existing := range r.tags {
		if existing.KeyLower == tag.KeyLower {
			return fmt.Errorf("標籤鍵已存在: %s", tag.Key)
		}
	}

	tagCopy := *tag
	r.tags[tag.ID] = &tagCopy

	return nil
}

// BatchCreate 批量建立標籤
func (r *MemoryRepository) BatchCreate(ctx context.Context, tags []*schema.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, tag := range tags {
		if _, exists := r.tags[tag.ID]; exists {
			return fmt.Errorf("標籤 ID 已存在: %s", tag.ID)
		}
		for _, existing := range r.tags {
			if existing.KeyLower == tag.KeyLower {
				return fmt.Errorf("標籤鍵已存在: %s", tag.Key)
			}
		}
	}

	for _, tag := range tags {
		tagCopy := *tag
		r.tags[tag.ID] = &tagCopy
	}

	return nil
}

// Update 更新標籤
func (r *MemoryRepository) Update(ctx context.Context, tag *schema.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[tag.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrTagNotFound, tag.ID)
	}

	tagCopy := *tag
	r.tags[tag.ID] = &tagCopy

	return nil
}

// Delete 刪除標籤
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[id]; !exists {
		return fmt.Errorf("%w: %s", ErrTagNotFound, id)
	}

	delete(r.tags, id)
	return nil
}

// GetByID 根據 ID 取得標籤
func (r *MemoryRepository) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tag, exists := r.tags[id]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTagNotFound, id)
	}

	tagCopy := *tag
	return &tagCopy, nil
}

// GetByKey 根據鍵取得標籤
func (r *MemoryRepository) GetByKey(ctx context.Context, key string) (*schema.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keyLower := NormalizeTagKey(key)
	for _, tag := range r.tags {
		if tag.KeyLower == keyLower {
			tagCopy := *tag
			return &tagCopy, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrTagNotFound, key)
}

// List 列出所有標籤
func (r *MemoryRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*schema.Tag, 0, len(r.tags))

	for _, tag := range r.tags {
		// 套用篩選條件
		if filter.Status != nil && tag.Status != *filter.Status {
			continue
		}
		if filter.DataType != nil && tag.DataType != *filter.DataType {
			continue
		}
		if filter.KeyPrefix != "" && !strings.HasPrefix(tag.KeyLower, strings.ToLower(filter.KeyPrefix)) {
			continue
		}
		if len(filter.Labels) > 0 {
			// 解析標籤 Labels JSON
			var tagLabels map[string]string
			if err := json.Unmarshal([]byte(tag.Labels), &tagLabels); err != nil {
				continue
			}
			match := true
			for k, v := range filter.Labels {
				if tagLabels[k] != v {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		tagCopy := *tag
		result = append(result, &tagCopy)
	}

	// 套用分頁
	if filter.Offset > 0 {
		if filter.Offset >= len(result) {
			return []*schema.Tag{}, nil
		}
		result = result[filter.Offset:]
	}
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}

	return result, nil
}

// UpdateStatus 更新標籤狀態
func (r *MemoryRepository) UpdateStatus(ctx context.Context, id string, status schema.TagStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tag, exists := r.tags[id]
	if !exists {
		return fmt.Errorf("%w: %s", ErrTagNotFound, id)
	}

	tag.Status = status
	tag.UpdatedAt = time.Now()

	return nil
}

// ExistsByKey 檢查鍵是否已存在
func (r *MemoryRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keyLower := NormalizeTagKey(key)
	for _, tag := range r.tags {
		if tag.KeyLower == keyLower {
			return true, nil
		}
	}

	return false, nil
}

// Count 計算標籤數量
func (r *MemoryRepository) Count(ctx context.Context) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tags)
}

// Clear 清空所有標籤 (用於測試)
func (r *MemoryRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags = make(map[string]*schema.Tag)
}
