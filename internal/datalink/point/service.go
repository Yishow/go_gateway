// Package point 提供點位目錄與配置管理功能。
//
// 本套件實作點位的 CRUD 操作、輪詢群組管理和狀態追蹤。
package point

import (
	"context"
	"sync"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 點位資料存取介面
type Repository interface {
	// Create 建立新點位
	Create(ctx context.Context, point *schema.Point) error

	// Update 更新點位
	Update(ctx context.Context, point *schema.Point) error

	// Delete 刪除點位
	Delete(ctx context.Context, id string) error

	// GetByID 根據 ID 取得點位
	GetByID(ctx context.Context, id string) (*schema.Point, error)

	// ListByDevice 列出設備的所有點位
	ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error)

	// ListByPollingGroup 列出輪詢群組的所有點位
	ListByPollingGroup(ctx context.Context, groupID string) ([]*schema.Point, error)

	// List 列出所有點位
	List(ctx context.Context, filter ListFilter) ([]*schema.Point, error)

	// UpdateReadResult 更新讀取結果
	UpdateReadResult(ctx context.Context, id string, value interface{}, errMsg string) error

	// BatchUpdateReadResult 批次更新讀取結果
	BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error
}

// ListFilter 點位列表篩選條件
type ListFilter struct {
	DeviceID       *string
	PollingGroupID *string
	Enabled        *bool
	DataType       *schema.DataType
	Limit          int
	Offset         int
}

// ReadResultUpdate 讀取結果更新
type ReadResultUpdate struct {
	PointID string
	Value   interface{}
	Error   string
}

// =============================================================================
// PollingGroupRepository 輪詢群組資料存取介面
// =============================================================================

// PollingGroupRepository 輪詢群組資料存取介面
type PollingGroupRepository interface {
	Create(ctx context.Context, group *schema.PollingGroup) error
	Update(ctx context.Context, group *schema.PollingGroup) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*schema.PollingGroup, error)
	List(ctx context.Context) ([]*schema.PollingGroup, error)
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 點位管理服務
type Service struct {
	repo      Repository
	groupRepo PollingGroupRepository
	mu        sync.RWMutex
}

// NewService 建立新的點位服務
func NewService(repo Repository, groupRepo PollingGroupRepository) *Service {
	return &Service{
		repo:      repo,
		groupRepo: groupRepo,
	}
}
