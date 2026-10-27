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

type partitionMaintenanceCall struct {
	before *MySQLPartitionEvidence
	after  *MySQLPartitionEvidence
}

type mysqlPartitionMaintenanceTargetFixture struct {
	layout           GroupSchemaLayout
	group            TargetTableInspection
	receipt          TargetTableInspection
	partitions       *MySQLPartitionEvidence
	after            *MySQLPartitionEvidence
	inspectErr       error
	reorganizeErr    error
	applyBeforeError bool
	reorganizations  []partitionMaintenanceCall
}

func (f *mysqlPartitionMaintenanceTargetFixture) target() MySQLPartitionMaintenanceTarget {
	return MySQLPartitionMaintenanceTarget{
		Inspect: func(_ context.Context, table string) (TargetTableInspection, error) {
			if f.inspectErr != nil {
				return TargetTableInspection{}, f.inspectErr
			}
			switch table {
			case f.layout.TableName:
				state := cloneTargetTableInspection(f.group)
				state.MySQL.Partitions = cloneMySQLPartitionEvidence(f.partitions)
				return state, nil
			case managedEffectReceiptTable:
				return cloneTargetTableInspection(f.receipt), nil
			default:
				return TargetTableInspection{}, ErrTargetInspectionUnconfirmed
			}
		},
		Reorganize: func(_ context.Context, before, after *MySQLPartitionEvidence) error {
			f.reorganizations = append(f.reorganizations, partitionMaintenanceCall{
				before: cloneMySQLPartitionEvidence(before), after: cloneMySQLPartitionEvidence(after),
			})
			if f.applyBeforeError {
				observed := after
				if f.after != nil {
					observed = f.after
				}
				f.partitions = cloneMySQLPartitionEvidence(observed)
			}
			return f.reorganizeErr
		},
	}
}

type mysqlPartitionMaintenanceFixture struct {
	db      *sql.DB
	service *Service
	request MySQLPartitionMaintenanceRequest
	base    *MySQLPartitionEvidence
	target  *mysqlPartitionMaintenanceTargetFixture
	repo    *SQLRepository
}

func newMySQLPartitionMaintenanceFixture(t *testing.T) *mysqlPartitionMaintenanceFixture {
	t.Helper()
	db := lineageDB(t)
	applyRecordingDDL(t, db)
	migration, err := os.ReadFile(filepath.Join("..", "schema", "migrations", "041_mysql_partition_maintenance_sqlite.up.sql"))
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), string(migration))
	require.NoError(t, err)

	scope := StorageScope{WorkspaceID: "ws", ConnectorID: "conn", Database: "recordings", TableSchema: "recordings",
		StorageMode: "native_monthly", LayoutVersion: "managed_long_v1"}
	base, err := NewMySQLNativeMonthlyPartition(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	layout, err := ManagedMySQLLongGroupLayout("gw_long_"+mysqlStorageIdentityHash(scope, "table")[:16],
		"_gw_owner_"+mysqlStorageIdentityHash(scope, "owner")[:24], MySQLGroupSchemaOptions{
			ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB",
			Collation: mysqlTestCollation, NativeMonthly: base,
		})
	require.NoError(t, err)
	group := mysqlRestoreCatalog(layout.Columns, layout.Indexes, layout.OwnerColumn, *layout.MySQL)
	group.MySQL.Partitions = cloneMySQLPartitionEvidence(base)
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
	}, time.Now().UTC()))

	repo := NewSQLRepository(db)
	require.NoError(t, repo.ensureMySQLPartitionMaintenanceState(t.Context(), &MySQLPartitionMaintenanceState{
		Scope: scope, BaseLayoutDigest: digestGroupSchemaLayout(layout), ConnectorRevision: "rev-1", Generation: 1,
		Current: cloneMySQLPartitionEvidence(base), Revision: 1, Status: "ready", UpdatedAt: time.Now().UTC(),
	}))
	return &mysqlPartitionMaintenanceFixture{
		db: db, repo: repo, service: NewService(repo),
		request: MySQLPartitionMaintenanceRequest{Scope: scope, GroupID: "group-1", ConnectorRevision: "rev-1", At: time.Date(2024, time.February, 3, 0, 0, 0, 0, time.UTC)},
		base:    base, target: &mysqlPartitionMaintenanceTargetFixture{layout: layout, group: group, receipt: receipt, partitions: cloneMySQLPartitionEvidence(base)},
	}
}

func (f *mysqlPartitionMaintenanceFixture) seedStaleOperation(t *testing.T, before, after *MySQLPartitionEvidence) *SchemaOperation {
	t.Helper()
	lineage, err := ReadStorageLineage(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	state, err := ReadMySQLPartitionMaintenanceState(t.Context(), f.db, f.request.Scope)
	require.NoError(t, err)
	detail, err := json.Marshal(mysqlPartitionMaintenanceDetail{
		Kind: mysqlPartitionMaintenanceKind, Scope: f.request.Scope, GroupID: f.request.GroupID,
		ConnectorRevision: f.request.ConnectorRevision, BaseLayoutDigest: lineage.LayoutDigest,
		Generation: lineage.Generation, StateRevision: state.Revision, At: f.request.At.UTC(),
		Before: cloneMySQLPartitionEvidence(before), After: cloneMySQLPartitionEvidence(after),
	})
	require.NoError(t, err)
	now := time.Now().UTC().Add(-2 * SchemaOperationLease)
	op := &SchemaOperation{
		OperationID: "stale-partition-op", Token: "stale-partition-token", WorkspaceID: f.request.Scope.WorkspaceID,
		ScopeKey: mysqlPartitionMaintenanceScopeKey(f.request.Scope, lineage.Layout), Owner: "claim-stale-owner",
		Action: MySQLPartitionMaintenanceAction, Status: SchemaOperationRunning, Detail: string(detail),
		CreatedAt: now, UpdatedAt: now,
	}
	claimed, outcome, err := f.repo.ClaimSchemaOperation(t.Context(), op)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	return claimed
}
