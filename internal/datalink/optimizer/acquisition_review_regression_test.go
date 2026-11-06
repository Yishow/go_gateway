package optimizer

import (
	"fmt"
	"slices"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestReviewQIQRAcquisitionXYUsesHexadecimal(t *testing.T) {
	in := randomInput("mc")
	in.Sources[0].Identity.ProtocolRevision = "mitsubishi-q-iqr-3e-binary-v1"
	in.Sources[0].Points = []EffectivePoint{
		{ID: "x10", Address: "X10", DataType: schema.DataTypeBool, NativeWidth: 1},
		{ID: "x20", Address: "X20", DataType: schema.DataTypeBool, NativeWidth: 1},
	}
	plan := BuildAcquisitionPlan(in)
	if len(plan.Errors()) != 0 {
		t.Fatalf("unexpected planning errors: %+v", plan.Errors())
	}
	for _, request := range plan.Requests() {
		for _, point := range request.Points {
			if point.Point.ID == "x10" && request.Kind == "random" {
				unit := request.Units[point.Offset]
				if unit.Address != 16 || point.BitIndex != 0 {
					t.Fatalf("X10 should use word X16 bit 0, got %+v bit %d", unit, point.BitIndex)
				}
				return
			}
		}
	}
	t.Fatal("expected a random request containing X10")
}

func TestReviewMCBitRandomDoesNotCrossRestrictionsOrBounds(t *testing.T) {
	tests := []struct {
		name       string
		coverage   ReadableCoverage
		bounds     string
		bound      NativeAddressBound
		unsafeID   string
		unsafeSpan NativeRange
	}{
		{name: "forbidden neighbor", coverage: ReadableCoverage{Forbidden: []NativeRange{{Area: "M", Start: 6, End: 7}}}, unsafeID: "m5", unsafeSpan: NativeRange{Area: "M", Start: 6, End: 7}},
		{name: "excluded neighbor", coverage: ReadableCoverage{Excluded: []NativeRange{{Area: "M", Start: 6, End: 7}}}, unsafeID: "m5", unsafeSpan: NativeRange{Area: "M", Start: 6, End: 7}},
		{name: "excluded active point", coverage: ReadableCoverage{Excluded: []NativeRange{{Area: "M", Start: 5, End: 6}}}, unsafeID: "m5", unsafeSpan: NativeRange{Area: "M", Start: 5, End: 6}},
		{name: "protocol lower bound", bounds: "protocol", bound: NativeAddressBound{Area: "M", StartInclusive: 1, EndExclusive: 64}, unsafeID: "m5"},
		{name: "client lower bound", bounds: "client", bound: NativeAddressBound{Area: "M", StartInclusive: 1, EndExclusive: 64}, unsafeID: "m5"},
		{name: "device upper bound", bounds: "device", bound: NativeAddressBound{Area: "M", StartInclusive: 0, EndExclusive: 63}, unsafeID: "m48"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := randomInput("mc")
			source := &in.Sources[0]
			source.Capability.RandomCrossArea = true
			source.Coverage = tc.coverage
			switch tc.bounds {
			case "protocol":
				source.Capability.ProtocolAddressBounds = []NativeAddressBound{tc.bound}
			case "client":
				source.Capability.ClientAddressBounds = []NativeAddressBound{tc.bound}
			case "device":
				source.Capability.DeviceAddressBounds = []NativeAddressBound{tc.bound}
			}
			source.Points = []EffectivePoint{
				{ID: "m5", Address: "M5", DataType: schema.DataTypeBool, NativeWidth: 1},
				{ID: "m32", Address: "M32", DataType: schema.DataTypeBool, NativeWidth: 1},
				{ID: "m48", Address: "M48", DataType: schema.DataTypeBool, NativeWidth: 1},
				{ID: "d100", Address: "D100", DataType: schema.DataTypeInt16, NativeWidth: 1},
				{ID: "d200", Address: "D200", DataType: schema.DataTypeInt16, NativeWidth: 1},
			}
			plan := BuildAcquisitionPlan(in)
			if len(plan.Errors()) != 0 {
				t.Fatalf("requested points are valid: %+v", plan.Errors())
			}
			seen := make(map[string]int)
			sawSafeRandom := false
			for _, request := range plan.Requests() {
				for _, point := range request.Points {
					seen[point.Point.ID]++
					if point.Point.ID == tc.unsafeID && (request.Kind != "contiguous" || request.FallbackReason == "") {
						t.Errorf("unsafe bit must have a non-random fallback: %+v", request)
					}
				}
				if request.Kind != "random" {
					continue
				}
				sawSafeRandom = true
				for _, unit := range request.Units {
					if unit.Area != "M" || unit.Bits != 16 {
						continue
					}
					if tc.unsafeSpan.End > tc.unsafeSpan.Start && unit.Address < tc.unsafeSpan.End && unit.Address+16 > tc.unsafeSpan.Start {
						t.Errorf("random unit %+v crosses restricted span %+v", unit, tc.unsafeSpan)
					}
					if tc.bounds != "" && (unit.Address < tc.bound.StartInclusive || unit.Address+16 > tc.bound.EndExclusive) {
						t.Errorf("random unit %+v crosses declared bound %+v", unit, tc.bound)
					}
				}
			}
			if !sawSafeRandom {
				t.Error("safe independent points should remain eligible for random reads")
			}
			for _, point := range source.Points {
				if seen[point.ID] != 1 {
					t.Errorf("point %s scheduled %d times, want once", point.ID, seen[point.ID])
				}
			}
		})
	}
}

func TestReviewMCBitRandomHonorsHalfOpenRestriction(t *testing.T) {
	in := randomInput("mc")
	in.Sources[0].Coverage.Excluded = []NativeRange{{Area: "M", Start: 16, End: 17}}
	in.Sources[0].Points = []EffectivePoint{
		{ID: "m5", Address: "M5", DataType: schema.DataTypeBool, NativeWidth: 1},
		{ID: "m32", Address: "M32", DataType: schema.DataTypeBool, NativeWidth: 1},
	}
	plan := BuildAcquisitionPlan(in)
	if len(plan.Errors()) != 0 || len(plan.Requests()) != 1 || plan.Requests()[0].Kind != "random" {
		t.Fatalf("an adjacent exclusion must not forbid safe random units: %+v, errors %+v", plan.Requests(), plan.Errors())
	}
}

func TestReviewQIQRSourceExclusionGeometryKeepsBitsSafe(t *testing.T) {
	for _, area := range []string{"X", "Y"} {
		t.Run(area, func(t *testing.T) {
			// The current source geometry API only receives "mc" and resolves X10
			// in the legacy radix. Q/iQ-R interprets that disabled point as 0x10.
			excluded, err := AcquisitionPointSpan("mc", EffectivePoint{Address: area + "10", DataType: schema.DataTypeBool})
			if err != nil {
				t.Fatal(err)
			}
			in := randomInput("mc")
			in.Sources[0].Identity.ProtocolRevision = "mitsubishi-q-iqr-3e-binary-v1"
			in.Sources[0].Coverage.Excluded = []NativeRange{excluded}
			in.Sources[0].Points = []EffectivePoint{
				{ID: "active11", Address: area + "11", DataType: schema.DataTypeBool, NativeWidth: 1},
				{ID: "active20", Address: area + "20", DataType: schema.DataTypeBool, NativeWidth: 1},
			}
			plan := BuildAcquisitionPlan(in)
			if len(plan.Errors()) != 0 {
				t.Fatalf("active points should remain readable: %+v", plan.Errors())
			}
			seen := make(map[string]int)
			for _, request := range plan.Requests() {
				for _, point := range request.Points {
					seen[point.Point.ID]++
				}
				if request.Kind != "random" {
					continue
				}
				for _, unit := range request.Units {
					if unit.Area == area && unit.Address <= 16 && unit.Address+16 > 16 {
						t.Errorf("random unit %+v reads disabled %s10; source exclusion %+v lacks the Q profile", unit, area, excluded)
					}
				}
			}
			for _, point := range in.Sources[0].Points {
				if seen[point.ID] != 1 {
					t.Errorf("point %s scheduled %d times, want once", point.ID, seen[point.ID])
				}
			}
		})
	}
}

func TestReviewQIQRSourceExclusionDoesNotFillContiguousGap(t *testing.T) {
	excluded, err := AcquisitionPointSpan("mc", EffectivePoint{Address: "X11", DataType: schema.DataTypeBool})
	if err != nil {
		t.Fatal(err)
	}
	in := randomInput("mc")
	in.Sources[0].Identity.ProtocolRevision = "mitsubishi-q-iqr-3e-binary-v1"
	in.Sources[0].Coverage = ReadableCoverage{
		Revision: "declared-coverage", Readable: []NativeRange{{Area: "X", Start: 16, End: 19}},
		Excluded: []NativeRange{excluded},
	}
	in.Sources[0].Points = []EffectivePoint{
		{ID: "x10", Address: "X10", DataType: schema.DataTypeBool, NativeWidth: 1},
		{ID: "x12", Address: "X12", DataType: schema.DataTypeBool, NativeWidth: 1},
	}
	plan := BuildAcquisitionPlan(in)
	if len(plan.Errors()) != 0 {
		t.Fatalf("active points should remain readable: %+v", plan.Errors())
	}
	seen := make(map[string]int)
	for _, request := range plan.Requests() {
		for _, point := range request.Points {
			seen[point.Point.ID]++
		}
		if request.Kind == "contiguous" && request.Group.Area == "X" && request.Start <= 17 && request.Start+request.Length > 17 {
			t.Errorf("contiguous request %+v reads disabled Q X11 at physical address 17", request)
		}
		for _, unit := range request.Units {
			if unit.Area == "X" && unit.Address <= 17 && unit.Address+16 > 17 {
				t.Errorf("random unit %+v reads disabled Q X11 at physical address 17", unit)
			}
		}
	}
	for _, point := range in.Sources[0].Points {
		if seen[point.ID] != 1 {
			t.Errorf("point %s scheduled %d times, want once", point.ID, seen[point.ID])
		}
	}
}

func TestReviewFatekRandomLimitCountsWeightedWords(t *testing.T) {
	for _, tc := range []struct {
		count, wantRequests int
		wantCosts           []int
	}{{24, 1, []int{48}}, {25, 2, []int{2, 48}}, {30, 2, []int{12, 48}}} {
		t.Run(fmt.Sprint(tc.count), func(t *testing.T) {
			in := randomInput(acquisitionProtocolFatek)
			in.Sources[0].Capability.RandomProtocolLimit = 48
			for i := range tc.count {
				in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("DR%d", i*10), DataType: schema.DataTypeInt32, NativeWidth: 1})
			}
			plan := BuildAcquisitionPlan(in)
			if len(plan.Errors()) != 0 || len(plan.Requests()) != tc.wantRequests {
				t.Fatalf("got %d requests and errors %+v, want %d requests", len(plan.Requests()), plan.Errors(), tc.wantRequests)
			}
			costs := make([]int, 0, len(plan.Requests()))
			seen := make(map[string]int)
			for _, request := range plan.Requests() {
				if request.Kind != "random" {
					t.Fatalf("expected random requests for sparse points: %+v", request)
				}
				cost := 0
				for _, unit := range request.Units {
					cost += max(1, unit.Bits/16)
				}
				costs = append(costs, cost)
				for _, point := range request.Points {
					seen[point.Point.ID]++
					if point.Length != 1 || point.Offset < 0 || point.Offset >= len(request.Units) {
						t.Errorf("point lost its complete native component: %+v", point)
					}
				}
			}
			slices.Sort(costs)
			if !slices.Equal(costs, tc.wantCosts) {
				t.Errorf("weighted words per request = %v, want %v", costs, tc.wantCosts)
			}
			for _, point := range in.Sources[0].Points {
				if seen[point.ID] != 1 {
					t.Errorf("point %s scheduled %d times, want once", point.ID, seen[point.ID])
				}
			}
		})
	}
}

func TestReviewFatekNativeComponentCapIsIndependentOfWordCap(t *testing.T) {
	in := randomInput(acquisitionProtocolFatek)
	in.Sources[0].Capability.RandomProtocolLimit = 48
	in.Sources[0].Capability.DeviceLimit = 10
	for i := range 10 {
		in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("DR%d", i*10), DataType: schema.DataTypeInt32, NativeWidth: 1})
	}
	plan := BuildAcquisitionPlan(in)
	requests := plan.Requests()
	if len(plan.Errors()) != 0 || len(requests) != 1 || requests[0].Kind != "random" || len(requests[0].Units) != 10 {
		t.Fatalf("10 components costing 20 words fit both limits: %+v, errors %+v", requests, plan.Errors())
	}
}
