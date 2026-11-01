package optimizer

import (
	"go-gateway/internal/datalink/decoding"
	"strings"
)

// pointNativeAddress uses the same Function/address precedence as the actual
// Modbus adapter and format resolver. The original Point snapshot stays intact.
func pointNativeAddress(protocol string, point EffectivePoint) (area string, start, nativeBits, bitIndex int, function, code string) {
	address := point.Address
	function = point.Function
	if protocol == acquisitionProtocolModbus {
		canonical, err := decoding.CanonicalModbusAddress(address, function)
		if err != nil {
			return "", 0, 0, 0, "", acquisitionReasonInvalidAddress
		}
		// The canonical wire aliases C/CO correspond to the native parser's
		// internal COIL area label. No address precedence is reimplemented here.
		switch {
		case strings.HasPrefix(canonical, "CO"):
			address = acquisitionAreaCoil + strings.TrimPrefix(canonical, "CO")
		case strings.HasPrefix(canonical, "C"):
			address = acquisitionAreaCoil + strings.TrimPrefix(canonical, "C")
		default:
			address = canonical
		}
	}
	area, start, nativeBits, bitIndex, code = nativeAddress(protocol, address)
	if code == "" && protocol == acquisitionProtocolModbus {
		function = map[string]string{"HR": "03", "IR": "04", acquisitionAreaCoil: "01", "DI": "02"}[area]
	}
	return
}
