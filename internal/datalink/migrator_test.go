package datalink

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMigrator_Migrate(t *testing.T) {
	// Arrange
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	ctx := t.Context()

	migrator := NewMigrator()

	// Act
	err = migrator.Migrate(db)
	require.NoError(t, err)
	err = migrator.Migrate(db)
	require.NoError(t, err)

	// Assert: Check if tables exist
	tables := []string{
		"devices",
		"points",
		"tags",
		"mappings",
		"polling_groups",
		"source_rules",
		"source_rule_links",
		"source_rule_candidate_snapshots",
		"system_settings",
		"timeseries",
		"database_connectors",
		"database_target_mappings",
	}
	for _, table := range tables {
		var name string
		err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name)
	}

	var colName string
	err = db.QueryRowContext(ctx, `SELECT name FROM pragma_table_info('points') WHERE name = 'data_format'`).Scan(&colName)
	require.NoError(t, err)
	assert.Equal(t, "data_format", colName)

	err = db.QueryRowContext(ctx, `SELECT name FROM pragma_table_info('source_rules') WHERE name = 'revision_id'`).Scan(&colName)
	require.NoError(t, err)
	assert.Equal(t, "revision_id", colName)

	for _, column := range []string{"share_enabled", "share_start_register", "share_stride"} {
		err = db.QueryRowContext(ctx, `SELECT name FROM pragma_table_info('source_rules') WHERE name = ?`, column).Scan(&colName)
		require.NoError(t, err)
		assert.Equal(t, column, colName)
	}

	for _, column := range []string{
		"status",
		"rule_candidate_id",
		"proposed_signature",
		"last_applied_signature",
		"blocking_reason",
	} {
		err = db.QueryRowContext(ctx, `SELECT name FROM pragma_table_info('mappings') WHERE name = ?`, column).Scan(&colName)
		require.NoError(t, err)
		assert.Equal(t, column, colName)
	}

	var modbusShareSettings string
	err = db.QueryRowContext(ctx, `SELECT value FROM system_settings WHERE key = 'modbus_share'`).Scan(&modbusShareSettings)
	require.NoError(t, err)
	assert.Contains(t, modbusShareSettings, `"enabled":false`)
	assert.Contains(t, modbusShareSettings, `"port":5020`)
	assert.Contains(t, modbusShareSettings, `"bind_address":"127.0.0.1"`)
	assert.Contains(t, modbusShareSettings, `"slave_id":1`)
	assert.Contains(t, modbusShareSettings, `"capacity_registers":32768`)

	// Assert: Check a specific column to ensure schema is correct (e.g., devices.protocol)
	// var protocol string // Unused for now
	// Insert a dummy row to check constraints (optional, but good for deeper verification)
	// For now, just checking table existence is sufficient for the first pass.
}

func TestMigrator_Migrate_LegacyMappingsTable(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()
	ctx := t.Context()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE mappings (
			id TEXT PRIMARY KEY,
			point_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			transform_pipeline TEXT NOT NULL DEFAULT '[]',
			enabled INTEGER NOT NULL DEFAULT 1,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE (point_id, tag_id)
		);
	`)
	require.NoError(t, err)

	migrator := NewMigrator()

	err = migrator.Migrate(db)
	require.NoError(t, err)

	for _, column := range []string{
		"status",
		"rule_candidate_id",
		"proposed_signature",
		"last_applied_signature",
		"blocking_reason",
	} {
		var colName string
		err = db.QueryRowContext(ctx, `SELECT name FROM pragma_table_info('mappings') WHERE name = ?`, column).Scan(&colName)
		require.NoError(t, err)
		assert.Equal(t, column, colName)
	}
}
