package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

func validateFatekReadProfile(req connector.ReadRequest) error {
	if err := connector.ValidateReadProfile(req); err != nil {
		return err
	}
	if req.Profile.Present() && req.Profile.Protocol() != schema.ProtocolFatekFBs {
		return connector.ErrReadProfile
	}
	if req.Profile.Present() && req.Profile.Resolution().Mode != schema.PointFormatLegacy {
		switch req.DataType {
		case schema.DataTypeBool, schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeString,
			schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32,
			schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		default:
			return decoding.ErrInapplicableFormat
		}
	}
	return nil
}

// Native compatibility deliberately retains the historical high-component-first
// assembler. An explicit resolved format instead operates on canonical 16-bit
// words, even when the wire component itself is 32 bits (DR/DD).
func decodeFatekReadComponents(values []int, width int, req connector.ReadRequest) (any, error) {
	r := req.Profile.Resolution()
	if !req.Profile.Present() || r.Mode == schema.PointFormatLegacy || r.NativeCompatibility || !r.Applicable {
		return convertFatekValue(values, req.DataType, width)
	}
	words, err := decoding.NormalizeFatekComponents(values, width, req.DataType)
	if err != nil {
		return nil, err
	}
	return decoding.DecodeScalarWords(words, req.DataType, r.EffectiveFormat)
}
