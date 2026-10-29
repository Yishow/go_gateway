package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

func TestHandoffUnpreparedAcceptedInputBlocksSourceWindow(t *testing.T) {
	store, db, _ := newTestStore(t)
	original := sampleAt("future-accepted", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), original))
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	closure := Closure{OwnershipView: restored.OwnershipView, Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{handoffBoundedBucket(testKey)}}

	require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, closure), ErrOwnershipConflict,
		"accepted successor-owned input cannot be consumed before its durable handoff exists")
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
	require.Zero(t, countRows(t, db, "wg_delivery_buckets"))
	require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed=1"))
}

func TestHandoffUnpreparedAcceptedInputBlocksTargetWindow(t *testing.T) {
	store, db, _ := newTestStore(t)
	original := sampleAt("future-accepted", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), original))
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	restored, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	closure := Closure{OwnershipView: restored.OwnershipView, Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{handoffBoundedBucket(otherKey)}}

	require.ErrorIs(t, store.CommitClosure(t.Context(), otherKey, closure), ErrOwnershipConflict,
		"target cannot finalize an empty window ahead of an already accepted predecessor input")
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
	require.Zero(t, countRows(t, db, "wg_delivery_buckets"))
	require.Zero(t, openCount(t, db, "wg_delivery_samples", "consumed=1"))
}

func seedPreparedHandoff(t *testing.T) (*Store, int64) {
	t.Helper()
	store, db, _ := newTestStore(t)
	store = store.WithClock(func() time.Time { return at(21) })
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, nil, false)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, nil, true)))
	sample := sampleAt("future-accepted", ownershipMember, 23, 7)
	sample.ReceivedAt = at(21)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), sample))
	var seq int64
	require.NoError(t, db.QueryRow("SELECT seq FROM wg_delivery_samples").Scan(&seq))
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	return store, seq
}

func TestHandoffAtomicCopyPreservesBytesAndGuardsStaleAssemblers(t *testing.T) {
	store, seq := seedPreparedHandoff(t)
	db := store.DB()
	var payload, digest, accepted string
	require.NoError(t, db.QueryRow("SELECT payload,payload_digest,created_at FROM wg_delivery_samples WHERE seq=?", seq).Scan(&payload, &digest, &accepted))
	oldView, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	targetView, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	for _, key := range []GroupKey{testKey, otherKey} {
		view := oldView
		if key == otherKey {
			view = targetView
		}
		closure := Closure{OwnershipView: view.OwnershipView, HandoffView: view.HandoffView, Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{handoffBoundedBucket(key)}}
		require.ErrorIs(t, store.CommitClosure(t.Context(), key, closure), ErrOwnershipConflict)
	}
	old, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Empty(t, old.Samples)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, target.Samples, 1)
	require.Equal(t, at(21), target.Samples[0].ReceivedAt)
	var gotPayload, gotDigest, gotAccepted string
	require.NoError(t, db.QueryRow("SELECT payload,payload_digest,created_at FROM wg_delivery_samples WHERE group_id=?", otherKey.GroupID).Scan(&gotPayload, &gotDigest, &gotAccepted))
	require.Equal(t, []string{payload, digest, accepted}, []string{gotPayload, gotDigest, gotAccepted})
	batch, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	require.Equal(t, 2, countRows(t, db, "wg_delivery_samples"))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_sample_handoffs"))
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, at(20), target.Samples[0]), ErrOwnershipConflict)
	recognized, err := store.RecognizeTransferredSample(t.Context(), otherKey, target.Samples[0])
	require.NoError(t, err)
	require.True(t, recognized.Found)
	require.Equal(t, target.Samples[0], *recognized.Sample)
	closure := Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView, Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{handoffBoundedBucket(otherKey)}}
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, closure))
	require.Zero(t, openCount(t, db, "wg_delivery_samples", "group_id='"+testKey.GroupID+"' AND consumed=1"))
}

func TestHandoffPendingIsNotTargetAcceptance(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.DB().Exec("DELETE FROM wg_runtime_versions WHERE group_id=?", otherKey.GroupID)
	require.NoError(t, err)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, 1, batch.Pending)
	source, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Empty(t, source.Samples)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Empty(t, target.Samples)
	require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_samples"))
	require.ErrorIs(t, store.CommitClosure(t.Context(), otherKey, Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView, Destination: testDestination, NextClose: at(30)}), ErrOwnershipConflict)
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView, Destination: testDestination, NextClose: at(20)}), "earlier unrelated window progresses")
}

func TestHandoffMixedLegacyTargetAssociatesOnlyOwnOriginals(t *testing.T) {
	store, _, _, targetProof := handoffCounterProofFixture(t, true)
	db := store.DB()
	seedOwnershipCutover(t, db, "counter-transfer", testKey, otherKey, at(2))
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, target.Samples, 1)
	require.True(t, target.Samples[0].CounterProof.IsZero())
	legacy := target.Samples[0]
	legacy.SampleID = "own-legacy"
	legacy.ObservedAt = at(9)
	payload, digest, err := encodeSample(legacy)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at) VALUES(?,?,?,?,?,?,?,?,?,0,?)`, otherKey.WorkspaceID, otherKey.GroupID, otherKey.GroupRevision, legacy.SampleID, legacy.MemberKey, timestamp(legacy.ObservedAt), timestamp(t0), string(payload), digest, timestamp(at(9)))
	require.NoError(t, err)
	ref, err := targetProof.Reference()
	require.NoError(t, err)
	require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), otherKey, []counterproof.CounterBindingReference{ref}))
	require.Equal(t, 1, openCount(t, db, "wg_counter_proof_associations", "group_id='"+otherKey.GroupID+"' AND sample_id='own-legacy'"))
	require.Zero(t, openCount(t, db, "wg_counter_proof_associations", "group_id='"+otherKey.GroupID+"' AND sample_id<>'own-legacy'"))
	target, err = store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, target.Samples, 2)
	require.True(t, target.AuthorizesHandoff(target.Samples[0], ref))
	require.False(t, target.AuthorizesHandoff(target.Samples[1], ref))
}

func TestHandoffCounterChainRetainsTwoOriginalProofIdentities(t *testing.T) {
	store, db, _ := newTestStore(t)
	origin, a, _ := journalProofFixture(t)
	a.ObservedAt = at(43)
	bFacts := origin.CounterBindingFacts
	bFacts.Member.GroupID, bFacts.Member.AppliedRevision = otherKey.GroupID, otherKey.GroupRevision
	bProof, err := counterproof.NewCounterBindingProof(bFacts)
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, &origin, false)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, &bProof, false)))
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(40), a))
	seedOwnershipCutover(t, db, "A-B", testKey, otherKey, t0.Add(22500*time.Millisecond))
	_, err = db.Exec("UPDATE write_group_versions SET payload=json_set(payload,'$.recording_profile.kind','cumulative')")
	require.NoError(t, err)
	_, err = db.Exec("UPDATE write_group_ownership_fences SET profile_kind='cumulative'")
	require.NoError(t, err)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	b := a
	b.ObservedAt = at(45)
	b.Value = measurement.NewInt64(9)
	b.CounterProof, err = bProof.Reference()
	require.NoError(t, err)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, at(40), b))
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "counter-C", GroupRevision: "rev-C"}
	cFacts := origin.CounterBindingFacts
	cFacts.Member.GroupID, cFacts.Member.AppliedRevision = third.GroupID, third.GroupRevision
	cProof, err := counterproof.NewCounterBindingProof(cFacts)
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), third, handoffProofDescriptor(t, third, &cProof, false)))
	seedOwnershipCutover(t, db, "B-C", otherKey, third, at(40))
	_, err = db.Exec("UPDATE write_group_versions SET payload=json_set(payload,'$.recording_profile.kind','cumulative')")
	require.NoError(t, err)
	_, err = db.Exec("UPDATE write_group_ownership_fences SET profile_kind='cumulative'")
	require.NoError(t, err)
	_, err = store.PrepareCutoverHandoffs(t.Context(), otherKey, 0, 128)
	require.NoError(t, err)
	target, err := store.Restore(t.Context(), third)
	require.NoError(t, err)
	require.Len(t, target.Samples, 2)
	ref, err := cProof.Reference()
	require.NoError(t, err)
	for _, sample := range target.Samples {
		require.True(t, target.AuthorizesHandoff(sample, ref))
		originalBytes, originalDigest, err := encodeSample(sample)
		require.NoError(t, err)
		sample.AcceptanceOrigin = snapshot.AcceptanceOrigin{}
		require.False(t, target.AuthorizesHandoff(sample, ref), "codec-ignored forged origin cannot reuse private authorization")
		bytes, digest, err := encodeSample(sample)
		require.NoError(t, err)
		require.Equal(t, originalBytes, bytes)
		require.Equal(t, originalDigest, digest)
		require.ErrorIs(t, store.AppendSample(t.Context(), third, at(40), sample), ErrOwnershipConflict)
	}
	forged := sampleAt("forged", ownershipMember, 46, 4)
	forged.AcceptanceOrigin = handoffSampleOrigin(testKey)
	forged.CounterProof = ref
	require.ErrorIs(t, store.AppendSample(t.Context(), third, at(40), forged), ErrOwnershipConflict)
}

// handoffBoundedBucket retains the fixture's original aligned identity while
// expressing the actual side of its 22.5-second ownership boundary.
func handoffBoundedBucket(key GroupKey) ClosedBucket {
	bucket := silentBucket("", at(20))
	bucket.Outcome.WindowStart, bucket.Outcome.WindowEnd = at(20), at(30)
	if key == testKey {
		bucket.Outcome.WindowEnd, bucket.Outcome.Truncated = t0.Add(22500*time.Millisecond), "end"
	} else {
		bucket.Outcome.WindowStart, bucket.Outcome.Truncated = t0.Add(22500*time.Millisecond), "start"
	}
	return bucket
}
