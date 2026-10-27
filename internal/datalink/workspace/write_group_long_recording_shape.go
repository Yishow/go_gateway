package workspace

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// LongMemberStoredType returns the exact type a long member's window value must
// be stored as. A mean gauge stores the reduction result (decimal for exact
// sources, float64 for floats); every other member stores its tag's own exact
// type, which is the latest reading or the raw counter.
func LongMemberStoredType(reduction WriteGroupSampleReduction, role WriteGroupValueRole, source schema.DataType) (measurement.ExactType, bool) {
	exact, ok := measurement.ExactTypeForTag(source)
	if !ok {
		return "", false
	}
	if reduction == WriteGroupSampleReductionMean && role == WriteGroupValueRoleGauge {
		switch exact {
		case measurement.ExactInt64, measurement.ExactUint64:
			return measurement.ExactDecimal, true
		case measurement.ExactFloat64:
			return measurement.ExactFloat64, true
		default:
			return "", false
		}
	}
	return exact, true
}

// LongRecordingShapeSupported reports whether a long group's frozen recording
// metadata is complete and consistent with the real tag types. The mean
// reduction applies only to instantaneous gauges; state and counter_raw members
// keep their latest typed value; every new-shape member states its frozen
// source and result type; and a counter_raw keeps its confirmed binding. A
// legacy long group without metadata stays supported.
func LongRecordingShapeSupported(group *WriteGroup, tagTypes map[string]schema.DataType) bool {
	if group == nil || !group.isLongLayout() || group.RecordingProfile == nil {
		return false
	}
	return recordingShapeSupported(group, tagTypes, false)
}

// recordingShapeSupported judges the reduction and the frozen recording facts of
// the members of a group that has a recording profile. With every set, a member
// that carries no metadata is refused too; without it, a legacy member under the
// latest reduction stays supported.
func recordingShapeSupported(group *WriteGroup, tagTypes map[string]schema.DataType, every bool) bool {
	reduction := group.RecordingProfile.EffectiveSampleReduction()
	if reduction != WriteGroupSampleReductionLatest && reduction != WriteGroupSampleReductionMean {
		return false
	}
	if reduction == WriteGroupSampleReductionMean && group.RecordingProfile.Kind != WriteGroupProfileInstantaneous {
		return false
	}
	mean := reduction == WriteGroupSampleReductionMean
	for i := range group.Members {
		member := &group.Members[i]
		hasMetadata := memberMetadataUnsupported(member)
		if !every && !mean && !hasMetadata {
			continue
		}
		tagType, ok := tagTypes[member.TagID]
		if !ok {
			return false
		}
		if member.ExactRecordingSource != nil {
			if !group.isLongLayout() || group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged ||
				ValidateExactRecordingSource(*member, tagType) != nil {
				return false
			}
			cumulative := group.RecordingProfile.Kind == WriteGroupProfileCumulative && member.ValueRole == ""
			counterRaw := group.RecordingProfile.Kind == WriteGroupProfileInstantaneous && member.ValueRole == WriteGroupValueRoleCounterRaw
			if !cumulative && !counterRaw {
				return false
			}
		} else {
			switch member.ValueRole {
			case WriteGroupValueRoleGauge, WriteGroupValueRoleCounterRaw, WriteGroupValueRoleState:
			default:
				return false
			}
			if group.RecordingProfile.Kind != WriteGroupProfileInstantaneous {
				return false
			}
		}
		// A new shape proves its frozen facts explicitly; an omitted type is not
		// guessed from the tag at runtime.
		if member.SourceValueType == "" || member.SourceValueType != tagType {
			return false
		}
		stored, ok := LongMemberStoredTypeForMember(reduction, *member, tagType)
		if !ok || member.ResultValueType == "" || member.ResultValueType != stored {
			return false
		}
		if member.ValueRole == WriteGroupValueRoleCounterRaw && (member.MeasurementID == nil || *member.MeasurementID == "") {
			// Structure only: the id must be present. Its confirmed counter meaning
			// is proved by the readiness service and the runtime counter provider,
			// never by the id alone.
			return false
		}
	}
	return true
}
