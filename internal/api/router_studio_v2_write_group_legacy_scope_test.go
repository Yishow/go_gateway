package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter_LegacyWriteAdapterConflictUnmigratedCRUD(t *testing.T) {
	f, targetPath := newLegacyWriteRouterFixture(t, false)
	createLegacyWriteTarget(t, targetPath)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES ('tag-unowned','unowned','unowned','Unowned','float32','active','{}')`)
	require.NoError(t, err)
	before := legacyWriteLocalSnapshot(t, f.db)
	created := performJSONRequest(t, f.router, http.MethodPost, legacyTargetMappingsPath, map[string]any{
		"tag_id": "tag-unowned", "connector_id": "connector-A", "table_schema": "main",
		"table_name": "raw_values", "column_name": "temperature", "write_mode": "insert", "enabled": true,
	})
	require.Equal(t, http.StatusCreated, created.Code, created.Body.String())
	id := decodeJSONBody(t, created)["data"].(map[string]any)["id"].(string)
	updated := performJSONRequest(t, f.router, http.MethodPut, legacyTargetMappingsPath+"/"+id, map[string]any{
		"column_name": "other_temperature", "enabled": false,
	})
	require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
	read := performJSONRequest(t, f.router, http.MethodGet, legacyTargetMappingsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	data := decodeJSONBody(t, read)["data"].(map[string]any)
	require.Equal(t, "other_temperature", data["column_name"])
	require.Equal(t, false, data["enabled"])
	require.NotContains(t, data, "canonical_group")
	deleted := performJSONRequest(t, f.router, http.MethodDelete, legacyTargetMappingsPath+"/"+id, nil)
	require.Equal(t, http.StatusOK, deleted.Code, deleted.Body.String())
	require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
}

func TestNewRouter_LegacyWriteAdapterConflictOriginalAndCurrentScopes(t *testing.T) {
	f, _ := newLegacyWriteRouterFixture(t, false)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO database_connectors
		(id,name,kind,connection_config,identity_revision,status,enabled)
		VALUES ('connector-B','Other target','sqlite','{}','connector-2','ready',1)`)
	require.NoError(t, err)
	list := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, list.Code, list.Body.String())
	group := decodeJSONBody(t, list)["data"].(map[string]any)["groups"].([]any)[0].(map[string]any)
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	data, group := groupSaveData(t, decodeJSONBody(t, read))
	group["destination"].(map[string]any)["connector_id"] = "connector-B"
	group["destination"].(map[string]any)["connector_revision"] = "connector-2"
	request := groupEditRequest(f, data, group)
	request["expected_connector_revision"] = "connector-2"
	updated := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), request)
	require.Equal(t, http.StatusOK, updated.Code, updated.Body.String())
	// The reviewed original intent remains authoritative for the adapter guard
	// even if its legacy storage row is unavailable after a canonical edit.
	_, err = f.db.ExecContext(t.Context(), `DELETE FROM database_target_mappings WHERE id='legacy-A'`)
	require.NoError(t, err)
	before := legacyWriteLocalSnapshot(t, f.db)
	for _, connectorID := range []string{"connector-A", "connector-B"} {
		response := performJSONRequest(t, f.router, http.MethodPost, legacyTargetMappingsPath, map[string]any{
			"tag_id": "tag-A", "connector_id": connectorID, "table_schema": "main",
			"table_name": "different_table", "column_name": "different_column", "write_mode": "insert",
		})
		assertLegacyWriteConflict(t, response.Code, decodeJSONBody(t, response))
		require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
	}
}

func TestNewRouter_LegacyWriteAdapterConflictNormalizedAndMissingLegacyIdentity(t *testing.T) {
	f, _ := newLegacyWriteRouterFixture(t, false)
	before := legacyWriteLocalSnapshot(t, f.db)
	create := performJSONRequest(t, f.router, http.MethodPost, legacyTargetMappingsPath, map[string]any{
		"tag_id": " tag-A ", "connector_id": " connector-A ", "table_schema": "main",
		"table_name": "different_table", "column_name": "different_column", "write_mode": "insert",
	})
	assertLegacyWriteConflict(t, create.Code, decodeJSONBody(t, create))
	require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM database_target_mappings WHERE id='legacy-A'`)
	require.NoError(t, err)
	before = legacyWriteLocalSnapshot(t, f.db)
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		response := performJSONRequest(t, f.router, method, legacyTargetMappingsPath+"/legacy-A", map[string]any{"column_name": "changed"})
		assertLegacyWriteConflict(t, response.Code, decodeJSONBody(t, response))
		require.Equal(t, before, legacyWriteLocalSnapshot(t, f.db))
	}
}
