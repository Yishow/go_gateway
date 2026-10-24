package dbtarget

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTimescaleExactChunkInterval reads the chunk interval back from the real
// catalog at the microsecond PostgreSQL stores it. A user who later changes the
// confirmed seven days by a microsecond or a fraction of a second must not be
// read as still seven days, and a value too large to represent must never
// alias the confirmed one.
func TestTimescaleExactChunkInterval(t *testing.T) {
	admin := postgresAdminDB(t, mustDSN(t))
	requireTimescaleServer(t, admin)
	db := timescaleDB(t, `CREATE EXTENSION timescaledb`, `CREATE SCHEMA app`)
	for _, statement := range []string{
		`CREATE TABLE app.long_table (bucket_start TIMESTAMPTZ NOT NULL, record_id TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`,
		`SELECT create_hypertable('app.long_table'::regclass, by_range('bucket_start'::name, INTERVAL '7 days'), create_default_indexes => false, if_not_exists => false, migrate_data => false)`,
	} {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}

	for name, tc := range map[string]struct {
		interval string
		want     time.Duration
	}{
		"exactly seven days":                 {`INTERVAL '7 days'`, TimescaleChunkInterval},
		"the same duration spelled in hours": {`INTERVAL '168 hours'`, TimescaleChunkInterval},
		"seven days plus a microsecond":      {`INTERVAL '7 days' + INTERVAL '1 microsecond'`, TimescaleChunkInterval + time.Microsecond},
		"seven days minus a microsecond":     {`INTERVAL '7 days' - INTERVAL '1 microsecond'`, TimescaleChunkInterval - time.Microsecond},
		"seven days plus 0.4 seconds":        {`INTERVAL '7 days 400 milliseconds'`, TimescaleChunkInterval + 400*time.Millisecond},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := db.ExecContext(t.Context(), `SELECT set_chunk_time_interval('app.long_table'::regclass, `+tc.interval+`)`)
			require.NoError(t, err)
			info, err := InspectHypertable(t.Context(), db, "app", "long_table")
			require.NoError(t, err)
			require.NotNil(t, info)
			require.Equal(t, tc.want, info.ChunkInterval, "the catalog value is read exactly, not rounded to a second")
		})
	}

	t.Run("a value too large for a duration never aliases the confirmed interval", func(t *testing.T) {
		_, err := db.ExecContext(t.Context(), `SELECT set_chunk_time_interval('app.long_table'::regclass, INTERVAL '360000 days')`)
		require.NoError(t, err)
		info, err := InspectHypertable(t.Context(), db, "app", "long_table")
		require.NoError(t, err)
		require.NotNil(t, info)
		require.NotEqual(t, TimescaleChunkInterval, info.ChunkInterval)
		require.Zero(t, info.ChunkInterval)
		require.True(t, info.ChunkIntervalUnknown, "an interval a Duration cannot hold is reported as unknown, not as a number")
	})

	t.Run("a representable interval is known", func(t *testing.T) {
		_, err := db.ExecContext(t.Context(), `SELECT set_chunk_time_interval('app.long_table'::regclass, INTERVAL '7 days')`)
		require.NoError(t, err)
		info, err := InspectHypertable(t.Context(), db, "app", "long_table")
		require.NoError(t, err)
		require.False(t, info.ChunkIntervalUnknown)
		require.Equal(t, TimescaleChunkInterval, info.ChunkInterval)
	})
}

// TestDurationFromMicrosText pins the exact conversion of the catalog's
// microsecond count: whole microseconds only, no float, and no overflow that
// wraps around to an unrelated (even the confirmed) interval.
func TestDurationFromMicrosText(t *testing.T) {
	// 2305843614013693952 microseconds is 125 * 2^64 + 7 days in nanoseconds, so a
	// plain Duration conversion wraps around to exactly seven days.
	const wraps = "2305843614013693952"
	wrapMicros := int64(2305843614013693952)
	require.Equal(t, TimescaleChunkInterval, time.Duration(wrapMicros)*time.Microsecond, "the hazard is real: a naive conversion aliases the confirmed interval")

	for name, tc := range map[string]struct {
		text string
		want time.Duration
		ok   bool
	}{
		"seven days":                     {"604800000000.000000", TimescaleChunkInterval, true},
		"seven days and a microsecond":   {"604800000001.000000", TimescaleChunkInterval + time.Microsecond, true},
		"seven days less a microsecond":  {"604799999999", TimescaleChunkInterval - time.Microsecond, true},
		"no fraction part":               {"604800000000", TimescaleChunkInterval, true},
		"zero":                           {"0.000000", 0, true},
		"the largest duration in micros": {"9223372036854775", 9223372036854775 * time.Microsecond, true},
		"one past the largest duration":  {"9223372036854776", 0, false},
		"a wrap-around to seven days":    {wraps, 0, false},
		"beyond int64":                   {"99999999999999999999", 0, false},
		"a fraction of a microsecond":    {"604800000000.000001", 0, false},
		"negative":                       {"-604800000000.000000", 0, false},
		"signed":                         {"+604800000000", 0, false},
		"scientific notation":            {"6.048e11", 0, false},
		"not a number":                   {"seven days", 0, false},
		"empty":                          {"", 0, false},
		"only a fraction":                {".000000", 0, false},
		"embedded whitespace":            {" 604800000000", 0, false},
	} {
		t.Run(name, func(t *testing.T) {
			got, ok := durationFromMicrosText(tc.text)
			require.Equal(t, tc.ok, ok)
			require.Equal(t, tc.want, got)
		})
	}
}
