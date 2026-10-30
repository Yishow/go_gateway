package optimizer

import (
	"cmp"
	"slices"
)

func canonicalAddressBounds(bounds []NativeAddressBound) []NativeAddressBound {
	result := append([]NativeAddressBound(nil), bounds...)
	slices.SortFunc(result, func(a, b NativeAddressBound) int { return cmp.Compare(canonicalJSON(a), canonicalJSON(b)) })
	return result
}

// declaredSpanError intersects every provided authority without treating a
// missing declaration as proof that the real device has unlimited memory.
// One bound per area per authority keeps intersection semantics unambiguous.
func declaredSpanError(protocol string, cap AcquisitionCapability, area string, start, width int) string {
	outside := false
	for _, bounds := range [][]NativeAddressBound{cap.ProtocolAddressBounds, cap.ClientAddressBounds, cap.DeviceAddressBounds} {
		areas := make(map[string]bool, len(bounds))
		for _, bound := range bounds {
			if !canonicalNativeArea(protocol, bound.Area) || bound.StartInclusive < 0 || bound.EndExclusive <= bound.StartInclusive || areas[bound.Area] {
				return "invalid-address-bounds"
			}
			areas[bound.Area] = true
			if bound.Area == area && (start < bound.StartInclusive || start+width > bound.EndExclusive) {
				outside = true
			}
		}
	}
	if outside {
		return "address-out-of-bounds"
	}
	return ""
}
