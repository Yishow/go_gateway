package measurement

import (
	"testing"
)

func TestGenerateTemplatePreview_ThreePhaseMeter(t *testing.T) {
	preview, err := GenerateTemplatePreview("three_phase_power_meter_v1", "dev-meter-1", "Main Power Meter", "40001")
	if err != nil {
		t.Fatalf("GenerateTemplatePreview failed: %v", err)
	}

	if preview.Confirmed {
		t.Errorf("expected preview.Confirmed to be false before review")
	}
	if !preview.NeedsReview {
		t.Errorf("expected preview.NeedsReview to be true")
	}
	if len(preview.ProposedItems) != 8 {
		t.Fatalf("expected 8 items for 3-phase meter, got %d", len(preview.ProposedItems))
	}

	// Active energy must be counter
	kwhDef := preview.Definitions[7]
	if kwhDef.SemanticKind != SemanticKindCounter || kwhDef.Unit != "kWh" {
		t.Errorf("expected kWh counter definition, got semantic=%s unit=%s", kwhDef.SemanticKind, kwhDef.Unit)
	}
}

func TestGenerateTemplatePreview_MixedSensor(t *testing.T) {
	preview, err := GenerateTemplatePreview("mixed_equipment_sensors_v1", "dev-plc-1", "CNC Machine 1", "")
	if err != nil {
		t.Fatalf("GenerateTemplatePreview failed: %v", err)
	}

	if len(preview.ProposedItems) != 7 {
		t.Fatalf("expected 7 items for mixed sensor, got %d", len(preview.ProposedItems))
	}

	// Verify temp, total flow, and alarm bitmask
	temp := preview.ProposedItems[0]
	if temp.SemanticKind != SemanticKindGauge || temp.ScaleMultiplier == nil || *temp.ScaleMultiplier != 0.1 {
		t.Errorf("expected temp gauge with scale 0.1, got %v", temp)
	}

	totalFlow := preview.ProposedItems[3]
	if totalFlow.SemanticKind != SemanticKindCounter || totalFlow.DataFormat != "CDAB" {
		t.Errorf("expected total flow counter with CDAB format, got %v", totalFlow)
	}
}
