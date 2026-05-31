package device

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	_ "modernc.org/sqlite"
)

func TestService_SetAvailabilityMarksDeviceUnavailableUntilExplicitlyCleared(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)

	lastSuccess := true
	now := time.Now()
	if err := repo.Create(context.Background(), &schema.Device{
		ID:               "dev-availability",
		Name:             "dev-availability",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2}`,
		LastTestSuccess:  &lastSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	unavailableDevice, err := svc.SetAvailability(context.Background(), "dev-availability", AvailabilityStatusUnavailable, "device form is invalid")
	if err != nil {
		t.Fatalf("set unavailable failed: %v", err)
	}
	if status, reason := AvailabilityOf(unavailableDevice); status != AvailabilityStatusUnavailable || reason != "device form is invalid" {
		t.Fatalf("expected unavailable device availability, got status=%s reason=%s", status, reason)
	}

	readiness, err := svc.CheckReadiness(context.Background(), "dev-availability")
	if err != nil {
		t.Fatalf("check readiness failed: %v", err)
	}
	if readiness.AvailabilityStatus != AvailabilityStatusUnavailable {
		t.Fatalf("expected readiness to preserve unavailable status, got %s", readiness.AvailabilityStatus)
	}
	if readiness.AvailabilityReason != "device form is invalid" {
		t.Fatalf("expected readiness to preserve unavailable reason, got %s", readiness.AvailabilityReason)
	}

	availableDevice, err := svc.SetAvailability(context.Background(), "dev-availability", AvailabilityStatusAvailable, "")
	if err != nil {
		t.Fatalf("clear availability failed: %v", err)
	}
	if status, reason := AvailabilityOf(availableDevice); status != AvailabilityStatusAvailable || reason != "" {
		t.Fatalf("expected available device availability, got status=%s reason=%s", status, reason)
	}
}

func TestService_SetAvailabilityPersistsWithSQLRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSQLRepository(db)
	svc := NewService(repo, nil)

	lastSuccess := true
	now := time.Now()
	if err := repo.Create(context.Background(), &schema.Device{
		ID:               "dev-sql-availability",
		Name:             "dev-sql-availability",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":2}`,
		LastTestSuccess:  &lastSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}

	if _, err := svc.SetAvailability(context.Background(), "dev-sql-availability", AvailabilityStatusUnavailable, "device form is invalid"); err != nil {
		t.Fatalf("set unavailable failed: %v", err)
	}

	savedDevice, err := repo.GetByID(context.Background(), "dev-sql-availability")
	if err != nil {
		t.Fatalf("reload device failed: %v", err)
	}
	if status, reason := AvailabilityOf(savedDevice); status != AvailabilityStatusUnavailable || reason != "device form is invalid" {
		t.Fatalf("expected persisted unavailable device availability, got status=%s reason=%s readiness=%s", status, reason, savedDevice.ReadinessStatus)
	}
}
