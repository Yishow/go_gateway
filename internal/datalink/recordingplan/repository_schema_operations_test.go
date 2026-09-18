package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// openLedgerDB opens the local SQLite file the way separate gateway processes do.
func openLedgerDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func ledgerOperation(id, token, scope, owner string) *SchemaOperation {
	now := time.Now().UTC()
	return &SchemaOperation{
		OperationID: id, Token: token, WorkspaceID: "ws-1", ScopeKey: scope, Owner: owner,
		Action: SchemaApplyAction, Status: SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	}
}

func TestSQLRepository_ConcurrentProcessesAcquireAnOperationOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	first := openLedgerDB(t, path)
	applyRecordingDDL(t, first)
	repos := []*SQLRepository{NewSQLRepository(first), NewSQLRepository(openLedgerDB(t, path))}

	const attempts = 8
	outcomes := make([]ClaimOutcome, attempts)
	errs := make([]error, attempts)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range attempts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			op := ledgerOperation("op-1", "tok-1", "scope-a", fmt.Sprintf("claim-%d", i))
			_, outcomes[i], errs[i] = repos[i%len(repos)].ClaimSchemaOperation(context.Background(), op)
		}()
	}
	close(start)
	wg.Wait()

	acquired := 0
	for i, outcome := range outcomes {
		if errs[i] != nil {
			t.Fatalf("claim %d failed: %v", i, errs[i])
		}
		switch outcome {
		case ClaimAcquired:
			acquired++
		case ClaimInProgress:
		default:
			t.Fatalf("claim %d got unexpected outcome %q", i, outcome)
		}
	}
	if acquired != 1 {
		t.Fatalf("exactly one process may acquire execution rights, got %d: %v", acquired, outcomes)
	}
}

func TestSQLRepository_RestartDoesNotReacquireOrReplay(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	before := openLedgerDB(t, path)
	applyRecordingDDL(t, before)
	if _, outcome, err := NewSQLRepository(before).ClaimSchemaOperation(t.Context(), ledgerOperation("op-1", "tok-1", "scope-a", "claim-before")); err != nil || outcome != ClaimAcquired {
		t.Fatalf("claim before restart: outcome=%q err=%v", outcome, err)
	}
	if err := before.Close(); err != nil {
		t.Fatalf("close before restart: %v", err)
	}

	restarted := NewSQLRepository(openLedgerDB(t, path))
	ctx := t.Context()
	if op, outcome, err := restarted.ClaimSchemaOperation(ctx, ledgerOperation("op-1", "tok-1", "scope-a", "claim-after")); err != nil || outcome != ClaimInProgress || op.Owner != "" {
		t.Fatalf("after restart the running operation must be returned, not acquired again: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if op, outcome, err := restarted.ClaimSchemaOperation(ctx, ledgerOperation("op-2", "tok-2", "scope-a", "claim-other")); err != nil || outcome != ClaimScopeBusy || op.OperationID != "op-1" {
		t.Fatalf("another operation must stay blocked on the running scope: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, err := restarted.FinishSchemaOperation(ctx, "op-1", "claim-after", SchemaOperationResult{Status: SchemaOperationSucceeded}); !errors.Is(err, ErrSchemaOperationNotOwned) {
		t.Fatalf("a duplicate claimant must not finish the operation, got %v", err)
	}
	finished, err := restarted.FinishSchemaOperation(ctx, "op-1", "claim-before", SchemaOperationResult{
		Status: SchemaOperationUnknown, ExecutedStatements: 1, Reason: "acknowledgement_lost", NextAction: "inspect the target",
	})
	if err != nil || finished.Status != SchemaOperationUnknown || finished.CompletedAt == nil {
		t.Fatalf("the owner must record the unresolved outcome: op=%+v err=%v", finished, err)
	}
	if op, outcome, err := restarted.ClaimSchemaOperation(ctx, ledgerOperation("op-1", "tok-1", "scope-a", "claim-retry")); err != nil ||
		outcome != ClaimCompleted || op.Status != SchemaOperationUnknown || op.ExecutedStatements != 1 {
		t.Fatalf("a retry must get the retained outcome instead of replaying: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if _, outcome, err := restarted.ClaimSchemaOperation(ctx, ledgerOperation("op-2", "tok-2", "scope-a", "claim-other")); err != nil || outcome != ClaimAcquired {
		t.Fatalf("a new confirmed preview may repair the scope: outcome=%q err=%v", outcome, err)
	}
}

func TestSQLRepository_SchemaOperationsAreWorkspaceScoped(t *testing.T) {
	repo := NewSQLRepository(recordingRepositoryDB(t))
	ctx := t.Context()
	if _, outcome, err := repo.ClaimSchemaOperation(ctx, ledgerOperation("op-1", "tok-1", "scope-a", "claim-1")); err != nil || outcome != ClaimAcquired {
		t.Fatalf("claim operation: outcome=%q err=%v", outcome, err)
	}

	got, err := repo.GetSchemaOperation(ctx, "ws-1", "op-1")
	if err != nil || got.Status != SchemaOperationRunning || got.Token != "tok-1" || got.ScopeKey != "scope-a" {
		t.Fatalf("running operation must be readable: op=%+v err=%v", got, err)
	}
	if _, err := repo.GetSchemaOperation(ctx, "ws-2", "op-1"); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("another workspace must not read the operation, got %v", err)
	}
	if _, err := repo.GetSchemaOperation(ctx, "ws-1", "op-missing"); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("an unknown operation must not be found, got %v", err)
	}
	if active, err := repo.FindActiveSchemaOperation(ctx, "ws-1", "scope-a"); err != nil || active.OperationID != "op-1" {
		t.Fatalf("the running operation must hold its scope: op=%+v err=%v", active, err)
	}
	if _, err := repo.FinishSchemaOperation(ctx, "op-1", "claim-1", SchemaOperationResult{
		Status: SchemaOperationPartial, ExecutedStatements: 3, VerifiedDigest: "digest-a", Reason: "statement_failed", NextAction: "preview again to repair",
	}); err != nil {
		t.Fatalf("finish operation: %v", err)
	}
	got, err = repo.GetSchemaOperation(ctx, "ws-1", "op-1")
	if err != nil || got.Status != SchemaOperationPartial || got.ExecutedStatements != 3 || got.VerifiedDigest != "digest-a" ||
		got.Reason != "statement_failed" || got.NextAction != "preview again to repair" || got.CompletedAt == nil {
		t.Fatalf("partial result must be persisted: op=%+v err=%v", got, err)
	}
	if _, err := repo.FindActiveSchemaOperation(ctx, "ws-1", "scope-a"); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("a finished operation must release its scope, got %v", err)
	}
}

func TestSQLRepository_TokenCannotBindASecondOperation(t *testing.T) {
	repo := NewSQLRepository(recordingRepositoryDB(t))
	ctx := t.Context()
	if _, _, err := repo.ClaimSchemaOperation(ctx, ledgerOperation("op-1", "tok-1", "scope-a", "claim-1")); err != nil {
		t.Fatalf("claim operation: %v", err)
	}
	if _, _, err := repo.ClaimSchemaOperation(ctx, ledgerOperation("op-2", "tok-1", "scope-b", "claim-2")); !errors.Is(err, ErrSchemaOperationMismatch) {
		t.Fatalf("a token must keep its issued operation, got %v", err)
	}
}
