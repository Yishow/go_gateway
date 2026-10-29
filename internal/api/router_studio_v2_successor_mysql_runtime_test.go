package api

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"testing"
	"time"
)

func TestMySQLSuccessorActivationRuntimeBoundaries(t *testing.T) {
	for _, transition := range mysqlSuccessorRuntimeTransitions {
		t.Run(transition.name, func(t *testing.T) {
			r := newMySQLSuccessorRig(t, transition.family, transition.family, transition.mode, false, true)
			// A real earlier finalized/delivered row proves transition never rewrites history.
			r.clock.Store(r.snapshot.EffectiveAt.Add(time.Second).UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("history", r.snapshot.EffectiveAt.Add(time.Second), float32(4))))
			r.clock.Store(r.snapshot.EffectiveAt.Add(16 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			require.Equal(t, []string{"4"}, mysqlSuccessorRuntimeValues(t, r.f.target, r.snapshot.Group))
			historical := mysqlSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group)
			oldOutbox := successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key")
			oldReceipts := mysqlSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.NotEmpty(t, oldOutbox)
			require.NotEmpty(t, oldReceipts)
			r.clock.Store(r.boundary.Add(-2500 * time.Millisecond).UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("before", r.boundary.Add(-time.Millisecond), float32(6))))
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("equal", r.boundary, float32(7))))
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("after", r.boundary.Add(500*time.Millisecond), float32(7))))
			old, err := r.f.groups.Get(t.Context(), r.groupID)
			require.NoError(t, err)
			event := r.activate(t)
			require.NotEmpty(t, event.ID)
			current, err := r.f.groups.Get(t.Context(), r.groupID)
			require.NoError(t, err)
			require.Equal(t, old.Group, current.Group, "schedulecommit must not disable or rewrite G1")
			var head string
			require.NoError(t, r.f.db.QueryRow(`SELECT current_group_id FROM write_group_lineages WHERE lineage_id=?`, r.prepared.LineageID).Scan(&head))
			require.Equal(t, r.prepared.Group.ID, head)
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("before-after-confirm", r.boundary.Add(-2*time.Millisecond), float32(6))), "G1 still receives before-B data after config head moved")
			original := successorRuntimeOriginals(t, r.successorRuntimeRig)
			require.ErrorIs(t, r.pipe.AcceptSample(t.Context(), r.sample("not-accepted", r.boundary.Add(time.Second), float32(999))), groupdelivery.ErrOwnershipUnavailable, "delayed reconcile must not ACK unowned post-B input")
			require.Equal(t, original, successorRuntimeOriginals(t, r.successorRuntimeRig))
			var targetCount int
			require.NoError(t, r.f.db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=?`, r.prepared.Group.ID).Scan(&targetCount))
			require.Zero(t, targetCount)
			// Exact spec example: B=12:00:12.500, observed=13, reconcile=18.
			r.clock.Store(time.Date(2026, 10, 28, 12, 0, 18, 0, time.UTC).UnixNano())
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			successorRuntimeHandoffs(t, r.successorRuntimeRig, original, false)
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("native-new", r.boundary.Add(1500*time.Millisecond), float32(7))))
			r.clock.Store(time.Date(2026, 10, 28, 12, 0, 28, 0, time.UTC).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			successorRuntimeHandoffs(t, r.successorRuntimeRig, original, true)
			r.deliver(t)
			require.Equal(t, []string{"4", "6"}, mysqlSuccessorRuntimeValues(t, r.f.target, r.snapshot.Group))
			require.Equal(t, []string{"7"}, mysqlSuccessorRuntimeValues(t, r.target, r.prepared.Group))
			successorRuntimePreservedRows(t, historical, mysqlSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			successorRuntimePreservedRows(t, oldOutbox, successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key"))
			successorRuntimePreservedRows(t, oldReceipts, mysqlSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			// Reconcile/closure/delivery replay cannot materialize another logical row.
			beforeOld, beforeNew := mysqlSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group), mysqlSuccessorRuntimeRows(t, r.target, r.prepared.Group)
			newReceipts := mysqlSuccessorRuntimeEvidenceRows(t, r.target, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.NotEmpty(t, newReceipts)
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			r.pipe.TickAll(t.Context())
			r.deliver(t)
			require.Equal(t, beforeOld, mysqlSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			require.Equal(t, beforeNew, mysqlSuccessorRuntimeRows(t, r.target, r.prepared.Group))
			require.Equal(t, newReceipts, mysqlSuccessorRuntimeEvidenceRows(t, r.target, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
		})
	}
}
