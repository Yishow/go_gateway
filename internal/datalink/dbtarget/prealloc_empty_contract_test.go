package dbtarget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPreallocatedReadResultsKeepEmptyNilContract(t *testing.T) {
	t.Run("connector rows", func(t *testing.T) {
		db := ownedRowFixture(t)
		rows, err := db.QueryContext(t.Context(), "SELECT 1 WHERE 0")
		require.NoError(t, err)
		defer func() { require.NoError(t, rows.Close()) }()
		values, err := scanConnectorRows(rows)
		require.NoError(t, err)
		require.Nil(t, values)
	})
	t.Run("mapping rows", func(t *testing.T) {
		db := ownedRowFixture(t)
		rows, err := db.QueryContext(t.Context(), "SELECT 1 WHERE 0")
		require.NoError(t, err)
		defer func() { require.NoError(t, rows.Close()) }()
		values, err := scanTargetMappingRows(rows)
		require.NoError(t, err)
		require.Nil(t, values)
	})
	t.Run("owned rows", func(t *testing.T) {
		db := ownedRowFixture(t)
		values, err := ReadOwnedRows(t.Context(), db, ownedRef("gw-test-absent"), []string{"note"})
		require.NoError(t, err)
		require.Nil(t, values)
	})
	t.Run("postgres tables", func(t *testing.T) {
		db := openPostgresInspectionStub(t, `CREATE TABLE information_schema.tables (table_schema TEXT, table_name TEXT, table_type TEXT)`)
		values, err := inspectPostgresTables(t.Context(), db)
		require.NoError(t, err)
		require.Nil(t, values)
	})
	t.Run("mysql columns", func(t *testing.T) {
		db := openPostgresInspectionStub(t, `CREATE TABLE information_schema.COLUMNS (COLUMN_NAME TEXT, DATA_TYPE TEXT, COLUMN_TYPE TEXT, IS_NULLABLE TEXT, COLUMN_DEFAULT TEXT, EXTRA TEXT, COLLATION_NAME TEXT, TABLE_SCHEMA TEXT, TABLE_NAME TEXT, ORDINAL_POSITION INTEGER)`)
		values, err := readColumns(t.Context(), db, mysqlManagedFamilyMySQL, "empty", "missing")
		require.NoError(t, err)
		require.Nil(t, values)
	})
	t.Run("concurrent helper no requests", func(t *testing.T) {
		results, errs := runConcurrentMySQLInserts(t, nil)
		require.Nil(t, results)
		require.Nil(t, errs)
	})
}
