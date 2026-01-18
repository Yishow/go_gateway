// Package device 提供設備註冊與管理功能。
//
// 本套件實作設備的 CRUD 操作、連線測試和狀態管理。
package device

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/common"
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
	mu      sync.RWMutex
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

// =============================================================================
// 狀態管理
// =============================================================================

// Activate 啟用設備
func (s *Service) Activate(ctx context.Context, id string) error {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得設備失敗: %w", err)
	}

	if device.Status == schema.DeviceStatusActive {
		return nil // 已經是啟用狀態
	}

	// 測試連線
	if err := s.TestConnection(ctx, id); err != nil {
		return fmt.Errorf("連線測試失敗，無法啟用設備: %w", err)
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.DeviceStatusActive); err != nil {
		return fmt.Errorf("更新設備狀態失敗: %w", err)
	}

	return nil
}

// Disable 停用設備
func (s *Service) Disable(ctx context.Context, id string) error {
	// 關閉現有連線
	s.connMgr.Close(id)

	if err := s.repo.UpdateStatus(ctx, id, schema.DeviceStatusDisabled); err != nil {
		return fmt.Errorf("更新設備狀態失敗: %w", err)
	}

	return nil
}

// =============================================================================
// 連線測試
// =============================================================================

// TestConnection 測試設備連線
func (s *Service) TestConnection(ctx context.Context, id string) error {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得設備失敗: %w", err)
	}

	// 嘗試建立連線
	conn, err := s.connMgr.GetOrCreate(ctx, id, device.Protocol, device.ConnectionConfig)
	if err != nil {
		// 記錄測試失敗
		_ = s.repo.UpdateTestResult(ctx, id, false, err.Error())
		return err
	}

	// 執行連線測試
	err = conn.Protocol.TestConnection(ctx)
	if err != nil {
		_ = s.repo.UpdateTestResult(ctx, id, false, err.Error())
		return err
	}

	// 記錄測試成功
	_ = s.repo.UpdateTestResult(ctx, id, true, "")

	return nil
}

// TestConnectionResult 連線測試結果
type TestConnectionResult struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	LatencyMs int64     `json:"latency_ms"`
}

// TestConnectionWithResult 測試設備連線並返回詳細結果
func (s *Service) TestConnectionWithResult(ctx context.Context, id string) (*TestConnectionResult, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	start := time.Now()
	result := &TestConnectionResult{
		Timestamp: start,
	}

	// 嘗試建立連線
	conn, err := s.connMgr.GetOrCreate(ctx, id, device.Protocol, device.ConnectionConfig)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.LatencyMs = time.Since(start).Milliseconds()
		_ = s.repo.UpdateTestResult(ctx, id, false, err.Error())
		return result, nil
	}

	// 執行連線測試
	err = conn.Protocol.TestConnection(ctx)
	result.LatencyMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		_ = s.repo.UpdateTestResult(ctx, id, false, err.Error())
	} else {
		result.Success = true
		_ = s.repo.UpdateTestResult(ctx, id, true, "")
	}

	return result, nil
}

// =============================================================================
// 輔助函數
// =============================================================================

// isValidProtocol 檢查協議類型是否有效
func isValidProtocol(protocol schema.ProtocolType) bool {
	return connector.IsRegistered(protocol)
}

// validateConnectionConfig 驗證連線配置
func validateConnectionConfig(protocol schema.ProtocolType, config map[string]interface{}) error {
	// 取得協議資訊
	info, ok := connector.GetProtocolInfo(protocol)
	if !ok {
		return fmt.Errorf("未知的協議類型: %s", protocol)
	}

	// 解析配置 Schema
	var configSchema map[string]interface{}
	if err := json.Unmarshal(info.ConfigSchema, &configSchema); err != nil {
		return nil // 無 Schema 驗證
	}

	// 檢查必填欄位
	required, ok := configSchema["required"].([]interface{})
	if ok {
		for _, field := range required {
			fieldName, ok := field.(string)
			if !ok {
				continue
			}
			if _, exists := config[fieldName]; !exists {
				return fmt.Errorf("缺少必填欄位: %s", fieldName)
			}
		}
	}

	return nil
}

// generateUUID 產生 UUID
// =============================================================================
// 連線配置輔助
// =============================================================================

// ParseConnectionConfig 解析連線配置為具體類型
func ParseConnectionConfig(device *schema.Device) (interface{}, error) {
	var config interface{}

	switch device.Protocol {
	case schema.ProtocolModbusTCP:
		var c schema.ConnectionConfigModbusTCP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolModbusRTU:
		var c schema.ConnectionConfigModbusRTU
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolModbusUDP:
		var c schema.ConnectionConfigModbusUDP
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolFatekFBs:
		var c schema.ConnectionConfigFatekFBs
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolMC3E:
		var c schema.ConnectionConfigMC3E
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	case schema.ProtocolMQTT:
		var c schema.ConnectionConfigMQTT
		if err := json.Unmarshal([]byte(device.ConnectionConfig), &c); err != nil {
			return nil, err
		}
		config = c
	default:
		return nil, fmt.Errorf("未知的協議類型: %s", device.Protocol)
	}

	return config, nil
}
