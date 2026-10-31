package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"testing"
)

func TestIndependentFatekWeightedLimit(t *testing.T) {
	in := randomInput("fatek")
	for i := 0; i < 33; i++ {
		in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("DR%d", i*10), DataType: schema.DataTypeInt32, NativeWidth: 1})
	}
	p := BuildAcquisitionPlan(in)
	rs := p.Requests()
	for _, r := range rs {
		n := 0
		for _, u := range r.Units {
			n += max(1, u.Bits/16)
		}
		t.Logf("kind=%s components=%d weighted=%dW limit=%d", r.Kind, len(r.Units), n, r.Limit)
		if n > 64 {
			t.Errorf("Cmd48 exceeds 64W with %dW", n)
		}
	}
	if len(p.Errors()) != 0 {
		t.Fatal(p.Errors())
	}
}
func TestIndependentFatekScalarStride(t *testing.T) {
	in := randomInput("fatek")
	for i := 0; i < 2; i++ {
		in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("DR%d", i*10), DataType: schema.DataTypeInt64, NativeWidth: 2})
	}
	p := BuildAcquisitionPlan(in)
	rs := p.Requests()
	if len(rs) != 1 {
		t.Fatalf("requests %+v errors %+v", rs, p.Errors())
	}
	t.Logf("64bit scalar expansion: %+v", rs[0].Units)
	for _, s := range rs[0].Points {
		base := rs[0].Units[s.Offset].Address
		if got := rs[0].Units[s.Offset+1].Address; got != base+2 {
			t.Errorf("64bit scalar %s second component DR%d, want DR%d", s.Point.Address, got, base+2)
		}
	}
}
func TestIndependentFatekOverlapping32Points(t *testing.T) {
	in := randomInput("fatek")
	for i := 0; i < 2; i++ {
		in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("DR%d", i), DataType: schema.DataTypeInt32, NativeWidth: 1})
	}
	p := BuildAcquisitionPlan(in)
	rs := p.Requests()
	t.Logf("two independent overlapping32bit points: %+v errors %+v", rs, p.Errors())
	if len(rs) == 1 && rs[0].Kind == "contiguous" {
		t.Error("DR0 + DR1 are different valid overlapping components, cannot map to contiguous Cmd46 DR0 count2 (DR0, DR2)")
	}
}
func TestIndependentFatekWordBitAddress(t *testing.T) {
	comp, _ := fatek.GetComponentType("WM")
	got, err := fatek.FormatAddress(comp, 8)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("WM8 format=%q", got)
	if got != "WM0008" {
		t.Errorf("want manual6char WM0008 got %q", got)
	}
	in := randomInput("fatek")
	in.Sources[0].Points = []EffectivePoint{{ID: "a", Address: "WY8", DataType: schema.DataTypeInt32, NativeWidth: 2}, {ID: "b", Address: "WY64", DataType: schema.DataTypeInt16, NativeWidth: 1}}
	p := BuildAcquisitionPlan(in)
	rs := p.Requests()
	t.Logf("word-bit units %+v errors %+v", rs[0].Units, p.Errors())
	for _, point := range rs[0].Points {
		if point.Point.ID == "a" {
			if got := rs[0].Units[point.Offset+1].Address; got != 24 {
				t.Errorf("WY8 second component WY%d, want WY24", got)
			}
		}
	}
}

func TestFatekNativeCrossPhaseScatterSlice(t *testing.T) {
	in := randomInput("fatek")
	in.Sources[0].Points = []EffectivePoint{
		{ID: "a", Address: "DR0", DataType: schema.DataTypeInt64, NativeWidth: 2},
		{ID: "b", Address: "DR1", DataType: schema.DataTypeInt32, NativeWidth: 1},
		{ID: "c", Address: "DR2", DataType: schema.DataTypeInt64, NativeWidth: 2},
	}
	plan := BuildAcquisitionPlan(in)
	rs := plan.Requests()
	if len(plan.Errors()) != 0 || len(rs) != 1 || rs[0].Kind != "random" {
		t.Fatalf("plan=%+v errors=%+v", rs, plan.Errors())
	}
	want := map[string][]int{"a": {0, 2}, "b": {1}, "c": {2, 4}}
	for _, p := range rs[0].Points {
		for i, addr := range want[p.Point.ID] {
			if got := rs[0].Units[p.Offset+i].Address; got != addr {
				t.Errorf("point %s component%d got%d want%d", p.Point.ID, i, got, addr)
			}
		}
	}
}
