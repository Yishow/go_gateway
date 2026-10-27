package workspace

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeSQLServerProbeInspector struct {
	capability *dbtarget.SQLServerRecordingCapability
}

func (f *fakeSQLServerProbeInspector) InspectTableAtRevision(context.Context, string, string, string, string) (*dbtarget.TableInspection, error) {
	return &dbtarget.TableInspection{Status: "missing"}, nil
}

func (f *fakeSQLServerProbeInspector) ProbeSQLServerRecordingAtRevision(context.Context, string, string) (*dbtarget.SQLServerRecordingCapability, error) {
	return f.capability, nil
}

func seedSQLServerConnector(t *testing.T, db *sql.DB, configJSON string) string {
	t.Helper()
	id := uuid.NewString()
	_, err := db.ExecContext(t.Context(), `
		INSERT INTO database_connectors (id, name, kind, connection_config, identity_revision, status, enabled)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, "SQL Server Target "+id[:8], "sqlserver", configJSON, "connector-revision-sqlserver", schema.DatabaseConnectorStatusReady, true)
	require.NoError(t, err)
	return id
}

func sqlServerTestCapability(connectorID, revision string) *dbtarget.SQLServerRecordingCapability {
	return &dbtarget.SQLServerRecordingCapability{
		ConnectorID: connectorID, ConnectorRevision: revision, Database: "recordings", Schema: "dbo",
		ServerVersion: "16.0.1000.6", ServerEdition: "Developer Edition (64-bit)", EngineEdition: 3,
		CompatibilityLevel: 160, CapabilityRevision: "sqlservercap-native-test",
		Connected: true, MetadataStatus: "verified", SupportsManagedLong: true, SupportsNativePartition: true,
		Permissions: dbtarget.SQLServerRecordingPermissions{Verified: true, CreateTable: true, AlterSchema: true, SelectSchema: true, InsertSchema: true, DeleteSchema: true},
	}
}

func ensureSQLServerGroup(t *testing.T, h *circuitProfileHarness, group *WriteGroup) (*WriteGroupSaveResult, error) {
	t.Helper()
	current, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	return h.service.EnsureCircuitProfile(t.Context(), WriteGroupMutation{
		WorkspaceID: h.fixture.workspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedConnectorRevision: group.Destination.ConnectorRevision, Group: group,
	})
}

func TestSQLServerNativeMonthlySchemaScopeKeepsBoundaries(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedSQLServerConnector(t, h.db, `{"host":"sql.example","database":"recordings","schema":"dbo","username":"u","password":"p"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "sqlserver-native", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-sqlserver"
	group.Destination.Database, group.Destination.TableSchema = "", ""
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	saved, err := ensureSQLServerGroup(t, h, group)
	require.NoError(t, err)
	require.NotNil(t, saved.Group)
	require.Equal(t, WriteGroupStorageModeNativeMonthly, saved.Group.StorageMode)

	capability := sqlServerTestCapability(connectorID, saved.Group.Destination.ConnectorRevision)
	h.service.WithManagedTableInspector(&fakeSQLServerProbeInspector{capability: capability})
	h.service.now = func() time.Time { return time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC) }

	current, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	mutation := WriteGroupMutation{
		WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: current.DatabaseSetupRevision,
		ExpectedGroupRevision: saved.Group.Revision, ExpectedConnectorRevision: saved.Group.Destination.ConnectorRevision, Group: saved.Group,
	}
	scope, err := h.service.ManagedSchemaScope(t.Context(), saved.Group.ID, mutation)
	require.NoError(t, err)
	require.NotNil(t, scope.GroupLayout.SQLServer.NativeMonthly)
	partitions := scope.GroupLayout.SQLServer.NativeMonthly.Partitions
	require.Len(t, partitions, 6)
	require.Equal(t, "2026-10-01T00:00:00.000000", partitions[0].LessThan)
	require.Equal(t, "2027-02-01T00:00:00.000000", partitions[4].LessThan)
	require.Equal(t, "MAXVALUE", partitions[5].LessThan)

	// Freeze the confirmation and advance the clock. A new preview must retain
	// the confirmed boundaries until scoped maintenance changes them.
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), h.db, recordingplan.StorageScope{
		WorkspaceID: saved.Group.WorkspaceID, ConnectorID: saved.Group.Destination.ConnectorID,
		Database: saved.Group.Destination.Database, TableSchema: saved.Group.Destination.TableSchema,
		StorageMode: string(WriteGroupStorageModeNativeMonthly), LayoutVersion: string(WriteGroupLayoutManagedLongV1),
	}, recordingplan.StorageConfirmation{
		Layout: *scope.GroupLayout, ConnectorRevision: saved.Group.Destination.ConnectorRevision,
		OperationID: "native-confirmed", VerifiedDigest: strings.Repeat("a", 64), CreatedTable: false, GroupID: saved.Group.ID,
	}, time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC)))
	h.service.now = func() time.Time { return time.Date(2027, time.April, 25, 0, 0, 0, 0, time.UTC) }
	repeated, err := h.service.ManagedSchemaScope(t.Context(), saved.Group.ID, mutation)
	require.NoError(t, err)
	require.Equal(t, partitions, repeated.GroupLayout.SQLServer.NativeMonthly.Partitions)

	capability.SupportsNativePartition = false
	_, err = h.service.ManagedSchemaScope(t.Context(), saved.Group.ID, mutation)
	require.ErrorIs(t, err, recordingplan.ErrSQLServerPartitionUnavailable)

	// A separate plain group remains selectable and does not inherit the
	// unavailable partition branch.
	plainSource := h.seedSource(h.fixture.deviceID, "sqlserver-plain-source")
	plainDraft := h.longDraft(WriteGroupProfileInstantaneous, "sqlserver-plain", plainSource)
	plainDraft.Destination.ConnectorID, plainDraft.Destination.ConnectorRevision = connectorID, "connector-revision-sqlserver"
	plainDraft.Destination.Database, plainDraft.Destination.TableSchema = "", ""
	plainDraft.StorageMode = WriteGroupStorageModePlain
	plainSaved, err := ensureSQLServerGroup(t, h, plainDraft)
	require.NoError(t, err)
	plainCurrent, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	plainScope, err := h.service.ManagedSchemaScope(t.Context(), plainSaved.Group.ID, WriteGroupMutation{
		WorkspaceID: plainSaved.Group.WorkspaceID, ExpectedWorkspaceRevision: plainCurrent.DatabaseSetupRevision,
		ExpectedGroupRevision: plainSaved.Group.Revision, ExpectedConnectorRevision: plainSaved.Group.Destination.ConnectorRevision, Group: plainSaved.Group,
	})
	require.NoError(t, err)
	require.Nil(t, plainScope.GroupLayout.SQLServer.NativeMonthly)
}
