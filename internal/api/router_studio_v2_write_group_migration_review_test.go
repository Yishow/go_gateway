package api

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const singleMappingMigrationReviewPath = writeGroupsPath + "/migrations/single-mappings/review"

func newSingleMigrationRouterFixture(t *testing.T) (fixture writeGroupRouterFixture, targetPath string) {
	t.Helper()
	f := newWriteGroupRouterFixture(t)
	targetPath = filepath.Join(t.TempDir(), "unopened-target.db")
	_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = 'connector-A'`, `{"dsn":"`+targetPath+`"}`)
	require.NoError(t, err)
	f.record, err = f.workspace.UpdateDatabaseSetup(t.Context(), f.record.DatabaseSetupRevision, func(_ context.Context, _ *sql.Tx, record *workspace.Record) error {
		record.DatabaseConnectorID = "connector-A"
		return nil
	})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `INSERT INTO database_target_mappings
		(id,tag_id,connector_id,table_schema,table_name,column_name,write_mode,timestamp_column,write_interval_seconds,enabled)
		VALUES ('legacy-A','tag-A','connector-A','main','raw_values','temperature','insert','observed_at',15,1)`)
	require.NoError(t, err)
	connectors := dbtarget.NewSQLConnectorRepository(f.db)
	mappings := dbtarget.NewSQLTargetMappingRepository(f.db)
	f.router = NewRouter(&DatalinkServices{
		Workspace: f.workspace, WriteGroups: workspace.NewWriteGroupService(f.workspace, workspace.NewSQLWriteGroupRepository(f.db)),
		DBTarget: dbtarget.NewConnectorService(connectors, mappings), DBMapping: dbtarget.NewMappingService(mappings, connectors, nil),
		SourceRule: sourcerule.NewService(sourcerule.NewSQLRepository(f.db), nil, nil, nil),
	})
	return f, targetPath
}

func migrationReviewRequest(t *testing.T, f writeGroupRouterFixture) map[string]any {
	t.Helper()
	preview := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationPreviewPath, map[string]any{
		"workspace_id": f.record.ID, "source_ids": []string{"legacy-A"},
	})
	require.Equal(t, http.StatusOK, preview.Code, preview.Body.String())
	data := decodeJSONBody(t, preview)["data"].(map[string]any)
	return map[string]any{
		"workspace_id": data["workspace_id"], "expected_workspace_revision": data["workspace_revision"],
		"expected_connector_revision": data["connector_revision"], "review_digest": data["review_digest"],
		"source_ids": []string{"legacy-A"}, "confirm_snapshot_conversion": true,
	}
}

func TestNewRouter_LegacySingleMappingMigrationReviewSavesDraft(t *testing.T) {
	f, targetPath := newSingleMigrationRouterFixture(t)
	request := migrationReviewRequest(t, f)
	request["group"] = map[string]any{"id": "forged", "applied_revision": "forged", "status": "running"}
	request["apply"] = true
	response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	require.Equal(t, true, body["success"])
	data := body["data"].(map[string]any)
	require.NotEqual(t, request["expected_workspace_revision"], data["workspace_revision"])
	groups := data["groups"].([]any)
	require.Len(t, groups, 1)
	group := groups[0].(map[string]any)
	require.NotEmpty(t, group["id"])
	require.NotEqual(t, "forged", group["id"])
	require.Equal(t, "draft", group["status"])
	require.Equal(t, "", group["applied_revision"])
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath+"/"+group["id"].(string), nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	_, reloaded := groupSaveData(t, decodeJSONBody(t, read))
	require.Equal(t, group, reloaded)
	var enabled int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT enabled FROM database_target_mappings WHERE id = 'legacy-A'`).Scan(&enabled))
	require.Equal(t, 1, enabled)
	_, err := os.Stat(targetPath)
	require.ErrorIs(t, err, os.ErrNotExist)
	group["members"].([]any)[0].(map[string]any)["target_column"] = "reviewed_temperature"
	edit := performJSONRequest(t, f.router, http.MethodPut, writeGroupsPath+"/"+group["id"].(string), groupEditRequest(f, data, group))
	require.Equal(t, http.StatusOK, edit.Code, edit.Body.String())
	legacyRead := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/db-targets/mappings/legacy-A", nil)
	require.Equal(t, http.StatusOK, legacyRead.Code, legacyRead.Body.String())
	legacyData := decodeJSONBody(t, legacyRead)["data"].(map[string]any)
	require.Equal(t, "reviewed_temperature", legacyData["column_name"])
	require.Equal(t, group["id"], legacyData["canonical_group"].(map[string]any)["id"])
	require.Equal(t, "temperature", legacyData["legacy_intent"].(map[string]any)["column_name"])
	workspaceRead := performJSONRequest(t, f.router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	require.Equal(t, http.StatusOK, workspaceRead.Code, workspaceRead.Body.String())
	rows := decodeJSONBody(t, workspaceRead)["data"].([]any)
	require.Len(t, rows, 1)
	require.Equal(t, "reviewed_temperature", rows[0].(map[string]any)["column_name"])
}

func TestNewRouter_LegacySingleMappingMigrationReviewRejectsInvalidRequest(t *testing.T) {
	f, _ := newSingleMigrationRouterFixture(t)
	valid := migrationReviewRequest(t, f)
	for _, field := range []string{"workspace_id", "expected_workspace_revision", "expected_connector_revision", "review_digest", "source_ids", "confirm_snapshot_conversion"} {
		t.Run(field, func(t *testing.T) {
			request := make(map[string]any, len(valid))
			for key, value := range valid {
				request[key] = value
			}
			delete(request, field)
			response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
		})
	}
	valid["confirm_snapshot_conversion"] = false
	response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, valid)
	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID")
}

func TestNewRouter_LegacySingleMappingMigrationReviewCASAndReplay(t *testing.T) {
	f, _ := newSingleMigrationRouterFixture(t)
	request := migrationReviewRequest(t, f)
	request["review_digest"] = "stale-preview"
	stale := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, stale), "revision_mismatch")
	request = migrationReviewRequest(t, f)
	response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	first := decodeJSONBody(t, response)["data"]
	stale = performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusConflict, stale.Code, stale.Body.String())
	request = migrationReviewRequest(t, f)
	response = performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Equal(t, first, decodeJSONBody(t, response)["data"])
}

func TestNewRouter_LegacySingleMappingMigrationReviewSafeFailure(t *testing.T) {
	f, _ := newSingleMigrationRouterFixture(t)
	request := migrationReviewRequest(t, f)
	request["workspace_id"] = "foreign"
	response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	request = migrationReviewRequest(t, f)
	_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_migration_projection BEFORE UPDATE ON system_settings
		WHEN NEW.key = 'studio_v2_workspace' BEGIN SELECT RAISE(ABORT, 'private-migration-diagnostic'); END`)
	require.NoError(t, err)
	response = performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationReviewPath, request)
	require.Equal(t, http.StatusInternalServerError, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), "private-migration-diagnostic")
	assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_UNAVAILABLE")
	read := performJSONRequest(t, f.router, http.MethodGet, writeGroupsPath, nil)
	require.Equal(t, http.StatusOK, read.Code, read.Body.String())
	data := decodeJSONBody(t, read)["data"].(map[string]any)
	require.Empty(t, data["groups"])
	require.Equal(t, request["expected_workspace_revision"], data["workspace_revision"])
}
