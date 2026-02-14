package hsllogic

import (
	"fmt"
	"strconv"
	"strings"
)

// =============================================================================
// FATEK 地址解析
// =============================================================================

// fatekDevices FATEK 設備碼列表
var fatekDevices = []string{
	"WX", "WY", "WM", "WS", "WT", "WC", // 字組設備 (複合)
	"DD", "DR", "DF", // 雙字組設備
	"X", "Y", "M", "S", "T", "C", // 位元/單字組設備
	"D", "R", "F", // 字組設備
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
