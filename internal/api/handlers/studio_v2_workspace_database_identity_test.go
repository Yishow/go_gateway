package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSavedConnectorIdentityWorkspaceSelection(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	before, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	pool, err := fixture.connectorSvc.Create(t.Context(), dbtarget.CreateConnectorRequest{
		Name: "saved pool", Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: dbtarget.ConnectionConfig{
			"host": "127.0.0.1", "port": 1, "database": "pool_metrics", "user": "pool_writer",
			"password": " test fixture password ", "schema": "public", "table": "pool_samples",
		},
	})
	require.NoError(t, err)
	payload, err := json.Marshal(pool)
	require.NoError(t, err)
	var persisted map[string]any
	require.NoError(t, json.Unmarshal(payload, &persisted))
	response := updateSavedConnectorConfig(t, fixture, map[string]any{
		"connector_id": pool.ID, "expected_connector_revision": persisted["identity_revision"], "expected_setup_revision": currentWorkspaceSetupRevision(t, fixture),
		"kind": "postgres", "name": "saved pool", "database": "pool_metrics", "username": "pool_writer",
		"host": "127.0.0.1", "port": 1, "schema": "public", "table": "pool_samples",
		"write_mode": "insert", "write_interval_seconds": 5, "timestamp_column": "ts",
	})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	body := decodeWorkspaceDatabaseBody(t, response)["data"].(map[string]any)
	require.Equal(t, pool.ID, body["id"])
	require.Equal(t, "postgres", body["kind"])
	require.NotEmpty(t, body["identity_revision"])
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, pool.ID, record.DatabaseConnectorID)
	old, err := fixture.connectorSvc.GetByID(t.Context(), before.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, schema.DatabaseConnectorKindSQLite, old.Kind)
	after, err := fixture.connectorSvc.GetByID(t.Context(), pool.ID)
	require.NoError(t, err)
	require.JSONEq(t, pool.ConnectionConfig, after.ConnectionConfig)
	require.NotContains(t, response.Body.String(), "test fixture password")
}

func TestSavedConnectorIdentityWorkspaceRejectsMissingTarget(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	before, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	response := updateSavedConnectorConfig(t, fixture, map[string]any{
		"connector_id": "missing-connector", "expected_connector_revision": "missing-revision",
		"kind": "sqlite", "name": "unknown", "database": fixture.targetDB,
		"table": "samples", "write_mode": "insert", "write_interval_seconds": 5,
	})
	require.Equal(t, http.StatusNotFound, response.Code, response.Body.String())
	after, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func TestSavedConnectorIdentityWorkspaceRejectsDisabledConnector(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	before, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	pool, err := fixture.connectorSvc.Create(t.Context(), dbtarget.CreateConnectorRequest{
		Name: "disabled pool", Kind: schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: dbtarget.ConnectionConfig{"dsn": fixture.targetDB, "database": fixture.targetDB, "table": "sensor_values"},
	})
	require.NoError(t, err)
	disabled := *pool
	disabled.Enabled = false
	require.NoError(t, fixture.connectorRepo.UpdateWithExpectedIdentityRevision(t.Context(), &disabled, pool.IdentityRevision))

	response := updateSavedConnectorConfig(t, fixture, map[string]any{
		"connector_id": pool.ID, "expected_connector_revision": pool.IdentityRevision, "expected_setup_revision": currentWorkspaceSetupRevision(t, fixture),
		"kind": "sqlite", "name": "disabled pool", "database": fixture.targetDB, "table": "sensor_values",
		"write_mode": "insert", "write_interval_seconds": 5,
	})

	require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
	require.Contains(t, response.Body.String(), "disabled")
	after, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func updateSavedConnectorConfig(t *testing.T, fixture workspaceDatabaseFixture, payload any) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request = newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/database-config", payload)
	fixture.handler.UpdateConfig(c)
	return response
}

func TestSavedConnectorIdentityExplicitClearWinsOnCreate(t *testing.T) {
	fixture := newWorkspaceDatabaseFixture(t)
	response := updateSavedConnectorConfig(t, fixture, map[string]any{
		"kind": "postgres", "name": "clear on create", "database": "fixture",
		"username": "writer", "host": "127.0.0.1", "port": 1, "table": "samples",
		"write_mode": "insert", "write_interval_seconds": 5,
		"password": "discard this fixture password", "clear_password": true,
	})
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(t.Context(), record.DatabaseConnectorID)
	require.NoError(t, err)
	var config map[string]any
	require.NoError(t, json.Unmarshal([]byte(connector.ConnectionConfig), &config))
	require.NotContains(t, config, "password")
	require.NotContains(t, response.Body.String(), "discard this fixture password")
}
