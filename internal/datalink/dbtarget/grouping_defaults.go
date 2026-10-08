package dbtarget

import "go-gateway/internal/datalink/common"

func defaultWriteIntervalSeconds(value *int) int {
	return common.Deref(value, 15)
}

func normalizeOptionalIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	return common.Ptr(*value)
}
