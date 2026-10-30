package decoding

import (
	"errors"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
)

// ErrInapplicableFormat means the point cannot use an explicit scalar order.
var ErrInapplicableFormat = errors.New("point type or address does not support scalar format; reset to inherited format")

// ScalarCapability describes transport-native units, not a changed read count.
type ScalarCapability struct {
	NativeBits int
	Components int
	Words      int
}

// FormatCapability validates a complete numeric scalar against the actual client
// address parser and wire span. Bit addresses never become word reads implicitly.
func FormatCapability(protocol schema.ProtocolType, address string, dtype schema.DataType) (ScalarCapability, error) {
	bits := scalarBits(dtype)
	if bits == 0 {
		return ScalarCapability{}, ErrInapplicableFormat
	}
	if protocol == schema.ProtocolFatekFBs {
		symbol, start, err := fatek.ParseComponentAddress(address)
		if err != nil {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		component, err := fatek.GetComponentType(symbol)
		if err != nil || (component.Width != 16 && component.Width != 32) {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		count := bits / component.Width
		if start < 0 || count-1 > int(^uint(0)>>1)-start {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		if _, err = fatek.FormatAddress(component, start); err != nil {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		if _, err = fatek.FormatAddress(component, start+count-1); err != nil {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		return ScalarCapability{component.Width, count, bits / 16}, nil
	}
	var parser hsllogic.ProtocolType
	maxAddress := 65535
	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		parser = hsllogic.ProtocolModbus
	case schema.ProtocolMC3E:
		parser = hsllogic.ProtocolMitsubishi
		maxAddress = 0xFFFFFF
	default:
		return ScalarCapability{}, ErrUnsupportedProtocol
	}
	parsed, err := hsllogic.ParseAddress(parser, address)
	if err != nil || parsed.BitIndex >= 0 {
		return ScalarCapability{}, ErrInapplicableFormat
	}
	bit := parsed.IsBitDevice
	if protocol == schema.ProtocolMC3E {
		device, err := mcprotocol.GetDeviceType(parsed.DeviceType)
		if err != nil {
			return ScalarCapability{}, ErrInapplicableFormat
		}
		bit = device.IsBit
	}
	words := bits / 16
	if bit || parsed.Offset < 0 || parsed.Offset > maxAddress-(words-1) {
		return ScalarCapability{}, ErrInapplicableFormat
	}
	return ScalarCapability{16, words, words}, nil
}

func scalarBits(dtype schema.DataType) int {
	switch dtype {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 32
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 64
	default:
		return 0
	}
}
