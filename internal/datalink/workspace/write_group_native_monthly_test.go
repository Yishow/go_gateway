package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func TestMySQLNativeMonthlySelectionIsExplicitAndScoped(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "partition-circuit", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-mysql"
	group.Destination.Database, group.Destination.TableSchema = "", ""
	group.StorageMode = WriteGroupStorageMode("native_monthly")
	saved, err := ensureMySQLGroup(t, h, group)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStorageMode("native_monthly"), saved.Group.StorageMode)
	reloaded, err := NewSQLWriteGroupRepository(h.db).Get(t.Context(), h.fixture.workspaceID, saved.Group.ID)
	require.NoError(t, err)
	require.Equal(t, saved.Group.StorageMode, reloaded.StorageMode)
	plain := *saved.Group
	plain.StorageMode = WriteGroupStorageModePlain
	plainTable, _, ok := LongStorageIdentity(&plain)
	require.True(t, ok)
	require.NotEqual(t, plainTable, saved.Group.Destination.TableName, "plain and monthly scopes cannot share an owned table")
}

func TestMySQLPartitionCapabilityGate(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connector := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "native-preview", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connector, "connector-revision-mysql"
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	saved := mustEnsureMySQLGroup(t, h, group)
	capability := &dbtarget.MySQLRecordingCapability{
		ConnectorID: connector, ConnectorRevision: saved.Destination.ConnectorRevision, Database: "partition_test",
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "partition-cap-1",
		Connected: true, MetadataStatus: "verified", TransactionalEngine: true,
	}
	h.service.WithManagedTableInspector(&fakeMySQLProbeInspector{capability: capability})
	h.service.now = func() time.Time { return time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC) }
	current, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	mutation := WriteGroupMutation{
		WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: saved.Revision, ExpectedConnectorRevision: saved.Destination.ConnectorRevision, Group: saved,
	}
	scope, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, mutation)
	require.NoError(t, err)
	require.NotNil(t, scope.GroupLayout.MySQL.NativeMonthly)
	partitions := scope.GroupLayout.MySQL.NativeMonthly.Partitions
	require.Len(t, partitions, 6)
	require.Equal(t, "2026-10-01T00:00:00Z", partitions[0].LessThan)
	require.Equal(t, "2027-02-01T00:00:00Z", partitions[4].LessThan)
	require.Equal(t, "MAXVALUE", partitions[5].LessThan)
	for _, invalid := range []func(){
		func() { capability.MetadataStatus = "forbidden" },
		func() { capability.ServerVersion = "8.0.1" },
		func() { capability.TransactionalEngine = false },
	} {
		original := *capability
		invalid()
		_, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, mutation)
		require.ErrorIs(t, err, ErrWriteGroupValidation)
		*capability = original
	}

	// Freeze the local confirmation bytes, then advance the clock. Readiness and
	// a new preview must keep the confirmed boundaries until scoped maintenance.
	encoded, err := json.Marshal(scope.GroupLayout)
	require.NoError(t, err)
	hash := sha256.Sum256(encoded)
	storage := longStorageScopeOf(saved)
	_, err = h.db.ExecContext(t.Context(), `UPDATE write_group_storage_scopes
		SET confirmed_layout = ?, layout_digest = ?, confirmed_operation_id = 'native-confirmed'
		WHERE workspace_id = ? AND connector_id = ? AND storage_mode = 'native_monthly'`,
		string(encoded), hex.EncodeToString(hash[:]), storage.WorkspaceID, storage.ConnectorID)
	require.NoError(t, err)
	h.service.now = func() time.Time { return time.Date(2027, 4, 25, 0, 0, 0, 0, time.UTC) }
	repeated, err := h.service.ManagedSchemaScope(t.Context(), saved.ID, mutation)
	require.NoError(t, err)
	require.Equal(t, partitions, repeated.GroupLayout.MySQL.NativeMonthly.Partitions)
}

func TestNativeMonthlyRefusesOtherConnectorKinds(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	group := h.longDraft(WriteGroupProfileInstantaneous, "sqlite-partition", h.primary())
	group.StorageMode = WriteGroupStorageMode("native_monthly")
	_, err := h.ensure(group)
	require.ErrorIs(t, err, ErrWriteGroupValidation, "explicit native MySQL partitioning cannot silently become SQLite plain storage")
}

func TestMySQLNativeMonthlyUnconfirmedScopeIdentity(t *testing.T) {
	for _, field := range []string{"table_name", "owner_column"} {
		t.Run(field, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			connector := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
			group := h.longDraft(WriteGroupProfileInstantaneous, "native-owner", h.primary())
			group.Destination.ConnectorID, group.Destination.ConnectorRevision = connector, "connector-revision-mysql"
			group.StorageMode = WriteGroupStorageModeNativeMonthly
			saved := mustEnsureMySQLGroup(t, h, group)
			// The field comes only from the fixed table above, never an API input.
			_, err := h.db.ExecContext(t.Context(), "UPDATE write_group_storage_scopes SET "+field+" = 'foreign' WHERE workspace_id = ?", saved.WorkspaceID)
			require.NoError(t, err)
			_, err = h.service.mysqlNativeMonthlyLayout(t.Context(), saved)
			require.ErrorIs(t, err, ErrWriteGroupValidation, "unconfirmed preview must reject a corrupted storage identity before generating DDL")
		})
	}
}
