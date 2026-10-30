package optimizer

import (
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestAcquisitionPlanDeclaredAddressBounds(t *testing.T) {
	for _, c := range []struct {
		name     string
		change   func(*AcquisitionSource)
		requests int
		code     string
	}{
		{"zero-only", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 1}}
		}, 1, ""},
		{"zero-only-reject-one", func(s *AcquisitionSource) {
			s.Points[0].Address = "HR1"
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 1}}
		}, 0, "address-out-of-bounds"},
		{"area-specific", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "IR", StartInclusive: 1, EndExclusive: 2}}
		}, 1, ""},
		{"empty", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 0}}
		}, 0, "invalid-address-bounds"},
		{"inverted", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 2, EndExclusive: 1}}
		}, 0, "invalid-address-bounds"},
		{"negative", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: -1, EndExclusive: 2}}
		}, 0, "invalid-address-bounds"},
		{"missing-area", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{StartInclusive: 0, EndExclusive: 2}}
		}, 0, "invalid-address-bounds"},
		{"duplicate", func(s *AcquisitionSource) {
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 2}, {Area: "HR", StartInclusive: 0, EndExclusive: 3}}
		}, 0, "invalid-address-bounds"},
		{"strictest-protocol", func(s *AcquisitionSource) {
			s.Points[0].Address = "HR3"
			s.Capability.ProtocolAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 3}}
			s.Capability.ClientAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 4}}
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 5}}
		}, 0, "address-out-of-bounds"},
		{"strictest-client-start", func(s *AcquisitionSource) {
			s.Capability.ProtocolAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 3}}
			s.Capability.ClientAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 1, EndExclusive: 4}}
		}, 0, "address-out-of-bounds"},
		{"wide-no-clipping", func(s *AcquisitionSource) {
			s.Points[0].DataType = "int64"
			s.Points[0].NativeWidth = 4
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 3}}
		}, 0, "address-out-of-bounds"},
		{"wire-never-widens", func(s *AcquisitionSource) {
			s.Points[0].Address = "HR65536"
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 100000}}
		}, 0, "invalid-address"},
	} {
		t.Run(c.name, func(t *testing.T) {
			in := acquisitionInput()
			c.change(&in.Sources[0])
			p := BuildAcquisitionPlan(in)
			if len(p.Requests()) != c.requests {
				t.Fatalf("requests%d want%d errors%+v", len(p.Requests()), c.requests, p.Errors())
			}
			if c.code != "" && (len(p.Errors()) != 1 || p.Errors()[0].Code != c.code) {
				t.Fatalf("expected%s errors%+v", c.code, p.Errors())
			}
		})
	}
}

func TestAcquisitionPlanBoundsDoNotAuthorizeGaps(t *testing.T) {
	in := sparseAcquisitionInput()
	s := &in.Sources[0]
	s.Coverage = ReadableCoverage{}
	s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 21}}
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 2 {
		t.Fatal("address bounds granted gap readability")
	}
}

func TestAcquisitionPlanBoundsFingerprintAndOwnership(t *testing.T) {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: "HR", StartInclusive: 0, EndExclusive: 21}, {Area: "IR", StartInclusive: 0, EndExclusive: 21}}
	original := BuildAcquisitionPlan(in)
	s.Capability.DeviceAddressBounds[0], s.Capability.DeviceAddressBounds[1] = s.Capability.DeviceAddressBounds[1], s.Capability.DeviceAddressBounds[0]
	if BuildAcquisitionPlan(in).Fingerprint() != original.Fingerprint() {
		t.Fatal("bound ordering changed fingerprint")
	}
	fingerprint := original.Fingerprint()
	s.Capability.DeviceAddressBounds[0].EndExclusive = 20
	if BuildAcquisitionPlan(in).Fingerprint() == fingerprint {
		t.Fatal("bound change did not invalidate fingerprint")
	}
	if original.Fingerprint() != fingerprint || len(original.Requests()) != 1 {
		t.Fatal("caller mutated plan")
	}
}

func TestAcquisitionPlanNativeDeclaredBoundsRegression(t *testing.T) {
	for _, c := range []struct {
		protocol, area, address, dt string
		start, end, width           int
		valid                       bool
	}{
		{"mc", "W", "W1A", "int16", 26, 27, 1, true},
		{"mc", "W", "W1B", "int16", 26, 27, 1, false},
		{"fatek", "DR", "DR5", "float32", 5, 6, 1, true},
		{"fatek", "D", "D5", "float32", 5, 6, 2, false},
		{"modbus", "COIL", "CO0", "bool", 0, 1, 1, true},
		{"modbus", "DI", "DI1", "bool", 0, 1, 1, false},
	} {
		t.Run(c.protocol+c.address, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Identity.Protocol = c.protocol
			s.Coverage = ReadableCoverage{}
			s.Points[0].Address = c.address
			s.Points[0].DataType = schema.DataType(c.dt)
			s.Points[0].NativeWidth = c.width
			s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: c.area, StartInclusive: c.start, EndExclusive: c.end}}
			p := BuildAcquisitionPlan(in)
			if c.valid {
				if len(p.Requests()) != 1 || len(p.Errors()) != 0 {
					t.Fatalf("valid native scalar rejected %+v", p.Errors())
				}
			} else if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "address-out-of-bounds" {
				t.Fatalf("native span escaped bound: %+v %+v", p.Requests(), p.Errors())
			}
		})
	}
}
