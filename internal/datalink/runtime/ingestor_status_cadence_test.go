package runtime

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type statusProjectionCounter struct {
	projection *workspace.RuntimeProjection
	reads      atomic.Int64
}

func (c *statusProjectionCounter) RuntimeProjection(context.Context) (*workspace.RuntimeProjection, error) {
	c.reads.Add(1)
	return c.projection, nil
}

// No status worker runs unless a test starts the existing worker explicitly.
// This keeps the per-ingest query budget deterministic, independent of tick timing.
func newStatusCadenceFixture(t *testing.T) (*Service, *statusProjectionCounter, []*schema.Point) {
	t.Helper()
	now := time.Now()
	groupID := "cadence-group"
	d := &schema.Device{ID: "cadence-device", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":1}`}
	points := []*schema.Point{
		{ID: "a", DeviceID: d.ID, Address: "40001", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: &groupID, LastReadAt: &now},
		{ID: "b", DeviceID: d.ID, Address: "40002", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: &groupID, LastReadAt: &now},
	}
	groups := []*schema.PollingGroup{{ID: groupID, Enabled: true, IntervalMs: int(time.Hour / time.Millisecond)}}
	cfg := collector.DefaultSchedulerConfig()
	cfg.DefaultRetryCount = 0
	cfg.BreakerConfig = health.BreakerConfig{ErrorThreshold: 0.2, UnstableThreshold: 0.5, WindowSize: 3, CooldownPeriod: time.Hour}
	scheduler := collector.NewScheduler(cfg, nil)
	scheduler.AddDevice(d)
	for _, p := range points {
		scheduler.AddPoint(p)
	}
	require.NoError(t, scheduler.Start(groups))
	projection := &workspace.RuntimeProjection{WorkspaceID: "cadence-workspace", Version: "v1", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: points, PollingGroups: groups}
	counter := &statusProjectionCounter{projection: projection}
	svc := &Service{writer: &mockWriter{}, scheduler: scheduler, workspace: counter, snapshot: Snapshot{Devices: projection.Devices, Points: points, PollingGroups: groups}, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}, stopCh: make(chan struct{})}
	for _, p := range points {
		svc.pointMetaIndex[p.ID] = pointMeta{DeviceID: p.DeviceID, Address: p.Address}
	}
	svc.markDeviceProjectionAligned(d.ID, projection)
	svc.running.Store(true)
	t.Cleanup(func() { require.NoError(t, svc.Stop(context.Background())) })
	return svc, counter, points
}

func TestIngestorDoesNotSynchronouslyRebuildStatus(t *testing.T) {
	for _, subscribed := range []bool{false, true} {
		t.Run(map[bool]string{false: "no-status-subscriber", true: "live-status-subscriber"}[subscribed], func(t *testing.T) {
			svc, counter, points := newStatusCadenceFixture(t)
			if subscribed {
				stream, cancel := svc.SubscribeStatusEvents(points[0].DeviceID)
				defer cancel()
				require.Equal(t, "running", receiveCadenceStatus(t, stream).Status)
			}
			values, cancel := svc.SubscribeValueEvents(points[0].DeviceID, nil)
			defer cancel()
			counter.reads.Store(0)
			for i := 0; i < 4; i++ {
				p := points[i%len(points)]
				quality := schema.QualityGood
				message := ""
				if i == 1 {
					quality = schema.QualityBad
					message = "read failed"
				}
				svc.handleCollectedValue(t.Context(), collector.CollectedValue{DeviceID: p.DeviceID, PointID: p.ID, Value: i, Quality: quality, Error: message, Timestamp: time.Now(), InstallationLease: svc.scheduler.CapturePointLease(p.ID, p.RevisionID), InstallationScope: svc.scheduler.InstallationScope()})
				select {
				case value := <-values:
					require.Equal(t, p.ID, value.PointID)
					require.Equal(t, quality, value.Quality)
				case <-time.After(time.Second):
					t.Fatal("point value/quality event was delayed by status cadence")
				}
			}
			require.Zero(t, counter.reads.Load(), "ingestion must not rebuild workspace status for each point, even with a live UI subscriber")
		})
	}
}
