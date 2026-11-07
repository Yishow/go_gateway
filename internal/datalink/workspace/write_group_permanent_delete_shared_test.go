package workspace

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteRefusesSharedCutoverProof(t *testing.T) {
	for _, tc := range []struct {
		name, boundary, code string
	}{
		{"finished", "2026-10-02T12:00:00Z", "WRITE_GROUP_DELETE_HAS_SUCCESSOR"},
		{"unfinished", "2026-10-02T12:01:00Z", "WRITE_GROUP_DELETE_OPERATION_PENDING"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, service, group := permanentDeleteSharedFixture(t)
			for _, owner := range []struct{ lineage, id string }{{"L1", group.ID}, {"L2", "G2"}} {
				_, err := db.ExecContext(t.Context(), `INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES(?,?,?)`, owner.lineage, group.WorkspaceID, owner.id)
				require.NoError(t, err)
				_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES(?,?,1)`, owner.id, owner.lineage)
				require.NoError(t, err)
			}
			_, err := db.ExecContext(t.Context(), `INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) VALUES(?,?,'basic',?,'L1')`, group.WorkspaceID, group.Members[0].DeviceID, group.ID)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_profile_keys(workspace_id,selection_scope,profile_kind,canonical_role,group_id,lineage_id) VALUES(?,?,'instantaneous','unassigned','G2','L2')`, group.WorkspaceID, group.Members[0].DeviceID)
			require.NoError(t, err)
			for _, owner := range []struct{ lineage, id string }{{"L1", group.ID}, {"L2", "G2"}} {
				lineage, found, err := service.repo.loadWriteGroupLineage(t.Context(), db, group.WorkspaceID, owner.id)
				require.NoError(t, err)
				require.True(t, found)
				require.Equal(t, owner.lineage, lineage.LineageID)
			}
			_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('C1',?,'forward',?,?, 'review','{}','{}')`, group.WorkspaceID, group.ID, tc.boundary)
			require.NoError(t, err)
			for _, owner := range []struct{ direction, id, revision string }{{"ended", group.ID, "R1"}, {"started", "G2", "R2"}} {
				_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_ownership_fences(fence_revision,workspace_id,direction,group_id,group_revision,device_id,point_id,tag_id,profile_kind,connector_id,connector_revision,database_name,table_schema,table_name) VALUES('C1',?, ?,?,?,'D1','P1','T1','','connector','revision','db','main','table')`, group.WorkspaceID, owner.direction, owner.id, owner.revision)
				require.NoError(t, err)
			}
			before := lifecycleReadState(t, db)
			const cutoverQuery = `SELECT json_array(id,workspace_id,kind,source_group_id,boundary,prior_cutover_id,review_digest,detail,created_at,intent) FROM write_group_cutovers WHERE id='C1'`
			var cutoverBefore string
			require.NoError(t, db.QueryRowContext(t.Context(), cutoverQuery).Scan(&cutoverBefore))
			_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
			require.ErrorContains(t, err, tc.code)
			if tc.name == "finished" {
				var refusal *WriteGroupDeleteError
				require.True(t, errors.As(err, &refusal))
				references, ok := refusal.Details.(WriteGroupDeleteReferences)
				require.True(t, ok)
				require.Equal(t, WriteGroupDeleteReferences{
					ReferencedBy: []WriteGroupDeleteReference{{ID: "G2", Name: "other recording"}},
					Total:        1,
				}, references)
			}
			require.Equal(t, before, lifecycleReadState(t, db))
			var cutoverAfter string
			require.NoError(t, db.QueryRowContext(t.Context(), cutoverQuery).Scan(&cutoverAfter))
			require.Equal(t, cutoverBefore, cutoverAfter)
		})
	}
}

func TestPermanentDeleteRefusesLegacyCutoverWithoutFenceRows(t *testing.T) {
	db, service, group := permanentDeleteSharedFixture(t)
	detail := `{"ended":[{"group_id":"G1","revision":"R1"}],"started":[{"group_id":"G2","revision":"R2"}]}`
	_, err := db.ExecContext(t.Context(), `INSERT INTO write_group_cutovers(id,workspace_id,kind,source_group_id,boundary,review_digest,detail,intent) VALUES('C1',?,'forward',?,'2026-10-02T12:00:00Z','review',?,'{}')`, group.WorkspaceID, group.ID, detail)
	require.NoError(t, err)
	before := lifecycleReadState(t, db)
	_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_HAS_SUCCESSOR")
	require.Equal(t, before, lifecycleReadState(t, db))
}

func TestPermanentDeleteRefusesSharedHandoffProof(t *testing.T) {
	for _, tc := range []struct {
		state, code string
		targetSeq   any
	}{
		{"pending", "WRITE_GROUP_DELETE_OPERATION_PENDING", nil},
		{"finalized", "WRITE_GROUP_DELETE_HAS_SUCCESSOR", 2},
	} {
		t.Run(tc.state, func(t *testing.T) {
			db, service, group := permanentDeleteSharedFixture(t)
			_, err := db.ExecContext(t.Context(), `INSERT INTO wg_delivery_sample_handoffs(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record) VALUES('H1',?,?,'R1',1,'G2','R2',?,1,?,'{}')`, group.WorkspaceID, group.ID, tc.targetSeq, tc.state)
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
			require.ErrorContains(t, err, tc.code)
			require.Equal(t, before, lifecycleReadState(t, db))
		})
	}
}

func permanentDeleteSharedFixture(t *testing.T) (*sql.DB, *WriteGroupService, *WriteGroup) {
	t.Helper()
	db, service, group := lifecyclePreservationFixture(t)
	service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
	require.NoError(t, err)
	other := cloneWriteGroup(disabled.Group)
	other.ID, other.Revision, other.AppliedRevision, other.Name = "G2", "R2", "", "other recording"
	require.NoError(t, service.repo.insertGroup(t.Context(), db, other))
	return db, service, disabled.Group
}

func TestPermanentDeleteRefusesPredecessorBacklog(t *testing.T) {
	db, service, predecessor := lifecyclePreservationFixture(t)
	seedLifecycleAcceptedData(t, db, predecessor)
	_, err := db.ExecContext(t.Context(), `DELETE FROM managed_schema_operations`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE write_groups SET status='deleted' WHERE id='G1'`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_lineages(lineage_id,workspace_id,current_group_id) VALUES('L1',?,'G1')`, predecessor.WorkspaceID)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_generations(group_id,lineage_id,generation) VALUES('G1','L1',1)`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_basic_keys(workspace_id,device_id,canonical_role,group_id,lineage_id) VALUES(?,?,'basic','G1','L1')`, predecessor.WorkspaceID, predecessor.Members[0].DeviceID)
	require.NoError(t, err)
	head := cloneWriteGroup(predecessor)
	head.ID, head.Revision, head.AppliedRevision, head.Status = "G2", "R2", "", WriteGroupStatusDisabled
	require.NoError(t, service.repo.insertGroup(t.Context(), db, head))
	_, err = db.ExecContext(t.Context(), `INSERT INTO write_group_generations(group_id,lineage_id,generation,predecessor_group_id) VALUES('G2','L1',2,'G1')`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE write_group_lineages SET current_group_id='G2',revision=2 WHERE lineage_id='L1'`)
	require.NoError(t, err)
	before := lifecycleReadState(t, db)
	_, err = service.Delete(t.Context(), head.ID, lifecycleReadMutation(t, service, head))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_BACKLOG")
	require.Equal(t, before, lifecycleReadState(t, db))
}
