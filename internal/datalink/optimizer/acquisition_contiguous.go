package optimizer

import (
	"cmp"
	"slices"
)

type contiguousCandidates struct {
	coverage ReadableCoverage
	requests []AcquisitionRequest
}

// appendContiguous preserves complete scalars. A limit boundary may repeat raw
// units in separate requests, but never divides one point's observation.
func appendContiguous(out []AcquisitionRequest, candidates contiguousCandidates) []AcquisitionRequest {
	requests := candidates.requests
	slices.SortFunc(requests, func(a, b AcquisitionRequest) int {
		if n := cmp.Compare(a.Start, b.Start); n != 0 {
			return n
		}
		if n := cmp.Compare(b.Length, a.Length); n != 0 {
			return n
		}
		return cmp.Compare(canonicalJSON(a), canonicalJSON(b))
	})
	var current *AcquisitionRequest
	for _, next := range requests {
		if current == nil {
			copy := next
			current = &copy
			continue
		}
		end := current.Start + current.Length
		unionEnd := max(end, next.Start+next.Length)
		if unionEnd-current.Start <= current.Limit && readableGap(candidates.coverage, current.Group.Area, end, next.Start) {
			for _, point := range next.Points {
				point.Offset += next.Start - current.Start
				current.Points = append(current.Points, point)
			}
			current.Length = unionEnd - current.Start
			current.FrameBytes, current.FrameCostKnown = contiguousFrameBytes(current.Group.Identity, current.Group.Area, current.Start, current.Length, current.Points[0].NativeBits)
			if current.FallbackReason == "scalar-plan" {
				current.FallbackReason = "safe-contiguous-plan"
			}
			continue
		}
		out = append(out, *current)
		copy := next
		current = &copy
	}
	if current != nil {
		out = append(out, *current)
	}
	return out
}

// readableGap checks an interval union, not individual cells. Active spans need
// no gap grant; unoccupied gaps require a versioned explicit readable statement.
func readableGap(coverage ReadableCoverage, area string, start, end int) bool {
	if start >= end {
		return true
	}
	if coverage.Revision == "" || intersectsRanges(coverage.Forbidden, area, start, end) || intersectsRanges(coverage.Excluded, area, start, end) {
		return false
	}
	ranges := append([]NativeRange(nil), coverage.Readable...)
	slices.SortFunc(ranges, func(a, b NativeRange) int { return cmp.Compare(a.Start, b.Start) })
	covered := start
	for _, r := range ranges {
		if r.Area != area || r.Start > covered || r.End <= covered {
			continue
		}
		covered = r.End
		if covered >= end {
			return true
		}
	}
	return false
}

func intersectsRanges(ranges []NativeRange, area string, start, end int) bool {
	for _, r := range ranges {
		if r.Area == area && r.Start < end && r.End > start && r.Start < r.End {
			return true
		}
	}
	return false
}
