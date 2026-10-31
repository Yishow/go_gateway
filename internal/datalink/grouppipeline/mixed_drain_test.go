package grouppipeline

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedLocalDrainNoPendingCreatesNoWindow(t *testing.T) {
	f, p, snap := mixedPipelineFixture(t)
	m := p.boundaries[boundaryKey(snap.Group.ID, snap.AppliedRevision)]
	require.NoError(t, p.beginMixedDrain(t.Context(), m))
	require.True(t, m.drainFinished.Load())
	f.clock = at(600)
	p.TickAll(t.Context())
	r, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Nil(t, r.Mixed)
	require.Empty(t, p.Status())
	require.True(t, m.boundary.ReplayState().Until.IsZero())
}

func TestMixedLocalDrainCapsLateTickAndCapturedIntake(t *testing.T) {
	f, p, snap := mixedPipelineFixture(t)
	f.clock = at(10)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(snap.Group, "accepted", at(10), 2.0)))
	m := p.boundaries[boundaryKey(snap.Group.ID, snap.AppliedRevision)]
	captured := []*managed{m}
	require.NoError(t, p.beginMixedDrain(t.Context(), m))
	require.Equal(t, StateRetiring, p.Status()[0].State)
	require.NoError(t, p.acceptWithOwnership(t.Context(), envelope(snap.Group, "after-disable", at(11), 7.0), captured))
	r, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Len(t, r.Samples, 1)
	f.clock = at(600)
	p.TickAll(t.Context())
	r, err = f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Equal(t, at(15), r.Mixed.Progress.ClosedThrough)
	require.True(t, r.Mixed.EndsAt.IsZero())
	require.Empty(t, p.Status())
	// A worker captured before pruning must not emit additional empty windows.
	p.reportTick(t.Context(), m, func() error { return p.tickManaged(t.Context(), m, at(900)) })
	r, err = f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Equal(t, at(15), r.Mixed.Progress.ClosedThrough)
}

func TestMixedLocalDrainRespectsLatenessAndNominalSlowWindow(t *testing.T) {
	f := newFixture(t)
	g, layout := mixedStorageProofFixture(t, false, "sqlite")
	g.Members[0].Cadence.Interval = 60
	g.RowPolicy.AllowedLatenessSeconds = 5
	snap := snapshotOf(g, base)
	snap.RuntimeLayout = layout
	raw, err := json.Marshal(snap)
	require.NoError(t, err)
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupKey(snap), raw))
	f.groups.groups = []*workspace.WriteGroup{liveGroup(g, workspace.WriteGroupStatusReady)}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{g.ID: {snap}}
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	f.clock = at(10)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(g, "slow", at(10), 2.0)))
	m := p.boundaries[boundaryKey(g.ID, g.AppliedRevision)]
	require.NoError(t, p.beginMixedDrain(t.Context(), m))
	require.Equal(t, at(65), m.drain.clockLimit)
	f.clock = at(60)
	p.TickAll(t.Context())
	require.False(t, m.drainFinished.Load())
	f.clock = at(65)
	p.TickAll(t.Context())
	r, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Equal(t, at(60), r.Mixed.Progress.ClosedThrough)
	require.True(t, r.Mixed.EndsAt.IsZero())
	require.Empty(t, p.Status())
}

func TestMixedLocalDrainPreparationFailureKeepsIntakeFenced(t *testing.T) {
	_, p, snap := mixedPipelineFixture(t)
	m := p.boundaries[boundaryKey(snap.Group.ID, snap.AppliedRevision)]
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	require.Error(t, p.beginMixedDrain(ctx, m))
	require.True(t, m.localDrain.Load())
	require.False(t, m.drainFinished.Load())
	require.Nil(t, m.drain)
	require.NoError(t, p.beginMixedDrain(t.Context(), m))
	require.True(t, m.drainFinished.Load())
}

func TestMixedLocalDrainKeepsSlowPendingAcrossFastClosure(t *testing.T) {
	f := newFixture(t)
	g, layout := mixedStorageProofFixture(t, false, dbtarget.SQLDialectSQLite)
	slow := g.Members[0]
	slow.PointID = "slow-point"
	slow.TagID = "slow-tag"
	slow.TargetColumn = "slow_value"
	slow.Cadence = &workspace.WriteGroupMemberCadence{Interval: 60, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
	g.Members = append(g.Members, slow)
	layout.TagTypes[slow.TagID] = schema.DataTypeFloat64
	expected, err := workspace.MixedManagedSchemaLayout(g, layout.TagTypes, string(layout.Dialect))
	require.NoError(t, err)
	layout.Columns = nil
	for _, c := range expected.Columns {
		layout.Columns = append(layout.Columns, dbtarget.ColumnInfo{Name: c.Name, DataType: c.SQLType, Nullable: c.Nullable, PrimaryKey: c.PrimaryKey})
	}
	snap := snapshotOf(g, base)
	snap.RuntimeLayout = layout
	raw, err := json.Marshal(snap)
	require.NoError(t, err)
	require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupKey(snap), raw))
	f.groups.groups = []*workspace.WriteGroup{liveGroup(g, workspace.WriteGroupStatusReady)}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{g.ID: {snap}}
	p := f.pipeline()
	require.NoError(t, p.Reconcile(t.Context()))
	f.clock = at(10)
	require.NoError(t, p.AcceptSample(t.Context(), envelope(g, "fast", at(10), 2.0)))
	e := envelope(g, "slow", at(10), 4.0)
	e.PointID = slow.PointID
	e.TagID = slow.TagID
	require.NoError(t, p.AcceptSample(t.Context(), e))
	m := p.boundaries[boundaryKey(g.ID, g.AppliedRevision)]
	require.NoError(t, p.beginMixedDrain(t.Context(), m))
	f.clock = at(15)
	p.TickAll(t.Context())
	require.False(t, m.drainFinished.Load())
	r, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Len(t, r.Samples, 1)
	require.Equal(t, "slow", r.Samples[0].SampleID)
	f.clock = at(600)
	p.TickAll(t.Context())
	r, err = f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Equal(t, at(60), r.Mixed.Progress.ClosedThrough)
	require.True(t, r.Mixed.EndsAt.IsZero())
	require.Empty(t, r.Samples)
	require.Empty(t, p.Status())
}
