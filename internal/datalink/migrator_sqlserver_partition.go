package datalink

import (
	"database/sql"
	"fmt"
)

func ensureSQLiteSQLServerPartitionMaintenance(db *sql.DB) error {
	exists, err := sqliteTableExists(db, "write_group_sqlserver_partition_maintenance")
	if err != nil {
		return fmt.Errorf("inspect sqlserver partition maintenance sidecar: %w", err)
	}
	if !exists {
		return applySQLiteMigrationOnce(db, "042_sqlserver_partition_maintenance_sqlite.up.sql")
	}
	return nil
}
