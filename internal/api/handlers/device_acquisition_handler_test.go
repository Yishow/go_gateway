package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupDeviceAcquisitionHandlerRig(t *testing.T) (devices *device.Service, points *point.Service, ws *workspace.Service, runtimeSyncer *stubDeviceRuntimeSyncer) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	devRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	devSvc := device.NewService(devRepo, nil)
	devSvc.SetAcquisitionIntervalCoordinator(pointSvc)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	syncer := &stubDeviceRuntimeSyncer{}

	now := time.Now()
	for _, id := range []string{"dev-a", "dev-b"} {
		require.NoError(t, devRepo.Create(ctx, &schema.Device{
			ID:               id,
			Name:             id,
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusActive,
			ConnectionConfig: `{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}`,
			CreatedAt:        now,
			UpdatedAt:        now,
		}))
		_, err := workspaceSvc.AttachDevice(ctx, id)
		require.NoError(t, err)
		group, err := pointSvc.SetDeviceAcquisitionInterval(ctx, id, 1000)
		require.NoError(t, err)
		_, err = pointSvc.Create(ctx, point.CreatePointRequest{
			DeviceID:       id,
			Name:           "P-" + id,
			Address:        "40001",
			DataType:       schema.DataTypeUint16,
			PollingGroupID: common.Ptr(group.ID),
		})
		require.NoError(t, err)
	}

	return devSvc, pointSvc, workspaceSvc, syncer
}

func TestDeviceHandler_Update_UpdatesAcquisitionSettingsAndRejectsInvalidValues(t *testing.T) {
	devSvc, pointSvc, _, syncer := setupDeviceAcquisitionHandlerRig(t)
	handler := NewDeviceHandler(devSvc, syncer)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/devices/dev-a", strings.NewReader(`{
		"acquisition_interval_ms": 500,
		"acquisition_auto_interval": true,
		"stale_threshold_ms": 3000,
		"frozen_threshold_ms": 15000
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
	c.Request = req

	handler.Update(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	require.Equal(t, float64(500), data["acquisition_interval_ms"])
	require.Equal(t, true, data["acquisition_auto_interval"])
	require.Equal(t, float64(3000), data["stale_threshold_ms"])
	require.Equal(t, float64(15000), data["frozen_threshold_ms"])
	require.NotNil(t, syncer.upserted)
	require.Equal(t, "dev-a", syncer.upserted.ID)

	groupA, err := pointSvc.GetPollingGroupByOwnerDeviceID(t.Context(), "dev-a")
	require.NoError(t, err)
	require.Equal(t, 500, groupA.IntervalMs)

	groupB, err := pointSvc.GetPollingGroupByOwnerDeviceID(t.Context(), "dev-b")
	require.NoError(t, err)
	require.Equal(t, 1000, groupB.IntervalMs)

	for _, payload := range []string{
		`{"acquisition_interval_ms": 99}`,
		`{"acquisition_interval_ms": 3600001}`,
		`{"stale_threshold_ms": 0}`,
		`{"frozen_threshold_ms": -1}`,
	} {
		badReq := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/devices/dev-a", strings.NewReader(payload))
		badReq.Header.Set("Content-Type", "application/json")
		badResp := httptest.NewRecorder()
		badCtx, _ := gin.CreateTestContext(badResp)
		badCtx.Params = gin.Params{{Key: "id", Value: "dev-a"}}
		badCtx.Request = badReq

		handler.Update(badCtx)
		require.Equal(t, http.StatusBadRequest, badResp.Code, badResp.Body.String())
		require.Contains(t, badResp.Body.String(), `"code":"validation"`)
	}
}

func TestStudioV2WorkspaceDevicesHandler_Update_UpdatesAcquisitionSettingsAndRejectsInvalidValues(t *testing.T) {
	devSvc, pointSvc, workspaceSvc, syncer := setupDeviceAcquisitionHandlerRig(t)
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, syncer)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-a", strings.NewReader(`{
		"acquisition_interval_ms": 500,
		"acquisition_auto_interval": true,
		"stale_threshold_ms": 2000,
		"frozen_threshold_ms": 5000
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
	c.Request = req

	handler.Update(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())

	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	require.Equal(t, float64(500), data["acquisition_interval_ms"])
	require.Equal(t, true, data["acquisition_auto_interval"])
	require.Equal(t, float64(2000), data["stale_threshold_ms"])
	require.Equal(t, float64(5000), data["frozen_threshold_ms"])
	require.Equal(t, "applied", data["runtime_apply_status"])

	groupA, err := pointSvc.GetPollingGroupByOwnerDeviceID(t.Context(), "dev-a")
	require.NoError(t, err)
	require.Equal(t, 500, groupA.IntervalMs)

	groupB, err := pointSvc.GetPollingGroupByOwnerDeviceID(t.Context(), "dev-b")
	require.NoError(t, err)
	require.Equal(t, 1000, groupB.IntervalMs)

	for _, payload := range []string{
		`{"acquisition_interval_ms": 99}`,
		`{"acquisition_interval_ms": 3600001}`,
		`{"stale_threshold_ms": 0}`,
		`{"frozen_threshold_ms": -1}`,
	} {
		badReq := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-a", strings.NewReader(payload))
		badReq.Header.Set("Content-Type", "application/json")
		badResp := httptest.NewRecorder()
		badCtx, _ := gin.CreateTestContext(badResp)
		badCtx.Params = gin.Params{{Key: "id", Value: "dev-a"}}
		badCtx.Request = badReq

		handler.Update(badCtx)
		require.Equal(t, http.StatusBadRequest, badResp.Code, badResp.Body.String())
		require.Contains(t, badResp.Body.String(), `"code":"validation"`)
	}
}

func TestSameTransportConnectionConfigIgnoresAcquisitionDeclaration(t *testing.T) {
	base := `{"host":"10.0.0.1","port":5000}`
	withProfile := `{"acquisition":{"profile":"safe-v1","version":1},"host":"10.0.0.1","port":5000}`
	require.True(t, sameTransportConnectionConfig(base, base))
	require.True(t, sameTransportConnectionConfig(base, withProfile))
	require.False(t, sameTransportConnectionConfig(base, `{"host":"10.0.0.2","port":5000}`))
	require.False(t, sameTransportConnectionConfig(base, `{"host":"10.0.0.1","port":5000,"timeout":3}`))
	require.False(t, sameTransportConnectionConfig(base, `not-json`))
}

func TestStudioV2WorkspaceDevicesHandler_UpdateAcquisitionOnlyAppliesWithoutRestart(t *testing.T) {
	devSvc, _, workspaceSvc, syncer := setupDeviceAcquisitionHandlerRig(t)
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, syncer)

	// 表單每次都會帶上未改動的 connection_config；只改採集設定不得要求重啟。
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-a", strings.NewReader(`{
		"connection_config": {"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5},
		"acquisition_interval_ms": 700
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
	c.Request = req

	handler.Update(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	var body map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	require.Equal(t, "applied", data["runtime_apply_status"])
	require.NotNil(t, syncer.upserted)
}

func TestRuntimeCollectorResponseExposesConfiguredAndEffectiveInterval(t *testing.T) {
	resp := mapRuntimeCollectorResponse(datalinkruntime.DeviceRuntimeStatus{DeviceID: "dev-a", ConfiguredIntervalMs: 1000, EffectiveIntervalMs: 1500})
	raw, err := json.Marshal(resp)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"configured_interval_ms":1000`)
	require.Contains(t, string(raw), `"effective_interval_ms":1500`)
}

func TestDeviceHandler_UpdateAcquisitionWithUnchangedConnectionResyncsRuntime(t *testing.T) {
	devSvc, _, _, syncer := setupDeviceAcquisitionHandlerRig(t)
	handler := NewDeviceHandler(devSvc, syncer)

	// 舊版 API 客戶端連同未改動的 connection_config 一起送出時，仍必須同步 runtime。
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/datalink/devices/dev-a", strings.NewReader(`{
		"connection_config": {"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5},
		"acquisition_interval_ms": 500
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
	c.Request = req

	handler.Update(c)
	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	require.NotNil(t, syncer.upserted)
}
