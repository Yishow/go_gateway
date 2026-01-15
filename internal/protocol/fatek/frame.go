package fatek

import (
	"fmt"
)

// CalculateLRC 計算縱向冗餘校驗 (LRC) checksum
//
// 計算範圍包含 STX、Station、Command 到 Body 的最後一個字元
// (不包含 Checksum 本身與 ETX)
//
// Args:
//   - data: 要計算 checksum 的位元組序列 (包含 STX)
//
// Returns:
//   - 兩個字元的十六進位字串
func CalculateLRC(data []byte) string {
	var lrc byte = 0
	for _, b := range data {
		lrc = (lrc + b) & 0xFF
	}
	return fmt.Sprintf("%02X", lrc)
}

// BuildFrame 構建完整的 FATEK ASCII 訊框
//
// 結構：STX + Station(2) + Command(2) + Body + LRC(2) + ETX
//
// Args:
//   - station: 站號 (0-255)
//   - command: 2 字元命令碼 (例如 '44')
//   - body: 命令的資料負載
//
// Returns:
//   - 編碼後的 ASCII 訊框，準備傳輸
func BuildFrame(station int, command, body string) []byte {
	stationStr := fmt.Sprintf("%02X", station)
	content := stationStr + command + body

	// 計算 LRC。注意：根據文檔，LRC 包含 STX
	lrcSum := byte(STX)
	for _, char := range content {
		lrcSum += byte(char)
	}
	lrcHex := fmt.Sprintf("%02X", lrcSum&0xFF)

	// 構建完整訊框
	frame := string(STX) + content + lrcHex + string(ETX)
	return []byte(frame)
}

// ParseResponse 驗證並解析回應訊框
//
// 檢查 STX、ETX、LRC，以及 PLC 回傳的錯誤碼
//
// Args:
//   - response: 原始回應訊框
//   - expectedCmd: 發送的命令碼
//
// Returns:
//   - 回應的 body (資料負載)
//
// Raises:
//   - FatekCommunicationError: 如果結構或 checksum 無效
//   - FatekProtocolError: 如果 PLC 回傳錯誤碼 (非零狀態)
func ParseResponse(response []byte, expectedCmd string) (string, error) {
	if len(response) < 6 {
		return "", NewFatekCommunicationError("response too short: %d bytes", len(response))
	}

	// 檢查 STX
	if response[0] != STX {
		return "", NewFatekCommunicationError("invalid STX")
	}

	// 檢查 ETX
	if response[len(response)-1] != ETX {
		return "", NewFatekCommunicationError("invalid ETX")
	}

	// 驗證 Checksum
	// 最後 3 個位元組之前 (LRC(2) + ETX(1)) 是內容
	contentWithSTX := response[:len(response)-3]
	receivedLRC := string(response[len(response)-3 : len(response)-1])

	lrcSum := byte(0)
	for _, b := range contentWithSTX {
		lrcSum += b
	}
	calculatedLRC := fmt.Sprintf("%02X", lrcSum&0xFF)

	if receivedLRC != calculatedLRC {
		return "", NewFatekCommunicationError("checksum mismatch. received: %s, calculated: %s", receivedLRC, calculatedLRC)
	}

	// 解析欄位
	// Station(2) -> [1:3]
	// Command(2) -> [3:5]
	cmd := string(response[3:5])

	if cmd != expectedCmd {
		return "", NewFatekCommunicationError("unexpected command in response: %s (expected: %s)", cmd, expectedCmd)
	}

	// 狀態碼 / 錯誤碼 (命令後的第一個位元組，索引 5)
	// 對於讀取命令 (例如 44, 46, 48)，成功是 '0' 後跟資料
	// 對於寫入命令 (例如 45, 47)，成功只是 '0' (無 body)
	// 對於錯誤，它是錯誤碼 '2', '4', 'A' 而不是 '0'
	statusCode := string(response[5:6])

	if statusCode != "0" {
		return "", NewFatekProtocolError(statusCode, cmd)
	}

	// 回傳 Body (跳過 STX(1)+Station(2)+Cmd(2)+StatusCode(1) = 6 個位元組)
	// 從索引 6 開始。Checksum 在 -3
	body := string(response[6 : len(response)-3])
	return body, nil
}
