package adapters

import (
	"errors"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

var errMCReadCount = errors.New("invalid MC scalar read count")

func validateMCReadProfile(req connector.ReadRequest) error {
	if err := connector.ValidateReadProfile(req); err != nil {
		return err
	}
	if req.Profile.Present() && req.Profile.Protocol() != schema.ProtocolMC3E {
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
func modernMCScalar(req connector.ReadRequest) bool {
	r := req.Profile.Resolution()
	return req.Profile.Present() && r.Mode != schema.PointFormatLegacy && r.Applicable && !r.NativeCompatibility
}
func decodeMCReadWords(words []uint16, req connector.ReadRequest, count int) (any, error) {
	if count < 1 {
		return nil, errMCReadCount
	}
	width := schema.RegisterCountForDataType(req.DataType)
	if width < 1 || count > 125/width || len(words) != count*width {
		return nil, decoding.ErrShortScalar
	}
	format := req.Profile.Resolution().EffectiveFormat
	if count == 1 {
		return decoding.DecodeScalarWords(words, req.DataType, format)
	}
	values := make([]any, count)
	for i := range values {
		value, err := decoding.DecodeScalarWords(words[i*width:(i+1)*width], req.DataType, format)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return values, nil
}
