package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteGroupSuccessorRoutesRejectMissingIntent(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	for _, suffix := range []string{"successor-review", "successor"} {
		t.Run(suffix, func(t *testing.T) {
			before := lifecycleAPIState(t, f.db)
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/group-missing/"+suffix, map[string]any{})
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.Equal(t, before, lifecycleAPIState(t, f.db))
		})
	}
}

func TestWriteGroupSuccessorReviewIsReadOnlyAndExact(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	inspector := &successorAPIForbiddenInspector{}
	f.f.groups.WithTableInspector(inspector).WithManagedTableInspector(inspector)
	before := successorAPIAllState(t, f.f.db)
	review := f.review(t)
	require.Len(t, review.ReviewDigest, 64)
	require.Equal(t, f.groupID, review.Envelope.PredecessorGroupID)
	require.EqualValues(t, 1, review.Envelope.Generation)
	require.Equal(t, f.snapshot.Group.Revision, review.Envelope.SourceAppliedRevision)
	require.True(t, f.request.Boundary.Equal(review.Envelope.Request.Boundary))
	require.Empty(t, review.Envelope.Request.Group.ID)
	require.Empty(t, review.Envelope.Request.Group.AppliedRevision)
	require.Empty(t, review.Envelope.Request.Group.Destination.SchemaDigest)
	require.NotNil(t, review.History)
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "opening or cancelling review must not save a token, reservation, ledger or any configuration")
	require.Zero(t, inspector.calls)
	again := f.review(t)
	require.Equal(t, review.ReviewDigest, again.ReviewDigest)
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
}
