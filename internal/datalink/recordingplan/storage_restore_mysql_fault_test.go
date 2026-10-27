package recordingplan

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLStorageRepairFailureDoesNotClaimHistoryRecovered(t *testing.T) {
	for _, partial := range []bool{false, true} {
		name := "create denied"
		if partial {
			name = "index denied after table committed"
		}
		t.Run(name, func(t *testing.T) {
			f := newMySQLRestoreFixture(t)
			index := 0
			if partial {
				index = 1
			}
			f.ddl.failStatement, f.ddl.failErr = f.ddl.statements[index], ErrTargetPermissionDenied
			result, err := f.restore(t.Context())
			require.NoError(t, err)
			require.Equal(t, StorageBlocked, result.Outcome)
			lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
			require.NoError(t, err)
			require.Equal(t, 1, lineage.Generation)
			require.False(t, lineage.HistoryMayBeLost)
			require.Zero(t, lineage.Reconstructions)
			if partial {
				require.Equal(t, SchemaOperationPartial, result.Operation.Status)
				require.Equal(t, StorageRepairStateBlocked, lineage.RepairState)
				calls := len(f.ddl.calls)
				_, err := f.restore(t.Context())
				require.NoError(t, err)
				require.Len(t, f.ddl.calls, calls, "a partial automatic repair waits for explicit confirmation")
			} else {
				require.Equal(t, SchemaOperationFailed, result.Operation.Status)
				require.Equal(t, StorageReasonPermissionDenied, result.Reason)
				require.Empty(t, lineage.RepairState)
			}
		})
	}
}

func TestMySQLStorageRepairSettlesLostDDLResponseFromCatalog(t *testing.T) {
	f := newMySQLRestoreFixture(t)
	f.ddl.loseAck = true
	result, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.Len(t, f.ddl.calls, 2, "each applied DDL is confirmed without executing it again")
}

func TestMySQLStorageRepairSkipsAnAlreadyRecreatedOwnedTable(t *testing.T) {
	f := newMySQLRestoreFixture(t)
	first := true
	inspect := func(ctx context.Context, table string) (TargetTableInspection, error) {
		state, err := f.ddl.inspect(ctx, table)
		if table == f.ddl.layout.TableName && first {
			first = false
			group := cloneTargetTableInspection(f.ddl.groupExact)
			group.Indexes = nil
			f.ddl.states[table] = group
		}
		return state, err
	}
	result, err := f.service.RestoreManagedStorage(t.Context(), f.request, StorageRestoreTarget{Inspect: inspect, Execute: f.ddl.execute})
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.Equal(t, [][]string{{f.ddl.statements[1]}}, f.ddl.calls, "only the missing index is executed under the acquired repair claim")
}

func TestMySQLStorageRepairStaleFinishRetainsStepEvidence(t *testing.T) {
	f := newMySQLRestoreFixture(t)
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER reject_mysql_repair_finish
		BEFORE UPDATE ON managed_schema_operations WHEN NEW.status='succeeded'
		BEGIN SELECT RAISE(ABORT, 'finish unavailable'); END`)
	require.NoError(t, err)
	first, err := f.restore(t.Context())
	require.ErrorIs(t, err, ErrSchemaOperationUnacknowledged)
	require.NotNil(t, first.Operation)
	operationID := first.Operation.OperationID
	before, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, operationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationRunning, before.Status)
	require.NotEmpty(t, before.Detail)
	_, err = f.db.ExecContext(t.Context(), "DROP TRIGGER reject_mysql_repair_finish")
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), "UPDATE managed_schema_operations SET updated_at=? WHERE operation_id=?",
		time.Now().UTC().Add(-2*SchemaOperationLease), operationID)
	require.NoError(t, err)
	resolved, err := f.service.ResolveStaleStorageRepair(t.Context(), f.request, StorageRestoreTarget{Inspect: f.ddl.inspect, Execute: func(context.Context, []string) (SchemaExecution, error) {
		return SchemaExecution{}, errors.New("stale resolution must never execute DDL")
	}})
	require.NoError(t, err)
	require.True(t, resolved)
	after, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, operationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationSucceeded, after.Status)
	require.JSONEq(t, before.Detail, after.Detail)
	var detail struct {
		Steps []mysqlApplyProgressTestStep `json:"mysql_steps"`
	}
	require.NoError(t, json.Unmarshal([]byte(after.Detail), &detail))
	require.Len(t, detail.Steps, 2)
	for _, step := range detail.Steps {
		require.Equal(t, mysqlSchemaStepVerified, step.Status)
	}
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, lineage.Generation)
	require.True(t, lineage.HistoryMayBeLost)
	require.Len(t, f.ddl.calls, 2)
}
