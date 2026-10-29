package workspace

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func successorReviewFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest) {
	t.Helper()
	h, r := recreateRepositoryFixture(t, "profile")
	return h, r.PredecessorGroupID, WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: r.ExpectedLineageRevision, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: "O2", Boundary: time.Date(2026, 10, 28, 10, 0, 0, 123, time.UTC), Group: r.Group}
}

func TestSuccessorReviewIsLocalCoherentAndNonMutating(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	before := lineageReadSnapshot(t, h)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	type reviewer interface {
		SuccessorReview(context.Context, string, WriteGroupSuccessorReviewRequest) (*WriteGroupSuccessorReview, error)
	}
	api, ok := any(h.service).(reviewer)
	require.True(t, ok, "service must expose a read-only successor review, not reservation")
	view, err := api.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	require.Equal(t, "successor-review-v1", view.Envelope.Version)
	require.Equal(t, id, view.Envelope.PredecessorGroupID)
	require.EqualValues(t, 1, view.Envelope.Generation)
	require.Equal(t, req.Boundary, view.Envelope.Request.Boundary)
	require.Len(t, view.ReviewDigest, 64)
	require.Len(t, view.Envelope.SourceDefinitionDigest, 64)
	require.NotNil(t, view.History)
	require.Empty(t, view.Envelope.Request.Group.ID)
	require.Empty(t, view.Envelope.Request.Group.Revision)
	require.Equal(t, before, lineageReadSnapshot(t, h))
	require.Zero(t, probe.calls)
}

func TestSuccessorReviewDeletedBasicKeepsUnallocatedIdentity(t *testing.T) {
	h, r := recreateBasicFixture(t, "sqlite")
	req := WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: "O2", Boundary: time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC), Group: r.Group}
	before := lineageReadSnapshot(t, h)
	got, err := h.service.SuccessorReview(t.Context(), r.PredecessorGroupID, req)
	require.NoError(t, err)
	require.Empty(t, got.Blockers)
	require.Empty(t, got.Envelope.Request.Group.ID)
	require.Empty(t, got.Envelope.Request.Group.Revision)
	require.Empty(t, got.Envelope.Request.Group.Destination.TableName)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorReviewRejectsCallerAllocatedWideStorage(t *testing.T) {
	h, r := recreateBasicFixture(t, "sqlite")
	req := WriteGroupSuccessorReviewRequest{WorkspaceID: r.WorkspaceID, ExpectedWorkspaceRevision: r.ExpectedWorkspaceRevision, LineageID: r.LineageID, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: r.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: r.ExpectedPredecessorConnectorRevision, ExpectedTargetConnectorRevision: r.ExpectedTargetConnectorRevision, OperationID: "O2", Boundary: time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC), Group: r.Group}
	req.Group.Destination.TableName = "another_owned_table"
	before := lineageReadSnapshot(t, h)
	got, err := h.service.SuccessorReview(t.Context(), r.PredecessorGroupID, req)
	require.True(t, err != nil || len(got.Blockers) > 0, "review must not approve caller-supplied physical wide ownership")
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
