package recordingplan

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecordingStartLedgerRequestIdentityAndAction(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	digest := strings.Repeat("a", 64)
	op, outcome, err := service.ClaimRecordingStart(t.Context(), "workspace", "request", digest, `{"stage":"save"}`)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, RecordingStartAction, op.Action)
	require.NotEmpty(t, op.Owner)
	repeated, outcome, err := service.ClaimRecordingStart(t.Context(), "workspace", "request", digest, `{"untrusted":"replacement"}`)
	require.NoError(t, err)
	require.Equal(t, ClaimInProgress, outcome)
	require.Equal(t, op.OperationID, repeated.OperationID)
	require.Equal(t, op.Detail, repeated.Detail)
	require.Empty(t, repeated.Owner)
	_, _, err = service.ClaimRecordingStart(t.Context(), "workspace", "request", strings.Repeat("b", 64), "{}")
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	_, err = service.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{
		Status: SchemaOperationPartial, PayloadDigest: digest, Detail: `{"stage":"activation","applied":true}`,
	})
	require.NoError(t, err)
	resumed, outcome, err := service.ResumeRecordingStart(t.Context(), "workspace", op.OperationID, digest)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.NotEqual(t, op.Owner, resumed.Owner)
	require.Contains(t, resumed.Detail, `"applied":true`)
	_, err = service.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded})
	require.ErrorIs(t, err, ErrSchemaOperationNotOwned)
	_, err = service.FinishSchemaOperation(t.Context(), op.OperationID, resumed.Owner, SchemaOperationResult{
		Status: SchemaOperationSucceeded, PayloadDigest: digest, Detail: resumed.Detail,
	})
	require.NoError(t, err)
	completed, outcome, err := service.ResumeRecordingStart(t.Context(), "workspace", op.OperationID, digest)
	require.NoError(t, err)
	require.Equal(t, ClaimCompleted, outcome)
	require.Equal(t, SchemaOperationSucceeded, completed.Status)
	require.Empty(t, completed.Owner)
	_, _, err = service.ResumeRecordingStart(t.Context(), "foreign", op.OperationID, digest)
	require.ErrorIs(t, err, ErrSchemaOperationNotFound)
}

func TestRecordingStartResumeHasOneOwnerAndRetainsLiveLease(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)
	digest := strings.Repeat("c", 64)
	op, _, err := service.ClaimRecordingStart(t.Context(), "workspace", "request", digest, "{}")
	require.NoError(t, err)
	_, outcome, err := service.ResumeRecordingStart(t.Context(), "workspace", op.OperationID, digest)
	require.NoError(t, err)
	require.Equal(t, ClaimInProgress, outcome)
	repo.mu.Lock()
	stored := repo.operations[op.OperationID]
	stored.UpdatedAt = time.Now().Add(-2 * SchemaOperationLease)
	repo.operations[op.OperationID] = stored
	repo.mu.Unlock()
	var wg sync.WaitGroup
	results := make(chan ClaimOutcome, 12)
	errors := make(chan error, 12)
	for range 12 {
		wg.Go(func() {
			_, outcome, err := service.ResumeRecordingStart(t.Context(), "workspace", op.OperationID, digest)
			results <- outcome
			errors <- err
		})
	}
	wg.Wait()
	close(results)
	close(errors)
	acquired := 0
	for result := range results {
		if result == ClaimAcquired {
			acquired++
		}
	}
	for err := range errors {
		require.NoError(t, err)
	}
	require.Equal(t, 1, acquired)
}

func TestRecordingStartResumeCannotReopenSchemaOrTestWrite(t *testing.T) {
	for _, action := range []string{"create_tables", "test_write"} {
		t.Run(action, func(t *testing.T) {
			repo := NewMemoryRepository()
			_, _, err := repo.ClaimSchemaOperation(t.Context(), &SchemaOperation{
				OperationID: "other", Token: "other", WorkspaceID: "workspace", Action: action,
				Status: SchemaOperationPartial, PayloadDigest: strings.Repeat("d", 64),
			})
			require.NoError(t, err)
			_, _, err = NewService(repo).ResumeRecordingStart(t.Context(), "workspace", "other", strings.Repeat("d", 64))
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
		})
	}
}
