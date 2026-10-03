package runtime

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

// ReconcileSourceRule applies the latest workspace projection after a source-rule lifecycle change.
func (s *Service) ReconcileSourceRule(ctx context.Context, req sourcerule.RuntimeReconcileRequest) sourcerule.RuntimeReconcileOutcome {
	if s == nil || !s.IsRunning() {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusNotRunning,
			Scope:   req.Scope,
			Message: "runtime is not running",
		}
	}
	if s.workspace == nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusDeferred,
			Scope:   req.Scope,
			Message: "workspace runtime projection unavailable",
		}
	}

	projection, err := s.workspace.RuntimeProjection(ctx)
	if err != nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusStale,
			Scope:   req.Scope,
			Code:    "runtime_projection_unavailable",
			Message: "workspace runtime projection is unavailable",
		}
	}
	if sourceRuleProjectionDeviceInactive(projection, req.Scope.DeviceID) {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusDeferred,
			Scope:   req.Scope,
			Message: "source rule device is not active",
		}
	}
	deviceStatus, statusFound := sourceRuleProjectionDeviceStatus(projection, req.Scope.DeviceID)
	if !statusFound || !sourceRuleProjectionDeviceStatusKnown(deviceStatus) {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusStale,
			Scope:   req.Scope,
			Code:    "runtime_projection_device_status_invalid",
			Message: "source rule device status is invalid in workspace runtime projection",
		}
	}
	if err := s.applySourceRuleProjection(ctx, projection, req); err != nil {
		return sourcerule.RuntimeReconcileOutcome{
			Status:  sourcerule.RuntimeReconcileStatusStale,
			Scope:   req.Scope,
			Code:    "runtime_projection_apply_failed",
			Message: "workspace runtime projection could not be applied",
		}
	}

	return sourcerule.RuntimeReconcileOutcome{
		Status: sourcerule.RuntimeReconcileStatusAligned,
		Scope:  req.Scope,
	}
}

func sourceRuleProjectionDeviceInactive(projection *workspace.RuntimeProjection, deviceID string) bool {
	status, ok := sourceRuleProjectionDeviceStatus(projection, deviceID)
	return ok && (status == schema.DeviceStatusDraft || status == schema.DeviceStatusDisabled)
}

func sourceRuleProjectionDeviceStatus(projection *workspace.RuntimeProjection, deviceID string) (schema.DeviceStatus, bool) {
	if projection == nil || strings.TrimSpace(deviceID) == "" {
		return "", false
	}
	for _, id := range projection.DeviceIDs {
		if id != deviceID {
			continue
		}
		for _, device := range projection.Devices {
			if device != nil && device.ID == deviceID {
				return device.Status, true
			}
		}
		return "", false
	}
	return "", false
}

func sourceRuleProjectionDeviceStatusKnown(status schema.DeviceStatus) bool {
	switch status {
	case schema.DeviceStatusActive, schema.DeviceStatusDraft, schema.DeviceStatusDisabled:
		return true
	default:
		return false
	}
}

func (s *Service) applySourceRuleProjection(
	ctx context.Context,
	projection *workspace.RuntimeProjection,
	req sourcerule.RuntimeReconcileRequest,
) error {
	if req.Operation == sourcerule.RuntimeReconcileOperationCreate {
		pointIDs := sourceRuleProjectionPointIDs(projection, req.Scope.RuleID)
		if len(pointIDs) > 0 {
			return s.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, []string{req.Scope.DeviceID}, pointIDs)
		}
		return nil
	}
	return s.ApplyWorkspaceProjectionForDevices(ctx, projection, []string{req.Scope.DeviceID})
}

func sourceRuleProjectionPointIDs(projection *workspace.RuntimeProjection, ruleID string) []string {
	if projection == nil || strings.TrimSpace(ruleID) == "" {
		return nil
	}
	seen := make(map[string]struct{})
	pointIDs := make([]string, 0)
	for _, link := range projection.RuleLinks {
		if link == nil || link.RuleID != ruleID {
			continue
		}
		pointID := strings.TrimSpace(link.PointID)
		if pointID == "" {
			continue
		}
		if _, ok := seen[pointID]; ok {
			continue
		}
		seen[pointID] = struct{}{}
		pointIDs = append(pointIDs, pointID)
	}
	return pointIDs
}
