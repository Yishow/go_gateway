package groupdelivery

import (
	"encoding/json"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/testutil/lostcommit"
)

func TestHandoffConvergenceSecondEdgeFailureAndLostAcknowledgement(t *testing.T) {
	for _, lost := range []bool{false, true} {
		t.Run(map[bool]string{false: "receipt-rollback", true: "lost-acknowledgement"}[lost], func(t *testing.T) {
			store, branch, target := seedConvergentHandoff(t)
			path := handoffDatabasePath(t, store)
			failure := errors.New("second-edge-receipt-fault")
			faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
			faults.Lose.Store(lost)
			if !lost {
				reject := func(query string) error {
					if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
						return failure
					}
					return nil
				}
				faults.Reject.Store(&reject)
			}
			wrapped := lostcommit.OpenFaults(path, faults)
			batch, err := NewStore(wrapped).PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
			if lost {
				require.ErrorIs(t, err, lostcommit.ErrResponseLost)
			} else {
				require.ErrorIs(t, err, failure)
			}
			require.Empty(t, batch.Participants, "unacknowledged transaction cannot publish a rebuild participant")
			require.NoError(t, wrapped.Close())
			require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"))
			require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
			if !lost {
				require.Equal(t, 3, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
				r, err := store.Restore(t.Context(), target)
				require.NoError(t, err)
				require.ErrorIs(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: r.OwnershipView, HandoffView: r.HandoffView, Destination: testDestination, NextClose: at(30)}), ErrOwnershipConflict)
				_, err = store.Reclaim(t.Context(), time.Hour)
				require.NoError(t, err)
				require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"))
			}
			require.NoError(t, store.DB().Close())
			store = NewStore(openStoreDB(t, path))
			for i := 0; i < 2; i++ {
				batch, err = store.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
				require.NoError(t, err)
				require.Equal(t, 1, batch.Accepted)
			}
			require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"))
			require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
			r, err := store.Restore(t.Context(), target)
			require.NoError(t, err)
			consumeConvergentHandoff(t, store, target, r)
		})
	}
}

func convergencePeerReceipts(t *testing.T, store *Store, target GroupKey) (handoffReceipt, handoffReceipt) {
	t.Helper()
	rs, err := queryHandoffs(t.Context(), store.DB(), "target_group_id=? AND target_revision=?", target.GroupID, target.GroupRevision)
	require.NoError(t, err)
	require.Len(t, rs, 2)
	if rs[0].CanonicalTransferID != "" {
		return rs[1], rs[0]
	}
	return rs[0], rs[1]
}

func replaceConvergenceReceipt(t *testing.T, store *Store, original string, r handoffReceipt) {
	t.Helper()
	r.ID = handoffID(r)
	raw, err := json.Marshal(r)
	require.NoError(t, err)
	_, err = store.DB().Exec(`UPDATE wg_delivery_sample_handoffs SET transfer_id=?,origin_seq=?,record=? WHERE transfer_id=?`, r.ID, r.OriginSeq, string(raw), original)
	require.NoError(t, err)
}

func TestHandoffConvergenceRejectsEveryCorruptPeerLiveAndAfterGC(t *testing.T) {
	mutations := map[string]func(*testing.T, *Store, handoffReceipt, handoffReceipt){
		"missing-canonical": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			_, err := s.DB().Exec(`DELETE FROM wg_delivery_sample_handoffs WHERE transfer_id=?`, canonical.ID)
			require.NoError(t, err)
		},
		"malformed-alias": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_sample_handoffs SET record='{}' WHERE transfer_id=?`, alias.ID)
			require.NoError(t, err)
		},
		"missing-reference": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.CanonicalTransferID = ""
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"reference-cycle": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.CanonicalTransferID = alias.ID
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"canonical-is-alias": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			canonical.CanonicalTransferID = alias.ID
			replaceConvergenceReceipt(t, s, canonical.ID, canonical)
		},
		"origin-sequence": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.OriginSeq++
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"origin-scope": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.Origin.GroupRevision = "different-origin"
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"sample-identity": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.SampleID = "different-sample"
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"raw-bytes": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.RawDigest = handoffHash("different bytes")
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"target-proof": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.TargetProof.Digest = "different-proof"
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"source-bucket": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			alias.SourceBucket = timestamp(at(50))
			replaceConvergenceReceipt(t, s, alias.ID, alias)
		},
		"branch-lineage": func(t *testing.T, s *Store, canonical, alias handoffReceipt) {
			_, err := s.DB().Exec(`DELETE FROM write_group_revision_ends WHERE group_id=? AND group_revision=?`, alias.Source.GroupID, alias.Source.GroupRevision)
			require.NoError(t, err)
		},
	}
	for name, mutate := range mutations {
		for _, reclaimed := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "/live", true: "/reclaimed"}[reclaimed], func(t *testing.T) {
				store, _, target, restored := finishConvergentHandoff(t)
				if reclaimed {
					consumeConvergentHandoff(t, store, target, restored)
				}
				canonical, alias := convergencePeerReceipts(t, store, target)
				mutate(t, store, canonical, alias)
				_, err := store.Restore(t.Context(), target)
				require.Error(t, err, "a healthy canonical cannot mask a corrupt peer")
				_, err = store.RecognizeTransferredSample(t.Context(), target, restored.Samples[0])
				require.Error(t, err, "recognition must inspect all shared peers before returning")
				_, err = incomingHandoff(t.Context(), store.DB(), canonical.TargetSeq)
				require.Error(t, err)
			})
		}
	}
}

func TestHandoffConvergenceMissingAliasAfterGCRejectsIncompleteFanout(t *testing.T) {
	for _, checkpoint := range []bool{false, true} {
		t.Run(map[bool]string{false: "no-branch-checkpoint", true: "legitimate-empty-branch-checkpoint"}[checkpoint], func(t *testing.T) {
			store, branch, target := seedConvergentHandoffAtInterval(t, 30)
			_, err := store.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
			require.NoError(t, err)
			restored, err := store.Restore(t.Context(), target)
			require.NoError(t, err)
			if checkpoint {
				prior, err := store.Restore(t.Context(), branch)
				require.NoError(t, err)
				require.Empty(t, prior.Samples)
				bucket := silentBucket("", at(0))
				bucket.Outcome.BucketEnd = at(30)
				bucket.Outcome.WindowStart, bucket.Outcome.WindowEnd = at(10), at(20)
				bucket.Outcome.Truncated = "end"
				require.NoError(t, store.CommitClosure(t.Context(), branch, Closure{OwnershipView: prior.OwnershipView, HandoffView: prior.HandoffView, Destination: testDestination, NextClose: at(30), Buckets: []ClosedBucket{bucket}}))
				require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"), "the empty pre-cutover branch window cannot consume the future observation")
			}
			consumeConvergentHandoff(t, store, target, restored)
			_, alias := convergencePeerReceipts(t, store, target)
			require.Equal(t, branch, alias.Source)
			_, err = store.DB().Exec(`DELETE FROM wg_delivery_sample_handoffs WHERE transfer_id=?`, alias.ID)
			require.NoError(t, err)
			_, err = store.Restore(t.Context(), target)
			require.Error(t, err, "the surviving canonical cannot hide a lost ancestor fanout continuation")
			_, err = store.RecognizeTransferredSample(t.Context(), target, restored.Samples[0])
			require.Error(t, err, "a predecessor empty-window checkpoint is not consumption authority for an out-of-bounds future observation")
		})
	}
}

func TestHandoffConvergenceConsumedTargetRequiresRelevantAliases(t *testing.T) {
	t.Run("missing-alias-after-consumption-before-gc", func(t *testing.T) {
		store, _, target, restored := finishConvergentHandoff(t)
		require.NoError(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}))
		require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"), "raw authority still exists before GC")
		require.Equal(t, 1, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
		_, alias := convergencePeerReceipts(t, store, target)
		_, err := store.DB().Exec(`DELETE FROM wg_delivery_sample_handoffs WHERE transfer_id=?`, alias.ID)
		require.NoError(t, err)
		_, err = store.Restore(t.Context(), target)
		require.Error(t, err, "consumption cannot hide a missing required converging edge")
		_, err = store.RecognizeTransferredSample(t.Context(), target, restored.Samples[0])
		require.Error(t, err)
	})
	t.Run("unconsumed-first-edge-remains-readable", func(t *testing.T) {
		store, _, target := seedConvergentHandoff(t)
		restored, err := store.Restore(t.Context(), target)
		require.NoError(t, err)
		require.Len(t, restored.Samples, 1)
		require.ErrorIs(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}), ErrOwnershipConflict)
		require.Equal(t, 3, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
		require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
	})
	t.Run("independent-pending-profile-is-not-a-consumption-obligation", func(t *testing.T) {
		store, _ := seedPreparedHandoff(t)
		independent := addHandoffFanoutTarget(t, store)
		batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
		require.ErrorIs(t, err, ErrOwnershipUnavailable)
		require.Equal(t, 1, batch.Accepted)
		require.Equal(t, 1, batch.Pending)
		target := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "independent-successor", GroupRevision: "rev-D"}
		require.NoError(t, store.SaveRuntimeVersion(t.Context(), target, handoffProofDescriptor(t, target, nil, true)))
		seedOwnershipCutover(t, store.DB(), "I-D-only", otherKey, target, at(23))
		_, err = store.PrepareCutoverHandoffs(t.Context(), otherKey, 0, 128)
		require.NoError(t, err)
		restored, err := store.Restore(t.Context(), target)
		require.NoError(t, err)
		require.Len(t, restored.Samples, 1)
		ownership, err := store.ResolveCutoverOwnership(t.Context(), sourceOf(testKey, restored.Samples[0]), restored.Samples[0].ObservedAt)
		require.NoError(t, err)
		require.False(t, handoffDescendant(ownership, target, independent, map[GroupKey]bool{}))
		require.NoError(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}))
		for i := 0; i < 2; i++ {
			after, err := store.Restore(t.Context(), target)
			require.NoError(t, err, "an unrelated profile's pending proof must not poison the consumed owner")
			require.Empty(t, after.Samples)
			recognized, err := store.RecognizeTransferredSample(t.Context(), target, restored.Samples[0])
			require.NoError(t, err)
			require.True(t, recognized.Found)
			require.True(t, recognized.Finalized)
			require.Nil(t, recognized.Sample)
			_, err = store.Reclaim(t.Context(), time.Hour)
			require.NoError(t, err)
		}
		require.Equal(t, 1, openCount(t, store.DB(), "wg_delivery_sample_handoffs", "state='pending'"))
		require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"), "the pending independent profile keeps its required raw pins")
	})
}
