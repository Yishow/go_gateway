package sourcerule

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

func qiqrBitAddress(identity optimizer.AcquisitionIdentity, address string) bool {
	address = strings.ToUpper(strings.TrimSpace(address))
	return identity.Protocol == "mc" && identity.ProtocolRevision == "mitsubishi-q-iqr-3e-binary-v1" &&
		(strings.HasPrefix(address, "X") || strings.HasPrefix(address, "Y"))
}

func offsetAddressForIdentity(address string, delta int, protocol schema.ProtocolType, identity optimizer.AcquisitionIdentity) (string, error) {
	if !qiqrBitAddress(identity, address) {
		return offsetAddress(address, delta, protocol)
	}
	span, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: address, DataType: schema.DataTypeBool})
	if err != nil || delta > mcClientMaxAddress-span.Start {
		return "", validationError("invalid Q/iQ-R bit address or offset")
	}
	return fmt.Sprintf("%s%X", span.Area, max(span.Start+delta, 0)), nil
}

func buildPlannedPointAddressesForIdentity(address string, count int, dataType schema.DataType, protocol schema.ProtocolType, identity optimizer.AcquisitionIdentity) ([]string, error) {
	if !qiqrBitAddress(identity, address) {
		return buildPlannedPointAddresses(address, count, dataType, protocol)
	}
	span, err := optimizer.AcquisitionPointSpanForIdentity(identity, optimizer.EffectivePoint{Address: address, DataType: dataType})
	if err != nil || count < 0 {
		return nil, validationError("invalid Q/iQ-R bit address or count")
	}
	step := getDataTypeCellSpan(dataType)
	if count > 0 && count-1 > (mcClientMaxAddress-span.Start)/step {
		return nil, validationError("Q/iQ-R address plan exceeds client range")
	}
	addresses := make([]string, 0, count)
	for i := range count {
		addresses = append(addresses, fmt.Sprintf("%s%X", span.Area, span.Start+i*step))
	}
	return addresses, nil
}
