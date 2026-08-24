package sourcerule

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

func buildPlannedPointAddresses(startAddress string, count int, dataType schema.DataType, protocol schema.ProtocolType) ([]string, error) {
	if count < 0 {
		return nil, validationError("count must not be negative")
	}
	span := getDataTypeCellSpan(dataType)
	addresses := make([]string, 0, count)
	for index := 0; index < count; index++ {
		address, err := offsetAddress(startAddress, index*span, protocol)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func getDataTypeCellSpan(dataType schema.DataType) int {
	switch dataType {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
		return 2
	case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return 4
	default:
		return 1
	}
}

func offsetAddress(address string, delta int, protocol schema.ProtocolType) (string, error) {
	address = strings.ToUpper(strings.TrimSpace(address))
	if address == "" {
		return "", validationError("address is required")
	}

	switch protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		if len(address) < 5 {
			return "", validationError(fmt.Sprintf("invalid modbus address: %s", address))
		}
		prefix := address[:1]
		if prefix != "0" && prefix != "1" && prefix != "3" && prefix != "4" {
			return "", validationError(fmt.Sprintf("invalid modbus address: %s", address))
		}
		number := address[1:]
		value, ok := parseDecimal(number)
		if !ok {
			return "", validationError(fmt.Sprintf("invalid modbus address: %s", address))
		}
		next := value + delta
		if next < 1 {
			next = 1
		}
		return fmt.Sprintf("%s%04d", prefix, next), nil
	case schema.ProtocolFatekFBs, schema.ProtocolMC3E:
		area, number := splitAlphaNumeric(address)
		if area == "" || number == "" {
			return "", validationError(fmt.Sprintf("invalid %s address: %s", strings.ToLower(string(protocol)), address))
		}
		if protocol == schema.ProtocolFatekFBs {
			if !isFatekArea(area) {
				return "", validationError(fmt.Sprintf("invalid fatek address: %s", address))
			}
			value, ok := parseDecimal(number)
			if !ok {
				return "", validationError(fmt.Sprintf("invalid fatek address: %s", address))
			}
			next := value + delta
			if next < 0 {
				next = 0
			}
			return fmt.Sprintf("%s%d", area, next), nil
		}

		if isMC3EHexArea(area) {
			value, ok := parseHex(number)
			if !ok {
				return "", validationError(fmt.Sprintf("invalid mc3e address: %s", address))
			}
			next := value + delta
			if next < 0 {
				next = 0
			}
			return fmt.Sprintf("%s%X", area, next), nil
		}
		if !isMC3EDecimalArea(area) {
			return "", validationError(fmt.Sprintf("invalid mc3e address: %s", address))
		}
		value, ok := parseDecimal(number)
		if !ok {
			return "", validationError(fmt.Sprintf("invalid mc3e address: %s", address))
		}
		next := value + delta
		if next < 0 {
			next = 0
		}
		return fmt.Sprintf("%s%d", area, next), nil
	default:
		return "", validationError(fmt.Sprintf("unsupported protocol address: %s", strings.ToLower(string(protocol))))
	}
}

func parseDecimal(raw string) (int, bool) {
	if raw == "" {
		return 0, false
	}
	value := 0
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return 0, false
		}
		value = value*10 + int(ch-'0')
	}
	return value, true
}

func parseHex(raw string) (int, bool) {
	if raw == "" {
		return 0, false
	}
	value := 0
	for _, ch := range raw {
		var digit int
		switch {
		case ch >= '0' && ch <= '9':
			digit = int(ch - '0')
		case ch >= 'A' && ch <= 'F':
			digit = int(ch-'A') + 10
		default:
			return 0, false
		}
		value = value*16 + digit
	}
	return value, true
}

func isFatekArea(area string) bool {
	switch area {
	case "X", "Y", "M", "S", "T", "C", "D", "R":
		return true
	default:
		return false
	}
}

func isMC3EHexArea(area string) bool {
	return area == "X" || area == "Y" || area == "B"
}

func isMC3EDecimalArea(area string) bool {
	return area == "D" || area == "W" || area == "M"
}

func splitAlphaNumeric(input string) (string, string) {
	index := 0
	for ; index < len(input); index++ {
		ch := input[index]
		if ch >= '0' && ch <= '9' {
			break
		}
	}
	if index == 0 || index >= len(input) {
		return "", ""
	}
	return input[:index], input[index:]
}
