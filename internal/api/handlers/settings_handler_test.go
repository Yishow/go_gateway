package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingSettingsRepository struct{}

func (failingSettingsRepository) Get(context.Context, string) (*settings.SettingItem, error) {
	return nil, errors.New("database secret should not be returned")
}
func (failingSettingsRepository) Set(context.Context, string, interface{}) error {
	return errors.New("database secret should not be returned")
}
func (failingSettingsRepository) List(context.Context) ([]*settings.SettingItem, error) {
	return nil, errors.New("database secret should not be returned")
}

func setupSettingsRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := settings.NewMemoryRepository()
	// Seed 預設設定
	_ = repo.Set(context.TODO(), "test_key", "test_value")

	svc := settings.NewService(repo)
	h := NewSettingsHandler(svc)

	r.GET("/datalink/settings", h.List)
	r.PUT("/datalink/settings/:key", h.Update)

	return r
}

func TestSettingsHandler_List(t *testing.T) {
	r := setupSettingsRouter()

	req, _ := http.NewRequestWithContext(t.Context(), "GET", "/datalink/settings", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestSettingsHandler_Update(t *testing.T) {
	r := setupSettingsRouter()

	reqBody := map[string]interface{}{
		"value": "updated_value",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), "PUT", "/datalink/settings/test_key", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "updated_value", data["value"])
}

func TestSettingsHandler_UnknownRepositoryErrorUsesSafeTypedEnvelope(t *testing.T) {
	h := NewSettingsHandler(settings.NewService(failingSettingsRepository{}))
	r := gin.New()
	r.GET("/settings", h.List)
	r.PUT("/settings/:key", h.Update)

	request := newHandlerTestRequest(http.MethodGet, "/settings", nil)
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)
	require.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), ErrCodeSettingsUnavailable)
	assert.NotContains(t, response.Body.String(), "database secret")

	body := bytes.NewBufferString(`{"value":"x"}`)
	request = newHandlerTestRequest(http.MethodPut, "/settings/example", body)
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	r.ServeHTTP(response, request)
	require.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Body.String(), ErrCodeSettingsUpdateFailed)
	assert.NotContains(t, response.Body.String(), "database secret")
}

func TestSettingsHandler_ModbusShareUpdatePersistsAndAppliesLifecycle(t *testing.T) {
	repo := settings.NewMemoryRepository()
	settingsSvc := settings.NewService(repo)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	shareSvc := modbusshare.NewService(tagSvc, 64)
	shareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	require.NoError(t, shareSvc.SetSettingsRepository(repo))
	t.Cleanup(func() { _ = shareSvc.Stop() })

	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())

	h := NewSettingsHandler(settingsSvc).WithModbusShare(shareSvc)
	r := gin.New()
	r.PUT("/datalink/settings/:key", h.Update)
	body, err := json.Marshal(map[string]any{"value": map[string]any{
		"enabled": true, "bind_address": "127.0.0.1", "port": port,
		"slave_id": 9, "capacity_registers": 8, "settings_revision": "handler-revision",
	}})
	require.NoError(t, err)
	req := newHandlerTestRequest(http.MethodPut, "/datalink/settings/modbus_share", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.True(t, shareSvc.Status().Enabled)
	assert.Equal(t, port, shareSvc.Status().Port)
	assert.Equal(t, uint8(9), shareSvc.Status().SlaveID)
}

func TestSettingsHandler_ModbusShareUpdateHydrationPendingSavesWithoutBinding(t *testing.T) {
	repo := settings.NewMemoryRepository()
	settingsSvc := settings.NewService(repo)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	shareSvc := modbusshare.NewService(tagSvc, 64)
	shareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStatePending, Readiness: false})
	require.NoError(t, shareSvc.SetSettingsRepository(repo))
	t.Cleanup(func() { _ = shareSvc.CloseRuntime() })

	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	require.NoError(t, ln.Close())

	h := NewSettingsHandler(settingsSvc).WithModbusShare(shareSvc)
	r := gin.New()
	r.PUT("/datalink/settings/:key", h.Update)
	body, err := json.Marshal(map[string]any{"value": map[string]any{
		"enabled": true, "bind_address": "127.0.0.1", "port": port,
		"slave_id": 1, "capacity_registers": 8,
	}})
	require.NoError(t, err)
	req := newHandlerTestRequest(http.MethodPut, "/datalink/settings/modbus_share", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.False(t, shareSvc.Status().Enabled)
	persisted, err := repo.Get(context.Background(), settings.KeyModbusShare)
	require.NoError(t, err)
	assert.True(t, persisted.Value.(modbusshare.Settings).Enabled)
}

func TestSettingsHandler_ModbusShareInvalidCapacityReturnsUnprocessableEntity(t *testing.T) {
	repo := settings.NewMemoryRepository()
	settingsSvc := settings.NewService(repo)
	shareSvc := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 64)
	shareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	require.NoError(t, shareSvc.SetSettingsRepository(repo))
	t.Cleanup(func() { _ = shareSvc.CloseRuntime() })

	h := NewSettingsHandler(settingsSvc).WithModbusShare(shareSvc)
	r := gin.New()
	r.PUT("/datalink/settings/:key", h.Update)
	body := bytes.NewBufferString(`{"value":{"enabled":true,"bind_address":"127.0.0.1","port":15020,"slave_id":1,"capacity_registers":9999}}`)
	req := newHandlerTestRequest(http.MethodPut, "/datalink/settings/modbus_share", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code, w.Body.String())
	assert.Contains(t, w.Body.String(), modbusshare.ErrCodeCapacityExceeded)
	_, err := repo.Get(context.Background(), settings.KeyModbusShare)
	assert.Error(t, err, "invalid settings must not be durably committed")
}
