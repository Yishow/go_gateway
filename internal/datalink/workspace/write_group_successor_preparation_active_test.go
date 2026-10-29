package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func successorActiveWideFixture(t *testing.T, layouts ...WriteGroupLayoutVersion) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest) {
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	draft := lineageTestDraft(h, "profile")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	if len(layouts) > 0 {
		draft.LayoutVersion = layouts[0]
		if draft.isLongLayout() {
			draft.StorageMode = WriteGroupStorageModePlain
		}
	}
	for i := range draft.Members {
		draft.Members[i].ValueRole = WriteGroupValueRoleGauge
	}
	saved, err := ensureLineageKind(h, "profile", h.mutation(t.Context(), draft))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, saved.Group)
	req := writeGroupRecreateReservationRequest{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: saved.Group.ID, ExpectedPredecessorRevision: saved.Group.Revision, ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: saved.Group.Destination.ConnectorRevision, OperationID: "active-O", Group: saved.Group}
	seedRecreateEngineEvidence(t, h, req)
	_, err = h.db.Exec("UPDATE write_groups SET status='ready' WHERE id=?", saved.Group.ID)
	require.NoError(t, err)
	h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	return h, saved.Group.ID, WriteGroupSuccessorReviewRequest{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: req.ExpectedWorkspaceRevision, LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: req.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: req.ExpectedTargetConnectorRevision, OperationID: req.OperationID, Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: cloneWriteGroup(saved.Group)}
}

func TestSuccessorPreparationActiveWideKeepsPredecessorAndNewOwner(t *testing.T) {
	h, id, req := successorActiveWideFixture(t)
	pred, err := h.service.repo.Get(t.Context(), req.WorkspaceID, id)
	require.NoError(t, err)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	c := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	before := recreateRepositorySnapshot(t, h)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, c)
	require.NoError(t, err)
	require.NotEqual(t, pred.Destination.TableName, saved.Group.Destination.TableName)
	require.Empty(t, saved.Group.Destination.SchemaDigest)
	require.Empty(t, saved.Group.AppliedRevision)
	after, err := h.service.repo.Get(t.Context(), req.WorkspaceID, id)
	require.NoError(t, err)
	require.Equal(t, pred, after)
	preserved := recreateRepositorySnapshot(t, h)
	for _, table := range []string{"write_group_versions", "wg_runtime_versions", "write_group_ownership_fences", "write_group_basic_keys", "write_group_profile_keys"} {
		require.Equal(t, before[table], preserved[table], table)
	}
	_, err = h.db.Exec("UPDATE write_groups SET name='new unapplied name',revision='draft-R2' WHERE id=?", id)
	require.NoError(t, err)
	replay, err := h.service.ReserveSuccessor(t.Context(), id, c)
	require.NoError(t, err)
	require.Equal(t, saved.Group.ID, replay.Group.ID)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	record, err := readWorkspaceScope(t.Context(), tx, h.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: c.Envelope, ReviewDigest: c.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: pred.Revision, PredecessorConnectorRevision: pred.Destination.ConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
	_, err = h.service.successorReservationInTx(t.Context(), tx, record, id, stage)
	require.Error(t, err, "replay is not current stage authority")
	stage.Current.PredecessorRevision = "draft-R2"
	_, err = h.service.successorReservationInTx(t.Context(), tx, record, id, stage)
	require.NoError(t, err)
}

func TestSuccessorPreparationActiveMissingBacklogAndSameLongScope(t *testing.T) {
	for _, mode := range []string{"missing-backlog", "same-long-scope"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := successorActiveWideFixture(t)
			if mode == "missing-backlog" {
				h.service.WithBacklogOwnershipGuard(nil)
			} else {
				// Use an already applied long source, whose physical owner is shared by scope.
				h0, r := recreateRepositoryFixture(t, "profile")
				h = h0
				id = r.PredecessorGroupID
				seedRecreateEngineEvidence(t, h, r)
				_, err := h.db.Exec("UPDATE write_groups SET status='ready' WHERE id=?", id)
				require.NoError(t, err)
				h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
				req = WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: "long-O", Boundary: req.Boundary, Group: r.Group}
			}
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorPreparationBidirectionalLayouts(t *testing.T) {
	for _, from := range []WriteGroupLayoutVersion{WriteGroupLayoutManagedWideProfileV1, WriteGroupLayoutManagedLongV1} {
		t.Run(string(from), func(t *testing.T) {
			h, id, req := successorActiveWideFixture(t, from)
			oldTable := req.Group.Destination.TableName
			if from == WriteGroupLayoutManagedWideProfileV1 {
				req.Group.LayoutVersion = WriteGroupLayoutManagedLongV1
				req.Group.StorageMode = WriteGroupStorageModePlain
			} else {
				req.Group.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
				req.Group.StorageMode = ""
				req.Group.WritePolicy.RecreateMissing = nil // Long-only recovery policy is not part of a wide target.
			}
			req.Group.Destination.TableName = ""
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, review.Blockers)
			before := recreateRepositorySnapshot(t, h)
			saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
			require.NoError(t, err)
			require.NotEqual(t, oldTable, saved.Group.Destination.TableName)
			require.Equal(t, req.Group.LayoutVersion, saved.Group.LayoutVersion)
			require.Empty(t, saved.Group.AppliedRevision)
			require.Empty(t, saved.Group.Destination.SchemaDigest)
			after := recreateRepositorySnapshot(t, h)
			for _, table := range []string{"write_group_versions", "wg_runtime_versions", "write_group_ownership_fences"} {
				require.Equal(t, before[table], after[table])
			}
		})
	}
}
