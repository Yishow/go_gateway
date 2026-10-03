package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	moderncsqlite "modernc.org/sqlite"
)

func openPostgresInspectionStub(t *testing.T, setup ...string) *sql.DB {
	t.Helper()
	dir := t.TempDir()
	driverName := fmt.Sprintf("sqlite-postgres-inspection-%d", time.Now().UnixNano())
	pgCatalogPath := filepath.Join(dir, "pg_catalog.db")
	informationSchemaPath := filepath.Join(dir, "information_schema.db")
	driverInstance := &moderncsqlite.Driver{}
	driverInstance.RegisterConnectionHook(func(conn moderncsqlite.ExecQuerierContext, _ string) error {
		for _, attached := range []struct {
			name string
			path string
		}{
			{name: "pg_catalog", path: pgCatalogPath},
			{name: "information_schema", path: informationSchemaPath},
		} {
			statement := "ATTACH DATABASE '" + strings.ReplaceAll(attached.path, "'", "''") + "' AS " + attached.name
			if _, err := conn.ExecContext(context.Background(), statement, nil); err != nil {
				return err
			}
		}
		return nil
	})
	sql.Register(driverName, driverInstance)
	db, err := sql.Open(driverName, filepath.Join(dir, "main.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	for _, statement := range setup {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	return db
}

func TestInspectPostgresTablesPreservesNumericPrecisionAndScale(t *testing.T) {
	db := openPostgresInspectionStub(t,
		`CREATE TABLE information_schema.tables (table_schema TEXT, table_name TEXT, table_type TEXT)`,
		`CREATE TABLE information_schema.columns (table_schema TEXT, table_name TEXT, column_name TEXT, data_type TEXT, numeric_precision INTEGER, numeric_scale INTEGER, is_nullable TEXT, ordinal_position INTEGER)`,
		`CREATE TABLE information_schema.table_constraints (constraint_name TEXT, table_schema TEXT, table_name TEXT, constraint_type TEXT)`,
		`CREATE TABLE information_schema.key_column_usage (constraint_name TEXT, table_schema TEXT, table_name TEXT, column_name TEXT)`,
		`INSERT INTO information_schema.tables VALUES ('public', 'sensor_values', 'BASE TABLE')`,
		`INSERT INTO information_schema.columns VALUES ('public', 'sensor_values', 'exact_value', 'numeric', 20, 0, 'NO', 1)`,
		`INSERT INTO information_schema.columns VALUES ('public', 'sensor_values', 'fractional_value', 'numeric', 20, 2, 'YES', 2)`,
		`INSERT INTO information_schema.columns VALUES ('public', 'sensor_values', 'small_value', 'numeric', 18, 0, 'YES', 3)`,
		`INSERT INTO information_schema.columns VALUES ('public', 'sensor_values', 'unbounded_value', 'numeric', NULL, NULL, 'YES', 4)`,
		`INSERT INTO information_schema.columns VALUES ('public', 'sensor_values', 'note', 'text', NULL, NULL, 'YES', 5)`,
	)

	tables, err := inspectPostgresTables(t.Context(), db)

	require.NoError(t, err)
	require.Len(t, tables, 1)
	require.Equal(t, []string{
		"numeric(20,0)",
		"numeric(20,2)",
		"numeric(18,0)",
		"numeric",
		"text",
	}, []string{
		tables[0].Columns[0].DataType,
		tables[0].Columns[1].DataType,
		tables[0].Columns[2].DataType,
		tables[0].Columns[3].DataType,
		tables[0].Columns[4].DataType,
	})
}

func TestPostgresColumnDataTypeRejectsIncompleteNumericMetadata(t *testing.T) {
	validPrecision := sql.NullInt64{Int64: 20, Valid: true}
	validScale := sql.NullInt64{Int64: 0, Valid: true}
	for _, tt := range []struct {
		name      string
		precision sql.NullInt64
		scale     sql.NullInt64
	}{
		{name: "missing scale", precision: validPrecision},
		{name: "missing precision", scale: validScale},
		{name: "zero precision", precision: sql.NullInt64{Int64: 0, Valid: true}, scale: validScale},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := postgresColumnDataType("numeric", tt.precision, tt.scale)
			require.Error(t, err)
		})
	}
}
