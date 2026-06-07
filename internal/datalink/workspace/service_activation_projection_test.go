package workspace

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestActivationService_ActivateEligibleAppliesPersistedRuntimeProjection(t *testing.T) {
	ctx := context.Background()
	workspaceSvc := NewService(NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(ctx, "dev-A"); err != nil {
		t.Fatalf("attach device failed: %v", err)
	}

	deviceSvc := &stubActivationDeviceService{
		devices: map[string]*schema.Device{
			"dev-A": {
				ID:     "dev-A",
				Name:   "Device A",
				Status: schema.DeviceStatusDraft,
			},
		},
		activateErrs: map[string]error{},
	}
	workspaceSvc.WithReadinessServices(deviceSvc, nil, nil, nil)
	workspaceSvc.WithRuntimeProjectionServices(
		deviceSvc,
		&runtimeProjectionRuleStub{},
		&runtimeProjectionPointStub{},
		&runtimeProjectionMappingStub{},
		&runtimeProjectionTagStub{},
		nil,
		nil,
		&runtimeProjectionGroupStub{},
	)
	runtimeSync := &projectionActivationRuntimeSyncer{}

	service := NewActivationService(workspaceSvc, deviceSvc, runtimeSync)
	response, err := service.ActivateEligible(ctx)
	if err != nil {
		t.Fatalf("activate eligible failed: %v", err)
	}

	if len(response.Results) != 1 || response.Results[0].Status != ActivationResultStatusSuccess {
		t.Fatalf("expected one successful activation result, got %+v", response.Results)
	}
	if len(runtimeSync.upserted) != 0 {
		t.Fatalf("expected projection sync instead of per-device upsert, got upserts %v", runtimeSync.upserted)
	}
	if len(runtimeSync.applied) != 1 {
		t.Fatalf("expected one applied runtime projection, got %d", len(runtimeSync.applied))
	}
	projection := runtimeSync.applied[0]
	if projection.WorkspaceID == "" || projection.Version == "" {
		t.Fatalf("expected persisted projection identity and version, got %+v", projection)
	}
	if len(projection.Devices) != 1 || projection.Devices[0].ID != "dev-A" || projection.Devices[0].Status != schema.DeviceStatusActive {
		t.Fatalf("expected projection to include activated dev-A, got %+v", projection.Devices)
	}
}

type projectionActivationRuntimeSyncer struct {
	stubActivationRuntimeSyncer
	applied []*RuntimeProjection
}

func (s *projectionActivationRuntimeSyncer) ApplyWorkspaceProjection(_ context.Context, projection *RuntimeProjection) error {
	s.applied = append(s.applied, projection)
	return nil
}
