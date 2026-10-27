package workspace

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

type sqlServerTestWriteLayoutReader interface {
	SQLServerTestWriteLayout(context.Context, string, string) (*WriteGroupRuntimeLayout, error)
}

func TestSQLServerTestWriteCapabilityRequiresExplicitFlags(t *testing.T) {
	group := &WriteGroup{
		StorageMode: WriteGroupStorageModePlain,
		Destination: WriteGroupDestination{ConnectorID: "conn-sql", ConnectorRevision: "rev-1", Database: "analytics", TableSchema: "dbo"},
	}
	evidence := &dbtarget.SQLServerManagedTableEvidence{
		CapabilityRevision: "cap-1", ServerVersion: "16.0.1", ServerEdition: "Standard", EngineEdition: 2, CompatibilityLevel: 160,
	}
	layout := &WriteGroupRuntimeLayout{SQLServer: &dbtarget.SQLServerManagedRowProof{
		Group: &dbtarget.SQLServerManagedInspection{Evidence: evidence},
	}}
	capability := &dbtarget.SQLServerRecordingCapability{
		ConnectorID: "conn-sql", ConnectorRevision: "rev-1", Database: "analytics", Schema: "dbo",
		CapabilityRevision: "cap-1", ServerVersion: "16.0.1", ServerEdition: "Standard", EngineEdition: 2, CompatibilityLevel: 160,
		Connected: true, MetadataStatus: "verified", SupportsManagedLong: true,
	}
	require.False(t, sqlServerTestWriteCapabilityMatches(group, layout, capability), "managed-long support alone must not authorize test writes")
	capability.SupportsGroupTestWrite = true
	require.True(t, sqlServerTestWriteCapabilityMatches(group, layout, capability))
	group.StorageMode = WriteGroupStorageModeNativeMonthly
	require.False(t, sqlServerTestWriteCapabilityMatches(group, layout, capability), "native-monthly also requires its capability")
	capability.SupportsNativePartition = true
	require.True(t, sqlServerTestWriteCapabilityMatches(group, layout, capability))
}

func TestSQLServerTestWriteLayoutRequiresConfirmedSchema(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	connectorID := seedSQLServerConnector(t, h.db, `{"host":"sql.example","database":"recordings","schema":"dbo","username":"u","password":"p"}`)
	group := h.longDraft(WriteGroupProfileInstantaneous, "sqlserver-test-write", h.primary())
	group.Destination.ConnectorID, group.Destination.ConnectorRevision = connectorID, "connector-revision-sqlserver"
	group.Destination.Database, group.Destination.TableSchema = "", ""
	saved, err := ensureSQLServerGroup(t, h, group)
	require.NoError(t, err)
	h.service.WithManagedTableInspector(&fakeSQLServerProbeInspector{
		capability: sqlServerTestCapability(connectorID, saved.Group.Destination.ConnectorRevision),
	})

	reader, ok := any(h.service).(sqlServerTestWriteLayoutReader)
	require.True(t, ok, "SQL Server test-write must use an explicit saved-scope workspace adapter")
	for _, revision := range []string{"", "stale"} {
		layout, err := reader.SQLServerTestWriteLayout(t.Context(), saved.Group.ID, revision)
		require.Error(t, err, "missing or stale revisions cannot authorize test rows")
		require.Nil(t, layout)
	}
	layout, err := reader.SQLServerTestWriteLayout(t.Context(), "foreign", saved.Group.Revision)
	require.Error(t, err)
	require.Nil(t, layout)
	layout, err = reader.SQLServerTestWriteLayout(t.Context(), saved.Group.ID, saved.Group.Revision)
	require.Error(t, err, "an unconfirmed destination cannot authorize a test row")
	require.Nil(t, layout)
}
