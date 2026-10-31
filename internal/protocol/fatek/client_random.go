package fatek

import (
	"fmt"
	"strings"
)

type RandomReadItem struct {
	Symbol string
	Addr   int
}

// ReadRandom (Cmd 48) counts components in N, with an independent 64-word
// budget: a bit/16-bit component costs one word, a 32-bit component costs two.
// FATEK Communication Protocol, Appendix 1 pp5,14.
func (c *FatekClient) ReadRandom(items []RandomReadItem) (map[string]any, error) {
	if len(items) < 1 || len(items) > 64 {
		return nil, fmt.Errorf("random item count must be 1 through 64")
	}

	countHex := IntToHex(len(items), 2)
	var sb strings.Builder
	sb.WriteString(countHex)

	comps := make([]ComponentType, len(items))
	expectedLength, wordCost := 0, 0

	for i, item := range items {
		comp, err := GetComponentType(item.Symbol)
		if err != nil {
			return nil, err
		}
		comps[i] = comp
		expectedLength += max(1, comp.Width/4)
		wordCost += max(1, comp.Width/16)
		if wordCost > 64 {
			return nil, fmt.Errorf("random request exceeds 64-word cost")
		}

		addrStr, err := FormatAddress(comp, item.Addr)
		if err != nil {
			return nil, err
		}
		sb.WriteString(addrStr)
	}

	dataStr, err := c.execute("48", sb.String())
	if err != nil {
		return nil, err
	}

	if len(dataStr) != expectedLength {
		return nil, fmt.Errorf("random response length mismatch: expected %d, got %d", expectedLength, len(dataStr))
	}

	// Parse mixed response
	results := make(map[string]any)
	ptr := 0

	for i, item := range items {
		comp := comps[i]
		key := fmt.Sprintf("%s%d", item.Symbol, item.Addr)

		if comp.IsDiscrete {
			// 1 char
			if ptr+1 > len(dataStr) {
				return nil, fmt.Errorf("malformed response for item %d", i)
			}
			valChar := dataStr[ptr]
			if valChar != '0' && valChar != '1' {
				return nil, fmt.Errorf("invalid random status byte for item %d", i)
			}
			results[key] = (valChar == '1')
			ptr++
		} else {
			// Register (4 or 8 chars)
			chars := comp.Width / 4
			if ptr+chars > len(dataStr) {
				return nil, fmt.Errorf("malformed response for item %d", i)
			}
			valHex := dataStr[ptr : ptr+chars]
			val, err := HexToInt(valHex)
			if err != nil {
				return nil, err
			}
			results[key] = val
			ptr += chars
		}
	}

	return results, nil
}
