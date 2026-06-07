package workspace

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

type stubActivationDeviceService struct {
	devices       map[string]*schema.Device
	readiness     map[string]*schema.DeviceReadiness
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

func (s *stubActivationDeviceService) CheckReadiness(_ context.Context, id string) (*schema.DeviceReadiness, error) {
	if s.readiness != nil {
		readiness, ok := s.readiness[id]
		if !ok {
			return nil, errors.New("readiness not found")
		}
		cloned := *readiness
		cloned.BlockingReasons = append([]string{}, readiness.BlockingReasons...)
		cloned.Checks = append([]schema.ReadinessCheck{}, readiness.Checks...)
		return &cloned, nil
	}

	record, ok := s.devices[id]
	if !ok {
		return nil, errors.New("device not found")
	}
	status, reason := device.AvailabilityOf(record)
	readiness := &schema.DeviceReadiness{
		DeviceID:           id,
		PlanningAllowed:    status == device.AvailabilityStatusAvailable,
		ActivationAllowed:  status == device.AvailabilityStatusAvailable,
		ApplyAllowed:       status == device.AvailabilityStatusAvailable,
		AvailabilityStatus: status,
		AvailabilityReason: reason,
	}
	if readiness.ActivationAllowed {
		readiness.ConnectStatus = schema.ReadinessStageStatusSuccess
		readiness.ProbeStatus = schema.ReadinessStageStatusSuccess
		return readiness, nil
	}
	readiness.ConnectStatus = schema.ReadinessStageStatusFailed
	readiness.ProbeStatus = schema.ReadinessStageStatusUnknown
	if strings.TrimSpace(reason) != "" {
		readiness.BlockingReasons = []string{reason}
	}
	return readiness, nil
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
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)

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

func TestActivationService_ActivateEligibleRejectsWorkspaceWhenNoDevicesAreReady(t *testing.T) {
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
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)

	service := NewActivationService(workspaceSvc, deviceSvc)
	response, err := service.ActivateEligible(context.Background())
	if response != nil {
		t.Fatalf("expected nil response when workspace readiness blocks activation, got %+v", response)
	}

	var blockedErr *ReadinessBlockedError
	if !errors.As(err, &blockedErr) {
		t.Fatalf("expected ReadinessBlockedError, got %v", err)
	}
	issues := blockedErr.BlockingIssues()
	if len(issues) != 1 {
		t.Fatalf("expected one blocking issue, got %+v", issues)
	}
	if issues[0].Scope != "dev-unavailable" {
		t.Fatalf("expected blocker scope dev-unavailable, got %+v", issues[0])
	}
	if len(deviceSvc.activateCalls) != 0 {
		t.Fatalf("expected no activation calls, got %v", deviceSvc.activateCalls)
	}
}

func TestActivationService_ActivateEligibleRejectsBlockingWorkspaceReadiness(t *testing.T) {
	workspaceSvc := NewService(NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-blocked"); err != nil {
		t.Fatalf("attach blocked device failed: %v", err)
	}

	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-blocked": {
				ID:     "dev-blocked",
				Name:   "Blocked PLC",
				Status: schema.DeviceStatusDraft,
			},
		},
		readiness: map[string]*schema.DeviceReadiness{
			"dev-blocked": {
				DeviceID:          "dev-blocked",
				ConnectStatus:     schema.ReadinessStageStatusSuccess,
				ProbeStatus:       schema.ReadinessStageStatusFailed,
				PlanningAllowed:   true,
				ActivationAllowed: false,
				ApplyAllowed:      false,
				BlockingReasons:   []string{"讀取探測失敗: timeout"},
			},
		},
		activateErrs: map[string]error{},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	runtimeSync := &stubActivationRuntimeSyncer{}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateEligible(context.Background())
	if response != nil {
		t.Fatalf("expected nil activation response when readiness blocks activation, got %+v", response)
	}

	var blockedErr *ReadinessBlockedError
	if !errors.As(err, &blockedErr) {
		t.Fatalf("expected ReadinessBlockedError, got %v", err)
	}
	issues := blockedErr.BlockingIssues()
	if len(issues) != 1 {
		t.Fatalf("expected one blocking issue, got %+v", issues)
	}
	if issues[0].Code != "device-probe-required" || issues[0].Scope != "dev-blocked" {
		t.Fatalf("expected device-probe-required for dev-blocked, got %+v", issues[0])
	}
	if len(deviceSvc.activateCalls) != 0 {
		t.Fatalf("expected no activation attempts, got %v", deviceSvc.activateCalls)
	}
	if len(runtimeSync.upserted) != 0 {
		t.Fatalf("expected no runtime sync attempts, got %v", runtimeSync.upserted)
	}
}
