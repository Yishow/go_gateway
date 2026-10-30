package mcprotocol

import (
	"encoding/binary"
	"fmt"
)

// BatchReadWord preserves the historical minimum-length response policy.
func (c *MCClient) BatchReadWord(device string, addr, count int) ([]int, error) {
	return c.batchReadWord(device, addr, count, false)
}

// BatchReadWordExact reads native 16-bit components with an exact payload length.
// It shares address/count/EndCode/deadline handling with the legacy entry point;
// scalar and future raw/scatter consumers must not silently truncate extra data.
func (c *MCClient) BatchReadWordExact(device string, addr, count int) ([]int, error) {
	return c.batchReadWord(device, addr, count, true)
}

func (c *MCClient) batchReadWord(device string, addr, count int, exact bool) ([]int, error) {
	// 輸入驗證
	if addr < 0 || addr > 0xFFFFFF {
		return nil, fmt.Errorf("address out of range: %d (valid: 0-16777215)", addr)
	}
	if count < 1 || count > 125 {
		return nil, fmt.Errorf("count out of range: %d (valid: 1-125 for word devices)", count)
	}

	devType, err := GetDeviceType(device)
	if err != nil {
		return nil, err
	}

	// 驗證設備類型為字組設備
	if devType.IsBit {
		return nil, fmt.Errorf("device type %s is a bit device, use BatchReadBit instead", device)
	}

	// Prepare Data: Head(3) + Code(1) + Count(2)
	data := make([]byte, 6)
	data[0] = byte(addr & 0xFF)
	data[1] = byte((addr >> 8) & 0xFF)
	data[2] = byte((addr >> 16) & 0xFF)
	data[3] = devType.Code
	countU16, err := checkedUint16(count, "count")
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint16(data[4:], countU16)

	// Build & Send
	req := c.frame.BuildPacket(CmdBatchRead, SubCmdWord, data)
	resp, err := c.transport.SendReceive(req)
	if err != nil {
		return nil, err
	}

	// 驗證回應最小長度
	if len(resp) < minResponseLen {
		return nil, fmt.Errorf("response too short: expected at least %d bytes, got %d", minResponseLen, len(resp))
	}

	// Check EndCode (First 2 bytes)
	endCode := binary.LittleEndian.Uint16(resp[0:])
	if endCode != 0 {
		return nil, &MCError{
			Code: int(endCode),
			Msg:  fmt.Sprintf("PLC returned error code 0x%04X for device %s at address %d (count: %d)", endCode, device, addr, count),
		}
	}

	// Data follows EndCode
	expectedLen := count * 2
	raw := resp[2:]
	if len(raw) < expectedLen {
		return nil, fmt.Errorf("response data too short: expected %d bytes, got %d", expectedLen, len(raw))
	}

	if exact && len(raw) != expectedLen {
		return nil, fmt.Errorf("invalid word response length: expected %d bytes, got %d", expectedLen, len(raw))
	}

	res := make([]int, count)
	for i := 0; i < count; i++ {
		res[i] = int(binary.LittleEndian.Uint16(raw[i*2:]))
	}
	return res, nil
}
