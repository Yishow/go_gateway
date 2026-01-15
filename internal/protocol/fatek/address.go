package fatek

import (
	"fmt"
	"strings"
)

// ComponentType defines the properties of a PLC component
type ComponentType struct {
	Name      string
	Width     int // 16 or 32 bits
	IsDiscrete bool
	FormatLen int // Length of the address string (e.g., 5 for X, 6 for R)
}

var (
	// Discrete (1-bit)
	TypeX = ComponentType{Name: SymbolX, Width: 1, IsDiscrete: true, FormatLen: 5}
	TypeY = ComponentType{Name: SymbolY, Width: 1, IsDiscrete: true, FormatLen: 5}
	TypeM = ComponentType{Name: SymbolM, Width: 1, IsDiscrete: true, FormatLen: 5}
	TypeS = ComponentType{Name: SymbolS, Width: 1, IsDiscrete: true, FormatLen: 5}
	TypeT = ComponentType{Name: SymbolT, Width: 1, IsDiscrete: true, FormatLen: 5} // Timer Status
	TypeC = ComponentType{Name: SymbolC, Width: 1, IsDiscrete: true, FormatLen: 5} // Counter Status

	// Registers (16-bit)
	TypeR  = ComponentType{Name: SymbolR, Width: 16, IsDiscrete: false, FormatLen: 6}
	TypeD  = ComponentType{Name: SymbolD, Width: 16, IsDiscrete: false, FormatLen: 6}
	TypeRT = ComponentType{Name: SymbolRT, Width: 16, IsDiscrete: false, FormatLen: 6} // Timer Value
	TypeRC = ComponentType{Name: SymbolRC, Width: 16, IsDiscrete: false, FormatLen: 6} // Counter Value

	// Registers (32-bit)
	TypeDR = ComponentType{Name: SymbolDR, Width: 32, IsDiscrete: false, FormatLen: 7}
	// Note: Fatek often uses D as 32-bit with specific commands, but explicit DR is safer for mixed reads.
)

func GetComponentType(symbol string) (ComponentType, error) {
	upper := strings.ToUpper(symbol)
	switch upper {
	case SymbolX:
		return TypeX, nil
	case SymbolY:
		return TypeY, nil
	case SymbolM:
		return TypeM, nil
	case SymbolS:
		return TypeS, nil
	case SymbolT:
		return TypeT, nil
	case SymbolC:
		return TypeC, nil
	case SymbolR:
		return TypeR, nil
	case SymbolD:
		return TypeD, nil
	case SymbolRT:
		return TypeRT, nil
	case SymbolRC:
		return TypeRC, nil
	case SymbolDR, "DD":
		return TypeDR, nil
	default:
		// Defaulting DWM, etc. to 32-bit if needed, but keeping it strict for now
		if strings.HasPrefix(upper, "D") {
			return ComponentType{Name: upper, Width: 32, IsDiscrete: false, FormatLen: 7}, nil
		}
		return ComponentType{}, fmt.Errorf("unknown component type: %s", symbol)
	}
}

// FormatAddress normalizes the address string (e.g., "X", 10 -> "X00010")
func FormatAddress(comp ComponentType, addr int) string {
	format := fmt.Sprintf("%%s%%0%dd", comp.FormatLen-len(comp.Name))
	return fmt.Sprintf(format, comp.Name, addr)
}