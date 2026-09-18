package recordingplan

import (
	"context"
	"errors"
	"testing"
	"time"
)

// ageLastSchemaOperation backdates the operation so its claim lease has run out.
func ageLastSchemaOperation(t *testing.T, repo *MemoryRepository, operationID string) {
	t.Helper()
	repo.mu.Lock()
	defer repo.mu.Unlock()
	op, ok := repo.operations[operationID]
	if !ok {
		t.Fatalf("operation %s is not stored", operationID)
	}
	op.UpdatedAt = time.Now().UTC().Add(-2 * SchemaOperationLease)
	repo.operations[operationID] = op
}

func TestApplySchemaPreview_StaleRunningResolvesFromTargetEvidence(t *testing.T) {
	repo := &failingFinishRepository{MemoryRepository: NewMemoryRepository()}
	svc := NewService(repo)
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)

	// The statements reached the target, but the result could not be recorded:
	// the operation stays running with all managed tables in place.
	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target)
	if !errors.Is(err, ErrSchemaOperationUnacknowledged) || op == nil {
		t.Fatalf("precondition: unrecorded result must stay unresolved: op=%+v err=%v", op, err)
	}
	if got := managedTableCount(t, applyTarget.db); got != len(managedTableOrder) {
		t.Fatalf("precondition: managed tables must exist on the target, got %d", got)
	}
	ageLastSchemaOperation(t, repo.MemoryRepository, token.OperationID)

	resolved, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target)

	if err != nil || outcome != ClaimCompleted || resolved.Status != SchemaOperationSucceeded {
		t.Fatalf("a stale running operation must resolve from the target: op=%+v outcome=%q err=%v", resolved, outcome, err)
	}
	if resolved.Reason != SchemaReasonRecovered || !hexDigest.MatchString(resolved.VerifiedDigest) || applyTarget.executions != 1 {
		t.Fatalf("recovery must use target evidence without re-executing: op=%+v executions=%d", resolved, applyTarget.executions)
	}
	if _, err := svc.repo.FindActiveSchemaOperation(t.Context(), token.WorkspaceID, schemaOperationScopeKey(token)); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("a recovered operation must release its scope, got %v", err)
	}
}

func TestApplySchemaPreview_StaleRunningWithoutEvidenceStaysUnknown(t *testing.T) {
	repo := &failingFinishRepository{MemoryRepository: NewMemoryRepository()}
	svc := NewService(repo)
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)
	if _, _, err := svc.ApplySchemaPreview(t.Context(), token, target); !errors.Is(err, ErrSchemaOperationUnacknowledged) {
		t.Fatalf("precondition: unrecorded result must stay unresolved, got %v", err)
	}
	ageLastSchemaOperation(t, repo.MemoryRepository, token.OperationID)
	// The target cannot be checked any more, so no evidence can resolve the run.
	broken := applyTarget.target()
	broken.Inspect = func(context.Context, string) (TargetTableInspection, error) {
		return TargetTableInspection{}, errors.New("connection lost")
	}

	resolved, outcome, err := svc.ApplySchemaPreview(t.Context(), token, broken)

	if err != nil || outcome != ClaimCompleted || resolved.Status != SchemaOperationUnknown || resolved.Reason != SchemaReasonUnverifiable {
		t.Fatalf("an uncheckable stale operation must resolve to unknown: op=%+v outcome=%q err=%v", resolved, outcome, err)
	}
	if applyTarget.executions != 1 {
		t.Fatalf("recovery must never re-execute: executions=%d", applyTarget.executions)
	}
}
