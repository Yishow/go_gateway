package snapshot

import "time"

// WithinActualBounds reports the exact half-open observation range of this
// assembler. It is an in-memory eligibility filter, not durable owner authority.
// Zero bounds preserve ordinary, interval-aligned assembly behavior.
func (a *Assembler) WithinActualBounds(observed time.Time) bool {
	return (a.config.StartsAt.IsZero() || !observed.Before(a.config.StartsAt)) &&
		(a.config.EndsAt.IsZero() || observed.Before(a.config.EndsAt))
}

// rebuildActualBounds retains the complete eligible accepted set, not merely
// the previous selected sample. This lets an earlier candidate become latest
// when a cutover removes a newer observation. Durable data is never changed.
func (a *Assembler) rebuildActualBounds() {
	for start, state := range a.buckets {
		for entity, previous := range state.selectors {
			next := a.selectorWithinActualBounds(start, entity, previous)
			if len(next.seen) == 0 {
				delete(state.selectors, entity)
			} else {
				state.selectors[entity] = next
			}
		}
		if len(state.selectors) == 0 {
			delete(a.buckets, start)
		}
	}
}

// selectorWithinActualBounds builds a closure-only projection without mutating
// the accepted buffer. Pure Assembler callers retain their existing Offer
// contract; accepting into this utility is not a durable recording ACK.
func (a *Assembler) selectorWithinActualBounds(start time.Time, entity string, previous *Selector) *Selector {
	if previous == nil || ((a.config.StartsAt.IsZero() || !start.Before(a.config.StartsAt)) &&
		(a.config.EndsAt.IsZero() || !start.Add(a.config.Interval).After(a.config.EndsAt))) {
		return previous
	}
	next := a.newSelector(start, entity)
	for _, sample := range previous.seen {
		if a.WithinActualBounds(sample.ObservedAt) {
			next.apply(sample)
		}
	}
	return next
}
