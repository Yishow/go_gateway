package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_BasicGroupReadinessIsScopedAndReadOnly(t *testing.T) {
	for _, strategy := range []string{"managed", "custom"} {
		t.Run(strategy, func(t *testing.T) {
			f := newWriteGroupRouterFixture(t)
			target := filepath.Join(t.TempDir(), "uncreated-target.db")
			config, err := json.Marshal(map[string]string{"database": target, "schema": "main"})
			require.NoError(t, err)
			_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = 'connector-A'`, string(config))
			require.NoError(t, err)
			request := f.createRequest()
			request["group"].(map[string]any)["destination"].(map[string]any)["storage_strategy"] = strategy
			created := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath, request)
			require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
			data, group := groupSaveData(t, decodeJSONBody(t, created))
			id := group["id"].(string)
			response := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id+"/readiness", nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			body := decodeJSONBody(t, response)
			require.Equal(t, true, body["success"])
			readiness := body["data"].(map[string]any)
			require.Equal(t, f.record.ID, readiness["workspace_id"])
			require.Equal(t, data["workspace_revision"], readiness["workspace_revision"])
			require.Equal(t, group["revision"], readiness["group_revision"])
			require.Equal(t, id, readiness["group_id"])
			require.Equal(t, "", readiness["applied_revision"])
			require.Equal(t, true, readiness["config_ready"])
			require.Equal(t, false, readiness["schema_ready"])
			require.Equal(t, false, readiness["ready"])
			require.NotContains(t, response.Body.String(), target)
			reloaded := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+id, nil)
			require.Equal(t, decodeJSONBody(t, created), decodeJSONBody(t, reloaded))
			require.NoFileExists(t, target)
			unknown := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/unknown/readiness", nil)
			require.Equal(t, http.StatusNotFound, unknown.Code)
			assertWriteGroupError(t, decodeJSONBody(t, unknown), "WRITE_GROUP_NOT_FOUND")
		})
	}
}
