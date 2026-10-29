package api

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

type successorActivationAPIFixture struct {
	successorAPIFixture
	reviewResult workspace.WriteGroupSuccessorReview
	prepared     workspace.WriteGroupSuccessorPreparation
	boundary     time.Time
}

func newSuccessorActivationAPIFixture(t *testing.T) successorActivationAPIFixture {
	t.Helper()
	f := newSuccessorAPIFixture(t)
	f.request.Group.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	f.request.Group.StorageMode = workspace.WriteGroupStorageModePlain
	review := f.review(t)
	prepared := f.reserve(t, review)
	token := successorAPIPreview(t, f, review, prepared)
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", successorAPIConfirmRequest(t, f, review, prepared, token))
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var body struct {
		Data recordingplan.SchemaOperation `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
	return successorActivationAPIFixture{successorAPIFixture: f, reviewResult: review, prepared: prepared, boundary: f.request.Boundary.Add(time.Minute)}
}

func (f successorActivationAPIFixture) activationRequest(t *testing.T) map[string]any {
	t.Helper()
	return map[string]any{"action": "activation-review", "stage": successorAPIStage(t, f.successorAPIFixture, f.reviewResult, f.prepared), "boundary": f.boundary}
}
func (f successorActivationAPIFixture) postActivation(t *testing.T, request any) map[string]any {
	t.Helper()
	response := performJSONRequest(t, f.f.router, http.MethodPost, writeGroupsPath+"/"+f.groupID+"/successor", request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "%s", response.Body.String())
	return data
}
