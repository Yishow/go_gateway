package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStudioV2WorkspaceDevicesHandler_CreateCleansUpDeviceWhenAttachFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceRepo := &studioV2AttachFailureWorkspaceRepository{
		MemoryRepository: workspace.NewMemoryRepository(),
		attachErr:        errors.New("attach diagnostic: workspace write refused"),
	}
	workspaceSvc := workspace.NewService(workspaceRepo)
	_, err := workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)

	deviceRepo := &studioV2CleanupDeviceRepository{MemoryRepository: device.NewMemoryRepository()}
	deviceSvc := device.NewService(deviceRepo, nil)
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc)
	resp := invokeStudioV2WorkspaceDeviceCreate(t, handler, "dev-attach-failure")

	require.Equal(t, 500, resp.Code, resp.Body.String())
	require.False(t, studioV2ResponseSuccess(t, resp))
	require.Equal(t, 2, workspaceRepo.saveCalls, "workspace attach must reach its failing save")
	require.Equal(t, 1, deviceRepo.deleteCalls, "created device must be cleaned up after attach failure")
	_, err = deviceRepo.GetByID(t.Context(), "dev-attach-failure")
	require.Error(t, err, "successful cleanup must remove the created device")
	require.NotContains(t, resp.Body.String(), "attach diagnostic")
}

func TestStudioV2WorkspaceDevicesHandler_CreateReportsCleanupFailureWithoutSavedDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceRepo := &studioV2AttachFailureWorkspaceRepository{
		MemoryRepository: workspace.NewMemoryRepository(),
		attachErr:        errors.New("attach diagnostic: workspace write refused"),
	}
	workspaceSvc := workspace.NewService(workspaceRepo)
	_, err := workspaceSvc.GetOrCreate(t.Context())
	require.NoError(t, err)

	deviceRepo := &studioV2CleanupDeviceRepository{
		MemoryRepository: device.NewMemoryRepository(),
		deleteErr:        errors.New("cleanup diagnostic: dsn=postgres://secret"),
	}
	deviceSvc := device.NewService(deviceRepo, nil)
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc)
	resp := invokeStudioV2WorkspaceDeviceCreate(t, handler, "dev-cleanup-failure")

	require.Equal(t, 500, resp.Code, resp.Body.String())
	require.False(t, studioV2ResponseSuccess(t, resp))
	require.Equal(t, 2, workspaceRepo.saveCalls, "workspace attach must reach its failing save")
	require.Equal(t, 1, deviceRepo.deleteCalls, "cleanup must still be attempted when it can fail")
	savedDevice, err := deviceRepo.GetByID(t.Context(), "dev-cleanup-failure")
	require.NoError(t, err, "failed cleanup must leave the created device observable")
	require.Equal(t, "dev-cleanup-failure", savedDevice.ID)
	require.NotContains(t, resp.Body.String(), "attach diagnostic")
	require.NotContains(t, resp.Body.String(), "cleanup diagnostic")
	require.NotContains(t, resp.Body.String(), "postgres://secret")
	require.NotContains(t, resp.Body.String(), "dev-cleanup-failure")
}

type studioV2AttachFailureWorkspaceRepository struct {
	*workspace.MemoryRepository
	attachErr error
	saveCalls int
}

func (r *studioV2AttachFailureWorkspaceRepository) Save(ctx context.Context, record *workspace.Record) error {
	r.saveCalls++
	if len(record.OrderedDeviceIDs) > 0 {
		return r.attachErr
	}
	return r.MemoryRepository.Save(ctx, record)
}

type studioV2CleanupDeviceRepository struct {
	*device.MemoryRepository
	deleteErr   error
	deleteCalls int
}

func (r *studioV2CleanupDeviceRepository) Delete(ctx context.Context, id string) error {
	r.deleteCalls++
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.MemoryRepository.Delete(ctx, id)
}

func invokeStudioV2WorkspaceDeviceCreate(t *testing.T, handler *StudioV2WorkspaceDevicesHandler, id string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", strings.NewReader(`{
		"id": "`+id+`",
		"name": "PLC",
		"protocol": "modbus_tcp",
		"connection_config": {
			"host": "127.0.0.1",
			"port": 502,
			"slave_id": 1,
			"timeout": 5
		}
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	handler.Create(c)
	return resp
}

func studioV2ResponseSuccess(t *testing.T, resp *httptest.ResponseRecorder) bool {
	t.Helper()
	var body struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.Empty(t, body.Data, "failure response must not include saved device data")
	return body.Success
}
