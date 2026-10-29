package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func newHandoffPipelineFixture(t *testing.T, shape averageShape) *ownershipPipelineFixture {
	t.Helper()
	f := newOwnershipPipelineFixture(t, shape)
	f.clock = at(21)
	f.store = f.store.WithClock(func() time.Time { return f.clock })
	deps, config := f.pipe.deps, f.pipe.config
	deps.Store = f.store
	f.pipe = New(deps, config)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	return f
}

func handoffOriginalRows(t *testing.T, f *ownershipPipelineFixture) []string {
	t.Helper()
	rows, err := f.local.Query(`SELECT sample_id||'|'||payload||'|'||payload_digest||'|'||created_at||'|'||consumed FROM wg_delivery_samples WHERE group_id=? ORDER BY seq`, f.source.GroupID)
	require.NoError(t, err)
	defer rows.Close()
	var original []string
	for rows.Next() {
		var row string
		require.NoError(t, rows.Scan(&row))
		original = append(original, row)
	}
	require.NoError(t, rows.Err())
	return original
}

func TestCutoverHandoffCallableAdapterRestoresOneConsumerBeforeLiveReconcile(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newHandoffPipelineFixture(t, shape)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "s-after", at(23), uint64(7))))
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "s-before", f.boundary.Add(-time.Millisecond), uint64(6))))
			original := handoffOriginalRows(t, f)
			f.commit()
			f.source.EndsAt = f.boundary
			// The explicit adapter still leaves cache publication to Reconcile.
			// Build an empty target without invoking that now-live migration path.
			f.groups.groups = append(f.groups.groups, f.target.Group)
			f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
			empty, reason := f.pipe.build(t.Context(), f.target)
			require.Empty(t, reason)
			f.pipe.installManaged(empty)
			before, err := f.store.Restore(t.Context(), groupKey(f.source))
			require.NoError(t, err)
			require.Len(t, before.Samples, 2, "building a worker alone cannot transfer accepted input")
			require.NoError(t, f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID))
			old, err := f.store.Restore(t.Context(), groupKey(f.source))
			require.NoError(t, err)
			require.Len(t, old.Samples, 1)
			require.Equal(t, "s-before", old.Samples[0].SampleID)
			target, err := f.store.Restore(t.Context(), groupKey(f.target))
			require.NoError(t, err)
			require.Len(t, target.Samples, 1)
			require.Equal(t, "s-after", target.Samples[0].SampleID)
			require.Equal(t, original, handoffOriginalRows(t, f), "transfer cannot mark source consumed or rewrite original bytes")
			var sameBytes int
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs h JOIN wg_delivery_samples s ON s.seq=h.source_seq JOIN wg_delivery_samples d ON d.seq=h.target_seq WHERE s.group_id=? AND d.group_id=? AND s.payload=d.payload AND s.payload_digest=d.payload_digest AND s.created_at=d.created_at`, f.source.GroupID, f.target.GroupID).Scan(&sameBytes))
			require.Equal(t, 1, sameBytes)
			var transferID, state string
			require.NoError(t, f.local.QueryRow(`SELECT transfer_id,state FROM wg_delivery_sample_handoffs WHERE json_extract(record,'$.SampleID')='s-after'`).Scan(&transferID, &state))
			require.Equal(t, "accepted", state)
			require.NoError(t, f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID))
			var repeated string
			var count int
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*),MIN(transfer_id) FROM wg_delivery_sample_handoffs`).Scan(&count, &repeated))
			require.Equal(t, 1, count)
			require.Equal(t, transferID, repeated)
			f.clock = at(40)
			require.Error(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			require.Error(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			for _, table := range []string{"wg_delivery_checkpoints", "wg_delivery_outbox", "wg_delivery_counter_baselines"} {
				var advanced int
				require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM `+table).Scan(&advanced))
				require.Zero(t, advanced)
			}
			// A restart creates both assemblers from the verified, transactionally
			// captured receipt/input/exclusion view. Late reconciliation is not a new acceptance.
			f.pipe = New(f.pipe.deps, f.pipe.config)
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			var sourceConsumed, targetConsumed int
			require.NoError(t, f.local.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='s-after'`, f.source.GroupID).Scan(&sourceConsumed))
			require.NoError(t, f.local.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='s-after'`, f.target.GroupID).Scan(&targetConsumed))
			require.Zero(t, sourceConsumed)
			require.Equal(t, 1, targetConsumed)
			_, err = f.store.Reclaim(t.Context(), 0)
			require.NoError(t, err)
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), at(50)), "normal GC must not permanently invalidate the ownership view")
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "s-after", at(23), uint64(7))), "finalized receipt recognizes resend without recreating reclaimed accepted bytes")
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE json_extract(payload,'$.sample_id')='s-after'`).Scan(&count))
			require.Zero(t, count)
		})
	}
}

func TestCutoverHandoffMissingTargetProofIsDurablePendingNotAcceptance(t *testing.T) {
	f := newHandoffPipelineFixture(t, shapeMeanLong)
	healthy := f.readyHealthy()
	require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "future-pending", at(23), uint64(7))))
	original := handoffOriginalRows(t, f)
	f.commit()
	f.source.EndsAt = f.boundary
	_, err := f.local.Exec(`DELETE FROM wg_runtime_versions WHERE group_id=?`, f.target.GroupID)
	require.NoError(t, err)
	for range 2 {
		require.ErrorIs(t, f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID), groupdelivery.ErrOwnershipUnavailable)
	}
	var pending, targetCount int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE state='pending'`).Scan(&pending))
	require.Equal(t, 1, pending)
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&targetCount))
	require.Zero(t, targetCount)
	require.Equal(t, original, handoffOriginalRows(t, f))
	for range 2 {
		require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(healthy.Group, "independent", at(23), uint64(9))))
	}
	f.clock = at(40)
	require.Error(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
	require.NoError(t, f.pipe.boundaries[boundaryKey(healthy.GroupID, healthy.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
	var consumed int
	require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=? AND consumed=1`, healthy.GroupID).Scan(&consumed))
	require.Equal(t, 1, consumed)
	require.Equal(t, original, handoffOriginalRows(t, f))
}

func TestCutoverHandoffOriginalIDDoesNotReserveNativeTargetIdentity(t *testing.T) {
	for _, second := range []int{35, 24} {
		t.Run((time.Duration(second) * time.Second).String(), func(t *testing.T) {
			f := newHandoffPipelineFixture(t, shapeMeanLong)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "s-1", at(23), uint64(7))))
			original := handoffOriginalRows(t, f)
			f.commit()
			f.source.EndsAt = f.boundary
			f.readyTarget()
			require.NoError(t, f.pipe.PrepareCutoverHandoffs(t.Context(), f.source.WorkspaceID))
			f.clock = at(second - 1)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "s-1", at(second), uint64(9))), "a new target-native observation has a distinct original acceptance scope")
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "s-1", at(23), uint64(7))), "resend still identifies the original transferred observation")
			restored, err := f.store.Restore(t.Context(), groupKey(f.target))
			require.NoError(t, err)
			require.Len(t, restored.Samples, 2)
			require.Equal(t, original, handoffOriginalRows(t, f))
			f.pipe = New(f.pipe.deps, f.pipe.config)
			require.NoError(t, f.pipe.Reconcile(t.Context()), "restart restores both accepted origins without a selector collision")
			f.clock = at(50)
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			var consumed int
			require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE group_id=? AND consumed=1`, f.target.GroupID).Scan(&consumed))
			require.Equal(t, 2, consumed, "both distinct original observations have one target consumer")
			_, err = f.store.Reclaim(t.Context(), 0)
			require.NoError(t, err)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "s-1", at(23), uint64(7))), "GC retains enough receipt identity to acknowledge the original resend")
		})
	}
}
