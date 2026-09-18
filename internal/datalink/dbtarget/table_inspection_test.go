package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newInspectionService(t *testing.T, kind schema.DatabaseConnectorKind, config string) (service *ConnectorService, connectorID string) {
	t.Helper()
	repo := NewSQLConnectorRepository(openMigratedTestDB(t))
	connector := &schema.DatabaseConnector{
		ID: "inspect-" + string(kind), Name: "inspection", Kind: kind, ConnectionConfig: config,
		Status: schema.DatabaseConnectorStatusUnreachable, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, repo.Create(t.Context(), connector))
	return NewConnectorService(repo), connector.ID
}

func newInspectionSQLiteTarget(t *testing.T, statements ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer db.Close()
	for _, statement := range statements {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
	}
	return path
}

// stubInspectionDatabase serves every external open from a SQLite database
// shaped by setup, so server catalog behavior is reproducible without a server.
func stubInspectionDatabase(t *testing.T, openErr error, setup ...string) {
	t.Helper()
	original := openExternalDBManagerFunc
	path := filepath.Join(t.TempDir(), "external.db")
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		if openErr != nil {
			return nil, openErr
		}
		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{Type: datalinkbase.DBTypeSQLite, DSN: path, MaxOpenConns: 1, MaxIdleConns: 1})
		if err := manager.Connect(); err != nil {
			return nil, err
		}
		for _, statement := range setup {
			if _, err := manager.DB().ExecContext(context.Background(), statement); err != nil {
				_ = manager.Close()
				return nil, err
			}
		}
		return manager, nil
	}
	t.Cleanup(func() { openExternalDBManagerFunc = original })
}

// postgresCatalogSetup models pg_catalog (visible to every account) and a
// privilege-filtered information_schema that starts empty.
func postgresCatalogSetup(dir string, rows ...string) []string {
	return append([]string{
		`ATTACH DATABASE '` + filepath.Join(dir, "pg_catalog.db") + `' AS pg_catalog`,
		`CREATE TABLE IF NOT EXISTS pg_catalog.pg_namespace (oid INTEGER, nspname TEXT)`,
		`CREATE TABLE IF NOT EXISTS pg_catalog.pg_class (relname TEXT, relnamespace INTEGER, relkind TEXT)`,
		`ATTACH DATABASE '` + filepath.Join(dir, "information_schema.db") + `' AS information_schema`,
		`CREATE TABLE IF NOT EXISTS information_schema.tables (table_schema TEXT, table_name TEXT, table_type TEXT)`,
		`INSERT INTO pg_catalog.pg_namespace (oid, nspname) VALUES (2200, 'public')`,
	}, rows...)
}

func mysqlCatalogSetup(dir string, rows ...string) []string {
	return append([]string{
		`ATTACH DATABASE '` + filepath.Join(dir, "information_schema.db") + `' AS information_schema`,
		`CREATE TABLE IF NOT EXISTS information_schema.tables (table_schema TEXT, table_name TEXT, table_type TEXT)`,
		`CREATE TABLE IF NOT EXISTS information_schema.columns (table_schema TEXT, table_name TEXT, column_name TEXT, data_type TEXT, is_nullable TEXT, ordinal_position INTEGER)`,
		`CREATE TABLE IF NOT EXISTS information_schema.statistics (table_schema TEXT, table_name TEXT, index_name TEXT, column_name TEXT, non_unique INTEGER, seq_in_index INTEGER)`,
	}, rows...)
}

func requireInspection(t *testing.T, result *TableInspection, err error, status TableInspectionStatus, reason string) {
	t.Helper()
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, status, result.Status)
	require.Equal(t, reason, result.Reason)
	if status != TableInspectionExists {
		require.Empty(t, result.Columns, "only an existing table may carry inspected columns")
	}
}

func TestInspectTable_SQLiteReportsTheSelectedTablesActualColumns(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE sensor_values (recorded_at TEXT NOT NULL PRIMARY KEY, reactor_temp REAL NOT NULL, note TEXT)`)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindSQLite, `{"dsn":"`+target+`"}`)

	result, err := service.InspectTable(t.Context(), id, "main", "sensor_values")

	requireInspection(t, result, err, TableInspectionExists, "")
	require.Equal(t, "main", result.Schema)
	require.Equal(t, "sensor_values", result.Table)
	require.Equal(t, []ColumnInfo{
		{Name: "recorded_at", DataType: "TEXT", Nullable: false, PrimaryKey: true, Unique: true},
		{Name: "reactor_temp", DataType: "REAL", Nullable: false},
		{Name: "note", DataType: "TEXT", Nullable: true},
	}, result.Columns)
}

func TestInspectTable_SQLiteAbsentTableIsMissing(t *testing.T) {
	target := newInspectionSQLiteTarget(t, `CREATE TABLE other_values (id INTEGER)`)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindSQLite, `{"dsn":"`+target+`"}`)

	result, err := service.InspectTable(t.Context(), id, "main", "sensor_values")

	requireInspection(t, result, err, TableInspectionMissing, "")
}

func TestInspectTable_PostgresTableHiddenByPrivilegesIsForbidden(t *testing.T) {
	stubInspectionDatabase(t, nil, postgresCatalogSetup(t.TempDir(),
		`INSERT INTO pg_catalog.pg_class (relname, relnamespace, relkind) VALUES ('sensor_values', 2200, 'r')`)...)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindPostgres, `{"host":"db.internal","user":"reader","database":"metrics"}`)

	result, err := service.InspectTable(t.Context(), id, "public", "sensor_values")

	requireInspection(t, result, err, TableInspectionForbidden, "permission_denied")
}

func TestInspectTable_PostgresAbsentFromCatalogIsMissing(t *testing.T) {
	stubInspectionDatabase(t, nil, postgresCatalogSetup(t.TempDir())...)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindPostgres, `{"host":"db.internal","user":"reader","database":"metrics"}`)

	result, err := service.InspectTable(t.Context(), id, "public", "sensor_values")

	requireInspection(t, result, err, TableInspectionMissing, "")
}

func TestInspectTable_UnreadableCatalogIsFailedNotMissing(t *testing.T) {
	stubInspectionDatabase(t, nil)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindPostgres, `{"host":"db.internal","user":"reader","database":"metrics"}`)

	result, err := service.InspectTable(t.Context(), id, "public", "sensor_values")

	requireInspection(t, result, err, TableInspectionFailed, "inspection_failed")
}

func TestInspectTable_MySQLUnlistedTableIsUnconfirmed(t *testing.T) {
	stubInspectionDatabase(t, nil, mysqlCatalogSetup(t.TempDir())...)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindMySQL, `{"host":"db.internal","user":"reader","database":"metrics"}`)

	result, err := service.InspectTable(t.Context(), id, "metrics", "sensor_values")

	requireInspection(t, result, err, TableInspectionFailed, "existence_unconfirmed")
}

func TestInspectTable_MySQLListedTableReportsColumns(t *testing.T) {
	stubInspectionDatabase(t, nil, mysqlCatalogSetup(t.TempDir(),
		`INSERT INTO information_schema.tables VALUES ('metrics', 'sensor_values', 'BASE TABLE')`,
		`INSERT INTO information_schema.columns VALUES ('metrics', 'sensor_values', 'recorded_at', 'datetime', 'NO', 1)`,
		`INSERT INTO information_schema.columns VALUES ('metrics', 'sensor_values', 'reactor_temp', 'double', 'YES', 2)`)...)
	service, id := newInspectionService(t, schema.DatabaseConnectorKindMySQL, `{"host":"db.internal","user":"reader","database":"metrics"}`)

	result, err := service.InspectTable(t.Context(), id, "metrics", "sensor_values")

	requireInspection(t, result, err, TableInspectionExists, "")
	require.Equal(t, []string{"recorded_at", "reactor_temp"}, []string{result.Columns[0].Name, result.Columns[1].Name})
}

func TestInspectTable_ConnectionFailuresAreFailed(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		reason string
	}{
		{name: "authentication", err: errors.New(`pq: password authentication failed for user "reader"`), reason: "authentication_failed"},
		{name: "unreachable", err: errors.New("dial tcp 127.0.0.1:5432: connect: connection refused"), reason: "connection_failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubInspectionDatabase(t, tt.err)
			service, id := newInspectionService(t, schema.DatabaseConnectorKindPostgres, `{"host":"db.internal","user":"reader","database":"metrics"}`)

			result, err := service.InspectTable(t.Context(), id, "public", "sensor_values")

			requireInspection(t, result, err, TableInspectionFailed, tt.reason)
		})
	}
}

func TestInspectTable_UnknownConnectorIsAnError(t *testing.T) {
	service, _ := newInspectionService(t, schema.DatabaseConnectorKindSQLite, `{"dsn":"unused.db"}`)

	_, err := service.InspectTable(t.Context(), "missing-connector", "main", "sensor_values")

	require.ErrorIs(t, err, ErrConnectorNotFound)
}

func TestClassifyInspectionError_PermissionDeniedIsForbidden(t *testing.T) {
	for _, err := range []error{
		&pq.Error{Code: "42501", Message: "permission denied for table sensor_values"},
		&pgconn.PgError{Code: "42501", Message: "permission denied for table sensor_values"},
		&mysqldriver.MySQLError{Number: 1142, Message: "SELECT command denied"},
	} {
		status, reason := classifyInspectionError(err)
		require.Equal(t, TableInspectionForbidden, status, "%T", err)
		require.Equal(t, "permission_denied", reason)
	}
	status, reason := classifyInspectionError(errors.New("i/o timeout"))
	require.Equal(t, TableInspectionFailed, status)
	require.Equal(t, "inspection_failed", reason)
}
