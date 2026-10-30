package optimizer

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"
)

// These oracles use real serializers and complete response fixtures. They do
// not repeat the planner's size equations (notably MC nibbles/FATEK ASCII).
func TestAcquisitionPlanFrameSizesAgainstSerializers(t *testing.T) {
	mcFrame := mcprotocol.NewRequestFrame(0, 255, 0)
	mcRequest := mcFrame.BuildPacket(mcprotocol.CmdBatchRead, mcprotocol.SubCmdWord, []byte{0, 0, 0, 0xA8, 3, 0})
	mcBitRequest := mcFrame.BuildPacket(mcprotocol.CmdBatchRead, mcprotocol.SubCmdBit, []byte{0, 0, 0, 0x90, 3, 0})
	mcResponse := []byte{0xD0, 0, 0, 255, 255, 3, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	mcBitResponse := []byte{0xD0, 0, 0, 255, 255, 3, 0, 4, 0, 0, 0, 0x10, 0x10}
	for _, response := range [][]byte{mcResponse, mcBitResponse} {
		n, err := mcprotocol.ParseResponseHeader(response[:9])
		if err != nil || n != len(response)-9 {
			t.Fatal("invalid MC response fixture")
		}
	}
	for _, c := range []struct {
		protocol, frame, area string
		dt                    schema.DataType
		request, response     []byte
	}{
		{"modbus", "tcp", "HR", schema.DataTypeInt16, modbus.BuildTCPFrame(1, 1, 3, modbus.BuildReadRequest(0, 3)), modbus.BuildTCPFrame(1, 1, 3, []byte{6, 0, 0, 0, 0, 0, 0})},
		{"modbus", "rtu", "HR", schema.DataTypeInt16, modbus.BuildRTUFrame(1, 3, modbus.BuildReadRequest(0, 3)), modbus.BuildRTUFrame(1, 3, []byte{6, 0, 0, 0, 0, 0, 0})},
		{"modbus", "tcp", "COIL", schema.DataTypeBool, modbus.BuildTCPFrame(1, 1, 1, modbus.BuildReadRequest(0, 3)), modbus.BuildTCPFrame(1, 1, 1, []byte{1, 5})},
		{"modbus", "rtu", "DI", schema.DataTypeBool, modbus.BuildRTUFrame(1, 2, modbus.BuildReadRequest(0, 3)), modbus.BuildRTUFrame(1, 2, []byte{1, 5})},
		{"mc", "3e-binary", "D", schema.DataTypeInt16, mcRequest, mcResponse},
		{"mc", "3e-binary", "M", schema.DataTypeBool, mcBitRequest, mcBitResponse},
		{"fatek", "ascii", "D", schema.DataTypeInt16, fatek.BuildFrame(1, "46", "03D00000"), fatek.BuildFrame(1, "46", "0000000000000")},
		{"fatek", "ascii", "DR", schema.DataTypeFloat32, fatek.BuildFrame(1, "46", "03DR00000"), fatek.BuildFrame(1, "46", "0000000000000000000000000")},
		{"fatek", "ascii", "M", schema.DataTypeBool, fatek.BuildFrame(1, "44", "03M0000"), fatek.BuildFrame(1, "44", "0101")},
	} {
		t.Run(c.protocol+c.frame+c.area, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Coverage = ReadableCoverage{}
			s.Identity.Protocol = c.protocol
			s.Identity.Frame = c.frame
			s.Points = nil
			for n := 0; n < 3; n++ {
				s.Points = append(s.Points, EffectivePoint{ID: fmt.Sprint(n), Address: fixtureNativeAddress(c.protocol, c.area, n), DataType: c.dt, NativeWidth: 1})
			}
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != 1 || !r[0].FrameCostKnown || r[0].FrameBytes.RequestBytes != len(c.request) || r[0].FrameBytes.ResponseBytes != len(c.response) {
				t.Fatalf("cost differs from serializer: %+v errors%+v, expected %d/%d", r, p.Errors(), len(c.request), len(c.response))
			}
		})
	}
}

func TestAcquisitionPlanFrameBudgetBoundaries(t *testing.T) {
	for _, c := range []struct {
		name     string
		change   func(*AcquisitionSource)
		requests int
		code     string
	}{
		{"response-split", func(s *AcquisitionSource) { s.Capability.DeviceFrameBudget.ResponseBytes = 11 }, 2, ""},
		{"request-exact", func(s *AcquisitionSource) { s.Capability.DeviceFrameBudget.RequestBytes = 12 }, 1, ""},
		{"request-small", func(s *AcquisitionSource) { s.Capability.DeviceFrameBudget.RequestBytes = 11 }, 0, "scalar-exceeds-byte-limit"},
		{"response-small", func(s *AcquisitionSource) { s.Capability.DeviceFrameBudget.ResponseBytes = 10 }, 0, "scalar-exceeds-byte-limit"},
		{"strictest", func(s *AcquisitionSource) {
			s.Capability.ProtocolFrameBudget.ResponseBytes = 100
			s.Capability.ClientFrameBudget.ResponseBytes = 11
			s.Capability.DeviceFrameBudget.ResponseBytes = 50
		}, 2, ""},
		{"negative", func(s *AcquisitionSource) { s.Capability.DeviceFrameBudget.RequestBytes = -1 }, 0, "invalid-capability"},
		{"unknown-frame-budget", func(s *AcquisitionSource) {
			s.Identity.Frame = "unknown"
			s.Capability.DeviceFrameBudget.ResponseBytes = 100
		}, 0, "unsupported-frame-cost"},
		{"unknown-frame-no-budget", func(s *AcquisitionSource) { s.Identity.Frame = "unknown" }, 1, ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			q := s.Points[0]
			q.ID = "q"
			q.Address = "HR1"
			s.Points = append(s.Points, q)
			c.change(s)
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			if len(r) != c.requests {
				t.Fatalf("requests%d want%d errors%+v", len(r), c.requests, p.Errors())
			}
			if c.code != "" {
				if len(p.Errors()) != 2 || p.Errors()[0].Code != c.code {
					t.Fatalf("expected%s errors%+v", c.code, p.Errors())
				}
			} else if len(p.Errors()) != 0 {
				t.Fatalf("unexpected errors%+v", p.Errors())
			}
			if c.name == "unknown-frame-no-budget" && (r[0].FrameCostKnown || r[0].FrameBytes != (ApplicationFrameBudget{})) {
				t.Fatal("unknown cost claimed compliance")
			}
		})
	}
}

func TestAcquisitionPlanBudgetFingerprint(t *testing.T) {
	in := acquisitionInput()
	original := BuildAcquisitionPlan(in).Fingerprint()
	in.Sources[0].Capability.DeviceFrameBudget.ResponseBytes = 100
	if BuildAcquisitionPlan(in).Fingerprint() == original {
		t.Fatal("frame budget not revision bound")
	}
}

func TestAcquisitionPlanEncodedResponseBudgetRegression(t *testing.T) {
	for _, c := range []struct {
		protocol, frame, area   string
		dt                      schema.DataType
		budget, first, requests int
	}{
		{"modbus", "rtu", "HR", schema.DataTypeInt16, 9, 2, 2},
		{"modbus", "tcp", "COIL", schema.DataTypeBool, 10, 3, 1},
		{"mc", "3e-binary", "D", schema.DataTypeInt16, 15, 2, 2},
		{"mc", "3e-binary", "M", schema.DataTypeBool, 12, 2, 2},
		{"fatek", "ascii", "D", schema.DataTypeInt16, 17, 2, 2},
		{"fatek", "ascii", "DR", schema.DataTypeFloat32, 17, 1, 3},
		{"fatek", "ascii", "M", schema.DataTypeBool, 11, 2, 2},
	} {
		t.Run(c.protocol+c.area, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Identity.Protocol = c.protocol
			s.Identity.Frame = c.frame
			s.Coverage = ReadableCoverage{}
			s.Points = nil
			s.Capability.ClientFrameBudget.ResponseBytes = c.budget
			for n := 0; n < 3; n++ {
				s.Points = append(s.Points, EffectivePoint{ID: fmt.Sprint(n), Address: fixtureNativeAddress(c.protocol, c.area, n), DataType: c.dt, NativeWidth: 1})
			}
			p := BuildAcquisitionPlan(in)
			r := p.Requests()
			if len(p.Errors()) != 0 || len(r) != c.requests || r[0].Length != c.first {
				t.Fatalf("encoded budget split: %+v errors%+v", r, p.Errors())
			}
			for _, request := range r {
				if !request.FrameCostKnown || request.FrameBytes.ResponseBytes > c.budget {
					t.Fatal("request exceeds encoded byte budget")
				}
			}
		})
	}
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Points[0].DataType = schema.DataTypeInt32
	s.Points[0].NativeWidth = 2
	s.Capability.DeviceFrameBudget.ResponseBytes = 12
	p := BuildAcquisitionPlan(in)
	if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != "scalar-exceeds-byte-limit" {
		t.Fatal("multiword scalar split to satisfy byte limit")
	}
}
