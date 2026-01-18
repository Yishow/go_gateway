package fatek

import (
	"fmt"
	"math"
	"strings"
)

// ComponentType defines the properties of a PLC component
type ComponentType struct {
	Name       string
	Width      int // 16 or 32 bits
	IsDiscrete bool
	FormatLen  int // Total length of the address string including Symbol (e.g., 5 for X, 6 for R)
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
	TypeF  = ComponentType{Name: "F", Width: 16, IsDiscrete: false, FormatLen: 6}      // File Register

	// 16-bit access to discrete (word access)
	// Use FormatLen=7 so word-access discrete symbols (e.g., WX) use 5-digit numbering like WX00012
	TypeWX = ComponentType{Name: SymbolWX, Width: 16, IsDiscrete: false, FormatLen: 7}
	TypeWY = ComponentType{Name: SymbolWY, Width: 16, IsDiscrete: false, FormatLen: 7}
	TypeWM = ComponentType{Name: SymbolWM, Width: 16, IsDiscrete: false, FormatLen: 7}
	TypeWS = ComponentType{Name: SymbolWS, Width: 16, IsDiscrete: false, FormatLen: 7}
	TypeWT = ComponentType{Name: SymbolWT, Width: 16, IsDiscrete: false, FormatLen: 7}
	TypeWC = ComponentType{Name: SymbolWC, Width: 16, IsDiscrete: false, FormatLen: 7}

	// Registers (32-bit)
	TypeDR = ComponentType{Name: SymbolDR, Width: 32, IsDiscrete: false, FormatLen: 7}
	TypeDD = ComponentType{Name: SymbolDD, Width: 32, IsDiscrete: false, FormatLen: 7}
	TypeDF = ComponentType{Name: SymbolDF, Width: 32, IsDiscrete: false, FormatLen: 7}
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
	case "F":
		return TypeF, nil
	case SymbolWX:
		return TypeWX, nil
	case SymbolWY:
		return TypeWY, nil
	case SymbolWM:
		return TypeWM, nil
	case SymbolWS:
		return TypeWS, nil
	case SymbolWT:
		return TypeWT, nil
	case SymbolWC:
		return TypeWC, nil
	case SymbolDR:
		return TypeDR, nil
	case SymbolDD:
		return TypeDD, nil
	case SymbolDF:
		return TypeDF, nil
	default:
		if strings.HasPrefix(upper, "DW") && len(upper) >= 3 {
			return ComponentType{Name: upper, Width: 32, IsDiscrete: false, FormatLen: 7}, nil
		}
		if strings.HasPrefix(upper, "D") && len(upper) >= 2 {
			return ComponentType{Name: upper, Width: 32, IsDiscrete: false, FormatLen: 7}, nil
		}
		return ComponentType{}, fmt.Errorf("unknown component type: %s", symbol)
	}
}

// FormatAddress normalizes the address string (e.g., "X", 10 -> "X00010")
// It now includes bounds checking to prevent protocol frame corruption.
func FormatAddress(comp ComponentType, addr int) (string, error) {
	// Calculate available digits for the number part
	// e.g., FormatLen(6) - Name("R", 1) = 5 digits. Max 99999.
	numDigits := comp.FormatLen - len(comp.Name)
	if numDigits <= 0 {
		return "", fmt.Errorf("invalid format length for component %s", comp.Name)
	}

	maxVal := int(math.Pow10(numDigits)) - 1
	if addr < 0 || addr > maxVal {
		return "", fmt.Errorf("address %d out of range for component %s (max %d)", addr, comp.Name, maxVal)
	}

	format := fmt.Sprintf("%%s%%0%dd", numDigits)
	return fmt.Sprintf(format, comp.Name, addr), nil
}
