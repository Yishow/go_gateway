package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDatabaseHandler_RowGroupsRoundTripWithTargetRefs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)

	sharedGroupConfig := validSQLiteConfigRequest(fixture)
	sharedGroupConfig.RowGroups = []workspaceDatabaseRowGroup{{
		ID: "group-shared-temp", TableSchema: "main", TableName: "sensor_values",
		MemberPointIDs: fixture.pointIDs, GroupKeyColumns: []string{"ts", "line_id"},
	}}
	updateWorkspaceDatabaseConfig(t, fixture, sharedGroupConfig)

	targetReq := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[0], workspaceDatabaseTargetRequest{
		ColumnName: "temperature_c", Enabled: true, RowGroupID: "group-shared-temp",
		ExpectedSetupRevision: currentWorkspaceSetupRevision(t, fixture),
	})
	targetResp := httptest.NewRecorder()
	targetCtx, _ := gin.CreateTestContext(targetResp)
	targetCtx.Request = targetReq
	targetCtx.Params = gin.Params{{Key: "point_id", Value: fixture.pointIDs[0]}}

	fixture.handler.UpsertTarget(targetCtx)

	require.Equal(t, http.StatusOK, targetResp.Code, targetResp.Body.String())
	targetBody := decodeWorkspaceDatabaseBody(t, targetResp)
	targetData := targetBody["data"].(map[string]any)
	require.Equal(t, "group-shared-temp", targetData["row_group_id"])

	mappings, err := fixture.dbMappingSvc.List(context.Background(), dbtarget.TargetMappingListFilter{})
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	require.NotNil(t, mappings[0].GroupKey)
	require.Equal(t, rowGroupTargetKey("group-shared-temp", fixture.pointIDs[0]), *mappings[0].GroupKey)

	configReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-config", http.NoBody)
	configResp := httptest.NewRecorder()
	configCtx, _ := gin.CreateTestContext(configResp)
	configCtx.Request = configReq

	fixture.handler.GetConfig(configCtx)

	require.Equal(t, http.StatusOK, configResp.Code, configResp.Body.String())
	configBody := decodeWorkspaceDatabaseBody(t, configResp)
	configData := configBody["data"].(map[string]any)
	rowGroups := configData["row_groups"].([]any)
	require.Len(t, rowGroups, 1)
	rowGroup := rowGroups[0].(map[string]any)
	require.Equal(t, "group-shared-temp", rowGroup["id"])
	require.NotEmpty(t, rowGroup["connector_id"])
	require.Equal(t, "main", rowGroup["table_schema"])
	require.Equal(t, "sensor_values", rowGroup["table_name"])
	require.Equal(t, []any{fixture.pointIDs[0], fixture.pointIDs[1]}, rowGroup["member_point_ids"])
	require.Equal(t, []any{"ts", "line_id"}, rowGroup["group_key_columns"])

	listReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", http.NoBody)
	listResp := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listResp)
	listCtx.Request = listReq

	fixture.handler.ListTargets(listCtx)

	require.Equal(t, http.StatusOK, listResp.Code, listResp.Body.String())
	listBody := decodeWorkspaceDatabaseBody(t, listResp)
	items := listBody["data"].([]any)
	require.Len(t, items, 1)
	require.Equal(t, "group-shared-temp", items[0].(map[string]any)["row_group_id"])
}

func TestStudioV2WorkspaceDatabaseHandler_LegacyTargetsRemainCompatible(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0])

	listReq := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", http.NoBody)
	listResp := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listResp)
	listCtx.Request = listReq

	fixture.handler.ListTargets(listCtx)

	require.Equal(t, http.StatusOK, listResp.Code, listResp.Body.String())
	listBody := decodeWorkspaceDatabaseBody(t, listResp)
	items := listBody["data"].([]any)
	require.Len(t, items, 1)
	require.Nil(t, items[0].(map[string]any)["row_group_id"])
}

func TestStudioV2WorkspaceDatabaseHandler_RowGroupTargetRequiresPointMembership(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)

	singleMemberConfig := validSQLiteConfigRequest(fixture)
	singleMemberConfig.RowGroups = []workspaceDatabaseRowGroup{{
		ID: "group-shared-temp", TableSchema: "main", TableName: "sensor_values",
		MemberPointIDs: []string{fixture.pointIDs[0]}, GroupKeyColumns: []string{"ts"},
	}}
	updateWorkspaceDatabaseConfig(t, fixture, singleMemberConfig)

	targetReq := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[1], workspaceDatabaseTargetRequest{
		ColumnName: "temperature_c", Enabled: true, RowGroupID: "group-shared-temp",
	})
	targetResp := httptest.NewRecorder()
	targetCtx, _ := gin.CreateTestContext(targetResp)
	targetCtx.Request = targetReq
	targetCtx.Params = gin.Params{{Key: "point_id", Value: fixture.pointIDs[1]}}

	fixture.handler.UpsertTarget(targetCtx)

	require.Equal(t, http.StatusBadRequest, targetResp.Code, targetResp.Body.String())
	require.Contains(t, targetResp.Body.String(), "database target row group does not contain point")

	mappings, err := fixture.dbMappingSvc.List(context.Background(), dbtarget.TargetMappingListFilter{})
	require.NoError(t, err)
	require.Empty(t, mappings)
}
