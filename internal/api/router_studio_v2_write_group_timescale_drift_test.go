package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// chunkIntervalMicros is the catalog's chunk interval of a table in whole
// microseconds, exactly as PostgreSQL stores it, or "" when it is no hypertable.
func chunkIntervalMicros(t *testing.T, f managedGroupFixture, table string) string {
	t.Helper()
	var micros string
	err := f.target.QueryRowContext(t.Context(), `
		SELECT (EXTRACT(EPOCH FROM time_interval) * 1000000)::numeric::text FROM timescaledb_information.dimensions
		WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&micros)
	if err != nil {
		return ""
	}
	return micros
}

// TestTimescaleReadinessRejectsADriftedHypertable: a confirmed hypertable is
// only ready while the catalog still shows exactly the confirmed layout. Changing
// its chunk interval (by a day, or by a single microsecond), enabling
// compression, binding a retention job, or replacing it with a plain table of
// the same name blocks readiness and Apply, and Apply does not touch the
// destination. Restating the very same seven days stays ready.
func TestTimescaleReadinessRejectsADriftedHypertable(t *testing.T) {
	interval := func(expression string) func(q string) []string {
		return func(q string) []string {
			return []string{`SELECT set_chunk_time_interval('` + q + `'::regclass, ` + expression + `)`}
		}
	}
	for name, tc := range map[string]struct {
		change     func(qualified string) []string
		stillReady bool
	}{
		"another chunk interval":         {change: interval(`INTERVAL '1 day'`)},
		"seven days plus a microsecond":  {change: interval(`INTERVAL '7 days' + INTERVAL '1 microsecond'`)},
		"seven days minus a microsecond": {change: interval(`INTERVAL '7 days' - INTERVAL '1 microsecond'`)},
		"seven days plus 0.4 seconds":    {change: interval(`INTERVAL '7 days 400 milliseconds'`)},
		"an interval beyond a duration":  {change: interval(`INTERVAL '360000 days'`)},
		"the same seven days restated":   {change: interval(`INTERVAL '168 hours'`), stillReady: true},
		"compression enabled": {change: func(q string) []string {
			return []string{`ALTER TABLE ` + q + ` SET (timescaledb.compress)`}
		}},
		"a retention job": {change: func(q string) []string {
			return []string{`SELECT add_retention_policy('` + q + `'::regclass, INTERVAL '365 days')`}
		}},
		"replaced by a plain table": {change: func(q string) []string {
			return []string{`DROP TABLE ` + q, `CREATE TABLE ` + q + ` (record_id text primary key)`}
		}},
	} {
		t.Run(name, func(t *testing.T) {
			f := newManagedDestinationFixture(t, "postgres")
			if !timescaleInstalled(t, f) {
				t.Skip("this destination database has no TimescaleDB extension; the case is NOT RUN here")
			}
			seedProfileSource(t, f.writeGroupRouterFixture, "voltage")
			_, created := ensureTimescaleGroup(t, f)
			_, group := groupSaveData(t, created)
			groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
			confirmLong(t, f, groupID, previewLong(t, f, groupID))
			ready, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, ready.Ready, "%+v", ready.Issues)
			require.Equal(t, "604800000000.000000", chunkIntervalMicros(t, f, table), "the confirmed seven days, in microseconds")

			qualified := `"` + f.namespace + `"."` + table + `"`
			for _, statement := range tc.change(qualified) {
				_, err := f.target.ExecContext(t.Context(), statement)
				require.NoError(t, err, statement)
			}
			drifted := chunkIntervalMicros(t, f, table)
			relations, _ := catalogObjects(t, f)
			ready, err = f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
			if tc.stillReady {
				require.True(t, ready.Ready, "%+v", ready.Issues)
				require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
			} else {
				require.False(t, ready.Ready, "a drifted layout must not stay ready")
				require.Equal(t, http.StatusConflict, applied.Code, applied.Body.String())
			}
			afterRelations, _ := catalogObjects(t, f)
			require.Equal(t, [2]any{drifted, relations}, [2]any{chunkIntervalMicros(t, f, table), afterRelations}, "readiness and Apply never modify the destination")
		})
	}
}
