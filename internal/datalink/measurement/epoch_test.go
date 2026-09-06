package measurement

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
)

func TestEpochTransition_RenamePreservesEpoch(t *testing.T) {
	oldDef := MeasurementDefinition{
		ID:           "meas-1",
		DeviceID:     "dev-1",
		PointID:      "pt-1",
		Name:         "Main Voltage A",
		Quantity:     "voltage",
		Unit:         "V",
		SemanticKind: SemanticKindGauge,
		SeriesEpoch:  "epoch-1",
	}

	newDef := oldDef
	newDef.Name = "Line Voltage L1-N" // Renaming only

	needsEpoch, nextEpoch, reason := CheckEpochTransition(oldDef, newDef)
	if needsEpoch {
		t.Fatalf("expected rename to not trigger new epoch, got nextEpoch=%s reason=%s", nextEpoch, reason)
	}
	if nextEpoch != "epoch-1" {
		t.Fatalf("expected epoch-1, got %s", nextEpoch)
	}
}

func TestEpochTransition_UnitOrSourceChangeCreatesEpoch(t *testing.T) {
	oldDef := MeasurementDefinition{
		ID:           "meas-1",
		DeviceID:     "dev-1",
		PointID:      "pt-1",
		Name:         "Active Energy",
		Quantity:     "energy",
		Unit:         "kWh",
		SemanticKind: SemanticKindCounter,
		SeriesEpoch:  "epoch-1",
	}

	// 1. Unit change (kWh -> MWh)
	newDef1 := oldDef
	newDef1.Unit = "MWh"
	needs, next, reason := CheckEpochTransition(oldDef, newDef1)
	if !needs || next != "epoch-2" {
		t.Fatalf("expected epoch-2 on unit change, got needs=%v next=%s reason=%s", needs, next, reason)
	}

	// 2. Point/Device replacement
	newDef2 := oldDef
	newDef2.PointID = "pt-2-new-meter"
	needs2, next2, reason2 := CheckEpochTransition(oldDef, newDef2)
	if !needs2 || next2 != "epoch-2" {
		t.Fatalf("expected epoch-2 on meter replacement, got needs=%v next=%s reason=%s", needs2, next2, reason2)
	}
}

func TestSampleEnvelope_LargeIntegerRoundTrip(t *testing.T) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	const maxUint64 uint64 = 18446744073709551615

	sample := SampleEnvelope{
		SampleID:           "s-max",
		MeasurementID:      "meas-c",
		SeriesEpoch:        "epoch-1",
		DefinitionRevision: "rev-1",
		AcquisitionID:      "acq-1",
		ObservedAt:         now,
		ReceivedAt:         now,
		TimeOrigin:         "device",
		ValueType:          "uint64",
		Value:              maxUint64,
		Quality:            schema.QualityGood,
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if valStr, ok := parsed["value_str"].(string); !ok || valStr != "18446744073709551615" {
		t.Errorf("expected value_str to be preserved as exact decimal string, got %v", parsed["value_str"])
	}
}
