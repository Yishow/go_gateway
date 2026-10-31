package workspace

import (
	"fmt"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// mixedApplyCutover schedules ownership only. Its common cadence is never a
// row policy or a member window; each member retains its original UTC phase.
func mixedApplyCutover(candidate *WriteGroup, active *WriteGroupAppliedSnapshot) (time.Time, error) {
	intervals, err := writeGroupSchedulingCadences(candidate)
	if err != nil {
		return time.Time{}, err
	}
	if active == nil {
		at, err := snapshot.CheckedWindowAdd(candidate.UpdatedAt, 0, "mixed activation")
		if err != nil {
			return time.Time{}, fmt.Errorf("%w: %w", ErrWriteGroupIntervalOverflow, err)
		}
		return at, nil
	}
	old, err := writeGroupSchedulingCadences(active.Group)
	if err != nil {
		return time.Time{}, err
	}
	intervals = append(intervals, old...)
	common := intervals[0]
	for _, interval := range intervals[1:] {
		common, err = writeGroupCommonInterval(common, interval)
		if err != nil {
			return time.Time{}, err
		}
	}
	_, at, err := snapshot.CheckedWindowBounds(candidate.UpdatedAt, common)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %w", ErrWriteGroupIntervalOverflow, err)
	}
	return at, nil
}
func writeGroupSchedulingCadences(group *WriteGroup) ([]time.Duration, error) {
	if group == nil {
		return nil, ErrWriteGroupValidation
	}
	if !isMixedWriteGroup(group) {
		interval, err := writeGroupIntervalDuration(group.RowPolicy.IntervalSeconds)
		return []time.Duration{interval}, err
	}
	if len(group.Members) == 0 || group.RowPolicy.IntervalSeconds != 0 {
		return nil, ErrWriteGroupValidation
	}
	intervals := make([]time.Duration, 0, len(group.Members))
	for _, m := range group.Members {
		if m.Cadence == nil {
			return nil, ErrWriteGroupValidation
		}
		interval, err := m.Cadence.Duration()
		if err != nil {
			return nil, err
		}
		intervals = append(intervals, interval)
	}
	return intervals, nil
}
