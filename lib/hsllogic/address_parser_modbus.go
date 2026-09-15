package hsllogic

import (
	"fmt"
	"strconv"
	"strings"
)

// =============================================================================
// Modbus 地址解析
// =============================================================================

const modbusCoilDeviceType = "COIL"

// parseModbusAddress 解析 Modbus 地址
// 支援格式:
//   - "100" 或 "40100" (Holding Register)
//   - "00100" (Coil)
//   - "10100" (Discrete Input)
//   - "30100" (Input Register)
//   - "100.5" (Holding Register 第 5 位元)
func parseModbusAddress(address string) (*ParsedAddress, error) {
	address = strings.ToUpper(strings.TrimSpace(address))

	// 處理特殊前綴格式: HR100, IR100, COIL100, DI100
	var deviceType string
	var addrPart string

	switch {
	case strings.HasPrefix(address, "HR"):
		deviceType = "HR"
		addrPart = address[2:]
	case strings.HasPrefix(address, "IR"):
		deviceType = "IR"
		addrPart = address[2:]
	case strings.HasPrefix(address, "COIL"):
		deviceType = modbusCoilDeviceType
		addrPart = address[4:]
	case strings.HasPrefix(address, "DI"):
		deviceType = "DI"
		addrPart = address[2:]
	default:
		addrPart = address
	}

	// 解析位元索引
	bitIndex := -1
	if dotIdx := strings.Index(addrPart, "."); dotIdx > 0 {
		bitPart := addrPart[dotIdx+1:]
		addrPart = addrPart[:dotIdx]
		bit, err := strconv.Atoi(bitPart)
		if err != nil || bit < 0 || bit > 15 {
			return nil, fmt.Errorf("無效的位元索引: %s", bitPart)
		}
		bitIndex = bit
	}

	// 解析地址數值
	addrNum, err := strconv.Atoi(addrPart)
	if err != nil {
		return nil, fmt.Errorf("無效的 Modbus 地址: %s", address)
	}

	// 根據五位數前綴判斷設備類型
	if deviceType == "" {
		switch {
		case addrNum >= 40000 && addrNum < 50000:
			deviceType = "HR"
			addrNum -= 40000
		case addrNum >= 30000 && addrNum < 40000:
			deviceType = "IR"
			addrNum -= 30000
		case addrNum >= 10000 && addrNum < 20000:
			deviceType = "DI"
			addrNum -= 10000
		case addrNum >= 0 && addrNum < 10000:
			deviceType = modbusCoilDeviceType
		default:
			// 預設為 Holding Register
			deviceType = "HR"
		}
	}

	isBit := deviceType == modbusCoilDeviceType || deviceType == "DI" || bitIndex >= 0

	return &ParsedAddress{
		Protocol:    ProtocolModbus,
		DeviceType:  deviceType,
		Offset:      addrNum,
		BitIndex:    bitIndex,
		DBNumber:    0,
		IsBitDevice: isBit,
		Raw:         address,
	}, nil
}
