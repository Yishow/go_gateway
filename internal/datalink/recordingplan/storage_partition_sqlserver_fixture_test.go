package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type sqlServerPartitionMaintenanceCall struct {
	before *SQLServerPartitionEvidence
	after  *SQLServerPartitionEvidence
}

type sqlServerPartitionMaintenanceTargetFixture struct {
	layout      GroupSchemaLayout
	group       TargetTableInspection
	receipt     TargetTableInspection
	partitions  *SQLServerPartitionEvidence
	inspectErr  error
	splitErr    error
	applyBefore bool
	calls       []sqlServerPartitionMaintenanceCall
}

func (f *sqlServerPartitionMaintenanceTargetFixture) target() SQLServerPartitionMaintenanceTarget {
	return SQLServerPartitionMaintenanceTarget{
		Inspect: func(_ context.Context, table string) (TargetTableInspection, error) {
			if f.inspectErr != nil {
				return TargetTableInspection{}, f.inspectErr
			}
			switch table {
			case f.layout.TableName:
				state := cloneSQLServerTargetInspectionForPartitionTest(f.group)
				state.SQLServer.NativeMonthly = sqlServerTargetPartitionFromPlan(f.partitions)
				return state, nil
			case managedEffectReceiptTable:
				return cloneSQLServerTargetInspectionForPartitionTest(f.receipt), nil
			default:
				return TargetTableInspection{}, ErrTargetInspectionUnconfirmed
			}
		},
		Split: func(_ context.Context, before, after *SQLServerPartitionEvidence) error {
			f.calls = append(f.calls, sqlServerPartitionMaintenanceCall{
				before: CloneSQLServerPartitionEvidence(before), after: CloneSQLServerPartitionEvidence(after),
			})
			if f.splitErr != nil && !f.applyBefore {
				return f.splitErr
			}
			f.partitions = CloneSQLServerPartitionEvidence(after)
			return f.splitErr
		},
	}
}

type sqlServerPartitionMaintenanceFixture struct {
	db      *sql.DB
	service *Service
	request SQLServerPartitionMaintenanceRequest
	base    *SQLServerPartitionEvidence
	layout  GroupSchemaLayout
	target  *sqlServerPartitionMaintenanceTargetFixture
	repo    *SQLRepository
}

func newSQLServerPartitionMaintenanceFixture(t *testing.T) *sqlServerPartitionMaintenanceFixture {
	t.Helper()
	db := lineageDB(t)
	applyRecordingDDL(t, db)
	migration, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "042_sqlserver_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)

	scope := StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "dbo",
		StorageMode: "native_monthly", LayoutVersion: "managed_long_v1"}
	base, err := NewSQLServerNativeMonthlyPartition(time.Date(2026, time.October, 15, 0, 0, 0, 0, time.UTC),
		"pf_"+"gw_long_"+mysqlStorageIdentityHash(scope, "table")[:16],
		"ps_"+"gw_long_"+mysqlStorageIdentityHash(scope, "table")[:16])
	require.NoError(t, err)
	options := SQLServerGroupSchemaOptions{
		Database: scope.Database, Schema: scope.TableSchema, Collation: "Latin1_General_100_BIN2", ServerVersion: "16.0.1000", ServerEdition: "Standard Edition",
		CapabilityRevision: "cap-1", EngineEdition: 2, CompatibilityLevel: 160, NativeMonthly: base,
	}
	layout, err := ManagedSQLServerLongGroupLayout("gw_long_"+mysqlStorageIdentityHash(scope, "table")[:16],
		"_gw_owner_"+mysqlStorageIdentityHash(scope, "owner")[:24], options)
	require.NoError(t, err)
	group := sqlServerRestoreCatalog(layout.Columns, layout.Indexes, layout.TableName, layout.OwnerColumn, 401, options)
	group.SQLServer.NativeMonthly = sqlServerTargetPartitionFromPlan(base)
	group.SQLServer.PartitionInspected = true
	group.SQLServer.Plain = false
	receipt := sqlServerRestoreReceiptCatalog(402, options)
	digest, err := VerifyManagedSQLServerCatalog(layout, scope.Database, group, receipt)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema, scope.StorageMode,
		scope.LayoutVersion, layout.TableName, layout.OwnerColumn, UnconfirmedSQLServerStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, RecordStorageConfirmation(t.Context(), db, scope, StorageConfirmation{
		Layout: layout, ConnectorRevision: "rev-1", OperationID: "op-confirm", VerifiedDigest: digest, CreatedTable: true, GroupID: "group-1",
	}, time.Now().UTC()))
	repo := NewSQLRepository(db)
	require.NoError(t, repo.ensureSQLServerPartitionMaintenanceState(t.Context(), &SQLServerPartitionMaintenanceState{
		Scope: scope, BaseLayoutDigest: digestGroupSchemaLayout(layout), ConnectorRevision: "rev-1", Generation: 1,
		Current: CloneSQLServerPartitionEvidence(base), Revision: 1, Status: serverMaintenanceReadyForTest(), UpdatedAt: time.Now().UTC(),
	}))
	return &sqlServerPartitionMaintenanceFixture{
		db: db, repo: repo, service: NewService(repo), layout: layout, base: base,
		request: SQLServerPartitionMaintenanceRequest{Scope: scope, GroupID: "group-1", ConnectorRevision: "rev-1",
			At: time.Date(2026, time.December, 20, 0, 0, 0, 0, time.UTC)},
		target: &sqlServerPartitionMaintenanceTargetFixture{layout: layout, group: group, receipt: receipt, partitions: CloneSQLServerPartitionEvidence(base)},
	}
}

func serverMaintenanceReadyForTest() string { return sqlServerPartitionMaintenanceReady }

func cloneSQLServerTargetInspectionForPartitionTest(state TargetTableInspection) TargetTableInspection {
	cloned := state
	cloned.Columns = append([]string(nil), state.Columns...)
	cloned.Indexes = append([]TargetIndexInspection(nil), state.Indexes...)
	if state.SQLServer != nil {
		evidence := *state.SQLServer
		evidence.PrimaryKey = append([]string(nil), state.SQLServer.PrimaryKey...)
		evidence.ColumnDefaults = cloneStringPtrMap(state.SQLServer.ColumnDefaults)
		evidence.ColumnCollations = map[string]string{}
		for key, value := range state.SQLServer.ColumnCollations {
			evidence.ColumnCollations[key] = value
		}
		if state.SQLServer.NativeMonthly != nil {
			native := *state.SQLServer.NativeMonthly
			native.Partitions = append([]SQLServerPartitionBoundary(nil), state.SQLServer.NativeMonthly.Partitions...)
			evidence.NativeMonthly = &native
		}
		cloned.SQLServer = &evidence
	}
	return cloned
}

func sqlServerTargetPartitionFromPlan(plan *SQLServerPartitionEvidence) *TargetSQLServerPartitionEvidence {
	if plan == nil {
		return nil
	}
	return &TargetSQLServerPartitionEvidence{
		Method: plan.Method, Column: plan.Column, FunctionName: plan.FunctionName, SchemeName: plan.SchemeName,
		Filegroup: plan.Filegroup, Partitions: append([]SQLServerPartitionBoundary(nil), plan.Partitions...),
		FunctionExists: true, SchemeExists: true, Exclusive: true, IndexesAligned: true,
	}
}

func (f *sqlServerPartitionMaintenanceFixture) seedStaleOperation(t *testing.T, before, after *SQLServerPartitionEvidence) *SchemaOperation {
	t.Helper()
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	state, err := ReadSQLServerPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	detail, err := json.Marshal(sqlServerPartitionMaintenanceDetail{
		Kind: sqlServerPartitionMaintenanceKind, Scope: f.request.Scope, GroupID: f.request.GroupID,
		ConnectorRevision: f.request.ConnectorRevision, BaseLayoutDigest: lineage.LayoutDigest,
		Generation: lineage.Generation, StateRevision: state.Revision, At: f.request.At.UTC(),
		Before: CloneSQLServerPartitionEvidence(before), After: CloneSQLServerPartitionEvidence(after),
	})
	require.NoError(t, err)
	now := time.Now().UTC().Add(-2 * SchemaOperationLease)
	op := &SchemaOperation{
		OperationID: "stale-sqlserver-partition-op", Token: "stale-sqlserver-partition-token", WorkspaceID: f.request.Scope.WorkspaceID,
		ScopeKey: sqlServerPartitionMaintenanceScopeKey(f.request.Scope, lineage.Layout), Owner: "claim-stale-owner",
		Action: SQLServerPartitionMaintenanceAction, Status: SchemaOperationRunning, Detail: string(detail),
		CreatedAt: now, UpdatedAt: now,
	}
	claimed, outcome, err := f.repo.ClaimSchemaOperation(t.Context(), op)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	return claimed
}
