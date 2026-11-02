package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const postgresDefaultSchemaName = "public"

const postgresInsufficientPrivilegeCode = "42501"

func scanPostgresColumn(rows *sql.Rows, column *ColumnInfo) error {
	var dataType string
	var precision sql.NullInt64
	var scale sql.NullInt64
	if err := rows.Scan(
		&column.Name,
		&dataType,
		&precision,
		&scale,
		&column.Nullable,
		&column.PrimaryKey,
		&column.Unique,
		&column.HasDefault,
	); err != nil {
		return err
	}

	var err error
	column.DataType, err = postgresColumnDataType(dataType, precision, scale)
	return err
}

func postgresColumnDataType(dataType string, precision, scale sql.NullInt64) (string, error) {
	if !strings.EqualFold(strings.TrimSpace(dataType), "numeric") {
		return dataType, nil
	}
	if !precision.Valid && !scale.Valid {
		return dataType, nil
	}
	if !precision.Valid || !scale.Valid || precision.Int64 <= 0 {
		return "", errors.New("numeric precision/scale metadata is incomplete")
	}
	return fmt.Sprintf("numeric(%d,%d)", precision.Int64, scale.Int64), nil
}

func inspectPostgresTables(ctx context.Context, db *sql.DB) ([]TableInfo, error) {
	rows, err := db.QueryContext(
		ctx,
		`
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_type = 'BASE TABLE'
			  AND table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("查詢 postgres 資料表失敗: %w", err)
	}
	defer rows.Close()

	tables := make([]TableInfo, 0, 8)
	for rows.Next() {
		var schemaName string
		var tableName string
		if err := rows.Scan(&schemaName, &tableName); err != nil {
			return nil, fmt.Errorf("掃描 postgres 資料表失敗: %w", err)
		}

		columnRows, err := db.QueryContext(
			ctx,
			`
				SELECT
					c.column_name,
					c.data_type,
					c.numeric_precision,
					c.numeric_scale,
					(c.is_nullable = 'YES') AS is_nullable,
					EXISTS (
						SELECT 1
						FROM information_schema.table_constraints tc
						JOIN information_schema.key_column_usage kcu
						  ON tc.constraint_name = kcu.constraint_name
						 AND tc.table_schema = kcu.table_schema
						 AND tc.table_name = kcu.table_name
						WHERE tc.constraint_type = 'PRIMARY KEY'
						  AND tc.table_schema = c.table_schema
						  AND tc.table_name = c.table_name
						  AND kcu.column_name = c.column_name
					) AS is_primary_key,
					EXISTS (
						SELECT 1
						FROM information_schema.table_constraints tc
						JOIN information_schema.key_column_usage kcu
						  ON tc.constraint_name = kcu.constraint_name
						 AND tc.table_schema = kcu.table_schema
						 AND tc.table_name = kcu.table_name
						WHERE tc.constraint_type IN ('PRIMARY KEY', 'UNIQUE')
						  AND tc.table_schema = c.table_schema
						  AND tc.table_name = c.table_name
						  AND kcu.column_name = c.column_name
						  AND (
							SELECT COUNT(*)
							FROM information_schema.key_column_usage kcu_count
							WHERE kcu_count.constraint_name = tc.constraint_name
							  AND kcu_count.table_schema = tc.table_schema
							  AND kcu_count.table_name = tc.table_name
						  ) = 1
					) AS is_unique, (c.column_default IS NOT NULL OR c.is_identity = 'YES' OR c.is_generated = 'ALWAYS') AS has_default
				FROM information_schema.columns c
				WHERE c.table_schema = $1
				  AND c.table_name = $2
				ORDER BY c.ordinal_position
			`,
			schemaName,
			tableName,
		)
		if err != nil {
			return nil, fmt.Errorf("查詢 postgres 欄位失敗: %w", err)
		}

		table := TableInfo{
			Schema: schemaName,
			Name:   tableName,
		}
		for columnRows.Next() {
			var column ColumnInfo
			if err := scanPostgresColumn(columnRows, &column); err != nil {
				columnRows.Close()
				return nil, fmt.Errorf("掃描 postgres 欄位失敗: %w", err)
			}
			table.Columns = append(table.Columns, column)
		}
		if err := columnRows.Close(); err != nil {
			return nil, fmt.Errorf("關閉 postgres 欄位查詢失敗: %w", err)
		}

		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷 postgres 資料表失敗: %w", err)
	}

	if len(tables) == 0 {
		tables = nil
	}
	return tables, nil
}
