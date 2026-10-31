package grouppipeline

import (
	"encoding/json"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedPipelineFrozenBuildDoesNotNeedCommonInterval(t *testing.T) {
	for _, long := range []bool{false, true} {
		t.Run(map[bool]string{false: "wide", true: "long"}[long], func(t *testing.T) {
			g, layout := mixedStorageProofFixture(t, long, dbtarget.SQLDialectSQLite)
			snap := &workspace.WriteGroupAppliedSnapshot{WorkspaceID: g.WorkspaceID, GroupID: g.ID, AppliedRevision: g.AppliedRevision, EffectiveAt: base, Group: g, RuntimeLayout: layout}
			store := newFixture(t).store
			raw, err := json.Marshal(snap)
			require.NoError(t, err)
			require.NoError(t, store.SaveRuntimeVersion(t.Context(), groupKey(snap), raw))
			pipe := New(Dependencies{Store: store}, Config{Now: func() time.Time { return at(15) }})
			loaded, err := pipe.loadRuntimeSnapshot(t.Context(), groupKey(snap))
			require.NoError(t, err)
			built, reason, cause := pipe.buildFrozenWithCause(t.Context(), loaded, time.Time{}, time.Time{})
			require.NoError(t, cause)
			require.Empty(t, reason)
			require.NotNil(t, built)
			require.True(t, built.boundary.NextClose().IsZero(), "mixed must not synthesize legacy cursor")
			require.NoError(t, built.boundary.Tick(t.Context(), at(15)))
			restored, err := store.Restore(t.Context(), groupKey(snap))
			require.NoError(t, err)
			require.NotNil(t, restored.Mixed)
			require.Equal(t, at(15), restored.Mixed.Progress.ClosedThrough)
			require.True(t, restored.NextClose.IsZero())
		})
	}
}

func mixedPipelineFixture(t *testing.T) (*fixture, *Pipeline, *workspace.WriteGroupAppliedSnapshot) {
	t.Helper()
	f := newFixture(t)
	g, layout := mixedStorageProofFixture(t, false, dbtarget.SQLDialectSQLite)
	snap := snapshotOf(g, base)
	snap.RuntimeLayout = layout
	raw, err := json.Marshal(snap)
	require.NoError(t, err)
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupKey(snap), raw))
	f.groups.groups = []*workspace.WriteGroup{liveGroup(g, workspace.WriteGroupStatusReady)}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{g.ID: {snap}}
	pipe := f.pipeline()
	require.NoError(t, pipe.Reconcile(t.Context()))
	return f, pipe, snap
}
func TestMixedPipelineDisableDrainsWithoutPermanentWindowEnd(t *testing.T) {
	f, pipe, snap := mixedPipelineFixture(t)
	f.clock = at(10)
	require.NoError(t, pipe.AcceptSample(t.Context(), envelope(snap.Group, "accepted", at(10), 2.0)))
	f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusDisabled)
	require.NoError(t, pipe.Reconcile(t.Context()))
	require.False(t, pipe.WantsSample("device-1", "point-1", "tag-t"), "disable stops intake instead of failing the absent common interval")
	f.clock = at(15)
	pipe.TickAll(t.Context())
	restored, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.NotNil(t, restored.Mixed)
	require.True(t, restored.Mixed.EndsAt.IsZero(), "local drain is not a persistent actual cutoff")
	require.Equal(t, at(15), restored.Mixed.Progress.ClosedThrough)
	require.Empty(t, pipe.Status(), "finished local drain must not create endless empty windows")
}
