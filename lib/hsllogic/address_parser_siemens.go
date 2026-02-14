package hsllogic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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
