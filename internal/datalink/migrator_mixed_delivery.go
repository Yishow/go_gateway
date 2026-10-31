package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// ensureSQLiteMixedDelivery retains the existing 050 predecessor and registers
// the additive 051 expansion without changing historical migration bytes.
func ensureSQLiteMixedDelivery(db *sql.DB) error {
	if err := ensureSQLiteMemberCadence(db); err != nil {
		return err
	}
	return expandSQLiteMixedDelivery(db, nil)
}

func expandSQLiteMixedDelivery(db *sql.DB, checkpoint func(string) error) error {
	return runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		raw, err := migrations.FS.ReadFile("051_mixed_delivery_sqlite.up.sql")
		if err != nil {
			return err
		}
		statements := strings.Split(strings.TrimSpace(string(raw)), ";")
		step := func(name string) error {
			if checkpoint != nil {
				return checkpoint(name)
			}
			return nil
		}
		for i, name := range []string{"wg_delivery_mixed_progress", "wg_delivery_mixed_outcomes"} {
			expected := strings.TrimSpace(statements[i])
			var actual string
			err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE name=? AND type='table'`, name).Scan(&actual)
			if errors.Is(err, sql.ErrNoRows) {
				if _, err = tx.ExecContext(ctx, expected); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else if canonicalMixedDeliveryDDL(actual) != canonicalMixedDeliveryDDL(expected) {
				return fmt.Errorf("unknown %s schema; migration refused", name)
			}
			if err = step(name); err != nil {
				return err
			}
		}
		legacy, err := migrations.FS.ReadFile("025_write_group_delivery_sqlite.up.sql")
		if err != nil {
			return err
		}
		const header = "CREATE TABLE IF NOT EXISTS wg_delivery_outbox ("
		_, body, ok := strings.Cut(string(legacy), header)
		if !ok {
			return fmt.Errorf("missing historical outbox declaration")
		}
		body, _, ok = strings.Cut(body, ";")
		if !ok {
			return fmt.Errorf("incomplete historical outbox declaration")
		}
		original := "CREATE TABLE wg_delivery_outbox (" + body
		expanded := strings.TrimSuffix(strings.TrimSpace(original), ")") + ", due_at TEXT)"
		var actual string
		if err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE type='table' AND name='wg_delivery_outbox'`).Scan(&actual); err != nil {
			return err
		}
		switch canonicalMixedDeliveryDDL(actual) {
		case canonicalMixedDeliveryDDL(original):
			if _, err = tx.ExecContext(ctx, strings.TrimSpace(statements[2])); err != nil {
				return err
			}
		case canonicalMixedDeliveryDDL(expanded):
		default:
			return fmt.Errorf("unknown wg_delivery_outbox schema; migration refused")
		}
		if err = step("outbox"); err != nil {
			return err
		}
		const indexName = "idx_wg_delivery_outbox_due_head"
		expectedIndex := strings.TrimSpace(statements[3])
		var indexSQL string
		err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE name=? AND type='index'`, indexName).Scan(&indexSQL)
		if errors.Is(err, sql.ErrNoRows) {
			if _, err = tx.ExecContext(ctx, expectedIndex); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if canonicalMixedDeliveryDDL(indexSQL) != canonicalMixedDeliveryDDL(expectedIndex) {
			return fmt.Errorf("unknown %s schema; migration refused", indexName)
		}
		if err = step("index"); err != nil {
			return err
		}
		if err = expandMixedOwnershipFences(ctx, tx, step); err != nil {
			return err
		}
		if err = checkSQLiteMixedForeignKeys(ctx, tx); err != nil {
			return err
		}
		return step("commit")
	})
}

func canonicalMixedDeliveryDDL(s string) string {
	// SQLite removes IF NOT EXISTS and the statement terminator when persisting
	// CREATE SQL. All constraints, column order, defaults and FK actions remain
	// part of this exact comparison; unknown schema is never guessed/repaired.
	s = strings.Replace(s, "CREATE TABLE IF NOT EXISTS ", "CREATE TABLE ", 1)
	return strings.Join(strings.Fields(strings.TrimSuffix(strings.TrimSpace(s), ";")), " ")
}
