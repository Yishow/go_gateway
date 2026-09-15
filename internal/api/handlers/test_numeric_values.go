package handlers

import (
	"fmt"
	"math"
)

func testRegisterValue(value any) (uint16, error) {
	switch v := value.(type) {
	case int:
		if v >= 0 && v <= math.MaxUint16 {
			return uint16(v), nil
		}
	case int64:
		if v >= 0 && v <= math.MaxUint16 {
			return uint16(v), nil
		}
	case float64:
		if v >= 0 && v <= math.MaxUint16 && math.Trunc(v) == v {
			return uint16(v), nil
		}
	}
	return 0, fmt.Errorf("register value must be an integer between 0 and 65535")
}

func testModbusUnitID(config map[string]interface{}) (byte, error) {
	value, exists := config["unitID"]
	if !exists {
		return 1, nil
	}
	parsed, err := testRegisterValue(value)
	if err != nil || parsed > math.MaxUint8 {
		return 0, fmt.Errorf("unitID must be an integer between 0 and 255")
	}
	return byte(parsed), nil
}
