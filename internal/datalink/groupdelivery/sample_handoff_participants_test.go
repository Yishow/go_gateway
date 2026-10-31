package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/testutil/lostcommit"
)

func sortedParticipants(keys ...GroupKey) []GroupKey {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].WorkspaceID != keys[j].WorkspaceID {
			return keys[i].WorkspaceID < keys[j].WorkspaceID
		}
		if keys[i].GroupID != keys[j].GroupID {
			return keys[i].GroupID < keys[j].GroupID
		}
		return keys[i].GroupRevision < keys[j].GroupRevision
	})
	return keys
}

func TestHandoffParticipantsAcceptedRetryAndFinalized(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	for i := 0; i < 2; i++ {
		batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
		require.NoError(t, err)
		require.Equal(t, sortedParticipants(testKey, otherKey), batch.Participants)
		require.Equal(t, sortedParticipants(testKey, otherKey), batch.ReplayParticipants)
	}
	restored, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}))
	// Exercise the finalized receipt before normal GC deletes its source row.
	require.NoError(t, store.inTx(t.Context(), "test-compact", func(ctx context.Context, tx *sql.Tx) error { return compactHandoffs(t.Context(), tx) }))
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, sortedParticipants(testKey, otherKey), batch.Participants)
	require.Equal(t, sortedParticipants(testKey, otherKey), batch.ReplayParticipants)
}

func TestHandoffParticipantsPendingFanoutAndQuarantine(t *testing.T) {
	for _, quarantine := range []bool{false, true} {
		t.Run(map[bool]string{false: "pending", true: "quarantined"}[quarantine], func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			third := addHandoffFanoutTarget(t, store)
			wantErr := error(ErrOwnershipUnavailable)
			if quarantine {
				descriptor := handoffProofDescriptor(t, third, nil, true)
				descriptor = []byte(strings.ReplaceAll(string(descriptor), "src-1", "changed-source"))
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), third, descriptor))
				wantErr = ErrHandoffProofConflict
			}
			for i := 0; i < 2; i++ {
				batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
				require.ErrorIs(t, err, wantErr)
				require.Equal(t, sortedParticipants(testKey, otherKey, third), batch.Participants)
				require.Equal(t, sortedParticipants(testKey, otherKey), batch.ReplayParticipants, "unproven sibling cannot join replay rebuild component")
			}
			if !quarantine {
				require.NoError(t, store.SaveRuntimeVersion(t.Context(), third, handoffProofDescriptor(t, third, nil, true)))
				batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
				require.NoError(t, err)
				require.Equal(t, sortedParticipants(testKey, otherKey, third), batch.ReplayParticipants, "proven promotion joins the replay component")
			}
		})
	}
}

func TestHandoffParticipantsOrdinarySourceEmpty(t *testing.T) {
	store, _, _ := newTestStore(t)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), sampleAt("ordinary", ownershipMember, 23, 7)))
	keys, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 128)
	require.NoError(t, err)
	require.Empty(t, keys, "ordinary accepted input is not handoff discovery work")
	faults := &lostcommit.Faults{}
	var writes atomic.Int32
	observe := func(query string) error {
		if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") {
			writes.Add(1)
		}
		return nil
	}
	faults.Reject.Store(&observe)
	db := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
	defer db.Close()
	store = NewStore(db)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Empty(t, batch.Participants)
	require.Empty(t, batch.ReplayParticipants)
	require.Zero(t, writes.Load(), "direct preparation also skips unrelated journals before the write lock")
}

func TestHandoffParticipantsUncertainTransactionNotPublished(t *testing.T) {
	for _, lost := range []bool{false, true} {
		t.Run(map[bool]string{false: "rollback", true: "lost-response"}[lost], func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
			faults.Lose.Store(lost)
			failure := errors.New("participants failure")
			if !lost {
				reject := func(query string) error {
					if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
						return failure
					}
					return nil
				}
				faults.Reject.Store(&reject)
			}
			db := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
			defer db.Close()
			batch, err := NewStore(db).PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			if lost {
				require.ErrorIs(t, err, lostcommit.ErrResponseLost)
			} else {
				require.ErrorIs(t, err, failure)
			}
			require.Empty(t, batch.Participants)
			require.Empty(t, batch.ReplayParticipants)
			retry, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, sortedParticipants(testKey, otherKey), retry.Participants)
			require.Equal(t, sortedParticipants(testKey, otherKey), retry.ReplayParticipants)
		})
	}
}

func TestHandoffParticipantsRetainEarlierCommittedBatchOnLaterFailure(t *testing.T) {
	store, db, _ := newTestStore(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, handoffProofDescriptor(t, testKey, nil, false)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), otherKey, handoffProofDescriptor(t, otherKey, nil, true)))
	for _, id := range []string{"one", "two"} {
		require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), sampleAt(id, ownershipMember, 23, 7)))
	}
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, at(22))
	var writes atomic.Int32
	failure := errors.New("second receipt rollback")
	faults := &lostcommit.Faults{}
	reject := func(query string) error {
		if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") && writes.Add(1) == 2 {
			return failure
		}
		return nil
	}
	faults.Reject.Store(&reject)
	wrapped := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
	defer wrapped.Close()
	batch, err := NewStore(wrapped).PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, failure)
	require.Equal(t, 1, batch.Accepted)
	require.Equal(t, sortedParticipants(testKey, otherKey), batch.Participants)
	require.Equal(t, sortedParticipants(testKey, otherKey), batch.ReplayParticipants)
	retry, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 2, retry.Accepted)
	require.Equal(t, batch.Participants, retry.Participants, "multiple original samples deduplicate complete scope identities")
}

func TestHandoffParticipantsPendingOnlyHasNoReplayComponent(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.DB().Exec("DELETE FROM wg_runtime_versions WHERE group_id=?", otherKey.GroupID)
	require.NoError(t, err)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, sortedParticipants(testKey, otherKey), batch.Participants)
	require.Empty(t, batch.ReplayParticipants)
}

func TestHandoffDiscoveryPreservesPendingFanoutAndHistoricalChain(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	pending := addHandoffFanoutTarget(t, store)
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "aaa-ordinary", GroupRevision: "rev-healthy"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(20), sampleAt("ordinary", "temp", 23, 6)))
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, 1, batch.Pending)
	keys, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 1)
	require.NoError(t, err)
	require.Equal(t, []GroupKey{testKey}, keys, "unrelated keys cannot consume the bounded page")
	keys, err = store.HandoffSources(t.Context(), testKey.WorkspaceID, testKey, 1)
	require.NoError(t, err)
	require.Equal(t, []GroupKey{otherKey}, keys)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), pending, handoffProofDescriptor(t, pending, nil, true)))
	batch, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 2, batch.Accepted, "pending sibling remains retryable")
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "chain-C", GroupRevision: "rev-C"}
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), third, handoffProofDescriptor(t, third, nil, true)))
	seedOwnershipCutover(t, store.DB(), "B-C", otherKey, third, at(23))
	_, err = store.DB().Exec("UPDATE write_groups SET status='disabled',applied_revision='new-revision' WHERE id=?", otherKey.GroupID)
	require.NoError(t, err)
	keys, err = store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 128)
	require.NoError(t, err)
	require.Equal(t, sortedParticipants(testKey, otherKey, pending), keys, "accepted intermediate and disabled historical revision remain candidates")
	batch, err = store.PrepareCutoverHandoffs(t.Context(), otherKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, sortedParticipants(otherKey, third), batch.ReplayParticipants)
	keys, err = store.HandoffSources(t.Context(), "foreign", GroupKey{}, 128)
	require.NoError(t, err)
	require.Empty(t, keys)
}

func TestHandoffDiscoveryPreservesDamagedAuthority(t *testing.T) {
	for name, damage := range map[string]string{
		"canonical-only":       "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends",
		"canonical-references": "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends; UPDATE write_group_cutovers SET source_group_id='unknown'",
		"moved-source-only":    "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends; UPDATE write_group_cutovers SET source_group_id='unknown',detail=json_set(detail,'$.ended',json('[{}]'),'$.started',json('[{}]'))",
		"source-version-only":  "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends; UPDATE write_group_cutovers SET source_group_id='group-H',detail=json_set(detail,'$.moved',json('[]'),'$.ended',json('[{}]'))",
		"orphan-end":           "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_cutovers",
		"orphan-fence":         "DELETE FROM write_group_revision_ends; DELETE FROM write_group_cutovers",
		"orphan-source-end":    "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_cutovers; UPDATE write_group_revision_ends SET group_id='group-H'",
		"orphan-source-fence":  "DELETE FROM write_group_revision_ends; DELETE FROM write_group_cutovers; DELETE FROM write_group_ownership_fences WHERE direction='ended'",
		"malformed-root":       "DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends; UPDATE write_group_cutovers SET detail='broken'",
	} {
		t.Run(name, func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			_, err := store.DB().Exec(damage)
			require.NoError(t, err)
			keys, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 128)
			require.NoError(t, err)
			require.Equal(t, []GroupKey{testKey}, keys, "unknown ownership remains discoverable, not silently skipped")
			_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
		})
	}
}

func TestHandoffDiscoveryPreservesReceiptOnlyEvidence(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	_, err = store.DB().Exec("DELETE FROM write_group_ownership_fences; DELETE FROM write_group_revision_ends; DELETE FROM write_group_cutovers; UPDATE wg_delivery_sample_handoffs SET record=json_set(record,'$.AcceptedAt','bad-time')")
	require.NoError(t, err)
	keys, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 128)
	require.NoError(t, err)
	require.Equal(t, sortedParticipants(testKey, otherKey), keys)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "damaged outgoing receipt still reaches proof validation")
	_, err = store.DB().Exec("DELETE FROM wg_delivery_sample_handoffs")
	require.NoError(t, err)
	keys, err = store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 128)
	require.NoError(t, err)
	require.Equal(t, sortedParticipants(testKey, otherKey), keys, "orphan copy and its original scope remain discoverable")
}
