package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerNativeRepairTestState struct {
	desired                      *SQLServerManagedPartitionEvidence
	functionExists, schemeExists bool
	functionOwned, schemeOwned   bool
	shared, tableExists          bool
	queryTransactions            []int
	allQueriesInTransaction      bool
}

func (state *sqlServerNativeRepairTestState) query(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
	transactionID int,
) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	state.queryTransactions = append(state.queryTransactions, transactionID)
	state.allQueriesInTransaction = state.allQueriesInTransaction || transactionID > 0
	rows := func(count int, values ...[]driver.Value) driver.Rows {
		columns := make([]string, count)
		for index := range columns {
			columns[index] = "catalog_value"
		}
		return sqlServerRecordingRows(columns, values...)
	}
	switch {
	case strings.Contains(query, "FROM sys.objects AS o"):
		if len(args) != 2 || args[0].Value != "dbo" || args[1].Value != "circuit_long" {
			return nil, errors.New("repair object lookup escaped the managed schema")
		}
		if !state.tableExists {
			return rows(11), nil
		}
		return rows(11, []driver.Value{int64(91), "dbo", "circuit_long", "USER_TABLE", true, false, int64(0), false, false, false, false}), nil
	case strings.Contains(query, "gw_native_partition_objects"):
		if len(args) != 3 || args[0].Value != "pf_circuit_long" || args[1].Value != "ps_circuit_long" ||
			args[2].Value != "gateway.recording.partition_owner" {
			return nil, errors.New("repair partition lookup escaped its derived object names")
		}
		var functionID, schemeID, schemeFunctionID, functionOwner, schemeOwner any
		if state.functionExists {
			functionID = int64(101)
		}
		if state.schemeExists {
			schemeID = int64(201)
			if state.functionExists {
				schemeFunctionID = int64(101)
			}
		}
		if state.functionOwned {
			functionOwner = "dbo.circuit_long"
		} else if state.functionExists {
			functionOwner = "dbo.foreign"
		}
		if state.schemeOwned {
			schemeOwner = "dbo.circuit_long"
		} else if state.schemeExists {
			schemeOwner = "dbo.foreign"
		}
		var schemeCount int64
		if state.functionExists && state.schemeExists {
			schemeCount = 1
			if state.shared {
				schemeCount = 2
			}
		}
		return rows(12, []driver.Value{functionID, true, "R", "datetime2", int64(6), int64(1), schemeID,
			schemeFunctionID, functionOwner, schemeOwner, schemeCount, true}), nil
	case strings.Contains(query, "gw_native_partition_bounds"):
		var values [][]driver.Value
		for index, boundary := range state.desired.Partitions[:len(state.desired.Partitions)-1] {
			at, err := time.ParseInLocation(sqlServerNativeBoundaryTimeLayout, boundary.LessThan, time.UTC)
			if err != nil {
				return nil, err
			}
			values = append(values, []driver.Value{int64(index + 1), int64(1), at})
		}
		return rows(3, values...), nil
	case strings.Contains(query, "gw_native_partition_destinations"):
		var values [][]driver.Value
		for index := range state.desired.Partitions {
			values = append(values, []driver.Value{int64(index + 1), "PRIMARY", "FG"})
		}
		return rows(3, values...), nil
	case strings.Contains(query, "gw_native_partition_indexes"):
		if state.shared {
			return rows(6, []driver.Value{int64(999), int64(1), int64(201), int64(1), int64(1), true}), nil
		}
		return rows(6), nil
	default:
		return nil, fmt.Errorf("unexpected native repair query: %s", query)
	}
}

func TestSQLServerManagedRepairBatchRechecksNativePartitionsInTransaction(t *testing.T) {
	for _, test := range []struct {
		name           string
		function       bool
		scheme         bool
		functionOwned  bool
		schemeOwned    bool
		shared         bool
		objectAppeared bool
		wantStatements int
		wantErr        error
	}{
		{name: "owned function and scheme", function: true, scheme: true, functionOwned: true, schemeOwned: true, wantStatements: 1},
		{name: "owned function only", function: true, functionOwned: true, wantStatements: 3},
		{name: "all objects missing", wantStatements: 5},
		{name: "scheme without function", scheme: true, schemeOwned: true, wantErr: ErrSQLServerManagedTableUnconfirmed},
		{name: "foreign function marker", function: true, scheme: true, schemeOwned: true, wantErr: ErrSQLServerManagedTableUnconfirmed},
		{name: "foreign scheme marker", function: true, scheme: true, functionOwned: true, wantErr: ErrSQLServerManagedTableUnconfirmed},
		{name: "shared partition objects", function: true, scheme: true, functionOwned: true, schemeOwned: true, shared: true, wantErr: ErrSQLServerManagedTableUnconfirmed},
		{name: "table appears before DDL", function: true, scheme: true, functionOwned: true, schemeOwned: true, objectAppeared: true, wantErr: ErrSQLServerManagedScopeChanged},
	} {
		t.Run(test.name, func(t *testing.T) {
			db, capability, script := sqlServerRepairTestDB(t)
			desired, statements := sqlServerNativeRepairTestDefinition()
			script.native = &sqlServerNativeRepairTestState{
				desired: desired, functionExists: test.function, schemeExists: test.scheme,
				functionOwned: test.functionOwned, schemeOwned: test.schemeOwned, shared: test.shared,
				tableExists: test.objectAppeared,
			}
			request := SQLServerManagedRepairRequest{
				TableName: "circuit_long", OperationID: "repair-operation-native", Statements: statements,
			}
			result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, request)
			require.Equal(t, []driver.TxOptions{{Isolation: driver.IsolationLevel(sql.LevelSerializable)}}, script.txOptions)
			require.Equal(t, []int{1, 1, 1, 1}, script.metadataQueryTransactionIDs,
				"native repair must re-read server scope and all permissions inside its serializable transaction")
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.True(t, result.RolledBack)
				require.Zero(t, script.execCalls, "catalog refusal must run no DDL or receipt marker")
				require.Equal(t, 1, script.rollbackCalls)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.wantStatements, result.Committed)
				require.Len(t, script.statements, test.wantStatements)
				require.Equal(t, "repair-operation-native", script.repairOperationID)
				require.NotEmpty(t, script.native.queryTransactions)
				require.True(t, script.native.allQueriesInTransaction, "ownership proof must be read inside DDL transaction")
				for _, transaction := range script.native.queryTransactions {
					require.Equal(t, 1, transaction)
				}
			}
		})
	}
}

func TestSQLServerNativeRepairRechecksMetadataAndRequiresLiveAlterAnyDataspaceOnlyForPFPSCreation(t *testing.T) {
	t.Run("permission snapshot changes inside transaction", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		desired, statements := sqlServerNativeRepairTestDefinition()
		script.native = &sqlServerNativeRepairTestState{desired: desired}
		changed := *script.metadata
		changed.deleteSchema = false
		script.txMetadata = &changed

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, SQLServerManagedRepairRequest{
			TableName: "circuit_long", OperationID: "repair-operation-native", Statements: statements,
		})
		require.ErrorIs(t, err, ErrSQLServerManagedScopeChanged)
		require.True(t, result.RolledBack)
		require.Equal(t, []int{1, 1, 1, 1}, script.metadataQueryTransactionIDs)
		require.Empty(t, script.native.queryTransactions, "permission drift is rejected before catalog-based DDL planning")
		require.Zero(t, script.execCalls)
		require.Equal(t, 1, script.rollbackCalls)
	})

	t.Run("destination scope changes inside transaction", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		desired, statements := sqlServerNativeRepairTestDefinition()
		script.native = &sqlServerNativeRepairTestState{desired: desired}
		changed := *script.metadata
		changed.database = "other_database"
		script.txMetadata = &changed

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, SQLServerManagedRepairRequest{
			TableName: "circuit_long", OperationID: "repair-operation-native", Statements: statements,
		})
		require.ErrorIs(t, err, ErrSQLServerManagedScopeChanged)
		require.True(t, result.RolledBack)
		require.Equal(t, []int{1}, script.metadataQueryTransactionIDs,
			"changed database is rejected by the first same-transaction metadata query")
		require.Empty(t, script.native.queryTransactions, "scope drift is rejected before catalog-based DDL planning")
		require.Zero(t, script.execCalls)
		require.Equal(t, 1, script.rollbackCalls)
	})

	t.Run("missing PF and PS require current AlterAnyDataspace", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.metadata.alterAnyDataspace = false
		sqlServerRepairTestReprojectCapability(t, db, capability, script)
		desired, statements := sqlServerNativeRepairTestDefinition()
		script.native = &sqlServerNativeRepairTestState{desired: desired}

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, SQLServerManagedRepairRequest{
			TableName: "circuit_long", OperationID: "repair-operation-native", Statements: statements,
		})
		require.ErrorIs(t, err, ErrSchemaExecutionPermissionDenied)
		require.True(t, result.RolledBack)
		require.Equal(t, []int{1, 1, 1, 1}, script.metadataQueryTransactionIDs)
		require.Zero(t, script.execCalls, "unprivileged partition-object creation must execute no DDL")
		require.Equal(t, 1, script.rollbackCalls)
	})

	t.Run("existing owned PF and PS do not require AlterAnyDataspace", func(t *testing.T) {
		db, capability, script := sqlServerRepairTestDB(t)
		script.metadata.alterAnyDataspace = false
		sqlServerRepairTestReprojectCapability(t, db, capability, script)
		desired, statements := sqlServerNativeRepairTestDefinition()
		script.native = &sqlServerNativeRepairTestState{
			desired: desired, functionExists: true, schemeExists: true,
			functionOwned: true, schemeOwned: true,
		}

		result, err := RunSQLServerManagedRepairBatch(t.Context(), db, capability, SQLServerManagedRepairRequest{
			TableName: "circuit_long", OperationID: "repair-operation-native", Statements: statements,
		})
		require.NoError(t, err)
		require.Equal(t, 1, result.Committed)
		require.Equal(t, []int{1, 1, 1, 1}, script.metadataQueryTransactionIDs)
		require.Len(t, script.statements, 1, "only the managed table is reconstructed")
		require.Equal(t, 1, script.commitCalls)
	})
}

func sqlServerRepairTestReprojectCapability(
	t *testing.T,
	db *sql.DB,
	capability *SQLServerRecordingCapability,
	script *sqlServerRepairScript,
) {
	t.Helper()
	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: capability.ConnectorRevision,
	}
	metadata := inspectSQLServerRecordingMetadata(t.Context(), db, capability.Database, capability.Schema)
	projected := projectSQLServerRecordingCapability(connector, metadata, true)
	projected.SupportsManagedLong = true
	projected.SupportsMissingRepair = true
	*capability = *projected
	script.queryCalls = 0
	script.queries = nil
}

func sqlServerNativeRepairTestDefinition() (evidence *SQLServerManagedPartitionEvidence, ddl []string) {
	values := []string{
		"2026-10-01T00:00:00.000000", "2026-11-01T00:00:00.000000", "2026-12-01T00:00:00.000000",
		"2027-01-01T00:00:00.000000", "2027-02-01T00:00:00.000000",
	}
	partitions := make([]SQLServerPartitionBoundary, 1, len(values)+1)
	partitions[0] = SQLServerPartitionBoundary{Name: "p_past", LessThan: values[0], Ordinal: 1}
	for index := 1; index < len(values); index++ {
		previous, _ := time.ParseInLocation(sqlServerNativeBoundaryTimeLayout, values[index-1], time.UTC)
		partitions = append(partitions, SQLServerPartitionBoundary{
			Name: "p_" + previous.Format("200601"), LessThan: values[index], Ordinal: index + 1,
		})
	}
	partitions = append(partitions, SQLServerPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: len(partitions) + 1})
	desired := &SQLServerManagedPartitionEvidence{
		Method: "RANGE RIGHT", Column: "bucket_start", FunctionName: "pf_circuit_long", SchemeName: "ps_circuit_long",
		Filegroup: "PRIMARY", Partitions: partitions,
	}
	quoted := make([]string, len(values))
	for index, value := range values {
		quoted[index] = "N'" + value + "'"
	}
	function := "[pf_circuit_long]"
	scheme := "[ps_circuit_long]"
	owner := "dbo.circuit_long"
	marker := "gateway.recording.partition_owner"
	statements := []string{
		"CREATE PARTITION FUNCTION " + function + " (datetime2(6)) AS RANGE RIGHT FOR VALUES (" + strings.Join(quoted, ", ") + ");",
		"EXEC sys.sp_addextendedproperty @name=N'" + marker + "', @value=N'" + owner + "', @level0type=N'PARTITION FUNCTION', @level0name=N'pf_circuit_long';",
		"CREATE PARTITION SCHEME " + scheme + " AS PARTITION " + function + " ALL TO ([PRIMARY]);",
		"EXEC sys.sp_addextendedproperty @name=N'" + marker + "', @value=N'" + owner + "', @level0type=N'PARTITION SCHEME', @level0name=N'ps_circuit_long';",
		"CREATE TABLE [dbo].[circuit_long] ([record_id] CHAR(64));",
	}
	return desired, statements
}
