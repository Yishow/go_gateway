package point_test

import (
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestFailedIntervalAssignmentDoesNotSaveOtherAcquisitionSettings(t *testing.T) {
	db, points := acquisitionSQLiteService(t)
	ctx := t.Context()
	shared, err := points.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
	require.NoError(t, err)
	for _, item := range []struct{ name, address string }{{"A1", "40001"}, {"A2", "40002"}} {
		_, err = points.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: item.name, Address: item.address, DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
		require.NoError(t, err)
	}
	_, err = db.ExecContext(ctx, `CREATE TRIGGER review_fail_assignment BEFORE UPDATE OF polling_group_id ON points WHEN OLD.device_id = 'device-a' AND OLD.name = 'A2' BEGIN SELECT RAISE(ABORT, 'review injected assignment failure'); END`)
	require.NoError(t, err)
	repo := device.NewSQLRepository(db)
	devices := device.NewService(repo, nil)
	devices.SetAcquisitionIntervalCoordinator(points)
	interval, enabled := 500, true
	_, err = devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval, AcquisitionAutoInterval: &enabled})
	require.ErrorContains(t, err, "review injected assignment failure")
	got, err := devices.GetByID(ctx, "device-a")
	require.NoError(t, err)
	require.False(t, got.AcquisitionAutoInterval, "a rejected acquisition save must not silently persist auto=true")
	group, err := points.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
	require.NoError(t, err)
	require.Nil(t, group, "the failed assignment must not save an owned group")
}

func TestFailedAcquisitionSettingsRollBackIntervalAndPointAssignment(t *testing.T) {
	db, points := acquisitionSQLiteService(t)
	ctx := t.Context()
	shared, err := points.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
	require.NoError(t, err)
	pt, err := points.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "A1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TRIGGER fail_auto_setting BEFORE UPDATE OF acquisition_auto_interval ON devices WHEN OLD.id = 'device-a' BEGIN SELECT RAISE(ABORT, 'injected setting failure'); END`)
	require.NoError(t, err)
	devices := device.NewService(device.NewSQLRepository(db), nil)
	devices.SetAcquisitionIntervalCoordinator(points)
	interval, enabled := 500, true
	_, err = devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval, AcquisitionAutoInterval: &enabled})
	require.ErrorContains(t, err, "injected setting failure")
	group, err := points.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
	require.NoError(t, err)
	require.Nil(t, group)
	gotPoint, err := points.GetByID(ctx, pt.ID)
	require.NoError(t, err)
	require.Equal(t, &shared.ID, gotPoint.PollingGroupID)
	require.Equal(t, pt.RevisionID, gotPoint.RevisionID)
	gotDevice, err := devices.GetByID(ctx, "device-a")
	require.NoError(t, err)
	require.False(t, gotDevice.AcquisitionAutoInterval)
}

func TestCombinedAcquisitionSettingsCommitWithOwnedGroup(t *testing.T) {
	db, points := acquisitionSQLiteService(t)
	ctx := t.Context()
	shared, err := points.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
	require.NoError(t, err)
	pt, err := points.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "A1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)
	other, err := points.Create(ctx, point.CreatePointRequest{DeviceID: "device-b", Name: "B1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)
	devices := device.NewService(device.NewSQLRepository(db), nil)
	devices.SetAcquisitionIntervalCoordinator(points)
	interval, enabled, frozen := 500, true, 600000
	updated, err := devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval, AcquisitionAutoInterval: &enabled, FrozenThresholdMs: &frozen})
	require.NoError(t, err)
	require.Equal(t, &interval, updated.AcquisitionIntervalMs)
	require.True(t, updated.AcquisitionAutoInterval)
	require.Equal(t, frozen, updated.FrozenThresholdMs)
	group, err := points.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
	require.NoError(t, err)
	require.Equal(t, interval, group.IntervalMs)
	gotPoint, err := points.GetByID(ctx, pt.ID)
	require.NoError(t, err)
	require.Equal(t, &group.ID, gotPoint.PollingGroupID)
	gotOther, err := points.GetByID(ctx, other.ID)
	require.NoError(t, err)
	require.Equal(t, other.PollingGroupID, gotOther.PollingGroupID)
	require.Equal(t, other.RevisionID, gotOther.RevisionID)
}

func TestCombinedAcquisitionSettingsRejectDifferentConfigDatabases(t *testing.T) {
	deviceDB, _ := acquisitionSQLiteService(t)
	_, points := acquisitionSQLiteService(t)
	devices := device.NewService(device.NewSQLRepository(deviceDB), nil)
	devices.SetAcquisitionIntervalCoordinator(points)
	interval, enabled := 500, true
	_, err := devices.Update(t.Context(), "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval, AcquisitionAutoInterval: &enabled})
	require.ErrorContains(t, err, "共用設定資料庫")
	group, err := points.GetPollingGroupByOwnerDeviceID(t.Context(), "device-a")
	require.NoError(t, err)
	require.Nil(t, group)
	got, err := devices.GetByID(t.Context(), "device-a")
	require.NoError(t, err)
	require.False(t, got.AcquisitionAutoInterval)
}
