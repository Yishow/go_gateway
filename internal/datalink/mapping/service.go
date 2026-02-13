// Package mapping 提供點位到標籤的映射管線功能。
//
// 本套件實作映射的 CRUD 操作、轉換管線執行和預覽功能。
package mapping

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 映射資料存取介面
type Repository interface {
	Create(ctx context.Context, mapping *schema.Mapping) error
	Update(ctx context.Context, mapping *schema.Mapping) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.Mapping, error)
	GetByPointID(ctx context.Context, pointID string) ([]*schema.Mapping, error)
	GetByTagID(ctx context.Context, tagID string) ([]*schema.Mapping, error)
	List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error)
}

// ListFilter 映射列表篩選條件
type ListFilter struct {
	PointID *string
	TagID   *string
	Enabled *bool
	Limit   int
	Offset  int
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 映射管理服務
type Service struct {
	repo Repository
	mu   sync.RWMutex
}

// NewService 建立新的映射服務
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}
