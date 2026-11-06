package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"reflect"
	"testing"
)

func TestAcquisitionRandomScalarDedupAndCompleteBoundary(t *testing.T) {
	in := randomInput("mc", 0, 1, 10)
	in.Sources[0].Points[0].DataType = schema.DataTypeInt64
	in.Sources[0].Points[0].NativeWidth = 4
	in.Sources[0].Points[0].Format = "ABCD"
	in.Sources[0].Points[1].Format = "CDAB"
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 1 || r[0].Kind != acquisitionKindRandom || len(r[0].Units) != 5 || len(r[0].Points) != 3 {
		t.Fatalf("nested span dedup: %+v errors%+v", r, p.Errors())
	}
	for _, point := range r[0].Points {
		_, start, _, _, _ := nativeAddress("mc", point.Point.Address)
		for n := 0; n < point.Length; n++ {
			if r[0].Units[point.Offset+n].Address != start+n {
				t.Fatal("scalar not contiguous in scatter indices")
			}
		}
	}
	in.Sources[0].Capability.DeviceLimit = 4
	r = BuildAcquisitionPlan(in).Requests()
	if len(r) != 2 {
		t.Fatalf("complete scalar split: %+v", r)
	}
	for _, request := range r {
		for _, point := range request.Points {
			if point.Point.ID == "0" && point.Length != 4 {
				t.Fatal("int64 split")
			}
		}
	}
}

func TestAcquisitionRandomFatekNativeWidth(t *testing.T) {
	in := randomInput(acquisitionProtocolFatek, 0, 10, 20)
	for i := range in.Sources[0].Points {
		p := &in.Sources[0].Points[i]
		p.Address = fmt.Sprintf("DR%d", i*10)
		p.DataType = schema.DataTypeInt64
		p.NativeWidth = 2
	}
	r := BuildAcquisitionPlan(in).Requests()
	if len(r) != 1 || r[0].Kind != acquisitionKindRandom || len(r[0].Units) != 6 || r[0].Length != 6 {
		t.Fatalf("DR uses native components: %+v", r)
	}
	for _, u := range r[0].Units {
		if u.Bits != 32 {
			t.Fatal("DR became16bits")
		}
	}
}

func TestAcquisitionRandomCrossAreaRequiresDeclaration(t *testing.T) {
	for _, protocol := range []string{"mc", acquisitionProtocolFatek} {
		in := randomInput(protocol, 0, 10)
		in.Sources[0].Points[1].Address = "R10"
		p := BuildAcquisitionPlan(in)
		if len(p.Requests()) != 2 {
			t.Fatal("cross-area without declaration")
		}
		old := p.Fingerprint()
		in.Sources[0].Capability.RandomCrossArea = true
		p = BuildAcquisitionPlan(in)
		r := p.Requests()
		if len(p.Errors()) != 0 || len(r) != 1 || r[0].Kind != acquisitionKindRandom || r[0].Group.Area != "" || len(r[0].Units) != 2 || p.Fingerprint() == old {
			t.Fatalf("explicit cross area: %+v errors%+v", r, p.Errors())
		}
	}
}

func TestAcquisitionRandomUnknownAndUnsupportedStaySafe(t *testing.T) {
	for _, support := range []string{"", "unknown", "unsupported", "typo"} {
		in := randomInput("mc", 0, 10, 20)
		in.Sources[0].Capability.RandomSupport = support
		r := BuildAcquisitionPlan(in).Requests()
		if len(r) != 3 {
			t.Fatal("unknown enabled random")
		}
		for _, request := range r {
			if request.Kind != "contiguous" {
				t.Fatal("unknown random")
			}
		}
	}
	in := randomInput("mc", 0, 10)
	for i := range in.Sources[0].Points {
		in.Sources[0].Points[i].Address = fmt.Sprintf("LTN%d", i*10)
	}
	longTimers := BuildAcquisitionPlan(in)
	if r := longTimers.Requests(); len(r) != 0 {
		t.Fatalf("MC cannot request long timers: %+v", r)
	}
	if errors := longTimers.Errors(); len(errors) != 2 || errors[0].Code != "unsupported-mc-long-timer-counter" || errors[1].Code != "unsupported-mc-long-timer-counter" {
		t.Fatalf("MC must report unsupported long timers: %+v", errors)
	}
	in = randomInput("mc", 0, 10)
	for i := range in.Sources[0].Points {
		in.Sources[0].Points[i].Address = fmt.Sprintf("M%d", i*10)
		in.Sources[0].Points[i].DataType = schema.DataTypeBool
	}
	if r := BuildAcquisitionPlan(in).Requests(); len(r) != 1 || r[0].Kind != acquisitionKindRandom || len(r[0].Units) != 1 || r[0].Units[0] != (NativeUnit{Area: "M", Address: 0, Bits: 16}) {
		t.Fatalf("MC random should read M0 and M10 via one aligned M0 word unit: %+v", r)
	}
}

func TestAcquisitionRandomBudgetAndDetached(t *testing.T) {
	in := randomInput("mc", 0, 10, 20)
	in.Sources[0].Capability.DeviceFrameBudget.RequestBytes = 25 // two complete random units
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 2 {
		t.Fatalf("budget packing: %+v", r)
	}
	for _, request := range r {
		if request.FrameBytes.RequestBytes > 25 {
			t.Fatal("random budget exceeded")
		}
	}
	before := p.Requests()
	for i := range r {
		if len(r[i].Units) > 0 {
			r[i].Units[0].Address = 999
		}
	}
	if !reflect.DeepEqual(before, p.Requests()) {
		t.Fatal("random units alias plan")
	}
	in.Sources[0].Capability.DeviceFrameBudget.RequestBytes = 21
	if r := BuildAcquisitionPlan(in).Requests(); len(r) != 3 {
		t.Fatal("single-unit budget")
	}
}

func TestAcquisitionRandomBoundedHybridAndDeterminism(t *testing.T) {
	offsets := make([]int, 260)
	for i := 0; i < 125; i++ {
		offsets[i] = i
	}
	for i := 125; i < len(offsets); i++ {
		offsets[i] = 1000 + i*10
	}
	in := randomInput("mc", offsets...)
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 2 {
		t.Fatalf("hybrid request count: %d", len(r))
	}
	kinds := map[string]int{}
	count := 0
	for _, request := range r {
		kinds[request.Kind]++
		count += len(request.Points)
	}
	if kinds[acquisitionKindRandom] != 1 || kinds["contiguous"] != 1 || count != 260 {
		t.Fatalf("bounded hybrid should save bytes at same count: %v %d", kinds, count)
	}
	for a, b := 0, len(in.Sources[0].Points)-1; a < b; a, b = a+1, b-1 {
		in.Sources[0].Points[a], in.Sources[0].Points[b] = in.Sources[0].Points[b], in.Sources[0].Points[a]
	}
	other := BuildAcquisitionPlan(in)
	if !reflect.DeepEqual(r, other.Requests()) || p.Fingerprint() != other.Fingerprint() {
		t.Fatal("input ordering changes random plan")
	}
}

func TestAcquisitionRandomMCWordDoubleWordAndBitInOne0403(t *testing.T) {
	in := randomInput("mc")
	in.Sources[0].Capability.RandomCrossArea = true
	in.Sources[0].Points = []EffectivePoint{
		{ID: "d39", Address: "D39", DataType: schema.DataTypeInt16, NativeWidth: 1, Format: "CDAB"},
		{ID: "d100", Address: "D100", DataType: schema.DataTypeInt32, NativeWidth: 2, Format: "CDAB"},
		{ID: "m5", Address: "M5", DataType: schema.DataTypeBool, NativeWidth: 1},
	}
	p := BuildAcquisitionPlan(in)
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected plan errors: %+v", p.Errors())
	}
	r := p.Requests()
	if len(r) != 1 || r[0].Kind != acquisitionKindRandom {
		t.Fatalf("expected 1 random request, got %+v", r)
	}
	wantUnits := []NativeUnit{
		{Area: "D", Address: 39, Bits: 16},
		{Area: "M", Address: 0, Bits: 16},
		{Area: "D", Address: 100, Bits: 32},
	}
	if !reflect.DeepEqual(r[0].Units, wantUnits) {
		t.Fatalf("got units %+v, want %+v", r[0].Units, wantUnits)
	}
	if r[0].Length != 3 || r[0].FrameBytes.RequestBytes != 29 || r[0].FrameBytes.ResponseBytes != 19 {
		t.Fatalf("unexpected cost: length=%d frame=%+v", r[0].Length, r[0].FrameBytes)
	}
	byID := map[string]PointSlice{}
	for _, slice := range r[0].Points {
		byID[slice.Point.ID] = slice
	}
	if s := byID["d39"]; s.Offset != 0 || s.Length != 1 || s.NativeBits != 16 || s.BitIndex != -1 {
		t.Fatalf("unexpected D39 slice: %+v", s)
	}
	if s := byID["m5"]; s.Offset != 1 || s.Length != 1 || s.NativeBits != 16 || s.BitIndex != 5 {
		t.Fatalf("unexpected M5 slice: %+v", s)
	}
	if s := byID["d100"]; s.Offset != 2 || s.Length != 1 || s.NativeBits != 32 || s.BitIndex != -1 {
		t.Fatalf("unexpected D100 slice: %+v", s)
	}
}

func TestAcquisitionRandomMCRejectsLongTimerAndFallsBackForString(t *testing.T) {
	in := randomInput("mc")
	in.Sources[0].Capability.RandomCrossArea = true
	in.Sources[0].Points = []EffectivePoint{
		{ID: "d10", Address: "D10", DataType: schema.DataTypeInt16, NativeWidth: 1},
		{ID: "d20", Address: "D20", DataType: schema.DataTypeInt16, NativeWidth: 1},
		{ID: "ltn0", Address: "LTN0", DataType: schema.DataTypeInt32, NativeWidth: 2},
		{ID: "str", Address: "D30", DataType: schema.DataTypeString, NativeWidth: 8},
	}
	p := BuildAcquisitionPlan(in)
	if errors := p.Errors(); len(errors) != 1 || errors[0].PointID != "ltn0" || errors[0].Code != "unsupported-mc-long-timer-counter" {
		t.Fatalf("expected unsupported long timer error: %+v", errors)
	}
	r := p.Requests()
	if len(r) != 2 {
		t.Fatalf("expected 1 random + 1 contiguous request, got %+v", r)
	}
	var randomReq *AcquisitionRequest
	contiguousByArea := map[string]AcquisitionRequest{}
	for i := range r {
		switch r[i].Kind {
		case acquisitionKindRandom:
			randomReq = &r[i]
		case "contiguous":
			contiguousByArea[r[i].Group.Area] = r[i]
		}
	}
	if randomReq == nil || len(randomReq.Units) != 2 || randomReq.Units[0].Address != 10 || randomReq.Units[1].Address != 20 {
		t.Fatalf("expected D10 and D20 in random request, got %+v", randomReq)
	}
	if ltn, ok := contiguousByArea["LTN"]; ok {
		t.Fatalf("LTN0 must not produce a request: %+v", ltn)
	}
	if strReq, ok := contiguousByArea["D"]; !ok || strReq.Start != 30 || strReq.Length != 8 || strReq.FallbackReason == "" {
		t.Fatalf("expected string D30 in non-random request with explicit reason, got %+v", strReq)
	}
}
