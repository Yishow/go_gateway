package snapshot

import (
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
)

// OutcomeKind is the result of closing one entity's bucket.
type OutcomeKind string

const (
	// OutcomeRow carries a complete row to hand to delivery.
	OutcomeRow OutcomeKind = "row"
	// OutcomeSkipped means samples existed but the row was not usable.
	OutcomeSkipped OutcomeKind = "skipped"
	// OutcomeNoData means the whole bucket had no sample for the entity.
	OutcomeNoData OutcomeKind = "no_data"
)

// Reasons attached to skipped and no-data outcomes.
const (
	ReasonIncompleteRequired  = "incomplete-required-member"
	ReasonNoSamples           = "no-samples"
	ReasonNoUsableMember      = "no-usable-member"
	ReasonIdentityUnavailable = "identity-unavailable"
	// ReasonBoundaryTruncated is the local reason recorded beside a row whose
	// window a cutover cut short (see Outcome.Truncated).
	ReasonBoundaryTruncated = "boundary-truncated"
	// ReasonTypeMismatch marks a reading whose value cannot be carried as the
	// member's declared exact type.
	ReasonTypeMismatch = "type-mismatch"
	// ReasonMeanInvalid marks a mean window whose observations cannot produce a
	// typed arithmetic mean; the member carries no Reduction value then.
	ReasonMeanInvalid = "mean-invalid"
)

// MemberStatus explains why a member is or is not usable.
type MemberStatus string

const (
	MemberOK      MemberStatus = "ok"
	MemberMissing MemberStatus = "missing"
	MemberBad     MemberStatus = "bad"
	MemberStale   MemberStatus = "stale"
	MemberInvalid MemberStatus = "invalid"
)

// MemberResult is one member's evaluation. Sample is the selected sample and
// is nil only when nothing was observed for the member in the bucket.
type MemberResult struct {
	// Window is populated only for mixed results.
	Window    *MemberWindow `json:",omitempty"`
	MemberKey string
	Status    MemberStatus
	Reason    string
	Sample    *Sample
	// Counter is the typed cumulative record of the window. It is set only in
	// ModeCumulativeProfile for a member with an eligible reading; Sample is then
	// the raw reading it reports.
	Counter *aggregation.CounterWindowResult
	// Reduction is the typed arithmetic mean of the window's eligible
	// observations. It is set only for a SampleReductionMean member; Sample stays
	// the newest eligible reading of the same real window and is never replaced
	// by the mean.
	Reduction *measurement.MeanResult
	// Evidence is the value-free evidence of a mean window: aggregation
	// metadata, actual counts and observed bounds, coverage and stable reasons.
	// It is set only for a SampleReductionMean member.
	Evidence *MeanEvidence
}

// Outcome is emitted exactly once per entity and bucket.
type Outcome struct {
	// DueAt is the mixed dispatch boundary, not a shared data window.
	DueAt     time.Time `json:",omitzero"`
	Kind      OutcomeKind
	EntityKey string
	// RecordID identifies the row (or the skipped/no-data bucket) stably.
	RecordID string
	// EffectKey is set only for rows and is scoped by the frozen destination.
	EffectKey   string
	BucketStart time.Time
	BucketEnd   time.Time
	Reason      string
	Members     []MemberResult
	// Partial is true for a row that carries NULL for at least one member.
	Partial bool
	// WindowStart and WindowEnd are the real bounds of a window that a forward
	// cutover cut short; Truncated says which side ("start", "end" or "both").
	// BucketStart and BucketEnd stay the aligned bounds record identity uses.
	// All three are zero for a window that was not cut.
	WindowStart, WindowEnd time.Time
	Truncated              string
}

// Window returns the real bounds of the window the outcome covers. Mixed wide
// has no shared window: use DueAt and each MemberResult.Window instead.
func (o Outcome) Window() (start, end time.Time) {
	if o.Truncated == "" {
		return o.BucketStart, o.BucketEnd
	}
	return o.WindowStart, o.WindowEnd
}

// Usable reports whether the member may be written as a value. Every other
// status is written as SQL NULL with its reason, never as zero or a stale value.
func (m MemberResult) Usable() bool { return m.Status == MemberOK }
