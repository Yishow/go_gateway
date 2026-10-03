package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type metadataGroupReader struct {
	result *workspace.WriteGroupListResult
}

func (r metadataGroupReader) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return r.result, nil
}

func TestWorkspaceDatabaseMetadata_UsesSavedGroupTable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	workspaceID, connector := boundWorkspaceConnector(t, f)
	target, err := sql.Open("sqlite", f.targetDB)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, target.Close()) })
	_, err = target.ExecContext(t.Context(), `CREATE TABLE group_b_values (counter TEXT NOT NULL)`)
	require.NoError(t, err)
	group := &workspace.WriteGroup{
		ID: "group-B", WorkspaceID: workspaceID, Revision: "group-1", Status: workspace.WriteGroupStatusDraft,
		Destination: workspace.WriteGroupDestination{ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision, TableSchema: "main", TableName: "group_b_values"},
	}
	f.handler.writeGroups = metadataGroupReader{&workspace.WriteGroupListResult{WorkspaceID: workspaceID, Groups: []*workspace.WriteGroup{group}}}
	query := url.Values{"group_id": {group.ID}, "expected_group_revision": {group.Revision}, "expected_connector_revision": {connector.IdentityRevision}}
	resp, body := serveGroupMetadata(t, f, query)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.Equal(t, "group_b_values", body.Data.Table)
	require.Len(t, body.Data.Columns, 1)
	require.Equal(t, "counter", body.Data.Columns[0].Name)
}

func TestWorkspaceDatabaseMetadata_RejectsUnprovenGroupScope(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		change func(*workspace.WriteGroup, url.Values)
		status int
	}{
		{"unknown", func(_ *workspace.WriteGroup, q url.Values) { q.Set("group_id", "missing") }, http.StatusNotFound},
		{"foreign workspace", func(g *workspace.WriteGroup, _ url.Values) { g.WorkspaceID = "foreign" }, http.StatusNotFound},
		{"deleted", func(g *workspace.WriteGroup, _ url.Values) { g.Status = workspace.WriteGroupStatusDeleted }, http.StatusNotFound},
		{"stale group", func(_ *workspace.WriteGroup, q url.Values) { q.Set("expected_group_revision", "old") }, http.StatusConflict},
		{"missing group revision", func(_ *workspace.WriteGroup, q url.Values) { q.Del("expected_group_revision") }, http.StatusBadRequest},
		{"stale connector", func(g *workspace.WriteGroup, _ url.Values) { g.Destination.ConnectorRevision = "old" }, http.StatusConflict},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			f := newWorkspaceDatabaseFixture(t)
			saveWorkspaceDatabaseConfig(t, f)
			workspaceID, connector := boundWorkspaceConnector(t, f)
			group := &workspace.WriteGroup{ID: "group-B", WorkspaceID: workspaceID, Revision: "group-1", Status: workspace.WriteGroupStatusDraft,
				Destination: workspace.WriteGroupDestination{ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision, TableName: "sensor_values"}}
			query := url.Values{"group_id": {group.ID}, "expected_group_revision": {group.Revision}, "expected_connector_revision": {connector.IdentityRevision}}
			scenario.change(group, query)
			f.handler.writeGroups = metadataGroupReader{&workspace.WriteGroupListResult{WorkspaceID: workspaceID, Groups: []*workspace.WriteGroup{group}}}
			resp, body := serveGroupMetadata(t, f, query)
			require.Equal(t, scenario.status, resp.Code, resp.Body.String())
			require.Nil(t, body.Data)
		})
	}
}

func serveGroupMetadata(t *testing.T, f workspaceDatabaseFixture, query url.Values) (*httptest.ResponseRecorder, workspaceDatabaseMetadataBody) {
	t.Helper()
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/database-metadata?"+query.Encode(), http.NoBody)
	f.handler.GetMetadata(c)
	var body workspaceDatabaseMetadataBody
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	return resp, body
}
