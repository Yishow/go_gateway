package fatek

import (
	"fmt"
)

// CalculateLRC calculates the Longitudinal Redundancy Check
// Logic: Sum of all bytes (including STX) modulo 256, returned as Hex string.
func CalculateLRC(data []byte) string {
	var sum byte = 0
	for _, b := range data {
		sum += b
	}
	return fmt.Sprintf("%02X", sum)
}

// BuildFrame constructs the ASCII frame
// Structure: STX + Station(2) + Command(2) + Body + LRC(2) + ETX
func BuildFrame(station int, cmd string, body string) []byte {
	// 1. Station to Hex String
	stationStr := fmt.Sprintf("%02X", station)
	
	// 2. Content = Station + Cmd + Body
	content := stationStr + cmd + body
	
	// 3. Prepare data for LRC (STX + Content)
	// Note: We don't prepend STX to the string yet because LRC needs the byte value of STX
	
	lrcData := make([]byte, 0, 1+len(content))
	lrcData = append(lrcData, STX)
	lrcData = append(lrcData, []byte(content)...)
	
	// 4. Calculate LRC
	lrc := CalculateLRC(lrcData)
	
	// 5. Final Frame: STX + Content + LRC + ETX
	frame := make([]byte, 0, len(lrcData)+3)
	frame = append(frame, lrcData...)
	frame = append(frame, []byte(lrc)...)
	frame = append(frame, ETX)
	
	return frame
}

// ParseResponse validates and extracts the body from a response frame
func ParseResponse(response []byte, expectedCmd string) (string, error) {
	// Min length: STX(1) + Station(2) + Cmd(2) + Status(1) + LRC(2) + ETX(1) = 9
	if len(response) < 9 {
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

// IntToHex converts int to Hex string with fixed width
func IntToHex(val int, width int) string {
	format := fmt.Sprintf("%%0%dX", width)
	return fmt.Sprintf(format, val)
}