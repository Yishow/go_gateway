package grouppipeline

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type sqlServerMaintenanceDeniedDestinations struct {
	*fakeDestinations
	opened *dbtarget.OpenedDestination
}

func (d *sqlServerMaintenanceDeniedDestinations) OpenDestination(context.Context, string, string) (*dbtarget.OpenedDestination, error) {
	return d.opened, nil
}

func TestSQLServerPartitionDeniedWorkerPersistsVisibleState(t *testing.T) {
	base := newFixture(t)
	pipeline := base.pipeline()
	snap := sqlServerNativeRuntimeSnapshot(t)
	proof := snap.RuntimeLayout.SQLServer
	native := sqlServerPlanEvidenceFromDBForTest(proof.Group.Evidence.NativeMonthly)
	layout := sqlServerLayoutForNativeTest(t, proof, native, snap.Group)
	scope := workspace.StorageScopeOf(snap.Group)
	_, err := base.store.DB().ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema,
		scope.StorageMode, scope.LayoutVersion, layout.TableName, layout.OwnerColumn, recordingplan.UnconfirmedSQLServerStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), base.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: snap.Group.Destination.ConnectorRevision, OperationID: snap.Group.Destination.SchemaRevision,
		VerifiedDigest: snap.Group.Destination.SchemaDigest, CreatedTable: true, GroupID: snap.Group.ID,
	}, time.Now().UTC()))
	capability := sqlServerNativeDeliveryFixtureCapability(t, false)
	applyNativeDeliveryCapabilityFixture(capability, proof)
	pipeline.deps.Destinations = &sqlServerMaintenanceDeniedDestinations{
		fakeDestinations: base.dest,
		opened: &dbtarget.OpenedDestination{
			DB: base.store.DB(), Kind: schema.DatabaseConnectorKindSQLServer, SQLServer: capability,
		},
	}
	splits := 0
	pipeline.sqlServerSplit = func(context.Context, *sql.DB, *dbtarget.SQLServerRecordingCapability,
		*dbtarget.SQLServerManagedInspection, *dbtarget.SQLServerManagedPartitionEvidence) error {
		splits++
		return nil
	}
	err = pipeline.maintainSQLServerPartitionScope(t.Context(), snap)
	require.ErrorContains(t, err, "partition_maintenance_blocked")
	state, err := recordingplan.ReadSQLServerPartitionMaintenanceState(t.Context(), base.store.DB(), scope)
	require.NoError(t, err)
	require.NotNil(t, state, "worker failures must survive restart and be available to readiness/UI")
	require.Equal(t, "blocked", state.Status)
	require.Equal(t, "permission_lost", state.Reason)
	require.Equal(t, native, state.Current)
	require.Zero(t, splits, "denied partition permission never reaches DDL")
}
