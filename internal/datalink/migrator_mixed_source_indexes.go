package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

func ensureSQLiteMixedSourceIndexes(db *sql.DB) error {
	if err := ensureSQLiteMixedDelivery(db); err != nil {
		return err
	}
	return expandSQLiteMixedSourceIndexes(db, nil)
}

// 052 only adds query indexes. It does not rebuild tables or rewrite retained
// outcomes. Unknown schema/objects and invalid JSON are errors, not repair input.
func expandSQLiteMixedSourceIndexes(db *sql.DB, checkpoint func(string) error) error {
	return runSQLiteConfigurationExpansion(db, func(ctx context.Context, tx *sql.Tx) error {
		predecessor, err := migrations.FS.ReadFile("051_mixed_delivery_sqlite.up.sql")
		if err != nil {
			return err
		}
		expectedTable := strings.TrimSpace(strings.Split(string(predecessor), ";")[1])
		var actual string
		if err = tx.QueryRowContext(ctx, `SELECT sql FROM sqlite_master WHERE name='wg_delivery_mixed_outcomes' AND type='table'`).Scan(&actual); err != nil {
			return err
		}
		if canonicalMixedDeliveryDDL(actual) != canonicalMixedDeliveryDDL(expectedTable) {
			return fmt.Errorf("unknown mixed outcomes schema; index migration refused")
		}
		var invalid int
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_mixed_outcomes WHERE CASE WHEN json_valid(outcome) THEN json_type(outcome) != 'object' ELSE 1 END`).Scan(&invalid); err != nil {
			return err
		}
		if invalid != 0 {
			return fmt.Errorf("invalid mixed outcome JSON; index migration refused")
		}
		raw, err := migrations.FS.ReadFile("052_mixed_source_indexes_sqlite.up.sql")
		if err != nil {
			return err
		}
		statements := strings.Split(strings.TrimSpace(string(raw)), ";")
		for i, name := range []string{"idx_wg_mixed_source_latest", "idx_wg_mixed_source_entity_latest"} {
			expected := strings.TrimSpace(statements[i])
			var kind, ddl string
			err = tx.QueryRowContext(ctx, `SELECT type,sql FROM sqlite_master WHERE name=?`, name).Scan(&kind, &ddl)
			if errors.Is(err, sql.ErrNoRows) {
				_, err = tx.ExecContext(ctx, expected)
			} else if err == nil && (kind != "index" || canonicalMixedDeliveryDDL(ddl) != canonicalMixedDeliveryDDL(expected)) {
				return fmt.Errorf("unknown %s schema; index migration refused", name)
			}
			if err != nil {
				return err
			}
			if checkpoint != nil {
				if err = checkpoint(name); err != nil {
					return err
				}
			}
		}
		if err = checkSQLiteMixedForeignKeys(ctx, tx); err != nil {
			return err
		}
		if checkpoint != nil {
			return checkpoint("commit")
		}
		return nil
	})
}
