package recordingplan

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

const (
	longTestTable = "gw_long_0123456789abcdef"
	longTestOwner = "_gw_owner_0123456789abcdef01234567"
)

// longColumnSpec is the documented fixed column contract:
// name|SQLite type|PostgreSQL type|N (NOT NULL) or n (nullable)|P (primary key) or -.
var longColumnSpec = []string{
	"bucket_start|TEXT|TIMESTAMPTZ|N|P",
	"bucket_end|TEXT|TIMESTAMPTZ|N|-",
	"record_id|TEXT|TEXT|N|P",
	"workspace_id|TEXT|TEXT|N|-",
	"group_id|TEXT|TEXT|N|-",
	"group_revision|TEXT|TEXT|N|-",
	"circuit_id|TEXT|TEXT|n|-",
	"device_id|TEXT|TEXT|N|-",
	"point_id|TEXT|TEXT|N|-",
	"tag_id|TEXT|TEXT|N|-",
	"measurement_id|TEXT|TEXT|n|-",
	"series_epoch|TEXT|TEXT|n|-",
	"profile_kind|TEXT|TEXT|N|-",
	"value_type|TEXT|TEXT|N|-",
	"unit|TEXT|TEXT|n|-",
	"status|TEXT|TEXT|N|-",
	"quality|TEXT|TEXT|n|-",
	"observed_at|TEXT|TIMESTAMPTZ|n|-",
	"received_at|TEXT|TIMESTAMPTZ|n|-",
	"value_int|INTEGER|BIGINT|n|-",
	"value_decimal|TEXT|NUMERIC|n|-",
	"value_float|REAL|DOUBLE PRECISION|n|-",
	"value_bool|INTEGER|BOOLEAN|n|-",
	"value_text|TEXT|TEXT|n|-",
	"delta_value|TEXT|NUMERIC|n|-",
	"delta_from|TEXT|TIMESTAMPTZ|n|-",
	"delta_to|TEXT|TIMESTAMPTZ|n|-",
	"provenance|TEXT|JSONB|N|-",
	longTestOwner + "|TEXT|TEXT|N|-",
}

func expectedLongColumns(dialect string) []GroupSchemaColumn {
	columns := make([]GroupSchemaColumn, 0, len(longColumnSpec))
	for _, line := range longColumnSpec {
		parts := strings.Split(line, "|")
		sqlType := parts[1]
		if dialect == dialectPostgres {
			sqlType = parts[2]
		}
		columns = append(columns, GroupSchemaColumn{Name: parts[0], SQLType: sqlType, Nullable: parts[3] == "n", PrimaryKey: parts[4] == "P"})
	}
	return columns
}

func TestLongLayoutExactTypes(t *testing.T) {
	for _, dialect := range []string{dialectSQLite, dialectPostgres} {
		t.Run(dialect+" fixed typed columns", func(t *testing.T) {
			layout, err := ManagedLongGroupLayout(dialect, longTestTable, longTestOwner)
			require.NoError(t, err)
			require.Equal(t, longTestTable, layout.TableName)
			require.Equal(t, longTestOwner, layout.OwnerColumn)
			require.Equal(t, expectedLongColumns(dialect), layout.Columns)
			require.Equal(t, []GroupSchemaIndex{{
				Name: "idx_" + longTestTable + "_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"},
			}}, layout.Indexes)
			_, err = ValidateGroupSchemaLayout(layout)
			require.NoError(t, err)
		})
	}

	t.Run("the layout never depends on sources", func(t *testing.T) {
		first, err := ManagedLongGroupLayout(dialectSQLite, longTestTable, longTestOwner)
		require.NoError(t, err)
		second, err := ManagedLongGroupLayout(dialectSQLite, longTestTable, longTestOwner)
		require.NoError(t, err)
		require.Equal(t, first, second)
		for _, column := range first.Columns {
			require.NotRegexp(t, `^v_[0-9a-f]{24}$`, column.Name, "a source must never become a column")
		}
	})

	t.Run("postgres DDL uses one composite key and an explicit index", func(t *testing.T) {
		layout, err := ManagedLongGroupLayout(dialectPostgres, longTestTable, longTestOwner)
		require.NoError(t, err)
		statements, err := GenerateGroupSchemaDDL(dialectPostgres, layout)
		require.NoError(t, err)
		require.Len(t, statements, 3)
		create := statements[0]
		require.Equal(t, 1, strings.Count(create, "PRIMARY KEY"), "exactly one table-level key: %s", create)
		require.Contains(t, create, `PRIMARY KEY ("bucket_start", "record_id")`)
		for _, fragment := range []string{
			`"bucket_start" TIMESTAMPTZ NOT NULL`, `"value_int" BIGINT`, `"value_decimal" NUMERIC`,
			`"value_float" DOUBLE PRECISION`, `"value_bool" BOOLEAN`, `"delta_value" NUMERIC`, `"provenance" JSONB NOT NULL`,
		} {
			require.Contains(t, create, fragment)
		}
		require.NotContains(t, create, "NUMERIC(", "decimal and delta must stay unbounded NUMERIC")
		require.Equal(t, `CREATE INDEX "idx_`+longTestTable+`_src_time" ON "`+longTestTable+`" ("workspace_id", "tag_id", "bucket_start");`, statements[1])
		require.Contains(t, statements[2], `"gw_effect_receipts"`)
		require.NotContains(t, strings.ToUpper(strings.Join(statements, "\n")), "RETENTION")
	})

	t.Run("sqlite catalog round trip keeps exact typed values", func(t *testing.T) {
		db := openLongSQLite(t)
		layout, err := ManagedLongGroupLayout(dialectSQLite, longTestTable, longTestOwner)
		require.NoError(t, err)
		statements, err := GenerateGroupSchemaDDL(dialectSQLite, layout)
		require.NoError(t, err)
		for _, statement := range statements {
			_, err := db.ExecContext(t.Context(), statement)
			require.NoError(t, err, statement)
		}

		// The real catalog matches the layout, including the composite key and index.
		state, err := sqliteInspection(t.Context(), db, longTestTable)
		require.NoError(t, err)
		require.NoError(t, ValidateManagedGroupTableInspection(layout, state))
		pending, _, tables, err := inspectGroupManagedTables(t.Context(), dialectSQLite, layout, func(ctx context.Context, table string) (TargetTableInspection, error) {
			return sqliteInspection(ctx, db, table)
		})
		require.NoError(t, err)
		require.False(t, pending)
		require.Equal(t, SchemaTableActionUnchanged, tables[0].Action)

		const insert = `INSERT INTO "` + longTestTable + `" (bucket_start, bucket_end, record_id, workspace_id, group_id, group_revision,
			circuit_id, device_id, point_id, tag_id, measurement_id, series_epoch, profile_kind, value_type, unit, status, quality,
			observed_at, received_at, value_int, value_decimal, value_float, value_bool, value_text, delta_value, delta_from, delta_to, provenance)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		const bucket = "2026-10-24T00:00:15.000000Z"
		row := func(recordID string, circuit any, tag string, typed [5]any, delta any) []any {
			return []any{bucket, "2026-10-24T00:00:30.000000Z", recordID, "ws", "g1", "r1", circuit, "dev", "pt", tag, nil, nil,
				"cumulative", "uint64", "kWh", "ok", "good", "2026-10-24T00:00:29.123456Z", "2026-10-24T00:00:29.500000Z",
				typed[0], typed[1], typed[2], typed[3], typed[4], delta, "2026-10-24T00:00:15.000000Z", "2026-10-24T00:00:30.000000Z", `{"type":"uint64"}`}
		}
		const maxUint64 = "18446744073709551615"
		_, err = db.ExecContext(t.Context(), insert, row("rec-a", "circuit-a", "kwh", [5]any{nil, maxUint64, nil, nil, nil}, "1")...)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), insert, row("rec-b", "circuit-b", "kwh", [5]any{int64(9007199254740993), nil, nil, nil, nil}, nil)...)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), insert, row("rec-c", nil, "flag", [5]any{nil, nil, nil, int64(1), nil}, nil)...)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), insert, row("rec-d", nil, "ratio", [5]any{nil, nil, 0.1, nil, nil}, nil)...)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), insert, row("rec-e", nil, "label", [5]any{nil, nil, nil, nil, "run"}, nil)...)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), insert, row("rec-f", nil, "dec", [5]any{nil, "12000.4", nil, nil, nil}, "0.3")...)
		require.NoError(t, err)

		var dec, deltaText, decType, deltaType, intType, boolType, floatType, textType, observed, provenance string
		var big int64
		var flag int64
		var ratio float64
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_decimal, typeof(value_decimal), delta_value, typeof(delta_value) FROM "`+longTestTable+`" WHERE record_id='rec-a'`).
			Scan(&dec, &decType, &deltaText, &deltaType))
		require.Equal(t, maxUint64, dec, "uint64 must not pass through float")
		require.Equal(t, "text", decType)
		require.Equal(t, "1", deltaText)
		require.Equal(t, "text", deltaType)
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_int, typeof(value_int) FROM "`+longTestTable+`" WHERE record_id='rec-b'`).Scan(&big, &intType))
		require.Equal(t, int64(9007199254740993), big)
		require.Equal(t, "integer", intType)
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_bool, typeof(value_bool) FROM "`+longTestTable+`" WHERE record_id='rec-c'`).Scan(&flag, &boolType))
		require.Equal(t, int64(1), flag)
		require.Equal(t, "integer", boolType)
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_float, typeof(value_float) FROM "`+longTestTable+`" WHERE record_id='rec-d'`).Scan(&ratio, &floatType))
		require.InDelta(t, 0.1, ratio, 0)
		require.Equal(t, "real", floatType)
		var text string
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_text, typeof(value_text) FROM "`+longTestTable+`" WHERE record_id='rec-e'`).Scan(&text, &textType))
		require.Equal(t, "run", text)
		require.Equal(t, "text", textType)
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_decimal FROM "`+longTestTable+`" WHERE record_id='rec-f'`).Scan(&dec))
		require.Equal(t, "12000.4", dec, "decimal keeps its exact text")
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT observed_at, provenance FROM "`+longTestTable+`" WHERE record_id='rec-a'`).Scan(&observed, &provenance))
		require.Equal(t, "2026-10-24T00:00:29.123456Z", observed)
		require.JSONEq(t, `{"type":"uint64"}`, provenance)

		// Same label on two circuits stays two rows; only one typed slot is set per row.
		var rows, typed int
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+longTestTable+`" WHERE tag_id='kwh'`).Scan(&rows))
		require.Equal(t, 2, rows)
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+longTestTable+`" WHERE
			(value_int IS NOT NULL)+(value_decimal IS NOT NULL)+(value_float IS NOT NULL)+(value_bool IS NOT NULL)+(value_text IS NOT NULL) > 1`).Scan(&typed))
		require.Zero(t, typed)

		// The composite key rejects a replayed record identity but allows the same ID in another window.
		_, err = db.ExecContext(t.Context(), insert, row("rec-a", "circuit-a", "kwh", [5]any{nil, maxUint64, nil, nil, nil}, "1")...)
		require.Error(t, err)
		other := row("rec-a", "circuit-a", "kwh", [5]any{nil, maxUint64, nil, nil, nil}, "1")
		other[0] = "2026-10-24T00:00:30.000000Z"
		_, err = db.ExecContext(t.Context(), insert, other...)
		require.NoError(t, err)
	})

	t.Run("a missing index or foreign layout is incompatible, not accepted", func(t *testing.T) {
		db := openLongSQLite(t)
		layout, err := ManagedLongGroupLayout(dialectSQLite, longTestTable, longTestOwner)
		require.NoError(t, err)
		statements, err := GenerateGroupSchemaDDL(dialectSQLite, layout)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), statements[0]) // table without its index
		require.NoError(t, err)
		state, err := sqliteInspection(t.Context(), db, longTestTable)
		require.NoError(t, err)
		require.ErrorIs(t, ValidateManagedGroupTableInspection(layout, state), ErrIncompatibleExistingTable)

		foreign, err := ManagedLongGroupLayout(dialectSQLite, longTestTable, "_gw_owner_ffffffffffffffffffffffff")
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), statements[1])
		require.NoError(t, err)
		state, err = sqliteInspection(t.Context(), db, longTestTable)
		require.NoError(t, err)
		require.NoError(t, ValidateManagedGroupTableInspection(layout, state))
		require.ErrorIs(t, ValidateManagedGroupTableInspection(foreign, state), ErrIncompatibleExistingTable, "another owner's marker must not validate")
	})
}

func openLongSQLite(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "long.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}
