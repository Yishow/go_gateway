package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestWordRule23PlannerFunctionParity(t *testing.T) {
	for _, tc := range []struct {
		address, function string
		valid             bool
	}{{"100", "03", true}, {"50000", "01", false}, {"HR100", "01", true}} {
		t.Run(tc.address+tc.function, func(t *testing.T) {
			in := acquisitionInput()
			p := &in.Sources[0].Points[0]
			p.Address = tc.address
			p.Function = tc.function
			p.DataType = schema.DataTypeInt32
			p.NativeWidth = 2
			out := BuildAcquisitionPlan(in)
			if tc.valid && len(out.Errors()) != 0 {
				t.Fatalf("valid transport word rejected: %v", out.Errors())
			}
			if !tc.valid && len(out.Errors()) == 0 {
				t.Fatalf("transport bit area accepted as numeric word: %v", out.Requests())
			}
		})
	}
}

// COIL is the planner's native area, not a legal adapter address alias.
func fixtureNativeAddress(protocol, area string, offset int) string {
	if protocol == "modbus" && area == "COIL" {
		area = "CO"
	}
	return fmt.Sprintf("%s%d", area, offset)
}

// Expected offsets/functions follow parseModbusAddress's wire precedence, not
// the context-free metadata parser. Keep Point identity and raw intent untouched.
func TestAcquisitionModbusCanonicalAddressFunctionTable(t *testing.T) {
	cases := []struct {
		address, function, area, wireFunction string
		start, bits                           int
	}{
		{"100", "03", "HR", "03", 100, 16}, {"100", "04", "IR", "04", 100, 16},
		{"50000", "01", "COIL", "01", 50000, 1}, {"50000", "02", "DI", "02", 50000, 1},
		{"HR100", "01", "HR", "03", 100, 16}, {"ir100", "03", "IR", "04", 100, 16},
		{"CO100", "04", "COIL", "01", 100, 1}, {"C100", "03", "COIL", "01", 100, 1}, {"DI100", "03", "DI", "02", 100, 1},
		{"40100", "01", "HR", "03", 99, 16}, {"30100", "03", "IR", "04", 99, 16},
		{"10100", "03", "DI", "02", 99, 1}, {"00100", "03", "COIL", "01", 99, 1},
		{"100", "", "COIL", "01", 99, 1}, {"0", "", "HR", "03", 0, 16},
		{"40000", "", "HR", "03", 40000, 16}, {"65535", "03", "HR", "03", 65535, 16},
		{" 100 ", " 03 ", "HR", "03", 100, 16}, {"HR100", "unknown", "HR", "03", 100, 16},
		{"100", "holding", "HR", "03", 100, 16}, {"100", "input", "IR", "04", 100, 16},
		{"100", "coil", "COIL", "01", 100, 1}, {"100", "discrete", "DI", "02", 100, 1},
	}
	for _, tc := range cases {
		t.Run(tc.address+"/"+tc.function, func(t *testing.T) {
			in := acquisitionInput()
			p := &in.Sources[0].Points[0]
			p.Address = tc.address
			p.Function = tc.function
			p.NativeWidth = 1
			p.DataType = schema.DataTypeInt16
			if tc.bits == 1 {
				p.DataType = schema.DataTypeBool
			}
			original := *p
			plan := BuildAcquisitionPlan(in)
			if len(plan.Errors()) != 0 || len(plan.Requests()) != 1 {
				t.Fatalf("valid transport request rejected: %v", plan.Errors())
			}
			r := plan.Requests()[0]
			if r.Start != tc.start || r.Group.Area != tc.area || r.Group.Function != tc.wireFunction || r.Points[0].NativeBits != tc.bits {
				t.Fatalf("wire mismatch: %+v", r)
			}
			if r.Points[0].Point != original {
				t.Fatal("normalization rewrote canonical point identity/intent")
			}
		})
	}
	for _, tc := range []struct{ address, function string }{{"HR8.5", "03"}, {"COIL7", "01"}, {"DI1.0", "02"}, {"65536", "03"}, {"-1", "03"}, {"100", "FC03"}, {"100", "unknown"}, {"", "03"}} {
		t.Run("reject/"+tc.address+tc.function, func(t *testing.T) {
			in := acquisitionInput()
			in.Sources[0].Points[0].Address = tc.address
			in.Sources[0].Points[0].Function = tc.function
			plan := BuildAcquisitionPlan(in)
			if len(plan.Requests()) != 0 || len(plan.Errors()) != 1 || plan.Errors()[0].Code != "invalid-address" {
				t.Fatalf("unsupported wire form accepted: %+v %+v", plan.Requests(), plan.Errors())
			}
		})
	}
}
