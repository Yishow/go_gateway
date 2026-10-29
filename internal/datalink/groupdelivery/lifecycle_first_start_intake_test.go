package groupdelivery

import (
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

func seedFirstStartIntake(t *testing.T, db *sql.DB, boundary time.Time) {
	t.Helper()
	seedLifecycleIntake(t, db, testKey)
	_, err := db.Exec(`UPDATE write_group_versions SET effective_at=? WHERE group_id=? AND group_revision=?`, boundary, testKey.GroupID, testKey.GroupRevision)
	require.NoError(t, err)
}

// First start has no predecessor version, end or fence. The journal itself must
// enforce the exact immutable boundary, including after a delayed reconcile or
// process restart; neither receive time nor the append clock chooses ownership.
func TestCanonicalLifecycleFirstStartObservedBoundary(t *testing.T) {
	boundary := at(12).Add(500 * time.Millisecond)
	for _, phase := range []string{"before-boundary", "delayed-reconcile", "restart"} {
		t.Run(phase, func(t *testing.T) {
			base, db, path := newTestStore(t)
			seedFirstStartIntake(t, db, boundary)
			now := at(10)
			if phase == "delayed-reconcile" {
				now = at(18)
			}
			if phase == "restart" {
				require.NoError(t, db.Close())
				db = openStoreDB(t, path)
				base = NewStore(db)
				now = at(3600)
			}
			store, err := base.WithCanonicalIntake().WithClock(func() time.Time { return now }).WithQuota(QuotaConfig{GlobalMaxBytes: 1 << 20})
			require.NoError(t, err)
			for _, tc := range []struct {
				name string
				when time.Time
				want error
			}{
				{"pre-B", boundary.Add(-time.Nanosecond), ErrOwnershipConflict},
				{"equal-B", boundary, nil},
				{"post-B", boundary.Add(time.Nanosecond), nil},
				{"equal-B-offset-zone", boundary.In(time.FixedZone("UTC+8", 8*60*60)), nil},
			} {
				t.Run(tc.name, func(t *testing.T) {
					sample := sampleAt(tc.name, ownershipMember, 0, 107)
					sample.ObservedAt, sample.ReceivedAt = tc.when, at(18)
					before := associationJournal(t, store)
					err := store.AppendSample(t.Context(), testKey, at(10), sample)
					if tc.want != nil {
						require.ErrorIs(t, err, tc.want)
						require.Equal(t, before, associationJournal(t, store), "pre-B intake cannot create or rewrite acceptance")
						return
					}
					require.NoError(t, err, "the bucket may begin before B; the observation controls acceptance")
				})
			}
			require.Equal(t, 3, countRows(t, db, "wg_delivery_samples"))
			require.Equal(t, 1, countRows(t, db, "write_group_versions"))
			for _, table := range []string{"write_group_cutovers", "write_group_revision_ends", "write_group_ownership_fences", "wg_delivery_sample_handoffs", "wg_delivery_outbox", "wg_delivery_receipts"} {
				require.Zero(t, countRows(t, db, table), "first intake must not fabricate historical ownership: %s", table)
			}
		})
	}
}

func TestCanonicalLifecycleFirstStartUsesExactImmutableVersion(t *testing.T) {
	base, db, _ := newTestStore(t)
	boundary := at(12).Add(500 * time.Millisecond)
	seedFirstStartIntake(t, db, boundary)
	// A saved draft/current pointer is not authority for the requested immutable
	// revision. No substitute latest-version lookup may erase its start bound.
	_, err := db.Exec(`UPDATE write_groups SET revision='new-draft',applied_revision='another-applied' WHERE id=?`, testKey.GroupID)
	require.NoError(t, err)
	store := base.WithCanonicalIntake()
	pre := sampleAt("pre", ownershipMember, 12, 100)
	require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(10), pre), ErrOwnershipConflict)
	post := sampleAt("post", ownershipMember, 13, 107)
	require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), post))
	missing := testKey
	missing.GroupRevision = "another-applied"
	require.ErrorIs(t, store.AppendSample(t.Context(), missing, at(10), post), ErrOwnershipUnavailable)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_samples"))
	// Opt-in must leave the standalone Store's original behavior untouched even
	// when the same database happens to contain canonical metadata.
	require.NoError(t, base.AppendSample(t.Context(), testKey, at(10), pre))
	require.Equal(t, 2, countRows(t, db, "wg_delivery_samples"))
}

func TestCanonicalLifecycleFirstStartCorruptAuthorityFailsClosed(t *testing.T) {
	for name, mutation := range map[string]string{
		"missing version":           `DELETE FROM write_group_versions`,
		"malformed version":         `UPDATE write_group_versions SET payload='not-json'`,
		"wrong workspace":           `UPDATE write_group_versions SET payload=json_set(payload,'$.workspace_id','foreign')`,
		"wrong revision":            `UPDATE write_group_versions SET payload=json_set(payload,'$.revision','foreign')`,
		"empty effective time":      `UPDATE write_group_versions SET effective_at=''`,
		"malformed effective time":  `UPDATE write_group_versions SET effective_at='private-invalid-time'`,
		"zero effective time":       `UPDATE write_group_versions SET effective_at='0001-01-01T00:00:00Z'`,
		"unreadable effective time": `ALTER TABLE write_group_versions RENAME COLUMN effective_at TO unavailable_effective_at`,
	} {
		t.Run(name, func(t *testing.T) {
			base, db, _ := newTestStore(t)
			seedFirstStartIntake(t, db, at(12))
			_, err := db.Exec(mutation)
			require.NoError(t, err)
			err = base.WithCanonicalIntake().AppendSample(t.Context(), testKey, at(10), sampleAt("refused", ownershipMember, 13, 107))
			require.ErrorIs(t, err, ErrOwnershipUnavailable)
			require.Equal(t, "recording-ownership-unavailable", err.Error(), "corrupt storage details remain private")
			require.Zero(t, countRows(t, db, "wg_delivery_samples"))
		})
	}
}

func TestCanonicalLifecycleFirstStartAcceptedRetryKeepsOriginalAuthority(t *testing.T) {
	for _, historical := range []bool{false, true} {
		name := "accepted-at-first-boundary"
		if historical {
			name = "already-accepted-before-guard"
		}
		t.Run(name, func(t *testing.T) {
			base, db, path := newTestStore(t)
			seedFirstStartIntake(t, db, at(12))
			store := base.WithCanonicalIntake()
			accepted := sampleAt("accepted", ownershipMember, 12, 107)
			if historical {
				// Existing durable acceptance is not reauthorized by a new guard.
				accepted.ObservedAt = at(11)
				require.NoError(t, base.AppendSample(t.Context(), testKey, at(10), accepted))
			} else {
				require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), accepted))
			}
			before := associationJournal(t, store)
			_, err := db.Exec(`UPDATE write_groups SET status='disabled'`)
			require.NoError(t, err)
			require.NoError(t, db.Close())
			db = openStoreDB(t, path)
			store = NewStore(db).WithCanonicalIntake().WithClock(func() time.Time { return at(3600) })
			require.NoError(t, store.AppendSample(t.Context(), testKey, at(10), accepted))
			changed := accepted
			changed.Value = measurement.NewInt64(108)
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(10), changed), ErrSampleConflict)
			changed = accepted
			changed.SourceRevision = "foreign-source"
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(10), changed), ErrSampleConflict)
			changed = accepted
			changed.ObservedAt = at(13)
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(10), changed), ErrSampleConflict)
			foreign := testKey
			foreign.WorkspaceID = "foreign-workspace"
			var proofError *counterproof.CounterBindingProofError
			require.ErrorAs(t, store.AppendSample(t.Context(), foreign, at(10), accepted), &proofError)
			fresh := sampleAt("fresh", ownershipMember, 13, 110)
			require.ErrorIs(t, store.AppendSample(t.Context(), testKey, at(10), fresh), ErrIntakeStopped)
			require.Equal(t, before, associationJournal(t, store))
		})
	}
}
