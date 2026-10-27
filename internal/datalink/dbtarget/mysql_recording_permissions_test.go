package dbtarget

import (
	"testing"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestMySQLCapabilityDoesNotOpenOperationsWhenPermissionEvidenceIsDenied(t *testing.T) {
	metadata := mysqlReleaseMetadata("mysql", "8.4.13")
	metadata.permissions.InsertSchema = false

	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.False(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
	require.Contains(t, capability.Reasons, "operation_permissions_denied")
}

func TestMySQLRecordingPermissionGatesAreIndependent(t *testing.T) {
	for _, denied := range []string{"CREATE", "ALTER", "INDEX", "SELECT", "INSERT", "DELETE"} {
		t.Run(denied, func(t *testing.T) {
			script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
			filtered := make([]mysqlRecordingPrivilegeFixture, 0, len(script.schemaPrivileges)-1)
			for _, privilege := range script.schemaPrivileges {
				if privilege.privilegeType != denied {
					filtered = append(filtered, privilege)
				}
			}
			script.schemaPrivileges = filtered

			metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
			capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

			require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
			require.True(t, capability.Permissions.Verified)
			require.Contains(t, capability.Reasons, mysqlRecordingPermissionDeniedReason)
			switch denied {
			case "CREATE", "INDEX":
				require.False(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
				require.True(t, capability.SupportsGroupTestWrite)
				require.False(t, capability.SupportsMissingRepair)
			case "ALTER":
				require.True(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
				require.True(t, capability.SupportsGroupTestWrite)
				require.True(t, capability.SupportsMissingRepair)
			case "SELECT":
				require.False(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
				require.False(t, capability.SupportsGroupTestWrite)
				require.False(t, capability.SupportsMissingRepair)
			case "INSERT":
				require.False(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
				require.False(t, capability.SupportsGroupTestWrite)
				require.True(t, capability.SupportsMissingRepair)
			case "DELETE":
				require.True(t, capability.SupportsManagedLong)
				require.True(t, capability.SupportsNativePartition)
				require.False(t, capability.SupportsGroupTestWrite)
				require.True(t, capability.SupportsMissingRepair)
			}
		})
	}
}

func TestMySQLRecordingGlobalPrivilegesAreAcceptedForOperationGates(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.currentUser = "root@localhost"
	script.schemaPrivileges = nil
	script.globalPrivileges = mysqlRecordingAllPrivilegeFixtures("'root'@'localhost'")

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.True(t, capability.Permissions.Verified)
	require.True(t, capability.SupportsManagedLong)
	require.True(t, capability.SupportsNativePartition)
	require.True(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
}

func TestMySQLRecordingActiveRoleKeepsOperationsUnverified(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.currentRole = "'recording_role'@'%'"

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
	require.False(t, capability.Permissions.Verified)
	require.False(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsGroupTestWrite)
	require.False(t, capability.SupportsMissingRepair)
	require.Contains(t, capability.Reasons, mysqlRecordingRoleUnknownReason)
	require.NotContains(t, capability.Reasons, "database_missing")
}

func TestMySQLRecordingPermissionDeniedIsNotMetadataMissing(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.errorQuery = mysqlRecordingGlobalPrivilegesQuery
	script.queryError = &mysqldriver.MySQLError{Number: 1142, Message: "SELECT command denied"}

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
	require.False(t, capability.Permissions.Verified)
	require.Contains(t, capability.Reasons, "operation_permissions_forbidden")
	require.NotContains(t, capability.Reasons, "database_missing")
	require.False(t, capability.SupportsManagedLong)
}

func TestMySQLRecordingMariaDBNullCurrentRoleDoesNotBlockDirectGrants(t *testing.T) {
	script := mysqlRecordingMetadataFixture("11.4.13-MariaDB-ubu2404", "MariaDB Server")
	script.currentRoleNull = true

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
	require.True(t, capability.SupportsManagedLong)
	require.True(t, capability.SupportsNativePartition)
	require.True(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
}

func TestMySQLRecordingMySQLNullCurrentRoleStaysUnverified(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.currentRoleNull = true

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.False(t, capability.Permissions.Verified)
	require.False(t, capability.SupportsManagedLong)
	require.Contains(t, capability.Reasons, mysqlRecordingPermissionUnknownReason)
}

func TestMySQLRecordingEscapedSchemaPrivilegeStaysBoundToExactDatabase(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.currentDatabase = "gw_circuit_test"
	script.visibleDatabase = script.currentDatabase
	for index := range script.schemaPrivileges {
		script.schemaPrivileges[index].schemaName = `gw\_circuit\_test`
	}

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), script.currentDatabase)
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
	require.True(t, capability.SupportsManagedLong)
	require.True(t, capability.SupportsNativePartition)
	require.True(t, capability.SupportsGroupTestWrite)
	require.True(t, capability.SupportsMissingRepair)
}

func TestMySQLRecordingOtherSchemaPrivilegesDoNotOpenExpectedDatabase(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.currentDatabase = "gw_circuit_test"
	script.visibleDatabase = script.currentDatabase
	for index := range script.schemaPrivileges {
		script.schemaPrivileges[index].schemaName = "otherdb"
	}

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), script.currentDatabase)
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.True(t, capability.Permissions.Verified)
	require.False(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsGroupTestWrite)
	require.False(t, capability.SupportsMissingRepair)
	require.Contains(t, capability.Reasons, mysqlRecordingPermissionDeniedReason)
}

func TestMySQLRecordingUnknownSchemaWildcardClosesPermissionGates(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	for index := range script.schemaPrivileges {
		script.schemaPrivileges[index].schemaName = "metrics%"
	}

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.False(t, capability.Permissions.Verified)
	require.False(t, capability.SupportsManagedLong)
	require.Contains(t, capability.Reasons, mysqlRecordingPermissionUnknownReason)
}

func TestMySQLRecordingMissingIndexPermissionClosesDDLOperations(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	filtered := make([]mysqlRecordingPrivilegeFixture, 0, len(script.schemaPrivileges)-1)
	for _, privilege := range script.schemaPrivileges {
		if privilege.privilegeType != "INDEX" {
			filtered = append(filtered, privilege)
		}
	}
	script.schemaPrivileges = filtered

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.False(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsMissingRepair)
	require.True(t, capability.SupportsGroupTestWrite)
}

func TestMySQLRecordingPartialRevokeClosesPermissionGates(t *testing.T) {
	script := mysqlRecordingMetadataFixture("8.4.13", "MySQL Community Server - GPL")
	script.showGrants = []string{
		"GRANT SELECT, INSERT, DELETE, CREATE, ALTER, INDEX ON *.* TO 'fixture'@'%'",
		"REVOKE INSERT ON `metrics`.* FROM 'fixture'@'%'",
	}

	metadata := inspectMySQLRecordingMetadata(t.Context(), openMySQLRecordingMetadataScript(t, script), "metrics")
	capability := projectMySQLRecordingCapability(mysqlRecordingTestConnector(), metadata, true)

	require.False(t, capability.SupportsManagedLong)
	require.False(t, capability.SupportsNativePartition)
	require.False(t, capability.SupportsGroupTestWrite)
	require.False(t, capability.SupportsMissingRepair)
	require.Contains(t, capability.Reasons, "operation_permissions_partial_revoke")
}

func TestMySQLRecordingGranteeMatchingPreservesAccountIdentity(t *testing.T) {
	require.False(t, mysqlRecordingGranteeMatches("'Reader'@'localhost'", "reader@localhost"))
	require.False(t, mysqlRecordingGranteeMatches("'reader''x'@'localhost'", "readerx@localhost"))
	require.True(t, mysqlRecordingGranteeMatches("'reader''x'@'localhost'", "reader'x@localhost"))
}
