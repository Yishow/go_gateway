package snapshot

import "time"

func (a *Assembler) initializeMixed() {
	a.mixedNext = make(map[string]time.Time, len(a.members))
	for key, m := range a.members {
		start, _, _ := CheckedWindowBounds(a.config.StartsAt, m.Interval)
		a.mixedNext[key] = start
	}
}

func (a *Assembler) mixedSelector(start time.Time, member Member) *Selector {
	// Both boundaries are checked by the routing/config/plan caller before this
	// existing selector constructor is used; its addition cannot wrap here.
	s := NewSelector(SelectorConfig{Start: start, Interval: member.Interval, Members: []string{member.Key}, MaxFutureSkew: a.config.MaxFutureSkew})
	_, end, _ := a.windowBounds(start, s.end)
	s.eligible = func(sample Sample) bool {
		return meanValueMatchesSource(member.SourceValueType, sample.Value) && evaluateMember(member, sample, end).Usable()
	}
	return s
}

func (a *Assembler) checkMixed(sample Sample, now time.Time) OfferResult {
	reject := func(reason string) OfferResult { return OfferResult{Outcome: OfferRejected, Reason: reason} }
	if sample.SampleID == "" || sample.MemberKey == "" || sample.ObservedAt.IsZero() {
		return reject(ReasonMissingIdentity)
	}
	m, ok := a.members[sample.MemberKey]
	if !ok {
		return reject(ReasonUnknownMember)
	}
	if (m.SourceRevision != "" && m.SourceRevision != sample.SourceRevision) || (m.MappingRevision != "" && m.MappingRevision != sample.MappingRevision) {
		return reject(ReasonRevisionMismatch)
	}
	start, _, err := CheckedWindowBounds(sample.ObservedAt, m.Interval)
	if err != nil {
		return reject(ReasonWindowArithmetic)
	}
	if _, err = CheckedWindowAdd(now, a.config.MaxFutureSkew, "future skew"); err != nil {
		return reject(ReasonWindowArithmetic)
	}
	if start.Before(a.mixedNext[m.Key]) {
		return OfferResult{Outcome: OfferLate, Reason: ReasonLate}
	}
	if !a.WithinActualBounds(sample.ObservedAt) {
		return reject(ReasonOutsideBucket)
	}
	if state := a.buckets[start]; state != nil {
		if s := state.selectors[m.Key]; s != nil {
			return s.evaluate(sample, now)
		}
	}
	return a.mixedSelector(start, m).evaluate(sample, now)
}

func (a *Assembler) offerMixed(sample Sample, now time.Time) OfferResult {
	result := a.checkMixed(sample, now)
	switch result.Outcome {
	case OfferLate:
		a.late++
	case OfferSelected, OfferNotSelected:
		m := a.members[sample.MemberKey]
		start, _, _ := CheckedWindowBounds(sample.ObservedAt, m.Interval)
		state := a.buckets[start]
		if state == nil {
			state = &bucketState{selectors: map[string]*Selector{}}
			a.buckets[start] = state
		}
		s := state.selectors[m.Key]
		if s == nil {
			s = a.mixedSelector(start, m)
			state.selectors[m.Key] = s
		}
		s.apply(sample)
	}
	return result
}

// mixedActualSelector reuses the accepted set and selector ordering, excluding
// observations removed by an exact cutover. It never mutates buffered evidence.
func (a *Assembler) mixedActualSelector(start time.Time, m Member) *Selector {
	state := a.buckets[start]
	if state == nil {
		return nil
	}
	previous := state.selectors[m.Key]
	if previous == nil {
		return nil
	}
	s := a.mixedSelector(start, m)
	for _, sample := range previous.seen {
		if a.WithinActualBounds(sample.ObservedAt) {
			s.apply(sample)
		}
	}
	return s
}
