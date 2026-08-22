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
	"go-gateway/internal/datalink/audit"
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
	handler       *StudioV2WorkspaceDatabaseHandler
	ruleSvc       *sourcerule.Service
	connectorSvc  *dbtarget.ConnectorService
	dbMappingSvc  *dbtarget.MappingService
	connectorRepo dbtarget.ConnectorRepository
	workspaceSvc  *workspace.Service
	auditSvc      *audit.Service
	targetDB      string
	pointIDs      []string
}

type workspaceDatabaseConfigRequest struct {
	Kind                 string                      `json:"kind"`
	Name                 string                      `json:"name"`
	Database             string                      `json:"database,omitempty"`
	Schema               string                      `json:"schema,omitempty"`
	Table                string                      `json:"table,omitempty"`
	WriteMode            string                      `json:"write_mode,omitempty"`
	WriteIntervalSeconds int                         `json:"write_interval_seconds,omitempty"`
	TimestampColumn      string                      `json:"timestamp_column,omitempty"`
	Host                 string                      `json:"host,omitempty"`
	Port                 int                         `json:"port,omitempty"`
	Username             string                      `json:"username,omitempty"`
	Password             string                      `json:"password,omitempty"`
	RowGroups            []workspaceDatabaseRowGroup `json:"row_groups,omitempty"`
}

type workspaceDatabaseRowGroup struct {
	ID              string   `json:"id"`
	TableSchema     string   `json:"table_schema"`
	TableName       string   `json:"table_name"`
	MemberPointIDs  []string `json:"member_point_ids"`
	GroupKeyColumns []string `json:"group_key_columns"`
}

type workspaceDatabaseTargetRequest struct {
	ColumnName string `json:"column_name"`
	Enabled    bool   `json:"enabled"`
	RowGroupID string `json:"row_group_id,omitempty"`
}

type workspaceDatabaseSchemaRequest struct {
	DryRun bool `json:"dry_run"`
}

func newWorkspaceDatabaseJSONRequest(t *testing.T, method string, target string, payload any) *http.Request {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)
	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestStudioV2WorkspaceDatabaseHandler_SaveConfigAndOneTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)

	configReq := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", workspaceDatabaseConfigRequest{
		Kind:                 "sqlite",
		Name:                 "Line A SQLite",
		Database:             fixture.targetDB,
		Schema:               "main",
		Table:                "sensor_values",
		WriteMode:            "insert",
		WriteIntervalSeconds: 5,
		TimestampColumn:      "ts",
	})
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

	targetReq := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[0], workspaceDatabaseTargetRequest{
		ColumnName: "line_a",
		Enabled:    true,
	})
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

	// 空 column_name 為真正無效的請求（handler 層硬擋），不應影響既有已存列。
	// 註：欄位「不存在」在 Studio V2 已改為放行（degraded），由建表流程補建，
	// 因此這裡用空欄位名而非 missing_column 來測無效情境。
	invalidReq := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+fixture.pointIDs[1], workspaceDatabaseTargetRequest{
		Enabled: true,
	})
	invalidResp := httptest.NewRecorder()
	invalidCtx, _ := gin.CreateTestContext(invalidResp)
	invalidCtx.Request = invalidReq
	invalidCtx.Params = gin.Params{{Key: "point_id", Value: fixture.pointIDs[1]}}

	fixture.handler.UpsertTarget(invalidCtx)

	require.Equal(t, http.StatusBadRequest, invalidResp.Code, invalidResp.Body.String())
	require.Contains(t, invalidResp.Body.String(), "column_name")

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

func TestStudioV2WorkspaceDatabaseHandler_PasswordPersistedAndPreservedOnUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()

	// 首次儲存帶密碼的 postgres connector（127.0.0.1:1 會立即 refused，probe 失敗但仍會儲存）。
	updateWorkspaceDatabaseConfig(t, fixture, workspaceDatabaseConfigRequest{
		Kind: "postgres", Name: "Line A PG", Host: "127.0.0.1", Port: 1,
		Database: "gateway", Username: "gw_writer", Password: "s3cret-pw",
		Schema: "public", Table: "sensor_values", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
	})

	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, record.DatabaseConnectorID)

	connector, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.Contains(t, connector.ConnectionConfig, "s3cret-pw", "password 應被寫入 ConnectionConfig")

	// 不帶 password 再次更新（例如只改 write interval），應保留既有密碼。
	updateWorkspaceDatabaseConfig(t, fixture, workspaceDatabaseConfigRequest{
		Kind: "postgres", Name: "Line A PG", Host: "127.0.0.1", Port: 1,
		Database: "gateway", Username: "gw_writer", Schema: "public",
		Table: "sensor_values", WriteMode: "insert", WriteIntervalSeconds: 10,
		TimestampColumn: "ts",
	})

	connectorAfter, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.Contains(t, connectorAfter.ConnectionConfig, "s3cret-pw", "未帶 password 的更新應保留既有密碼")
}

func TestStudioV2WorkspaceDatabaseHandler_SwitchToSQLiteClearsPersistedPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()

	updateWorkspaceDatabaseConfig(t, fixture, workspaceDatabaseConfigRequest{
		Kind: "postgres", Name: "Line A PG", Host: "127.0.0.1", Port: 1,
		Database: "gateway", Username: "gw_writer", Password: "s3cret-pw",
		Schema: "public", Table: "sensor_values", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
	})

	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	updateWorkspaceDatabaseConfig(t, fixture, workspaceDatabaseConfigRequest{
		Kind: "sqlite", Name: "Line A SQLite", Database: fixture.targetDB,
		Schema: "main", Table: "sensor_values", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
	})

	connectorAfter, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.NotContains(t, connectorAfter.ConnectionConfig, "s3cret-pw", "切換到 sqlite 時應清掉既有密碼")
}

func TestStudioV2WorkspaceDatabaseHandler_InvalidRowGroupsDoNotPartiallyUpdateConnector(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixture(t)
	ctx := context.Background()

	updateWorkspaceDatabaseConfig(t, fixture, workspaceDatabaseConfigRequest{
		Kind: "sqlite", Name: "Line A SQLite", Database: fixture.targetDB,
		Schema: "main", Table: "sensor_values", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
	})

	record, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	connectorBefore, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, "sensor_values", connectorConfigString(connectorBefore, "table"))

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", workspaceDatabaseConfigRequest{
		Kind: "sqlite", Name: "Line A SQLite", Database: fixture.targetDB,
		Schema: "main", Table: "sensor_values_v2", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
		RowGroups: []workspaceDatabaseRowGroup{{
			ID: "group-stale-scope", TableSchema: "main", TableName: "sensor_values",
			MemberPointIDs: []string{fixture.pointIDs[0]}, GroupKeyColumns: []string{"ts"},
		}},
	})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.UpdateConfig(c)

	require.Equal(t, http.StatusBadRequest, resp.Code, resp.Body.String())
	require.Contains(t, resp.Body.String(), "database row group scope must match workspace database config")

	connectorAfter, err := fixture.connectorSvc.GetByID(ctx, record.DatabaseConnectorID)
	require.NoError(t, err)
	require.Equal(t, "sensor_values", connectorConfigString(connectorAfter, "table"), "row-group 驗證失敗時不應部分更新 connector")

	recordAfter, err := fixture.workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Empty(t, recordAfter.DatabaseRowGroups)
}

func TestStudioV2WorkspaceDatabaseHandler_GenerateSchemaCreatesTable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixture := newWorkspaceDatabaseFixtureEmptyTarget(t)
	saveWorkspaceDatabaseConfig(t, fixture)
	saveValidTarget(t, fixture, fixture.pointIDs[0], "line_a")

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/database-schema/generate", workspaceDatabaseSchemaRequest{DryRun: false})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.GenerateSchema(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	body := decodeWorkspaceDatabaseBody(t, resp)
	data := body["data"].(map[string]any)
	require.Greater(t, data["executed"].(float64), float64(0), "應有資料表結構被建立")

	// 目標 SQLite 真的應有 sensor_values 表。
	targetDB, err := sql.Open("sqlite", fixture.targetDB)
	require.NoError(t, err)
	defer targetDB.Close()
	var name string
	err = targetDB.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='sensor_values'`).Scan(&name)
	require.NoError(t, err, "sensor_values 表應已被建立")
	require.Equal(t, "sensor_values", name)
}

func updateWorkspaceDatabaseConfig(t *testing.T, fixture workspaceDatabaseFixture, payload any) {
	t.Helper()

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", payload)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.UpdateConfig(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
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

	auditSvc := audit.NewService(audit.NewMemoryRepository())
	handler := NewStudioV2WorkspaceDatabaseHandler(workspaceSvc, deviceSvc, ruleSvc, dbConnectorSvc, dbMappingSvc).WithAudit(auditSvc)
	return workspaceDatabaseFixture{
		handler:       handler,
		ruleSvc:       ruleSvc,
		connectorSvc:  dbConnectorSvc,
		dbMappingSvc:  dbMappingSvc,
		connectorRepo: dbConnectorRepo,
		workspaceSvc:  workspaceSvc,
		auditSvc:      auditSvc,
		targetDB:      targetDB,
		pointIDs:      pointIDs,
	}
}

func newWorkspaceDatabaseFixtureEmptyTarget(t *testing.T) workspaceDatabaseFixture {
	t.Helper()

	fixture := newWorkspaceDatabaseFixture(t)
	// 指向一個尚未建立 sensor_values 表的空 SQLite，讓 GenerateSchema 真正建表。
	emptyTarget := filepath.Join(t.TempDir(), "empty-target.db")
	emptyDB, err := sql.Open("sqlite", emptyTarget)
	require.NoError(t, err)
	require.NoError(t, emptyDB.Ping())
	require.NoError(t, emptyDB.Close())
	fixture.targetDB = emptyTarget
	return fixture
}

func saveWorkspaceDatabaseConfig(t *testing.T, fixture workspaceDatabaseFixture) {
	t.Helper()

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", workspaceDatabaseConfigRequest{
		Kind: "sqlite", Name: "Line A SQLite", Database: fixture.targetDB,
		Schema: "main", Table: "sensor_values", WriteMode: "insert",
		WriteIntervalSeconds: 5, TimestampColumn: "ts",
	})
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.UpdateConfig(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
}

func saveValidTarget(t *testing.T, fixture workspaceDatabaseFixture, pointID string, columnName string) {
	t.Helper()

	req := newWorkspaceDatabaseJSONRequest(t, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-targets/"+pointID, workspaceDatabaseTargetRequest{
		ColumnName: columnName,
		Enabled:    true,
	})
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
