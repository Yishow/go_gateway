package optimizer

import (
	"fmt"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestAcquisitionRandomLimitsCountComponentsNotPoints(t *testing.T) {
	for _, c := range []struct {
		protocol, area            string
		count, width, bits, limit int
	}{{"mc", "D", 64, 4, 16, 255}, {"fatek", "DR", 33, 2, 32, 64}} {
		in := randomInput(c.protocol)
		for i := 0; i < c.count; i++ {
			in.Sources[0].Points = append(in.Sources[0].Points, EffectivePoint{ID: fmt.Sprint(i), Address: fmt.Sprintf("%s%d", c.area, i*10), DataType: schema.DataTypeInt64, NativeWidth: c.width})
		}
		plan := BuildAcquisitionPlan(in)
		r := plan.Requests()
		if len(plan.Errors()) != 0 || len(r) != 2 {
			t.Fatalf("%s native capacity: %+v errors%+v", c.protocol, r, plan.Errors())
		}
		seen := map[string]bool{}
		for _, request := range r {
			if request.Kind != "random" || len(request.Units) > c.limit {
				t.Fatalf("component overflow: %+v", request)
			}
			for _, point := range request.Points {
				if seen[point.Point.ID] || point.Length != c.width || point.Offset+point.Length > len(request.Units) {
					t.Fatal("split or duplicate scalar")
				}
				seen[point.Point.ID] = true
				_, start, _, _, _ := nativeAddress(c.protocol, point.Point.Address)
				for n := 0; n < c.width; n++ {
					u := request.Units[point.Offset+n]
					if u.Address != start+n || u.Area != c.area || u.Bits != c.bits {
						t.Fatalf("invalid scatter: %+v", u)
					}
				}
			}
		}
		if len(seen) != c.count {
			t.Fatal("lost scalar")
		}
	}
}

func TestAcquisitionRandomForbiddenAndBudgetNeverBroadenRead(t *testing.T) {
	in := randomInput("mc", 0, 10, 20)
	in.Sources[0].Coverage = ReadableCoverage{Revision: "cov", Forbidden: []NativeRange{{Area: "D", Start: 5, End: 10}}}
	plan := BuildAcquisitionPlan(in)
	if len(plan.Requests()) != 1 || plan.Requests()[0].Kind != "random" {
		t.Fatal("sparse random should skip forbidden gap")
	}
	for _, unit := range plan.Requests()[0].Units {
		if unit.Address >= 5 && unit.Address < 10 {
			t.Fatal("random reads forbidden gap")
		}
	}
	in.Sources[0].Coverage.Forbidden[0].End = 11
	plan = BuildAcquisitionPlan(in)
	if len(plan.Errors()) != 1 || plan.Errors()[0].PointID != "1" || plan.Errors()[0].Code != "forbidden-address" {
		t.Fatalf("forbidden scalar not rejected: %+v", plan.Errors())
	}
	for _, request := range plan.Requests() {
		for _, point := range request.Points {
			if point.Point.ID == "1" {
				t.Fatal("forbidden scalar reintroduced by random")
			}
		}
	}
}
