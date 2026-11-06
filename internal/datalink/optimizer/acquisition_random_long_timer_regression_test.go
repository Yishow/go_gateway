package optimizer

import (
	"slices"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestAcquisitionMCLongTimerRejectedAlongsideRandomWords(t *testing.T) {
	for _, protocol := range []string{"mc", "mitsubishi"} {
		t.Run(protocol, func(t *testing.T) {
			in := randomInput(protocol, 10, 20)
			source := &in.Sources[0]
			source.Identity.ProtocolRevision = mcQIQRProfile
			source.Capability.RandomCrossArea = true
			source.Capability.RandomProtocolLimit = 192
			longTimer := EffectivePoint{ID: "ltn0", Address: "LTN0", DataType: schema.DataTypeInt32, NativeWidth: 2}
			source.Points = append(source.Points, longTimer)

			plan := BuildAcquisitionPlan(in)
			if got := plan.Errors(); len(got) != 1 || got[0].PointID != "ltn0" || got[0].Code != "unsupported-mc-long-timer-counter" {
				t.Errorf("LTN0 must have an explicit unsupported error: %+v", got)
			}
			want := []NativeUnit{{Area: "D", Address: 10, Bits: 16}, {Area: "D", Address: 20, Bits: 16}}
			if requests := plan.Requests(); len(requests) != 1 || requests[0].Kind != "random" || !slices.Equal(requests[0].Units, want) {
				t.Errorf("only D10 and D20 may be requested: %+v", requests)
			}
			_, _, err := AcquisitionPointShapeForIdentity(source.Identity, longTimer)
			if err == nil || err.Error() != "unsupported-mc-long-timer-counter" {
				t.Errorf("LTN0 shape must be rejected before scheduling: %v", err)
			}
		})
	}
}

func TestAcquisitionMCLongTimerCounterCurrentValuesAreUnsupported(t *testing.T) {
	for _, tc := range []struct {
		address string
		typeID  schema.DataType
		width   int
	}{
		{"LTN0", schema.DataTypeInt32, 2},
		{"LSTN0", schema.DataTypeInt32, 2},
		{"LCN0", schema.DataTypeInt32, 2},
	} {
		t.Run(tc.address, func(t *testing.T) {
			in := randomInput("mc")
			point := EffectivePoint{ID: tc.address, Address: tc.address, DataType: tc.typeID, NativeWidth: tc.width}
			in.Sources[0].Points = []EffectivePoint{point}
			plan := BuildAcquisitionPlan(in)
			if got := plan.Requests(); len(got) != 0 {
				t.Errorf("unsupported area produced a request: %+v", got)
			}
			if got := plan.Errors(); len(got) != 1 || got[0].PointID != tc.address || got[0].Code != "unsupported-mc-long-timer-counter" {
				t.Errorf("expected unsupported area error: %+v", got)
			}
			if _, _, err := AcquisitionPointShapeForIdentity(in.Sources[0].Identity, point); err == nil || err.Error() != "unsupported-mc-long-timer-counter" {
				t.Errorf("unsupported area shape accepted: %v", err)
			}
		})
	}
}

func TestAcquisitionMCLongTimerCounterBitsStayNonRandom(t *testing.T) {
	for _, protocol := range []string{"mc", "mitsubishi"} {
		for _, area := range []string{"LTS", "LTC", "LSTS", "LSTC", "LCS", "LCC"} {
			t.Run(protocol+"/"+area, func(t *testing.T) {
				in := randomInput(protocol)
				point := EffectivePoint{ID: area, Address: area + "0", DataType: schema.DataTypeBool, NativeWidth: 1}
				in.Sources[0].Points = []EffectivePoint{point}
				plan := BuildAcquisitionPlan(in)
				if got := plan.Errors(); len(got) != 0 {
					t.Errorf("long bit must not be rejected by planning: %+v", got)
				}
				requests := plan.Requests()
				if len(requests) != 1 || requests[0].Kind != "contiguous" || requests[0].Group.Area != area || requests[0].Start != 0 || requests[0].Length != 1 || requests[0].FallbackReason != "random-unsupported-device" || len(requests[0].Points) != 1 || requests[0].Points[0].NativeBits != 1 {
					t.Errorf("long bit must retain its non-random plan: %+v", requests)
				}
				bits, width, err := AcquisitionPointShapeForIdentity(in.Sources[0].Identity, point)
				if err != nil || bits != 1 || width != 1 {
					t.Errorf("long bit shape: bits=%d width=%d err=%v", bits, width, err)
				}
			})
		}
	}
}

func TestAcquisitionMCRandomLimitCountsWordAndDoubleWordTogether(t *testing.T) {
	in := randomInput("mc")
	source := &in.Sources[0]
	source.Capability.RandomProtocolLimit = 2
	source.Points = []EffectivePoint{
		{ID: "d10", Address: "D10", DataType: schema.DataTypeInt16, NativeWidth: 1},
		{ID: "d20", Address: "D20", DataType: schema.DataTypeInt32, NativeWidth: 2},
		{ID: "d40", Address: "D40", DataType: schema.DataTypeInt16, NativeWidth: 1},
		{ID: "d60", Address: "D60", DataType: schema.DataTypeInt32, NativeWidth: 2},
	}
	plan := BuildAcquisitionPlan(in)
	if got := plan.Errors(); len(got) != 0 {
		t.Fatalf("unexpected errors: %+v", got)
	}
	requests := plan.Requests()
	if len(requests) != 2 {
		t.Fatalf("four mixed-width access points with limit 2 need two requests: %+v", requests)
	}
	want := []NativeUnit{
		{Area: "D", Address: 10, Bits: 16}, {Area: "D", Address: 20, Bits: 32},
		{Area: "D", Address: 40, Bits: 16}, {Area: "D", Address: 60, Bits: 32},
	}
	seen := make(map[NativeUnit]bool)
	for i, request := range requests {
		if request.Kind != "random" || request.Limit != 2 || len(request.Units) > 2 {
			t.Errorf("request %d exceeded the combined limit: %+v", i, request)
		}
		for _, unit := range request.Units {
			if !slices.Contains(want, unit) || seen[unit] {
				t.Errorf("unexpected or duplicate unit: %+v", unit)
			}
			seen[unit] = true
		}
	}
	if len(seen) != len(want) {
		t.Errorf("missing word or double word access points: %+v", seen)
	}
}
