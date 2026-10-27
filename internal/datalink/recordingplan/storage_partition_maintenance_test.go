package recordingplan

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLPartitionRollover(t *testing.T) {
	current, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 14, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	got, err := PlanMySQLNativeMonthlyExtension(current, time.Date(2024, time.February, 3, 8, 0, 0, 0, time.FixedZone("offset", 8*60*60)))
	require.NoError(t, err)
	require.Equal(t, 7, len(got.Partitions))
	require.Equal(t, MySQLPartitionBoundary{Name: "p_202405", LessThan: "2024-06-01T00:00:00Z", Ordinal: 6}, got.Partitions[5])
	require.Equal(t, MySQLPartitionBoundary{Name: "p_future", LessThan: mysqlPartitionMaxValue, Ordinal: 7}, got.Partitions[6])
	require.Equal(t, current.Partitions[:len(current.Partitions)-1], got.Partitions[:len(current.Partitions)-1])
}

func TestMySQLPartitionLongDowntime(t *testing.T) {
	current, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	got, err := PlanMySQLNativeMonthlyExtension(current, time.Date(2026, time.January, 22, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Equal(t, len(current.Partitions)+5, len(got.Partitions))
	require.Equal(t, MySQLPartitionBoundary{Name: "p_gap_202601", LessThan: "2026-01-01T00:00:00Z", Ordinal: 6}, got.Partitions[5])
	for month := range 4 {
		at := time.Date(2026, time.January+time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		want := MySQLPartitionBoundary{Name: fmt.Sprintf("p_%04d%02d", at.Year(), at.Month()), LessThan: at.AddDate(0, 1, 0).Format(time.RFC3339), Ordinal: 7 + month}
		require.Equal(t, want, got.Partitions[6+month])
	}
	require.Equal(t, MySQLPartitionBoundary{Name: "p_future", LessThan: mysqlPartitionMaxValue, Ordinal: len(got.Partitions)}, got.Partitions[len(got.Partitions)-1])
	for i := range len(current.Partitions) - 1 {
		require.Equal(t, current.Partitions[i], got.Partitions[i])
	}
}

func TestMySQLPartitionPermissionLost(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	f.target.reorganizeErr = ErrTargetPermissionDenied

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Equal(t, "partition_maintenance_blocked", result.Code)
	require.Equal(t, "permission_lost", result.Reason)
	require.False(t, result.Changed)
	require.Equal(t, f.base, result.State.Current)
	require.Equal(t, "blocked", result.State.Status)
	require.NotEmpty(t, result.State.Reason)
	require.Len(t, f.target.reorganizations, 1)
}

func TestMySQLPartitionAckLostUsesExactAfterCatalog(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	f.target.reorganizeErr = errors.New("acknowledgement lost")
	f.target.applyBeforeError = true

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Empty(t, result.Code)
	require.True(t, result.Changed)
	require.Equal(t, "succeeded", string(result.Operation.Status))
	require.Equal(t, f.target.reorganizations[0].after, result.State.Current)
	require.Equal(t, result.Operation.OperationID, result.State.LastOperationID)
}

func TestMySQLPartitionUnverifiableCatalogDoesNotAdopt(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	f.target.after = &MySQLPartitionEvidence{Method: mysqlPartitionMethodRangeColumns, Expression: mysqlPartitionExpression,
		Partitions: []MySQLPartitionBoundary{{Name: "foreign", LessThan: "2025-01-01T00:00:00Z", Ordinal: 1}, {Name: "p_future", LessThan: mysqlPartitionMaxValue, Ordinal: 2}}}
	f.target.applyBeforeError = true

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Equal(t, "partition_maintenance_blocked", result.Code)
	require.Equal(t, "catalog_unverifiable", result.Reason)
	require.Equal(t, f.base, result.State.Current)
	require.Equal(t, "blocked", result.State.Status)
}

func TestMySQLPartitionStaleAfterStateRecoversWithoutDDL(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	planned, err := PlanMySQLNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	f.target.partitions = planned
	op := f.seedStaleOperation(t, f.base, planned)

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.True(t, result.Changed)
	require.Equal(t, op.OperationID, result.State.LastOperationID)
	require.Equal(t, planned, result.State.Current)
	require.Empty(t, f.target.reorganizations, "stale recovery inspects but never reruns DDL")
}

func TestMySQLPartitionCatalogMismatchRefusesDDL(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	foreign, err := PlanMySQLNativeMonthlyExtension(f.base, f.request.At)
	require.NoError(t, err)
	f.target.partitions = foreign

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Equal(t, "partition_maintenance_blocked", result.Code)
	require.Equal(t, "catalog_mismatch", result.Reason)
	require.Empty(t, f.target.reorganizations)
	require.Equal(t, f.base, result.State.Current)
}

func TestMySQLPartitionBootstrapRequiresExactConfirmedBase(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM `+mysqlPartitionMaintenanceStateTable)
	require.NoError(t, err)
	f.target.partitions = &MySQLPartitionEvidence{Method: mysqlPartitionMethodRangeColumns, Expression: mysqlPartitionExpression}
	_, err = ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)

	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.NoError(t, err)
	require.Equal(t, "partition_maintenance_blocked", result.Code)
	require.Nil(t, result.State, "unverified catalog must not bootstrap sidecar ownership")
	require.Empty(t, f.target.reorganizations)
}

func TestMySQLPartitionPlanRejectsNonCanonicalCurrentState(t *testing.T) {
	base, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	cases := map[string]func(*MySQLPartitionEvidence){
		"lost future":    func(e *MySQLPartitionEvidence) { e.Partitions = e.Partitions[:len(e.Partitions)-1] },
		"changed prefix": func(e *MySQLPartitionEvidence) { e.Partitions[1].Name = "p_foreign" },
		"duplicate gap": func(e *MySQLPartitionEvidence) {
			e.Partitions = append(e.Partitions[:len(e.Partitions)-1], MySQLPartitionBoundary{Name: "p_gap_202501", LessThan: "2025-01-01T00:00:00Z", Ordinal: 6}, e.Partitions[len(e.Partitions)-1])
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			bad := cloneMySQLPartitionEvidence(base)
			mutate(bad)
			_, err := PlanMySQLNativeMonthlyExtension(bad, time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC))
			require.ErrorIs(t, err, ErrInvalidGroupSchemaLayout)
		})
	}
}

func TestMySQLPartitionEvidenceIsCopied(t *testing.T) {
	base, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	planned, err := PlanMySQLNativeMonthlyExtension(base, time.Date(2024, time.June, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.NotSame(t, base, planned)
	require.True(t, slices.Equal(base.Partitions[:len(base.Partitions)-1], planned.Partitions[:len(base.Partitions)-1]))
	planned.Partitions[0].Name = "mutated"
	require.Equal(t, "p_past", base.Partitions[0].Name)
}

func TestMySQLPartitionMaintenanceInvalidRequest(t *testing.T) {
	f := newMySQLPartitionMaintenanceFixture(t)
	f.request.At = time.Time{}
	result, err := f.service.MaintainMySQLPartitions(t.Context(), f.request, f.target.target())
	require.ErrorIs(t, err, ErrMySQLPartitionMaintenanceRequest)
	require.Nil(t, result.State)
}
