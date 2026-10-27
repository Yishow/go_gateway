package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mysqlRestoreFixture struct {
	db      *sql.DB
	service *Service
	request StorageRestoreRequest
	ddl     *mysqlDDLApplyFixture
}

func newMySQLRestoreFixture(t *testing.T) *mysqlRestoreFixture {
	t.Helper()
	db := lineageDB(t)
	applyRecordingDDL(t, db)
	scope := StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "recordings", StorageMode: "plain", LayoutVersion: "managed_long_v1"}
	layout, err := ManagedMySQLLongGroupLayout("gw_long_"+mysqlStorageIdentityHash(scope, "table")[:16],
		"_gw_owner_"+mysqlStorageIdentityHash(scope, "owner")[:24], MySQLGroupSchemaOptions{
			ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB", Collation: mysqlTestCollation,
		})
	require.NoError(t, err)
	group := mysqlRestoreCatalog(layout.Columns, layout.Indexes, layout.OwnerColumn, *layout.MySQL)
	receipt := mysqlRestoreCatalog(mysqlReceiptColumns(), nil, managedEffectReceiptMarker, *layout.MySQL)
	digest, err := VerifyManagedMySQLCatalog(layout, scope.Database, group, receipt)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode,
		scope.LayoutVersion, layout.TableName, layout.OwnerColumn, UnconfirmedMySQLStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm", VerifiedDigest: digest, CreatedTable: true, GroupID: "group-1",
	}, time.Now()))
	statements, err := GenerateGroupSchemaDDL(dialectMySQL, layout)
	require.NoError(t, err)
	ddl := &mysqlDDLApplyFixture{
		layout: layout, statements: statements, groupExact: group, receiptExact: receipt,
		states: map[string]TargetTableInspection{
			layout.TableName:          {Status: inspectionMissing, MySQL: cloneMySQLTableEvidence(group.MySQL)},
			managedEffectReceiptTable: receipt,
		},
	}
	return &mysqlRestoreFixture{
		db: db, service: NewService(NewSQLRepository(db)), ddl: ddl,
		request: StorageRestoreRequest{Scope: scope, Dialect: dialectMySQL, ConnectorRevision: "rev-1", TableName: layout.TableName, OwnerColumn: layout.OwnerColumn, GroupID: "group-1"},
	}
}

func mysqlRestoreCatalog(columns []GroupSchemaColumn, indexes []GroupSchemaIndex, owner string, options MySQLGroupSchemaOptions) TargetTableInspection {
	state := TargetTableInspection{
		Status: inspectionExists, ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}, IndexesInspected: true,
		MySQL: &TargetMySQLTableEvidence{
			ServerFamily: options.ServerFamily, ServerVersion: options.ServerVersion, CapabilityRevision: options.CapabilityRevision, Database: "recordings",
			Engine: options.Engine, Charset: options.Charset, Collation: options.Collation, TableType: "BASE TABLE",
			ColumnCollations: map[string]string{}, ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
		},
	}
	for _, column := range columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name], state.ColumnNullable[column.Name], state.ColumnPrimaryKey[column.Name] = column.SQLType, column.Nullable, column.PrimaryKey
		state.MySQL.ColumnDefaults[column.Name] = nil
		if column.PrimaryKey {
			state.MySQL.PrimaryKey = append(state.MySQL.PrimaryKey, column.Name)
		}
		if mysqlCollatedSQLType(column.SQLType) {
			state.MySQL.ColumnCollations[column.Name] = options.Collation
		}
		if column.SQLType == mysqlJSONSQLType {
			state.MySQL.JSONColumns[column.Name] = true
		}
	}
	marker := managedOwnerDefault
	state.MySQL.ColumnDefaults[owner] = &marker
	for _, index := range indexes {
		state.Indexes = append(state.Indexes, TargetIndexInspection{Name: index.Name, Columns: index.Columns})
	}
	return state
}

func (f *mysqlRestoreFixture) restore(ctx context.Context) (StorageRestoreResult, error) {
	return f.service.RestoreManagedStorage(ctx, f.request, StorageRestoreTarget{Inspect: f.ddl.inspect, Execute: f.ddl.execute})
}

func TestMySQLStorageRepairKeepsReceiptAndPersistsEachDDLStep(t *testing.T) {
	f := newMySQLRestoreFixture(t)
	beforeReceipt, err := json.Marshal(f.ddl.states[managedEffectReceiptTable])
	require.NoError(t, err)
	result, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageRestored, result.Outcome)
	require.Len(t, f.ddl.calls, 2, "table and index each execute once, never the receipt")
	for _, call := range f.ddl.calls {
		require.Len(t, call, 1)
	}
	afterReceipt, err := json.Marshal(f.ddl.states[managedEffectReceiptTable])
	require.NoError(t, err)
	require.Equal(t, beforeReceipt, afterReceipt)
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	require.Equal(t, 2, lineage.Generation)
	require.True(t, lineage.HistoryMayBeLost)
	require.Equal(t, 1, lineage.Reconstructions)
	operation, err := f.service.GetSchemaOperation(t.Context(), f.request.Scope.WorkspaceID, result.Operation.OperationID)
	require.NoError(t, err)
	var detail struct {
		Kind  string                       `json:"kind"`
		Steps []mysqlApplyProgressTestStep `json:"mysql_steps"`
	}
	require.NoError(t, json.Unmarshal([]byte(operation.Detail), &detail))
	require.Equal(t, StorageRepairAction, detail.Kind)
	require.Len(t, detail.Steps, 2)
	for _, step := range detail.Steps {
		require.Equal(t, mysqlSchemaStepVerified, step.Status)
		require.NotEmpty(t, step.SQLHash)
		require.NotEmpty(t, step.Fingerprint)
	}
	again, err := f.restore(t.Context())
	require.NoError(t, err)
	require.Equal(t, StorageNotMissing, again.Outcome)
	require.Len(t, f.ddl.calls, 2)
}

func TestMySQLStorageRepairRefusesUnconfirmedOrOccupiedScope(t *testing.T) {
	for name, mutate := range map[string]func(*testing.T, *mysqlRestoreFixture){
		"connector revision": func(_ *testing.T, f *mysqlRestoreFixture) { f.request.ConnectorRevision = "changed" },
		"recovery disabled": func(_ *testing.T, f *mysqlRestoreFixture) {
			f.request.Deny = &StorageDenial{Code: StorageCodeRepairBlocked, Reason: StorageReasonRecoveryDisabled}
		},
		"foreign owner": func(_ *testing.T, f *mysqlRestoreFixture) {
			foreign := cloneTargetTableInspection(f.ddl.groupExact)
			marker := "foreign"
			foreign.MySQL.ColumnDefaults[f.ddl.layout.OwnerColumn] = &marker
			f.ddl.states[f.ddl.layout.TableName] = foreign
		},
		"missing receipt": func(_ *testing.T, f *mysqlRestoreFixture) {
			f.ddl.states[managedEffectReceiptTable] = TargetTableInspection{Status: inspectionMissing, MySQL: cloneMySQLTableEvidence(f.ddl.receiptExact.MySQL)}
		},
		"unexplained missing index": func(_ *testing.T, f *mysqlRestoreFixture) {
			group := cloneTargetTableInspection(f.ddl.groupExact)
			group.Indexes = nil
			f.ddl.states[f.ddl.layout.TableName] = group
		},
		"changed layout digest": func(t *testing.T, f *mysqlRestoreFixture) {
			_, err := f.db.ExecContext(t.Context(), "UPDATE write_group_storage_scopes SET layout_digest='changed'")
			require.NoError(t, err)
		},
	} {
		t.Run(name, func(t *testing.T) {
			f := newMySQLRestoreFixture(t)
			mutate(t, f)
			result, err := f.restore(t.Context())
			require.NoError(t, err)
			require.Equal(t, StorageBlocked, result.Outcome)
			require.Empty(t, f.ddl.calls)
			lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
			require.NoError(t, err)
			require.Equal(t, 1, lineage.Generation)
			require.False(t, lineage.HistoryMayBeLost)
		})
	}
}
