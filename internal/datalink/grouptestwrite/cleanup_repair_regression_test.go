package grouptestwrite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func setupSQLServerCleanupTest(t *testing.T, denyDelete, ambiguousCommit bool) (*Service, *sqlServerRunScript, *Preview, Workspace, string) {
	group, layout, tagTypes := sqlServerTestWriteFixture(t)
	columns := make([]string, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		columns = append(columns, column.Name)
	}
	script := &sqlServerRunScript{
		sqlServerExecutionScript: &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}},
		denyDelete:               denyDelete,
		ambiguousCommit:          ambiguousCommit,
	}
	db := sql.OpenDB(sqlServerRunConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })

	groups := &sqlServerTestWriteGroups{
		fakeGroups: &fakeGroups{group: group}, layout: layout, testOnlySupportsGroupTestWrite: true,
	}
	svc := New(Dependencies{
		Groups: groups, Tags: fakeTags(tagTypes), Destinations: sqlServerRunDestinations{db: db},
		Ledger: recordingplan.NewService(recordingplan.NewMemoryRepository()),
	}, Config{StepTimeout: time.Second})
	svc.insertRow = func(_ context.Context, _ *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		script.inserts++
		digest, err := groupdelivery.RowPayloadDigest(req.Row)
		require.NoError(t, err)
		script.receipts[req.Row.EffectKey] = digest
		values := make([]driver.Value, len(columns))
		for i, column := range columns {
			for _, cell := range req.Row.Cells {
				if cell.Column == column {
					values[i] = cell.Value
				}
			}
		}
		script.rows = append(script.rows, values)
		return dbtarget.GroupInsertResult{}, nil
	}

	ws := Workspace{ID: group.WorkspaceID, Revision: testRevision}
	preview, err := svc.Preview(t.Context(), ws, group.ID)
	require.NoError(t, err)
	return svc, script, preview, ws, group.ID
}

func TestSQLServerCleanupRepairBehavioralRegression(t *testing.T) {
	svc, script, preview, ws, groupID := setupSQLServerCleanupTest(t, true, false)

	// Step 1: Initial confirm runs write + readback (succeeds) and cleanup (fails due to denyDelete).
	initialReq := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: groupID}
	initialOutcome, err := svc.Confirm(t.Context(), ws, initialReq)
	require.NoError(t, err)
	require.Equal(t, WriteVerified, initialOutcome.Operation.WriteOutcome)
	require.Equal(t, CleanupFailed, initialOutcome.Operation.CleanupStatus)
	require.NotEmpty(t, script.rows)
	initialInserts := script.inserts

	// Normal Confirm replay without action must just return retained result and not clean.
	deletesBeforeReplay := script.deletes
	replayOutcome, err := svc.Confirm(t.Context(), ws, initialReq)
	require.NoError(t, err)
	require.Equal(t, CleanupFailed, replayOutcome.Operation.CleanupStatus)
	require.Equal(t, deletesBeforeReplay, script.deletes)
	require.Equal(t, initialInserts, script.inserts)

	// Step 2: Resolve delete permission, then request cleanup repair.
	script.denyDelete = false
	repairReq := Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "cleanup_repair",
	}
	repairOutcome, err := svc.Confirm(t.Context(), ws, repairReq)
	require.NoError(t, err)
	require.Equal(t, CleanupCleaned, repairOutcome.Operation.CleanupStatus)
	require.Equal(t, WriteVerified, repairOutcome.Operation.WriteOutcome)
	require.Empty(t, script.rows, "target test rows must be cleaned on repair")
	require.Empty(t, script.receipts, "receipts must be cleaned on repair")
	require.Equal(t, initialInserts, script.inserts, "repair must NEVER write again")

	// Step 3: Double repair when already cleaned returns ClaimCompleted with CleanupCleaned.
	deletesAfterClean := script.deletes
	doubleOutcome, err := svc.Confirm(t.Context(), ws, repairReq)
	require.NoError(t, err)
	require.Equal(t, recordingplan.ClaimCompleted, doubleOutcome.Claim)
	require.Equal(t, CleanupCleaned, doubleOutcome.Operation.CleanupStatus)
	require.Equal(t, deletesAfterClean, script.deletes, "no additional deletes on double retry")
	require.Equal(t, initialInserts, script.inserts)
}

func TestSQLServerCleanupRepairAmbiguousUnknownRecovery(t *testing.T) {
	svc, script, preview, ws, groupID := setupSQLServerCleanupTest(t, false, true)

	// Cleanup commit fails ambiguously -> CleanupUnknown.
	initialReq := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: groupID}
	initialOutcome, err := svc.Confirm(t.Context(), ws, initialReq)
	require.NoError(t, err)
	require.Equal(t, WriteVerified, initialOutcome.Operation.WriteOutcome)
	require.Equal(t, CleanupUnknown, initialOutcome.Operation.CleanupStatus)
	initialInserts := script.inserts

	// Remove ambiguity and request cleanup repair.
	script.ambiguousCommit = false
	repairReq := Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "cleanup_repair",
	}
	repairOutcome, err := svc.Confirm(t.Context(), ws, repairReq)
	require.NoError(t, err)
	require.Equal(t, CleanupCleaned, repairOutcome.Operation.CleanupStatus)
	require.Empty(t, script.rows)
	require.Empty(t, script.receipts)
	require.Equal(t, initialInserts, script.inserts, "repair must never write again")
}

func TestSQLServerCleanupRepairRejectsUnknownActionsAndStaleScope(t *testing.T) {
	svc, script, preview, ws, groupID := setupSQLServerCleanupTest(t, true, false)

	initialReq := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: groupID}
	_, err := svc.Confirm(t.Context(), ws, initialReq)
	require.NoError(t, err)
	initialInserts := script.inserts

	// Unknown nonempty action must be rejected immediately, never falling through.
	_, err = svc.Confirm(t.Context(), ws, Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "some_unauthorized_mutation",
	})
	var unsupported *UnsupportedError
	require.ErrorAs(t, err, &unsupported)
	require.Equal(t, initialInserts, script.inserts, "unknown action must never write")

	// Foreign group ID must be rejected.
	_, err = svc.Confirm(t.Context(), ws, Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     "group-different",
		Action:      "cleanup_repair",
	})
	require.ErrorIs(t, err, ErrGroupMismatch)

	// Stale workspace revision must be rejected.
	staleWS := Workspace{ID: ws.ID, Revision: "rev-stale-999"}
	_, err = svc.Confirm(t.Context(), staleWS, Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "cleanup_repair",
	})
	require.ErrorIs(t, err, recordingplan.ErrPreviewTokenStale)
}

func TestSQLServerCleanupRepairFreshUnconfirmedPreviewRefusal(t *testing.T) {
	svc, script, preview, ws, groupID := setupSQLServerCleanupTest(t, false, false)

	// Calling cleanup_repair on a freshly issued preview token that was never confirmed
	// must be safely refused without panicking, and must execute zero INSERTs and zero DELETEs.
	repairReq := Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "cleanup_repair",
	}
	outcome, err := svc.Confirm(t.Context(), ws, repairReq)
	require.Nil(t, outcome)
	var unsupported *UnsupportedError
	require.ErrorAs(t, err, &unsupported)
	require.Equal(t, "cleanup-repair-operation-unconfirmed", unsupported.Reason)
	require.Equal(t, 0, script.inserts, "fresh unconfirmed repair must perform 0 INSERTs")
	require.Equal(t, 0, script.deletes, "fresh unconfirmed repair must perform 0 DELETEs")
}

func TestSQLServerCleanupRepairMalformedProofRefusal(t *testing.T) {
	svc, script, preview, ws, groupID := setupSQLServerCleanupTest(t, true, false)

	// 1. Initial confirm runs write + readback and fails cleanup -> recorded in ledger.
	initialReq := Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: groupID}
	initialOutcome, err := svc.Confirm(t.Context(), ws, initialReq)
	require.NoError(t, err)
	require.Equal(t, WriteVerified, initialOutcome.Operation.WriteOutcome)
	initialInserts := script.inserts
	initialDeletes := script.deletes

	// 2. Legitimate seed of malformed detail on terminal op via ClaimTestWriteCleanup -> FinishSchemaOperation:
	p, err := svc.buildPlan(t.Context(), groupID)
	require.NoError(t, err)
	claimed, claim, err := svc.deps.Ledger.ClaimTestWriteCleanup(t.Context(), ws.ID, preview.Token, preview.OperationID, svc.scopeFor(ws, p), p.digest)
	require.NoError(t, err)
	require.Equal(t, recordingplan.ClaimAcquired, claim)
	require.NotEmpty(t, claimed.Owner)

	_, err = svc.deps.Ledger.FinishSchemaOperation(t.Context(), preview.OperationID, claimed.Owner, recordingplan.SchemaOperationResult{
		Status:        recordingplan.SchemaOperationSucceeded,
		PayloadDigest: initialOutcome.Operation.PayloadDigest,
		WriteOutcome:  WriteVerified,
		CleanupStatus: CleanupFailed,
		Detail:        `{"phase":"planning"}`, // malformed detail: wrong phase, missing record/effect IDs
	})
	require.NoError(t, err)

	// 3. Request cleanup repair with malformed detail -> must be rejected BEFORE reclaim without any writes/deletes.
	repairReq := Confirmation{
		Token:       preview.Token,
		OperationID: preview.OperationID,
		GroupID:     groupID,
		Action:      "cleanup_repair",
	}
	outcome, err := svc.Confirm(t.Context(), ws, repairReq)
	require.Nil(t, outcome)
	var unsupported *UnsupportedError
	require.ErrorAs(t, err, &unsupported)
	require.Equal(t, ReasonProofMalformed, unsupported.Reason)
	require.Equal(t, initialInserts, script.inserts, "malformed proof repair must not perform any new INSERTs")
	require.Equal(t, initialDeletes, script.deletes, "malformed proof repair must not perform any DELETEs")
}
