package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrator_SQLiteSourceRuleMigration_PartialSharePresenceTable(t *testing.T) {
	tests := []struct {
		name             string
		shareColumns     []string
		legacyNewColumns []string
	}{
		{name: "none", legacyNewColumns: []string{"target_data_type"}},
		{name: "share_enabled", shareColumns: []string{"share_enabled"}, legacyNewColumns: []string{"scale_multiplier"}},
		{name: "share_start_register", shareColumns: []string{"share_start_register"}, legacyNewColumns: []string{"scale_offset"}},
		{name: "share_enabled_and_start", shareColumns: []string{"share_enabled", "share_start_register"}, legacyNewColumns: []string{"data_format"}},
		{name: "share_stride", shareColumns: []string{"share_stride"}, legacyNewColumns: []string{"revision_id"}},
		{name: "share_enabled_and_stride", shareColumns: []string{"share_enabled", "share_stride"}, legacyNewColumns: []string{"target_data_type", "scale_multiplier"}},
		{name: "share_start_and_stride", shareColumns: []string{"share_start_register", "share_stride"}, legacyNewColumns: []string{"scale_offset", "data_format", "revision_id"}},
		{name: "all_share_columns", shareColumns: []string{"share_enabled", "share_start_register", "share_stride"}, legacyNewColumns: []string{"target_data_type", "scale_multiplier", "scale_offset", "data_format", "revision_id"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sql.Open("sqlite", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			require.NoError(t, createLegacySourceRuleSchema(t.Context(), db, "rule-partial-"+tt.name))
			for _, column := range append(append([]string{}, tt.shareColumns...), tt.legacyNewColumns...) {
				require.NoError(t, addLegacySourceRuleColumn(t.Context(), db, column))
				require.NoError(t, setLegacySourceRuleColumnValue(t.Context(), db, column))
			}

			migrator := NewMigrator()
			require.NoError(t, migrator.Migrate(db))

			allColumns := []string{
				"target_data_type", "scale_multiplier", "scale_offset", "data_format", "revision_id",
				"share_enabled", "share_start_register", "share_stride",
			}
			for _, column := range allColumns {
				exists, existsErr := sqliteColumnExists(db, "source_rules", column)
				require.NoError(t, existsErr, "inspect %s", column)
				require.True(t, exists, "column %s should exist after migration", column)
			}

			var id, deviceID, startAddress, dataType, namingPrefix, origin, skippedAddresses, templateName string
			var count, enabled, locked int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT id, device_id, start_address, count, data_type, naming_prefix, enabled, locked, origin, template_name, skipped_addresses FROM source_rules WHERE id = ?`, "rule-partial-"+tt.name).Scan(
				&id, &deviceID, &startAddress, &count, &dataType, &namingPrefix, &enabled, &locked, &origin, &templateName, &skippedAddresses))
			require.Equal(t, "rule-partial-"+tt.name, id)
			require.Equal(t, "device-partial", deviceID)
			require.Equal(t, "40001", startAddress)
			require.Equal(t, 2, count)
			require.Equal(t, "int16", dataType)
			require.Equal(t, "LEGACY_", namingPrefix)
			require.Equal(t, 1, enabled)
			require.Equal(t, 1, locked)
			require.Equal(t, "template", origin)
			require.Equal(t, "TPL_LEGACY", templateName)
			require.Equal(t, `["40002"]`, skippedAddresses)

			legacyValues := map[string]any{"target_data_type": "float32", "scale_multiplier": 2.5, "scale_offset": -1.5, "data_format": "DCBA", "revision_id": "legacy-revision"}
			for column, expected := range legacyValues {
				assertMigratedSourceRuleValue(t, db, id, column, containsString(tt.legacyNewColumns, column), expected)
			}
			shareValues := map[string]any{"share_enabled": 1, "share_start_register": 40101, "share_stride": 2}
			for column, expected := range shareValues {
				assertMigratedSourceRuleValue(t, db, id, column, containsString(tt.shareColumns, column), expected)
			}

			var columnsBefore int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM pragma_table_info('source_rules')`).Scan(&columnsBefore))
			require.NoError(t, migrator.Migrate(db))
			require.NoError(t, migrator.Migrate(db))
			var columnsAfter int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count(*) FROM pragma_table_info('source_rules')`).Scan(&columnsAfter))
			require.Equal(t, columnsBefore, columnsAfter)
			var countAfterRerun int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT count FROM source_rules WHERE id = ?`, "rule-partial-"+tt.name).Scan(&countAfterRerun))
			require.Equal(t, 2, countAfterRerun)
		})
	}
}

func createLegacySourceRuleSchema(ctx context.Context, db *sql.DB, ruleID string) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE devices (id TEXT PRIMARY KEY, name TEXT NOT NULL, protocol TEXT NOT NULL, connection_config TEXT NOT NULL DEFAULT '{}', status TEXT NOT NULL DEFAULT 'draft', created_at TEXT NOT NULL DEFAULT (datetime('now')), updated_at TEXT NOT NULL DEFAULT (datetime('now')));
		INSERT INTO devices (id, name, protocol) VALUES ('device-partial', 'Legacy PLC', 'modbus_tcp');
		CREATE TABLE source_rules (id TEXT PRIMARY KEY, device_id TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE, start_address TEXT NOT NULL, count INTEGER NOT NULL CHECK (count > 0), data_type TEXT NOT NULL, naming_prefix TEXT NOT NULL, enabled INTEGER NOT NULL DEFAULT 1, locked INTEGER NOT NULL DEFAULT 0, origin TEXT NOT NULL DEFAULT 'manual', template_name TEXT, skipped_addresses TEXT NOT NULL DEFAULT '[]', created_at TEXT NOT NULL DEFAULT (datetime('now')), updated_at TEXT NOT NULL DEFAULT (datetime('now')));
		INSERT INTO source_rules (id, device_id, start_address, count, data_type, naming_prefix, enabled, locked, origin, template_name, skipped_addresses)
		VALUES (?, 'device-partial', '40001', 2, 'int16', 'LEGACY_', 1, 1, 'template', 'TPL_LEGACY', '["40002"]');
	`, ruleID)
	return err
}

func addLegacySourceRuleColumn(ctx context.Context, db *sql.DB, column string) error {
	ddls := map[string]string{
		"target_data_type":     "ALTER TABLE source_rules ADD COLUMN target_data_type TEXT",
		"scale_multiplier":     "ALTER TABLE source_rules ADD COLUMN scale_multiplier REAL",
		"scale_offset":         "ALTER TABLE source_rules ADD COLUMN scale_offset REAL",
		"data_format":          "ALTER TABLE source_rules ADD COLUMN data_format TEXT",
		"revision_id":          "ALTER TABLE source_rules ADD COLUMN revision_id TEXT NOT NULL DEFAULT ''",
		"share_enabled":        "ALTER TABLE source_rules ADD COLUMN share_enabled INTEGER NOT NULL DEFAULT 0",
		"share_start_register": "ALTER TABLE source_rules ADD COLUMN share_start_register INTEGER",
		"share_stride":         "ALTER TABLE source_rules ADD COLUMN share_stride INTEGER",
	}
	ddl, ok := ddls[column]
	if !ok {
		return fmt.Errorf("unknown source-rule column %q", column)
	}
	_, err := db.ExecContext(ctx, ddl)
	return err
}

func setLegacySourceRuleColumnValue(ctx context.Context, db *sql.DB, column string) error {
	values := map[string]any{
		"target_data_type":     "float32",
		"scale_multiplier":     2.5,
		"scale_offset":         -1.5,
		"data_format":          "DCBA",
		"revision_id":          "legacy-revision",
		"share_enabled":        1,
		"share_start_register": 40101,
		"share_stride":         2,
	}
	value, ok := values[column]
	if !ok {
		return fmt.Errorf("unknown source-rule column %q", column)
	}
	updates := map[string]string{
		"target_data_type":     "UPDATE source_rules SET target_data_type = ? WHERE id LIKE 'rule-partial-%'",
		"scale_multiplier":     "UPDATE source_rules SET scale_multiplier = ? WHERE id LIKE 'rule-partial-%'",
		"scale_offset":         "UPDATE source_rules SET scale_offset = ? WHERE id LIKE 'rule-partial-%'",
		"data_format":          "UPDATE source_rules SET data_format = ? WHERE id LIKE 'rule-partial-%'",
		"revision_id":          "UPDATE source_rules SET revision_id = ? WHERE id LIKE 'rule-partial-%'",
		"share_enabled":        "UPDATE source_rules SET share_enabled = ? WHERE id LIKE 'rule-partial-%'",
		"share_start_register": "UPDATE source_rules SET share_start_register = ? WHERE id LIKE 'rule-partial-%'",
		"share_stride":         "UPDATE source_rules SET share_stride = ? WHERE id LIKE 'rule-partial-%'",
	}
	statement, ok := updates[column]
	if !ok {
		return fmt.Errorf("unknown source-rule column %q", column)
	}
	_, err := db.ExecContext(ctx, statement, value)
	return err
}

func containsString(columns []string, target string) bool {
	for _, column := range columns {
		if column == target {
			return true
		}
	}
	return false
}

func assertMigratedSourceRuleValue(t *testing.T, db *sql.DB, ruleID, column string, existed bool, expected any) {
	t.Helper()
	var value sql.NullString
	queries := map[string]string{
		"target_data_type":     "SELECT CAST(target_data_type AS TEXT) FROM source_rules WHERE id = ?",
		"scale_multiplier":     "SELECT CAST(scale_multiplier AS TEXT) FROM source_rules WHERE id = ?",
		"scale_offset":         "SELECT CAST(scale_offset AS TEXT) FROM source_rules WHERE id = ?",
		"data_format":          "SELECT CAST(data_format AS TEXT) FROM source_rules WHERE id = ?",
		"revision_id":          "SELECT CAST(revision_id AS TEXT) FROM source_rules WHERE id = ?",
		"share_enabled":        "SELECT CAST(share_enabled AS TEXT) FROM source_rules WHERE id = ?",
		"share_start_register": "SELECT CAST(share_start_register AS TEXT) FROM source_rules WHERE id = ?",
		"share_stride":         "SELECT CAST(share_stride AS TEXT) FROM source_rules WHERE id = ?",
	}
	query, ok := queries[column]
	if !ok {
		t.Fatalf("unknown source-rule column %q", column)
	}
	require.NoError(t, db.QueryRowContext(t.Context(), query, ruleID).Scan(&value))
	if !existed && column != "revision_id" {
		if column == "share_enabled" {
			require.Equal(t, "0", value.String)
			return
		}
		require.False(t, value.Valid, "%s should use NULL when introduced by migration", column)
		return
	}
	if !existed {
		expected = ruleID + ":legacy"
	}
	require.True(t, value.Valid)
	require.Equal(t, fmt.Sprint(expected), value.String)
}
