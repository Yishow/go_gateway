package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestWriteGroupSuccessorActivationStagesRejectAmbiguousInput(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	review := f.postActivation(t, f.activationRequest(t))
	before := successorAPIAllState(t, f.f.db)
	for _, request := range []map[string]any{
		{"action": "activate", "envelope": review["envelope"], "review_digest": review["review_digest"]},
		{"action": "activation-confirm", "envelope": review["envelope"], "review_digest": review["review_digest"]},
		{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"], "allow_reserved": true},
		{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"], "readiness_digest": "precomputed"},
		{"action": "activation-review", "stage": f.activationRequest(t)["stage"], "boundary": f.boundary, "token": "schema-token"},
		{"action": "activation-confirm", "envelope": f.reviewResult.Envelope, "review_digest": f.reviewResult.ReviewDigest, "readiness_digest": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	} {
		response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		require.Equal(t, before, successorAPIAllState(t, f.f.db))
	}
}

func TestWriteGroupSuccessorActivationRechecksStoredProofAndCAS(t *testing.T) {
	for _, field := range []string{"connector", "schema", "head", "source", "target"} {
		t.Run(field, func(t *testing.T) {
			f := newSuccessorActivationAPIFixture(t)
			request := f.activationConfirmation(t)
			var err error
			switch field {
			case "connector":
				_, err = f.f.db.Exec(`UPDATE database_connectors SET identity_revision='changed' WHERE id='connector-A'`)
			case "schema":
				_, err = f.f.db.Exec(`UPDATE write_groups SET destination_schema_digest='changed' WHERE id=?`, f.prepared.Group.ID)
			case "head":
				_, err = f.f.db.Exec(`UPDATE write_group_lineages SET revision=revision+1 WHERE lineage_id=?`, f.prepared.LineageID)
			case "source":
				_, err = f.f.db.Exec(`UPDATE write_groups SET revision='changed' WHERE id=?`, f.groupID)
			case "target":
				_, err = f.f.db.Exec(`UPDATE write_groups SET revision='changed' WHERE id=?`, f.prepared.Group.ID)
			}
			require.NoError(t, err)
			before := successorAPIAllState(t, f.f.db)
			response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, f.f.db), "stale evidence must not write activation metadata")
			require.NotContains(t, response.Body.String(), f.f.file)
		})
	}
}

func TestWriteGroupSuccessorActivationCannotUseLegacyApplyOrRollback(t *testing.T) {
	f := newSuccessorActivationAPIFixture(t)
	event := f.postActivation(t, f.activationConfirmation(t))
	before := successorAPIAllState(t, f.f.db)
	for _, id := range []string{f.groupID, f.prepared.Group.ID} {
		response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f.f, id))
		require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
		require.Equal(t, before, successorAPIAllState(t, f.f.db))
	}
	current, err := f.f.groups.Get(t.Context(), f.groupID)
	require.NoError(t, err)
	request := workspace.CutoverRequest{Kind: workspace.WriteGroupCutoverRollback, WorkspaceID: f.f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, SourceGroupID: f.groupID, ExpectedSourceRevision: current.Group.Revision, ExpectedConnectorRevision: "connector-1", CutoverID: event["id"].(string), Boundary: f.boundary.Add(time.Minute)}
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/cutover-review", request)
	if response.Code == http.StatusOK {
		var body struct {
			Data struct {
				Applicable bool `json:"applicable"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		require.False(t, body.Data.Applicable, response.Body.String())
	} else {
		require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	}
	request.ReviewDigest = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	response = performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/cutover", request)
	require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, f.f.db))
}
