package fatek

import (
	"bytes"
	"fmt"
)

// hexChars 預定義十六進位字元表，用於高效能轉換
var hexChars = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

// CalculateLRC calculates the Longitudinal Redundancy Check
// Logic: Sum of all bytes (including STX) modulo 256, returned as Hex string.
func CalculateLRC(data []byte) string {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return fmt.Sprintf("%02X", sum)
}

// BuildFrame constructs the ASCII frame and returns a new byte slice.
//
// Deprecated: Use BuildFrameToBuffer for better performance.
func BuildFrame(station int, cmd, body string) []byte {
	buf := GetBuffer()
	defer PutBuffer(buf)

	BuildFrameToBuffer(buf, station, cmd, body)

	// We must copy the result because the buffer goes back to the pool
	res := make([]byte, buf.Len())
	copy(res, buf.Bytes())
	return res
}

// BuildFrameToBuffer writes the ASCII frame directly into the provided buffer.
// Structure: STX + Station(2) + Command(2) + Body + LRC(2) + ETX
func BuildFrameToBuffer(buf *bytes.Buffer, station int, cmd, body string) {
	// 1. Write STX
	buf.WriteByte(STX)

	// 2. Write Station (2 chars)
	// Optimization: fmt.Fprintf is slower than manual string manipulation, but safe.
	// For max performance, we could implement specialized IntToHex writer.
	fmt.Fprintf(buf, "%02X", station)

	// 3. Write Command
	buf.WriteString(cmd)

	// 4. Write Body
	buf.WriteString(body)

	// 5. Calculate LRC (Current content of buffer)
	// buf.Bytes() returns the slice from STX to end of Body
	lrc := CalculateLRC(buf.Bytes())

	// 6. Write LRC
	buf.WriteString(lrc)

	// 7. Write ETX
	buf.WriteByte(ETX)
}

// ParseResponse validates and extracts the body from a response frame
func ParseResponse(response []byte, expectedCmd string) (string, error) {
	// Min length check
	// Normal: STX(1) + Station(2) + Cmd(2) + Status(1) + LRC(2) + ETX(1) = 9
	// Loopback(4E): STX(1) + Station(2) + Cmd(2) + LRC(2) + ETX(1) = 8 (Empty body)

	minLen := 9
	if expectedCmd == "4E" {
		minLen = 8
	}

	if len(response) < minLen {
		return "", ErrResponseTooShort
	}

	if response[0] != STX {
		return "", ErrInvalidSTX
	}
	if response[len(response)-1] != ETX {
		return "", ErrInvalidETX
	}

	// Validate LRC
	// Content to check is everything before LRC(2) + ETX(1)
	contentLen := len(response) - 3
	contentWithSTX := response[:contentLen]
	receivedLRC := string(response[contentLen : contentLen+2])

	calcLRC := CalculateLRC(contentWithSTX)
	if receivedLRC != calcLRC {
		return "", fmt.Errorf("%w: received %s, calculated %s", ErrChecksumMismatch, receivedLRC, calcLRC)
	}

	// Extract Command
	// Station is [1:3], Command is [3:5]
	cmd := string(response[3:5])

	// Special handling for Loopback (4E)
	// The response echoes the command 4E, but there is NO Status Code '0'.
	// Structure: STX + Station + 4E + Data + LRC + ETX
	if expectedCmd == "4E" {
		if cmd != "4E" {
			return "", ErrInvalidCommand
		}
		// Body is everything after Command and before LRC
		// Start index = 5 (STX+Station+Cmd)
		if len(response) <= 8 { // No data?
			return "", nil
		}
		body := string(response[5 : len(response)-3])
		return body, nil
	}

	if cmd != expectedCmd {
		return "", fmt.Errorf("%w: expected %s, got %s", ErrInvalidCommand, expectedCmd, cmd)
	}

	// Status Code (Index 5)
	statusCode := string(response[5])
	if statusCode != "0" {
		return "", NewProtocolError(statusCode, cmd)
	}

	// Body (Index 6 to End-3)
	// If response is just Ack (like write), len might be 9.
	if len(response) == 9 {
		return "", nil
	}

	body := string(response[6 : len(response)-3])
	return body, nil
}

// HexToInt converts Hex string to int
func HexToInt(hexStr string) (int, error) {
	// Use ParseUint to handle potential large values if needed,
	// but standard Fatek data fits in int (up to 32-bit)
	// Using Sscanf is easy
	var val int
	_, err := fmt.Sscanf(hexStr, "%x", &val)
	return val, err
}

// IntToHex 將整數轉換為固定寬度的十六進位字串
// 使用查表法避免 fmt.Sprintf 的效能開銷
// Args:
//   - val: 要轉換的整數值
//   - width: 輸出字串的固定寬度 (1-8)
//
// Returns:
//   - 大寫十六進位字串，左側補零至指定寬度
func IntToHex(val, width int) string {
	// 防護性預設：確保寬度在合理範圍內
	if width <= 0 {
		width = 1
	} else if width > 8 {
		width = 8
	}

	// 使用固定大小陣列避免堆分配（小於 64 bytes 的陣列通常在棧上分配）
	buf := make([]byte, width)

	// 從低位到高位填充
	for i := width - 1; i >= 0; i-- {
		buf[i] = hexChars[val&0x0F]
		val >>= 4
	}

	return string(buf)
}
