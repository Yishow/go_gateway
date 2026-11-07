package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Use a disposable Modbus peer instead of assuming a device exists on port 502.
func setupActivationFixture(t *testing.T) (*gin.Engine, *device.MemoryRepository, string) {
	t.Helper()
	bank := memory.NewMemoryBank(64)
	require.NoError(t, bank.WriteWord(0, 123))
	peer := virtualmodbus.NewServer(bank)
	require.NoError(t, peer.Start(0))
	t.Cleanup(func() { require.NoError(t, peer.Stop()) })
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	repo := device.NewMemoryRepository()
	handler := NewDeviceHandler(device.NewService(repo, manager))
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/datalink/devices", handler.Create)
	router.POST("/datalink/devices/:id/activate", handler.Activate)
	router.POST("/datalink/devices/:id/disable", handler.Disable)
	payload, err := json.Marshal(device.CreateDeviceRequest{
		Name: t.Name(), Protocol: "modbus_tcp",
		ConnectionConfig: map[string]any{"host": "127.0.0.1", "port": peer.Port(), "slave_id": 1, "timeout": 2},
	})
	require.NoError(t, err)
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/datalink/devices", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	require.NotEmpty(t, result.Data.ID)
	return router, repo, result.Data.ID
}

func activationRequest(t *testing.T, router *gin.Engine, id, action string) {
	t.Helper()
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/datalink/devices/"+id+"/"+action, http.NoBody)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	var result struct {
		Success bool `json:"success"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
	require.True(t, result.Success)
}

func TestDeviceHandler_Activate(t *testing.T) {
	router, repo, id := setupActivationFixture(t)
	activationRequest(t, router, id, "activate")
	stored, err := repo.GetByID(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, schema.DeviceStatusActive, stored.Status)
	require.NotNil(t, stored.LastTestSuccess)
	require.True(t, *stored.LastTestSuccess)
}

func TestDeviceHandler_ActivateThenDisable(t *testing.T) {
	router, repo, id := setupActivationFixture(t)
	activationRequest(t, router, id, "activate")
	active, err := repo.GetByID(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, schema.DeviceStatusActive, active.Status)
	activationRequest(t, router, id, "disable")
	disabled, err := repo.GetByID(t.Context(), id)
	require.NoError(t, err)
	require.Equal(t, schema.DeviceStatusDisabled, disabled.Status)
}
