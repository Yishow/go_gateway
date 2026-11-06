package connector

import (
	"testing"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

func TestReviewAcquisitionProfileControlsMCXYRadix(t *testing.T) {
	for _, tc := range []struct {
		name, config string
		addresses    map[string]int
	}{
		{
			name:      "Q and iQ-R hexadecimal",
			config:    `{"acquisition":{"version":1,"profile":"mitsubishi-q-iqr-3e-binary-v1"}}`,
			addresses: map[string]int{"X10": 16, "Y10": 16, "X1A": 26, "Y2F": 47},
		},
		{
			name:      "FX5 existing radix",
			config:    `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`,
			addresses: map[string]int{"X10": 8, "Y10": 8, "X17": 15, "Y20": 16, "X18": 18},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, tc.config)
			if err != nil {
				t.Fatalf("resolve profile: %v", err)
			}
			source := optimizer.AcquisitionSource{Identity: config.Identity("device"), Capability: config.Capability(1), Coverage: config.Coverage()}
			for address := range tc.addresses {
				source.Points = append(source.Points, optimizer.EffectivePoint{ID: address, Address: address, DataType: schema.DataTypeBool, NativeWidth: 1})
			}
			plan := optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{Sources: []optimizer.AcquisitionSource{source}})
			if len(plan.Errors()) != 0 {
				t.Fatalf("unexpected planning errors: %+v", plan.Errors())
			}
			seen := make(map[string]bool)
			for _, request := range plan.Requests() {
				for _, point := range request.Points {
					address := request.Start + point.Offset*request.AddressStep
					if request.Kind == "random" {
						if point.Offset < 0 || point.Offset >= len(request.Units) || point.BitIndex < 0 {
							t.Fatalf("invalid random point slice: %+v", point)
						}
						address = request.Units[point.Offset].Address + point.BitIndex
					}
					if address != tc.addresses[point.Point.ID] {
						t.Errorf("%s resolved to %d, want %d", point.Point.ID, address, tc.addresses[point.Point.ID])
					}
					seen[point.Point.ID] = true
				}
			}
			for address := range tc.addresses {
				if !seen[address] {
					t.Errorf("point %s was not planned", address)
				}
			}
		})
	}
}

func TestReviewFX5StillRejectsHexOnlyXYAddress(t *testing.T) {
	config, err := ResolveAcquisitionConfig(schema.ProtocolMC3E, `{"acquisition":{"version":1,"profile":"mitsubishi-fx5-cpu-3e-binary-v1"}}`)
	if err != nil {
		t.Fatalf("resolve profile: %v", err)
	}
	source := optimizer.AcquisitionSource{
		Identity:   config.Identity("device"),
		Capability: config.Capability(1),
		Points:     []optimizer.EffectivePoint{{ID: "x1a", Address: "X1A", DataType: schema.DataTypeBool, NativeWidth: 1}},
	}
	plan := optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{Sources: []optimizer.AcquisitionSource{source}})
	if len(plan.Errors()) != 1 || plan.Errors()[0].Code != "invalid-address" || len(plan.Requests()) != 0 {
		t.Fatalf("FX5 behavior changed: requests %+v, errors %+v", plan.Requests(), plan.Errors())
	}
}
