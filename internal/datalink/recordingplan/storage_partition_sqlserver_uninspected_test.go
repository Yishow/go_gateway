package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionMaintenanceCannotPromoteUninspectedCatalog(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	f.target.group.SQLServer.PartitionInspected = false
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.False(t, result.Changed, "cached partition fields do not prove a live inspection")
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Empty(t, f.target.calls, "uninspected catalog must not authorize SPLIT")
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, f.base, state.Current)
}

func TestSQLServerPartitionNormalizationCannotCreateInspectionAuthority(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	group, err := f.target.target().Inspect(t.Context(), f.layout.TableName)
	require.NoError(t, err)
	receipt, err := f.target.target().Inspect(t.Context(), managedEffectReceiptTable)
	require.NoError(t, err)
	group.SQLServer.PartitionInspected = false
	_, err = NormalizeSQLServerPartitionCatalog(t.Context(), f.db, f.request.Scope, f.layout, group, receipt)
	require.Error(t, err, "normalization may replace authorized boundaries, never missing live inspection")
	require.False(t, group.SQLServer.PartitionInspected)
}

func TestSQLServerPartitionFirstFailurePreservesPermissionReason(t *testing.T) {
	f := newSQLServerPartitionMaintenanceFixture(t)
	_, err := f.db.ExecContext(t.Context(), "DELETE FROM "+sqlServerPartitionMaintenanceStateTable)
	require.NoError(t, err)
	f.target.inspectErr = ErrTargetPermissionDenied
	result, err := f.maintain(t)
	require.NoError(t, err)
	require.Equal(t, sqlServerPartitionMaintenanceCode, result.Code)
	require.Equal(t, "permission_lost", result.Reason)
	require.Empty(t, f.target.calls)
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.NotNil(t, state, "the first failed inspection must remain visible")
	require.Equal(t, "blocked", state.Status)
	require.Equal(t, "permission_lost", state.Reason)
}
