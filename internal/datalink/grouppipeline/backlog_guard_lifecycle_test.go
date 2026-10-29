package grouppipeline

import (
	"database/sql"
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestBacklogGuardDeleteMissingR0PreservesG1AndQ1(t *testing.T) {
	for _, kind := range []string{"outbox", "sample"} {
		t.Run(kind, func(t *testing.T) {
			db, service, mutation := backlogLifecycleFixture(t, kind)
			_, err := db.Exec(`UPDATE ` + backlogTable(kind) + ` SET group_revision='R0'`)
			require.NoError(t, err)
			before := backlogOwnershipState(t, db)
			result, err := service.Delete(t.Context(), "G1", mutation)
			require.Nil(t, result)
			require.ErrorIs(t, err, workspace.ErrWriteGroupLifecycleBlocked)
			require.ErrorIs(t, err, ErrBacklogOwnershipUnproven)
			require.Equal(t, before, backlogOwnershipState(t, db), "failed remove preserves G1, Q1, immutable versions, receipts and workspace CAS")
		})
	}
}

func TestBacklogGuardDeleteRetainsHistoricalPendingOwner(t *testing.T) {
	for _, kind := range []string{"outbox", "sample"} {
		t.Run(kind, func(t *testing.T) {
			db, service, mutation := backlogLifecycleFixture(t, kind)
			// The old accepted R1 remains frozen while the saved group advances.
			_, err := db.Exec(`UPDATE write_groups SET revision='R2', applied_revision='R2', destination_connector_revision='C2', destination_table_name='D2'`)
			require.NoError(t, err)
			_, err = db.Exec(`INSERT INTO write_group_versions(group_id,group_revision,effective_at,payload)
			 SELECT group_id,'R2','2026-01-01T00:00:10Z',json_set(payload,'$.revision','R2','$.applied_revision','R2','$.destination.connector_revision','C2','$.destination.table_name','D2') FROM write_group_versions WHERE group_revision='R1'`)
			require.NoError(t, err)
			_, err = db.Exec(`UPDATE database_connectors SET identity_revision='C2'`)
			require.NoError(t, err)
			mutation.ExpectedGroupRevision = "R2"
			mutation.ExpectedConnectorRevision = "C2"
			before := backlogOwnershipState(t, db)
			result, err := service.Delete(t.Context(), "G1", mutation)
			require.NoError(t, err)
			require.Equal(t, workspace.WriteGroupStatusDeleted, result.Group.Status)
			require.Equal(t, "R2", result.Group.AppliedRevision)
			after := backlogOwnershipState(t, db)
			for _, table := range []string{"write_group_versions", "wg_runtime_versions", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts"} {
				require.Equal(t, before[table], after[table], table+" must remain byte-for-byte unchanged")
			}
			reloaded, err := service.Get(t.Context(), "G1")
			require.NoError(t, err)
			require.Equal(t, result.Group, reloaded.Group, "tombstone remains readable")
			intake, err := service.IntakeEligibility(t.Context(), "G1")
			require.NoError(t, err)
			require.False(t, intake.Accepting)
			require.Equal(t, "group-deleted", intake.Reason)
		})
	}
}

func backlogLifecycleFixture(t *testing.T, kind string) (*sql.DB, *workspace.WriteGroupService, workspace.WriteGroupMutation) {
	t.Helper()
	db, group := backlogOwnershipFixture(t, kind)
	ws := workspace.NewService(workspace.NewSQLRepository(db))
	record, err := ws.GetOrCreate(t.Context())
	require.NoError(t, err)
	group.WorkspaceID = record.ID
	for _, table := range []string{"write_groups", "wg_delivery_samples", "wg_delivery_outbox"} {
		_, err := db.Exec(`UPDATE `+table+` SET workspace_id=?`, record.ID)
		require.NoError(t, err)
	}
	raw, err := json.Marshal(group)
	require.NoError(t, err)
	_, err = db.Exec(`UPDATE write_group_versions SET payload=?`, string(raw))
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO database_connectors(id,name,kind,identity_revision,connection_config,status) VALUES(?,'test','sqlite','C1','{}','ready')`, group.Destination.ConnectorID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO wg_delivery_receipts(effect_key,payload_digest,committed_at) VALUES('historic-receipt','historic-digest','2026-01-01T00:00:00Z')`)
	require.NoError(t, err)
	service := workspace.NewWriteGroupService(ws, workspace.NewSQLWriteGroupRepository(db)).WithBacklogOwnershipGuard(BacklogGuard{})
	return db, service, workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedGroupRevision: "R1", ExpectedConnectorRevision: "C1"}
}
