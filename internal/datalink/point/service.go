// Package point 提供點位目錄與配置管理功能。
//
// 本套件實作點位的 CRUD 操作、輪詢群組管理和狀態追蹤。
package point

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/common"
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

// =============================================================================
// 點位 CRUD 操作
// =============================================================================

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
	DeviceID       string               `json:"device_id"`
	PollingGroupID string               `json:"polling_group_id"`
	DataType       schema.DataType      `json:"data_type"`
	Enabled        bool                 `json:"enabled"`
	Points         []BatchPointItem     `json:"points"`
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

// =============================================================================
// 輪詢群組操作
// =============================================================================

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

// =============================================================================
// 狀態更新
// =============================================================================

// UpdateReadResult 更新點位讀取結果
func (s *Service) UpdateReadResult(ctx context.Context, pointID string, value interface{}, errMsg string) error {
	return s.repo.UpdateReadResult(ctx, pointID, value, errMsg)
}

// BatchUpdateReadResult 批次更新讀取結果
func (s *Service) BatchUpdateReadResult(ctx context.Context, results []ReadResultUpdate) error {
	return s.repo.BatchUpdateReadResult(ctx, results)
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
// =============================================================================
// 點位詳細資訊
// =============================================================================

// PointDetail 點位詳細資訊 (包含設備和群組資訊)
type PointDetail struct {
	*schema.Point
	DeviceName       string `json:"device_name,omitempty"`
	DeviceProtocol   string `json:"device_protocol,omitempty"`
	PollingGroupName string `json:"polling_group_name,omitempty"`
}

// ToJSON 將最後讀取值轉換為 JSON
func (s *Service) GetLastValue(ctx context.Context, pointID string) (interface{}, error) {
	point, err := s.repo.GetByID(ctx, pointID)
	if err != nil {
		return nil, err
	}

	if point.LastValue == nil {
		return nil, nil
	}

	var value interface{}
	if err := json.Unmarshal([]byte(*point.LastValue), &value); err != nil {
		// 非 JSON 格式，直接返回字串
		return *point.LastValue, nil
	}

	return value, nil
}
