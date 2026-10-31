package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/schema/migrations"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/testutil/lostcommit"
)

func handoffDatabasePath(t *testing.T, store *Store) string {
	t.Helper()
	var seq int
	var name, path string
	require.NoError(t, store.DB().QueryRow("PRAGMA database_list").Scan(&seq, &name, &path))
	return path
}

func TestHandoffTransactionFailureAndLostResponseAreAtomic(t *testing.T) {
	for _, lost := range []bool{false, true} {
		t.Run(map[bool]string{false: "before-commit", true: "lost-response"}[lost], func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			faults := &lostcommit.Faults{Lose: &atomic.Bool{}}
			faults.Lose.Store(lost)
			failure := errors.New("injected transfer failure")
			if !lost {
				reject := func(query string) error {
					if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
						return failure
					}
					return nil
				}
				faults.Reject.Store(&reject)
			}
			wrapped := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
			defer wrapped.Close()
			_, err := NewStore(wrapped).PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			if lost {
				require.ErrorIs(t, err, lostcommit.ErrResponseLost)
			} else {
				require.ErrorIs(t, err, failure)
			}
			if lost {
				require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"))
				require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
			} else {
				require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_samples"))
				require.Zero(t, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
			}
			_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.NoError(t, err)
			require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"))
			require.Equal(t, 1, countRows(t, store.DB(), "wg_delivery_sample_handoffs"))
		})
	}
}

func TestHandoffMissingFenceBlocksEmptyTargetClosure(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.DB().Exec("DELETE FROM write_group_ownership_fences WHERE direction='started'")
	require.NoError(t, err)
	require.ErrorIs(t, store.CommitClosure(t.Context(), otherKey, Closure{Destination: testDestination, NextClose: at(30)}), ErrOwnershipUnavailable)
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_checkpoints"))
}

func TestHandoffDamagedReceiptCannotHideOriginalOrAdvanceClosure(t *testing.T) {
	for _, damage := range []string{
		"UPDATE wg_delivery_sample_handoffs SET record=json_set(record,'$.AcceptedAt','2000-01-01T00:00:00Z')",
		"UPDATE wg_delivery_samples SET payload=json_set(payload,'$.received_at','2000-01-01T00:00:00Z') WHERE group_id='" + testKey.GroupID + "'",
		"DELETE FROM write_group_ownership_fences WHERE direction='started'",
	} {
		t.Run(damage, func(t *testing.T) {
			store, _ := seedPreparedHandoff(t)
			_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
			require.NoError(t, err)
			_, err = store.DB().Exec(damage)
			require.NoError(t, err)
			_, err = store.Restore(t.Context(), testKey)
			require.Error(t, err)
			_, err = store.PendingSampleCount(t.Context(), testKey)
			require.Error(t, err)
			require.Error(t, store.CommitClosure(t.Context(), otherKey, Closure{Destination: testDestination, NextClose: at(30)}))
			_, err = store.Reclaim(t.Context(), time.Hour)
			require.Error(t, err)
			require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"))
		})
	}
}

func TestHandoffBoundedDiscoveryAndHealthyUnrelatedProgress(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	keys, err := store.HandoffSources(t.Context(), testKey.WorkspaceID, GroupKey{}, 1)
	require.NoError(t, err)
	require.Equal(t, []GroupKey{testKey}, keys)
	keys, err = store.HandoffSources(t.Context(), "foreign", GroupKey{}, 1)
	require.NoError(t, err)
	require.Empty(t, keys)
	_, err = store.HandoffSources(t.Context(), "foreign", testKey, 1)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 129)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	_, err = store.DB().Exec("DELETE FROM wg_runtime_versions WHERE group_id=?", otherKey.GroupID)
	require.NoError(t, err)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "independent", GroupRevision: "rev-healthy"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(20), sampleAt("healthy", "temp", 23, 6)))
	require.NoError(t, store.CommitClosure(t.Context(), healthy, Closure{Destination: testDestination, NextClose: at(30)}))
}

func TestHandoffMissingReceiptCannotTurnCopiedRowIntoOrdinaryAcceptance(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	_, err = store.DB().Exec("DELETE FROM wg_delivery_sample_handoffs")
	require.NoError(t, err)
	_, err = store.Restore(t.Context(), otherKey)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "an unreceipted copy cannot become ordinary target acceptance")
}

func TestHandoffTransactionSerializesStaleClosure(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	started, release, closureAttempt := make(chan struct{}), make(chan struct{}), make(chan struct{})
	faults := &lostcommit.Faults{}
	reject := func(query string) error {
		if strings.Contains(query, "INSERT INTO wg_delivery_sample_handoffs") {
			close(started)
			<-release
		}
		return nil
	}
	faults.Reject.Store(&reject)
	writer := lostcommit.OpenFaults(handoffDatabasePath(t, store)+"?_pragma=busy_timeout(5000)", faults)
	defer writer.Close()
	transferred := make(chan error, 1)
	go func() {
		_, err := NewStore(writer).PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
		transferred <- err
	}()
	<-started
	closureFaults := &lostcommit.Faults{}
	observe := func(query string) error {
		if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at = updated_at WHERE 0") {
			close(closureAttempt)
		}
		return nil
	}
	closureFaults.Reject.Store(&observe)
	closer := lostcommit.OpenFaults(handoffDatabasePath(t, store)+"?_pragma=busy_timeout(5000)", closureFaults)
	defer closer.Close()
	closed := make(chan error, 1)
	go func() {
		closed <- NewStore(closer).CommitClosure(t.Context(), otherKey, Closure{Destination: testDestination, NextClose: at(30)})
	}()
	<-closureAttempt
	close(release)
	require.NoError(t, <-transferred)
	require.ErrorIs(t, <-closed, ErrOwnershipConflict)
	require.Zero(t, countRows(t, store.DB(), "wg_delivery_checkpoints"))
	require.Zero(t, openCount(t, store.DB(), "wg_delivery_samples", "consumed=1"))
}

func addHandoffFanoutTarget(t *testing.T, store *Store) GroupKey {
	t.Helper()
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "pending-cumulative", GroupRevision: "rev-C"}
	// A pending runtime proof still belongs to a real canonical group. Keep
	// the fixture's immutable revision FK valid before exercising lost ACKs.
	_, err := store.DB().Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'connector-1','crev-1',?)`, third.GroupID, third.WorkspaceID, third.GroupRevision, third.GroupRevision, third.GroupID, third.GroupID)
	require.NoError(t, err)
	_, err = store.DB().Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) SELECT ?,?,effective_at,json_set(payload,'$.id',?,'$.revision',?,'$.applied_revision',?,'$.recording_profile.kind','cumulative','$.destination.table_name',?) FROM write_group_versions WHERE group_id=? AND group_revision=?`, third.GroupID, third.GroupRevision, third.GroupID, third.GroupRevision, third.GroupRevision, third.GroupID, otherKey.GroupID, otherKey.GroupRevision)
	require.NoError(t, err)
	_, err = store.DB().Exec(`UPDATE write_group_cutovers SET detail=json_insert(detail,'$.started[#]',json_object('group_id',?,'revision',?))`, third.GroupID, third.GroupRevision)
	require.NoError(t, err)
	_, err = store.DB().Exec(`INSERT INTO write_group_ownership_fences SELECT fence_revision,workspace_id,direction,?,?,device_id,point_id,tag_id,'cumulative',connector_id,connector_revision,database_name,table_schema,? FROM write_group_ownership_fences WHERE direction='started'`, third.GroupID, third.GroupRevision, third.GroupID)
	require.NoError(t, err)
	var invalid int
	require.NoError(t, store.DB().QueryRow(`SELECT COUNT(*) FROM pragma_foreign_key_check`).Scan(&invalid))
	require.Zero(t, invalid, "normal fanout fixture must preserve real FK authority")
	return third
}

func TestHandoffFanoutMissingReceiptCannotReleaseOriginal(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	third := addHandoffFanoutTarget(t, store)
	batch, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, 1, batch.Accepted)
	require.Equal(t, 1, batch.Pending)
	_, err = store.DB().Exec("DELETE FROM wg_delivery_sample_handoffs WHERE target_group_id=?", third.GroupID)
	require.NoError(t, err)
	_, err = store.Restore(t.Context(), testKey)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "surviving sibling cannot prove complete disposition")
}

func TestHandoffQuotaAllowsNetNeutralMoveAtHardLimit(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	limited, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: 1})
	require.NoError(t, err)
	before, err := limited.Usage(t.Context(), "")
	require.NoError(t, err)
	_, err = limited.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err, "already accepted net-neutral ownership migration is not new intake")
	after, err := limited.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func TestHandoffQuotaRetainsPendingFanoutCharge(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	addHandoffFanoutTarget(t, store)
	before, err := store.Usage(t.Context(), "")
	require.NoError(t, err)
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.NoError(t, store.CommitClosure(t.Context(), otherKey, Closure{OwnershipView: target.OwnershipView, HandoffView: target.HandoffView, Destination: testDestination, NextClose: at(30)}))
	after, err := store.Usage(t.Context(), "")
	require.NoError(t, err)
	require.Equal(t, before, after, "pending sibling still owns durable unfinalized bytes")
}

func TestHandoffPrematureFinalizedStateCannotReleaseOpenTarget(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	rs, err := outgoingHandoffs(t.Context(), store.DB(), 1)
	require.NoError(t, err)
	require.Len(t, rs, 1)
	require.NoError(t, store.inTx(t.Context(), "test corruption", func(ctx context.Context, tx *sql.Tx) error {
		r := rs[0]
		r.State = "finalized"
		return writeHandoffReceipt(t.Context(), tx, r)
	}))
	_, err = store.Restore(t.Context(), otherKey)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "state text alone cannot attest target consumption")
}

func TestHandoffRestoreReadFailurePreservesSafeCause(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	cause := errors.New("private-db-path-and-raw-SQL")
	faults := &lostcommit.Faults{}
	reject := func(query string) error {
		if strings.Contains(query, "SELECT seq FROM wg_delivery_samples") {
			return cause
		}
		return nil
	}
	faults.Reject.Store(&reject)
	db := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
	defer db.Close()
	_, err := NewStore(db).Restore(t.Context(), testKey)
	require.ErrorIs(t, err, cause)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.NotContains(t, err.Error(), cause.Error())
}

func TestHandoffTransferredIDDoesNotReserveNativeTargetObservation(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	newObservation := sampleAt("future-accepted", ownershipMember, 35, 9)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, at(30), newObservation),
		"an original source ID in transferred history must not reserve the target's native ID namespace")
}

func TestHandoffPhysicalAddressesDoNotReserveNativeIDs(t *testing.T) {
	store, seq := seedPreparedHandoff(t)
	row, err := readHandoffSample(t.Context(), store.DB(), seq)
	require.NoError(t, err)
	receipt := newHandoffReceipt(row, row, otherKey, "f-2")
	receipt.Targets = []GroupKey{otherKey}
	receipt.ID = handoffID(receipt)
	native := sampleAt(receipt.ID, ownershipMember, 35, 9)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, at(30), native))
	_, err = store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	var address string
	require.NoError(t, store.DB().QueryRow("SELECT sample_id FROM wg_delivery_samples WHERE handoff_id=?", receipt.ID).Scan(&address))
	require.NotEqual(t, receipt.ID, address, "preexisting arbitrary native ID keeps its physical row")
	otherNative := sampleAt(address, ownershipMember, 36, 10)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, at(30), otherNative), "even an allocated handoff address is a legal native logical ID")
	require.NoError(t, store.AppendSample(t.Context(), otherKey, at(30), native))
	changed := native
	changed.ObservedAt = at(37)
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, at(30), changed), ErrSampleConflict, "native logical identity semantics remain strict")
	target, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, target.Samples, 3)
}

func TestHandoffDiscoveryKeepsDamagedScopeWithoutPoisoningHealthy(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy-journal", GroupRevision: "rev-healthy"}
	require.NoError(t, store.AppendSample(t.Context(), healthy, at(20), sampleAt("healthy", "temp", 23, 6)))
	_, err = store.DB().Exec("UPDATE wg_delivery_sample_handoffs SET record=json_set(record,'$.AcceptedAt','bad-time')")
	require.NoError(t, err)
	journals, err := store.OpenJournals(t.Context())
	require.NoError(t, err)
	found := map[GroupKey]bool{}
	for _, journal := range journals {
		found[journal.Key] = true
	}
	require.True(t, found[testKey], "damaged original remains discoverable")
	require.True(t, found[healthy])
	_, err = store.Restore(t.Context(), testKey)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	restored, err := store.Restore(t.Context(), healthy)
	require.NoError(t, err)
	require.Len(t, restored.Samples, 1)
}

func TestHandoffDiscoveryReadFailureRetainsSafeCause(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.PrepareCutoverHandoffs(t.Context(), testKey, 0, 128)
	require.NoError(t, err)
	cause := errors.New("private-source-sql-read")
	faults := &lostcommit.Faults{}
	reject := func(query string) error {
		if strings.Contains(query, "FROM wg_delivery_sample_handoffs") {
			return cause
		}
		return nil
	}
	faults.Reject.Store(&reject)
	db := lostcommit.OpenFaults(handoffDatabasePath(t, store), faults)
	defer db.Close()
	_, err = NewStore(db).OpenJournals(t.Context())
	require.ErrorIs(t, err, cause)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.NotContains(t, err.Error(), cause.Error())
}

func TestHandoffMigrationPreservesLegacyNativeRowsAndSequence(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "pre045.db"))
	require.NoError(t, err)
	defer db.Close()
	oldSchema, err := migrations.FS.ReadFile("025_write_group_delivery_sqlite.up.sql")
	require.NoError(t, err)
	_, err = db.Exec(string(oldSchema))
	require.NoError(t, err)
	sample := sampleAt("original-native", ownershipMember, 23, 7)
	payload, digest, err := encodeSample(sample)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO wg_delivery_samples(seq,workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at) VALUES(41,?,?,?,?,?,?,?,?,?,?)`, testKey.WorkspaceID, testKey.GroupID, testKey.GroupRevision, sample.SampleID, sample.MemberKey, timestamp(sample.ObservedAt), timestamp(at(20)), string(payload), digest, timestamp(at(21)))
	require.NoError(t, err)
	_, err = db.Exec("UPDATE sqlite_sequence SET seq=99 WHERE name='wg_delivery_samples'")
	require.NoError(t, err)
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	var seq int64
	var physical, gotPayload, gotDigest, accepted, kind, logical string
	require.NoError(t, db.QueryRow("SELECT seq,sample_id,payload,payload_digest,created_at,acceptance_kind,origin_sample_id FROM wg_delivery_samples").Scan(&seq, &physical, &gotPayload, &gotDigest, &accepted, &kind, &logical))
	require.Equal(t, int64(41), seq)
	require.Equal(t, sample.SampleID, physical)
	require.Equal(t, string(payload), gotPayload)
	require.Equal(t, digest, gotDigest)
	require.Equal(t, timestamp(at(21)), accepted)
	require.Equal(t, "native", kind)
	require.Empty(t, logical)
	require.NoError(t, NewStore(db).AppendSample(t.Context(), testKey, at(20), sampleAt("next-native", ownershipMember, 24, 8)))
	require.NoError(t, db.QueryRow("SELECT seq FROM wg_delivery_samples WHERE origin_sample_id='next-native'").Scan(&seq))
	require.Equal(t, int64(100), seq)
	_, err = db.Exec("UPDATE wg_delivery_samples SET acceptance_kind='handoff' WHERE seq=41")
	require.Error(t, err, "handoff origin proof metadata cannot be omitted")
	_, err = db.Exec(`INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id) SELECT workspace_id,group_id,group_revision,'another-physical-key',member_key,observed_at,bucket_start,payload,payload_digest,created_at,workspace_id,group_id,group_revision,'original-native' FROM wg_delivery_samples WHERE seq=41`)
	require.Error(t, err, "native logical uniqueness survives arbitrary physical addressing")
}

func TestHandoffRestoreCannotHideForeignWorkspaceJournalRow(t *testing.T) {
	store, _ := seedPreparedHandoff(t)
	_, err := store.DB().Exec("UPDATE wg_delivery_samples SET workspace_id='foreign',origin_workspace_id='foreign'")
	require.NoError(t, err)
	_, err = store.Restore(t.Context(), testKey)
	require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict, "a foreign row cannot be silently omitted from an original group revision")
}
