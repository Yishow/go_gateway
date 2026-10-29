package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestRecreatePendingReplayPrecedesFreshCAS(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	beforeWorkspace, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	first, err := h.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err)
	require.NotEqual(t, req.ExpectedWorkspaceRevision, first.WorkspaceRevision)
	afterWorkspace, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, beforeWorkspace.DatabaseTargetRefs, afterWorkspace.DatabaseTargetRefs, "reservation must not retarget predecessor refs")
	require.Len(t, afterWorkspace.DatabaseRowGroups, len(beforeWorkspace.DatabaseRowGroups)+1)
	before := lineageReadSnapshot(t, h)
	replay, err := h.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err, "identical pending replay must succeed before stale new-mutation CAS")
	require.Equal(t, first, replay)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestRecreatePendingReplayRejectsChangedIntent(t *testing.T) {
	changes := map[string]func(*writeGroupRecreateReservationRequest){
		"operation":           func(r *writeGroupRecreateReservationRequest) { r.OperationID = "O3" },
		"name":                func(r *writeGroupRecreateReservationRequest) { r.Group.Name = "changed" },
		"workspace_cas":       func(r *writeGroupRecreateReservationRequest) { r.ExpectedWorkspaceRevision = "new-cas" },
		"lineage_cas":         func(r *writeGroupRecreateReservationRequest) { r.ExpectedLineageRevision++ },
		"predecessor_cas":     func(r *writeGroupRecreateReservationRequest) { r.ExpectedPredecessorRevision = "new-cas" },
		"saved_connector_cas": func(r *writeGroupRecreateReservationRequest) { r.ExpectedPredecessorConnectorRevision = "new-cas" },
		"target_connector_cas": func(r *writeGroupRecreateReservationRequest) {
			r.ExpectedTargetConnectorRevision = "new-cas"
			r.Group.Destination.ConnectorRevision = "new-cas"
		},
		"destination": func(r *writeGroupRecreateReservationRequest) { r.Group.Destination.Database = "different" },
		"source":      func(r *writeGroupRecreateReservationRequest) { r.Group.Members[0].SourceRevision = "different" },
		"policy":      func(r *writeGroupRecreateReservationRequest) { r.Group.RowPolicy.AllowedLatenessSeconds++ },
		"foreign_workspace": func(r *writeGroupRecreateReservationRequest) {
			r.WorkspaceID = "foreign"
			r.Group.WorkspaceID = "foreign"
		},
		"foreign_lineage": func(r *writeGroupRecreateReservationRequest) { r.LineageID = "foreign" },
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			_, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			req.Group = cloneWriteGroup(req.Group)
			change(&req)
			_, err = h.service.reserveRecreatedDraft(t.Context(), req)
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			var conflict *writeGroupReservationConflictError
			require.ErrorAs(t, err, &conflict)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestRecreateReplayIgnoresLiveTargetAndSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "recreate.db")
	h, req := recreateServiceFixture(t, path)
	first, err := h.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err)
	_, err = h.workspace.UpdateDatabaseSetup(t.Context(), first.WorkspaceRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE database_connectors SET identity_revision='edited',enabled=0,connection_config='{}'`)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE mappings SET enabled=0`)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	require.NoError(t, h.db.Close())
	restarted := newCircuitProfileHarnessOn(t, openWorkspaceTestDB(t, path), h)
	restarted.db.SetMaxOpenConns(1)
	replay, err := restarted.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, first.Group, replay.Group)
	require.Equal(t, first.Generation, replay.Generation)
	require.NotEqual(t, first.WorkspaceRevision, replay.WorkspaceRevision)
	require.Equal(t, before, lineageReadSnapshot(t, restarted))
}

func TestRecreateReplayRejectsCorruptPendingMetadata(t *testing.T) {
	cases := map[string]string{
		"status":        `UPDATE write_groups SET status='ready' WHERE id=(SELECT reserved_group_id FROM write_group_lineages)`,
		"applied":       `UPDATE write_groups SET applied_revision='bad' WHERE id=(SELECT reserved_group_id FROM write_group_lineages)`,
		"revision":      `UPDATE write_group_lineages SET revision=revision+1`,
		"predecessor":   `UPDATE write_group_generations SET predecessor_group_id='missing' WHERE generation=2`,
		"partial_tuple": `UPDATE write_group_lineages SET reservation_operation_id=NULL`,
	}
	for name, corruption := range cases {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			_, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.NoError(t, err)
			_, err = h.db.Exec(`PRAGMA foreign_keys=OFF; PRAGMA ignore_check_constraints=ON`)
			require.NoError(t, err)
			_, err = h.db.Exec(corruption)
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.reserveRecreatedDraft(t.Context(), req)
			require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
			require.NotErrorIs(t, err, ErrWriteGroupNotFound)
			require.NotErrorIs(t, err, ErrNotFound)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestRecreateFrozenRetryAfterTransferredHeadCannotAllocate(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	first, err := h.service.reserveRecreatedDraft(t.Context(), req)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=reserved_group_id,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1`)
	require.NoError(t, err)
	before := lineageReadSnapshot(t, h)
	_, err = h.service.reserveRecreatedDraft(t.Context(), req)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Equal(t, before, lineageReadSnapshot(t, h))
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_group_generations`))
	require.NotEmpty(t, first.Group.ID)
}

func recreateServiceFixture(t *testing.T, path string) (*circuitProfileHarness, writeGroupRecreateReservationRequest) {
	t.Helper()
	h := newCircuitProfileHarness(t, path)
	h.db.SetMaxOpenConns(1)
	draft := lineageTestDraft(h, "profile")
	draft.LayoutVersion = WriteGroupLayoutManagedLongV1
	saved := h.mustEnsure(draft)
	removed, err := h.service.Delete(t.Context(), saved.ID, lifecycleReadMutation(t, h.service, saved))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, removed.Group)
	return h, writeGroupRecreateReservationRequest{WorkspaceID: removed.Group.WorkspaceID, ExpectedWorkspaceRevision: removed.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: removed.Group.ID, ExpectedPredecessorRevision: removed.Group.Revision,
		ExpectedPredecessorConnectorRevision: removed.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: removed.Group.Destination.ConnectorRevision,
		OperationID: "O2", Group: cloneWriteGroup(removed.Group)}
}
