package grouppipeline

import (
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestReconcileAfterEffectiveBoundaryRetiresPreviousRevisionWithoutDoubleIntake(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	first := f.groups.snapshots["group-G"][0].Group
	require.NoError(t, p.Reconcile(t.Context()))
	second := groupAt("rev-2")
	// A row-policy edit leaves the actual source/mapping identity unchanged.
	for i := range second.Members {
		second.Members[i].SourceRevision = first.Members[i].SourceRevision
	}
	f.groups.groups = []*workspace.WriteGroup{liveGroup(second, workspace.WriteGroupStatusReady)}
	f.groups.snapshots["group-G"] = append(f.groups.snapshots["group-G"], snapshotOf(second, at(20)))
	// No reconcile occurred before t=20: the first one observes the new revision.
	f.clock = at(22)
	require.NoError(t, p.Reconcile(t.Context()))
	require.NoError(t, p.AcceptSample(t.Context(), envelope(second, "new-only", at(21), 2.0)))
	var old, current int
	require.NoError(t, rowCountWhere(f.store, &old, "wg_delivery_samples", "group_revision = 'rev-1'"))
	require.NoError(t, rowCountWhere(f.store, &current, "wg_delivery_samples", "group_revision = 'rev-2'"))
	require.Zero(t, old, "the old boundary must not ACK samples after its replacement became effective")
	require.Equal(t, 1, current)
	f.clock = at(40)
	p.TickAll(t.Context())
	require.Len(t, p.Status(), 1)
	require.Equal(t, "rev-2", p.Status()[0].Revision)
}

func TestRevisionHandoverKeepsJournaledRowsAcrossRestart(t *testing.T) {
	f := newFixture(t)
	p := f.pipeline()
	first := f.groups.snapshots["group-G"][0].Group
	require.NoError(t, p.Reconcile(t.Context()))

	f.clock = at(12)
	oldTemperature := envelope(first, "old-temperature", at(11), 21.5)
	oldPressure := envelope(first, "old-pressure", at(11), int64(101))
	oldPressure.PointID = "point-2"
	oldPressure.TagID = "tag-p"
	oldPressure.SourceRevision = first.Members[1].SourceRevision
	require.NoError(t, p.AcceptSample(t.Context(), oldTemperature))
	require.NoError(t, p.AcceptSample(t.Context(), oldPressure))

	second := groupAt("rev-2")
	f.groups.groups = []*workspace.WriteGroup{liveGroup(second, workspace.WriteGroupStatusReady)}
	f.groups.snapshots["group-G"] = append(f.groups.snapshots["group-G"], snapshotOf(second, at(20)))
	f.clock = at(21)
	require.NoError(t, p.Reconcile(t.Context()))
	newTemperature := envelope(second, "new-temperature", at(21), 22.5)
	newPressure := envelope(second, "new-pressure", at(21), int64(102))
	newPressure.PointID = "point-2"
	newPressure.TagID = "tag-p"
	newPressure.SourceRevision = second.Members[1].SourceRevision
	require.NoError(t, p.AcceptSample(t.Context(), newTemperature))
	require.NoError(t, p.AcceptSample(t.Context(), newPressure))

	restarted := f.pipeline()
	require.NoError(t, restarted.Reconcile(t.Context()), "restart must rebuild the active revision and drain the historical journal")
	f.clock = at(32)
	restarted.TickAll(t.Context())

	var oldRows, newRows int
	require.NoError(t, rowCountWhere(f.store, &oldRows, "wg_delivery_outbox", "group_revision = 'rev-1'"))
	require.NoError(t, rowCountWhere(f.store, &newRows, "wg_delivery_outbox", "group_revision = 'rev-2'"))
	require.Equal(t, 1, oldRows, "the superseded revision remains queued under its original identity")
	require.Equal(t, 1, newRows, "the active revision closes its own bucket after restart")
	var consumed int
	require.NoError(t, rowCountWhere(f.store, &consumed, "wg_delivery_samples", "sample_id IN ('old-temperature', 'old-pressure', 'new-temperature', 'new-pressure') AND consumed = 1"))
	require.Equal(t, 4, consumed, "restart closure consumes each accepted sample exactly once")
}
