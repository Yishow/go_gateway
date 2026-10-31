package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// expandMixedOwnershipFences runs only inside the existing owned-connection
// configuration expansion transaction. 044 remains immutable.
func expandMixedOwnershipFences(ctx context.Context, tx *sql.Tx, step func(string) error) error {
	raw, err := migrations.FS.ReadFile("044_write_group_ownership_fences_sqlite.up.sql")
	if err != nil {
		return err
	}
	const table = "write_group_ownership_fences"
	const header = "CREATE TABLE IF NOT EXISTS " + table + " ("
	_, tail, ok := strings.Cut(string(raw), header)
	if !ok {
		return fmt.Errorf("missing historical ownership fence schema")
	}
	body, _, ok := strings.Cut(tail, ";")
	if !ok {
		return fmt.Errorf("incomplete historical ownership fence schema")
	}
	legacy := "CREATE TABLE " + table + " (" + body
	expanded := strings.Replace(legacy, "profile_kind IN ('', 'instantaneous', 'cumulative')", "profile_kind IN ('', 'instantaneous', 'cumulative', 'mixed')", 1)
	normalize := func(s string) string {
		return canonicalMixedDeliveryDDL(strings.Replace(s, `CREATE TABLE "`+table+`"`, "CREATE TABLE "+table, 1))
	}
	var actual string
	if err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&actual); err != nil {
		return err
	}
	if normalize(actual) == normalize(expanded) {
		return nil
	}
	if normalize(actual) != normalize(legacy) {
		return fmt.Errorf("unknown ownership fence schema; migration refused")
	}
	const temporary = "write_group_ownership_fences_mixed_expand"
	var count int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE name=?`, temporary).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("ownership fence migration temporary table exists")
	}
	before, err := mixedFenceDependencies(ctx, tx)
	if err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT sql FROM sqlite_master WHERE tbl_name=? AND type IN ('index','trigger') AND sql IS NOT NULL ORDER BY type,name`, table)
	if err != nil {
		return err
	}
	var objects []string
	for rows.Next() {
		var ddl string
		if err = rows.Scan(&ddl); err != nil {
			rows.Close()
			return err
		}
		objects = append(objects, ddl)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err = rows.Close(); err != nil {
		return err
	}
	if err = step("fence_copy"); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, strings.Replace(expanded, "CREATE TABLE "+table, "CREATE TABLE "+temporary, 1)); err != nil {
		return err
	}
	columns := strings.Fields("fence_revision workspace_id direction group_id group_revision device_id point_id tag_id profile_kind connector_id connector_revision database_name table_schema table_name")
	fields := strings.Join(columns, ",")
	if _, err = tx.ExecContext(ctx, "INSERT INTO "+temporary+" ("+fields+") SELECT "+fields+" FROM "+table); err != nil {
		return err
	}
	var comparisons, joins []string
	for _, column := range columns {
		comparisons = append(comparisons, "a."+column+" IS NOT b."+column+" OR typeof(a."+column+")<>typeof(b."+column+")")
	}
	for _, column := range strings.Fields("fence_revision direction group_id group_revision device_id point_id tag_id") {
		joins = append(joins, "a."+column+"=b."+column)
	}
	query := "SELECT (SELECT COUNT(*) FROM " + table + " a LEFT JOIN " + temporary + " b ON " + strings.Join(joins, " AND ") + " WHERE b.fence_revision IS NULL OR " + strings.Join(comparisons, " OR ") + ") + ABS((SELECT COUNT(*) FROM " + table + ")-(SELECT COUNT(*) FROM " + temporary + "))"
	if err = tx.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("ownership fence copy verification failed")
	}
	if err = step("fence_replace"); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "DROP TABLE "+table+"; ALTER TABLE "+temporary+" RENAME TO "+table); err != nil {
		return err
	}
	if err = step("fence_indexes"); err != nil {
		return err
	}
	for _, ddl := range objects {
		if _, err = tx.ExecContext(ctx, ddl); err != nil {
			return err
		}
	}
	after, err := mixedFenceDependencies(ctx, tx)
	if err != nil {
		return err
	}
	if !maps.Equal(before, after) {
		return fmt.Errorf("ownership fence dependencies changed")
	}
	if err = step("fence_foreign_keys"); err != nil {
		return err
	}
	if err = checkSQLiteMixedForeignKeys(ctx, tx); err != nil {
		return err
	}
	return step("fence_commit")
}

// Comparing every other schema object preserves referencing views and FK table
// declarations as well as the original indexes/triggers, not just row counts.
func mixedFenceDependencies(ctx context.Context, tx *sql.Tx) (map[string]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT type,name,tbl_name,COALESCE(sql,'') FROM sqlite_master WHERE NOT(type='table' AND name='write_group_ownership_fences') ORDER BY type,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]string{}
	for rows.Next() {
		var kind, name, table, ddl string
		if err = rows.Scan(&kind, &name, &table, &ddl); err != nil {
			return nil, err
		}
		result[kind+":"+name] = table + "\x00" + ddl
	}
	return result, rows.Err()
}
