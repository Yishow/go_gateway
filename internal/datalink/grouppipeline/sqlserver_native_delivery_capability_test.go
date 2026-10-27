package grouppipeline

import (
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeDeliveryCapabilityProofBoundary(t *testing.T) {
	snap := sqlServerNativeRuntimeSnapshot(t)
	proof := snap.RuntimeLayout.SQLServer
	capability := &dbtarget.SQLServerRecordingCapability{
		ConnectorID: proof.ConnectorID, ConnectorRevision: proof.ConnectorRevision,
		Database: proof.Group.Evidence.Database, Schema: proof.Group.Evidence.Schema,
		ServerVersion: proof.Group.Evidence.ServerVersion, ServerEdition: proof.Group.Evidence.ServerEdition,
		EngineEdition: proof.Group.Evidence.EngineEdition, CompatibilityLevel: proof.Group.Evidence.CompatibilityLevel,
		Connected: true, MetadataStatus: "verified", SupportsManagedLong: true,
		CapabilityRevision: proof.Group.Evidence.CapabilityRevision,
		Permissions: dbtarget.SQLServerRecordingPermissions{
			Verified: true, DatabaseViewDefinition: true, SchemaViewDefinition: true,
			CreateTable: true, AlterSchema: true, SelectSchema: true, InsertSchema: true, DeleteSchema: true,
		},
	}
	require.True(t, sqlServerCapabilityMatchesProof(capability, proof))
	capability.Permissions.AlterAnyDataspace = false
	require.True(t, sqlServerCapabilityMatchesProof(capability, proof), "exact frozen revision remains valid when maintenance permission is absent")

	changedVersion := *capability
	changedVersion.ServerVersion = "17.0.1"
	require.False(t, sqlServerCapabilityMatchesProof(&changedVersion, proof))
	changedRevision := *capability
	changedRevision.CapabilityRevision = "other-revision"
	require.False(t, sqlServerCapabilityMatchesProof(&changedRevision, proof))
}
