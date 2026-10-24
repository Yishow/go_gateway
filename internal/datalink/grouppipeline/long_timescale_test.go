package grouppipeline

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// hypertableFacts reads the catalog's view of the destination table: the number
// of dimensions and chunks, the time column and interval, and the jobs bound to it.
type hypertableFacts struct {
	found                    bool
	dimensions, chunks, jobs int
	timeColumn               string
	intervalSeconds          float64
	compression              bool
}

func (f *longPipeline) hypertable() hypertableFacts {
	f.t.Helper()
	admin, schemaName := f.postgres.admin, f.postgres.schema
	var facts hypertableFacts
	err := admin.QueryRowContext(f.t.Context(), `
		SELECT num_dimensions, num_chunks, compression_enabled FROM timescaledb_information.hypertables
		WHERE hypertable_schema = $1 AND hypertable_name = $2`, schemaName, f.table).Scan(&facts.dimensions, &facts.chunks, &facts.compression)
	if err == sql.ErrNoRows {
		return facts
	}
	require.NoError(f.t, err)
	facts.found = true
	require.NoError(f.t, admin.QueryRowContext(f.t.Context(), `
		SELECT column_name, EXTRACT(EPOCH FROM time_interval)::float8 FROM timescaledb_information.dimensions
		WHERE hypertable_schema = $1 AND hypertable_name = $2 ORDER BY dimension_number LIMIT 1`, schemaName, f.table).Scan(&facts.timeColumn, &facts.intervalSeconds))
	require.NoError(f.t, admin.QueryRowContext(f.t.Context(), `
		SELECT COUNT(*) FROM timescaledb_information.jobs WHERE hypertable_schema = $1 AND hypertable_name = $2`, schemaName, f.table).Scan(&facts.jobs))
	return facts
}

// keyColumns lists the primary-key columns and the non-key index definitions of
// the destination table, as PostgreSQL reports them.
func (f *longPipeline) keyAndIndexes() (key, indexes []string) {
	f.t.Helper()
	admin := f.postgres.admin
	rows, err := admin.QueryContext(f.t.Context(), `
		SELECT a.attname FROM pg_index i
		JOIN pg_class c ON c.oid = i.indrelid JOIN pg_namespace n ON n.oid = c.relnamespace
		JOIN LATERAL unnest(i.indkey::int2[]) WITH ORDINALITY k(attnum, ord) ON true
		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = k.attnum
		WHERE n.nspname = $1 AND c.relname = $2 AND i.indisprimary ORDER BY k.ord`, f.postgres.schema, f.table)
	require.NoError(f.t, err)
	for rows.Next() {
		var column string
		require.NoError(f.t, rows.Scan(&column))
		key = append(key, column)
	}
	require.NoError(f.t, rows.Err())
	require.NoError(f.t, rows.Close())
	listed, err := admin.QueryContext(f.t.Context(), `SELECT indexname FROM pg_indexes WHERE schemaname = $1 AND tablename = $2 ORDER BY indexname`, f.postgres.schema, f.table)
	require.NoError(f.t, err)
	defer listed.Close()
	for listed.Next() {
		var name string
		require.NoError(f.t, listed.Scan(&name))
		indexes = append(indexes, name)
	}
	require.NoError(f.t, listed.Err())
	return key, indexes
}

// TestTimescalePipelineWritesLikePlain runs one cumulative scenario through the
// production pipeline into a plain table and into an owned hypertable of the same
// typed layout, and requires the same logical records. It also proves from the
// catalog that the timescale table is exactly the confirmed hypertable and that
// the plain one is not.
func TestTimescalePipelineWritesLikePlain(t *testing.T) {
	scenario := func(f *longPipeline) {
		f.closeFirstWindow()
		f.deliver()
		f.acceptSecondWindow()
		f.closeUntil(at(1200))
		f.deliver()
		f.requireDestinationExact()
	}
	plain := newLongPipelinePostgres(t)
	scenario(plain)
	hyper := newLongPipelineTimescale(t)
	scenario(hyper)

	require.NotEqual(t, plain.table, hyper.table, "the mode is part of the storage scope, so the tables differ")
	require.False(t, plain.hypertable().found, "plain storage stays an ordinary table")

	facts := hyper.hypertable()
	require.True(t, facts.found, "the timescale destination is a hypertable")
	require.Equal(t, 1, facts.dimensions)
	require.Equal(t, "bucket_start", facts.timeColumn)
	require.InDelta(t, 7*24*3600, facts.intervalSeconds, 0.5)
	require.False(t, facts.compression)
	require.Zero(t, facts.jobs, "no retention, compression or other job exists")
	key, indexes := hyper.keyAndIndexes()
	require.Equal(t, []string{"bucket_start", "record_id"}, key, "the composite unique key includes the partition column")
	require.Equal(t, []string{hyper.table + "_pkey", "idx_" + hyper.table + "_src_time"}, indexes,
		"only the confirmed key and source/time index; no default time index was added")

	plainRows, hyperRows := plain.targetRows(), hyper.targetRows()
	require.Len(t, hyperRows, 4)
	require.Equal(t, plainRows, hyperRows, "identical inputs give equivalent typed raw values, deltas, status and provenance in both modes")
}

// TestTimescaleSpansChunksAndReplayAddsNoRow writes records seven days and more
// apart through the production encoder, outbox and sender into the hypertable:
// they land in separate chunks, read back in order with exact values, and neither
// a repeated closure nor a delivery whose acknowledgement was lost adds a row.
func TestTimescaleSpansChunksAndReplayAddsNoRow(t *testing.T) {
	f := newLongPipelineTimescale(t)
	member := f.group.Members[0]
	memberKey := fmt.Sprintf("[%q,%q,%q]", member.DeviceID, member.PointID, member.TagID)
	layout, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: dbtarget.SQLDialectPostgres, Columns: f.columns, WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID,
		GroupRevision: f.group.Revision, ProfileKind: "cumulative", ProfileInterval: 10, ProfileUnit: "minute",
		Members: []dbtarget.LongRowMember{{
			MemberKey: memberKey, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID, Type: measurement.ExactUint64,
		}},
	})
	require.Empty(t, issues)
	scope, err := snapshot.DestinationScope("connector-1", "crev-1", "", f.group.Destination.TableSchema, f.table)
	require.NoError(t, err)

	starts := []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 9, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 17, 6, 0, 0, 0, time.UTC), time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC),
	}
	buckets := make([]groupdelivery.ClosedBucket, 0, len(starts))
	for i, start := range starts {
		recordID, err := snapshot.ProfileRecordID(f.group.WorkspaceID, f.group.ID, f.group.Revision, memberKey, "cumulative", start)
		require.NoError(t, err)
		effect, err := snapshot.EffectKey(scope, recordID)
		require.NoError(t, err)
		observed := start.Add(9 * time.Minute)
		outcome := snapshot.Outcome{
			Kind: snapshot.OutcomeRow, EntityKey: memberKey, RecordID: recordID, EffectKey: effect,
			BucketStart: start, BucketEnd: start.Add(10 * time.Minute),
			Members: []snapshot.MemberResult{{MemberKey: memberKey, Status: snapshot.MemberOK, Sample: &snapshot.Sample{
				SampleID: fmt.Sprintf("s-%d", i), MemberKey: memberKey, ObservedAt: observed, ReceivedAt: observed,
				Quality: schema.QualityGood, Value: measurement.NewUint64(9007199254740993 + uint64(i)),
			}}},
		}
		row, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		buckets = append(buckets, groupdelivery.ClosedBucket{Outcome: outcome, Row: &row})
	}
	closure := groupdelivery.Closure{
		Destination: groupdelivery.FrozenDestination{
			Scope: scope, ConnectorID: "connector-1", ConnectorRevision: "crev-1", TableSchema: f.group.Destination.TableSchema,
			TableName: f.table, DedupeCapability: groupdelivery.DedupeReceipt,
		},
		NextClose: starts[len(starts)-1].Add(10 * time.Minute), Buckets: buckets,
	}
	key := groupdelivery.GroupKey{WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, GroupRevision: f.group.Revision}
	require.NoError(t, f.store.CommitClosure(t.Context(), key, closure))
	f.deliver()

	admin, qualified := f.postgres.admin, fmt.Sprintf("%q.%q", f.postgres.schema, f.table)
	count := func(query string) int {
		var n int
		require.NoError(t, admin.QueryRowContext(t.Context(), query).Scan(&n))
		return n
	}
	require.Equal(t, 4, count(`SELECT COUNT(*) FROM `+qualified))
	require.GreaterOrEqual(t, f.hypertable().chunks, 4, "records seven days and more apart live in separate chunks (catalog)")
	require.Equal(t, 4, count(`SELECT COUNT(DISTINCT tableoid) FROM `+qualified), "and every row is physically in its own chunk")

	rows, err := admin.QueryContext(t.Context(), `SELECT CAST(value_decimal AS TEXT), to_char(bucket_start AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI') FROM `+qualified+` ORDER BY bucket_start`)
	require.NoError(t, err)
	var read []string
	for rows.Next() {
		var value, bucket string
		require.NoError(t, rows.Scan(&value, &bucket))
		read = append(read, bucket+"="+value)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.Equal(t, []string{
		"2026-01-01T00:00=9007199254740993", "2026-01-09T12:00=9007199254740994",
		"2026-01-17T06:00=9007199254740995", "2026-01-26T00:00=9007199254740996",
	}, read, "ordered queries across chunks keep every exact value")

	// The same closure again is a no-op, and a delivery whose acknowledgement was
	// lost (every effect back to pending) finds its receipt and writes nothing.
	require.NoError(t, f.store.CommitClosure(t.Context(), key, closure))
	_, err = f.plainLocal().ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET state = 'pending', committed_at = ''`)
	require.NoError(t, err)
	f.deliver()
	require.Equal(t, 4, count(`SELECT COUNT(*) FROM `+qualified), "a replay adds no row")
	require.Equal(t, 4, count(fmt.Sprintf(`SELECT COUNT(*) FROM %q.%q`, f.postgres.schema, "gw_effect_receipts")), "and one receipt per effect remains")
	require.Equal(t, 4, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_outbox WHERE state = 'sql_committed'`))

	// The composite key itself refuses a second row for an identity.
	_, err = admin.ExecContext(t.Context(), fmt.Sprintf(`INSERT INTO %s SELECT * FROM %s LIMIT 1`, qualified, qualified))
	require.Error(t, err, "the unique key includes bucket_start and record_id, so a duplicate is refused by the database")
}
