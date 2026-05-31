package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

type stubActivationDeviceService struct {
	devices       map[string]*schema.Device
	activateCalls []string
	disableCalls  []string
	activateErrs  map[string]error
}

func (s *stubActivationDeviceService) GetByID(_ context.Context, id string) (*schema.Device, error) {
	record, ok := s.devices[id]
	if !ok {
		return nil, errors.New("device not found")
	}
	cloned := *record
	return &cloned, nil
}

func (s *stubActivationDeviceService) Activate(_ context.Context, id string) error {
	s.activateCalls = append(s.activateCalls, id)
	if err := s.activateErrs[id]; err != nil {
		return err
	}
	record := s.devices[id]
	record.Status = schema.DeviceStatusActive
	return nil
}

func (s *stubActivationDeviceService) Disable(_ context.Context, id string) error {
	s.disableCalls = append(s.disableCalls, id)
	record := s.devices[id]
	record.Status = schema.DeviceStatusDisabled
	return nil
}

type stubActivationRuntimeSyncer struct {
	upserted   []string
	removed    []string
	upsertErrs map[string]error
}

func (s *stubActivationRuntimeSyncer) UpsertDevice(_ context.Context, record *schema.Device) error {
	if record != nil {
		s.upserted = append(s.upserted, record.ID)
	}
	if record == nil {
		return nil
	}
	return s.upsertErrs[record.ID]
}

func (s *stubActivationRuntimeSyncer) RemoveDevice(deviceID string) {
	s.removed = append(s.removed, deviceID)
}

func TestActivationService_ActivateEligibleSkipsAlreadyRunningDevices(t *testing.T) {
	workspaceSvc := NewService(NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-running"); err != nil {
		t.Fatalf("attach running device failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-draft"); err != nil {
		t.Fatalf("attach draft device failed: %v", err)
	}

	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-running": {
				ID:     "dev-running",
				Name:   "Running PLC",
				Status: schema.DeviceStatusActive,
			},
			"dev-draft": {
				ID:     "dev-draft",
				Name:   "Draft PLC",
				Status: schema.DeviceStatusDraft,
			},
		},
		activateErrs: map[string]error{},
	}

	service := NewActivationService(workspaceSvc, deviceSvc)
	response, err := service.ActivateEligible(context.Background())
	if err != nil {
		t.Fatalf("activate eligible failed: %v", err)
	}

	if len(response.Results) != 1 {
		t.Fatalf("expected exactly one attempted device, got %+v", response.Results)
	}
	if response.Results[0].DeviceID != "dev-draft" || response.Results[0].Status != ActivationResultStatusSuccess {
		t.Fatalf("expected dev-draft success result, got %+v", response.Results[0])
	}
	if slices.Contains(deviceSvc.activateCalls, "dev-running") {
		t.Fatalf("expected already-running device to be skipped, got activate calls %v", deviceSvc.activateCalls)
	}
	if !slices.Contains(deviceSvc.activateCalls, "dev-draft") {
		t.Fatalf("expected draft device activation attempt, got %v", deviceSvc.activateCalls)
	}
}

func TestActivationService_ActivateEligibleReturnsActionableEmptyResultWhenNoDevicesCanStart(t *testing.T) {
	workspaceSvc := NewService(NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-running"); err != nil {
		t.Fatalf("attach running device failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-unavailable"); err != nil {
		t.Fatalf("attach unavailable device failed: %v", err)
	}

	readiness := schema.DeviceReadiness{
		DeviceID:           "dev-unavailable",
		AvailabilityStatus: device.AvailabilityStatusUnavailable,
		AvailabilityReason: "device form is invalid",
	}
	readinessJSON, err := json.Marshal(readiness)
	if err != nil {
		t.Fatalf("marshal readiness failed: %v", err)
	}

	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-running": {
				ID:     "dev-running",
				Name:   "Running PLC",
				Status: schema.DeviceStatusActive,
			},
			"dev-unavailable": {
				ID:              "dev-unavailable",
				Name:            "Broken PLC",
				Status:          schema.DeviceStatusDraft,
				ReadinessStatus: string(readinessJSON),
			},
		},
		activateErrs: map[string]error{},
	}

	service := NewActivationService(workspaceSvc, deviceSvc)
	response, err := service.ActivateEligible(context.Background())
	if err != nil {
		t.Fatalf("activate eligible failed: %v", err)
	}

	if len(response.Results) != 0 {
		t.Fatalf("expected no attempted devices, got %+v", response.Results)
	}
	if response.Message == "" {
		t.Fatalf("expected actionable empty result message, got %+v", response)
	}
	if len(deviceSvc.activateCalls) != 0 {
		t.Fatalf("expected no activation calls, got %v", deviceSvc.activateCalls)
	}
}
