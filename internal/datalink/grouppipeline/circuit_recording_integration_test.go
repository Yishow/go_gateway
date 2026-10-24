package grouppipeline

import (
	"slices"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

// runStandardScenario drives two cumulative windows through the production
// pipeline and verifies exact raw and delta values.
func runStandardScenario(f *longPipeline) {
	f.t.Helper()
	f.closeFirstWindow()
	f.deliver()
	f.acceptSecondWindow()
	f.closeUntil(at(1200))
	f.deliver()
	f.requireDestinationExact()
}

// requireLogicalRowsEquivalent checks that two destinations hold equivalent
// logical rows for all counters and windows, ignoring mode-specific metadata.
func requireLogicalRowsEquivalent(t *testing.T, expected, actual map[string]targetRow) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "row count must match")
	for key, expRow := range expected {
		actRow, ok := actual[key]
		require.True(t, ok, "key %s must exist in actual rows", key)
		require.Equal(t, expRow.status, actRow.status, "%s status", key)
		require.Equal(t, expRow.raw, actRow.raw, "%s raw", key)
		require.Equal(t, expRow.delta, actRow.delta, "%s delta", key)
		require.Equal(t, expRow.deltaFrom, actRow.deltaFrom, "%s deltaFrom", key)
		require.Equal(t, expRow.deltaTo, actRow.deltaTo, "%s deltaTo", key)
		require.JSONEq(t, expRow.provenance.String, actRow.provenance.String, "%s provenance", key)
	}
}

// TestCircuitRecordingIntegration is the primary integration matrix entry for
// circuit recording profiles across SQLite, plain PostgreSQL, and TimescaleDB.
func TestCircuitRecordingIntegration(t *testing.T) {
	t.Run("sqlite", func(t *testing.T) {
		f := newLongPipeline(t)
		runStandardScenario(f)
	})

	t.Run("postgres", func(t *testing.T) {
		f := newLongPipelinePostgres(t)
		runStandardScenario(f)
	})

	t.Run("timescale", func(t *testing.T) {
		f := newLongPipelineTimescale(t)
		runStandardScenario(f)
	})

	t.Run("logical_equivalence", func(t *testing.T) {
		sqlitePipe := newLongPipeline(t)
		runStandardScenario(sqlitePipe)
		pgPipe := newLongPipelinePostgres(t)
		runStandardScenario(pgPipe)
		sqliteRows := sqlitePipe.targetRows()
		pgRows := pgPipe.targetRows()
		require.Len(t, sqliteRows, 4)
		requireLogicalRowsEquivalent(t, sqliteRows, pgRows)
		t.Run("timescale", func(t *testing.T) {
			tsPipe := newLongPipelineTimescale(t)
			runStandardScenario(tsPipe)
			requireLogicalRowsEquivalent(t, pgRows, tsPipe.targetRows())
		})
	})

	t.Run("timescale_catalog_and_hypertable", func(t *testing.T) {
		hyper := newLongPipelineTimescale(t)
		runStandardScenario(hyper)

		facts := hyper.hypertable()
		require.True(t, facts.found, "the timescale destination must be a hypertable")
		require.Equal(t, 1, facts.dimensions, "single time dimension")
		require.Equal(t, "bucket_start", facts.timeColumn, "partition dimension column")
		var intervalMicros int64
		require.NoError(t, hyper.postgres.admin.QueryRowContext(t.Context(), `SELECT (EXTRACT(EPOCH FROM time_interval) * 1000000)::bigint FROM timescaledb_information.dimensions WHERE hypertable_schema=$1 AND hypertable_name=$2`, hyper.postgres.schema, hyper.table).Scan(&intervalMicros))
		require.Equal(t, dbtarget.TimescaleChunkInterval.Microseconds(), intervalMicros, "exact confirmed chunk interval")
		require.False(t, facts.compression, "compression disabled")
		require.Zero(t, facts.jobs, "zero background retention/compression jobs")

		key, indexes := hyper.keyAndIndexes()
		require.True(t, slices.Equal([]string{"bucket_start", "record_id"}, key), "composite primary key")
		require.Contains(t, indexes, hyper.table+"_pkey")
		require.Contains(t, indexes, "idx_"+hyper.table+"_src_time")

		// Spanning multiple 7-day chunks.
		t.Run("cross_chunk_readback_and_receipts", TestTimescaleSpansChunksAndReplayAddsNoRow)
	})

	t.Run("sql_ack_lost_and_replay_dedup", func(t *testing.T) {
		for _, name := range []string{"sqlite", "postgres", "timescale"} {
			t.Run(name, func(t *testing.T) {
				var f *longPipeline
				switch name {
				case "sqlite":
					f = newLongPipeline(t)
				case "postgres":
					f = newLongPipelinePostgres(t)
				case "timescale":
					f = newLongPipelineTimescale(t)
				}
				testCommitAckLost(t, f)
			})
		}
	})

	t.Run("counter_reset_within_window", func(t *testing.T) {
		for _, name := range []string{"sqlite", "postgres", "timescale"} {
			t.Run(name, func(t *testing.T) {
				var f *longPipeline
				switch name {
				case "sqlite":
					f = newLongPipeline(t)
				case "postgres":
					f = newLongPipelinePostgres(t)
				case "timescale":
					f = newLongPipelineTimescale(t)
				}
				testCounterResetWithinWindow(t, f)
			})
		}
	})

	t.Run("gap_span_total", testCircuitIntegrationGap)

	t.Run("closure_crash_and_recovery", func(t *testing.T) {
		t.Run("sqlite", func(t *testing.T) {
			f := newLongPipeline(t)
			testClosureCrashRecovery(t, f)
		})
		t.Run("postgres", func(t *testing.T) {
			f := newLongPipelinePostgres(t)
			testClosureCrashRecovery(t, f)
		})
	})
}

func testCommitAckLost(t *testing.T, f *longPipeline) {
	t.Helper()
	f.closeFirstWindow()
	f.deliver()
	f.acceptSecondWindow()
	f.closeUntil(at(1200))
	frozen := f.outboxDigests()
	_, baselineBefore := f.baseline(0)

	f.targetFlt.Lose.Store(true)
	f.deliverCycles(1)
	require.NotEmpty(t, f.targetRows())

	f.targetFlt.Lose.Store(false)
	f.restart()
	f.deliver()
	f.requireDestinationExact()
	require.Equal(t, frozen, f.outboxDigests(), "retries resend the frozen payloads")
	require.Equal(t, 4, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))
	_, baselineAfter := f.baseline(0)
	require.Equal(t, baselineBefore, baselineAfter, "delivery retries never move the baseline")
	f.requireBaselinesAt(1)
}

func testCounterResetWithinWindow(t *testing.T, f *longPipeline) {
	t.Helper()
	f.closeFirstWindow()
	f.deliver()

	// In window 1: accept 100 at 12m, decrease to 10 at 15m, then 110 at 18m.
	require.NoError(t, f.reading(0, "m1-r1", at(12*60), firstStart+100))
	require.NoError(t, f.reading(0, "m1-r2", at(15*60), firstStart+10))
	require.NoError(t, f.reading(0, "m1-r3", at(18*60), firstStart+110))
	require.NoError(t, f.reading(1, "m2-r1", at(18*60), secondStart+30))

	f.closeUntil(at(1200))
	f.deliver()

	rows := f.targetRows()
	w1Key := windowKey(longKwhFirst, 1)
	row, ok := rows[w1Key]
	require.True(t, ok, "window 1 row must exist")
	require.Equal(t, "reset_or_uncertain", row.status.String, "unproven drop marks status reset_or_uncertain")
	require.Equal(t, "9007199254741102", row.raw.String, "raw reading is preserved")
	require.False(t, row.delta.Valid, "delta is NULL on unproven reset")
	require.Contains(t, row.provenance.String, "decrease-unproven", "reset evidence survives SQL encoding")

	// Member 2 without reset progressed normally.
	m2Key := windowKey(longKwhSecond, 1)
	require.Equal(t, "ok", rows[m2Key].status.String)
	require.Equal(t, "30", rows[m2Key].delta.String)
}

func testClosureCrashRecovery(t *testing.T, f *longPipeline) {
	t.Helper()
	f.closeFirstWindow()
	f.deliver()
	f.requireBaselinesAt(0)
	f.acceptSecondWindow()

	f.failClosureBaselineWrite()
	f.closeUntil(at(1200))
	require.NotEmpty(t, f.reportedErrors(), "failed closure reported")
	local := f.plainLocal()
	require.Equal(t, 2, f.count(local, `SELECT COUNT(*) FROM wg_delivery_outbox`))
	require.Equal(t, 2, f.count(local, `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed = 0`))
	f.requireBaselinesAt(0)

	f.localFault.Reject.Store(nil)
	f.closeUntil(at(1200))
	f.deliver()
	f.requireBaselinesAt(1)
	f.requireDestinationExact()
	require.Equal(t, 4, f.count(local, `SELECT COUNT(DISTINCT effect_key) FROM wg_delivery_outbox`))
}
