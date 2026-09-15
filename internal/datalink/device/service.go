// Package device 提供設備註冊與管理功能。
//
// 本套件實作設備的 CRUD 操作、連線測試和狀態管理。
package device

import (
	"context"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 設備資料存取介面
type Repository interface {
	// Create 建立新設備
	Create(ctx context.Context, device *schema.Device) error

	// Update 更新設備
	Update(ctx context.Context, device *schema.Device) error

	// Delete 刪除設備
	Delete(ctx context.Context, id string) error

	// GetByID 根據 ID 取得設備
	GetByID(ctx context.Context, id string) (*schema.Device, error)

	// List 列出所有設備
	List(ctx context.Context, filter ListFilter) ([]*schema.Device, error)

	// UpdateTestResult 更新連線測試結果
	UpdateTestResult(ctx context.Context, id string, success bool, errMsg string) error
	// ClearTestResult 清除已失效的連線測試結果
	ClearTestResult(ctx context.Context, id string) error

	// UpdateStatus 更新設備狀態
	UpdateStatus(ctx context.Context, id string, status schema.DeviceStatus) error
}

// ListFilter 設備列表篩選條件
type ListFilter struct {
	// Protocol 篩選協議類型
	Protocol *schema.ProtocolType

	// Status 篩選狀態
	Status *schema.DeviceStatus

	// Limit 限制數量
	Limit int

	// Offset 偏移量
	Offset int
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 設備管理服務
type Service struct {
	repo    Repository
	connMgr *connector.ConnectionManager
}

// NewService 建立新的設備服務
func NewService(repo Repository, connMgr *connector.ConnectionManager) *Service {
	if connMgr == nil {
		connMgr = connector.GetConnectionManager()
	}
	return &Service{
		repo:    repo,
		connMgr: connMgr,
	}
}
