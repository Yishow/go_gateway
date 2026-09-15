package device

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

const (
	// AvailabilityStatusAvailable means the device can still be treated as usable by Studio V2.
	AvailabilityStatusAvailable = "available"
	// AvailabilityStatusUnavailable means the current Studio V2 device draft makes the device unusable.
	AvailabilityStatusUnavailable = "unavailable"
)

// AvailabilityOf resolves the explicit Studio V2 availability state stored on a device record.
func AvailabilityOf(record *schema.Device) (availabilityStatus, availabilityReason string) {
	readiness := decodeDeviceReadiness(record)
	if readiness == nil || strings.TrimSpace(readiness.AvailabilityStatus) == "" {
		return AvailabilityStatusAvailable, ""
	}

	status := strings.TrimSpace(readiness.AvailabilityStatus)
	if status != AvailabilityStatusUnavailable {
		return AvailabilityStatusAvailable, ""
	}
	return AvailabilityStatusUnavailable, strings.TrimSpace(readiness.AvailabilityReason)
}

// SetAvailability updates the explicit Studio V2 availability state on a device record.
func (s *Service) SetAvailability(ctx context.Context, id, status, reason string) (*schema.Device, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得設備失敗: %w", err)
	}

	nextStatus := strings.TrimSpace(status)
	switch nextStatus {
	case "", AvailabilityStatusAvailable:
		nextStatus = AvailabilityStatusAvailable
		reason = ""
	case AvailabilityStatusUnavailable:
		reason = strings.TrimSpace(reason)
		if reason == "" {
			return nil, validationError("availability_reason 不能為空")
		}
	default:
		return nil, validationError("不支援的 availability_status")
	}

	readiness := decodeDeviceReadiness(record)
	if readiness == nil {
		readiness = &schema.DeviceReadiness{
			DeviceID:           record.ID,
			ConnectStatus:      schema.ReadinessStageStatusUnknown,
			ProbeStatus:        schema.ReadinessStageStatusUnknown,
			AvailabilityStatus: AvailabilityStatusAvailable,
		}
	}
	readiness.AvailabilityStatus = nextStatus
	readiness.AvailabilityReason = reason

	readinessJSON, err := json.Marshal(readiness)
	if err != nil {
		return nil, fmt.Errorf("序列化設備可用狀態失敗: %w", err)
	}
	record.ReadinessStatus = string(readinessJSON)
	if err := s.repo.Update(ctx, record); err != nil {
		return nil, fmt.Errorf("更新設備可用狀態失敗: %w", err)
	}

	return record, nil
}

func decodeDeviceReadiness(record *schema.Device) *schema.DeviceReadiness {
	if record == nil || strings.TrimSpace(record.ReadinessStatus) == "" {
		return nil
	}

	var readiness schema.DeviceReadiness
	if err := json.Unmarshal([]byte(record.ReadinessStatus), &readiness); err != nil {
		return nil
	}
	if readiness.DeviceID == "" {
		readiness.DeviceID = record.ID
	}
	if strings.TrimSpace(readiness.AvailabilityStatus) == "" {
		readiness.AvailabilityStatus = AvailabilityStatusAvailable
	}
	return &readiness
}
