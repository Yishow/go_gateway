package groupdelivery

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func seedLifecycleIntake(t *testing.T, db *sql.DB, key GroupKey) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO write_groups(id,workspace_id,revision,applied_revision,name,status,destination_connector_id,destination_connector_revision,destination_table_name) VALUES(?,?,?,?,?,'ready','connector-1','C1','D1')`, key.GroupID, key.WorkspaceID, key.GroupRevision, key.GroupRevision, key.GroupID)
	require.NoError(t, err)
	payload, err := json.Marshal(map[string]any{
		"id": key.GroupID, "workspace_id": key.WorkspaceID, "revision": key.GroupRevision, "applied_revision": key.GroupRevision,
		"members":     []map[string]string{{"device_id": "device-1", "point_id": "point-1", "tag_id": "tag-1"}},
		"destination": map[string]string{"connector_id": "connector-1", "connector_revision": "C1", "table_name": "D1"},
	})
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload) VALUES(?,?,?,?)`, key.GroupID, key.GroupRevision, t0, string(payload))
	require.NoError(t, err)
}

func TestCanonicalLifecycleNewIntakeFailsClosed(t *testing.T) {
	for name, mutation := range map[string]string{
		"disabled":                    `UPDATE write_groups SET status='disabled'`,
		"deleted":                     `UPDATE write_groups SET status='deleted'`,
		"missing group":               `DELETE FROM write_groups`,
		"foreign workspace":           `UPDATE write_groups SET workspace_id='other-workspace'`,
		"unknown status":              `UPDATE write_groups SET status='secret-driver-text'`,
		"no applied revision":         `UPDATE write_groups SET applied_revision=''`,
		"unknown revision":            `DELETE FROM write_group_versions`,
		"malformed immutable version": `UPDATE write_group_versions SET payload='not-json'`,
		"wrong immutable identity":    `UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','other')`,
		"unreadable authority":        `ALTER TABLE write_groups RENAME COLUMN status TO unavailable_status`,
	} {
		t.Run(name, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			seedLifecycleIntake(t, db, testKey)
			store := base.WithCanonicalIntake()
			prior := sampleAt("accepted", "temp", 3, 30)
			require.NoError(t, store.AppendSample(t.Context(), testKey, t0, prior))
			before := associationJournal(t, store)
			_, err := db.Exec(mutation)
			require.NoError(t, err)
			require.NoError(t, store.AppendSample(t.Context(), testKey, t0, prior), "durable exact retry does not reauthorize intake")
			err = store.AppendSample(t.Context(), testKey, t0, sampleAt("new", "temp", 4, 40))
			var typed *OwnershipError
			require.ErrorAs(t, err, &typed)
			if name == "disabled" || name == "deleted" {
				require.Equal(t, "recording-intake-stopped", err.Error())
			} else {
				require.ErrorIs(t, err, ErrOwnershipUnavailable)
			}
			require.Equal(t, before, associationJournal(t, store))
			changed := prior
			changed.Value = measurement.NewInt64(99)
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, changed), ErrSampleConflict)
		})
	}
}

func TestCanonicalLifecycleActiveAndStandaloneCompatibility(t *testing.T) {
	for _, status := range []string{"draft", "ready", "running"} {
		t.Run(status, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			seedLifecycleIntake(t, db, testKey)
			_, err := db.Exec(`UPDATE write_groups SET revision='new-draft',applied_revision='later-applied',status=?`, status)
			require.NoError(t, err)
			store, err := base.WithCanonicalIntake().WithClock(func() time.Time { return at(5) }).WithQuota(QuotaConfig{GlobalMaxBytes: 1 << 20})
			require.NoError(t, err)
			require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("historical", "temp", 3, 30)), "known historical revisions retain existing cutover authority")
			invented := testKey
			invented.GroupRevision = "invented"
			require.ErrorIs(t, store.AppendSample(t.Context(), invented, t0, sampleAt("invented", "temp", 3, 30)), ErrOwnershipUnavailable)
			_, err = db.Exec(`UPDATE write_groups SET status='disabled'`)
			require.NoError(t, err)
			require.Error(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("new", "temp", 4, 40)), "both clones retain the canonical guard")
			require.NoError(t, base.AppendSample(t.Context(), otherKey, t0, sampleAt("standalone", "temp", 3, 30)), "opting in must not mutate the original standalone ledger")
		})
	}
}

func TestCanonicalLifecycleAuthorityReadFailureIsTyped(t *testing.T) {
	_, db, path := newTestStore(t)
	seedLifecycleIntake(t, db, testKey)
	faults := &lostcommit.Faults{}
	faultDB := lostcommit.OpenFaults(path+"?_pragma=busy_timeout(5000)", faults)
	t.Cleanup(func() { require.NoError(t, faultDB.Close()) })
	reject := func(query string) error {
		if strings.Contains(query, "FROM write_groups") {
			return errors.New("private DSN must not escape")
		}
		return nil
	}
	faults.Reject.Store(&reject)
	err := NewStore(faultDB).WithCanonicalIntake().AppendSample(t.Context(), testKey, t0, sampleAt("refused", "temp", 3, 30))
	require.ErrorIs(t, err, ErrOwnershipUnavailable)
	require.Equal(t, "recording-ownership-unavailable", err.Error())
	require.Zero(t, countRows(t, db, "wg_delivery_samples"))
}

func TestCanonicalLifecyclePreservesCounterProofAndQuota(t *testing.T) {
	base, db, _ := newTestStore(t)
	seedLifecycleIntake(t, db, testKey)
	store, err := base.WithCanonicalIntake().WithClock(func() time.Time { return at(5) }).WithQuota(QuotaConfig{GlobalMaxBytes: 1 << 20})
	require.NoError(t, err)
	_, sample, payload := journalProofFixture(t)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	usage, ok := store.usage.get("", at(5))
	require.True(t, ok)
	require.Positive(t, usage)
	_, err = db.Exec(`UPDATE write_groups SET status='disabled'`)
	require.NoError(t, err)
	before := associationJournal(t, store)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
	newSample := sample
	newSample.SampleID = "new-counter"
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, t0, newSample), ErrIntakeStopped)
	newSample.CounterProof.Digest = strings.Repeat("b", 64)
	var proofError *counterproof.CounterBindingProofError
	require.ErrorAs(t, store.AppendSample(t.Context(), testKey, t0, newSample), &proofError, "stop must not bypass immutable proof validation")
	_, err = db.Exec(`UPDATE wg_runtime_versions SET payload=?`, strings.Replace(string(payload), "kWh", "MWh", 1))
	require.NoError(t, err)
	require.ErrorAs(t, store.AppendSample(t.Context(), testKey, t0, sample), &proofError, "an accepted retry still verifies its immutable proof")
	require.Equal(t, before, associationJournal(t, store))
	after, ok := store.usage.get("", at(5))
	require.True(t, ok)
	require.Equal(t, usage, after, "retries and refusals do not charge accepted bytes")
}
