package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema/migrations"
)

// applySQLiteMigrationOnce reads a migration from the embedded filesystem and
// executes it on the SQLite target after the caller's existence check reported
// the migration is still missing.
func applySQLiteMigrationOnce(db *sql.DB, migrationName string) error {
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

func ensureSQLiteDatabaseConnectorIdentityRevision(db *sql.DB) error {
	const migrationName = "021_database_connector_identity_sqlite.up.sql"

	exists, err := sqliteColumnExists(db, "database_connectors", "identity_revision")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite connector identity revision column: %w", err)
	}
	if !exists {
		if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
			return err
		}
	}

	rows, err := db.QueryContext(context.Background(), `
		SELECT id FROM database_connectors
		WHERE identity_revision IS NULL OR TRIM(identity_revision) = ''
	`)
	if err != nil {
		return fmt.Errorf("failed to inspect empty connector identity revisions: %w", err)
	}
	var pendingIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return fmt.Errorf("failed to scan empty connector identity revision: %w", err)
		}
		pendingIDs = append(pendingIDs, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("failed to iterate empty connector identity revisions: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("failed to close connector identity revision rows: %w", err)
	}

	for _, id := range pendingIDs {
		revision, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("failed to generate connector identity revision: %w", err)
		}
		if _, err := db.ExecContext(context.Background(), `
			UPDATE database_connectors
			SET identity_revision = ?
			WHERE id = ? AND (identity_revision IS NULL OR TRIM(identity_revision) = '')
		`, revision, strings.TrimSpace(id)); err != nil {
			return fmt.Errorf("failed to backfill connector identity revision: %w", err)
		}
	}

	return nil
}
