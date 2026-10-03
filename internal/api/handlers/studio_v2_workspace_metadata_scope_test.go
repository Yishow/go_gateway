package handlers

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type metadataChangingGroupReader struct {
	result *workspace.WriteGroupListResult
	reads  int
}

func (r *metadataChangingGroupReader) List(context.Context) (*workspace.WriteGroupListResult, error) {
	r.reads++
	result := *r.result
	group := *r.result.Groups[0]
	if r.reads > 1 {
		group.Revision = "changed-during-inspection"
	}
	result.Groups = []*workspace.WriteGroup{&group}
	return &result, nil
}

func TestWorkspaceGroupMetadataRejectsScopeChangedDuringInspection(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	workspaceID, connector := boundWorkspaceConnector(t, f)
	group := &workspace.WriteGroup{
		ID: "metadata-race", WorkspaceID: workspaceID, Revision: "saved-group-revision",
		Destination: workspace.WriteGroupDestination{ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision, TableSchema: "main", TableName: "sensor_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom},
	}
	f.handler.writeGroups = &metadataChangingGroupReader{result: &workspace.WriteGroupListResult{WorkspaceID: workspaceID, Groups: []*workspace.WriteGroup{group}}}
	resp, _ := serveGroupMetadata(t, f, url.Values{
		"group_id": {group.ID}, "expected_group_revision": {group.Revision}, "expected_connector_revision": {connector.IdentityRevision},
	})
	require.Equal(t, http.StatusConflict, resp.Code, "metadata must not report old revision with columns read for a changed scope: %s", resp.Body.String())
	require.NotContains(t, resp.Body.String(), "line_a")
}
