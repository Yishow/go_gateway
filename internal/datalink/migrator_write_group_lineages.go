package datalink

import (
	"context"
	"database/sql"
	"fmt"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema/migrations"
)

type sqliteLineageReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// sqliteWriteGroupLineageState uses the key references as durable markers, so
// missing registries cannot be mistaken for an old database and silently rebuilt.
func sqliteWriteGroupLineageState(ctx context.Context, reader sqliteLineageReader) (initialized bool, err error) {
	var count int
	err = reader.QueryRowContext(ctx, `SELECT
		EXISTS(SELECT 1 FROM pragma_table_info('write_group_basic_keys') WHERE name='lineage_id') +
		EXISTS(SELECT 1 FROM pragma_table_info('write_group_profile_keys') WHERE name='lineage_id') +
		EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='write_group_lineages') +
		EXISTS(SELECT 1 FROM sqlite_master WHERE type='table' AND name='write_group_generations')`).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("inspect write-group lineage migration state: %w", err)
	}
	switch count {
	case 0:
		return false, nil
	case 4:
		return true, nil
	default:
		return false, fmt.Errorf("write-group lineage migration is incomplete: %d of 4 artifacts", count)
	}
}

// ensureSQLiteWriteGroupLineages atomically expands canonical keys once. Reentry
// only validates; it never repairs references, resets heads, or invents history.
func ensureSQLiteWriteGroupLineages(db *sql.DB) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin write-group lineage migration: %w", err)
	}
	defer tx.Rollback()
	// Take the write lock before schema inspection, including on concurrent startup.
	if _, err := tx.ExecContext(ctx, `UPDATE write_groups SET revision=revision WHERE 0`); err != nil {
		return fmt.Errorf("lock write-group lineage migration: %w", err)
	}
	initialized, err := sqliteWriteGroupLineageState(ctx, tx)
	if err != nil {
		return err
	}
	if err := validateSQLiteLineageKeyScopes(ctx, tx); err != nil {
		return err
	}
	if !initialized {
		if err := initializeSQLiteWriteGroupLineages(ctx, tx); err != nil {
			return err
		}
	}
	if err := validateSQLiteWriteGroupLineages(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit write-group lineage migration: %w", err)
	}
	return nil
}

func initializeSQLiteWriteGroupLineages(ctx context.Context, tx *sql.Tx) error {
	raw, err := migrations.FS.ReadFile("047_write_group_lineages_sqlite.up.sql")
	if err != nil {
		return fmt.Errorf("read write-group lineage migration: %w", err)
	}
	if _, err := tx.ExecContext(ctx, string(raw)); err != nil {
		return fmt.Errorf("expand write-group lineage schema: %w", err)
	}
	rows, err := tx.QueryContext(ctx, `SELECT workspace_id,group_id FROM write_group_basic_keys
		UNION SELECT workspace_id,group_id FROM write_group_profile_keys ORDER BY group_id`)
	if err != nil {
		return fmt.Errorf("read write-group lineage genesis keys: %w", err)
	}
	type anchor struct{ workspace, group string }
	var anchors []anchor
	for rows.Next() {
		var a anchor
		if err := rows.Scan(&a.workspace, &a.group); err != nil {
			rows.Close()
			return fmt.Errorf("scan write-group lineage genesis key: %w", err)
		}
		anchors = append(anchors, a)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate write-group lineage genesis keys: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close write-group lineage genesis keys: %w", err)
	}
	for _, a := range anchors {
		id, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("generate write-group lineage identity: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES (?,?,?)`, id, a.workspace, a.group); err != nil {
			return fmt.Errorf("insert write-group lineage: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES (?,?,1)`, a.group, id); err != nil {
			return fmt.Errorf("insert write-group genesis: %w", err)
		}
		for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys"} {
			if _, err := tx.ExecContext(ctx, `UPDATE `+table+` SET lineage_id=? WHERE group_id=?`, id, a.group); err != nil {
				return fmt.Errorf("bind write-group lineage key: %w", err)
			}
		}
	}
	return nil
}
