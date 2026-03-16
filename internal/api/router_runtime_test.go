package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"

	_ "go-gateway/internal/datalink/connector/adapters"
)

func TestNewRouter_RuntimeStatusEndpoint(t *testing.T) {
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)

	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)

	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewService(mapping.NewMemoryRepository())
	pollingGroupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())
	settingsSvc := settings.NewService(settings.NewMemoryRepository())

	ctx := t.Context()

	createdDevice, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		Name:     "Runtime Mixer",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     502,
			"slave_id": 1,
		},
	})
	if err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	createdDevice.Status = schema.DeviceStatusActive
	if err := deviceRepo.Update(ctx, createdDevice); err != nil {
		t.Fatalf("activate device failed: %v", err)
	}

	createdGroup, err := pollingGroupSvc.Create(ctx, pollinggroup.CreateRequest{
		Name:       "runtime-fast",
		IntervalMs: 1000,
		Priority:   1,
	})
	if err != nil {
		t.Fatalf("create polling group failed: %v", err)
	}

	if _, err := pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID:       createdDevice.ID,
		Name:           "Flow",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &createdGroup.ID,
	}); err != nil {
		t.Fatalf("create point failed: %v", err)
	}

	createdPoint, err := pointSvc.List(ctx, point.ListFilter{})
	if err != nil || len(createdPoint) != 1 {
		t.Fatalf("list points failed: %v", err)
	}
	createdPoint[0].PollingGroupID = &createdGroup.ID
	readAt := time.Now().Add(-2 * time.Second)
	createdPoint[0].LastReadAt = &readAt
	if err := pointRepo.Update(ctx, createdPoint[0]); err != nil {
		t.Fatalf("update point failed: %v", err)
	}

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pollingGroupSvc,
		Settings:     settingsSvc,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/runtime/status?device_id="+createdDevice.ID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Running       bool `json:"running"`
			UptimeSeconds int  `json:"uptime_seconds"`
			Collectors    []struct {
				DeviceID     string `json:"device_id"`
				DeviceName   string `json:"device_name"`
				PointsTotal  int    `json:"points_total"`
				PointsStale  int    `json:"points_stale"`
				BreakerState string `json:"breaker_state"`
			} `json:"collectors"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}

	if !body.Success {
		t.Fatalf("expected success body, got %s", resp.Body.String())
	}
	if len(body.Data.Collectors) != 1 {
		t.Fatalf("expected one collector, got %+v", body.Data.Collectors)
	}
	if body.Data.Collectors[0].DeviceID != createdDevice.ID {
		t.Fatalf("unexpected device id: %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].PointsTotal != 1 {
		t.Fatalf("expected points_total=1, got %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].PointsStale < 1 {
		t.Fatalf("expected stale point count, got %+v", body.Data.Collectors[0])
	}
}

func TestNewRouter_RuntimeStreamEndpointRequiresDeviceID(t *testing.T) {
	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewService(mapping.NewMemoryRepository())
	pollingGroupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())
	settingsSvc := settings.NewService(settings.NewMemoryRepository())

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pollingGroupSvc,
		Settings:     settingsSvc,
		Runtime:      &datalinkruntime.Service{},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/runtime/stream", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing device_id, got %d body=%s", resp.Code, resp.Body.String())
	}
}

func TestNewRouter_RuntimeStatusUsesRuntimeServiceState(t *testing.T) {
	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewService(mapping.NewMemoryRepository())
	pollingGroupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())
	settingsSvc := settings.NewService(settings.NewMemoryRepository())

	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.Config{
		Writer:   storage.NewMemoryStorage(8),
		Snapshot: datalinkruntime.Snapshot{},
	})
	if err != nil {
		t.Fatalf("new runtime service failed: %v", err)
	}

	if err := runtimeSvc.Start(context.Background()); err != nil {
		t.Fatalf("start runtime service failed: %v", err)
	}
	t.Cleanup(func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := runtimeSvc.Stop(stopCtx); err != nil {
			t.Fatalf("stop runtime service failed: %v", err)
		}
	})

	time.Sleep(1100 * time.Millisecond)

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pollingGroupSvc,
		Settings:     settingsSvc,
		Runtime:      runtimeSvc,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/runtime/status", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Running       bool  `json:"running"`
			UptimeSeconds int64 `json:"uptime_seconds"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}

	if !body.Data.Running {
		t.Fatalf("expected runtime status running=true, got %+v", body.Data)
	}
	if body.Data.UptimeSeconds < 1 {
		t.Fatalf("expected runtime uptime to be populated, got %+v", body.Data)
	}
}
