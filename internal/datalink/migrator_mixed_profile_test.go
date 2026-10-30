package datalink

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema/migrations"
	"path/filepath"
	"sync"
	"testing"
)

// Build a genuine pre-expansion table from immutable historical migrations,
// including external references and user-owned schema objects.
func legacyMixedStorageDB(t *testing.T, paths ...string) *sql.DB {
	t.Helper()
	dsn := ":memory:"
	if len(paths) > 0 {
		dsn = paths[0]
	}
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	for _, name := range []string{"024_write_groups_sqlite.up.sql", "029_write_group_circuit_profiles_sqlite.up.sql"} {
		raw, err := migrations.FS.ReadFile(name)
		require.NoError(t, err)
		_, err = db.Exec(string(raw))
		require.NoError(t, err)
	}
	_, err = db.Exec(`PRAGMA foreign_keys=ON;
 INSERT INTO write_groups(id,workspace_id,revision,name,destination_connector_id,destination_connector_revision,destination_table_name,profile_kind,profile_interval,profile_unit)
 VALUES('legacy','ws','revision-exact','Legacy','connector','connector-revision','history','instantaneous',15,'second');
 INSERT INTO write_group_members(group_id,member_index,device_id,point_id,tag_id,target_column) VALUES('legacy',0,'device','point','tag','value');
 CREATE INDEX preserve_name ON write_groups(name);
 CREATE TABLE audit_test (value TEXT);
 CREATE TRIGGER preserve_trigger AFTER UPDATE OF name ON write_groups BEGIN INSERT INTO audit_test VALUES(new.name); END;
 CREATE VIEW preserve_view AS SELECT id,revision FROM write_groups;`)
	require.NoError(t, err)
	return db
}

func TestMixedProfileMigrationPreservesSchemaAndValues(t *testing.T) {
	db := legacyMixedStorageDB(t)
	before := mixedStorageRow(t, db)
	require.NoError(t, ensureSQLiteMemberCalculationIntent(db))
	after := mixedStorageRow(t, db)
	for key, value := range before {
		require.Equal(t, value, after[key], key)
	}
	// The same canonical table must admit mixed; migration cannot route storage
	// to a shadow authority, rewrite the legacy kind, or drop referenced rows.
	_, err := db.Exec(`UPDATE write_groups SET profile_kind='mixed',profile_interval=NULL,profile_unit=NULL,result_shape='member-calculation-v1' WHERE id='legacy'`)
	require.NoError(t, err)
	var count, fk int
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM write_group_members WHERE group_id='legacy'`).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk))
	require.Equal(t, 1, fk)
	var revision string
	require.NoError(t, db.QueryRow(`SELECT revision FROM preserve_view WHERE id='legacy'`).Scan(&revision))
	require.Equal(t, "revision-exact", revision)
	_, err = db.Exec(`UPDATE write_groups SET name='After' WHERE id='legacy'`)
	require.NoError(t, err)
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM audit_test WHERE value='After'`).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name='preserve_name'`).Scan(&count))
	require.Equal(t, 1, count)
	_, err = db.Exec(`DELETE FROM write_groups WHERE id='legacy'`)
	require.Error(t, err, "foreign key still restricts delete")
	require.NoError(t, ensureSQLiteMemberCalculationIntent(db))
}

func mixedStorageRow(t *testing.T, db *sql.DB) map[string]any {
	t.Helper()
	rows, err := db.Query("SELECT * FROM write_groups WHERE id='legacy'")
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	require.True(t, rows.Next())
	values := make([]any, len(columns))
	targets := make([]any, len(columns))
	for i := range targets {
		targets[i] = &values[i]
	}
	require.NoError(t, rows.Scan(targets...))
	result := make(map[string]any, len(columns))
	for i, column := range columns {
		result[column] = values[i]
	}
	require.False(t, rows.Next())
	require.NoError(t, rows.Err())
	return result
}

func TestMixedProfileMigrationRollsBackEveryStage(t *testing.T) {
	for _, stage := range []string{"copy", "replace", "indexes", "foreign_keys", "commit"} {
		for _, fk := range []int{0, 1} {
			t.Run(fmt.Sprintf("%s/fk%d", stage, fk), func(t *testing.T) {
				db := legacyMixedStorageDB(t)
				_, err := db.Exec(fmt.Sprintf("PRAGMA foreign_keys=%d", fk))
				require.NoError(t, err)
				before := mixedStorageRow(t, db)
				var schemaBefore string
				require.NoError(t, db.QueryRow(`SELECT sql FROM sqlite_master WHERE name='write_groups'`).Scan(&schemaBefore))
				injected := errors.New("migration injected failure")
				err = expandSQLiteMixedProfile(db, func(at string) error {
					if at == stage {
						return injected
					}
					return nil
				})
				require.ErrorIs(t, err, injected)
				require.Equal(t, before, mixedStorageRow(t, db))
				var schemaAfter string
				require.NoError(t, db.QueryRow(`SELECT sql FROM sqlite_master WHERE name='write_groups'`).Scan(&schemaAfter))
				require.Equal(t, schemaBefore, schemaAfter)
				var actual, count int
				require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&actual))
				require.Equal(t, fk, actual)
				require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE name='write_groups_member_calculation_expand'`).Scan(&count))
				require.Zero(t, count)
				require.NoError(t, ensureSQLiteMemberCalculationIntent(db), "a failed migration remains safely retryable")
			})
		}
	}
}

func TestMixedProfileMigrationRefusesUnknownSchemaAndForeignKeys(t *testing.T) {
	t.Run("unknown column", func(t *testing.T) {
		db := legacyMixedStorageDB(t)
		_, err := db.Exec(`ALTER TABLE write_groups ADD COLUMN unrecognized TEXT`)
		require.NoError(t, err)
		require.ErrorContains(t, ensureSQLiteMemberCalculationIntent(db), "unknown write_groups column")
		var count int
		require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('write_groups') WHERE name='result_shape'`).Scan(&count))
		require.Zero(t, count)
	})
	t.Run("preexisting foreign key violation is not repaired", func(t *testing.T) {
		db := legacyMixedStorageDB(t)
		_, err := db.Exec(`PRAGMA foreign_keys=OFF; UPDATE write_group_members SET group_id='missing'`)
		require.NoError(t, err)
		require.ErrorContains(t, ensureSQLiteMemberCalculationIntent(db), "foreign key check failed")
		var fk int
		require.NoError(t, db.QueryRow(`PRAGMA foreign_keys`).Scan(&fk))
		require.Zero(t, fk)
	})
}

func TestMixedProfileMigrationConcurrentStartup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "concurrent.db")
	seed := legacyMixedStorageDB(t, path)
	require.NoError(t, seed.Close())
	const workers = 4
	var wait sync.WaitGroup
	results := make(chan error, workers)
	for i := 0; i < workers; i++ {
		db, err := sql.Open("sqlite", path)
		require.NoError(t, err)
		db.SetMaxOpenConns(1)
		_, err = db.Exec("PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON")
		require.NoError(t, err)
		wait.Add(1)
		go func() { defer wait.Done(); defer db.Close(); results <- ensureSQLiteMemberCalculationIntent(db) }()
	}
	wait.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	var kind string
	require.NoError(t, db.QueryRow("SELECT profile_kind FROM write_groups WHERE id='legacy'").Scan(&kind))
	require.Equal(t, "instantaneous", kind)
}
