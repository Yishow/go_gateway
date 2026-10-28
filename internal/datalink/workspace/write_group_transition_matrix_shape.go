package workspace

import (
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// transitionComparableMember compares the effective raw type of cumulative
// long's legacy omission with wide's required explicit type. It only changes a
// value copy. It never infers an instantaneous role, alters a scaled proof, or
// claims that the caller's type evidence establishes frozen counter meaning.
func transitionComparableMember(group *WriteGroup, member WriteGroupMember, types map[string]schema.DataType) WriteGroupMember {
	if group.isLongLayout() && group.RecordingProfile.Kind == WriteGroupProfileCumulative &&
		member.ValueRole == "" && member.SourceValueType == "" && member.ResultValueType == "" && member.ExactRecordingSource == nil {
		member.SourceValueType = types[member.TagID]
		member.ResultValueType, _ = measurement.ExactTypeForTag(member.SourceValueType)
	}
	return member
}

func transitionShapeSupported(group *WriteGroup, types map[string]schema.DataType) bool {
	if group.RecordingProfile.validate() != nil {
		return false
	}
	for _, member := range group.Members {
		exact, known := measurement.ExactTypeForTag(types[member.TagID])
		if !known {
			return false
		}
		// The canonical wide contract has no exact-scaled representation.
		// Even a raw result type must not hide an attached scaled-source proof.
		if group.isWideProfile() && member.ExactRecordingSource != nil {
			return false
		}
		if group.RecordingProfile.Kind == WriteGroupProfileCumulative || member.ValueRole == WriteGroupValueRoleCounterRaw {
			if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
				return false
			}
			switch exact {
			case measurement.ExactInt64, measurement.ExactUint64, measurement.ExactFloat64:
			default:
				return false
			}
		}
	}
	if group.isWideProfile() {
		return WideProfileShapeSupported(group, types)
	}
	return LongRecordingShapeSupported(group, types)
}
