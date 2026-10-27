package datalink

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLServerPartitionMaintenanceMigration(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "sqlserver-partitions.db"))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	exists, err := sqliteTableExists(db, "write_group_sqlserver_partition_maintenance")
	require.NoError(t, err)
	require.True(t, exists, "startup must install the separate SQL Server authorized partition-state table")
	const current = `{"method":"RANGE RIGHT","column":"bucket_start"}`
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_sqlserver_partition_maintenance (
		workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,
		base_layout_digest,connector_revision,generation,current_evidence,revision,status,updated_at)
		VALUES ('workspace','connector','recordings','dbo','native_monthly','managed_long_v1',
		'base-digest','revision',1,?,1,'ready','2026-10-25T00:00:00Z')`, current)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	var saved string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT current_evidence FROM write_group_sqlserver_partition_maintenance`).Scan(&saved))
	require.Equal(t, current, saved, "startup preserves the independent SQL Server state bytes")
	_, err = db.ExecContext(t.Context(), `DROP TABLE write_group_sqlserver_partition_maintenance`)
	require.NoError(t, err)
	require.NoError(t, migrator.Migrate(db))
	exists, err = sqliteTableExists(db, "write_group_sqlserver_partition_maintenance")
	require.NoError(t, err)
	require.True(t, exists, "a missing additive SQL Server sidecar is reinstalled")
}
