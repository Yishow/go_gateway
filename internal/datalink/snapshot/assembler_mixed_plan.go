package snapshot

import (
	"go-gateway/internal/datalink/aggregation"
	"maps"
	"time"
)

type mixedClosure struct {
	member string
	start  time.Time
}

func (a *Assembler) nextMixedDue(next map[string]time.Time) (time.Time, error) {
	var due time.Time
	for key, start := range next {
		if !a.config.EndsAt.IsZero() && !start.Before(a.config.EndsAt) {
			continue
		}
		end, err := CheckedWindowAdd(start, a.members[key].Interval, "due")
		if err != nil {
			return time.Time{}, err
		}
		if due.IsZero() || end.Before(due) {
			due = end
		}
	}
	return due, nil
}

func (a *Assembler) planMixed(now time.Time) Plan {
	if _, err := CheckedWindowAdd(now, 0, "closure time"); err != nil {
		return Plan{Err: err}
	}
	progress, _ := a.MixedProgress()
	p := Plan{mixedBase: a.mixedThrough, mixedProgress: &progress, mixedNext: maps.Clone(a.mixedNext), counters: map[string]*aggregation.CounterCheckpoint{}}
	limit := a.config.MaxBucketsPerTick
	if a.config.Mode == ModeMixedLong {
		limit = max(1, limit/len(a.members))
	}
	for event := 0; event < limit; event++ {
		due, err := a.nextMixedDue(p.mixedNext)
		if err != nil {
			return Plan{Err: err}
		}
		if due.IsZero() {
			break
		}
		closesAt, err := CheckedWindowAdd(due, a.config.AllowedLateness, "lateness")
		if err != nil {
			return Plan{Err: err}
		}
		if now.Before(closesAt) {
			break
		}
		wide := Outcome{Kind: OutcomeRow, DueAt: due}
		for _, m := range a.config.Members {
			start := p.mixedNext[m.Key]
			end, err := CheckedWindowAdd(start, m.Interval, "member end")
			if err != nil {
				return Plan{Err: err}
			}
			isDue := end.Equal(due) && (a.config.EndsAt.IsZero() || start.Before(a.config.EndsAt))
			if !isDue {
				if !a.config.EndsAt.IsZero() && !start.Before(a.config.EndsAt) {
					end = time.Time{}
				}
				if a.config.Mode == ModeMixedWide {
					wide.Members = append(wide.Members, MemberResult{MemberKey: m.Key, Status: MemberNotDue, Reason: ReasonNotDue, Window: &MemberWindow{Calculation: m.Calculation, Interval: m.Interval, NextDue: end}})
					wide.Partial = true
				}
				continue
			}
			nextDue, err := CheckedWindowAdd(end, m.Interval, "next due")
			if err != nil {
				return Plan{Err: err}
			}
			if !a.config.EndsAt.IsZero() && !end.Before(a.config.EndsAt) {
				nextDue = time.Time{}
			}
			result := a.mixedMemberResult(&p, m, start, end, nextDue)
			p.mixedClosures = append(p.mixedClosures, mixedClosure{m.Key, start})
			p.mixedNext[m.Key] = end
			if a.config.Mode == ModeMixedWide {
				wide.Members = append(wide.Members, result)
				if !result.Usable() {
					wide.Partial = true
					if m.Required {
						wide.Kind, wide.Reason = OutcomeSkipped, ReasonIncompleteRequired
					}
				}
			} else {
				out := Outcome{Kind: OutcomeRow, EntityKey: m.Key, DueAt: due, BucketStart: start, BucketEnd: end, Members: []MemberResult{result}}
				if !result.Usable() {
					out.Partial = true
					if m.Required {
						out.Kind, out.Reason = OutcomeSkipped, ReasonIncompleteRequired
					}
				}
				out = a.markTruncation(out)
				a.mixedIdentity(&out, m.Key, start)
				p.Outcomes = append(p.Outcomes, out)
			}
		}
		if a.config.Mode == ModeMixedWide {
			a.mixedIdentity(&wide, "", due)
			p.Outcomes = append(p.Outcomes, wide)
		}
		p.mixedProgress.ClosedThrough = due
	}
	return p
}
