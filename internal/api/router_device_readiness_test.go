package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Exercise the same NewRouter entry used by cmd/test_ui, never a test-only
// readiness registration. Unknown saved diagnostics must stay unknown, and
// neither device activation nor diagnostic evidence may change.
func TestNewRouter_DeviceReadinessUsesSavedDiagnosticsWithoutActivation(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	svc := device.NewService(device.NewSQLRepository(f.db), nil)
	router := NewRouter(&DatalinkServices{Device: svc, Workspace: f.workspace})
	for _, status := range []schema.DeviceStatus{schema.DeviceStatusDraft, schema.DeviceStatusDisabled} {
		t.Run(string(status), func(t *testing.T) {
			_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET description = '', status = ?, connection_config = ?, last_test_success = NULL WHERE id = 'device-A'`, status, `{"host":"127.0.0.1","port":1,"slave_id":1}`)
			require.NoError(t, err)
			before, err := svc.GetByID(t.Context(), "device-A")
			require.NoError(t, err)
			workspaceBefore, err := f.workspace.ReadSnapshot(t.Context())
			require.NoError(t, err)
			response := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/devices/device-A/readiness", nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			body := decodeJSONBody(t, response)
			require.Equal(t, true, body["success"])
			data := body["data"].(map[string]any)
			require.Equal(t, "device-A", data["device_id"])
			require.Equal(t, "unknown", data["connect_status"])
			require.Equal(t, "unknown", data["probe_status"])
			require.Equal(t, false, data["planning_allowed"])
			require.Equal(t, false, data["activation_allowed"])
			require.Equal(t, false, data["apply_allowed"])
			after, err := svc.GetByID(t.Context(), "device-A")
			require.NoError(t, err)
			require.Equal(t, before.Status, after.Status)
			require.Equal(t, before.ConnectionConfig, after.ConnectionConfig)
			require.Equal(t, before.LastTestSuccess, after.LastTestSuccess)
			require.Equal(t, before.LastTestError, after.LastTestError)
			require.NotEmpty(t, after.ReadinessStatus) // Existing metadata refresh is permitted.
			workspaceAfter, err := f.workspace.ReadSnapshot(t.Context())
			require.NoError(t, err)
			require.Equal(t, workspaceBefore, workspaceAfter)
		})
	}
	response := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/devices/unknown/readiness", nil)
	require.Equal(t, http.StatusInternalServerError, response.Code) // Existing handler's missing-ID contract.
	require.Equal(t, false, decodeJSONBody(t, response)["success"])
}
