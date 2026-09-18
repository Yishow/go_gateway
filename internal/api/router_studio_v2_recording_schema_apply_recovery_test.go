package api

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

// A confirmation whose operation was claimed but never finished must stop
// answering 202 forever: once the claim lease has run out, a retry resolves the
// operation from the actual target instead of replaying it blindly.
func TestStudioV2RecordingRoutes_SchemaApplyResolvesStaleRunningFromTarget(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-recovery", "")
	token := previewForRoute(t, fixture, "plan-recovery")
	ctx := context.Background()
	stale := time.Now().UTC().Add(-2 * recordingplan.SchemaOperationLease)
	claim := &recordingplan.SchemaOperation{
		OperationID: token.OperationID, Token: token.Token, WorkspaceID: fixture.workspaceID,
		ScopeKey: token.ScopeKey(), Owner: "claim-lost", Action: recordingplan.SchemaApplyAction,
		Status: recordingplan.SchemaOperationRunning, CreatedAt: stale, UpdatedAt: stale,
	}
	if _, outcome, err := fixture.planRepo.ClaimSchemaOperation(ctx, claim); err != nil || outcome != recordingplan.ClaimAcquired {
		t.Fatalf("precondition: stale claim must be acquirable: outcome=%q err=%v", outcome, err)
	}

	w := fixture.apply(fixture.applyBody(t, token, nil))

	if w.Code != http.StatusOK {
		t.Fatalf("a stale running operation must resolve on retry: %d %s", w.Code, w.Body.String())
	}
	var body struct {
		Success bool `json:"success"`
		Data    struct {
			OperationID string `json:"operation_id"`
			Status      string `json:"status"`
			Reason      string `json:"reason"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode resolved operation: %v; body=%s", err, w.Body.String())
	}
	if !body.Success || body.Data.OperationID != token.OperationID {
		t.Fatalf("recovery must return the same operation: %s", w.Body.String())
	}
	// The lost claim never executed, so the tables are missing and the honest
	// answer is unknown with a repair action, never a silent re-run.
	if body.Data.Status != "unknown" || body.Data.Reason == "" {
		t.Fatalf("a lost claim without evidence must resolve to unknown: %s", w.Body.String())
	}
	if _, err := fixture.planSvc.GetSchemaOperation(ctx, fixture.workspaceID, token.OperationID); err != nil {
		t.Fatalf("recovered operation must stay queryable: %v", err)
	}
}

// A fresh running claim still holds its scope: a retry keeps answering 202
// while the owning execution may still be alive.
func TestStudioV2RecordingRoutes_SchemaApplyKeepsFreshRunningClaim(t *testing.T) {
	fixture := newRecordingRouteFixture(t)
	seedRecordingRoutePlan(t, fixture, "plan-fresh-claim", "")
	token := previewForRoute(t, fixture, "plan-fresh-claim")
	claim := &recordingplan.SchemaOperation{
		OperationID: token.OperationID, Token: token.Token, WorkspaceID: fixture.workspaceID,
		ScopeKey: token.ScopeKey(), Owner: "claim-live", Action: recordingplan.SchemaApplyAction,
		Status:    recordingplan.SchemaOperationRunning,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if _, outcome, err := fixture.planRepo.ClaimSchemaOperation(context.Background(), claim); err != nil || outcome != recordingplan.ClaimAcquired {
		t.Fatalf("precondition: fresh claim must be acquirable: outcome=%q err=%v", outcome, err)
	}

	w := fixture.apply(fixture.applyBody(t, token, nil))

	if w.Code != http.StatusAccepted {
		t.Fatalf("a fresh running claim must keep returning 202: %d %s", w.Code, w.Body.String())
	}
}
