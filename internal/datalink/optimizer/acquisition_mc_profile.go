package optimizer

import (
	"slices"
	"strconv"
	"strings"
)

const mcQIQRProfile = "mitsubishi-q-iqr-3e-binary-v1"

// Source-rule exclusions can still use legacy octal X/Y geometry. Without
// profile-bound geometry, their positions cannot authorize additional reads.
func unverifiedQIQRXYExclusions(identity AcquisitionIdentity, coverage ReadableCoverage, area string) bool {
	return (identity.Protocol == "mc" || identity.Protocol == acquisitionProtocolMitsubishi) &&
		identity.ProtocolRevision == mcQIQRProfile && (area == "X" || area == "Y") &&
		slices.ContainsFunc(coverage.Excluded, func(r NativeRange) bool { return r.Area == area })
}

// pointAddressForIdentity keeps historical MC parsing unless the declared
// Q/iQ-R profile supplies the hexadecimal X/Y address convention.
func pointAddressForIdentity(identity AcquisitionIdentity, point EffectivePoint) nativePointAddress {
	if (identity.Protocol == "mc" || identity.Protocol == acquisitionProtocolMitsubishi) && identity.ProtocolRevision == mcQIQRProfile {
		address := strings.ToUpper(strings.TrimSpace(point.Address))
		if address != "" && (address[0] == 'X' || address[0] == 'Y') {
			digits := address[1:]
			if digits == "" {
				return nativePointAddress{code: acquisitionReasonInvalidAddress}
			}
			for _, digit := range digits {
				if (digit < '0' || digit > '9') && (digit < 'A' || digit > 'F') {
					return nativePointAddress{code: acquisitionReasonInvalidAddress}
				}
			}
			offset, err := strconv.ParseUint(digits, 16, 24)
			if err != nil {
				return nativePointAddress{code: acquisitionReasonInvalidAddress}
			}
			return nativePointAddress{area: address[:1], start: int(offset), nativeBits: 1, bitIndex: -1, function: point.Function}
		}
	}
	resolved := pointNativeAddress(identity.Protocol, point)
	if resolved.code == "" && (identity.Protocol == "mc" || identity.Protocol == acquisitionProtocolMitsubishi) && slices.Contains([]string{"LTN", "LSTN", "LCN"}, resolved.area) {
		return nativePointAddress{code: "unsupported-mc-long-timer-counter"}
	}
	return resolved
}
