package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

type apiFailingTargetWriter struct {
	err   error
	calls atomic.Int64
}

func (w *apiFailingTargetWriter) WriteTagValue(context.Context, string, any, time.Time) error {
	w.calls.Add(1)
	return w.err
}

func TestNewRouter_RuntimeStatusIncludesDatabaseDeliveryFailureStages(t *testing.T) {
	bank := memory.NewMemoryBank(4096)
	if err := bank.WriteWord(0, 21); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}

	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual modbus tcp failed: %v", err)
	}
	defer server.Stop()

	cfgJSON, err := json.Marshal(schema.ConnectionConfigModbusTCP{
		Host:    "127.0.0.1",
		Port:    server.Port(),
		SlaveID: 1,
		Timeout: 2,
	})
	if err != nil {
		t.Fatalf("marshal device config failed: %v", err)
	}

	ctx := context.Background()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, connMgr)
	pointRepo := point.NewMemoryRepository()
	groupRepo := pollinggroup.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, groupRepo)
	tagRepo := tag.NewMemoryRepository()
	tagSvc := tag.NewService(tagRepo)
	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewService(mappingRepo)
	pollingGroupSvc := pollinggroup.NewService(groupRepo)
	settingsSvc := settings.NewService(settings.NewMemoryRepository())

	deviceRecord := &schema.Device{
		ID:               "dev-A",
		Name:             "Runtime DB delivery device",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: string(cfgJSON),
	}
	if err := deviceRepo.Create(ctx, deviceRecord); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	groupRecord := &schema.PollingGroup{ID: "group-fast", Name: "fast", IntervalMs: 100, Priority: 1, Enabled: true}
	if err := groupRepo.Create(ctx, groupRecord); err != nil {
		t.Fatalf("create polling group failed: %v", err)
	}

	pointRecord := &schema.Point{
		ID:             "pt-A",
		DeviceID:       deviceRecord.ID,
		Name:           "Pressure",
		Address:        "0",
		Function:       "03",
		DataType:       schema.DataTypeUint16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupRecord.ID,
		Enabled:        true,
	}
	if err := pointRepo.Create(ctx, pointRecord); err != nil {
		t.Fatalf("create point failed: %v", err)
	}

	tagRecord := &schema.Tag{
		ID:          "tag-A",
		Key:         "runtime.db.delivery",
		KeyLower:    "runtime.db.delivery",
		DisplayName: "Runtime DB Delivery",
		DataType:    schema.DataTypeUint16,
		Status:      schema.TagStatusActive,
	}
	if err := tagRepo.Create(ctx, tagRecord); err != nil {
		t.Fatalf("create tag failed: %v", err)
	}

	mappingRecord := &schema.Mapping{
		ID:                "map-A",
		PointID:           pointRecord.ID,
		TagID:             tagRecord.ID,
		TransformPipeline: "[]",
		Enabled:           true,
	}
	if err := mappingRepo.Create(ctx, mappingRecord); err != nil {
		t.Fatalf("create mapping failed: %v", err)
	}

	targetWriter := &apiFailingTargetWriter{err: errors.New("permission denied")}
	runtimeSvc, err := datalinkruntime.NewService(datalinkruntime.DefaultConfig(), datalinkruntime.Dependencies{
		Scheduler:           collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr),
		Writer:              storage.NewMemoryStorage(16),
		TargetWriter:        targetWriter,
		DeviceService:       deviceSvc,
		PointService:        pointSvc,
		MappingService:      mappingSvc,
		TagService:          tagSvc,
		PollingGroupService: pollingGroupSvc,
	})
	if err != nil {
		t.Fatalf("new runtime service failed: %v", err)
	}

	if err := runtimeSvc.Start(ctx); err != nil {
		t.Fatalf("start runtime failed: %v", err)
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := runtimeSvc.Stop(stopCtx); err != nil {
			t.Fatalf("stop runtime failed: %v", err)
		}
	}()

	deadline := time.Now().Add(2 * time.Second)
	for targetWriter.calls.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(25 * time.Millisecond)
	}
	if targetWriter.calls.Load() == 0 {
		t.Fatal("expected database target writer to be called")
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/runtime/status?device_id="+deviceRecord.ID, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			DatabaseDelivery []struct {
				DeviceID    string   `json:"device_id"`
				PointID     string   `json:"point_id"`
				TagID       string   `json:"tag_id"`
				Status      string   `json:"status"`
				Stages      []string `json:"stages"`
				FailedStage string   `json:"failed_stage"`
				Error       string   `json:"error"`
			} `json:"database_delivery"`
			Diagnostics []struct {
				Scope                 string `json:"scope"`
				DeviceID              string `json:"device_id"`
				PointID               string `json:"point_id"`
				TagID                 string `json:"tag_id"`
				LatestSuccessfulStage string `json:"latest_successful_stage"`
				FailureStage          string `json:"failure_stage"`
				FailureReason         string `json:"failure_reason"`
				LastFailureAt         string `json:"last_failure_at"`
				Stages                []struct {
					Stage  string `json:"stage"`
					Status string `json:"status"`
					Reason string `json:"reason"`
				} `json:"stages"`
			} `json:"diagnostics"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}
	if !body.Success {
		t.Fatalf("expected success response, got %s", resp.Body.String())
	}
	if len(body.Data.DatabaseDelivery) != 1 {
		t.Fatalf("expected one database delivery diagnostic, got %s", resp.Body.String())
	}

	diagnostic := body.Data.DatabaseDelivery[0]
	if diagnostic.DeviceID != deviceRecord.ID || diagnostic.PointID != pointRecord.ID || diagnostic.TagID != tagRecord.ID {
		t.Fatalf("unexpected database delivery scope: %+v", diagnostic)
	}
	if diagnostic.Status != "failed" || diagnostic.FailedStage != "db_write" || diagnostic.Error != "permission denied" {
		t.Fatalf("unexpected database delivery failure: %+v", diagnostic)
	}
	expectedStages := []string{"collected", "mapped", "db_write_failed"}
	if len(diagnostic.Stages) != len(expectedStages) {
		t.Fatalf("unexpected database delivery stages: %+v", diagnostic.Stages)
	}
	for idx := range expectedStages {
		if diagnostic.Stages[idx] != expectedStages[idx] {
			t.Fatalf("unexpected database delivery stages: %+v", diagnostic.Stages)
		}
	}

	if len(body.Data.Diagnostics) != 1 {
		t.Fatalf("expected one runtime flow diagnostic, got %s", resp.Body.String())
	}
	flowDiagnostic := body.Data.Diagnostics[0]
	if flowDiagnostic.Scope != "device:"+deviceRecord.ID ||
		flowDiagnostic.DeviceID != deviceRecord.ID ||
		flowDiagnostic.PointID != pointRecord.ID ||
		flowDiagnostic.TagID != tagRecord.ID {
		t.Fatalf("unexpected flow diagnostic scope: %+v", flowDiagnostic)
	}
	if flowDiagnostic.LatestSuccessfulStage != "runtime_projection" ||
		flowDiagnostic.FailureStage != "database_delivery" ||
		flowDiagnostic.FailureReason != "permission denied" ||
		flowDiagnostic.LastFailureAt == "" {
		t.Fatalf("unexpected flow diagnostic failure context: %+v", flowDiagnostic)
	}
	if len(flowDiagnostic.Stages) != 4 {
		t.Fatalf("unexpected flow diagnostic stages: %+v", flowDiagnostic.Stages)
	}
	if flowDiagnostic.Stages[3].Stage != "database_delivery" ||
		flowDiagnostic.Stages[3].Status != "failed" ||
		flowDiagnostic.Stages[3].Reason != "permission denied" {
		t.Fatalf("unexpected flow diagnostic failed stage: %+v", flowDiagnostic.Stages)
	}
}
