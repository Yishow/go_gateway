package memory

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// =============================================================================
// 記憶體區域定義
// =============================================================================

// MemoryArea 記憶體區域類型
type MemoryArea int

const (
	// Modbus 標準區域
	AreaCoil            MemoryArea = iota // 00001-09999 Coils (讀寫 bit)
	AreaDiscreteInput                     // 10001-19999 Discrete Inputs (唯讀 bit)
	AreaInputRegister                     // 30001-39999 Input Registers (唯讀 16-bit)
	AreaHoldingRegister                   // 40001-49999 Holding Registers (讀寫 16-bit)

	// PLC 風格區域
	AreaDataRegister  // D 區 (16-bit 資料暫存器)
	AreaMarker        // M 區 (bit 標記)
	AreaInternalRelay // R 區 (16-bit 內部繼電器)
	AreaTimer         // T 區 (計時器)
	AreaCounter       // C 區 (計數器)
	AreaInput         // X 區 (輸入)
	AreaOutput        // Y 區 (輸出)
)

// 錯誤定義
var (
	ErrInvalidAddress = errors.New("無效的記憶體地址格式")
)

// =============================================================================
// AddressMapper 地址映射器
// =============================================================================

// AddressMapper 將 PLC/Modbus 地址映射到記憶體偏移量
type AddressMapper struct {
	// 正則表達式模式
	modbusPattern *regexp.Regexp
	plcPattern    *regexp.Regexp
}

// NewAddressMapper 建立新的地址映射器
func NewAddressMapper() *AddressMapper {
	return &AddressMapper{
		modbusPattern: regexp.MustCompile(`^(\d{5})$`),
		plcPattern:    regexp.MustCompile(`^([A-Za-z]+)(\d+)$`),
	}
}

// Map 將地址映射到字節偏移量
func (m *AddressMapper) Map(address string) (int, error) {
	_, offset, err := m.MapWithArea(address)
	return offset, err
}

// MapWithArea 將地址映射到記憶體區域和偏移量
func (m *AddressMapper) MapWithArea(address string) (MemoryArea, int, error) {
	address = strings.TrimSpace(address)

	// 嘗試 Modbus 風格 (5 位數字)
	if m.modbusPattern.MatchString(address) {
		return m.parseModbusAddress(address)
	}

	// 嘗試 PLC 風格 (字母+數字)
	if m.plcPattern.MatchString(address) {
		return m.parsePLCAddress(address)
	}

	return 0, 0, ErrInvalidAddress
}

// parseModbusAddress 解析 Modbus 地址 (如 40001)
func (m *AddressMapper) parseModbusAddress(address string) (MemoryArea, int, error) {
	num, err := strconv.Atoi(address)
	if err != nil {
		return 0, 0, ErrInvalidAddress
	}

	switch {
	case num >= 1 && num <= 9999:
		// Coils: 00001-09999 -> bit offset
		return AreaCoil, (num - 1) / 8, nil

	case num >= 10001 && num <= 19999:
		// Discrete Inputs: 10001-19999 -> bit offset
		return AreaDiscreteInput, (num - 10001) / 8, nil

	case num >= 30001 && num <= 39999:
		// Input Registers: 30001-39999 -> word offset * 2
		return AreaInputRegister, (num - 30001) * 2, nil

	case num >= 40001 && num <= 49999:
		// Holding Registers: 40001-49999 -> word offset * 2
		return AreaHoldingRegister, (num - 40001) * 2, nil

	default:
		return 0, 0, ErrInvalidAddress
	}
}

// parsePLCAddress 解析 PLC 風格地址 (如 D100, M0)
func (m *AddressMapper) parsePLCAddress(address string) (MemoryArea, int, error) {
	matches := m.plcPattern.FindStringSubmatch(address)
	if len(matches) != 3 {
		return 0, 0, ErrInvalidAddress
	}

	prefix := strings.ToUpper(matches[1])
	num, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, ErrInvalidAddress
	}

	switch prefix {
	case "D":
		// D 區: 16-bit 暫存器，每個佔 2 bytes
		return AreaDataRegister, num * 2, nil

	case "M":
		// M 區: bit 標記，每 8 個佔 1 byte
		return AreaMarker, num / 8, nil

	case "R":
		// R 區: 16-bit 內部繼電器
		return AreaInternalRelay, num * 2, nil

	case "T":
		// T 區: 計時器 (16-bit)
		return AreaTimer, num * 2, nil

	case "C":
		// C 區: 計數器 (16-bit)
		return AreaCounter, num * 2, nil

	case "X":
		// X 區: 輸入 (bit)
		return AreaInput, num / 8, nil

	case "Y":
		// Y 區: 輸出 (bit)
		return AreaOutput, num / 8, nil

	default:
		return 0, 0, ErrInvalidAddress
	}
}

// GetBitOffset 獲取 bit 地址的 bit 偏移 (0-7)
func (m *AddressMapper) GetBitOffset(address string) (int, error) {
	address = strings.TrimSpace(address)

	// Modbus Coil
	if m.modbusPattern.MatchString(address) {
		num, _ := strconv.Atoi(address)
		if num >= 1 && num <= 9999 {
			return (num - 1) % 8, nil
		}
		if num >= 10001 && num <= 19999 {
			return (num - 10001) % 8, nil
		}
	}

	// PLC bit 區域
	if m.plcPattern.MatchString(address) {
		matches := m.plcPattern.FindStringSubmatch(address)
		prefix := strings.ToUpper(matches[1])
		num, _ := strconv.Atoi(matches[2])

		switch prefix {
		case "M", "X", "Y":
			return num % 8, nil
		}
	}

	return 0, ErrInvalidAddress
}
