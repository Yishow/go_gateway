package device

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新設備
func (s *Service) Create(ctx context.Context, req CreateDeviceRequest) (*schema.Device, error) {
	if err := validateDeviceName(req.Name); err != nil {
		return nil, err
	}

	// 驗證協議類型
	if !isValidProtocol(req.Protocol) {
		return nil, validationError(fmt.Sprintf("不支援的協議類型: %s", req.Protocol))
	}

	// 驗證連線配置
	if err := validateConnectionConfig(req.Protocol, req.ConnectionConfig); err != nil {
		return nil, connectionConfigValidationError(err)
	}

	// 序列化連線配置
	configJSON, err := json.Marshal(req.ConnectionConfig)
	if err != nil {
		return nil, fmt.Errorf("序列化連線配置失敗: %w", err)
	}

	id := strings.TrimSpace(req.ID)
	if id == "" {
		id, err = common.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("建立設備 ID 失敗: %w", err)
		}
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
	ID               string              `json:"id,omitempty"`
	Name             string              `json:"name"`
	Description      string              `json:"description,omitempty"`
	Protocol         schema.ProtocolType `json:"protocol"`
	ConnectionConfig map[string]any      `json:"connection_config"`
}

const (
	// MinAcquisitionIntervalMs 設備採集頻率最小毫秒數 (100ms)
	MinAcquisitionIntervalMs = 100
	// MaxAcquisitionIntervalMs 設備採集頻率最大毫秒數 (1h)
	MaxAcquisitionIntervalMs = 3600000
)

// Update 更新設備
func (s *Service) Update(ctx context.Context, id string, req UpdateDeviceRequest) (*schema.Device, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	if req.Name != nil {
		if err := validateDeviceName(*req.Name); err != nil {
			return nil, err
		}
		device.Name = *req.Name
	}
	if req.Description != nil {
		device.Description = *req.Description
	}
	targetProtocol := device.Protocol
	protocolChanged := false
	if req.Protocol != nil {
		if !isValidProtocol(*req.Protocol) {
			return nil, validationError(fmt.Sprintf("不支援的協議類型: %s", *req.Protocol))
		}
		targetProtocol = *req.Protocol
		protocolChanged = targetProtocol != device.Protocol
		if protocolChanged && req.ConnectionConfig == nil {
			return nil, validationError("切換協議時必須提供連線配置")
		}
	}
	connectionChanged := false
	if req.ConnectionConfig != nil {
		// 驗證連線配置
		if err := validateConnectionConfig(targetProtocol, req.ConnectionConfig); err != nil {
			return nil, connectionConfigValidationError(err)
		}
		configJSON, err := json.Marshal(req.ConnectionConfig)
		if err != nil {
			return nil, fmt.Errorf("序列化連線配置失敗: %w", err)
		}
		connectionChanged = device.ConnectionConfig != string(configJSON)
		device.ConnectionConfig = string(configJSON)
	}
	device.Protocol = targetProtocol
	if req.AcquisitionIntervalMs != nil {
		if *req.AcquisitionIntervalMs < MinAcquisitionIntervalMs || *req.AcquisitionIntervalMs > MaxAcquisitionIntervalMs {
			return nil, fieldValidationError("acquisition_interval_ms", "acquisition_interval_ms 必須介於 100 到 3600000 毫秒之間")
		}
	}
	if req.StaleThresholdMs != nil {
		if *req.StaleThresholdMs <= 0 {
			return nil, fieldValidationError("stale_threshold_ms", "stale_threshold_ms 必須大於 0，或留空使用預設值")
		}
		device.StaleThresholdMs = common.Ptr(*req.StaleThresholdMs)
	}
	if req.FrozenThresholdMs != nil {
		if *req.FrozenThresholdMs < 0 {
			return nil, fieldValidationError("frozen_threshold_ms", "frozen_threshold_ms 不能小於 0（0 表示停用）")
		}
		device.FrozenThresholdMs = *req.FrozenThresholdMs
	}
	if req.AcquisitionAutoInterval != nil {
		device.AcquisitionAutoInterval = *req.AcquisitionAutoInterval
	}

	if protocolChanged || connectionChanged {
		device.LastTestAt = nil
		device.LastTestSuccess = nil
		device.LastTestError = ""
		device.ReadinessStatus = ""
		s.connMgr.Close(id) //nolint:contextcheck // Single-device retirement intentionally completes independently of request cancellation; see TestDeviceMutationCleanupDoesNotInheritRequestCancellation.
	}

	device.UpdatedAt = time.Now()

	// 先清除舊探測結果再寫入新設定：兩個寫入無法在同一個交易內完成，順序
	// 決定了失敗時停在哪一側。先清除的話，中途失敗只會讓設備顯示為未測試；
	// 反過來則會留下「新協議 + 舊探測成功」這個本來就要消滅的狀態。
	if protocolChanged || connectionChanged {
		if err := s.repo.ClearTestResult(ctx, id); err != nil {
			return nil, fmt.Errorf("清除設備測試結果失敗: %w", err)
		}
	}
	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("更新設備失敗: %w", err)
	}
	if err := s.updateDeviceAcquisition(ctx, device, req); err != nil {
		return nil, err
	}
	if req.HasAcquisitionChanges() {
		device, err = s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("重新讀取設備採集設定失敗: %w", err)
		}
		s.hydrateDeviceAcquisitionInterval(ctx, device)
	}

	return device, nil
}

// UpdateDeviceRequest 更新設備請求
type UpdateDeviceRequest struct {
	Name                    *string              `json:"name,omitempty"`
	Description             *string              `json:"description,omitempty"`
	Protocol                *schema.ProtocolType `json:"protocol,omitempty"`
	ConnectionConfig        map[string]any       `json:"connection_config,omitempty"`
	AcquisitionIntervalMs   *int                 `json:"acquisition_interval_ms,omitempty"`
	AcquisitionAutoInterval *bool                `json:"acquisition_auto_interval,omitempty"`
	StaleThresholdMs        *int                 `json:"stale_threshold_ms,omitempty"`
	FrozenThresholdMs       *int                 `json:"frozen_threshold_ms,omitempty"`
}

// HasAcquisitionChanges 回報請求是否包含採集間隔、自動調整或異常門檻欄位。
func (r UpdateDeviceRequest) HasAcquisitionChanges() bool {
	return r.AcquisitionIntervalMs != nil || r.AcquisitionAutoInterval != nil ||
		r.StaleThresholdMs != nil || r.FrozenThresholdMs != nil
}

// Delete 刪除設備
func (s *Service) Delete(ctx context.Context, id string) error {
	// 關閉現有連線
	s.connMgr.Close(id) //nolint:contextcheck // Single-device retirement intentionally completes independently of request cancellation; see TestDeviceMutationCleanupDoesNotInheritRequestCancellation.

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
	s.hydrateDeviceAcquisitionInterval(ctx, device)
	return device, nil
}

// List 列出設備
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	devices, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出設備失敗: %w", err)
	}
	s.hydrateDevicesAcquisitionInterval(ctx, devices)
	return devices, nil
}

// hydrateDevicesAcquisitionInterval 只讀一次輪詢群組，依擁有設備對應取值間隔，避免每台設備各查一次。
func (s *Service) hydrateDevicesAcquisitionInterval(ctx context.Context, devices []*schema.Device) {
	if s == nil || s.acquisitionCoordinator == nil || len(devices) == 0 {
		return
	}
	groups, err := s.acquisitionCoordinator.ListPollingGroups(ctx)
	if err != nil {
		return
	}
	intervalByOwner := make(map[string]int, len(groups))
	for _, group := range groups {
		if group != nil && group.OwnerDeviceID != nil {
			intervalByOwner[*group.OwnerDeviceID] = group.IntervalMs
		}
	}
	for _, dev := range devices {
		if dev == nil {
			continue
		}
		if interval, ok := intervalByOwner[dev.ID]; ok {
			dev.AcquisitionIntervalMs = common.Ptr(interval)
		}
	}
}

func (s *Service) hydrateDeviceAcquisitionInterval(ctx context.Context, dev *schema.Device) {
	if s == nil || s.acquisitionCoordinator == nil || dev == nil {
		return
	}
	if group, err := s.acquisitionCoordinator.GetPollingGroupByOwnerDeviceID(ctx, dev.ID); err == nil && group != nil {
		dev.AcquisitionIntervalMs = common.Ptr(group.IntervalMs)
	}
}
