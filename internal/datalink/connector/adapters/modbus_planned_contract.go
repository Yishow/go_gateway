package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"strconv"
)

func validateModbusPlannedRequest(r optimizer.AcquisitionRequest, bindings []connector.PlannedPointBinding, protocol schema.ProtocolType, frame string, station byte) error {
	if station == 0 {
		station = 1
	}
	function := modbusPlannedFunction(r.Group.Function)
	area := map[byte]string{1: "COIL", 2: "DI", 3: "HR", 4: "IR"}[function]
	limit, nativeBits := 125, 16
	if function == 1 || function == 2 {
		limit, nativeBits = 2000, 1
	}
	if r.Group.Identity.Protocol != "modbus" || r.Group.Identity.Frame != frame || r.Group.Identity.Station != strconv.Itoa(int(station)) || r.Kind != "contiguous" || r.AddressStep != 1 || function == 0 || r.Group.Area != area || len(r.Units) != 0 || r.Length < 1 || r.Length > limit || r.Length > r.Limit || r.Start < 0 || r.Start > 65536-r.Length {
		return connector.ErrPlannedRead
	}
	for _, p := range r.Points {
		address, pointFunction, err := parseModbusAddress(p.Point.Address, p.Point.Function)
		width := schema.RegisterCountForDataType(p.Point.DataType)
		if err != nil || modbusPlannedFunction(pointFunction) != function || p.NativeBits != nativeBits || p.Length != width || p.Length != p.Point.NativeWidth || p.Offset < 0 || p.Offset > r.Length-p.Length || r.Start+p.Offset != int(address) || (nativeBits == 1 && p.Point.DataType != schema.DataTypeBool) {
			return connector.ErrPlannedRead
		}
	}
	for _, b := range bindings {
		if b.Profile.Protocol() != protocol {
			return connector.ErrReadProfile
		}
		if err := validateModbusReadProfile(connector.ReadRequest{Address: b.Point.Address, Function: b.Point.Function, DataType: b.Point.DataType, Profile: b.Profile}); err != nil {
			return err
		}
	}
	return nil
}

func modbusPlannedFunction(function string) byte {
	switch function {
	case "01", "coil", "fc01":
		return 1
	case "02", "discrete", "fc02":
		return 2
	case "03", "holding", "fc03", "":
		return 3
	case "04", "input", "fc04":
		return 4
	default:
		return 0
	}
}
