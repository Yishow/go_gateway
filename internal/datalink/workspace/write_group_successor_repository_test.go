package workspace

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func recreateRepositoryFixture(t *testing.T, kind string) (*circuitProfileHarness, writeGroupRecreateReservationRequest) {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	group := lineageTestDraft(h, kind)
	// Both canonical aliases can belong to a supported profile layout.
	if kind != "basic" {
		group.LayoutVersion = WriteGroupLayoutManagedLongV1
	}
	saved, err := ensureLineageKind(h, kind, h.mutation(t.Context(), group))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, saved.Group)
	_, err = h.db.Exec("UPDATE write_groups SET status='deleted' WHERE id=?", saved.Group.ID)
	require.NoError(t, err)
	candidate := cloneWriteGroup(saved.Group)
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	return h, writeGroupRecreateReservationRequest{
		WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: saved.Group.ID,
		ExpectedPredecessorRevision:          saved.Group.Revision,
		ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision,
		ExpectedTargetConnectorRevision:      saved.Group.Destination.ConnectorRevision,
		OperationID:                          "O2", Group: candidate,
	}
}

func reserveRepositoryTest(t *testing.T, h *circuitProfileHarness, req writeGroupRecreateReservationRequest) (*writeGroupReservedDraft, error) {
	t.Helper()
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	result, err := h.service.repo.reserveRecreatedDraftInTx(t.Context(), tx, req, "intent-O2")
	if err == nil {
		err = tx.Commit()
	}
	return result, err
}

func TestRecreateRepositoryReservesNextDraftWithoutMovingAliases(t *testing.T) {
	for _, kind := range []string{"profile", "dual"} {
		t.Run(kind, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, kind)
			before, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			preserved := recreateRepositorySnapshot(t, h)
			result, err := reserveRepositoryTest(t, h, req)
			require.NoError(t, err)
			require.NotEqual(t, req.PredecessorGroupID, result.Group.ID)
			require.NotEqual(t, req.ExpectedPredecessorRevision, result.Group.Revision)
			require.Equal(t, WriteGroupStatusDraft, result.Group.Status)
			require.Empty(t, result.Group.AppliedRevision)
			require.EqualValues(t, 2, result.Generation)
			require.EqualValues(t, 2, result.LineageRevision)
			require.Equal(t, req.LineageID, result.LineageID)
			require.Equal(t, "O2", result.OperationID)
			require.Equal(t, "intent-O2", result.IntentDigest)
			require.Equal(t, req.PredecessorGroupID, result.PredecessorGroupID)
			var head, reserved, predecessor string
			var generation int64
			require.NoError(t, h.db.QueryRow(`SELECT l.current_group_id,l.reserved_group_id,g.predecessor_group_id,g.generation
 FROM write_group_lineages l JOIN write_group_generations g ON g.group_id=l.reserved_group_id`).Scan(&head, &reserved, &predecessor, &generation))
			require.Equal(t, req.PredecessorGroupID, head)
			require.Equal(t, result.Group.ID, reserved)
			require.Equal(t, req.PredecessorGroupID, predecessor)
			require.EqualValues(t, 2, generation)
			for _, table := range lineageKeyTables(kind) {
				require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM "+table))
				var anchor string
				require.NoError(t, h.db.QueryRow("SELECT group_id FROM "+table).Scan(&anchor))
				require.Equal(t, req.PredecessorGroupID, anchor)
			}
			after, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			require.Equal(t, before, after)
			current := recreateRepositorySnapshot(t, h)
			for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys", "write_group_versions", "wg_runtime_versions", "write_group_ownership_fences"} {
				require.Equal(t, preserved[table], current[table], "reservation preserves %s bytes and timestamps", table)
			}
			require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_groups"))
			require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_versions WHERE group_id=?", result.Group.ID))
			readback, err := h.service.repo.Get(t.Context(), req.WorkspaceID, result.Group.ID)
			require.NoError(t, err)
			require.Equal(t, result.Group, readback)
			_, err = reserveRepositoryTest(t, h, req)
			require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
			require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_groups"))
		})
	}
}

// Compile-time references keep the direct repository contract transaction-only.
var _ func(context.Context, *sql.Tx, writeGroupRecreateReservationRequest, string) (*writeGroupReservedDraft, error) = (*SQLWriteGroupRepository)(nil).reserveRecreatedDraftInTx

func TestRecreateRepositoryCallerRollbackIncludesNewStorageScope(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	recreateRepositoryNewTarget(t, h, &req)
	before := recreateRepositorySnapshot(t, h)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	result, err := h.service.repo.reserveRecreatedDraftInTx(t.Context(), tx, req, "intent-O2")
	require.NoError(t, err)
	var count int
	require.NoError(t, tx.QueryRow("SELECT COUNT(*) FROM write_group_storage_scopes WHERE connector_id='recreate-target'").Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, tx.QueryRow("SELECT COUNT(*) FROM write_group_generations WHERE group_id=?", result.Group.ID).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, tx.Rollback())
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}
