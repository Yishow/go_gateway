package snapshot

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"go-gateway/internal/datalink/schema"
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
		byEntity[member.EntityKey] = append(byEntity[member.EntityKey], member)
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
		selector, ok := state.selectors[member.EntityKey]
		if !ok {
			selector = a.newSelector(start, member.EntityKey)
			state.selectors[member.EntityKey] = selector
		}
		selector.apply(sample)
	}
	return result
}

func (a *Assembler) newSelector(start time.Time, entity string) *Selector {
	return NewSelector(SelectorConfig{
		Start:         start,
		Interval:      a.config.Interval,
		MaxFutureSkew: a.config.MaxFutureSkew,
		Members:       a.memberKeys(entity),
	})
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
	if state, ok := a.buckets[start]; ok {
		if selector, ok := state.selectors[member.EntityKey]; ok {
			return selector.evaluate(sample, now)
		}
	}
	return a.newSelector(start, member.EntityKey).evaluate(sample, now)
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
	plan := Plan{next: a.nextClose}
	for len(plan.buckets) < a.config.MaxBucketsPerTick {
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
			plan.Outcomes = append(plan.Outcomes, a.finalize(entity, plan.next, end, selector))
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

func (a *Assembler) finalize(entity entityGroup, start, end time.Time, selector *Selector) Outcome {
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

func evaluateMember(member Member, sample Sample, end time.Time) MemberResult {
	selected := sample
	result := MemberResult{MemberKey: member.Key, Sample: &selected}
	switch {
	case sample.Quality != schema.QualityGood:
		result.Status = MemberBad
		result.Reason = sample.QualityReason
		if result.Reason == "" {
			result.Reason = fmt.Sprintf("quality-%s", sample.Quality)
		}
	case sample.Value.Type() == "":
		result.Status, result.Reason = MemberInvalid, "missing-value"
	case end.Sub(sample.ObservedAt) > member.MaxAge:
		result.Status, result.Reason = MemberStale, "max-age-exceeded"
	default:
		result.Status = MemberOK
	}
	return result
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
}

// NextClose is the checkpoint a committed plan advances to.
func (p Plan) NextClose() time.Time { return p.next }

// Empty reports whether the plan closes nothing.
func (p Plan) Empty() bool { return len(p.buckets) == 0 }
