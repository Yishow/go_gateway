package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type workspaceDatabaseFixture struct {
	handler  *StudioV2WorkspaceDatabaseHandler
	ruleSvc  *sourcerule.Service
	targetDB string
	pointIDs []string
}

func TestStudioV2WorkspaceDatabaseHandler_SaveConfigAndOneTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)

	configReq := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", bytes.NewBufferString(`{
		"kind":"sqlite",
		"name":"Line A SQLite",
		"database":"`+fixture.targetDB+`",
		"schema":"main",
		"table":"sensor_values",
		"write_mode":"insert",
		"write_interval_seconds":5,
		"timestamp_column":"ts"
	}`))
	configReq.Header.Set("Content-Type", "application/json")
	configResp := httptest.NewRecorder()
	configCtx, _ := gin.CreateTestContext(configResp)
	configCtx.Request = configReq

	fixture.handler.UpdateConfig(configCtx)

	require.Equal(t, http.StatusOK, configResp.Code, configResp.Body.String())
	configBody := decodeWorkspaceDatabaseBody(t, configResp)
	configData := configBody["data"].(map[string]any)
	require.Equal(t, "sqlite", configData["kind"])
	require.Equal(t, "sensor_values", configData["table"])
	require.Equal(t, "saved", configData["save_state"])
	require.Equal(t, "not_running", configData["runtime_apply_status"])

	targetReq := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[0], bytes.NewBufferString(`{
		"column_name":"line_a",
		"enabled":true
	}`))
	targetReq.Header.Set("Content-Type", "application/json")
	targetResp := httptest.NewRecorder()
	targetCtx, _ := gin.CreateTestContext(targetResp)
	targetCtx.Request = targetReq
	targetCtx.Params = gin.Params{{Key: "point_id", Value: fixture.pointIDs[0]}}

	fixture.handler.UpsertTarget(targetCtx)

	require.Equal(t, http.StatusOK, targetResp.Code, targetResp.Body.String())
	targetBody := decodeWorkspaceDatabaseBody(t, targetResp)
	targetData := targetBody["data"].(map[string]any)
	require.Equal(t, fixture.pointIDs[0], targetData["point_id"])
	require.Equal(t, "line_a", targetData["column_name"])
	require.Equal(t, "saved", targetData["save_state"])
	require.Equal(t, "not_running", targetData["runtime_apply_status"])

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	listResp := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listResp)
	listCtx.Request = listReq

	fixture.handler.ListTargets(listCtx)

	require.Equal(t, http.StatusOK, listResp.Code, listResp.Body.String())
	listBody := decodeWorkspaceDatabaseBody(t, listResp)
	items := listBody["data"].([]any)
	require.Len(t, items, 1)
	require.Equal(t, fixture.pointIDs[0], items[0].(map[string]any)["point_id"])
}

func TestStudioV2WorkspaceDatabaseHandler_InvalidTargetDoesNotTouchSavedRows(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	saveWorkspaceDatabaseConfig(t, fixture)

	saveValidTarget(t, fixture, fixture.pointIDs[0], "line_a")

	invalidReq := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[1], bytes.NewBufferString(`{
		"column_name":"missing_column",
		"enabled":true
	}`))
	invalidReq.Header.Set("Content-Type", "application/json")
	invalidResp := httptest.NewRecorder()
	invalidCtx, _ := gin.CreateTestContext(invalidResp)
	invalidCtx.Request = invalidReq
	invalidCtx.Params = gin.Params{{Key: "point_id", Value: fixture.pointIDs[1]}}

	fixture.handler.UpsertTarget(invalidCtx)

	require.Equal(t, http.StatusBadRequest, invalidResp.Code, invalidResp.Body.String())
	require.Contains(t, invalidResp.Body.String(), "missing_column")

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/database-targets", nil)
	listResp := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listResp)
	listCtx.Request = listReq

	fixture.handler.ListTargets(listCtx)

	require.Equal(t, http.StatusOK, listResp.Code, listResp.Body.String())
	listBody := decodeWorkspaceDatabaseBody(t, listResp)
	items := listBody["data"].([]any)
	require.Len(t, items, 1)
	require.Equal(t, fixture.pointIDs[0], items[0].(map[string]any)["point_id"])
}

func decodeWorkspaceDatabaseBody(t *testing.T, resp *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	return body
}

func newWorkspaceDatabaseFixture(t *testing.T) workspaceDatabaseFixture {
	t.Helper()

	mainDB := openWorkspaceDatabaseTestDB(t)
	targetDB := createWorkspaceTargetSQLite(t)

	deviceRepo := device.NewSQLRepository(mainDB)
	pointRepo := point.NewSQLRepository(mainDB)
	tagRepo := tag.NewSQLRepository(mainDB)
	mappingRepo := mapping.NewSQLRepository(mainDB)
	sourceRuleRepo := sourcerule.NewSQLRepository(mainDB)
	workspaceRepo := workspace.NewSQLRepository(mainDB)
	dbConnectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	dbMappingRepo := dbtarget.NewSQLTargetMappingRepository(mainDB)

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	tagSvc := tag.NewService(tagRepo)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)
	ruleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, pointSvc, nil)
	workspaceSvc := workspace.NewService(workspaceRepo)
	dbConnectorSvc := dbtarget.NewConnectorService(dbConnectorRepo, dbMappingRepo)
	dbMappingSvc := dbtarget.NewMappingService(dbMappingRepo, dbConnectorRepo, tagSvc)

	ctx := context.Background()
	require.NoError(t, deviceRepo.Create(ctx, &schema.Device{
		ID:               "dev-A",
		Name:             "dev-A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
	_, err := workspaceSvc.AttachDevice(ctx, "dev-A")
	require.NoError(t, err)

	ruleRecord, err := ruleSvc.Create(ctx, sourcerule.CreateRuleRequest{
		ID:           "rule-A",
		DeviceID:     "dev-A",
		StartAddress: "40001",
		Count:        2,
		DataType:     schema.DataTypeFloat64,
		NamingPrefix: "LINE_",
		Enabled:      true,
	})
	require.NoError(t, err)

	links, err := ruleSvc.ListLinks(ctx, ruleRecord.ID)
	require.NoError(t, err)
	pointIDs := make([]string, 0, len(links))
	for index, link := range links {
		tagRecord, err := tagSvc.Create(ctx, tag.CreateTagRequest{
			Key:         "line.a.tag." + strconv.Itoa(index+1),
			DisplayName: "Line A Tag",
			DataType:    schema.DataTypeFloat64,
			Unit:        "C",
		})
		require.NoError(t, err)
		mappingRecord, err := mappingSvc.Create(ctx, mapping.CreateMappingRequest{
			PointID: link.PointID,
			TagID:   tagRecord.ID,
			Enabled: boolPtr(true),
		})
		require.NoError(t, err)
		link.TagID = &tagRecord.ID
		link.MappingID = &mappingRecord.ID
		pointIDs = append(pointIDs, link.PointID)
	}
	require.NoError(t, ruleSvc.ReplaceLinks(ctx, ruleRecord.ID, links))

	handler := NewStudioV2WorkspaceDatabaseHandler(workspaceSvc, deviceSvc, ruleSvc, dbConnectorSvc, dbMappingSvc)
	return workspaceDatabaseFixture{
		handler:  handler,
		ruleSvc:  ruleSvc,
		targetDB: targetDB,
		pointIDs: pointIDs,
	}
}

func saveWorkspaceDatabaseConfig(t *testing.T, fixture workspaceDatabaseFixture) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", bytes.NewBufferString(`{
		"kind":"sqlite",
		"name":"Line A SQLite",
		"database":"`+fixture.targetDB+`",
		"schema":"main",
		"table":"sensor_values",
		"write_mode":"insert",
		"write_interval_seconds":5,
		"timestamp_column":"ts"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.UpdateConfig(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}

func saveValidTarget(t *testing.T, fixture workspaceDatabaseFixture, pointID string, columnName string) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+pointID, bytes.NewBufferString(`{
		"column_name":"`+columnName+`",
		"enabled":true
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "point_id", Value: pointID}}

	fixture.handler.UpsertTarget(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}

func openWorkspaceDatabaseTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "main.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	migrator := datalinkbase.NewMigrator()
	require.NoError(t, migrator.Migrate(db))
	return db
}

func createWorkspaceTargetSQLite(t *testing.T) string {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE sensor_values (
			ts DATETIME PRIMARY KEY,
			line_a REAL NOT NULL,
			line_b REAL NOT NULL
		)
	`)
	require.NoError(t, err)

	return dbPath
}
