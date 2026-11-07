package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteChecksAndRemovesOwnedSchemaOperations(t *testing.T) {
	for _, tc := range []struct {
		name, status, write, cleanup string
		executed                     int
		blocked                      bool
	}{
		{name: "running", status: "running", blocked: true},
		{name: "succeeded", status: "succeeded"},
		{name: "succeeded_test_write_cleaned", status: "succeeded", write: "written_verified", cleanup: "cleaned"},
		{name: "failed_no_effect", status: "failed"},
		{name: "failed_test_write", status: "failed", write: "failed", cleanup: "not_attempted"},
		{name: "failed_test_write_cleaned", status: "failed", write: "written_unverified", cleanup: "cleaned"},
		{name: "partial", status: "partial", blocked: true},
		{name: "unknown", status: "unknown", blocked: true},
		{name: "failed_with_effects", status: "failed", executed: 1, blocked: true},
		{name: "failed_cleanup_unknown", status: "failed", write: "unknown", cleanup: "unknown", blocked: true},
		{name: "succeeded_cleanup_failed", status: "succeeded", write: "written_verified", cleanup: "failed", blocked: true},
		{name: "succeeded_cleanup_not_attempted", status: "succeeded", write: "written_verified", cleanup: "not_attempted", blocked: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			h.db.SetMaxOpenConns(1)
			created, err := h.service.Create(t.Context(), h.mutation(t.Context(), newBasicWriteGroup(h.fixture)))
			require.NoError(t, err)
			_, err = h.db.ExecContext(t.Context(), `INSERT INTO managed_schema_preview_tokens(token,workspace_id,plan_id,plan_revision,connector_id,table_prefix,statements,expires_at,operation_id,action) VALUES('T1',?,?,?,'connector','table','[]',CURRENT_TIMESTAMP,'O1','schema_apply')`, h.fixture.workspaceID, created.Group.ID, created.Group.Revision)
			require.NoError(t, err)
			_, err = h.db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,executed_statements,write_outcome,cleanup_status,created_at,updated_at) VALUES('O1','T1',?,'scope','owner','schema_apply',?,?,?,?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, h.fixture.workspaceID, tc.status, tc.executed, tc.write, tc.cleanup)
			require.NoError(t, err)
			before := lifecycleReadState(t, h.db)
			_, err = h.service.Delete(t.Context(), created.Group.ID, h.mutationFor(created))
			if tc.blocked {
				require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
				require.Equal(t, before, lifecycleReadState(t, h.db))
				require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM managed_schema_preview_tokens WHERE token='T1'`))
				require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='O1'`))
				return
			}
			require.NoError(t, err)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM managed_schema_preview_tokens WHERE token='T1'`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='O1'`))
		})
	}
}
