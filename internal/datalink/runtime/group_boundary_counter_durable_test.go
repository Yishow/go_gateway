package runtime

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// outboxCell reads one committed record of the durable outbox by source and
// window, decoded from its frozen payload.
func outboxCell(t *testing.T, db *sql.DB, bucket time.Duration) map[string]any {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `SELECT payload FROM wg_delivery_outbox ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var payload string
		require.NoError(t, rows.Scan(&payload))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		cells := map[string]any{}
		for _, cell := range row.Cells {
			cells[cell.Column] = cell.Value
		}
		if cells["tag_id"] == "tag-uint" && cells["bucket_start"] == bucketText(bucket) {
			return cells
		}
	}
	require.NoError(t, rows.Err())
	require.FailNow(t, "no committed record", "tag-uint at %s", bucketText(bucket))
	return nil
}

func baselineCount(t *testing.T, db *sql.DB) int {
	return countTable(t, db, "wg_delivery_counter_baselines")
}

// TestCumulativeProfileBaselineSurvivesRestartThroughTheLedger closes a window,
// restarts the boundary from the same local database and closes the next one:
// the exact difference is computed against the baseline committed with the
// first closure, not against a fresh start.
func TestCumulativeProfileBaselineSurvivesRestartThroughTheLedger(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	db := openDurableDB(t, path)
	f := newCounterFixture(t)
	f.config.Ledger = longLedger(db)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(1, "uint-0", at(counterWindow-time.Minute), uint64(9007199254740992)))
	f.clock = at(counterWindow)
	require.NoError(t, boundary.Tick(t.Context(), f.clock))
	require.Equal(t, 1, baselineCount(t, db), "the baseline is committed with the closure that produced the record")
	require.Equal(t, 3, countTable(t, db, "wg_delivery_outbox"))

	require.NoError(t, db.Close()) // the process dies after the closure
	reopened := openDurableDB(t, path)
	again := newCounterFixture(t)
	again.config.Ledger = longLedger(reopened)
	again.clock = at(2 * counterWindow)
	recovered := again.build(t)
	again.accept(t, recovered, again.envelope(1, "uint-1", at(2*counterWindow-time.Minute), uint64(9007199254740993)))
	require.NoError(t, recovered.Tick(t.Context(), again.clock))

	first := outboxCell(t, reopened, 0)
	require.Equal(t, "baseline", first["status"])
	second := outboxCell(t, reopened, counterWindow)
	require.Equal(t, "ok", second["status"])
	require.Equal(t, "9007199254740993", second["value_decimal"])
	require.Equal(t, "1", second["delta_value"], "exact against the restored baseline")
	require.Equal(t, 1, baselineCount(t, reopened), "one baseline per member, advanced in place")
}

// ledgerWithoutBaselines is a ledger that journals and closes but does not say
// it persists counter baselines.
type ledgerWithoutBaselines struct{ inner GroupLedger }

func (l ledgerWithoutBaselines) AppendSample(ctx context.Context, start time.Time, sample snapshot.Sample) error {
	return l.inner.AppendSample(ctx, start, sample)
}
func (l ledgerWithoutBaselines) Restore(ctx context.Context) (groupdelivery.Restored, error) {
	return l.inner.Restore(ctx)
}
func (l ledgerWithoutBaselines) CommitClosure(ctx context.Context, closure groupdelivery.Closure) error {
	return l.inner.CommitClosure(ctx, closure)
}

func TestCumulativeProfileNeedsALedgerThatPersistsBaselines(t *testing.T) {
	f := newCounterFixture(t)
	f.config.Ledger = ledgerWithoutBaselines{inner: longLedger(openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db")))}
	freezeDurableFixture(t, f)
	_, err := NewGroupBoundary(f.config)
	var blocked *GroupBoundaryError
	require.ErrorAs(t, err, &blocked)
	require.Equal(t, "counter-baseline-not-durable", blocked.Code, "a closure that cannot carry the baseline would restart from nothing")

	f.config.Ledger = longLedger(openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db")))
	freezeDurableFixture(t, f)
	_, err = NewGroupBoundary(f.config)
	require.NoError(t, err)
}

// seedBaseline commits one closure whose baseline payload is rewritten after the
// fact, the way damaged or foreign durable state would look.
func seedBaseline(t *testing.T, db *sql.DB, rewrite func(payload string) string) {
	t.Helper()
	f := newCounterFixture(t)
	f.config.Ledger = longLedger(db)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(1, "uint-0", at(counterWindow-time.Minute), uint64(100)))
	f.clock = at(counterWindow)
	require.NoError(t, boundary.Tick(t.Context(), f.clock))
	var payload string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines`).Scan(&payload))
	rewritten := rewrite(payload)
	sum := sha256.Sum256([]byte(rewritten))
	_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET payload = ?, payload_digest = ?`, rewritten, hex.EncodeToString(sum[:]))
	require.NoError(t, err)
}

func TestCumulativeProfileTellsDamagedBaselinesFromRevisedBindings(t *testing.T) {
	activate := func(t *testing.T, db *sql.DB, mutate func(*longFixture)) (*GroupBoundary, *longFixture, error) {
		t.Helper()
		f := newCounterFixture(t)
		f.config.Ledger = longLedger(db)
		if mutate != nil {
			mutate(f)
		}
		freezeDurableFixture(t, f)
		boundary, err := NewGroupBoundary(f.config)
		return boundary, f, err
	}
	code := func(t *testing.T, err error) string {
		t.Helper()
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked)
		return blocked.Code
	}

	t.Run("a damaged baseline blocks with an actionable code and is never reinterpreted", func(t *testing.T) {
		for name, rewrite := range map[string]func(string) string{
			"not json":       func(string) string { return "{" },
			"trailing bytes": func(p string) string { return p + "x" },
			"empty object":   func(string) string { return "{}" },
			"another member": func(p string) string {
				return strings.ReplaceAll(p, `"member_key":"[\"device-A\",\"point-uint\",\"tag-uint\"]"`, `"member_key":"[\"device-A\",\"point-int\",\"tag-int\"]"`)
			},
		} {
			db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
			seedBaseline(t, db, rewrite)
			_, _, err := activate(t, db, nil)
			require.Equal(t, "counter-baseline-invalid", code(t, err), name)
		}
	})

	t.Run("a bad digest is the same diagnostic", func(t *testing.T) {
		db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
		seedBaseline(t, db, func(p string) string { return p })
		_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET payload = replace(payload, '"100"', '"101"')`)
		require.NoError(t, err)
		_, _, activateErr := activate(t, db, nil)
		require.Equal(t, "counter-baseline-invalid", code(t, activateErr))
	})

	t.Run("a checkpoint of another format version is named", func(t *testing.T) {
		db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
		seedBaseline(t, db, func(p string) string { return strings.Replace(p, `"version":2`, `"version":1`, 1) })
		_, _, err := activate(t, db, nil)
		require.Equal(t, "counter-baseline-version-unsupported", code(t, err))
	})

	t.Run("a same-revision epoch replacement is refused without resetting the baseline", func(t *testing.T) {
		db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
		seedBaseline(t, db, func(p string) string { return p })
		before, err := longLedger(db).Restore(t.Context())
		require.NoError(t, err)
		boundary, _, err := activate(t, db, func(f *longFixture) {
			source := f.config.Counters["tag-uint"]
			source.SeriesEpoch = "epoch-2"
			f.config.Counters["tag-uint"] = source
		})
		require.Nil(t, boundary)
		require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
		after, err := longLedger(db).Restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, before, after)
	})

	t.Run("an intact baseline restores", func(t *testing.T) {
		db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
		seedBaseline(t, db, func(p string) string { return p })
		boundary, f, err := activate(t, db, nil)
		require.NoError(t, err)
		f.clock = at(2 * counterWindow)
		f.accept(t, boundary, f.envelope(1, "uint-1", at(2*counterWindow-time.Minute), uint64(105)))
		require.NoError(t, boundary.Tick(t.Context(), f.clock))
		require.Equal(t, "5", outboxCell(t, db, counterWindow)["delta_value"])
	})
}

// duplicatingLedger restores every baseline twice, the way an ambiguous durable
// state would look, and says it persists baselines.
type duplicatingLedger struct{ inner GroupLedger }

func (l duplicatingLedger) PersistsCounterBaselines() bool { return true }
func (l duplicatingLedger) AppendSample(ctx context.Context, start time.Time, sample snapshot.Sample) error {
	return l.inner.AppendSample(ctx, start, sample)
}
func (l duplicatingLedger) CommitClosure(ctx context.Context, closure groupdelivery.Closure) error {
	return l.inner.CommitClosure(ctx, closure)
}
func (l duplicatingLedger) Restore(ctx context.Context) (groupdelivery.Restored, error) {
	restored, err := l.inner.Restore(ctx)
	restored.Baselines = append(restored.Baselines, restored.Baselines...)
	return restored, err
}

// TestCumulativeProfileBaselineScopeIsTheGroupsOwn pins that the workspace,
// group and profile revision of a restored checkpoint are the boundary's fixed
// scope: a valid checkpoint of another scope, an ambiguous restore and a
// baseline ahead of the closed cursor are damaged state. Changed measurement
// semantics cannot replace the facts of the same immutable revision.
func TestCumulativeProfileBaselineScopeIsTheGroupsOwn(t *testing.T) {
	blockedCode := func(t *testing.T, config GroupBoundaryConfig) string {
		t.Helper()
		_, err := NewGroupBoundary(config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked)
		return blocked.Code
	}
	seeded := func(t *testing.T, rewrite func(string) string) (*sql.DB, *longFixture) {
		t.Helper()
		db := openDurableDB(t, filepath.Join(t.TempDir(), "gateway.db"))
		seedBaseline(t, db, rewrite)
		f := newCounterFixture(t)
		f.config.Ledger = longLedger(db)
		freezeDurableFixture(t, f)
		return db, f
	}

	for name, replace := range map[string][2]string{
		"another workspace":               {`"workspace_id":"workspace-A"`, `"workspace_id":"workspace-B"`},
		"another group":                   {`"group_id":"group-L"`, `"group_id":"group-other"`},
		"another profile revision":        {`"profile_revision":"rev-1"`, `"profile_revision":"rev-0"`},
		"a window that is not closed yet": {`"window_end":"2026-01-01T00:10:00Z"`, `"window_end":"2026-01-01T01:10:00Z"`},
	} {
		t.Run("a valid checkpoint with "+name+" is refused", func(t *testing.T) {
			_, f := seeded(t, func(payload string) string {
				require.Contains(t, payload, replace[0])
				return strings.Replace(payload, replace[0], replace[1], 1)
			})
			require.Equal(t, "counter-baseline-invalid", blockedCode(t, f.config))
		})
	}

	t.Run("an ambiguous restore is refused", func(t *testing.T) {
		_, f := seeded(t, func(payload string) string { return payload })
		f.config.Ledger = duplicatingLedger{inner: f.config.Ledger}
		require.Equal(t, "counter-baseline-invalid", blockedCode(t, f.config))
	})

	for name, mutate := range map[string]func(*CounterSource){
		"definition revision": func(s *CounterSource) { s.DefinitionRevision = "def-2" },
		"source binding":      func(s *CounterSource) { s.SourceBindingRevision = "bind-2" },
		"unit":                func(s *CounterSource) { s.Unit = "MWh" },
		"rollover policy": func(s *CounterSource) {
			modulus, rate := measurement.NewUint64(65536), measurement.NewUint64(100)
			s.RolloverModulus, s.MaxRatePerSecond = &modulus, &rate
		},
	} {
		t.Run("same-revision "+name+" replacement is refused", func(t *testing.T) {
			db, f := seeded(t, func(payload string) string { return payload })
			before, err := longLedger(db).Restore(t.Context())
			require.NoError(t, err)
			source := f.config.Counters["tag-uint"]
			mutate(&source)
			f.config.Counters["tag-uint"] = source
			boundary, err := NewGroupBoundary(f.config)
			require.Nil(t, boundary)
			require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
			after, err := longLedger(db).Restore(t.Context())
			require.NoError(t, err)
			require.Equal(t, before, after)
		})
	}
}
