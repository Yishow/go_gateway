package datalink

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// Migrator handles database migrations
type Migrator struct{}

// NewMigrator creates a new migrator
func NewMigrator() *Migrator {
	return &Migrator{}
}

// Migrate executes the embedded SQL migration scripts
func (m *Migrator) Migrate(db *sql.DB) error {
	// For this prototype/phase, we prioritize the sqlite file if using sqlite driver (implied by this codebase context)
	// In a full implementation we'd check driver name.
	// For now, let's look specifically for "001_initial_schema_sqlite.sql" first.

	targetFile := "001_initial_schema_sqlite.sql"

	content, err := migrations.FS.ReadFile(targetFile)
	if err == nil {
		if err := ensureSQLiteMappingLifecycleColumns(db); err != nil {
			return err
		}

		log.Printf("Executing SQLite migration: %s", targetFile)
		if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", targetFile, err)
		}
		if err := ensureSQLiteDeviceReadinessColumn(db); err != nil {
			return err
		}

		needPointUniqueMigration, err := needsSQLitePointUniqueMigration(db)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite points schema: %w", err)
		}
		if needPointUniqueMigration {
			const sqlitePointUniqueMigration = "003_align_points_unique_function_sqlite.up.sql"
			content, err := migrations.FS.ReadFile(sqlitePointUniqueMigration)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", sqlitePointUniqueMigration, err)
			}
			log.Printf("Executing SQLite migration: %s", sqlitePointUniqueMigration)
			if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", sqlitePointUniqueMigration, err)
			}
		}

		const sqliteDatabaseTargetMigration = "004_database_targets_sqlite.up.sql"
		content, err = migrations.FS.ReadFile(sqliteDatabaseTargetMigration)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", sqliteDatabaseTargetMigration, err)
		}
		log.Printf("Executing SQLite migration: %s", sqliteDatabaseTargetMigration)
		if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", sqliteDatabaseTargetMigration, err)
		}
		if err := ensureSQLiteDatabaseTargetGroupingColumns(db); err != nil {
			return err
		}
		if err := ensureSQLiteDatabaseDeliveryOutcomeColumns(db); err != nil {
			return err
		}
		if err := ensureMySQLTargetMappingSchema(db); err != nil {
			return err
		}

		const sqliteSourceRuleMigration = "005_source_rules_sqlite.up.sql"
		content, err = migrations.FS.ReadFile(sqliteSourceRuleMigration)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", sqliteSourceRuleMigration, err)
		}
		log.Printf("Executing SQLite migration: %s", sqliteSourceRuleMigration)
		if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", sqliteSourceRuleMigration, err)
		}

		if err := ensureSQLiteSourceRuleTargetDatatypeColumns(db); err != nil {
			return err
		}
		if err := ensureSQLiteSourceRuleRevisionColumn(db); err != nil {
			return err
		}
		if err := ensureSQLiteSourceRuleCandidateSnapshotsTable(db); err != nil {
			return err
		}
		if err := ensureSQLiteSourceRuleTagReviewDecisionsTable(db); err != nil {
			return err
		}
		if err := ensureSQLiteSourceRuleTagReviewDecisionStaleColumns(db); err != nil {
			return err
		}
		if err := ensureSQLiteSourceRuleModbusShareColumns(db); err != nil {
			return err
		}
		if err := ensureSQLiteModbusShareSettings(db); err != nil {
			return err
		}

		if err := ensureSQLitePointsDataFormatColumn(db); err != nil {
			return err
		}
		if err := ensureSQLiteWorkspaceAuditHistoryTable(db); err != nil {
			return err
		}
		if err := ensureSQLiteTelemetryRecordingMigrations(db); err != nil {
			return err
		}

		return nil
	}

	// Fallback to iterating if specific file not found (legacy behavior or different structure)
	files, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, file := range files {
		// Skip non-sql, down migrations, or known postgres files if we are in sqlite mode
		// Ideally we should have a cleaner way to distinguish, but for now filtering out .up.sql if it overlaps with sqlite.sql
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") || strings.Contains(file.Name(), ".down.") {
			continue
		}

		// Skip generated postgres files if we just want sqlite
		if file.Name() == "001_initial_schema.up.sql" {
			continue
		}

		log.Printf("Executing migration: %s", file.Name())

		content, err := migrations.FS.ReadFile(file.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Execute SQL
		_, err = db.ExecContext(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}

	return nil
}

func ensureSQLiteModbusShareSettings(db *sql.DB) error {
	const migrationName = "018_modbus_share_settings"
	_, err := db.ExecContext(context.Background(), `
		INSERT OR IGNORE INTO system_settings (key, value, description)
		VALUES (?, ?, ?)
	`, "modbus_share", `{"enabled":false,"bind_address":"127.0.0.1","port":5020,"slave_id":1,"capacity_registers":32768,"settings_revision":"migration-018"}`, "Local Modbus Share listener settings")
	if err != nil {
		return fmt.Errorf("failed to materialize %s: %w", migrationName, err)
	}
	var raw string
	if err := db.QueryRowContext(context.Background(), `SELECT value FROM system_settings WHERE key = ?`, "modbus_share").Scan(&raw); err != nil {
		return fmt.Errorf("failed to read %s settings: %w", migrationName, err)
	}
	var value map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return fmt.Errorf("failed to decode %s settings: %w", migrationName, err)
	}
	if capacity, ok := value["capacity_registers"].(float64); ok && capacity > 32768 {
		value["capacity_registers"] = 32768
		payload, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to encode %s settings: %w", migrationName, err)
		}
		if _, err := db.ExecContext(context.Background(), `UPDATE system_settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE key = ?`, string(payload), "modbus_share"); err != nil {
			return fmt.Errorf("failed to clamp %s settings capacity: %w", migrationName, err)
		}
	}
	return nil
}

func needsSQLitePointUniqueMigration(db *sql.DB) (bool, error) {
	var ddl string
	err := db.QueryRowContext(context.Background(), `
		SELECT sql
		FROM sqlite_master
		WHERE type = 'table' AND name = 'points'
	`).Scan(&ddl)
	if err != nil {
		return false, err
	}

	normalized := strings.ToLower(ddl)
	if strings.Contains(normalized, "unique (device_id, address, function)") {
		return false, nil
	}
	return strings.Contains(normalized, "unique (device_id, address)"), nil
}

func ensureSQLiteDeviceReadinessColumn(db *sql.DB) error {
	const migrationName = "002_add_device_collection_stats.up.sql"

	exists, err := sqliteColumnExists(db, "devices", "readiness_status")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite column readiness_status for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	log.Printf("Executing SQLite migration: %s (devices.readiness_status)", migrationName)
	if _, err := db.ExecContext(context.Background(), `ALTER TABLE devices ADD COLUMN readiness_status TEXT`); err != nil {
		return fmt.Errorf("failed to execute migration %s for devices.readiness_status: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteSourceRuleTargetDatatypeColumns(db *sql.DB) error {
	const migrationName = "006_source_rule_target_datatype_sqlite.up.sql"

	columns := []struct {
		name string
		ddl  string
	}{
		{
			name: "target_data_type",
			ddl:  "ALTER TABLE source_rules ADD COLUMN target_data_type TEXT",
		},
		{
			name: "scale_multiplier",
			ddl:  "ALTER TABLE source_rules ADD COLUMN scale_multiplier REAL",
		},
		{
			name: "scale_offset",
			ddl:  "ALTER TABLE source_rules ADD COLUMN scale_offset REAL",
		},
		{
			name: "data_format",
			ddl:  "ALTER TABLE source_rules ADD COLUMN data_format TEXT",
		},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "source_rules", column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (%s)", migrationName, column.name)
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}

// ensureSQLitePointsDataFormatColumn 為 SQLite points 表補上 data_format 欄位（與 PostgreSQL 008 對齊）。
func ensureSQLitePointsDataFormatColumn(db *sql.DB) error {
	const migrationName = "008_point_data_format_sqlite"

	exists, err := sqliteColumnExists(db, "points", "data_format")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite column data_format for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	log.Printf("Executing SQLite migration: %s", migrationName)
	if _, err := db.ExecContext(context.Background(), `ALTER TABLE points ADD COLUMN data_format TEXT`); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteSourceRuleRevisionColumn(db *sql.DB) error {
	const migrationName = "009_source_rule_revision_sqlite.up.sql"

	exists, err := sqliteColumnExists(db, "source_rules", "revision_id")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite column revision_id for migration %s: %w", migrationName, err)
	}
	if !exists {
		content, err := migrations.FS.ReadFile(migrationName)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
		}

		log.Printf("Executing SQLite migration: %s", migrationName)
		if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
		}
	}

	if _, err := db.ExecContext(context.Background(), `
		UPDATE source_rules
		SET revision_id = id || ':legacy'
		WHERE trim(coalesce(revision_id, '')) = ''
	`); err != nil {
		return fmt.Errorf("failed to backfill sqlite source-rule revision ids for migration %s: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteSourceRuleCandidateSnapshotsTable(db *sql.DB) error {
	const migrationName = "010_source_rule_candidate_snapshot_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "source_rule_candidate_snapshots")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table source_rule_candidate_snapshots for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	content, err := migrations.FS.ReadFile(migrationName)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
	}

	log.Printf("Executing SQLite migration: %s", migrationName)
	if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteSourceRuleModbusShareColumns(db *sql.DB) error {
	const migrationName = "017_source_rule_modbus_share_sqlite.up.sql"

	columns := []struct {
		name string
		ddl  string
	}{
		{
			name: "share_enabled",
			ddl:  "ALTER TABLE source_rules ADD COLUMN share_enabled INTEGER NOT NULL DEFAULT 0",
		},
		{
			name: "share_start_register",
			ddl:  "ALTER TABLE source_rules ADD COLUMN share_start_register INTEGER",
		},
		{
			name: "share_stride",
			ddl:  "ALTER TABLE source_rules ADD COLUMN share_stride INTEGER",
		},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, "source_rules", column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (%s)", migrationName, column.name)
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}

func ensureSQLiteWorkspaceAuditHistoryTable(db *sql.DB) error {
	const migrationName = "016_workspace_audit_history_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "workspace_audit_history")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table workspace_audit_history for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}

	content, err := migrations.FS.ReadFile(migrationName)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
	}

	log.Printf("Executing SQLite migration: %s", migrationName)
	if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
	}

	return nil
}

func ensureSQLiteDatabaseTargetGroupingColumns(db *sql.DB) error {
	const migrationName = "014_database_target_grouping_sqlite.up.sql"

	columns := []struct {
		table string
		name  string
		ddl   string
	}{
		{
			table: "database_connectors",
			name:  "default_write_interval_seconds",
			ddl:   "ALTER TABLE database_connectors ADD COLUMN default_write_interval_seconds INTEGER NOT NULL DEFAULT 15",
		},
		{
			table: "database_target_mappings",
			name:  "group_key",
			ddl:   "ALTER TABLE database_target_mappings ADD COLUMN group_key TEXT",
		},
		{
			table: "database_target_mappings",
			name:  "write_interval_seconds",
			ddl:   "ALTER TABLE database_target_mappings ADD COLUMN write_interval_seconds INTEGER",
		},
	}

	for _, column := range columns {
		exists, err := sqliteColumnExists(db, column.table, column.name)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s for migration %s: %w", column.name, migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (%s.%s)", migrationName, column.table, column.name)
		if _, err := db.ExecContext(context.Background(), column.ddl); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s: %w", migrationName, column.name, err)
		}
	}

	return nil
}

func sqliteColumnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.QueryContext(context.Background(), fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &pk); err != nil {
			return false, err
		}
		if strings.EqualFold(name, columnName) {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}

func sqliteTableExists(db *sql.DB, tableName string) (bool, error) {
	var name sql.NullString
	err := db.QueryRowContext(context.Background(), `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = ?
	`, tableName).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return name.Valid, nil
}
