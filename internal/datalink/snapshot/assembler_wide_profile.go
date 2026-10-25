package snapshot

import "time"

// latestWindowResult is a member's value under the latest rule: the newest
// eligible reading of the window, or the reason there is none. Nothing from
// another window and no zero is ever substituted.
func latestWindowResult(member Member, selector *Selector, end time.Time) MemberResult {
	result := MemberResult{MemberKey: member.Key, Status: MemberMissing, Reason: ReasonNoSamples}
	if selector == nil {
		return result
	}
	sample, ok := selector.Selected(member.Key)
	if !ok {
		return result
	}
	result = evaluateMember(member, sample, end)
	if result.Status == MemberBad && result.Reason == ReasonTypeMismatch {
		result.Status = MemberInvalid
	}
	return result
}

// finalizeWideProfile closes one fixed UTC window of a wide profile group into
// exactly one record. Each member keeps its own reduction: a numeric gauge
// averages the window's eligible observations, a state or counter_raw keeps its
// newest eligible reading. An optional member without a usable value is written
// as NULL with its reason, so a window where nothing was observed is still one
// row of NULLs rather than a silent gap. Only a required member that has no
// usable value skips the whole row, and says so. In the cumulative mode every
// member with an eligible reading is closed against its own counter baseline by
// the same engine as the per-member profile, so the raw reading and its typed
// difference are exactly what a long record would hold.
func (a *Assembler) finalizeWideProfile(plan *Plan, entity entityGroup, start, end time.Time, selector *Selector) Outcome {
	outcome := Outcome{Kind: OutcomeRow, EntityKey: entity.key, BucketStart: start, BucketEnd: end}
	usable, requiredGap := 0, false
	for _, member := range entity.members {
		var result MemberResult
		if member.SampleReduction == SampleReductionMean {
			result = a.meanMemberResult(member, selector, start, end)
		} else {
			result = latestWindowResult(member, selector, end)
			if a.config.Mode == ModeWideCumulative && result.Usable() {
				a.applyCounter(plan, member, start, end, selector, &result)
			}
		}
		if result.Usable() {
			usable++
		} else if member.Required {
			requiredGap = true
		}
		outcome.Members = append(outcome.Members, result)
	}
	recordID, idErr := RecordID(a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, entity.key, start)
	outcome.RecordID = recordID
	switch {
	case idErr != nil:
		outcome.Kind, outcome.Reason = OutcomeSkipped, ReasonIdentityUnavailable
	case requiredGap:
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
