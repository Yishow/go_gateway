package handlers

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Reproduce the persisted legacy shape: blank connector schema and a canonical
// group's compatibility projection in SQLite's main namespace.
func legacySQLiteNamespaceFixture(t *testing.T) (workspaceDatabaseFixture, workspaceDatabaseConfigRequest) {
	t.Helper()
	fixture := newWorkspaceDatabaseFixture(t)
	payload := validSQLiteConfigRequest(fixture)
	updateWorkspaceDatabaseConfig(t, fixture, payload)
	record, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(context.Background(), record.DatabaseConnectorID)
	require.NoError(t, err)
	var config map[string]any
	require.NoError(t, json.Unmarshal([]byte(connector.ConnectionConfig), &config))
	config["schema"] = ""
	raw, err := json.Marshal(config)
	require.NoError(t, err)
	connector.ConnectionConfig = string(raw)
	require.NoError(t, fixture.connectorRepo.Update(context.Background(), connector))
	links, err := fixture.ruleSvc.ListLinks(context.Background(), "rule-A")
	require.NoError(t, err)
	groupsService := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	created, err := groupsService.Create(context.Background(), workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, ExpectedConnectorRevision: connector.IdentityRevision,
		Group: &workspace.WriteGroup{WorkspaceID: record.ID, Name: "Existing recording", Members: []workspace.WriteGroupMember{{DeviceID: "dev-A", PointID: links[0].PointID, TagID: *links[0].TagID, TargetColumn: "line_a", Required: true}},
			Destination: workspace.WriteGroupDestination{ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision, TableName: payload.Table, StorageStrategy: workspace.WriteGroupStorageStrategyCustom}},
	})
	require.NoError(t, err)
	require.Equal(t, "main", created.Group.Destination.TableSchema)
	fixture.handler.writeGroups = groupsService
	record, err = fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	groups := record.DatabaseRowGroups
	payload.Schema = ""
	payload.RowGroups = []workspaceDatabaseRowGroup{{ID: groups[0].ID, TableSchema: "main", TableName: payload.Table, MemberPointIDs: groups[0].MemberPointIDs}}
	return fixture, payload
}

func TestWorkspaceSQLiteDefaultNamespace_NameOnlyPreservesIdentityAndGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture, payload := legacySQLiteNamespaceFixture(t)
	before, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connectorBefore, err := fixture.connectorSvc.GetByID(context.Background(), before.DatabaseConnectorID)
	require.NoError(t, err)
	groupsService := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	groupBefore, err := groupsService.Get(context.Background(), payload.RowGroups[0].ID)
	require.NoError(t, err)
	payload.Name = "Linux QA SQLite long connector name"
	payload.ExpectedSetupRevision = currentWorkspaceSetupRevision(t, fixture)
	resp := serveWorkspaceDatabaseConfig(t, fixture, payload)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	after, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connectorAfter, err := fixture.connectorSvc.GetByID(context.Background(), after.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, connectorBefore.ID, connectorAfter.ID)
	require.Equal(t, connectorBefore.IdentityRevision, connectorAfter.IdentityRevision)
	require.Equal(t, connectorConfigString(connectorBefore, "database"), connectorConfigString(connectorAfter, "database"))
	require.Equal(t, connectorConfigString(connectorBefore, "dsn"), connectorConfigString(connectorAfter, "dsn"))
	require.Equal(t, before.DatabaseRowGroups, after.DatabaseRowGroups)
	groupAfter, err := groupsService.Get(context.Background(), payload.RowGroups[0].ID)
	require.NoError(t, err)
	require.Equal(t, groupBefore.Group, groupAfter.Group, "name-only saves must preserve the canonical recording and its proof")
	require.Equal(t, payload.Name, connectorAfter.Name)
	require.Equal(t, "main", connectorConfigString(connectorAfter, "schema"))
}

func TestWorkspaceSQLiteDefaultNamespace_ReadbackCanonicalizesLegacyBlank(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture, _ := legacySQLiteNamespaceFixture(t)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newHandlerTestRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-config", http.NoBody)
	fixture.handler.GetConfig(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	body := decodeWorkspaceDatabaseBody(t, resp)
	require.Equal(t, "main", body["data"].(map[string]any)["schema"])
}

func TestWorkspaceSQLiteDefaultNamespace_ExplicitForeignScopeRemainsRejected(t *testing.T) {
	for _, field := range []string{"schema", "table", "setup_revision"} {
		t.Run(field, func(t *testing.T) {
			fixture, payload := legacySQLiteNamespaceFixture(t)
			before, err := fixture.workspaceSvc.GetOrCreate(context.Background())
			require.NoError(t, err)
			connectorBefore, err := fixture.connectorSvc.GetByID(context.Background(), before.DatabaseConnectorID)
			require.NoError(t, err)
			payload.ExpectedSetupRevision = currentWorkspaceSetupRevision(t, fixture)
			expected := http.StatusBadRequest
			switch field {
			case "schema":
				payload.Schema = "other_namespace"
			case "table":
				payload.Table = "other_table"
			case "setup_revision":
				payload.ExpectedSetupRevision = "stale"
				expected = http.StatusConflict
			}
			resp := serveWorkspaceDatabaseConfig(t, fixture, payload)
			require.Equal(t, expected, resp.Code, resp.Body.String())
			after, err := fixture.workspaceSvc.GetOrCreate(context.Background())
			require.NoError(t, err)
			connectorAfter, err := fixture.connectorSvc.GetByID(context.Background(), before.DatabaseConnectorID)
			require.NoError(t, err)
			require.Equal(t, before, after)
			require.Equal(t, connectorBefore, connectorAfter)
		})
	}
}

func TestWorkspaceSQLiteDefaultNamespace_OtherKindsRemainUnchanged(t *testing.T) {
	for _, kind := range []schema.DatabaseConnectorKind{schema.DatabaseConnectorKindPostgres, schema.DatabaseConnectorKindMySQL, schema.DatabaseConnectorKindSQLServer} {
		for _, value := range []string{"", "  ", "explicit"} {
			require.Equal(t, value, workspaceDatabaseAPISchema(kind, value))
		}
	}
}

func TestWorkspaceSQLiteDefaultNamespace_UnboundLongPathSavesNewIdentity(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	payload := validSQLiteConfigRequest(fixture)
	updateWorkspaceDatabaseConfig(t, fixture, payload)
	before, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connectorBefore, err := fixture.connectorSvc.GetByID(context.Background(), before.DatabaseConnectorID)
	require.NoError(t, err)
	payload.Database = filepath.Join(t.TempDir(), "long-linux-qa-sqlite-destination-path.db")
	payload.Schema = ""
	payload.ExpectedSetupRevision = currentWorkspaceSetupRevision(t, fixture)
	resp := serveWorkspaceDatabaseConfig(t, fixture, payload)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	after, err := fixture.workspaceSvc.GetOrCreate(context.Background())
	require.NoError(t, err)
	connectorAfter, err := fixture.connectorSvc.GetByID(context.Background(), before.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, before.DatabaseConnectorID, after.DatabaseConnectorID)
	require.NotEqual(t, connectorBefore.IdentityRevision, connectorAfter.IdentityRevision)
	require.Equal(t, payload.Database, connectorConfigString(connectorAfter, "database"))
	require.Empty(t, after.DatabaseRowGroups)
}

func TestWorkspaceSQLiteDefaultNamespace_BoundPathChangeDoesNotRebindCanonicalGroup(t *testing.T) {
	fixture, payload := legacySQLiteNamespaceFixture(t)
	ctx := context.Background()
	service := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	before, err := service.Get(ctx, payload.RowGroups[0].ID)
	require.NoError(t, err)
	payload.Database = filepath.Join(t.TempDir(), "different-target.db")
	payload.ExpectedSetupRevision = currentWorkspaceSetupRevision(t, fixture)
	resp := serveWorkspaceDatabaseConfig(t, fixture, payload)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	after, err := service.Get(ctx, before.Group.ID)
	require.NoError(t, err)
	require.Equal(t, before.Group, after.Group, "a connection save must not migrate the recording's destination, members, or proof")
	connector, err := fixture.connectorSvc.GetByID(ctx, before.Group.Destination.ConnectorID)
	require.NoError(t, err)
	require.NotEqual(t, before.Group.Destination.ConnectorRevision, connector.IdentityRevision)
	readiness, err := service.Readiness(ctx, before.Group.ID)
	require.NoError(t, err)
	require.False(t, readiness.Ready)
	require.False(t, readiness.ConfigReady)
	codes := make([]string, 0, len(readiness.Issues))
	for _, issue := range readiness.Issues {
		codes = append(codes, issue.Code)
	}
	require.Contains(t, codes, "connector-revision-stale")
}
