package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
)

func TestStudioV2WorkspaceRecordingPlans_ForeignCRUDUsesScopedLookup(t *testing.T) {
	router, planSvc, planRepo, workspaceID := newStudioV2WorkspaceRecordingPlanRouter(t)
	foreign := &recordingplan.RecordingPlan{
		ID:           "plan-foreign-scope",
		WorkspaceID:  "workspace-other",
		Name:         "Foreign plan",
		Revision:     "rev-1",
		Streams:      []recordingplan.PlanStream{{StreamID: "stream-foreign", MeasurementID: "measurement-foreign", Mode: recordingplan.StreamModeRawHistory}},
		Destinations: []recordingplan.PlanDestination{{DestinationID: "destination-foreign", ConnectorID: "conn-c1", ConnectorRevision: "identity-c1"}},
	}
	if err := planSvc.CreatePlan(context.Background(), foreign); err != nil {
		t.Fatalf("seed foreign plan: %v", err)
	}

	getResponse := serveRecordingScopedRequest(router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-foreign-scope", "")
	assertScopedNotFound(t, getResponse)

	updateBody := `{"name":"must not move plan","streams":[{"stream_id":"stream-foreign","measurement_id":"measurement-foreign","mode":"raw_history"}],"destinations":[{"destination_id":"destination-foreign","connector_id":"conn-c1","connector_revision":"identity-c1"}]}`
	updateResponse := serveRecordingScopedRequest(router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-foreign-scope", updateBody)
	assertScopedNotFound(t, updateResponse)

	deleteResponse := serveRecordingScopedRequest(router, http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/recording-plans/plan-foreign-scope", "")
	assertScopedNotFound(t, deleteResponse)

	remaining, err := planRepo.GetPlanByID(context.Background(), foreign.ID)
	if err != nil {
		t.Fatalf("foreign plan should remain readable to fixture: %v", err)
	}
	if remaining.WorkspaceID != "workspace-other" || remaining.Name != foreign.Name {
		t.Fatalf("foreign plan was changed: %+v", remaining)
	}
	if workspaceID == remaining.WorkspaceID {
		t.Fatalf("fixture workspaces must be distinct: %q", workspaceID)
	}
}

func TestStudioV2WorkspaceRecordingPlans_CreateRejectsForeignWorkspaceBeforeWrites(t *testing.T) {
	router, _, planRepo, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	body := `{"id":"plan-foreign-create","workspace_id":"workspace-other","name":"foreign create","streams":[{"stream_id":"stream-1","measurement_id":"measurement-1","mode":"raw_history"}],"destinations":[{"destination_id":"destination-1","connector_id":"conn-c1","connector_revision":"identity-c1"}]}`

	response := serveRecordingScopedRequest(router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/recording-plans", body)
	assertScopedNotFound(t, response)
	if _, err := planRepo.GetPlanByID(context.Background(), "plan-foreign-create"); err == nil {
		t.Fatal("foreign workspace create must not persist a plan")
	}
}

func serveRecordingScopedRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequestWithContext(context.Background(), method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func assertScopedNotFound(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if response.Code != http.StatusNotFound {
		t.Fatalf("foreign resource must return safe 404, got %d: %s", response.Code, response.Body.String())
	}
	var payload struct {
		Success bool                           `json:"success"`
		Error   handlersTypedErrorEnvelopeTest `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode safe 404: %v; body=%s", err, response.Body.String())
	}
	if payload.Success || payload.Error.Code != "RECORDING_PLAN_NOT_FOUND" {
		t.Fatalf("unexpected safe 404 payload: %s", response.Body.String())
	}
}

type handlersTypedErrorEnvelopeTest struct {
	Code string `json:"code"`
}
