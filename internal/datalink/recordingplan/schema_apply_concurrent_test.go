package recordingplan

import (
	"testing"
)

func TestApplySchemaPreview_ConcurrentDuplicateGetsRunningWhenTargetChangesDuringInspection(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	token := previewFor(t, svc, applyTarget.target())

	// Simulate first request acquiring claim and starting execution:
	firstOp, outcome, err := svc.ClaimSchemaApply(t.Context(), token)
	if err != nil || outcome != ClaimAcquired {
		t.Fatalf("first apply claim must acquire: %v %v", outcome, err)
	}

	// First request executes DDL on the target, so target now contains a table:
	statements, err := GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := applyTarget.db.ExecContext(t.Context(), statements[0]); err != nil {
		t.Fatal(err)
	}

	// Concurrent duplicate request arrives while first request is running:
	// Even though target changed (some tables now exist), the concurrent duplicate
	// refers to the running operation and gets ClaimInProgress (202):
	duplicateOp, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())
	if err != nil {
		t.Fatalf("concurrent duplicate must not fail with stale error: %v", err)
	}
	if outcome != ClaimInProgress || duplicateOp.OperationID != firstOp.OperationID || duplicateOp.Status != SchemaOperationRunning {
		t.Fatalf("concurrent duplicate must return running operation: outcome=%v op=%+v", outcome, duplicateOp)
	}
}

func TestApplySchemaPreview_ConcurrentDuplicateGetsCompletedWhenFirstFinished(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	token := previewFor(t, svc, applyTarget.target())

	// First apply completes successfully:
	op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())
	if err != nil || outcome != ClaimAcquired || op.Status != SchemaOperationSucceeded {
		t.Fatalf("first apply must succeed: %v %v", outcome, err)
	}

	// Duplicate request arriving after target has all tables returns ClaimCompleted (200) with retained result:
	dup, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())
	if err != nil {
		t.Fatalf("duplicate must not fail: %v", err)
	}
	if outcome != ClaimCompleted || dup.OperationID != op.OperationID || dup.Status != SchemaOperationSucceeded {
		t.Fatalf("duplicate must return completed operation: outcome=%v op=%+v", outcome, dup)
	}
}
