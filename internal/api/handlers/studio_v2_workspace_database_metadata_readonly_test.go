package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestWorkspaceDatabaseMetadata_MissingSQLiteFileDoesNotCreate(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	_, connector, target := pointMetadataAtMissingSQLiteFile(t, fixture)

	resp, body := serveWorkspaceDatabaseMetadata(t, fixture, connector.IdentityRevision)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.NotNil(t, body.Data)
	require.Equal(t, "failed", body.Data.InspectionStatus)
	require.NoFileExists(t, target)
}

func TestWorkspaceDatabaseMetadata_CustomGroupMissingSQLiteFileDoesNotCreate(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	workspaceID, connector, target := pointMetadataAtMissingSQLiteFile(t, fixture)
	group := &workspace.WriteGroup{
		ID: "group-missing", WorkspaceID: workspaceID, Revision: "group-revision",
		Destination: workspace.WriteGroupDestination{
			ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision,
			TableSchema: "main", TableName: "group_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom,
		},
	}
	fixture.handler.writeGroups = metadataGroupReader{result: &workspace.WriteGroupListResult{
		WorkspaceID: workspaceID, Groups: []*workspace.WriteGroup{group},
	}}
	query := url.Values{
		"group_id":                    {group.ID},
		"expected_group_revision":     {group.Revision},
		"expected_connector_revision": {connector.IdentityRevision},
	}
	resp, body := serveGroupMetadata(t, fixture, query)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.NotNil(t, body.Data)
	require.Equal(t, "failed", body.Data.InspectionStatus)
	require.NoFileExists(t, target)
}

func TestWorkspaceDatabaseMetadata_ManagedGroupMissingSQLiteFileIsMissing(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	workspaceID, connector := boundWorkspaceConnector(t, fixture)
	target := filepath.Join(t.TempDir(), "missing-managed-metadata.db")
	config, err := json.Marshal(map[string]string{"dsn": target, "schema": "main", "table": "managed_values"})
	require.NoError(t, err)
	_, err = fixture.mainDB.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), connector.ID)
	require.NoError(t, err)

	var deviceID, pointID, tagID string
	err = fixture.mainDB.QueryRowContext(t.Context(), `
		SELECT p.device_id, m.point_id, m.tag_id
		FROM mappings m JOIN points p ON p.id = m.point_id
		WHERE m.enabled = 1 ORDER BY m.id LIMIT 1
	`).Scan(&deviceID, &pointID, &tagID)
	require.NoError(t, err)
	groups := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	created, err := groups.Create(t.Context(), workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedConnectorRevision: connector.IdentityRevision,
		Group: &workspace.WriteGroup{
			WorkspaceID: workspaceID, Name: "managed metadata",
			Members: []workspace.WriteGroupMember{{DeviceID: deviceID, PointID: pointID, TagID: tagID, Required: true}},
			Destination: workspace.WriteGroupDestination{
				ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision,
				Database: target, TableSchema: "main", TableName: "managed_values",
				StorageStrategy: workspace.WriteGroupStorageStrategyManaged,
			},
			RowPolicy: workspace.WriteGroupRowPolicy{IntervalSeconds: 15},
		},
	})
	require.NoError(t, err)
	fixture.handler.WithWriteGroups(groups)
	query := url.Values{
		"group_id":                    {created.Group.ID},
		"expected_group_revision":     {created.Group.Revision},
		"expected_connector_revision": {connector.IdentityRevision},
	}
	resp, body := serveGroupMetadata(t, fixture, query)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.NotNil(t, body.Data)
	require.Equal(t, "missing", body.Data.InspectionStatus)
	require.NoFileExists(t, target)
}

func pointMetadataAtMissingSQLiteFile(t *testing.T, fixture workspaceDatabaseFixture) (workspaceID string, connector *schema.DatabaseConnector, target string) {
	t.Helper()
	saveWorkspaceDatabaseConfig(t, fixture)
	workspaceID, connector = boundWorkspaceConnector(t, fixture)
	target = filepath.Join(t.TempDir(), "missing-metadata.db")
	config, err := json.Marshal(map[string]string{"dsn": target, "schema": "main", "table": "sensor_values"})
	require.NoError(t, err)
	_, err = fixture.mainDB.ExecContext(t.Context(), `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), connector.ID)
	require.NoError(t, err)
	return workspaceID, connector, target
}
