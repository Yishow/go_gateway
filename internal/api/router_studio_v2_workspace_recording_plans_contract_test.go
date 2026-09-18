package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/gin-gonic/gin"
)

const recordingTestWriteNotImplementedMessage = "recording test write is not implemented"

type recordingOperationErrorResponse struct {
	Success bool                           `json:"success"`
	Error   handlers.TypedAPIErrorEnvelope `json:"error"`
	Data    json.RawMessage                `json:"data"`
}

func serveRecordingPlanMutation(router *gin.Engine, path string, body *string, requestID string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(*body)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func assertRecordingValidationError(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected validation status 400, got %d: %s", w.Code, w.Body.String())
	}
	var response recordingOperationErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode validation response: %v; body=%s", err, w.Body.String())
	}
	if response.Success {
		t.Fatalf("validation response must not report success: %s", w.Body.String())
	}
	if response.Error.Code != "validation" {
		t.Fatalf("expected existing validation code, got %q: %s", response.Error.Code, w.Body.String())
	}
}

// assertRecordingTestWriteNotImplemented checks the one recording operation
// that is still unavailable: the response stays safe, non-retryable and free of
// fabricated evidence.
func assertRecordingTestWriteNotImplemented(t *testing.T, w *httptest.ResponseRecorder, wantRequestID string) {
	t.Helper()

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected status 501, got %d: %s", w.Code, w.Body.String())
	}
	var response recordingOperationErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode not-implemented response: %v; body=%s", err, w.Body.String())
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &fields); err != nil {
		t.Fatalf("decode not-implemented response fields: %v; body=%s", err, w.Body.String())
	}
	assertJSONBooleanField(t, fields, "success", false)
	errorRaw, ok := fields["error"]
	if !ok {
		t.Fatalf("not-implemented response is missing error envelope: %s", w.Body.String())
	}
	var errorFields map[string]json.RawMessage
	if err := json.Unmarshal(errorRaw, &errorFields); err != nil {
		t.Fatalf("decode not-implemented error fields: %v; body=%s", err, w.Body.String())
	}
	assertJSONBooleanField(t, errorFields, "retryable", false)
	if response.Success {
		t.Fatalf("not-implemented response must not report success: %s", w.Body.String())
	}
	if response.Error.Code != "RECORDING_TEST_WRITE_NOT_IMPLEMENTED" {
		t.Fatalf("expected the test write not-implemented code, got %q: %s", response.Error.Code, w.Body.String())
	}
	if response.Error.Message != recordingTestWriteNotImplementedMessage {
		t.Fatalf("expected safe message %q, got %q: %s", recordingTestWriteNotImplementedMessage, response.Error.Message, w.Body.String())
	}
	if response.Error.Retryable {
		t.Fatalf("not-implemented response must be non-retryable: %s", w.Body.String())
	}
	if response.Error.Action != "wait_for_supported_operation" {
		t.Fatalf("expected wait action, got %q: %s", response.Error.Action, w.Body.String())
	}
	if wantRequestID != "" && response.Error.RequestID != wantRequestID {
		t.Fatalf("expected request_id %q, got %q: %s", wantRequestID, response.Error.RequestID, w.Body.String())
	}
	if wantRequestID == "" && response.Error.RequestID == "" {
		t.Fatalf("not-implemented response must include a generated request_id: %s", w.Body.String())
	}
	if response.Data != nil && string(response.Data) != "null" {
		t.Fatalf("not-implemented response must not contain data: %s", w.Body.String())
	}
	for _, forbidden := range []string{`"applied":true`, `"written_verified"`, `"record_id"`} {
		if strings.Contains(w.Body.String(), forbidden) {
			t.Fatalf("not-implemented response contains fabricated evidence %q: %s", forbidden, w.Body.String())
		}
	}
}

func assertJSONBooleanField(t *testing.T, fields map[string]json.RawMessage, name string, want bool) {
	t.Helper()

	raw, ok := fields[name]
	if !ok {
		t.Fatalf("JSON response is missing %q", name)
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatalf("decode JSON field %q: %v", name, err)
	}
	got, ok := value.(bool)
	if !ok {
		t.Fatalf("JSON field %q must be boolean, got %s", name, string(raw))
	}
	if got != want {
		t.Fatalf("JSON field %q: got %v, want %v", name, got, want)
	}
}

func seedRecordingPlanForContract(t *testing.T, planSvc *recordingplan.Service, workspaceIDs ...string) *recordingplan.RecordingPlan {
	t.Helper()
	workspaceID := "ws-c1-contract"
	if len(workspaceIDs) > 0 && strings.TrimSpace(workspaceIDs[0]) != "" {
		workspaceID = strings.TrimSpace(workspaceIDs[0])
	}

	plan := &recordingplan.RecordingPlan{
		ID:          "plan-c1-contract",
		WorkspaceID: workspaceID,
		Name:        "C1 Contract Plan",
		Streams: []recordingplan.PlanStream{{
			StreamID:      "stream-c1",
			MeasurementID: "measurement-c1",
			Mode:          recordingplan.StreamModeRawHistory,
			RawPolicy:     recordingplan.RawPolicyEverySample,
		}},
		Destinations: []recordingplan.PlanDestination{{
			DestinationID:     "destination-c1-contract",
			ConnectorID:       "conn-c1",
			ConnectorRevision: "identity-c1",
		}},
	}
	if err := planSvc.CreatePlan(context.Background(), plan); err != nil {
		t.Fatalf("seed recording plan: %v", err)
	}
	return plan
}

func TestStudioV2WorkspaceRecordingPlans_MutationValidationRemains400(t *testing.T) {
	const schemaApplyPath = "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply"
	const testWritePath = "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write"

	tests := []struct {
		name string
		path string
		body *string
	}{
		{name: "schema apply missing body", path: schemaApplyPath},
		{name: "schema apply empty object", path: schemaApplyPath, body: stringPointer(`{}`)},
		{name: "schema apply malformed json", path: schemaApplyPath, body: stringPointer(`{"token":`)},
		{name: "schema apply missing token", path: schemaApplyPath, body: stringPointer(`{"other":"value"}`)},
		{name: "schema apply empty token", path: schemaApplyPath, body: stringPointer(`{"token":""}`)},
		{name: "test write missing body", path: testWritePath},
		{name: "test write empty object", path: testWritePath, body: stringPointer(`{}`)},
		{name: "test write malformed json", path: testWritePath, body: stringPointer(`{"plan_id":`)},
		{name: "test write missing plan id", path: testWritePath, body: stringPointer(`{"stream_id":"stream-c1"}`)},
		{name: "test write empty plan id", path: testWritePath, body: stringPointer(`{"plan_id":""}`)},
	}

	router, _, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertRecordingValidationError(t, serveRecordingPlanMutation(router, tt.path, tt.body, ""))
		})
	}
}

// A confirmation without its operation and revisions is rejected as incomplete
// before any token is read or consumed.
func TestStudioV2WorkspaceRecordingPlans_SchemaApply_IncompleteConfirmationRejectsWithoutRepoMutation(t *testing.T) {
	router, planSvc, planRepo, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	plan := seedRecordingPlanForContract(t, planSvc)

	unseededBody := `{"token":"unseeded-arbitrary-token"}`
	unseededResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply", &unseededBody, "")
	assertRecordingValidationError(t, unseededResponse)

	token := &recordingplan.SchemaPreviewToken{
		Token:        "arbitrary-nonempty-token",
		WorkspaceID:  plan.WorkspaceID,
		PlanID:       plan.ID,
		PlanRevision: plan.Revision,
		ConnectorID:  "connector-c1",
		TablePrefix:  "gw_record_",
		Statements:   []string{"CREATE TABLE should_not_run (id TEXT)"},
		ExpiresAt:    time.Now().UTC().Add(time.Hour),
		CreatedAt:    time.Now().UTC(),
	}
	if err := planRepo.SavePreviewToken(context.Background(), token); err != nil {
		t.Fatalf("seed preview token: %v", err)
	}
	beforePlan, err := planSvc.GetPlan(context.Background(), plan.ID)
	if err != nil {
		t.Fatalf("read plan before schema apply: %v", err)
	}
	beforeToken, err := planRepo.GetPreviewToken(context.Background(), token.Token)
	if err != nil {
		t.Fatalf("read token before schema apply: %v", err)
	}

	body := `{"token":"arbitrary-nonempty-token"}`
	w := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply", &body, "req-c1-full-correlation-id")
	assertRecordingValidationError(t, w)

	afterPlan, err := planSvc.GetPlan(context.Background(), plan.ID)
	if err != nil {
		t.Fatalf("read plan after schema apply: %v", err)
	}
	if !reflect.DeepEqual(beforePlan, afterPlan) {
		t.Fatalf("schema apply changed recording plan: before=%+v after=%+v", beforePlan, afterPlan)
	}
	afterToken, err := planRepo.GetPreviewToken(context.Background(), token.Token)
	if err != nil {
		t.Fatalf("schema apply consumed preview token: %v", err)
	}
	if !reflect.DeepEqual(beforeToken, afterToken) {
		t.Fatalf("schema apply changed preview token: before=%+v after=%+v", beforeToken, afterToken)
	}
}

func TestStudioV2WorkspaceRecordingPlans_TestWrite_UnimplementedRejectsSeededAndMissingPlansWithoutRepoMutation(t *testing.T) {
	router, planSvc, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)
	plan := seedRecordingPlanForContract(t, planSvc)
	beforePlan, err := planSvc.GetPlan(context.Background(), plan.ID)
	if err != nil {
		t.Fatalf("read plan before test write: %v", err)
	}

	tests := []struct {
		name    string
		body    string
		request string
	}{
		{name: "seeded plan", body: `{"plan_id":"plan-c1-contract","stream_id":"stream-c1"}`, request: "req-c1-seeded-test-write"},
		{name: "nonexistent plan", body: `{"plan_id":"nonexistent-plan"}`, request: "req-c1-missing-test-write"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", &tt.body, tt.request)
			assertRecordingTestWriteNotImplemented(t, w, tt.request)

			afterPlan, err := planSvc.GetPlan(context.Background(), plan.ID)
			if err != nil {
				t.Fatalf("read plan after test write: %v", err)
			}
			if !reflect.DeepEqual(beforePlan, afterPlan) {
				t.Fatalf("test write changed recording plan: before=%+v after=%+v", beforePlan, afterPlan)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
