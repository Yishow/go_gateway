package snapshot

import (
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
)

func (a *Assembler) mixedMemberResult(p *Plan, m Member, start, end, nextDue time.Time) MemberResult {
	selector := a.mixedActualSelector(start, m)
	actualStart, actualEnd, cut := a.windowBounds(start, end)
	var result MemberResult
	if m.Calculation == CalculationMean {
		result = a.meanMemberResult(m, selector, start, end)
	} else {
		result = latestWindowResult(m, selector, actualEnd)
		if result.Sample != nil && !meanValueMatchesSource(m.SourceValueType, result.Sample.Value) {
			result.Status, result.Reason = MemberInvalid, ReasonTypeMismatch
		}
		if m.Calculation == CalculationDelta && result.Usable() {
			a.applyCounter(p, m, start, end, selector, &result)
			// Raw evidence is not a selected delta. Keep the engine's baseline/span
			// evidence and progress, but only a reliable single-window delta is usable.
			if result.Usable() && result.Counter != nil && (result.Counter.Status != aggregation.CounterOK || result.Counter.Delta == nil) {
				result.Status = MemberMissing
				result.Reason = result.Counter.Reason
				if result.Reason == "" {
					result.Reason = string(result.Counter.Status)
				}
			}
		}
	}
	result.Window = &MemberWindow{Calculation: m.Calculation, Interval: m.Interval, AlignedStart: start, AlignedEnd: end, ActualStart: actualStart, ActualEnd: actualEnd, Truncated: cut, NextDue: nextDue}
	return result
}

func (a *Assembler) mixedIdentity(out *Outcome, member string, at time.Time) {
	var id string
	var err error
	if a.config.Mode == ModeMixedLong {
		id, err = ProfileRecordID(a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, member, "mixed", at)
	} else {
		id, err = measurement.HashStringTuple("write-group-mixed-dispatch-v1", a.config.WorkspaceID, a.config.GroupID, a.config.GroupRevision, at.UTC().Format(time.RFC3339Nano))
	}
	if err != nil {
		out.Kind, out.Reason = OutcomeSkipped, ReasonIdentityUnavailable
		return
	}
	out.RecordID = id
	if out.Kind == OutcomeRow {
		out.EffectKey, err = EffectKey(a.config.DestinationScope, id)
		if err != nil {
			out.Kind, out.Reason = OutcomeSkipped, ReasonIdentityUnavailable
		}
	}
}

// enforceMixedCounterType preserves the engine's full-sequence uncertainty and
// reliable-tail decision. Filtering wrong-type observations before reduction
// would hide a discontinuity and could fabricate normal usage. A wrong-type
// tail is never durable state; a reliable correct-type tail may still advance,
// but the affected window cannot expose a selected delta.
func (a *Assembler) enforceMixedCounterType(member Member, readings []aggregation.CounterReading, counter *aggregation.CounterWindowResult, next, prev *aggregation.CounterCheckpoint, result *MemberResult) *aggregation.CounterCheckpoint {
	mismatch := false
	for _, reading := range readings {
		if !meanValueMatchesSource(member.SourceValueType, reading.Value) {
			mismatch = true
			break
		}
	}
	if !mismatch {
		return next
	}
	result.Status, result.Reason = MemberInvalid, ReasonTypeMismatch
	counter.Delta = nil
	switch counter.Status {
	case aggregation.CounterOK, aggregation.CounterBaseline, aggregation.CounterGap:
		counter.Status, counter.Reason = aggregation.CounterUncertain, ReasonTypeMismatch
	}
	if next != nil && next != prev && !meanValueMatchesSource(member.SourceValueType, next.Reading.Value) {
		return prev
	}
	return next
}
