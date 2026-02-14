package point

import (
	"context"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// Create 建立新點位
func (s *Service) Create(ctx context.Context, req CreatePointRequest) (*schema.Point, error) {
	// 驗證資料型別
	if !isValidDataType(req.DataType) {
		return nil, fmt.Errorf("不支援的資料型別: %s", req.DataType)
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
		Address:        req.Address,
		Function:       req.Function,
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
	Points         []BatchPointItem `json:"points"`
}

// BatchPointItem 批次建立的單點資訊
type BatchPointItem struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// BatchCreateResult 批次建立結果
type BatchCreateResult struct {
	CreatedCount int             `json:"created_count"`
	Points       []*schema.Point `json:"points"`
	Errors       []string        `json:"errors,omitempty"`
}

// BatchCreate 批次建立新點位
func (s *Service) BatchCreate(ctx context.Context, req BatchCreatePointsRequest) (*BatchCreateResult, error) {
	result := &BatchCreateResult{
		Points: make([]*schema.Point, 0),
		Errors: make([]string, 0),
	}

	for _, item := range req.Points {
		// 建構單點建立請求
		createReq := CreatePointRequest{
			DeviceID:       req.DeviceID,
			Name:           item.Name,
			Address:        item.Address,
			DataType:       req.DataType,
			Mode:           schema.PointModeReadOnly, // 預設唯讀
			PollingGroupID: nil,
		}

		if req.PollingGroupID != "" {
			groupID := req.PollingGroupID
			createReq.PollingGroupID = &groupID
		}

		point, err := s.Create(ctx, createReq)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("點位 %s (%s) 建立失敗: %s", item.Name, item.Address, err.Error()))
			continue
		}

		result.Points = append(result.Points, point)
		result.CreatedCount++
	}

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
		point.Address = *req.Address
	}
	if req.Function != nil {
		point.Function = *req.Function
	}
	if req.DataType != nil {
		if !isValidDataType(*req.DataType) {
			return nil, fmt.Errorf("不支援的資料型別: %s", *req.DataType)
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
