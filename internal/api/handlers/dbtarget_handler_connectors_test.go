package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/dbtarget"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupDBTargetConnectorRouter(t *testing.T) (*gin.Engine, *dbtarget.ConnectorService) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	mainDB := openMigratedDBTargetMainDB(t)
	connectorRepo := dbtarget.NewSQLConnectorRepository(mainDB)
	mappingRepo := dbtarget.NewSQLTargetMappingRepository(mainDB)
	connectorSvc := dbtarget.NewConnectorService(connectorRepo, mappingRepo)
	mappingSvc := dbtarget.NewMappingService(mappingRepo, connectorRepo, nil)

	handler := NewDatabaseTargetHandler(connectorSvc, mappingSvc)
	router := gin.Default()
	router.GET("/datalink/db-targets/connectors", handler.ListConnectors)
	router.POST("/datalink/db-targets/connectors", handler.CreateConnector)
	router.PUT("/datalink/db-targets/connectors/:id", handler.UpdateConnector)
	router.DELETE("/datalink/db-targets/connectors/:id", handler.DeleteConnector)
	router.POST("/datalink/db-targets/connectors/:id/test", handler.TestConnector)

	return router, connectorSvc
}

func TestDatabaseTargetHandler_ConnectorCrudAndTestRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	router, connectorSvc := setupDBTargetConnectorRouter(t)
	targetPath := filepath.Join(t.TempDir(), "connector-handler-target.db")

	createBody, err := json.Marshal(map[string]any{
		"name":  "handler-sqlite",
		"kind":  "sqlite",
		"enabled": true,
		"default_write_interval_seconds": 5,
		"connection_config": map[string]any{
			"path":  targetPath,
			"table": "sensor_values",
		},
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/datalink/db-targets/connectors", bytes.NewBuffer(createBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusCreated, resp.Code)

	var createPayload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &createPayload))
	assert.Equal(t, true, createPayload["success"])
	connectorID := createPayload["data"].(map[string]any)["id"].(string)

	req, err = http.NewRequest(http.MethodGet, "/datalink/db-targets/connectors", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var listPayload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &listPayload))
	listData := listPayload["data"].([]any)
	require.Len(t, listData, 1)

	updateBody, err := json.Marshal(map[string]any{
		"connection_config": map[string]any{
			"path":  targetPath,
			"table": "sensor_values_v2",
		},
	})
	require.NoError(t, err)

	req, err = http.NewRequest(http.MethodPut, "/datalink/db-targets/connectors/"+connectorID, bytes.NewBuffer(updateBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	saved, err := connectorSvc.GetByID(ctx, connectorID)
	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Contains(t, saved.ConnectionConfig, "sensor_values_v2")

	req, err = http.NewRequest(http.MethodPost, "/datalink/db-targets/connectors/"+connectorID+"/test", nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var testPayload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &testPayload))
	assert.Equal(t, true, testPayload["success"])
	testData := testPayload["data"].(map[string]any)
	assert.Equal(t, "ready", testData["status"])

	req, err = http.NewRequest(http.MethodDelete, "/datalink/db-targets/connectors/"+connectorID, nil)
	require.NoError(t, err)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	_, err = connectorSvc.GetByID(ctx, connectorID)
	require.Error(t, err)
}
