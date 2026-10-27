package dbtarget

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerPartitionMaintenanceDenied(t *testing.T) {
	t.Run("capability lacks dataspace permission", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		capability.Permissions.AlterAnyDataspace = false
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.ErrorIs(t, err, ErrSchemaExecutionPermissionDenied)
		require.Zero(t, script.beginCalls)
		require.Empty(t, script.execStatements)
	})

	t.Run("live dataspace permission revoked", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		script.native.base.metadata.alterAnyDataspace = false
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.ErrorIs(t, err, ErrSchemaExecutionPermissionDenied)
		require.Equal(t, 1, script.beginCalls)
		require.Equal(t, 1, script.rollbackCalls)
		require.Empty(t, script.execStatements)
	})
}

func TestSQLServerPartitionSameScopeConcurrency(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
	script.expectedStatements = []string{
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-03-01T00:00:00.000000')",
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-04-01T00:00:00.000000')",
	}
	script.splitBoundaries = []string{"2027-03-01T00:00:00.000000", "2027-04-01T00:00:00.000000"}
	script.resetTrace()

	errorsFound := make(chan error, 2)
	var workers sync.WaitGroup
	for range 2 {
		workers.Go(func() {
			errorsFound <- SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		})
	}
	workers.Wait()
	close(errorsFound)
	for err := range errorsFound {
		require.NoError(t, err)
	}
	require.Equal(t, script.expectedStatements, script.execStatements, "same-scope calls must not duplicate already applied splits")
	require.Equal(t, 2, script.beginCalls)
	require.Equal(t, 2, script.commitCalls)
	require.Zero(t, script.rollbackCalls)
}

func TestSQLServerPartitionRestartReplayDoesNotSplitAgain(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
	_, script.expectedStatements, script.splitBoundaries = sqlServerNativePartitionMaintenancePlan(before.Evidence.NativeMonthly,
		[]SQLServerPartitionBoundary{
			{Name: "p_202702", LessThan: "2027-03-01T00:00:00.000000", Ordinal: 6},
			{Name: "p_202703", LessThan: "2027-04-01T00:00:00.000000", Ordinal: 7},
		})
	script.resetTrace()

	require.NoError(t, SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after))
	require.Len(t, script.native.bounds, 7)
	script.resetTrace()

	require.NoError(t, SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after))
	require.Empty(t, script.execStatements, "replayed authorized current state is a verified no-op")
	require.Equal(t, 1, script.beginCalls)
	require.Equal(t, 1, script.commitCalls)
}

func TestSQLServerPartitionLostAcknowledgmentResolvesFromCatalog(t *testing.T) {
	t.Run("commit applied before acknowledgment loss", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		_, script.expectedStatements, script.splitBoundaries = sqlServerNativePartitionMaintenancePlan(before.Evidence.NativeMonthly,
			[]SQLServerPartitionBoundary{
				{Name: "p_202702", LessThan: "2027-03-01T00:00:00.000000", Ordinal: 6},
				{Name: "p_202703", LessThan: "2027-04-01T00:00:00.000000", Ordinal: 7},
			})
		script.commitError = errors.New("commit acknowledgment lost")
		script.applyCommitOnError = true
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.NoError(t, err)
		require.Len(t, script.native.bounds, 7)
		require.Equal(t, 1, script.commitCalls)
	})

	t.Run("commit not applied and retry completes", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		_, script.expectedStatements, script.splitBoundaries = sqlServerNativePartitionMaintenancePlan(before.Evidence.NativeMonthly,
			[]SQLServerPartitionBoundary{
				{Name: "p_202702", LessThan: "2027-03-01T00:00:00.000000", Ordinal: 6},
				{Name: "p_202703", LessThan: "2027-04-01T00:00:00.000000", Ordinal: 7},
			})
		script.commitError = errors.New("commit acknowledgment lost")
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.ErrorIs(t, err, ErrSQLServerNativePartitionMaintenanceBlocked)
		require.False(t, errors.Is(err, ErrSQLServerNativePartitionOutcomeUnknown))
		require.Len(t, script.native.bounds, 5, "the exact before proof resolves this attempt as unapplied")

		script.commitError = nil
		script.applyCommitOnError = false
		script.resetTrace()
		require.NoError(t, SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after))
		require.Len(t, script.native.bounds, 7)
	})
}

func TestSQLServerPartitionRejectsStaleConnectorAndForeignOwnership(t *testing.T) {
	t.Run("connector revision changed", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		capability.ConnectorRevision = "revision-stale"
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.ErrorIs(t, err, ErrSQLServerManagedScopeChanged)
		require.Empty(t, script.execStatements)
	})

	t.Run("foreign partition owner", func(t *testing.T) {
		db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
		after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
		script.native.objects[8] = "dbo.foreign"
		script.resetTrace()

		err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
		require.ErrorIs(t, err, ErrSQLServerManagedPhysicalMismatch)
		require.ErrorIs(t, err, ErrSQLServerNativePartitionOutcomeUnknown)
		require.Empty(t, script.execStatements)
	})
}

func TestSQLServerPartitionDesiredFlagsCannotReplacePostProof(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
	require.False(t, after.FunctionExists)
	require.False(t, after.SchemeExists)
	require.False(t, after.Exclusive)
	require.False(t, after.IndexesAligned)
	_, script.expectedStatements, script.splitBoundaries = sqlServerNativePartitionMaintenancePlan(before.Evidence.NativeMonthly,
		[]SQLServerPartitionBoundary{
			{Name: "p_202702", LessThan: "2027-03-01T00:00:00.000000", Ordinal: 6},
			{Name: "p_202703", LessThan: "2027-04-01T00:00:00.000000", Ordinal: 7},
		})
	script.foreignOwnerAfterFinalSplit = true
	script.resetTrace()

	err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
	require.ErrorIs(t, err, ErrSQLServerNativePartitionMaintenanceBlocked)
	require.ErrorIs(t, err, ErrSQLServerNativePartitionOutcomeUnknown)
	require.Len(t, script.native.bounds, 5, "failed after-proof rolls back the transaction")
	require.Len(t, script.execStatements, 4, "a desired layout cannot stand in for actual ownership proof")
}

func TestSQLServerPartitionPartialCatalogRemainsUnknown(t *testing.T) {
	db, capability, before, script := sqlServerNativePartitionMaintenanceFixture(t)
	after := sqlServerNativePartitionAfterTwoMonths(before.Evidence.NativeMonthly)
	script.expectedStatements = []string{
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-03-01T00:00:00.000000')",
		"ALTER PARTITION SCHEME [ps_circuit_long] NEXT USED [PRIMARY]",
		"ALTER PARTITION FUNCTION [pf_circuit_long]() SPLIT RANGE ('2027-04-01T00:00:00.000000')",
	}
	script.splitBoundaries = []string{"2027-03-01T00:00:00.000000", "2027-04-01T00:00:00.000000"}
	script.failExecAt = 3
	script.execError = errors.New("second NEXT USED failed")
	script.rollbackError = errors.New("rollback acknowledgment lost")
	script.applyOnRollbackError = true
	script.resetTrace()

	err := SplitOpenSQLServerNativeMonthly(t.Context(), db, capability, before, after)
	require.ErrorIs(t, err, ErrSQLServerNativePartitionMaintenanceBlocked)
	require.ErrorIs(t, err, ErrSQLServerNativePartitionOutcomeUnknown)
	require.Len(t, script.native.bounds, 6, "one observed split is neither the exact before nor exact after proof")
	require.Equal(t, 3, script.execAttempts, "partial proof must not trigger a blind append retry")
}
