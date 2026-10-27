package grouppipeline

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLNativeRecoveryUsesConfirmedInitialLayout(t *testing.T) {
	f := newFixture(t)
	group, frozen := mysqlNativeMonthlyFixture(t)
	scope := workspace.StorageScopeOf(group)
	table, owner, ok := workspace.LongStorageIdentity(group)
	require.True(t, ok)
	evidence := frozen.MySQL.Group.Evidence
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: evidence.ServerFamily, ServerVersion: evidence.ServerVersion, CapabilityRevision: evidence.CapabilityRevision,
		Engine: evidence.Engine, Charset: evidence.Charset, Collation: evidence.Collation,
		NativeMonthly: mirrorMySQLPartitionEvidence(evidence.Partitions),
	})
	require.NoError(t, err)
	_, err = f.store.DB().ExecContext(t.Context(), `INSERT INTO write_group_storage_scopes
		(workspace_id,connector_id,database_name,table_schema,storage_mode,layout_version,table_name,owner_column,layout_digest,generation)
		VALUES (?,?,?,?,?,?,?,?,?,1)`, scope.WorkspaceID, scope.ConnectorID, scope.Database, scope.TableSchema,
		scope.StorageMode, scope.LayoutVersion, table, owner, recordingplan.UnconfirmedMySQLStorageDigest(scope))
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), f.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: group.Destination.ConnectorRevision, OperationID: group.Destination.SchemaRevision,
		VerifiedDigest: group.Destination.SchemaDigest, GroupID: group.ID,
	}, time.Now().UTC()))
	recovery := &StorageRecovery{store: f.store}
	fallback, lineage, err := recovery.mysqlConfirmedFallbackLayout(t.Context(), scope, table, owner)
	require.NoError(t, err)
	require.NotNil(t, lineage)
	require.NotNil(t, fallback, "native missing repair needs the exact immutable initial layout, preserving receipt ownership")
	require.Equal(t, layout, *fallback)
	_, err = f.store.DB().ExecContext(t.Context(), `UPDATE write_group_storage_scopes SET confirmed_layout = REPLACE(confirmed_layout, 'native_monthly', 'foreign')`)
	require.NoError(t, err)
	fallback, _, err = recovery.mysqlConfirmedFallbackLayout(t.Context(), scope, table, owner)
	require.True(t, err != nil || fallback == nil, "changed native proof cannot authorize DDL")
}
