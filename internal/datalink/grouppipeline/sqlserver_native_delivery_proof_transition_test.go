package grouppipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestSQLServerNativeDeliveryProofKeepsCurrentRevisionOnAlterAnyLoss(t *testing.T) {
	base := newFixture(t)
	snap := sqlServerNativeRuntimeSnapshot(t)
	initial := dbtarget.CloneSQLServerManagedRowProof(snap.RuntimeLayout.SQLServer)
	frozenCapability := sqlServerNativeDeliveryFixtureCapability(t, true)
	currentCapability := sqlServerNativeDeliveryFixtureCapability(t, false)
	for _, capability := range []*dbtarget.SQLServerRecordingCapability{frozenCapability, currentCapability} {
		applyNativeDeliveryCapabilityFixture(capability, initial)
	}
	initial.Group.Evidence.CapabilityRevision = frozenCapability.CapabilityRevision
	initial.Receipt.Evidence.CapabilityRevision = frozenCapability.CapabilityRevision
	snap.RuntimeLayout.SQLServer = initial

	native := sqlServerPlanEvidenceFromDBForTest(initial.Group.Evidence.NativeMonthly)
	layout := sqlServerLayoutForNativeTest(t, initial, native, snap.Group)
	groupTarget := sqlServerRestoreTargetInspection(initial.Group)
	receiptTarget := sqlServerRestoreTargetInspection(initial.Receipt)
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, snap.Group.Destination.Database, groupTarget, receiptTarget)
	require.NoError(t, err)
	snap.Group.Destination.SchemaDigest = digest
	snap.RuntimeLayout.SchemaDigest = digest
	scope := workspace.StorageScopeOf(snap.Group)
	_, err = base.store.DB().ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema,
		scope.StorageMode, scope.LayoutVersion, layout.TableName, layout.OwnerColumn, recordingplan.UnconfirmedSQLServerStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), base.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: snap.Group.Destination.ConnectorRevision, OperationID: snap.Group.Destination.SchemaRevision,
		VerifiedDigest: digest, CreatedTable: true, GroupID: snap.Group.ID,
	}, time.Now().UTC()))

	actualGroup := dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: initial.Group}).Group
	actualReceipt := dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: initial.Receipt}).Group
	actualGroup.Evidence.CapabilityRevision = currentCapability.CapabilityRevision
	actualReceipt.Evidence.CapabilityRevision = currentCapability.CapabilityRevision
	repo := recordingplan.NewSQLRepository(base.store.DB())
	recovery := NewStorageRecovery(base.store, recordingplan.NewService(repo), nil)
	recovery.sqlServerOps = &sqlServerStorageRestoreTargetOps{
		Inspect: func(_ context.Context, _ *sql.DB, _ *dbtarget.SQLServerRecordingCapability, table string) (*dbtarget.SQLServerManagedInspection, error) {
			if table == dbtarget.EffectReceiptTable {
				return dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: actualReceipt}).Group, nil
			}
			return dbtarget.CloneSQLServerManagedRowProof(&dbtarget.SQLServerManagedRowProof{Group: actualGroup}).Group, nil
		},
	}
	frozenJSON, err := json.Marshal(initial)
	require.NoError(t, err)
	require.NoError(t, workspace.ValidateSQLServerRuntimeLayout(snap.Group, snap.RuntimeLayout))
	require.True(t, dbtarget.SQLServerNativeDeliveryCapabilityRevisionMatches(currentCapability, frozenCapability.CapabilityRevision))
	require.True(t, sqlServerCapabilityMatchesProof(currentCapability, initial))
	_, _, err = recovery.sqlServerConfirmedFallbackLayout(t.Context(), sqlServerRestoreRequest(snap))
	require.NoError(t, err)
	verifiedGroup, verifiedReceipt, err := normalizeSQLServerDeliveryInspectionsForVerification(actualGroup, actualReceipt, frozenCapability.CapabilityRevision)
	require.NoError(t, err)
	verifiedGroupState, err := recordingplan.NormalizeSQLServerPartitionCatalog(t.Context(), base.store.DB(), scope, layout, sqlServerRestoreTargetInspection(verifiedGroup), sqlServerRestoreTargetInspection(verifiedReceipt))
	require.NoError(t, err)
	require.NoError(t, recovery.ledger.VerifySQLServerStorageGeneration(t.Context(), sqlServerRestoreRequest(snap), verifiedGroupState, sqlServerRestoreTargetInspection(verifiedReceipt)))
	proof, err := recovery.sqlServerCurrentProof(t.Context(), snap, groupdelivery.Target{SQLServerCapability: currentCapability}, false)
	require.NoError(t, err)
	require.Equal(t, currentCapability.CapabilityRevision, proof.Group.Evidence.CapabilityRevision)
	require.Equal(t, currentCapability.CapabilityRevision, proof.Receipt.Evidence.CapabilityRevision)
	afterJSON, err := json.Marshal(initial)
	require.NoError(t, err)
	require.Equal(t, frozenJSON, afterJSON, "frozen proof bytes must not be rewritten")

	versionDrift := sqlServerNativeDeliveryFixtureCapability(t, false)
	applyNativeDeliveryCapabilityFixture(versionDrift, initial)
	versionDrift.ServerVersion = "17.0.1"
	versionDrift.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(versionDrift)
	_, err = recovery.sqlServerCurrentProof(t.Context(), snap, groupdelivery.Target{SQLServerCapability: versionDrift}, false)
	require.Error(t, err)

	permissionDrift := sqlServerNativeDeliveryFixtureCapability(t, false)
	applyNativeDeliveryCapabilityFixture(permissionDrift, initial)
	permissionDrift.Permissions.SelectSchema = false
	permissionDrift.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(permissionDrift)
	_, err = recovery.sqlServerCurrentProof(t.Context(), snap, groupdelivery.Target{SQLServerCapability: permissionDrift}, false)
	require.Error(t, err)
}

func applyNativeDeliveryCapabilityFixture(capability *dbtarget.SQLServerRecordingCapability, proof *dbtarget.SQLServerManagedRowProof) {
	capability.ConnectorID, capability.ConnectorRevision = proof.ConnectorID, proof.ConnectorRevision
	capability.Database, capability.Schema = proof.Group.Evidence.Database, proof.Group.Evidence.Schema
	capability.ServerVersion, capability.ServerEdition = proof.Group.Evidence.ServerVersion, proof.Group.Evidence.ServerEdition
	capability.EngineEdition, capability.CompatibilityLevel = proof.Group.Evidence.EngineEdition, proof.Group.Evidence.CompatibilityLevel
	capability.CapabilityRevision = sqlServerNativeDeliveryCapabilityRevisionForTest(capability)
}

func sqlServerLayoutForNativeTest(t *testing.T, proof *dbtarget.SQLServerManagedRowProof, native *recordingplan.SQLServerPartitionEvidence, group *workspace.WriteGroup) recordingplan.GroupSchemaLayout {
	t.Helper()
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	options := recordingplan.SQLServerGroupSchemaOptions{
		Database: proof.Group.Evidence.Database, Schema: proof.Group.Evidence.Schema,
		Collation: proof.Group.Evidence.ColumnCollations["record_id"], ServerVersion: proof.Group.Evidence.ServerVersion,
		ServerEdition: proof.Group.Evidence.ServerEdition, CapabilityRevision: proof.Group.Evidence.CapabilityRevision,
		EngineEdition: proof.Group.Evidence.EngineEdition, CompatibilityLevel: proof.Group.Evidence.CompatibilityLevel,
		NativeMonthly: native,
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	require.NoError(t, err)
	return layout
}

func sqlServerPlanEvidenceFromDBForTest(native *dbtarget.SQLServerManagedPartitionEvidence) *recordingplan.SQLServerPartitionEvidence {
	if native == nil {
		return nil
	}
	partitions := make([]recordingplan.SQLServerPartitionBoundary, len(native.Partitions))
	for i, partition := range native.Partitions {
		partitions[i] = recordingplan.SQLServerPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return &recordingplan.SQLServerPartitionEvidence{
		Method: native.Method, Column: native.Column, FunctionName: native.FunctionName, SchemeName: native.SchemeName,
		Filegroup: native.Filegroup, Partitions: partitions,
	}
}
