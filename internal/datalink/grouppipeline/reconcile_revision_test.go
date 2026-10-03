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
