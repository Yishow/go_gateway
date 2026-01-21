package hsllogic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// =============================================================================
// ProtocolType 協議類型定義
// =============================================================================

// ProtocolType 協議類型常數
type ProtocolType string

const (
	// ProtocolModbus Modbus 協議
	ProtocolModbus ProtocolType = "modbus"
	// ProtocolMitsubishi 三菱 MC 協議
	ProtocolMitsubishi ProtocolType = "mitsubishi"
	// ProtocolFatek 永宏 FATEK 協議
	ProtocolFatek ProtocolType = "fatek"
	// ProtocolSiemens 西門子 S7 協議
	ProtocolSiemens ProtocolType = "siemens"
)

// =============================================================================
// ParsedAddress 解析後的地址結構
// =============================================================================

// ParsedAddress 解析後的地址結構，統一表示各協議的設備地址
type ParsedAddress struct {
	// Protocol 協議類型
	Protocol ProtocolType `json:"protocol"`
	// DeviceType 設備/區域類型 (如 "D", "M", "HR", "DB" 等)
	DeviceType string `json:"device_type"`
	// Offset 主偏移量 (暫存器地址或起始位元)
	Offset int `json:"offset"`
	// BitIndex 位元索引，用於位元設備 (如 M10.1 中的 1)，-1 表示非位元存取
	BitIndex int `json:"bit_index"`
	// DBNumber 資料區塊編號，用於 Siemens S7 (如 DB1.DBX0.0 中的 1)
	DBNumber int `json:"db_number"`
	// IsBitDevice 是否為位元設備
	IsBitDevice bool `json:"is_bit_device"`
	// Raw 原始地址字串
	Raw string `json:"raw"`
}

// =============================================================================
// 地址解析器
// =============================================================================

// AddressParser 地址解析器介面
type AddressParser interface {
	// Parse 解析地址字串
	Parse(address string) (*ParsedAddress, error)
	// Protocol 取得協議類型
	Protocol() ProtocolType
}

// =============================================================================
// ParseAddress 統一地址解析入口
// =============================================================================

// ParseAddress 統一解析各協議地址，根據協議類型自動選擇解析器
func ParseAddress(protocol ProtocolType, address string) (*ParsedAddress, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return nil, fmt.Errorf("地址不可為空")
	}

	switch protocol {
	case ProtocolModbus:
		return parseModbusAddress(address)
	case ProtocolMitsubishi:
		return parseMitsubishiAddress(address)
	case ProtocolFatek:
		return parseFatekAddress(address)
	case ProtocolSiemens:
		return parseSiemensAddress(address)
	default:
		return nil, fmt.Errorf("不支援的協議類型: %s", protocol)
	}
}

// =============================================================================
// Modbus 地址解析
// =============================================================================

// modbusAddressPattern Modbus 地址正則表達式
// 格式: 0xxxx (Coil), 1xxxx (DI), 3xxxx (IR), 4xxxx (HR) 或簡化格式
var modbusAddressPattern = regexp.MustCompile(`^([0134])?(\d+)(?:\.(\d+))?$`)

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
		deviceType = "COIL"
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
			deviceType = "COIL"
		default:
			// 預設為 Holding Register
			deviceType = "HR"
		}
	}

	isBit := deviceType == "COIL" || deviceType == "DI" || bitIndex >= 0

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

// =============================================================================
// 三菱 MC Protocol 地址解析
// =============================================================================

// mitsubishiDevices 三菱設備碼列表 (按長度排序以優先匹配較長的名稱)
var mitsubishiDevices = []string{
	"ZR", "SD", "SW", "SB", "SM", // 特殊設備
	"CC", "TC", "ST",             // 計時器/計數器 (雙字元)
	"D", "W", "R", "B", "F",      // 字組設備
	"M", "L", "S", "X", "Y",      // 位元設備
	"T", "C",                     // 計時器/計數器
}

// mitsubishiBitDevices 三菱位元設備集合
var mitsubishiBitDevices = map[string]bool{
	"M": true, "L": true, "S": true,
	"X": true, "Y": true, "B": true,
	"F": true, "SB": true, "SM": true,
}

// parseMitsubishiAddress 解析三菱 MC Protocol 地址
// 支援格式: "D100", "M0", "X0", "Y0", "R100", "W100", etc.
// 注意: W, B, SW, SB 設備使用 16 進位編址；X, Y 設備使用 8 進位編址
func parseMitsubishiAddress(address string) (*ParsedAddress, error) {
	address = strings.ToUpper(strings.TrimSpace(address))
	if len(address) < 2 {
		return nil, fmt.Errorf("無效的三菱地址: %s", address)
	}

	// Hex 編址設備 (Link Register/Relay)
	hexDevices := map[string]bool{
		"W": true, "B": true, "SW": true, "SB": true,
	}

	for _, device := range mitsubishiDevices {
		if strings.HasPrefix(address, device) {
			addrPart := address[len(device):]
			var offset int
			var err error

			switch {
			// X, Y 設備使用八進位
			case device == "X" || device == "Y":
				offset64, parseErr := strconv.ParseInt(addrPart, 8, 32)
				if parseErr != nil {
					// 嘗試十進位
					offset64, parseErr = strconv.ParseInt(addrPart, 10, 32)
					if parseErr != nil {
						return nil, fmt.Errorf("無效的三菱地址數字: %s", addrPart)
					}
				}
				offset = int(offset64)

			// W, B, SW, SB 設備強制使用十六進位
			case hexDevices[device]:
				offset64, parseErr := strconv.ParseInt(addrPart, 16, 32)
				if parseErr != nil {
					return nil, fmt.Errorf("無效的三菱 Hex 地址數字: %s (設備 %s 應使用十六進位)", addrPart, device)
				}
				offset = int(offset64)

			// 其他設備使用十進位
			default:
				offset, err = strconv.Atoi(addrPart)
				if err != nil {
					return nil, fmt.Errorf("無效的三菱地址數字: %s", addrPart)
				}
			}

			return &ParsedAddress{
				Protocol:    ProtocolMitsubishi,
				DeviceType:  device,
				Offset:      offset,
				BitIndex:    -1,
				DBNumber:    0,
				IsBitDevice: mitsubishiBitDevices[device],
				Raw:         address,
			}, nil
		}
	}

	return nil, fmt.Errorf("無法識別的三菱設備碼: %s", address)
}

// =============================================================================
// FATEK 地址解析
// =============================================================================

// fatekDevices FATEK 設備碼列表
var fatekDevices = []string{
	"WX", "WY", "WM", "WS", "WT", "WC", // 字組設備 (複合)
	"DD", "DR", "DF",                   // 雙字組設備
	"X", "Y", "M", "S", "T", "C",       // 位元/單字組設備
	"D", "R", "F",                      // 字組設備
}

// fatekBitDevices FATEK 位元設備集合
var fatekBitDevices = map[string]bool{
	"X": true, "Y": true, "M": true, "S": true,
	"T": true, "C": true,
}

// parseFatekAddress 解析 FATEK 地址
// 支援格式: "D0100", "R0", "M100", "X0", "Y0", etc.
func parseFatekAddress(address string) (*ParsedAddress, error) {
	address = strings.ToUpper(strings.TrimSpace(address))
	if len(address) < 2 {
		return nil, fmt.Errorf("無效的 FATEK 地址: %s", address)
	}

	for _, device := range fatekDevices {
		if strings.HasPrefix(address, device) {
			addrPart := address[len(device):]
			offset, err := strconv.Atoi(addrPart)
			if err != nil {
				return nil, fmt.Errorf("無效的 FATEK 地址數字: %s", addrPart)
			}

			return &ParsedAddress{
				Protocol:    ProtocolFatek,
				DeviceType:  device,
				Offset:      offset,
				BitIndex:    -1,
				DBNumber:    0,
				IsBitDevice: fatekBitDevices[device],
				Raw:         address,
			}, nil
		}
	}

	return nil, fmt.Errorf("無法識別的 FATEK 設備碼: %s", address)
}

// =============================================================================
// Siemens S7 地址解析
// =============================================================================

// siemensAddressPattern Siemens S7 地址正則表達式
// 格式: DB1.DBX0.0, DB1.DBW10, M0.0, I0.0, Q0.0, MW100, IW0, etc.
var siemensAddressPattern = regexp.MustCompile(
	`^(DB(\d+)\.)?(DB[XBWD]|[MIQV][BWD]?|PE[BWD]?)(\d+)(?:\.(\d+))?$`,
)

// parseSiemensAddress 解析 Siemens S7 地址
// 支援格式:
//   - "DB1.DBX0.0" (資料區塊位元)
//   - "DB1.DBW10" (資料區塊字組)
//   - "M0.0" (Marker 位元)
//   - "MW100" (Marker 字組)
//   - "I0.0" (Input)
//   - "Q0.0" (Output)
func parseSiemensAddress(address string) (*ParsedAddress, error) {
	address = strings.ToUpper(strings.TrimSpace(address))

	matches := siemensAddressPattern.FindStringSubmatch(address)
	if matches == nil {
		return nil, fmt.Errorf("無效的 Siemens S7 地址: %s", address)
	}

	// matches[2] = DB 編號 (如果有)
	// matches[3] = 設備類型 (DBX, DBW, M, MW, I, Q, etc.)
	// matches[4] = 偏移量
	// matches[5] = 位元索引 (如果有)

	dbNumber := 0
	if matches[2] != "" {
		db, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("無效的 DB 編號: %s", matches[2])
		}
		dbNumber = db
	}

	deviceType := matches[3]
	offset, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("無效的偏移量: %s", matches[4])
	}

	bitIndex := -1
	if matches[5] != "" {
		bit, err := strconv.Atoi(matches[5])
		if err != nil || bit < 0 || bit > 7 {
			return nil, fmt.Errorf("無效的位元索引: %s", matches[5])
		}
		bitIndex = bit
	}

	// 判斷是否為位元設備
	isBit := bitIndex >= 0 ||
		deviceType == "DBX" ||
		(len(deviceType) == 1 && (deviceType == "M" || deviceType == "I" || deviceType == "Q"))

	return &ParsedAddress{
		Protocol:    ProtocolSiemens,
		DeviceType:  deviceType,
		Offset:      offset,
		BitIndex:    bitIndex,
		DBNumber:    dbNumber,
		IsBitDevice: isBit,
		Raw:         address,
	}, nil
}

// =============================================================================
// 輔助方法
// =============================================================================

// String 取得地址的字串表示
func (p *ParsedAddress) String() string {
	if p.DBNumber > 0 {
		if p.BitIndex >= 0 {
			return fmt.Sprintf("DB%d.%s%d.%d", p.DBNumber, p.DeviceType, p.Offset, p.BitIndex)
		}
		return fmt.Sprintf("DB%d.%s%d", p.DBNumber, p.DeviceType, p.Offset)
	}
	if p.BitIndex >= 0 {
		return fmt.Sprintf("%s%d.%d", p.DeviceType, p.Offset, p.BitIndex)
	}
	return fmt.Sprintf("%s%d", p.DeviceType, p.Offset)
}

// ToByteOffset 計算位元組偏移量 (用於連續讀取優化)
func (p *ParsedAddress) ToByteOffset() int {
	switch p.Protocol {
	case ProtocolSiemens:
		// Siemens 地址已經是位元組偏移
		return p.Offset
	default:
		// 字組協議：暫存器地址 * 2
		return p.Offset * 2
	}
}
