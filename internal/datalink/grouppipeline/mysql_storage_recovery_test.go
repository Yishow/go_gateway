package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestMySQLDestinationResolverCarriesOpenedCapabilityForRepair(t *testing.T) {
	frozen := mysqlRuntimeSnapshot(t)
	actual := mysqlRuntimeCapability(frozen.RuntimeLayout.MySQL)
	resolver, item, _ := newMySQLResolverFixture(t, frozen, actual)
	resolver.loadRuntimeSnapshot = func(_ context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
		require.Equal(t, item.Key, key)
		return frozen, nil
	}

	target, err := resolver.Resolve(t.Context(), item)

	require.NoError(t, err)
	require.Same(t, frozen.RuntimeLayout.MySQL, target.MySQL, "the row adapter still receives the frozen layout proof")
	require.Same(t, actual, target.MySQLCapability, "repair inspection binds to the capability from this opened destination")
	require.NoError(t, target.Close())
}

func TestMySQLRestoreTargetUsesActualCapabilityAndPreservesCatalogEvidence(t *testing.T) {
	snapshot := mysqlRuntimeSnapshot(t)
	capability := mysqlRuntimeCapability(snapshot.RuntimeLayout.MySQL)
	db := &sql.DB{}
	probe := &mysqlStorageRestoreProbe{
		inspection: &dbtarget.MySQLManagedInspection{
			Inspection: &dbtarget.TableInspection{
				Status: dbtarget.TableInspectionExists, Schema: "gwtest", Table: snapshot.Group.Destination.TableName,
				Columns: []dbtarget.ColumnInfo{{Name: "bucket_start", DataType: "DATETIME(6)", PrimaryKey: true}},
				Indexes: []dbtarget.IndexInfo{{Name: "idx_source", Columns: []string{"workspace_id", "tag_id"}}}, IndexesInspected: true,
			},
			Evidence: &dbtarget.MySQLManagedTableEvidence{
				ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Database: "gwtest",
				Engine: "InnoDB", Charset: "utf8mb4", Collation: "utf8mb4_0900_bin", TableType: "BASE TABLE",
				PrimaryKey: []string{"bucket_start"}, ColumnCollations: map[string]string{"bucket_start": "utf8mb4_0900_bin"},
				ColumnDefaults: map[string]*string{}, JSONColumns: map[string]bool{},
			},
		},
	}
	target := mysqlStorageRestoreTargetWithOps(groupdelivery.Target{
		DB: db, Kind: schema.DatabaseConnectorKindMySQL, MySQLCapability: capability,
	}, "gwtest", mysqlStorageRestoreTargetOps{Inspect: probe.inspect, Execute: probe.execute})

	state, err := target.Inspect(t.Context(), snapshot.Group.Destination.TableName)

	require.NoError(t, err)
	require.Nil(t, target.CanCreate, "the exact DDL itself is the MySQL permission check")
	require.Same(t, db, probe.inspectedDB)
	require.Same(t, capability, probe.inspectedCapability)
	require.Equal(t, "gwtest", probe.inspectedDatabase)
	require.Equal(t, snapshot.Group.Destination.TableName, probe.inspectedTable)
	require.Equal(t, "exists", state.Status)
	require.Equal(t, []string{"bucket_start"}, state.Columns)
	require.Equal(t, "DATETIME(6)", state.ColumnTypes["bucket_start"])
	require.True(t, state.ColumnPrimaryKey["bucket_start"])
	require.True(t, state.IndexesInspected)
	require.Equal(t, "idx_source", state.Indexes[0].Name)
	require.Equal(t, "mysql", state.MySQL.ServerFamily)
	require.Equal(t, "utf8mb4_0900_bin", state.MySQL.Collation)
}

func TestMySQLRestoreTargetRunsOnlyOneDDLWithoutClaimingRollback(t *testing.T) {
	db := &sql.DB{}
	capability := &dbtarget.MySQLRecordingCapability{Connected: true}
	permissionErr := &mysqldriver.MySQLError{Number: 1142}
	probe := &mysqlStorageRestoreProbe{executeErr: permissionErr}
	target := mysqlStorageRestoreTargetWithOps(groupdelivery.Target{
		DB: db, Kind: schema.DatabaseConnectorKindMySQL, MySQLCapability: capability,
	}, "gwtest", mysqlStorageRestoreTargetOps{Inspect: probe.inspect, Execute: probe.execute})

	execution, err := target.Execute(t.Context(), []string{"CREATE TABLE `gw_long_x` (`id` BIGINT)"})

	require.ErrorIs(t, err, recordingplan.ErrTargetPermissionDenied)
	require.Equal(t, 0, execution.Committed)
	require.False(t, execution.RolledBack, "MySQL DDL is not wrapped in a row or DDL transaction")
	require.Equal(t, 1, probe.executeCalls)
	require.Same(t, db, probe.executedDB)
	require.Same(t, capability, probe.executedCapability)
	require.Equal(t, "gwtest", probe.executedDatabase)
	require.Equal(t, []string{"CREATE TABLE `gw_long_x` (`id` BIGINT)"}, probe.executedStatements)

	_, err = target.Execute(t.Context(), []string{"CREATE TABLE `one` (`id` BIGINT)", "CREATE INDEX `two` ON `one` (`id`)"})

	require.ErrorIs(t, err, dbtarget.ErrMySQLManagedStatementRefused)
	require.Equal(t, 1, probe.executeCalls, "a multi-statement request is rejected before the target is touched")
}

func TestMySQLRestoreTargetWithNilCapabilityDoesNoInspectionOrDDL(t *testing.T) {
	probe := &mysqlStorageRestoreProbe{}
	target := mysqlStorageRestoreTargetWithOps(groupdelivery.Target{
		DB: &sql.DB{}, Kind: schema.DatabaseConnectorKindMySQL,
	}, "gwtest", mysqlStorageRestoreTargetOps{Inspect: probe.inspect, Execute: probe.execute})

	_, inspectErr := target.Inspect(t.Context(), "gw_long_x")
	execution, executeErr := target.Execute(t.Context(), []string{"CREATE TABLE `gw_long_x` (`id` BIGINT)"})

	require.ErrorIs(t, inspectErr, dbtarget.ErrMySQLManagedServerMismatch)
	require.ErrorIs(t, executeErr, dbtarget.ErrMySQLManagedServerMismatch)
	require.Zero(t, probe.inspectCalls)
	require.Zero(t, probe.executeCalls)
	require.Zero(t, execution.Committed)
	require.False(t, execution.RolledBack)
}

func TestMySQLRecoveryKeepsFrozenScopePolicyAndRevisionGuards(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*workspace.WriteGroupAppliedSnapshot, *groupdelivery.OutboxItem)
		code   string
		reason string
	}{
		{
			name: "policy absent",
			mutate: func(snapshot *workspace.WriteGroupAppliedSnapshot, _ *groupdelivery.OutboxItem) {
				snapshot.Group.WritePolicy.RecreateMissing = nil
			},
			code: recordingplan.StorageCodeRepairBlocked, reason: recordingplan.StorageReasonPolicyUnconfirmed,
		},
		{
			name: "policy disabled",
			mutate: func(snapshot *workspace.WriteGroupAppliedSnapshot, _ *groupdelivery.OutboxItem) {
				disabled := false
				snapshot.Group.WritePolicy.RecreateMissing = &disabled
			},
			code: recordingplan.StorageCodeRepairBlocked, reason: recordingplan.StorageReasonRecoveryDisabled,
		},
		{
			name: "connector revision changed",
			mutate: func(_ *workspace.WriteGroupAppliedSnapshot, item *groupdelivery.OutboxItem) {
				item.Destination.ConnectorRevision = "changed-revision"
			},
			code: recordingplan.StorageCodeOwnerMismatch, reason: recordingplan.StorageReasonOwnershipMismatch,
		},
		{
			name: "database scope changed",
			mutate: func(_ *workspace.WriteGroupAppliedSnapshot, item *groupdelivery.OutboxItem) {
				item.Destination.Database = "other-database"
			},
			code: recordingplan.StorageCodeOwnerMismatch, reason: recordingplan.StorageReasonOwnershipMismatch,
		},
		{
			name: "group confirmation missing",
			mutate: func(snapshot *workspace.WriteGroupAppliedSnapshot, _ *groupdelivery.OutboxItem) {
				enabled := true
				snapshot.Group.WritePolicy.RecreateMissing = &enabled
				snapshot.Group.Destination.SchemaRevision = ""
			},
			code: recordingplan.StorageCodeOwnerMismatch, reason: recordingplan.StorageReasonConfirmationMissing,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			snapshot := mysqlRuntimeSnapshot(t)
			capability := mysqlRuntimeCapability(snapshot.RuntimeLayout.MySQL)
			_, item, _ := newMySQLResolverFixture(t, snapshot, capability)
			table, _, ok := workspace.LongStorageIdentity(snapshot.Group)
			require.True(t, ok)
			test.mutate(snapshot, &item)

			denial := (&StorageRecovery{}).denial(t.Context(), snapshot, item, table, "mysql")

			require.NotNil(t, denial)
			require.Equal(t, test.code, denial.Code)
			require.Equal(t, test.reason, denial.Reason)
		})
	}
}

type mysqlStorageRestoreProbe struct {
	inspection *dbtarget.MySQLManagedInspection
	inspectErr error
	executeErr error

	inspectedDB         *sql.DB
	inspectedCapability *dbtarget.MySQLRecordingCapability
	inspectedDatabase   string
	inspectedTable      string
	inspectCalls        int

	executedDB         *sql.DB
	executedCapability *dbtarget.MySQLRecordingCapability
	executedDatabase   string
	executedStatements []string
	executeCalls       int
}

func (p *mysqlStorageRestoreProbe) inspect(
	_ context.Context,
	db *sql.DB,
	capability *dbtarget.MySQLRecordingCapability,
	database, table string,
) (*dbtarget.MySQLManagedInspection, error) {
	p.inspectCalls++
	p.inspectedDB, p.inspectedCapability = db, capability
	p.inspectedDatabase, p.inspectedTable = database, table
	return p.inspection, p.inspectErr
}

func (p *mysqlStorageRestoreProbe) execute(
	_ context.Context,
	db *sql.DB,
	capability *dbtarget.MySQLRecordingCapability,
	database, statement string,
) error {
	p.executeCalls++
	p.executedDB, p.executedCapability = db, capability
	p.executedDatabase = database
	p.executedStatements = append(p.executedStatements, statement)
	return p.executeErr
}

func TestMySQLRecoveryFallbackRequiresConfirmedCanonicalLineage(t *testing.T) {
	fixture := newFixture(t)
	snapshot := mysqlRuntimeSnapshot(t)
	scope := workspace.StorageScopeOf(snapshot.Group)
	table, owner, ok := workspace.LongStorageIdentity(snapshot.Group)
	require.True(t, ok)
	evidence := snapshot.RuntimeLayout.MySQL.Group.Evidence
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion,
		CapabilityRevision: evidence.CapabilityRevision, Engine: evidence.Engine,
		Charset: evidence.Charset, Collation: evidence.Collation,
	})
	require.NoError(t, err)
	_, err = fixture.store.DB().ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema,
		scope.StorageMode, scope.LayoutVersion, table, owner, recordingplan.UnconfirmedMySQLStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), fixture.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: snapshot.Group.Destination.ConnectorRevision,
		OperationID: snapshot.Group.Destination.SchemaRevision, VerifiedDigest: snapshot.Group.Destination.SchemaDigest,
		GroupID: snapshot.Group.ID,
	}, time.Now()))
	recovery := &StorageRecovery{store: fixture.store}

	fallback, lineage, err := recovery.mysqlConfirmedFallbackLayout(t.Context(), scope, table, owner)

	require.NoError(t, err)
	require.NotNil(t, lineage)
	require.NotNil(t, fallback)
	require.Equal(t, layout, *fallback)
	_, err = fixture.store.DB().ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET layout_digest = 'tampered'
		WHERE workspace_id = ? AND connector_id = ? AND database_name = ? AND table_schema = ? AND storage_mode = ? AND layout_version = ?`,
		scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode, scope.LayoutVersion)
	require.NoError(t, err)

	fallback, _, err = recovery.mysqlConfirmedFallbackLayout(t.Context(), scope, table, owner)

	require.NoError(t, err)
	require.Nil(t, fallback, "a confirmed layout whose stored digest changed is not a fallback")
}

func TestMySQLRestoreTargetMapsForbiddenCatalogAsPermissionNotMissing(t *testing.T) {
	probe := &mysqlStorageRestoreProbe{
		inspection: &dbtarget.MySQLManagedInspection{
			Inspection: &dbtarget.TableInspection{Status: dbtarget.TableInspectionForbidden, Schema: "gwtest", Table: "gw_long_x"},
			Evidence:   &dbtarget.MySQLManagedTableEvidence{ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Database: "gwtest"},
		},
		inspectErr: dbtarget.ErrMySQLManagedTableForbidden,
	}
	snapshot := mysqlRuntimeSnapshot(t)
	target := mysqlStorageRestoreTargetWithOps(groupdelivery.Target{
		DB: &sql.DB{}, Kind: schema.DatabaseConnectorKindMySQL,
		MySQLCapability: mysqlRuntimeCapability(snapshot.RuntimeLayout.MySQL),
	}, "gwtest", mysqlStorageRestoreTargetOps{Inspect: probe.inspect, Execute: probe.execute})

	state, err := target.Inspect(t.Context(), "gw_long_x")

	require.ErrorIs(t, err, recordingplan.ErrTargetPermissionDenied)
	require.ErrorIs(t, err, dbtarget.ErrMySQLManagedTableForbidden)
	require.Equal(t, "forbidden", state.Status)
	require.NotNil(t, state.MySQL)
	require.False(t, errors.Is(err, dbtarget.ErrMySQLManagedTableUnconfirmed))
}
