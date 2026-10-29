package workspace

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSuccessorPreparationReservesReviewedIntentAndReplaysAfterProof(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	confirmation := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	saved, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	require.EqualValues(t, 2, saved.Generation)
	require.NotEqual(t, id, saved.Group.ID)
	require.Empty(t, saved.Group.AppliedRevision)
	require.Equal(t, WriteGroupStatusDraft, saved.Group.Status)
	original := lineageReadSnapshot(t, h)
	replay, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	require.Equal(t, saved, replay)
	require.Equal(t, original, lineageReadSnapshot(t, h))
	_, err = h.db.Exec("UPDATE write_groups SET revision='schema-R2',destination_schema_revision='proof-R',destination_schema_digest='proof-D' WHERE id=?", saved.Group.ID)
	require.NoError(t, err)
	_, err = h.service.workspaceSvc.UpdateDatabaseSetup(t.Context(), saved.WorkspaceRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
	require.NoError(t, err)
	replay, err = h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	require.Equal(t, saved.Group.ID, replay.Group.ID)
	require.Equal(t, saved, replay)
	changed := confirmation
	changed.Envelope.Request.Boundary = changed.Envelope.Request.Boundary.Add(1)
	_, err = h.service.ReserveSuccessor(t.Context(), id, changed)
	require.Error(t, err)
	require.Equal(t, 2, h.countRows("SELECT COUNT(*) FROM write_groups"))
}

func TestSuccessorPreparationRejectsStaleReviewWithoutEffects(t *testing.T) {
	for _, mode := range []string{"digest", "boundary", "workspace", "foreign", "source"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := successorReviewFixture(t)
			review, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			c := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
			switch mode {
			case "digest":
				c.ReviewDigest = "foreign"
			case "boundary":
				c.Envelope.Request.Boundary = c.Envelope.Request.Boundary.Add(1)
			case "workspace":
				_, err = h.service.workspaceSvc.UpdateDatabaseSetup(t.Context(), req.ExpectedWorkspaceRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
				require.NoError(t, err)
			case "foreign":
				id = "foreign"
			case "source":
				c.Envelope.SourceDefinitionDigest = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			}
			before := lineageReadSnapshot(t, h)
			_, err = h.service.ReserveSuccessor(t.Context(), id, c)
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}
