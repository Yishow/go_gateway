package dbtarget

import (
	"encoding/json"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

// MixedSourceEvidence validates and projects the encoder's own provenance.
// Non-row outcomes may explain a required failure; they cannot emit SQL rows.
func MixedSourceEvidence(member GroupRowMember, outcome snapshot.Outcome, result snapshot.MemberResult, wide bool) ([]byte, error) {
	w, g, r := "", "", ""
	if member.Counter != nil {
		w, g, r = member.Counter.WorkspaceID, member.Counter.GroupID, member.Counter.ProfileRevision
	}
	if !validMixedIntent(member.mixedIntent(), w, g, r) {
		return nil, mixedRowError()
	}
	switch outcome.Kind {
	case snapshot.OutcomeRow:
	case snapshot.OutcomeSkipped, snapshot.OutcomeNoData:
		member.Required = false
	default:
		return nil, mixedRowError()
	}
	if err := validateMixedResult(member.mixedIntent(), outcome, result, wide); err != nil {
		return nil, err
	}
	return json.Marshal(mixedProvenance(member.mixedIntent(), result, outcome.DueAt))
}

// MixedSelectedValue is the same exact selector used by both SQL encoders.
// It never replaces an unavailable reduction with a raw reading or span delta.
func MixedSelectedValue(method snapshot.Calculation, result snapshot.MemberResult) (measurement.ExactValue, error) {
	return mixedMemberValue(method, result)
}
