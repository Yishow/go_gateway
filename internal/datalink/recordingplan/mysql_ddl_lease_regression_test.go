package recordingplan

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLSchemaApplyHeartbeatKeepsLeaseDuringBlockedDDL(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		repo := NewMemoryRepository()
		fixture := newMySQLDDLApplyFixture(t, repo)
		service := NewService(repo)

		var firstStatementExecutions atomic.Int32
		firstExecutorStarted := make(chan struct{})
		releaseFirstExecutor := make(chan struct{})
		apply := func(ctx context.Context, statements []string) (SchemaExecution, error) {
			if len(statements) == 1 && statements[0] == fixture.statements[0] {
				switch firstStatementExecutions.Add(1) {
				case 1:
					close(firstExecutorStarted)
					select {
					case <-releaseFirstExecutor:
					case <-ctx.Done():
						return SchemaExecution{}, ctx.Err()
					}
				case 2:
					// The unchanged implementation reaches this branch after a
					// retry takes over the stale operation.
				}
			}
			return fixture.execute(ctx, statements)
		}
		target := SchemaApplyTarget{Inspect: fixture.inspect, Execute: apply}

		type applyResult struct {
			op      *SchemaOperation
			outcome ClaimOutcome
			err     error
		}
		firstResult := make(chan applyResult, 1)
		go func() {
			op, outcome, err := service.ApplySchemaPreview(t.Context(), fixture.token, target)
			firstResult <- applyResult{op: op, outcome: outcome, err: err}
		}()
		<-firstExecutorStarted
		synctest.Wait()

		claimed, err := repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
		require.NoError(t, err)
		ownerBefore := claimed.Owner
		require.NotEmpty(t, ownerBefore)

		// Move fake time past the real ten-minute lease while the first DDL
		// remains blocked. The retry must observe the heartbeat-refreshed
		// owner rather than start a second executor.
		time.Sleep(SchemaOperationLease + time.Second)
		synctest.Wait()

		retryResult := make(chan applyResult, 1)
		go func() {
			op, outcome, err := service.ApplySchemaPreview(t.Context(), fixture.token, target)
			retryResult <- applyResult{op: op, outcome: outcome, err: err}
		}()
		synctest.Wait()
		retry := <-retryResult

		close(releaseFirstExecutor)
		synctest.Wait()
		first := <-firstResult

		current, err := repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
		require.NoError(t, err)
		require.Equal(t, int32(1), firstStatementExecutions.Load(),
			"a live DDL heartbeat must prevent a second service executor")
		require.Equal(t, ownerBefore, current.Owner,
			"a live DDL heartbeat must keep the original operation owner")
		require.NoError(t, retry.err)
		require.Equal(t, ClaimInProgress, retry.outcome)
		require.NotNil(t, retry.op)
		require.NoError(t, first.err)
		require.Equal(t, SchemaOperationSucceeded, first.op.Status)
		found, err := inspectMySQLGroupTables(t.Context(), fixture.layout, fixture.token.Database, fixture.inspect)
		require.NoError(t, err)
		require.False(t, found.pendingGroup)
		require.False(t, found.pendingReceipt)
	})
}
