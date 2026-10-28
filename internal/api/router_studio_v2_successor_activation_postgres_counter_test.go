package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgresSuccessorActivationRuntimeCounterStartsIndependentBaseline(t *testing.T) {
	for _, transition := range postgresSuccessorRuntimeTransitions {
		t.Run(transition.name, func(t *testing.T) {
			r := newPostgresSuccessorRuntimeRig(t, transition.sourceLong, transition.targetLong, true, transition.timescale)
			first := r.snapshot.EffectiveAt.Add(time.Second)
			r.clock.Store(first.UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("history-100", first, uint64(100))))
			r.clock.Store(r.snapshot.EffectiveAt.Add(61 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			oldRows := postgresSuccessorRuntimeCounterRows(t, r, r.snapshot.Group)
			require.Len(t, oldRows, 1)
			require.Equal(t, "baseline", oldRows[0].status)
			require.Equal(t, "100", oldRows[0].raw)
			require.False(t, oldRows[0].hasDelta)
			finalized := postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group)
			oldOutbox := successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key")
			oldReceipts := postgresSuccessorRuntimeEvidenceRows(t, r.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.NotEmpty(t, oldOutbox)
			require.NotEmpty(t, oldReceipts)
			r.clock.Store(r.boundary.Add(-2500 * time.Millisecond).UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("before", r.boundary.Add(-time.Millisecond), uint64(100))))
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("equal", r.boundary, uint64(107))))
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("after", r.boundary.Add(500*time.Millisecond), uint64(107))))
			original := successorRuntimeOriginals(t, r.successorRuntimeRig)
			r.activate(t)
			r.clock.Store(time.Date(2026, 10, 28, 12, 0, 18, 0, time.UTC).UnixNano())
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			successorRuntimeHandoffs(t, r.successorRuntimeRig, original, false)
			r.clock.Store(r.boundary.Add(61 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			successorRuntimeHandoffs(t, r.successorRuntimeRig, original, true)
			newRows := postgresSuccessorRuntimeCounterRows(t, r, r.prepared.Group)
			require.Len(t, newRows, 1)
			require.Equal(t, "baseline", newRows[0].status, "G1 baseline100 must not authorize G2 delta7")
			require.Equal(t, "107", newRows[0].raw)
			require.False(t, newRows[0].hasDelta, "first successor value stores SQL NULL delta, never0 or7")
			require.Empty(t, newRows[0].delta)
			require.True(t, newRows[0].deltaFrom.IsZero())
			require.True(t, newRows[0].deltaTo.IsZero())
			successorRuntimePreservedRows(t, finalized, postgresSuccessorRuntimeRows(t, r.target, r.snapshot.Group))
			successorRuntimePreservedRows(t, oldOutbox, successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key"))
			successorRuntimePreservedRows(t, oldReceipts, postgresSuccessorRuntimeEvidenceRows(t, r.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			// A later successor window establishes ordinary in-generation counting.
			later := r.boundary.Add(62 * time.Second)
			r.clock.Store(later.UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("next-110", later, uint64(110))))
			r.clock.Store(r.boundary.Add(121 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			newRows = postgresSuccessorRuntimeCounterRows(t, r, r.prepared.Group)
			require.Len(t, newRows, 2)
			require.Equal(t, "baseline", newRows[0].status)
			require.False(t, newRows[0].hasDelta)
			require.Equal(t, "ok", newRows[1].status)
			require.Equal(t, "110", newRows[1].raw)
			require.True(t, newRows[1].hasDelta)
			require.Equal(t, "3", newRows[1].delta)
		})
	}
}
