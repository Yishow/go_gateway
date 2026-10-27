package dbtarget

import (
	"database/sql/driver"
	"errors"
	"slices"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedRepairBatchUsesOneTransaction(t *testing.T) {
	db, capability, script := sqlServerRepairTestDB(t)
	statements := []string{
		"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))",
		"CREATE INDEX [dbo].[idx_circuit_long] ON [dbo].[circuit_long] ([record_id])",
	}
	request := sqlServerManagedRepairTestRequest(statements)

	result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
	require.NoError(t, err)
	require.Equal(t, len(statements), result.Committed)
	require.False(t, result.RolledBack)
	require.Equal(t, statements, script.statements)
	require.Equal(t, append(slices.Clone(statements), sqlServerRepairMarkerAddStatement), script.execStatements)
	require.Equal(t, 1, script.beginCalls)
	require.Equal(t, []driver.TxOptions{{}}, script.txOptions, "plain repair retains the default transaction isolation")
	require.Equal(t, len(statements)+1, script.execCalls)
	require.Equal(t, []int{1, 1, 1}, script.execTransactionIDs)
	require.Equal(t, 1, script.commitCalls)
	require.Zero(t, script.rollbackCalls)
	require.Equal(t, request.OperationID, script.repairOperationID)
	require.Equal(t, []driver.NamedValue{
		{Ordinal: 1, Value: sqlServerManagedRepairOperationPropertyName},
		{Ordinal: 2, Value: request.OperationID},
		{Ordinal: 3, Value: capability.Schema},
		{Ordinal: 4, Value: request.TableName},
	}, script.markerArgs)
	markerSQL := script.execStatements[len(script.execStatements)-1]
	require.Contains(t, markerSQL, "sys.sp_addextendedproperty")
	require.Contains(t, markerSQL, "@p1")
	require.Contains(t, markerSQL, "@p2")
	require.Contains(t, markerSQL, "@p3")
	require.Contains(t, markerSQL, "@p4")
	require.NotContains(t, markerSQL, request.OperationID)
	require.NotContains(t, markerSQL, capability.Schema)
	require.NotContains(t, markerSQL, request.TableName)
}

func TestSQLServerManagedRepairBatchRollsBackAndPreservesUnknownCommit(t *testing.T) {
	t.Run("statement failure rolls back", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.failExecAt = 2
		script.execError = errors.New("scripted DDL failure")
		request := sqlServerManagedRepairTestRequest([]string{
			"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))",
			"CREATE INDEX [dbo].[idx_circuit_long] ON [dbo].[circuit_long] ([record_id])",
		})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.Error(t, err)
		require.True(t, result.RolledBack)
		require.Equal(t, 2, script.execCalls)
		require.Zero(t, script.commitCalls)
		require.Equal(t, 1, script.rollbackCalls)
		require.Empty(t, script.repairOperationID)
	})

	t.Run("typed permission denial is classified", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		serverError := mssql.Error{Number: 229}
		script.failExecAt = 1
		script.execError = serverError
		request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.ErrorIs(t, err, ErrSchemaExecutionPermissionDenied)
		var numbered interface{ SQLErrorNumber() int32 }
		require.ErrorAs(t, err, &numbered)
		require.EqualValues(t, 229, numbered.SQLErrorNumber())
		require.True(t, result.RolledBack)
		require.Equal(t, 1, script.rollbackCalls)
	})

	t.Run("marker rollback uncertainty stays unknown", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.failExecAt = 3
		script.execError = errors.New("scripted marker failure")
		script.rollbackError = errors.New("scripted rollback acknowledgment loss")
		request := sqlServerManagedRepairTestRequest([]string{
			"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))",
			"CREATE INDEX [dbo].[idx_circuit_long] ON [dbo].[circuit_long] ([record_id])",
		})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.Error(t, err)
		require.False(t, result.RolledBack)
		require.Zero(t, result.Committed)
		require.Zero(t, script.commitCalls)
		require.Equal(t, 1, script.rollbackCalls)
	})

	t.Run("commit acknowledgment loss retains marker", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.commitError = errors.New("scripted commit acknowledgment loss")
		request := sqlServerManagedRepairTestRequest([]string{"CREATE TABLE [dbo].[circuit_long] ([record_id] char(64))"})

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
		require.Error(t, err)
		require.False(t, result.RolledBack)
		require.Zero(t, result.Committed)
		require.Equal(t, 1, script.commitCalls)
		require.Zero(t, script.rollbackCalls)
		require.Equal(t, request.OperationID, script.repairOperationID)
		require.Len(t, script.markerArgs, 4)
	})
}
