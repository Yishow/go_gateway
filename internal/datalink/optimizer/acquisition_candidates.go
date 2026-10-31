package optimizer

import (
	"cmp"
	"slices"
)

// chooseAcquisitionRequests compares at most six candidates, never exponential
// partitions. Random uses original complete scalars, not pre-fixed dense runs.
func chooseAcquisitionRequests(group contiguousCandidates) []AcquisitionRequest {
	scalars := append([]AcquisitionRequest(nil), group.requests...)
	sortScalarRequests(scalars)
	var contiguous []AcquisitionRequest
	for start := 0; start < len(scalars); {
		end := start + 1
		for end < len(scalars) && scalars[end].Group.Area == scalars[start].Group.Area && scalars[end].AddressStep == scalars[start].AddressStep && scalars[end].Start%scalars[end].AddressStep == scalars[start].Start%scalars[start].AddressStep {
			end++
		}
		contiguous = appendContiguous(contiguous, contiguousCandidates{coverage: group.coverage, requests: append([]AcquisitionRequest(nil), scalars[start:end]...)})
		start = end
	}
	best := contiguous
	random, valid := packRandom(scalars, group.capability)
	if !valid {
		return best
	}
	best = betterRequests(best, random)
	// Four fixed thresholds retain long, cheap dense reads while packing remaining
	// sparse scalars together. Points are assigned once; raw overlap may repeat
	// across requests, but no individual scalar is divided.
	for _, threshold := range []int{2, 4, 8, 16} {
		var fixed []AcquisitionRequest
		selected := make(map[string]bool)
		for _, request := range contiguous {
			if request.Length < threshold {
				continue
			}
			fixed = append(fixed, request)
			for _, p := range request.Points {
				selected[canonicalJSON(p.Point)] = true
			}
		}
		if len(fixed) == 0 || len(fixed) == len(contiguous) {
			continue
		}
		var remainder []AcquisitionRequest
		for _, scalar := range scalars {
			if !selected[canonicalJSON(scalar.Points[0].Point)] {
				remainder = append(remainder, scalar)
			}
		}
		packed, ok := packRandom(remainder, group.capability)
		if ok {
			best = betterRequests(best, append(fixed, packed...))
		}
	}
	return best
}

func betterRequests(a, b []AcquisitionRequest) []AcquisitionRequest {
	if len(a) < len(b) {
		return a
	}
	if len(b) < len(a) {
		return b
	}
	cost := func(requests []AcquisitionRequest) (int, bool) {
		total := 0
		for _, request := range requests {
			if !request.FrameCostKnown {
				return 0, false
			}
			total += request.FrameBytes.RequestBytes + request.FrameBytes.ResponseBytes
		}
		return total, true
	}
	ac, ak := cost(a)
	bc, bk := cost(b)
	if !ak || !bk {
		return a
	} // Unknown cost never establishes an improvement.
	if ac < bc {
		return a
	}
	if bc < ac {
		return b
	}
	stable := func(requests []AcquisitionRequest) string {
		copy := append([]AcquisitionRequest(nil), requests...)
		slices.SortFunc(copy, func(a, b AcquisitionRequest) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
		return canonicalJSON(copy)
	}
	if stable(b) < stable(a) {
		return b
	}
	return a
}
