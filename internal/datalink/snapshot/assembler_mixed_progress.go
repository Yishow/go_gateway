package snapshot

import (
	"fmt"
	"maps"
	"time"
)

func (a *Assembler) mixedDone(next map[string]time.Time) bool {
	if a.config.EndsAt.IsZero() {
		return false
	}
	for _, start := range next {
		if start.Before(a.config.EndsAt) {
			return false
		}
	}
	return true
}

// MixedProgress returns typed dispatch progress only for a mixed assembler.
func (a *Assembler) MixedProgress() (DispatchWatermark, bool) {
	if !a.config.Mode.mixed() {
		return DispatchWatermark{}, false
	}
	return DispatchWatermark{a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, a.mixedThrough}, true
}

// MixedProgress is the watermark installed atomically with this plan's results.
func (p Plan) MixedProgress() (DispatchWatermark, bool) {
	if p.mixedProgress == nil {
		return DispatchWatermark{}, false
	}
	return *p.mixedProgress, true
}

// MemberProgress returns a fresh projection, never another writable authority.
func (a *Assembler) MemberProgress() (map[string]MemberProgress, error) {
	if !a.config.Mode.mixed() {
		return nil, fmt.Errorf("%w: not mixed", ErrInvalidConfig)
	}
	result := make(map[string]MemberProgress, len(a.mixedNext))
	for key, start := range a.mixedNext {
		if !a.config.EndsAt.IsZero() && !start.Before(a.config.EndsAt) {
			result[key] = MemberProgress{Start: start}
			continue
		}
		end, err := CheckedWindowAdd(start, a.members[key].Interval, "member progress")
		if err != nil {
			return nil, err
		}
		result[key] = MemberProgress{start, end}
	}
	return result, nil
}

// RestoreMixedProgress is allowed before accepting samples or committing any
// closure. The durable watermark derives all member starts on their old phase;
// callers replay pending accepted samples after restoring it, not before.
func (a *Assembler) RestoreMixedProgress(p DispatchWatermark) error {
	invalid := func(s string) error { return fmt.Errorf("%w: %s", ErrInvalidConfig, s) }
	if !a.config.Mode.mixed() || len(a.buckets) > 0 || !a.mixedThrough.IsZero() {
		return invalid("mixed progress restore requires untouched assembler")
	}
	if p.WorkspaceID != a.config.WorkspaceID || p.GroupID != a.config.GroupID || p.GroupRevision != a.config.GroupRevision {
		return invalid("foreign mixed progress")
	}
	if p.ClosedThrough.IsZero() {
		return nil
	}
	if !p.ClosedThrough.After(a.config.StartsAt) {
		return invalid("dispatch must follow activation")
	}
	next := make(map[string]time.Time, len(a.members))
	onBoundary := false
	for key, m := range a.members {
		start, _, err := CheckedWindowBounds(p.ClosedThrough, m.Interval)
		if err != nil {
			return err
		}
		memberBoundary := start.Equal(p.ClosedThrough)
		// Ended members stop at their last aligned end, not arbitrary later events.
		if !a.config.EndsAt.IsZero() && !start.Before(a.config.EndsAt) {
			bound, end, err := CheckedWindowBounds(a.config.EndsAt, m.Interval)
			if err != nil {
				return err
			}
			if !bound.Equal(a.config.EndsAt) {
				bound = end
			}
			if p.ClosedThrough.After(bound) {
				memberBoundary = false
			}
			start = bound
		}
		onBoundary = onBoundary || memberBoundary
		next[key] = start
	}
	if !onBoundary {
		return invalid("dispatch is not a member boundary")
	}
	if !a.config.EndsAt.IsZero() {
		last := time.Time{}
		for _, m := range a.members {
			start, end, err := CheckedWindowBounds(a.config.EndsAt, m.Interval)
			if err != nil {
				return err
			}
			if start.Equal(a.config.EndsAt) {
				end = start
			}
			if end.After(last) {
				last = end
			}
		}
		if p.ClosedThrough.After(last) {
			return invalid("dispatch exceeds last closure")
		}
	}
	a.mixedNext = next
	a.mixedThrough = p.ClosedThrough.UTC()
	return nil
}

// SetMixedEndsAt narrows exact mixed coverage without inventing a common aligned
// Until. An invalid or already-closed cutover is rejected without mutation.
func (a *Assembler) SetMixedEndsAt(end time.Time) error {
	if !a.config.Mode.mixed() {
		return fmt.Errorf("%w: not mixed", ErrInvalidConfig)
	}
	if !end.After(a.config.StartsAt) || (!a.mixedThrough.IsZero() && end.Before(a.mixedThrough)) {
		return fmt.Errorf("%w: end precedes active progress", ErrInvalidConfig)
	}
	for _, m := range a.members {
		if _, _, err := CheckedWindowBounds(end, m.Interval); err != nil {
			return err
		}
	}
	if !a.config.EndsAt.IsZero() && !end.Before(a.config.EndsAt) {
		return nil
	}
	a.config.EndsAt = end.UTC()
	for start, state := range a.buckets {
		for key := range state.selectors {
			s := a.mixedActualSelector(start, a.members[key])
			if s == nil || len(s.seen) == 0 {
				delete(state.selectors, key)
			} else {
				state.selectors[key] = s
			}
		}
		if len(state.selectors) == 0 {
			delete(a.buckets, start)
		}
	}
	return nil
}

func (a *Assembler) commitMixed(p Plan) {
	if p.Err != nil || p.mixedProgress == nil || len(p.mixedClosures) == 0 {
		return
	}
	if !p.mixedBase.Equal(a.mixedThrough) || p.mixedProgress.WorkspaceID != a.config.WorkspaceID || p.mixedProgress.GroupID != a.config.GroupID || p.mixedProgress.GroupRevision != a.config.GroupRevision {
		return
	}
	for _, closure := range p.mixedClosures {
		if state := a.buckets[closure.start]; state != nil {
			delete(state.selectors, closure.member)
			if len(state.selectors) == 0 {
				delete(a.buckets, closure.start)
			}
		}
	}
	a.mixedNext = maps.Clone(p.mixedNext)
	a.mixedThrough = p.mixedProgress.ClosedThrough
	for key, checkpoint := range p.counters {
		a.counters[key] = checkpoint
	}
}
