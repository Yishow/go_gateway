package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func ensureSQLiteWriteGroupsTables(db *sql.DB) error {
	const migrationName = "024_write_groups_sqlite.up.sql"
	// Check markers before legacy helpers could recreate a missing key table.
	if _, err := sqliteWriteGroupLineageState(context.Background(), db); err != nil {
		return err
	}

	groupsExists, err := sqliteTableExists(db, "write_groups")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_groups for migration %s: %w", migrationName, err)
	}
	membersExists, err := sqliteTableExists(db, "write_group_members")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_group_members for migration %s: %w", migrationName, err)
	}
	versionsExists, err := sqliteTableExists(db, "write_group_versions")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_group_versions for migration %s: %w", migrationName, err)
	}
	migrationMapsExists, err := sqliteTableExists(db, "write_group_migration_maps")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_group_migration_maps for migration %s: %w", migrationName, err)
	}
	if !groupsExists || !membersExists || !versionsExists || !migrationMapsExists {
		if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
			return err
		}
	}
	if err := ensureSQLiteWriteGroupMemberEntityKey(db); err != nil {
		return err
	}
	if err := ensureSQLiteBasicManagedKeysTable(db); err != nil {
		return err
	}
	if err := ensureSQLiteCircuitRecordingProfiles(db); err != nil {
		return err
	}
	if err := ensureSQLiteAdditiveColumns(db, "038_write_group_recording_reduction_sqlite.up.sql"); err != nil {
		return err
	}
	if err := ensureSQLiteAdditiveColumns(db, "040_write_group_exact_recording_source_sqlite.up.sql"); err != nil {
		return err
	}
	if err := ensureSQLiteMySQLPartitionMaintenance(db); err != nil {
		return err
	}
	if err := ensureSQLiteSQLServerPartitionMaintenance(db); err != nil {
		return err
	}
	if err := ensureSQLiteWriteGroupCutovers(db); err != nil {
		return err
	}
	return ensureSQLiteWriteGroupLineages(db)
}

// ensureSQLiteWriteGroupCutovers adds the forward cutover ledger beside the
// write-group tables of an installation that already has them.
func ensureSQLiteWriteGroupCutovers(db *sql.DB) error {
	const migrationName = "036_write_group_cutovers_sqlite.up.sql"
	for _, table := range []string{"write_group_cutovers", "write_group_revision_ends"} {
		exists, err := sqliteTableExists(db, table)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite table %s for migration %s: %w", table, migrationName, err)
		}
		if !exists {
			if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
				return err
			}
			break
		}
	}
	if err := ensureSQLiteAdditiveColumns(db, "037_write_group_cutover_intent_sqlite.up.sql"); err != nil {
		return err
	}
	// A review is applied once: its digest names at most one cutover.
	if _, err := db.ExecContext(context.Background(), `CREATE UNIQUE INDEX IF NOT EXISTS idx_write_group_cutovers_review ON write_group_cutovers(workspace_id, review_digest)`); err != nil {
		return fmt.Errorf("failed to index write_group_cutovers review digest: %w", err)
	}
	return applySQLiteMigrationOnce(db, "044_write_group_ownership_fences_sqlite.up.sql")
}

func ensureSQLiteBasicManagedKeysTable(db *sql.DB) error {
	const migrationName = "028_write_group_basic_keys_sqlite.up.sql"
	exists, err := sqliteTableExists(db, "write_group_basic_keys")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table write_group_basic_keys for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, migrationName)
}

// ensureSQLiteWriteGroupMemberEntityKey upgrades a pre-024 member table in
// place. CREATE TABLE IF NOT EXISTS cannot add columns to an existing table,
// and the member payload must survive that upgrade unchanged.
func ensureSQLiteWriteGroupMemberEntityKey(db *sql.DB) error {
	exists, err := sqliteColumnExists(db, "write_group_members", "entity_key")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite write-group member entity_key: %w", err)
	}
	if exists {
		return nil
	}
	log.Printf("Executing SQLite migration: 024_write_groups_sqlite.up.sql (write_group_members.entity_key)")
	if _, err := db.ExecContext(context.Background(), `ALTER TABLE write_group_members ADD COLUMN entity_key TEXT NOT NULL DEFAULT ''`); err != nil {
		return fmt.Errorf("failed to add sqlite write-group member entity_key: %w", err)
	}
	return nil
}
