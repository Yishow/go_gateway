package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	runtimepkg "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

type telemetryGateWindowRepo struct {
	point.Repository
	hook func()
}

func (r *telemetryGateWindowRepo) UpdateReadResult(ctx context.Context, id string, value any, message string) error {
	r.hook()
	return r.Repository.UpdateReadResult(ctx, id, value, message)
}
func TestTelemetryCASWriteWindowKeepsNewRevisionState(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &lateManualProfileProtocol{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1"}`, Status: schema.DeviceStatusActive}
	scheduler.AddDevice(d)
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()
	_, err := db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('d','D','modbus_tcp','active','{}','{}')`)
	require.NoError(t, err)
	wrapped := &telemetryGateWindowRepo{Repository: point.NewSQLRepository(db)}
	svc := point.NewService(wrapped, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() {}
	wrapped.hook = func() {
		_, err := db.ExecContext(t.Context(), `UPDATE points SET revision_id='new-durable-revision',format_mode='override',format_override='DCBA' WHERE id=?`, p.ID)
		require.NoError(t, err)
	}

	h := NewPointHandler(svc, nil).WithPolling(scheduler, nil, nil)
	got, err := h.pollPoint(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}

func (r *telemetryGateWindowRepo) UpdateReadResultAtRevision(ctx context.Context, id, deviceID, revision string, value any, message string) error {
	r.hook()
	return r.Repository.(point.ReadResultRevisionWriter).UpdateReadResultAtRevision(ctx, id, deviceID, revision, value, message)
}

func TestTelemetryCASBatchWriteWindowKeepsNewRevisionState(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &lateManualProfileProtocol{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1"}`, Status: schema.DeviceStatusActive}
	scheduler.AddDevice(d)
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()
	_, err := db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('d','D','modbus_tcp','active','{}','{}')`)
	require.NoError(t, err)
	wrapped := &telemetryGateWindowRepo{Repository: point.NewSQLRepository(db)}
	svc := point.NewService(wrapped, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() {}
	wrapped.hook = func() {
		_, err := db.ExecContext(t.Context(), `UPDATE points SET revision_id='new-durable-revision',format_mode='override',format_override='DCBA' WHERE id=?`, p.ID)
		require.NoError(t, err)
	}

	h := NewPointHandler(svc, nil).WithPolling(scheduler, nil, nil)
	results, err := h.pollBatch(t.Context(), []string{p.ID})
	require.Len(t, results, 1)
	got := results[0]
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}

func TestTelemetryCASDirectWriteWindowKeepsNewRevisionState(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &lateManualProfileProtocol{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1"}`, Status: schema.DeviceStatusActive}
	scheduler.AddDevice(d)
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()
	_, err := db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('d','D','modbus_tcp','active','{}','{}')`)
	require.NoError(t, err)
	wrapped := &telemetryGateWindowRepo{Repository: point.NewSQLRepository(db)}
	svc := point.NewService(wrapped, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() {}
	wrapped.hook = func() {
		_, err := db.ExecContext(t.Context(), `UPDATE points SET revision_id='new-durable-revision',format_mode='override',format_override='DCBA' WHERE id=?`, p.ID)
		require.NoError(t, err)
	}

	dr := device.NewMemoryRepository()
	require.NoError(t, dr.Create(t.Context(), d))
	live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
	require.NoError(t, err)
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		current, err := svc.GetByID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		return &workspace.RuntimeProjection{WorkspaceID: "W", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{current}}, nil
	}))
	h := NewPointHandler(svc, nil).WithDirectReader(NewPointDirectReader(svc, device.NewService(dr, manager)).WithReadProfileResolver(live))
	got, err := h.pollPoint(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}

func TestTelemetryCASDirectBatchWriteWindowKeepsNewRevisionState(t *testing.T) {
	old, exists := connector.GetFactory(schema.ProtocolModbusTCP)
	probe := &lateManualProfileProtocol{}
	connector.Register(schema.ProtocolModbusTCP, func() connector.Protocol { return probe })
	t.Cleanup(func() {
		if exists {
			connector.Register(schema.ProtocolModbusTCP, old)
		} else {
			connector.Unregister(schema.ProtocolModbusTCP)
		}
	})
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "d", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1"}`, Status: schema.DeviceStatusActive}
	scheduler.AddDevice(d)
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()
	_, err := db.ExecContext(t.Context(), `INSERT INTO devices(id,name,protocol,status,connection_config,readiness_status) VALUES('d','D','modbus_tcp','active','{}','{}')`)
	require.NoError(t, err)
	wrapped := &telemetryGateWindowRepo{Repository: point.NewSQLRepository(db)}
	svc := point.NewService(wrapped, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() {}
	wrapped.hook = func() {
		_, err := db.ExecContext(t.Context(), `UPDATE points SET revision_id='new-durable-revision',format_mode='override',format_override='DCBA' WHERE id=?`, p.ID)
		require.NoError(t, err)
	}

	dr := device.NewMemoryRepository()
	require.NoError(t, dr.Create(t.Context(), d))
	live, err := runtimepkg.NewService(runtimepkg.Config{Writer: storage.NewSQLiteWriter(db)})
	require.NoError(t, err)
	live.SetWorkspaceProjectionReader(pointProfileProjectionFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
		current, err := svc.GetByID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		return &workspace.RuntimeProjection{WorkspaceID: "W", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{current}}, nil
	}))
	h := NewPointHandler(svc, nil).WithDirectReader(NewPointDirectReader(svc, device.NewService(dr, manager)).WithReadProfileResolver(live))
	results, err := h.pollBatch(t.Context(), []string{p.ID})
	require.Len(t, results, 1)
	got := results[0]
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}
