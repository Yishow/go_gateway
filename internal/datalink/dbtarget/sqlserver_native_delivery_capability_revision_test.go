package dbtarget

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeDeliveryCapabilityRevisionAllowsOnlyMaintenancePermissionLoss(t *testing.T) {
	frozen := sqlServerNativeDeliveryCapabilityFixture(true, nil, nil)
	current := sqlServerNativeDeliveryCapabilityFixture(false, nil, nil)

	require.True(t, SQLServerNativeDeliveryCapabilityRevisionMatches(frozen, frozen.CapabilityRevision))
	require.False(t, current.SupportsNativePartition, "maintenance permission loss does not enable or retain partition maintenance")
	require.True(t, SQLServerNativeDeliveryCapabilityRevisionMatches(current, frozen.CapabilityRevision))
	require.False(t, SQLServerNativeDeliveryCapabilityRevisionMatches(current, frozen.CapabilityRevision+"stale"))

	inconsistent := *current
	inconsistent.CapabilityRevision = frozen.CapabilityRevision
	require.False(t, SQLServerNativeDeliveryCapabilityRevisionMatches(&inconsistent, frozen.CapabilityRevision),
		"the revision stored on a capability must match its recomputed full hash")
}

func TestSQLServerNativeDeliveryCapabilityRevisionRejectsOtherDifferences(t *testing.T) {
	frozen := sqlServerNativeDeliveryCapabilityFixture(true, nil, nil)
	mutations := []struct {
		name string
		edit func(*schema.DatabaseConnector, *sqlServerRecordingMetadata)
	}{
		{name: "connector id", edit: func(connector *schema.DatabaseConnector, _ *sqlServerRecordingMetadata) {
			connector.ID = "other-connector"
		}},
		{name: "connector revision", edit: func(connector *schema.DatabaseConnector, _ *sqlServerRecordingMetadata) {
			connector.IdentityRevision = "other-revision"
		}},
		{name: "server version", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.serverVersion = "16.0.1001.1"
		}},
		{name: "server edition", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.serverEdition = "Standard Edition"
		}},
		{name: "engine edition", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) { metadata.engineEdition++ }},
		{name: "compatibility level", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) { metadata.compatibilityLevel-- }},
		{name: "database scope", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.database, metadata.actualDatabase = "other_db", "other_db"
		}},
		{name: "schema scope", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.schema, metadata.actualSchema = "other_schema", "other_schema"
		}},
		{name: "database visibility permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.DatabaseViewDefinition = false
		}},
		{name: "schema visibility permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.SchemaViewDefinition = false
		}},
		{name: "create permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.CreateTable = false
		}},
		{name: "alter schema permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.AlterSchema = false
		}},
		{name: "select permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.SelectSchema = false
		}},
		{name: "insert permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.InsertSchema = false
		}},
		{name: "delete permission", edit: func(_ *schema.DatabaseConnector, metadata *sqlServerRecordingMetadata) {
			metadata.permissions.DeleteSchema = false
		}},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			current := sqlServerNativeDeliveryCapabilityFixture(false, mutation.edit, nil)
			require.False(t, SQLServerNativeDeliveryCapabilityRevisionMatches(current, frozen.CapabilityRevision))
		})
	}
}

func TestSQLServerNativeDeliveryCapabilityRevisionRequiresVerifiedWritableScope(t *testing.T) {
	frozen := sqlServerNativeDeliveryCapabilityFixture(true, nil, nil)
	tests := []struct {
		name   string
		mutate func(*SQLServerRecordingCapability)
	}{
		{name: "disconnected", mutate: func(capability *SQLServerRecordingCapability) { capability.Connected = false }},
		{name: "metadata unverified", mutate: func(capability *SQLServerRecordingCapability) {
			capability.MetadataStatus = sqlServerRecordingMetadataUnconfirmed
		}},
		{name: "managed long unsupported", mutate: func(capability *SQLServerRecordingCapability) { capability.SupportsManagedLong = false }},
		{name: "permissions unverified", mutate: func(capability *SQLServerRecordingCapability) { capability.Permissions.Verified = false }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current := sqlServerNativeDeliveryCapabilityFixture(false, nil, nil)
			test.mutate(current)
			require.False(t, SQLServerNativeDeliveryCapabilityRevisionMatches(current, frozen.CapabilityRevision))
		})
	}

	for _, test := range []struct {
		name string
		edit func(*SQLServerRecordingPermissions)
	}{
		{name: "no metadata visibility", edit: func(permissions *SQLServerRecordingPermissions) {
			permissions.DatabaseViewDefinition, permissions.SchemaViewDefinition = false, false
		}},
		{name: "no select", edit: func(permissions *SQLServerRecordingPermissions) { permissions.SelectSchema = false }},
		{name: "no insert", edit: func(permissions *SQLServerRecordingPermissions) { permissions.InsertSchema = false }},
	} {
		t.Run(test.name, func(t *testing.T) {
			current := sqlServerNativeDeliveryCapabilityFixture(false, nil, test.edit)
			require.False(t, SQLServerNativeDeliveryCapabilityRevisionMatches(current, frozen.CapabilityRevision))
		})
	}
}

func sqlServerNativeDeliveryCapabilityFixture(
	alterAnyDataspace bool,
	editMetadata func(*schema.DatabaseConnector, *sqlServerRecordingMetadata),
	editPermissions func(*SQLServerRecordingPermissions),
) *SQLServerRecordingCapability {
	connector := &schema.DatabaseConnector{ID: "connector-1", IdentityRevision: "revision-1"}
	metadata := sqlServerRecordingMetadata{
		serverVersion: "16.0.1000.6", serverEdition: "Developer Edition (64-bit)",
		engineEdition: 3, engineEditionValid: true, compatibilityLevel: 160, compatibilityValid: true,
		database: "fixture_db", actualDatabase: "fixture_db", schema: "dbo", actualSchema: "dbo",
		metadataStatus: sqlServerRecordingMetadataVerified,
		permissions: SQLServerRecordingPermissions{
			Verified: true, DatabaseViewDefinition: true, SchemaViewDefinition: true,
			CreateTable: true, AlterSchema: true, SelectSchema: true, InsertSchema: true,
			DeleteSchema: true, AlterAnyDataspace: alterAnyDataspace,
		},
	}
	if editMetadata != nil {
		editMetadata(connector, &metadata)
	}
	if editPermissions != nil {
		editPermissions(&metadata.permissions)
	}
	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	capability.SupportsManagedLong = true
	return capability
}
