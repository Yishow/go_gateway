package measurement

import (
	"context"
	"testing"
)

func TestMeasurementService_CRUD_And_EpochIncrement(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	def := &MeasurementDefinition{
		ID:           "meas-1",
		WorkspaceID:  "ws-1",
		DeviceID:     "dev-1",
		PointID:      "pt-1",
		Name:         "Total Power",
		Quantity:     "power",
		Unit:         "kW",
		SemanticKind: SemanticKindGauge,
	}

	if err := svc.CreateMeasurement(ctx, def); err != nil {
		t.Fatalf("CreateMeasurement failed: %v", err)
	}

	got, err := svc.GetMeasurement(ctx, "meas-1")
	if err != nil {
		t.Fatalf("GetMeasurement failed: %v", err)
	}
	if got.SeriesEpoch != "epoch-1" || got.DefinitionRevision != "rev-1" {
		t.Errorf("expected epoch-1 rev-1, got %s %s", got.SeriesEpoch, got.DefinitionRevision)
	}

	// Cosmetic rename: should NOT increment epoch, but increments revision
	got.Name = "Total Active Power"
	needsEpoch, err := svc.UpdateMeasurement(ctx, got)
	if err != nil {
		t.Fatalf("UpdateMeasurement failed: %v", err)
	}
	if needsEpoch {
		t.Errorf("expected cosmetic rename to not trigger epoch transition")
	}
	if got.SeriesEpoch != "epoch-1" || got.DefinitionRevision != "rev-2" {
		t.Errorf("expected epoch-1 rev-2, got %s %s", got.SeriesEpoch, got.DefinitionRevision)
	}

	// Comparability breaking change: unit changed to MW
	got.Unit = "MW"
	needsEpoch2, err := svc.UpdateMeasurement(ctx, got)
	if err != nil {
		t.Fatalf("UpdateMeasurement failed: %v", err)
	}
	if !needsEpoch2 {
		t.Errorf("expected unit change to trigger epoch transition")
	}
	if got.SeriesEpoch != "epoch-2" || got.DefinitionRevision != "rev-3" {
		t.Errorf("expected epoch-2 rev-3, got %s %s", got.SeriesEpoch, got.DefinitionRevision)
	}
}

func TestMeasurementService_ApplyTemplate_RequiresConfirmation(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	preview, err := GenerateTemplatePreview("three_phase_power_meter_v1", "meter-1", "Power Meter 1", "40001")
	if err != nil {
		t.Fatalf("GenerateTemplatePreview failed: %v", err)
	}

	// Try to apply unconfirmed preview -> should fail
	err = svc.ApplyTemplate(ctx, "ws-1", preview)
	if err == nil {
		t.Fatalf("expected error when applying unconfirmed template preview")
	}

	// Confirm preview and apply
	preview.Confirmed = true
	err = svc.ApplyTemplate(ctx, "ws-1", preview)
	if err != nil {
		t.Fatalf("ApplyTemplate failed with confirmed preview: %v", err)
	}

	list, err := svc.ListByDevice(ctx, "ws-1", "meter-1")
	if err != nil {
		t.Fatalf("ListByDevice failed: %v", err)
	}
	if len(list) != 8 {
		t.Fatalf("expected 8 measurements created, got %d", len(list))
	}
}
