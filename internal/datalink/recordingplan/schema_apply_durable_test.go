package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

// failingFinishOnce lets the first result go unrecorded, as if the local
// database were unavailable exactly when the outcome had to be stored.
type failingFinishOnce struct {
	Repository
	failed bool
}

func (r *failingFinishOnce) FinishSchemaOperation(ctx context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error) {
	if !r.failed {
		r.failed = true
		return nil, errors.New("local database unavailable")
	}
	return r.Repository.FinishSchemaOperation(ctx, operationID, owner, result)
}

func ledgerOwner(t *testing.T, db *sql.DB, operationID string) string {
	t.Helper()
	var owner string
	if err := db.QueryRowContext(t.Context(), `SELECT owner FROM managed_schema_operations WHERE operation_id = ?`, operationID).Scan(&owner); err != nil {
		t.Fatalf("read operation owner: %v", err)
	}
	return owner
}

// The target may already be changed when the local acknowledgement is lost.
// A restart must not run the batch again; it reports the unresolved operation
// until its owner records the verified outcome.
func TestApplySchemaPreview_RestartDoesNotReplayAnUnacknowledgedResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	ledger := openLedgerDB(t, path)
	applyRecordingDDL(t, ledger)
	svc := NewService(&failingFinishOnce{Repository: NewSQLRepository(ledger)})
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	ctx := t.Context()

	token, err := svc.PrepareSchemaPreview(ctx, previewScope(), applyTarget.target().Inspect)
	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	op, _, err := svc.ApplySchemaPreview(ctx, token, applyTarget.target())
	if !errors.Is(err, ErrSchemaOperationUnacknowledged) || op == nil {
		t.Fatalf("a lost acknowledgement must leave the operation unresolved: op=%+v err=%v", op, err)
	}
	if got := managedTableCount(t, applyTarget.db); got != len(managedTableOrder) {
		t.Fatalf("the target was changed before the acknowledgement was lost, got %d tables", got)
	}

	restarted := NewService(NewSQLRepository(openLedgerDB(t, path)))
	again, outcome, err := restarted.ApplySchemaPreview(ctx, token, applyTarget.target())
	if err != nil || outcome != ClaimInProgress || again.OperationID != op.OperationID {
		t.Fatalf("after a restart the unresolved operation must be reported, not replayed: op=%+v outcome=%q err=%v", again, outcome, err)
	}
	if applyTarget.executions != 1 {
		t.Fatalf("a restart must not execute the batch again, got %d executions", applyTarget.executions)
	}

	owner := ledgerOwner(t, ledger, op.OperationID)
	if _, err := restarted.FinishSchemaOperation(ctx, op.OperationID, owner, SchemaOperationResult{
		Status: SchemaOperationSucceeded, ExecutedStatements: len(token.Statements), VerifiedDigest: token.Digest,
	}); err != nil {
		t.Fatalf("record the verified outcome after recovery: %v", err)
	}
	retained, outcome, err := restarted.ApplySchemaPreview(ctx, token, applyTarget.target())
	if err != nil || outcome != ClaimCompleted || retained.Status != SchemaOperationSucceeded || applyTarget.executions != 1 {
		t.Fatalf("a later confirmation must return the retained result: op=%+v outcome=%q err=%v executions=%d",
			retained, outcome, err, applyTarget.executions)
	}
}
