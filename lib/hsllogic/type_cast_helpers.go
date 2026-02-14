package hsllogic

// ToFloat64 將任意數值轉換為 float64
func ToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case int8:
		return float64(v), true
	case uint8:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// ToInt64 將任意數值轉換為 int64
func ToInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case int8:
		return int64(v), true
	case uint8:
		return int64(v), true
	case int16:
		return int64(v), true
	case uint16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint32:
		return int64(v), true
	case int64:
		return v, true
	case uint64:
		return int64(v), true
	case int:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
