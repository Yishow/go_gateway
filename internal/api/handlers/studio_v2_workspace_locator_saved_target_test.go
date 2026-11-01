package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type locatorSavedTargetReader struct {
	connector *schema.DatabaseConnector
	err       error
}

func (r locatorSavedTargetReader) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return r.connector, r.err
}

func TestStudioLocatorEditorUsesExactSavedTarget(t *testing.T) {
	f := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, f)
	target := filepath.Join(t.TempDir(), "alternate.db")
	alternate, err := f.connectorSvc.Create(t.Context(), dbtarget.CreateConnectorRequest{Name: "Alternate", Kind: schema.DatabaseConnectorKindSQLite, ConnectionConfig: dbtarget.ConnectionConfig{"dsn": target, "password": "never-return-this", "table": "connector_default"}})
	require.NoError(t, err)
	record, err := f.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	originalDefault := record.DatabaseConnectorID
	var d, p, tagID string
	require.NoError(t, f.mainDB.QueryRowContext(t.Context(), `SELECT p.device_id,m.point_id,m.tag_id FROM mappings m JOIN points p ON p.id=m.point_id WHERE m.enabled=1 ORDER BY m.id LIMIT 1`).Scan(&d, &p, &tagID))
	groups := workspace.NewWriteGroupService(f.workspaceSvc, workspace.NewSQLWriteGroupRepository(f.mainDB))
	created, err := groups.Create(t.Context(), workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: alternate.IdentityRevision, Group: &workspace.WriteGroup{WorkspaceID: record.ID, Name: "Saved alternate", Members: []workspace.WriteGroupMember{{DeviceID: d, PointID: p, TagID: tagID, Required: true}}, Destination: workspace.WriteGroupDestination{ConnectorID: alternate.ID, ConnectorRevision: alternate.IdentityRevision, TableSchema: "main", TableName: "saved_group_table", StorageStrategy: workspace.WriteGroupStorageStrategyManaged}, RowPolicy: workspace.WriteGroupRowPolicy{IntervalSeconds: 15, RecordKeyColumn: "record_id"}}})
	require.NoError(t, err)
	record, err = f.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	before, err := json.Marshal(record)
	require.NoError(t, err)
	h := NewStudioV2WorkspaceHandler(f.workspaceSvc).WithLocatorEntryServices(StudioV2LocatorReaders{Connectors: f.connectorSvc, Targets: f.dbMappingSvc, Devices: device.NewSQLRepository(f.mainDB), Points: point.NewSQLRepository(f.mainDB), Tags: tag.NewSQLRepository(f.mainDB), Mappings: mapping.NewSQLRepository(f.mainDB), Rules: f.ruleSvc, Groups: groups})
	// Isolate connector hydration; source-binding guards have independent locator tests.
	editorGroup := *created.Group
	editorGroup.Members = nil
	observation, err := h.readStudioLocatorEditor(t.Context(), record, &editorGroup)
	require.NoError(t, err)
	require.Equal(t, alternate.ID, observation.Connector.ID)
	entry := &studioLocatorEntry{}
	require.NoError(t, projectStudioLocatorEditor(record, observation, entry))
	require.Equal(t, alternate.ID, entry.DatabaseConfig.ID)
	require.Equal(t, alternate.IdentityRevision, entry.DatabaseConfig.IdentityRevision)
	require.Equal(t, created.Group.Destination.Database, entry.DatabaseConfig.Database)
	require.Equal(t, created.Group.Destination.TableSchema, entry.DatabaseConfig.Schema)
	require.Equal(t, created.Group.Destination.TableName, entry.DatabaseConfig.Table)
	require.Equal(t, record.DatabaseSetupRevision, entry.DatabaseConfig.SetupRevision)
	raw, err := json.Marshal(entry)
	require.NoError(t, err)
	require.NotContains(t, string(raw), "never-return-this")
	require.NotContains(t, string(raw), "password")
	require.Empty(t, entry.DatabaseConfig.RowGroups)
	scope, err := groups.ManagedSchemaScope(t.Context(), created.Group.ID, workspace.WriteGroupMutation{WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedGroupRevision: created.Group.Revision, ExpectedConnectorRevision: alternate.IdentityRevision})
	require.NoError(t, err)
	require.Equal(t, alternate.ID, scope.ConnectorID)
	require.Equal(t, entry.DatabaseConfig.Database, scope.Database)
	require.Equal(t, entry.DatabaseConfig.Table, scope.TablePrefix)
	f.handler.WithWriteGroups(groups)
	response, body := serveGroupMetadata(t, f, url.Values{"group_id": {created.Group.ID}, "expected_group_revision": {created.Group.Revision}, "expected_connector_revision": {alternate.IdentityRevision}})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.NotNil(t, body.Data)
	require.Equal(t, "missing", body.Data.InspectionStatus)
	require.NoFileExists(t, target)
	afterRecord, err := f.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)
	after, err := json.Marshal(afterRecord)
	require.NoError(t, err)
	require.Equal(t, string(before), string(after))
	require.Equal(t, originalDefault, afterRecord.DatabaseConnectorID)
	for name, reader := range map[string]locatorSavedTargetReader{
		"missing": {err: sql.ErrNoRows},
		"foreign": {connector: &schema.DatabaseConnector{ID: originalDefault, IdentityRevision: alternate.IdentityRevision}},
		"stale":   {connector: &schema.DatabaseConnector{ID: alternate.ID, IdentityRevision: "stale"}},
	} {
		t.Run(name, func(t *testing.T) {
			denied := *h
			denied.locatorReaders.Connectors = reader
			_, err := denied.readStudioLocatorEditor(t.Context(), record, &editorGroup)
			require.Error(t, err)
		})
	}
}
