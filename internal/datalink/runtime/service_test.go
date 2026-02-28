package runtime_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

func TestService_StartConsumeAndStop(t *testing.T) {
	bank := memory.NewMemoryBank(4096)
	if err := bank.WriteWord(0, 321); err != nil {
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

	groupID := "group-1"
	memWriter := storage.NewMemoryStorage(1000)

	svc, err := runtime.NewService(runtime.Config{
		Writer: memWriter,
		Snapshot: runtime.Snapshot{
			Devices: []*schema.Device{
				{
					ID:               "dev-1",
					Name:             "dev-1",
					Protocol:         schema.ProtocolModbusTCP,
					Status:           schema.DeviceStatusActive,
					ConnectionConfig: string(cfgJSON),
				},
			},
			Points: []*schema.Point{
				{
					ID:             "point-1",
					DeviceID:       "dev-1",
					Name:           "p1",
					Address:        "0",
					Function:       "03",
					DataType:       schema.DataTypeUint16,
					Mode:           schema.PointModeReadOnly,
					PollingGroupID: &groupID,
					Enabled:        true,
				},
			},
			PollingGroups: []*schema.PollingGroup{
				{ID: groupID, Name: "g1", IntervalMs: 100, Priority: 1, Enabled: true},
			},
			Mappings: []*schema.Mapping{
				{
					ID:                "map-1",
					PointID:           "point-1",
					TagID:             "tag-1",
					TransformPipeline: "[]",
					Enabled:           true,
				},
			},
			Tags: []*schema.Tag{
				{ID: "tag-1", Key: "test.tag.1", DataType: schema.DataTypeUint16, Status: schema.TagStatusActive},
			},
		},
	})
	if err != nil {
		t.Fatalf("new runtime service failed: %v", err)
	}

	ctx := context.Background()
	if err := svc.Start(ctx); err != nil {
		t.Fatalf("start runtime failed: %v", err)
	}

	time.Sleep(350 * time.Millisecond)

	stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := svc.Stop(stopCtx); err != nil {
		t.Fatalf("stop runtime failed: %v", err)
	}

	records, err := memWriter.Query(context.Background(), storage.TimeSeriesQuery{
		TagID: "tag-1",
		Limit: 5,
		Order: "desc",
	})
	if err != nil {
		t.Fatalf("query memory writer failed: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("expected at least one timeseries record")
	}
	if records[0].ValueNum == nil || *records[0].ValueNum != 321 {
		t.Fatalf("unexpected latest value: %+v", records[0].ValueNum)
	}

	metrics := svc.Metrics()
	if metrics.TotalReads == 0 || metrics.TotalWrites == 0 {
		t.Fatalf("expected non-zero runtime metrics, got %+v", metrics)
	}
}

