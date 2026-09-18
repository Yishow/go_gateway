package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

type recordingPlanTestResolver struct {
	connector *schema.DatabaseConnector
}

func (r recordingPlanTestResolver) ResolveSavedTarget(_ context.Context, _, _ string) (*schema.DatabaseConnector, error) {
	return r.connector, nil
}

// ExecuteSchemaStatements refuses: this fixture has no real target database.
func (r recordingPlanTestResolver) ExecuteSchemaStatements(context.Context, string, string, []string) (*dbtarget.SchemaStatementExecution, error) {
	return nil, errors.New("fixture target cannot execute schema statements")
}

// InspectTable reports every managed table as missing so previews list them all.
func (r recordingPlanTestResolver) InspectTable(_ context.Context, _, schemaName, tableName string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{Status: dbtarget.TableInspectionMissing, Schema: schemaName, Table: tableName}, nil
}

type studioV2WorkspaceRecordingPlanFixture struct {
	router            *gin.Engine
	planSvc           *recordingplan.Service
	planRepo          *recordingplan.MemoryRepository
	workspaceID       string
	workspaceRevision string
}

func newStudioV2WorkspaceRecordingPlanRouter(t *testing.T) (engine *gin.Engine, planService *recordingplan.Service, planRepository *recordingplan.MemoryRepository, workspaceID string) {
	t.Helper()
	fixture := newStudioV2WorkspaceRecordingPlanFixture(t)
	return fixture.router, fixture.planSvc, fixture.planRepo, fixture.workspaceID
}

func newStudioV2WorkspaceRecordingPlanFixture(t *testing.T) studioV2WorkspaceRecordingPlanFixture {
	t.Helper()

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	workspaceRecord, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("create workspace fixture: %v", err)
	}
	if _, err := workspaceSvc.BindDatabaseConnector(context.Background(), "conn-c1"); err != nil {
		t.Fatalf("bind connector fixture: %v", err)
	}
	for _, deviceID := range []string{"device-a", "device-b"} {
		if _, err := workspaceSvc.AttachDevice(context.Background(), deviceID); err != nil {
			t.Fatalf("attach device fixture %s: %v", deviceID, err)
		}
	}
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)

	handler := handlers.NewStudioV2WorkspaceRecordingPlansHandler(workspaceSvc, planSvc, recordingPlanTestResolver{connector: &schema.DatabaseConnector{
		ID:               "conn-c1",
		Kind:             schema.DatabaseConnectorKindSQLite,
		IdentityRevision: "identity-c1",
		Enabled:          true,
	}})
	membership := newRecordingMembershipStub(workspaceRecord.ID)
	handler.SetRecordingMembershipServices(membership, membership)
	router := gin.New()
	router.GET("/api/v1/datalink/studio-v2/workspace/recording-plans", handler.List)
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans", handler.Create)
	router.GET("/api/v1/datalink/studio-v2/workspace/recording-plans/:id", handler.Get)
	router.PUT("/api/v1/datalink/studio-v2/workspace/recording-plans/:id", handler.Update)
	router.DELETE("/api/v1/datalink/studio-v2/workspace/recording-plans/:id", handler.Delete)
	router.GET("/api/v1/datalink/studio-v2/workspace/recording-plans/capabilities", handler.Capabilities)
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans/schema-preview", handler.SchemaPreview)
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply", handler.SchemaApply)
	router.POST("/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", handler.TestWrite)

	current, err := workspaceSvc.GetOrCreate(context.Background())
	if err != nil {
		t.Fatalf("read workspace revision: %v", err)
	}
	return studioV2WorkspaceRecordingPlanFixture{
		router: router, planSvc: planSvc, planRepo: planRepo,
		workspaceID: workspaceRecord.ID, workspaceRevision: current.DatabaseSetupRevision,
	}
}

func TestStudioV2WorkspaceRecordingPlans_CRUD(t *testing.T) {
	router, _, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)

	// 1. Create plan
	createBody := `{
		"id": "plan-meter",
		"name": "Power Meter Plan",
		"members": [{"member_id": "member-kw", "measurement_id": "meas-kw", "equipment_id": "equipment-meter"}],
		"timezone": "Asia/Taipei",
		"streams": [
			{
				"stream_id": "s-kw",
				"measurement_id": "meas-kw",
				"mode": "raw_history",
				"raw_policy": "every_sample"
			}
		],
		"destinations": [{
			"destination_id": "destination-meter",
			"connector_id": "conn-c1",
			"connector_revision": "identity-c1"
		}]
	}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create plan failed with status %d: %s", w.Code, w.Body.String())
	}

	// 2. List plans
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list plans failed: %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "plan-meter") {
		t.Fatalf("expected plan-meter in list: %s", w.Body.String())
	}

	// 3. Get plan
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("get plan failed: %d %s", w.Code, w.Body.String())
	}

	// 4. Update plan
	updateBody := `{
		"name": "Power Meter Plan (Renamed)",
		"members": [{"member_id": "member-kw", "measurement_id": "meas-kw", "equipment_id": "equipment-meter"}],
		"timezone": "Asia/Taipei",
		"streams": [
			{
				"stream_id": "s-kw",
				"measurement_id": "meas-kw",
				"mode": "raw_history",
				"raw_policy": "every_sample"
			}
		],
		"destinations": [{
			"destination_id": "destination-meter",
			"connector_id": "conn-c1",
			"connector_revision": "identity-c1"
		}]
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("update plan failed: %d %s", w.Code, w.Body.String())
	}

	// 5. Delete plan
	req = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("delete plan failed: %d %s", w.Code, w.Body.String())
	}
}

func TestStudioV2WorkspaceRecordingPlans_Capabilities_And_SchemaPreview(t *testing.T) {
	fixture := newStudioV2WorkspaceRecordingPlanFixture(t)
	router, planSvc, workspaceID := fixture.router, fixture.planSvc, fixture.workspaceID

	// 1. Capabilities
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/capabilities", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("capabilities query failed: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "sqlite") || !strings.Contains(w.Body.String(), "postgres") {
		t.Fatalf("expected sqlite and postgres in capabilities: %s", w.Body.String())
	}

	// Seed plan
	plan := &recordingplan.RecordingPlan{
		ID:          "plan-schema-test",
		WorkspaceID: workspaceID,
		Name:        "Schema Test Plan",
		Streams:     []recordingplan.PlanStream{{StreamID: "s1", MeasurementID: "m1", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID:     "destination-schema-test",
			ConnectorID:       "conn-c1",
			ConnectorRevision: "identity-c1",
		}},
	}
	_ = planSvc.CreatePlan(context.Background(), plan)

	// 2. Schema Preview
	previewBody := `{
		"plan_id": "plan-schema-test",
		"connector_id": "conn-c1",
		"expected_connector_revision": "identity-c1",
		"expected_workspace_revision": "` + fixture.workspaceRevision + `",
		"expected_plan_revision": "rev-1",
		"table_prefix": "gw_record_",
		"dialect": "sqlite"
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-preview", strings.NewReader(previewBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "tok-") {
		t.Fatalf("expected token in preview response: %s", w.Body.String())
	}

	// 3. Test Write
	testWriteBody := `{
		"plan_id": "plan-schema-test",
		"stream_id": "s1",
		"table_prefix": "gw_record_"
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", strings.NewReader(testWriteBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "req-c1-existing-seeded-plan")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertRecordingTestWriteNotImplemented(t, w, "req-c1-existing-seeded-plan")
}
