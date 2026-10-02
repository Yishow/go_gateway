package snapshot

import "time"

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
	MemberKey string
	Status    MemberStatus
	Reason    string
	Sample    *Sample
}

// Outcome is emitted exactly once per entity and bucket.
type Outcome struct {
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
}

// Usable reports whether the member may be written as a value. Every other
// status is written as SQL NULL with its reason, never as zero or a stale value.
func (m MemberResult) Usable() bool { return m.Status == MemberOK }
