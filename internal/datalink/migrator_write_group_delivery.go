package datalink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema/migrations"
)

// ensureSQLiteWriteGroupDeliveryTables creates the durable write-group
// delivery tables once; every table is CREATE IF NOT EXISTS.
func ensureSQLiteWriteGroupDeliveryTables(db *sql.DB) error {
	const migrationName = "025_write_group_delivery_sqlite.up.sql"

	for _, table := range []string{
		"wg_delivery_samples", "wg_delivery_checkpoints", "wg_delivery_buckets",
		"wg_delivery_outbox", "wg_delivery_receipts",
	} {
		exists, err := sqliteTableExists(db, table)
		if err != nil {
			return fmt.Errorf("failed to inspect sqlite table %s for migration %s: %w", table, migrationName, err)
		}
		if !exists {
			if err := applySQLiteMigrationOnce(db, migrationName); err != nil {
				return err
			}
			break
		}
	}
	if err := ensureSQLiteCounterBaselineTable(db); err != nil {
		return err
	}
	if err := ensureSQLiteCounterProofAssociations(db); err != nil {
		return err
	}
	return ensureSQLiteSampleHandoffs(db)
}

// ensureSQLiteCounterBaselineTable adds the durable counter baseline table
// beside the delivery tables of an installation that already has them.
func ensureSQLiteCounterBaselineTable(db *sql.DB) error {
	const migrationName = "033_write_group_counter_baselines_sqlite.up.sql"
	exists, err := sqliteTableExists(db, "wg_delivery_counter_baselines")
	if err != nil {
		return fmt.Errorf("failed to inspect sqlite table wg_delivery_counter_baselines for migration %s: %w", migrationName, err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, migrationName)
}

// ensureSQLiteCounterProofAssociations adds audit metadata without rewriting history.
func ensureSQLiteCounterProofAssociations(db *sql.DB) error {
	exists, err := sqliteTableExists(db, "wg_counter_proof_associations")
	if err != nil {
		return fmt.Errorf("inspect counter proof associations: %w", err)
	}
	if exists {
		return nil
	}
	return applySQLiteMigrationOnce(db, "043_counter_proof_associations_sqlite.up.sql")
}

// ensureSQLiteSampleHandoffs handles the mixed additive-column/table migration
// transactionally and idempotently, including a retry after interrupted startup.
func ensureSQLiteSampleHandoffs(db *sql.DB) error {
	const name = "045_write_group_sample_handoffs_sqlite.up.sql"
	raw, err := migrations.FS.ReadFile(name)
	if err != nil {
		return err
	}
	var lines []string
	for _, line := range strings.Split(string(raw), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "--") {
			lines = append(lines, line)
		}
	}
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec("UPDATE wg_delivery_checkpoints SET updated_at=updated_at WHERE 0"); err != nil {
		return err
	}
	for _, part := range strings.Split(strings.Join(lines, "\n"), ";") {
		statement := strings.TrimSpace(part)
		if statement == "" {
			continue
		}
		// 046 supersedes only this receipt index. Never recreate its old unique
		// form during 045 reentry: accepted aliases may already share a target.
		if strings.HasPrefix(statement, "CREATE UNIQUE INDEX IF NOT EXISTS idx_wg_handoff_target_seq") {
			if err = ensureSQLiteHandoffConvergenceIndex(tx); err != nil {
				return err
			}
			continue
		}
		if match := sqliteAddColumnStatement.FindStringSubmatch(statement); match != nil {
			rows, err := tx.Query("PRAGMA table_info(" + match[1] + ")")
			if err != nil {
				return err
			}
			exists := false
			for rows.Next() {
				var ordinal, notNull, primary int
				var column, kind string
				var defaultValue any
				if err = rows.Scan(&ordinal, &column, &kind, &notNull, &defaultValue, &primary); err != nil {
					rows.Close()
					return err
				}
				if column == match[2] {
					exists = true
				}
			}
			if err = rows.Err(); err != nil {
				rows.Close()
				return err
			}
			if err = rows.Close(); err != nil {
				return err
			}
			if exists {
				continue
			}
		}
		if _, err = tx.Exec(statement); err != nil {
			return fmt.Errorf("apply sample handoff migration: %w", err)
		}
	}
	return tx.Commit()
}

// ensureSQLiteHandoffConvergenceIndex preserves every journal and receipt while
// allowing independently verified transfer edges to share a canonical target.
// Its DDL runs in the same transaction as 045, including startup retries.
func ensureSQLiteHandoffConvergenceIndex(tx *sql.Tx) error {
	var unique int
	err := tx.QueryRow(`SELECT "unique" FROM pragma_index_list('wg_delivery_sample_handoffs') WHERE name='idx_wg_handoff_target_seq'`).Scan(&unique)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("inspect handoff convergence index: %w", err)
	}
	if err == nil && unique == 0 {
		return nil
	}
	raw, err := migrations.FS.ReadFile("046_write_group_handoff_convergence_sqlite.up.sql")
	if err != nil {
		return err
	}
	for _, part := range strings.Split(string(raw), ";") {
		if statement := strings.TrimSpace(part); statement != "" {
			if _, err = tx.Exec(statement); err != nil {
				return fmt.Errorf("apply handoff convergence migration: %w", err)
			}
		}
	}
	return nil
}
