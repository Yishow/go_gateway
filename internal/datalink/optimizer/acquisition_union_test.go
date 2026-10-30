package optimizer

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestAcquisitionPlanNestedUnion(t *testing.T) {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Coverage = ReadableCoverage{}
	s.Identity.Protocol = "mc"
	s.Points = []EffectivePoint{
		{ID: "wide", Address: "D0", DataType: schema.DataTypeInt64, NativeWidth: 4, Format: "ABCD"},
		{ID: "inner", Address: "D1", DataType: schema.DataTypeInt16, NativeWidth: 1, Format: "CDAB"},
	}
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 1 || r[0].Start != 0 || r[0].Length != 4 || len(r[0].Points) != 2 {
		t.Fatalf("nested union shrank/split: %+v", r)
	}
	slices := map[string]PointSlice{}
	for _, point := range r[0].Points {
		slices[point.Point.ID] = point
	}
	if slices["wide"].Offset != 0 || slices["wide"].Length != 4 || slices["inner"].Offset != 1 || slices["inner"].Length != 1 || slices["inner"].Point.Format != "CDAB" {
		t.Fatalf("lost decoder identity: %+v", slices)
	}
}

func TestAcquisitionPlanContiguousAcceptanceVectors(t *testing.T) {
	data, err := os.ReadFile("testdata/acquisition_acceptance.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures struct{ Cases []acceptanceCase }
	if err = json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	selected := map[string]bool{"nested-overlap": true, "adjacent": true, "gap-not-authorized": true, "gap-explicit-readable": true, "gap-forbidden": true, "limit-125": true, "limit-126": true, "multiword-boundary": true, "mixed-decoder-identities": true, "record-selection-is-not-source-set": true, "manual-exact-subset": true}
	for _, c := range fixtures.Cases {
		if !selected[c.ID] {
			continue
		}
		t.Run(c.ID, func(t *testing.T) {
			in := acquisitionInput()
			s := &in.Sources[0]
			s.Capability.ProtocolLimit = c.Limit
			s.Coverage = ReadableCoverage{Revision: c.CoverageRevision}
			s.Points = nil
			for _, r := range c.Readable {
				s.Coverage.Readable = append(s.Coverage.Readable, NativeRange{Area: "HR", Start: r[0], End: r[1]})
			}
			for _, r := range c.Forbidden {
				s.Coverage.Forbidden = append(s.Coverage.Forbidden, NativeRange{Area: "HR", Start: r[0], End: r[1]})
			}
			count := 0
			for _, p := range c.Points {
				repeat := p.Repeat
				if repeat == 0 {
					repeat = 1
				}
				for n := 0; n < repeat; n++ {
					id := p.ID
					if p.Repeat > 0 {
						id = fmt.Sprintf("%s%d", p.ID, n)
					}
					s.Points = append(s.Points, EffectivePoint{ID: id, Address: fmt.Sprintf("HR%d", p.Start+n*p.Stride), DataType: schema.DataType(p.Type), NativeWidth: p.Width, Format: p.Format})
					count++
				}
			}
			plan := BuildAcquisitionPlan(in)
			r := plan.Requests()
			if len(plan.Errors()) != 0 || len(r) != len(c.Requests) {
				t.Fatalf("requests %+v errors %+v expected %+v", r, plan.Errors(), c.Requests)
			}
			actualPoints := 0
			for n, want := range c.Requests {
				if r[n].Start != want.Start || r[n].Length != want.Length {
					t.Fatalf("request %+v expected %+v", r[n], want)
				}
				actualPoints += len(r[n].Points)
				for _, p := range r[n].Points {
					if p.Offset < 0 || p.Offset+p.Length > r[n].Length {
						t.Fatalf("partial scalar %+v in %+v", p, r[n])
					}
				}
			}
			if actualPoints != count {
				t.Fatalf("point identity lost: %d != %d", actualPoints, count)
			}
		})
	}
}

func TestAcquisitionPlanTwoHundredRegisters(t *testing.T) {
	in := acquisitionInput()
	s := &in.Sources[0]
	s.Points = nil
	for n := 0; n < 200; n++ {
		s.Points = append(s.Points, EffectivePoint{ID: fmt.Sprint(n), Address: fmt.Sprintf("HR%d", n), DataType: schema.DataTypeInt16, NativeWidth: 1})
	}
	p := BuildAcquisitionPlan(in)
	r := p.Requests()
	if len(r) != 2 || r[0].Start != 0 || r[0].Length != 125 || r[1].Start != 125 || r[1].Length != 75 {
		t.Fatalf("expected125+75: %+v", r)
	}
}
