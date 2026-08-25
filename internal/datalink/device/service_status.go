package device

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 狀態管理
// =============================================================================

// Activate 啟用設備
func (s *Service) Activate(ctx context.Context, id string) error {
	return s.ProbeAndActivate(ctx, id)
}

// CheckReadiness 檢查設備就緒狀態
func (s *Service) CheckReadiness(ctx context.Context, id string) (*schema.DeviceReadiness, error) {
	device, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	connectStatus, probeStatus, blockingReasons := deriveReadinessDiagnostics(device)
	readiness := &schema.DeviceReadiness{
		DeviceID:           device.ID,
		AvailabilityStatus: AvailabilityStatusAvailable,
		ConnectStatus:      connectStatus,
		ProbeStatus:        probeStatus,
		PlanningAllowed:    connectStatus == schema.ReadinessStageStatusSuccess,
		ActivationAllowed:  connectStatus == schema.ReadinessStageStatusSuccess && probeStatus == schema.ReadinessStageStatusSuccess,
		ApplyAllowed:       connectStatus == schema.ReadinessStageStatusSuccess && probeStatus == schema.ReadinessStageStatusSuccess,
		BlockingReasons:    blockingReasons,
	}
	if currentReadiness := decodeDeviceReadiness(device); currentReadiness != nil && currentReadiness.AvailabilityStatus == AvailabilityStatusUnavailable {
		readiness.AvailabilityStatus = AvailabilityStatusUnavailable
		readiness.AvailabilityReason = currentReadiness.AvailabilityReason
	}
	readiness.Status = deriveLegacyReadinessStatus(readiness)
	readiness.Checks = buildReadinessChecks(device, readiness)

	// Update device readiness status in DB
	readinessJSON, err := json.Marshal(readiness)
	if err != nil {
		return nil, fmt.Errorf("序列化設備就緒狀態失敗: %w", err)
	}
	device.ReadinessStatus = string(readinessJSON)
	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("更新設備就緒狀態失敗: %w", err)
	}

	return readiness, nil
}

func deriveReadinessDiagnostics(device *schema.Device) (
	schema.ReadinessStageStatus,
	schema.ReadinessStageStatus,
	[]string,
) {
	if strings.TrimSpace(device.ConnectionConfig) == "" {
		return schema.ReadinessStageStatusUnknown, schema.ReadinessStageStatusUnknown, []string{
			"missing connection configuration",
		}
	}

	if device.LastTestSuccess == nil {
		return schema.ReadinessStageStatusUnknown, schema.ReadinessStageStatusUnknown, []string{
			"connect diagnostics have not succeeded yet",
		}
	}

	if *device.LastTestSuccess {
		probeStatus, reasons := deriveSuccessfulProbeStatus(device)
		return schema.ReadinessStageStatusSuccess, probeStatus, reasons
	}

	if isProbeFailure(device.LastTestError) {
		return schema.ReadinessStageStatusSuccess, schema.ReadinessStageStatusFailed, []string{
			safeReadinessDiagnostic(firstNonEmpty(device.LastTestError, "probe diagnostics failed")),
		}
	}

	return schema.ReadinessStageStatusFailed, schema.ReadinessStageStatusUnknown, []string{
		safeReadinessDiagnostic(firstNonEmpty(device.LastTestError, "connect diagnostics failed")),
	}
}

func deriveSuccessfulProbeStatus(device *schema.Device) (schema.ReadinessStageStatus, []string) {
	target, err := buildReadProbeTarget(device)
	if err != nil {
		return schema.ReadinessStageStatusUnknown, []string{
			fmt.Sprintf("probe diagnostics unavailable: %v", err),
		}
	}
	if !target.enabled {
		return schema.ReadinessStageStatusSkipped, []string{
			"probe diagnostics are not supported for this device",
		}
	}
	return schema.ReadinessStageStatusSuccess, nil
}

func deriveLegacyReadinessStatus(readiness *schema.DeviceReadiness) string {
	switch {
	case readiness.ApplyAllowed:
		return "ready"
	case readiness.PlanningAllowed:
		return "warning"
	default:
		return "error"
	}
}

func buildReadinessChecks(device *schema.Device, readiness *schema.DeviceReadiness) []schema.ReadinessCheck {
	checks := []schema.ReadinessCheck{
		{
			Name:    "Connection Configuration",
			Pass:    strings.TrimSpace(device.ConnectionConfig) != "",
			Message: readinessConfigurationMessage(device),
		},
		{
			Name:    "Connect Diagnostics",
			Pass:    readiness.ConnectStatus == schema.ReadinessStageStatusSuccess,
			Message: readinessStageMessage("connect", readiness.ConnectStatus, device.LastTestError),
		},
		{
			Name:    "Probe Diagnostics",
			Pass:    readiness.ProbeStatus == schema.ReadinessStageStatusSuccess,
			Message: readinessStageMessage("probe", readiness.ProbeStatus, device.LastTestError),
		},
		{
			Name:    "Planning Eligibility",
			Pass:    readiness.PlanningAllowed,
			Message: readinessEligibilityMessage("planning", readiness.PlanningAllowed, readiness.BlockingReasons),
		},
		{
			Name:    "Activation Eligibility",
			Pass:    readiness.ActivationAllowed,
			Message: readinessEligibilityMessage("activation", readiness.ActivationAllowed, readiness.BlockingReasons),
		},
		{
			Name:    "Apply Eligibility",
			Pass:    readiness.ApplyAllowed,
			Message: readinessEligibilityMessage("apply", readiness.ApplyAllowed, readiness.BlockingReasons),
		},
	}
	return checks
}

func readinessConfigurationMessage(device *schema.Device) string {
	if strings.TrimSpace(device.ConnectionConfig) == "" {
		return "missing connection configuration"
	}
	return "configuration present"
}

func readinessStageMessage(stage string, status schema.ReadinessStageStatus, failure string) string {
	switch status {
	case schema.ReadinessStageStatusSuccess:
		return stage + " diagnostics succeeded"
	case schema.ReadinessStageStatusFailed:
		return safeReadinessDiagnostic(firstNonEmpty(failure, stage+" diagnostics failed"))
	case schema.ReadinessStageStatusSkipped:
		return stage + " diagnostics skipped"
	default:
		return stage + " diagnostics not yet successful"
	}
}

func readinessEligibilityMessage(stage string, allowed bool, reasons []string) string {
	if allowed {
		return stage + " allowed"
	}
	if len(reasons) == 0 {
		return stage + " blocked"
	}
	return stage + " blocked: " + reasons[0]
}

func isProbeFailure(message string) bool {
	return strings.Contains(message, "讀取探測失敗")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
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
