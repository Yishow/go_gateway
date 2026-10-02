package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const singleMappingMigrationPreviewPath = writeGroupsPath + "/migrations/single-mappings/preview"

func TestNewRouter_LegacySingleMappingMigrationPreviewPreservesOriginal(t *testing.T) {
	f := newWriteGroupRouterFixture(t)
	targetPath := filepath.Join(t.TempDir(), "untouched-target.db")
	config, err := json.Marshal(map[string]string{"dsn": targetPath})
	require.NoError(t, err)
	_, err = f.db.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = 'connector-A'`, string(config))
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

	request := map[string]any{"workspace_id": f.record.ID, "source_ids": []string{"legacy-A"}}
	response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationPreviewPath, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeJSONBody(t, response)
	require.Equal(t, true, body["success"])
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, f.record.ID, data["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, data["workspace_revision"])
	require.Equal(t, "connector-1", data["connector_revision"])
	require.NotEmpty(t, data["review_digest"])
	items, ok := data["items"].([]any)
	require.True(t, ok)
	require.Len(t, items, 1)
	item := items[0].(map[string]any)
	require.Equal(t, "legacy-A", item["source_id"])
	require.Equal(t, "needs_review", item["status"])
	require.NotEmpty(t, item["source_revision"])
	require.NotEmpty(t, item["differences"])
	second := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationPreviewPath, request)
	require.Equal(t, http.StatusOK, second.Code, second.Body.String())
	require.Equal(t, body, decodeJSONBody(t, second))

	current, err := f.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, f.record, current)
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM write_groups`).Scan(&count))
	require.Zero(t, count)
	var mode, timestamp string
	var interval, enabled int
	require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT write_mode,timestamp_column,write_interval_seconds,enabled FROM database_target_mappings WHERE id = 'legacy-A'`).Scan(&mode, &timestamp, &interval, &enabled))
	require.Equal(t, "insert", mode)
	require.Equal(t, "observed_at", timestamp)
	require.Equal(t, 15, interval)
	require.Equal(t, 1, enabled)
	_, err = os.Stat(targetPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestNewRouter_LegacySingleMappingMigrationPreviewRejectsInvalidScope(t *testing.T) {
	for _, request := range []map[string]any{
		{"workspace_id": "foreign", "source_ids": []string{"missing"}},
		{"workspace_id": "unknown", "source_ids": []string{"missing"}},
	} {
		f := newWriteGroupRouterFixture(t)
		response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationPreviewPath, request)
		require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_NOT_FOUND")
	}
	for _, request := range []map[string]any{
		{}, {"workspace_id": "workspace"}, {"workspace_id": "workspace", "source_ids": []string{}},
	} {
		f := newWriteGroupRouterFixture(t)
		response := performJSONRequest(t, f.router, http.MethodPost, singleMappingMigrationPreviewPath, request)
		require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
		assertWriteGroupError(t, decodeJSONBody(t, response), "WRITE_GROUP_INVALID_REQUEST")
	}
}
