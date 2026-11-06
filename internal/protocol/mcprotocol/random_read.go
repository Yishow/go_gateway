package mcprotocol

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// RandomReadItem 定義隨機讀取項目
type RandomReadItem struct {
	Device     string
	Addr       int
	DoubleWord bool
	DataType   string
}

// RandomWriteItem 定義隨機寫入字組項目
type RandomWriteItem struct {
	Device string
	Addr   int
	Value  int // 16-bit word value
}

// RandomWriteBitItem 定義隨機寫入位元項目
type RandomWriteBitItem struct {
	Device string
	Addr   int
	Value  bool // bit value
}

// ExtractBitFromWord extracts one bit (0-15) from a 16-bit word returned by MC random read.
func ExtractBitFromWord(word, bitOffset int) (bool, error) {
	if bitOffset < 0 || bitOffset > 15 {
		return false, fmt.Errorf("bit offset out of range: %d (valid: 0-15)", bitOffset)
	}
	return ((uint16(word) >> bitOffset) & 1) != 0, nil
}

type randomReadEntry struct {
	index int
	addr  int
	code  byte
}

// RandomRead reads word and double-word devices via MC 3E binary command 0403 (subcommand 0000).
// Bit devices are permitted in the word section when 16-point aligned (addr%16 == 0).
func (c *MCClient) RandomRead(items []RandomReadItem) ([]int, error) {
	count := len(items)
	if count < 1 || count > 255 {
		return nil, fmt.Errorf("count out of range: %d (valid: 1-255 for random read)", count)
	}

	wordEntries := make([]randomReadEntry, 0, count)
	dwordEntries := make([]randomReadEntry, 0)
	for i, item := range items {
		if strings.EqualFold(strings.TrimSpace(item.DataType), "string") {
			return nil, fmt.Errorf("random read does not support string data type")
		}
		if item.Addr < 0 || item.Addr > 0xFFFFFF {
			return nil, fmt.Errorf("address out of range: %d (valid: 0-16777215)", item.Addr)
		}
		devType, err := GetDeviceType(item.Device)
		if err != nil {
			return nil, err
		}
		if IsLongTimerOrCounter(devType.Name) {
			return nil, fmt.Errorf("device type %s is a long timer/counter and is not supported in random read", item.Device)
		}
		if devType.IsBit {
			if item.DoubleWord {
				return nil, fmt.Errorf("device type %s is a bit device and cannot use double word random read", item.Device)
			}
			if item.Addr%16 != 0 {
				return nil, fmt.Errorf("bit device %s address %d must be 16-point aligned for random read", item.Device, item.Addr)
			}
		}
		entry := randomReadEntry{index: i, addr: item.Addr, code: devType.Code}
		if item.DoubleWord {
			dwordEntries = append(dwordEntries, entry)
		} else {
			wordEntries = append(wordEntries, entry)
		}
	}

	// Binary 0403 (SH-080008): wordCount(1) + dwordCount(1) + word entries(4*W) + dword entries(4*D).
	data := make([]byte, 2+count*4)
	data[0] = byte(len(wordEntries))
	data[1] = byte(len(dwordEntries))
	offset := 2
	for _, entry := range wordEntries {
		data[offset] = byte(entry.addr & 0xFF)
		data[offset+1] = byte((entry.addr >> 8) & 0xFF)
		data[offset+2] = byte((entry.addr >> 16) & 0xFF)
		data[offset+3] = entry.code
		offset += 4
	}
	for _, entry := range dwordEntries {
		data[offset] = byte(entry.addr & 0xFF)
		data[offset+1] = byte((entry.addr >> 8) & 0xFF)
		data[offset+2] = byte((entry.addr >> 16) & 0xFF)
		data[offset+3] = entry.code
		offset += 4
	}

	req := c.frame.BuildPacket(CmdRandomRead, SubCmdWord, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return nil, err
	}

	if len(resp) < minResponseLen {
		return nil, fmt.Errorf("response too short: expected at least %d bytes, got %d", minResponseLen, len(resp))
	}

	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return nil, &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for random read (count: %d)", endCode, count),
		}
	}

	raw := resp[2:]
	expectedLen := len(wordEntries)*2 + len(dwordEntries)*4
	if len(raw) < expectedLen {
		return nil, fmt.Errorf("response data too short: expected %d bytes, got %d bytes for %d items", expectedLen, len(raw), count)
	}

	res := make([]int, count)
	for i, entry := range wordEntries {
		res[entry.index] = int(binary.LittleEndian.Uint16(raw[i*2:]))
	}
	dwordBase := len(wordEntries) * 2
	for i, entry := range dwordEntries {
		res[entry.index] = int(binary.LittleEndian.Uint32(raw[dwordBase+i*4:]))
	}
	return res, nil
}

// RandomWrite 隨機寫入多個字組和位元設備
//
// 注意：此功能需要 Command 0x1402，目前未實作。
// 參考: MC Protocol 3E Binary Frame Specification
//
// Args:
//   - wordItems: 要寫入的字組設備列表
//   - bitItems: 要寫入的位元設備列表
//
// Returns:
//   - 錯誤（目前總是返回未實作錯誤）
func (c *MCClient) RandomWrite(wordItems []RandomWriteItem, bitItems []RandomWriteBitItem) error {
	// TODO: Implement Command 0x1402 (Random Write)
	// 參考: MC Protocol 3E Binary Frame Specification
	// 需要構建複雜的請求結構，包含字組和位元設備的混合寫入
	return fmt.Errorf("random write not implemented: requires Command 0x1402")
}
