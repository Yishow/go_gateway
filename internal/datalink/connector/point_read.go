package connector

import (
	"fmt"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// PointReadCount returns the protocol-specific quantity for one typed point.
// MC counts typed values, while the existing Modbus/FATEK adapters count words.
// Direct batch callers retain their own Count contract.
func PointReadCount(protocol schema.ProtocolType, dataType schema.DataType) int {
	if protocol == schema.ProtocolMC3E {
		return 1
	}
	return schema.RegisterCountForDataType(dataType)
}

// PointReadCountForAddress returns the protocol-specific quantity for one typed point.
// FATEK counts native components, whose width depends on the addressed symbol.
func PointReadCountForAddress(protocol schema.ProtocolType, address string, dataType schema.DataType) (int, error) {
	if protocol != schema.ProtocolFatekFBs {
		return PointReadCount(protocol, dataType), nil
	}

	symbol, componentAddress, err := fatek.ParseComponentAddress(address)
	if err != nil {
		return 0, fmt.Errorf("invalid FATEK scalar address %q: %w", address, err)
	}
	component, err := fatek.GetComponentType(symbol)
	if err != nil {
		return 0, fmt.Errorf("resolve FATEK component %q: %w", symbol, err)
	}
	if _, err := fatek.FormatAddress(component, componentAddress); err != nil {
		return 0, fmt.Errorf("validate FATEK scalar address %q: %w", address, err)
	}
	if component.IsDiscrete {
		return 1, nil
	}

	bits := schema.RegisterCountForDataType(dataType) * 16
	count := (bits + component.Width - 1) / component.Width
	return max(count, 1), nil
}
