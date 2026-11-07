package workspace

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func seedOwnedAdditionalDeliveryRows(t *testing.T, db *sql.DB, group *WriteGroup) {
	t.Helper()
	ctx := t.Context()
	_, err := db.ExecContext(ctx, `INSERT INTO wg_delivery_mixed_progress(workspace_id,group_id,group_revision,version,descriptor_digest,closed_through,checkpoint,checkpoint_digest,last_closure_digest,updated_at) VALUES(?,'G1','R1',1,'descriptor','2026-10-02T12:00:00Z','{}','checkpoint','closure','2026-10-02T12:00:00Z')`, group.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO wg_delivery_mixed_outcomes(workspace_id,group_id,group_revision,record_id,due_at,kind,outcome,outcome_digest,created_at) VALUES(?,'G1','R1','mixed-record','2026-10-02T12:00:00Z','row','{}','outcome','2026-10-02T12:00:00Z')`, group.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO wg_delivery_sample_handoffs(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record) VALUES('H1',?,'G1','R1',1,'G2','R2',2,1,'finalized','{}')`, group.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('C1',?,'forward','G1','2026-10-02T12:00:00Z','review','{}','{}')`, group.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('C1',?,'ended','G1','R1','D1','P1','T1','','connector','revision','db','main','table')`, group.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,'G1','R1','2026-10-02T12:00:00Z','C1')`, group.WorkspaceID)
	require.NoError(t, err)
	for _, table := range []string{"wg_runtime_versions", "wg_delivery_mixed_progress", "wg_delivery_mixed_outcomes", "wg_delivery_sample_handoffs", "write_group_cutovers", "write_group_ownership_fences", "write_group_revision_ends"} {
		var count int
		require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&count), table)
		require.Equal(t, 1, count, table)
	}
}
