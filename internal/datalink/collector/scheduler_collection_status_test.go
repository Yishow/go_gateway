package collector

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

func TestCollectionStatusObservesOnlyCurrentPollableSlots(t *testing.T) {
	var absent *Scheduler
	require.Empty(t, absent.CollectionStatus())
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	d := &schema.Device{ID: "status-device", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"127.0.0.1","port":1}`}
	s.AddDevice(d)
	group := &schema.PollingGroup{ID: "status-group", Enabled: true, IntervalMs: int(time.Hour / time.Millisecond)}
	p := &schema.Point{ID: "status-point", DeviceID: d.ID, Address: "40001", DataType: schema.DataTypeUint32, Enabled: true, PollingGroupID: &group.ID, RevisionID: "revision", FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD", FormatVersion: 1}
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	require.NotZero(t, s.AddPointWithObservationLease(p, profile))
	require.False(t, s.CollectionStatus()[p.ID].Running)
	require.True(t, s.CollectionStatus()[p.ID].Readable)
	require.NoError(t, s.Start([]*schema.PollingGroup{group}))
	t.Cleanup(func() { require.NoError(t, s.StopContext(context.Background())) })
	snapshot := s.CollectionStatus()
	require.Equal(t, CollectionPointStatus{DeviceID: d.ID, Running: true, Readable: true}, snapshot[p.ID])
	delete(snapshot, p.ID)
	require.Len(t, s.CollectionStatus(), 1, "snapshot is detached from the scheduler")
	proof, ok := s.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	lease, ok := s.StagePointFormatRecovery(p, proof, profile)
	require.True(t, ok)
	require.Equal(t, CollectionPointStatus{DeviceID: d.ID}, s.CollectionStatus()[p.ID], "staged slot remains visible, never running/readable")
	require.True(t, s.ActivatePointFormatRecovery(p.ID, lease, profile))
	require.True(t, s.CollectionStatus()[p.ID].Running)
	s.RemovePollingGroup(group.ID)
	require.False(t, s.CollectionStatus()[p.ID].Running)
	disabled := *group
	disabled.Enabled = false
	s.AddPollingGroup(&disabled)
	require.False(t, s.CollectionStatus()[p.ID].Running)
	s.AddPollingGroup(group)
	require.True(t, s.CollectionStatus()[p.ID].Running)
	require.NoError(t, s.Stop())
	require.False(t, s.CollectionStatus()[p.ID].Running)
	require.True(t, s.CollectionStatus()[p.ID].Readable, "stop does not invalidate retained health facts")
	s.RemoveDevice(d.ID)
	require.Empty(t, s.CollectionStatus())
}
