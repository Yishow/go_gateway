package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

var sqliteAddColumnStatement = regexp.MustCompile(`(?i)^ALTER\s+TABLE\s+(\w+)\s+ADD\s+COLUMN\s+(\w+)\s`)

// ensureSQLiteCircuitRecordingProfiles adds the nullable profile and storage
// metadata and their key tables. Existing rows keep NULL, which is read as wide.
func ensureSQLiteCircuitRecordingProfiles(db *sql.DB) error {
	for _, step := range []struct{ columns, table, tableMigration string }{
		{"029_write_group_circuit_profiles_sqlite.up.sql", "write_group_profile_keys", "030_write_group_profile_keys_sqlite.up.sql"},
		{"031_write_group_storage_mode_sqlite.up.sql", "write_group_storage_scopes", "032_write_group_storage_scopes_sqlite.up.sql"},
	} {
		if err := ensureSQLiteAdditiveColumns(db, step.columns); err != nil {
			return err
		}
		exists, err := sqliteTableExists(db, step.table)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite table %s for migration %s: %w", step.table, step.tableMigration, err)
		}
		if exists {
			continue
		}
		if err := applySQLiteMigrationOnce(db, step.tableMigration); err != nil {
			return err
		}
	}
	return ensureSQLiteStorageLineage(db)
}

// ensureSQLiteStorageLineage adds the lineage evidence columns to the storage
// scopes (which the loop above guarantees exist) and the diagnostics table.
func ensureSQLiteStorageLineage(db *sql.DB) error {
	if err := ensureSQLiteAdditiveColumns(db, "034_write_group_storage_lineage_sqlite.up.sql"); err != nil {
		return err
	}
	const eventsMigration = "035_write_group_storage_events_sqlite.up.sql"
	exists, err := sqliteTableExists(db, "write_group_storage_events")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_group_storage_events for migration %s: %w", eventsMigration, err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, eventsMigration)
}

// ensureSQLiteAdditiveColumns runs a one-statement-per-line ALTER file,
// skipping columns that already exist so a re-run, or a member table recreated
// by the 024 repair path, is completed rather than rejected.
func ensureSQLiteAdditiveColumns(db *sql.DB, migrationName string) error {
	content, err := migrations.FS.ReadFile(migrationName)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		statement := strings.TrimSpace(line)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}
		match := sqliteAddColumnStatement.FindStringSubmatch(statement)
		if match == nil {
			return fmt.Errorf("unsupported statement in migration %s: %q", migrationName, statement)
		}
		exists, err := sqliteColumnExists(db, match[1], match[2])
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite column %s.%s for migration %s: %w", match[1], match[2], migrationName, err)
		}
		if exists {
			continue
		}
		log.Printf("Executing SQLite migration: %s (%s.%s)", migrationName, match[1], match[2])
		if _, err := db.ExecContext(context.Background(), statement); err != nil {
			return fmt.Errorf("failed to execute migration %s for column %s.%s: %w", migrationName, match[1], match[2], err)
		}
	}
	return nil
}
