package mapping

import (
	"sort"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 輔助函數
// =============================================================================

func toFloat64Value(v interface{}) (float64, bool) {
	f, err := checkedFloat64(v)
	return f, err == nil
}

func toUint16Value(v interface{}) (uint16, bool) {
	value, err := checkedCast(v, schema.DataTypeUint16)
	if err != nil {
		return 0, false
	}
	converted, ok := value.(uint16)
	return converted, ok
}

func toUint32Value(v interface{}) (uint32, bool) {
	value, err := checkedCast(v, schema.DataTypeUint32)
	if err != nil {
		return 0, false
	}
	converted, ok := value.(uint32)
	return converted, ok
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
