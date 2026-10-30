package optimizer

import (
	"fmt"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/lib/hsllogic"
)

// Check fixture offsets/widths against existing protocol parsers, not the legacy
// decimal-only merger parser. MC native word width is NOT adapter Count (typed values).
func validateAcceptanceNativePoint(t *testing.T, c acceptanceCase, p acceptancePoint) {
	t.Helper()
	address := p.Address
	if address == "" {
		prefix := "D"
		if c.Protocol == "modbus" {
			prefix = "HR"
		}
		address = fmt.Sprintf("%s%d", prefix, p.Start)
	}
	var offset int
	var err error
	nativeBits := 16
	if c.Protocol == "fatek" {
		var symbol string
		symbol, offset, err = fatek.ParseComponentAddress(address)
		if err == nil {
			var component fatek.ComponentType
			component, err = fatek.GetComponentType(symbol)
			if err == nil {
				nativeBits = component.Width
			}
		}
	} else {
		protocol := hsllogic.ProtocolModbus
		if c.Protocol == "mc" {
			protocol = hsllogic.ProtocolMitsubishi
		}
		var parsed *hsllogic.ParsedAddress
		parsed, err = hsllogic.ParseAddress(protocol, address)
		if err == nil {
			offset = parsed.Offset
		}
	}
	if c.Reason == "invalid-address" {
		if err == nil {
			t.Fatal("invalid-address vector unexpectedly parsed")
		}
		return
	}
	if err != nil || offset != p.Start {
		t.Fatalf("native parser %q: offset=%d err=%v", address, offset, err)
	}
	if c.NativeBits != 0 && c.NativeBits != nativeBits {
		t.Fatal("native component width differs from source")
	}
	switch schema.DataType(p.Type) {
	case schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat32, schema.DataTypeFloat64:
	default:
		t.Fatalf("unsupported fixture type %q", p.Type)
	}
	width := (schema.RegisterCountForDataType(schema.DataType(p.Type))*16 + nativeBits - 1) / nativeBits
	if width != p.Width {
		t.Fatalf("native width: source=%d fixture=%d", width, p.Width)
	}
}
