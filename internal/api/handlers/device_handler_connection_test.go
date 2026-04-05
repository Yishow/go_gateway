package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceHandler_TestConnection_ReturnsPlanningHints(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))

	server := virtualmodbus.NewServer(bank)
	require.NoError(t, server.Start(0))
	defer server.Stop()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	h := NewDeviceHandler(svc)
	r.POST("/datalink/devices", h.Create)
	r.POST("/datalink/devices/:id/test", h.TestConnection)

	createReq := device.CreateDeviceRequest{
		Name:     "planning-hints",
		Protocol: "modbus_tcp",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     server.Port(),
			"slave_id": 1,
			"timeout":  2,
		},
	}

	body, err := json.Marshal(createReq)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/datalink/devices", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	deviceID := createResp["data"].(map[string]interface{})["id"].(string)

	req, err = http.NewRequest("POST", "/datalink/devices/"+deviceID+"/test", nil)
	require.NoError(t, err)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	require.Contains(t, data, "planning_hints")
	hints := data["planning_hints"].(map[string]interface{})
	assert.Equal(t, true, hints["source_rule_planning_supported"])
	assert.Equal(t, true, hints["probe_supported"])
	assert.Empty(t, hints["planning_blocked_reason"])
}
