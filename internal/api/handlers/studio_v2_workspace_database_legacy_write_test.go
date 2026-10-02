package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type legacyRowGroupGuardStub struct {
	preflightErr  error
	checkErr      error
	preflightCall int
	checkCall     int
}

func (s *legacyRowGroupGuardStub) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return &workspace.WriteGroupListResult{}, nil
}

func (s *legacyRowGroupGuardStub) PreflightLegacyRowGroupReplacement(context.Context, string, []workspace.DatabaseRowGroup) error {
	s.preflightCall++
	return s.preflightErr
}

func (s *legacyRowGroupGuardStub) CheckLegacyRowGroupReplacementInTx(context.Context, *sql.Tx, *workspace.Record, string, []workspace.DatabaseRowGroup) error {
	s.checkCall++
	return s.checkErr
}

func TestRenderLegacyWriteConflictUsesSafeActionableEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)

	renderStudioV2WorkspaceDatabaseError(c, fmt.Errorf("dsn=postgres://secret: %w", workspace.ErrWriteGroupLegacyWriteConflict))

	require.Equal(t, http.StatusConflict, resp.Code)
	var body struct {
		Success bool `json:"success"`
		Error   struct {
			Code      string `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
			Action    string `json:"action"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.False(t, body.Success)
	require.Equal(t, legacyWriteConflictCode, body.Error.Code)
	require.Equal(t, "open_write_groups", body.Error.Action)
	require.NotEmpty(t, body.Error.RequestID)
	require.NotContains(t, resp.Body.String(), "postgres://secret")
}

func TestStudioV2WorkspaceDatabaseHandler_RowGroupConflictStopsBeforeConnectorPreparation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	guard := &legacyRowGroupGuardStub{preflightErr: workspace.ErrWriteGroupLegacyWriteConflict}
	fixture.handler.writeGroups = guard

	payload := validSQLiteConfigRequest(fixture)
	payload.RowGroups = []workspaceDatabaseRowGroup{{
		ID:             "owned-row-group",
		TableSchema:    "main",
		TableName:      "sensor_values",
		MemberPointIDs: []string{fixture.pointIDs[0]},
	}}
	resp := serveWorkspaceDatabaseConfig(t, fixture, payload)

	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	require.Equal(t, 1, guard.preflightCall)
	require.Zero(t, guard.checkCall)
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Empty(t, record.DatabaseConnectorID)
	require.Empty(t, record.DatabaseRowGroups)
	require.Contains(t, resp.Body.String(), legacyWriteConflictCode)
	require.Contains(t, resp.Body.String(), "open_write_groups")
}

func TestStudioV2WorkspaceDatabaseHandler_ExplicitEmptyRowGroupsStillPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	guard := &legacyRowGroupGuardStub{preflightErr: workspace.ErrWriteGroupLegacyWriteConflict}
	fixture.handler.writeGroups = guard

	payload := validSQLiteConfigRequest(fixture)
	payload.RowGroups = []workspaceDatabaseRowGroup{}
	encoded, err := json.Marshal(payload)
	require.NoError(t, err)
	var rawPayload map[string]any
	require.NoError(t, json.Unmarshal(encoded, &rawPayload))
	// The shared fixture type uses omitempty, so put the empty array back into
	// the wire payload to distinguish an explicit replacement from nil.
	rawPayload["row_groups"] = []any{}
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", rawPayload)
	fixture.handler.UpdateConfig(c)

	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	require.Equal(t, 1, guard.preflightCall)
	require.Zero(t, guard.checkCall)
}

func TestStudioV2WorkspaceDatabaseHandler_RowGroupConflictRechecksInsideTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fixture := newWorkspaceDatabaseFixture(t)
	guard := &legacyRowGroupGuardStub{checkErr: errors.Join(workspace.ErrWriteGroupLegacyWriteConflict, errors.New("credential=secret"))}
	fixture.handler.writeGroups = guard

	payload := validSQLiteConfigRequest(fixture)
	payload.RowGroups = []workspaceDatabaseRowGroup{{
		ID:             "owned-row-group",
		TableSchema:    "main",
		TableName:      "sensor_values",
		MemberPointIDs: []string{fixture.pointIDs[0]},
	}}
	resp := serveWorkspaceDatabaseConfig(t, fixture, payload)

	require.Equal(t, http.StatusConflict, resp.Code, resp.Body.String())
	require.Equal(t, 1, guard.preflightCall)
	require.Equal(t, 1, guard.checkCall)
	require.NotContains(t, resp.Body.String(), "credential=secret")
	record, err := fixture.workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)
	require.Empty(t, record.DatabaseConnectorID)
	require.Empty(t, record.DatabaseRowGroups)
}
