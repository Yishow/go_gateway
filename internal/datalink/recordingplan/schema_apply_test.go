package recordingplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func openApplyTarget(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "target.db"))
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func sqliteTableState(ctx context.Context, db *sql.DB, table string) (TargetTableInspection, error) {
	rows, err := db.QueryContext(ctx, `SELECT name FROM pragma_table_info(?)`, table)
	if err != nil {
		return TargetTableInspection{}, err
	}
	defer func() { _ = rows.Close() }()
	var columns []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return TargetTableInspection{}, err
		}
		columns = append(columns, name)
	}
	if err := rows.Err(); err != nil {
		return TargetTableInspection{}, err
	}
	if len(columns) == 0 {
		return TargetTableInspection{Status: "missing"}, nil
	}
	return TargetTableInspection{Status: "exists", Columns: columns}, nil
}

// recordingApplyTarget inspects and changes a real SQLite target and counts executions.
type recordingApplyTarget struct {
	db         *sql.DB
	executions int
}

func (r *recordingApplyTarget) target() SchemaApplyTarget {
	return SchemaApplyTarget{
		Inspect: func(ctx context.Context, table string) (TargetTableInspection, error) {
			return sqliteTableState(ctx, r.db, table)
		},
		Execute: func(ctx context.Context, statements []string) (SchemaExecution, error) {
			r.executions++
			tx, err := r.db.BeginTx(ctx, nil)
			if err != nil {
				return SchemaExecution{RolledBack: true}, err
			}
			for _, statement := range statements {
				if _, err := tx.ExecContext(ctx, statement); err != nil {
					_ = tx.Rollback()
					return SchemaExecution{RolledBack: true}, err
				}
			}
			if err := tx.Commit(); err != nil {
				return SchemaExecution{}, err
			}
			return SchemaExecution{Committed: len(statements)}, nil
		},
	}
}

func managedTableCount(t *testing.T, db *sql.DB) int {
	t.Helper()
	var count int
	if err := db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name LIKE 'gw_record_%'`).Scan(&count); err != nil {
		t.Fatalf("count managed tables: %v", err)
	}
	return count
}

func previewFor(t *testing.T, svc *Service, target SchemaApplyTarget) *SchemaPreviewToken {
	t.Helper()
	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), target.Inspect)
	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	return token
}

func TestApplySchemaPreview_CreatesAndVerifiesTheManagedSchemaOnce(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	token := previewFor(t, svc, applyTarget.target())

	op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())

	if err != nil || outcome != ClaimAcquired || op.Status != SchemaOperationSucceeded || op.ExecutedStatements != len(token.Statements) {
		t.Fatalf("legitimate apply must succeed: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	if !hexDigest.MatchString(op.VerifiedDigest) || op.CompletedAt == nil || op.Owner != "" {
		t.Fatalf("success must carry the verified schema digest: %+v", op)
	}
	if got := managedTableCount(t, applyTarget.db); got != len(managedTableOrder) {
		t.Fatalf("all managed tables must exist after apply, got %d", got)
	}
	again, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())
	if err != nil || outcome != ClaimCompleted || again.Status != SchemaOperationSucceeded || applyTarget.executions != 1 {
		t.Fatalf("a repeated confirmation must return the retained result without executing: op=%+v outcome=%q err=%v executions=%d",
			again, outcome, err, applyTarget.executions)
	}
}

func TestApplySchemaPreview_CompatibleSchemaIsAVerifiedNoOp(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	statements, err := GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatal(err)
	}
	for _, statement := range statements {
		if _, err := applyTarget.db.ExecContext(t.Context(), statement); err != nil {
			t.Fatal(err)
		}
	}
	token := previewFor(t, svc, applyTarget.target())

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())

	if err != nil || op.Status != SchemaOperationSucceeded || op.ExecutedStatements != 0 || !hexDigest.MatchString(op.VerifiedDigest) || applyTarget.executions != 0 {
		t.Fatalf("a compatible schema must succeed only after verification and without execution: op=%+v err=%v executions=%d", op, err, applyTarget.executions)
	}
}

func TestApplySchemaPreview_RejectsTargetChangedAfterPreviewBeforeClaiming(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	token := previewFor(t, svc, applyTarget.target())
	statements, err := GenerateManagedSchemaDDL("sqlite", "gw_record_")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := applyTarget.db.ExecContext(t.Context(), statements[0]); err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())

	if !errors.Is(err, ErrPreviewTokenStale) || !errors.Is(err, ErrSchemaTargetChanged) || applyTarget.executions != 0 {
		t.Fatalf("a changed target must make the preview stale before execution: err=%v executions=%d", err, applyTarget.executions)
	}
	if _, err := svc.GetSchemaOperation(t.Context(), "ws-1", token.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("a stale preview must not claim an operation, got %v", err)
	}
}

func TestApplySchemaPreview_RejectsExpiredLegacyAndUnverifiedAdapters(t *testing.T) {
	svc, repo := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	ctx := t.Context()

	expired := *previewFor(t, svc, applyTarget.target())
	expired.ExpiresAt = expired.CreatedAt.Add(-1)
	if _, _, err := svc.ApplySchemaPreview(ctx, &expired, applyTarget.target()); !errors.Is(err, ErrPreviewTokenExpired) {
		t.Fatalf("an expired preview must be rejected, got %v", err)
	}
	legacy := &SchemaPreviewToken{Token: "tok-legacy", WorkspaceID: "ws-1", PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", TablePrefix: "gw_record_"}
	if _, _, err := svc.ApplySchemaPreview(ctx, legacy, applyTarget.target()); !errors.Is(err, ErrPreviewTokenLegacy) {
		t.Fatalf("a legacy preview must be rejected, got %v", err)
	}
	unverified := &SchemaPreviewToken{
		Token: "tok-mysql", OperationID: "op-mysql", Action: SchemaApplyAction, WorkspaceID: "ws-1", WorkspaceRevision: "setup-1",
		PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", ConnectorRevision: "identity-1", Dialect: "mysql",
		Database: "metrics", Schema: "metrics", TablePrefix: "gw_record_", Digest: strings.Repeat("a", 64),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	if _, _, err := svc.ApplySchemaPreview(ctx, unverified, applyTarget.target()); !errors.Is(err, ErrSchemaAdapterUnverified) {
		t.Fatalf("an unverified adapter must stay unavailable, got %v", err)
	}
	if applyTarget.executions != 0 || repo.saves != 1 {
		t.Fatalf("rejected applies must not execute: executions=%d saves=%d", applyTarget.executions, repo.saves)
	}
	if _, err := svc.GetSchemaOperation(ctx, "ws-1", unverified.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("an unverified adapter must not claim an operation, got %v", err)
	}
}

func TestApplySchemaPreview_RolledBackFailureIsVerifiedFailed(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)
	target.Execute = func(ctx context.Context, statements []string) (SchemaExecution, error) {
		tx, err := applyTarget.db.BeginTx(ctx, nil)
		if err != nil {
			return SchemaExecution{RolledBack: true}, err
		}
		if _, err := tx.ExecContext(ctx, statements[0]); err != nil {
			t.Fatalf("execute first statement: %v", err)
		}
		_ = tx.Rollback()
		return SchemaExecution{RolledBack: true}, errors.New("disk full")
	}

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target)

	if err != nil || op.Status != SchemaOperationFailed || op.ExecutedStatements != 0 || op.Reason != SchemaReasonRolledBack || op.NextAction == "" {
		t.Fatalf("a rolled back batch must be a verified failure with a next action: op=%+v err=%v", op, err)
	}
	if got := managedTableCount(t, applyTarget.db); got != 0 {
		t.Fatalf("a rolled back batch must leave no managed table, got %d", got)
	}
}

func TestApplySchemaPreview_PermissionDeniedIsAVerifiedFailure(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)
	target.Execute = func(context.Context, []string) (SchemaExecution, error) {
		return SchemaExecution{RolledBack: true}, fmt.Errorf("%w: create denied", ErrTargetPermissionDenied)
	}

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target)

	if err != nil || op.Status != SchemaOperationFailed || op.Reason != SchemaReasonPermissionDenied || op.NextAction == "" {
		t.Fatalf("a refused batch must be a verified permission failure: op=%+v err=%v", op, err)
	}
	if got := managedTableCount(t, applyTarget.db); got != 0 {
		t.Fatalf("a refused batch must leave no managed table, got %d", got)
	}
}

func TestApplySchemaPreview_PartialFailureIsNotReportedAsRollback(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)
	target.Execute = func(ctx context.Context, statements []string) (SchemaExecution, error) {
		for _, statement := range statements[:3] {
			if _, err := applyTarget.db.ExecContext(ctx, statement); err != nil {
				t.Fatalf("execute statement: %v", err)
			}
		}
		return SchemaExecution{Committed: 3}, errors.New("connection reset")
	}

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target)

	if err != nil || op.Status != SchemaOperationPartial || op.ExecutedStatements != 3 || op.Reason != SchemaReasonPartiallyApplied || op.NextAction == "" {
		t.Fatalf("a partly applied batch must be reported partial with a repair action: op=%+v err=%v", op, err)
	}
	if got := managedTableCount(t, applyTarget.db); got != 1 {
		t.Fatalf("partial apply must leave exactly the committed table, got %d", got)
	}
}

func TestApplySchemaPreview_UnverifiableResultIsUnknown(t *testing.T) {
	svc, _ := newPreviewService()
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	target := applyTarget.target()
	token := previewFor(t, svc, target)
	inspect := target.Inspect
	target.Inspect = func(ctx context.Context, table string) (TargetTableInspection, error) {
		if applyTarget.executions > 0 {
			return TargetTableInspection{}, errors.New("connection lost")
		}
		return inspect(ctx, table)
	}

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target)

	if err != nil || op.Status != SchemaOperationUnknown || op.Reason != SchemaReasonUnverifiable || op.NextAction == "" {
		t.Fatalf("an unverifiable result must stay unknown: op=%+v err=%v", op, err)
	}
}

type failingFinishRepository struct {
	*MemoryRepository
}

func (r failingFinishRepository) FinishSchemaOperation(context.Context, string, string, SchemaOperationResult) (*SchemaOperation, error) {
	return nil, errors.New("local database unavailable")
}

func TestApplySchemaPreview_UnrecordedResultStaysUnresolvedWithoutReplay(t *testing.T) {
	svc := NewService(failingFinishRepository{MemoryRepository: NewMemoryRepository()})
	applyTarget := &recordingApplyTarget{db: openApplyTarget(t)}
	token := previewFor(t, svc, applyTarget.target())

	op, _, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())

	if !errors.Is(err, ErrSchemaOperationUnacknowledged) || op == nil || op.OperationID != token.OperationID || op.Owner != "" {
		t.Fatalf("an unrecorded result must name the unresolved operation: op=%+v err=%v", op, err)
	}
	again, outcome, err := svc.ApplySchemaPreview(t.Context(), token, applyTarget.target())
	if err != nil || outcome != ClaimInProgress || again.Status != SchemaOperationRunning || applyTarget.executions != 1 {
		t.Fatalf("a retry must not replay an unresolved operation: op=%+v outcome=%q err=%v executions=%d", again, outcome, err, applyTarget.executions)
	}
}
