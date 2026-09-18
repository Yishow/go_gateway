package dbtarget

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

var postgresIdentifierPattern = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

// pgIdentifier quotes a generated test identifier; the pattern keeps the test
// from building an identifier out of anything unexpected.
func pgIdentifier(t *testing.T, name string) string {
	t.Helper()
	require.True(t, postgresIdentifierPattern.MatchString(name), "unsafe test identifier %q", name)
	return `"` + name + `"`
}

func postgresKeywordDSN(t *testing.T) (dsn string, fields map[string]string) {
	t.Helper()
	dsn = strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL schema execution test: POSTGRES_DSN not set")
	}
	fields = map[string]string{}
	for _, part := range strings.Fields(dsn) {
		key, value, found := strings.Cut(part, "=")
		if found {
			fields[key] = value
		}
	}
	return dsn, fields
}

// postgresAdminDB connects as the DSN's own role for test setup and checks.
func postgresAdminDB(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	require.NoError(t, db.PingContext(t.Context()))
	return db
}

// postgresScopedSchema creates an isolated schema for one test and drops it,
// with everything in it, afterwards.
func postgresScopedSchema(t *testing.T, admin *sql.DB) string {
	t.Helper()
	name := fmt.Sprintf("gw_apply_%d", time.Now().UnixNano())
	_, err := admin.ExecContext(t.Context(), `CREATE SCHEMA `+pgIdentifier(t, name))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP SCHEMA `+pgIdentifier(t, name)+` CASCADE`)
		// A batch that escaped the scoped schema must not be left in public.
		_, _ = admin.ExecContext(context.Background(), `DROP TABLE IF EXISTS public.`+pgIdentifier(t, managedProbeTable))
	})
	return name
}

func postgresConnectorConfig(fields map[string]string, schemaName, user, password string) ConnectionConfig {
	config := ConnectionConfig{
		"host": fields["host"], "port": fields["port"], "user": user, "password": password,
		"database": fields["dbname"], "schema": schemaName, "sslmode": "disable",
	}
	if config["port"] == "" {
		config["port"] = "5432"
	}
	return config
}

// managedProbeTable is the single table these execution tests create.
const managedProbeTable = "gw_record_a"

func postgresTableExists(t *testing.T, admin *sql.DB, schemaName string) bool {
	t.Helper()
	var count int
	require.NoError(t, admin.QueryRowContext(t.Context(), `
		SELECT COUNT(*) FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1 AND c.relname = $2 AND c.relkind IN ('r', 'p')
	`, schemaName, managedProbeTable).Scan(&count))
	return count == 1
}

// postgresExecutionTarget prepares an isolated schema and a connector bound to it.
func postgresExecutionTarget(t *testing.T) (*ConnectorService, *schema.DatabaseConnector, *sql.DB, string) {
	t.Helper()
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)
	service := NewConnectorService(NewSQLConnectorRepository(openMigratedTestDB(t)))
	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "managed postgres target", Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: postgresConnectorConfig(fields, schemaName, fields["user"], fields["password"]),
	})
	require.NoError(t, err)
	return service, connector, admin, schemaName
}

func TestExecuteSchemaStatements_PostgresCommitsIntoTheScopedSchema(t *testing.T) {
	service, connector, admin, schemaName := postgresExecutionTarget(t)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, schemaName, []string{
		`CREATE TABLE IF NOT EXISTS gw_record_a (id TEXT PRIMARY KEY)`,
		`CREATE INDEX IF NOT EXISTS idx_gw_record_a_id ON gw_record_a(id)`,
	})

	require.NoError(t, err)
	require.Equal(t, &SchemaStatementExecution{Committed: 2}, result)
	require.True(t, postgresTableExists(t, admin, schemaName), "the table must be created in the scoped schema")
	require.False(t, postgresTableExists(t, admin, "public"), "the confirmed scope must not leak into public")
}

func TestExecuteSchemaStatements_PostgresRollsBackAFailedBatch(t *testing.T) {
	service, connector, admin, schemaName := postgresExecutionTarget(t)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, schemaName, []string{
		`CREATE TABLE gw_record_a (id TEXT PRIMARY KEY)`,
		`CREATE TABLE gw_record_a (id TEXT PRIMARY KEY)`,
	})

	require.Error(t, err)
	require.Equal(t, &SchemaStatementExecution{RolledBack: true}, result)
	require.False(t, postgresTableExists(t, admin, schemaName), "a failed batch must leave no table behind")
}

func TestExecuteSchemaStatements_PostgresReportsPermissionDenied(t *testing.T) {
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	schemaName := postgresScopedSchema(t, admin)

	secret := make([]byte, 16)
	_, err := rand.Read(secret)
	require.NoError(t, err)
	roleName := fmt.Sprintf("gw_apply_role_%d", time.Now().UnixNano())
	password := hex.EncodeToString(secret)
	_, err = admin.ExecContext(t.Context(), fmt.Sprintf(`CREATE ROLE %s LOGIN PASSWORD '%s'`, pgIdentifier(t, roleName), password))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP OWNED BY `+pgIdentifier(t, roleName))
		_, _ = admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS `+pgIdentifier(t, roleName))
	})
	// The role may look into the schema but may not create anything in it.
	_, err = admin.ExecContext(t.Context(), fmt.Sprintf(`GRANT USAGE ON SCHEMA %s TO %s`, pgIdentifier(t, schemaName), pgIdentifier(t, roleName)))
	require.NoError(t, err)

	service := NewConnectorService(NewSQLConnectorRepository(openMigratedTestDB(t)))
	connector, err := service.Create(t.Context(), CreateConnectorRequest{
		Name: "restricted postgres target", Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: postgresConnectorConfig(fields, schemaName, roleName, password),
	})
	require.NoError(t, err)

	result, err := service.ExecuteSchemaStatements(t.Context(), connector.ID, schemaName, []string{
		`CREATE TABLE IF NOT EXISTS gw_record_a (id TEXT PRIMARY KEY)`,
	})

	require.True(t, errors.Is(err, ErrSchemaExecutionPermissionDenied), "a refused batch must be reported as permission denied, got %v", err)
	require.Equal(t, &SchemaStatementExecution{RolledBack: true}, result)
	require.False(t, postgresTableExists(t, admin, schemaName))
}
