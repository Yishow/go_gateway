package fatek

import (
	"math"
	"strings"
)

// AddressStep is the physical-address footprint/advance of one component.
// DR0 contains R0/R1; its next non-overlapping component is DR2. Word access
// to discrete memory advances 16 bits (W*) or 32 bits (DW*), not one bit.
func (c ComponentType) AddressStep() int {
	if strings.HasPrefix(c.Name, "DW") {
		return 32
	}
	if strings.HasPrefix(c.Name, "W") {
		return 16
	}
	return max(1, c.Width/16)
}

// AddressLimit is the exclusive physical-memory limit from the protocol's
// component table. A complete component must fit, not merely its start digits.
// Device declarations may further restrict this protocol limit.
func (c ComponentType) AddressLimit() int {
	switch c.Name {
	case "D", "R", "F", "DR", "DD", "DF":
		return 65536
	}
	return int(math.Pow10(c.FormatLen - len(c.Name)))
}
func (c ComponentType) requiresEightAlignedStart() bool {
	switch c.Name {
	case "WX", "WY", "WM", "WS", "DWX", "DWY", "DWM", "DWS":
		return true
	}
	return false
}
