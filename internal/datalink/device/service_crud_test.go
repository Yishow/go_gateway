package device

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
)

func TestService_UpdateKeepsDormantDeviceNotRunning(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	now := time.Now()
	if err := repo.Create(context.Background(), &schema.Device{
		ID:               "dev-draft",
		Name:             "Line A PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	name := "Line A Saved"
	updated, err := svc.Update(context.Background(), "dev-draft", UpdateDeviceRequest{
		Name: &name,
		ConnectionConfig: map[string]interface{}{
			"host":     "192.168.10.20",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if err != nil {
		t.Fatalf("update device failed: %v", err)
	}

	if updated.Status != schema.DeviceStatusDraft {
		t.Fatalf("expected draft device to remain dormant, got %s", updated.Status)
	}
}

func TestService_UpdatePersistsProtocolSwitchAndClearsStaleProbe(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSQLRepository(db)
	svc := NewService(repo, nil)
	now := time.Now()
	probeSucceeded := true
	if err := repo.Create(context.Background(), &schema.Device{
		ID:               "dev-mc",
		Name:             "Line A PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}`,
		LastTestAt:       &now,
		LastTestSuccess:  &probeSucceeded,
		LastTestError:    "stale probe result",
		ReadinessStatus:  "ready",
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	protocol := schema.ProtocolMC3E
	updated, err := svc.Update(context.Background(), "dev-mc", UpdateDeviceRequest{
		Protocol: &protocol,
		ConnectionConfig: map[string]interface{}{
			"host":       "192.168.10.10",
			"port":       6000,
			"station_no": 0,
			"network_no": 0,
			"pc_no":      255,
			"io_no":      1023,
			"timeout":    5,
		},
	})
	if err != nil {
		t.Fatalf("update device failed: %v", err)
	}

	if updated.Protocol != schema.ProtocolMC3E {
		t.Fatalf("expected MC 3E protocol, got %s", updated.Protocol)
	}
	persisted, err := repo.GetByID(context.Background(), "dev-mc")
	if err != nil {
		t.Fatalf("get updated device failed: %v", err)
	}
	if persisted.LastTestAt != nil || persisted.LastTestSuccess != nil || persisted.LastTestError != "" || persisted.ReadinessStatus != "" {
		t.Fatalf("expected stale probe state to be cleared, got %+v", persisted)
	}
}
