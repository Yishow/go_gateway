package datalink

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func openMigratedMigratorTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "main.db"))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, NewMigrator().Migrate(db))
	return db
}

func insertMigratorConnector(t *testing.T, db *sql.DB, id, kind, connectionConfig string) {
	t.Helper()

	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO database_connectors (id, name, kind, connection_config, status, enabled, default_write_interval_seconds, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'ready', 1, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		id, "connector-"+id, kind, connectionConfig,
	)
	require.NoError(t, err)
}

func insertMigratorMapping(t *testing.T, db *sql.DB, id, connectorID, tableSchema string) {
	t.Helper()

	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO database_target_mappings (id, tag_id, connector_id, table_schema, table_name, column_name, write_mode, enabled, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 'sensor_readings', 'value', 'insert', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		id, "tag-"+id, connectorID, tableSchema,
	)
	require.NoError(t, err)
}

func mappingSchema(t *testing.T, db *sql.DB, id string) string {
	t.Helper()

	var tableSchema string
	require.NoError(t, db.QueryRowContext(
		context.Background(),
		`SELECT table_schema FROM database_target_mappings WHERE id = ?`,
		id,
	).Scan(&tableSchema))
	return tableSchema
}

// 舊版寫入的 MySQL 映射帶著 SQLite 的 "main" 預設值，會產生指向不存在資料庫的語句。
func TestEnsureMySQLTargetMappingSchema_NormalizesLegacyMainSchema(t *testing.T) {
	db := openMigratedMigratorTestDB(t)
	insertMigratorConnector(t, db, "conn-mysql", "mysql", `{"host":"127.0.0.1","port":"3306","user":"gateway","database":"gateway_metrics"}`)
	insertMigratorMapping(t, db, "map-legacy", "conn-mysql", "main")
	insertMigratorMapping(t, db, "map-explicit", "conn-mysql", "reporting")

	require.NoError(t, ensureMySQLTargetMappingSchema(db))

	assert.Equal(t, "gateway_metrics", mappingSchema(t, db, "map-legacy"))
	assert.Equal(t, "reporting", mappingSchema(t, db, "map-explicit"), "使用者明確指定的 schema 不得被改寫")
}

func TestEnsureMySQLTargetMappingSchema_LeavesOtherKindsUntouched(t *testing.T) {
	db := openMigratedMigratorTestDB(t)
	insertMigratorConnector(t, db, "conn-sqlite", "sqlite", `{"dsn":"gateway.db"}`)
	insertMigratorMapping(t, db, "map-sqlite", "conn-sqlite", "main")

	require.NoError(t, ensureMySQLTargetMappingSchema(db))

	assert.Equal(t, "main", mappingSchema(t, db, "map-sqlite"))
}

func TestEnsureMySQLTargetMappingSchema_KeepsValueWhenDatabaseNameUnknown(t *testing.T) {
	db := openMigratedMigratorTestDB(t)
	insertMigratorConnector(t, db, "conn-mysql-dsn", "mysql", `{"dsn":"gateway:secret@tcp(127.0.0.1:3306)/gateway_metrics"}`)
	insertMigratorMapping(t, db, "map-dsn", "conn-mysql-dsn", "main")

	require.NoError(t, ensureMySQLTargetMappingSchema(db))

	assert.Equal(t, "main", mappingSchema(t, db, "map-dsn"), "無法可靠取得資料庫名稱時不得猜測改寫")
}

func TestEnsureMySQLTargetMappingSchema_IsIdempotent(t *testing.T) {
	db := openMigratedMigratorTestDB(t)
	insertMigratorConnector(t, db, "conn-mysql", "mysql", `{"host":"127.0.0.1","database":"gateway_metrics"}`)
	insertMigratorMapping(t, db, "map-legacy", "conn-mysql", "main")

	require.NoError(t, ensureMySQLTargetMappingSchema(db))
	require.NoError(t, ensureMySQLTargetMappingSchema(db))

	assert.Equal(t, "gateway_metrics", mappingSchema(t, db, "map-legacy"))
}
