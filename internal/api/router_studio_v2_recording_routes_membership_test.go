package api

import (
	"context"
	"database/sql"
	"net/http"
	"path/filepath"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "modernc.org/sqlite"
)

// recordingRouteFixture wires the production router with the real saved
// connector resolver so route behavior matches the deployed service.
type recordingRouteFixture struct {
	router            http.Handler
	planRepo          *recordingplan.MemoryRepository
	planSvc           *recordingplan.Service
	connector         *schema.DatabaseConnector
	repo              *dbtarget.SQLConnectorRepository
	workspaceID       string
	workspaceRevision string
	targetPath        string
}

func newRecordingRouteFixture(t *testing.T) recordingRouteFixture {
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
	repo := dbtarget.NewSQLConnectorRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(repo)
	targetPath := filepath.Join(t.TempDir(), "target.db")
	connector, err := connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name:             "Recording route target",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{"dsn": targetPath},
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
	router := NewRouter(&DatalinkServices{
		Workspace:     workspaceSvc,
		RecordingPlan: planSvc,
		DBTarget:      connectorSvc,
	})
	return recordingRouteFixture{
		router: router, planRepo: planRepo, planSvc: planSvc, connector: connector, repo: repo,
		workspaceID: record.ID, workspaceRevision: record.DatabaseSetupRevision, targetPath: targetPath,
	}
}

func (f recordingRouteFixture) previewBody(planID string) string {
	return `{"plan_id":"` + planID + `","connector_id":"` + f.connector.ID + `","expected_connector_revision":"` + f.connector.IdentityRevision +
		`","expected_workspace_revision":"` + f.workspaceRevision + `","expected_plan_revision":"rev-1"}`
}

func (f recordingRouteFixture) planBody(id string) string {
	return `{"id":"` + id + `","name":"Route plan","members":[{"member_id":"member-a","measurement_id":"meas-a","equipment_id":"equipment-a"}],` +
		`"streams":[{"stream_id":"stream-a","measurement_id":"meas-a","mode":"raw_history"}],` +
		`"destinations":[{"destination_id":"destination-1","connector_id":"` + f.connector.ID + `","connector_revision":"` + f.connector.IdentityRevision + `"}]}`
}

// The production route wiring must leave absent membership readers unset so
// recording-plan creation fails closed instead of calling a nil service.
func TestStudioV2RecordingRoutes_MissingMembershipReadersFailClosed(t *testing.T) {
	fixture := newRecordingRouteFixture(t)

	response := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath, fixture.planBody("plan-wiring"))

	assertRecordingPlanError(t, response, http.StatusServiceUnavailable, "RECORDING_PLAN_UNAVAILABLE", true)
	if _, err := fixture.planRepo.GetPlanByID(context.Background(), "plan-wiring"); err == nil {
		t.Fatal("missing membership readers must not persist a plan")
	}
}

// A disabled saved connector exists, so the operator must be asked to enable
// it rather than being told the recording plan is missing.
func TestStudioV2RecordingRoutes_DisabledSavedConnectorIsUnavailable(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	ctx := context.Background()
	if err := fixture.planSvc.CreatePlan(ctx, &recordingplan.RecordingPlan{
		ID: "plan-disabled", WorkspaceID: fixture.workspaceID, Name: "Disabled target",
		Streams: []recordingplan.PlanStream{{StreamID: "stream-a", MeasurementID: "meas-a", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID: "destination-1", ConnectorID: fixture.connector.ID, ConnectorRevision: fixture.connector.IdentityRevision,
		}},
	}); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	disabled := *fixture.connector
	disabled.Enabled = false
	if err := fixture.repo.UpdateWithExpectedIdentityRevision(ctx, &disabled, fixture.connector.IdentityRevision); err != nil {
		t.Fatalf("disable connector: %v", err)
	}

	createResponse := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath, fixture.planBody("plan-created"))
	assertRecordingPlanError(t, createResponse, http.StatusUnprocessableEntity, "RECORDING_CONNECTOR_UNAVAILABLE", false)
	previewResponse := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", fixture.previewBody("plan-disabled"))
	assertRecordingPlanError(t, previewResponse, http.StatusUnprocessableEntity, "RECORDING_CONNECTOR_UNAVAILABLE", false)
}
