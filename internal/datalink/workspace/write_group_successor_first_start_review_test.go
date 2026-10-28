package workspace

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSuccessorFirstStartReviewDoesNotInventAppliedPredecessor(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Now().Add(time.Hour)})
	require.NoError(t, err)
	require.NotContains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "successor_first_activation_unavailable", NextAction: "review_lifecycle"})
	raw, err := json.Marshal(view.Envelope)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"activation_kind":"first_start"`)
	require.Empty(t, view.Envelope.SourceAppliedRevision)
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_versions"))
}

func TestSuccessorFirstStartHasDistinctConfirmation(t *testing.T) {
	type firstStarter interface {
		ConfirmSuccessorFirstStart(context.Context, string, WriteGroupSuccessorActivationConfirmation) (*WriteGroupSuccessorFirstStartResult, error)
	}
	_, ok := any(&WriteGroupService{}).(firstStarter)
	require.True(t, ok, "first-start must use its own typed action, not an empty predecessor cutover")
}
