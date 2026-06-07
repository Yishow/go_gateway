package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDatabaseHandler_UpdateConfigKeepsAppliedStatusForWarnings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mainDB := openWorkspaceDatabaseTestDB(t)
	defer mainDB.Close()

	deviceRepo := device.NewSQLRepository(mainDB)
	pointRepo := point.NewSQLRepository(mainDB)
	sourceRuleRepo := sourcerule.NewSQLRepository(mainDB)
	workspaceRepo := workspace.NewSQLRepository(mainDB)
	dbConnectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	dbMappingRepo := dbtarget.NewSQLTargetMappingRepository(mainDB)
	tagRepo := tag.NewSQLRepository(mainDB)

	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	ruleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, pointSvc, nil)
	tagSvc := tag.NewService(tagRepo)
	workspaceSvc := workspace.NewService(workspaceRepo)
	dbConnectorSvc := dbtarget.NewConnectorService(dbConnectorRepo, dbMappingRepo)
	dbMappingSvc := dbtarget.NewMappingService(dbMappingRepo, dbConnectorRepo, tagSvc)
	workspaceSvc.WithReadinessServices(deviceSvc, ruleSvc, dbConnectorSvc, dbMappingSvc)
	handler := NewStudioV2WorkspaceDatabaseHandler(workspaceSvc, deviceSvc, ruleSvc, dbConnectorSvc, dbMappingSvc)

	lastTestSuccess := true
	require.NoError(t, deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-A",
		Name:             "dev-A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		LastTestSuccess:  &lastTestSuccess,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}))
	_, err := workspaceSvc.AttachDevice(context.Background(), "dev-A")
	require.NoError(t, err)

	configReq := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/database-config", bytes.NewBufferString(`{
		"kind":"postgres",
		"name":"Line A PG",
		"host":"127.0.0.1",
		"port":1,
		"database":"gateway",
		"username":"gw_writer",
		"password":"s3cret-pw",
		"schema":"public",
		"table":"sensor_values",
		"write_mode":"insert",
		"write_interval_seconds":5,
		"timestamp_column":"ts"
	}`))
	configReq.Header.Set("Content-Type", "application/json")
	configResp := httptest.NewRecorder()
	configCtx, _ := gin.CreateTestContext(configResp)
	configCtx.Request = configReq

	handler.UpdateConfig(configCtx)

	require.Equal(t, http.StatusOK, configResp.Code, configResp.Body.String())
	require.Contains(t, configResp.Body.String(), `"runtime_apply_status":"applied"`)
	require.Contains(t, configResp.Body.String(), "database-connector-unreachable")
}
