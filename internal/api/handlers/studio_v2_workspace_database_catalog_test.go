package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogReturnsSavedSQLiteCatalog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	wireCatalogManagedGuard(fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)

	query := url.Values{
		"workspace_id":                {record.ID},
		"expected_workspace_revision": {record.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
	}
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	fixture.handler.GetCatalog(ctx)

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeWorkspaceDatabaseBody(t, response)
	data := body["data"].(map[string]any)
	require.Equal(t, record.ID, data["workspace_id"])
	require.Equal(t, record.DatabaseSetupRevision, data["workspace_revision"])
	require.Equal(t, connector.ID, data["connector_id"])
	require.Equal(t, connector.IdentityRevision, data["connector_revision"])
	require.Equal(t, "sqlite", data["kind"])
	require.Equal(t, "main", data["default_schema"])
	schemas := data["schemas"].([]any)
	require.NotEmpty(t, schemas)
	require.Equal(t, "main", schemas[0].(map[string]any)["name"])
	require.IsType(t, []any{}, schemas[0].(map[string]any)["tables"])
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogHonorsWorkspaceAndConnectorRevision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	wireCatalogManagedGuard(fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)

	for _, test := range []struct {
		name                 string
		expectedWorkspaceRev string
		expectedConnectorRev string
		wantStatus           int
	}{
		{name: "workspace revision", expectedWorkspaceRev: "stale-workspace", expectedConnectorRev: connector.IdentityRevision, wantStatus: http.StatusConflict},
		{name: "connector revision", expectedWorkspaceRev: record.DatabaseSetupRevision, expectedConnectorRev: "stale-connector", wantStatus: http.StatusConflict},
	} {
		t.Run(test.name, func(t *testing.T) {
			query := url.Values{
				"workspace_id":                {record.ID},
				"expected_workspace_revision": {test.expectedWorkspaceRev},
				"connector_id":                {connector.ID},
				"expected_connector_revision": {test.expectedConnectorRev},
			}
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			fixture.handler.GetCatalog(ctx)

			require.Equal(t, test.wantStatus, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), fixture.targetDB)
			var envelope map[string]any
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
			require.Equal(t, false, envelope["success"])
		})
	}
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogSupportsLocatorScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	wireCatalogManagedGuard(fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	query := url.Values{
		"workspace_id":                {record.ID},
		"expected_workspace_revision": {record.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
		"locator_workspace_id":        {record.ID},
		"locator_workspace_revision":  {record.DatabaseSetupRevision},
	}
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	fixture.handler.GetCatalog(ctx)

	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogRequiresGroupRevision(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	wireCatalogManagedGuard(fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	query := url.Values{
		"workspace_id":                {record.ID},
		"expected_workspace_revision": {record.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
		"group_id":                    {"group-1"},
	}
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	fixture.handler.GetCatalog(ctx)

	require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
	require.NotContains(t, response.Body.String(), fixture.targetDB)
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogUsesSavedGroupScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	var deviceID, pointID, tagID string
	require.NoError(t, fixture.mainDB.QueryRowContext(t.Context(), `
		SELECT p.device_id, m.point_id, m.tag_id
		FROM mappings m JOIN points p ON p.id = m.point_id
		WHERE m.enabled = 1 ORDER BY m.id LIMIT 1`).Scan(&deviceID, &pointID, &tagID))
	groups := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	saved, err := groups.Create(t.Context(), workspace.WriteGroupMutation{
		WorkspaceID: record.ID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision,
		ExpectedConnectorRevision: connector.IdentityRevision,
		Group: &workspace.WriteGroup{
			WorkspaceID: record.ID, Name: "Catalog saved group",
			Members: []workspace.WriteGroupMember{{
				DeviceID: deviceID, PointID: pointID, TagID: tagID, TargetColumn: "line_a", Required: true,
			}},
			Destination: workspace.WriteGroupDestination{
				ConnectorID: connector.ID, ConnectorRevision: connector.IdentityRevision,
				TableSchema: "main", TableName: "sensor_values", StorageStrategy: workspace.WriteGroupStorageStrategyCustom,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, saved.Group)
	fixture.handler.WithWriteGroups(groups)
	current, err := fixture.workspaceSvc.ReadSnapshot(t.Context())
	require.NoError(t, err)

	serve := func(query url.Values) *httptest.ResponseRecorder {
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
		response := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(response)
		ctx.Request = request
		fixture.handler.GetCatalog(ctx)
		return response
	}
	base := url.Values{
		"workspace_id":                {current.ID},
		"expected_workspace_revision": {current.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
		"group_id":                    {saved.Group.ID},
		"expected_group_revision":     {saved.Group.Revision},
	}
	require.Equal(t, http.StatusOK, serve(base).Code)

	staleGroup := cloneCatalogQuery(base)
	staleGroup.Set("expected_group_revision", "stale-group-revision")
	staleResponse := serve(staleGroup)
	require.Equal(t, http.StatusConflict, staleResponse.Code, staleResponse.Body.String())
	require.NotContains(t, staleResponse.Body.String(), fixture.targetDB)

	alternate, err := fixture.connectorSvc.Create(t.Context(), dbtarget.CreateConnectorRequest{
		Name: "Catalog alternate", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{"dsn": fixture.targetDB, "schema": "main"},
	})
	require.NoError(t, err)
	mismatchedConnector := cloneCatalogQuery(base)
	mismatchedConnector.Set("connector_id", alternate.ID)
	mismatchedConnector.Set("expected_connector_revision", alternate.IdentityRevision)
	mismatchedResponse := serve(mismatchedConnector)
	require.Equal(t, http.StatusConflict, mismatchedResponse.Code, mismatchedResponse.Body.String())
	require.NotContains(t, mismatchedResponse.Body.String(), fixture.targetDB)
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogFailsClosedWithoutSQLiteManagedGuard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	query := url.Values{
		"workspace_id":                {record.ID},
		"expected_workspace_revision": {record.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
	}
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	fixture.handler.GetCatalog(ctx)

	require.Equal(t, http.StatusServiceUnavailable, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "WRITE_GROUP_UNAVAILABLE")
	require.NotContains(t, response.Body.String(), fixture.targetDB)
}

func TestStudioV2WorkspaceDatabaseHandler_GetCatalogRejectsInternalSQLiteWhenManagedGuardIsWired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	var internalPath string
	require.NoError(t, fixture.mainDB.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(new(int), new(string), &internalPath))
	connector := &schema.DatabaseConnector{
		ID: "catalog-internal-sqlite", Name: "Internal SQLite", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"` + internalPath + `"}`, IdentityRevision: "internal-revision",
		Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, fixture.connectorRepo.Create(t.Context(), connector))
	groups := workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB))
	fixture.handler.WithWriteGroups(groups)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	query := url.Values{
		"workspace_id":                {record.ID},
		"expected_workspace_revision": {record.DatabaseSetupRevision},
		"connector_id":                {connector.ID},
		"expected_connector_revision": {connector.IdentityRevision},
	}
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-catalog?"+query.Encode(), http.NoBody)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	fixture.handler.GetCatalog(ctx)

	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "WRITE_GROUP_INTERNAL_DATABASE")
	require.NotContains(t, response.Body.String(), internalPath)
}

func wireCatalogManagedGuard(fixture workspaceDatabaseFixture) {
	fixture.handler.WithWriteGroups(workspace.NewWriteGroupService(fixture.workspaceSvc, workspace.NewSQLWriteGroupRepository(fixture.mainDB)))
}

func cloneCatalogQuery(source url.Values) url.Values {
	clone := make(url.Values, len(source))
	for key, values := range source {
		clone[key] = append([]string(nil), values...)
	}
	return clone
}
