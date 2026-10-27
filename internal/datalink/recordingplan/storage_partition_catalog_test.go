package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLPartitionCatalogUsesOnlyAuthorizedExtension(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	f.target.applyBeforeError = true
	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.True(t, result.Changed)
	actual, err := f.target.target().Inspect(t.Context(), f.target.layout.TableName)
	require.NoError(t, err)
	normalized, err := NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.NoError(t, err)
	require.Equal(t, f.base, normalized.MySQL.Partitions)
	require.Equal(t, result.State.Current, actual.MySQL.Partitions, "normalization must not rewrite actual catalog evidence")
	digest, err := VerifyManagedMySQLCatalog(f.target.layout, f.request.Scope.Database, normalized, f.target.receipt)
	require.NoError(t, err)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, lineage.ConfirmedDigest, digest, "authorized extension preserves the immutable confirmation digest")

	actual.MySQL.Partitions.Partitions[0].LessThan = "2023-01-01T00:00:00Z"
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.ErrorIs(t, err, ErrIncompatibleExistingTable)
}

func TestMySQLPartitionCatalogCannotAdoptExternalExtension(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM write_group_mysql_partition_maintenance`)
	require.NoError(t, err)
	actual, err := f.target.target().Inspect(t.Context(), f.target.layout.TableName)
	require.NoError(t, err)
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.NoError(t, err, "the exact initial layout needs no maintenance sidecar yet")
	actual.MySQL.Partitions, err = PlanMySQLNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.ErrorIs(t, err, ErrIncompatibleExistingTable, "plausible month names alone cannot authorize an unrecorded catalog change")
}

func TestMySQLPartitionCatalogRejectsStaleGeneration(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, err := f.db.ExecContext(t.Context(), `UPDATE write_group_mysql_partition_maintenance SET generation = generation + 1`)
	require.NoError(t, err)
	actual, err := f.target.target().Inspect(t.Context(), f.target.layout.TableName)
	require.NoError(t, err)
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.ErrorIs(t, err, ErrMySQLPartitionMaintenanceStateMismatch)
}

func TestMySQLPartitionCatalogAllowsOnlyInitialUnconfirmedLayout(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM write_group_mysql_partition_maintenance`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET
		confirmed_layout='', confirmed_operation_id='', layout_digest=?`, UnconfirmedMySQLStorageDigest(f.request.Scope))
	require.NoError(t, err)
	actual, err := f.target.target().Inspect(t.Context(), f.target.layout.TableName)
	require.NoError(t, err)
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.NoError(t, err, "a partial CREATE can resume before the first confirmed layout exists")
	actual.MySQL.Partitions, err = PlanMySQLNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	_, err = NormalizeMySQLPartitionCatalog(t.Context(), f.db, f.request.Scope, f.target.layout, actual)
	require.ErrorIs(t, err, ErrMySQLPartitionMaintenanceStateMismatch)
}
