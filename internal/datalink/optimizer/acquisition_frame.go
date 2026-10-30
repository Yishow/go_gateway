package optimizer

import "go-gateway/internal/protocol/fatek"

// contiguousFrameBytes describes complete application frames. Supported frame
// identifiers are modbus tcp/rtu, MC 3e-binary, and FATEK ascii. Transport names
// such as "tcp" do not imply an MC encoding. Unknown encoding stays explicit.
// Encodings match modbus.BuildTCPFrame/BuildRTUFrame; MC RequestFrame.BuildPacket,
// ParseResponseHeader and PackBits; FATEK BuildFrame and ReadRegisters/ReadStatus.
func contiguousFrameBytes(identity AcquisitionIdentity, area string, start, count, nativeBits int) (ApplicationFrameBudget, bool) {
	switch identity.Protocol {
	case "modbus":
		payload := count * 2
		if nativeBits == 1 {
			payload = (count + 7) / 8
		}
		switch identity.Frame {
		case "tcp":
			return ApplicationFrameBudget{RequestBytes: 12, ResponseBytes: 9 + payload}, true
		case "rtu":
			return ApplicationFrameBudget{RequestBytes: 8, ResponseBytes: 5 + payload}, true
		}
	case "mc", "mitsubishi":
		if identity.Frame != "3e-binary" {
			break
		}
		payload := count * 2
		if nativeBits == 1 {
			payload = (count + 1) / 2
		}
		return ApplicationFrameBudget{RequestBytes: 21, ResponseBytes: 11 + payload}, true
	case "fatek":
		if identity.Frame != "ascii" {
			break
		}
		component, err := fatek.GetComponentType(area)
		if err != nil {
			break
		}
		address, err := fatek.FormatAddress(component, start)
		if err != nil {
			break
		}
		payload := count * (nativeBits / 4)
		if nativeBits == 1 {
			payload = count
		}
		return ApplicationFrameBudget{RequestBytes: 10 + len(address), ResponseBytes: 9 + payload}, true
	}
	return ApplicationFrameBudget{}, false
}

func strictFrameBudget(cap AcquisitionCapability) (ApplicationFrameBudget, bool) {
	var result ApplicationFrameBudget
	for _, budget := range []ApplicationFrameBudget{cap.ProtocolFrameBudget, cap.ClientFrameBudget, cap.DeviceFrameBudget} {
		if budget.RequestBytes < 0 || budget.ResponseBytes < 0 {
			return ApplicationFrameBudget{}, false
		}
		if budget.RequestBytes > 0 && (result.RequestBytes == 0 || budget.RequestBytes < result.RequestBytes) {
			result.RequestBytes = budget.RequestBytes
		}
		if budget.ResponseBytes > 0 && (result.ResponseBytes == 0 || budget.ResponseBytes < result.ResponseBytes) {
			result.ResponseBytes = budget.ResponseBytes
		}
	}
	return result, true
}

func fitsFrameBudget(cost, budget ApplicationFrameBudget) bool {
	return (budget.RequestBytes == 0 || cost.RequestBytes <= budget.RequestBytes) &&
		(budget.ResponseBytes == 0 || cost.ResponseBytes <= budget.ResponseBytes)
}

// frameLimitedUnits uses the monotone encoded frame size, not bytes/nativeBits:
// MC packs bit nibbles and FATEK encodes word values as ASCII hex characters.
func frameLimitedUnits(source AcquisitionSource, area string, start, width, nativeBits, limit int) (int, ApplicationFrameBudget, bool, string) {
	budget, valid := strictFrameBudget(source.Capability)
	if !valid {
		return 0, ApplicationFrameBudget{}, false, "invalid-capability"
	}
	cost, known := contiguousFrameBytes(source.Identity, area, start, width, nativeBits)
	if !known {
		if budget != (ApplicationFrameBudget{}) {
			return 0, cost, false, "unsupported-frame-cost"
		}
		return limit, cost, false, ""
	}
	if !fitsFrameBudget(cost, budget) {
		return 0, cost, true, "scalar-exceeds-byte-limit"
	}
	low, high := width, limit
	for low < high {
		mid := low + (high-low+1)/2
		candidate, _ := contiguousFrameBytes(source.Identity, area, start, mid, nativeBits)
		if fitsFrameBudget(candidate, budget) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low, cost, true, ""
}
