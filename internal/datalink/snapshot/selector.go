package snapshot

import (
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// Sample is the typed observation a member contributes to a bucket.
type Sample struct {
	SampleID      string
	MemberKey     string
	ObservedAt    time.Time
	Quality       schema.QualityFlag
	QualityReason string
	Value         measurement.ExactValue
	// SourceRevision and MappingRevision identify the configuration that
	// produced the value; they take part in duplicate detection.
	SourceRevision  string
	MappingRevision string
}

// OfferOutcome describes what the selector did with one sample.
type OfferOutcome string

const (
	OfferSelected    OfferOutcome = "selected"
	OfferNotSelected OfferOutcome = "not_selected"
	OfferDuplicate   OfferOutcome = "duplicate"
	OfferConflict    OfferOutcome = "identity_conflict"
	OfferRejected    OfferOutcome = "rejected"
)

// Reasons attached to rejected or conflicting offers.
const (
	ReasonFutureSkew       = "future-skew"
	ReasonUnknownMember    = "unknown-member"
	ReasonOutsideBucket    = "outside-bucket"
	ReasonMissingIdentity  = "missing-identity"
	ReasonIdentityConflict = "identity-conflict"
	ReasonRevisionMismatch = "revision-mismatch"
)

// OfferResult is the deterministic result of Selector.Offer.
type OfferResult struct {
	Outcome OfferOutcome
	Reason  string
}

// SelectorConfig fixes the bucket and members one Selector serves.
type SelectorConfig struct {
	Start         time.Time
	Interval      time.Duration
	MaxFutureSkew time.Duration
	Members       []string
}

// Selector keeps, per member, the sample with the greatest observed_at in one
// bucket; the greatest sample_id breaks an observed_at tie. Arrival order never
// changes the result. The newest observation wins even when it is bad, so an
// older good reading cannot hide a failed read at the end of the bucket.
type Selector struct {
	start   time.Time
	end     time.Time
	skew    time.Duration
	members map[string]struct{}
	chosen  map[string]Sample
	seen    map[string]Sample
}

// BucketStart returns the UTC start of the half-open bucket containing t, or
// the zero time for a non-positive interval.
func BucketStart(t time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		return time.Time{}
	}
	ns := t.UnixNano()
	step := int64(interval)
	bucket := ns / step
	if ns%step < 0 {
		bucket--
	}
	return time.Unix(0, bucket*step).UTC()
}

// NewSelector builds a selector for one bucket.
func NewSelector(config SelectorConfig) *Selector {
	members := make(map[string]struct{}, len(config.Members))
	for _, member := range config.Members {
		members[member] = struct{}{}
	}
	start := config.Start.UTC()
	return &Selector{
		start:   start,
		end:     start.Add(config.Interval),
		skew:    config.MaxFutureSkew,
		members: members,
		chosen:  make(map[string]Sample, len(members)),
		seen:    make(map[string]Sample),
	}
}

// Offer applies a sample at gateway time now.
func (s *Selector) Offer(sample Sample, now time.Time) OfferResult {
	result := s.evaluate(sample, now)
	if result.Outcome == OfferSelected || result.Outcome == OfferNotSelected {
		s.apply(sample)
	}
	return result
}

// evaluate decides what Offer would do without changing the selector.
func (s *Selector) evaluate(sample Sample, now time.Time) OfferResult {
	if sample.SampleID == "" || sample.MemberKey == "" || sample.ObservedAt.IsZero() {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonMissingIdentity}
	}
	if previous, ok := s.seen[sample.SampleID]; ok {
		if samePayload(previous, sample) {
			return OfferResult{Outcome: OfferDuplicate}
		}
		return OfferResult{Outcome: OfferConflict, Reason: ReasonIdentityConflict}
	}
	if _, ok := s.members[sample.MemberKey]; !ok {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonUnknownMember}
	}
	observed := sample.ObservedAt.UTC()
	if observed.Before(s.start) || !observed.Before(s.end) {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonOutsideBucket}
	}
	if observed.After(now.UTC().Add(s.skew)) {
		return OfferResult{Outcome: OfferRejected, Reason: ReasonFutureSkew}
	}
	if current, ok := s.chosen[sample.MemberKey]; ok {
		sample.ObservedAt = observed
		if !supersedes(sample, current) {
			return OfferResult{Outcome: OfferNotSelected}
		}
	}
	return OfferResult{Outcome: OfferSelected}
}

// apply stores a sample that evaluate accepted.
func (s *Selector) apply(sample Sample) {
	sample.ObservedAt = sample.ObservedAt.UTC()
	s.seen[sample.SampleID] = sample
	if current, ok := s.chosen[sample.MemberKey]; !ok || supersedes(sample, current) {
		s.chosen[sample.MemberKey] = sample
	}
}

// Selected returns the chosen sample for a member.
func (s *Selector) Selected(member string) (Sample, bool) {
	sample, ok := s.chosen[member]
	return sample, ok
}

func supersedes(candidate, current Sample) bool {
	if !candidate.ObservedAt.Equal(current.ObservedAt) {
		return candidate.ObservedAt.After(current.ObservedAt)
	}
	return candidate.SampleID > current.SampleID
}

func samePayload(a, b Sample) bool {
	if a.MemberKey != b.MemberKey || !a.ObservedAt.Equal(b.ObservedAt) ||
		a.Quality != b.Quality || a.QualityReason != b.QualityReason ||
		a.SourceRevision != b.SourceRevision || a.MappingRevision != b.MappingRevision {
		return false
	}
	if a.Value.Type() == "" || b.Value.Type() == "" {
		return a.Value.Type() == b.Value.Type()
	}
	return a.Value.Equal(b.Value)
}
