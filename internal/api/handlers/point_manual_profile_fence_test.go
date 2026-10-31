package handlers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type lateManualProfileProtocol struct {
	directProfileProtocol
	hook func()
}

func (r *lateManualProfileProtocol) Read(_ context.Context, _ connector.ReadRequest) (connector.ReadResult, error) {
	r.hook()
	return connector.ReadResult{Value: uint32(77), Quality: schema.QualityGood}, nil
}
func TestManualHTTPRejectsReinstalledRawResult(t *testing.T) {
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
	svc := point.NewService(point.NewMemoryRepository(), nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() { require.NotEqual(t, oldLease, scheduler.AddPointWithObservationLease(p, profile)) }
	h := NewPointHandler(svc, nil).WithPolling(scheduler, nil, nil)
	got, err := h.pollPoint(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}

func TestManualHTTPBatchRejectsReinstalledRawResult(t *testing.T) {
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
	svc := point.NewService(point.NewMemoryRepository(), nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() { require.NotEqual(t, oldLease, scheduler.AddPointWithObservationLease(p, profile)) }
	h := NewPointHandler(svc, nil).WithPolling(scheduler, nil, nil)
	results, err := h.pollBatch(t.Context(), []string{p.ID})
	require.NoError(t, err)
	require.Len(t, results, 1)
	got := results[0]
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}

func TestManualHTTPRejectsChangedDurableRevision(t *testing.T) {
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
	svc := point.NewService(point.NewMemoryRepository(), nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: d.ID, Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	oldLease := scheduler.AddPointWithObservationLease(p, profile)
	require.NotZero(t, oldLease)
	probe.hook = func() {
		p.RevisionID = "new-durable-revision"
		p.FormatMode = schema.PointFormatOverride
		p.FormatOverride = "DCBA"
		require.NoError(t, svc.Restore(t.Context(), p))
	}
	h := NewPointHandler(svc, nil).WithPolling(scheduler, nil, nil)
	got, err := h.pollPoint(t.Context(), p.ID)
	require.NoError(t, err)
	require.NotEqual(t, pollQualityGood, got.Quality, "late raw manual result must not be presented as current good data")
	current, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Nil(t, current.LastValue)
}
