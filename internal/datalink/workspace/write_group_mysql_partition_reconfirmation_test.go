package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go-gateway/internal/datalink/recordingplan"
)

func TestNormalizeMySQLPartitionInspectionAllowsExplicitCapabilityReconfirmation(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "native-reconfirm", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	saved := mustEnsureMySQLGroup(t, h, group)
	scope := longStorageScopeOf(saved)
	storageScope := StorageScopeOf(saved)
	initial := initialNativePartitions(t)
	expanded := expandedNativePartitions(t, initial)

	oldLayout, err := recordingplan.ManagedMySQLLongGroupLayout(scope.tableName(), scope.ownerColumn(), recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-1", Engine: "InnoDB",
		Collation: "utf8mb4_0900_bin", NativeMonthly: initial,
	})
	require.NoError(t, err)
	baseDigest := seedConfirmedStorageScope(t, h.db, storageScope, oldLayout, "op-old")
	seedPartitionMaintenanceState(t, h.db, storageScope, baseDigest, "connector-revision-mysql", 1, expanded, 1, "op-maint", "ready", "", time.Now())

	newLayout, err := recordingplan.ManagedMySQLLongGroupLayout(scope.tableName(), scope.ownerColumn(), recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "cap-2", Engine: "InnoDB",
		Collation: "utf8mb4_0900_bin", NativeMonthly: initial,
	})
	require.NoError(t, err)
	previewScope := recordingplan.SchemaPreviewScope{
		WorkspaceID: saved.WorkspaceID, ConnectorID: saved.Destination.ConnectorID,
		ConnectorRevision: saved.Destination.ConnectorRevision, Database: saved.Destination.Database,
		Schema: saved.Destination.TableSchema, GroupLayout: &newLayout,
	}
	actual := recordingPlanMySQLInspection(makeMySQLPhysical(scope.tableName(), scope.ownerColumn(), newLayout, expanded))
	normalized, err := h.service.NormalizeMySQLPartitionInspection(t.Context(), previewScope, actual)
	require.NoError(t, err)
	require.Equal(t, initial, normalized.MySQL.Partitions)
	require.Equal(t, "cap-2", normalized.MySQL.CapabilityRevision)
	state, err := recordingplan.ReadMySQLPartitionMaintenanceState(t.Context(), h.db, storageScope)
	require.NoError(t, err)
	require.Equal(t, expanded, state.Current, "normalization does not mutate the authorized current sidecar")
}
