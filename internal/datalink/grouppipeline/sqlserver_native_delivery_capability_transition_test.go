package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func sqlServerNativeDeliveryFixtureCapability(t *testing.T, alterAny bool) *dbtarget.SQLServerRecordingCapability {
	t.Helper()
	capability := &dbtarget.SQLServerRecordingCapability{
		ConnectorID: "connector-1", ConnectorRevision: "revision-1", Database: "fixture_db", Schema: "dbo",
		ServerVersion: "16.0.1000.6", ServerEdition: "Developer Edition (64-bit)", EngineEdition: 3,
		CompatibilityLevel: 160, Connected: true, MetadataStatus: "verified", SupportsManagedLong: true,
		Permissions: dbtarget.SQLServerRecordingPermissions{
			Verified: true, DatabaseViewDefinition: true, SchemaViewDefinition: true,
			CreateTable: true, AlterSchema: true, SelectSchema: true, InsertSchema: true, DeleteSchema: true,
			AlterAnyDataspace: alterAny,
		},
	}
	capability.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(capability)
	return capability
}

func sqlServerNativeDeliveryCapabilityRevisionForTest(capability *dbtarget.SQLServerRecordingCapability) string {
	p := capability.Permissions
	canonical := fmt.Appendf(nil,
		"%q|%q|%q|%q|%q|%q|%d|%t|%d|%t|%q|%q|%q|%q|%t|%t|%t|%t|%t|%t|%t|%t|%t|%t",
		"sqlserver-recording-capability-v1", capability.ConnectorID, capability.ConnectorRevision,
		capability.ServerVersion, capability.ServerEdition, capability.Database,
		capability.EngineEdition, true, capability.CompatibilityLevel, true, capability.Database,
		capability.Schema, capability.Schema, capability.MetadataStatus,
		p.Verified, p.DatabaseViewDefinition, p.SchemaViewDefinition, p.CreateTable, p.AlterSchema,
		p.SelectSchema, p.InsertSchema, p.DeleteSchema, p.AlterAnyDataspace, capability.Connected,
	)
	sum := sha256.Sum256(canonical)
	return "sqlservercap-v1-" + hex.EncodeToString(sum[:])
}

func TestSQLServerNativeDeliveryCapabilityRevisionTransition(t *testing.T) {
	snap := sqlServerNativeRuntimeSnapshot(t)
	frozenCapability := sqlServerNativeDeliveryFixtureCapability(t, true)
	currentCapability := sqlServerNativeDeliveryFixtureCapability(t, false)
	proof := dbtarget.CloneSQLServerManagedRowProof(snap.RuntimeLayout.SQLServer)
	for _, capability := range []*dbtarget.SQLServerRecordingCapability{frozenCapability, currentCapability} {
		capability.ConnectorID, capability.ConnectorRevision = proof.ConnectorID, proof.ConnectorRevision
		capability.Database, capability.Schema = proof.Group.Evidence.Database, proof.Group.Evidence.Schema
		capability.ServerVersion, capability.ServerEdition = proof.Group.Evidence.ServerVersion, proof.Group.Evidence.ServerEdition
		capability.EngineEdition, capability.CompatibilityLevel = proof.Group.Evidence.EngineEdition, proof.Group.Evidence.CompatibilityLevel
		capability.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(capability)
	}
	proof.Group.Evidence.CapabilityRevision = frozenCapability.CapabilityRevision
	proof.Receipt.Evidence.CapabilityRevision = frozenCapability.CapabilityRevision
	require.True(t, sqlServerCapabilityMatchesProof(frozenCapability, proof))
	require.True(t, sqlServerCapabilityMatchesProof(currentCapability, proof), "AlterAnyDataspace loss must preserve native delivery")
	frozenJSON, err := json.Marshal(proof)
	require.NoError(t, err)
	require.True(t, sqlServerCapabilityMatchesProof(currentCapability, proof))
	afterJSON, err := json.Marshal(proof)
	require.NoError(t, err)
	require.Equal(t, frozenJSON, afterJSON, "runtime compatibility check must not mutate the frozen proof")

	versionDrift := sqlServerNativeDeliveryFixtureCapability(t, false)
	versionDrift.ServerVersion = "17.0.1"
	versionDrift.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(versionDrift)
	require.False(t, sqlServerCapabilityMatchesProof(versionDrift, proof))

	permissionDrift := sqlServerNativeDeliveryFixtureCapability(t, false)
	permissionDrift.Permissions.SelectSchema = false
	permissionDrift.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(permissionDrift)
	require.False(t, sqlServerCapabilityMatchesProof(permissionDrift, proof))

	unknownRevision := sqlServerNativeDeliveryFixtureCapability(t, false)
	unknownRevision.CapabilityRevision = frozenCapability.CapabilityRevision + "stale"
	require.False(t, sqlServerCapabilityMatchesProof(unknownRevision, proof))

	actualGroup := dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: proof.Group}).Group
	actualReceipt := dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: proof.Receipt}).Group
	actualGroup.Evidence.CapabilityRevision = currentCapability.CapabilityRevision
	actualReceipt.Evidence.CapabilityRevision = currentCapability.CapabilityRevision
	verifiedGroup, verifiedReceipt, err := normalizeSQLServerDeliveryInspectionsForVerification(actualGroup, actualReceipt, frozenCapability.CapabilityRevision)
	require.NoError(t, err)
	require.Equal(t, frozenCapability.CapabilityRevision, verifiedGroup.Evidence.CapabilityRevision)
	require.Equal(t, frozenCapability.CapabilityRevision, verifiedReceipt.Evidence.CapabilityRevision)
	require.Equal(t, currentCapability.CapabilityRevision, actualGroup.Evidence.CapabilityRevision)
	require.Equal(t, currentCapability.CapabilityRevision, actualReceipt.Evidence.CapabilityRevision)
}
