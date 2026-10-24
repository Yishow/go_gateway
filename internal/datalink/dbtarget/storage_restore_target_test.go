package dbtarget

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func openRestoreSQLite(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "restore.db")+"?_pragma=busy_timeout(5000)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// TestSchemaShapeErrorsAreTheOnlyOnesThatLookAtTheCatalog: only a missing
// relation or column, or a column of the wrong type, can start a look at the
// destination catalog. Privileges, credentials, constraints, locks and a lost
// connection never do.
func TestSchemaShapeErrorsAreTheOnlyOnesThatLookAtTheCatalog(t *testing.T) {
	for code, want := range map[string]bool{
		"42P01": true, "42703": true, "42804": true,
		"42501": false, "28P01": false, "3D000": false, "3F000": false, "23505": false, "23502": false, "22P02": false,
		"40001": false, "40P01": false, "53300": false, "57P01": false, "08006": false, "42601": false, "42P07": false,
	} {
		err := fmt.Errorf("insert destination row: %w", &pgconn.PgError{Code: code, Message: "secret detail"})
		require.Equal(t, want, IsSchemaShapeError(err), code)
	}
	require.False(t, IsSchemaShapeError(nil))
	require.False(t, IsSchemaShapeError(errors.New("no such table: looks like text but is not a driver error")))
	require.False(t, IsSchemaShapeError(&GroupInsertError{Phase: InsertPhaseCommit, Cause: errors.New("lost")}))

	db := openRestoreSQLite(t)
	_, err := db.ExecContext(t.Context(), `CREATE TABLE readings (id TEXT PRIMARY KEY, counter INTEGER NOT NULL CHECK (counter >= 0))`)
	require.NoError(t, err)
	for name, tc := range map[string]struct {
		statement string
		want      bool
	}{
		"a missing table":      {`INSERT INTO absent (id) VALUES ('a')`, true},
		"a missing column":     {`INSERT INTO readings (id, nope) VALUES ('a', 1)`, true},
		"a duplicate key":      {`INSERT INTO readings (id, counter) VALUES ('a', 1), ('a', 2)`, false},
		"a check violation":    {`INSERT INTO readings (id, counter) VALUES ('b', -1)`, false},
		"a not null violation": {`INSERT INTO readings (id, counter) VALUES ('c', NULL)`, false},
		"a syntax error":       {`INSERT INTO readings (`, false},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := db.ExecContext(t.Context(), tc.statement)
			require.Error(t, err)
			require.Equal(t, tc.want, IsSchemaShapeError(err), err.Error())
		})
	}
}

func TestInspectOpenManagedTableOnSQLite(t *testing.T) {
	db := openRestoreSQLite(t)
	ctx := t.Context()
	kind := schema.DatabaseConnectorKindSQLite

	missing, err := InspectOpenManagedTable(ctx, db, kind, "main", "gw_long_x", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionMissing, missing.Status)

	for _, statement := range []string{
		`CREATE TABLE gw_long_x (record_id TEXT NOT NULL, bucket_start TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`,
		`CREATE INDEX idx_gw_long_x ON gw_long_x (record_id, bucket_start)`,
		`CREATE TABLE source (id TEXT)`,
		`CREATE VIEW as_view AS SELECT id FROM source`,
	} {
		_, err := db.ExecContext(ctx, statement)
		require.NoError(t, err, statement)
	}
	table, err := InspectOpenManagedTable(ctx, db, kind, "main", "gw_long_x", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, table.Status)
	require.Len(t, table.Columns, 2)
	require.True(t, table.IndexesInspected)
	require.Equal(t, []IndexInfo{{Name: "idx_gw_long_x", Columns: []string{"record_id", "bucket_start"}}}, table.Indexes)

	shallow, err := InspectOpenManagedTable(ctx, db, kind, "main", "gw_long_x", false)
	require.NoError(t, err)
	require.False(t, shallow.IndexesInspected, "indexes are read only when asked for")

	view, err := InspectOpenManagedTable(ctx, db, kind, "main", "as_view", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, view.Status, "a name held by a view is occupied, not free")
	require.Empty(t, view.Columns, "and it carries no managed columns")

	index, err := InspectOpenManagedTable(ctx, db, kind, "main", "idx_gw_long_x", false)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, index.Status, "an index name is occupied too")

	other, err := InspectOpenManagedTable(ctx, db, schema.DatabaseConnectorKindMySQL, "main", "x", false)
	require.NoError(t, err)
	require.Equal(t, TableInspectionFailed, other.Status, "only SQLite and PostgreSQL have a managed destination")
}

func TestRunManagedSchemaBatchIsAllOrNothing(t *testing.T) {
	db := openRestoreSQLite(t)
	kind := schema.DatabaseConnectorKindSQLite

	execution, err := RunManagedSchemaBatch(t.Context(), db, kind, "main", []string{
		`CREATE TABLE kept (id TEXT)`, `CREATE INDEX idx_kept ON kept (id)`,
	})
	require.NoError(t, err)
	require.Equal(t, 2, execution.Committed)

	execution, err = RunManagedSchemaBatch(t.Context(), db, kind, "main", []string{
		`CREATE TABLE half (id TEXT)`, `CREATE TABLE kept (id TEXT)`,
	})
	require.Error(t, err)
	require.True(t, execution.RolledBack, "a failed batch is rolled back whole")
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sqlite_master WHERE name = 'half'`).Scan(&count))
	require.Zero(t, count, "no statement of the failed batch stays applied")

	execution, err = RunManagedSchemaBatch(t.Context(), db, schema.DatabaseConnectorKindMySQL, "main", []string{`CREATE TABLE nope (id TEXT)`})
	require.ErrorIs(t, err, ErrSchemaExecutionUnsupported)
	require.True(t, execution.RolledBack)

	ok, err := CanCreateManagedTable(t.Context(), db, kind, "main")
	require.NoError(t, err)
	require.True(t, ok, "SQLite has no per-role privilege to check")
}

func TestStorageRestoreTargetOnPostgres(t *testing.T) {
	dsn, _ := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	namespace := postgresScopedSchema(t, admin)
	ctx := t.Context()
	kind := schema.DatabaseConnectorKindPostgres
	quoted := pgIdentifier(t, namespace)

	missing, err := InspectOpenManagedTable(ctx, admin, kind, namespace, "gw_long_x", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionMissing, missing.Status)

	execution, err := RunManagedSchemaBatch(ctx, admin, kind, namespace, []string{
		`CREATE TABLE gw_long_x (record_id TEXT NOT NULL, bucket_start TIMESTAMPTZ NOT NULL, PRIMARY KEY (bucket_start, record_id))`,
		`CREATE INDEX idx_gw_long_x ON gw_long_x (record_id, bucket_start)`,
	})
	require.NoError(t, err)
	require.Equal(t, 2, execution.Committed)
	_, err = admin.ExecContext(ctx, `CREATE VIEW `+quoted+`.as_view AS SELECT record_id FROM `+quoted+`.gw_long_x`)
	require.NoError(t, err)

	table, err := InspectOpenManagedTable(ctx, admin, kind, namespace, "gw_long_x", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, table.Status)
	require.Len(t, table.Columns, 2)
	require.Equal(t, []IndexInfo{{Name: "idx_gw_long_x", Columns: []string{"record_id", "bucket_start"}}}, table.Indexes)

	view, err := InspectOpenManagedTable(ctx, admin, kind, namespace, "as_view", true)
	require.NoError(t, err)
	require.Equal(t, TableInspectionExists, view.Status, "a view holding the name is occupied")
	require.Empty(t, view.Columns)

	allowed, err := CanCreateManagedTable(ctx, admin, kind, namespace)
	require.NoError(t, err)
	require.True(t, allowed)
	allowed, err = CanCreateManagedTable(ctx, admin, kind, "no_such_schema")
	require.NoError(t, err)
	require.False(t, allowed, "a schema that does not exist cannot be created in")
}
