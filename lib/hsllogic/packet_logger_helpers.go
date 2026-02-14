package hsllogic

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// =============================================================================
// 輔助函數
// =============================================================================

// formatHexBytes 格式化位元組為十六進位字串
func formatHexBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(data))
}

// FormatHexWithSpaces 格式化位元組為帶空格的十六進位字串
func FormatHexWithSpaces(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, " ")
}

// ParseHexString 解析十六進位字串為位元組
func ParseHexString(hexStr string) ([]byte, error) {
	// 移除空格和前綴
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	return hex.DecodeString(hexStr)
}
