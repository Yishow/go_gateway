package workspace

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSuccessorReviewRejectsChangedCASWithoutEffects(t *testing.T) {
	mutations := map[string]func(*WriteGroupSuccessorReviewRequest){
		"workspace":          func(r *WriteGroupSuccessorReviewRequest) { r.WorkspaceID = "foreign"; r.Group.WorkspaceID = "foreign" },
		"workspace-revision": func(r *WriteGroupSuccessorReviewRequest) { r.ExpectedWorkspaceRevision = "old" },
		"lineage":            func(r *WriteGroupSuccessorReviewRequest) { r.LineageID = "foreign" },
		"lineage-revision":   func(r *WriteGroupSuccessorReviewRequest) { r.ExpectedLineageRevision++ },
		"group-revision":     func(r *WriteGroupSuccessorReviewRequest) { r.ExpectedPredecessorRevision = "old" },
		"source-connector":   func(r *WriteGroupSuccessorReviewRequest) { r.ExpectedPredecessorConnectorRevision = "old" },
		"target-connector": func(r *WriteGroupSuccessorReviewRequest) {
			r.ExpectedTargetConnectorRevision = "old"
			r.Group.Destination.ConnectorRevision = "old"
		},
	}
	for name, change := range mutations {
		t.Run(name, func(t *testing.T) {
			h, id, req := successorReviewFixture(t)
			change(&req)
			before := lineageReadSnapshot(t, h)
			_, err := h.service.SuccessorReview(t.Context(), id, req)
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorReviewRequiresAppliedOwnerForActiveProfile(t *testing.T) {
	for _, mode := range []string{"no-applied", "no-owner", "owner-error", "trusted-owner", "bad-source"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := successorReviewFixture(t)
			_, err := h.db.Exec(`UPDATE write_groups SET status='ready' WHERE id=?`, id)
			require.NoError(t, err)
			old := writeGroupRecreateReservationRequest{WorkspaceID: req.WorkspaceID, PredecessorGroupID: id, ExpectedPredecessorRevision: req.ExpectedPredecessorRevision, Group: req.Group}
			if mode != "no-applied" {
				seedRecreateEngineEvidence(t, h, old)
			}
			if mode == "bad-source" {
				_, err = h.db.Exec(`UPDATE write_group_versions SET payload='{}'`)
				require.NoError(t, err)
			}
			if mode == "trusted-owner" || mode == "owner-error" {
				h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error {
					if mode == "owner-error" {
						return errors.New("private password target path")
					}
					return nil
				}))
			}
			// A distinct destination is required for an active same-layout replacement.
			target := old
			target.Group = cloneWriteGroup(req.Group)
			recreateRepositoryNewTarget(t, h, &target)
			req.Group = target.Group
			req.ExpectedTargetConnectorRevision = target.ExpectedTargetConnectorRevision
			before := lineageReadSnapshot(t, h)
			got, err := h.service.SuccessorReview(t.Context(), id, req)
			if mode == "bad-source" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if mode == "trusted-owner" {
					require.Empty(t, got.Blockers)
				} else {
					require.NotEmpty(t, got.Blockers)
				}
			}
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorReviewDoesNotAuthorizeActiveBasic(t *testing.T) {
	h, r := recreateBasicFixture(t, "sqlite")
	_, err := h.db.Exec(`UPDATE write_groups SET status='ready' WHERE id=?`, r.PredecessorGroupID)
	require.NoError(t, err)
	req := WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: "O2", Boundary: time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC), Group: r.Group}
	before := lineageReadSnapshot(t, h)
	got, err := h.service.SuccessorReview(t.Context(), r.PredecessorGroupID, req)
	require.NoError(t, err)
	require.Contains(t, got.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "active_basic_unsupported", NextAction: "review_lifecycle"})
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorReviewUsesCallerSnapshot(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	original, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	store := h.service.workspaceSvc.repo.(DatabaseSetupStore)
	st, err := store.BeginDatabaseSetup(ctx)
	require.NoError(t, err)
	defer st.Rollback()
	record, err := st.Get(ctx)
	require.NoError(t, err)
	_, err = st.SQLTx().ExecContext(ctx, `UPDATE write_groups SET name='uncommitted source name' WHERE id=?`, id)
	require.NoError(t, err)
	view, err := h.service.reviewSuccessorInTx(ctx, st.SQLTx(), record, id, req)
	require.NoError(t, err)
	require.NotEqual(t, original.Envelope.SourceDefinitionDigest, view.Envelope.SourceDefinitionDigest, "review must read its caller snapshot, not a second DB connection")
}

func TestSuccessorReviewRejectsActiveLongNoopDestination(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	old := writeGroupRecreateReservationRequest{WorkspaceID: req.WorkspaceID, PredecessorGroupID: id, ExpectedPredecessorRevision: req.ExpectedPredecessorRevision, Group: req.Group}
	seedRecreateEngineEvidence(t, h, old)
	_, err := h.db.Exec(`UPDATE write_groups SET status='ready' WHERE id=?`, id)
	require.NoError(t, err)
	h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	before := lineageReadSnapshot(t, h)
	got, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Contains(t, got.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "destination-unchanged", NextAction: "review_lifecycle"})
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
