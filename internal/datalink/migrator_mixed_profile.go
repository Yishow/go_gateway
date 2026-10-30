package datalink

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// expandSQLiteMixedProfile preserves the canonical table and all references.
// The dedicated connection never returns to the pool with altered pragmas.
func expandSQLiteMixedProfile(db *sql.DB, checkpoint func(string) error) (err error) {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	var foreignKeys, legacyAlter int
	if err = conn.QueryRowContext(ctx, "PRAGMA foreign_keys").Scan(&foreignKeys); err != nil {
		return err
	}
	if err = conn.QueryRowContext(ctx, "PRAGMA legacy_alter_table").Scan(&legacyAlter); err != nil {
		return err
	}
	hadTemp, err := sqliteMigrationHasTemp(ctx, conn)
	if err != nil {
		return err
	}
	var tempStore int
	if err = conn.QueryRowContext(ctx, "PRAGMA temp_store").Scan(&tempStore); err != nil {
		return err
	}
	defer func() {
		tempErr := restoreSQLiteMigrationTemp(ctx, conn, hadTemp, tempStore)
		restore := restoreSQLiteMigrationPragmas(
			[]sqliteMigrationPragma{{"legacy_alter_table", legacyAlter}, {"foreign_keys", foreignKeys}},
			func(query string) error { _, e := conn.ExecContext(ctx, query); return e },
			func(query string) (int, error) {
				var actual int
				e := conn.QueryRowContext(ctx, query).Scan(&actual)
				return actual, e
			},
			func() error { return conn.Raw(func(any) error { return driver.ErrBadConn }) },
		)
		if tempErr != nil {
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
			restore = errors.Join(restore, tempErr)
		}
		if restore != nil {
			err = errors.Join(err, fmt.Errorf("restore mixed migration connection: %w", restore))
		}
	}()
	if _, err = conn.ExecContext(ctx, "PRAGMA foreign_keys=OFF; PRAGMA legacy_alter_table=ON"); err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			if e := tx.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
				err = errors.Join(err, e)
			}
		}
	}()
	if _, err = tx.ExecContext(ctx, "UPDATE write_groups SET id=id WHERE 0"); err != nil {
		return err
	}
	if err = addSQLiteMemberCalculationColumns(ctx, tx); err != nil {
		return err
	}
	if err = rebuildMixedProfileTable(ctx, tx, checkpoint); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

func rebuildMixedProfileTable(ctx context.Context, tx *sql.Tx, checkpoint func(string) error) error {
	var original string
	if err := tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name='write_groups'`).Scan(&original); err != nil {
		return err
	}
	const old = "profile_kind TEXT CHECK (profile_kind IS NULL OR profile_kind IN ('instantaneous', 'cumulative'))"
	const expanded = "profile_kind TEXT CHECK (profile_kind IS NULL OR profile_kind IN ('instantaneous', 'cumulative', 'mixed'))"
	columns, err := mixedProfileColumns(ctx, tx)
	if err != nil {
		return err
	}
	if strings.Count(original, expanded) == 1 {
		return checkSQLiteMixedForeignKeys(ctx, tx)
	}
	if strings.Count(original, old) != 1 {
		return fmt.Errorf("unknown write_groups profile constraint; migration refused")
	}
	const temporary = "write_groups_member_calculation_expand"
	var exists int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE name=?`, temporary).Scan(&exists); err != nil {
		return err
	}
	if exists != 0 {
		return fmt.Errorf("mixed migration temporary table already exists")
	}
	// Only these exact canonical CREATE headers and the known historical column
	// declaration are accepted; arbitrary SQL is never searched and rewritten.
	header := ""
	for _, candidate := range []string{"CREATE TABLE write_groups (", "CREATE TABLE \"write_groups\" ("} {
		if strings.HasPrefix(original, candidate) {
			header = candidate
			break
		}
	}
	if header == "" {
		return fmt.Errorf("unknown write_groups table declaration")
	}
	ddl := strings.Replace(original, old, expanded, 1)
	ddl = "CREATE TABLE \"" + temporary + "\" (" + strings.TrimPrefix(ddl, header)
	objects, err := mixedProfileObjects(ctx, tx)
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
	if _, err = tx.ExecContext(ctx, ddl); err != nil {
		return err
	}
	fields := strings.Join(columns, ",")
	if _, err = tx.ExecContext(ctx, "INSERT INTO "+temporary+" ("+fields+") SELECT "+fields+" FROM write_groups"); err != nil {
		return err
	}
	var different int
	var comparisons []string
	for _, column := range columns {
		comparisons = append(comparisons, "a."+column+" IS NOT b."+column+" OR typeof(a."+column+")<>typeof(b."+column+")")
	}
	query := "SELECT (SELECT COUNT(*) FROM write_groups a LEFT JOIN " + temporary + " b ON a.id=b.id WHERE b.id IS NULL OR " + strings.Join(comparisons, " OR ") + ") + ABS((SELECT COUNT(*) FROM write_groups)-(SELECT COUNT(*) FROM " + temporary + "))"

	if err = tx.QueryRowContext(ctx, query).Scan(&different); err != nil {
		return err
	}
	if different != 0 {
		return fmt.Errorf("mixed migration copy verification failed")
	}
	if err = step("replace"); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "DROP TABLE write_groups; ALTER TABLE "+temporary+" RENAME TO write_groups"); err != nil {
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

func checkSQLiteMixedForeignKeys(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, "PRAGMA foreign_key_check")
	if err != nil {
		return err
	}
	violated := rows.Next()
	rowErr := rows.Err()
	closeErr := rows.Close()
	if rowErr != nil || closeErr != nil {
		return errors.Join(rowErr, closeErr)
	}
	if violated {
		return fmt.Errorf("mixed migration foreign key check failed")
	}
	return nil
}
