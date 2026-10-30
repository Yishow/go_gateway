package optimizer

import (
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"strings"
)

// validNativeSpan supplements permissive parsers with actual wire bounds.
func validNativeSpan(protocol, area string, start, width int) bool {
	if start < 0 || width <= 0 || width-1 > int(^uint(0)>>1)-start {
		return false
	}
	end := start + width - 1
	switch protocol {
	case "modbus":
		return end <= 65535
	case "mc", "mitsubishi":
		// The existing MC binary client encodes a 24-bit address (client.go).
		return end <= 0xFFFFFF
	case "fatek":
		component, err := fatek.GetComponentType(area)
		if err != nil {
			return false
		}
		if _, err = fatek.FormatAddress(component, start); err != nil {
			return false
		}
		_, err = fatek.FormatAddress(component, end)
		return err == nil
	default:
		return false
	}
}

// nativeClientLimit reflects the existing contiguous read implementations, not
// a universal protocol or random-read limit. See the three protocol clients.
func nativeClientLimit(protocol string, nativeBits int) int {
	switch protocol {
	case "modbus", "mc", "mitsubishi":
		if nativeBits == 1 {
			return 2000
		}
		return 125
	case "fatek":
		if nativeBits == 1 {
			return 255
		}
		if nativeBits == 32 {
			return 32
		}
		return 64
	default:
		return 0
	}
}

func validCoverage(protocol string, coverage ReadableCoverage) bool {
	for _, ranges := range [][]NativeRange{coverage.Readable, coverage.Forbidden, coverage.Excluded} {
		for _, r := range ranges {
			if !canonicalNativeArea(protocol, r.Area) || r.Start < 0 || r.End <= r.Start || !validNativeSpan(protocol, r.Area, r.Start, r.End-r.Start) {
				return false
			}
		}
	}
	return true
}

// canonicalNativeArea validates declaration keys without silently normalizing
// restrictions into a different scope. Client support, not a permissive string
// prefix, is authoritative. Valid unrelated areas remain independent.
func canonicalNativeArea(protocol, area string) bool {
	if area == "" || strings.TrimSpace(area) != area || strings.ToUpper(area) != area {
		return false
	}
	switch protocol {
	case "modbus":
		return area == "HR" || area == "IR" || area == "COIL" || area == "DI"
	case "mc", "mitsubishi":
		_, err := mcprotocol.GetDeviceType(area)
		return err == nil
	case "fatek":
		// GetComponentType permits generic D* names; require an existing parser
		// symbol or a named client component omitted by the general address parser.
		symbol, _, err := fatek.ParseComponentAddress(area + "0")
		if err == nil && symbol == area {
			_, err = fatek.GetComponentType(area)
			return err == nil
		}
		switch area {
		case "F", "RT", "RC", "WT", "WC", "DF":
			_, err = fatek.GetComponentType(area)
			return err == nil
		}
	}
	return false
}
