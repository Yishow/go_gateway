package datalink

import (
	"database/sql"
	"fmt"
)

func ensureSQLiteMySQLPartitionMaintenance(db *sql.DB) error {
	exists, err := sqliteTableExists(db, "write_group_mysql_partition_maintenance")
	if err != nil {
		return fmt.Errorf("inspect mysql partition maintenance sidecar: %w", err)
	}
	if !exists {
		return applySQLiteMigrationOnce(db, "041_mysql_partition_maintenance_sqlite.up.sql")
	}
	return nil
}
