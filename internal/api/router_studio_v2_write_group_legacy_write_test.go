package api

import (
	"net/http"
	"os"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/stretchr/testify/require"
)

const legacyTargetMappingsPath = "/api/v1/datalink/db-targets/mappings"

func TestNewRouter_LegacyWriteAdapterConflictDelete(t *testing.T) {
	f, targetPath := newSingleMigrationRouterFixture(t)
	review := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, migrationReviewRequest(t, f))
	require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	before, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	repository := dbtarget.NewSQLTargetMappingRepository(f.db)
	legacy, err := repository.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	groupsBefore := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, groupsBefore.Code, groupsBefore.Body.String())

	response := performJSONRequest(t, f.router, http.MethodDelete, legacyTargetMappingsPath+"/legacy-A", nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	assertWriteGroupError(t, body, "WRITE_GROUP_LEGACY_WRITE_CONFLICT")
	require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
	current, err := repository.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	require.Equal(t, legacy, current)
	after, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
	groupsAfter := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, decodeJSONBody(t, groupsBefore), decodeJSONBody(t, groupsAfter))
	var maps int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_migration_maps`).Scan(&maps))
	require.Equal(t, 1, maps)
	_, err = os.Stat(targetPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestNewRouter_LegacyWriteAdapterConflictTargetEntries(t *testing.T) {
	for _, test := range []struct {
		name, method, path string
		rowGroup           bool
		body               map[string]any
	}{
		{"single-update", http.MethodPut, legacyTargetMappingsPath + "/legacy-A", false, map[string]any{"column_name": "other_temperature"}},
		{"single-create", http.MethodPost, legacyTargetMappingsPath, false, map[string]any{
			"tag_id": "tag-A", "connector_id": "connector-A", "table_schema": "main", "table_name": "other_table", "column_name": "other_column", "write_mode": "insert",
		}},
		{"row-member-update", http.MethodPut, legacyTargetMappingsPath + "/legacy-B", true, map[string]any{"column_name": "other_value"}},
		{"row-member-delete", http.MethodDelete, legacyTargetMappingsPath + "/legacy-B", true, nil},
		{"studio-target-update", http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/point-A", false, map[string]any{"column_name": "other_temperature", "enabled": true}},
		{"studio-row-target-update", http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/point-A", true, map[string]any{"column_name": "other_value", "row_group_id": "legacy-g", "enabled": true}},
	} {
		t.Run(test.name, func(t *testing.T) {
			f, targetPath := newLegacyWriteRouterFixture(t, test.rowGroup)
			if test.body != nil {
				test.body["expected_setup_revision"] = f.record.DatabaseSetupRevision
			}
			before := legacyWriteLocalSnapshot(t, f.db)
			response := performJSONRequest(t, f.router, test.method, test.path, test.body)
			assertLegacyWriteConflict(t, response.Code, decodeJSONBody(t, response))
			require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
			require.NotContains(t, response.Body.String(), targetPath)
			_, err := os.Stat(targetPath)
			require.ErrorIs(t, err, os.ErrNotExist)
		})
	}
}

func TestNewRouter_LegacyWriteAdapterConflictConnectorDelete(t *testing.T) {
	f, _ := newSingleMigrationRouterFixture(t)
	review := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, migrationReviewRequest(t, f))
	require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	repository := dbtarget.NewSQLTargetMappingRepository(f.db)
	legacy, err := repository.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)

	response := performJSONRequest(t, f.router, http.MethodDelete, "/api/v1/datalink/db-targets/connectors/"+legacy.ConnectorID, nil)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_LEGACY_WRITE_CONFLICT")
	current, err := repository.GetByID(t.Context(), "legacy-A")
	require.NoError(t, err)
	require.Equal(t, legacy, current, "the legacy mapping and its connector are untouched")
	var connectors int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM database_connectors WHERE id = ?`, legacy.ConnectorID).Scan(&connectors))
	require.Equal(t, 1, connectors)
}
