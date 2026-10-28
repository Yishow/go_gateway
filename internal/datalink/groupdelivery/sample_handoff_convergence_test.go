package groupdelivery

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Keep the Store fixture on the existing immutable-input comparator. The public
// API companion covers the genuine counter_raw / cumulative proof matrix.
func seedConvergentHandoff(t *testing.T) (*Store, GroupKey, GroupKey) {
	t.Helper()
	return seedConvergentHandoffAtInterval(t, 10)
}

func seedConvergentHandoffAtInterval(t *testing.T, branchSeconds int) (*Store, GroupKey, GroupKey) {
	t.Helper()
	store, db, _ := newTestStore(t)
	store = store.WithClock(func() time.Time { return at(8) })
	for _, key := range []GroupKey{testKey, otherKey} {
		require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, handoffProofDescriptor(t, key, nil, true)))
	}
	sample := sampleAt("same-origin", ownershipMember, 25, 7)
	sample.ReceivedAt = at(8)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(20), sample))
	seedOwnershipCutover(t, db, "forward", testKey, otherKey, at(10))
	branch := addHandoffFanoutTarget(t, store)
	descriptor := handoffProofDescriptor(t, branch, nil, true)
	descriptor = []byte(strings.ReplaceAll(string(descriptor), `"interval_seconds":10`, fmt.Sprintf(`"interval_seconds":%d`, branchSeconds)))
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), branch, descriptor))
	_, err := db.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.row_policy.interval_seconds',?) WHERE group_id=? AND group_revision=?`, branchSeconds, branch.GroupID, branch.GroupRevision)
	require.NoError(t, err)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 2, batch.Accepted)
	target := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: "restored-3"}
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), target, handoffProofDescriptor(t, target, nil, true)))
	seedOwnershipCutover(t, db, "rollback", otherKey, target, at(20))
	_, err = db.Exec(`UPDATE write_group_cutovers SET detail=json_insert(detail,'$.ended[#]',json_object('group_id',?,'revision',?)) WHERE id='rollback'`, branch.GroupID, branch.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?, 'rollback')`, branch.WorkspaceID, branch.GroupID, branch.GroupRevision, timestamp(at(20)))
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences SELECT 'rollback',workspace_id,'ended',group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name FROM write_group_ownership_fences WHERE fence_revision='forward' AND group_id=?`, branch.GroupID)
	require.NoError(t, err)
	batch, err = store.PrepareCutoverHandoffs(t.Context(), otherKey, 0, 128)
	require.NoError(t, err)
	require.Equal(t, 1, batch.Accepted)
	return store, branch, target
}

func TestHandoffConvergenceOneCanonicalJournal(t *testing.T) {
	store, branch, target := seedConvergentHandoff(t)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
	require.NoError(t, err, "second independently proven edge must reuse the same original's target journal")
	require.Equal(t, 1, batch.Accepted)
	require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"))
	require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
	var targets int
	require.NoError(t, store.DB().QueryRow(`SELECT COUNT(DISTINCT target_seq) FROM wg_delivery_sample_handoffs WHERE target_group_id=? AND target_revision=?`, target.GroupID, target.GroupRevision).Scan(&targets))
	require.Equal(t, 1, targets)
	restored, err := store.Restore(t.Context(), target)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	require.Equal(t, at(8), restored.Samples[0].ReceivedAt)
	for _, key := range []GroupKey{testKey, otherKey, branch} {
		prior, err := store.Restore(t.Context(), key)
		require.NoError(t, err)
		require.Empty(t, prior.Samples)
	}
	require.NoError(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}))
	require.Equal(t, 1, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
}

func finishConvergentHandoff(t *testing.T) (*Store, GroupKey, GroupKey, Restored) {
	t.Helper()
	store, branch, target := seedConvergentHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
	require.NoError(t, err)
	restored, err := store.Restore(t.Context(), target)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
	return store, branch, target, restored
}

func consumeConvergentHandoff(t *testing.T, store *Store, target GroupKey, restored Restored) {
	t.Helper()
	require.NoError(t, store.CommitClosure(t.Context(), target, Closure{OwnershipView: restored.OwnershipView, HandoffView: restored.HandoffView, Destination: testDestination, NextClose: at(30)}))
	require.Equal(t, 1, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"), "only the shared leaf consumes")
	_, err := store.Reclaim(t.Context(), time.Hour)
	require.NoError(t, err)
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_samples"))
	require.Equal(t, 4, openCount(t, store.DB(), "wg_delivery_sample_handoffs", "state='finalized'"))
}

func TestHandoffConvergenceReopenReclaimAndStableView(t *testing.T) {
	store, branch, target, restored := finishConvergentHandoff(t)
	path := handoffDatabasePath(t, store)
	var originalPayload, originalDigest, originalTime string
	require.NoError(t, store.DB().QueryRow(`SELECT payload,payload_digest,created_at FROM wg_delivery_samples WHERE acceptance_kind='native'`).Scan(&originalPayload, &originalDigest, &originalTime))
	var identical int
	require.NoError(t, store.DB().QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples WHERE payload=? AND payload_digest=? AND created_at=?`, originalPayload, originalDigest, originalTime).Scan(&identical))
	require.Equal(t, 4, identical, "raw bytes, proof, original received/accepted times survive every branch")
	require.NoError(t, store.DB().Close())
	store = NewStore(openStoreDB(t, path))
	for _, key := range []GroupKey{testKey, otherKey, branch} {
		_, err := store.PrepareCutoverHandoffs(t.Context(), key, 0, 128)
		require.NoError(t, err)
	}
	_, err := store.Reclaim(t.Context(), time.Hour)
	require.NoError(t, err)
	require.Equal(t, 4, countRows(t, store.DB(), "wg_delivery_samples"), "unconsumed target retains every pin")
	consumeConvergentHandoff(t, store, target, restored)
	after, err := store.Restore(t.Context(), target)
	require.NoError(t, err)
	require.Equal(t, restored.HandoffView, after.HandoffView)
	require.Empty(t, after.Samples)
	require.NoError(t, store.DB().Close())
	store = NewStore(openStoreDB(t, path))
	for _, key := range []GroupKey{testKey, otherKey, branch, target} {
		_, err := store.PrepareCutoverHandoffs(t.Context(), key, 0, 128)
		require.NoError(t, err)
		_, err = store.Restore(t.Context(), key)
		require.NoError(t, err)
	}
	recognized, err := store.RecognizeTransferredSample(t.Context(), target, restored.Samples[0])
	require.NoError(t, err)
	require.True(t, recognized.Found)
	require.True(t, recognized.Finalized)
	require.Nil(t, recognized.Sample)
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_samples"))
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_outbox"), "GC and replay create no synthetic effects")
}

func TestHandoffConvergenceQuotaAtCapacityCountsCanonicalOnce(t *testing.T) {
	store, branch, target := seedConvergentHandoff(t)
	before, err := store.Usage(t.Context(), target.GroupID)
	require.NoError(t, err)
	limited, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: 1, GroupMaxBytes: before, WarningRatio: 1, HardRatio: 1})
	require.NoError(t, err)
	_, err = limited.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
	require.NoError(t, err, "the second edge creates no target payload charge")
	after, err := limited.Usage(t.Context(), target.GroupID)
	require.NoError(t, err)
	require.Equal(t, before, after)
	global, err := limited.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, after, global, "both intermediate copies are disposed, not double charged")
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy", GroupRevision: "v1"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(20), sampleAt("healthy", "temperature", 23, 2)))
	r, err := store.Restore(t.Context(), healthy)
	require.NoError(t, err)
	require.Len(t, r.Samples, 1)
}

func TestHandoffConvergencePreservesIndependentSourceBuckets(t *testing.T) {
	store, branch, target := seedConvergentHandoffAtInterval(t, 5)
	_, err := store.PrepareCutoverHandoffs(t.Context(), branch, 0, 128)
	require.NoError(t, err)
	canonical, alias := convergencePeerReceipts(t, store, target)
	require.NotEqual(t, canonical.SourceBucket, alias.SourceBucket)
	require.Equal(t, canonical.TargetBucket, alias.TargetBucket)
	restored, err := store.Restore(t.Context(), target)
	require.NoError(t, err)
	consumeConvergentHandoff(t, store, target, restored)
	_, err = store.Restore(t.Context(), target)
	require.NoError(t, err, "each compact branch proves its own source bucket, not its peer's bucket")
}

func TestHandoffConvergenceCanonicalIsIndependentOfSQLOrder(t *testing.T) {
	store, _, target, restored := finishConvergentHandoff(t)
	store.DB().SetMaxOpenConns(1)
	canonical, alias := convergencePeerReceipts(t, store, target)
	_, err := store.DB().Exec(`PRAGMA reverse_unordered_selects=ON`)
	require.NoError(t, err)
	tx, err := store.DB().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	peers, err := queryHandoffs(t.Context(), tx, "target_seq=?", canonical.TargetSeq)
	require.NoError(t, err)
	require.Equal(t, alias.ID, peers[0].ID, "fixture must actually put the alias first")
	incoming, err := incomingHandoff(t.Context(), tx, canonical.TargetSeq)
	require.NoError(t, err)
	require.Equal(t, canonical.ID, incoming.ID)
	require.NoError(t, tx.Commit())
	consumeConvergentHandoff(t, store, target, restored)
	tx, err = store.DB().BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	incoming, err = incomingHandoff(t.Context(), tx, canonical.TargetSeq)
	require.NoError(t, err)
	require.Equal(t, canonical.ID, incoming.ID, "compact references retain the same primary after raw-row GC")
}
