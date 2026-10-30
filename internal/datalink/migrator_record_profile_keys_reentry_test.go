package datalink

import (
	"context"
	"database/sql"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordProfileKeysForeignKeyFailureRollsBack(t *testing.T) {
	db := recordProfileDB(t)
	_, err := db.Exec(`PRAGMA foreign_keys=OFF; UPDATE write_group_profile_keys SET lineage_id='missing' WHERE group_id='G1'`)
	require.NoError(t, err)
	schema := recordProfileSchema(t, db)
	require.ErrorContains(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx), "foreign key check failed")
	require.Equal(t, schema, recordProfileSchema(t, db))
	var lineage string
	require.NoError(t, db.QueryRow(`SELECT lineage_id FROM write_group_profile_keys WHERE group_id='G1'`).Scan(&lineage))
	require.Equal(t, "missing", lineage, "no attempted repair")
}

func TestRecordProfileKeysConcurrentReentry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "concurrent.db")
	db := openLineageMigrationDB(t, path)
	createPre047LineageSchema(t, db)
	seedLegacyLineageRows(t, db)
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	columns := lineageSnapshotColumns(t, db, "write_group_profile_keys", "write_group_lineages", "write_group_generations")
	before := lineageMigrationSnapshot(t, db, columns)
	connections := make([]*sql.DB, 3)
	for i := range connections {
		connections[i] = openLineageMigrationDB(t, path)
		_, err := connections[i].Exec("PRAGMA busy_timeout=5000; PRAGMA foreign_keys=ON")
		require.NoError(t, err)
	}
	var wg sync.WaitGroup
	results := make(chan error, len(connections))
	for _, connection := range connections {
		wg.Add(1)
		go func(db *sql.DB) {
			defer wg.Done()
			results <- runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx)
		}(connection)
	}
	wg.Wait()
	close(results)
	for err := range results {
		require.NoError(t, err)
	}
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
}

func TestRecordProfileKeysCallerTransactionIsAtomic(t *testing.T) {
	db := recordProfileDB(t)
	schema := recordProfileSchema(t, db)
	err := runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `ALTER TABLE write_group_members ADD COLUMN record_test_intent TEXT`); err != nil {
			return err
		}
		if err := migrateMixedRecordProfileKeysTx(ctx, tx); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, `INSERT INTO table_does_not_exist VALUES(1)`)
		return err
	})
	require.Error(t, err)
	require.Equal(t, schema, recordProfileSchema(t, db))
	var count int
	require.NoError(t, db.QueryRow(`SELECT count(*) FROM pragma_table_info('write_group_members') WHERE name='record_test_intent'`).Scan(&count))
	require.Zero(t, count)
}
