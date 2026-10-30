package datalink

import (
	"context"
	"database/sql"
	"fmt"
)

func sqliteMigrationHasTemp(ctx context.Context, conn *sql.Conn) (bool, error) {
	rows, err := conn.QueryContext(ctx, "PRAGMA database_list")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	hasTemp := false
	for rows.Next() {
		var seq int
		var name, path string
		if err = rows.Scan(&seq, &name, &path); err != nil {
			return false, err
		}
		hasTemp = hasTemp || name == "temp"
	}
	return hasTemp, rows.Err()
}

// ALTER TABLE RENAME opens an empty temp schema. Leaving it in a pooled
// connection changes database_list and prevents the existing destination guard
// from proving internal-file identity. Clear only a newly created EMPTY temp
// schema, preserving all pre-existing temporary data and the original setting.
func restoreSQLiteMigrationTemp(ctx context.Context, conn *sql.Conn, hadTemp bool, original int) error {
	if hadTemp {
		return nil
	}
	hasTemp, err := sqliteMigrationHasTemp(ctx, conn)
	if err != nil || !hasTemp {
		return err
	}
	var objects int
	if err = conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM temp.sqlite_master").Scan(&objects); err != nil {
		return err
	}
	if objects != 0 {
		return fmt.Errorf("migration unexpectedly created temporary objects")
	}
	alternate := 2
	if original == 2 {
		alternate = 1
	}
	if _, err = conn.ExecContext(ctx, fmt.Sprintf("PRAGMA temp_store=%d", alternate)); err != nil {
		return err
	}
	if _, err = conn.ExecContext(ctx, fmt.Sprintf("PRAGMA temp_store=%d", original)); err != nil {
		return err
	}
	var actual int
	if err = conn.QueryRowContext(ctx, "PRAGMA temp_store").Scan(&actual); err != nil {
		return err
	}
	if actual != original {
		return fmt.Errorf("migration temp_store restore mismatch")
	}
	hasTemp, err = sqliteMigrationHasTemp(ctx, conn)
	if err != nil {
		return err
	}
	if hasTemp {
		return fmt.Errorf("migration left a temporary schema")
	}
	return nil
}
