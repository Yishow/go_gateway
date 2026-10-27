package recordingplan

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionUnknownAckRecoversSameOperationAfterRestart(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	first := seedSQLServerUnknownPartitionCommit(t, f)
	committed := CloneSQLServerPartitionEvidence(f.target.partitions)

	restarted := NewService(NewSQLRepository(f.db))
	stillUnknown, err := restarted.MaintainSQLServerPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Equal(t, SchemaOperationUnknown, stillUnknown.Operation.Status)
	require.Equal(t, first.Operation.OperationID, stillUnknown.Operation.OperationID)
	require.Equal(t, f.base, stillUnknown.State.Current)

	f.target.inspectErr = nil
	recovered, err := restarted.MaintainSQLServerPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Empty(t, recovered.Code, "restored visibility must resolve the recorded exact after catalog")
	require.NotNil(t, recovered.Operation)
	require.Equal(t, first.Operation.OperationID, recovered.Operation.OperationID)
	require.Equal(t, SchemaOperationSucceeded, recovered.Operation.Status)
	require.Len(t, f.target.calls, 1, "recovery inspects the committed operation without issuing another SPLIT")
	require.Equal(t, committed, recovered.State.Current)
	require.Equal(t, sqlServerPartitionMaintenanceReady, recovered.State.Status)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, lineage.Generation)
	require.False(t, lineage.HistoryMayBeLost)
}

func seedSQLServerUnknownPartitionCommit(t *testing.T, f *sqlServerPartitionMaintenanceFixture) SQLServerPartitionMaintenanceResult {
	t.Helper()
	target := f.target.target()
	split := target.Split
	target.Split = func(ctx context.Context, before, after *SQLServerPartitionEvidence) error {
		require.NoError(t, split(ctx, before, after))
		f.target.inspectErr = errors.New("scripted destination offline after commit")
		return errors.New("scripted lost commit acknowledgement")
	}
	first, err := f.service.MaintainSQLServerPartitions(t.Context(), f.request, target)
	require.NoError(t, err)
	require.NotNil(t, first.Operation)
	require.Equal(t, SchemaOperationUnknown, first.Operation.Status)
	require.Equal(t, sqlServerPartitionMaintenanceCode, first.Code)
	require.Len(t, f.target.calls, 1)
	return first
}

func TestSQLServerPartitionUnknownAckExactBeforeIsAbandoned(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	first := seedSQLServerUnknownPartitionCommit(t, f)
	f.target.inspectErr = nil
	f.target.partitions = CloneSQLServerPartitionEvidence(f.base)
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.Equal(t, first.Operation.OperationID, result.Operation.OperationID)
	require.Equal(t, SchemaOperationFailed, result.Operation.Status)
	require.Equal(t, "attempt_abandoned", result.Reason)
	require.Equal(t, f.base, result.State.Current)
	require.Len(t, f.target.calls, 1)
}

func TestSQLServerPartitionUnknownAckDoesNotAdoptForeignCatalog(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	first := seedSQLServerUnknownPartitionCommit(t, f)
	f.target.inspectErr = nil
	f.target.partitions.Partitions[2].Name = "p_foreign"
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.Equal(t, first.Operation.OperationID, result.Operation.OperationID)
	require.Equal(t, SchemaOperationUnknown, result.Operation.Status)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Equal(t, f.base, result.State.Current)
	require.Len(t, f.target.calls, 1)
}
