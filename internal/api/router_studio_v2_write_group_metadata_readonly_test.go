package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManagedGroupMetadataNeverCreatesSQLiteDestination(t *testing.T) {
	for _, groupScoped := range []bool{false, true} {
		t.Run(map[bool]string{false: "connector-default", true: "saved-group"}[groupScoped], func(t *testing.T) {
			f := newManagedGroupFixture(t, "sqlite")
			query := url.Values{"expected_connector_revision": {"connector-1"}}
			if groupScoped {
				query.Set("group_id", f.group["id"].(string))
				query.Set("expected_group_revision", f.group["revision"].(string))
			}
			require.NoFileExists(t, f.file)
			response := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-metadata?"+query.Encode(), nil)
			require.Equal(t, http.StatusOK, response.Code, response.Body.String())
			require.NoFileExists(t, f.file, "read-only metadata cannot pre-create the file before explicit schema confirmation")
			body := decodeJSONBody(t, response)["data"].(map[string]any)
			require.NotEqual(t, "exists", body["inspection_status"])
			if groupScoped {
				require.Equal(t, "missing", body["inspection_status"])
			}
		})
	}
}

func TestManagedGroupMetadataRejectsInternalSQLiteDestination(t *testing.T) {
	f := newManagedGroupFixture(t, "sqlite")
	var sequence int
	var name, internalFile string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&sequence, &name, &internalFile))
	config, err := json.Marshal(map[string]string{"dsn": internalFile})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config=? WHERE id='connector-A'`, string(config))
	require.NoError(t, err)
	query := url.Values{"expected_connector_revision": {"connector-1"}, "group_id": {f.group["id"].(string)}, "expected_group_revision": {f.group["revision"].(string)}}
	response := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-metadata?"+query.Encode(), nil)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "WRITE_GROUP_INTERNAL_DATABASE")
	require.NotContains(t, response.Body.String(), internalFile)
}
