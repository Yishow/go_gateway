package workspace

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSuccessorActivationHasExplicitAtomicConfirmation(t *testing.T) {
	h, _, _ := successorActiveWideFixture(t)
	method := reflect.ValueOf(h.service).MethodByName("ConfirmSuccessorActivation")
	require.True(t, method.IsValid(), "prepared successor needs explicit atomic confirmation, never ordinary Apply")
}

func TestHistoricalCanonicalGenerationCannotPassOrdinaryApplyGuard(t *testing.T) {
	h, id, req := successorActiveWideFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?, reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1 WHERE lineage_id=?`, saved.Group.ID, saved.LineageID)
	require.NoError(t, err)
	tx, err := h.db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	// The ordinary entry guard itself must deny a historical canonical owner,
	// even if its successor is disabled or has no overlapping active writer.
	record, err := readWorkspaceScope(t.Context(), tx, h.service.repo.query, req.WorkspaceID)
	require.NoError(t, err)
	_, err = h.service.applyGroupInTx(t.Context(), tx, record, id, WriteGroupMutation{WorkspaceID: req.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, ExpectedGroupRevision: req.ExpectedPredecessorRevision, ExpectedConnectorRevision: req.ExpectedPredecessorConnectorRevision}, &WriteGroupReadiness{}, func(*WriteGroup, *WriteGroupAppliedSnapshot) (time.Time, error) { return req.Boundary, nil }, false)
	var action *WriteGroupLifecycleActionError
	require.True(t, errors.As(err, &action), "historical owner must be rejected by explicit ownership gate, got %v", err)
	if action != nil {
		require.Equal(t, "historical_generation", action.ReasonCode)
	}
}
