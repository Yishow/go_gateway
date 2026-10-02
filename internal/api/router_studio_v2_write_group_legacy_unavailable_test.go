package api

import (
	"net/http"
	"os"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_LegacyWriteAdapterConflictOwnershipReadUnavailable(t *testing.T) {
	for _, entry := range []string{"global-post", "global-put", "global-delete", "studio-target", "studio-config"} {
		t.Run(entry, func(t *testing.T) {
			f, targetPath := newLegacyWriteRouterFixture(t, entry == "studio-config")
			repository := dbtarget.NewSQLTargetMappingRepository(f.db)
			legacy, err := repository.GetByID(t.Context(), "legacy-A")
			require.NoError(t, err)
			before, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			groupsBefore := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
			require.Equal(t, http.StatusOK, groupsBefore.Code, groupsBefore.Body.String())
			_, err = f.db.ExecContext(t.Context(), `DROP TABLE write_group_migration_maps`)
			require.NoError(t, err)
			path := legacyTargetMappingsPath + "/legacy-A"
			body := map[string]any{"column_name": "changed"}
			method := http.MethodPut
			switch entry {
			case "global-post":
				method = http.MethodPost
				path = legacyTargetMappingsPath
				body = map[string]any{"tag_id": "tag-A", "connector_id": "connector-A", "table_name": "other_table", "column_name": "other_column", "write_mode": "insert"}
			case "global-delete":
				method = http.MethodDelete
			case "studio-target":
				path = "/api/v1/datalink/studio-v2/workspace/database-targets/point-A"
				body["enabled"] = true
				body["expected_setup_revision"] = f.record.DatabaseSetupRevision
			case "studio-config":
				path = legacyDatabaseConfigPath
				body = legacyDatabaseConfigRequest(f, targetPath)
				body["row_groups"] = []any{}
			}
			response := performJSONRequest(t, f.router, method, path, body)
			require.Equal(t, http.StatusInternalServerError, response.Code, response.Body.String())
			decoded := decodeJSONBody(t, response)
			assertWriteGroupError(t, decoded, "internal")
			require.Equal(t, true, decoded["error"].(map[string]any)["retryable"])
			require.Equal(t, "Review the request and retry", decoded["error"].(map[string]any)["action"])
			require.NotContains(t, response.Body.String(), "write_group_migration_maps")
			require.NotContains(t, response.Body.String(), targetPath)
			after, err := f.workspace.GetOrCreate(t.Context())
			require.NoError(t, err)
			require.Equal(t, before, after)
			current, err := repository.GetByID(t.Context(), "legacy-A")
			require.NoError(t, err)
			require.Equal(t, legacy, current)
			groupsAfter := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
			require.Equal(t, http.StatusOK, groupsAfter.Code, groupsAfter.Body.String())
			require.Equal(t, decodeJSONBody(t, groupsBefore), decodeJSONBody(t, groupsAfter))
			_, err = os.Stat(targetPath)
			require.ErrorIs(t, err, os.ErrNotExist)
		})
	}
}
