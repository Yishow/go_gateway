package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"

	"github.com/gin-gonic/gin"
)

type stubDeviceRuntimeSyncer struct {
	upserted  *schema.Device
	upsertErr error
	removed   string
}

func (s *stubDeviceRuntimeSyncer) UpsertDevice(_ context.Context, device *schema.Device) error {
	s.upserted = device
	return s.upsertErr
}

func (s *stubDeviceRuntimeSyncer) RemoveDevice(deviceID string) {
	s.removed = deviceID
}

func TestDeviceHandler_ActivateSyncsRuntimeDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 123); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}

	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual server failed: %v", err)
	}
	defer server.Stop()

	cfgBytes, err := json.Marshal(map[string]any{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	})
	if err != nil {
		t.Fatalf("marshal config failed: %v", err)
	}

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	now := time.Now()
	dev := &schema.Device{
		ID:               "dev-1",
		Name:             "dev-1",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(cfgBytes),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := repo.Create(context.Background(), dev); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	syncer := &stubDeviceRuntimeSyncer{}
	handler := NewDeviceHandler(svc, syncer)
	req := httptest.NewRequest(http.MethodPost, "/devices/dev-1/activate", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: "dev-1"}}
	c.Request = req

	handler.Activate(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if syncer.upserted == nil || syncer.upserted.ID != "dev-1" {
		t.Fatalf("expected runtime upsert for dev-1, got %+v", syncer.upserted)
	}
}

func TestDeviceHandler_DisableRemovesRuntimeDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := device.NewMemoryRepository()
	svc := device.NewService(repo, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	now := time.Now()
	dev := &schema.Device{
		ID:               "dev-1",
		Name:             "dev-1",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := repo.Create(context.Background(), dev); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	syncer := &stubDeviceRuntimeSyncer{}
	handler := NewDeviceHandler(svc, syncer)
	req := httptest.NewRequest(http.MethodPost, "/devices/dev-1/disable", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: "dev-1"}}
	c.Request = req

	handler.Disable(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if syncer.removed != "dev-1" {
		t.Fatalf("expected runtime remove for dev-1, got %s", syncer.removed)
	}
}

func TestDeviceHandler_ActivateRollsBackWhenRuntimeSyncFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 123); err != nil {
		t.Fatalf("write memory failed: %v", err)
	}

	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start virtual server failed: %v", err)
	}
	defer server.Stop()

	cfgBytes, err := json.Marshal(map[string]any{
		"host":     "127.0.0.1",
		"port":     server.Port(),
		"slave_id": 1,
		"timeout":  2,
	})
	if err != nil {
		t.Fatalf("marshal config failed: %v", err)
	}

	repo := device.NewMemoryRepository()
	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	svc := device.NewService(repo, connMgr)
	now := time.Now()
	dev := &schema.Device{
		ID:               "dev-rollback",
		Name:             "dev-rollback",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: string(cfgBytes),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := repo.Create(context.Background(), dev); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	syncer := &stubDeviceRuntimeSyncer{upsertErr: errors.New("runtime sync failed")}
	handler := NewDeviceHandler(svc, syncer)
	req := httptest.NewRequest(http.MethodPost, "/devices/dev-rollback/activate", nil)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Params = gin.Params{{Key: "id", Value: "dev-rollback"}}
	c.Request = req

	handler.Activate(c)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	updated, err := svc.GetByID(context.Background(), "dev-rollback")
	if err != nil {
		t.Fatalf("get device failed: %v", err)
	}
	if updated.Status != schema.DeviceStatusDisabled {
		t.Fatalf("expected device to be rolled back to disabled, got %s", updated.Status)
	}
	if syncer.removed != "dev-rollback" {
		t.Fatalf("expected runtime remove for rollback, got %s", syncer.removed)
	}
}
