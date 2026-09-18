package dbtarget

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func sqliteExecutionTarget(t *testing.T) (*ConnectorService, *schema.DatabaseConnector, *sql.DB) {
	t.Helper()
	service := NewConnectorService(NewSQLConnectorRepository(openMigratedTestDB(t)))
	targetPath := filepath.Join(t.TempDir(), "target.db")
	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "managed target", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: ConnectionConfig{"dsn": targetPath},
	})
	require.NoError(t, err)
	target, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = target.Close() })
	return service, connector, target
}

func sqliteTableExists(t *testing.T, db *sql.DB, table string) bool {
	t.Helper()
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, table).Scan(&count))
	return count == 1
}

func TestExecuteSchemaStatements_SQLiteCommitsTheWholeBatch(t *testing.T) {
	service, connector, target := sqliteExecutionTarget(t)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, "main", []string{
		`CREATE TABLE gw_record_a (id TEXT PRIMARY KEY)`,
		`CREATE INDEX idx_gw_record_a_id ON gw_record_a(id)`,
	})

	require.NoError(t, err)
	require.Equal(t, &SchemaStatementExecution{Committed: 2}, result)
	require.True(t, sqliteTableExists(t, target, "gw_record_a"))
}

func TestExecuteSchemaStatements_SQLiteRollsBackAFailedBatch(t *testing.T) {
	service, connector, target := sqliteExecutionTarget(t)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, "main", []string{
		`CREATE TABLE gw_record_a (id TEXT PRIMARY KEY)`,
		`CREATE TABLE gw_record_a (id TEXT PRIMARY KEY)`,
	})

	require.Error(t, err)
	require.Equal(t, &SchemaStatementExecution{RolledBack: true}, result)
	require.False(t, sqliteTableExists(t, target, "gw_record_a"), "a failed batch must leave no table behind")
}

func TestExecuteSchemaStatements_EmptyBatchDoesNothing(t *testing.T) {
	service, connector, target := sqliteExecutionTarget(t)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, "main", nil)

	require.NoError(t, err)
	require.Equal(t, &SchemaStatementExecution{}, result)
	require.False(t, sqliteTableExists(t, target, "gw_record_a"))
}

func TestExecuteSchemaStatements_UnverifiedKindsAndUnknownConnectorsDoNotOpenTargets(t *testing.T) {
	probeCalls := installIdentityProbeStub(t)
	service := NewConnectorService(NewSQLConnectorRepository(openMigratedTestDB(t)))
	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "mysql target", Kind: schema.DatabaseConnectorKindMySQL,
		ConnectionConfig: ConnectionConfig{"host": "db.internal", "port": 3306, "user": "writer", "database": "metrics"},
	})
	require.NoError(t, err)
	opened := *probeCalls

	_, err = service.ExecuteSchemaStatements(t.Context(), connector.ID, "main", []string{`CREATE TABLE gw_record_a (id TEXT)`})
	require.True(t, errors.Is(err, ErrSchemaExecutionUnsupported), "mysql execution must stay unavailable, got %v", err)
	require.Equal(t, opened, *probeCalls, "an unverified kind must not open the target")

	_, err = service.ExecuteSchemaStatements(t.Context(), "connector-missing", "main", []string{`CREATE TABLE gw_record_a (id TEXT)`})
	require.Error(t, err)
	require.Equal(t, opened, *probeCalls, "an unknown connector must not open a target")
}
