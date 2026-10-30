package optimizer

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
)

// BuildAcquisitionPlan plans without device I/O or runtime mutation. The caller
// must not mutate input concurrently with this call; subsequent mutation is safe.
// Compatible complete scalars share safe contiguous reads; uncovered gaps and
// native limits remain hard boundaries. No runtime I/O or mutable connection state is used.
func BuildAcquisitionPlan(input AcquisitionInput) AcquisitionPlan {
	canonical := canonicalAcquisitionInput(input)
	digest := sha256.Sum256([]byte("acquisition-plan-v4\n" + canonicalJSON(canonical)))
	plan := AcquisitionPlan{fingerprint: hex.EncodeToString(digest[:])}
	groups := make(map[string]*contiguousCandidates)
	for _, source := range canonical.Sources {
		for _, point := range source.Points {
			request, code := planScalar(source, point)
			if code != "" {
				plan.errors = append(plan.errors, AcquisitionPointError{Identity: source.Identity, PointID: point.ID, Code: code})
				continue
			}
			// Include complete capability/coverage declarations as well as revisions:
			// a malformed caller reusing a revision must not merge conflicting scopes.
			groupIdentity := request.Group
			if source.Capability.RandomCrossArea && randomEligible(source.Identity, source.Capability, request.Points[0].NativeBits) {
				groupIdentity.Area = ""
			}
			key := canonicalJSON(struct {
				Group      AcquisitionGroup
				Capability AcquisitionCapability
				Coverage   ReadableCoverage
			}{groupIdentity, source.Capability, source.Coverage})
			group := groups[key]
			if group == nil {
				group = &contiguousCandidates{coverage: source.Coverage, capability: source.Capability}
				groups[key] = group
			}
			group.requests = append(group.requests, request)
		}
	}
	for _, group := range groups {
		plan.requests = append(plan.requests, chooseAcquisitionRequests(*group)...)
	}
	slices.SortFunc(plan.requests, func(a, b AcquisitionRequest) int {
		if n := cmp.Compare(canonicalJSON(a.Group), canonicalJSON(b.Group)); n != 0 {
			return n
		}
		if n := cmp.Compare(a.Start, b.Start); n != 0 {
			return n
		}
		return cmp.Compare(canonicalJSON(a), canonicalJSON(b))
	})
	slices.SortFunc(plan.errors, func(a, b AcquisitionPointError) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
	return plan
}

// canonicalJSON only receives the closed, JSON-safe value types in this file.
// No user-defined marshalers, maps, floating-point values or pointers are accepted.
func canonicalJSON[T any](value T) string {
	data, err := json.Marshal(value)
	if err != nil {
		panic("optimizer: non-serializable internal acquisition contract")
	}
	return string(data)
}

func canonicalAcquisitionInput(input AcquisitionInput) AcquisitionInput {
	result := input
	result.Sources = append([]AcquisitionSource(nil), input.Sources...)
	for i := range result.Sources {
		s := &result.Sources[i]
		s.Capability.ProtocolAddressBounds = canonicalAddressBounds(s.Capability.ProtocolAddressBounds)
		s.Capability.ClientAddressBounds = canonicalAddressBounds(s.Capability.ClientAddressBounds)
		s.Capability.DeviceAddressBounds = canonicalAddressBounds(s.Capability.DeviceAddressBounds)
		s.Points = append([]EffectivePoint(nil), s.Points...)
		slices.SortFunc(s.Points, func(a, b EffectivePoint) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
		s.Coverage.Readable = canonicalRanges(s.Coverage.Readable)
		s.Coverage.Forbidden = canonicalRanges(s.Coverage.Forbidden)
		s.Coverage.Excluded = canonicalRanges(s.Coverage.Excluded)
	}
	slices.SortFunc(result.Sources, func(a, b AcquisitionSource) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
	return result
}

func canonicalRanges(ranges []NativeRange) []NativeRange {
	result := append([]NativeRange(nil), ranges...)
	slices.SortFunc(result, func(a, b NativeRange) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
	return result
}

func planScalar(source AcquisitionSource, point EffectivePoint) (AcquisitionRequest, string) {
	area, start, nativeBits, bitIndex, function, code := pointNativeAddress(source.Identity.Protocol, point)
	if code != "" {
		return AcquisitionRequest{}, code
	}
	if !validCoverage(source.Identity.Protocol, source.Coverage) {
		return AcquisitionRequest{}, "invalid-coverage"
	}
	bits := scalarBits(point.DataType)
	if bits == 0 || ((nativeBits == 1 || bitIndex >= 0) && point.DataType != schema.DataTypeBool) {
		return AcquisitionRequest{}, "unsupported-type"
	}
	width := (bits + nativeBits - 1) / nativeBits
	if point.NativeWidth != width {
		return AcquisitionRequest{}, "invalid-native-width"
	}
	if !validNativeSpan(source.Identity.Protocol, area, start, width) {
		return AcquisitionRequest{}, "invalid-address"
	}
	if code := declaredSpanError(source.Identity.Protocol, source.Capability, area, start, width); code != "" {
		return AcquisitionRequest{}, code
	}
	if intersectsRanges(source.Coverage.Forbidden, area, start, start+width) {
		return AcquisitionRequest{}, "forbidden-address"
	}
	cap := source.Capability
	limit := 0
	for _, n := range []int{cap.ProtocolLimit, cap.ClientLimit, cap.DeviceLimit} {
		if n < 0 {
			return AcquisitionRequest{}, "invalid-capability"
		}
		if n > 0 && (limit == 0 || n < limit) {
			limit = n
		}
	}
	if limit == 0 {
		return AcquisitionRequest{}, "invalid-capability"
	}
	limit = min(limit, nativeClientLimit(source.Identity.Protocol, nativeBits))
	if width > limit {
		return AcquisitionRequest{}, "scalar-exceeds-limit"
	}
	limit, frameBytes, frameKnown, code := frameLimitedUnits(source, area, start, width, nativeBits, limit)
	if code != "" {
		return AcquisitionRequest{}, code
	}
	reason := "scalar-plan"
	if (source.Identity.Protocol == "mc" || source.Identity.Protocol == "mitsubishi" || source.Identity.Protocol == "fatek") && (cap.RandomSupport == "" || cap.RandomSupport == "unknown") {
		reason = "random-capability-unknown"
	}
	return AcquisitionRequest{
		Group: AcquisitionGroup{Identity: source.Identity, CapabilityRevision: cap.Revision, CoverageRevision: source.Coverage.Revision, Function: function, Area: area},
		Kind:  "contiguous", Start: start, Length: width, Limit: limit, FallbackReason: reason,
		FrameBytes: frameBytes, FrameCostKnown: frameKnown,
		Points: []PointSlice{{Point: point, Offset: 0, Length: width, NativeBits: nativeBits, BitIndex: bitIndex}},
	}, ""
}

func scalarBits(dt schema.DataType) int {
	switch dt {
	case schema.DataTypeBool:
		return 1
	case schema.DataTypeInt16, schema.DataTypeUint16:
		return 16
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 32
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 64
	default:
		return 0 // String has no immutable explicit length contract yet.
	}
}

func nativeAddress(protocol, address string) (area string, start, nativeBits, bitIndex int, code string) {
	nativeBits, bitIndex = 16, -1
	if protocol == "fatek" {
		symbol, offset, err := fatek.ParseComponentAddress(address)
		if err != nil {
			return "", 0, 0, 0, "invalid-address"
		}
		component, err := fatek.GetComponentType(symbol)
		if err != nil {
			return "", 0, 0, 0, "invalid-address"
		}
		return symbol, offset, component.Width, -1, ""
	}
	var parserProtocol hsllogic.ProtocolType
	switch protocol {
	case "mc", "mitsubishi":
		parserProtocol = hsllogic.ProtocolMitsubishi
	case "modbus":
		parserProtocol = hsllogic.ProtocolModbus
	default:
		return "", 0, 0, 0, "unsupported-protocol"
	}
	parsed, err := hsllogic.ParseAddress(parserProtocol, address)
	if err != nil {
		return "", 0, 0, 0, "invalid-address"
	}
	// A coil/discrete input is already one native bit. The shared parser also
	// accepts dotted syntax there; reject it rather than inventing a word slice.
	if parserProtocol == hsllogic.ProtocolModbus && parsed.BitIndex >= 0 &&
		(parsed.DeviceType == "COIL" || parsed.DeviceType == "DI") {
		return "", 0, 0, 0, "invalid-address"
	}
	if parserProtocol == hsllogic.ProtocolMitsubishi {
		device, err := mcprotocol.GetDeviceType(parsed.DeviceType)
		if err != nil {
			return "", 0, 0, 0, "unsupported-area"
		}
		parsed.IsBitDevice = device.IsBit
	}
	if parsed.IsBitDevice && parsed.BitIndex < 0 {
		nativeBits = 1
	}
	return parsed.DeviceType, parsed.Offset, nativeBits, parsed.BitIndex, ""
}
