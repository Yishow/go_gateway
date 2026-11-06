package optimizer

import (
	"cmp"
	"slices"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
)

func randomEligible(identity AcquisitionIdentity, capability AcquisitionCapability, bits int) bool {
	if capability.RandomSupport != "supported" {
		return false
	}
	switch identity.Protocol {
	case "mc", acquisitionProtocolMitsubishi:
		return identity.Frame == acquisitionFrame3EBinary && (bits == 1 || bits == 16)
	case acquisitionProtocolFatek:
		return identity.Frame == acquisitionFrameASCII
	}
	return false
}

func randomScalarEligible(identity AcquisitionIdentity, capability AcquisitionCapability, coverage ReadableCoverage, scalar AcquisitionRequest) bool {
	if len(scalar.Points) == 0 {
		return false
	}
	p := scalar.Points[0]
	if p.Point.DataType == schema.DataTypeString {
		return false
	}
	if identity.Protocol == "mc" || identity.Protocol == acquisitionProtocolMitsubishi {
		if mcprotocol.IsLongTimerOrCounter(scalar.Group.Area) {
			return false
		}
	}
	if slices.Contains(capability.RandomExcludedAreas, scalar.Group.Area) {
		return false
	}
	return randomEligible(identity, capability, p.NativeBits) && safeMCRandomBitUnit(identity, capability, coverage, scalar)
}

func safeMCRandomBitUnit(identity AcquisitionIdentity, capability AcquisitionCapability, coverage ReadableCoverage, scalar AcquisitionRequest) bool {
	if (identity.Protocol != "mc" && identity.Protocol != acquisitionProtocolMitsubishi) || scalar.Points[0].NativeBits != 1 {
		return true
	}
	start := scalar.Start / 16 * 16
	area := scalar.Group.Area
	if unverifiedQIQRXYExclusions(identity, coverage, area) {
		return false
	}
	return validNativeSpan(identity.Protocol, area, start, 16) &&
		declaredSpanError(identity.Protocol, capability, area, start, 16) == "" &&
		!intersectsRanges(coverage.Forbidden, area, start, start+16) &&
		!intersectsRanges(coverage.Excluded, area, start, start+16)
}

func randomLimit(identity AcquisitionIdentity, capability AcquisitionCapability) int {
	limit := 64
	if identity.Protocol == "mc" || identity.Protocol == acquisitionProtocolMitsubishi {
		limit = 255
	}
	for _, n := range []int{capability.ProtocolLimit, capability.RandomProtocolLimit, capability.ClientLimit, capability.DeviceLimit} {
		if n > 0 {
			limit = min(limit, n)
		}
	}
	return limit
}

func randomFrameBytes(identity AcquisitionIdentity, units []NativeUnit) (ApplicationFrameBudget, bool) {
	switch identity.Protocol {
	case "mc", acquisitionProtocolMitsubishi:
		if identity.Frame != acquisitionFrame3EBinary {
			return ApplicationFrameBudget{}, false
		}
		responsePayload := 0
		for _, u := range units {
			if u.Bits == 32 {
				responsePayload += 4
			} else {
				responsePayload += 2
			}
		}
		return ApplicationFrameBudget{RequestBytes: 17 + 4*len(units), ResponseBytes: 11 + responsePayload}, true
	case acquisitionProtocolFatek:
		if identity.Frame != acquisitionFrameASCII {
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
func packRandom(scalars []AcquisitionRequest, capability AcquisitionCapability, coverage ReadableCoverage) ([]AcquisitionRequest, bool) {
	var result []AcquisitionRequest
	var current *AcquisitionRequest
	budget, valid := strictFrameBudget(capability)
	if !valid {
		return nil, false
	}
	for _, scalar := range scalars {
		if !randomScalarEligible(scalar.Group.Identity, capability, coverage, scalar) {
			return nil, false
		}
		candidate := appendRandomScalar(current, scalar, capability)
		if !randomWithinLimits(candidate, capability) || !fitsFrameBudget(candidate.FrameBytes, budget) {
			if current == nil {
				return nil, false
			}
			result = append(result, *current)
			candidate = appendRandomScalar(nil, scalar, capability)
			if !randomWithinLimits(candidate, capability) || !fitsFrameBudget(candidate.FrameBytes, budget) {
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

func appendRandomScalar(current *AcquisitionRequest, scalar AcquisitionRequest, capability AcquisitionCapability) AcquisitionRequest {
	var candidate AcquisitionRequest
	if current == nil {
		candidate = scalar
		candidate.Kind = acquisitionKindRandom
		candidate.AddressStep = 0
		candidate.Points = nil
		candidate.Units = nil
		candidate.Limit = randomLimit(scalar.Group.Identity, capability)
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
	isMC := scalar.Group.Identity.Protocol == "mc" || scalar.Group.Identity.Protocol == acquisitionProtocolMitsubishi
	if isMC && point.NativeBits == 1 {
		wordAddr := (scalar.Start / 16) * 16
		unit := NativeUnit{Area: scalar.Group.Area, Address: wordAddr, Bits: 16}
		point.Offset = insertMCUnit(&candidate, unit)
		point.Length = 1
		point.NativeBits = 16
		point.BitIndex = scalar.Start % 16
	} else if isMC && point.NativeBits == 16 && scalarBits(point.Point.DataType) == 32 {
		unit := NativeUnit{Area: scalar.Group.Area, Address: scalar.Start, Bits: 32}
		point.Offset = insertMCUnit(&candidate, unit)
		point.Length = 1
		point.NativeBits = 32
	} else {
		for n := range scalar.Length {
			address := scalar.Start + n*scalar.AddressStep
			unit := NativeUnit{Area: scalar.Group.Area, Address: address, Bits: point.NativeBits}
			var index int
			if isMC {
				index = insertMCUnit(&candidate, unit)
			} else {
				index = slices.Index(candidate.Units, unit)
				if index < 0 {
					index = len(candidate.Units)
					candidate.Units = append(candidate.Units, unit)
				}
			}
			if point.Offset < 0 {
				point.Offset = index
			}
		}
	}
	candidate.Points = append(candidate.Points, point)
	candidate.Length = len(candidate.Units)
	candidate.FrameBytes, candidate.FrameCostKnown = randomFrameBytes(candidate.Group.Identity, candidate.Units)
	return candidate
}

func insertMCUnit(candidate *AcquisitionRequest, unit NativeUnit) int {
	if index := slices.Index(candidate.Units, unit); index >= 0 {
		return index
	}
	if unit.Bits == 32 {
		candidate.Units = append(candidate.Units, unit)
		return len(candidate.Units) - 1
	}
	insertAt := len(candidate.Units)
	for i, existing := range candidate.Units {
		if existing.Bits == 32 {
			insertAt = i
			break
		}
	}
	candidate.Units = slices.Insert(candidate.Units, insertAt, unit)
	for i := range candidate.Points {
		if candidate.Points[i].Offset >= insertAt {
			candidate.Points[i].Offset++
		}
	}
	return insertAt
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
func randomWithinLimits(r AcquisitionRequest, capability AcquisitionCapability) bool {
	if len(r.Units) > r.Limit {
		return false
	}
	if r.Group.Identity.Protocol == acquisitionProtocolFatek {
		cost := 0
		for _, u := range r.Units {
			cost += max(1, u.Bits/16)
		}
		wordLimit := 64
		if capability.RandomProtocolLimit > 0 {
			wordLimit = min(wordLimit, capability.RandomProtocolLimit)
		}
		return cost <= wordLimit
	}
	return true
}
