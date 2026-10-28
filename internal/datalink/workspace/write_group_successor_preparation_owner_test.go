package workspace

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessorPreparationAllowsOnlyVerifiedPredecessorOverlap(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	seedRecreateEngineEvidence(t, h, req)
	_, err := h.db.Exec("UPDATE write_groups SET status='ready' WHERE id=?", req.PredecessorGroupID)
	require.NoError(t, err)
	predecessor, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
	require.NoError(t, err)
	candidate := cloneWriteGroup(req.Group)
	candidate.ID = ""
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	require.NoError(t, h.service.repo.validateSuccessorDefinitionInTx(t.Context(), tx, predecessor, candidate))
	require.ErrorIs(t, h.service.repo.validate(t.Context(), tx, candidate), ErrWriteGroupCircuitProfileConflict, "ordinary create must still refuse overlap")
}

func TestSuccessorPreparationRejectsUnrelatedOwnerAndInvalidAuthority(t *testing.T) {
	for _, mode := range []string{"third-owner", "basic-active", "missing-applied", "foreign-head", "subset"} {
		t.Run(mode, func(t *testing.T) {
			kind := "profile"
			if mode == "basic-active" {
				kind = "basic"
			}
			h, req := recreateRepositoryFixture(t, kind)
			seedRecreateEngineEvidence(t, h, req)
			_, err := h.db.Exec("UPDATE write_groups SET status='ready' WHERE id=?", req.PredecessorGroupID)
			require.NoError(t, err)
			pred, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			candidate := cloneWriteGroup(req.Group)
			candidate.ID = ""
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			switch mode {
			case "third-owner":
				other := cloneWriteGroup(pred)
				other.ID = "unrelated-owner"
				other.Revision = "other-rev"
				other.AppliedRevision = ""
				require.NoError(t, h.service.repo.insertGroup(t.Context(), tx, other))
			case "missing-applied":
				pred.AppliedRevision = ""
			case "foreign-head":
				pred.ID = "foreign"
			case "subset":
				candidate.Members = nil
			}
			err = h.service.repo.validateSuccessorDefinitionInTx(t.Context(), tx, pred, candidate)
			require.Error(t, err)
			if mode == "third-owner" {
				require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
			}
		})
	}
}

func TestSuccessorPreparationRepositoryRejectsUnmintedAuthority(t *testing.T) {
	h, id, request := successorActiveWideFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, request)
	require.NoError(t, err)
	envelope, _, digest, err := normalizeSuccessorEnvelope(review.Envelope)
	require.NoError(t, err)
	req := successorRecreateRequest(envelope)
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = h.service.repo.reserveTransitionDraftInTx(t.Context(), tx, req, digest, successorTransitionAuthority{})
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	var groups int
	require.NoError(t, tx.QueryRow("SELECT COUNT(*) FROM write_groups").Scan(&groups))
	require.Equal(t, 1, groups)
}
