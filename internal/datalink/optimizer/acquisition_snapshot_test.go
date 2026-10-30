package optimizer

import (
	"go-gateway/internal/datalink/schema"
	"reflect"
	"sync"
	"testing"
)

func TestAcquisitionPlanSnapshotRegression(t *testing.T) {
	in := acquisitionInput()
	in.Sources[0].Points[0].Address = "bad"
	p := BuildAcquisitionPlan(in)
	errors := p.Errors()
	errors[0].Code = "mutated"
	if p.Errors()[0].Code != "invalid-address" {
		t.Fatal("errors accessor aliases internal storage")
	}

	in = acquisitionInput()
	in.Sources[0].Coverage.Readable = append(in.Sources[0].Coverage.Readable, NativeRange{Area: "HR", Start: 30, End: 40})
	original := BuildAcquisitionPlan(in)
	in.Sources[0].Coverage.Readable[0], in.Sources[0].Coverage.Readable[1] = in.Sources[0].Coverage.Readable[1], in.Sources[0].Coverage.Readable[0]
	if BuildAcquisitionPlan(in).Fingerprint() != original.Fingerprint() {
		t.Fatal("coverage ordering is not canonical")
	}
	// Readers may mutate their copies concurrently without modifying the installed plan.
	var wg sync.WaitGroup
	for n := 0; n < 8; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := original.Requests()
			r[0].Points[0].Point.Format = "mutated"
			if original.Requests()[0].Points[0].Point.Format != "ABCD" {
				t.Error("concurrent reader mutated plan")
			}
		}()
	}
	wg.Wait()
}

func TestAcquisitionPlanScalarTypesRegression(t *testing.T) {
	for _, c := range []struct {
		dt    schema.DataType
		width int
	}{
		{schema.DataTypeBool, 1}, {schema.DataTypeInt16, 1}, {schema.DataTypeUint16, 1},
		{schema.DataTypeInt32, 2}, {schema.DataTypeUint32, 2}, {schema.DataTypeFloat32, 2},
		{schema.DataTypeInt64, 4}, {schema.DataTypeUint64, 4}, {schema.DataTypeFloat64, 4},
	} {
		in := acquisitionInput()
		in.Sources[0].Points[0].DataType = c.dt
		in.Sources[0].Points[0].NativeWidth = c.width
		p := BuildAcquisitionPlan(in)
		r := p.Requests()
		if len(r) != 1 || r[0].Length != c.width || r[0].Points[0].Length != c.width {
			t.Fatalf("%s scalar split/rejected: %+v %+v", c.dt, r, p.Errors())
		}
	}
	for _, protocol := range []string{"mc", "mitsubishi", "fatek"} {
		for _, support := range []string{"", "unknown", "supported", "unsupported"} {
			in := acquisitionInput()
			in.Sources[0].Coverage = ReadableCoverage{}
			in.Sources[0].Identity.Protocol = protocol
			in.Sources[0].Points[0].Address = "D0"
			in.Sources[0].Capability.RandomSupport = support
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			reason := "scalar-plan"
			if support == "" || support == "unknown" {
				reason = "random-capability-unknown"
			}
			if len(r) != 1 || r[0].Kind != "contiguous" || r[0].FallbackReason != reason {
				t.Fatalf("unsafe capability fallback: %+v", r)
			}
			in.Sources[0].Points[0].Address = "garbage"
			if got := BuildAcquisitionPlan(in).Errors(); len(got) != 1 || got[0].Code != "invalid-address" {
				t.Fatalf("malformed %s: %+v", protocol, got)
			}
		}
	}
	// Planning is pure and must not sort the caller's mutable slices in place.
	in := acquisitionInput()
	q := in.Sources[0].Points[0]
	q.ID = "a"
	in.Sources[0].Points = append(in.Sources[0].Points, q)
	before := append([]EffectivePoint(nil), in.Sources[0].Points...)
	BuildAcquisitionPlan(in)
	if !reflect.DeepEqual(before, in.Sources[0].Points) {
		t.Fatal("planner mutated caller points")
	}
}
