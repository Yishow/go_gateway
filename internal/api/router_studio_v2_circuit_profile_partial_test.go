package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestCircuitProfilesPartialPolicyUsesTheLongLayout(t *testing.T) {
	for _, backend := range []string{"sqlite", "postgres"} {
		t.Run(backend, func(t *testing.T) {
			f := newManagedDestinationFixture(t, backend, "uint64")
			insertCounterMeasurement(t, f, "meas-energy", "counter", "kWh")
			pair := profilePair{f: f}
			for _, kind := range []string{"instantaneous", "cumulative"} {
				interval, unit := 15, "second"
				if kind == "cumulative" {
					interval, unit = 1, "minute"
				}
				request := profileRequest(t, f.writeGroupRouterFixture, kind, interval, unit, "main", "")
				group := request["group"].(map[string]any)
				group["layout_version"] = "managed_long_v1"
				group["destination"].(map[string]any)["table_schema"] = f.namespace
				group["row_policy"].(map[string]any)["incomplete_policy"] = "partial"
				if kind == "cumulative" {
					group["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
				}
				response := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
				require.Equal(t, http.StatusOK, response.Code, response.Body.String())
				_, saved := groupSaveData(t, decodeJSONBody(t, response))
				id := saved["id"].(string)
				require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
				ready, err := f.groups.Readiness(t.Context(), id)
				require.NoError(t, err)
				require.True(t, ready.Ready, "%+v", ready.Issues)
				if kind == "instantaneous" {
					pair.instant = saved
				} else {
					pair.cumulative = saved
				}
			}
			svc := newRecordingStartServiceFixture(t, f, recordingStartSuccess())
			op, err := svc.start.Start(t.Context(), pair.request(t, "partial-long-start"))
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status, "%+v", op)
			require.Equal(t, 1, svc.activation.Calls())
		})
	}
}

func TestManagedWidePolicyEditRequiresNoChangeConfirmation(t *testing.T) {
	for _, backend := range []string{"sqlite", "postgres"} {
		t.Run(backend, func(t *testing.T) {
			f := newManagedGroupFixture(t, backend)
			id := f.group["id"].(string)
			first := previewLong(t, f, id)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, first).Status)
			current, err := f.groups.Get(t.Context(), id)
			require.NoError(t, err)
			current.Group.Members[0].Required = false
			current.Group.RowPolicy.IncompletePolicy = "partial"
			request := scopeBody(t, f, id)
			request["group"] = current.Group
			updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+id, request)
			require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
			ready, err := f.groups.Readiness(t.Context(), id)
			require.NoError(t, err)
			require.False(t, ready.Ready, "a policy-only edit still needs its own schema confirmation")
			preview := previewLong(t, f, id)
			require.Empty(t, preview.Statements)
			require.NotEqual(t, first.Token, preview.Token)
			require.NotEqual(t, first.OperationID, preview.OperationID)
			blocked := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/apply", scopeBody(t, f, id))
			require.Equal(t, http.StatusConflict, blocked.Code, blocked.Body.String())
			op := confirmLong(t, f, id, preview)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, op.Status)
			require.Zero(t, op.ExecutedStatements)
			ready, err = f.groups.Readiness(t.Context(), id)
			require.NoError(t, err)
			require.True(t, ready.Ready, "%+v", ready.Issues)
		})
	}
}
