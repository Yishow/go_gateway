package mapping

import (
	"encoding/json"
	"sort"
	"strconv"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 輔助函數
// =============================================================================

func toFloat64Value(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return f, true
		}
	}
	return 0, false
}

func toUint16Value(v interface{}) (uint16, bool) {
	f, ok := toFloat64Value(v)
	if ok {
		return uint16(f), true
	}
	return 0, false
}

func toUint32Value(v interface{}) (uint32, bool) {
	f, ok := toFloat64Value(v)
	if ok {
		return uint32(f), true
	}
	return 0, false
}

func toBoolValue(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int, int16, int32, int64, uint, uint16, uint32, uint64, float32, float64:
		f, _ := toFloat64Value(v)
		return f != 0
	case string:
		return val == "true" || val == "1" || val == "on" || val == "yes"
	}
	return false
}

func normalizeTransformSteps(steps []schema.TransformStep) []schema.TransformStep {
	if len(steps) == 0 {
		return steps
	}

	needsOrdering := false
	for _, step := range steps {
		if step.Order != 0 {
			needsOrdering = true
			break
		}
	}
	if !needsOrdering {
		return steps
	}

	ordered := make([]schema.TransformStep, len(steps))
	copy(ordered, steps)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Order < ordered[j].Order
	})
	return ordered
}
