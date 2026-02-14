package device

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新設備
func (s *Service) Create(ctx context.Context, req CreateDeviceRequest) (*schema.Device, error) {
	// 驗證協議類型
	if !isValidProtocol(req.Protocol) {
		return nil, fmt.Errorf("不支援的協議類型: %s", req.Protocol)
	}

	// 驗證連線配置
	if err := validateConnectionConfig(req.Protocol, req.ConnectionConfig); err != nil {
		return nil, fmt.Errorf("連線配置無效: %w", err)
	}

	// 序列化連線配置
	configJSON, err := json.Marshal(req.ConnectionConfig)
	if err != nil {
		return nil, fmt.Errorf("序列化連線配置失敗: %w", err)
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立設備 ID 失敗: %w", err)
	}

	device := &schema.Device{
		ID:               id,
		Name:             req.Name,
		Description:      req.Description,
		Protocol:         req.Protocol,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(configJSON),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, fmt.Errorf("建立設備失敗: %w", err)
	}

	return device, nil
}

// CreateDeviceRequest 建立設備請求
type CreateDeviceRequest struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	Protocol         schema.ProtocolType    `json:"protocol"`
	ConnectionConfig map[string]interface{} `json:"connection_config"`
}

// Update 更新設備
func (s *Service) Update(ctx context.Context, id string, req UpdateDeviceRequest) (*schema.Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	if req.Name != nil {
		device.Name = *req.Name
	}
	if req.Description != nil {
		device.Description = *req.Description
	}
	if req.ConnectionConfig != nil {
		// 驗證連線配置
		if err := validateConnectionConfig(device.Protocol, req.ConnectionConfig); err != nil {
			return nil, fmt.Errorf("連線配置無效: %w", err)
		}
		configJSON, err := json.Marshal(req.ConnectionConfig)
		if err != nil {
			return nil, fmt.Errorf("序列化連線配置失敗: %w", err)
		}
		device.ConnectionConfig = string(configJSON)

		// 如果配置變更，關閉現有連線
		s.connMgr.Close(id)
	}

	device.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("更新設備失敗: %w", err)
	}

	return device, nil
}

// UpdateDeviceRequest 更新設備請求
type UpdateDeviceRequest struct {
	Name             *string                `json:"name,omitempty"`
	Description      *string                `json:"description,omitempty"`
	ConnectionConfig map[string]interface{} `json:"connection_config,omitempty"`
}

// Delete 刪除設備
func (s *Service) Delete(ctx context.Context, id string) error {
	// 關閉現有連線
	s.connMgr.Close(id)

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除設備失敗: %w", err)
	}

	return nil
}

// GetByID 根據 ID 取得設備
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}
	return device, nil
}

// List 列出設備
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	devices, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出設備失敗: %w", err)
	}
	return devices, nil
}
