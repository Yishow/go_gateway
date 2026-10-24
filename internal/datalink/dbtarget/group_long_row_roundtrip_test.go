package dbtarget

import (
	"database/sql"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// insertLongRows commits every case as one effect and reads each back with its
// SQLite storage class, proving the encoded cells land in the right typed slot.
func TestLongRowSQLiteRealCatalogInsertAndReadback(t *testing.T) {
	cases := longCases(t)
	statements := longDDL(t, "sqlite")
	target := newInspectionSQLiteTarget(t, statements...)
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", indexTestTable)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, inspection.Status)
	layout, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, inspection.Columns, cases))
	require.Empty(t, issues, "the real catalog satisfies the fixed layout")

	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	for _, c := range cases {
		row := mustEncode(t, layout, c.member, snapshot.MemberOK, "", longSample("s-"+c.member.MemberKey, c.value))
		result, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{
			Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: indexTestTable, Row: row,
			PayloadDigest: "digest-" + c.member.MemberKey, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
		})
		require.NoError(t, err, c.member.MemberKey)
		require.False(t, result.AlreadyCommitted)
	}
	missing := mustEncode(t, layout, cases[0].member, snapshot.MemberMissing, "no-samples", nil)
	missing.RecordID, missing.EffectKey = "rec-missing", "eff-missing"
	for i := range missing.Cells {
		switch missing.Cells[i].Column {
		case "record_id":
			missing.Cells[i].Value = "rec-missing"
		case "tag_id":
			missing.Cells[i].Value = "t-missing"
		}
	}
	_, err = InsertGroupRow(t.Context(), db, GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: indexTestTable, Row: missing,
		PayloadDigest: "digest-missing", Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
	})
	require.NoError(t, err)

	type stored struct {
		class, text string
	}
	read := func(tag, column string) stored {
		var s stored
		var text sql.NullString
		require.NoError(t, db.QueryRowContext(t.Context(),
			fmt.Sprintf(`SELECT typeof(%[1]s), CAST(%[1]s AS TEXT) FROM %[2]s WHERE tag_id = ?`, column, indexTestTable), tag).Scan(&s.class, &text))
		s.text = text.String
		return s
	}
	require.Equal(t, stored{"integer", "9007199254740993"}, read("t-int", "value_int"))
	require.Equal(t, stored{"text", "18446744073709551615"}, read("t-uint", "value_decimal"))
	require.Equal(t, stored{"text", "12000.4"}, read("t-dec", "value_decimal"))
	require.Equal(t, stored{"real", "0.1"}, read("t-float", "value_float"))
	require.Equal(t, stored{"integer", "1"}, read("t-bool", "value_bool"))
	require.Equal(t, stored{"text", "run-42"}, read("t-text", "value_text"))
	require.Equal(t, stored{"text", "2026-10-24T00:00:29.123456Z"}, read("t-int", "observed_at"))
	require.Equal(t, stored{"null", ""}, read("t-missing", "value_int"))
	require.Equal(t, stored{"text", "missing"}, read("t-missing", "status"))
	require.Equal(t, stored{"null", ""}, read("t-missing", "observed_at"))
	require.Equal(t, stored{"text", "gateway"}, read("t-int", indexTestOwner), "the row carries its owner marker")

	// A resent effect is recognized by its receipt and writes nothing again.
	row := mustEncode(t, layout, cases[0].member, snapshot.MemberOK, "", longSample("s-"+cases[0].member.MemberKey, cases[0].value))
	again, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: indexTestTable, Row: row,
		PayloadDigest: "digest-" + cases[0].member.MemberKey, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:31Z",
	})
	require.NoError(t, err)
	require.True(t, again.AlreadyCommitted)
	var rows int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM `+indexTestTable).Scan(&rows))
	require.Equal(t, len(cases)+1, rows)
}

func TestLongRowPostgresRealCatalogInsertAndReadback(t *testing.T) {
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

	cases := longCases(t)
	layout, issues := NewLongRowLayout(longSpec(SQLDialectPostgres, inspection.Columns, cases))
	require.Empty(t, issues, "the real PostgreSQL catalog satisfies the fixed layout")

	for _, c := range cases {
		row := mustEncode(t, layout, c.member, snapshot.MemberOK, "", longSample("s-"+c.member.MemberKey, c.value))
		_, err := InsertGroupRow(t.Context(), admin, GroupInsertRequest{
			Kind: schema.DatabaseConnectorKindPostgres, SchemaName: schemaName, TableName: indexTestTable, Row: row,
			PayloadDigest: "digest-" + c.member.MemberKey, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
		})
		require.NoError(t, err, c.member.MemberKey)
	}
	missing := mustEncode(t, layout, cases[0].member, snapshot.MemberMissing, "no-samples", nil)
	missing.EffectKey = "eff-missing"
	for i := range missing.Cells {
		if missing.Cells[i].Column == "record_id" {
			missing.Cells[i].Value = "rec-missing"
		}
		if missing.Cells[i].Column == "tag_id" {
			missing.Cells[i].Value = "t-missing"
		}
	}
	_, err = InsertGroupRow(t.Context(), admin, GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindPostgres, SchemaName: schemaName, TableName: indexTestTable, Row: missing,
		PayloadDigest: "digest-missing", Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T00:00:30Z",
	})
	require.NoError(t, err)

	table := pgIdentifier(t, schemaName) + "." + pgIdentifier(t, indexTestTable)
	text := func(tag, expression string) string {
		var value sql.NullString
		require.NoError(t, admin.QueryRowContext(t.Context(), fmt.Sprintf(`SELECT %s FROM %s WHERE tag_id = $1`, expression, table), tag).Scan(&value))
		if !value.Valid {
			return "<NULL>"
		}
		return value.String
	}
	require.Equal(t, "9007199254740993", text("t-int", "value_int::text"))
	require.Equal(t, "18446744073709551615", text("t-uint", "value_decimal::text"), "uint64 is preserved by NUMERIC")
	require.Equal(t, "12000.4", text("t-dec", "value_decimal::text"))
	require.Equal(t, "0.1", text("t-float", "value_float::text"))
	require.Equal(t, "true", text("t-bool", "value_bool::text"))
	require.Equal(t, "run-42", text("t-text", "value_text"))
	require.Equal(t, "2026-10-24 00:00:29.123456+00", text("t-int", "(observed_at AT TIME ZONE 'UTC')::text||'+00'"))
	require.Equal(t, "instantaneous", text("t-int", "provenance->'profile'->>'kind'"))
	require.Equal(t, "<NULL>", text("t-missing", "value_int::text"))
	require.Equal(t, "missing", text("t-missing", "status"))
	require.Equal(t, "<NULL>", text("t-missing", "observed_at::text"))
}
