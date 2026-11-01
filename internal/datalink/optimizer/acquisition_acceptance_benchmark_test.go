package optimizer

import (
	"fmt"
	"testing"
)

// BenchmarkAcquisitionPlan measures the canonical production planner, including
// fingerprinting and safe request selection, not the legacy BlockMerger.
// Fixture construction and validation are deliberately outside the timer.
func BenchmarkAcquisitionPlan(b *testing.B) {
	for _, count := range []int{100, 1000, 10000} {
		for _, stride := range []int{1, 4} {
			b.Run(fmt.Sprintf("points=%d/stride=%d", count, stride), func(b *testing.B) {
				input := acquisitionInput()
				source := &input.Sources[0]
				source.Coverage = ReadableCoverage{Revision: "acceptance-no-gap"}
				point := source.Points[0]
				source.Points = make([]EffectivePoint, count)
				for i := range count {
					p := point
					p.ID = fmt.Sprintf("point-%05d", i)
					// Keep the largest sparse case within Modbus address bounds.
					p.Address = fmt.Sprintf("HR%d", i*stride)
					p.SourceIdentity = p.ID
					source.Points[i] = p
				}
				plan := BuildAcquisitionPlan(input)
				if len(plan.Errors()) != 0 || len(plan.Requests()) == 0 {
					b.Fatalf("invalid benchmark fixture: %+v", plan.Errors())
				}
				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					plan = BuildAcquisitionPlan(input)
				}
				b.StopTimer()
				b.ReportMetric(float64(len(plan.Requests())), "requests/plan")
			})
		}
	}
}
