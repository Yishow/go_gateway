package datalink

import (
	"context"
	"database/sql"
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
		log.Printf("Executing SQLite migration: %s", targetFile)
		if _, err := db.ExecContext(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", targetFile, err)
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
