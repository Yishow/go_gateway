package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create_SavesPlanningButKeepsRuleDisabledWhenActivationReadinessFails(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-create-readiness")
	require.NoError(t, err)
	require.NoError(t, markDeviceProbeFailure(ctx, deviceRepo, dev.ID))

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-create-readiness",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	assert.False(t, rule.Enabled)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	derivedPoint, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.False(t, derivedPoint.Enabled)
}

func TestService_Enable_UsesReadinessContractInsteadOfDeviceStatus(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-enable-readiness")
	require.NoError(t, err)
	require.NoError(t, markDeviceProbeFailure(ctx, deviceRepo, dev.ID))

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-enable-readiness",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      false,
	})
	require.NoError(t, err)

	err = svc.Enable(ctx, rule.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "讀取探測失敗")
}

func TestService_SyncDerivedPointState_DisablesEnabledRuleWhenActivationReadinessFails(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-sync-readiness")
	require.NoError(t, err)
	require.NoError(t, markDeviceProbeSuccess(ctx, deviceRepo, dev.ID))

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rule-sync-readiness",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	require.True(t, rule.Enabled)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	derivedPoint, err := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.True(t, derivedPoint.Enabled)

	require.NoError(t, markDeviceProbeFailure(ctx, deviceRepo, dev.ID))
	require.NoError(t, svc.SyncDerivedPointState(ctx))

	derivedPoint, err = pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, err)
	assert.False(t, derivedPoint.Enabled)
}

func markDeviceProbeFailure(ctx context.Context, repo *device.MemoryRepository, id string) error {
	record, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	success := false
	record.LastTestSuccess = &success
	record.LastTestError = "讀取探測失敗: host=127.0.0.1 port=502 slave=1 function=03 address=49999, read 失敗: timeout"
	return repo.Update(ctx, record)
}

func markDeviceProbeSuccess(ctx context.Context, repo *device.MemoryRepository, id string) error {
	record, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	success := true
	record.LastTestSuccess = &success
	record.LastTestError = ""
	return repo.Update(ctx, record)
}

// A device that passed its diagnostics but has not been activated yet must not
// have its derived points switched off by a rule edit: nothing re-enables them
// at activation, so the device would then be "activated" with nothing to collect.
func TestService_Update_KeepsPointsEnabledForReadyDeviceThatIsNotActiveYet(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	svc := NewService(NewMemoryRepository(), deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-update-before-activation")
	require.NoError(t, err)
	require.NoError(t, deviceRepo.UpdateStatus(ctx, dev.ID, schema.DeviceStatusDraft))

	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rule-update-before-activation", DeviceID: dev.ID, StartAddress: "40001", Count: 8,
		DataType: schema.DataTypeInt16, NamingPrefix: "SRC", Enabled: true,
	})
	require.NoError(t, err)
	require.True(t, rule.Enabled)

	count := 2
	prefix := "EDITED"
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{Count: &count, NamingPrefix: &prefix})
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 2)
	for _, link := range links {
		derived, err := pointSvc.GetByID(ctx, link.PointID)
		require.NoError(t, err)
		assert.True(t, derived.Enabled, "an edit before activation must keep the point enabled")
	}
}

func TestService_Update_DisablesPointsWhenTheDeviceIsNoLongerActivationReady(t *testing.T) {
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	deviceSvc := device.NewService(deviceRepo, nil)
	svc := NewService(NewMemoryRepository(), deviceSvc, pointSvc, nil)

	dev, err := seedActiveDevice(ctx, deviceRepo, "device-update-not-ready")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rule-update-not-ready", DeviceID: dev.ID, StartAddress: "40001", Count: 2,
		DataType: schema.DataTypeInt16, NamingPrefix: "SRC", Enabled: true,
	})
	require.NoError(t, err)
	require.NoError(t, markDeviceProbeFailure(ctx, deviceRepo, dev.ID))

	prefix := "EDITED"
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})
	require.NoError(t, err)

	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	for _, link := range links {
		derived, err := pointSvc.GetByID(ctx, link.PointID)
		require.NoError(t, err)
		assert.False(t, derived.Enabled, "a device that is not activation-ready must not collect")
	}
}
