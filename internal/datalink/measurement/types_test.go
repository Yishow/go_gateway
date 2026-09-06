package measurement

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
)

func TestMeasurementDefinition_Validation(t *testing.T) {
	tests := []struct {
		name    string
		def     MeasurementDefinition
		wantErr bool
	}{
		{
			name: "valid gauge definition",
			def: MeasurementDefinition{
				ID:           "meas-1",
				WorkspaceID:  "ws-1",
				DeviceID:     "dev-1",
				PointID:      "pt-1",
				EquipmentID:  "eq-1",
				Name:         "Temperature",
				Quantity:     "temperature",
				Unit:         "°C",
				SemanticKind: SemanticKindGauge,
			},
			wantErr: false,
		},
		{
			name: "empty id",
			def: MeasurementDefinition{
				DeviceID:     "dev-1",
				PointID:      "pt-1",
				SemanticKind: SemanticKindGauge,
			},
			wantErr: true,
		},
		{
			name: "invalid semantic kind",
			def: MeasurementDefinition{
				ID:           "meas-1",
				DeviceID:     "dev-1",
				PointID:      "pt-1",
				SemanticKind: SemanticKind("unsupported"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.def.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if tt.def.SeriesEpoch == "" {
					t.Errorf("expected default series_epoch to be set")
				}
				if tt.def.EquipmentID == "" {
					t.Errorf("expected default equipment_id to be set")
				}
			}
		})
	}
}

func TestSampleEnvelope_JSONSerialization_LargeIntegers(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	const bigUint uint64 = 9007199254740993 // > 2^53 - 1

	sample := SampleEnvelope{
		SampleID:           "s-1",
		MeasurementID:      "meas-1",
		SeriesEpoch:        "epoch-1",
		DefinitionRevision: "rev-1",
		AcquisitionID:      "acq-1",
		ObservedAt:         now,
		ReceivedAt:         now,
		TimeOrigin:         "device",
		ValueType:          "uint64",
		Value:              bigUint,
		Quality:            schema.QualityGood,
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	rawJSON := string(data)
	if !strings.Contains(rawJSON, `"value_str":"9007199254740993"`) {
		t.Errorf("expected value_str to contain full precision string, got %s", rawJSON)
	}
	if !strings.Contains(rawJSON, `"value":"9007199254740993"`) {
		t.Errorf("expected value to be formatted as string for safe JS handling, got %s", rawJSON)
	}
}

func TestSampleEnvelope_JSONSerialization_NaNInf(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	sample := SampleEnvelope{
		SampleID:      "s-2",
		MeasurementID: "meas-2",
		ObservedAt:    now,
		ReceivedAt:    now,
		ValueType:     "float64",
		Value:         math.NaN(),
		Quality:       schema.QualityGood,
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	rawJSON := string(data)
	if !strings.Contains(rawJSON, `"quality":"bad"`) {
		t.Errorf("expected quality to be bad for NaN value, got %s", rawJSON)
	}
}
