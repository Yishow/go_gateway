package hsllogic

import (
	"fmt"
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
