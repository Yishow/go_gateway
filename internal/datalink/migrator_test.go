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
		err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		assert.NoError(t, err, "Table %s should exist", table)
		assert.Equal(t, table, name)
	}

	var colName string
	err = db.QueryRow(`SELECT name FROM pragma_table_info('points') WHERE name = 'data_format'`).Scan(&colName)
	require.NoError(t, err)
	assert.Equal(t, "data_format", colName)

	err = db.QueryRow(`SELECT name FROM pragma_table_info('source_rules') WHERE name = 'revision_id'`).Scan(&colName)
	require.NoError(t, err)
	assert.Equal(t, "revision_id", colName)

	// Assert: Check a specific column to ensure schema is correct (e.g., devices.protocol)
	// var protocol string // Unused for now
	// Insert a dummy row to check constraints (optional, but good for deeper verification)
	// For now, just checking table existence is sufficient for the first pass.
}
