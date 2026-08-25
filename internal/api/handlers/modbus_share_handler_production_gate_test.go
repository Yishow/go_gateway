package handlers

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModbusShareHandler_ProductionGateRejectsDirectProjectionMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	defer svc.Stop()
	port := reservePortForProductionGate(t)
	require.NoError(t, svc.ApplySettings(context.Background(), modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 128, SettingsRevision: "settings-1"}))
	svc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: "ws-1", WorkspaceRevision: "workspace-1", SettingsRevision: "settings-1", Readiness: true})
	svc.SetOwnershipChecker(func(context.Context, string, string) bool { return true })
	handler := NewModbusShareHandler(svc).WithGate(svc)
	router := gin.New()
	router.PUT("/mappings/:tagId", handler.UpsertMapping)
	router.POST("/write-tag-value", handler.WriteTagValue)
	router.POST("/sync", handler.SyncFromMappings)
	req := newHandlerTestRequest(http.MethodPut, "/mappings/tag-1", bytes.NewBufferString(`{"register":0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Contains(t, w.Body.String(), modbusshare.ErrCodeProjectionRequired)
	assert.Empty(t, svc.ListMappings())

	for _, path := range []string{"/write-tag-value", "/sync"} {
		req := newHandlerTestRequest(http.MethodPost, path, bytes.NewBufferString(`{"tag_id":"tag-1","value":1}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusUnprocessableEntity, w.Code, path)
		assert.Contains(t, w.Body.String(), modbusshare.ErrCodeProjectionRequired, path)
	}
}

func TestModbusShareHandler_ProductionStartUsesPersistedPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := modbusshare.NewService(tag.NewService(tag.NewMemoryRepository()), 65536)
	defer svc.Stop()
	port := reservePortForProductionGate(t)
	require.NoError(t, svc.ApplySettings(context.Background(), modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 128, SettingsRevision: "settings-1"}))
	_ = svc.Stop()
	require.NoError(t, svc.ApplySettings(context.Background(), modbusshare.Settings{Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 128, SettingsRevision: "settings-2"}))
	svc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	handler := NewModbusShareHandler(svc).WithGate(svc)
	router := gin.New()
	router.POST("/start", handler.Start)
	req := newHandlerTestRequest(http.MethodPost, "/start", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, port, svc.Status().Port)
}

func reservePortForProductionGate(t *testing.T) int {
	t.Helper()
	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}
