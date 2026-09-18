package api

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

var postgresRouteIdentifier = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func postgresRouteIdentifierQuoted(t *testing.T, name string) string {
	t.Helper()
	if !postgresRouteIdentifier.MatchString(name) {
		t.Fatalf("unsafe test identifier %q", name)
	}
	return `"` + name + `"`
}

// postgresRouteTarget prepares an isolated PostgreSQL schema and the production
// router bound to a connector that writes only into that schema.
func postgresRouteTarget(t *testing.T, role, password string) (recordingRouteFixture, *sql.DB, string) {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL recording route test: POSTGRES_DSN not set")
	}
	fields := map[string]string{}
	for _, part := range strings.Fields(dsn) {
		if key, value, found := strings.Cut(part, "="); found {
			fields[key] = value
		}
	}
	admin, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	t.Cleanup(func() { _ = admin.Close() })
	ctx := context.Background()
	if err := admin.PingContext(ctx); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}
	schemaName := fmt.Sprintf("gw_route_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+postgresRouteIdentifierQuoted(t, schemaName)); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP SCHEMA `+postgresRouteIdentifierQuoted(t, schemaName)+` CASCADE`)
		// Managed tables that escaped the scoped schema must not be left in public.
		for _, table := range []string{"samples", "intervals", "events", "snapshots", "definitions", "receipts"} {
			_, _ = admin.ExecContext(context.Background(), `DROP TABLE IF EXISTS public.`+postgresRouteIdentifierQuoted(t, "gw_record_"+table))
		}
	})

	if role == "" {
		role, password = fields["user"], fields["password"]
	}
	mainDB, err := sql.Open("sqlite", t.TempDir()+"/main.db")
	if err != nil {
		t.Fatalf("open main database: %v", err)
	}
	t.Cleanup(func() { _ = mainDB.Close() })
	if err := datalinkbase.NewMigrator().Migrate(mainDB); err != nil {
		t.Fatalf("migrate main database: %v", err)
	}
	repo := dbtarget.NewSQLConnectorRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(repo)
	port := fields["port"]
	if port == "" {
		port = "5432"
	}
	connector, err := connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "Recording postgres target", Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"host": fields["host"], "port": port, "user": role, "password": password,
			"database": fields["dbname"], "schema": schemaName, "sslmode": "disable",
		},
	})
	if err != nil {
		t.Fatalf("create connector: %v", err)
	}

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	record, err := workspaceSvc.BindDatabaseConnector(ctx, connector.ID)
	if err != nil {
		t.Fatalf("bind connector: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(ctx, "device-a"); err != nil {
		t.Fatalf("attach device: %v", err)
	}
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)
	router := NewRouter(&DatalinkServices{Workspace: workspaceSvc, RecordingPlan: planSvc, DBTarget: connectorSvc})
	return recordingRouteFixture{
		router: router, planRepo: planRepo, planSvc: planSvc, connector: connector, repo: repo,
		workspaceID: record.ID, workspaceRevision: record.DatabaseSetupRevision,
	}, admin, schemaName
}

func postgresRouteTables(t *testing.T, admin *sql.DB, schemaName string) []string {
	t.Helper()
	rows, err := admin.QueryContext(context.Background(), `
		SELECT c.relname FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1 AND c.relkind IN ('r', 'p') ORDER BY c.relname
	`, schemaName)
	if err != nil {
		t.Fatalf("list schema tables: %v", err)
	}
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan table: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate tables: %v", err)
	}
	return names
}

func TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema(t *testing.T) {
	fixture, admin, schemaName := postgresRouteTarget(t, "", "")
	seedRecordingRoutePlan(t, fixture, "plan-postgres", "")
	token := previewForRoute(t, fixture, "plan-postgres")

	w := fixture.apply(fixture.applyBody(t, token, nil))

	created := decodeRecordingOperation(t, w).Data
	if w.Code != http.StatusOK || created.Status != recordingplan.SchemaOperationSucceeded ||
		created.ExecutedStatements != len(token.Statements) || len(created.VerifiedDigest) != 64 {
		t.Fatalf("postgres apply must create and verify the schema: %d %s", w.Code, w.Body.String())
	}
	if tables := postgresRouteTables(t, admin, schemaName); len(tables) != 6 {
		t.Fatalf("all managed tables must exist in the confirmed schema, got %v", tables)
	}

	// A second preview of the now compatible schema is an explained no-op that
	// still verifies the existing tables when confirmed.
	compatible := previewForRoute(t, fixture, "plan-postgres")
	if len(compatible.Statements) != 0 || compatible.NoChangeReason != recordingplan.NoChangeSchemaCompatible {
		t.Fatalf("a compatible postgres schema must preview as a no-op: %+v", compatible)
	}
	noop := fixture.apply(fixture.applyBody(t, compatible, nil))
	verified := decodeRecordingOperation(t, noop).Data
	if noop.Code != http.StatusOK || verified.Status != recordingplan.SchemaOperationSucceeded ||
		verified.ExecutedStatements != 0 || len(verified.VerifiedDigest) != 64 {
		t.Fatalf("a no-op confirmation must still be verified: %d %s", noop.Code, noop.Body.String())
	}
}

func TestStudioV2RecordingRoutes_PostgresIncompatibleTableBlocksPreviewWithoutChange(t *testing.T) {
	fixture, admin, schemaName := postgresRouteTarget(t, "", "")
	seedRecordingRoutePlan(t, fixture, "plan-incompatible", "")
	create := fmt.Sprintf(`CREATE TABLE %s.gw_record_samples (id TEXT PRIMARY KEY, note TEXT)`, postgresRouteIdentifierQuoted(t, schemaName))
	if _, err := admin.ExecContext(context.Background(), create); err != nil {
		t.Fatalf("seed incompatible table: %v", err)
	}

	w := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", fixture.previewBody("plan-incompatible"))

	assertRecordingPlanError(t, w, http.StatusUnprocessableEntity, "RECORDING_SCHEMA_INCOMPATIBLE", false)
	if tables := postgresRouteTables(t, admin, schemaName); len(tables) != 1 || tables[0] != "gw_record_samples" {
		t.Fatalf("a blocked preview must leave the existing table alone, got %v", tables)
	}
	var columns int
	if err := admin.QueryRowContext(context.Background(), `
		SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'gw_record_samples'
	`, schemaName).Scan(&columns); err != nil {
		t.Fatalf("read existing columns: %v", err)
	}
	if columns != 2 {
		t.Fatalf("the existing table must keep its own columns, got %d", columns)
	}
}

func TestStudioV2RecordingRoutes_PostgresPermissionDeniedIsRecordedWithoutChange(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL recording route test: POSTGRES_DSN not set")
	}
	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		t.Fatalf("generate role password: %v", err)
	}
	role, password := fmt.Sprintf("gw_route_role_%d", time.Now().UnixNano()), hex.EncodeToString(secret)
	fixture, admin, schemaName := postgresRouteTarget(t, role, password)
	ctx := context.Background()
	if _, err := admin.ExecContext(ctx, fmt.Sprintf(`CREATE ROLE %s LOGIN PASSWORD '%s'`, postgresRouteIdentifierQuoted(t, role), password)); err != nil {
		t.Fatalf("create restricted role: %v", err)
	}
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP OWNED BY `+postgresRouteIdentifierQuoted(t, role))
		_, _ = admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS `+postgresRouteIdentifierQuoted(t, role))
	})
	// The role may look into the schema but may not create anything in it.
	if _, err := admin.ExecContext(ctx, fmt.Sprintf(`GRANT USAGE ON SCHEMA %s TO %s`,
		postgresRouteIdentifierQuoted(t, schemaName), postgresRouteIdentifierQuoted(t, role))); err != nil {
		t.Fatalf("grant usage: %v", err)
	}
	seedRecordingRoutePlan(t, fixture, "plan-denied", "")
	token := previewForRoute(t, fixture, "plan-denied")

	w := fixture.apply(fixture.applyBody(t, token, nil))

	denied := decodeRecordingOperation(t, w).Data
	if w.Code != http.StatusOK || denied.Status != recordingplan.SchemaOperationFailed ||
		denied.Reason != recordingplan.SchemaReasonPermissionDenied || denied.ExecutedStatements != 0 {
		t.Fatalf("a denied confirmation must be recorded as a permission failure: %d %s", w.Code, w.Body.String())
	}
	if denied.NextAction == "" {
		t.Fatalf("a denied confirmation must name the next step: %+v", denied)
	}
	if tables := postgresRouteTables(t, admin, schemaName); len(tables) != 0 {
		t.Fatalf("a denied confirmation must create nothing, got %v", tables)
	}
}
