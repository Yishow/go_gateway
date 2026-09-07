package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go-gateway/internal/datalink/schema/migrations"
)

func ensureSQLiteTelemetryRecordingMigrations(db *sql.DB) error {
	if err := ensureSQLiteMeasurementSemanticsTable(db); err != nil {
		return err
	}
	if err := ensureSQLiteSourceRuleMixedLayoutColumns(db); err != nil {
		return err
	}
	if err := ensureSQLiteRecordingPlansTable(db); err != nil {
		return err
	}
	return ensureSQLiteDurableDeliveryTables(db)
}

func ensureSQLiteDurableDeliveryTables(db *sql.DB) error {
	const migrationName = "020_durable_delivery_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "gw_delivery_outbox")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table gw_delivery_outbox: %w", err)
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

func ensureSQLiteRecordingPlansTable(db *sql.DB) error {
	const migrationName = "019_recording_plans_sqlite.up.sql"

	exists, err := sqliteTableExists(db, "recording_plans")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table recording_plans for migration %s: %w", migrationName, err)
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
