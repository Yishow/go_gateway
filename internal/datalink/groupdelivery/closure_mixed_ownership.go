package groupdelivery

import (
	"go-gateway/internal/datalink/aggregation"
	"time"
)

// ValidateMixedBounds checks exact ownership coverage, not a shared window.
func (v OwnershipView) ValidateMixedBounds(key GroupKey, startsAt, endsAt time.Time) error {
	if v == (OwnershipView{}) {
		return nil
	}
	if !v.mixed || v.key != key || !startsAt.Equal(v.start) {
		return ErrOwnershipConflict
	}
	if !v.end.IsZero() && (endsAt.IsZero() || endsAt.After(v.end)) {
		return ErrOwnershipConflict
	}
	return nil
}
func (v OwnershipView) validateMixedClosure(key GroupKey, c Closure) error {
	if err := v.ValidateMixedBounds(key, c.Mixed.After.StartsAt, c.Mixed.After.EndsAt); err != nil {
		return err
	}
	for _, b := range c.Buckets {
		for _, m := range b.Outcome.Members {
			if m.Window == nil || m.Window.ActualStart.IsZero() {
				continue
			}
			w := m.Window
			if w.ActualStart.Before(v.start) || (!v.end.IsZero() && w.ActualEnd.After(v.end)) {
				return ErrOwnershipConflict
			}
		}
	}
	for _, baseline := range c.Baselines {
		cp, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if err != nil || cp.WindowEnd.Before(v.start) || (!v.end.IsZero() && cp.WindowEnd.After(v.end)) || cp.Reading.ObservedAt.Before(v.start) || (!v.end.IsZero() && !cp.Reading.ObservedAt.Before(v.end)) {
			return ErrOwnershipConflict
		}
	}
	return nil
}

// AuthorizeMixedEnd proves an exact canonical cutover end from the immutable
// ownership view captured by Restore. Absence of a fence is not authority to
// rewrite an already-persisted checkpoint's coverage.
func (v OwnershipView) AuthorizeMixedEnd(key GroupKey, startsAt, endsAt time.Time) error {
	if v == (OwnershipView{}) || !v.mixed || v.key != key || !startsAt.Equal(v.start) || v.end.IsZero() || !endsAt.Equal(v.end) {
		return ErrOwnershipConflict
	}
	return nil
}
