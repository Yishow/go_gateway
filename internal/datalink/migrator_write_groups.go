package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func ensureSQLiteWriteGroupsTables(db *sql.DB) error {
	const migrationName = "024_write_groups_sqlite.up.sql"

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
	if groupsExists && membersExists && versionsExists && migrationMapsExists {
		return ensureSQLiteWriteGroupMemberEntityKey(db)
	}
	if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
		return err
	}
	return ensureSQLiteWriteGroupMemberEntityKey(db)
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
