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
	if len(r) != 1 || r[0].Kind != "random" || len(r[0].Units) != 5 || len(r[0].Points) != 3 {
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
	in := randomInput("fatek", 0, 10, 20)
	for i := range in.Sources[0].Points {
		p := &in.Sources[0].Points[i]
		p.Address = fmt.Sprintf("DR%d", i*10)
		p.DataType = schema.DataTypeInt64
		p.NativeWidth = 2
	}
	r := BuildAcquisitionPlan(in).Requests()
	if len(r) != 1 || r[0].Kind != "random" || len(r[0].Units) != 6 || r[0].Length != 6 {
		t.Fatalf("DR uses native components: %+v", r)
	}
	for _, u := range r[0].Units {
		if u.Bits != 32 {
			t.Fatal("DR became16bits")
		}
	}
}

func TestAcquisitionRandomCrossAreaRequiresDeclaration(t *testing.T) {
	for _, protocol := range []string{"mc", "fatek"} {
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
		if len(p.Errors()) != 0 || len(r) != 1 || r[0].Kind != "random" || r[0].Group.Area != "" || len(r[0].Units) != 2 || p.Fingerprint() == old {
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
		in.Sources[0].Points[i].Address = fmt.Sprintf("M%d", i*10)
		in.Sources[0].Points[i].DataType = schema.DataTypeBool
	}
	if r := BuildAcquisitionPlan(in).Requests(); len(r) != 2 || r[0].Kind != "contiguous" {
		t.Fatal("MC random cannot read bits")
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
	if kinds["random"] != 1 || kinds["contiguous"] != 1 || count != 260 {
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
