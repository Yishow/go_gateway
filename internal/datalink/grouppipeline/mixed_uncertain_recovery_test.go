package grouppipeline

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/sqltx"
	"go-gateway/internal/datalink/workspace"
)

type mixedLostClosureAck struct {
	*groupdelivery.Ledger
	calls int
}

func (l *mixedLostClosureAck) CommitClosure(ctx context.Context, c groupdelivery.Closure) error {
	l.calls++
	if err := l.Ledger.CommitClosure(ctx, c); err != nil {
		return err
	}
	return sqltx.ErrCommitOutcomeUnknown
}

func TestMixedPipelineReconcileRestoresUnknownClosureWithoutReencoding(t *testing.T) {
	f, pipe, snap := mixedPipelineFixture(t)
	key := boundaryKey(snap.GroupID, snap.AppliedRevision)
	old := pipe.boundaries[key]
	ledger := &mixedLostClosureAck{Ledger: f.store.Ledger(groupKey(snap))}
	layout := snap.RuntimeLayout
	boundary, err := runtime.NewGroupBoundaryContext(t.Context(), runtime.GroupBoundaryConfig{
		Group: snap.Group, TagTypes: layout.TagTypes, Columns: layout.Columns, Dialect: layout.Dialect,
		CounterBindings: snap.CounterBindings, ReceiptTableReady: layout.ReceiptTableReady,
		StartsAt: snap.EffectiveAt, Ledger: ledger, Clock: func() time.Time { return f.clock },
	})
	require.NoError(t, err)
	old.boundary = boundary
	f.clock = at(10)
	require.NoError(t, pipe.AcceptSample(t.Context(), envelope(snap.Group, "accepted", at(10), 2.0)))
	f.clock = at(15)
	require.ErrorIs(t, old.boundary.Tick(t.Context(), at(15)), sqltx.ErrCommitOutcomeUnknown)
	require.True(t, old.boundary.NeedsDurableRestore())
	var original string
	require.NoError(t, f.store.DB().QueryRow(`SELECT payload FROM wg_delivery_outbox`).Scan(&original))
	require.NoError(t, pipe.Reconcile(t.Context()))
	fresh := pipe.boundaries[key]
	require.NotSame(t, old, fresh, "unknown acknowledgment requires a new restored boundary")
	require.True(t, old.superseded, "captured old worker cannot resume")
	require.False(t, fresh.boundary.NeedsDurableRestore())
	pipe.TickAll(t.Context())
	require.Equal(t, 1, ledger.calls, "unknown old closure must never be retried")
	var stored string
	require.NoError(t, f.store.DB().QueryRow(`SELECT payload FROM wg_delivery_outbox`).Scan(&stored))
	require.Equal(t, original, stored, "accepted SQL payload is not recomputed")
	var count int
	require.NoError(t, f.store.DB().QueryRow(`SELECT COUNT(*) FROM wg_delivery_outbox`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestMixedPipelineReadyReplacesLocalDrainFromDurableProgress(t *testing.T) {
	f, pipe, snap := mixedPipelineFixture(t)
	f.clock = at(10)
	require.NoError(t, pipe.AcceptSample(t.Context(), envelope(snap.Group, "accepted", at(10), 2.0)))
	f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusDisabled)
	require.NoError(t, pipe.Reconcile(t.Context()))
	old := pipe.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
	require.True(t, old.localDrain.Load())
	f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusReady)
	require.NoError(t, pipe.Reconcile(t.Context()))
	fresh := pipe.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
	require.NotSame(t, old, fresh)
	require.True(t, old.superseded)
	require.False(t, fresh.localDrain.Load())
	require.True(t, pipe.WantsSample("device-1", "point-1", "tag-t"))
}
