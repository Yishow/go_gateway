package dbtarget

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func mysqlReleaseMetadata(family, version string) mysqlRecordingMetadata {
	return mysqlRecordingMetadata{
		serverFamily: family, serverVersion: version, database: "fixture", actualDatabase: "fixture",
		metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true,
		permissions: mysqlRecordingAllPermissions(),
	}
}

func mysqlRecordingAllPermissions() MySQLRecordingPermissions {
	return MySQLRecordingPermissions{
		Verified: true, CreateTable: true, AlterSchema: true, IndexSchema: true, SelectSchema: true,
		InsertSchema: true, DeleteSchema: true,
	}
}

func TestMySQLRecordingReleaseGateProjection(t *testing.T) {
	connector := mysqlRecordingTestConnector()
	for _, tc := range []struct {
		name, family, version string
	}{
		{name: "mysql", family: "mysql", version: "8.4.13"},
		{name: "mariadb", family: "mariadb", version: "11.4.13-MariaDB-ubu2404"},
	} {
		t.Run(tc.name+" opens verified operations", func(t *testing.T) {
			capability := projectMySQLRecordingCapability(connector, mysqlReleaseMetadata(tc.family, tc.version), true)
			require.True(t, capability.Connected)
			require.Equal(t, mysqlRecordingMetadataVerified, capability.MetadataStatus)
			require.True(t, capability.TransactionalEngine)
			require.True(t, capability.SupportsManagedLong)
			require.True(t, capability.SupportsNativePartition)
			require.True(t, capability.SupportsGroupTestWrite)
			require.True(t, capability.SupportsMissingRepair)
			require.NotContains(t, capability.Reasons, "adapter_unverified")
			require.NotContains(t, capability.Reasons, "managed_long_unverified")
			require.NotContains(t, capability.Reasons, "native_partition_unverified")
			require.NotContains(t, capability.Reasons, "group_test_write_unverified")
			require.NotContains(t, capability.Reasons, "missing_repair_unverified")
		})
	}

	for _, tc := range []struct {
		name, family, version string
	}{
		{name: "old mysql patch", family: "mysql", version: "8.4.4"},
		{name: "new mysql patch", family: "mysql", version: "8.4.14"},
		{name: "old mariadb patch", family: "mariadb", version: "11.4.2-MariaDB-ubu2204"},
		{name: "new mariadb patch", family: "mariadb", version: "11.4.14-MariaDB-ubu2404"},
		{name: "foreign mariadb distro", family: "mariadb", version: "11.4.13-MariaDB-otherdistro"},
		{name: "unknown family", family: "unknown", version: "8.4.13"},
	} {
		t.Run(tc.name+" stays closed", func(t *testing.T) {
			capability := projectMySQLRecordingCapability(connector, mysqlReleaseMetadata(tc.family, tc.version), true)
			require.False(t, capability.SupportsManagedLong)
			require.False(t, capability.SupportsNativePartition)
			require.False(t, capability.SupportsGroupTestWrite)
			require.False(t, capability.SupportsMissingRepair)
			require.Contains(t, capability.Reasons, "adapter_unverified")
			require.Contains(t, capability.Reasons, "managed_long_unverified")
			require.Contains(t, capability.Reasons, "native_partition_unverified")
			require.Contains(t, capability.Reasons, "group_test_write_unverified")
			require.Contains(t, capability.Reasons, "missing_repair_unverified")
		})
	}

	t.Run("metadata engine connection and scope facts stay closed", func(t *testing.T) {
		for name, mutate := range map[string]func(*mysqlRecordingMetadata){
			"forbidden":   func(m *mysqlRecordingMetadata) { m.metadataStatus = mysqlRecordingMetadataForbidden },
			"unconfirmed": func(m *mysqlRecordingMetadata) { m.metadataStatus = mysqlRecordingMetadataUnconfirmed },
			"non innodb":  func(m *mysqlRecordingMetadata) { m.transactionalEngine = false },
			"wrong scope": func(m *mysqlRecordingMetadata) { m.actualDatabase = "other" },
		} {
			t.Run(name, func(t *testing.T) {
				metadata := mysqlReleaseMetadata("mysql", "8.4.13")
				mutate(&metadata)
				capability := projectMySQLRecordingCapability(connector, metadata, true)
				require.False(t, capability.SupportsManagedLong)
				require.False(t, capability.SupportsNativePartition)
				require.False(t, capability.SupportsGroupTestWrite)
				require.False(t, capability.SupportsMissingRepair)
			})
		}
		capability := projectMySQLRecordingCapability(connector, mysqlReleaseMetadata("mysql", "8.4.13"), false)
		require.False(t, capability.SupportsManagedLong)
		require.False(t, capability.SupportsNativePartition)
		require.False(t, capability.SupportsGroupTestWrite)
		require.False(t, capability.SupportsMissingRepair)
	})

	t.Run("capability hash preserves the legacy fully-authorized projection", func(t *testing.T) {
		goldenConnector := &schema.DatabaseConnector{ID: "saved-mysql", IdentityRevision: "saved-rev-1"}
		capability := projectMySQLRecordingCapability(goldenConnector, mysqlReleaseMetadata("mysql", "8.4.13"), true)
		require.Equal(t, "mysqlcap-v1-99ba817ad82e64d16e1b45ef19f4223ae726fba0be71fb6515115cc7da9adc06", capability.CapabilityRevision)
	})
}
