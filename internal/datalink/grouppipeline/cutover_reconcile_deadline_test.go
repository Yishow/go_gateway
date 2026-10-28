package grouppipeline

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

func newDeadlinePipeline(t *testing.T, shape averageShape) *ownershipPipelineFixture {
	t.Helper()
	f := newOwnershipPipelineFixture(t, shape)
	raw, err := json.Marshal(f.target)
	require.NoError(t, err)
	var target workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(raw, &target))
	// Freeze a new original revision once, before acceptance. The previously
	// frozen fixture revision is never overwritten or reinterpreted.
	target.AppliedRevision, target.GroupRevision = "deadline-rev", "deadline-rev"
	target.Group.Revision, target.Group.AppliedRevision = target.AppliedRevision, target.AppliedRevision
	target.EffectiveAt = at(10)
	target.Group.RowPolicy.AllowedLatenessSeconds = 1
	target.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
	for i := range target.Group.Members {
		target.Group.Members[i].ResultValueType = measurement.ExactUint64
	}
	f.target, f.boundary = &target, at(10)
	require.NoError(t, f.pipe.saveRuntimeSnapshot(t.Context(), f.target))
	f.store = f.store.WithClock(func() time.Time { return f.clock })
	deps, config := f.pipe.deps, f.pipe.config
	deps.Store, config.Sender.Recovery = f.store, nil
	f.pipe = New(deps, config)
	f.clock = at(9)
	f.commit() // Future boundary10 commits at9, before any target acceptance.
	f.source.EndsAt = f.boundary
	f.groups.groups = append(f.groups.groups, f.target.Group)
	f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)])
	return f
}

func deadlineAcceptedBytes(t *testing.T, f *ownershipPipelineFixture, id string, accepted time.Time) [3]string {
	t.Helper()
	var result [3]string
	require.NoError(t, f.local.QueryRow(`SELECT payload,payload_digest,created_at FROM wg_delivery_samples
 WHERE group_id=? AND group_revision=? AND COALESCE(NULLIF(origin_sample_id,''),sample_id)=?`, f.target.GroupID, f.target.AppliedRevision, id).Scan(&result[0], &result[1], &result[2]))
	created, err := time.Parse(time.RFC3339Nano, result[2])
	require.NoError(t, err)
	require.Equal(t, accepted, created, "real Store.WithClock acceptance time, not envelope ReceivedAt or wall clock")
	return result
}

func TestCutoverDeadlineOriginalAccepted14RemainsEligibleAfterReopen20(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newDeadlinePipeline(t, shape)
			f.clock = at(14)
			input := envelope(f.target.Group, "s-1", at(11), uint64(7))
			input.ReceivedAt = at(13)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), input))
			original := deadlineAcceptedBytes(t, f, "s-1", at(14))
			descriptors := cutoverCrashDescriptors(t, f)
			restored, err := f.store.Restore(t.Context(), groupKey(f.target))
			require.NoError(t, err)
			require.Len(t, restored.Samples, 1)
			require.Equal(t, at(11), restored.Samples[0].ObservedAt)
			require.Equal(t, at(13), restored.Samples[0].ReceivedAt)
			f.clock = at(15)
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			require.Empty(t, cutoverCrashEffects(t, f.local), "lateness1 keeps [10,15) open until deadline16")
			f.clock = at(20)
			reopenCutoverPipeline(t, f) // Physical Close/Open/Migrate; no old workers survive.
			fresh := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary
			require.NoError(t, fresh.Tick(t.Context(), f.clock))
			require.Equal(t, at(15), fresh.NextClose())
			require.Equal(t, original, deadlineAcceptedBytes(t, f, "s-1", at(14)))
			f.deliver()
			value, evidence := readReductionRow(t, f, f.target, at(10), "latest")
			require.Equal(t, "7", value)
			require.Equal(t, "s-1", evidence.SampleID)
			require.Equal(t, at(10).Format("2006-01-02T15:04:05.000000Z"), evidence.Window.Start)
			require.Equal(t, at(15).Format("2006-01-02T15:04:05.000000Z"), evidence.Window.End)
			observed, err := time.Parse(time.RFC3339Nano, evidence.ObservedAt)
			require.NoError(t, err)
			require.Equal(t, at(11), observed)
			var consumed int
			require.NoError(t, f.local.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&consumed))
			require.Equal(t, 1, consumed)
			effects := cutoverCrashEffects(t, f.local)
			reopenCutoverPipeline(t, f)
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			require.Equal(t, effects, cutoverCrashEffects(t, f.local), "another recovery cannot reopen finalized output")
			require.Equal(t, descriptors, cutoverCrashDescriptors(t, f))
		})
	}
}

func TestCutoverDeadlineFirstArrival17CannotReopenFinalizedWindow(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			f := newDeadlinePipeline(t, shape)
			f.clock = at(12)
			require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.target.Group, "already-accepted", at(11), uint64(3))))
			f.clock = at(16)
			worker := f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary
			require.NoError(t, worker.Tick(t.Context(), f.clock))
			require.Equal(t, at(15), worker.NextClose(), "the deadline really finalized the original window")
			f.deliver()
			originals, consumers := cutoverCrashJournal(t, f.local)
			effects := cutoverCrashEffects(t, f.local)
			require.NotEmpty(t, effects)
			f.clock = at(17)
			late := envelope(f.target.Group, "first-arrives-17", at(11), uint64(99))
			late.ReceivedAt = at(17)
			refuse := func() {
				t.Helper()
				var rejected *runtime.GroupSampleError
				require.ErrorAs(t, f.pipe.AcceptSample(t.Context(), late), &rejected)
				require.Equal(t, snapshot.OfferLate, rejected.Outcome)
				require.Equal(t, snapshot.ReasonLate, rejected.Reason)
			}
			refuse()
			after, afterConsumers := cutoverCrashJournal(t, f.local)
			require.Equal(t, originals, after, "first late arrival has no durable acceptance and cannot ACK")
			require.Equal(t, consumers, afterConsumers)
			require.Equal(t, effects, cutoverCrashEffects(t, f.local))
			f.clock = at(20)
			reopenCutoverPipeline(t, f)
			refuse()
			require.NoError(t, f.pipe.boundaries[boundaryKey(f.target.GroupID, f.target.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
			f.deliver()
			value, evidence := readReductionRow(t, f, f.target, at(10), "latest")
			require.Equal(t, "3", value)
			require.Equal(t, "already-accepted", evidence.SampleID)
			require.Equal(t, effects, cutoverCrashEffects(t, f.local), "late diagnostics never reopen or rewrite finalized bytes/destination")
			after, afterConsumers = cutoverCrashJournal(t, f.local)
			require.Equal(t, originals, after)
			require.Equal(t, consumers, afterConsumers)
		})
	}
}
