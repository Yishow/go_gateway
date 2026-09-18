package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

// applyBody builds a complete confirmation for token. An override replaces a
// field; an empty override removes it.
func (f recordingRouteFixture) applyBody(t *testing.T, token *recordingplan.SchemaPreviewToken, overrides map[string]string) string {
	t.Helper()
	fields := map[string]string{
		"token":                       token.Token,
		"operation_id":                token.OperationID,
		"expected_workspace_revision": f.workspaceRevision,
		"expected_plan_revision":      token.PlanRevision,
		"expected_connector_revision": f.connector.IdentityRevision,
	}
	for key, value := range overrides {
		if value == "" {
			delete(fields, key)
			continue
		}
		fields[key] = value
	}
	body, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("encode apply body: %v", err)
	}
	return string(body)
}

func (f recordingRouteFixture) apply(body string) *httptest.ResponseRecorder {
	return serveRecordingScopedRequest(f.router, http.MethodPost, recordingPlansPath+"/schema-apply", body)
}

func managedRecordingTables(tables []string) int {
	count := 0
	for _, table := range tables {
		if strings.HasPrefix(table, "gw_record_") {
			count++
		}
	}
	return count
}

func TestStudioV2RecordingRoutes_SchemaApplyRejectsUnsafeConfirmationsWithoutMutation(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-apply", "")
	token := previewForRoute(t, fixture, "plan-apply")
	ctx := context.Background()
	save := func(preview *recordingplan.SchemaPreviewToken) {
		t.Helper()
		if err := fixture.planRepo.SavePreviewToken(ctx, preview); err != nil {
			t.Fatalf("store preview: %v", err)
		}
	}

	for _, field := range []string{"token", "operation_id", "expected_workspace_revision", "expected_plan_revision", "expected_connector_revision"} {
		assertRecordingValidationError(t, fixture.apply(fixture.applyBody(t, token, map[string]string{field: ""})))
	}
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, map[string]string{"expected_workspace_revision": "workspace-stale"})),
		http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, map[string]string{"expected_plan_revision": "rev-0"})),
		http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, map[string]string{"expected_connector_revision": "identity-old"})),
		http.StatusConflict, "RECORDING_CONNECTOR_REVISION_CONFLICT", false)

	expired := *token
	expired.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	save(&expired)
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, nil)), http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_EXPIRED", false)
	tampered := *token
	tampered.Statements = []string{"DROP TABLE sentinel"}
	save(&tampered)
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, nil)), http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)
	save(token)

	foreign := *token
	foreign.Token, foreign.OperationID, foreign.WorkspaceID = "tok-foreign", "op-foreign", "workspace-foreign"
	save(&foreign)
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, &foreign, nil)), http.StatusNotFound, "RECORDING_SCHEMA_PREVIEW_NOT_FOUND", false)

	plan, err := fixture.planSvc.GetPlanByWorkspace(ctx, "plan-apply", fixture.workspaceID)
	if err != nil {
		t.Fatalf("read plan: %v", err)
	}
	plan.Name = "Renamed after preview"
	if err := fixture.planSvc.UpdatePlanByWorkspace(ctx, plan, fixture.workspaceID); err != nil {
		t.Fatalf("update plan: %v", err)
	}
	assertRecordingPlanError(t, fixture.apply(fixture.applyBody(t, token, map[string]string{"expected_plan_revision": plan.Revision})),
		http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)

	if tables := recordingTargetTables(t, openRecordingTarget(t, fixture.targetPath)); len(tables) != 0 {
		t.Fatalf("rejected confirmations must not touch the target: %v", tables)
	}
	if _, err := fixture.planSvc.GetSchemaOperation(ctx, fixture.workspaceID, token.OperationID); !errors.Is(err, recordingplan.ErrSchemaOperationNotFound) {
		t.Fatalf("rejected confirmations must not claim an operation, got %v", err)
	}
}

func TestStudioV2RecordingRoutes_SchemaApplyCreatesVerifiedSQLiteSchemaOnce(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-create", "")
	token := previewForRoute(t, fixture, "plan-create")
	body := fixture.applyBody(t, token, nil)

	w := fixture.apply(body)

	created := decodeRecordingOperation(t, w).Data
	if w.Code != http.StatusOK || created.OperationID != token.OperationID || created.Status != recordingplan.SchemaOperationSucceeded ||
		created.ExecutedStatements != len(token.Statements) || len(created.VerifiedDigest) != 64 || created.CompletedAt == nil {
		t.Fatalf("a current confirmation must create and verify the schema: %d %s", w.Code, w.Body.String())
	}
	if got := managedRecordingTables(recordingTargetTables(t, openRecordingTarget(t, fixture.targetPath))); got != 6 {
		t.Fatalf("all managed tables must exist after apply, got %d", got)
	}

	again := fixture.apply(body)
	retained := decodeRecordingOperation(t, again).Data
	if again.Code != http.StatusOK || retained.Status != recordingplan.SchemaOperationSucceeded || retained.CompletedAt == nil ||
		!retained.CompletedAt.Equal(*created.CompletedAt) {
		t.Fatalf("a repeated confirmation must return the retained result: %d %s", again.Code, again.Body.String())
	}
	status := serveRecordingScopedRequest(fixture.router, http.MethodGet, databaseOperationsPath+token.OperationID, "")
	if status.Code != http.StatusOK || decodeRecordingOperation(t, status).Data.Status != recordingplan.SchemaOperationSucceeded {
		t.Fatalf("the finished operation must be readable: %d %s", status.Code, status.Body.String())
	}
}

func TestStudioV2RecordingRoutes_SchemaApplyRejectsTargetChangedAfterPreview(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-changed", "")
	token := previewForRoute(t, fixture, "plan-changed")
	statements, err := recordingplan.GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatalf("generate managed schema: %v", err)
	}
	target := openRecordingTarget(t, fixture.targetPath)
	if _, err := target.ExecContext(context.Background(), statements[0]); err != nil {
		t.Fatalf("change target: %v", err)
	}

	w := fixture.apply(fixture.applyBody(t, token, nil))

	assertRecordingPlanError(t, w, http.StatusConflict, "RECORDING_SCHEMA_PREVIEW_STALE", false)
	if got := managedRecordingTables(recordingTargetTables(t, target)); got != 1 {
		t.Fatalf("a stale preview must not execute any statement, found %d managed tables", got)
	}
	if _, err := fixture.planSvc.GetSchemaOperation(context.Background(), fixture.workspaceID, token.OperationID); !errors.Is(err, recordingplan.ErrSchemaOperationNotFound) {
		t.Fatalf("a stale preview must not claim an operation, got %v", err)
	}
}
