package workspace

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

type mysqlTestWriteLayoutReader interface {
	MySQLTestWriteLayout(context.Context, string, string) (*WriteGroupRuntimeLayout, error)
}

func TestMySQLGroupTestWriteLayoutRequiresConfirmedSchema(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/test_write?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "test-write", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, group)
	reader, ok := any(h.service).(mysqlTestWriteLayoutReader)
	require.True(t, ok, "the explicit test-write boundary must verify canonical saved schema without authorizing Apply")
	for _, revision := range []string{"", "stale", saved.Revision} {
		layout, err := reader.MySQLTestWriteLayout(t.Context(), saved.ID, revision)
		require.Error(t, err, "an unconfirmed destination cannot authorize a test row")
		require.Nil(t, layout)
	}
	layout, err := reader.MySQLTestWriteLayout(t.Context(), "foreign", saved.Revision)
	require.Error(t, err)
	require.Nil(t, layout)
	ready, err := h.service.Readiness(t.Context(), saved.ID)
	require.NoError(t, err)
	require.False(t, ready.Ready)
}

func TestMySQLGroupTestWriteLayoutUsesAuthorizedCurrentPartitions(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			ensurePartitionMaintenanceTable(t, h.db)
			h.service.now = func() time.Time { return time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC) }
			connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
			_, err := h.workspace.BindDatabaseConnector(t.Context(), connectorID)
			require.NoError(t, err)
			group := h.longDraft(WriteGroupProfileInstantaneous, "test-write-current", h.primary())
			group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
			group.StorageMode = mode
			saved := mustEnsureMySQLGroup(t, h, group)
			var initial, current *recordingplan.MySQLPartitionEvidence
			if mode == WriteGroupStorageModeNativeMonthly {
				initial = initialNativePartitions(t)
				current = expandedNativePartitions(t, initial)
			}
			_, owner, ok := LongStorageIdentity(saved)
			require.True(t, ok)
			layout, err := recordingplan.ManagedMySQLLongGroupLayout(saved.Destination.TableName, owner, recordingplan.MySQLGroupSchemaOptions{
				ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1",
				Engine: "InnoDB", Collation: "utf8mb4_0900_bin", Charset: "utf8mb4", NativeMonthly: initial,
			})
			require.NoError(t, err)
			scope := StorageScopeOf(saved)
			baseDigest := seedConfirmedStorageScope(t, h.db, scope, layout, "op-confirm-test-write")
			if current != nil {
				seedPartitionMaintenanceState(t, h.db, scope, baseDigest, saved.Destination.ConnectorRevision, 1, current, 1, "op-maint-test-write", "ready", "", h.service.now())
			}
			physical := makeMySQLPhysical(saved.Destination.TableName, owner, layout, current)
			receipt := makeMySQLReceiptPhysical()
			initialInspection := recordingPlanMySQLInspection(physical)
			initialInspection.MySQL.Partitions = initial
			digest, err := recordingplan.VerifyManagedMySQLCatalog(layout, saved.Destination.Database, initialInspection, recordingPlanMySQLInspection(receipt))
			require.NoError(t, err)
			_, err = h.db.ExecContext(t.Context(), `UPDATE write_groups SET destination_schema_revision='schema-test-write', destination_schema_digest=? WHERE id=?`, digest, saved.ID)
			require.NoError(t, err)
			_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_digest=? WHERE workspace_id=? AND connector_id=?`, digest, saved.WorkspaceID, connectorID)
			require.NoError(t, err)
			h.service.WithManagedTableInspector(&fakeMySQLManagedRuntimeInspector{
				capability: &dbtarget.MySQLRecordingCapability{ConnectorID: connectorID, ConnectorRevision: saved.Destination.ConnectorRevision, Database: "partition_test",
					ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1", Connected: true, MetadataStatus: "verified", TransactionalEngine: true,
					SupportsManagedLong: true, SupportsNativePartition: true, SupportsGroupTestWrite: true, SupportsMissingRepair: true},
				group: physical, receipt: receipt,
			})
			before, err := recordingplan.ReadStorageLineage(t.Context(), h.db, scope)
			require.NoError(t, err)
			proof, err := h.service.MySQLTestWriteLayout(t.Context(), saved.ID, saved.Revision)
			require.NoError(t, err)
			require.Equal(t, digest, proof.SchemaDigest)
			if mode == WriteGroupStorageModeNativeMonthly {
				require.Len(t, proof.MySQL.Group.Evidence.Partitions.Partitions, 7)
				require.Len(t, before.Layout.MySQL.NativeMonthly.Partitions, 6)
			} else {
				require.Nil(t, proof.MySQL.Group.Evidence.Partitions)
			}
			after, err := recordingplan.ReadStorageLineage(t.Context(), h.db, scope)
			require.NoError(t, err)
			require.Equal(t, before, after)
			ready, err := h.service.Readiness(t.Context(), saved.ID)
			require.NoError(t, err)
			require.True(t, ready.SchemaReady)
			require.True(t, ready.Ready, "verified MySQL operations must open readiness without changing test-write behavior")
			stale, err := h.service.MySQLTestWriteLayout(t.Context(), saved.ID, "stale")
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			require.Nil(t, stale)
			if current != nil {
				_, err := h.db.ExecContext(t.Context(), `DELETE FROM write_group_mysql_partition_maintenance WHERE workspace_id=?`, saved.WorkspaceID)
				require.NoError(t, err)
				seedPartitionMaintenanceState(t, h.db, scope, baseDigest, saved.Destination.ConnectorRevision, 2, current, 2, "op-stale", "blocked", "permission_lost", h.service.now())
				invalid, err := h.service.MySQLTestWriteLayout(t.Context(), saved.ID, saved.Revision)
				require.Error(t, err, "stale generation cannot authorize a physical partition overlay")
				require.Nil(t, invalid)
			}
		})
	}
}
