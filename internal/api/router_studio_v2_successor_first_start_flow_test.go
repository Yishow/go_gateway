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

func TestSuccessorFirstStartHTTPHonestResultAndImmutableRecovery(t *testing.T) {
	for _, sourceLong := range []bool{false, true} {
		name := "wide-to-long"
		if sourceLong {
			name = "long-to-wide"
		}
		t.Run(name, func(t *testing.T) {
			r := newFirstStartAPIRig(t, "sqlite", sourceLong, !sourceLong, false, workspace.WriteGroupStorageModePlain)
			old, err := r.f.groups.Get(t.Context(), r.groupID)
			require.NoError(t, err)
			before := successorAPIAllState(t, r.f.db)
			targetBefore := successorAPIAllState(t, r.target)
			request := r.confirmation(t)
			require.Equal(t, before, successorAPIAllState(t, r.f.db), "review and readiness are read-only")
			result := r.activationAPI().postActivation(t, request)
			require.Equal(t, "succeeded", result["status"])
			require.Equal(t, r.request.OperationID, result["root_operation_id"])
			require.Equal(t, recordingplan.SuccessorFirstStartOperationID(r.f.record.ID, recordingplan.SuccessorReserveOperationID(r.f.record.ID, r.request.OperationID)), result["operation_id"])
			require.Equal(t, r.prepared.Group.ID, result["group_id"])
			require.Equal(t, r.groupID, result["predecessor_group_id"])
			require.EqualValues(t, 2, result["generation"])
			require.NotEmpty(t, result["applied_revision"])
			for _, field := range []string{"ended", "moved", "ownership_fence_revision", "id", "sql_committed", "runtime_running"} {
				require.NotContains(t, result, field)
			}
			snapshot, err := r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary)
			require.NoError(t, err)
			require.True(t, r.boundary.Equal(snapshot.EffectiveAt))
			_, err = r.f.groups.ResolveAppliedAt(t.Context(), r.prepared.Group.ID, r.boundary.Add(-time.Nanosecond))
			require.ErrorIs(t, err, workspace.ErrWriteGroupLifecycleBlocked)
			preserved, err := r.f.groups.Get(t.Context(), r.groupID)
			require.NoError(t, err)
			require.Equal(t, old.Group, preserved.Group, "first-start never fabricates or modifies G1 runtime history")
			require.Equal(t, targetBefore, successorAPIAllState(t, r.target), "first-start does not perform SQL writes or DDL")
			status := r.activationAPI().postActivation(t, r.statusRequest())
			require.Equal(t, "not_submitted", status["cutover"].(map[string]any)["status"])
			require.Equal(t, result, status["first_start"].(map[string]any)["result"])
			var original map[string]any
			raw, err := json.Marshal(r.prepared)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &original))
			require.Equal(t, original, status["reserve"].(map[string]any)["result"], "schema/start must not overwrite the immutable reserve result")
			after := successorAPIAllState(t, r.f.db)
			require.Equal(t, result, r.activationAPI().postActivation(t, request))
			require.Equal(t, status, r.activationAPI().postActivation(t, r.statusRequest()))
			require.Equal(t, after, successorAPIAllState(t, r.f.db))
			require.Equal(t, r.prepared, r.reserve(t, r.preparationReview), "reserve exact replay survives cleared reservation/head move")
			require.Equal(t, after, successorAPIAllState(t, r.f.db))
		})
	}
}

func TestSuccessorFirstStartHTTPRejectsCutoverExchangeAndTampering(t *testing.T) {
	r := newFirstStartAPIRig(t, "sqlite", false, true, false, workspace.WriteGroupStorageModePlain)
	original := r.confirmation(t)
	before := successorAPIAllState(t, r.f.db)
	for _, field := range []string{"action", "activation_kind", "legacy_kind", "review_digest", "readiness_digest", "boundary", "group_revision"} {
		t.Run(field, func(t *testing.T) {
			raw, err := json.Marshal(original)
			require.NoError(t, err)
			var request map[string]any
			require.NoError(t, json.Unmarshal(raw, &request))
			envelope := request["envelope"].(map[string]any)
			switch field {
			case "action":
				request["action"] = "activation-confirm"
			case "activation_kind":
				envelope[field] = "cutover"
			case "legacy_kind":
				delete(envelope, "activation_kind")
				envelope["version"] = "successor-activation-review-v1"
			case "review_digest", "readiness_digest":
				request[field] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			case "boundary":
				envelope["request"].(map[string]any)["boundary"] = r.boundary.Add(time.Second)
			case "group_revision":
				envelope["request"].(map[string]any)["stage"].(map[string]any)["current"].(map[string]any)[field] = "foreign"
			}
			response := performJSONRequest(t, r.f.router, http.MethodPost, writeGroupsPath+"/"+r.groupID+"/successor", request)
			require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
			require.Less(t, response.Code, 500, response.Body.String())
			require.Equal(t, before, successorAPIAllState(t, r.f.db))
		})
	}
	cutover := newSuccessorActivationAPIFixture(t)
	request := cutover.activationConfirmation(t)
	request["action"] = "first-start-confirm"
	before = successorAPIAllState(t, cutover.f.db)
	response := performJSONRequest(t, cutover.f.router, http.MethodPost, writeGroupsPath+"/"+cutover.groupID+"/successor", request)
	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	require.Equal(t, before, successorAPIAllState(t, cutover.f.db))
}
