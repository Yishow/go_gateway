package recordingplan

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func testWriteScope() TestWriteScope {
	return TestWriteScope{
		WorkspaceID: "ws-1", WorkspaceRevision: "setup-1", GroupID: "group-1", GroupRevision: "rev-1",
		ConnectorID: "db-1", ConnectorRevision: "identity-1", Dialect: "sqlite", Database: "/tmp/line.db",
		Schema: "main", Table: "readings",
	}
}

func newTestWriteService(t *testing.T) *Service {
	t.Helper()
	return NewService(NewMemoryRepository())
}

func TestWritePreviewIsBoundToActionScopeAndContent(t *testing.T) {
	svc := newTestWriteService(t)
	token, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	if token.Action != TestWriteAction || token.OperationID == "" || token.Token == "" || token.PlanID != "group-1" || token.TablePrefix != "readings" {
		t.Fatalf("a test-write token binds its action, operation, group and table: %+v", token)
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", token.Token, testWriteScope(), "content-1"); err != nil {
		t.Fatalf("an unchanged scope validates: %v", err)
	}

	for name, mutate := range map[string]func(*TestWriteScope){
		"workspace revision": func(s *TestWriteScope) { s.WorkspaceRevision = "setup-2" },
		"group revision":     func(s *TestWriteScope) { s.GroupRevision = "rev-2" },
		"connector revision": func(s *TestWriteScope) { s.ConnectorRevision = "identity-2" },
		"table":              func(s *TestWriteScope) { s.Table = "other" },
		"group":              func(s *TestWriteScope) { s.GroupID = "group-2" },
	} {
		current := testWriteScope()
		mutate(&current)
		if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", token.Token, current, "content-1"); !errors.Is(err, ErrPreviewTokenStale) {
			t.Fatalf("%s change must make the token stale, got %v", name, err)
		}
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", token.Token, testWriteScope(), "content-2"); !errors.Is(err, ErrPreviewTokenStale) {
		t.Fatalf("changed preview content must make the token stale, got %v", err)
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-other", token.Token, testWriteScope(), "content-1"); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("a foreign workspace sees an unknown token, got %v", err)
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", "tok-unknown", testWriteScope(), "content-1"); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("an unknown token is not found, got %v", err)
	}
}

func TestWriteTokensNeverCrossKinds(t *testing.T) {
	svc := newTestWriteService(t)
	schemaToken := &SchemaPreviewToken{
		Token: "tok-schema", OperationID: "op-schema", Action: SchemaApplyAction, WorkspaceID: "ws-1", WorkspaceRevision: "setup-1",
		PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", ConnectorRevision: "identity-1", Dialect: "sqlite",
		Database: "/tmp/line.db", Schema: "main", TablePrefix: "gw_record_", Statements: []string{"CREATE TABLE x (a TEXT)"},
		ExpiresAt: time.Now().Add(time.Minute), CreatedAt: time.Now(),
	}
	digest, err := previewDigest(schemaToken)
	if err != nil {
		t.Fatal(err)
	}
	schemaToken.Digest = digest
	if err := svc.repo.SavePreviewToken(t.Context(), schemaToken); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", "tok-schema", testWriteScope(), "content-1"); !errors.Is(err, ErrPreviewTokenKind) {
		t.Fatalf("a schema token cannot confirm a test write, got %v", err)
	}
	if _, _, err := svc.ResolveTestWriteReplay(t.Context(), "ws-1", "tok-schema", "op-schema"); !errors.Is(err, ErrPreviewTokenKind) {
		t.Fatalf("a schema token cannot even replay as a test write, got %v", err)
	}

	testToken, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(t.Context(), testToken.Token, SchemaPreviewScope{WorkspaceID: "ws-1"}); err == nil {
		t.Fatal("a test-write token must never confirm a schema apply")
	}
}

func TestWriteExpiredTokenIsRefusedBeforeClaimButDoesNotErasePersistedResult(t *testing.T) {
	svc := newTestWriteService(t)
	token, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := svc.repo.GetPreviewToken(t.Context(), token.Token)
	if err != nil {
		t.Fatal(err)
	}
	stored.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	if err := svc.repo.SavePreviewToken(t.Context(), stored); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidateTestWriteToken(t.Context(), "ws-1", token.Token, testWriteScope(), "content-1"); !errors.Is(err, ErrPreviewTokenExpired) {
		t.Fatalf("an unclaimed expired token is refused, got %v", err)
	}
	if _, outcome, err := svc.ResolveTestWriteReplay(t.Context(), "ws-1", token.Token, token.OperationID); err != nil || outcome != ClaimNone {
		t.Fatalf("an expired token with no operation has nothing to replay: outcome=%q err=%v", outcome, err)
	}
}

func TestWriteClaimRunsOnceAndReplaysRunningAndSavedResults(t *testing.T) {
	svc := newTestWriteService(t)
	token, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	op, outcome, err := svc.ClaimTestWrite(t.Context(), token)
	if err != nil || outcome != ClaimAcquired || op.Action != TestWriteAction || op.PayloadDigest != "content-1" {
		t.Fatalf("first claim: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, outcome, err := svc.ClaimTestWrite(t.Context(), token); err != nil || outcome != ClaimInProgress {
		t.Fatalf("a duplicate while running is in progress (202), got outcome=%q err=%v", outcome, err)
	}
	replay, outcome, err := svc.ResolveTestWriteReplay(t.Context(), "ws-1", token.Token, token.OperationID)
	if err != nil || outcome != ClaimInProgress || replay.Owner != "" {
		t.Fatalf("replay of a running operation: outcome=%q err=%v owner=%q", outcome, err, replay.Owner)
	}

	// A second preview of the same table is a different operation on the same scope.
	other, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	if busy, outcome, err := svc.ClaimTestWrite(t.Context(), other); err != nil || outcome != ClaimScopeBusy || busy.OperationID != token.OperationID {
		t.Fatalf("another operation on the same table is busy and names the holder: busy=%+v outcome=%q err=%v", busy, outcome, err)
	}

	finished, err := svc.FinishSchemaOperation(t.Context(), token.OperationID, op.Owner, SchemaOperationResult{
		Status: SchemaOperationSucceeded, WriteOutcome: "written_verified", CleanupStatus: "cleaned", PayloadDigest: "content-1",
	})
	if err != nil || finished.WriteOutcome != "written_verified" {
		t.Fatalf("finish: %+v err=%v", finished, err)
	}
	saved, outcome, err := svc.ResolveTestWriteReplay(t.Context(), "ws-1", token.Token, token.OperationID)
	if err != nil || outcome != ClaimCompleted || saved.CleanupStatus != "cleaned" {
		t.Fatalf("a retained result is replayed (200): outcome=%q err=%v op=%+v", outcome, err, saved)
	}
	if _, outcome, err := svc.ClaimTestWrite(t.Context(), other); err != nil || outcome != ClaimAcquired {
		t.Fatalf("the scope is free once the first operation ended: outcome=%q err=%v", outcome, err)
	}
	if _, _, err := svc.ResolveTestWriteReplay(t.Context(), "ws-1", token.Token, "op-other"); !errors.Is(err, ErrSchemaOperationMismatch) {
		t.Fatalf("an operation id that is not the token's is a mismatch, got %v", err)
	}
}

func TestWriteTakeOverOnlyAfterTheLeaseAndOnlyOnce(t *testing.T) {
	svc := newTestWriteService(t)
	token, err := svc.PrepareTestWritePreview(t.Context(), testWriteScope(), "content-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := svc.ClaimTestWrite(t.Context(), token); err != nil {
		t.Fatal(err)
	}
	if adopted, err := svc.TakeOverTestWrite(t.Context(), token.OperationID); err != nil || adopted != nil {
		t.Fatalf("a live lease is never taken over: op=%v err=%v", adopted, err)
	}
	repo := svc.repo.(*MemoryRepository)
	repo.mu.Lock()
	op := repo.operations[token.OperationID]
	op.UpdatedAt = time.Now().UTC().Add(-2 * TestWriteLease)
	repo.operations[token.OperationID] = op
	repo.mu.Unlock()

	adopted, err := svc.TakeOverTestWrite(t.Context(), token.OperationID)
	if err != nil || adopted == nil || adopted.Owner == "" || strings.Contains(adopted.Owner, "claim-") && adopted.Owner == op.Owner {
		t.Fatalf("an expired lease is adopted by a new owner: op=%+v err=%v", adopted, err)
	}
	if again, err := svc.TakeOverTestWrite(t.Context(), token.OperationID); err != nil || again != nil {
		t.Fatalf("the adopted lease is live again, so nobody else may take it: op=%v err=%v", again, err)
	}
}
