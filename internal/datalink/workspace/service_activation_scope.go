package workspace

import (
	"context"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

type activationRuntimeScopedProjectionSyncer interface {
	ApplyWorkspaceProjectionForDevices(context.Context, *RuntimeProjection, []string) error
}

type activationRuntimeScopedGroupProjectionSyncer interface {
	ApplyWorkspaceProjectionForDevicesAndPoints(context.Context, *RuntimeProjection, []string, []string) error
}

const (
	scopedActivationOperation = "activation"
	scopedActivatedMessage    = "activated"
)

// ActivateScope activates only the persisted device scope supplied by the caller.
func (s *ActivationService) ActivateScope(ctx context.Context, deviceIDs, groupIDs []string) (*ActivationResponse, error) {
	if s == nil || s.workspaceSvc == nil || s.deviceSvc == nil {
		return nil, ErrReadinessUnavailable
	}
	selection, readiness, err := s.workspaceSvc.activationReadinessScope(ctx, deviceIDs, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("evaluate scoped readiness before activation: %w", err)
	}
	if readiness != nil && readiness.BlockingCount > 0 {
		return nil, &ReadinessBlockedError{Operation: scopedActivationOperation, Summary: readiness}
	}

	response := &ActivationResponse{
		WorkspaceID: selection.record.ID,
		Results:     make([]ActivationResult, 0, len(selection.deviceIDs)),
	}
	projectionSyncer, scopedProjection := s.runtimeSync.(activationRuntimeScopedProjectionSyncer)
	groupProjectionSyncer, scopedGroupProjection := s.runtimeSync.(activationRuntimeScopedGroupProjectionSyncer)
	selectedPointIDs := make([]string, 0)
	if len(selection.groupIDs) > 0 {
		seenPointIDs := make(map[string]bool)
		for _, groupID := range selection.groupIDs {
			for _, member := range selection.groups[groupID].Members {
				pointID := strings.TrimSpace(member.PointID)
				if pointID == "" {
					return nil, fmt.Errorf("selected group %s has no persisted point identity", groupID)
				}
				if !seenPointIDs[pointID] {
					seenPointIDs[pointID] = true
					selectedPointIDs = append(selectedPointIDs, pointID)
				}
			}
		}
		if len(selectedPointIDs) == 0 {
			return nil, fmt.Errorf("selected group scope has no persisted point identity")
		}
	}
	if len(selection.groupIDs) > 0 && (!scopedProjection || !scopedGroupProjection) {
		return nil, fmt.Errorf("selected group runtime projection is unavailable")
	}
	projectionIDs := make([]string, 0, len(selection.deviceIDs))
	activatedIDs := make([]string, 0, len(selection.deviceIDs))
	alreadyActiveIDs := make([]string, 0, len(selection.deviceIDs))
	for _, deviceID := range selection.deviceIDs {
		savedDevice, err := s.deviceSvc.GetByID(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("read scoped workspace device %s: %w", deviceID, err)
		}
		if savedDevice == nil || strings.TrimSpace(savedDevice.ID) != deviceID {
			return nil, fmt.Errorf("read scoped workspace device %s: %w", deviceID, ErrReadinessScopeInvalid)
		}
		if savedDevice.Status == schema.DeviceStatusActive {
			if scopedProjection {
				projectionIDs = append(projectionIDs, deviceID)
				alreadyActiveIDs = append(alreadyActiveIDs, deviceID)
				continue
			}
			if s.runtimeSync != nil {
				if err := s.runtimeSync.UpsertDevice(ctx, savedDevice); err != nil {
					response.Results = append(response.Results, ActivationResult{
						DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: err.Error(),
					})
					continue
				}
			}
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID, Status: ActivationResultStatusSuccess, Message: "already active",
			})
			continue
		}
		if !isEligibleForFirstActivation(savedDevice) {
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: "device is not available for activation",
			})
			continue
		}
		if err := s.deviceSvc.Activate(ctx, deviceID); err != nil {
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: err.Error(),
			})
			continue
		}
		activatedDevice, err := s.deviceSvc.GetByID(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("reload scoped activated device %s: %w", deviceID, err)
		}
		if activatedDevice == nil || strings.TrimSpace(activatedDevice.ID) != deviceID {
			return nil, fmt.Errorf("reload scoped activated device %s: %w", deviceID, ErrReadinessScopeInvalid)
		}
		if s.runtimeSync != nil && !scopedProjection {
			if err := s.runtimeSync.UpsertDevice(ctx, activatedDevice); err != nil {
				s.runtimeSync.RemoveDevice(deviceID)
				message := err.Error()
				if disableErr := s.deviceSvc.Disable(ctx, deviceID); disableErr != nil {
					message = fmt.Sprintf("%s; rollback device status: %v", message, disableErr)
				}
				response.Results = append(response.Results, ActivationResult{
					DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: message,
				})
				continue
			}
		}
		if scopedProjection {
			projectionIDs = append(projectionIDs, deviceID)
			activatedIDs = append(activatedIDs, deviceID)
			continue
		}
		response.Results = append(response.Results, ActivationResult{
			DeviceID: deviceID, Status: ActivationResultStatusSuccess, Message: scopedActivatedMessage,
		})
	}

	if scopedProjection && len(projectionIDs) > 0 {
		s.applyScopedActivationProjection(ctx, projectionSyncer, groupProjectionSyncer, projectionIDs, selectedPointIDs, activatedIDs, alreadyActiveIDs, response)
	}
	if len(response.Results) == 0 {
		response.Message = noEligibleActivationMessage
	}
	return response, nil
}

func (s *ActivationService) applyScopedActivationProjection(
	ctx context.Context,
	projectionSyncer activationRuntimeScopedProjectionSyncer,
	groupProjectionSyncer activationRuntimeScopedGroupProjectionSyncer,
	projectionIDs []string,
	pointIDs []string,
	activatedIDs []string,
	alreadyActiveIDs []string,
	response *ActivationResponse,
) {
	projection, err := s.workspaceSvc.RuntimeProjection(ctx)
	if err == nil {
		if len(pointIDs) > 0 {
			err = groupProjectionSyncer.ApplyWorkspaceProjectionForDevicesAndPoints(ctx, projection, projectionIDs, pointIDs)
		} else {
			err = projectionSyncer.ApplyWorkspaceProjectionForDevices(ctx, projection, projectionIDs)
		}
	}
	if err != nil {
		for _, deviceID := range activatedIDs {
			s.runtimeSync.RemoveDevice(deviceID)
			message := err.Error()
			if rollbackErr := s.deviceSvc.Disable(ctx, deviceID); rollbackErr != nil {
				message = fmt.Sprintf("%s; rollback device status: %v", message, rollbackErr)
			}
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: message,
			})
		}
		for _, deviceID := range alreadyActiveIDs {
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID, Status: ActivationResultStatusFailed, Message: err.Error(),
			})
		}
		return
	}
	for _, deviceID := range alreadyActiveIDs {
		response.Results = append(response.Results, ActivationResult{
			DeviceID: deviceID, Status: ActivationResultStatusSuccess, Message: "already active",
		})
	}
	for _, deviceID := range activatedIDs {
		response.Results = append(response.Results, ActivationResult{
			DeviceID: deviceID, Status: ActivationResultStatusSuccess, Message: scopedActivatedMessage,
		})
	}
}
