package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// The caller owns the dedicated connection, FK/rename pragmas and transaction.
// No commit or repair of existing lineage metadata is performed here.
func migrateMixedRecordProfileKeysTx(ctx context.Context, tx *sql.Tx) error {
	return rebuildRecordProfileKeys(ctx, tx, nil)
}

func rebuildRecordProfileKeys(ctx context.Context, tx *sql.Tx, checkpoint func(string) error) error {
	legacy, expanded, err := recordProfileKeySchemas()
	if err != nil {
		return err
	}
	var original string
	if err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name='write_group_profile_keys'`).Scan(&original); err != nil {
		return err
	}
	normalize := func(s string) string {
		s = strings.Replace(s, `CREATE TABLE "write_group_profile_keys"`, "CREATE TABLE write_group_profile_keys", 1)
		return strings.Join(strings.Fields(s), " ")
	}
	if normalize(original) == normalize(expanded) {
		return checkSQLiteMixedForeignKeys(ctx, tx)
	}
	if normalize(original) != normalize(legacy) {
		return fmt.Errorf("unknown write_group_profile_keys schema; migration refused")
	}
	const temporary = "write_group_profile_keys_record_expand"
	var count int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE name=?`, temporary).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("record profile migration temporary table already exists")
	}
	objects, err := recordProfileKeyObjects(ctx, tx)
	if err != nil {
		return err
	}
	step := func(name string) error {
		if checkpoint != nil {
			return checkpoint(name)
		}
		return nil
	}
	if err = step("copy"); err != nil {
		return err
	}
	ddl := strings.Replace(expanded, "CREATE TABLE write_group_profile_keys", "CREATE TABLE "+temporary, 1)
	if _, err = tx.ExecContext(ctx, ddl); err != nil {
		return err
	}
	columns := strings.Fields("workspace_id circuit_id selection_scope profile_kind canonical_role group_id created_at lineage_id")
	fields := strings.Join(columns, ",")
	if _, err = tx.ExecContext(ctx, "INSERT INTO "+temporary+" ("+fields+") SELECT "+fields+" FROM write_group_profile_keys"); err != nil {
		return err
	}
	var comparisons []string
	for _, column := range columns {
		comparisons = append(comparisons, "a."+column+" IS NOT b."+column+" OR typeof(a."+column+")<>typeof(b."+column+")")
	}
	query := "SELECT (SELECT COUNT(*) FROM write_group_profile_keys a LEFT JOIN " + temporary + " b ON a.group_id=b.group_id WHERE b.group_id IS NULL OR " + strings.Join(comparisons, " OR ") + ") + ABS((SELECT COUNT(*) FROM write_group_profile_keys)-(SELECT COUNT(*) FROM " + temporary + "))"
	if err = tx.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("record profile migration copy verification failed")
	}
	if err = step("replace"); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "DROP TABLE write_group_profile_keys; ALTER TABLE "+temporary+" RENAME TO write_group_profile_keys"); err != nil {
		return err
	}
	if err = step("indexes"); err != nil {
		return err
	}
	for _, object := range objects {
		if _, err = tx.ExecContext(ctx, object); err != nil {
			return err
		}
	}
	if err = step("foreign_keys"); err != nil {
		return err
	}
	if err = checkSQLiteMixedForeignKeys(ctx, tx); err != nil {
		return err
	}
	return step("commit")
}

func recordProfileKeyObjects(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT sql FROM sqlite_master WHERE tbl_name='write_group_profile_keys' AND type IN ('index','trigger') AND sql IS NOT NULL ORDER BY type,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var objects []string
	for rows.Next() {
		var ddl string
		if err = rows.Scan(&ddl); err != nil {
			return nil, err
		}
		objects = append(objects, ddl)
	}
	return objects, rows.Err()
}
