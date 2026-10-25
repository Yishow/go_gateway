package snapshot

import (
	"cmp"
	"slices"
	"time"

	"go-gateway/internal/datalink/aggregation"
)

// ReasonLate marks a sample whose bucket already closed.
const ReasonLate = "late-after-close"

// OfferLate is returned for a sample that belongs to a closed bucket.
const OfferLate OfferOutcome = "late"

type entityGroup struct {
	key     string
	members []Member
}

type bucketState struct {
	selectors map[string]*Selector
}

// Assembler turns samples into per-bucket outcomes for one applied group
// revision. State lives in memory only: it is not durable and is lost on
// restart until the delivery change journals it.
type Assembler struct {
	config    Config
	entities  []entityGroup
	members   map[string]Member
	late      uint64
	buckets   map[time.Time]*bucketState
	nextClose time.Time
	// counters are the committed cumulative baselines by member key; they move
	// only through Commit, like nextClose.
	counters map[string]*aggregation.CounterCheckpoint
}

// NewAssembler validates the configuration.
func NewAssembler(config Config) (*Assembler, error) {
	config, err := config.normalized()
	if err != nil {
		return nil, err
	}
	byEntity := make(map[string][]Member)
	members := make(map[string]Member, len(config.Members))
	for _, member := range config.Members {
		byEntity[config.entityOf(member)] = append(byEntity[config.entityOf(member)], member)
		members[member.Key] = member
	}
	entities := make([]entityGroup, 0, len(byEntity))
	for key, members := range byEntity {
		entities = append(entities, entityGroup{key: key, members: members})
	}
	slices.SortFunc(entities, func(a, b entityGroup) int { return cmp.Compare(a.key, b.key) })
	return &Assembler{
		config:    config,
		entities:  entities,
		members:   members,
		buckets:   make(map[time.Time]*bucketState),
		nextClose: config.FirstBucket,
		counters:  make(map[string]*aggregation.CounterCheckpoint),
	}, nil
}

// Offer routes a sample to its bucket at gateway time now.
func (a *Assembler) Offer(sample Sample, now time.Time) OfferResult {
	result := a.Check(sample, now)
	switch result.Outcome {
	case OfferLate:
		a.late++
	case OfferSelected, OfferNotSelected:
		member := a.members[sample.MemberKey]
		start := BucketStart(sample.ObservedAt, a.config.Interval)
		state, ok := a.buckets[start]
		if !ok {
			state = &bucketState{selectors: make(map[string]*Selector)}
			a.buckets[start] = state
		}
		entity := a.config.entityOf(member)
		selector, ok := state.selectors[entity]
		if !ok {
			selector = a.newSelector(start, entity)
			state.selectors[entity] = selector
		}
		selector.apply(sample)
	}
	return result
}

func (a *Assembler) newSelector(start time.Time, entity string) *Selector {
	selector := NewSelector(SelectorConfig{
		Start:         start,
		Interval:      a.config.Interval,
		MaxFutureSkew: a.config.MaxFutureSkew,
		Members:       a.memberKeys(entity),
	})
	if a.config.Mode.perMember() || a.config.Mode.wide() {
		selector.eligible = a.eligibility(start.Add(a.config.Interval))
	}
	return selector
}

// Check reports what Offer would do without changing any state, including the
// late counter, so a caller can persist a sample before accepting it.
func (a *Assembler) Check(sample Sample, now time.Time) OfferResult {
	if sample.SampleID == "" || sample.MemberKey == "" || sample.ObservedAt.IsZero() {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonMissingIdentity}
	}
	member, ok := a.members[sample.MemberKey]
	if !ok {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonUnknownMember}
	}
	if (member.SourceRevision != "" && sample.SourceRevision != member.SourceRevision) ||
		(member.MappingRevision != "" && sample.MappingRevision != member.MappingRevision) {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonRevisionMismatch}
	}
	start := BucketStart(sample.ObservedAt, a.config.Interval)
	if start.Before(a.nextClose) {
		return OfferResult{Outcome: OfferLate, Reason: ReasonLate}
	}
	entity := a.config.entityOf(member)
	if state, ok := a.buckets[start]; ok {
		if selector, ok := state.selectors[entity]; ok {
			return selector.evaluate(sample, now)
		}
	}
	return a.newSelector(start, entity).evaluate(sample, now)
}

// Replay re-offers a sample that was already durably accepted, bypassing the
// future-skew check because its acceptance was decided when it first arrived.
func (a *Assembler) Replay(sample Sample) OfferResult {
	return a.Offer(sample, sample.ObservedAt)
}

// Tick closes every bucket whose end plus allowed lateness has passed,
// including buckets that never received a sample. At most MaxBucketsPerTick
// buckets close per call; the rest follow on later ticks, in order.
func (a *Assembler) Tick(now time.Time) []Outcome {
	plan := a.Plan(now)
	a.Commit(plan)
	return plan.Outcomes
}

// Plan computes the closure Tick would perform at now without applying it.
func (a *Assembler) Plan(now time.Time) Plan {
	plan := Plan{next: a.nextClose, counters: make(map[string]*aggregation.CounterCheckpoint)}
	for limit := a.bucketsPerTick(); len(plan.buckets) < limit; {
		if !a.config.Until.IsZero() && !plan.next.Before(a.config.Until) {
			break
		}
		end := plan.next.Add(a.config.Interval)
		if now.Before(end.Add(a.config.AllowedLateness)) {
			break
		}
		state := a.buckets[plan.next]
		for _, entity := range a.entities {
			var selector *Selector
			if state != nil {
				selector = state.selectors[entity.key]
			}
			outcome := a.finalize(&plan, entity, plan.next, end, selector)
			plan.Outcomes = append(plan.Outcomes, a.markTruncation(outcome))
		}
		plan.buckets = append(plan.buckets, plan.next)
		plan.next = end
	}
	return plan
}

// Commit applies a plan produced by Plan on the same assembler state.
func (a *Assembler) Commit(plan Plan) {
	for _, bucket := range plan.buckets {
		delete(a.buckets, bucket)
	}
	if !plan.Empty() {
		a.nextClose = plan.next
	}
	for member, checkpoint := range plan.counters {
		a.counters[member] = checkpoint
	}
}

// SetUntil ends this assembler's responsibility at an interval-aligned UTC
// boundary (a superseding revision's effective time, or a disable). It never
// extends a boundary that was already set earlier.
func (a *Assembler) SetUntil(until time.Time) {
	until = until.UTC()
	if !BucketStart(until, a.config.Interval).Equal(until) {
		return
	}
	if a.config.Until.IsZero() || until.Before(a.config.Until) {
		a.config.Until = until
	}
}

// SetEndsAt ends this assembler's responsibility at an exact instant that need
// not be interval-aligned; the window containing it closes with its real end.
// It never extends an earlier end.
func (a *Assembler) SetEndsAt(ends time.Time) {
	ends = ends.UTC()
	if !a.config.EndsAt.IsZero() && !ends.Before(a.config.EndsAt) {
		return
	}
	end := BucketStart(ends, a.config.Interval)
	if !end.Equal(ends) {
		end = end.Add(a.config.Interval)
	}
	if !a.config.Until.IsZero() && end.After(a.config.Until) {
		return
	}
	a.config.EndsAt, a.config.Until = ends, end
}

// markTruncation records the real bounds of a window that a cutover cut short.
// What the window would have written is still written, with the cut named: the
// accepted readings are not dropped, and a cut window is never passed off as a
// complete one. A window already decided not to be written (an unusable required
// member, no data) keeps that decision.
func (a *Assembler) markTruncation(outcome Outcome) Outcome {
	start, end, cut := a.windowBounds(outcome.BucketStart, outcome.BucketEnd)
	if cut == "" {
		return outcome
	}
	outcome.WindowStart, outcome.WindowEnd, outcome.Truncated = start, end, cut
	return outcome
}

// Done reports whether an Until was configured and every bucket before it has
// closed.
func (a *Assembler) Done() bool {
	return !a.config.Until.IsZero() && !a.nextClose.Before(a.config.Until)
}

// NextClose is the start of the next bucket that has not closed.
func (a *Assembler) NextClose() time.Time { return a.nextClose }

// OpenBuckets is the number of buckets currently holding samples.
func (a *Assembler) OpenBuckets() int { return len(a.buckets) }

func (a *Assembler) memberKeys(entity string) []string {
	for _, group := range a.entities {
		if group.key == entity {
			keys := make([]string, len(group.members))
			for i, member := range group.members {
				keys[i] = member.Key
			}
			return keys
		}
	}
	return nil
}

func (a *Assembler) finalize(plan *Plan, entity entityGroup, start, end time.Time, selector *Selector) Outcome {
	if a.config.Mode.perMember() {
		return a.finalizeProfile(plan, entity.members[0], start, end, selector)
	}
	if a.config.Mode.wide() {
		return a.finalizeWideProfile(plan, entity, start, end, selector)
	}
	outcome := Outcome{Kind: OutcomeRow, EntityKey: entity.key, BucketStart: start, BucketEnd: end}
	observed, usable, requiredGap := 0, 0, false
	for _, member := range entity.members {
		result := MemberResult{MemberKey: member.Key, Status: MemberMissing, Reason: ReasonNoSamples}
		if selector != nil {
			if sample, ok := selector.Selected(member.Key); ok {
				observed++
				result = evaluateMember(member, sample, end)
			}
		}
		if result.Usable() {
			usable++
		} else if member.Required {
			requiredGap = true
		}
		outcome.Members = append(outcome.Members, result)
	}
	partial := a.config.IncompletePolicy == IncompletePartial
	recordID, idErr := RecordID(a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, entity.key, start)
	outcome.RecordID = recordID
	switch {
	case idErr != nil:
		outcome.Kind, outcome.Reason = OutcomeSkipped, ReasonIdentityUnavailable
	case observed == 0:
		outcome.Kind, outcome.Reason = OutcomeNoData, ReasonNoSamples
	case partial && usable == 0 && !requiredGap:
		outcome.Kind, outcome.Reason = OutcomeSkipped, ReasonNoUsableMember
	case requiredGap || (!partial && usable < len(entity.members)):
		outcome.Kind, outcome.Reason = OutcomeSkipped, ReasonIncompleteRequired
	case usable < len(entity.members):
		outcome.Partial = true
	}
	if outcome.Kind == OutcomeRow {
		effectKey, err := EffectKey(a.config.DestinationScope, recordID)
		if err != nil {
			outcome.Kind, outcome.Reason, outcome.Partial = OutcomeSkipped, ReasonIdentityUnavailable, false
		} else {
			outcome.EffectKey = effectKey
		}
	}
	return outcome
}

// Late is the number of samples refused because their bucket had closed.
func (a *Assembler) Late() uint64 { return a.late }

// Plan is the set of buckets a Tick would close, computed without mutating
// the assembler. A caller that must persist closure first applies it with
// Commit only after the durable write succeeded.
type Plan struct {
	Outcomes []Outcome
	buckets  []time.Time
	next     time.Time
	// counters are the baselines Commit installs: only the members whose
	// baseline this plan advances.
	counters map[string]*aggregation.CounterCheckpoint
}

// NextClose is the checkpoint a committed plan advances to.
func (p Plan) NextClose() time.Time { return p.next }

// Empty reports whether the plan closes nothing.
func (p Plan) Empty() bool { return len(p.buckets) == 0 }
