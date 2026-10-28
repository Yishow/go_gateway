package snapshot

import (
	"maps"
	"slices"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// windowBounds returns the real bounds of a window that a cutover cut short and
// which side was cut ("start", "end", "both" or ""). BucketStart and BucketEnd
// stay the aligned bounds that record identity uses.
func (a *Assembler) windowBounds(start, end time.Time) (realStart, realEnd time.Time, cut string) {
	realStart, realEnd = start, end
	if !a.config.StartsAt.IsZero() && a.config.StartsAt.After(start) {
		realStart, cut = a.config.StartsAt, "start"
	}
	if !a.config.EndsAt.IsZero() && a.config.EndsAt.Before(end) {
		realEnd = a.config.EndsAt
		if cut == "" {
			cut = "end"
		} else {
			cut = "both"
		}
	}
	return realStart, realEnd, cut
}

// numericExactType reports whether an exact value can enter an arithmetic mean.
func numericExactType(valueType measurement.ExactType) bool {
	switch valueType {
	case measurement.ExactInt64, measurement.ExactUint64, measurement.ExactDecimal, measurement.ExactFloat64:
		return true
	}
	return false
}

// meanValueMatchesSource reports whether a reading agrees with the frozen Tag
// type. Without a frozen type nothing is inferred. A Sample value is already
// normalized to the wanted exact type by ExactFromGo before it enters a Sample,
// so the carried type must equal it; no range or conversion is invented here.
func meanValueMatchesSource(source schema.DataType, value measurement.ExactValue) bool {
	if source == "" {
		return true
	}
	wanted, ok := measurement.ExactTypeForTag(source)
	if !ok {
		return false
	}
	return value.Type() == wanted
}

// meanContract returns the typed result contract of a known frozen source type.
// A type without a numeric mean (bool, text, unknown) has no result contract.
func meanContract(source schema.DataType) (measurement.MeanBasis, measurement.ExactType) {
	wanted, ok := measurement.ExactTypeForTag(source)
	if !ok {
		return "", ""
	}
	switch wanted {
	case measurement.ExactInt64, measurement.ExactUint64:
		return measurement.MeanBasisExact, measurement.ExactDecimal
	case measurement.ExactFloat64:
		return measurement.MeanBasisFloat64, measurement.ExactFloat64
	}
	return "", ""
}

// meanObservation classifies one in-window observation for the mean evidence.
// It is deliberately separate from the legacy latest/counter status: only
// bounded public reasons are emitted, so an arbitrary raw quality reason never
// reaches operator-visible evidence.
func meanObservation(member Member, sample Sample, end time.Time) (string, bool) {
	evaluated := evaluateMember(member, sample, end)
	if evaluated.Status != MemberOK {
		if evaluated.Reason == ReasonTypeMismatch {
			return ReasonTypeMismatch, false
		}
		switch sample.Quality {
		case schema.QualityBad:
			return MeanReasonBad, false
		case schema.QualityInvalid:
			return MeanReasonInvalid, false
		case schema.QualityStale:
			return MeanReasonStale, false
		}
		switch evaluated.Status {
		case MemberStale:
			return MeanReasonStale, false
		case MemberInvalid:
			if evaluated.Reason == reasonMissingValue {
				return reasonMissingValue, false
			}
			return MeanReasonInvalid, false
		}
		return MeanReasonBad, false
	}
	if !meanValueMatchesSource(member.SourceValueType, sample.Value) {
		return ReasonTypeMismatch, false
	}
	if !numericExactType(sample.Value.Type()) {
		return MeanReasonNonNumeric, false
	}
	return "", true
}

// meanMissingStatus maps a stable reason to the member status of a window that
// carries no mean value.
func meanMissingStatus(reason string) MemberStatus {
	switch reason {
	case ReasonNoSamples:
		return MemberMissing
	case MeanReasonBad:
		return MemberBad
	case MeanReasonStale:
		return MemberStale
	default:
		return MemberInvalid
	}
}

// meanMemberResult derives a mean member's status, source sample and evidence
// only from the observations inside the real window. The mean never replaces
// Sample: a valid mean reports the newest eligible observation of that window,
// an unusable window reports its own evidence, and a window with no observation
// at all is missing. A reading outside the real window is never passed off as
// the window's value or evidence, and coverage is never complete without a
// whole-window proof.
func (a *Assembler) meanMemberResult(member Member, selector *Selector, start, end time.Time) MemberResult {
	basis, resultType := meanContract(member.SourceValueType)
	evidence := &MeanEvidence{
		Version:         measurement.MeanAggregationVersion,
		Method:          measurement.MeanAggregationMethod,
		Basis:           basis,
		SourceValueType: member.SourceValueType,
		ResultValueType: resultType,
		Coverage:        CoverageMissing,
		Reason:          ReasonNoSamples,
	}
	result := MemberResult{
		MemberKey: member.Key, Status: MemberMissing, Reason: ReasonNoSamples, Evidence: evidence,
	}
	if selector == nil {
		return result
	}
	windowStart, windowEnd, cut := a.windowBounds(start, end)
	accepted := selector.Accepted(member.Key)
	values := make([]measurement.ExactValue, 0, len(accepted))
	rejections := make(map[string]struct{})
	newestActual, newestEligible := -1, -1
	newestActualReason := ReasonNoSamples
	for index := range accepted {
		observed := accepted[index].ObservedAt.UTC()
		if observed.Before(windowStart) || !observed.Before(windowEnd) {
			continue
		}
		if newestActual < 0 {
			evidence.FirstObservedAt = observed
		}
		evidence.LastObservedAt = observed
		newestActual = index
		reason, eligible := meanObservation(member, accepted[index], windowEnd)
		if !eligible {
			evidence.RejectedCount++
			rejections[reason] = struct{}{}
			newestActualReason = reason
			continue
		}
		values = append(values, accepted[index].Value)
		newestEligible = index
	}
	evidence.ValidCount = int64(len(values))
	evidence.RejectionReasons = slices.Sorted(maps.Keys(rejections))
	if newestActual < 0 {
		return result
	}
	if newestEligible < 0 {
		// Observations exist inside the real window but none can supply a value:
		// report that evidence rather than the newest reading of another window.
		selected := accepted[newestActual]
		evidence.Reason = newestActualReason
		result.Status, result.Reason = meanMissingStatus(newestActualReason), newestActualReason
		result.Sample = &selected
		return result
	}
	selected := accepted[newestEligible]
	result.Sample = &selected
	reduction, err := measurement.Mean(values)
	if err != nil {
		result.Status, result.Reason = MemberInvalid, ReasonMeanInvalid
		evidence.Reason = ReasonMeanInvalid
		return result
	}
	result.Status, result.Reason, result.Reduction = MemberOK, "", &reduction
	evidence.Basis = reduction.Basis
	evidence.ResultValueType = reduction.Value.Type()
	switch {
	case cut != "":
		evidence.Coverage, evidence.Reason = CoveragePartial, ReasonBoundaryTruncated
	case evidence.RejectedCount > 0:
		evidence.Coverage = CoveragePartial
		evidence.Reason = evidence.RejectionReasons[0]
	default:
		evidence.Coverage, evidence.Reason = CoverageUnknown, ReasonCoverageUnverified
	}
	return result
}
