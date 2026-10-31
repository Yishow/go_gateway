package grouppipeline

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/sqltx"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type mixedLostDrainAck struct {
	*groupdelivery.Ledger
	appendLost bool
	calls      int
}

func (l *mixedLostDrainAck) AppendSample(ctx context.Context, start time.Time, s snapshot.Sample) error {
	if err := l.Ledger.AppendSample(ctx, start, s); err != nil {
		return err
	}
	if l.appendLost {
		l.calls++
		return sqltx.ErrCommitOutcomeUnknown
	}
	return nil
}
func (l *mixedLostDrainAck) CommitClosure(ctx context.Context, c groupdelivery.Closure) error {
	if err := l.Ledger.CommitClosure(ctx, c); err != nil {
		return err
	}
	if !l.appendLost {
		l.calls++
		return sqltx.ErrCommitOutcomeUnknown
	}
	return nil
}
func mixedUncertainDrainFixture(t *testing.T, appendLost bool) (*fixture, *Pipeline, *workspace.WriteGroupAppliedSnapshot, *managed, *mixedLostDrainAck) {
	t.Helper()
	f, p, snap := mixedPipelineFixture(t)
	old := p.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
	ledger := &mixedLostDrainAck{Ledger: f.store.Ledger(groupKey(snap)), appendLost: appendLost}
	l := snap.RuntimeLayout
	b, err := runtime.NewGroupBoundaryContext(t.Context(), runtime.GroupBoundaryConfig{Group: snap.Group, TagTypes: l.TagTypes, Columns: l.Columns, Dialect: l.Dialect, CounterBindings: snap.CounterBindings, ReceiptTableReady: l.ReceiptTableReady, StartsAt: snap.EffectiveAt, Ledger: ledger, Clock: func() time.Time { return f.clock }})
	require.NoError(t, err)
	old.boundary = b
	f.clock = at(10)
	err = p.AcceptSample(t.Context(), envelope(snap.Group, "accepted", at(10), 2.0))
	if appendLost {
		require.ErrorIs(t, err, sqltx.ErrCommitOutcomeUnknown)
	} else {
		require.NoError(t, err)
		require.ErrorIs(t, b.Tick(t.Context(), at(15)), sqltx.ErrCommitOutcomeUnknown)
	}
	require.True(t, b.NeedsDurableRestore())
	return f, p, snap, old, ledger
}
func TestMixedUnknownDisabledOrMissingRestoresBeforeDrain(t *testing.T) {
	for _, missing := range []bool{false, true} {
		for _, appendLost := range []bool{false, true} {
			name := map[bool]string{false: "disabled", true: "missing"}[missing] + "/" + map[bool]string{false: "closure", true: "append"}[appendLost]
			t.Run(name, func(t *testing.T) {
				f, p, snap, old, ledger := mixedUncertainDrainFixture(t, appendLost)
				var original string
				if !appendLost {
					require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox`).Scan(&original))
				}
				if missing {
					f.groups.groups = nil
				} else {
					f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusDisabled)
				}
				require.NoError(t, p.Reconcile(t.Context()))
				fresh := p.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
				require.NotSame(t, old, fresh)
				require.True(t, old.superseded)
				require.True(t, fresh.localDrain.Load())
				require.False(t, fresh.boundary.NeedsDurableRestore())
				require.True(t, old.boundary.NeedsDurableRestore(), "never clear the old uncertainty")
				f.clock = at(15)
				for range 2 {
					p.TickAll(t.Context())
					require.NoError(t, p.Reconcile(t.Context()))
				}
				r, err := f.store.Restore(t.Context(), groupKey(snap))
				require.NoError(t, err)
				require.Empty(t, r.Samples)
				require.True(t, r.Mixed.EndsAt.IsZero())
				var count int
				require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_outbox`).Scan(&count))
				require.Equal(t, 1, count)
				require.Equal(t, 1, ledger.calls, "do not retry the old ledger's unknown operation")
				if !appendLost {
					var stored string
					require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_outbox`).Scan(&stored))
					require.Equal(t, original, stored)
				}
			})
		}
	}
}
func TestMixedUnknownHandoffEqualityCannotRetainUncertainWorker(t *testing.T) {
	f, p, snap, old, ledger := mixedUncertainDrainFixture(t, true)
	key := groupKey(snap)
	require.NoError(t, p.rebuildHandoffComponent(t.Context(), []groupdelivery.GroupKey{key}, map[groupdelivery.GroupKey]bool{key: true}))
	fresh := p.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
	require.NotSame(t, old, fresh, "equal progress/view is not equal transaction certainty")
	require.True(t, old.superseded)
	f.clock = at(15)
	p.TickAll(t.Context())
	r, err := f.store.Restore(t.Context(), key)
	require.NoError(t, err)
	require.Empty(t, r.Samples)
	require.Equal(t, 1, ledger.calls)
}
func TestMixedUnknownDrainForeignFrozenScopeStaysFenced(t *testing.T) {
	f, p, snap, old, _ := mixedUncertainDrainFixture(t, true)
	var raw string
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions`).Scan(&raw))
	_, err := f.store.DB().ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=json_set(payload,'$.workspace_id','foreign')`)
	require.NoError(t, err)
	f.groups.groups[0] = liveGroup(snap.Group, workspace.WriteGroupStatusDisabled)
	var observed []error
	p.config.OnError = func(err error) { observed = append(observed, err) }
	require.NoError(t, p.Reconcile(t.Context()))
	require.Same(t, old, p.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)])
	require.NotEmpty(t, observed)
	require.True(t, old.localDrain.Load())
	require.True(t, old.boundary.NeedsDurableRestore())
	var pending int
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=0`).Scan(&pending))
	require.Equal(t, 1, pending)
	_, err = f.store.DB().ExecContext(t.Context(), `UPDATE wg_runtime_versions SET payload=?`, raw)
	require.NoError(t, err)
	require.NoError(t, p.Reconcile(t.Context()))
	f.clock = at(15)
	p.TickAll(t.Context())
	r, err := f.store.Restore(t.Context(), groupKey(snap))
	require.NoError(t, err)
	require.Empty(t, r.Samples)
}
