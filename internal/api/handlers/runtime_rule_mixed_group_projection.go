package handlers

import (
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

func runtimeMixedGroupFacts(group *workspace.WriteGroup, facts *runtimeview.GroupFacts) error {
	if group.RecordID == "" || group.RecordingProfile.ResultShape != workspace.WriteGroupResultShapeMemberCalculationV1 || group.RowPolicy.IntervalSeconds != 0 {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if err := workspace.ValidateMemberCalculationIntent(group); err != nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	facts.MixedRecording = &runtimeview.MixedRecordingFacts{RecordID: group.RecordID, ResultShape: string(group.RecordingProfile.ResultShape)}
	facts.Reduction = "member_calculation"
	return nil
}
func runtimeMixedMemberFacts(member workspace.WriteGroupMember, out *runtimeview.Member) error {
	if member.Cadence == nil || member.SourceRevision == "" || member.MappingRevision == "" || member.ResultValueType == "" || out.SourceValueType == nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	if _, err := member.Cadence.Duration(); err != nil {
		return runtimeEvidenceError(runtimeview.ReasonUnavailable)
	}
	out.Calculation = &runtimeview.MemberCalculationFacts{Method: string(member.Calculation), Interval: member.Cadence.Interval, Unit: string(member.Cadence.Unit), Alignment: member.Cadence.Alignment, Required: member.Required, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, MeasurementID: member.MeasurementID, ResultType: string(member.ResultValueType)}
	return nil
}
