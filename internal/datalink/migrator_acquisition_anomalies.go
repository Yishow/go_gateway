package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

const (
	acquisitionAnomaliesMigration = "054_acquisition_anomalies_sqlite.up.sql"
	sqliteTableObject             = "table"
	sqliteIndexObject             = "index"
)

// ensureSQLiteAcquisitionAnomalies creates the acquisition_anomalies table and
// its indexes. Existing objects must match the migration exactly; an unknown
// schema is refused instead of being repaired or guessed.
func ensureSQLiteAcquisitionAnomalies(db *sql.DB) error {
	return runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		raw, err := migrations.FS.ReadFile(acquisitionAnomaliesMigration)
		if err != nil {
			return err
		}
		statements := strings.Split(strings.TrimSpace(string(raw)), ";")
		objects := []struct{ kind, name, statement string }{
			{sqliteTableObject, "acquisition_anomalies", statements[0]},
			{sqliteIndexObject, "idx_acquisition_anomalies_active", statements[1]},
			{sqliteIndexObject, "idx_acquisition_anomalies_device_started", statements[2]},
		}
		for _, object := range objects {
			if err := expandSQLiteExactObject(ctx, tx, object.kind, object.name, strings.TrimSpace(object.statement)); err != nil {
				return err
			}
		}
		return nil
	})
}

func expandSQLiteExactObject(ctx context.Context, tx *sql.Tx, kind, name, expected string) error {
	var actualKind, ddl string
	err := tx.QueryRowContext(ctx, `SELECT type, sql FROM sqlite_master WHERE name = ?`, name).Scan(&actualKind, &ddl)
	if errors.Is(err, sql.ErrNoRows) {
		if _, err := tx.ExecContext(ctx, expected); err != nil {
			return fmt.Errorf("create %s %s: %w", kind, name, err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect %s: %w", name, err)
	}
	if actualKind != kind || canonicalMixedDeliveryDDL(ddl) != canonicalMixedDeliveryDDL(expected) {
		return fmt.Errorf("unknown acquisition_anomalies schema (%s); migration refused", name)
	}
	return nil
}
