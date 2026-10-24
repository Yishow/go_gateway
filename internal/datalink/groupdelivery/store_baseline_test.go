package groupdelivery

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

// baselineAt is the serialized baseline of one member: reading value taken at
// windowEnd-1s of a window that ends at windowEnd.
func baselineAt(t *testing.T, member string, windowEnd time.Time, value uint64) CounterBaseline {
	t.Helper()
	binding := aggregation.CounterBinding{
		WorkspaceID: "workspace-A", GroupID: "group-G", ProfileRevision: "rev-1", MemberKey: member,
		MeasurementID: "meas-" + member, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", Kind: aggregation.CounterKindMonotonic, Unit: "kWh",
	}
	payload, err := aggregation.EncodeCounterCheckpoint(aggregation.CounterCheckpoint{
		Binding: binding, WindowEnd: windowEnd,
		Reading: aggregation.CounterReading{
			SampleID: "s-" + member + "-" + windowEnd.Format("1504"), ObservedAt: windowEnd.Add(-time.Second), Value: measurement.NewUint64(value),
		},
	})
	require.NoError(t, err)
	return CounterBaseline{MemberKey: member, Payload: payload}
}

func closureWith(t *testing.T, end time.Time, counter uint64, baselines ...CounterBaseline) Closure {
	t.Helper()
	destination := testDestination
	destination.TableSchema = ""
	return Closure{
		Destination: destination, NextClose: end, Baselines: baselines,
		Buckets: []ClosedBucket{rowBucket("", end.Add(-10*time.Second), counter)},
	}
}

func baselineRow(t *testing.T, db *sql.DB, member string) (payload string) {
	t.Helper()
	err := db.QueryRowContext(t.Context(),
		`SELECT payload FROM wg_delivery_counter_baselines WHERE group_id = ? AND group_revision = ? AND member_key = ?`,
		testKey.GroupID, testKey.GroupRevision, member).Scan(&payload)
	require.NoError(t, err)
	return payload
}

func TestCommitClosureStoresBaselinesWithTheClosure(t *testing.T) {
	store, db, _ := newTestStore(t)
	end := t0.Add(10 * time.Second)
	alpha, beta := baselineAt(t, "alpha", end, 100), baselineAt(t, "beta", end, 200)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, alpha, beta)))

	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))
	require.Equal(t, 2, countRows(t, db, "wg_delivery_counter_baselines"))
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Equal(t, end, restored.NextClose)
	require.ElementsMatch(t, []CounterBaseline{alpha, beta}, restored.Baselines, "baselines come back byte for byte")

	other, err := store.Restore(t.Context(), otherKey)
	require.NoError(t, err)
	require.Empty(t, other.Baselines, "a baseline belongs to one group revision")
}

func TestCommitClosureBaselineIsAtomicWithTheRest(t *testing.T) {
	assertNothingPersisted := func(t *testing.T, store *Store, db *sql.DB) {
		t.Helper()
		for _, table := range []string{"wg_delivery_outbox", "wg_delivery_buckets", "wg_delivery_checkpoints", "wg_delivery_counter_baselines"} {
			require.Zero(t, countRows(t, db, table), table)
		}
		restored, err := store.Restore(t.Context(), testKey)
		require.NoError(t, err)
		require.True(t, restored.NextClose.IsZero())
	}
	end := t0.Add(10 * time.Second)

	t.Run("a failure while storing the baseline rolls back the records and the checkpoint", func(t *testing.T) {
		store, db, _ := newTestStore(t)
		require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("s1", "alpha", 3, 5)))
		failOn(t, db, "wg_delivery_counter_baselines")
		err := store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 100)))
		require.Error(t, err)
		require.Zero(t, countRows(t, db, "wg_delivery_outbox"), "no outbox intent without its baseline")
		require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"), "the closure cursor does not move")
		require.Zero(t, countRows(t, db, "wg_delivery_counter_baselines"))
		restored, err := store.Restore(t.Context(), testKey)
		require.NoError(t, err)
		require.Len(t, restored.Samples, 1, "the journaled sample is still unconsumed")
	})

	t.Run("a baseline that cannot be decoded is refused before anything is written", func(t *testing.T) {
		store, db, _ := newTestStore(t)
		for name, payload := range map[string][]byte{
			"not json":        []byte("{"),
			"trailing bytes":  append(baselineAt(t, "alpha", end, 1).Payload, 'x'),
			"unknown version": []byte(`{"version":9}`),
		} {
			err := store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, CounterBaseline{MemberKey: "alpha", Payload: payload}))
			require.ErrorIs(t, err, ErrInvalidClosure, name)
			assertNothingPersisted(t, store, db)
		}
		err := store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, CounterBaseline{MemberKey: "", Payload: baselineAt(t, "alpha", end, 1).Payload}))
		require.ErrorIs(t, err, ErrInvalidClosure, "a baseline needs its member")
		err = store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 1), baselineAt(t, "alpha", end, 2)))
		require.ErrorIs(t, err, ErrInvalidClosure, "one baseline per member in a closure")
		assertNothingPersisted(t, store, db)
	})

	t.Run("a baseline of another member than its payload names is refused", func(t *testing.T) {
		store, db, _ := newTestStore(t)
		mismatched := baselineAt(t, "alpha", end, 1)
		mismatched.MemberKey = "beta"
		require.ErrorIs(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, mismatched)), ErrInvalidClosure)
		assertNothingPersisted(t, store, db)
	})
}

func TestCommitClosureBaselineMovesForwardOnce(t *testing.T) {
	store, db, _ := newTestStore(t)
	first, second := t0.Add(10*time.Second), t0.Add(20*time.Second)
	one, two := baselineAt(t, "alpha", first, 100), baselineAt(t, "alpha", second, 130)

	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, first, 1, one)))
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, first, 1, one)), "repeating a committed closure is a no-op")
	payload := baselineRow(t, db, "alpha")
	require.Equal(t, string(one.Payload), payload)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"))

	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, second, 2, two)))
	payload = baselineRow(t, db, "alpha")
	require.Equal(t, string(two.Payload), payload, "the newer window replaces the baseline")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_counter_baselines"), "one row per member")

	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, first, 1, one)), "an older closure replayed late")
	payload = baselineRow(t, db, "alpha")
	require.Equal(t, string(two.Payload), payload, "a baseline never moves back")
	require.Equal(t, 2, countRows(t, db, "wg_delivery_outbox"))
}

func TestCommitClosureRefusesADifferentBaselineForTheSameWindow(t *testing.T) {
	store, db, _ := newTestStore(t)
	end := t0.Add(10 * time.Second)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 100))))

	err := store.CommitClosure(t.Context(), testKey, closureWith(t, end.Add(10*time.Second), 2, baselineAt(t, "alpha", end, 101)))
	require.ErrorIs(t, err, ErrBaselineConflict)
	require.Equal(t, 1, countRows(t, db, "wg_delivery_outbox"), "the conflicting closure left nothing behind")
	restored, restoreErr := store.Restore(t.Context(), testKey)
	require.NoError(t, restoreErr)
	require.Equal(t, end, restored.NextClose)
}

func TestRestoreRefusesABaselineThatChangedAfterCommit(t *testing.T) {
	store, db, _ := newTestStore(t)
	end := t0.Add(10 * time.Second)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 100))))
	_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET payload = replace(payload, '"100"', '"900"')`)
	require.NoError(t, err)

	_, err = store.Restore(t.Context(), testKey)
	require.ErrorIs(t, err, ErrBaselineCorrupt, "damaged durable state is reported, never guessed")
}

func TestReclaimKeepsCounterBaselines(t *testing.T) {
	store, db, _ := newTestStore(t)
	end := t0.Add(10 * time.Second)
	require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sampleAt("s1", "alpha", 3, 5)))
	require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 100))))
	removed, err := store.Reclaim(t.Context(), 0)
	require.NoError(t, err)
	require.NotZero(t, removed, "the consumed sample was reclaimed")
	require.Equal(t, 1, countRows(t, db, "wg_delivery_counter_baselines"),
		"a baseline outlives the journal rows it was derived from")
	restored, err := store.Restore(t.Context(), testKey)
	require.NoError(t, err)
	require.Len(t, restored.Baselines, 1)
}

// foreignBaseline is a perfectly valid checkpoint of another scope: only the
// workspace, group or profile revision it names differs from the closure's key.
func foreignBaseline(t *testing.T, member string, end time.Time, mutate func(*aggregation.CounterBinding)) CounterBaseline {
	t.Helper()
	binding := aggregation.CounterBinding{
		WorkspaceID: "workspace-A", GroupID: "group-G", ProfileRevision: "rev-1", MemberKey: member,
		MeasurementID: "meas-" + member, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", Kind: aggregation.CounterKindMonotonic, Unit: "kWh",
	}
	mutate(&binding)
	payload, err := aggregation.EncodeCounterCheckpoint(aggregation.CounterCheckpoint{
		Binding: binding, WindowEnd: end,
		Reading: aggregation.CounterReading{SampleID: "s", ObservedAt: end.Add(-time.Second), Value: measurement.NewUint64(1)},
	})
	require.NoError(t, err)
	return CounterBaseline{MemberKey: member, Payload: payload}
}

// TestCommitClosureRefusesABaselineOfAnotherScope pins that the workspace, group
// and profile revision of a checkpoint are the fixed scope of the closure's own
// key: a valid checkpoint of another scope is foreign state, not a binding
// that merely changed, and nothing of the closure is written.
func TestCommitClosureRefusesABaselineOfAnotherScope(t *testing.T) {
	end := t0.Add(10 * time.Second)
	for name, mutate := range map[string]func(*aggregation.CounterBinding){
		"another workspace":        func(b *aggregation.CounterBinding) { b.WorkspaceID = "workspace-B" },
		"another group":            func(b *aggregation.CounterBinding) { b.GroupID = "group-other" },
		"another profile revision": func(b *aggregation.CounterBinding) { b.ProfileRevision = "rev-0" },
	} {
		t.Run(name, func(t *testing.T) {
			store, db, _ := newTestStore(t)
			err := store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, foreignBaseline(t, "alpha", end, mutate)))
			require.ErrorIs(t, err, ErrInvalidClosure)
			for _, table := range []string{"wg_delivery_outbox", "wg_delivery_buckets", "wg_delivery_checkpoints", "wg_delivery_counter_baselines"} {
				require.Zero(t, countRows(t, db, table), table)
			}
		})
	}

	t.Run("a revised measurement is still a legitimate baseline", func(t *testing.T) {
		for name, mutate := range map[string]func(*aggregation.CounterBinding){
			"epoch":          func(b *aggregation.CounterBinding) { b.SeriesEpoch = "epoch-2" },
			"definition":     func(b *aggregation.CounterBinding) { b.DefinitionRevision = "def-2" },
			"source binding": func(b *aggregation.CounterBinding) { b.SourceBindingRevision = "bind-2" },
			"unit":           func(b *aggregation.CounterBinding) { b.Unit = "MWh" },
		} {
			store, db, _ := newTestStore(t)
			require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, foreignBaseline(t, "alpha", end, mutate))), name)
			require.Equal(t, 1, countRows(t, db, "wg_delivery_counter_baselines"), name)
		}
	})
}

func TestCommitClosureRefusesABaselineBeyondTheClosedCursor(t *testing.T) {
	store, db, _ := newTestStore(t)
	end := t0.Add(10 * time.Second)
	err := store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end.Add(10*time.Second), 5)))
	require.ErrorIs(t, err, ErrInvalidClosure, "a baseline cannot claim a window the closure has not closed")
	require.Zero(t, countRows(t, db, "wg_delivery_counter_baselines"))
	require.Zero(t, countRows(t, db, "wg_delivery_checkpoints"))
}

func TestRestoreRefusesBaselinesOutsideTheirScopeOrCursor(t *testing.T) {
	end := t0.Add(10 * time.Second)
	seed := func(t *testing.T) (*Store, *sql.DB) {
		store, db, _ := newTestStore(t)
		require.NoError(t, store.CommitClosure(t.Context(), testKey, closureWith(t, end, 1, baselineAt(t, "alpha", end, 5))))
		return store, db
	}
	digest := func(t *testing.T, db *sql.DB) {
		t.Helper()
		var payload string
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&payload))
		sum := sha256.Sum256([]byte(payload))
		_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET payload_digest = ?`, hex.EncodeToString(sum[:]))
		require.NoError(t, err)
	}

	t.Run("a row filed under another workspace", func(t *testing.T) {
		store, db := seed(t)
		_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET workspace_id = 'workspace-B'`)
		require.NoError(t, err)
		_, err = store.Restore(t.Context(), testKey)
		require.ErrorIs(t, err, ErrBaselineCorrupt)
	})
	t.Run("a baseline ahead of the closed cursor", func(t *testing.T) {
		store, db := seed(t)
		_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET window_end = ?`, timestamp(end.Add(time.Hour)))
		require.NoError(t, err)
		_, err = store.Restore(t.Context(), testKey)
		require.ErrorIs(t, err, ErrBaselineCorrupt, "a baseline past the checkpoint is an inconsistent state, not an unfinished window")
	})
	t.Run("an intact state restores", func(t *testing.T) {
		store, db := seed(t)
		digest(t, db)
		restored, err := store.Restore(t.Context(), testKey)
		require.NoError(t, err)
		require.Len(t, restored.Baselines, 1)
	})
}
