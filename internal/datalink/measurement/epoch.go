package measurement

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CheckEpochTransition 評估量測定義變更是否破壞數值可比性，若破壞則產生新 series_epoch。
func CheckEpochTransition(oldDef, newDef MeasurementDefinition) (bool, string, string) {
	currentEpoch := oldDef.SeriesEpoch
	if strings.TrimSpace(currentEpoch) == "" {
		currentEpoch = "epoch-1"
	}

	reasons := make([]string, 0)

	// 檢查破壞可比性的關鍵欄位
	if oldDef.DeviceID != newDef.DeviceID || oldDef.PointID != newDef.PointID {
		reasons = append(reasons, "source point or device binding changed")
	}
	if oldDef.Quantity != newDef.Quantity && oldDef.Quantity != "" && newDef.Quantity != "" {
		reasons = append(reasons, fmt.Sprintf("quantity changed from %s to %s", oldDef.Quantity, newDef.Quantity))
	}
	if oldDef.Unit != newDef.Unit && oldDef.Unit != "" && newDef.Unit != "" {
		reasons = append(reasons, fmt.Sprintf("unit changed from %s to %s", oldDef.Unit, newDef.Unit))
	}
	if oldDef.SemanticKind != newDef.SemanticKind && oldDef.SemanticKind != "" && newDef.SemanticKind != "" {
		reasons = append(reasons, fmt.Sprintf("semantic_kind changed from %s to %s", oldDef.SemanticKind, newDef.SemanticKind))
	}

	// 檢查 CounterPolicy 變更
	if (oldDef.CounterPolicy == nil && newDef.CounterPolicy != nil) ||
		(oldDef.CounterPolicy != nil && newDef.CounterPolicy == nil) {
		reasons = append(reasons, "counter policy modified")
	} else if oldDef.CounterPolicy != nil && newDef.CounterPolicy != nil {
		if oldDef.CounterPolicy.AllowNegative != newDef.CounterPolicy.AllowNegative {
			reasons = append(reasons, "counter policy allow_negative changed")
		}
		if !floatPtrEqual(oldDef.CounterPolicy.MaxRollOver, newDef.CounterPolicy.MaxRollOver) {
			reasons = append(reasons, "counter policy max_roll_over changed")
		}
		if !floatPtrEqual(oldDef.CounterPolicy.ResetThreshold, newDef.CounterPolicy.ResetThreshold) {
			reasons = append(reasons, "counter policy reset_threshold changed")
		}
	}

	if len(reasons) == 0 {
		return false, currentEpoch, ""
	}

	nextEpoch := IncrementEpoch(currentEpoch)
	return true, nextEpoch, strings.Join(reasons, "; ")
}

// IncrementEpoch 將 "epoch-1" 遞增為 "epoch-2"，或建立新 epoch 序號。
func IncrementEpoch(epoch string) string {
	epoch = strings.TrimSpace(epoch)
	if strings.HasPrefix(epoch, "epoch-") {
		numStr := strings.TrimPrefix(epoch, "epoch-")
		if n, err := strconv.Atoi(numStr); err == nil {
			return fmt.Sprintf("epoch-%d", n+1)
		}
	}
	return "epoch-2"
}

func floatPtrEqual(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return math.Abs(*a-*b) < 1e-9
}
