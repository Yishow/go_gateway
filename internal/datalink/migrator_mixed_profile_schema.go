package datalink

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Known columns are the immutable historical write-group schema plus nullable
// expansions. Unknown columns/generation rules need an explicit new migration.
func mixedProfileColumns(ctx context.Context, tx *sql.Tx) ([]string, error) {
	known := strings.Fields("id workspace_id revision applied_revision name status destination_connector_id destination_connector_revision destination_database destination_table_schema destination_table_name destination_storage_strategy destination_schema_revision destination_schema_digest row_policy write_policy migration created_at updated_at layout_version profile_kind profile_interval profile_unit storage_mode sample_reduction result_shape")
	allowed := make(map[string]bool, len(known))
	for _, name := range known {
		allowed[name] = true
	}
	rows, err := tx.QueryContext(ctx, "PRAGMA table_xinfo(write_groups)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []string
	seen := make(map[string]bool)
	for rows.Next() {
		var cid, notNull, pk, hidden int
		var name, kind string
		var defaultValue any
		if err = rows.Scan(&cid, &name, &kind, &notNull, &defaultValue, &pk, &hidden); err != nil {
			return nil, err
		}
		if !allowed[name] || hidden != 0 {
			return nil, fmt.Errorf("unknown write_groups column %q", name)
		}
		seen[name] = true
		columns = append(columns, "\""+name+"\"")
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	for _, name := range known[:19] {
		if !seen[name] {
			return nil, fmt.Errorf("incomplete write_groups schema: %s", name)
		}
	}
	return columns, nil
}

func mixedProfileObjects(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT sql FROM sqlite_master WHERE tbl_name='write_groups' AND type IN ('index','trigger') AND sql IS NOT NULL ORDER BY type,name`)
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
