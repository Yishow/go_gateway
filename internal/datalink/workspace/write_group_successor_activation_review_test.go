package workspace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func activationReviewFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest) {
	t.Helper()
	h, id, req := successorActiveWideFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	h.service.now = func() time.Time { return time.Date(2026, 10, 28, 11, 0, 0, 0, time.UTC) }
	stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
	return h, id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}
}

func TestSuccessorActivationReviewIsLocalAndRequiresPreparedSchema(t *testing.T) {
	h, id, req := activationReviewFixture(t)
	type api interface {
		SuccessorActivationReview(context.Context, string, WriteGroupSuccessorActivationRequest) (*WriteGroupSuccessorActivationReview, error)
	}
	svc, ok := any(h.service).(api)
	require.True(t, ok, "activation review must exist independently of preparation")
	before := lineageReadSnapshot(t, h)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	view, err := svc.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-required", NextAction: "review_lifecycle"})
	require.Len(t, view.ReviewDigest, 64)
	require.Zero(t, probe.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
