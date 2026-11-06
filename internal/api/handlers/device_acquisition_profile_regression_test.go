package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const existingAcquisitionProfile = `{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1","limits":{"random_points":180}}}`

func acquisitionProfileHandlerRig(t *testing.T) (*device.Service, *workspace.Service, *stubDeviceRuntimeSyncer) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	repo := device.NewMemoryRepository()
	devSvc := device.NewService(repo, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	now := time.Now()
	require.NoError(t, repo.Create(t.Context(), &schema.Device{
		ID: "profile-device", Name: "profile-device", Protocol: schema.ProtocolMC3E,
		Status: schema.DeviceStatusActive, ConnectionConfig: existingAcquisitionProfile,
		CreatedAt: now, UpdatedAt: now,
	}))
	_, err := workspaceSvc.AttachDevice(t.Context(), "profile-device")
	require.NoError(t, err)
	return devSvc, workspaceSvc, &stubDeviceRuntimeSyncer{}
}

func runAcquisitionProfileRequest(t *testing.T, method, body string, handle func(*gin.Context)) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequestWithContext(t.Context(), method, "/devices/profile-device", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Params = gin.Params{{Key: "id", Value: "profile-device"}}
	c.Request = request
	handle(c)
	return response
}

func TestDeviceHandler_AcquisitionProfileChangeResyncsRuntime(t *testing.T) {
	devSvc, _, syncer := acquisitionProfileHandlerRig(t)
	handler := NewDeviceHandler(devSvc, syncer)
	body := `{"connection_config":{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1","limits":{"random_points":120}}}}`
	response := runAcquisitionProfileRequest(t, http.MethodPut, body, handler.Update)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.NotNil(t, syncer.upserted)
	require.Contains(t, syncer.upserted.ConnectionConfig, `"random_points":120`)
}

func TestDeviceHandler_TransportChangeDoesNotApplyProfileWithoutRestart(t *testing.T) {
	devSvc, _, syncer := acquisitionProfileHandlerRig(t)
	handler := NewDeviceHandler(devSvc, syncer)
	body := `{"connection_config":{"host":"127.0.0.2","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1","limits":{"random_points":120}}}}`
	response := runAcquisitionProfileRequest(t, http.MethodPut, body, handler.Update)
	require.Equal(t, http.StatusOK, response.Code, response.Body.String())
	require.Nil(t, syncer.upserted)
}

func TestAcquisitionProfileLimitFieldErrorOnDeviceUpdate(t *testing.T) {
	for _, route := range []string{"device", "workspace"} {
		t.Run(route, func(t *testing.T) {
			devSvc, workspaceSvc, syncer := acquisitionProfileHandlerRig(t)
			handle := NewDeviceHandler(devSvc, syncer).Update
			if route == "workspace" {
				handle = NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, syncer).Update
			}
			body := `{"connection_config":{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1","limits":{"random_points":300}}}}`
			response := runAcquisitionProfileRequest(t, http.MethodPut, body, handle)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), `"field":"connection_config.acquisition.limits.random_points"`)
			stored, err := devSvc.GetByID(t.Context(), "profile-device")
			require.NoError(t, err)
			require.JSONEq(t, existingAcquisitionProfile, stored.ConnectionConfig)
			require.Nil(t, syncer.upserted)
		})
	}
}

func TestAcquisitionProfileLimitFieldErrorOnDeviceCreate(t *testing.T) {
	for _, route := range []string{"device", "workspace"} {
		t.Run(route, func(t *testing.T) {
			devSvc, workspaceSvc, syncer := acquisitionProfileHandlerRig(t)
			handle := NewDeviceHandler(devSvc, syncer).Create
			if route == "workspace" {
				handle = NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, syncer).Create
			}
			body := `{"name":"new-device","protocol":"mc_3e","connection_config":{"host":"127.0.0.1","port":5000,"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1","limits":{"random_points":300}}}}`
			response := runAcquisitionProfileRequest(t, http.MethodPost, body, handle)
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			require.Contains(t, response.Body.String(), `"field":"connection_config.acquisition.limits.random_points"`)
			devices, err := devSvc.List(t.Context(), device.ListFilter{})
			require.NoError(t, err)
			require.Len(t, devices, 1)
		})
	}
}
