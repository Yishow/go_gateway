package dbtarget

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

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
