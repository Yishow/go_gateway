package point

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// Create 建立新點位
func (s *Service) Create(ctx context.Context, req CreatePointRequest) (*schema.Point, error) {
	if err := validatePointDataType(req.DataType); err != nil {
		return nil, err
	}

	normalizedAddress, err := validateAndNormalizeAddress(req.Address)
	if err != nil {
		return nil, err
	}
	normalizedFunction, err := validateAndNormalizeFunction(req.Function)
	if err != nil {
		return nil, err
	}

	// 驗證模式
	if req.Mode != "" && req.Mode != schema.PointModeReadOnly && req.Mode != schema.PointModeReadWrite {
		return nil, fmt.Errorf("不支援的點位模式: %s", req.Mode)
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立點位 ID 失敗: %w", err)
	}

	point := &schema.Point{
		ID:             id,
		DeviceID:       req.DeviceID,
		Name:           req.Name,
		Description:    req.Description,
		Address:        normalizedAddress,
		Function:       normalizedFunction,
		DataType:       req.DataType,
		Mode:           req.Mode,
		PollingGroupID: req.PollingGroupID,
		Enabled:        true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 設定預設模式
	if point.Mode == "" {
		point.Mode = schema.PointModeReadOnly
	}

	if err := s.repo.Create(ctx, point); err != nil {
		return nil, fmt.Errorf("建立點位失敗: %w", err)
	}

	return point, nil
}

// BatchCreatePointsRequest 批次建立點位請求
type BatchCreatePointsRequest struct {
	DeviceID       string           `json:"device_id"`
	PollingGroupID string           `json:"polling_group_id"`
	DataType       schema.DataType  `json:"data_type"`
	Enabled        bool             `json:"enabled"`
	DryRun         bool             `json:"dry_run,omitempty"`
	ApplyIfClean   bool             `json:"apply_if_clean,omitempty"`
	Points         []BatchPointItem `json:"points"`
}

// BatchPointItem 批次建立的單點資訊
type BatchPointItem struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Function string `json:"function,omitempty"`
}

// BatchCreateResult 批次建立結果
type BatchCreateResult struct {
	CreatedCount int             `json:"created_count"`
	Points       []*schema.Point `json:"points"`
	Errors       []string        `json:"errors,omitempty"`
	DryRun       bool            `json:"dry_run"`
	Applied      bool            `json:"applied"`
}

// BatchCreate 批次建立新點位
func (s *Service) BatchCreate(ctx context.Context, req BatchCreatePointsRequest) (*BatchCreateResult, error) {
	result, err := s.BatchDryRun(ctx, req)
	if err != nil {
		return nil, err
	}
	result.DryRun = req.DryRun

	if len(result.Errors) > 0 {
		return result, nil
	}
	if req.DryRun && !req.ApplyIfClean {
		return result, nil
	}

	for _, item := range req.Points {
		createReq := CreatePointRequest{
			DeviceID:       req.DeviceID,
			Name:           strings.TrimSpace(item.Name),
			Address:        strings.TrimSpace(item.Address),
			Function:       strings.TrimSpace(item.Function),
			DataType:       req.DataType,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: nil,
		}
		if req.PollingGroupID != "" {
			groupID := req.PollingGroupID
			createReq.PollingGroupID = &groupID
		}

		point, createErr := s.Create(ctx, createReq)
		if createErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("點位 %s (%s) 建立失敗: %s", item.Name, item.Address, createErr.Error()))
			continue
		}
		result.Points = append(result.Points, point)
		result.CreatedCount++
	}
	result.Applied = result.CreatedCount > 0

	return result, nil
}

// BatchDryRun 批次建立預檢（只驗證，不落 DB）。
func (s *Service) BatchDryRun(ctx context.Context, req BatchCreatePointsRequest) (*BatchCreateResult, error) {
	result := &BatchCreateResult{
		Points: make([]*schema.Point, 0),
		Errors: make([]string, 0),
		DryRun: true,
	}

	errors, err := s.validateBatchCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	result.Errors = append(result.Errors, errors...)
	return result, nil
}

// CreatePointRequest 建立點位請求
type CreatePointRequest struct {
	DeviceID       string           `json:"device_id"`
	Name           string           `json:"name"`
	Description    string           `json:"description,omitempty"`
	Address        string           `json:"address"`
	Function       string           `json:"function,omitempty"`
	DataType       schema.DataType  `json:"data_type"`
	Mode           schema.PointMode `json:"mode,omitempty"`
	PollingGroupID *string          `json:"polling_group_id,omitempty"`
}

// Update 更新點位
func (s *Service) Update(ctx context.Context, id string, req UpdatePointRequest) (*schema.Point, error) {
	point, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得點位失敗: %w", err)
	}

	if req.Name != nil {
		point.Name = *req.Name
	}
	if req.Description != nil {
		point.Description = *req.Description
	}
	if req.Address != nil {
		address, err := validateAndNormalizeAddress(*req.Address)
		if err != nil {
			return nil, err
		}
		point.Address = address
	}
	if req.Function != nil {
		function, err := validateAndNormalizeFunction(*req.Function)
		if err != nil {
			return nil, err
		}
		point.Function = function
	}
	if req.DataType != nil {
		if err := validatePointDataType(*req.DataType); err != nil {
			return nil, err
		}
		point.DataType = *req.DataType
	}
	if req.Mode != nil {
		point.Mode = *req.Mode
	}
	if req.PollingGroupID != nil {
		point.PollingGroupID = req.PollingGroupID
	}
	if req.Enabled != nil {
		point.Enabled = *req.Enabled
	}

	point.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, point); err != nil {
		return nil, fmt.Errorf("更新點位失敗: %w", err)
	}

	return point, nil
}

// UpdatePointRequest 更新點位請求
type UpdatePointRequest struct {
	Name           *string           `json:"name,omitempty"`
	Description    *string           `json:"description,omitempty"`
	Address        *string           `json:"address,omitempty"`
	Function       *string           `json:"function,omitempty"`
	DataType       *schema.DataType  `json:"data_type,omitempty"`
	Mode           *schema.PointMode `json:"mode,omitempty"`
	PollingGroupID *string           `json:"polling_group_id,omitempty"`
	Enabled        *bool             `json:"enabled,omitempty"`
}

// Delete 刪除點位
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除點位失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得點位
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	point, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得點位失敗: %w", err)
	}
	return point, nil
}

// ListByDevice 列出設備的所有點位
func (s *Service) ListByDevice(ctx context.Context, deviceID string) ([]*schema.Point, error) {
	points, err := s.repo.ListByDevice(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("列出點位失敗: %w", err)
	}
	return points, nil
}

// List 列出點位
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Point, error) {
	points, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出點位失敗: %w", err)
	}
	return points, nil
}
