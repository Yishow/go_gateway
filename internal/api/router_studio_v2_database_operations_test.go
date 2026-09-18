package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

const databaseOperationsPath = "/api/v1/datalink/studio-v2/workspace/database-operations/"

type recordingOperationResponse struct {
	Success bool                          `json:"success"`
	Data    recordingplan.SchemaOperation `json:"data"`
	Error   struct {
		Code        string `json:"code"`
		OperationID string `json:"operation_id"`
	} `json:"error"`
}

func decodeRecordingOperation(t *testing.T, w *httptest.ResponseRecorder) recordingOperationResponse {
	t.Helper()
	var response recordingOperationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode operation response: %v; body=%s", err, w.Body.String())
	}
	return response
}

// previewForRoute issues a real preview through the router and returns the stored token.
func previewForRoute(t *testing.T, fixture recordingRouteFixture, planID string) *recordingplan.SchemaPreviewToken {
	t.Helper()
	w := serveRecordingScopedRequest(fixture.router, http.MethodPost, recordingPlansPath+"/schema-preview", fixture.previewBody(planID))
	if w.Code != http.StatusOK {
		t.Fatalf("schema preview failed: %d %s", w.Code, w.Body.String())
	}
	stored, err := fixture.planRepo.GetPreviewToken(context.Background(), decodeRecordingPreview(t, w).Token)
	if err != nil {
		t.Fatalf("read stored preview: %v", err)
	}
	return stored
}

func TestStudioV2DatabaseOperations_StatusIsWorkspaceScoped(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-status", "")
	seedRecordingRoutePlan(t, fixture, "plan-status-other", "gw_other_")
	ctx := context.Background()
	token := previewForRoute(t, fixture, "plan-status")
	op, outcome, err := fixture.planSvc.ClaimSchemaApply(ctx, token)
	if err != nil || outcome != recordingplan.ClaimAcquired {
		t.Fatalf("claim operation: outcome=%q err=%v", outcome, err)
	}
	statusOf := func(operationID string) *httptest.ResponseRecorder {
		return serveRecordingScopedRequest(fixture.router, http.MethodGet, databaseOperationsPath+operationID, "")
	}

	w := statusOf(op.OperationID)
	if got := decodeRecordingOperation(t, w); w.Code != http.StatusOK || got.Data.OperationID != op.OperationID || got.Data.Status != recordingplan.SchemaOperationRunning {
		t.Fatalf("running operation must be readable: %d %s", w.Code, w.Body.String())
	}
	if body := w.Body.String(); strings.Contains(body, token.Token) || strings.Contains(body, op.Owner) || strings.Contains(body, "scope_key") {
		t.Fatalf("status must not expose the token, owner or scope: %s", body)
	}
	if _, err := fixture.planSvc.FinishSchemaOperation(ctx, op.OperationID, op.Owner, recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationPartial, ExecutedStatements: 2, Reason: "statement_failed", NextAction: "preview again to repair",
	}); err != nil {
		t.Fatalf("finish operation: %v", err)
	}
	w = statusOf(op.OperationID)
	if got := decodeRecordingOperation(t, w).Data; w.Code != http.StatusOK || got.Status != recordingplan.SchemaOperationPartial ||
		got.ExecutedStatements != 2 || got.NextAction != "preview again to repair" || got.CompletedAt == nil {
		t.Fatalf("partial operation must be readable with its repair action: %d %s", w.Code, w.Body.String())
	}

	other, _, err := fixture.planSvc.ClaimSchemaApply(ctx, previewForRoute(t, fixture, "plan-status-other"))
	if err != nil {
		t.Fatalf("claim other operation: %v", err)
	}
	if _, err := fixture.planSvc.FinishSchemaOperation(ctx, other.OperationID, other.Owner, recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationUnknown, Reason: "acknowledgement_lost",
	}); err != nil {
		t.Fatalf("finish other operation: %v", err)
	}
	if w := statusOf(other.OperationID); w.Code != http.StatusOK || decodeRecordingOperation(t, w).Data.Status != recordingplan.SchemaOperationUnknown {
		t.Fatalf("unknown operation must stay readable: %d %s", w.Code, w.Body.String())
	}

	now := time.Now().UTC()
	if _, _, err := fixture.planRepo.ClaimSchemaOperation(ctx, &recordingplan.SchemaOperation{
		OperationID: "op-foreign", Token: "tok-foreign", WorkspaceID: "workspace-foreign", ScopeKey: "scope-foreign", Owner: "claim-foreign",
		Action: recordingplan.SchemaApplyAction, Status: recordingplan.SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("seed foreign operation: %v", err)
	}
	for _, id := range []string{"op-foreign", "op-missing"} {
		w := statusOf(id)
		assertRecordingPlanError(t, w, http.StatusNotFound, "RECORDING_OPERATION_NOT_FOUND", false)
		if strings.Contains(w.Body.String(), "workspace-foreign") || strings.Contains(w.Body.String(), "running") {
			t.Fatalf("foreign or unknown operation must not leak details: %s", w.Body.String())
		}
	}
}

// A repeated confirmation reports the operation already recorded for its token
// or the one holding its scope instead of executing again.
func TestStudioV2RecordingRoutes_SchemaApplyReplaysTheIssuedOperation(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-replay", "")
	first := previewForRoute(t, fixture, "plan-replay")
	second := previewForRoute(t, fixture, "plan-replay")
	apply := func(token *recordingplan.SchemaPreviewToken, overrides map[string]string) *httptest.ResponseRecorder {
		return fixture.apply(fixture.applyBody(t, token, overrides))
	}

	assertRecordingPlanError(t, apply(first, map[string]string{"operation_id": "op-substituted"}), http.StatusConflict, "RECORDING_SCHEMA_OPERATION_MISMATCH", false)
	assertRecordingPlanError(t, apply(first, map[string]string{"token": "tok-unknown"}), http.StatusNotFound, "RECORDING_SCHEMA_PREVIEW_NOT_FOUND", false)

	op, _, err := fixture.planSvc.ClaimSchemaApply(context.Background(), first)
	if err != nil {
		t.Fatalf("claim operation: %v", err)
	}
	running := apply(first, nil)
	if running.Code != http.StatusAccepted || decodeRecordingOperation(t, running).Data.OperationID != op.OperationID {
		t.Fatalf("a running duplicate must return 202 with the same operation: %d %s", running.Code, running.Body.String())
	}
	busy := apply(second, nil)
	assertRecordingPlanError(t, busy, http.StatusConflict, "RECORDING_SCHEMA_OPERATION_BUSY", true)
	if got := decodeRecordingOperation(t, busy).Error.OperationID; got != op.OperationID {
		t.Fatalf("busy response must name the occupying operation, got %q", got)
	}

	if _, err := fixture.planSvc.FinishSchemaOperation(context.Background(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationSucceeded, ExecutedStatements: 9,
	}); err != nil {
		t.Fatalf("finish operation: %v", err)
	}
	done := apply(first, nil)
	if got := decodeRecordingOperation(t, done).Data; done.Code != http.StatusOK || got.Status != recordingplan.SchemaOperationSucceeded || got.ExecutedStatements != 9 {
		t.Fatalf("a terminal duplicate must return 200 with the retained result: %d %s", done.Code, done.Body.String())
	}
	if tables := recordingTargetTables(t, openRecordingTarget(t, fixture.targetPath)); len(tables) != 0 {
		t.Fatalf("replaying a confirmation must not touch the target: %v", tables)
	}
}
