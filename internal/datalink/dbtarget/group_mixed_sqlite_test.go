package dbtarget

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"testing"

	"github.com/stretchr/testify/require"
)

func mixedSparseOutcomes(t *testing.T, mode snapshot.Mode) []snapshot.Outcome {
	t.Helper()
	c := mixedEncodeConfig(mode)
	c.StartsAt = mixedEncodeAt(120)
	a, err := snapshot.NewAssembler(c)
	require.NoError(t, err)
	mixedEncodeOffer(t, a, "a", 125, 23)
	mixedEncodeOffer(t, a, "b", 125, 13)
	p := a.Plan(mixedEncodeAt(135))
	require.NoError(t, p.Err)
	return p.Outcomes
}
func mixedInsert(t *testing.T, db *sql.DB, table string, row EncodedRow) {
	t.Helper()
	// The managed long table has a composite identity, not a record_id-only UNIQUE key.
	strategy := GroupEffectUniqueKey
	if table == indexTestTable {
		strategy = GroupEffectNone
	}
	r, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{Kind: schema.DatabaseConnectorKindSQLite, SchemaName: "main", TableName: table, Row: row, Strategy: strategy, RecordKeyColumn: "record_id"})
	var insertErr *GroupInsertError
	if errors.As(err, &insertErr) {
		t.Logf("insert cause: %v", insertErr.Cause)
	}
	require.NoError(t, err)
	require.False(t, r.AlreadyCommitted)
}

func TestMixedSQLiteWideActualCatalogInsertReadback(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE mixed_results(record_id TEXT PRIMARY KEY NOT NULL,due_at TEXT NOT NULL,provenance TEXT NOT NULL,a TEXT,b TEXT,c TEXT)`)
	f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspection, err := f.inspector.InspectTable(t.Context(), f.connector.ID, "main", "mixed_results")
	require.NoError(t, err)
	s := mixedWideSpec(SQLDialectSQLite)
	s.Columns = inspection.Columns
	l, issues := NewGroupRowLayout(s)
	require.Empty(t, issues)
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	outcomes := append(mixedEncodeFixture(t, snapshot.ModeMixedWide), mixedSparseOutcomes(t, snapshot.ModeMixedWide)...)
	for _, out := range outcomes {
		row, err := l.EncodeRow(out)
		require.NoError(t, err)
		mixedInsert(t, db, "mixed_results", row)
	}
	var a, b, c, class, provenance string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT a,b,c,typeof(c),provenance FROM mixed_results WHERE due_at=?`, mixedEncodeAt(120).Format("2006-01-02T15:04:05.999999999Z07:00")).Scan(&a, &b, &c, &class, &provenance))
	require.Equal(t, []string{"22", "12", "8", "text"}, []string{a, b, c, class})
	var nullable sql.NullString
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT a,b,c,typeof(c),provenance FROM mixed_results WHERE due_at=?`, mixedEncodeAt(135).Format("2006-01-02T15:04:05.999999999Z07:00")).Scan(&a, &b, &nullable, &class, &provenance))
	require.Equal(t, "23", a)
	require.Equal(t, "13", b)
	require.False(t, nullable.Valid)
	require.Equal(t, "null", class)
	var entries []map[string]any
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	require.Equal(t, "not_due", entries[2]["reason"])
	require.Nil(t, entries[2]["raw"])
	require.Nil(t, entries[2]["window"])
}
func TestMixedSQLiteLongActualCatalogDueOnlyReadback(t *testing.T) {
	target := newInspectionSQLiteTarget(t, longDDL(t, "sqlite")...)
	f := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": target})
	inspection, err := f.inspector.InspectTable(t.Context(), f.connector.ID, "main", indexTestTable)
	require.NoError(t, err)
	s := mixedLongSpec(t, SQLDialectSQLite)
	s.Columns = inspection.Columns
	l, issues := NewLongRowLayout(s)
	require.Empty(t, issues)
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	full := mixedEncodeFixture(t, snapshot.ModeMixedLong)
	sparse := mixedSparseOutcomes(t, snapshot.ModeMixedLong)
	require.Len(t, sparse, 2)
	for _, out := range append(full, sparse...) {
		row, err := l.EncodeRow(out)
		require.NoError(t, err)
		mixedInsert(t, db, indexTestTable, row)
	}
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+indexTestTable).Scan(&count))
	require.Equal(t, 5, count)
	var value, class, profile, rawJSON string
	var oldDelta sql.NullString
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT value_decimal,typeof(value_decimal),profile_kind,delta_value,provenance FROM `+indexTestTable+` WHERE tag_id='tag-c'`).Scan(&value, &class, &profile, &oldDelta, &rawJSON))
	require.Equal(t, "8", value)
	require.Equal(t, "text", class)
	require.Equal(t, "mixed", profile)
	require.False(t, oldDelta.Valid)
	var e map[string]any
	require.NoError(t, json.Unmarshal([]byte(rawJSON), &e))
	require.Equal(t, "1008", e["raw"].(map[string]any)["value"])
	require.Equal(t, "8", e["counter"].(map[string]any)["delta"])
}

func TestMixedSQLiteExactDecimalReadback(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE mixed_results(record_id TEXT PRIMARY KEY NOT NULL,due_at TEXT NOT NULL,provenance TEXT NOT NULL,a TEXT,b TEXT,c TEXT)`)
	db, err := sql.Open("sqlite", target)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	layout, issues := NewGroupRowLayout(mixedWideSpec(SQLDialectSQLite))
	require.Empty(t, issues)
	row, err := layout.EncodeRow(mixedExactOutcome(t, snapshot.ModeMixedWide)[0])
	require.NoError(t, err)
	mixedInsert(t, db, "mixed_results", row)
	var a, b, c, typ string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT a,b,c,typeof(a) FROM mixed_results`).Scan(&a, &b, &c, &typ))
	require.Equal(t, []string{"9007199254740994", "1.5", "1", "text"}, []string{a, b, c, typ})
}
