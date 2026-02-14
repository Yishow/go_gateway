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
