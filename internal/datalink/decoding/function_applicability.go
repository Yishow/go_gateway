package decoding

import (
	"go-gateway/internal/datalink/schema"
	"strconv"
	"strings"
)

// pointFormatCapability preserves transport address precedence. Modbus aliases
// and recognized five-digit references select their function before an explicit
// Function; only raw numeric offsets use Function. This mirrors the adapter's
// parseModbusAddress without changing its legacy read path or the context-free
// FormatCapability API.
func pointFormatCapability(protocol schema.ProtocolType, p schema.Point) (ScalarCapability, error) {
	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		address, err := CanonicalModbusAddress(p.Address, p.Function)
		if err != nil {
			return ScalarCapability{}, err
		}
		return FormatCapability(protocol, address, p.DataType)
	default:
		return FormatCapability(protocol, p.Address, p.DataType)
	}
}

// CanonicalModbusAddress maps a read address and Function to the actual wire area
// alias, preserving transport precedence and rejecting unsupported access modes.
func CanonicalModbusAddress(address, function string) (string, error) {
	address = strings.TrimSpace(address)
	upper := strings.ToUpper(address)
	// These prefixes select the wire area regardless of Function.
	for _, prefix := range []string{"HR", "IR", "CO", "C", "DI"} {
		if strings.HasPrefix(upper, prefix) {
			offset, err := strconv.ParseUint(address[len(prefix):], 10, 16)
			if err != nil {
				return "", ErrInapplicableFormat
			}
			return prefix + strconv.FormatUint(offset, 10), nil
		}
	}
	offset, err := strconv.ParseUint(address, 10, 32)
	if err != nil {
		return "", ErrInapplicableFormat
	}
	if len(address) == 5 && offset%10000 >= 1 && offset%10000 <= 9999 {
		prefix := map[byte]string{'4': "HR", '3': "IR", '1': "DI", '0': "C"}[address[0]]
		if prefix != "" {
			return prefix + strconv.FormatUint(offset%10000-1, 10), nil
		}
	}
	function = strings.TrimSpace(strings.ToLower(function))
	if function == "" && offset >= 1 && offset <= 9999 {
		return "C" + strconv.FormatUint(offset-1, 10), nil
	}
	if offset > 65535 {
		return "", ErrInapplicableFormat
	}
	if function == "" {
		function = "03"
	}
	prefix := map[string]string{"03": "HR", "04": "IR", "01": "C", "02": "DI", "holding": "HR", "input": "IR", "coil": "C", "discrete": "DI"}[function]
	if prefix == "" {
		return "", ErrInapplicableFormat
	}
	return prefix + strconv.FormatUint(offset, 10), nil
}
