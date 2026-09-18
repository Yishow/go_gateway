package recordingplan

import (
	"errors"
	"testing"
)

func preparedPreview(t *testing.T, svc *Service, scope SchemaPreviewScope) *SchemaPreviewToken {
	t.Helper()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, (&inspectorStub{}).inspect)
	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	return token
}

func TestClaimSchemaApply_AcquiresOnceAndRetainsTheResult(t *testing.T) {
	svc, _ := newPreviewService()
	ctx := t.Context()
	token := preparedPreview(t, svc, previewScope())

	op, outcome, err := svc.ClaimSchemaApply(ctx, token)
	if err != nil || outcome != ClaimAcquired || op.OperationID != token.OperationID || op.Status != SchemaOperationRunning || op.Owner == "" {
		t.Fatalf("first confirmation must acquire the issued operation: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	again, outcome, err := svc.ClaimSchemaApply(ctx, token)
	if err != nil || outcome != ClaimInProgress || again.OperationID != op.OperationID || again.Owner != "" {
		t.Fatalf("a running duplicate must get the same operation without execution rights: op=%+v outcome=%q err=%v", again, outcome, err)
	}

	if _, err := svc.FinishSchemaOperation(ctx, op.OperationID, "claim-other", SchemaOperationResult{Status: SchemaOperationSucceeded}); !errors.Is(err, ErrSchemaOperationNotOwned) {
		t.Fatalf("only the claiming execution may finish the operation, got %v", err)
	}
	finished, err := svc.FinishSchemaOperation(ctx, op.OperationID, op.Owner, SchemaOperationResult{
		Status: SchemaOperationPartial, ExecutedStatements: 2, Reason: "statement_failed", NextAction: "preview again to repair",
	})
	if err != nil || finished.Status != SchemaOperationPartial || finished.ExecutedStatements != 2 || finished.CompletedAt == nil {
		t.Fatalf("the owner must record the terminal result: op=%+v err=%v", finished, err)
	}
	retained, outcome, err := svc.ClaimSchemaApply(ctx, token)
	if err != nil || outcome != ClaimCompleted || retained.Status != SchemaOperationPartial || retained.Reason != "statement_failed" {
		t.Fatalf("a terminal duplicate must return the retained result: op=%+v outcome=%q err=%v", retained, outcome, err)
	}
	if _, err := svc.FinishSchemaOperation(ctx, op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded}); !errors.Is(err, ErrSchemaOperationNotOwned) {
		t.Fatalf("a terminal operation must not be rewritten, got %v", err)
	}
}

func TestClaimSchemaApply_AnotherOperationOnTheSameScopeIsBusy(t *testing.T) {
	svc, _ := newPreviewService()
	ctx := t.Context()
	first := preparedPreview(t, svc, previewScope())
	second := preparedPreview(t, svc, previewScope())
	otherScope := previewScope()
	otherScope.TablePrefix = "gw_other_"
	other := preparedPreview(t, svc, otherScope)

	running, _, err := svc.ClaimSchemaApply(ctx, first)
	if err != nil {
		t.Fatalf("claim first operation: %v", err)
	}
	busy, outcome, err := svc.ClaimSchemaApply(ctx, second)
	if err != nil || outcome != ClaimScopeBusy || busy.OperationID != first.OperationID {
		t.Fatalf("another operation must not run on an occupied scope: op=%+v outcome=%q err=%v", busy, outcome, err)
	}
	if _, outcome, err := svc.ClaimSchemaApply(ctx, other); err != nil || outcome != ClaimAcquired {
		t.Fatalf("a different scope must not be blocked: outcome=%q err=%v", outcome, err)
	}
	if _, err := svc.FinishSchemaOperation(ctx, running.OperationID, running.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, ExecutedStatements: 9}); err != nil {
		t.Fatalf("finish first operation: %v", err)
	}
	if op, outcome, err := svc.ClaimSchemaApply(ctx, second); err != nil || outcome != ClaimAcquired || op.OperationID != second.OperationID {
		t.Fatalf("the scope must be free once the operation ends: op=%+v outcome=%q err=%v", op, outcome, err)
	}
}

func TestClaimSchemaApply_RejectsTokensWithoutIssuedOperation(t *testing.T) {
	svc, _ := newPreviewService()
	legacy := &SchemaPreviewToken{Token: "tok-legacy", WorkspaceID: "ws-1", PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", TablePrefix: "gw_record_"}
	if _, _, err := svc.ClaimSchemaApply(t.Context(), legacy); !errors.Is(err, ErrPreviewTokenLegacy) {
		t.Fatalf("a token without an issued operation must require a new preview, got %v", err)
	}
	if _, _, err := svc.ClaimSchemaApply(t.Context(), nil); !errors.Is(err, ErrPreviewTokenLegacy) {
		t.Fatalf("a missing token must be rejected, got %v", err)
	}
}

func TestFinishSchemaOperation_RequiresATerminalResult(t *testing.T) {
	svc, _ := newPreviewService()
	op, _, err := svc.ClaimSchemaApply(t.Context(), preparedPreview(t, svc, previewScope()))
	if err != nil {
		t.Fatalf("claim operation: %v", err)
	}
	for _, result := range []SchemaOperationResult{
		{Status: SchemaOperationRunning}, {Status: SchemaOperationPending}, {Status: "done"},
		{Status: SchemaOperationSucceeded, ExecutedStatements: -1},
	} {
		if _, err := svc.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, result); !errors.Is(err, ErrSchemaOperationResult) {
			t.Fatalf("result %+v must be rejected, got %v", result, err)
		}
	}
}

func TestGetSchemaOperation_IsWorkspaceScopedAndKeepsUnresolvedStates(t *testing.T) {
	svc, _ := newPreviewService()
	ctx := t.Context()
	op, _, err := svc.ClaimSchemaApply(ctx, preparedPreview(t, svc, previewScope()))
	if err != nil {
		t.Fatalf("claim operation: %v", err)
	}

	if got, err := svc.GetSchemaOperation(ctx, "ws-1", op.OperationID); err != nil || got.Status != SchemaOperationRunning || got.Owner != "" {
		t.Fatalf("a running operation must be readable without its owner: op=%+v err=%v", got, err)
	}
	if _, err := svc.GetSchemaOperation(ctx, "ws-2", op.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("another workspace must not see the operation, got %v", err)
	}
	if _, err := svc.GetSchemaOperation(ctx, "ws-1", "op-unknown"); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("an unknown operation must not be found, got %v", err)
	}
	if _, err := svc.FinishSchemaOperation(ctx, op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationUnknown, Reason: "acknowledgement_lost"}); err != nil {
		t.Fatalf("record unknown outcome: %v", err)
	}
	if got, err := svc.GetSchemaOperation(ctx, "ws-1", op.OperationID); err != nil || got.Status != SchemaOperationUnknown || got.Reason != "acknowledgement_lost" {
		t.Fatalf("an unknown outcome must stay readable: op=%+v err=%v", got, err)
	}
}

func TestResolveSchemaApplyReplay_ReportsLedgerStateWithoutExecuting(t *testing.T) {
	svc, repo := newPreviewService()
	ctx := t.Context()
	first := preparedPreview(t, svc, previewScope())
	second := preparedPreview(t, svc, previewScope())

	if op, outcome, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", first.Token, first.OperationID); err != nil || outcome != ClaimNone || op != nil {
		t.Fatalf("an unclaimed preview has no operation yet: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, err := svc.GetSchemaOperation(ctx, "ws-1", first.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("resolving a replay must not create an operation, got %v", err)
	}
	if _, _, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", first.Token, "op-substituted"); !errors.Is(err, ErrSchemaOperationMismatch) {
		t.Fatalf("a substituted operation id must be rejected, got %v", err)
	}
	if _, _, err := svc.ResolveSchemaApplyReplay(ctx, "ws-2", first.Token, first.OperationID); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("a foreign token must look unknown, got %v", err)
	}
	if _, _, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", "tok-unknown", first.OperationID); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("an unknown token must be rejected, got %v", err)
	}

	running, _, err := svc.ClaimSchemaApply(ctx, first)
	if err != nil {
		t.Fatalf("claim first operation: %v", err)
	}
	if op, outcome, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", first.Token, first.OperationID); err != nil || outcome != ClaimInProgress || op.OperationID != first.OperationID {
		t.Fatalf("a running operation must be reported in progress: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if op, outcome, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", second.Token, second.OperationID); err != nil || outcome != ClaimScopeBusy || op.OperationID != first.OperationID {
		t.Fatalf("another preview of the same scope must see the busy operation: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, err := svc.FinishSchemaOperation(ctx, running.OperationID, running.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded}); err != nil {
		t.Fatalf("finish first operation: %v", err)
	}
	if op, outcome, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", first.Token, first.OperationID); err != nil || outcome != ClaimCompleted || op.Status != SchemaOperationSucceeded {
		t.Fatalf("a finished operation must be reported with its result: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, outcome, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", second.Token, second.OperationID); err != nil || outcome != ClaimNone {
		t.Fatalf("the scope must be free after the first operation ends: outcome=%q err=%v", outcome, err)
	}

	legacy := &SchemaPreviewToken{Token: "tok-legacy", WorkspaceID: "ws-1", PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", TablePrefix: "gw_record_"}
	if err := repo.MemoryRepository.SavePreviewToken(ctx, legacy); err != nil {
		t.Fatalf("store legacy token: %v", err)
	}
	if _, _, err := svc.ResolveSchemaApplyReplay(ctx, "ws-1", "tok-legacy", "op-any"); !errors.Is(err, ErrPreviewTokenLegacy) {
		t.Fatalf("a legacy token must require a new preview, got %v", err)
	}
}
