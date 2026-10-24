package api

import (
	"maps"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// TestCircuitProfileCounterHTTPConfirmation proves through real HTTP routes that:
// 1. A cumulative profile naming a gauge or unconfirmed measurement fails readiness and blocks start.
// 2. The UI's exact PUT payload updates the definition to a confirmed counter with unit and policy.
// 3. Readback via GET /measurements confirms exact revisions, series epoch, and counter policy.
// 4. Readiness re-evaluation via GET /readiness becomes ready=true without issues.
// 5. Apply (POST /apply) and recording start succeed and activate the device.
// 6. When source binding proof is missing, readiness remains safely blocked without faking revisions.
func TestCircuitProfileCounterHTTPConfirmation(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind, "uint64")

			// 1. Initial measurement definition via real HTTP POST:
			// Configured as a gauge with empty unit, but with valid binding proof.
			initialPayload := map[string]any{
				"id":                      "meas-energy",
				"workspace_id":            f.record.ID,
				"device_id":               "device-A",
				"point_id":                "point-A",
				"tag_id":                  "tag-A",
				"equipment_id":            "main",
				"name":                    "Energy",
				"quantity":                "energy",
				"unit":                    "",
				"semantic_kind":           "gauge",
				"definition_revision":     "def-1",
				"source_binding_revision": "bind-1",
				"series_epoch":            "epoch-1",
				"numeric_encoding":        "uint64",
			}
			createRes := performJSONRequest(t, f.router, http.MethodPost, measurementsPath, initialPayload)
			require.Equal(t, http.StatusCreated, createRes.Code, createRes.Body.String())

			// 2. Ensure cumulative profile group pointing to this measurement
			cumGroup := ensureProfileGroup(t, f, "cumulative", 10, "minute", "meas-energy")
			cumID := cumGroup["id"].(string)

			// Confirm schema explicitly (Step 4 managed long table preparation)
			token := previewLong(t, f, cumID)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, cumID, token).Status)

			// 3. Readiness check: must report NOT ready due to gauge semantic
			readinessRes := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID+"/readiness", nil)
			require.Equal(t, http.StatusOK, readinessRes.Code, readinessRes.Body.String())
			readinessBody := decodeJSONBody(t, readinessRes)["data"].(map[string]any)
			require.Equal(t, false, readinessBody["ready"], "unconfirmed gauge must not be ready")

			// Start service must refuse activation while readiness is blocked
			currentWSBefore, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)

			svc := newRecordingStartServiceFixture(t, managedGroupFixture{writeGroupRouterFixture: f.writeGroupRouterFixture, groups: f.groups}, recordingStartSuccess())
			startReq := workspace.RecordingStartRequest{
				RequestID:                 "req-blocked-gauge",
				WorkspaceID:               f.record.ID,
				ExpectedWorkspaceRevision: currentWSBefore.DatabaseSetupRevision,
				DeviceIDs:                 []string{"device-A"},
				Groups: []workspace.RecordingStartGroupIntent{
					{GroupID: cumID, ExpectedGroupRevision: cumGroup["revision"].(string), ExpectedConnectorRevision: "connector-1"},
				},
			}
			opBlocked, err := svc.start.Start(t.Context(), startReq)
			require.NoError(t, err)
			require.NotEqual(t, recordingplan.SchemaOperationSucceeded, opBlocked.Status)
			require.Equal(t, "readiness", opBlocked.Stage)
			require.Equal(t, "preparation_required", opBlocked.Reason)
			require.Zero(t, svc.activation.Calls(), "device must never activate while counter is unconfirmed")
			require.Empty(t, appliedRevisionOf(t, f, cumID), "group must not be applied while counter is unconfirmed")

			// 4. Narrow confirmation: clone saved data via maps.Clone, mutate only unit, semantics, and policy
			listBefore := performJSONRequest(t, f.router, http.MethodGet, measurementsPath, nil)
			require.Equal(t, http.StatusOK, listBefore.Code, listBefore.Body.String())
			savedList := decodeJSONBody(t, listBefore)["data"].([]any)
			require.NotEmpty(t, savedList)
			confirmPayload := maps.Clone(savedList[0].(map[string]any))
			confirmPayload["semantic_kind"] = "counter"
			confirmPayload["unit"] = "kWh"
			confirmPayload["counter_policy"] = map[string]any{
				"allow_negative": false,
			}

			updateRes := performJSONRequest(t, f.router, http.MethodPut, measurementsPath+"/meas-energy", confirmPayload)
			require.Equal(t, http.StatusOK, updateRes.Code, updateRes.Body.String())

			// 5. Readback verification via real HTTP GET /measurements
			getRes := performJSONRequest(t, f.router, http.MethodGet, measurementsPath, nil)
			require.Equal(t, http.StatusOK, getRes.Code, getRes.Body.String())
			listData := decodeJSONBody(t, getRes)["data"].([]any)
			require.NotEmpty(t, listData)
			getDef := listData[0].(map[string]any)
			require.Equal(t, "counter", getDef["semantic_kind"])
			require.Equal(t, "kWh", getDef["unit"])
			require.Equal(t, "uint64", getDef["numeric_encoding"])
			require.Equal(t, "bind-1", getDef["source_binding_revision"])
			require.NotEmpty(t, getDef["series_epoch"])
			policy := getDef["counter_policy"].(map[string]any)
			require.Equal(t, false, policy["allow_negative"])

			// 6. Readiness check: now ready == true
			readinessRes2 := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID+"/readiness", nil)
			require.Equal(t, http.StatusOK, readinessRes2.Code, readinessRes2.Body.String())
			readinessBody2 := decodeJSONBody(t, readinessRes2)["data"].(map[string]any)
			require.Equal(t, true, readinessBody2["ready"], "confirmed counter must be ready: %+v", readinessBody2["issues"])

			// 7. Apply & Recording Start succeeds
			applyRes := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+cumID+"/apply", scopeBody(t, f, cumID))
			require.Equal(t, http.StatusOK, applyRes.Code, applyRes.Body.String())

			reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID, nil)
			require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
			_, savedGroup := groupSaveData(t, decodeJSONBody(t, reloaded))

			currentWS, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)

			startReqSuccess := workspace.RecordingStartRequest{
				RequestID:                 "req-confirmed-success",
				WorkspaceID:               f.record.ID,
				ExpectedWorkspaceRevision: currentWS.DatabaseSetupRevision,
				DeviceIDs:                 []string{"device-A"},
				Groups: []workspace.RecordingStartGroupIntent{
					{GroupID: cumID, ExpectedGroupRevision: savedGroup["revision"].(string), ExpectedConnectorRevision: "connector-1"},
				},
			}
			opSuccess, err := svc.start.Start(t.Context(), startReqSuccess)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, opSuccess.Status)
			require.Equal(t, 1, svc.activation.Calls())

			// 8. Negative case: missing binding proof remains blocked
			_, err = f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET source_binding_revision = '' WHERE id = 'meas-energy'`)
			require.NoError(t, err)

			readinessBlocked := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID+"/readiness", nil)
			require.Equal(t, http.StatusOK, readinessBlocked.Code, readinessBlocked.Body.String())
			blockedData := decodeJSONBody(t, readinessBlocked)["data"].(map[string]any)
			require.Equal(t, false, blockedData["ready"], "definition without binding proof must never be ready")
		})
	}
}

// TestCircuitProfileCounterHTTPExistingPolicyPreserved verifies that a counter with existing
// exact rollover modulus, max rate per second, and numeric encoding preserves those exact string
// values intact when only unit is completed. It confirms series_epoch and revisions follow contract.
func TestCircuitProfileCounterHTTPExistingPolicyPreserved(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind, "uint64")

			// Initial definition: counter with exact rollover policy and numeric encoding, missing unit
			initialPayload := map[string]any{
				"id":                      "meas-counter-exact",
				"workspace_id":            f.record.ID,
				"device_id":               "device-A",
				"point_id":                "point-A",
				"tag_id":                  "tag-A",
				"equipment_id":            "main",
				"name":                    "Exact Energy",
				"quantity":                "energy",
				"unit":                    "", // only missing unit
				"semantic_kind":           "counter",
				"definition_revision":     "def-1",
				"source_binding_revision": "bind-1",
				"series_epoch":            "epoch-1",
				"numeric_encoding":        "uint64",
				"counter_policy": map[string]any{
					"allow_negative":      false,
					"rollover_modulus":    "18446744073709551616",
					"max_rate_per_second": "0.125",
				},
			}
			createRes := performJSONRequest(t, f.router, http.MethodPost, measurementsPath, initialPayload)
			require.Equal(t, http.StatusCreated, createRes.Code, createRes.Body.String())

			// Ensure cumulative profile group pointing to this measurement
			cumGroup := ensureProfileGroup(t, f, "cumulative", 10, "minute", "meas-counter-exact")
			cumID := cumGroup["id"].(string)

			// Confirm schema
			token := previewLong(t, f, cumID)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, cumID, token).Status)

			// Readiness check before unit is confirmed: must be false due to missing unit
			readinessRes := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID+"/readiness", nil)
			require.Equal(t, http.StatusOK, readinessRes.Code, readinessRes.Body.String())
			readinessBody := decodeJSONBody(t, readinessRes)["data"].(map[string]any)
			require.Equal(t, false, readinessBody["ready"], "counter missing unit must not be ready")

			// Blocked start check with fresh workspace revision
			currentWSBefore, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)

			svc := newRecordingStartServiceFixture(t, managedGroupFixture{writeGroupRouterFixture: f.writeGroupRouterFixture, groups: f.groups}, recordingStartSuccess())
			startReq := workspace.RecordingStartRequest{
				RequestID:                 "req-blocked-unit",
				WorkspaceID:               f.record.ID,
				ExpectedWorkspaceRevision: currentWSBefore.DatabaseSetupRevision,
				DeviceIDs:                 []string{"device-A"},
				Groups: []workspace.RecordingStartGroupIntent{
					{GroupID: cumID, ExpectedGroupRevision: cumGroup["revision"].(string), ExpectedConnectorRevision: "connector-1"},
				},
			}
			opBlocked, err := svc.start.Start(t.Context(), startReq)
			require.NoError(t, err)
			require.NotEqual(t, recordingplan.SchemaOperationSucceeded, opBlocked.Status)
			require.Equal(t, "readiness", opBlocked.Stage)
			require.Equal(t, "preparation_required", opBlocked.Reason)
			require.Zero(t, svc.activation.Calls())
			require.Empty(t, appliedRevisionOf(t, f, cumID))

			// Read back actual saved definition from GET list
			listRes := performJSONRequest(t, f.router, http.MethodGet, measurementsPath, nil)
			require.Equal(t, http.StatusOK, listRes.Code, listRes.Body.String())
			savedDefs := decodeJSONBody(t, listRes)["data"].([]any)
			require.NotEmpty(t, savedDefs)
			savedDef := savedDefs[0].(map[string]any)

			// Clone saved definition via maps.Clone and update ONLY unit
			confirmPayload := maps.Clone(savedDef)
			confirmPayload["unit"] = "kWh"

			// Send real PUT with cloned payload
			updateRes := performJSONRequest(t, f.router, http.MethodPut, measurementsPath+"/meas-counter-exact", confirmPayload)
			require.Equal(t, http.StatusOK, updateRes.Code, updateRes.Body.String())

			// Readback via real GET list
			readbackRes := performJSONRequest(t, f.router, http.MethodGet, measurementsPath, nil)
			require.Equal(t, http.StatusOK, readbackRes.Code, readbackRes.Body.String())
			readbackDefs := decodeJSONBody(t, readbackRes)["data"].([]any)
			require.NotEmpty(t, readbackDefs)
			updated := readbackDefs[0].(map[string]any)

			// Assert exact policy preservation: rollover_modulus and max_rate_per_second must be exact
			policy := updated["counter_policy"].(map[string]any)
			require.Equal(t, "18446744073709551616", policy["rollover_modulus"], "exact rollover string preserved")
			require.Equal(t, "0.125", policy["max_rate_per_second"], "exact max rate string preserved")
			require.Equal(t, false, policy["allow_negative"])

			// Assert encoding and binding preserved
			require.Equal(t, "uint64", updated["numeric_encoding"])
			require.Equal(t, "bind-1", updated["source_binding_revision"])
			require.Equal(t, "kWh", updated["unit"])
			require.Equal(t, "counter", updated["semantic_kind"])

			// Assert revisions: definition_revision incremented per service contract
			require.Equal(t, "rev-2", updated["definition_revision"])
			// Initial unit completion does NOT break comparability, so series_epoch remains epoch-1
			require.Equal(t, "epoch-1", updated["series_epoch"])

			// Readiness re-evaluation: ready == true
			readinessRes2 := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID+"/readiness", nil)
			require.Equal(t, http.StatusOK, readinessRes2.Code, readinessRes2.Body.String())
			readinessBody2 := decodeJSONBody(t, readinessRes2)["data"].(map[string]any)
			require.Equal(t, true, readinessBody2["ready"], "confirmed counter must be ready: %+v", readinessBody2["issues"])

			// Apply succeeds
			applyRes := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+cumID+"/apply", scopeBody(t, f, cumID))
			require.Equal(t, http.StatusOK, applyRes.Code, applyRes.Body.String())

			// Recording Start succeeds
			reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+cumID, nil)
			require.Equal(t, http.StatusOK, reloaded.Code, reloaded.Body.String())
			_, savedGroup := groupSaveData(t, decodeJSONBody(t, reloaded))

			currentWS, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)

			startReqSuccess := workspace.RecordingStartRequest{
				RequestID:                 "req-exact-confirmed",
				WorkspaceID:               f.record.ID,
				ExpectedWorkspaceRevision: currentWS.DatabaseSetupRevision,
				DeviceIDs:                 []string{"device-A"},
				Groups: []workspace.RecordingStartGroupIntent{
					{GroupID: cumID, ExpectedGroupRevision: savedGroup["revision"].(string), ExpectedConnectorRevision: "connector-1"},
				},
			}
			opSuccess, err := svc.start.Start(t.Context(), startReqSuccess)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, opSuccess.Status)
			require.Equal(t, 1, svc.activation.Calls())
		})
	}
}
