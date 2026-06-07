package workspace

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

const noEligibleActivationMessage = "目前沒有可啟動設備"

type ActivationResultStatus string

const (
	ActivationResultStatusSuccess ActivationResultStatus = "success"
	ActivationResultStatusFailed  ActivationResultStatus = "failed"
)

type ActivationResult struct {
	DeviceID string                 `json:"device_id"`
	Status   ActivationResultStatus `json:"status"`
	Message  string                 `json:"message"`
}

type ActivationResponse struct {
	WorkspaceID string             `json:"workspace_id"`
	Results     []ActivationResult `json:"results"`
	Message     string             `json:"message,omitempty"`
}

type activationDeviceService interface {
	GetByID(ctx context.Context, id string) (*schema.Device, error)
	Activate(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error
}

type activationRuntimeSyncer interface {
	UpsertDevice(ctx context.Context, device *schema.Device) error
	RemoveDevice(deviceID string)
}

type activationRuntimeProjectionSyncer interface {
	ApplyWorkspaceProjection(ctx context.Context, projection *RuntimeProjection) error
}

type ActivationService struct {
	workspaceSvc *Service
	deviceSvc    activationDeviceService
	runtimeSync  activationRuntimeSyncer
}

func NewActivationService(workspaceSvc *Service, deviceSvc activationDeviceService, runtimeSync ...activationRuntimeSyncer) *ActivationService {
	var syncer activationRuntimeSyncer
	if len(runtimeSync) > 0 {
		syncer = runtimeSync[0]
	}
	return &ActivationService{
		workspaceSvc: workspaceSvc,
		deviceSvc:    deviceSvc,
		runtimeSync:  syncer,
	}
}

func (s *ActivationService) ActivateEligible(ctx context.Context) (*ActivationResponse, error) {
	readiness, err := s.workspaceSvc.Readiness(ctx)
	if err != nil {
		return nil, fmt.Errorf("evaluate workspace readiness before activation: %w", err)
	}
	if readiness != nil && readiness.BlockingCount > 0 {
		return nil, &ReadinessBlockedError{
			Operation: "activation",
			Summary:   readiness,
		}
	}

	record, err := s.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return nil, err
	}

	response := &ActivationResponse{
		WorkspaceID: record.ID,
		Results:     []ActivationResult{},
	}
	projectionSyncer, syncProjection := s.runtimeSync.(activationRuntimeProjectionSyncer)
	activatedIDs := make([]string, 0, len(record.OrderedDeviceIDs))

	for _, deviceID := range record.OrderedDeviceIDs {
		savedDevice, err := s.deviceSvc.GetByID(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("read workspace device %s: %w", deviceID, err)
		}
		if !isEligibleForFirstActivation(savedDevice) {
			continue
		}

		if err := s.deviceSvc.Activate(ctx, deviceID); err != nil {
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID,
				Status:   ActivationResultStatusFailed,
				Message:  err.Error(),
			})
			continue
		}

		activatedDevice, err := s.deviceSvc.GetByID(ctx, deviceID)
		if err != nil {
			return nil, fmt.Errorf("reload activated device %s: %w", deviceID, err)
		}
		if s.runtimeSync != nil && !syncProjection {
			if err := s.runtimeSync.UpsertDevice(ctx, activatedDevice); err != nil {
				s.runtimeSync.RemoveDevice(deviceID)
				if disableErr := s.deviceSvc.Disable(ctx, deviceID); disableErr != nil {
					err = fmt.Errorf("%w; rollback device status: %w", err, disableErr)
				}
				response.Results = append(response.Results, ActivationResult{
					DeviceID: deviceID,
					Status:   ActivationResultStatusFailed,
					Message:  err.Error(),
				})
				continue
			}
		}
		if syncProjection {
			activatedIDs = append(activatedIDs, deviceID)
			continue
		}

		response.Results = append(response.Results, ActivationResult{
			DeviceID: deviceID,
			Status:   ActivationResultStatusSuccess,
			Message:  "activated",
		})
	}

	if syncProjection && len(activatedIDs) > 0 {
		s.applyActivationProjection(ctx, projectionSyncer, activatedIDs, response)
	}

	if len(response.Results) == 0 {
		response.Message = noEligibleActivationMessage
	}

	return response, nil
}

func (s *ActivationService) applyActivationProjection(
	ctx context.Context,
	projectionSyncer activationRuntimeProjectionSyncer,
	activatedIDs []string,
	response *ActivationResponse,
) {
	projection, err := s.workspaceSvc.RuntimeProjection(ctx)
	if err == nil {
		err = projectionSyncer.ApplyWorkspaceProjection(ctx, projection)
	}
	if err != nil {
		for _, deviceID := range activatedIDs {
			s.runtimeSync.RemoveDevice(deviceID)
			rollbackErr := s.deviceSvc.Disable(ctx, deviceID)
			message := err.Error()
			if rollbackErr != nil {
				message = fmt.Sprintf("%s; rollback device status: %v", message, rollbackErr)
			}
			response.Results = append(response.Results, ActivationResult{
				DeviceID: deviceID,
				Status:   ActivationResultStatusFailed,
				Message:  message,
			})
		}
		return
	}

	for _, deviceID := range activatedIDs {
		response.Results = append(response.Results, ActivationResult{
			DeviceID: deviceID,
			Status:   ActivationResultStatusSuccess,
			Message:  "activated",
		})
	}
}

func isEligibleForFirstActivation(savedDevice *schema.Device) bool {
	if savedDevice == nil || savedDevice.Status == schema.DeviceStatusActive {
		return false
	}

	status, _ := device.AvailabilityOf(savedDevice)
	return status == device.AvailabilityStatusAvailable
}
