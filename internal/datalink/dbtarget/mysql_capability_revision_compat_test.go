package dbtarget

import (
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMySQLCapabilityRevisionPreservesPreparedTargetsWhenOperationsOpen(t *testing.T) {
	connector := &schema.DatabaseConnector{ID: "saved-mysql", IdentityRevision: "saved-rev-1"}
	metadata := mysqlRecordingMetadata{serverFamily: "mysql", serverVersion: "8.4.13", database: "fixture", actualDatabase: "fixture", metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true, permissions: mysqlRecordingAllPermissions()}
	prepared := projectMySQLRecordingCapability(connector, metadata, true)
	// Fully authorized permission facts preserve the pre-probe physical hash.
	require.Equal(t, "mysqlcap-v1-99ba817ad82e64d16e1b45ef19f4223ae726fba0be71fb6515115cc7da9adc06", prepared.CapabilityRevision)
	enabled := *prepared
	enabled.SupportsManagedLong, enabled.SupportsNativePartition = true, true
	enabled.SupportsGroupTestWrite, enabled.SupportsMissingRepair = true, true
	require.Equal(t, prepared.CapabilityRevision, mysqlRecordingCapabilityRevision(connector, metadata, &enabled),
		"release availability cannot retarget an already confirmed physical schema or accepted effect")
}

func TestMySQLCapabilityRevisionInvalidatesChangedPermissions(t *testing.T) {
	connector := &schema.DatabaseConnector{ID: "saved-mysql", IdentityRevision: "saved-rev-1"}
	metadata := mysqlRecordingMetadata{serverFamily: "mysql", serverVersion: "8.4.13", database: "fixture", actualDatabase: "fixture", metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true, permissions: mysqlRecordingAllPermissions()}
	before := projectMySQLRecordingCapability(connector, metadata, true).CapabilityRevision
	metadata.permissions.DeleteSchema = false
	after := projectMySQLRecordingCapability(connector, metadata, true).CapabilityRevision
	require.NotEqual(t, before, after)
	require.True(t, strings.HasPrefix(after, "mysqlcap-v2-"))
}

func TestMySQLCapabilityRevisionStillInvalidatesChangedPhysicalIdentity(t *testing.T) {
	connector := &schema.DatabaseConnector{ID: "saved-mysql", IdentityRevision: "saved-rev-1"}
	metadata := mysqlRecordingMetadata{serverFamily: "mysql", serverVersion: "8.4.13", database: "fixture", actualDatabase: "fixture", metadataStatus: mysqlRecordingMetadataVerified, transactionalEngine: true, permissions: mysqlRecordingAllPermissions()}
	before := projectMySQLRecordingCapability(connector, metadata, true).CapabilityRevision
	for _, mutate := range []func(*mysqlRecordingMetadata){
		func(m *mysqlRecordingMetadata) { m.serverVersion = "8.4.14" },
		func(m *mysqlRecordingMetadata) { m.serverFamily = "mariadb" },
		func(m *mysqlRecordingMetadata) { m.actualDatabase = "foreign" },
		func(m *mysqlRecordingMetadata) { m.metadataStatus = mysqlRecordingMetadataForbidden },
		func(m *mysqlRecordingMetadata) { m.transactionalEngine = false },
	} {
		changed := metadata
		mutate(&changed)
		require.NotEqual(t, before, projectMySQLRecordingCapability(connector, changed, true).CapabilityRevision)
	}
	changedConnector := *connector
	changedConnector.IdentityRevision = "saved-rev-2"
	require.NotEqual(t, before, projectMySQLRecordingCapability(&changedConnector, metadata, true).CapabilityRevision)
}
