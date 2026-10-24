package snapshot

import (
	"fmt"
	"time"

	"go-gateway/internal/datalink/schema"
)

// Profile kinds as they appear in record identities.
const (
	profileKindInstantaneous = "instantaneous"
	profileKindCumulative    = "cumulative"
)

func (m Mode) profileKind() string {
	if m == ModeCumulativeProfile {
		return profileKindCumulative
	}
	return profileKindInstantaneous
}

// entityOf returns the partition a member's rows belong to. A wide group
// partitions by its entity key; a profile group makes every member its own
// partition, so each source has an independent record in every window.
func (c Config) entityOf(member Member) string {
	if c.Mode.perMember() {
		return member.Key
	}
	return member.EntityKey
}

// bucketsPerTick bounds one tick's work. A profile window yields one record per
// member, so the same record budget that a wide group spends on windows is
// divided across the members instead of growing with their number.
func (a *Assembler) bucketsPerTick() int {
	if !a.config.Mode.perMember() {
		return a.config.MaxBucketsPerTick
	}
	return max(1, a.config.MaxBucketsPerTick/len(a.entities))
}

// eligibility reports whether a sample may supply a window's value: good
// quality, a carried exact value, and fresh enough at the window end.
func (a *Assembler) eligibility(end time.Time) func(Sample) bool {
	return func(sample Sample) bool {
		member, ok := a.members[sample.MemberKey]
		return ok && evaluateMember(member, sample, end).Status == MemberOK
	}
}

// finalizeProfile closes one member's window. A record is always produced:
// with the newest eligible reading, or with no value and the reason why. A
// reading that cannot supply a value is never replaced by zero, by an older
// window's value, or by an invented time.
func (a *Assembler) finalizeProfile(plan *Plan, member Member, start, end time.Time, selector *Selector) Outcome {
	result := MemberResult{MemberKey: member.Key, Status: MemberMissing, Reason: ReasonNoSamples}
	if selector != nil {
		if sample, ok := selector.Selected(member.Key); ok {
			result = evaluateMember(member, sample, end)
			if result.Status == MemberBad && result.Reason == ReasonTypeMismatch {
				result.Status = MemberInvalid
			}
			if a.config.Mode == ModeCumulativeProfile && result.Usable() {
				a.applyCounter(plan, member, start, end, selector, &result)
			}
		}
	}
	outcome := Outcome{
		Kind: OutcomeRow, EntityKey: member.Key, BucketStart: start, BucketEnd: end,
		Members: []MemberResult{result},
	}
	recordID, err := ProfileRecordID(a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, member.Key, a.config.Mode.profileKind(), start)
	if err != nil {
		outcome.Kind, outcome.Reason = OutcomeSkipped, ReasonIdentityUnavailable
		return outcome
	}
	outcome.RecordID = recordID
	if outcome.EffectKey, err = EffectKey(a.config.DestinationScope, recordID); err != nil {
		outcome.Kind, outcome.Reason, outcome.EffectKey = OutcomeSkipped, ReasonIdentityUnavailable, ""
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
