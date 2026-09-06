package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func inspectTablesByKind(
	ctx context.Context,
	db *sql.DB,
	kind schema.DatabaseConnectorKind,
	config ConnectionConfig,
	extraSchemas []string,
) ([]TableInfo, error) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		return inspectSQLiteTables(ctx, db)
	case schema.DatabaseConnectorKindPostgres:
		return inspectPostgresTables(ctx, db)
	case schema.DatabaseConnectorKindMySQL:
		databaseName := defaultString(
			stringConfigValue(config, "database"),
			stringConfigValue(config, "dbname"),
		)
		if databaseName == "" {
			return nil, fmt.Errorf("mysql 連接設定缺少 database")
		}
		return inspectMySQLTables(ctx, db, mergeSchemaNames(databaseName, extraSchemas))
	default:
		return nil, fmt.Errorf("不支援的資料庫類型: %s", kind)
	}
}

// mergeSchemaNames 把連接器自身的資料庫與映射實際引用的資料庫合併去重。
// 映射的 table_schema 是使用者可設定的，只查連接器自己的資料庫會讓跨庫映射
// 永遠對不上既有資料表，並在每次 schema 產生時重複輸出建庫 / 建表語句。
func mergeSchemaNames(databaseName string, extraSchemas []string) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0, len(extraSchemas)+1)
	for _, name := range append([]string{databaseName}, extraSchemas...) {
		normalized := strings.TrimSpace(name)
		if normalized == "" {
			continue
		}
		key := strings.ToLower(normalized)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		names = append(names, normalized)
	}
	sort.Strings(names)
	return names
}

func inspectMySQLTables(ctx context.Context, db *sql.DB, databaseNames []string) ([]TableInfo, error) {
	if len(databaseNames) == 0 {
		return nil, fmt.Errorf("mysql 連接設定缺少 database")
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(databaseNames)), ",")
	args := make([]any, 0, len(databaseNames))
	for _, name := range databaseNames {
		args = append(args, name)
	}

	uniqueColumns, err := inspectMySQLSingleColumnUniqueKeys(ctx, db, placeholders, args)
	if err != nil {
		return nil, err
	}

	// #nosec G202 -- placeholders 只由 "?," 重複組成，schema 名稱一律以參數繫結傳入
	rows, err := db.QueryContext(
		ctx,
		`SELECT
			c.table_schema,
			c.table_name,
			c.column_name,
			c.data_type,
			c.is_nullable
		FROM information_schema.columns c
		JOIN information_schema.tables t
		  ON t.table_schema = c.table_schema
		 AND t.table_name = c.table_name
		WHERE t.table_type = 'BASE TABLE'
		  AND c.table_schema IN (`+placeholders+`)
		ORDER BY c.table_schema, c.table_name, c.ordinal_position`,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("查詢 mysql 資料表欄位失敗: %w", err)
	}
	defer rows.Close()

	tables := make([]TableInfo, 0)
	for rows.Next() {
		var schemaName string
		var tableName string
		var column ColumnInfo
		var nullable string
		if err := rows.Scan(
			&schemaName,
			&tableName,
			&column.Name,
			&column.DataType,
			&nullable,
		); err != nil {
			return nil, fmt.Errorf("掃描 mysql 資料表欄位失敗: %w", err)
		}
		column.Nullable = strings.EqualFold(nullable, "YES")
		if keyInfo, ok := uniqueColumns[mysqlColumnKey(schemaName, tableName, column.Name)]; ok {
			column.Unique = true
			column.PrimaryKey = keyInfo.primary
		}

		last := len(tables) - 1
		if last < 0 || tables[last].Schema != schemaName || tables[last].Name != tableName {
			tables = append(tables, TableInfo{Schema: schemaName, Name: tableName})
			last++
		}
		tables[last].Columns = append(tables[last].Columns, column)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷 mysql 資料表欄位失敗: %w", err)
	}
	return tables, nil
}

type mysqlUniqueKeyInfo struct {
	primary bool
}

func mysqlColumnKey(schemaName, tableName, columnName string) string {
	return strings.ToLower(schemaName) + "\x00" + strings.ToLower(tableName) + "\x00" + strings.ToLower(columnName)
}

// inspectMySQLSingleColumnUniqueKeys 只回報「單獨構成唯一索引」的欄位。
// information_schema.columns.column_key 會把複合主鍵 / 複合唯一鍵的每個欄位都
// 標成 PRI / UNI，直接採用會讓 upsert 誤判成可用（實際上 ON DUPLICATE KEY
// 永遠撞不到那個複合鍵，每次輪詢都新增一列）。
func inspectMySQLSingleColumnUniqueKeys(
	ctx context.Context,
	db *sql.DB,
	placeholders string,
	args []any,
) (map[string]mysqlUniqueKeyInfo, error) {
	queryArgs := make([]any, 0, len(args)*2)
	queryArgs = append(queryArgs, args...)
	queryArgs = append(queryArgs, args...)

	// #nosec G202 -- placeholders 只由 "?," 重複組成，schema 名稱一律以參數繫結傳入
	rows, err := db.QueryContext(
		ctx,
		`SELECT s.table_schema, s.table_name, s.column_name,
		        MAX(CASE WHEN s.index_name = 'PRIMARY' THEN 1 ELSE 0 END) AS is_primary
		 FROM information_schema.statistics s
		 JOIN (
		     SELECT table_schema, table_name, index_name
		     FROM information_schema.statistics
		     WHERE table_schema IN (`+placeholders+`) AND non_unique = 0
		     GROUP BY table_schema, table_name, index_name
		     HAVING COUNT(*) = 1
		 ) single_col
		   ON single_col.table_schema = s.table_schema
		  AND single_col.table_name = s.table_name
		  AND single_col.index_name = s.index_name
		 WHERE s.table_schema IN (`+placeholders+`) AND s.non_unique = 0
		 GROUP BY s.table_schema, s.table_name, s.column_name`,
		queryArgs...,
	)
	if err != nil {
		return nil, fmt.Errorf("查詢 mysql 唯一索引失敗: %w", err)
	}
	defer rows.Close()

	uniqueColumns := map[string]mysqlUniqueKeyInfo{}
	for rows.Next() {
		var schemaName, tableName, columnName string
		var isPrimary int
		if err := rows.Scan(&schemaName, &tableName, &columnName, &isPrimary); err != nil {
			return nil, fmt.Errorf("掃描 mysql 唯一索引失敗: %w", err)
		}
		uniqueColumns[mysqlColumnKey(schemaName, tableName, columnName)] = mysqlUniqueKeyInfo{primary: isPrimary == 1}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷 mysql 唯一索引失敗: %w", err)
	}
	return uniqueColumns, nil
}

func defaultSchemaForConnector(connector *schema.DatabaseConnector) string {
	if connector.Kind == schema.DatabaseConnectorKindMySQL {
		config, err := parseConnectionConfig(connector.ConnectionConfig)
		if err == nil {
			if databaseName := defaultString(
				stringConfigValue(config, "database"),
				stringConfigValue(config, "dbname"),
			); databaseName != "" {
				return databaseName
			}
		}
	}
	return defaultSchemaForKind(connector.Kind)
}
