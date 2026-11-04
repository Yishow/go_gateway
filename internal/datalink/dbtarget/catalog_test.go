package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestConnectorServiceReadCatalogSQLiteFiltersSystemTablesAndReturnsEmptySchema(t *testing.T) {
	t.Parallel()

	targetPath := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `
		CREATE TABLE gw_records (id INTEGER PRIMARY KEY AUTOINCREMENT, value TEXT);
		CREATE TABLE user_records (id INTEGER PRIMARY KEY, value TEXT);
		CREATE TABLE sqliteX_records (id INTEGER PRIMARY KEY, value TEXT);
	`)
	require.NoError(t, err)
	require.NoError(t, db.Close())

	mainDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(mainDB)
	connector := &schema.DatabaseConnector{
		ID:               "catalog-sqlite",
		Name:             "Catalog SQLite",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: mustJSONConnectionConfig(t, ConnectionConfig{"dsn": targetPath}),
		IdentityRevision: "connector-revision-1",
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, repo.Create(t.Context(), connector))

	service := NewConnectorService(repo)
	catalog, err := service.ReadCatalog(t.Context(), connector.ID, connector.IdentityRevision)

	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorKindSQLite, catalog.Kind)
	require.Equal(t, "main", catalog.DefaultSchema)
	require.Equal(t, []CatalogSchema{{Name: "main", Tables: []string{"gw_records", "sqliteX_records", "user_records"}}}, catalog.Schemas)
	encoded, err := json.Marshal(catalog)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "sqlite_sequence")
}

func TestConnectorServiceReadCatalogSQLiteMissingFileDoesNotCreateDestination(t *testing.T) {
	t.Parallel()

	targetPath := filepath.Join(t.TempDir(), "missing", "target.db")
	mainDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(mainDB)
	connector := &schema.DatabaseConnector{
		ID:               "catalog-missing-sqlite",
		Name:             "Missing Catalog SQLite",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: mustJSONConnectionConfig(t, ConnectionConfig{"path": targetPath}),
		IdentityRevision: "connector-revision-1",
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, repo.Create(t.Context(), connector))

	service := NewConnectorService(repo)
	_, err := service.ReadCatalog(t.Context(), connector.ID, connector.IdentityRevision)

	require.Error(t, err)
	require.False(t, errors.Is(err, context.Canceled))
	_, statErr := os.Stat(targetPath)
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestCatalogSchemasPreserveUserSchemaWithoutTablesAndFilterBuiltins(t *testing.T) {
	t.Parallel()

	got := catalogSchemasFromRows(
		[]catalogSchemaRow{{Name: "pg_catalog"}, {Name: "information_schema"}, {Name: "pg_temp_3"}, {Name: "PG_USER"}, {Name: "public"}, {Name: "gw_archive"}},
		[]catalogTableRow{{Schema: "public", Name: "sensor_values"}, {Schema: "gw_archive", Name: "gw_records"}, {Schema: "pg_catalog", Name: "pg_type"}},
		catalogFilterPostgres,
	)

	require.Equal(t, []CatalogSchema{
		{Name: "PG_USER", Tables: []string{}},
		{Name: "gw_archive", Tables: []string{"gw_records"}},
		{Name: "public", Tables: []string{"sensor_values"}},
	}, got)
}

func TestCatalogPreservesCaseSensitiveUserSchemas(t *testing.T) {
	t.Parallel()
	require.False(t, catalogSchemaAllowed("mysql", catalogFilterMySQL))
	require.True(t, catalogSchemaAllowed("MySQL", catalogFilterMySQL))
	require.False(t, catalogSchemaAllowed("sys", catalogFilterSQLServer))
	require.False(t, catalogSchemaAllowed("INFORMATION_SCHEMA", catalogFilterSQLServer))
	require.True(t, catalogSchemaAllowed("Sys", catalogFilterSQLServer))
}

func TestDefaultSchemaForCatalogUsesSavedDialectSettings(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name      string
		connector *schema.DatabaseConnector
		want      string
	}{
		{name: "postgres", connector: &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindPostgres}, want: "public"},
		{name: "postgres-saved-schema", connector: &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindPostgres, ConnectionConfig: mustJSONConnectionConfig(t, ConnectionConfig{"schema": "analytics"})}, want: "analytics"},
		{name: "sqlite", connector: &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindSQLite}, want: "main"},
		{name: "mysql", connector: &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindMySQL, ConnectionConfig: mustJSONConnectionConfig(t, ConnectionConfig{"database": "gateway_metrics"})}, want: "gateway_metrics"},
		{name: "sqlserver", connector: &schema.DatabaseConnector{Kind: schema.DatabaseConnectorKindSQLServer, ConnectionConfig: mustJSONConnectionConfig(t, ConnectionConfig{"schema": "dbo"})}, want: "dbo"},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := ConnectionConfig{}
			if test.connector.ConnectionConfig != "" {
				var err error
				config, err = parseConnectionConfig(test.connector.ConnectionConfig)
				require.NoError(t, err)
			}
			require.Equal(t, test.want, defaultSchemaForCatalogConnector(test.connector.Kind, config))
		})
	}
}

func TestReadCatalogPostgresIncludesEmptyAccessibleSchemasAndKeepsCaseSensitiveUserNames(t *testing.T) {
	t.Parallel()

	db := openPostgresInspectionStub(t,
		`CREATE TABLE information_schema.schemata (schema_name TEXT)`,
		`CREATE TABLE information_schema.tables (table_schema TEXT, table_name TEXT, table_type TEXT)`,
		`INSERT INTO information_schema.schemata VALUES ('pg_catalog'), ('information_schema'), ('pg_temp_3'), ('analytics')`,
		`INSERT INTO information_schema.tables VALUES ('analytics', 'gw_values', 'BASE TABLE'), ('pg_catalog', 'pg_type', 'BASE TABLE')`,
	)

	catalog, err := readCatalogByKind(t.Context(), db, schema.DatabaseConnectorKindPostgres, ConnectionConfig{"schema": "analytics"})

	require.NoError(t, err)
	require.Equal(t, "analytics", catalog.DefaultSchema)
	require.Equal(t, []CatalogSchema{
		{Name: "analytics", Tables: []string{"gw_values"}},
	}, catalog.Schemas)
}

func TestReadCatalogMySQLFiltersSystemDatabasesAndPreservesEmptyDatabase(t *testing.T) {
	t.Parallel()

	db := newMySQLInspectionDB(t, "catalog-mysql")
	_, err := db.ExecContext(t.Context(), `
		CREATE TABLE information_schema.schemata (schema_name TEXT);
		INSERT INTO information_schema.schemata VALUES ('mysql'), ('information_schema'), ('performance_schema'), ('sys'), (' mysql '), ('gateway_metrics'), ('empty_db');
	`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `
		INSERT INTO information_schema.tables VALUES ('gateway_metrics', 'gw_values', 'BASE TABLE'), ('mysql', 'user', 'BASE TABLE'), (' mysql ', 'kept_table', 'BASE TABLE');
	`)
	require.NoError(t, err)

	catalog, err := readCatalogByKind(t.Context(), db, schema.DatabaseConnectorKindMySQL, ConnectionConfig{"database": "gateway_metrics"})

	require.NoError(t, err)
	require.Equal(t, "gateway_metrics", catalog.DefaultSchema)
	require.Equal(t, []CatalogSchema{
		{Name: " mysql ", Tables: []string{"kept_table"}},
		{Name: "empty_db", Tables: []string{}},
		{Name: "gateway_metrics", Tables: []string{"gw_values"}},
	}, catalog.Schemas)
}

func TestReadCatalogSQLServerFiltersSystemSchemasAndUsesCurrentDefault(t *testing.T) {
	t.Parallel()

	script := &catalogSQLServerScript{}
	db := sql.OpenDB(catalogSQLServerConnector{script: script})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	catalog, err := readCatalogByKind(t.Context(), db, schema.DatabaseConnectorKindSQLServer, ConnectionConfig{})

	require.NoError(t, err)
	require.Equal(t, "analytics", catalog.DefaultSchema)
	require.Equal(t, []CatalogSchema{
		{Name: "analytics", Tables: []string{}},
		{Name: "dbo", Tables: []string{"gateway_values"}},
	}, catalog.Schemas)
	var objectQuery string
	for _, query := range script.queries {
		if strings.Contains(query, "FROM sys.objects") {
			objectQuery = query
			break
		}
	}
	require.NotEmpty(t, objectQuery)
	require.Contains(t, objectQuery, "ISNULL(is_ms_shipped, 0) = 0")
	require.Contains(t, script.queries, "SELECT SCHEMA_NAME()")
}

type catalogSQLServerScript struct {
	queries []string
}

type catalogSQLServerConnector struct {
	script *catalogSQLServerScript
}

func (c catalogSQLServerConnector) Connect(context.Context) (driver.Conn, error) {
	return &catalogSQLServerConn{script: c.script}, nil
}

func (catalogSQLServerConnector) Driver() driver.Driver { return catalogSQLServerDriver{} }

type catalogSQLServerDriver struct{}

func (catalogSQLServerDriver) Open(string) (driver.Conn, error) {
	return nil, fmt.Errorf("catalog SQL Server script requires sql.OpenDB")
}

type catalogSQLServerConn struct {
	script *catalogSQLServerScript
}

func (*catalogSQLServerConn) Prepare(string) (driver.Stmt, error) {
	return nil, fmt.Errorf("catalog SQL Server script does not prepare statements")
}

func (*catalogSQLServerConn) Close() error { return nil }

func (*catalogSQLServerConn) Begin() (driver.Tx, error) {
	return catalogSQLServerTx{}, nil
}

func (c *catalogSQLServerConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	normalized := strings.Join(strings.Fields(query), " ")
	c.script.queries = append(c.script.queries, normalized)
	switch {
	case normalized == "SELECT name FROM sys.schemas ORDER BY name":
		return sqlServerRecordingRows([]string{"name"},
			[]driver.Value{"sys"}, []driver.Value{"INFORMATION_SCHEMA"}, []driver.Value{"dbo"}, []driver.Value{"analytics"}), nil
	case strings.Contains(normalized, "FROM sys.objects"):
		return sqlServerRecordingRows([]string{"schema_name", "table_name"},
			[]driver.Value{"sys", "internal_values"}, []driver.Value{"dbo", "gateway_values"}), nil
	case normalized == "SELECT SCHEMA_NAME()":
		return sqlServerRecordingRows([]string{"schema_name"}, []driver.Value{"analytics"}), nil
	default:
		return nil, fmt.Errorf("unexpected catalog SQL Server query: %s", normalized)
	}
}

func (*catalogSQLServerConn) ExecContext(context.Context, string, []driver.NamedValue) (driver.Result, error) {
	return nil, fmt.Errorf("catalog SQL Server read executed a statement")
}

type catalogSQLServerTx struct{}

func (catalogSQLServerTx) Commit() error   { return nil }
func (catalogSQLServerTx) Rollback() error { return nil }

var _ driver.QueryerContext = (*catalogSQLServerConn)(nil)
var _ driver.ExecerContext = (*catalogSQLServerConn)(nil)

func mustJSONConnectionConfig(t *testing.T, config ConnectionConfig) string {
	t.Helper()
	raw, err := json.Marshal(config)
	require.NoError(t, err)
	return string(raw)
}
