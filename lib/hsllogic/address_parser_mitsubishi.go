package hsllogic

import (
	"fmt"
	"strconv"
	"strings"
)

// =============================================================================
// 三菱 MC Protocol 地址解析
// =============================================================================

// mitsubishiDevices 三菱設備碼列表 (按長度排序以優先匹配較長的名稱)
var mitsubishiDevices = []string{
	"ZR", "SD", "SW", "SB", "SM", // 特殊設備
	"CC", "TC", "ST", // 計時器/計數器 (雙字元)
	"D", "W", "R", "B", "F", // 字組設備
	"M", "L", "S", "X", "Y", // 位元設備
	"T", "C", // 計時器/計數器
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
