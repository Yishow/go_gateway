package datalink

import (
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrator_SQLiteSourceRuleRollbackRequiresBackupRestore(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	_, err = db.ExecContext(t.Context(), `
		INSERT INTO devices (id, name, protocol) VALUES ('device-rollback-boundary', 'PLC', 'modbus_tcp');
		INSERT INTO source_rules (id, device_id, start_address, count, data_type, naming_prefix)
		VALUES ('rule-rollback-boundary', 'device-rollback-boundary', '40001', 1, 'int16', 'LEGACY_');
	`)
	require.NoError(t, err)

	downSQL, err := migrations.FS.ReadFile("017_source_rule_modbus_share.down.sql")
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(downSQL))
	t.Logf("sqlite down migration capability probe error: %v", err)
	require.Error(t, err, "the repository does not provide a runnable SQLite down migration")

	var rowCount int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM source_rules WHERE id = 'rule-rollback-boundary'`).Scan(&rowCount))
	require.Equal(t, 1, rowCount, "failed down migration must not be treated as rollback evidence")
	for _, column := range []string{"share_enabled", "share_start_register", "share_stride"} {
		exists, existsErr := sqliteColumnExists(db, "source_rules", column)
		require.NoError(t, existsErr)
		require.True(t, exists, "failed down migration should leave %s intact", column)
	}
}

func TestMigrator_SQLiteSourceRuleRevisionPartialUpgradeBackfillsBlankValues(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, createLegacySourceRuleSchema(t.Context(), db, "rule-revision-partial"))
	// Simulate a migration interrupted immediately after ALTER TABLE and before
	// the 009 data backfill statement.
	require.NoError(t, addLegacySourceRuleColumn(t.Context(), db, "revision_id"))

	require.NoError(t, NewMigrator().Migrate(db))

	var revisionID string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT revision_id FROM source_rules WHERE id = 'rule-revision-partial'`).Scan(&revisionID))
	require.Equal(t, "rule-revision-partial:legacy", revisionID)
}

func TestMigrator_SQLiteSourceRuleMigration_PreservesLegacyDataAndPartialRecovery(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE devices (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			protocol TEXT NOT NULL,
			connection_config TEXT NOT NULL DEFAULT '{}',
			status TEXT NOT NULL DEFAULT 'draft',
			readiness_status TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now'))
		);
		INSERT INTO devices (id, name, protocol) VALUES ('dev-legacy-1', 'PLC 1', 'modbus_tcp');

		CREATE TABLE source_rules (
			id                TEXT PRIMARY KEY,
			device_id         TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
			start_address     TEXT NOT NULL,
			count             INTEGER NOT NULL CHECK (count > 0),
			data_type         TEXT NOT NULL CHECK (data_type IN ('bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string')),
			naming_prefix     TEXT NOT NULL,
			enabled           INTEGER NOT NULL DEFAULT 1,
			locked            INTEGER NOT NULL DEFAULT 0,
			origin            TEXT NOT NULL DEFAULT 'manual' CHECK (origin IN ('manual', 'template')),
			template_name     TEXT,
			skipped_addresses TEXT NOT NULL DEFAULT '[]',
			created_at        TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at        TEXT NOT NULL DEFAULT (datetime('now'))
		);
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), `
		INSERT INTO source_rules (
			id, device_id, start_address, count, data_type, naming_prefix, enabled, locked, origin, template_name, skipped_addresses
		) VALUES (
			'rule-legacy-1', 'dev-legacy-1', '40001', 8, 'int32', 'LINE1_', 1, 1, 'template', 'TPL_PUMP', '["40003"]'
		);
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), `
		ALTER TABLE source_rules ADD COLUMN target_data_type TEXT;
		ALTER TABLE source_rules ADD COLUMN scale_multiplier REAL;
		UPDATE source_rules SET target_data_type = 'float32', scale_multiplier = 2.5 WHERE id = 'rule-legacy-1';
	`)
	require.NoError(t, err)

	migrator := NewMigrator()
	require.NoError(t, migrator.Migrate(db))

	expectedCols := []string{
		"target_data_type", "scale_multiplier", "scale_offset", "data_format",
		"revision_id", "share_enabled", "share_start_register", "share_stride",
	}
	for _, col := range expectedCols {
		var colName string
		err = db.QueryRowContext(t.Context(), `SELECT name FROM pragma_table_info('source_rules') WHERE name = ?`, col).Scan(&colName)
		require.NoError(t, err, "column %s should exist", col)
		assert.Equal(t, col, colName)
	}

	var (
		id                 string
		deviceID           string
		startAddress       string
		count              int
		dataType           string
		namingPrefix       string
		enabled            int
		locked             int
		origin             string
		templateName       sql.NullString
		skippedAddrs       string
		targetDataType     sql.NullString
		scaleMultiplier    sql.NullFloat64
		scaleOffset        sql.NullFloat64
		dataFormat         sql.NullString
		shareEnabled       int
		shareStartRegister sql.NullInt64
		shareStride        sql.NullInt64
	)
	row := db.QueryRowContext(t.Context(), `
		SELECT id, device_id, start_address, count, data_type, naming_prefix, enabled, locked, origin, template_name, skipped_addresses,
		       target_data_type, scale_multiplier, scale_offset, data_format, share_enabled, share_start_register, share_stride
		FROM source_rules WHERE id = 'rule-legacy-1'
	`)
	err = row.Scan(
		&id, &deviceID, &startAddress, &count, &dataType, &namingPrefix, &enabled, &locked, &origin, &templateName, &skippedAddrs,
		&targetDataType, &scaleMultiplier, &scaleOffset, &dataFormat, &shareEnabled, &shareStartRegister, &shareStride,
	)
	require.NoError(t, err)

	assert.Equal(t, "rule-legacy-1", id)
	assert.Equal(t, "dev-legacy-1", deviceID)
	assert.Equal(t, "40001", startAddress)
	assert.Equal(t, 8, count)
	assert.Equal(t, "int32", dataType)
	assert.Equal(t, "LINE1_", namingPrefix)
	assert.Equal(t, 1, enabled)
	assert.Equal(t, 1, locked)
	assert.Equal(t, "template", origin)
	assert.True(t, templateName.Valid)
	assert.Equal(t, "TPL_PUMP", templateName.String)
	assert.Equal(t, `["40003"]`, skippedAddrs)
	assert.True(t, targetDataType.Valid)
	assert.Equal(t, "float32", targetDataType.String)
	assert.True(t, scaleMultiplier.Valid)
	assert.Equal(t, 2.5, scaleMultiplier.Float64)
	assert.False(t, scaleOffset.Valid)
	assert.False(t, dataFormat.Valid)
	assert.Equal(t, 0, shareEnabled)
	assert.False(t, shareStartRegister.Valid)
	assert.False(t, shareStride.Valid)

	require.NoError(t, migrator.Migrate(db))
	require.NoError(t, migrator.Migrate(db))
	var postShareEnabled int
	err = db.QueryRowContext(t.Context(), `SELECT share_enabled FROM source_rules WHERE id = 'rule-legacy-1'`).Scan(&postShareEnabled)
	require.NoError(t, err)
	assert.Equal(t, 0, postShareEnabled)
}
