package mcprotocol

import (
	"encoding/binary"
	"fmt"
)

// Header Constants
// SubHeader 以 Big Endian 寫入線上封包（0x5000 → bytes [0x50, 0x00]）
const (
	ReqSubHeader = 0x5000
	ResSubHeader = 0xD000
)

// Frame3E represents the 3E Frame structure helpers
// Request: SubHeader(2) + Net(1) + PC(1) + IO(2) + Station(1) + Len(2) + Timer(2) + Cmd(2) + Sub(2) + Data...
type RequestFrame struct {
	NetworkNo byte
	PCNo      byte
	IONo      uint16
	StationNo byte
	Timer     uint16 // 250ms units
}

func NewRequestFrame(net, pc, station byte) RequestFrame {
	return RequestFrame{
		NetworkNo: net,
		PCNo:      pc,
		IONo:      0x03FF, // Default own station IO
		StationNo: station,
		Timer:     0x0010, // 4 seconds (approx)
	}
}

// BuildPacket constructs the full byte slice.
//
// Callers must validate the payload size so the encoded length fits in uint16.
func (f RequestFrame) BuildPacket(cmd, subCmd uint16, data []byte) []byte {
	// Payload = Timer(2) + Cmd(2) + Sub(2) + Data(...)
	payloadLen := 2 + 2 + 2 + len(data)

	// Total Size = Header(9) + Payload
	// Header: Sub(2)+Net(1)+PC(1)+IO(2)+Station(1)+Len(2)
	totalSize := 9 + payloadLen

	buf := make([]byte, totalSize)

	// Header：SubHeader 為 Big Endian（線上位元組序 [0x50, 0x00]），其餘欄位為 Little Endian
	binary.BigEndian.PutUint16(buf[0:], ReqSubHeader)
	buf[2] = f.NetworkNo
	buf[3] = f.PCNo
	binary.LittleEndian.PutUint16(buf[4:], f.IONo)
	buf[6] = f.StationNo
	binary.LittleEndian.PutUint16(buf[7:], uint16(payloadLen)) // #nosec G115 -- callers validate request sizes before building a wire frame.

	// Payload
	binary.LittleEndian.PutUint16(buf[9:], f.Timer)
	binary.LittleEndian.PutUint16(buf[11:], cmd)
	binary.LittleEndian.PutUint16(buf[13:], subCmd)
	copy(buf[15:], data)

	return buf
}

// ParseResponseHeader parses the response header and checks for errors
// Response: Sub(2)+Net(1)+PC(1)+IO(2)+Station(1)+Len(2)+EndCode(2)+Data...
// Returns: data (after EndCode), error
// Note: 回應標頭使用原始位元組序列 [0xD0, 0x00]，不是 Little Endian 格式
func ParseResponseHeader(header []byte) (int, error) {
	if len(header) < 9 {
		return 0, fmt.Errorf("response too short")
	}

	// 檢查 Subheader: 應該是 [0xD0, 0x00]
	// Python 代碼直接比較: header_res[:2] != b'\xD0\x00'
	if header[0] != 0xD0 || header[1] != 0x00 {
		received := uint16(header[0])<<8 | uint16(header[1])
		return 0, fmt.Errorf("invalid response subheader: 0x%04X (expected 0xD000)", received)
	}

	// 數據長度使用 Little Endian
	dataLen := binary.LittleEndian.Uint16(header[7:])
	return int(dataLen), nil
}

// PackBits packs bools into bytes (2 bits per byte for 3E Binary)
// 1st bit -> High nibble, 2nd bit -> Low nibble
func PackBits(values []bool) []byte {
	count := len(values)
	byteCount := (count + 1) / 2
	buf := make([]byte, byteCount)

	for i := 0; i < count; i += 2 {
		val1 := 0
		if values[i] {
			val1 = 1
		}

		val2 := 0
		if i+1 < count {
			if values[i+1] { // #nosec G602 -- i+1 is checked against count before indexing the optional second bit.
				val2 = 1
			}
		}

		// First device occupies the high nibble; second device occupies the low nibble.
		buf[i/2] = byte(val1<<4 | val2) // #nosec G115 -- val1 and val2 are each limited to the wire values 0 or 1.
	}
	return buf
}

// UnpackBits unpacks bytes into bools
func UnpackBits(data []byte, count int) []bool {
	res := make([]bool, count)

	for i := 0; i < count; i++ {
		byteIdx := i / 2
		isHigh := (i % 2) == 0 // i=0 -> High, i=1 -> Low

		if byteIdx >= len(data) {
			break
		}

		val := data[byteIdx]
		var bitVal int

		if isHigh {
			bitVal = int((val >> 4) & 0x0F)
		} else {
			bitVal = int(val & 0x0F)
		}

		res[i] = (bitVal == 1)
	}
	return res
}
