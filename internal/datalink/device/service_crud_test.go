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
