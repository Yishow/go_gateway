package dbtarget

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLServerRecordingReleaseCapability2025(t *testing.T) {
	connector := &schema.DatabaseConnector{
		ID: "connector-17", Name: "sqlserver-2025", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "rev-17", Enabled: true,
	}
	metadata := sqlServerRecordingMetadata{
		database:           "gw_circuit_acceptance",
		actualDatabase:     "gw_circuit_acceptance",
		schema:             "dbo",
		actualSchema:       "dbo",
		serverVersion:      "17.0.5035.5",
		serverEdition:      "Enterprise Developer Edition (64-bit)",
		engineEdition:      3,
		engineEditionValid: true,
		compatibilityLevel: 170,
		compatibilityValid: true,
		metadataStatus:     sqlServerRecordingMetadataVerified,
		permissions: SQLServerRecordingPermissions{
			Verified:               true,
			DatabaseViewDefinition: true,
			SchemaViewDefinition:   true,
			CreateTable:            true,
			AlterSchema:            true,
			SelectSchema:           true,
			InsertSchema:           true,
			DeleteSchema:           true,
			AlterAnyDataspace:      true,
		},
	}

	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	require.True(t, capability.Connected)
	require.Equal(t, sqlServerRecordingMetadataVerified, capability.MetadataStatus)
	require.True(t, capability.SupportsManagedLong, "SQL Server 2025 (17.x/170) must support managed long when permissions are verified")
	require.True(t, capability.SupportsNativePartition, "AlterAnyDataspace must enable native partition for 2025")
	require.True(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
	require.NotContains(t, capability.Reasons, "server_version_unverified")
	require.NotContains(t, capability.Reasons, "compatibility_unverified")
	require.NotContains(t, capability.Reasons, "managed_long_unverified")
}

func TestSQLServerRecordingReleaseCapabilityMissingDataspacePermission(t *testing.T) {
	connector := &schema.DatabaseConnector{
		ID: "connector-17", Name: "sqlserver-2025", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "rev-17", Enabled: true,
	}
	metadata := sqlServerRecordingMetadata{
		database:           "gw_circuit_acceptance",
		actualDatabase:     "gw_circuit_acceptance",
		schema:             "dbo",
		actualSchema:       "dbo",
		serverVersion:      "17.0.5035.5",
		serverEdition:      "Enterprise Developer Edition (64-bit)",
		engineEdition:      3,
		engineEditionValid: true,
		compatibilityLevel: 170,
		compatibilityValid: true,
		metadataStatus:     sqlServerRecordingMetadataVerified,
		permissions: SQLServerRecordingPermissions{
			Verified:               true,
			DatabaseViewDefinition: true,
			SchemaViewDefinition:   true,
			CreateTable:            true,
			AlterSchema:            true,
			SelectSchema:           true,
			InsertSchema:           true,
			DeleteSchema:           true,
			AlterAnyDataspace:      false, // denied
		},
	}

	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	require.True(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition, "missing AlterAnyDataspace must keep native partition disabled")
}

func TestSQLServerRecordingReleaseCapabilityMissingDatabaseViewDefinitionPermission(t *testing.T) {
	connector := &schema.DatabaseConnector{
		ID: "connector-17", Name: "sqlserver-2025", Kind: schema.DatabaseConnectorKindSQLServer,
		IdentityRevision: "rev-17", Enabled: true,
	}
	metadata := sqlServerRecordingMetadata{
		database:           "gw_circuit_acceptance",
		actualDatabase:     "gw_circuit_acceptance",
		schema:             "dbo",
		actualSchema:       "dbo",
		serverVersion:      "17.0.5035.5",
		serverEdition:      "Enterprise Developer Edition (64-bit)",
		engineEdition:      3,
		engineEditionValid: true,
		compatibilityLevel: 170,
		compatibilityValid: true,
		metadataStatus:     sqlServerRecordingMetadataVerified,
		permissions: SQLServerRecordingPermissions{
			Verified:               true,
			DatabaseViewDefinition: false, // schema-only VIEW DEFINITION
			SchemaViewDefinition:   true,
			CreateTable:            true,
			AlterSchema:            true,
			SelectSchema:           true,
			InsertSchema:           true,
			DeleteSchema:           true,
			AlterAnyDataspace:      true,
		},
	}

	capability := projectSQLServerRecordingCapability(connector, metadata, true)
	require.True(t, capability.SupportsManagedLong, "schema-only VIEW DEFINITION must preserve plain managed long")
	require.False(t, capability.SupportsNativePartition, "missing DatabaseViewDefinition must keep native partition disabled")
	require.Contains(t, capability.Reasons, "native_partition_permissions_denied")
}
