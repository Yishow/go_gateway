package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func randomInput(protocol string, offsets ...int) AcquisitionInput {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Identity.Protocol = protocol
	s.Identity.Frame = "3e-binary"
	if protocol == "fatek" {
		s.Identity.Frame = "ascii"
	}
	s.Capability.ProtocolLimit = 1000
	s.Capability.RandomSupport = "supported"
	s.Coverage = ReadableCoverage{}
	s.Points = nil
	for i, offset := range offsets {
		s.Points = append(s.Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("D%d", offset), DataType: schema.DataTypeInt16, NativeWidth: 1})
	}
	return in
}

func TestAcquisitionRandomSparseAndDense(t *testing.T) {
	for _, protocol := range []string{"mc", "fatek"} {
		t.Run(protocol, func(t *testing.T) {
			p := BuildAcquisitionPlan(randomInput(protocol, 0, 1, 10))
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != 1 || r[0].Kind != "random" {
				t.Fatalf("sparse whole-group candidate: %+v errors%+v", r, p.Errors())
			}
			offsets := make([]int, 21)
			for i := range offsets {
				offsets[i] = i
			}
			p = BuildAcquisitionPlan(randomInput(protocol, offsets...))
			r = p.Requests()
			if len(r) != 1 || r[0].Kind != "contiguous" {
				t.Fatalf("dense must use fewer encoded bytes: %+v", r)
			}
		})
	}
}

func TestAcquisitionRandomNativeLimits(t *testing.T) {
	for _, c := range []struct {
		protocol string
		limit    int
	}{{"mc", 255}, {"fatek", 64}} {
		for _, n := range []int{c.limit, c.limit + 1} {
			offsets := make([]int, n)
			for i := range offsets {
				offsets[i] = i * 10
			}
			p := BuildAcquisitionPlan(randomInput(c.protocol, offsets...))
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != (n+c.limit-1)/c.limit {
				t.Fatalf("%s %d native units: %d requests errors%+v", c.protocol, n, len(r), p.Errors())
			}
		}
	}
}
