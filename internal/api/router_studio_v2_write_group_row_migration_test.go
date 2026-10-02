package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const rowGroupMigrationPreviewPath = writeGroupsPath + "/migrations/row-groups/preview"
const rowGroupMigrationReviewPath = writeGroupsPath + "/migrations/row-groups/review"

func newRowGroupMigrationRouterFixture(t *testing.T) (fixture writeGroupRouterFixture, targetPath string) {
	t.Helper()
	f, targetPath := newSingleMigrationRouterFixture(t)
	_, err := f.db.ExecContext(t.Context(), `INSERT INTO points (id,device_id,name,address,function,data_type,mode,enabled)
		VALUES ('point-B','device-A','Pressure','40002','FC03','float32','read',1)`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO tags (id,key,key_lower,display_name,data_type,status,labels)
		VALUES ('tag-B','pressure','pressure','Pressure','float32','active','{}')`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO mappings (id,point_id,tag_id,transform_pipeline,status,enabled)
		VALUES ('mapping-B','point-B','tag-B','[]','active',1)`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_target_mappings SET column_name='value',group_key='legacy-g:point-A'
		WHERE id='legacy-A'`)
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,group_key,write_interval_seconds,enabled)
		VALUES ('legacy-B','tag-B','connector-A','main','raw_values','value','insert','observed_at','legacy-g:point-B',15,1)`)
	require.NoError(t, err)
	f.record, err = f.workspace.UpdateDatabaseSetup(t.Context(), f.record.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *workspace.Record) error {
		record.DatabaseRowGroups = []workspace.DatabaseRowGroup{{
			ID: "legacy-g", ConnectorID: "connector-A", TableSchema: "main", TableName: "raw_values",
			MemberPointIDs: []string{"point-A", "point-B"}, GroupKeyColumns: []string{"entity_id"},
			UniqueKeyColumns: []string{"entity_id", "observed_at"},
		}}
		record.DatabaseTargetRefs = []workspace.DatabaseTargetRef{{PointID: "point-A", RowGroupID: "legacy-g"}, {PointID: "point-B", RowGroupID: "legacy-g"}}
		return nil
	})
	require.NoError(t, err)
	return f, targetPath
}

func rowGroupMigrationRequest(t *testing.T, f writeGroupRouterFixture) map[string]any {
	t.Helper()
	preview := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationPreviewPath, map[string]any{
		"workspace_id": f.record.ID, "source_ids": []string{"legacy-g"},
	})
	require.Equal(t, http.StatusOK, preview.Code, preview.Body.String())
	data := decodeJSONBody(t, preview)["data"].(map[string]any)
	items := data["items"].([]any)
	require.Len(t, items, 1)
	require.Equal(t, "needs_review", items[0].(map[string]any)["status"])
	return map[string]any{
		"workspace_id": data["workspace_id"], "expected_workspace_revision": data["workspace_revision"],
		"expected_connector_revision": data["connector_revision"], "review_digest": data["review_digest"],
		"source_ids": []string{"legacy-g"}, "confirm_snapshot_conversion": true,
	}
}

func TestNewRouter_LegacyRowGroupMigrationPreviewAndReview(t *testing.T) {
	f, targetPath := newRowGroupMigrationRouterFixture(t)
	request := rowGroupMigrationRequest(t, f)
	request["apply"] = true
	review := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	data := decodeJSONBody(t, review)["data"].(map[string]any)
	groups := data["groups"].([]any)
	require.Len(t, groups, 1)
	group := groups[0].(map[string]any)
	require.NotEmpty(t, group["id"])
	require.Equal(t, "draft", group["status"])
	require.Equal(t, "", group["applied_revision"])
	members := group["members"].([]any)
	require.Len(t, members, 2)
	require.Equal(t, "legacy-g:point-A", members[0].(map[string]any)["entity_key"])
	require.Equal(t, "legacy-g:point-B", members[1].(map[string]any)["entity_key"])
	policy := group["row_policy"].(map[string]any)
	require.Equal(t, []any{"entity_id"}, policy["group_key_columns"])
	require.Equal(t, []any{"entity_id", "observed_at"}, policy["unique_key_columns"])
	migration := group["migration"].(map[string]any)
	require.Equal(t, "legacy-g", migration["legacy_row_group_id"])
	require.Equal(t, []any{"legacy-A", "legacy-B"}, migration["source_ids"])
	require.Equal(t, map[string]any{"legacy-A": "point-A", "legacy-B": "point-B"}, migration["target_mapping_points"])
	reloaded, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, f.record.DatabaseRowGroups, reloaded.DatabaseRowGroups)
	require.Equal(t, f.record.DatabaseTargetRefs, reloaded.DatabaseTargetRefs)
	var legacyEnabled int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM database_target_mappings WHERE enabled=1 AND column_name='value'`).Scan(&legacyEnabled))
	require.Equal(t, 2, legacyEnabled)
	_, err = os.Stat(targetPath)
	require.ErrorIs(t, err, os.ErrNotExist)
	members[0].(map[string]any)["target_column"] = "reviewed_value"
	edit := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	legacy := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings/legacy-A", nil)
	require.Equal(t, http.StatusOK, legacy.Code, legacy.Body.String())
	legacyData := decodeJSONBody(t, legacy)["data"].(map[string]any)
	require.Equal(t, "reviewed_value", legacyData["column_name"])
	require.Equal(t, "legacy-g:point-A", legacyData["group_key"])
	require.Equal(t, "value", legacyData["legacy_intent"].(map[string]any)["column_name"])
	read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	rows := decodeJSONBody(t, read)["data"].([]any)
	require.Len(t, rows, 2)
	for _, row := range rows {
		require.Equal(t, "legacy-g", row.(map[string]any)["row_group_id"])
	}
	replay := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, rowGroupMigrationRequest(t, f))
	require.Equal(t, http.StatusOK, replay.Code, replay.Body.String())
	replayed := decodeJSONBody(t, replay)["data"].(map[string]any)["groups"].([]any)[0].(map[string]any)
	require.Equal(t, group["id"], replayed["id"])
	require.Equal(t, "reviewed_value", replayed["members"].([]any)[0].(map[string]any)["target_column"])
}

func TestNewRouter_LegacyRowGroupMigrationReviewBoundaries(t *testing.T) {
	for _, field := range []string{"workspace_id", "expected_workspace_revision", "expected_connector_revision", "review_digest", "source_ids", "confirm_snapshot_conversion"} {
		t.Run("missing-"+field, func(t *testing.T) {
			f, _ := newRowGroupMigrationRouterFixture(t)
			request := rowGroupMigrationRequest(t, f)
			delete(request, field)
			response := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
		})
	}
	f, _ := newRowGroupMigrationRouterFixture(t)
	request := rowGroupMigrationRequest(t, f)
	request["confirm_snapshot_conversion"] = false
	response := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	request["confirm_snapshot_conversion"] = true
	request["review_digest"] = "stale-preview"
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "revision_mismatch")
	request = rowGroupMigrationRequest(t, f)
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	first := decodeJSONBody(t, response)["data"]
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	request = rowGroupMigrationRequest(t, f)
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, first, decodeJSONBody(t, response)["data"])
}

func TestNewRouter_LegacyRowGroupMigrationBlockedAndSafeFailure(t *testing.T) {
	f, _ := newRowGroupMigrationRouterFixture(t)
	for _, sourceID := range []string{"missing-group", "foreign-group"} {
		unknown := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationPreviewPath, map[string]any{
			"workspace_id": f.record.ID, "source_ids": []string{sourceID},
		})
		require.Equal(t, http.StatusNotFound, unknown.Code, unknown.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, unknown), "WRITE_GROUP_NOT_FOUND")
	}
	request := rowGroupMigrationRequest(t, f)
	_, err := f.db.ExecContext(t.Context(), `UPDATE database_target_mappings SET group_key='legacy-g:point-A' WHERE id='legacy-B'`)
	require.NoError(t, err)
	response := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
	preview := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationPreviewPath, map[string]any{
		"workspace_id": f.record.ID, "source_ids": []string{"legacy-g"},
	})
	require.Equal(t, http.StatusOK, preview.Code, preview.Body.String())
	data := decodeJSONBody(t, preview)["data"].(map[string]any)
	item := data["items"].([]any)[0].(map[string]any)
	require.Equal(t, "blocked", item["status"])
	require.NotContains(t, item, "candidate_group")
	request["review_digest"] = data["review_digest"]
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	request["workspace_id"] = "foreign"
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	f, _ = newRowGroupMigrationRouterFixture(t)
	request = rowGroupMigrationRequest(t, f)
	_, err = f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_row_migration BEFORE UPDATE ON system_settings
		WHEN NEW.key='studio_v2_workspace' BEGIN SELECT RAISE(ABORT,'private-row-migration-diagnostic'); END`)
	require.NoError(t, err)
	response = performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusInternalServerError, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), "private-row-migration-diagnostic")
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	data = decodeJSONBody(t, read)["data"].(map[string]any)
	require.Empty(t, data["groups"])
	require.Equal(t, request["expected_workspace_revision"], data["workspace_revision"])
	var maps int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT count(*) FROM write_group_migration_maps`).Scan(&maps))
	require.Zero(t, maps)
}

func TestNewRouter_LegacyRowGroupMigrationReadCannotHideRemovedMember(t *testing.T) {
	f, _ := newRowGroupMigrationRouterFixture(t)
	request := rowGroupMigrationRequest(t, f)
	review := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	data := decodeJSONBody(t, review)["data"].(map[string]any)
	group := data["groups"].([]any)[0].(map[string]any)
	group["members"] = group["members"].([]any)[:1]
	edit := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	for _, path := range []string{"/api/v1/datalink/db-targets/mappings/legacy-B", "/api/v1/datalink/studio-v2/workspace/database-targets"} {
		read := performJSONRequest(t, f.router, http.MethodGet, path, nil)
		require.Equal(t, http.StatusConflict, read.Code, read.Body.String())
		body := decodeJSONBody(t, read)
		assertWriteGroupError(t, body, "WRITE_GROUP_LEGACY_READ_CONFLICT")
		require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
	}
}

func TestNewRouter_LegacyRowGroupMigrationReadCannotHideMissingLegacyTarget(t *testing.T) {
	f, _ := newRowGroupMigrationRouterFixture(t)
	review := performJSONRequest(t, f.router, http.MethodPost, rowGroupMigrationReviewPath, rowGroupMigrationRequest(t, f))
	require.Equal(t, http.StatusOK, review.Code, review.Body.String())
	_, err := f.db.ExecContext(t.Context(), `DELETE FROM database_target_mappings WHERE id='legacy-B'`)
	require.NoError(t, err)
	for _, path := range []string{"/api/v1/datalink/db-targets/mappings", "/api/v1/datalink/studio-v2/workspace/database-targets"} {
		t.Run(path, func(t *testing.T) {
			read := performJSONRequest(t, f.router, http.MethodGet, path, nil)
			require.Equal(t, http.StatusConflict, read.Code, read.Body.String())
			body := decodeJSONBody(t, read)
			assertWriteGroupError(t, body, "WRITE_GROUP_LEGACY_READ_CONFLICT")
			require.Equal(t, "open_write_groups", body["error"].(map[string]any)["action"])
		})
	}
	read := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings?connector_id=unrelated", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	require.Empty(t, decodeJSONBody(t, read)["data"])
	read = performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings?tag_id=tag-A", nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	require.Len(t, decodeJSONBody(t, read)["data"], 1)
	read = performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings?tag_id=tag-B", nil)
	require.Equal(t, http.StatusConflict, read.Code, read.Body.String())
}
