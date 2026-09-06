package measurement

import (
	"math"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestValidateMixedLayout_NoOverlap(t *testing.T) {
	// D0..D7: 7 items (D3-D4 is uint32)
	mult := 0.1
	items := []MixedItem{
		{
			ItemID:          "item-temp",
			Name:            "Temperature",
			Address:         "D0",
			RegisterOffset:  0,
			DataType:        schema.DataTypeInt16,
			ScaleMultiplier: &mult,
			SemanticKind:    SemanticKindGauge,
		},
		{
			ItemID:         "item-press",
			Name:           "Pressure",
			Address:        "D1",
			RegisterOffset: 1,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindGauge,
		},
		{
			ItemID:         "item-flow",
			Name:           "FlowRate",
			Address:        "D2",
			RegisterOffset: 2,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindRate,
		},
		{
			ItemID:         "item-total-flow",
			Name:           "TotalFlow",
			Address:        "D3",
			RegisterOffset: 3,
			DataType:       schema.DataTypeUint32, // takes 2 words: D3, D4
			SemanticKind:   SemanticKindCounter,
		},
		{
			ItemID:         "item-status",
			Name:           "Status",
			Address:        "D5",
			RegisterOffset: 5,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindState,
		},
		{
			ItemID:         "item-alarm",
			Name:           "AlarmBitmask",
			Address:        "D6",
			RegisterOffset: 6,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindEvent,
		},
		{
			ItemID:         "item-prod",
			Name:           "ProductionCount",
			Address:        "D7",
			RegisterOffset: 7,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindCounter,
		},
	}

	if err := ValidateMixedLayout(items); err != nil {
		t.Fatalf("expected valid layout, got error: %v", err)
	}
}

func TestValidateMixedLayout_DetectsOverlap(t *testing.T) {
	items := []MixedItem{
		{
			ItemID:         "item-1",
			Address:        "D3",
			RegisterOffset: 3,
			DataType:       schema.DataTypeUint32, // takes D3, D4
			SemanticKind:   SemanticKindCounter,
		},
		{
			ItemID:         "item-2",
			Address:        "D4",
			RegisterOffset: 4,
			DataType:       schema.DataTypeUint16, // overlaps with D4
			SemanticKind:   SemanticKindGauge,
		},
	}

	err := ValidateMixedLayout(items)
	if err == nil {
		t.Fatalf("expected error for overlapping items, got nil")
	}
}

func TestValidateMixedLayout_BitSlicesSharingWord(t *testing.T) {
	bit0 := 0
	bit1 := 1
	bitLen1 := 1
	bit2 := 2
	bitLen4 := 4

	items := []MixedItem{
		{
			ItemID:         "item-b0",
			Address:        "D6.0",
			RegisterOffset: 6,
			BitOffset:      &bit0,
			BitLength:      &bitLen1,
			DataType:       schema.DataTypeBool,
			SemanticKind:   SemanticKindState,
		},
		{
			ItemID:         "item-b1",
			Address:        "D6.1",
			RegisterOffset: 6,
			BitOffset:      &bit1,
			BitLength:      &bitLen1,
			DataType:       schema.DataTypeBool,
			SemanticKind:   SemanticKindEvent,
		},
		{
			ItemID:         "item-nibble",
			Address:        "D6.2-5",
			RegisterOffset: 6,
			BitOffset:      &bit2,
			BitLength:      &bitLen4,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindState,
		},
	}

	if err := ValidateMixedLayout(items); err != nil {
		t.Fatalf("expected valid bit slice layout, got error: %v", err)
	}
}

func TestDecodeMixedRegisters_MCScenario(t *testing.T) {
	// Raw registers for D0..D7:
	// D0 = 0xFF85 (int16 -123)
	// D1 = 500 (uint16)
	// D2 = 120 (uint16)
	// D3, D4 = uint32 100000 -> in CDAB (Mitsubishi default): low word D3=0x86A0 (34464), high word D4=0x0001 (1) -> 65536 + 34464 = 100000
	// D5 = 1 (Running)
	// D6 = 0x0005 (Alarm bit 0 and bit 2 active)
	// D7 = 42 (Production count)
	rawRegs := []uint16{
		0xFF85,
		500,
		120,
		0x86A0, 0x0001,
		1,
		0x0005,
		42,
	}

	tempScale := 0.1
	pressScale := 0.01

	items := []MixedItem{
		{
			ItemID:          "temp",
			Address:         "D0",
			RegisterOffset:  0,
			DataType:        schema.DataTypeInt16,
			ScaleMultiplier: &tempScale,
			SemanticKind:    SemanticKindGauge,
		},
		{
			ItemID:          "press",
			Address:         "D1",
			RegisterOffset:  1,
			DataType:        schema.DataTypeUint16,
			ScaleMultiplier: &pressScale,
			SemanticKind:    SemanticKindGauge,
		},
		{
			ItemID:         "flow",
			Address:        "D2",
			RegisterOffset: 2,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindRate,
		},
		{
			ItemID:         "vol",
			Address:        "D3",
			RegisterOffset: 3,
			DataType:       schema.DataTypeUint32,
			DataFormat:     "CDAB",
			SemanticKind:   SemanticKindCounter,
		},
		{
			ItemID:         "status",
			Address:        "D5",
			RegisterOffset: 5,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindState,
		},
		{
			ItemID:         "alarm",
			Address:        "D6",
			RegisterOffset: 6,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindEvent,
		},
		{
			ItemID:         "prod",
			Address:        "D7",
			RegisterOffset: 7,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindCounter,
		},
	}

	results, err := DecodeMixedRegisters(rawRegs, items, "CDAB")
	if err != nil {
		t.Fatalf("DecodeMixedRegisters failed: %v", err)
	}

	if len(results) != len(items) {
		t.Fatalf("expected %d results, got %d", len(items), len(results))
	}

	// 1. Temp: -123 * 0.1 = -12.3
	tempRes := results[0]
	if tempRes.Quality != schema.QualityGood {
		t.Errorf("temp quality = %v, want good", tempRes.Quality)
	}
	if v, ok := tempRes.Value.(float64); !ok || math.Abs(v-(-12.3)) > 1e-6 {
		t.Errorf("temp value = %v (%T), want -12.3", tempRes.Value, tempRes.Value)
	}

	// 2. Press: 500 * 0.01 = 5.0
	pressRes := results[1]
	if v, ok := pressRes.Value.(float64); !ok || math.Abs(v-5.0) > 1e-6 {
		t.Errorf("press value = %v, want 5.0", pressRes.Value)
	}

	// 4. Total Volume: uint32 100000
	volRes := results[3]
	if v, ok := volRes.Value.(uint32); !ok || v != 100000 {
		t.Errorf("volume value = %v (%T), want uint32(100000)", volRes.Value, volRes.Value)
	}
}

func TestDecodeMixedRegisters_SingleItemFailureIsolation(t *testing.T) {
	// Only 2 registers provided, but item 2 needs offset 3 (D3)
	rawRegs := []uint16{100, 200}
	items := []MixedItem{
		{
			ItemID:         "item-1",
			RegisterOffset: 0,
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindGauge,
		},
		{
			ItemID:         "item-2",
			RegisterOffset: 3, // out of bounds
			DataType:       schema.DataTypeUint16,
			SemanticKind:   SemanticKindGauge,
		},
	}

	results, err := DecodeMixedRegisters(rawRegs, items, "ABCD")
	if err != nil {
		t.Fatalf("unexpected fatal error: %v", err)
	}

	if results[0].Quality != schema.QualityGood {
		t.Errorf("item 1 should be good, got %v", results[0].Quality)
	}
	if results[1].Quality != schema.QualityMissing && results[1].Quality != schema.QualityBad {
		t.Errorf("item 2 should be missing/bad, got %v", results[1].Quality)
	}
}
