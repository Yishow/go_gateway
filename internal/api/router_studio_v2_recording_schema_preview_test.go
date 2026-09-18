package api

import (
	"context"
	"database/sql"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	_ "modernc.org/sqlite"
)

func seedRecordingRoutePlan(t *testing.T, fixture recordingRouteFixture, id, tablePrefix string) {
	t.Helper()
	if err := fixture.planSvc.CreatePlan(context.Background(), &recordingplan.RecordingPlan{
		ID: id, WorkspaceID: fixture.workspaceID, Name: "Preview target",
		Streams: []recordingplan.PlanStream{{StreamID: "stream-a", MeasurementID: "meas-a", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID: "destination-1", ConnectorID: fixture.connector.ID, ConnectorRevision: fixture.connector.IdentityRevision,
			TablePrefix: tablePrefix,
		}},
	}); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
}

func openRecordingTarget(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open target database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func recordingTargetTables(t *testing.T, db *sql.DB) []string {
	t.Helper()
	rows, err := db.QueryContext(context.Background(), `SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name`)
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

// A preview only reads the saved SQLite target; the operator must confirm
// before any managed table appears there.
func TestStudioV2RecordingRoutes_SchemaPreviewHasNoTargetSideEffects(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-preview", "")

	w := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", fixture.previewBody("plan-preview"))

	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	data := decodeRecordingPreview(t, w)
	if data.Dialect != "sqlite" || data.Schema != "main" || data.Database != fixture.targetPath || data.WorkspaceRevision != fixture.workspaceRevision {
		t.Fatalf("preview must be bound to the saved SQLite target: %+v", data)
	}
	if len(data.Statements) == 0 || len(data.Tables) != 6 || data.Tables[0].Action != recordingplan.SchemaTableActionCreate {
		t.Fatalf("preview must describe the tables it would create: %+v", data)
	}
	if tables := recordingTargetTables(t, openRecordingTarget(t, fixture.targetPath)); len(tables) != 0 {
		t.Fatalf("preview must not create target tables: %v", tables)
	}
	stored, err := fixture.planRepo.GetPreviewToken(context.Background(), data.Token)
	if err != nil || stored.Digest != data.Digest || stored.OperationID != data.OperationID || stored.WorkspaceRevision != fixture.workspaceRevision {
		t.Fatalf("preview must persist its protected scope: stored=%+v err=%v", stored, err)
	}
}

func TestStudioV2RecordingRoutes_SchemaPreviewExplainsCompatibleExistingSchema(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-compatible", "")
	target := openRecordingTarget(t, fixture.targetPath)
	statements, err := recordingplan.GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatalf("generate managed schema: %v", err)
	}
	for _, statement := range statements {
		if _, err := target.ExecContext(context.Background(), statement); err != nil {
			t.Fatalf("seed managed schema: %v", err)
		}
	}
	before := recordingTargetTables(t, target)

	w := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", fixture.previewBody("plan-compatible"))

	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	data := decodeRecordingPreview(t, w)
	if len(data.Statements) != 0 || data.NoChangeReason != recordingplan.NoChangeSchemaCompatible || len(data.Tables) != 6 {
		t.Fatalf("compatible schema must be an explained no-op: %+v", data)
	}
	if after := recordingTargetTables(t, target); !reflect.DeepEqual(after, before) {
		t.Fatalf("compatible preview changed target tables: before=%v after=%v", before, after)
	}
}

// The saved plan destination owns the table scope; a client that still shows
// another prefix must preview again instead of creating a different set.
func TestStudioV2RecordingRoutes_SchemaPreviewUsesSavedDestinationPrefix(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-prefix", "gw_line_")
	body := fixture.previewBody("plan-prefix")
	stale := strings.TrimSuffix(body, "}") + `,"table_prefix":"gw_record_"}`

	staleResponse := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", stale)
	assertRecordingPlanError(t, staleResponse, http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)

	w := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", body)
	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	data := decodeRecordingPreview(t, w)
	if data.TablePrefix != "gw_line_" || len(data.Tables) == 0 || data.Tables[0].Name != "gw_line_samples" {
		t.Fatalf("preview must use the saved destination prefix: %+v", data)
	}
}
