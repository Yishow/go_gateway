package dbtarget

import (
	"context"
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// information_schema 以 SQLite 模擬：欄位清單與單欄唯一索引查詢分開，
// 複合鍵不得被當成單欄唯一鍵（MySQL 的 column_key 對複合鍵每欄都標 PRI）。
func newMySQLInspectionDB(t *testing.T, name string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "file:"+name+"?mode=memory&cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	db.SetMaxOpenConns(1)

	_, err = db.ExecContext(context.Background(), `
		ATTACH DATABASE ':memory:' AS information_schema;
		CREATE TABLE information_schema.tables (
			table_schema TEXT,
			table_name TEXT,
			table_type TEXT
		);
		CREATE TABLE information_schema.columns (
			table_schema TEXT,
			table_name TEXT,
			column_name TEXT,
			data_type TEXT,
			is_nullable TEXT,
			ordinal_position INTEGER
		);
		CREATE TABLE information_schema.statistics (
			table_schema TEXT,
			table_name TEXT,
			index_name TEXT,
			column_name TEXT,
			non_unique INTEGER,
			seq_in_index INTEGER
		);
	`)
	require.NoError(t, err)
	return db
}

func TestInspectTablesByKind_MySQLReturnsCurrentDatabaseColumns(t *testing.T) {
	t.Parallel()

	db := newMySQLInspectionDB(t, "mysql-inspection")
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO information_schema.tables VALUES
			('gateway_metrics', 'sensor_readings', 'BASE TABLE'),
			('other_database', 'ignored_table', 'BASE TABLE');
		INSERT INTO information_schema.columns VALUES
			('gateway_metrics', 'sensor_readings', 'ts', 'datetime', 'NO', 1),
			('gateway_metrics', 'sensor_readings', 'temp_out_c', 'double', 'YES', 2),
			('other_database', 'ignored_table', 'value', 'double', 'YES', 1);
		INSERT INTO information_schema.statistics VALUES
			('gateway_metrics', 'sensor_readings', 'PRIMARY', 'ts', 0, 1);
	`)
	require.NoError(t, err)

	tables, err := inspectTablesByKind(
		ctx,
		db,
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{"database": "gateway_metrics"},
		nil,
	)

	require.NoError(t, err)
	require.Equal(t, []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime", Nullable: false, PrimaryKey: true, Unique: true},
			{Name: "temp_out_c", DataType: "double", Nullable: true, PrimaryKey: false, Unique: false},
		},
	}}, tables)
}

// 複合主鍵的欄位不得被當成 upsert 可用的唯一鍵，否則 ON DUPLICATE KEY 永遠
// 撞不到那個鍵，每次輪詢都會新增一列。
func TestInspectTablesByKind_MySQLRejectsCompositeKeysAsUnique(t *testing.T) {
	t.Parallel()

	db := newMySQLInspectionDB(t, "mysql-composite")
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO information_schema.tables VALUES
			('gateway_metrics', 'composite', 'BASE TABLE');
		INSERT INTO information_schema.columns VALUES
			('gateway_metrics', 'composite', 'device_id', 'varchar', 'NO', 1),
			('gateway_metrics', 'composite', 'ts', 'datetime', 'NO', 2),
			('gateway_metrics', 'composite', 'value', 'double', 'YES', 3);
		INSERT INTO information_schema.statistics VALUES
			('gateway_metrics', 'composite', 'PRIMARY', 'device_id', 0, 1),
			('gateway_metrics', 'composite', 'PRIMARY', 'ts', 0, 2);
	`)
	require.NoError(t, err)

	tables, err := inspectTablesByKind(
		ctx,
		db,
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{"database": "gateway_metrics"},
		nil,
	)

	require.NoError(t, err)
	require.Len(t, tables, 1)
	for _, column := range tables[0].Columns {
		require.False(t, column.Unique, "複合鍵欄位 %s 不得被視為唯一", column.Name)
		require.False(t, column.PrimaryKey, "複合鍵欄位 %s 不得被視為單欄主鍵", column.Name)
	}
}

// 映射引用到連接器自身資料庫以外的 schema 時，檢查範圍必須涵蓋它。
func TestInspectTablesByKind_MySQLIncludesReferencedSchemas(t *testing.T) {
	t.Parallel()

	db := newMySQLInspectionDB(t, "mysql-extra-schema")
	ctx := context.Background()
	_, err := db.ExecContext(ctx, `
		INSERT INTO information_schema.tables VALUES
			('gateway_metrics', 'sensor_readings', 'BASE TABLE'),
			('reporting', 'daily_rollup', 'BASE TABLE'),
			('unrelated', 'noise', 'BASE TABLE');
		INSERT INTO information_schema.columns VALUES
			('gateway_metrics', 'sensor_readings', 'ts', 'datetime', 'NO', 1),
			('reporting', 'daily_rollup', 'ts', 'datetime', 'NO', 1),
			('unrelated', 'noise', 'ts', 'datetime', 'NO', 1);
		INSERT INTO information_schema.statistics VALUES
			('reporting', 'daily_rollup', 'uk_ts', 'ts', 0, 1);
	`)
	require.NoError(t, err)

	tables, err := inspectTablesByKind(
		ctx,
		db,
		schema.DatabaseConnectorKindMySQL,
		ConnectionConfig{"database": "gateway_metrics"},
		[]string{"reporting", "reporting", " "},
	)

	require.NoError(t, err)
	schemas := make([]string, 0, len(tables))
	for _, table := range tables {
		schemas = append(schemas, table.Schema)
	}
	require.ElementsMatch(t, []string{"gateway_metrics", "reporting"}, schemas)
	for _, table := range tables {
		if table.Schema == "reporting" {
			require.True(t, table.Columns[0].Unique)
			require.False(t, table.Columns[0].PrimaryKey)
		}
	}
}

func TestDefaultSchemaForConnector_MySQLUsesDatabaseName(t *testing.T) {
	t.Parallel()

	connector := &schema.DatabaseConnector{
		Kind:             schema.DatabaseConnectorKindMySQL,
		ConnectionConfig: `{"database":"gateway_metrics"}`,
	}

	require.Equal(t, "gateway_metrics", defaultSchemaForConnector(connector))
}

func TestQuoteIdentifier_MySQLUsesBackticks(t *testing.T) {
	t.Parallel()

	require.Equal(t, "`gateway``metrics`", quoteIdentifier(
		schema.DatabaseConnectorKindMySQL,
		"gateway`metrics",
	))
}
