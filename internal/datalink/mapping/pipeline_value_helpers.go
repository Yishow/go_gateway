package mapping

import (
	"cmp"
	"slices"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 輔助函數
// =============================================================================

func toFloat64Value(v any) (float64, bool) {
	f, err := checkedFloat64(v)
	return f, err == nil
}

func toUint16Value(v any) (uint16, bool) {
	value, err := checkedCast(v, schema.DataTypeUint16)
	if err != nil {
		return 0, false
	}
	converted, ok := value.(uint16)
	return converted, ok
}

func toUint32Value(v any) (uint32, bool) {
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

	needsOrdering := slices.ContainsFunc(steps, func(step schema.TransformStep) bool {
		return step.Order != 0
	})
	if !needsOrdering {
		return steps
	}

	ordered := slices.Clone(steps)
	slices.SortStableFunc(ordered, func(a, b schema.TransformStep) int {
		return cmp.Compare(a.Order, b.Order)
	})
	return ordered
}
