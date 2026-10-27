package recordingplan

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func (f *sqlServerPartitionMaintenanceFixture) maintain(t *testing.T) (SQLServerPartitionMaintenanceResult, error) {
	t.Helper()
	return f.service.MaintainSQLServerPartitions(t.Context(), f.request, f.target.target())
}

func TestSQLServerPartitionMaintenanceRolloverCaller(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	before, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Empty(t, result.Code)
	require.Len(t, f.target.calls, 1)
	require.Equal(t, f.base, f.target.calls[0].before)
	require.True(t, sameSQLServerPartitionPrefix(f.base, f.target.calls[0].after))
	after, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, f.target.partitions, after.Current)
	require.Equal(t, sqlServerPartitionMaintenanceReady, after.Status)
	require.Equal(t, int64(2), after.Revision)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, before.Generation, lineage.Generation)
	require.Equal(t, before.HistoryMayBeLost, lineage.HistoryMayBeLost)
	require.Equal(t, before.ConfirmedDigest, lineage.ConfirmedDigest)

	again, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, again.Changed)
	require.Len(t, f.target.calls, 1, "idempotent maintenance must not split again")
}

func TestSQLServerPartitionDowntimeCallerIsBounded(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	f.request.At = time.Date(2027, time.May, 20, 0, 0, 0, 0, time.UTC)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Len(t, f.target.calls, 1)
	require.Len(t, f.target.calls[0].after.Partitions, len(f.base.Partitions)+5)
	require.Equal(t, "p_gap_202705", f.target.calls[0].after.Partitions[len(f.base.Partitions)-1].Name)
	require.Equal(t, "p_future", f.target.calls[0].after.Partitions[len(f.target.calls[0].after.Partitions)-1].Name)
}

func TestSQLServerPartitionMaintenanceConcurrentScopeClaimsOnce(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	target := f.target.target()
	innerSplit := target.Split
	target.Split = func(ctx context.Context, before, after *SQLServerPartitionEvidence) error {
		once.Do(func() { close(started) })
		<-release
		return innerSplit(ctx, before, after)
	}
	resultCh := make(chan SQLServerPartitionMaintenanceResult, 1)
	errCh := make(chan error, 1)
	go func() {
		result, err := f.service.MaintainSQLServerPartitions(t.Context(), f.request, target)
		resultCh <- result
		errCh <- err
	}()
	<-started
	blocked, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, blocked.Changed)
	require.Equal(t, sqlServerPartitionMaintenanceCode, blocked.Code)
	close(release)
	require.NoError(t, <-errCh)
	require.True(t, (<-resultCh).Changed)
	require.Len(t, f.target.calls, 1)
}

func TestSQLServerPartitionMaintenanceLostAckBeforeCatalogAbandons(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	after, err := PlanSQLServerNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	stale := f.seedStaleOperation(t, f.base, after)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, result.Changed)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Equal(t, stale.OperationID, result.Operation.OperationID)
	require.Empty(t, f.target.calls, "stale recovery must not run DDL")
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, f.base, state.Current)
}

func TestSQLServerPartitionMaintenanceRejectsStaleBinding(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	f.request.ConnectorRevision = "rev-stale"
	result, err := f.maintain(t)
	require.Error(t, err)
	require.False(t, result.Changed)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Empty(t, f.target.calls)
}

func TestSQLServerPartitionMaintenanceNormalizesOnlyAuthorizedCurrentCatalog(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, result.Changed)
	target := f.target.target()
	group, err := target.Inspect(t.Context(), f.layout.TableName)
	require.NoError(t, err)
	receipt, err := target.Inspect(t.Context(), managedEffectReceiptTable)
	require.NoError(t, err)
	normalized, err := NormalizeSQLServerPartitionCatalog(t.Context(), f.db, f.request.Scope, f.layout, group, receipt)
	require.NoError(t, err)
	require.Equal(t, f.base.Partitions, normalized.SQLServer.NativeMonthly.Partitions)
	require.False(t, normalized.SQLServer.Plain)

	foreign := cloneSQLServerTargetInspectionForPartitionTest(group)
	foreign.SQLServer.NativeMonthly.Partitions[2].Name = "p_foreign"
	_, err = NormalizeSQLServerPartitionCatalog(t.Context(), f.db, f.request.Scope, f.layout, foreign, receipt)
	require.ErrorIs(t, err, ErrSQLServerPartitionMaintenanceStateMismatch)
}

func TestSQLServerPartitionMaintenanceDenied(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	f.target.splitErr = ErrTargetPermissionDenied
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, result.Changed)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Equal(t, "permission_lost", result.Reason)
	require.Len(t, f.target.calls, 1)
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, sqlServerPartitionMaintenanceBlocked, state.Status)
	require.Equal(t, f.base, state.Current)
}

func TestSQLServerPartitionMaintenanceForeignCatalogBlocksWithoutSplit(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	foreign := CloneSQLServerPartitionEvidence(f.base)
	foreign.Partitions[2].Name = "p_foreign"
	f.target.partitions = foreign
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, result.Changed)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Empty(t, f.target.calls)
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, sqlServerPartitionMaintenanceBlocked, state.Status)
}

func TestSQLServerPartitionMaintenanceLostAckResolvesSameOperation(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	after, err := PlanSQLServerNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	stale := f.seedStaleOperation(t, f.base, after)
	f.target.partitions = CloneSQLServerPartitionEvidence(after)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Equal(t, stale.OperationID, result.Operation.OperationID)
	require.Empty(t, f.target.calls, "stale recovery must not run DDL")
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, after, state.Current)
	require.Equal(t, sqlServerPartitionMaintenanceReady, state.Status)
}
