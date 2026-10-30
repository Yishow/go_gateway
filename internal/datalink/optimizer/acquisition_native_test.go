package optimizer

import (
	"encoding/json"
	"go-gateway/internal/datalink/schema"
	"os"
	"strconv"
	"testing"
)

func TestAcquisitionPlanNativeAcceptanceVectors(t *testing.T) {
	data, err := os.ReadFile("testdata/acquisition_acceptance.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct{ Cases []acceptanceCase }
	if err = json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	selected := map[string]bool{"invalid-address": true, "scalar-over-limit": true, "mc-unknown-capability": true, "fatek-unknown-capability": true, "fatek-D-float32": true, "fatek-DR-float32": true, "mc-hex-W1A": true, "gap-not-authorized": true}
	for _, c := range fixtures.Cases {
		if !selected[c.ID] {
			continue
		}
		t.Run(c.ID, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Coverage = ReadableCoverage{}
			s.Identity.Protocol = c.Protocol
			s.Capability.ProtocolLimit = c.Limit
			s.Points = nil
			for _, p := range c.Points {
				address := p.Address
				if address == "" {
					prefix := "D"
					if c.Protocol == "modbus" {
						prefix = "HR"
					}
					address = prefix + strconv.Itoa(p.Start)
				}
				s.Points = append(s.Points, EffectivePoint{ID: p.ID, Address: address, DataType: schema.DataType(p.Type), NativeWidth: p.Width})
			}
			plan := BuildAcquisitionPlan(in)
			if len(plan.Errors()) != len(c.Errors) {
				t.Fatalf("errors %+v expected %v", plan.Errors(), c.Errors)
			}
			if len(c.Errors) > 0 && plan.Errors()[0].Code != c.Reason {
				t.Fatalf("wrong reason: %+v", plan.Errors())
			}
			r := plan.Requests()
			if len(r) != len(c.Requests) {
				t.Fatalf("requests %+v expected %+v", r, c.Requests)
			}
			for n, want := range c.Requests {
				if r[n].Start != want.Start || r[n].Length != want.Length || r[n].Kind != want.Kind {
					t.Fatalf("request %+v expected %+v", r[n], want)
				}
				if c.NativeBits != 0 && r[n].Points[0].NativeBits != c.NativeBits {
					t.Fatal("wrong native decoder width")
				}
				if c.Reason == "random-capability-unknown" && r[n].FallbackReason != c.Reason {
					t.Fatal("missing safe fallback reason")
				}
			}
		})
	}
}

func TestAcquisitionPlanRejectsUnsafeScalar(t *testing.T) {
	cases := []struct {
		name   string
		change func(*AcquisitionSource)
		code   string
	}{
		{"negative-address", func(s *AcquisitionSource) { s.Points[0].Address = "HR-1" }, "invalid-address"},
		{"out-of-range", func(s *AcquisitionSource) { s.Points[0].Address = "HR65536" }, "invalid-address"},
		{"cross-end", func(s *AcquisitionSource) {
			s.Points[0].Address = "HR65535"
			s.Points[0].DataType = schema.DataTypeInt32
			s.Points[0].NativeWidth = 2
		}, "invalid-address"},
		{"wrong-width", func(s *AcquisitionSource) { s.Points[0].NativeWidth = 2 }, "invalid-native-width"},
		{"unknown-protocol", func(s *AcquisitionSource) { s.Identity.Protocol = "unknown" }, "unsupported-protocol"},
		{"missing-limit", func(s *AcquisitionSource) { s.Capability.ProtocolLimit = 0 }, "invalid-capability"},
		{"negative-limit", func(s *AcquisitionSource) { s.Capability.DeviceLimit = -1 }, "invalid-capability"},
		{"client-limit", func(s *AcquisitionSource) {
			s.Capability.ClientLimit = 1
			s.Points[0].DataType = schema.DataTypeFloat32
			s.Points[0].NativeWidth = 2
		}, "scalar-exceeds-limit"},
		{"bit-numeric", func(s *AcquisitionSource) { s.Points[0].Address = "CO0" }, "unsupported-type"},
		{"bit-index-numeric", func(s *AcquisitionSource) { s.Points[0].Address = "HR0.5" }, "invalid-address"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := acquisitionInput()
			c.change(&in.Sources[0])
			p := BuildAcquisitionPlan(in)
			if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != c.code {
				t.Fatalf("unsafe scalar accepted/wrong error: %+v %+v", p.Requests(), p.Errors())
			}
		})
	}
}

func TestAcquisitionPlanBitAndNormalizedOffsets(t *testing.T) {
	for _, c := range []struct {
		address          string
		dt               schema.DataType
		start, bits, bit int
	}{{"40100", schema.DataTypeInt16, 99, 16, -1}, {"CO7", schema.DataTypeBool, 7, 1, -1}, {"DI8", schema.DataTypeBool, 8, 1, -1}} {
		in := acquisitionInput()
		s := &in.Sources[0]
		s.Points[0].Address = c.address
		s.Points[0].DataType = c.dt
		p := BuildAcquisitionPlan(in)
		r := p.Requests()
		if len(r) != 1 || r[0].Start != c.start || r[0].Points[0].NativeBits != c.bits || r[0].Points[0].BitIndex != c.bit {
			t.Fatalf("%s: %+v errors %+v", c.address, r, p.Errors())
		}
	}
}

func TestAcquisitionPlanAllPointsPreservedAndLimits(t *testing.T) {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Capability.ClientLimit = 120
	s.Capability.DeviceLimit = 100
	p := s.Points[0]
	for n := 1; n < 21; n++ {
		q := p
		q.ID = "p" + strconv.Itoa(n)
		q.Address = "HR" + strconv.Itoa(n)
		s.Points = append(s.Points, q)
	}
	plan := BuildAcquisitionPlan(in)
	if len(plan.Requests()) != 1 || len(plan.Requests()[0].Points) != 21 {
		t.Fatal("acquisition set pruned")
	}
	for _, r := range plan.Requests() {
		if r.Limit != 100 || r.FallbackReason != "safe-contiguous-plan" {
			t.Fatalf("limits/fallback missing: %+v", r)
		}
	}
}

// Register bit extraction must not reinterpret a one-bit Modbus area as a word.
func TestAcquisitionPlanRejectsBitIndexOnNativeBitArea(t *testing.T) {
	for _, address := range []string{"COIL7.5", "DI7.5", "00007.5", "10007.5"} {
		t.Run(address, func(t *testing.T) {
			in := acquisitionInput()
			in.Sources[0].Points[0].Address = address
			in.Sources[0].Points[0].DataType = schema.DataTypeBool
			p := BuildAcquisitionPlan(in)
			if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "invalid-address" {
				t.Fatalf("native bit area accepted word extraction: %+v errors %+v", p.Requests(), p.Errors())
			}
		})
	}
	for _, address := range []string{"CO7", "DI7"} {
		t.Run(address, func(t *testing.T) {
			in := acquisitionInput()
			in.Sources[0].Points[0].Address = address
			in.Sources[0].Points[0].DataType = schema.DataTypeBool
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != 1 || r[0].Start != 7 || r[0].Points[0].NativeBits != 1 || r[0].Points[0].BitIndex != -1 {
				t.Fatalf("valid native bit rejected: %+v errors %+v", r, p.Errors())
			}
		})
	}
}
