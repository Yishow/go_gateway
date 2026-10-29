package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestPostgresSuccessorUnsupportedShapesHaveNoEffects(t *testing.T) {
	for _, shape := range []string{"native_monthly", "cross_mode", "cross_engine", "custom", "partial"} {
		t.Run(shape, func(t *testing.T) {
			r := newPostgresSuccessorRig(t, true, true, false, false, false, meanSource{name: "voltage", tagType: "float32", role: "gauge", required: true})
			switch shape {
			case "native_monthly":
				r.request.Group.StorageMode = "native_monthly"
			case "cross_mode":
				r.request.Group.StorageMode = "timescale"
			case "cross_engine":
				_, err := r.f.db.Exec(`UPDATE database_connectors SET kind='sqlite' WHERE id='connector-B'`)
				require.NoError(t, err)
			case "custom":
				r.request.Group.Destination.StorageStrategy = workspace.WriteGroupStorageStrategyCustom
			case "partial":
				require.Len(t, r.request.Group.Members, 2)
				r.request.Group.Members = r.request.Group.Members[:1]
			}
			before := successorAPIAllState(t, r.f.db)
			target := r.targetState(t)
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor-review", r.request)
			if response.Code == http.StatusOK {
				var body struct {
					Data workspace.WriteGroupSuccessorReview `json:"data"`
				}
				require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
				require.NotEmpty(t, body.Data.Blockers, response.Body.String())
			} else {
				require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
				require.Less(t, response.Code, 500, response.Body.String())
			}
			require.Equal(t, before, successorAPIAllState(t, r.f.db))
			require.Equal(t, target, r.targetState(t))
		})
	}
}

func TestPostgresSuccessorTimescaleMissingExtensionDoesNotDowngrade(t *testing.T) {
	r := newPostgresSuccessorRig(t, true, true, false, true, false)
	// The separate plain test database intentionally has no Timescale extension.
	// No schema or table in that database is created by this negative case.
	_, err := r.f.db.Exec(`UPDATE database_connectors SET connection_config=replace(connection_config,'gateway_timescale_test','gateway_test') WHERE id='connector-B'`)
	require.NoError(t, err)
	review := r.review(t)
	prepared := r.reserve(t, review)
	request := successorAPIStage(t, r.successorAPIFixture, review, prepared)
	request["action"] = "schema-preview"
	before := successorAPIAllState(t, r.f.db)
	target := r.targetState(t)
	response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "extension-not-installed", "must specifically refuse missing extension, not an unrelated schema error")
	require.Equal(t, before, successorAPIAllState(t, r.f.db))
	require.Equal(t, target, r.targetState(t))
	current, err := r.f.groups.Get(t.Context(), prepared.Group.ID)
	require.NoError(t, err)
	require.Equal(t, workspace.WriteGroupStorageModeTimescale, current.Group.StorageMode)
	require.Empty(t, current.Group.AppliedRevision)
	require.Empty(t, current.Group.Destination.SchemaDigest)
}
