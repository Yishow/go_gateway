package api

import (
	"context"
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "modernc.org/sqlite"
)

func TestStudioV2WorkspaceRecordingPlans_UnimplementedMutationsLeaveConfiguredTargetUntouched(t *testing.T) {
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
	t.Cleanup(func() { _ = targetDB.Close() })
	if _, err := targetDB.ExecContext(ctx, `CREATE TABLE sentinel (id INTEGER PRIMARY KEY, value TEXT NOT NULL)`); err != nil {
		t.Fatalf("create target sentinel table: %v", err)
	}
	if _, err := targetDB.ExecContext(ctx, `INSERT INTO sentinel (id, value) VALUES (1, 'before')`); err != nil {
		t.Fatalf("seed target sentinel row: %v", err)
	}

	connectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo)
	connector, err := connectorSvc.Create(ctx, dbtarget.CreateConnectorRequest{
		Name: "C1 target witness",
		Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"dsn": targetPath,
		},
	})
	if err != nil {
		t.Fatalf("create target connector: %v", err)
	}

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)
	router := NewRouter(&DatalinkServices{
		Workspace:     workspaceSvc,
		RecordingPlan: planSvc,
		DBTarget:      connectorSvc,
	})

	plan := &recordingplan.RecordingPlan{
		ID:          "plan-c1-target-witness",
		WorkspaceID: "ws-c1-target-witness",
		Name:        "C1 Target Witness Plan",
		Streams: []recordingplan.PlanStream{{
			StreamID:      "stream-c1-target",
			MeasurementID: "measurement-c1-target",
			Mode:          recordingplan.StreamModeRawHistory,
			RawPolicy:     recordingplan.RawPolicyEverySample,
		}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID: "destination-c1-target",
			ConnectorID:   connector.ID,
			TablePrefix:   "gw_record_",
		}},
	}
	if err := planSvc.CreatePlan(ctx, plan); err != nil {
		t.Fatalf("seed target recording plan: %v", err)
	}
	previewToken := &recordingplan.SchemaPreviewToken{
		Token:        "target-witness-token",
		WorkspaceID:  plan.WorkspaceID,
		PlanID:       plan.ID,
		PlanRevision: plan.Revision,
		ConnectorID:  connector.ID,
		TablePrefix:  "gw_record_",
		Statements:   []string{"CREATE TABLE should_not_run (id TEXT)"},
		ExpiresAt:    time.Now().UTC().Add(time.Hour),
		CreatedAt:    time.Now().UTC(),
	}
	if err := planRepo.SavePreviewToken(ctx, previewToken); err != nil {
		t.Fatalf("seed target preview token: %v", err)
	}

	beforeTables, beforeValue := readTargetWitnessState(ctx, t, targetDB)
	schemaBody := `{"token":"target-witness-token"}`
	schemaResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply", &schemaBody, "req-c1-target-schema")
	assertRecordingValidationError(t, schemaResponse)

	testWriteBody := `{"plan_id":"plan-c1-target-witness","stream_id":"stream-c1-target"}`
	testWriteResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", &testWriteBody, "req-c1-target-write")
	assertRecordingTestWriteNotImplemented(t, testWriteResponse, "req-c1-target-write")

	afterTables, afterValue := readTargetWitnessState(ctx, t, targetDB)
	if !reflect.DeepEqual(beforeTables, afterTables) || beforeValue != afterValue {
		t.Fatalf("unimplemented recording mutations changed target: tables before=%v after=%v sentinel before=%q after=%q", beforeTables, afterTables, beforeValue, afterValue)
	}
}

func readTargetWitnessState(ctx context.Context, t *testing.T, db *sql.DB) (tables []string, value string) {
	t.Helper()

	rows, err := db.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`)
	if err != nil {
		t.Fatalf("read target tables: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan target table: %v", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate target tables: %v", err)
	}
	if err := db.QueryRowContext(ctx, `SELECT value FROM sentinel WHERE id = 1`).Scan(&value); err != nil {
		t.Fatalf("read target sentinel: %v", err)
	}
	return tables, value
}
