package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type workspaceDatabaseMetadataBody struct {
	Success bool `json:"success"`
	Data    *struct {
		WorkspaceID       string `json:"workspace_id"`
		ConnectorID       string `json:"connector_id"`
		ConnectorRevision string `json:"connector_revision"`
		Database          string `json:"database"`
		Schema            string `json:"schema"`
		Table             string `json:"table"`
		InspectionStatus  string `json:"inspection_status"`
		Reason            string `json:"reason"`
		Columns           []struct {
			Name       string `json:"name"`
			DataType   string `json:"data_type"`
			PrimaryKey bool   `json:"primary_key"`
		} `json:"columns"`
	} `json:"data"`
}

func serveWorkspaceDatabaseMetadata(t *testing.T, fixture workspaceDatabaseFixture, expectedRevision string) (*httptest.ResponseRecorder, workspaceDatabaseMetadataBody) {
	t.Helper()
	target := "/api/v1/datalink/studio-v2/workspace/database-metadata"
	if expectedRevision != "" {
		target += "?expected_connector_revision=" + url.QueryEscape(expectedRevision)
	}
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, target, http.NoBody)
	fixture.handler.GetMetadata(c)
	var body workspaceDatabaseMetadataBody
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	return resp, body
}

func boundWorkspaceConnector(t *testing.T, fixture workspaceDatabaseFixture) (workspaceID string, connector *schema.DatabaseConnector) {
	t.Helper()
	record, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connector, err = fixture.connectorSvc.GetByID(context.Background(), record.DatabaseConnectorID)
	require.NoError(t, err)
	return record.ID, connector
}

func TestWorkspaceDatabaseMetadata_ReportsSavedTargetColumns(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	workspaceID, connector := boundWorkspaceConnector(t, fixture)

	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, connector.IdentityRevision)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.True(t, body.Success)
	require.NotNil(t, body.Data)
	require.Equal(t, workspaceID, body.Data.WorkspaceID)
	require.Equal(t, connector.ID, body.Data.ConnectorID)
	require.Equal(t, connector.IdentityRevision, body.Data.ConnectorRevision)
	require.Equal(t, fixture.targetDB, body.Data.Database)
	require.Equal(t, "main", body.Data.Schema)
	require.Equal(t, "sensor_values", body.Data.Table)
	require.Equal(t, "exists", body.Data.InspectionStatus)
	require.Len(t, body.Data.Columns, 3)
	require.Equal(t, []string{"ts", "line_a", "line_b"}, []string{body.Data.Columns[0].Name, body.Data.Columns[1].Name, body.Data.Columns[2].Name})
	require.Equal(t, []string{"DATETIME", "REAL", "REAL"}, []string{body.Data.Columns[0].DataType, body.Data.Columns[1].DataType, body.Data.Columns[2].DataType})
	require.True(t, body.Data.Columns[0].PrimaryKey)
}

func TestWorkspaceDatabaseMetadata_NamesTargetLikePreviewScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	_, connector := boundWorkspaceConnector(t, fixture)

	// Rewrite the saved config so the SQLite file lives under "path" instead of
	// "database" and no schema is configured. Metadata must name the target
	// exactly like the schema preview scope does: the file location and the
	// adapter default schema.
	var config map[string]any
	require.NoError(t, json.Unmarshal([]byte(connector.ConnectionConfig), &config))
	config["path"] = fixture.targetDB
	delete(config, "database")
	delete(config, "schema")
	configJSON, err := json.Marshal(config)
	require.NoError(t, err)
	connector.ConnectionConfig = string(configJSON)
	require.NoError(t, fixture.connectorRepo.Update(context.Background(), connector))
	_, connector = boundWorkspaceConnector(t, fixture)

	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, connector.IdentityRevision)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.True(t, body.Success)
	require.NotNil(t, body.Data)
	require.Equal(t, fixture.targetDB, body.Data.Database)
	require.Equal(t, "main", body.Data.Schema)
	require.Equal(t, "sensor_values", body.Data.Table)
	require.Equal(t, "exists", body.Data.InspectionStatus)
}

func TestWorkspaceDatabaseMetadata_AbsentSavedTableIsMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	missing := validSQLiteConfigRequest(fixture)
	missing.Table = "sensor_values_missing"
	updateWorkspaceDatabaseConfig(t, fixture, missing)
	_, connector := boundWorkspaceConnector(t, fixture)

	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, connector.IdentityRevision)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Equal(t, "missing", body.Data.InspectionStatus)
	require.Equal(t, "sensor_values_missing", body.Data.Table)
	require.Empty(t, body.Data.Columns)
}

func TestWorkspaceDatabaseMetadata_UnreadableTargetIsFailedNotMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	corrupt := filepath.Join(t.TempDir(), "corrupt.db")
	require.NoError(t, os.WriteFile(corrupt, []byte("this is not a sqlite database file at all, only text"), 0o600))
	broken := validSQLiteConfigRequest(fixture)
	broken.Database = corrupt
	updateWorkspaceDatabaseConfig(t, fixture, broken)
	_, connector := boundWorkspaceConnector(t, fixture)

	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, connector.IdentityRevision)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Equal(t, "failed", body.Data.InspectionStatus)
	require.NotEmpty(t, body.Data.Reason)
	require.Empty(t, body.Data.Columns)
}

func TestWorkspaceDatabaseMetadata_RejectsStaleOrMissingRevisionAndUnboundWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	unbound := newWorkspaceDatabaseFixture(t)
	resp, _ := serveWorkspaceDatabaseMetadata(t, unbound, "any-revision")
	require.Equal(t, http.StatusNotFound, resp.Code, resp.Body.String())

	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, "stale-revision")
	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	require.Nil(t, body.Data)
	resp, body = serveWorkspaceDatabaseMetadata(t, fixture, "")
	require.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
	require.Nil(t, body.Data)
}
