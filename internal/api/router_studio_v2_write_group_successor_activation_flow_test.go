package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func (f successorActivationAPIFixture) activationConfirmation(t *testing.T) map[string]any {
	t.Helper()
	review := f.postActivation(t, f.activationRequest(t))
	require.Empty(t, review["blockers"])
	ready := f.postActivation(t, map[string]any{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"]})
	readiness, ok := ready["readiness"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, true, readiness["ready"], "%+v", readiness)
	require.NotEmpty(t, ready["readiness_digest"])
	return map[string]any{"action": "activation-confirm", "envelope": review["envelope"], "review_digest": review["review_digest"], "readiness_digest": ready["readiness_digest"]}
}

func TestWriteGroupSuccessorActivationUsesExplicitFreshReadiness(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	before := successorAPIAllState(t, f.f.db)
	target := f.target(t)
	beforeTarget := successorAPIAllState(t, target)
	request := f.activationConfirmation(t)
	require.Equal(t, before, successorAPIAllState(t, f.f.db), "review/readiness must not change reservation, versions or operations")
	require.Equal(t, beforeTarget, successorAPIAllState(t, target), "readiness must not execute DDL")
	previous, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	result := f.postActivation(t, request)
	require.NotEmpty(t, result["id"])
	var head string
	var pending any
	require.NoError(t, f.f.db.QueryRow(`SELECT current_group_id,reserved_group_id FROM write_group_lineages WHERE lineage_id=?`, f.prepared.LineageID).Scan(&head, &pending))
	require.Equal(t, f.prepared.Group.ID, head)
	require.Nil(t, pending)
	old, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	require.Equal(t, previous.Group.Status, old.Group.Status, "scheduling must not stop old pre-boundary intake")
	require.Equal(t, previous.Group.AppliedRevision, old.Group.AppliedRevision)
	next, err := f.f.groups.Get(t.Context(), f.prepared.Group.ID)
	require.NoError(t, err)
	require.NotEmpty(t, next.Group.AppliedRevision)
	var beforeBoundary, atBoundary *workspace.WriteGroupAppliedSnapshot
	beforeBoundary, err = f.f.groups.ResolveAppliedAt(t.Context(), next.Group.ID, f.boundary.Add(-time.Nanosecond))
	require.ErrorIs(t, err, workspace.ErrWriteGroupLifecycleBlocked)
	require.Nil(t, beforeBoundary)
	atBoundary, err = f.f.groups.ResolveAppliedAt(t.Context(), next.Group.ID, f.boundary)
	require.NoError(t, err)
	require.NotNil(t, atBoundary)
	require.True(t, f.boundary.Equal(atBoundary.EffectiveAt))
	after := successorAPIAllState(t, f.f.db)
	replay := f.postActivation(t, request)
	require.Equal(t, result["id"], replay["id"])
	require.Equal(t, after, successorAPIAllState(t, f.f.db), "full intent retry after reservation clear must return same event without effects")
	require.Equal(t, beforeTarget, successorAPIAllState(t, target), "activation must not mutate SQL history or perform DDL")
}

func TestWriteGroupSuccessorActivationRejectsChangedConfirmation(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	original := f.activationConfirmation(t)
	before := successorAPIAllState(t, f.f.db)
	for _, field := range []string{"review_digest", "readiness_digest", "boundary", "group_revision", "schema_digest", "source_applied_revision"} {
		t.Run(field, func(t *testing.T) {
			raw, err := json.Marshal(original)
			require.NoError(t, err)
			var request map[string]any
			require.NoError(t, json.Unmarshal(raw, &request))
			envelope := request["envelope"].(map[string]any)
			switch field {
			case "review_digest", "readiness_digest":
				request[field] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			case "boundary":
				envelope["request"].(map[string]any)["boundary"] = f.boundary.Add(time.Second)
			case "group_revision":
				envelope["request"].(map[string]any)["stage"].(map[string]any)["current"].(map[string]any)[field] = "foreign"
			default:
				envelope[field] = "foreign"
			}
			response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Less(t, response.Code, 500, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, f.f.db))
		})
	}
}
