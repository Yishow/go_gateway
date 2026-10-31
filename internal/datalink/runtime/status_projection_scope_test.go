package runtime

import (
	"context"
	"maps"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestCollectionScopePartialMergeFullReplaceAndStop(t *testing.T) {
	s, points := newStatusScopeService(t)
	p := *points[0]
	p.ID = "status-c"
	s.snapshot.Points = append(s.snapshot.Points, &p)
	projection := &workspace.RuntimeProjection{WorkspaceID: "status-ws", Version: "v2", DeviceIDs: []string{"status-device"}, Devices: s.snapshot.Devices, Points: s.snapshot.Points, PollingGroups: s.snapshot.PollingGroups}
	s.SetWorkspaceProjectionReader(workspaceProjectionStub{projection: projection})
	require.NoError(t, s.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), projection, projection.DeviceIDs, []string{p.ID}))
	state := s.deviceProjectionState("status-device")
	require.Equal(t, map[string]struct{}{"status-a": {}, "status-b": {}, "status-c": {}}, state.AppliedPointIDs)
	delete(state.AppliedPointIDs, "status-a")
	scopes := s.appliedCollectionScopes()
	delete(scopes["status-device"], "status-b")
	require.Len(t, s.deviceProjectionState("status-device").AppliedPointIDs, 3, "returned collections cannot mutate canonical projection")
	projection.Points = []*schema.Point{points[0]}
	projection.Version = "v3"
	require.NoError(t, s.ApplyWorkspaceProjection(t.Context(), projection))
	require.Equal(t, map[string]struct{}{"status-a": {}}, s.deviceProjectionState("status-device").AppliedPointIDs)
	require.NoError(t, s.Stop(context.Background()))
	require.Equal(t, map[string]struct{}{"status-a": {}}, s.deviceProjectionState("status-device").AppliedPointIDs, "stop preserves last applied scope")
	snapshot, err := s.RuntimeStatusSnapshot(t.Context(), "status-device")
	require.NoError(t, err)
	require.False(t, snapshot.Collectors[0].Running)
	require.Equal(t, 1, snapshot.Collectors[0].PointsTotal)
	s.RemoveDevice("status-device")
	require.Nil(t, s.deviceProjectionState("status-device").AppliedPointIDs)
}

func TestCollectionScopeFailedApplyDoesNotPublishRequestedPoints(t *testing.T) {
	for _, partial := range []bool{false, true} {
		t.Run(map[bool]string{false: "full", true: "partial"}[partial], func(t *testing.T) {
			s, points := newStatusScopeService(t)
			before := s.deviceProjectionState("status-device")
			next := *points[0]
			next.ID = "unapplied-point"
			projection := &workspace.RuntimeProjection{WorkspaceID: "status-ws", Version: "unapplied-v2", DeviceIDs: []string{"status-device"}, Devices: s.snapshot.Devices, Points: []*schema.Point{&next}, PollingGroups: s.snapshot.PollingGroups, Mappings: []*schema.Mapping{{ID: "missing-tag-mapping", PointID: next.ID, TagID: "absent-tag", Enabled: true}}}
			s.SetWorkspaceProjectionReader(workspaceProjectionStub{projection: projection})
			var err error
			if partial {
				err = s.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), projection, projection.DeviceIDs, []string{next.ID})
			} else {
				err = s.ApplyWorkspaceProjection(t.Context(), projection)
			}
			require.Error(t, err)
			after := s.deviceProjectionState("status-device")
			require.Equal(t, before.RuntimeVersion, after.RuntimeVersion)
			require.True(t, maps.Equal(before.AppliedPointIDs, after.AppliedPointIDs))
			require.NotContains(t, after.AppliedPointIDs, next.ID)
			snapshot, err := s.RuntimeStatusSnapshot(t.Context(), "status-device")
			require.NoError(t, err)
			require.NotEqual(t, "ready", snapshot.SnapshotState.State, "failed scope may not look confirmed healthy")
		})
	}
}

func TestCollectionScopeRestartBootstrapsFreshAppliedProjection(t *testing.T) {
	old, points := newStatusScopeService(t)
	require.NoError(t, old.Stop(context.Background()))
	projection := &workspace.RuntimeProjection{WorkspaceID: "status-ws", Version: "restart-v2", DeviceIDs: []string{"status-device"}, Devices: old.snapshot.Devices, Points: []*schema.Point{points[1]}, PollingGroups: old.snapshot.PollingGroups}
	restarted, err := NewService(Config{Writer: &mockWriter{}, Snapshot: old.snapshot})
	require.NoError(t, err)
	restarted.SetWorkspaceProjectionReader(workspaceProjectionStub{projection: projection})
	// This is the same bootstrap invoked by production service startup, not a
	// restoration of the previous process's in-memory status/counters.
	require.NoError(t, restarted.bootstrapFromWorkspaceProjection(t.Context(), projection))
	t.Cleanup(func() { require.NoError(t, restarted.Stop(context.Background())) })
	restarted.running.Store(true)
	require.Equal(t, map[string]struct{}{"status-b": {}}, restarted.deviceProjectionState("status-device").AppliedPointIDs)
	snapshot, err := restarted.RuntimeStatusSnapshot(t.Context(), "status-device")
	require.NoError(t, err)
	require.True(t, snapshot.Collectors[0].Running)
	require.Equal(t, 1, snapshot.Collectors[0].PointsTotal)
	require.Zero(t, snapshot.Metrics.CollectedTotal, "no historical metric is required to establish scheduling")
}
