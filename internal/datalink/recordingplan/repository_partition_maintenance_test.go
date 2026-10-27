package recordingplan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newPartitionMaintenanceRepositoryFixture(t *testing.T) (*SQLRepository, *MySQLPartitionMaintenanceState) {
	t.Helper()
	db := lineageDB(t)
	applyRecordingDDL(t, db)
	ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "041_mysql_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(ddl))
	require.NoError(t, err)
	current, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	state := &MySQLPartitionMaintenanceState{
		Scope: StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "recordings",
			StorageMode: "native_monthly", LayoutVersion: "managed_long_v1"},
		BaseLayoutDigest:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		ConnectorRevision: "connector-rev-1", Generation: 1, Current: current, Revision: 1,
		Status: "ready", UpdatedAt: time.Now().UTC(),
	}
	repo := NewSQLRepository(db)
	require.NoError(t, repo.ensureMySQLPartitionMaintenanceState(t.Context(), state))
	return repo, state
}

func TestReadMySQLPartitionMaintenanceStateAbsentAndPersisted(t *testing.T) {
	repo, want := newPartitionMaintenanceRepositoryFixture(t)
	got, err := ReadMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), want.Scope)
	require.NoError(t, err)
	require.Equal(t, want, got)

	absent := want.Scope
	absent.WorkspaceID = "other-workspace"
	got, err = ReadMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), absent)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestMySQLPartitionMaintenanceStateCAS(t *testing.T) {
	repo, state := newPartitionMaintenanceRepositoryFixture(t)
	planned, err := PlanMySQLNativeMonthlyExtension(state.Current, time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	next := *state
	next.Current, next.Revision, next.UpdatedAt = planned, state.Revision+1, time.Now().UTC()
	require.NoError(t, updateMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), state.Revision, &next))

	stale := next
	stale.Revision++
	stale.Current = cloneMySQLPartitionEvidence(state.Current)
	require.ErrorIs(t, updateMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), state.Revision, &stale), ErrMySQLPartitionMaintenanceStateMismatch)
	got, err := ReadMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), state.Scope)
	require.NoError(t, err)
	require.Equal(t, next.Revision, got.Revision)
	require.Equal(t, planned, got.Current)
}

func TestMySQLPartitionMaintenanceConcurrentStateCAS(t *testing.T) {
	repo, state := newPartitionMaintenanceRepositoryFixture(t)
	first, err := PlanMySQLNativeMonthlyExtension(state.Current, time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	second, err := PlanMySQLNativeMonthlyExtension(state.Current, time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	states := []*MySQLPartitionMaintenanceState{{
		Scope: state.Scope, BaseLayoutDigest: state.BaseLayoutDigest, ConnectorRevision: state.ConnectorRevision,
		Generation: state.Generation, Current: first, Revision: state.Revision + 1, Status: "ready", UpdatedAt: time.Now().UTC(),
	}, {
		Scope: state.Scope, BaseLayoutDigest: state.BaseLayoutDigest, ConnectorRevision: state.ConnectorRevision,
		Generation: state.Generation, Current: second, Revision: state.Revision + 1, Status: "ready", UpdatedAt: time.Now().UTC(),
	}}
	start := make(chan struct{})
	results := make(chan error, len(states))
	var workers sync.WaitGroup
	for _, next := range states {
		workers.Go(func() {
			<-start
			results <- updateMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), state.Revision, next)
		})
	}
	close(start)
	workers.Wait()
	close(results)
	var applied, conflicted int
	for err := range results {
		switch {
		case err == nil:
			applied++
		case errors.Is(err, ErrMySQLPartitionMaintenanceStateMismatch):
			conflicted++
		default:
			t.Fatalf("unexpected CAS result: %v", err)
		}
	}
	require.Equal(t, 1, applied)
	require.Equal(t, 1, conflicted)
}

func TestMySQLPartitionMaintenanceClaimsSharePhysicalScope(t *testing.T) {
	repo, state := newPartitionMaintenanceRepositoryFixture(t)
	first := partitionMaintenanceLedgerOperation(state.Scope, "partition-op-1", "partition-token-1", "same-physical-scope")
	_, outcome, err := repo.ClaimSchemaOperation(t.Context(), first)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	second := partitionMaintenanceLedgerOperation(state.Scope, "partition-op-2", "partition-token-2", "same-physical-scope")
	_, outcome, err = repo.ClaimSchemaOperation(t.Context(), second)
	require.NoError(t, err)
	require.Equal(t, ClaimScopeBusy, outcome)
}

func TestMySQLPartitionMaintenanceFinishRollsBackSidecarAndLedgerTogether(t *testing.T) {
	repo, state := newPartitionMaintenanceRepositoryFixture(t)
	op := partitionMaintenanceLedgerOperation(state.Scope, "partition-op-rollback", "partition-token-rollback", "scope-rollback")
	claimed, outcome, err := repo.ClaimSchemaOperation(t.Context(), op)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	planned, err := PlanMySQLNativeMonthlyExtension(state.Current, time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	rollbackErr := errors.New("forced callback failure")
	_, err = repo.finishMySQLPartitionMaintenance(t.Context(), claimed.OperationID, claimed.Owner,
		SchemaOperationResult{Status: SchemaOperationSucceeded, ExecutedStatements: 1}, func(ctx context.Context, q DBTX) error {
			next := *state
			next.Current, next.Revision, next.LastOperationID = planned, state.Revision+1, claimed.OperationID
			next.UpdatedAt = time.Now().UTC()
			if err := updateMySQLPartitionMaintenanceState(ctx, q, state.Revision, &next); err != nil {
				return err
			}
			return rollbackErr
		})
	require.ErrorIs(t, err, rollbackErr)
	stored, err := ReadMySQLPartitionMaintenanceState(t.Context(), repo.storageDB(), state.Scope)
	require.NoError(t, err)
	require.Equal(t, state.Revision, stored.Revision)
	require.Equal(t, state.Current, stored.Current)
	operation, err := repo.GetSchemaOperation(t.Context(), state.Scope.WorkspaceID, claimed.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationRunning, operation.Status)
}

func partitionMaintenanceLedgerOperation(scope StorageScope, operationID, token, scopeKey string) *SchemaOperation {
	now := time.Now().UTC()
	return &SchemaOperation{
		OperationID: operationID, Token: token, WorkspaceID: scope.WorkspaceID, ScopeKey: scopeKey,
		Owner: "claim-" + operationID, Action: MySQLPartitionMaintenanceAction, Status: SchemaOperationRunning,
		CreatedAt: now, UpdatedAt: now,
	}
}
