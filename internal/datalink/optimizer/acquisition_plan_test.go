package optimizer

import (
	"go-gateway/internal/datalink/schema"
	"reflect"
	"testing"
	"time"
)

func acquisitionInput() AcquisitionInput {
	return AcquisitionInput{ConfigurationRevision: "cfg1", Sources: []AcquisitionSource{{
		Identity:   AcquisitionIdentity{DeviceID: "device", DeviceRevision: "device1", Station: "1", ConnectionRevision: "conn1", Protocol: "modbus", Frame: "tcp", ProtocolRevision: "protocol1", PollingGroup: "fast", Cadence: 500 * time.Millisecond},
		Capability: AcquisitionCapability{Revision: "cap1", ProtocolLimit: 125},
		Coverage:   ReadableCoverage{Revision: "coverage1", Readable: []NativeRange{{Area: "HR", Start: 0, End: 21}}},
		Points:     []EffectivePoint{{ID: "p", SourceIdentity: "source-p", SourceRevision: "source1", Address: "HR0", Function: "read-holding", DataType: schema.DataTypeInt16, NativeWidth: 1, Format: "ABCD", FormatSource: "point", FormatRevision: "format1"}},
	}}}
}

func TestAcquisitionPlanScalarAndErrors(t *testing.T) {
	in := acquisitionInput()
	good := in.Sources[0].Points[0]
	bad := good
	bad.ID = "bad"
	bad.Address = "not-an-address"
	unsupported := good
	unsupported.ID = "unsupported"
	unsupported.DataType = "imaginary"
	in.Sources[0].Points = append(in.Sources[0].Points, bad, unsupported)
	p := BuildAcquisitionPlan(in)
	if got := p.Errors(); len(got) != 2 || got[0].PointID != "bad" || got[0].Code != "invalid-address" || got[1].PointID != "unsupported" || got[1].Code != "unsupported-type" {
		t.Fatalf("explicit errors: %+v", got)
	}
	r := p.Requests()
	if len(r) != 1 || r[0].Start != 0 || r[0].Length != 1 || r[0].Limit != 125 || r[0].Group.Area != "HR" || r[0].Points[0].Point != good {
		t.Fatalf("valid scalar lost: %+v", r)
	}
	if p.Fingerprint() == "" {
		t.Fatal("missing fingerprint")
	}
}

func TestAcquisitionPlanDeterministicAndDetached(t *testing.T) {
	in := acquisitionInput()
	in.Sources[0].Coverage.Readable[0].Area = "IR" // unrelated coverage cannot authorize HR gaps
	second := in.Sources[0].Points[0]
	second.ID = "q"
	second.Address = "HR10"
	in.Sources[0].Points = append(in.Sources[0].Points, second)
	p := BuildAcquisitionPlan(in)
	before := p.Requests()
	fingerprint := p.Fingerprint()
	in.Sources[0].Points[0], in.Sources[0].Points[1] = in.Sources[0].Points[1], in.Sources[0].Points[0]
	reordered := BuildAcquisitionPlan(in)
	if !reflect.DeepEqual(before, reordered.Requests()) || fingerprint != reordered.Fingerprint() {
		t.Fatal("input order changed plan")
	}
	if len(before) != 2 || before[0].Start != 0 || before[1].Start != 10 {
		t.Fatalf("not sorted scalar requests: %+v", before)
	}
	in.Sources[0].Points[0].Format = "CDAB"
	in.Sources[0].Coverage.Readable[0].End = 999
	before[0].Points[0].Point.Format = "DCBA"
	before[0].Group.Identity.DeviceID = "mutated"
	if p.Requests()[0].Points[0].Point.Format != "ABCD" || p.Requests()[0].Group.Identity.DeviceID != "device" || p.Fingerprint() != fingerprint {
		t.Fatal("plan aliases mutable caller storage")
	}
}

func TestAcquisitionPlanRevisionInvalidation(t *testing.T) {
	changes := map[string]func(*AcquisitionInput){
		"config":        func(i *AcquisitionInput) { i.ConfigurationRevision = "cfg2" },
		"device":        func(i *AcquisitionInput) { i.Sources[0].Identity.DeviceRevision = "device2" },
		"connection":    func(i *AcquisitionInput) { i.Sources[0].Identity.ConnectionRevision = "conn2" },
		"capability":    func(i *AcquisitionInput) { i.Sources[0].Capability.Revision = "cap2" },
		"coverage":      func(i *AcquisitionInput) { i.Sources[0].Coverage.Revision = "coverage2" },
		"format":        func(i *AcquisitionInput) { i.Sources[0].Points[0].Format = "CDAB" },
		"format-source": func(i *AcquisitionInput) { i.Sources[0].Points[0].FormatSource = "device" },
		"source":        func(i *AcquisitionInput) { i.Sources[0].Points[0].SourceRevision = "source2" },
		"cadence":       func(i *AcquisitionInput) { i.Sources[0].Identity.Cadence = 2 * time.Second },
	}
	original := BuildAcquisitionPlan(acquisitionInput()).Fingerprint()
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			in := acquisitionInput()
			change(&in)
			if BuildAcquisitionPlan(in).Fingerprint() == original {
				t.Fatal("revision did not invalidate fingerprint")
			}
		})
	}
}

func TestAcquisitionPlanCompatibilityIsolation(t *testing.T) {
	changes := []func(*AcquisitionSource){
		func(s *AcquisitionSource) { s.Identity.DeviceID = "other" },
		func(s *AcquisitionSource) { s.Identity.Station = "2" },
		func(s *AcquisitionSource) { s.Identity.ConnectionRevision = "conn2" },
		func(s *AcquisitionSource) { s.Identity.Cadence = 2 * time.Second },
		func(s *AcquisitionSource) { s.Identity.PollingGroup = "slow" },
		func(s *AcquisitionSource) { s.Points[0].Function = "other-function" },
		func(s *AcquisitionSource) { s.Points[0].Address = "IR0" },
		func(s *AcquisitionSource) { s.Capability.Revision = "cap2" },
		func(s *AcquisitionSource) { s.Identity.Frame = "rtu" },
	}
	for n, change := range changes {
		in := acquisitionInput()
		other := acquisitionInput().Sources[0]
		change(&other)
		in.Sources = append(in.Sources, other)
		p := BuildAcquisitionPlan(in)
		r := p.Requests()
		if len(r) != 2 || r[0].Group == r[1].Group {
			t.Fatalf("case %d crossed boundary: %+v", n, r)
		}
		in.Sources[0], in.Sources[1] = in.Sources[1], in.Sources[0]
		if BuildAcquisitionPlan(in).Fingerprint() != p.Fingerprint() {
			t.Fatal("source reordering changed fingerprint")
		}
	}
}
