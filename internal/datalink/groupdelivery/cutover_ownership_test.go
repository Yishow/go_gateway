package groupdelivery

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/testutil/lostcommit"
)

var ownershipMember = `["device-1","point-1","tag-1"]`

func seedOwnershipCutover(t *testing.T, db *sql.DB, cutoverID string, oldKey, newKey GroupKey, boundary time.Time) {
	t.Helper()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	for i, key := range []GroupKey{oldKey, newKey} {
		kind := ""
		effective := t0
		if i == 1 {
			kind = "instantaneous"
			effective = boundary
		}
		payload := map[string]any{"id": key.GroupID, "workspace_id": key.WorkspaceID, "revision": key.GroupRevision, "applied_revision": key.GroupRevision,
			"members":           []map[string]string{{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1"}},
			"destination":       map[string]string{"connector_id": "connector-1", "connector_revision": "crev-1", "database": "db", "table_schema": "main", "table_name": key.GroupID},
			"recording_profile": map[string]string{"kind": kind}, "row_policy": map[string]int{"interval_seconds": 10}}
		encoded, err := json.Marshal(payload)
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(id) DO NOTHING`, key.GroupID, key.WorkspaceID, key.GroupRevision, key.GroupRevision, key.GroupID, "connector-1", "crev-1", key.GroupID)
		require.NoError(t, err)
		_, err = tx.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,?,?,?) ON CONFLICT(group_id,group_revision) DO NOTHING`, key.GroupID, key.GroupRevision, effective, string(encoded))
		require.NoError(t, err)
	}
	detail, err := json.Marshal(map[string]any{"ownership_fence_revision": cutoverID, "ended": []map[string]string{{"group_id": oldKey.GroupID, "revision": oldKey.GroupRevision}}, "started": []map[string]string{{"group_id": newKey.GroupID, "revision": newKey.GroupRevision}}, "moved": []map[string]string{{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1"}}})
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES(?,?,'forward',?,?,?,?,'{}')`, cutoverID, oldKey.WorkspaceID, oldKey.GroupID, boundary.Format(time.RFC3339Nano), cutoverID, string(detail))
	require.NoError(t, err)
	_, err = tx.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,?)`, oldKey.WorkspaceID, oldKey.GroupID, oldKey.GroupRevision, boundary.Format(time.RFC3339Nano), cutoverID)
	require.NoError(t, err)
	for i, key := range []GroupKey{oldKey, newKey} {
		kind, direction := "", "ended"
		if i == 1 {
			kind, direction = "instantaneous", "started"
		}
		require.NoError(t, tx.QueryRow(`SELECT COALESCE(json_extract(payload,'$.recording_profile.kind'),'') FROM write_group_versions WHERE group_id=? AND group_revision=?`, key.GroupID, key.GroupRevision).Scan(&kind))
		_, err = tx.Exec(`INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES(?,?,?,?,?,'device-1','point-1','tag-1',?,'connector-1','crev-1','db','main',?)`, cutoverID, key.WorkspaceID, direction, key.GroupID, key.GroupRevision, kind, key.GroupID)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
}

func TestCutoverOwnershipAppendHalfOpenWithoutReconcile(t *testing.T) {
	store, db, _ := newTestStore(t)
	boundary := t0.Add(22500 * time.Millisecond)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, boundary)
	for i, when := range []time.Time{t0.Add(23 * time.Second), boundary, t0.Add(22499 * time.Millisecond)} {
		sample := sampleAt([]string{"after", "exact", "before"}[i], ownershipMember, 0, []int64{7, 8, 6}[i])
		sample.ObservedAt = when
		oldErr := store.AppendSample(t.Context(), testKey, t0, sample)
		newErr := store.AppendSample(t.Context(), otherKey, t0, sample)
		if i == 2 {
			require.NoError(t, oldErr)
			require.ErrorContains(t, newErr, "recording-ownership-conflict")
		} else {
			require.ErrorContains(t, oldErr, "recording-ownership-conflict")
			require.NoError(t, newErr)
		}
	}
	old, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, old.Samples, 1)
	require.Equal(t, measurement.NewInt64(6), old.Samples[0].Value)
	next, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Len(t, next.Samples, 2)
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
}

func TestCutoverOwnershipCannotReappendAcceptedPredecessor(t *testing.T) {
	store, db, _ := newTestStore(t)
	sample := sampleAt("accepted-before-cutover", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	var payload, digest string
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&payload, &digest))
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sample), ErrOwnershipConflict)
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, sample), ErrOwnershipConflict, "accepted data requires separate durable transfer authority")
	var afterPayload, afterDigest string
	require.NoError(t, db.QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples`).Scan(&afterPayload, &afterDigest))
	require.Equal(t, payload, afterPayload)
	require.Equal(t, digest, afterDigest)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipUnavailableNeverMeansUnfenced(t *testing.T) {
	cases := map[string]string{
		"all fence rows missing":  `DELETE FROM write_group_ownership_fences`,
		"target fence missing":    `DELETE FROM write_group_ownership_fences WHERE direction='started'`,
		"source fence missing":    `DELETE FROM write_group_ownership_fences WHERE direction='ended'`,
		"root token missing":      `UPDATE write_group_cutovers SET detail=json_remove(detail,'$.ownership_fence_revision')`,
		"malformed root":          `UPDATE write_group_cutovers SET detail='not json'`,
		"wrong destination":       `UPDATE write_group_ownership_fences SET table_name='another-table' WHERE direction='started'`,
		"end missing":             `DELETE FROM write_group_revision_ends`,
		"target snapshot missing": `DELETE FROM write_group_versions WHERE group_id='group-H'`,
		"orphan revision end":     `DELETE FROM write_group_ownership_fences; DELETE FROM write_group_cutovers`,
	}
	for name, mutation := range cases {
		t.Run(name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
			_, err := db.Exec(mutation)
			require.NoError(t, err)
			source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
			_, err = store.ResolveCutoverOwnership(t.Context(), source, at(23))
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("old", ownershipMember, 23, 7)), ErrOwnershipUnavailable)
			require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, sampleAt("new", ownershipMember, 23, 7)), ErrOwnershipUnavailable)
			healthy := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "healthy", GroupRevision: "healthy-1"}
			require.NoError(t, store.AppendSample(t.Context(), healthy, t0, sampleAt("healthy", `["device-2","point-2","tag-2"]`, 23, 4)), "demonstrably unrelated scope progresses")
			require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
			require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
		})
	}
}

func TestCutoverOwnershipMissingTablesFailsClosed(t *testing.T) {
	for _, table := range []string{"write_group_cutovers", "write_group_ownership_fences", "write_group_revision_ends"} {
		t.Run(table, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			_, err := db.Exec(`DROP TABLE ` + table)
			require.NoError(t, err)
			err = store.AppendSample(t.Context(), testKey, t0, sampleAt("never-ack", "temp", 2, 1))
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			var typed *OwnershipError
			require.ErrorAs(t, err, &typed)
			require.Equal(t, "recording-ownership-unavailable", err.Error())
			require.Zero(t, countRows(t, db, "wg_delivery_samples"))
		})
	}
}

func TestCutoverOwnershipResendNewTimeIsNotTransfer(t *testing.T) {
	store, db, _ := newTestStore(t)
	prior := sampleAt("reused-id", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, prior))
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	next := prior
	next.ObservedAt = at(24)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, t0, next))
	next.ObservedAt = prior.ObservedAt
	next.Value = measurement.NewInt64(999)
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, next), ErrOwnershipConflict, "changed payload cannot disguise predecessor observation")
	require.Equal(t, 2, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipDamagedIndependentProfileDoesNotBlockHealthyKey(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	_, err := db.Exec(`DELETE FROM write_group_ownership_fences`)
	require.NoError(t, err)
	independent := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "independent-cumulative", GroupRevision: "rev-I"}
	require.NoError(t, store.AppendSample(t.Context(), independent, t0, sampleAt("same-source", ownershipMember, 23, 7)))
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, sampleAt("affected-source", ownershipMember, 23, 7)), ErrOwnershipUnavailable)
}

func TestCutoverOwnershipConcurrentCommitBeforeAcceptanceRefusesStalePlan(t *testing.T) {
	store, db, path := newTestStore(t)
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
	before, err := store.ResolveCutoverOwnership(t.Context(), source, at(23))
	require.NoError(t, err)
	require.Empty(t, before.Affected)
	entered, release := make(chan struct{}), make(chan struct{})
	faults := &lostcommit.Faults{}
	hook := func(query string) error {
		if strings.Contains(query, "UPDATE wg_delivery_checkpoints SET updated_at") {
			close(entered)
			<-release
		}
		return nil
	}
	faults.Reject.Store(&hook)
	concurrentDB := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", faults)
	defer concurrentDB.Close()
	result := make(chan error, 1)
	go func() {
		result <- NewStore(concurrentDB).AppendSample(t.Context(), testKey, t0, sampleAt("stale-plan", ownershipMember, 23, 7))
	}()
	<-entered
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	close(release)
	require.ErrorIs(t, <-result, ErrOwnershipConflict)
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
}

func TestCutoverOwnershipNanosAndMalformedMember(t *testing.T) {
	store, db, _ := newTestStore(t)
	boundary := t0.Add(22500000001 * time.Nanosecond)
	seedOwnershipCutover(t, db, "f-nanos", testKey, otherKey, boundary)
	before := sampleAt("one-nanosecond-before", ownershipMember, 0, 6)
	before.ObservedAt = boundary.Add(-time.Nanosecond)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, before))
	exact := before
	exact.SampleID = "exact-nanos"
	exact.ObservedAt = boundary
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, exact), ErrOwnershipConflict)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, t0, exact))
	malformed := sampleAt("malformed-member", "temp", 23, 7)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, malformed), ErrOwnershipConflict)
	require.Equal(t, 2, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipRollbackUsesNewRevisionAtExactBoundary(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "forward", testKey, otherKey, t0.Add(22500*time.Millisecond))
	restored := testKey
	restored.GroupRevision = "restored-revision"
	seedOwnershipCutover(t, db, "rollback", otherKey, restored, at(40))
	_, err := db.Exec(`UPDATE write_group_cutovers SET kind='rollback',prior_cutover_id='forward',source_group_id=? WHERE id='rollback'`, testKey.GroupID)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_versions SET payload=json_set(payload,'$.recording_profile.kind','') WHERE group_id=? AND group_revision=?`, restored.GroupID, restored.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_ownership_fences SET profile_kind='' WHERE fence_revision='rollback' AND direction='started'`)
	require.NoError(t, err)
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
	owners, err := store.ResolveCutoverOwnership(t.Context(), source, at(40))
	require.NoError(t, err)
	require.Equal(t, map[GroupKey]bool{restored: true}, owners.Owners)
	exact := sampleAt("rollback-exact", ownershipMember, 40, 8)
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, exact), ErrOwnershipConflict)
	require.NoError(t, store.AppendSample(t.Context(), restored, t0, exact))
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipIndependentProfileTargetsAreNotDeduplicated(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	cumulative := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "cumulative-target", GroupRevision: "cumulative-1"}
	_, err := db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'connector-1','crev-1',?)`, cumulative.GroupID, cumulative.WorkspaceID, cumulative.GroupRevision, cumulative.GroupRevision, cumulative.GroupID, cumulative.GroupID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) SELECT ?,?,effective_at,json_set(payload,'$.id',?,'$.revision',?,'$.applied_revision',?,'$.destination.table_name',?,'$.recording_profile.kind','cumulative') FROM write_group_versions WHERE group_id=?`, cumulative.GroupID, cumulative.GroupRevision, cumulative.GroupID, cumulative.GroupRevision, cumulative.GroupRevision, cumulative.GroupID, otherKey.GroupID)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_cutovers SET detail=json_insert(detail,'$.started[#]',json_object('group_id',?,'revision',?))`, cumulative.GroupID, cumulative.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences SELECT fence_revision,workspace_id,direction,?,?,device_id,point_id,tag_id,'cumulative',connector_id,connector_revision,database_name,table_schema,? FROM write_group_ownership_fences WHERE direction='started'`, cumulative.GroupID, cumulative.GroupRevision, cumulative.GroupID)
	require.NoError(t, err)
	sample := sampleAt("same-observation", ownershipMember, 23, 7)
	require.NoError(t, store.AppendSample(t.Context(), otherKey, t0, sample))
	require.NoError(t, store.AppendSample(t.Context(), cumulative, t0, sample))
	require.Equal(t, 2, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipMalformedUnrelatedJSONDoesNotPoisonSource(t *testing.T) {
	for _, bad := range []string{`["scalar"]`, `{"moved":["bad-json"],"ended":["bad-json"],"started":["bad-json"]}`} {
		t.Run(bad, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
			_, err := db.Exec(`UPDATE write_group_cutovers SET detail=?`, bad)
			require.NoError(t, err)
			unrelated := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "unrelated-device", PointID: "unrelated-point", TagID: "unrelated-tag"}
			owners, err := store.ResolveCutoverOwnership(t.Context(), unrelated, at(23))
			require.NoError(t, err)
			require.Empty(t, owners.Affected)
			unrelatedKey := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "unrelated", GroupRevision: "rev-U"}
			require.NoError(t, store.AppendSample(t.Context(), unrelatedKey, t0, sampleAt("healthy", `["unrelated-device","unrelated-point","unrelated-tag"]`, 23, 4)))
			_, err = store.ResolveCutoverOwnership(t.Context(), OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}, at(23))
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
		})
	}
}

func TestCutoverOwnershipDuplicateRevalidationDoesNotInflateQuota(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	quota, err := store.WithQuota(QuotaConfig{GlobalMaxBytes: 1000000})
	require.NoError(t, err)
	quota.now = func() time.Time { return t0 }
	sample := sampleAt("before-boundary", ownershipMember, 22, 6)
	require.NoError(t, quota.AppendSample(t.Context(), testKey, t0, sample))
	before, present := quota.usage.get("", t0)
	require.True(t, present)
	for range 5 {
		require.NoError(t, quota.AppendSample(t.Context(), testKey, t0, sample))
	}
	after, present := quota.usage.get("", t0)
	require.True(t, present)
	require.Equal(t, before, after, "durable duplicate revalidation is not new accepted data")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipOrphanRootAmbiguityDoesNotInventOwner(t *testing.T) {
	for _, tc := range []struct {
		name      string
		canonical bool
		effective time.Time
		refused   bool
	}{
		{"same-time canonical remains ambiguous", true, t0.Add(22500 * time.Millisecond), true},
		{"different-time canonical proves independence", true, at(5), false},
		{"bare key cannot prove independence", false, time.Time{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
			independent := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "independent", GroupRevision: "rev-I"}
			if tc.canonical {
				_, err := db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'connector-1','crev-1',?)`, independent.GroupID, independent.WorkspaceID, independent.GroupRevision, independent.GroupRevision, independent.GroupID, independent.GroupID)
				require.NoError(t, err)
				_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) SELECT ?,?,?,json_set(payload,'$.id',?,'$.revision',?,'$.applied_revision',?) FROM write_group_versions WHERE group_id=?`, independent.GroupID, independent.GroupRevision, tc.effective, independent.GroupID, independent.GroupRevision, independent.GroupRevision, otherKey.GroupID)
				require.NoError(t, err)
			}
			_, err := db.Exec(`DELETE FROM write_group_ownership_fences; DELETE FROM write_group_cutovers`)
			require.NoError(t, err)
			err = store.AppendSample(t.Context(), independent, t0, sampleAt("new", ownershipMember, 23, 7))
			if tc.refused {
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
				require.Zero(t, countRows(t, db, "wg_delivery_samples"))
			} else {
				require.NoError(t, err)
				require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
			}
		})
	}
}

func TestCutoverOwnershipOrphanFencesCannotMasqueradeAsNoCutover(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	_, err := db.Exec(`DELETE FROM write_group_revision_ends; DELETE FROM write_group_cutovers`)
	require.NoError(t, err)
	_, err = store.ResolveCutoverOwnership(t.Context(), OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}, at(23))
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, sampleAt("orphan", ownershipMember, 23, 7)), ErrOwnershipUnavailable)
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipMalformedEmptyIDsStayScoped(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	_, err := db.Exec(`UPDATE write_group_cutovers SET detail=json_set(detail,'$.moved',json('[{"device_id":"","point_id":"","tag_id":""}]'))`)
	require.NoError(t, err)
	unrelated := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "unrelated", GroupRevision: "rev-U"}
	require.NoError(t, store.AppendSample(t.Context(), unrelated, t0, sampleAt("healthy", `["different-device","different-point","different-tag"]`, 23, 4)), "empty discovery params are not identities")
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, sampleAt("affected", ownershipMember, 23, 7)), ErrOwnershipUnavailable, "malformed moved membership is not trusted authority")
}

func TestCutoverOwnershipTransitiveAcceptedPredecessorCannotReappend(t *testing.T) {
	store, db, _ := newTestStore(t)
	original := sampleAt("accepted-A", ownershipMember, 45, 7)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, original))
	seedOwnershipCutover(t, db, "A-to-B", testKey, otherKey, t0.Add(22500*time.Millisecond))
	third := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: "third-owner", GroupRevision: "rev-C"}
	seedOwnershipCutover(t, db, "B-to-C", otherKey, third, at(40))
	require.ErrorIs(t, store.AppendSample(t.Context(), third, t0, original), ErrOwnershipConflict)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
	original.SampleID = "new-C"
	require.NoError(t, store.AppendSample(t.Context(), third, t0, original))
}

func TestCutoverOwnershipInterestIsScopedAndAdvisory(t *testing.T) {
	store, db, _ := newTestStore(t)
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
	interested, err := store.HasCutoverSource(t.Context(), source)
	require.NoError(t, err)
	require.False(t, interested, "unconfigured source is not activated")
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	interested, err = store.HasCutoverSource(t.Context(), source)
	require.NoError(t, err)
	require.True(t, interested, "an applied durable owner is interesting without any worker")
	foreign := source
	foreign.WorkspaceID = "foreign-workspace"
	interested, err = store.HasCutoverSource(t.Context(), foreign)
	require.NoError(t, err)
	require.False(t, interested)
	unrelated := source
	unrelated.TagID = "unrelated-tag"
	interested, err = store.HasCutoverSource(t.Context(), unrelated)
	require.NoError(t, err)
	require.False(t, interested)
	for _, status := range []string{"disabled", "deleted"} {
		_, err = db.Exec(`UPDATE write_groups SET status=? WHERE workspace_id=?`, status, testKey.WorkspaceID)
		require.NoError(t, err)
		interested, err = store.HasCutoverSource(t.Context(), source)
		require.NoError(t, err)
		require.False(t, interested, "interest cannot reverse lifecycle %s", status)
	}
	_, err = db.Exec(`UPDATE write_groups SET status='ready',applied_revision='another-revision' WHERE id=?`, otherKey.GroupID)
	require.NoError(t, err)
	interested, err = store.HasCutoverSource(t.Context(), source)
	require.NoError(t, err)
	require.False(t, interested, "historical cutover does not activate a superseded revision")
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
}

func TestCutoverOwnershipInterestReadFailureIsNotDisinterest(t *testing.T) {
	store, db, _ := newTestStore(t)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, t0.Add(22500*time.Millisecond))
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
	_, err := db.Exec(`DELETE FROM write_groups WHERE id=?`, otherKey.GroupID)
	require.NoError(t, err)
	interested, err := store.HasCutoverSource(t.Context(), source)
	require.False(t, interested)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	_, err = db.Exec(`DROP TABLE write_group_ownership_fences`)
	require.NoError(t, err)
	interested, err = store.HasCutoverSource(t.Context(), source)
	require.False(t, interested)
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
}

func TestCutoverOwnershipRetainedSourceNotListedAsMoved(t *testing.T) {
	store, db, _ := newTestStore(t)
	boundary := t0.Add(22500 * time.Millisecond)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, boundary)
	retained := GroupKey{WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: "remaining-revision"}
	member := `{"device_id":"retained-device","point_id":"retained-point","tag_id":"retained-tag"}`
	_, err := db.Exec(`UPDATE write_group_versions SET payload=json_insert(payload,'$.members[#]',json(?)) WHERE group_id=? AND group_revision=?`, member, testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) SELECT ?,?,?,json_set(payload,'$.revision',?,'$.applied_revision',?,'$.members',json_array(json(?))) FROM write_group_versions WHERE group_id=? AND group_revision=?`, retained.GroupID, retained.GroupRevision, boundary, retained.GroupRevision, retained.GroupRevision, member, testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_cutovers SET detail=json_insert(detail,'$.started[#]',json_object('group_id',?,'revision',?))`, retained.GroupID, retained.GroupRevision)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences SELECT fence_revision,workspace_id,direction,group_id,group_revision,'retained-device','retained-point','retained-tag',profile_kind,connector_id,connector_revision,database_name,table_schema,table_name FROM write_group_ownership_fences WHERE direction='ended'`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_ownership_fences SELECT fence_revision,workspace_id,'started',group_id,?,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name FROM write_group_ownership_fences WHERE direction='ended' AND tag_id='retained-tag'`, retained.GroupRevision)
	require.NoError(t, err)
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "retained-device", PointID: "retained-point", TagID: "retained-tag"}
	owners, err := store.ResolveCutoverOwnership(t.Context(), source, boundary)
	require.NoError(t, err)
	require.Equal(t, map[GroupKey]bool{retained: true}, owners.Owners)
	sample := sampleAt("retained-sample", `["retained-device","retained-point","retained-tag"]`, 23, 6)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, sample), ErrOwnershipConflict)
	require.NoError(t, store.AppendSample(t.Context(), retained, t0, sample))
	_, err = db.Exec(`DELETE FROM write_group_ownership_fences WHERE tag_id='retained-tag'`)
	require.NoError(t, err)
	_, err = store.ResolveCutoverOwnership(t.Context(), source, boundary)
	require.ErrorIs(t, err, ErrOwnershipUnavailable, "immutable source membership discovers missing retained scopes without Moved")
}

func TestCutoverOwnershipInterestIncludesFutureEndedSource(t *testing.T) {
	store, db, _ := newTestStore(t)
	boundary := t0.Add(22500 * time.Millisecond)
	seedOwnershipCutover(t, db, "f-2", testKey, otherKey, boundary)
	_, err := db.Exec("UPDATE write_groups SET status='disabled' WHERE id=?", otherKey.GroupID)
	require.NoError(t, err)
	source := OwnershipSource{WorkspaceID: testKey.WorkspaceID, DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"}
	interested, err := store.HasCutoverSource(t.Context(), source)
	require.NoError(t, err)
	require.True(t, interested, "a clock-free hint must not hide the active source owning pre-boundary observations")
	before := sampleAt("before", ownershipMember, 0, 6)
	before.ObservedAt = boundary.Add(-time.Nanosecond)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, before))
	require.ErrorIs(t, store.AppendSample(t.Context(), otherKey, t0, before), ErrOwnershipConflict)
	exact := sampleAt("exact", ownershipMember, 0, 7)
	exact.ObservedAt = boundary
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, exact), ErrOwnershipConflict)
}
