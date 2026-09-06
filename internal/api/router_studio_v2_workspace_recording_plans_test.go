package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func newStudioV2WorkspaceRecordingPlanRouter(t *testing.T) (*gin.Engine, *workspace.Service, *recordingplan.Service) {
	t.Helper()

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	planRepo := recordingplan.NewMemoryRepository()
	planSvc := recordingplan.NewService(planRepo)

	router := NewRouter(&DatalinkServices{
		Workspace:     workspaceSvc,
		RecordingPlan: planSvc,
	})

	return router, workspaceSvc, planSvc
}

func TestStudioV2WorkspaceRecordingPlans_CRUD(t *testing.T) {
	router, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)

	// 1. Create plan
	createBody := `{
		"id": "plan-meter",
		"name": "Power Meter Plan",
		"timezone": "Asia/Taipei",
		"streams": [
			{
				"stream_id": "s-kw",
				"measurement_id": "meas-kw",
				"mode": "raw_history",
				"raw_policy": "every_sample"
			}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create plan failed with status %d: %s", w.Code, w.Body.String())
	}

	// 2. List plans
	req = httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list plans failed: %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "plan-meter") {
		t.Fatalf("expected plan-meter in list: %s", w.Body.String())
	}

	// 3. Get plan
	req = httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("get plan failed: %d %s", w.Code, w.Body.String())
	}

	// 4. Update plan
	updateBody := `{
		"name": "Power Meter Plan (Renamed)",
		"timezone": "Asia/Taipei",
		"streams": [
			{
				"stream_id": "s-kw",
				"measurement_id": "meas-kw",
				"mode": "raw_history",
				"raw_policy": "every_sample"
			}
		]
	}`
	req = httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("update plan failed: %d %s", w.Code, w.Body.String())
	}

	// 5. Delete plan
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-meter", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("delete plan failed: %d %s", w.Code, w.Body.String())
	}
}

func TestStudioV2WorkspaceRecordingPlans_Capabilities_And_SchemaPreview(t *testing.T) {
	router, _, planSvc := newStudioV2WorkspaceRecordingPlanRouter(t)

	// 1. Capabilities
	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/capabilities", nil)
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
		WorkspaceID: "ws-test",
		Name:        "Schema Test Plan",
		Streams:     []recordingplan.PlanStream{{StreamID: "s1", MeasurementID: "m1", Mode: recordingplan.StreamModeRawHistory}},
	}
	_ = planSvc.CreatePlan(context.Background(), plan)

	// 2. Schema Preview
	previewBody := `{
		"plan_id": "plan-schema-test",
		"connector_id": "conn-sqlite",
		"table_prefix": "gw_record_",
		"dialect": "sqlite"
	}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-preview", strings.NewReader(previewBody))
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
	req = httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", strings.NewReader(testWriteBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("test write failed: %d %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "written_verified") {
		t.Fatalf("expected written_verified in test write response: %s", w.Body.String())
	}
}
