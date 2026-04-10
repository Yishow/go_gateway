package dbtarget

func defaultWriteIntervalSeconds(value *int) int {
	if value == nil {
		return 15
	}
	return *value
}

func normalizeOptionalIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	normalized := *value
	return &normalized
}
