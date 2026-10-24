package dbtarget

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// The owner column is part of each long record. Writing must not depend on a
// column default that the catalog layout proof does not require.
func TestLongRowWritesItsOwnerWithoutAColumnDefault(t *testing.T) {
	for _, dialect := range []SQLDialect{SQLDialectSQLite, SQLDialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			statements := longDDL(t, string(dialect))
			for i := range statements {
				if strings.Contains(statements[i], indexTestOwner) {
					statements[i] = strings.ReplaceAll(statements[i], " DEFAULT 'gateway'", "")
				}
			}
			db, namespace, inspection := longOwnerTarget(t, dialect, statements)
			for _, column := range inspection.Columns {
				if column.Name == indexTestOwner {
					require.False(t, column.HasDefault, "the real catalog has no owner default")
				}
			}
			cases := longCases(t)[:1]
			layout, issues := NewLongRowLayout(longSpec(dialect, inspection.Columns, cases))
			require.Empty(t, issues)
			row := mustEncode(t, layout, cases[0].member, snapshot.MemberOK, "", longSample("owner-sample", cases[0].value))
			kind := schema.DatabaseConnectorKindSQLite
			table := `"` + indexTestTable + `"`
			if dialect == SQLDialectPostgres {
				kind = schema.DatabaseConnectorKindPostgres
				table = pgIdentifier(t, namespace) + "." + pgIdentifier(t, indexTestTable)
			}
			result, err := InsertGroupRow(t.Context(), db, GroupInsertRequest{
				Kind: kind, SchemaName: namespace, TableName: indexTestTable, Row: row,
				PayloadDigest: "owner-payload", Strategy: GroupEffectReceipt, CommittedAt: "2026-10-25T00:00:30Z",
			})
			require.NoError(t, err, "a verified long layout remains writable without an owner default: %v", errors.Unwrap(err))
			require.False(t, result.AlreadyCommitted)
			var owner, raw string
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT "`+indexTestOwner+`", CAST(value_int AS TEXT) FROM `+table).Scan(&owner, &raw))
			require.Equal(t, "gateway", owner)
			require.Equal(t, "9007199254740993", raw)
		})
	}
}

func longOwnerTarget(t *testing.T, dialect SQLDialect, statements []string) (*sql.DB, string, *TableInspection) {
	t.Helper()
	if dialect == SQLDialectSQLite {
		path := newInspectionSQLiteTarget(t, statements...)
		fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": path})
		inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, "main", indexTestTable)
		require.NoError(t, err)
		db, err := sql.Open("sqlite", path)
		require.NoError(t, err)
		t.Cleanup(func() { _ = db.Close() })
		return db, "main", inspection
	}
	dsn, fields := postgresKeywordDSN(t)
	db := postgresAdminDB(t, dsn)
	namespace := postgresScopedSchema(t, db)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer func() { _ = tx.Rollback() }()
	_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO `+pgIdentifier(t, namespace))
	require.NoError(t, err)
	for _, statement := range statements {
		_, err := tx.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	require.NoError(t, tx.Commit())
	fixture := newManagedDestinationFixture(t, schema.DatabaseConnectorKindPostgres,
		postgresConnectorConfig(fields, namespace, fields["user"], fields["password"]))
	inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, namespace, indexTestTable)
	require.NoError(t, err)
	return db, namespace, inspection
}
