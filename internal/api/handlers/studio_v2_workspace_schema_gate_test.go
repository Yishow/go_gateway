package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func workspaceTargetTableExists(t *testing.T, path string) bool {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	var count int
	require.NoError(t, db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sensor_values'`).Scan(&count))
	return count == 1
}

// workspaceSchemaStatements reads the read-only plan that stays available.
func workspaceSchemaStatements(t *testing.T, fixture workspaceDatabaseFixture) []string {
	t.Helper()
	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/database-schema/generate", workspaceDatabaseSchemaRequest{DryRun: true})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	fixture.handler.GenerateSchema(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	var payload struct {
		Data struct {
			Statements []string `json:"statements"`
			Executed   int      `json:"executed"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	require.Equal(t, 0, payload.Data.Executed, "a read-only plan must execute nothing")
	return payload.Data.Statements
}

// Activation must report what is missing instead of quietly creating tables.
func TestEnsureWorkspaceSchema_ReportsPreparationInsteadOfCreating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixtureEmptyTarget(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0])

	err := fixture.handler.EnsureWorkspaceSchema(context.Background())

	require.True(t, errors.Is(err, errWorkspaceSchemaPreparationRequired), "missing tables must be reported, got %v", err)
	require.False(t, workspaceTargetTableExists(t, fixture.targetDB), "activation must not create the target table")

	statements := workspaceSchemaStatements(t, fixture)
	require.NotEmpty(t, statements, "the read-only plan must still report what is missing")
	targetDB, err := sql.Open("sqlite", fixture.targetDB)
	require.NoError(t, err)
	defer func() { _ = targetDB.Close() }()
	for _, statement := range statements {
		_, execErr := targetDB.ExecContext(context.Background(), statement)
		require.NoError(t, execErr)
	}

	require.NoError(t, fixture.handler.EnsureWorkspaceSchema(context.Background()), "a prepared schema must pass the check")
}

// A Local Modbus-only workspace has no database target, so nothing is required.
func TestEnsureWorkspaceSchema_LocalModbusOnlyNeedsNoDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixtureEmptyTarget(t)

	require.NoError(t, fixture.handler.EnsureWorkspaceSchema(context.Background()))
	require.False(t, workspaceTargetTableExists(t, fixture.targetDB))
}

// Readiness follows the live projection: a disabled mapping's missing table
// must not block activation once nothing enabled is left to write.
func TestEnsureWorkspaceSchema_IgnoresDisabledMappingTables(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixtureEmptyTarget(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0])

	// Prepare the table so the workspace starts out ready.
	targetDB, err := sql.Open("sqlite", fixture.targetDB)
	require.NoError(t, err)
	for _, statement := range workspaceSchemaStatements(t, fixture) {
		_, execErr := targetDB.ExecContext(context.Background(), statement)
		require.NoError(t, execErr)
	}
	require.NoError(t, fixture.handler.EnsureWorkspaceSchema(context.Background()))

	// Disable the only target and remove its table behind the scenes.
	disableWorkspaceTargetMapping(t, fixture)
	_, err = targetDB.ExecContext(context.Background(), `DROP TABLE sensor_values`)
	require.NoError(t, err)
	require.NoError(t, targetDB.Close())

	require.NoError(t, fixture.handler.EnsureWorkspaceSchema(context.Background()),
		"a disabled mapping's missing table must not block activation")
}

// disableWorkspaceTargetMapping flips every target mapping to disabled with a
// direct store update, simulating an operator who disabled a mapping whose
// column type no longer passes live validation.
func disableWorkspaceTargetMapping(t *testing.T, fixture workspaceDatabaseFixture) {
	t.Helper()

	result, err := fixture.mainDB.ExecContext(context.Background(),
		`UPDATE database_target_mappings SET enabled = 0`)
	require.NoError(t, err)
	rows, err := result.RowsAffected()
	require.NoError(t, err)
	require.Positive(t, rows)
}

func TestStudioV2WorkspaceActivation_ReportsSchemaPreparationWithoutActivating(t *testing.T) {
	gin.SetMode(gin.TestMode)
	activator := &stubWorkspaceActivator{response: &workspace.ActivationResponse{WorkspaceID: "workspace-1"}}
	handler := NewStudioV2WorkspaceActivationHandler(activator, &stubWorkspaceSchemaEnsurer{err: errWorkspaceSchemaPreparationRequired})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newValidActivationRequest(t)

	handler.Activate(c)

	require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), workspaceSchemaPreparationRequiredCode)
	require.False(t, activator.called, "activation must not run while the schema is unprepared")
}

// The workspace schema endpoint keeps its read-only plan and refuses to create
// without a confirmed preview.
func TestStudioV2WorkspaceDatabaseHandler_GenerateSchemaNeedsConfirmation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixtureEmptyTarget(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0])

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/database-schema/generate", workspaceDatabaseSchemaRequest{DryRun: false})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.GenerateSchema(c)

	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), schemaConfirmationRequiredCode)
	require.False(t, workspaceTargetTableExists(t, fixture.targetDB), "an unconfirmed request must not create tables")
	require.NotEmpty(t, workspaceSchemaStatements(t, fixture), "the read-only plan must stay available")
}
