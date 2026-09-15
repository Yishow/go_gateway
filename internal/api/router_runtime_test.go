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

	otherDevice, err := deviceSvc.Create(ctx, device.CreateDeviceRequest{
		Name:     "Runtime Filler",
		Protocol: schema.ProtocolModbusTCP,
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.2",
			"port":     502,
			"slave_id": 1,
		},
	})
	if err != nil {
		t.Fatalf("create other device failed: %v", err)
	}

	otherDevice.Status = schema.DeviceStatusActive
	if err := deviceRepo.Update(ctx, otherDevice); err != nil {
		t.Fatalf("activate other device failed: %v", err)
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

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/datalink/runtime/status?device_id="+createdDevice.ID, http.NoBody)
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
			Metrics       *struct {
				CollectedTotal       uint64 `json:"collected_total"`
				WriteSuccessTotal    uint64 `json:"write_success_total"`
				WriteErrorTotal      uint64 `json:"write_error_total"`
				MappingErrorTotal    uint64 `json:"mapping_error_total"`
				PointStateErrorTotal uint64 `json:"point_state_error_total"`
			} `json:"metrics"`
			Collectors []struct {
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
	if body.Data.Metrics == nil {
		t.Fatalf("expected top-level metrics payload, got %s", resp.Body.String())
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

func TestNewRouter_RuntimeStatusIncludesUnavailableCollectorState(t *testing.T) {
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	pollingGroupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())

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
	createdDevice.ReadinessStatus = `{"availability_status":"unavailable","availability_reason":"device form is invalid"}`
	if err := deviceRepo.Update(ctx, createdDevice); err != nil {
		t.Fatalf("update device failed: %v", err)
	}

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		PollingGroup: pollingGroupSvc,
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
	})

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/datalink/runtime/status?device_id="+createdDevice.ID, http.NoBody)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			Collectors []struct {
				AvailabilityStatus string  `json:"availability_status"`
				AvailabilityReason *string `json:"availability_reason"`
				Running            bool    `json:"running"`
			} `json:"collectors"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}

	if len(body.Data.Collectors) != 1 {
		t.Fatalf("expected one collector, got %+v", body.Data.Collectors)
	}
	if body.Data.Collectors[0].AvailabilityStatus != "unavailable" {
		t.Fatalf("expected availability_status=unavailable, got %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].AvailabilityReason == nil || *body.Data.Collectors[0].AvailabilityReason != "device form is invalid" {
		t.Fatalf("expected availability_reason to be populated, got %+v", body.Data.Collectors[0])
	}
	if body.Data.Collectors[0].Running {
		t.Fatalf("expected running=false for unavailable collector, got %+v", body.Data.Collectors[0])
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

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/runtime/stream", http.NoBody)
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

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/runtime/status", http.NoBody)
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

func TestNewRouter_RuntimeStatusMarksSelectedDeviceEmptySnapshot(t *testing.T) {
	deviceSvc := device.NewService(device.NewMemoryRepository(), nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewService(mapping.NewMemoryRepository())
	pollingGroupSvc := pollinggroup.NewService(pollinggroup.NewMemoryRepository())
	settingsSvc := settings.NewService(settings.NewMemoryRepository())

	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.Config{
		Writer: storage.NewMemoryStorage(8),
		Snapshot: datalinkruntime.Snapshot{
			Devices: []*schema.Device{
				{
					ID:       "dev-A",
					Name:     "Mixer PLC",
					Protocol: schema.ProtocolModbusTCP,
					Status:   schema.DeviceStatusActive,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("new runtime service failed: %v", err)
	}

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pollingGroupSvc,
		Settings:     settingsSvc,
		Runtime:      runtimeSvc,
	})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/runtime/status?device_id=dev-A", http.NoBody)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			SnapshotState struct {
				State string `json:"state"`
				Empty bool   `json:"empty"`
			} `json:"snapshot_state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	if !body.Success {
		t.Fatalf("expected success body, got %s", resp.Body.String())
	}
	if body.Data.SnapshotState.State != "empty" || !body.Data.SnapshotState.Empty {
		t.Fatalf("expected explicit empty snapshot state, got %s", resp.Body.String())
	}
}
