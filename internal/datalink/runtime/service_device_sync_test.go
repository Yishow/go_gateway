package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

func TestService_UpsertDeviceAddsDeviceAndEnabledPoints(t *testing.T) {
	ctx := context.Background()
	runtimeSvc := newRuntimeServiceForDeviceSyncTest(t)
	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 321); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}
	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual server failed: %v", err)
	}
	defer server.Stop()
	cfgJSON, err := json.Marshal(map[string]any{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	})
	if err != nil {
		t.Fatalf("marshal config failed: %v", err)
	}

	dev := &schema.Device{
		ID:               "device-1",
		Name:             "device-1",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: string(cfgJSON),
	}

	createdPoint, err := runtimeSvc.pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID: "device-1",
		Name:     "p1",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
	})
	if err != nil {
		t.Fatalf("create point failed: %v", err)
	}

	if err := runtimeSvc.UpsertDevice(ctx, dev); err != nil {
		t.Fatalf("upsert device failed: %v", err)
	}

	if _, exists := runtimeSvc.scheduler.GetDeviceBreakerState("device-1"); !exists {
		t.Fatal("expected device breaker to be initialized after device sync")
	}
	results := runtimeSvc.scheduler.PollNow([]string{createdPoint.ID})
	if len(results) != 1 {
		t.Fatalf("expected 1 poll result after device sync, got %d", len(results))
	}
	if got := fmt.Sprint(results[0].Value); got != "321" {
		t.Fatalf("unexpected poll value after device sync: %+v", got)
	}
	if meta := runtimeSvc.lookupPointMeta(createdPoint.ID); meta.DeviceID != "device-1" {
		t.Fatalf("unexpected point meta after device sync: %+v", meta)
	}
	updatedPoint, err := runtimeSvc.pointSvc.GetByID(ctx, createdPoint.ID)
	if err != nil {
		t.Fatalf("get point after device sync failed: %v", err)
	}
	if updatedPoint.PollingGroupID == nil || *updatedPoint.PollingGroupID == "" {
		t.Fatal("expected point polling group to be auto-assigned during device sync")
	}
}

func TestService_UpsertDeviceStartsAutoCreatedPollingGroupOnRunningScheduler(t *testing.T) {
	ctx := context.Background()
	runtimeSvc := newRuntimeServiceForDeviceSyncTest(t)
	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 456); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}
	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual server failed: %v", err)
	}
	defer server.Stop()

	cfgJSON, err := json.Marshal(map[string]any{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	})
	if err != nil {
		t.Fatalf("marshal config failed: %v", err)
	}

	dev := &schema.Device{
		ID:               "device-running",
		Name:             "device-running",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: string(cfgJSON),
	}

	createdPoint, err := runtimeSvc.pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID: "device-running",
		Name:     "p-running",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
	})
	if err != nil {
		t.Fatalf("create point failed: %v", err)
	}

	if err := runtimeSvc.scheduler.Start(nil); err != nil {
		t.Fatalf("start scheduler failed: %v", err)
	}
	defer runtimeSvc.scheduler.Stop()

	if err := runtimeSvc.UpsertDevice(ctx, dev); err != nil {
		t.Fatalf("upsert device failed: %v", err)
	}

	select {
	case collected := <-runtimeSvc.scheduler.ValueChannel():
		if collected.PointID != createdPoint.ID {
			t.Fatalf("unexpected point collected: %+v", collected)
		}
	case <-time.After(1500 * time.Millisecond):
		t.Fatal("expected auto-created polling group to start ticking on running scheduler")
	}
}

func TestService_RemoveDeviceRemovesSchedulerStateAndPointMeta(t *testing.T) {
	ctx := context.Background()
	runtimeSvc := newRuntimeServiceForDeviceSyncTest(t)
	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 321); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}
	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual server failed: %v", err)
	}
	defer server.Stop()
	cfgJSON, err := json.Marshal(map[string]any{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	})
	if err != nil {
		t.Fatalf("marshal config failed: %v", err)
	}

	dev := &schema.Device{
		ID:               "device-1",
		Name:             "device-1",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: string(cfgJSON),
	}

	createdPoint, err := runtimeSvc.pointSvc.Create(ctx, point.CreatePointRequest{
		DeviceID: "device-1",
		Name:     "p1",
		Address:  "40001",
		DataType: schema.DataTypeInt16,
		Mode:     schema.PointModeReadOnly,
	})
	if err != nil {
		t.Fatalf("create point failed: %v", err)
	}
	if err := runtimeSvc.UpsertDevice(ctx, dev); err != nil {
		t.Fatalf("upsert device failed: %v", err)
	}

	runtimeSvc.RemoveDevice("device-1")

	if _, exists := runtimeSvc.scheduler.GetDeviceBreakerState("device-1"); exists {
		t.Fatal("expected device breaker to be removed from scheduler")
	}
	results := runtimeSvc.scheduler.PollNow([]string{createdPoint.ID})
	if len(results) != 0 {
		t.Fatalf("expected no poll results after device removal, got %d", len(results))
	}
	if meta := runtimeSvc.lookupPointMeta(createdPoint.ID); meta.DeviceID != "" {
		t.Fatalf("expected point meta to be cleared, got %+v", meta)
	}
}

func newRuntimeServiceForDeviceSyncTest(t *testing.T) *Service {
	t.Helper()

	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	devSvc := device.NewService(device.NewMemoryRepository(), connMgr)
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	mappingSvc := mapping.NewService(mapping.NewMemoryRepository())
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	groupSvc := pollinggroup.NewService(groupRepo)
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)

	svc, err := NewService(DefaultConfig(), Dependencies{
		Scheduler:           scheduler,
		Writer:              storage.NewMemoryStorage(16),
		DeviceService:       devSvc,
		PointService:        pointSvc,
		MappingService:      mappingSvc,
		TagService:          tagSvc,
		PollingGroupService: groupSvc,
	})
	if err != nil {
		t.Fatalf("new runtime service failed: %v", err)
	}

	return svc
}
