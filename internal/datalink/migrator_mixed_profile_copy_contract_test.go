package datalink

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func migrationCopyContractRows(t *testing.T, db *sql.DB, query string) [][]any {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), query)
	require.NoError(t, err)
	defer func() { require.NoError(t, rows.Close()) }()
	columns, err := rows.Columns()
	require.NoError(t, err)
	result := make([][]any, 0, 2)
	for rows.Next() {
		values := make([]any, len(columns))
		targets := make([]any, len(columns))
		for i := range values {
			targets[i] = &values[i]
		}
		require.NoError(t, rows.Scan(targets...))
		result = append(result, values)
	}
	require.NoError(t, rows.Err())
	return result
}

func TestMixedProfileFixedCopyMatchesOriginalExplicitColumns(t *testing.T) {
	for _, fixture := range []struct{ name, ddl, copySQL string }{
		{name: "minimal historical optional subset", ddl: `CREATE TABLE write_groups (
id TEXT PRIMARY KEY,
workspace_id TEXT,
revision TEXT,
applied_revision TEXT,
name BLOB,
status TEXT,
destination_connector_id TEXT,
destination_connector_revision TEXT,
destination_database TEXT,
destination_table_schema TEXT,
destination_table_name TEXT,
destination_storage_strategy TEXT,
destination_schema_revision TEXT,
destination_schema_digest TEXT,
row_policy TEXT,
write_policy TEXT,
migration BLOB,
created_at INTEGER,
updated_at REAL,
profile_kind TEXT CHECK (profile_kind IS NULL OR profile_kind IN ('instantaneous', 'cumulative'))
)`, copySQL: `INSERT INTO copy_reference ("id","workspace_id","revision","applied_revision","name","status","destination_connector_id","destination_connector_revision","destination_database","destination_table_schema","destination_table_name","destination_storage_strategy","destination_schema_revision","destination_schema_digest","row_policy","write_policy","migration","created_at","updated_at","profile_kind") SELECT "id","workspace_id","revision","applied_revision","name","status","destination_connector_id","destination_connector_revision","destination_database","destination_table_schema","destination_table_name","destination_storage_strategy","destination_schema_revision","destination_schema_digest","row_policy","write_policy","migration","created_at","updated_at","profile_kind" FROM write_groups`},
		{name: "reordered declaration and optional fields", ddl: `CREATE TABLE write_groups (
profile_interval INTEGER DEFAULT 15,
name BLOB,
id TEXT PRIMARY KEY,
result_shape TEXT,
workspace_id TEXT,
revision TEXT,
applied_revision TEXT,
status TEXT,
destination_connector_id TEXT,
destination_connector_revision TEXT,
destination_database TEXT,
destination_table_schema TEXT,
destination_table_name TEXT,
destination_storage_strategy TEXT,
destination_schema_revision TEXT,
destination_schema_digest TEXT,
row_policy TEXT,
write_policy TEXT,
migration BLOB,
created_at INTEGER,
updated_at REAL,
profile_kind TEXT CHECK (profile_kind IS NULL OR profile_kind IN ('instantaneous', 'cumulative')),
layout_version TEXT DEFAULT 'legacy-layout'
)`, copySQL: `INSERT INTO copy_reference ("profile_interval","name","id","result_shape","workspace_id","revision","applied_revision","status","destination_connector_id","destination_connector_revision","destination_database","destination_table_schema","destination_table_name","destination_storage_strategy","destination_schema_revision","destination_schema_digest","row_policy","write_policy","migration","created_at","updated_at","profile_kind","layout_version") SELECT "profile_interval","name","id","result_shape","workspace_id","revision","applied_revision","status","destination_connector_id","destination_connector_revision","destination_database","destination_table_schema","destination_table_name","destination_storage_strategy","destination_schema_revision","destination_schema_digest","row_policy","write_policy","migration","created_at","updated_at","profile_kind","layout_version" FROM write_groups`},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "copy.db"))
			require.NoError(t, err)
			defer func() { require.NoError(t, db.Close()) }()
			db.SetMaxOpenConns(1)
			_, err = db.ExecContext(t.Context(), fixture.ddl)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `INSERT INTO write_groups(rowid,id,workspace_id,revision,name,migration,created_at,updated_at,profile_kind)
   VALUES(40,'first','workspace','rev-a',?,NULL,'17',1.25,'instantaneous'),
   (80,'second','workspace','rev-b',NULL,?,NULL,2.5,'cumulative')`, []byte{0, 255, 1}, []byte{255, 0, 2})
			require.NoError(t, err)
			referenceDDL := strings.Replace(fixture.ddl, "CREATE TABLE write_groups (", "CREATE TABLE copy_reference (", 1)
			_, err = db.ExecContext(t.Context(), referenceDDL)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), fixture.copySQL)
			require.NoError(t, err)
			expected := migrationCopyContractRows(t, db, "SELECT rowid,typeof(name),typeof(migration),typeof(created_at),typeof(updated_at),* FROM copy_reference ORDER BY id")
			require.Equal(t, int64(1), expected[0][0], "reference is the original explicit INSERT behavior, not source-rowid preservation")
			require.NoError(t, runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error { return rebuildMixedProfileTable(ctx, tx, nil) }))
			actual := migrationCopyContractRows(t, db, "SELECT rowid,typeof(name),typeof(migration),typeof(created_at),typeof(updated_at),* FROM write_groups ORDER BY id")
			require.Equal(t, expected, actual, "same DDL column order, original explicit INSERT semantics, NULL/BLOB and affinity bytes")
			var count int
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM sqlite_master WHERE name='write_groups_member_calculation_expand'").Scan(&count))
			require.Zero(t, count)
		})
	}
}

func TestMixedProfileFixedCopyStillRejectsUnknownAndHiddenColumns(t *testing.T) {
	for _, tc := range []struct{ name, alter string }{
		{"unknown", `ALTER TABLE write_groups ADD COLUMN unexpected TEXT`},
		{"known but generated hidden", `ALTER TABLE write_groups ADD COLUMN result_shape TEXT GENERATED ALWAYS AS ('hidden') VIRTUAL`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := legacyMixedStorageDB(t)
			_, err := db.ExecContext(t.Context(), tc.alter)
			require.NoError(t, err)
			before := mixedStorageRow(t, db)
			var ddl string
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT sql FROM sqlite_master WHERE name='write_groups'").Scan(&ddl))
			err = runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error { return rebuildMixedProfileTable(ctx, tx, nil) })
			require.ErrorContains(t, err, "unknown write_groups column")
			require.Equal(t, before, mixedStorageRow(t, db))
			var after string
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT sql FROM sqlite_master WHERE name='write_groups'").Scan(&after))
			require.Equal(t, ddl, after)
		})
	}
}
