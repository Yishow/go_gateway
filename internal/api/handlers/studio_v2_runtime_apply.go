package handlers

import (
	"context"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

func resolveStudioV2RuntimeApplyStatus(ctx context.Context, deviceSvc *device.Service, deviceID string) (string, string) {
	if deviceSvc == nil || deviceID == "" {
		return "not_running", ""
	}

	record, err := deviceSvc.GetByID(ctx, deviceID)
	if err != nil {
		return "apply_failed", err.Error()
	}
	if record.Status != schema.DeviceStatusActive {
		return "not_running", ""
	}
	return "applied", ""
}

func resolveStudioV2WorkspaceRuntimeApplyStatus(ctx context.Context, deviceSvc *device.Service, deviceIDs []string) (string, string) {
	if len(deviceIDs) == 0 {
		return "not_running", ""
	}

	for _, deviceID := range deviceIDs {
		status, message := resolveStudioV2RuntimeApplyStatus(ctx, deviceSvc, deviceID)
		if status == "apply_failed" {
			return status, message
		}
		if status == "applied" {
			return status, message
		}
	}

	return "not_running", ""
}
