package point

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// CreatePollingGroup 建立輪詢群組
func (s *Service) CreatePollingGroup(ctx context.Context, req CreatePollingGroupRequest) (*schema.PollingGroup, error) {
	if req.IntervalMs < 100 {
		return nil, fmt.Errorf("輪詢間隔不能小於 100ms")
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立輪詢群組 ID 失敗: %w", err)
	}

	group := &schema.PollingGroup{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IntervalMs:  req.IntervalMs,
		Priority:    req.Priority,
		Enabled:     true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if group.Priority == 0 {
		group.Priority = 100 // 預設優先級
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, fmt.Errorf("建立輪詢群組失敗: %w", err)
	}

	return group, nil
}

// CreatePollingGroupRequest 建立輪詢群組請求
type CreatePollingGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IntervalMs  int    `json:"interval_ms"`
	Priority    int    `json:"priority,omitempty"`
}

// UpdatePollingGroup 更新輪詢群組
func (s *Service) UpdatePollingGroup(ctx context.Context, id string, req UpdatePollingGroupRequest) (*schema.PollingGroup, error) {
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得輪詢群組失敗: %w", err)
	}

	if req.Name != nil {
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

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, fmt.Errorf("更新輪詢群組失敗: %w", err)
	}

	return group, nil
}

// UpdatePollingGroupRequest 更新輪詢群組請求
type UpdatePollingGroupRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	IntervalMs  *int    `json:"interval_ms,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
}

// DeletePollingGroup 刪除輪詢群組
func (s *Service) DeletePollingGroup(ctx context.Context, id string) error {
	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除輪詢群組失敗: %w", err)
	}
	return nil
}

// ListPollingGroups 列出所有輪詢群組
func (s *Service) ListPollingGroups(ctx context.Context) ([]*schema.PollingGroup, error) {
	groups, err := s.groupRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("列出輪詢群組失敗: %w", err)
	}
	return groups, nil
}

// AssignToGroup 將點位分配到輪詢群組
func (s *Service) AssignToGroup(ctx context.Context, pointID, groupID string) error {
	point, err := s.repo.GetByID(ctx, pointID)
	if err != nil {
		return fmt.Errorf("取得點位失敗: %w", err)
	}

	// 驗證群組存在
	if _, err := s.groupRepo.GetByID(ctx, groupID); err != nil {
		return fmt.Errorf("輪詢群組不存在: %w", err)
	}

	point.PollingGroupID = &groupID
	point.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, point); err != nil {
		return fmt.Errorf("更新點位失敗: %w", err)
	}

	return nil
}

// RemoveFromGroup 將點位從輪詢群組移除
func (s *Service) RemoveFromGroup(ctx context.Context, pointID string) error {
	point, err := s.repo.GetByID(ctx, pointID)
	if err != nil {
		return fmt.Errorf("取得點位失敗: %w", err)
	}

	point.PollingGroupID = nil
	point.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, point); err != nil {
		return fmt.Errorf("更新點位失敗: %w", err)
	}

	return nil
}
