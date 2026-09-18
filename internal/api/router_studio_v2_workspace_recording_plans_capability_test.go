package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
)

func TestStudioV2WorkspaceRecordingPlans_CapabilitiesExposeVerifiedManagedSchema(t *testing.T) {
	router, _, _, _ := newStudioV2WorkspaceRecordingPlanRouter(t)

	listRequest := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/capabilities", http.NoBody)
	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, listRequest)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("capabilities list failed: %d %s", listResponse.Code, listResponse.Body.String())
	}
	var listBody struct {
		Success bool                                `json:"success"`
		Data    []recordingplan.ConnectorCapability `json:"data"`
	}
	if err := json.Unmarshal(listResponse.Body.Bytes(), &listBody); err != nil {
		t.Fatalf("decode capabilities list: %v; body=%s", err, listResponse.Body.String())
	}
	if !listBody.Success {
		t.Fatalf("capabilities list must succeed: %s", listResponse.Body.String())
	}
	if len(listBody.Data) != 5 {
		t.Fatalf("expected five capability entries, got %d: %s", len(listBody.Data), listResponse.Body.String())
	}

	capabilities := make(map[string]recordingplan.ConnectorCapability, len(listBody.Data))
	for _, capability := range listBody.Data {
		capabilities[capability.Kind] = capability
	}
	// Managed schema execution is verified for SQLite/PostgreSQL, so the
	// capability matrix must advertise it; test writes stay masked until wired.
	for _, kind := range []string{"sqlite", "postgres"} {
		capability, ok := capabilities[kind]
		if !ok {
			t.Fatalf("capabilities list is missing %s: %s", kind, listResponse.Body.String())
		}
		if !capability.Supported || !capability.SupportsManagedSchema || capability.SupportsTestWrites {
			t.Fatalf("%s capability must expose verified managed schema without test writes: %+v", kind, capability)
		}
		if !capability.SupportsTransactions || !capability.SupportsReceipts {
			t.Fatalf("%s capability lost existing transaction/receipt support: %+v", kind, capability)
		}
		if len(capability.SupportedModes) != 2 || capability.SupportedModes[0] != "managed_recording" || capability.SupportedModes[1] != "custom_table" {
			t.Fatalf("%s capability must retain managed/custom query modes: %+v", kind, capability)
		}
	}
	for _, kind := range []string{"mysql"} {
		capability, ok := capabilities[kind]
		if !ok {
			t.Fatalf("capabilities list is missing %s: %s", kind, listResponse.Body.String())
		}
		if !capability.Supported || capability.SupportsManagedSchema || capability.SupportsTestWrites {
			t.Fatalf("%s capability must keep connection support while masking unverified mutations: %+v", kind, capability)
		}
		if capability.Notes == "" {
			t.Fatalf("%s capability must explain why managed schema is unavailable: %+v", kind, capability)
		}
	}
	for _, kind := range []string{"sqlserver", "oracle"} {
		capability, ok := capabilities[kind]
		if !ok {
			t.Fatalf("capabilities list is missing %s: %s", kind, listResponse.Body.String())
		}
		if capability.Supported || capability.SupportsManagedSchema || capability.SupportsTestWrites || len(capability.SupportedModes) != 0 {
			t.Fatalf("%s must remain unsupported: %+v", kind, capability)
		}
	}

	aliasRequest := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/recording-plans/capabilities?kind=sqlite3", http.NoBody)
	aliasResponse := httptest.NewRecorder()
	router.ServeHTTP(aliasResponse, aliasRequest)
	if aliasResponse.Code != http.StatusOK {
		t.Fatalf("sqlite3 capability query failed: %d %s", aliasResponse.Code, aliasResponse.Body.String())
	}
	var aliasBody struct {
		Success bool                              `json:"success"`
		Data    recordingplan.ConnectorCapability `json:"data"`
	}
	if err := json.Unmarshal(aliasResponse.Body.Bytes(), &aliasBody); err != nil {
		t.Fatalf("decode sqlite3 capability: %v; body=%s", err, aliasResponse.Body.String())
	}
	if !aliasBody.Success || aliasBody.Data.Kind != "sqlite3" || !aliasBody.Data.Supported || !aliasBody.Data.SupportsManagedSchema || aliasBody.Data.SupportsTestWrites {
		t.Fatalf("sqlite3 alias capability mismatch: %+v", aliasBody)
	}
}

func TestStudioV2WorkspaceRecordingPlans_PreviewRemainsAvailableWithMaskedMutations(t *testing.T) {
	fixture := newStudioV2WorkspaceRecordingPlanFixture(t)
	router, planSvc, workspaceID := fixture.router, fixture.planSvc, fixture.workspaceID
	plan := seedRecordingPlanForContract(t, planSvc, workspaceID)

	previewBody := `{"plan_id":"plan-c1-contract","connector_id":"conn-c1","expected_connector_revision":"identity-c1",` +
		`"expected_workspace_revision":"` + fixture.workspaceRevision + `","expected_plan_revision":"rev-1","table_prefix":"gw_record_","dialect":"sqlite"}`
	previewResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-preview", &previewBody, "req-c1-capability-preview")
	if previewResponse.Code != http.StatusOK {
		t.Fatalf("schema preview must remain available: %d %s", previewResponse.Code, previewResponse.Body.String())
	}
	var previewBodyDecoded struct {
		Success bool                             `json:"success"`
		Data    recordingplan.SchemaPreviewToken `json:"data"`
	}
	if err := json.Unmarshal(previewResponse.Body.Bytes(), &previewBodyDecoded); err != nil {
		t.Fatalf("decode schema preview: %v; body=%s", err, previewResponse.Body.String())
	}
	if !previewBodyDecoded.Success || previewBodyDecoded.Data.Token == "" || len(previewBodyDecoded.Data.Statements) == 0 {
		t.Fatalf("schema preview response lost usable preview data: %+v", previewBodyDecoded)
	}

	schemaBody := `{"token":"unseeded-capability-token"}`
	schemaResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply", &schemaBody, "req-c1-capability-schema")
	assertRecordingValidationError(t, schemaResponse)

	testWriteBody := `{"plan_id":"plan-c1-contract","stream_id":"stream-c1"}`
	testWriteResponse := serveRecordingPlanMutation(router, "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write", &testWriteBody, "req-c1-capability-write")
	assertRecordingTestWriteNotImplemented(t, testWriteResponse, "req-c1-capability-write")

	if _, err := planSvc.GetPlan(context.Background(), plan.ID); err != nil {
		t.Fatalf("preview/capability checks must keep seeded plan queryable: %v", err)
	}
}
