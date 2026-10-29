package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
)

func TestSQLServerSuccessorActivationRuntimeBoundaries(t *testing.T) {
	sqlServerSuccessorFixtures(t)
	for _, transition := range sqlServerSuccessorRuntimeTransitions {
		t.Run(transition.name, func(t *testing.T) {
			r := newSQLServerSuccessorRig(t, transition.mode, false, true)
			sqlServerSuccessorRegisterEffectCleanup(t, r)
			// A real earlier finalized/delivered row proves transition never rewrites history.
			r.clock.Store(r.snapshot.EffectiveAt.Add(time.Second).UnixNano())
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("history", r.snapshot.EffectiveAt.Add(time.Second), float32(4))))
			r.clock.Store(r.snapshot.EffectiveAt.Add(16 * time.Second).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			r.deliver(t)
			require.Equal(t, []string{"4"}, sqlServerSuccessorRuntimeValues(t, r.f.target, r.snapshot.Group))
			historical := sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group)
			oldOutbox := successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key")
			oldReceipts := sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
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
			require.Equal(t, old.Group, current.Group, "schedule commit must not disable or rewrite G1")
			var head string
			require.NoError(t, r.f.db.QueryRow(`SELECT current_group_id FROM write_group_lineages WHERE lineage_id=?`, r.prepared.LineageID).Scan(&head))
			require.Equal(t, r.prepared.Group.ID, head)
			r.clock.Store(r.boundary.Add(-time.Millisecond).UnixNano())
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
			require.NoError(t, r.pipe.AcceptSample(t.Context(), r.sample("late-unfinalized", r.boundary.Add(1500*time.Millisecond), float32(7))))
			var accepted string
			require.NoError(t, r.f.db.QueryRow(`SELECT created_at FROM wg_delivery_samples WHERE group_id=? AND sample_id=?`, r.prepared.Group.ID, "late-unfinalized").Scan(&accepted))
			require.Equal(t, time.Date(2026, 10, 28, 12, 0, 18, 0, time.UTC).Format(time.RFC3339Nano), accepted, "late observation remains in its unfinalized successor window")
			r.clock.Store(time.Date(2026, 10, 28, 12, 0, 28, 0, time.UTC).UnixNano())
			r.pipe.TickAll(t.Context())
			require.Empty(t, r.reports)
			successorRuntimeHandoffs(t, r.successorRuntimeRig, original, true)
			r.deliver(t)
			require.Equal(t, []string{"4", "6"}, sqlServerSuccessorRuntimeValues(t, r.f.target, r.snapshot.Group))
			require.Equal(t, []string{"7"}, sqlServerSuccessorRuntimeValues(t, r.targetDB, r.prepared.Group))
			successorRuntimePreservedRows(t, historical, sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			successorRuntimePreservedRows(t, oldOutbox, successorRuntimeEvidenceRows(t, r.f.db, "wg_delivery_outbox", "effect_key"))
			successorRuntimePreservedRows(t, oldReceipts, sqlServerSuccessorRuntimeEvidenceRows(t, r.f.target, r.snapshot.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
			// Reconcile/closure/delivery replay cannot materialize another logical row.
			beforeOld, beforeNew := sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group), sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group)
			newReceipts := sqlServerSuccessorRuntimeEvidenceRows(t, r.targetDB, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key")
			require.NotEmpty(t, newReceipts)
			require.NoError(t, r.pipe.Reconcile(t.Context()))
			r.pipe.TickAll(t.Context())
			r.deliver(t)
			require.Equal(t, beforeOld, sqlServerSuccessorRuntimeRows(t, r.f.target, r.snapshot.Group))
			require.Equal(t, beforeNew, sqlServerSuccessorRuntimeRows(t, r.targetDB, r.prepared.Group))
			require.Equal(t, newReceipts, sqlServerSuccessorRuntimeEvidenceRows(t, r.targetDB, r.prepared.Group.Destination.TableSchema, "gw_effect_receipts", "effect_key"))
		})
	}
}
