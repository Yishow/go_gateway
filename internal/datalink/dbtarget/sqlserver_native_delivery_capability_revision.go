package dbtarget

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

// SQLServerNativeDeliveryCapabilityRevisionMatches accepts the frozen SQL
// Server delivery revision, or the same verified scope with only ALTER ANY
// DATASPACE revoked. It is for runtime delivery checks, never DDL preparation.
func SQLServerNativeDeliveryCapabilityRevisionMatches(
	capability *SQLServerRecordingCapability,
	expectedRevision string,
) bool {
	if strings.TrimSpace(expectedRevision) == "" || validateSQLServerRecordingScope(capability) != nil ||
		!capability.SupportsManagedLong || !sqlServerRecordingPermissionsAllowManagedLong(capability.Permissions) ||
		!capability.Permissions.DatabaseViewDefinition && !capability.Permissions.SchemaViewDefinition {
		return false
	}

	connector := &schema.DatabaseConnector{
		ID: capability.ConnectorID, IdentityRevision: capability.ConnectorRevision,
	}
	// Verified metadata is produced only after checking the live database and
	// schema equal the requested scope and engine/compatibility values are valid.
	metadata := sqlServerRecordingMetadata{
		serverVersion: capability.ServerVersion, serverEdition: capability.ServerEdition,
		engineEdition: capability.EngineEdition, engineEditionValid: true,
		compatibilityLevel: capability.CompatibilityLevel, compatibilityValid: true,
		database: capability.Database, actualDatabase: capability.Database,
		schema: capability.Schema, actualSchema: capability.Schema,
		metadataStatus: capability.MetadataStatus, permissions: capability.Permissions,
	}
	currentRevision := sqlServerRecordingCapabilityRevision(connector, metadata, capability.Connected)
	if currentRevision != capability.CapabilityRevision {
		return false
	}
	if currentRevision == expectedRevision {
		return true
	}
	if capability.Permissions.AlterAnyDataspace {
		return false
	}

	metadata.permissions.AlterAnyDataspace = true
	return sqlServerRecordingCapabilityRevision(connector, metadata, capability.Connected) == expectedRevision
}
