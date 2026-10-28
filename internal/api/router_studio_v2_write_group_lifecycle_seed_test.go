package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestWriteGroupLifecycleSuccessorSeedUsesAppliedDefinitionAndFreshReview(t *testing.T) {
	f := newSuccessorAPIFixture(t)
	inspector := &successorAPIForbiddenInspector{}
	f.f.groups.WithTableInspector(inspector)
	f.f.groups.WithManagedTableInspector(inspector)
	// Fixture isolation: preserve ready while changing only saved fields so this
	// test examines immutable basis selection, not the real Update -> draft gate.
	_, err := f.f.db.Exec(`UPDATE write_groups SET name='edited draft',row_policy=json_set(row_policy,'$.interval_seconds',77),revision='draft-R2' WHERE id=?`, f.groupID)
	require.NoError(t, err)
	before := successorAPIAllState(t, f.f.db)
	response := performJSONRequest(t, f.f.router, http.MethodGet, writeGroupsPath+"/"+f.groupID+"/lifecycle", nil)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data workspace.WriteGroupLifecycle `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	view := body.Data
	require.True(t, view.Actions[4].Allowed)
	require.True(t, view.Actions[4].RequiresReview)
	require.False(t, view.Actions[3].Allowed)
	require.Equal(t, "draft-R2", view.GroupRevision)
	seed := view.SuccessorSeed
	require.NotNil(t, seed)
	require.Equal(t, "applied", seed.BasisKind)
	require.Equal(t, f.snapshot.GroupRevision, seed.BasisRevision)
	require.Equal(t, f.snapshot.Group.RowPolicy, seed.Group.RowPolicy)
	require.Empty(t, seed.Group.Destination.SchemaRevision)
	require.Empty(t, seed.Group.Destination.SchemaDigest)
	require.Empty(t, seed.Group.ID)
	require.Empty(t, seed.Group.Destination.TableName)
	require.Zero(t, inspector.calls)
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
	for _, mode := range []string{"valid", "foreign", "stale", "tampered"} {
		t.Run(mode, func(t *testing.T) {
			request := f.request
			raw, e := json.Marshal(seed.Group)
			require.NoError(t, e)
			request.Group = &workspace.WriteGroup{}
			require.NoError(t, json.Unmarshal(raw, request.Group))
			request.ExpectedPredecessorRevision = view.GroupRevision
			switch mode {
			case "foreign":
				request.Group.WorkspaceID = "foreign"
			case "stale":
				request.ExpectedPredecessorRevision = "stale"
			case "tampered":
				request.Group.Members[0].SourceRevision = "tampered"
			}
			review := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor-review", request)
			if mode == "valid" {
				require.Equal(t, http.StatusOK, review.Code, review.Body.String())
				require.Empty(t, decodeJSONBody(t, review)["data"].(map[string]any)["blockers"])
			} else if review.Code == http.StatusOK {
				require.NotEmpty(t, decodeJSONBody(t, review)["data"].(map[string]any)["blockers"])
			} else {
				require.GreaterOrEqual(t, review.Code, 400)
			}
			require.Equal(t, before, successorAPIAllState(t, f.f.db))
			require.Zero(t, inspector.calls)
		})
	}
}
