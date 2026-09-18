package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// installWorkspaceSaveFailure makes the local workspace write fail, which is
// the second step of every Step 4 database setup save.
func installWorkspaceSaveFailure(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, statement := range []string{
		`CREATE TRIGGER fail_workspace_update BEFORE UPDATE ON system_settings WHEN NEW.key = 'studio_v2_workspace' BEGIN SELECT RAISE(ABORT, 'injected workspace save failure'); END`,
		`CREATE TRIGGER fail_workspace_insert BEFORE INSERT ON system_settings WHEN NEW.key = 'studio_v2_workspace' BEGIN SELECT RAISE(ABORT, 'injected workspace save failure'); END`,
	} {
		_, err := db.ExecContext(context.Background(), statement)
		require.NoError(t, err)
	}
}

func serveWorkspaceDatabaseConfig(t *testing.T, fixture workspaceDatabaseFixture, payload workspaceDatabaseConfigRequest) *httptest.ResponseRecorder {
	t.Helper()
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", payload)
	fixture.handler.UpdateConfig(c)
	return resp
}

func serveWorkspaceDatabaseTarget(t *testing.T, fixture workspaceDatabaseFixture, pointID string, payload workspaceDatabaseTargetRequest) *httptest.ResponseRecorder {
	t.Helper()
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+pointID, payload)
	c.Params = gin.Params{{Key: "point_id", Value: pointID}}
	fixture.handler.UpsertTarget(c)
	return resp
}

func workspaceSetupRevisionFrom(t *testing.T, resp *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct {
		Data struct {
			SetupRevision string `json:"setup_revision"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	return body.Data.SetupRevision
}

// currentWorkspaceSetupRevision reads the revision exactly as the browser does
// before it sends a Step 4 database mutation.
func currentWorkspaceSetupRevision(t *testing.T, fixture workspaceDatabaseFixture) string {
	t.Helper()
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-config", http.NoBody)
	fixture.handler.GetConfig(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	return workspaceSetupRevisionFrom(t, resp)
}

func assertWorkspaceSetupRevisionConflict(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	var body struct {
		Success bool `json:"success"`
		Error   struct {
			Code      string `json:"code"`
			RequestID string `json:"request_id"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body), resp.Body.String())
	require.False(t, body.Success)
	require.Equal(t, "revision_mismatch", body.Error.Code)
	require.NotEmpty(t, body.Error.RequestID)
}

func TestAtomicSetupSave_ConfigSecondStepFailureKeepsPriorSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()
	saveWorkspaceDatabaseConfig(t, fixture)
	before, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	connectorBefore, err := fixture.connectorSvc.GetByID(ctx, before.DatabaseConnectorID)
	require.NoError(t, err)
	revisionBefore := currentWorkspaceSetupRevision(t, fixture)
	installWorkspaceSaveFailure(t, fixture.mainDB)

	changed := validSQLiteConfigRequest(fixture)
	changed.Table = "sensor_values_v2"
	changed.RowGroups = []workspaceDatabaseRowGroup{{
		ID: "group-v2", TableSchema: "main", TableName: "sensor_values_v2",
		MemberPointIDs: []string{fixture.pointIDs[0]}, GroupKeyColumns: []string{"ts"},
	}}
	changed.ExpectedSetupRevision = revisionBefore
	resp := serveWorkspaceDatabaseConfig(t, fixture, changed)

	require.Equal(t, http.StatusInternalServerError, resp.Code, resp.Body.String())
	require.NotContains(t, resp.Body.String(), "injected")
	connectorAfter, err := fixture.connectorSvc.GetByID(ctx, before.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, "sensor_values", connectorConfigString(connectorAfter, "table"), "a failed second step must roll back the connector")
	require.Equal(t, connectorBefore.IdentityRevision, connectorAfter.IdentityRevision)
	after, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, before, after, "reload must return the prior consistent workspace")
	require.Equal(t, revisionBefore, currentWorkspaceSetupRevision(t, fixture))
}

func TestAtomicSetupSave_TargetReferenceFailureKeepsPriorTargets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()
	saveWorkspaceDatabaseConfig(t, fixture)
	before, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	revisionBefore := currentWorkspaceSetupRevision(t, fixture)
	installWorkspaceSaveFailure(t, fixture.mainDB)

	resp := serveWorkspaceDatabaseTarget(t, fixture, fixture.pointIDs[0], workspaceDatabaseTargetRequest{
		ColumnName: "line_a", Enabled: true, ExpectedSetupRevision: revisionBefore,
	})

	require.Equal(t, http.StatusInternalServerError, resp.Code, resp.Body.String())
	rows, err := fixture.dbMappingSvc.List(ctx, dbtarget.TargetMappingListFilter{ConnectorID: &before.DatabaseConnectorID})
	require.NoError(t, err)
	require.Empty(t, rows, "a failed workspace reference must roll back the target mapping")
	after, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, before, after)
}

func TestAtomicSetupSave_StaleSetupRevisionIsRejectedBeforeMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()

	first := serveWorkspaceDatabaseConfig(t, fixture, validSQLiteConfigRequest(fixture))
	require.Equal(t, http.StatusOK, first.Code, first.Body.String())
	staleRevision := workspaceSetupRevisionFrom(t, first)
	require.NotEmpty(t, staleRevision, "database config responses must expose the setup revision")

	update := validSQLiteConfigRequest(fixture)
	update.WriteIntervalSeconds = 10
	update.ExpectedSetupRevision = staleRevision
	second := serveWorkspaceDatabaseConfig(t, fixture, update)
	require.Equal(t, http.StatusOK, second.Code, second.Body.String())
	currentRevision := workspaceSetupRevisionFrom(t, second)
	require.NotEmpty(t, currentRevision)
	require.NotEqual(t, staleRevision, currentRevision)

	stale := validSQLiteConfigRequest(fixture)
	stale.Table = "sensor_values_v2"
	stale.ExpectedSetupRevision = staleRevision
	assertWorkspaceSetupRevisionConflict(t, serveWorkspaceDatabaseConfig(t, fixture, stale))
	missing := validSQLiteConfigRequest(fixture)
	missing.Table = "sensor_values_v2"
	assertWorkspaceSetupRevisionConflict(t, serveWorkspaceDatabaseConfig(t, fixture, missing))
	assertWorkspaceSetupRevisionConflict(t, serveWorkspaceDatabaseTarget(t, fixture, fixture.pointIDs[0], workspaceDatabaseTargetRequest{
		ColumnName: "line_a", Enabled: true, ExpectedSetupRevision: staleRevision,
	}))

	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	connector, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, "sensor_values", connectorConfigString(connector, "table"))
	require.Equal(t, 10, connector.DefaultWriteIntervalSeconds)
	rows, err := fixture.dbMappingSvc.List(ctx, dbtarget.TargetMappingListFilter{ConnectorID: &record.DatabaseConnectorID})
	require.NoError(t, err)
	require.Empty(t, rows)
	require.Equal(t, currentRevision, currentWorkspaceSetupRevision(t, fixture))

	saved := serveWorkspaceDatabaseTarget(t, fixture, fixture.pointIDs[0], workspaceDatabaseTargetRequest{
		ColumnName: "line_a", Enabled: true, ExpectedSetupRevision: currentRevision,
	})
	require.Equal(t, http.StatusOK, saved.Code, saved.Body.String())
	targetRevision := workspaceSetupRevisionFrom(t, saved)
	require.NotEmpty(t, targetRevision)
	require.NotEqual(t, currentRevision, targetRevision)
	require.Equal(t, targetRevision, currentWorkspaceSetupRevision(t, fixture))
}

func TestAtomicSetupSave_FirstSaveFailureLeavesNoConnector(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()
	installWorkspaceSaveFailure(t, fixture.mainDB)

	resp := serveWorkspaceDatabaseConfig(t, fixture, validSQLiteConfigRequest(fixture))

	require.Equal(t, http.StatusInternalServerError, resp.Code, resp.Body.String())
	connectors, err := fixture.connectorSvc.List(ctx, dbtarget.ConnectorListFilter{})
	require.NoError(t, err)
	require.Empty(t, connectors, "a failed first save must not leave an unbound connector")
	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Empty(t, record.DatabaseConnectorID)
}
