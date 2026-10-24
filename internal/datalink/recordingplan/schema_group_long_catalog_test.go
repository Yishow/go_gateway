package recordingplan

import (
	"context"
	"database/sql"
	"fmt"
)

// sqliteInspection reads the real SQLite catalog the way the production
// inspector does: table_info for columns and keys, index_list for indexes.
func sqliteInspection(ctx context.Context, db *sql.DB, table string) (TargetTableInspection, error) {
	state := TargetTableInspection{
		ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{},
	}
	rows, err := db.QueryContext(ctx, fmt.Sprintf(`PRAGMA table_info('%s')`, table))
	if err != nil {
		return state, err
	}
	for rows.Next() {
		var cid, notNull, pk int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			_ = rows.Close()
			return state, err
		}
		state.Columns = append(state.Columns, name)
		state.ColumnTypes[name], state.ColumnNullable[name], state.ColumnPrimaryKey[name] = dataType, notNull == 0, pk > 0
	}
	if err := rows.Close(); err != nil {
		return state, err
	}
	if len(state.Columns) == 0 {
		state.Status = inspectionMissing
		return state, nil
	}
	state.Status = inspectionExists

	state.IndexesInspected = true
	listed, err := db.QueryContext(ctx, fmt.Sprintf(`PRAGMA index_list('%s')`, table))
	if err != nil {
		return state, err
	}
	type listedIndex struct {
		name   string
		unique bool
	}
	var indexes []listedIndex
	for listed.Next() {
		var seq, unique, partial int
		var name, origin string
		if err := listed.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			_ = listed.Close()
			return state, err
		}
		if origin == "c" && partial == 0 {
			indexes = append(indexes, listedIndex{name: name, unique: unique == 1})
		}
	}
	if err := listed.Close(); err != nil {
		return state, err
	}
	for _, index := range indexes {
		info, err := db.QueryContext(ctx, fmt.Sprintf(`PRAGMA index_info('%s')`, index.name))
		if err != nil {
			return state, err
		}
		entry := TargetIndexInspection{Name: index.name, Unique: index.unique}
		for info.Next() {
			var seqno, cid int
			var name sql.NullString
			if err := info.Scan(&seqno, &cid, &name); err != nil {
				_ = info.Close()
				return state, err
			}
			entry.Columns = append(entry.Columns, name.String)
		}
		if err := info.Close(); err != nil {
			return state, err
		}
		state.Indexes = append(state.Indexes, entry)
	}
	return state, nil
}
