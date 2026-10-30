package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func sparseAcquisitionInput() AcquisitionInput {
	in := acquisitionInput()
	s := &in.Sources[0]
	q := s.Points[0]
	q.ID = "q"
	q.Address = "HR10"
	s.Points = append(s.Points, q)
	return in
}

func TestAcquisitionPlanGapCoverageBoundaries(t *testing.T) {
	for _, c := range []struct {
		name     string
		change   func(*AcquisitionSource)
		requests int
	}{
		{"explicit", func(s *AcquisitionSource) {}, 1},
		{"no-revision", func(s *AcquisitionSource) { s.Coverage.Revision = "" }, 2},
		{"no-coverage", func(s *AcquisitionSource) { s.Coverage.Readable = nil }, 2},
		{"wrong-area", func(s *AcquisitionSource) { s.Coverage.Readable[0].Area = "IR" }, 2},
		{"partial-coverage", func(s *AcquisitionSource) { s.Coverage.Readable[0].End = 9 }, 2},
		{"forbidden-gap", func(s *AcquisitionSource) { s.Coverage.Forbidden = []NativeRange{{Area: "HR", Start: 5, End: 6}} }, 2},
		{"excluded-gap", func(s *AcquisitionSource) { s.Coverage.Excluded = []NativeRange{{Area: "HR", Start: 5, End: 6}} }, 2},
		{"fragment-union", func(s *AcquisitionSource) {
			s.Coverage.Readable = []NativeRange{{Area: "HR", Start: 5, End: 10}, {Area: "HR", Start: 1, End: 5}}
		}, 1},
		{"uncovered-cell", func(s *AcquisitionSource) {
			s.Coverage.Readable = []NativeRange{{Area: "HR", Start: 1, End: 5}, {Area: "HR", Start: 6, End: 10}}
		}, 2},
		{"excluded-active", func(s *AcquisitionSource) { s.Coverage.Excluded = []NativeRange{{Area: "HR", Start: 0, End: 1}} }, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			in := sparseAcquisitionInput()
			c.change(&in.Sources[0])
			p := BuildAcquisitionPlan(in)
			if len(p.Errors()) != 0 || len(p.Requests()) != c.requests {
				t.Fatalf("requests=%d errors=%+v, want%d", len(p.Requests()), p.Errors(), c.requests)
			}
		})
	}
}

func TestAcquisitionPlanRejectsForbiddenScalar(t *testing.T) {
	in := acquisitionInput()
	in.Sources[0].Coverage.Forbidden = []NativeRange{{Area: "HR", Start: 0, End: 1}}
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "forbidden-address" {
		t.Fatalf("forbidden scalar read: %+v %+v", p.Requests(), p.Errors())
	}
}

func TestAcquisitionPlanNativeWireBounds(t *testing.T) {
	for _, c := range []struct {
		protocol, address string
		dt                schema.DataType
		width             int
		valid             bool
	}{
		{"mc", "WFFFFFF", schema.DataTypeInt16, 1, true},
		{"mc", "W1000000", schema.DataTypeInt16, 1, false},
		{"mc", "D16777215", schema.DataTypeInt32, 2, false},
		{"fatek", "D99999", schema.DataTypeInt16, 1, true},
		{"fatek", "D100000", schema.DataTypeInt16, 1, false},
		{"fatek", "D99999", schema.DataTypeFloat32, 2, false},
		{"fatek", "DR99999", schema.DataTypeFloat32, 1, true},
		{"fatek", "DR100000", schema.DataTypeFloat32, 1, false},
		{"fatek", "DR99999", schema.DataTypeFloat64, 2, false},
		{"fatek", "M9999", schema.DataTypeBool, 1, true},
		{"fatek", "M10000", schema.DataTypeBool, 1, false},
	} {
		t.Run(c.protocol+c.address+string(c.dt), func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Coverage = ReadableCoverage{}
			s.Identity.Protocol = c.protocol
			s.Points[0].Address = c.address
			s.Points[0].DataType = c.dt
			s.Points[0].NativeWidth = c.width
			p := BuildAcquisitionPlan(in)
			if c.valid {
				if len(p.Requests()) != 1 || len(p.Errors()) != 0 {
					t.Fatalf("valid bound rejected: %+v", p.Errors())
				}
			} else if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "invalid-address" {
				t.Fatalf("wire-overflow accepted: %+v %+v", p.Requests(), p.Errors())
			}
		})
	}
}

func TestAcquisitionPlanHardNativeLimits(t *testing.T) {
	for _, c := range []struct {
		protocol, area string
		dt             schema.DataType
		limit          int
	}{
		{"modbus", "HR", schema.DataTypeInt16, 125}, {"modbus", "COIL", schema.DataTypeBool, 2000},
		{"mc", "D", schema.DataTypeInt16, 125}, {"mc", "M", schema.DataTypeBool, 2000},
		{"fatek", "D", schema.DataTypeInt16, 64}, {"fatek", "DR", schema.DataTypeFloat32, 32}, {"fatek", "M", schema.DataTypeBool, 255},
	} {
		t.Run(c.protocol+c.area, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Coverage = ReadableCoverage{}
			s.Identity.Protocol = c.protocol
			s.Capability.ProtocolLimit = 10000
			s.Points = nil
			for n := 0; n <= c.limit; n++ {
				s.Points = append(s.Points, EffectivePoint{ID: fmt.Sprint(n), Address: fmt.Sprintf("%s%d", c.area, n), DataType: c.dt, NativeWidth: 1})
			}
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != 2 || r[0].Length != c.limit || r[1].Length != 1 || r[0].Limit != c.limit {
				t.Fatalf("wire client limit missing: requests%d errors%+v", len(r), p.Errors())
			}
		})
	}
}

func TestAcquisitionPlanActiveSpansOverrideExcludedOnly(t *testing.T) {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Points = []EffectivePoint{
		{ID: "wide", Address: "HR0", DataType: schema.DataTypeInt64, NativeWidth: 4},
		{ID: "next", Address: "HR4", DataType: schema.DataTypeInt16, NativeWidth: 1},
	}
	s.Coverage = ReadableCoverage{Excluded: []NativeRange{{Area: "HR", Start: 2, End: 3}}}
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 1 || r[0].Length != 5 {
		t.Fatal("active coverage incorrectly excluded")
	}
}

func TestAcquisitionPlanSameRevisionDifferentDeclarationsStaySeparate(t *testing.T) {
	in := acquisitionInput()
	other := acquisitionInput().Sources[0]
	other.Points[0].ID = "other"
	other.Points[0].Address = "HR1"
	other.Capability.DeviceLimit = 1
	in.Sources = append(in.Sources, other)
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 2 {
		t.Fatal("conflicting capability declarations merged")
	}
}

func TestAcquisitionPlanMCNativeClientClassification(t *testing.T) {
	for _, area := range []string{"TC", "CC"} {
		in := acquisitionInput()
		s := &in.Sources[0]
		s.Coverage = ReadableCoverage{}
		s.Identity.Protocol = "mc"
		s.Identity.Frame = "3e-binary"
		s.Points[0].Address = area + "0"
		s.Points[0].DataType = schema.DataTypeBool
		p := BuildAcquisitionPlan(in)
		r := p.Requests()
		if len(r) != 1 || r[0].Points[0].NativeBits != 1 {
			t.Errorf("%s must use existing MC bit command classification: %+v", area, r)
		}
		s.Points[0].DataType = schema.DataTypeInt16
		p = BuildAcquisitionPlan(in)
		if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "unsupported-type" {
			t.Errorf("%s numeric bit read allowed", area)
		}
	}
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Coverage = ReadableCoverage{}
	s.Identity.Protocol = "mc"
	s.Points[0].Address = "SM0"
	s.Points[0].DataType = schema.DataTypeBool
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "unsupported-area" {
		t.Errorf("unsupported MC client area accepted: %+v %+v", p.Requests(), p.Errors())
	}
}

func TestAcquisitionPlanInvalidCoverageCannotAuthorizeGap(t *testing.T) {
	in := sparseAcquisitionInput()
	in.Sources[0].Coverage.Readable[0].Start = -1
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 0 || len(p.Errors()) != 2 || p.Errors()[0].Code != "invalid-coverage" {
		t.Errorf("malformed coverage authorized: %+v %+v", p.Requests(), p.Errors())
	}
}
