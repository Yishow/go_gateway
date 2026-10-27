package recordingplan

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeMissingTableRepairUsesOwnedPartitions(t *testing.T) {
	fixture := newSQLServerNativeRestoreFixture(t)
	require.EqualValues(t, 101, fixture.request.SQLServer.Group.SQLServer.ObjectID,
		"lineage fixture starts from the exact confirmed native table proof")
	current := extendSQLServerNativeRestoreFixture(t, fixture)
	setSQLServerNativeMissingGroup(fixture, current, true, true, true)
	before, err := ReadStorageLineage(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 1, before.Generation)
	require.False(t, before.HistoryMayBeLost)
	require.Equal(t, fixture.layout, *before.Layout)
	frozenGroup, err := json.Marshal(fixture.request.SQLServer.Group)
	require.NoError(t, err)
	frozenReceipt, err := json.Marshal(fixture.request.SQLServer.Receipt)
	require.NoError(t, err)
	verifiedDigest, err := VerifyManagedSQLServerCatalog(*before.Layout, fixture.request.Scope.Database,
		fixture.request.SQLServer.Group, fixture.request.SQLServer.Receipt)
	require.NoError(t, err)
	require.Equal(t, before.ConfirmedDigest, verifiedDigest)

	var offered []string
	result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request,
		sqlServerNativeRestoreTarget(fixture, current, &offered))
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome, "code=%s reason=%s operation=%+v offered=%v executed=%v", result.Code, result.Reason, result.Operation, offered, fixture.calls)
	require.NotNil(t, result.Operation)
	require.NotEmpty(t, result.Operation.VerifiedDigest)
	require.Contains(t, strings.Join(offered, "\n"), "CREATE PARTITION FUNCTION",
		"the target transaction receives the desired definition for its ownership recheck")
	require.NotContains(t, strings.Join(fixture.calls, "\n"), "CREATE PARTITION FUNCTION",
		"an existing owned function must not be recreated")
	require.NotContains(t, strings.Join(fixture.calls, "\n"), "CREATE PARTITION SCHEME",
		"an existing owned scheme must not be recreated")
	require.Contains(t, strings.Join(fixture.calls, "\n"), "ON [ps_"+fixture.layout.TableName+"] ([bucket_start])")

	after, err := ReadStorageLineage(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, after.Generation)
	require.True(t, after.HistoryMayBeLost)
	require.Equal(t, before.LayoutDigest, after.LayoutDigest)
	require.Equal(t, before.ConfirmedDigest, after.ConfirmedDigest, "repair retains the frozen confirmation digest")
	require.Equal(t, result.Operation.OperationID, after.LastReconstructionOperationID)
	currentState, err := ReadSQLServerPartitionMaintenanceState(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	require.NotNil(t, currentState)
	require.Equal(t, 2, currentState.Generation)
	require.Equal(t, current, currentState.Current, "authorized current boundaries follow reconstruction")
	require.Equal(t, result.Operation.OperationID, currentState.LastOperationID)
	afterGroup, err := json.Marshal(fixture.request.SQLServer.Group)
	require.NoError(t, err)
	afterReceipt, err := json.Marshal(fixture.request.SQLServer.Receipt)
	require.NoError(t, err)
	require.Equal(t, frozenGroup, afterGroup, "frozen initial group proof stays byte-identical")
	require.Equal(t, frozenReceipt, afterReceipt, "receipt proof stays byte-identical")
	var detail storageRepairDetail
	require.NoError(t, json.Unmarshal([]byte(result.Operation.Detail), &detail))
	require.Equal(t, targetSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly, true, true, true, true),
		detail.SQLServer.Group.SQLServer.NativeMonthly,
		"the operation retains the frozen initial proof rather than current repair boundaries")
	lineageAfter := mustReadStorageLineage(t, fixture)
	maintenanceAfter, err := ReadSQLServerPartitionMaintenanceState(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	require.True(t, sqlServerStateMatchesLineage(maintenanceAfter, lineageAfter, SQLServerPartitionMaintenanceRequest{
		Scope: fixture.request.Scope, ConnectorRevision: lineageAfter.ConnectorRevision,
	}, lineageAfter.LayoutDigest), "state=%+v lineage=%+v", maintenanceAfter, lineageAfter)
	actualCurrent, err := sqlServerPartitionEvidenceFromTarget(fixture.states[fixture.layout.TableName].SQLServer)
	require.NoError(t, err)
	require.Equal(t, maintenanceAfter.Current, actualCurrent)
	require.NoError(t, fixture.service.VerifySQLServerStorageGeneration(t.Context(), fixture.request,
		fixture.states[fixture.layout.TableName], fixture.receipt))
}

func TestSQLServerNativeMissingTableRepairCreatesOnlyAbsentObjects(t *testing.T) {
	for _, test := range []struct {
		name                   string
		functionExists         bool
		schemeExists           bool
		wantFunctionStatements int
		wantSchemeStatements   int
	}{
		{name: "function present scheme missing", functionExists: true, wantFunctionStatements: 0, wantSchemeStatements: 1},
		{name: "function and scheme missing", wantFunctionStatements: 1, wantSchemeStatements: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			fixture := newSQLServerNativeRestoreFixture(t)
			current := CloneSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly)
			setSQLServerNativeMissingGroup(fixture, current, test.functionExists, test.schemeExists, true)
			result, err := fixture.service.RestoreManagedStorage(t.Context(), fixture.request,
				sqlServerNativeRestoreTarget(fixture, current, nil))
			require.NoError(t, err)
			require.Equal(t, StorageRestored, result.Outcome, "code=%s reason=%s", result.Code, result.Reason)
			require.Equal(t, test.wantFunctionStatements, strings.Count(strings.Join(fixture.calls, "\n"), "CREATE PARTITION FUNCTION"))
			require.Equal(t, test.wantSchemeStatements, strings.Count(strings.Join(fixture.calls, "\n"), "CREATE PARTITION SCHEME"))
			after, err := ReadStorageLineage(t.Context(), fixture.db, fixture.request.Scope)
			require.NoError(t, err)
			require.Equal(t, 2, after.Generation)
			var detail storageRepairDetail
			require.NoError(t, json.Unmarshal([]byte(result.Operation.Detail), &detail))
			require.Equal(t, targetSQLServerPartitionEvidence(fixture.layout.SQLServer.NativeMonthly, true, true, true, true),
				detail.SQLServer.Group.SQLServer.NativeMonthly)
		})
	}
}

func newSQLServerNativeRestoreFixture(t *testing.T) *sqlServerRestoreFixture {
	t.Helper()
	fixture := newSQLServerRestoreFixture(t)
	maintenanceDDL, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "042_sqlserver_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = fixture.db.ExecContext(t.Context(), string(maintenanceDDL))
	require.NoError(t, err)
	scope := fixture.request.Scope
	scope.StorageMode = sqlServerPartitionStorageMode
	table := "gw_long_" + mysqlStorageIdentityHash(scope, "table")[:16]
	owner := "_gw_owner_" + mysqlStorageIdentityHash(scope, "owner")[:24]
	function, scheme, err := SQLServerNativeMonthlyObjectNames(table)
	require.NoError(t, err)
	native, err := NewSQLServerNativeMonthlyPartition(time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC), function, scheme)
	require.NoError(t, err)
	options := fixture.options
	options.NativeMonthly = native
	layout, err := ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	group := sqlServerRestoreCatalog(layout.Columns, layout.Indexes, layout.TableName, layout.OwnerColumn, 101, options)
	group.SQLServer.Plain = false
	group.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(native, true, true, true, true)
	group.SQLServer.PartitionInspected = true

	digest, err := VerifyManagedSQLServerCatalog(layout, scope.Database, group, fixture.receipt)
	require.NoError(t, err)
	encodedLayout, err := ValidateGroupSchemaLayout(layout)
	require.NoError(t, err)
	layoutJSON, err := json.Marshal(encodedLayout)
	require.NoError(t, err)
	layoutDigest := sha256.Sum256(layoutJSON)
	_, err = fixture.db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode,
		scope.LayoutVersion, layout.TableName, layout.OwnerColumn, hex.EncodeToString(layoutDigest[:]))
	require.NoError(t, err)
	require.NoError(t, RecordStorageConfirmation(t.Context(), fixture.db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm-native", VerifiedDigest: digest,
		CreatedTable: true, GroupID: fixture.request.GroupID,
	}, time.Now()))

	missing := group
	missing.Status = inspectionMissing
	missing.Columns = nil
	missing.Indexes = nil
	encodedEvidence, err := json.Marshal(group.SQLServer)
	require.NoError(t, err)
	var missingEvidence TargetSQLServerTableEvidence
	require.NoError(t, json.Unmarshal(encodedEvidence, &missingEvidence))
	missing.SQLServer = &missingEvidence
	missing.SQLServer.ObjectID = 0
	missing.SQLServer.TableType = ""
	missing.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(native, true, true, true, false)
	fixture.layout, fixture.options, fixture.group = layout, options, group
	fixture.request.Scope, fixture.request.TableName, fixture.request.OwnerColumn = scope, layout.TableName, layout.OwnerColumn
	fixture.request.SQLServer = &SQLServerStorageRestoreEvidence{Group: group, Receipt: fixture.receipt}
	fixture.states = map[string]TargetTableInspection{
		layout.TableName:          missing,
		managedEffectReceiptTable: fixture.receipt,
	}
	return fixture
}

func extendSQLServerNativeRestoreFixture(t *testing.T, fixture *sqlServerRestoreFixture) *SQLServerPartitionEvidence {
	t.Helper()
	lineage, err := ReadStorageLineage(t.Context(), fixture.db, fixture.request.Scope)
	require.NoError(t, err)
	current, err := PlanSQLServerNativeMonthlyExtension(fixture.layout.SQLServer.NativeMonthly,
		time.Date(2027, time.February, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	repository := NewSQLRepository(fixture.db)
	require.NoError(t, repository.ensureSQLServerPartitionMaintenanceState(t.Context(), &SQLServerPartitionMaintenanceState{
		Scope: fixture.request.Scope, BaseLayoutDigest: lineage.LayoutDigest, ConnectorRevision: lineage.ConnectorRevision,
		Generation: lineage.Generation, Current: CloneSQLServerPartitionEvidence(current), Revision: 1,
		Status: sqlServerPartitionMaintenanceReady, UpdatedAt: time.Now().UTC(),
	}))
	encoded, err := json.Marshal(fixture.group)
	require.NoError(t, err)
	var group TargetTableInspection
	require.NoError(t, json.Unmarshal(encoded, &group))
	group.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, true, true, true, true)
	fixture.states[fixture.layout.TableName] = group
	return current
}

func setSQLServerNativeMissingGroup(
	fixture *sqlServerRestoreFixture,
	current *SQLServerPartitionEvidence,
	functionExists, schemeExists, exclusive bool,
) {
	missing := fixture.states[fixture.layout.TableName]
	missing.Status = inspectionMissing
	missing.Columns = nil
	missing.Indexes = nil
	missing.SQLServer.ObjectID = 0
	missing.SQLServer.TableType = ""
	missing.SQLServer.Plain = false
	missing.SQLServer.PartitionInspected = true
	if !functionExists && !schemeExists {
		missing.SQLServer.NativeMonthly = nil
	} else {
		missing.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, functionExists, schemeExists, exclusive, false)
		if !schemeExists {
			missing.SQLServer.NativeMonthly.Filegroup = ""
		}
	}
	fixture.states[fixture.layout.TableName] = missing
}

func sqlServerNativeRestoreTarget(
	fixture *sqlServerRestoreFixture,
	current *SQLServerPartitionEvidence,
	offered *[]string,
) StorageRestoreTarget {
	return StorageRestoreTarget{
		Inspect: fixture.inspect,
		Execute: func(context.Context, []string) (SchemaExecution, error) {
			return SchemaExecution{}, ErrSchemaAdapterUnverified
		},
		ExecuteSQLServerRepair: func(_ context.Context, _, operationID string, statements []string) (SchemaExecution, error) {
			if offered != nil {
				*offered = append(*offered, statements...)
			}
			observed := fixture.states[fixture.layout.TableName].SQLServer.NativeMonthly
			executed := sqlServerNativeExecutedStatements(statements, observed)
			fixture.calls = append(fixture.calls, executed...)
			recreated := sqlServerRestoreCatalog(fixture.layout.Columns, fixture.layout.Indexes,
				fixture.layout.TableName, fixture.layout.OwnerColumn, 201, fixture.options)
			recreated.SQLServer.Plain = false
			recreated.SQLServer.NativeMonthly = targetSQLServerPartitionEvidence(current, true, true, true, true)
			recreated.SQLServer.PartitionInspected = true
			recreated.SQLServer.RepairOperationID = operationID
			fixture.states[fixture.layout.TableName] = recreated
			return SchemaExecution{Committed: len(executed)}, nil
		},
	}
}

func sqlServerNativeExecutedStatements(statements []string, evidence *TargetSQLServerPartitionEvidence) []string {
	switch {
	case evidence == nil:
		return slices.Clone(statements)
	case evidence.FunctionExists && evidence.SchemeExists:
		return slices.Clone(statements[4:])
	case evidence.FunctionExists && !evidence.SchemeExists:
		return slices.Clone(statements[2:])
	default:
		return slices.Clone(statements)
	}
}

func targetSQLServerPartitionEvidence(
	native *SQLServerPartitionEvidence,
	functionExists, schemeExists, exclusive, indexesAligned bool,
) *TargetSQLServerPartitionEvidence {
	partitions := slices.Clone(native.Partitions)
	return &TargetSQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: partitions, FunctionExists: functionExists, SchemeExists: schemeExists,
		Exclusive: exclusive, IndexesAligned: indexesAligned,
	}
}
