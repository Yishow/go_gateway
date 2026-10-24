package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestStorageRestoreFailureModes pins what each way the statements can end
// means: a clean refusal may be retried by a new attempt, anything nobody can
// vouch for blocks the scope until a confirmation checks the table again, and
// a statement is never run twice for one attempt.
func TestStorageRestoreFailureModes(t *testing.T) {
	t.Run("a rolled back batch is a failed attempt and the next one may try again", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		f.runner = func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{RolledBack: true}, errors.New("disk full")
		}
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageBlocked, result.Outcome)
		require.Equal(t, [2]string{StorageCodeRepairBlocked, StorageReasonStatementFailed}, [2]string{result.Code, result.Reason})
		require.Equal(t, SchemaOperationFailed, result.Operation.Status)
		require.Empty(t, f.lineage().RepairState, "nothing was left half done, so nothing blocks the scope")
		require.Equal(t, 1, f.lineage().Generation)

		f.runner = f.realRunner
		retried, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageRestored, retried.Outcome, "%+v", retried)
		require.Equal(t, 2, f.repairOperations(), "each attempt is its own operation")
		require.Equal(t, 2, int(f.ddlRuns.Load()))
	})

	t.Run("a permission refusal says so", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		f.runner = func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{RolledBack: true}, ErrTargetPermissionDenied
		}
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageReasonPermissionDenied, result.Reason)
		require.Equal(t, SchemaOperationFailed, result.Operation.Status)
	})

	t.Run("a lost commit acknowledgement is settled by the catalog", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		f.runner = func(ctx context.Context, statements []string) (SchemaExecution, error) {
			_, _ = f.realRunner(ctx, statements)
			return SchemaExecution{}, errors.New("commit acknowledgement lost")
		}
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageRestored, result.Outcome, "%+v", result)
		require.Equal(t, SchemaReasonRecovered, result.Operation.Reason)
		require.Equal(t, 1, int(f.ddlRuns.Load()), "the statements are not run again to find out")
		require.Equal(t, 2, f.lineage().Generation)
	})

	t.Run("a result nobody can vouch for blocks the scope and is never retried blindly", func(t *testing.T) {
		for name, tc := range map[string]struct {
			runner func(*restoreFixture) func(context.Context, []string) (SchemaExecution, error)
			// again is what the next restore says: the scope waits for a
			// confirmation, or the half built table is seen for what it is.
			again string
		}{
			"the commit result is unknown and nothing is there": {func(*restoreFixture) func(context.Context, []string) (SchemaExecution, error) {
				return func(context.Context, []string) (SchemaExecution, error) {
					return SchemaExecution{}, errors.New("connection lost during commit")
				}
			}, StorageReasonRepairNeedsConfirm},
			"only the first statement stayed applied": {func(f *restoreFixture) func(context.Context, []string) (SchemaExecution, error) {
				return func(ctx context.Context, statements []string) (SchemaExecution, error) {
					_, err := f.target.ExecContext(ctx, statements[0])
					require.NoError(t, err)
					return SchemaExecution{Committed: 1}, errors.New("second statement failed after the first committed")
				}
			}, StorageReasonTableOccupied},
		} {
			t.Run(name, func(t *testing.T) {
				f := newRestoreFixture(t)
				f.drop(f.layout.TableName)
				f.runner = tc.runner(f)
				result, err := f.restore(t.Context())
				require.NoError(t, err)
				require.Equal(t, StorageBlocked, result.Outcome)
				require.Equal(t, StorageCodeRepairBlocked, result.Code)
				require.Equal(t, SchemaOperationPartial, result.Operation.Status)
				lineage := f.lineage()
				require.Equal(t, StorageRepairStateBlocked, lineage.RepairState)
				require.Equal(t, 1, lineage.Generation, "no generation is claimed for a result nobody can vouch for")
				require.False(t, lineage.HistoryMayBeLost)

				f.runner = f.realRunner
				ran := f.ddlRuns.Load()
				blocked, err := f.restore(t.Context())
				require.NoError(t, err)
				require.Equal(t, StorageBlocked, blocked.Outcome)
				require.Equal(t, tc.again, blocked.Reason)
				require.Equal(t, ran, f.ddlRuns.Load(), "nothing is run again until a confirmation checks the table")
				require.Equal(t, 1, f.repairOperations())
			})
		}
	})

	t.Run("the table came back between the first look and the claim", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		first := true
		f.inspect = func(ctx context.Context, table string) (TargetTableInspection, error) {
			state, err := f.sqliteInspector(ctx, table)
			if table == f.layout.TableName && first {
				first = false
				statements, ddlErr := GenerateGroupSchemaDDL("sqlite", f.layout)
				require.NoError(t, ddlErr)
				for _, statement := range statements[:len(statements)-1] {
					_, execErr := f.target.ExecContext(ctx, statement)
					require.NoError(t, execErr)
				}
			}
			return state, err
		}
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageNotMissing, result.Outcome, "%+v", result)
		require.Zero(t, f.ddlRuns.Load())
		require.Equal(t, 1, f.lineage().Generation, "someone else's table is not a reconstruction")
	})
}

// TestStorageRestoreResolvesAnAbandonedAttemptFromTheCatalog: an attempt whose
// owner never came back is closed from the destination alone. A table that is
// there proves the statements committed; a table that is not closes the attempt
// as abandoned and a new, separately claimed one may follow; nothing is run
// to find out.
func TestStorageRestoreResolvesAnAbandonedAttemptFromTheCatalog(t *testing.T) {
	abandon := func(f *restoreFixture) {
		f.t.Helper()
		lineage := f.lineage()
		claim, outcome, err := f.svc.claimStorageRepair(f.t.Context(), f.svc.repo.(storageRepairRepository), f.request, lineage, *lineage.Layout)
		require.NoError(f.t, err)
		require.Equal(f.t, ClaimAcquired, outcome)
		old := time.Now().UTC().Add(-SchemaOperationLease - time.Minute)
		_, err = f.local.ExecContext(f.t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, claim.OperationID)
		require.NoError(f.t, err)
	}

	t.Run("the statements had committed", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		abandon(f)
		statements, err := GenerateGroupSchemaDDL("sqlite", f.layout)
		require.NoError(t, err)
		for _, statement := range statements[:len(statements)-1] {
			_, err := f.target.ExecContext(t.Context(), statement)
			require.NoError(t, err)
		}

		resolved, err := f.svc.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.NoError(t, err)
		require.True(t, resolved)
		require.Zero(t, f.ddlRuns.Load())
		lineage := f.lineage()
		require.Equal(t, 2, lineage.Generation, "the rebuild is recorded even though its owner died before it could")
		require.True(t, lineage.HistoryMayBeLost)

		again, err := f.svc.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.NoError(t, err)
		require.False(t, again, "an operation is resolved once")
		require.Equal(t, 2, f.lineage().Generation)
	})

	t.Run("nothing had been created", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		abandon(f)
		before := f.repairOperations()

		resolved, err := f.svc.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.inspect})
		require.NoError(t, err)
		require.True(t, resolved)
		require.False(t, f.exists(f.layout.TableName), "resolving never creates anything")
		require.Zero(t, f.ddlRuns.Load())
		require.Equal(t, 1, f.lineage().Generation)
		require.Empty(t, f.lineage().RepairState)

		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageRestored, result.Outcome, "a new attempt with its own identity may follow: %+v", result)
		require.Equal(t, before+1, f.repairOperations())
	})

	t.Run("restoring meets the abandoned attempt and resolves it first", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		abandon(f)
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageRestored, result.Outcome, "%+v", result)
		require.Equal(t, 1, int(f.ddlRuns.Load()))
		require.Equal(t, 2, f.repairOperations())
		require.Equal(t, 2, f.lineage().Generation)
	})

	t.Run("a live attempt is waited for", func(t *testing.T) {
		f := newRestoreFixture(t)
		f.drop(f.layout.TableName)
		lineage := f.lineage()
		_, outcome, err := f.svc.claimStorageRepair(t.Context(), f.svc.repo.(storageRepairRepository), f.request, lineage, *lineage.Layout)
		require.NoError(t, err)
		require.Equal(t, ClaimAcquired, outcome)
		result, err := f.restore(t.Context())
		require.NoError(t, err)
		require.Equal(t, StorageInProgress, result.Outcome)
		require.Zero(t, f.ddlRuns.Load())
	})
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	return string(encoded)
}
