package recordingplan

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// TestLongLayoutExactTypesPostgres runs the generated long-table DDL on a live
// disposable PostgreSQL and checks that exact typed values survive a
// parameterized write and readback. It is skipped without POSTGRES_DSN and is
// then reported as not run, never as passed.
func TestLongLayoutExactTypesPostgres(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL long layout test: POSTGRES_DSN not set")
	}
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	namespace := fmt.Sprintf("gw_long_exact_%d", time.Now().UnixNano())
	_, err = db.ExecContext(t.Context(), `CREATE SCHEMA "`+namespace+`"`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = db.ExecContext(context.Background(), `DROP SCHEMA "`+namespace+`" CASCADE`)
	})

	layout, err := ManagedLongGroupLayout(dialectPostgres, longTestTable, longTestOwner)
	require.NoError(t, err)
	statements, err := GenerateGroupSchemaDDL(dialectPostgres, layout)
	require.NoError(t, err)
	conn, err := db.Conn(t.Context())
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	_, err = conn.ExecContext(t.Context(), `SET search_path TO "`+namespace+`"`)
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := conn.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}

	const insert = `INSERT INTO "` + longTestTable + `" (bucket_start, bucket_end, record_id, workspace_id, group_id, group_revision,
		circuit_id, device_id, point_id, tag_id, profile_kind, value_type, status, observed_at, received_at,
		value_int, value_decimal, value_float, value_bool, value_text, delta_value, delta_from, delta_to, provenance)
		VALUES ($1::timestamptz,$2::timestamptz,$3,'ws','g1','r1',$4,'dev','pt',$5,'cumulative','uint64','ok',$6::timestamptz,$6::timestamptz,
		$7::bigint,$8::numeric,$9::double precision,$10::boolean,$11,$12::numeric,$1::timestamptz,$2::timestamptz,$13::jsonb)`
	const bucket = "2026-10-24T00:00:15Z"
	const end = "2026-10-24T00:00:30Z"
	const maxUint64 = "18446744073709551615"
	rows := []struct {
		id      string
		circuit any
		tag     string
		args    [6]any // value_int, value_decimal, value_float, value_bool, value_text, delta_value
	}{
		{"rec-a", "circuit-a", "kwh", [6]any{nil, maxUint64, nil, nil, nil, "1"}},
		{"rec-b", "circuit-b", "kwh", [6]any{int64(9007199254740993), nil, nil, nil, nil, nil}},
		{"rec-c", nil, "flag", [6]any{nil, nil, nil, true, nil, nil}},
		{"rec-d", nil, "ratio", [6]any{nil, nil, 0.1, nil, nil, nil}},
		{"rec-e", nil, "label", [6]any{nil, nil, nil, nil, "run", nil}},
		{"rec-f", nil, "dec", [6]any{nil, "12000.4", nil, nil, nil, "0.3"}},
	}
	for _, row := range rows {
		_, err := conn.ExecContext(t.Context(), insert, bucket, end, row.id, row.circuit, row.tag,
			"2026-10-24T00:00:29.123456Z", row.args[0], row.args[1], row.args[2], row.args[3], row.args[4], row.args[5], `{"type":"uint64"}`)
		require.NoError(t, err, row.id)
	}

	var decimal, delta, kind string
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_decimal::text, delta_value::text, pg_typeof(value_decimal)::text FROM "`+longTestTable+`" WHERE record_id='rec-a'`).Scan(&decimal, &delta, &kind))
	require.Equal(t, maxUint64, decimal, "uint64 must not pass through float")
	require.Equal(t, "1", delta)
	require.Equal(t, "numeric", kind)
	var big int64
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_int FROM "`+longTestTable+`" WHERE record_id='rec-b'`).Scan(&big))
	require.Equal(t, int64(9007199254740993), big)
	var flag bool
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_bool FROM "`+longTestTable+`" WHERE record_id='rec-c'`).Scan(&flag))
	require.True(t, flag)
	var ratio float64
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_float FROM "`+longTestTable+`" WHERE record_id='rec-d'`).Scan(&ratio))
	require.InDelta(t, 0.1, ratio, 0)
	var text string
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_text FROM "`+longTestTable+`" WHERE record_id='rec-e'`).Scan(&text))
	require.Equal(t, "run", text)
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT value_decimal::text, delta_value::text FROM "`+longTestTable+`" WHERE record_id='rec-f'`).Scan(&decimal, &delta))
	require.Equal(t, "12000.4", decimal)
	require.Equal(t, "0.3", delta)
	var observed time.Time
	var provenance string
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT observed_at, provenance::text FROM "`+longTestTable+`" WHERE record_id='rec-a'`).Scan(&observed, &provenance))
	require.Equal(t, time.Date(2026, 10, 24, 0, 0, 29, 123456000, time.UTC), observed.UTC(), "microseconds survive in UTC")
	require.JSONEq(t, `{"type":"uint64"}`, provenance)

	// A replayed record identity is rejected by the composite key; the same ID in another window is not.
	_, err = conn.ExecContext(t.Context(), insert, bucket, end, "rec-a", "circuit-a", "kwh",
		"2026-10-24T00:00:29.123456Z", nil, maxUint64, nil, nil, nil, "1", `{"type":"uint64"}`)
	require.Error(t, err)
	_, err = conn.ExecContext(t.Context(), `SELECT 1`)
	require.NoError(t, err)
	var keyColumns string
	require.NoError(t, conn.QueryRowContext(t.Context(), `SELECT string_agg(a.attname::text, ',' ORDER BY k.ord) FROM pg_index ix
		JOIN pg_class t ON t.oid = ix.indrelid CROSS JOIN LATERAL unnest(ix.indkey::int2[]) WITH ORDINALITY AS k(attnum, ord)
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum
		WHERE t.relname = $1 AND ix.indisprimary`, longTestTable).Scan(&keyColumns))
	require.Equal(t, "bucket_start,record_id", keyColumns)
}
