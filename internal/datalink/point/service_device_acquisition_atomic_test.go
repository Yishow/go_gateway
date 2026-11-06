package point_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func acquisitionSQLiteService(t *testing.T) (*sql.DB, *point.Service) {
	t.Helper()
	dsn := "file:" + filepath.Join(t.TempDir(), "acquisition.db") + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_, err = db.ExecContext(t.Context(), `PRAGMA journal_mode=WAL`)
	require.NoError(t, err)
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err = db.ExecContext(t.Context(), `INSERT INTO devices (id, name, protocol, description) VALUES ('device-a', 'A', 'modbus_tcp', ''), ('device-b', 'B', 'modbus_tcp', '')`)
	require.NoError(t, err)
	service := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	return db, service
}

func TestDeviceIntervalFailureRollsBackGroupAndAllPointAssignments(t *testing.T) {
	for _, existing := range []bool{false, true} {
		name := "create group"
		if existing {
			name = "update group"
		}
		t.Run(name, func(t *testing.T) {
			db, service := acquisitionSQLiteService(t)
			ctx := t.Context()
			shared, err := service.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "shared", IntervalMs: 1000})
			require.NoError(t, err)
			if existing {
				_, err = service.SetDeviceAcquisitionInterval(ctx, "device-a", 1000)
				require.NoError(t, err)
			}
			beforeGroup, err := service.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
			require.NoError(t, err)
			before := make([]*schema.Point, 0, 3)
			for _, item := range []struct{ device, name, address string }{
				{"device-a", "A1", "40001"},
				{"device-a", "A2", "40002"},
				{"device-b", "B1", "40001"},
			} {
				created, createErr := service.Create(ctx, point.CreatePointRequest{DeviceID: item.device, Name: item.name, Address: item.address, DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
				require.NoError(t, createErr)
				before = append(before, created)
			}
			_, err = db.ExecContext(ctx, `CREATE TRIGGER fail_second_assignment BEFORE UPDATE OF polling_group_id ON points
				WHEN OLD.device_id = 'device-a' AND OLD.name = 'A2'
				BEGIN SELECT RAISE(ABORT, 'injected assignment failure'); END`)
			require.NoError(t, err)

			devices := device.NewService(device.NewSQLRepository(db), nil)
			devices.SetAcquisitionIntervalCoordinator(service)
			interval := 500
			_, err = devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval})
			require.ErrorContains(t, err, "injected assignment failure")

			group, err := service.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
			require.NoError(t, err)
			if beforeGroup == nil {
				assert.Nil(t, group, "a failed save must not create an owned group")
			} else if assert.NotNil(t, group) {
				assert.Equal(t, beforeGroup.ID, group.ID)
				assert.Equal(t, beforeGroup.IntervalMs, group.IntervalMs, "a failed save must not change the owned interval")
				assert.Equal(t, beforeGroup.Enabled, group.Enabled)
			}
			for _, original := range before {
				got, getErr := service.GetByID(ctx, original.ID)
				require.NoError(t, getErr)
				if assert.NotNil(t, got.PollingGroupID) {
					assert.Equal(t, *original.PollingGroupID, *got.PollingGroupID, "failed save must leave %s assigned to its original group", original.Name)
				}
				assert.Equal(t, original.RevisionID, got.RevisionID)
			}
		})
	}
}

func TestDeviceIntervalSavesThroughSQLiteDeviceService(t *testing.T) {
	db, service := acquisitionSQLiteService(t)
	ctx := t.Context()
	shared, err := service.CreatePollingGroup(ctx, point.CreatePollingGroupRequest{Name: "device:device-a", IntervalMs: 1000})
	require.NoError(t, err)
	a1, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "A1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)
	a2, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-a", Name: "A2", Address: "40002", DataType: schema.DataTypeUint16})
	require.NoError(t, err)
	disabled := false
	a2, err = service.Update(ctx, a2.ID, point.UpdatePointRequest{Enabled: &disabled})
	require.NoError(t, err)
	b1, err := service.Create(ctx, point.CreatePointRequest{DeviceID: "device-b", Name: "B1", Address: "40001", DataType: schema.DataTypeUint16, PollingGroupID: &shared.ID})
	require.NoError(t, err)

	devices := device.NewService(device.NewSQLRepository(db), nil)
	devices.SetAcquisitionIntervalCoordinator(service)
	interval := 500
	updated, err := devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval})
	require.NoError(t, err)
	require.Equal(t, &interval, updated.AcquisitionIntervalMs)
	group, err := service.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
	require.NoError(t, err)
	require.Equal(t, "device:device-a#2", group.Name)
	for _, original := range []*schema.Point{a1, a2} {
		got, getErr := service.GetByID(ctx, original.ID)
		require.NoError(t, getErr)
		assert.Equal(t, &group.ID, got.PollingGroupID)
		assert.Equal(t, original.Name, got.Name)
		assert.Equal(t, original.Enabled, got.Enabled)
		assert.NotEqual(t, original.RevisionID, got.RevisionID)
	}
	gotB, err := service.GetByID(ctx, b1.ID)
	require.NoError(t, err)
	require.Equal(t, b1.PollingGroupID, gotB.PollingGroupID)
	require.Equal(t, b1.RevisionID, gotB.RevisionID)

	interval = 750
	updated, err = devices.Update(ctx, "device-a", device.UpdateDeviceRequest{AcquisitionIntervalMs: &interval})
	require.NoError(t, err)
	require.Equal(t, &interval, updated.AcquisitionIntervalMs)
	again, err := service.GetPollingGroupByOwnerDeviceID(ctx, "device-a")
	require.NoError(t, err)
	require.Equal(t, group.ID, again.ID)
	require.Equal(t, 750, again.IntervalMs)
}
