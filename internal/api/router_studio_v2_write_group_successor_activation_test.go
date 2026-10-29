package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWriteGroupSuccessorActivationReviewIsLocalAndNonMutating(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	before := successorAPIAllState(t, f.f.db)
	inspector := &successorAPIForbiddenInspector{}
	f.f.groups.WithTableInspector(inspector).WithManagedTableInspector(inspector)
	data := f.postActivation(t, f.activationRequest(t))
	require.NotEmpty(t, data["review_digest"])
	require.NotNil(t, data["envelope"])
	require.Empty(t, data["blockers"])
	require.Zero(t, inspector.calls, "initial activation review must not probe target")
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "review/Cancel must not mutate any local table")
}
