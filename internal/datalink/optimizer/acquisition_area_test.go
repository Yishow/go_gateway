package optimizer

import "testing"

func TestAcquisitionPlanInvalidRestrictionAreas(t *testing.T) {
	for _, area := range []string{"", "hr", "TYPO", "HR ", "ZZ"} {
		for _, kind := range []string{"readable", "forbidden", "excluded", "bound"} {
			t.Run(kind+"/"+area, func(t *testing.T) {
				in := sparseAcquisitionInput()
				s := &in.Sources[0]
				r := NativeRange{Area: area, Start: 5, End: 6}
				code := "invalid-coverage"
				switch kind {
				case "readable":
					s.Coverage.Readable = append(s.Coverage.Readable, r)
				case "forbidden":
					s.Coverage.Forbidden = []NativeRange{r}
				case "excluded":
					s.Coverage.Excluded = []NativeRange{r}
				case "bound":
					s.Capability.DeviceAddressBounds = []NativeAddressBound{{Area: area, StartInclusive: 0, EndExclusive: 1}}
					code = "invalid-address-bounds"
				}
				p := BuildAcquisitionPlan(in)
				if len(p.Requests()) != 0 || len(p.Errors()) != 2 || p.Errors()[0].Code != code {
					t.Fatalf("malformed restriction silently ignored: requests%d errors%+v", len(p.Requests()), p.Errors())
				}
			})
		}
	}
}

func TestAcquisitionPlanProtocolRestrictionAreas(t *testing.T) {
	for _, c := range []struct {
		protocol, point, area string
		valid                 bool
	}{
		{"modbus", "HR0", "IR", true}, {"modbus", "HR0", "COIL", true}, {"modbus", "HR0", "D", false},
		{"mc", "D0", "W", true}, {"mc", "D0", "TC", true}, {"mc", "D0", "SM", false}, {"mc", "D0", "d", false},
		{"fatek", "D0", "DR", true}, {"fatek", "D0", "WM", true}, {"fatek", "D0", "TYPO", false}, {"fatek", "D0", "dr", false},
	} {
		for _, kind := range []string{"coverage", "bounds"} {
			t.Run(c.protocol+c.area+kind, func(t *testing.T) {
				in := acquisitionInput()
				s := &in.Sources[0]
				s.Identity.Protocol = c.protocol
				s.Points[0].Address = c.point
				s.Coverage = ReadableCoverage{}
				code := "invalid-coverage"
				if kind == "coverage" {
					s.Coverage.Forbidden = []NativeRange{{Area: c.area, Start: 5, End: 6}}
				} else {
					s.Capability.ClientAddressBounds = []NativeAddressBound{{Area: c.area, StartInclusive: 0, EndExclusive: 1}}
					code = "invalid-address-bounds"
				}
				p := BuildAcquisitionPlan(in)
				if c.valid {
					if len(p.Requests()) != 1 || len(p.Errors()) != 0 {
						t.Fatalf("valid independent area rejected: %+v", p.Errors())
					}
				} else if len(p.Requests()) != 0 || len(p.Errors()) != 1 || p.Errors()[0].Code != code {
					t.Fatalf("invalid protocol area ignored: %+v %+v", p.Requests(), p.Errors())
				}
			})
		}
	}
}
