package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermanentDeleteChecksRecordingStartLedger(t *testing.T) {
	for _, tc := range []struct {
		status  string
		blocked bool
	}{
		{"pending", true}, {"running", true}, {"partial", true}, {"unknown", true},
		{"succeeded", false}, {"failed", false},
	} {
		t.Run(tc.status, func(t *testing.T) {
			db, service, group := lifecyclePreservationFixture(t)
			service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
			require.NoError(t, err)
			detail, err := json.Marshal(recordingStartProgress{Version: 1, Request: RecordingStartRequest{
				WorkspaceID: group.WorkspaceID, Groups: []RecordingStartGroupIntent{{GroupID: group.ID}},
			}})
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('start-1','start-token',?,'start-scope','owner','recording_start',?,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, group.WorkspaceID, tc.status, string(detail))
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
			if tc.blocked {
				require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
				require.Equal(t, before, lifecycleReadState(t, db))
			} else {
				require.NoError(t, err)
			}
			var operations int
			require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='start-1'`).Scan(&operations))
			require.Equal(t, 1, operations, "start evidence remains even after a successful delete")
		})
	}
}

func TestPermanentDeleteRecordingStartScopeFailsClosed(t *testing.T) {
	for _, tc := range []struct {
		name, groups string
		blocked      bool
	}{
		{"other recording", `[{"group_id":"other"}]`, false},
		{"device only", `null`, false},
		{"missing group identity", `[{"group_id":""}]`, true},
		{"malformed groups", `{}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, service, group := lifecyclePreservationFixture(t)
			service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
			disabled, err := service.Disable(t.Context(), group.ID, lifecycleReadMutation(t, service, group))
			require.NoError(t, err)
			detail := `{"version":1,"request":{"workspace_id":"` + group.WorkspaceID + `","groups":` + tc.groups + `}}`
			_, err = db.ExecContext(t.Context(), `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('start-1','start-token',?,'start-scope','owner','recording_start','running',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, group.WorkspaceID, detail)
			require.NoError(t, err)
			before := lifecycleReadState(t, db)
			_, err = service.Delete(t.Context(), group.ID, lifecycleReadMutation(t, service, disabled.Group))
			if tc.blocked {
				require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
				require.Equal(t, before, lifecycleReadState(t, db))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPermanentDeleteRacingRecordingStartRollsBack(t *testing.T) {
	dsn := "file:" + filepath.Join(t.TempDir(), "recording-start-race.db") + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(1000)&_pragma=foreign_keys(1)"
	h := newCircuitProfileHarness(t, dsn)
	h.db.SetMaxOpenConns(1)
	other, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, other.Close()) })
	require.NoError(t, other.PingContext(t.Context()))
	group := h.mustEnsure(h.longDraft(WriteGroupProfileInstantaneous, "main", h.primary()))
	current, err := h.service.Get(t.Context(), group.ID)
	require.NoError(t, err)
	detail := `{"version":1,"request":{"workspace_id":"` + group.WorkspaceID + `","groups":[{"group_id":"` + group.ID + `"}]}}`
	_, err = h.workspace.UpdateDatabaseSetup(t.Context(), current.WorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		inTx, err := h.service.repo.GetInTx(ctx, tx, group.WorkspaceID, group.ID)
		if err != nil {
			return err
		}
		if err := h.service.checkPermanentDeleteOperations(ctx, tx, inTx); err != nil {
			return err
		}
		// Another connection commits a start after this transaction's read.
		_, err = other.ExecContext(ctx, `INSERT INTO managed_schema_operations(operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail) VALUES('racing-start','racing-token',?,'racing-scope','owner','recording_start','running',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`, group.WorkspaceID, detail)
		if err != nil {
			return err
		}
		_, err = h.service.permanentDeleteInTx(ctx, tx, record, inTx)
		return err
	})
	require.Error(t, err, "SQLite must not commit deletion from a stale read snapshot")
	require.ErrorContains(t, err, "database is locked", "the concurrent commit must invalidate the delete transaction")
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups WHERE id=?`, group.ID))
	require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE operation_id='racing-start'`))
	_, err = h.service.Delete(t.Context(), group.ID, h.mutationFor(current))
	require.ErrorContains(t, err, "WRITE_GROUP_DELETE_OPERATION_PENDING")
}
