package device

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

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

// CheckReadiness 檢查設備就緒狀態
func (s *Service) CheckReadiness(ctx context.Context, id string) (*schema.DeviceReadiness, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	readiness := &schema.DeviceReadiness{
		DeviceID: device.ID,
		Status:   "ready", // Default to ready, downgrade if issues found
		Checks:   []schema.ReadinessCheck{},
	}

	// 1. Check Device Status
	statusCheck := schema.ReadinessCheck{
		Name: "Device Status",
		Pass: device.Status == schema.DeviceStatusActive,
	}
	if !statusCheck.Pass {
		statusCheck.Message = fmt.Sprintf("Device is in %s state", device.Status)
		readiness.Status = "warning"
	} else {
		statusCheck.Message = "Device is active"
	}
	readiness.Checks = append(readiness.Checks, statusCheck)

	// 2. Check Connection Config
	// Since we validate on create/update, this should be pass if exists, but let's double check
	connCheck := schema.ReadinessCheck{
		Name: "Connection Configuration",
		Pass: device.ConnectionConfig != "",
	}
	if !connCheck.Pass {
		connCheck.Message = "Missing connection configuration"
		readiness.Status = "error"
	} else {
		connCheck.Message = "Configuration present"
	}
	readiness.Checks = append(readiness.Checks, connCheck)

	// TODO: 3. Check Points (Need Point Repository)
	// For now we skip point checks to avoid circular dependencies or need to inject PointRepo

	// Update device readiness status in DB
	readinessJSON, _ := json.Marshal(readiness)
	device.ReadinessStatus = string(readinessJSON)
	_ = s.repo.Update(ctx, device)

	return readiness, nil
}

// UpdateCollectionStats 更新設備收集統計
func (s *Service) UpdateCollectionStats(ctx context.Context, id string, success bool) error {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	device.LastCollectedAt = &now
	device.CollectionCount++
	if !success {
		device.ErrorCount++
	}

	return s.repo.Update(ctx, device)
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
