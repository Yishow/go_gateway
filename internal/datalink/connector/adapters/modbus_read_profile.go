package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

func decodeModbusReadWords(words []uint16, req connector.ReadRequest, connectionFormat string) (any, error) {
	if err := validateModbusReadProfile(req); err != nil {
		return nil, err
	}
	profile := req.Profile
	resolution := profile.Resolution()
	if !profile.Present() || resolution.Mode == schema.PointFormatLegacy {
		return convertModbusValue(words, req.DataType, effectiveModbusDataFormat(req.DataFormat, connectionFormat)), nil
	}
	if resolution.Applicable {
		return decoding.DecodeScalarWords(words, req.DataType, resolution.EffectiveFormat)
	}
	// New inherit of single-register/bool/string values keeps the native converter
	// semantics. The resolved default, not a stale legacy projection, selects order.
	return convertModbusValue(words, req.DataType, resolution.EffectiveFormat), nil
}

func validateModbusReadProfile(req connector.ReadRequest) error {
	if err := connector.ValidateReadProfile(req); err != nil {
		return err
	}
	if !req.Profile.Present() {
		return nil
	}
	switch req.Profile.Protocol() {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		return nil
	default:
		return connector.ErrReadProfile
	}
}
