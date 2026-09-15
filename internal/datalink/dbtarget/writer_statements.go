package dbtarget

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

func buildGroupedWriteStatement(
	kind schema.DatabaseConnectorKind,
	key groupedWriteKey,
	values map[string]any,
) (statement string, parameters []any, buildErr error) {
	if len(values) == 0 {
		return "", nil, fmt.Errorf("grouped write has no buffered values")
	}

	columnNames := make([]string, 0, len(values))
	for columnName := range values {
		columnNames = append(columnNames, columnName)
	}
	sort.Strings(columnNames)

	columns := make([]string, 0, len(columnNames)+1)
	args := make([]any, 0, len(columnNames)+1)
	if key.WriteMode == schema.DatabaseWriteModeUpsert && strings.TrimSpace(key.TimestampColumn) != "" {
		columns = append(columns, quoteIdentifier(kind, key.TimestampColumn))
		args = append(args, key.BucketStart)
	}
	for _, columnName := range columnNames {
		columns = append(columns, quoteIdentifier(kind, columnName))
		args = append(args, values[columnName])
	}

	placeholders := buildPlaceholders(kind, len(args))
	tableRef := qualifiedTableName(kind, key.SchemaName, key.TableName)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableRef,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if key.WriteMode == schema.DatabaseWriteModeUpsert {
		if strings.TrimSpace(key.TimestampColumn) == "" {
			return "", nil, fmt.Errorf("upsert 模式缺少 timestamp_column")
		}
		assignments := make([]string, 0, len(columnNames))
		for _, columnName := range columnNames {
			quoted := quoteIdentifier(kind, columnName)
			if kind == schema.DatabaseConnectorKindMySQL {
				assignments = append(assignments, fmt.Sprintf("%s = VALUES(%s)", quoted, quoted))
				continue
			}
			assignments = append(assignments, fmt.Sprintf("%s = EXCLUDED.%s", quoted, quoted))
		}
		if kind == schema.DatabaseConnectorKindMySQL {
			query += " ON DUPLICATE KEY UPDATE " + strings.Join(assignments, ", ")
		} else {
			query += fmt.Sprintf(
				" ON CONFLICT (%s) DO UPDATE SET %s",
				quoteIdentifier(kind, key.TimestampColumn),
				strings.Join(assignments, ", "),
			)
		}
	}
	return query, args, nil
}

func buildWriteStatement(
	kind schema.DatabaseConnectorKind,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) (statement string, parameters []any, buildErr error) {
	if mapping.TableName == "" || mapping.ColumnName == "" {
		return "", nil, fmt.Errorf("資料庫目標映射缺少資料表或欄位資訊")
	}

	columns := []string{quoteIdentifier(kind, mapping.ColumnName)}
	args := []any{normalizeDBValue(value)}
	if mapping.WriteMode == schema.DatabaseWriteModeUpsert && mapping.TimestampColumn != nil {
		columns = append(columns, quoteIdentifier(kind, *mapping.TimestampColumn))
		args = append(args, observedAt.UTC())
	}

	placeholders := buildPlaceholders(kind, len(args))
	tableRef := qualifiedTableName(kind, mapping.TableSchema, mapping.TableName)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableRef,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if mapping.WriteMode == schema.DatabaseWriteModeUpsert {
		if mapping.TimestampColumn == nil {
			return "", nil, fmt.Errorf("upsert 模式缺少 timestamp_column")
		}

		valueColumn := quoteIdentifier(kind, mapping.ColumnName)
		timestampColumn := quoteIdentifier(kind, *mapping.TimestampColumn)
		if kind == schema.DatabaseConnectorKindMySQL {
			query += fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s = VALUES(%s)", valueColumn, valueColumn)
		} else {
			query += fmt.Sprintf(
				" ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s",
				timestampColumn,
				valueColumn,
				valueColumn,
			)
		}
	}

	return query, args, nil
}

func buildPlaceholders(kind schema.DatabaseConnectorKind, count int) []string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		switch kind {
		case schema.DatabaseConnectorKindPostgres:
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		default:
			placeholders = append(placeholders, "?")
		}
	}
	return placeholders
}

func quoteIdentifier(kind schema.DatabaseConnectorKind, name string) string {
	if kind == schema.DatabaseConnectorKindMySQL {
		escaped := strings.ReplaceAll(strings.TrimSpace(name), "`", "``")
		return "`" + escaped + "`"
	}
	escaped := strings.ReplaceAll(strings.TrimSpace(name), `"`, `""`)
	return `"` + escaped + `"`
}

func qualifiedTableName(kind schema.DatabaseConnectorKind, schemaName, tableName string) string {
	quotedTable := quoteIdentifier(kind, tableName)
	if strings.TrimSpace(schemaName) == "" || kind == schema.DatabaseConnectorKindSQLite {
		return quotedTable
	}
	return quoteIdentifier(kind, schemaName) + "." + quotedTable
}

func normalizeDBValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]any, []any:
		data, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(data)
	case time.Time:
		return typed.UTC()
	default:
		return value
	}
}
