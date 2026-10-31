package optimizer

import (
	"cmp"
	"go-gateway/internal/protocol/fatek"
	"slices"
)

func randomEligible(identity AcquisitionIdentity, cap AcquisitionCapability, bits int) bool {
	if cap.RandomSupport != "supported" {
		return false
	}
	switch identity.Protocol {
	case "mc", "mitsubishi":
		return identity.Frame == "3e-binary" && bits == 16
	case "fatek":
		return identity.Frame == "ascii"
	}
	return false
}

func randomLimit(identity AcquisitionIdentity, cap AcquisitionCapability) int {
	limit := 64
	if identity.Protocol == "mc" || identity.Protocol == "mitsubishi" {
		limit = 255
	}
	for _, n := range []int{cap.ProtocolLimit, cap.ClientLimit, cap.DeviceLimit} {
		if n > 0 {
			limit = min(limit, n)
		}
	}
	return limit
}

func randomFrameBytes(identity AcquisitionIdentity, units []NativeUnit) (ApplicationFrameBudget, bool) {
	switch identity.Protocol {
	case "mc", "mitsubishi":
		if identity.Frame != "3e-binary" {
			return ApplicationFrameBudget{}, false
		}
		return ApplicationFrameBudget{RequestBytes: 17 + 4*len(units), ResponseBytes: 11 + 2*len(units)}, true
	case "fatek":
		if identity.Frame != "ascii" {
			return ApplicationFrameBudget{}, false
		}
		result := ApplicationFrameBudget{RequestBytes: 10, ResponseBytes: 9}
		for _, u := range units {
			component, err := fatek.GetComponentType(u.Area)
			if err != nil {
				return ApplicationFrameBudget{}, false
			}
			address, err := fatek.FormatAddress(component, u.Address)
			if err != nil {
				return ApplicationFrameBudget{}, false
			}
			result.RequestBytes += len(address)
			result.ResponseBytes += max(1, u.Bits/4)
		}
		return result, true
	}
	return ApplicationFrameBudget{}, false
}

// packRandom adds complete scalars, deduplicating only identical native units.
// A request has at most 255 units, bounding each candidate copy/lookup.
func packRandom(scalars []AcquisitionRequest, cap AcquisitionCapability) ([]AcquisitionRequest, bool) {
	var result []AcquisitionRequest
	var current *AcquisitionRequest
	budget, valid := strictFrameBudget(cap)
	if !valid {
		return nil, false
	}
	for _, scalar := range scalars {
		if !randomEligible(scalar.Group.Identity, cap, scalar.Points[0].NativeBits) {
			return nil, false
		}
		candidate := appendRandomScalar(current, scalar, cap)
		if !randomWithinLimits(candidate) || !fitsFrameBudget(candidate.FrameBytes, budget) {
			if current == nil {
				return nil, false
			}
			result = append(result, *current)
			candidate = appendRandomScalar(nil, scalar, cap)
			if !randomWithinLimits(candidate) || !fitsFrameBudget(candidate.FrameBytes, budget) {
				return nil, false
			}
		}
		current = &candidate
	}
	if current != nil {
		result = append(result, *current)
	}
	return result, true
}

func appendRandomScalar(current *AcquisitionRequest, scalar AcquisitionRequest, cap AcquisitionCapability) AcquisitionRequest {
	var candidate AcquisitionRequest
	if current == nil {
		candidate = scalar
		candidate.Kind = "random"
		candidate.AddressStep = 0
		candidate.Points = nil
		candidate.Units = nil
		candidate.Limit = randomLimit(scalar.Group.Identity, cap)
		candidate.FallbackReason = ""
	} else {
		candidate = *current
		candidate.Points = append([]PointSlice(nil), current.Points...)
		candidate.Units = append([]NativeUnit(nil), current.Units...)
		if candidate.Group.Area != scalar.Group.Area {
			candidate.Group.Area = ""
		}
	}
	point := scalar.Points[0]
	point.Offset = -1
	for n := range scalar.Length {
		address := scalar.Start + n*scalar.AddressStep
		unit := NativeUnit{Area: scalar.Group.Area, Address: address, Bits: point.NativeBits}
		index := slices.Index(candidate.Units, unit)
		if index < 0 {
			index = len(candidate.Units)
			candidate.Units = append(candidate.Units, unit)
		}
		if point.Offset < 0 {
			point.Offset = index
		}
	}
	candidate.Points = append(candidate.Points, point)
	candidate.Length = len(candidate.Units)
	candidate.FrameBytes, candidate.FrameCostKnown = randomFrameBytes(candidate.Group.Identity, candidate.Units)
	return candidate
}

func sortScalarRequests(requests []AcquisitionRequest) {
	slices.SortFunc(requests, func(a, b AcquisitionRequest) int {
		if n := cmp.Compare(a.Group.Area, b.Group.Area); n != 0 {
			return n
		}
		if n := cmp.Compare(a.Start%a.AddressStep, b.Start%b.AddressStep); n != 0 {
			return n
		}
		if n := cmp.Compare(a.Start, b.Start); n != 0 {
			return n
		}
		if n := cmp.Compare(b.Length, a.Length); n != 0 {
			return n
		}
		return cmp.Compare(canonicalJSON(a), canonicalJSON(b))
	})
}

// Component caps are independent of FATEK's protocol word-cost budget.
func randomWithinLimits(r AcquisitionRequest) bool {
	if len(r.Units) > r.Limit {
		return false
	}
	if r.Group.Identity.Protocol == "fatek" {
		cost := 0
		for _, u := range r.Units {
			cost += max(1, u.Bits/16)
		}
		return cost <= 64
	}
	return true
}
