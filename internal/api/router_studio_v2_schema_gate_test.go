package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	_ "modernc.org/sqlite"
)

type schemaGateRouteFixture struct {
	router      http.Handler
	connectorID string
	targetPath  string
}

// newSchemaGateRouteFixture wires the production router with every service the
// schema endpoints need, and a saved SQLite target that starts out empty.
func newSchemaGateRouteFixture(t *testing.T) schemaGateRouteFixture {
	t.Helper()
	ctx := context.Background()
	mainDB, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "main.db"))
	if err != nil {
		t.Fatalf("open main database: %v", err)
	}
	t.Cleanup(func() { _ = mainDB.Close() })
	if err := datalinkbase.NewMigrator().Migrate(mainDB); err != nil {
		t.Fatalf("migrate main database: %v", err)
	}
	targetPath := filepath.Join(t.TempDir(), "target.db")
	targetDB, err := sql.Open("sqlite", targetPath)
	if err != nil {
		t.Fatalf("open target database: %v", err)
	}
	if err := targetDB.PingContext(ctx); err != nil {
		t.Fatalf("reach target database: %v", err)
	}
	if err := targetDB.Close(); err != nil {
		t.Fatalf("close target database: %v", err)
	}

	deviceSvc := device.NewService(device.NewSQLRepository(mainDB), nil)
	pointSvc := point.NewService(point.NewSQLRepository(mainDB), nil)
	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	ruleSvc := sourcerule.NewService(sourcerule.NewSQLRepository(mainDB), deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspace.NewSQLRepository(mainDB))
	connectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	mappingRepo := dbtarget.NewSQLTargetMappingRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := dbtarget.NewMappingService(mappingRepo, connectorRepo, tagSvc)

	connector, err := connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "Schema gate target", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{"dsn": targetPath},
	})
	if err != nil {
		t.Fatalf("create connector: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(ctx, connector.ID); err != nil {
		t.Fatalf("bind connector: %v", err)
	}

	router := NewRouter(&DatalinkServices{
		Workspace: workspaceSvc, Device: deviceSvc, Point: pointSvc, Tag: tagSvc,
		SourceRule: ruleSvc, DBTarget: connectorSvc, DBMapping: mappingSvc,
	})
	return schemaGateRouteFixture{router: router, connectorID: connector.ID, targetPath: targetPath}
}

func (f schemaGateRouteFixture) generate(t *testing.T, path string, dryRun bool) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(map[string]bool{"dry_run": dryRun})
	if err != nil {
		t.Fatalf("encode schema request: %v", err)
	}
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	f.router.ServeHTTP(w, req)
	return w
}

func schemaGateTargetTables(t *testing.T, path string) []string {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open target database: %v", err)
	}
	defer func() { _ = db.Close() }()
	rows, err := db.QueryContext(context.Background(), `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	if err != nil {
		t.Fatalf("list target tables: %v", err)
	}
	defer func() { _ = rows.Close() }()
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan target table: %v", err)
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate target tables: %v", err)
	}
	return names
}

func assertSchemaConfirmationRequired(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusConflict {
		t.Fatalf("creating without a confirmed preview must be refused, got %d: %s", w.Code, w.Body.String())
	}
	var response struct {
		Success bool `json:"success"`
		Error   struct {
			Code   string `json:"code"`
			Action string `json:"action"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode refusal: %v; body=%s", err, w.Body.String())
	}
	if response.Success || response.Error.Code != "SCHEMA_CONFIRMATION_REQUIRED" || response.Error.Action == "" {
		t.Fatalf("refusal must name the confirmed flow: %s", w.Body.String())
	}
}

// Both public schema endpoints keep their read-only plan and refuse to create
// anything without a confirmed preview, through the production router.
func TestStudioV2Routes_SchemaGenerationNeedsConfirmation(t *testing.T) {
	fixture := newSchemaGateRouteFixture(t)
	before := schemaGateTargetTables(t, fixture.targetPath)

	for _, endpoint := range []struct{ name, path string }{
		{name: "workspace", path: "/api/v1/datalink/studio-v2/workspace/database-schema/generate"},
		{name: "generic", path: "/api/v1/datalink/db-targets/connectors/" + fixture.connectorID + "/schema/generate"},
	} {
		t.Run(endpoint.name, func(t *testing.T) {
			assertSchemaConfirmationRequired(t, fixture.generate(t, endpoint.path, false))

			readOnly := fixture.generate(t, endpoint.path, true)
			if readOnly.Code != http.StatusOK {
				t.Fatalf("the read-only plan must stay available, got %d: %s", readOnly.Code, readOnly.Body.String())
			}
			if after := schemaGateTargetTables(t, fixture.targetPath); len(after) != len(before) {
				t.Fatalf("no schema request may change the target: before=%v after=%v", before, after)
			}
		})
	}
}
