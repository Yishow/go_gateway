package main

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"go-gateway/internal/datalink/modbusshare"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/protocol/modbus"
)

func TestProductionTargetWriter_BlocksBeforeHydration(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	created, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "pending.temperature", DataType: schema.DataTypeInt16})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	share := modbusshare.NewService(tagSvc, 4096)
	if _, err := share.UpsertMapping(ctx, created.ID, 25); err != nil {
		t.Fatalf("upsert mapping: %v", err)
	}
	err = (modbusShareTargetWriter{share: share}).WriteTagValue(ctx, created.ID, int16(1), time.Now())
	var shareErr *modbusshare.Error
	if !errors.As(err, &shareErr) || shareErr.Code != modbusshare.ErrCodeHydrationRequired {
		t.Fatalf("expected hydration gate, got %v", err)
	}
}

func TestProductionTargetWriter_BlocksWhenHydrationReadFails(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	created, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "hydration.error", DataType: schema.DataTypeInt16})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	share := modbusshare.NewService(tagSvc, 4096)
	if _, err := share.UpsertMapping(ctx, created.ID, 25); err != nil {
		t.Fatalf("upsert mapping: %v", err)
	}
	share.SetHydrationError(errors.New("hydration read failed"))
	err = (modbusShareTargetWriter{share: share}).WriteTagValue(ctx, created.ID, int16(1), time.Now())
	var deliveryErr *datalinkruntime.DeliveryError
	if !errors.As(err, &deliveryErr) || deliveryErr.Stage != "hydration" {
		t.Fatalf("expected hydration delivery error, got %v", err)
	}
}

func TestProductionTargetWriter_SkipsDisabledDurableShareWithoutNoise(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	created, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "disabled.temperature", DataType: schema.DataTypeInt16})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	share := modbusshare.NewService(tagSvc, 4096)
	if _, err := share.UpsertMapping(ctx, created.ID, 25); err != nil {
		t.Fatalf("upsert mapping: %v", err)
	}
	target := newProductionTargetWriter(nil, share).(fanoutTargetWriter)
	if got := target.WriteTagValueOutcomes(ctx, created.ID, int16(1), time.Now()); len(got) != 0 {
		t.Fatalf("expected disabled Share target to be skipped, got %d outcomes", len(got))
	}
}

func TestProductionTargetWriter_PublishesCollectedValueToModbusMemory(t *testing.T) {
	ctx := context.Background()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	created, err := tagSvc.Create(ctx, tag.CreateTagRequest{Key: "production.temperature", DataType: schema.DataTypeInt16})
	if err != nil {
		t.Fatalf("create tag: %v", err)
	}
	share := modbusshare.NewService(tagSvc, 4096)
	share.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	if _, err := share.UpsertMapping(ctx, created.ID, 25); err != nil {
		t.Fatalf("upsert mapping: %v", err)
	}
	port := freeTCPPort(t)
	if err := share.ApplySettings(ctx, modbusshare.Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              port,
		SlaveID:           1,
		CapacityRegisters: 2048,
	}); err != nil {
		t.Fatalf("enable share: %v", err)
	}
	t.Cleanup(func() { _ = share.CloseRuntime() })

	target := newProductionTargetWriter(nil, share)
	if err := target.WriteTagValue(ctx, created.ID, int16(4321), time.Now()); err != nil {
		t.Fatalf("publish collected value: %v", err)
	}
	client := modbus.CreateTCPClient("127.0.0.1", port, 1, time.Second)
	if err := client.Connect(); err != nil {
		t.Fatalf("connect Modbus client: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	words, err := client.ReadHoldingRegisters(25, 1)
	if err != nil {
		t.Fatalf("read mapped register: %v", err)
	}
	if len(words) != 1 || words[0] != 4321 {
		t.Fatalf("expected mapped register 4321, got %v", words)
	}
}

type recordingTargetWriter struct {
	err   error
	calls int
}

func (w *recordingTargetWriter) WriteTagValue(context.Context, string, any, time.Time) error {
	w.calls++
	return w.err
}

func TestProductionTargetWriter_ReportsEachTargetAfterOneFails(t *testing.T) {
	database := &recordingTargetWriter{err: errors.New("database unavailable")}
	share := &recordingTargetWriter{}
	writer := fanoutTargetWriter{targets: []datalinkruntime.TargetWriter{database, share}}

	outcomes := writer.WriteTagValueOutcomes(context.Background(), "tag-1", int16(2), time.Now())
	if database.calls != 1 || share.calls != 1 {
		t.Fatalf("expected both target attempts, database=%d share=%d", database.calls, share.calls)
	}
	if len(outcomes) != 2 || outcomes[0].Target != datalinkruntime.TargetDatabase || outcomes[0].Err == nil || outcomes[1].Target != datalinkruntime.TargetModbusShare || outcomes[1].Err != nil {
		t.Fatalf("unexpected per-target outcomes: %+v", outcomes)
	}
}

func freeTCPPort(t *testing.T) int {
	t.Helper()
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release port: %v", err)
	}
	return port
}
