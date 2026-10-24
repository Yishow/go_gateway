package dbtarget

import (
	"database/sql"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// TestLongRowCounterPostgresReadback commits cumulative records into the real
// PostgreSQL long table and reads each stored value back as text, proving NUMERIC
// keeps the source's digits and the exact deltas for decimal and large integers.
// It needs POSTGRES_DSN like the other live PostgreSQL fixtures.
func TestLongRowCounterPostgresReadback(t *testing.T) {
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

	insert := func(t *testing.T, member LongRowMember, prefix string, outcomes []snapshot.Outcome) {
		t.Helper()
		layout, issues := NewLongRowLayout(counterSpecFor(SQLDialectPostgres, inspection.Columns, member))
		require.Empty(t, issues, "the real PostgreSQL catalog satisfies the fixed layout")
		for _, outcome := range outcomes {
			outcome.RecordID, outcome.EffectKey = prefix+outcome.RecordID, prefix+outcome.EffectKey
			row := outcomesEncoded(t, layout, outcome)
			_, err := InsertGroupRow(t.Context(), admin, GroupInsertRequest{
				Kind: schema.DatabaseConnectorKindPostgres, SchemaName: schemaName, TableName: indexTestTable, Row: row,
				PayloadDigest: "digest-" + row.RecordID, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
			})
			require.NoError(t, err, row.RecordID)
		}
	}
	decimal := counterMemberFor(measurement.ExactDecimal)
	insert(t, decimal, "dec-", counterOutcomes(t, decimal, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), exactDecimal(t, "12000.7")))
	large := counterMemberFor(measurement.ExactUint64)
	insert(t, large, "big-", counterOutcomes(t, large, exactUint(9007199254740992), exactUint(9007199254740993)))
	signedInt := counterMemberFor(measurement.ExactInt64)
	first, second := measurement.NewInt64(9007199254740992), measurement.NewInt64(9007199254740993)
	insert(t, signedInt, "int-", counterOutcomes(t, signedInt, &first, &second))
	unknown := counterOutcomes(t, large, exactUint(100), exactUint(10), nil)
	insert(t, large, "unk-", unknown)

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
	require.Equal(t, "12000.0", text("dec-rec-0", "value_decimal::text"), "NUMERIC keeps the scale of the source digits")
	require.Equal(t, "<NULL>", text("dec-rec-0", "delta_value::text"))
	require.Equal(t, "baseline", text("dec-rec-0", "status"))
	require.Equal(t, "12000.4", text("dec-rec-1", "value_decimal::text"))
	require.Equal(t, "0.4", text("dec-rec-1", "delta_value::text"))
	require.Equal(t, "0.3", text("dec-rec-2", "delta_value::text"))
	require.Equal(t, "12000.4", text("dec-rec-1", "provenance->'counter'->>'raw_text'"))
	require.Equal(t, "2026-10-24 10:09:00.123456", text("dec-rec-1", "(delta_from AT TIME ZONE 'UTC')::text"))
	require.Equal(t, "2026-10-24 10:19:00.123456", text("dec-rec-1", "(delta_to AT TIME ZONE 'UTC')::text"))

	require.Equal(t, "9007199254740993", text("big-rec-1", "value_decimal::text"))
	require.Equal(t, "1", text("big-rec-1", "delta_value::text"), "uint64 above 2^53 differs by exactly 1")
	require.Equal(t, "9007199254740993", text("int-rec-1", "value_int::text"))
	require.Equal(t, "1", text("int-rec-1", "delta_value::text"))
	require.Equal(t, "epoch-1", text("int-rec-1", "series_epoch"))
	require.Equal(t, "kWh", text("int-rec-1", "unit"))

	require.Equal(t, "reset_or_uncertain", text("unk-rec-1", "status"))
	require.Equal(t, "10", text("unk-rec-1", "value_decimal::text"), "an unexplained decrease keeps its raw reading")
	require.Equal(t, "<NULL>", text("unk-rec-1", "delta_value::text"), "and has no delta, not an absolute value or zero")
	require.Equal(t, "missing", text("unk-rec-2", "status"))
	require.Equal(t, "<NULL>", text("unk-rec-2", "value_decimal::text"))
	require.Equal(t, "<NULL>", text("unk-rec-2", "delta_value::text"))
}
