// Package tag 提供全域標籤字典管理功能。
//
// 本套件實作標籤的 CRUD 操作、元資料管理和生命週期狀態控制。
package tag

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 標籤鍵驗證
// =============================================================================

// 標籤鍵正則表達式：允許 ASCII 字母、數字、底線、連字號、點號、斜線
// 長度 1-128 字元
var tagKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-./]{1,128}$`)

// ValidateTagKey 驗證標籤鍵格式
func ValidateTagKey(key string) error {
	if key == "" {
		return fmt.Errorf("標籤鍵不能為空")
	}
	if len(key) > 128 {
		return fmt.Errorf("標籤鍵長度超過 128 字元")
	}
	if !tagKeyPattern.MatchString(key) {
		return fmt.Errorf("標籤鍵格式無效，只允許 ASCII 字母、數字、底線、連字號、點號、斜線")
	}
	return nil
}

// NormalizeTagKey 正規化標籤鍵 (轉小寫)
func NormalizeTagKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 標籤資料存取介面
type Repository interface {
	// Create 建立新標籤
	Create(ctx context.Context, tag *schema.Tag) error

	// Update 更新標籤
	Update(ctx context.Context, tag *schema.Tag) error

	// Delete 刪除標籤
	Delete(ctx context.Context, id string) error

	// GetByID 根據 ID 取得標籤
	GetByID(ctx context.Context, id string) (*schema.Tag, error)

	// GetByKey 根據鍵取得標籤 (大小寫不敏感)
	GetByKey(ctx context.Context, key string) (*schema.Tag, error)

	// List 列出所有標籤
	List(ctx context.Context, filter ListFilter) ([]*schema.Tag, error)

	// UpdateStatus 更新標籤狀態
	UpdateStatus(ctx context.Context, id string, status schema.TagStatus) error

	// ExistsByKey 檢查鍵是否已存在
	ExistsByKey(ctx context.Context, key string) (bool, error)
}

// ListFilter 標籤列表篩選條件
type ListFilter struct {
	// Status 篩選狀態
	Status *schema.TagStatus

	// DataType 篩選資料型別
	DataType *schema.DataType

	// KeyPrefix 鍵前綴篩選
	KeyPrefix string

	// Labels 標籤屬性篩選
	Labels map[string]string

	// Limit 限制數量
	Limit int

	// Offset 偏移量
	Offset int
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 標籤管理服務
type Service struct {
	repo Repository
	mu   sync.RWMutex
}

// NewService 建立新的標籤服務
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新標籤
func (s *Service) Create(ctx context.Context, req CreateTagRequest) (*schema.Tag, error) {
	// 驗證標籤鍵
	if err := ValidateTagKey(req.Key); err != nil {
		return nil, err
	}

	// 檢查鍵是否已存在
	exists, err := s.repo.ExistsByKey(ctx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("檢查標籤鍵失敗: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("標籤鍵已存在: %s", req.Key)
	}

	// 驗證資料型別
	if !isValidDataType(req.DataType) {
		return nil, fmt.Errorf("不支援的資料型別: %s", req.DataType)
	}

	tag := &schema.Tag{
		ID:          generateUUID(),
		Key:         req.Key,
		KeyLower:    NormalizeTagKey(req.Key),
		DisplayName: req.DisplayName,
		Description: req.Description,
		Unit:        req.Unit,
		DataType:    req.DataType,
		Status:      schema.TagStatusDraft,
		Labels:      req.Labels,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 設定預設顯示名稱
	if tag.DisplayName == "" {
		tag.DisplayName = req.Key
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("建立標籤失敗: %w", err)
	}

	return tag, nil
}

// CreateTagRequest 建立標籤請求
type CreateTagRequest struct {
	Key         string            `json:"key"`
	DisplayName string            `json:"display_name,omitempty"`
	Description string            `json:"description,omitempty"`
	Unit        string            `json:"unit,omitempty"`
	DataType    schema.DataType   `json:"data_type"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Update 更新標籤
func (s *Service) Update(ctx context.Context, id string, req UpdateTagRequest) (*schema.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}

	// 標籤鍵不可修改 (若提供則驗證一致性)
	if req.Key != nil && NormalizeTagKey(*req.Key) != tag.KeyLower {
		return nil, fmt.Errorf("標籤鍵不可修改")
	}

	if req.DisplayName != nil {
		tag.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		tag.Description = *req.Description
	}
	if req.Unit != nil {
		tag.Unit = *req.Unit
	}
	if req.DataType != nil {
		if !isValidDataType(*req.DataType) {
			return nil, fmt.Errorf("不支援的資料型別: %s", *req.DataType)
		}
		tag.DataType = *req.DataType
	}
	if req.Labels != nil {
		tag.Labels = req.Labels
	}

	tag.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("更新標籤失敗: %w", err)
	}

	return tag, nil
}

// UpdateTagRequest 更新標籤請求
type UpdateTagRequest struct {
	Key         *string           `json:"key,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Unit        *string           `json:"unit,omitempty"`
	DataType    *schema.DataType  `json:"data_type,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Delete 刪除標籤
func (s *Service) Delete(ctx context.Context, id string) error {
	// TODO: 檢查是否有映射使用此標籤
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除標籤失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得標籤
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	return tag, nil
}

// GetByKey 根據鍵取得標籤
func (s *Service) GetByKey(ctx context.Context, key string) (*schema.Tag, error) {
	tag, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	return tag, nil
}

// List 列出標籤
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Tag, error) {
	tags, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出標籤失敗: %w", err)
	}
	return tags, nil
}

// =============================================================================
// 生命週期管理
// =============================================================================

// Activate 啟用標籤
func (s *Service) Activate(ctx context.Context, id string) error {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得標籤失敗: %w", err)
	}

	if tag.Status == schema.TagStatusActive {
		return nil
	}

	if tag.Status == schema.TagStatusRetired {
		return fmt.Errorf("已退役的標籤不能重新啟用")
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.TagStatusActive); err != nil {
		return fmt.Errorf("更新標籤狀態失敗: %w", err)
	}

	return nil
}

// Retire 退役標籤
func (s *Service) Retire(ctx context.Context, id string) error {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得標籤失敗: %w", err)
	}

	if tag.Status == schema.TagStatusRetired {
		return nil
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.TagStatusRetired); err != nil {
		return fmt.Errorf("更新標籤狀態失敗: %w", err)
	}

	return nil
}

// =============================================================================
// 輔助函數
// =============================================================================

// isValidDataType 檢查資料型別是否有效
func isValidDataType(dataType schema.DataType) bool {
	validTypes := map[schema.DataType]bool{
		schema.DataTypeBool:    true,
		schema.DataTypeInt16:   true,
		schema.DataTypeUint16:  true,
		schema.DataTypeInt32:   true,
		schema.DataTypeUint32:  true,
		schema.DataTypeInt64:   true,
		schema.DataTypeUint64:  true,
		schema.DataTypeFloat32: true,
		schema.DataTypeFloat64: true,
		schema.DataTypeString:  true,
	}
	return validTypes[dataType]
}

// generateUUID 產生 UUID
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// =============================================================================
// 記憶體 Repository
// =============================================================================

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

// Update 更新標籤
func (r *MemoryRepository) Update(ctx context.Context, tag *schema.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tags[tag.ID]; !exists {
		return fmt.Errorf("標籤不存在: %s", tag.ID)
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
		return fmt.Errorf("標籤不存在: %s", id)
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
		return nil, fmt.Errorf("標籤不存在: %s", id)
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

	return nil, fmt.Errorf("標籤不存在: %s", key)
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
			match := true
			for k, v := range filter.Labels {
				if tag.Labels[k] != v {
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
		return fmt.Errorf("標籤不存在: %s", id)
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
