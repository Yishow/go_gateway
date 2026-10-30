package fatek

import (
	"fmt"
	"strconv"
	"strings"
)

var supportedComponentSymbols = []string{
	"DR", "DD", "DWM", "DWS", "DWX", "DWY",
	"D", "R", "WM", "WS", "WX", "WY",
	"M", "S", "T", "C", "X", "Y",
}

// ParseComponentAddress splits a supported FATEK component symbol from its address.
func ParseComponentAddress(addressStr string) (symbol string, address int, err error) {
	normalized := strings.TrimSpace(strings.ToUpper(addressStr))
	if len(normalized) < 2 {
		return "", 0, fmt.Errorf("無效的 FATEK 地址: %s", normalized)
	}

	for _, candidate := range supportedComponentSymbols {
		if !strings.HasPrefix(normalized, candidate) {
			continue
		}
		addressPart := normalized[len(candidate):]
		parsed, parseErr := strconv.Atoi(addressPart)
		if parseErr != nil {
			return "", 0, fmt.Errorf("無效的 FATEK 地址數字: %s", addressPart)
		}
		return candidate, parsed, nil
	}

	return "", 0, fmt.Errorf("無法識別的 FATEK 元件符號: %s", normalized)
}
