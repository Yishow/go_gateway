package runtimeview

import "time"

func afterTime(a, b string) bool {
	x, _ := time.Parse(time.RFC3339Nano, a)
	y, _ := time.Parse(time.RFC3339Nano, b)
	return x.After(y)
}
func sameTime(a, b string) bool {
	x, _ := time.Parse(time.RFC3339Nano, a)
	y, _ := time.Parse(time.RFC3339Nano, b)
	return x.Equal(y)
}

// Validate only provided facts; never synthesize defaults or missing interval/reduction.
func validateSummaryProfile(f *GroupFacts) error {
	p := f.RecordingProfile
	if p == nil {
		return nil
	}
	seconds := p.Interval
	if p.Kind == "instantaneous" {
		if p.Unit != "second" || p.Interval > 300 || (f.Reduction != "unavailable" && f.Reduction != "latest" && f.Reduction != "mean") {
			return errUnavailable
		}
	} else {
		seconds *= 60
		if p.Unit != "minute" || (f.Reduction != "unavailable" && f.Reduction != "cumulative_delta") {
			return errUnavailable
		}
	}
	if f.IntervalSeconds != nil && *f.IntervalSeconds != seconds {
		return errUnavailable
	}
	return nil
}

func validateGroup(g *GroupView, scope Scope, observedAt *string) error {
	revision := g.EffectiveRevision
	// This fallback only validates legacy consistency; it never establishes effective proof.
	if revision == nil && g.Effective != nil {
		revision = &g.Effective.Binding.EntityRevision
	}
	if g.AsOf != nil && observedAt != nil && !sameTime(*g.AsOf, *observedAt) {
		return errStale
	}
	if g.Effective != nil && g.Effective.EffectiveAt != nil && g.AsOf != nil && afterTime(*g.Effective.EffectiveAt, *g.AsOf) {
		return errStale
	}

	selected := false
	for _, pair := range []struct {
		facts    *GroupFacts
		revision *string
	}{{g.Saved, g.SavedRevision}, {g.Effective, revision}} {
		if pair.facts == nil {
			continue
		}
		if err := validateSummaryProfile(pair.facts); err != nil {
			return err
		}
		if err := checkBinding(pair.facts.Binding, pair.revision, scope); err != nil {
			return err
		}
		members := map[string]bool{}
		for _, m := range pair.facts.Members {
			key := m.DeviceID + ":" + m.PointID + ":" + m.TagID
			if members[key] {
				return errUnavailable
			}
			members[key] = true
			if m.DeviceID == scope.DeviceID {
				selected = true
			}
		}
	}
	if !selected && (g.Saved != nil || g.Effective != nil) {
		return errForeign
	}
	if g.State == "saved-only" && (g.Saved == nil || g.Effective != nil) {
		return errUnavailable
	}
	if g.State != "effective" && g.Effective != nil {
		return errUnavailable
	}
	// Validate mixed revisions before degrading proof, so missing proof cannot hide mixed inputs.
	if g.Delivery != nil {
		if g.Delivery.CounterReplayQuarantines != nil {
			seen := map[string]bool{}
			for _, q := range *g.Delivery.CounterReplayQuarantines {
				if q.WorkspaceID != scope.WorkspaceID || q.GroupID != g.GroupID {
					return errForeign
				}
				if seen[q.GroupRevision] {
					return errUnavailable
				}
				seen[q.GroupRevision] = true
				if q.Reason != ReasonCounterMissing && q.Reason != ReasonCounterConflict && q.Reason != ReasonCounterUnsupported {
					return errUnavailable
				}
			}
		}
		for _, b := range g.Delivery.Backlog {
			if err := checkBinding(b.Binding, &b.Binding.EntityRevision, scope); err != nil {
				return err
			}
			if revision != nil && b.Binding.EntityRevision == *revision {
				return errStale
			}
		}
		if g.Delivery.Current != nil {
			if err := checkBinding(g.Delivery.Current.Binding, revision, scope); err != nil {
				return err
			}
			c := g.Delivery.Current
			if c.Intake != nil && c.Intake.State == "active" && g.Effective != nil && (!g.Effective.Enabled || (g.Effective.EndsAt != nil && g.AsOf != nil && !afterTime(*g.Effective.EndsAt, *g.AsOf))) {
				return errUnavailable
			}
			if c.Intake != nil && c.Intake.GroupRevision != c.Binding.EntityRevision {
				return errStale
			}
			receipts := []*Receipt{c.LastCommit}
			for _, source := range c.Sources {
				for _, r := range []*SourceResult{source.Latest, source.LastCommitted} {
					if r != nil {
						receipts = append(receipts, r.Receipt)
					}
				}
			}
			identities := map[string]Receipt{}
			for _, r := range receipts {
				if r == nil {
					continue
				}
				if r.GroupRevision != *revision {
					return errStale
				}
				for _, key := range []string{"effect:" + r.EffectKey, "record:" + r.RecordID} {
					if prior, exists := identities[key]; exists && prior != *r {
						return errStale
					}
					identities[key] = *r
				}
			}
			if g.Effective != nil {
				if err := validateCurrent(c, g.Effective); err != nil {
					return err
				}
			}
		}
	}
	if g.State == "effective" && (g.EffectiveRevision == nil || g.AsOf == nil || observedAt == nil || g.Proof != "applied-snapshot-and-layout" || g.Effective == nil || revision == nil || g.Effective.EffectiveAt == nil || g.Effective.Layout == "unavailable" || g.Effective.IntervalSeconds == nil) {
		g.State = "unavailable"
		g.Effective = nil
		g.Proof = "unavailable"
	}
	if g.Delivery == nil || g.Delivery.Current == nil {
		return nil
	}
	if g.Effective == nil {
		g.Delivery.Current = nil
		return nil
	}
	return nil
}

func checkReceipt(r *Receipt, facts *GroupFacts) error {
	if r == nil {
		return nil
	}
	if r.GroupRevision != facts.Binding.EntityRevision || r.ConnectorRevision != facts.Destination.ConnectorRevision {
		return errStale
	}
	return nil
}
func validateCurrent(c *CurrentDelivery, facts *GroupFacts) error {
	if err := checkReceipt(c.LastCommit, facts); err != nil {
		return err
	}
	if (c.Stages.SQLCommitted > 0) != (c.LastCommit != nil) {
		return errUnavailable
	}
	devices := map[string]bool{}
	members := map[string]bool{}
	for _, m := range facts.Members {
		devices[m.DeviceID] = true
		members[m.DeviceID+":"+m.PointID+":"+m.TagID] = true
	}
	if c.Shared != (len(devices) > 1) || (c.SourceState == "unavailable" && len(c.Sources) > 0) {
		return errUnavailable
	}
	seen := map[string]bool{}
	for _, s := range c.Sources {
		key := s.DeviceID + ":" + s.PointID + ":" + s.TagID
		if !members[key] {
			return errForeign
		}
		if seen[key] {
			return errUnavailable
		}
		seen[key] = true
		for _, r := range []*SourceResult{s.Latest, s.LastCommitted} {
			if r == nil {
				continue
			}
			if err := checkReceipt(r.Receipt, facts); err != nil {
				return err
			}
			if (r.WindowStart == nil) != (r.WindowEnd == nil) {
				return errUnavailable
			}
			if r.WindowStart != nil {
				a, _ := time.Parse(time.RFC3339Nano, *r.WindowStart)
				b, _ := time.Parse(time.RFC3339Nano, *r.WindowEnd)
				if !a.Before(b) {
					return errUnavailable
				}
			}
			if (r.DeltaKnown != nil && !*r.DeltaKnown && r.State != "unknown") || (r.Coverage != nil && *r.Coverage == "partial" && r.State != "partial") {
				return errUnavailable
			}
			if r.Receipt != nil && (r.WindowStart == nil || r.State == "unavailable" || c.LastCommit == nil) {
				return errUnavailable
			}
		}
		if s.LastCommitted != nil && s.LastCommitted.Receipt == nil {
			return errUnavailable
		}
	}
	return nil
}
