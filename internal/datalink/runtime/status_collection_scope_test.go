package runtime

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestCollectionStatusSeparatesLifecycleFromHealth(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		mutate                 func(*testing.T, *Service, []*schema.Point)
		status                 string
		running                bool
		healthy, stale, errors int
	}{
		{"new collector awaiting first read", func(_ *testing.T, _ *Service, p []*schema.Point) {
			for _, pt := range p {
				pt.LastReadAt = nil
			}
		}, "warning", true, 0, 2, 0},
		{"fresh", nil, "running", true, 2, 0, 0},
		{"stale remains running", func(_ *testing.T, _ *Service, p []*schema.Point) {
			at := time.Now().Add(-2 * time.Hour)
			p[0].LastReadAt = &at
		}, "warning", true, 1, 1, 0},
		{"partial point error remains visible", func(_ *testing.T, _ *Service, p []*schema.Point) { p[0].LastError = "read failed" }, "error", true, 1, 0, 1},
		{"lost selected slot cannot become good", func(t *testing.T, s *Service, p []*schema.Point) { s.scheduler.RemovePoint(p[0].ID) }, "warning", true, 1, 1, 0},
		{"missing selected record cannot become good", func(_ *testing.T, s *Service, p []*schema.Point) { s.snapshot.Points = p[1:] }, "warning", true, 1, 1, 0},
		{"missing all selected records does not stop scheduling", func(_ *testing.T, s *Service, _ []*schema.Point) { s.snapshot.Points = nil }, "warning", true, 0, 2, 0},
		{"quarantined selected point is not healthy", func(t *testing.T, s *Service, p []*schema.Point) { stageStatusPoint(t, s, p[0]) }, "warning", true, 1, 1, 0},
		{"all quarantined points cannot run", func(t *testing.T, s *Service, p []*schema.Point) {
			for _, pt := range p {
				stageStatusPoint(t, s, pt)
			}
		}, "warning", false, 0, 2, 0},
		{"stopped scheduler ignores historical good", func(t *testing.T, s *Service, _ []*schema.Point) { require.NoError(t, s.scheduler.Stop()) }, "idle", false, 2, 0, 0},
		{"removed ticker ignores historical good", func(t *testing.T, s *Service, _ []*schema.Point) { s.scheduler.RemovePollingGroup("status-group") }, "idle", false, 2, 0, 0},
		{"removed device ignores historical good", func(t *testing.T, s *Service, _ []*schema.Point) { s.RemoveDevice("status-device") }, "idle", false, 0, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc, points := newStatusScopeService(t)
			if tc.mutate != nil {
				tc.mutate(t, svc, points)
			}
			snapshot, err := svc.RuntimeStatusSnapshot(t.Context(), "status-device")
			require.NoError(t, err)
			require.Len(t, snapshot.Collectors, 1)
			c := snapshot.Collectors[0]
			require.Equal(t, tc.running, c.Running)
			require.Equal(t, tc.status, c.Status)
			require.Equal(t, tc.healthy, c.PointsHealthy)
			require.Equal(t, tc.stale, c.PointsStale)
			require.Equal(t, tc.errors, c.PointsError)
			stream, cancel := svc.SubscribeStatusEvents("status-device")
			defer cancel()
			select {
			case event := <-stream:
				require.Equal(t, c, event, "SSE and HTTP snapshot share scope and lifecycle semantics")
			case <-time.After(time.Second):
				t.Fatal("missing initial status event")
			}
		})
	}
}

func newStatusScopeService(t *testing.T) (*Service, []*schema.Point) {
	t.Helper()
	now := time.Now()
	groupID := "status-group"
	d := &schema.Device{ID: "status-device", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":1}`}
	points := []*schema.Point{
		{ID: "status-a", DeviceID: d.ID, Address: "40001", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: &groupID, LastReadAt: &now},
		{ID: "status-b", DeviceID: d.ID, Address: "40002", DataType: schema.DataTypeInt16, Enabled: true, PollingGroupID: &groupID, LastReadAt: &now},
	}
	groups := []*schema.PollingGroup{{ID: groupID, Enabled: true, IntervalMs: int(time.Hour / time.Millisecond)}}
	svc, err := NewService(Config{Writer: &mockWriter{}, Snapshot: Snapshot{Devices: []*schema.Device{d}, Points: points, PollingGroups: groups}})
	require.NoError(t, err)
	require.NoError(t, svc.Start(t.Context()))
	t.Cleanup(func() { require.NoError(t, svc.Stop(context.Background())) })
	projection := &workspace.RuntimeProjection{WorkspaceID: "status-ws", Version: "status-v1", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: points, PollingGroups: groups}
	svc.SetWorkspaceProjectionReader(workspaceProjectionStub{projection: projection})
	require.NoError(t, svc.ApplyWorkspaceProjection(t.Context(), projection))
	return svc, points
}

func stageStatusPoint(t *testing.T, s *Service, p *schema.Point) {
	t.Helper()
	p.DataType = schema.DataTypeUint32
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.FormatVersion = 1
	p.RevisionID = "quarantined-revision"
	d := s.snapshot.Devices[0]
	profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(d.ConnectionConfig)})
	require.NoError(t, err)
	proof, ok := s.scheduler.CaptureOrdinaryPoint(p.ID)
	require.True(t, ok)
	_, ok = s.scheduler.StagePointFormatRecovery(p, proof, profile)
	require.True(t, ok)
}
