package dbtarget

import (
	"database/sql"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// TestLongRowCounterDiscontinuitiesPostgresReadback commits the 2.3 scenarios into
// the real PostgreSQL long table and reads status, signed NUMERIC deltas, span
// bounds and the JSONB events back. It needs POSTGRES_DSN like the other live
// PostgreSQL fixtures.
func TestLongRowCounterDiscontinuitiesPostgresReadback(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	tx, err := admin.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO `+pgIdentifier(t, schemaName))
	require.NoError(t, err)
	for _, statement := range longDDL(t, "postgres") {
		_, err = tx.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	require.NoError(t, tx.Commit())
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres,
		postgresConnectorConfig(fields, schemaName, fields["user"], fields["password"]))
	inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, schemaName, indexTestTable)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, inspection.Status)

	for name, scenario := range discontinuityScenarios(t) {
		layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectPostgres, inspection.Columns, scenario.member))
		require.Empty(t, issues, "the real PostgreSQL catalog satisfies the fixed layout")
		for _, outcome := range scenario.outcomes {
			outcome.RecordID, outcome.EffectKey = name+"-"+outcome.RecordID, name+"-"+outcome.EffectKey
			row := outcomesEncoded(t, layout, outcome)
			_, err := InsertGroupRow(t.Context(), admin, GroupInsertRequest{
				Kind: schema.DatabaseConnectorKindPostgres, SchemaName: schemaName, TableName: indexTestTable, Row: row,
				PayloadDigest: "digest-" + row.RecordID, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
			})
			require.NoError(t, err, row.RecordID)
		}
	}
	table := pgIdentifier(t, schemaName) + "." + pgIdentifier(t, indexTestTable)
	text := func(record, expression string) string {
		var value sql.NullString
		require.NoError(t, admin.QueryRowContext(t.Context(),
			fmt.Sprintf(`SELECT %s FROM %s WHERE record_id = $1`, expression, table), record).Scan(&value))
		if !value.Valid {
			return "<NULL>"
		}
		return value.String
	}
	require.Equal(t, "reset_or_uncertain", text("reset-rec-1", "status"))
	require.Equal(t, "110", text("reset-rec-1", "value_decimal::text"))
	require.Equal(t, "<NULL>", text("reset-rec-1", "delta_value::text"))
	require.Equal(t, "decrease", text("reset-rec-1", "provenance->'counter'->'events'->0->>'kind'"))
	require.Equal(t, "100", text("reset-rec-1", "provenance->'counter'->'events'->0->>'from_value'"))
	require.Equal(t, "10", text("reset-rec-1", "provenance->'counter'->'events'->0->>'to_value'"))
	require.Equal(t, "15", text("reset-rec-2", "delta_value::text"))

	require.Equal(t, "-3", text("signed-rec-1", "delta_value::text"), "NUMERIC keeps the negative exact delta")
	require.Equal(t, "7", text("signed-rec-1", "value_int::text"))
	require.Equal(t, "10", text("rollover-rec-1", "delta_value::text"))
	require.Equal(t, "rollover", text("rollover-rec-1", "provenance->'counter'->'events'->0->>'kind'"))
	require.Equal(t, "65536", text("rollover-rec-1", "provenance->'counter'->>'rollover_modulus'"))

	require.Equal(t, "missing", text("gap-rec-1", "status"))
	require.Equal(t, "<NULL>", text("gap-rec-1", "delta_value::text"))
	require.Equal(t, "gap", text("gap-rec-2", "status"))
	require.Equal(t, "2", text("gap-rec-2", "delta_value::text"))
	require.Equal(t, "2026-10-24 10:00:00", text("gap-rec-2", "(delta_from AT TIME ZONE 'UTC')::text"))
	require.Equal(t, "2026-10-24 10:20:00", text("gap-rec-2", "(delta_to AT TIME ZONE 'UTC')::text"))
	require.Equal(t, "true", text("gap-rec-2", "provenance->'counter'->>'span_total'"))

	require.Equal(t, "16", text("flapping-rec-1", "jsonb_array_length(provenance->'counter'->'events')::text"))
	require.Equal(t, "4", text("flapping-rec-1", "provenance->'counter'->>'events_omitted'"))
}
