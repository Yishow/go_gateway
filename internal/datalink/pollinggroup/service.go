// Package pollinggroup 提供輪詢群組管理功能。
//
// 本套件實作輪詢群組的 CRUD 操作和配置管理。
package pollinggroup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// Repository 介面定義
// =============================================================================

// Repository 輪詢群組資料存取介面
type Repository interface {
	// Create 建立新輪詢群組
	Create(ctx context.Context, group *schema.PollingGroup) error

	// Update 更新輪詢群組
	Update(ctx context.Context, group *schema.PollingGroup) error

	// Delete 刪除輪詢群組
	Delete(ctx context.Context, id string) error

	// GetByID 根據 ID 取得輪詢群組
	GetByID(ctx context.Context, id string) (*schema.PollingGroup, error)

	// List 列出所有輪詢群組
	List(ctx context.Context) ([]*schema.PollingGroup, error)
}

// =============================================================================
// Service 服務層
// =============================================================================

// Service 輪詢群組管理服務
type Service struct {
	repo Repository
	mu   sync.RWMutex
}

// NewService 建立新的輪詢群組服務
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// CreateRequest 建立輪詢群組請求
type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IntervalMs  int    `json:"interval_ms"`
	Priority    int    `json:"priority,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
}

// Create 建立新輪詢群組
func (s *Service) Create(ctx context.Context, req CreateRequest) (*schema.PollingGroup, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("輪詢群組名稱不能為空")
	}
	if req.IntervalMs < 100 {
		return nil, fmt.Errorf("輪詢間隔不能小於 100ms")
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立輪詢群組 ID 失敗: %w", err)
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	group := &schema.PollingGroup{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IntervalMs:  req.IntervalMs,
		Priority:    req.Priority,
		Enabled:     enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if group.Priority == 0 {
		group.Priority = 100 // 預設優先級
	}

	if err := s.repo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("建立輪詢群組失敗: %w", err)
	}

	return group, nil
}

// UpdateRequest 更新輪詢群組請求
type UpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IntervalMs  *int    `json:"interval_ms,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

// Update 更新輪詢群組
func (s *Service) Update(ctx context.Context, id string, req UpdateRequest) (*schema.PollingGroup, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得輪詢群組失敗: %w", err)
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, fmt.Errorf("輪詢群組名稱不能為空")
		}
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.IntervalMs != nil {
		if *req.IntervalMs < 100 {
			return nil, fmt.Errorf("輪詢間隔不能小於 100ms")
		}
		group.IntervalMs = *req.IntervalMs
	}
	if req.Priority != nil {
		group.Priority = *req.Priority
	}
	if req.Enabled != nil {
		group.Enabled = *req.Enabled
	}

	group.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("更新輪詢群組失敗: %w", err)
	}

	return group, nil
}

// Delete 刪除輪詢群組
func (s *Service) Delete(ctx context.Context, id string) error {
	// TODO: 檢查是否有點位使用此群組
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除輪詢群組失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得輪詢群組
func (s *Service) GetByID(ctx context.Context, id string) (*schema.PollingGroup, error) {
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得輪詢群組失敗: %w", err)
	}
	return group, nil
}

// List 列出所有輪詢群組
func (s *Service) List(ctx context.Context) ([]*schema.PollingGroup, error) {
	groups, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("列出輪詢群組失敗: %w", err)
	}
	return groups, nil
}
