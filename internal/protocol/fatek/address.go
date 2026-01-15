package fatek

import (
	"fmt"
	"strings"
)

// NormalizeAddress 格式化元件位址為協定要求的固定長度
//
// Args:
//   - symbolType: 元件類型 (例如 'X', 'R', 'D')
//   - number: 位址編號
//
// Returns:
//   - 格式化後的位址字串 (例如 'X00010', 'D00100')
//
// Raises:
//   - error: 如果元件類型未知
func NormalizeAddress(symbolType string, number int) (string, error) {
	symbol := strings.ToUpper(symbolType)

	// 5 字元：X, Y, M, S, T, C (單點)
	if symbol == "X" || symbol == "Y" || symbol == "M" || symbol == "S" || symbol == "T" || symbol == "C" {
		return fmt.Sprintf("%s%04d", symbol, number), nil
	}

	// 6 字元：R, D, RT, RC, F (16-bit 暫存器)
	if symbol == "R" || symbol == "D" || symbol == "RT" || symbol == "RC" || symbol == "F" {
		return fmt.Sprintf("%s%05d", symbol, number), nil
	}

	// 6 字元：WX, WY, WM, WS, WT, WC (16-bit 存取離散點)
	if symbol == "WX" || symbol == "WY" || symbol == "WM" || symbol == "WS" || symbol == "WT" || symbol == "WC" {
		return fmt.Sprintf("%s%04d", symbol, number), nil
	}

	// 7 字元：32-bit 暫存器 (DR, DD, DF, DWX, DWY 等)
	if symbol == "DR" || symbol == "DD" || symbol == "DF" {
		return fmt.Sprintf("%s%05d", symbol, number), nil
	}

	// DWX, DWY, DWM, DWS, DWT, DWC
	if strings.HasPrefix(symbol, "DW") && len(symbol) == 3 {
		lastChar := symbol[2]
		if lastChar == 'X' || lastChar == 'Y' || lastChar == 'M' || lastChar == 'S' || lastChar == 'T' || lastChar == 'C' {
			return fmt.Sprintf("%s%04d", symbol, number), nil
		}
	}

	return "", fmt.Errorf("unknown component type: %s", symbolType)
}
