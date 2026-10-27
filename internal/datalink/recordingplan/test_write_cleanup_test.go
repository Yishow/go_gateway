package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newTestSQLRepository(t *testing.T) *SQLRepository {
	t.Helper()
	return NewSQLRepository(recordingRepositoryDB(t))
}

func TestClaimTestWriteCleanup(t *testing.T) {
	for name, repo := range map[string]Repository{
		"memory": NewMemoryRepository(),
		"sql":    newTestSQLRepository(t),
	} {
		t.Run(name, func(t *testing.T) {
			svc := NewService(repo)
			scope := testWriteScope()
			token, err := svc.PrepareTestWritePreview(t.Context(), scope, "content-1")
			require.NoError(t, err)

			op, claim, err := svc.ClaimTestWrite(t.Context(), token)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, claim)

			detail := `{"record_ids":["r1","r2"],"effect_keys":["k1","k2"]}`
			// Finish with write verified and cleanup failed
			_, err = svc.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{
				Status:        SchemaOperationSucceeded,
				PayloadDigest: "content-1",
				WriteOutcome:  writeOutcomeVerified,
				CleanupStatus: "failed",
				CleanupReason: "cleanup-denied",
				Detail:        detail,
			})
			require.NoError(t, err)

			// 1. Reclaim cleanup atomically
			reclaimed, claim, err := svc.ClaimTestWriteCleanup(t.Context(), "ws-1", token.Token, token.OperationID, scope, "content-1")
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, claim)
			require.Equal(t, SchemaOperationRunning, reclaimed.Status)

			// Concurrent active retry before Finish returns ClaimInProgress with same opID
			inProgress, claim, err := svc.ClaimTestWriteCleanup(t.Context(), "ws-1", token.Token, token.OperationID, scope, "content-1")
			require.NoError(t, err)
			require.Equal(t, ClaimInProgress, claim)
			require.Equal(t, reclaimed.OperationID, inProgress.OperationID)

			// 2. Finish cleanup by first owner still succeeds
			_, err = svc.FinishSchemaOperation(t.Context(), reclaimed.OperationID, reclaimed.Owner, SchemaOperationResult{
				Status:        SchemaOperationSucceeded,
				PayloadDigest: "content-1",
				WriteOutcome:  writeOutcomeVerified,
				CleanupStatus: cleanupStatusCleaned,
				Detail:        detail,
			})
			require.NoError(t, err)

			// 3. Reclaim when already cleaned returns ClaimCompleted
			alreadyCleaned, claim, err := svc.ClaimTestWriteCleanup(t.Context(), "ws-1", token.Token, token.OperationID, scope, "content-1")
			require.NoError(t, err)
			require.Equal(t, ClaimCompleted, claim)
			require.Equal(t, cleanupStatusCleaned, alreadyCleaned.CleanupStatus)
		})
	}
}

func TestClaimTestWriteCleanupPrerequisites(t *testing.T) {
	svc := newTestWriteService(t)
	scope := testWriteScope()
	token, err := svc.PrepareTestWritePreview(t.Context(), scope, "content-1")
	require.NoError(t, err)

	op, _, err := svc.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)

	// Finish with WriteFailed -> cleanup repair cannot proceed
	_, err = svc.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{
		Status:        SchemaOperationFailed,
		PayloadDigest: "content-1",
		WriteOutcome:  "failed",
		CleanupStatus: "not_attempted",
		Detail:        `{"record_ids":["r1"]}`,
	})
	require.NoError(t, err)

	_, _, err = svc.ClaimTestWriteCleanup(t.Context(), "ws-1", token.Token, token.OperationID, scope, "content-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot repair cleanup when write is not verified")

	// Content digest mismatch
	_, _, err = svc.ClaimTestWriteCleanup(t.Context(), "ws-1", token.Token, token.OperationID, scope, "content-changed")
	require.ErrorIs(t, err, ErrPreviewTokenStale)
}
