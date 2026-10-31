package workspace

import (
	"fmt"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// writeGroupIncompletePolicyPartial is the one incomplete-row policy of the wide
// profile layout: a window writes SQL NULL for members without a usable value.
const writeGroupIncompletePolicyPartial = "partial"

// normalizeWideProfilePolicy gives the wide profile layout its only incomplete-row
// policy, partial, and refuses any other explicit choice. A member that must not
// be missing is marked required, which skips the whole row instead. Only this new
// layout gets the default; legacy wide and custom groups keep theirs.
func normalizeWideProfilePolicy(group *WriteGroup) error {
	switch strings.ToLower(strings.TrimSpace(group.RowPolicy.IncompletePolicy)) {
	case "", writeGroupIncompletePolicyPartial:
		group.RowPolicy.IncompletePolicy = writeGroupIncompletePolicyPartial
		return nil
	}
	return fmt.Errorf("%w: the %s layout writes partial rows; mark a member required to skip incomplete windows",
		ErrWriteGroupValidation, WriteGroupLayoutManagedWideProfileV1)
}

// isWideProfile reports whether the group uses the managed wide profile layout:
// one row per window with one column per member, typed by the frozen result.
func (g *WriteGroup) isWideProfile() bool {
	return g != nil && g.EffectiveLayoutVersion() == WriteGroupLayoutManagedWideProfileV1
}

// WideProfileShapeSupported reports whether a managed wide profile group's frozen
// recording facts are complete and consistent with the real tag types. Every
// member states its role and its source and result types; the mean applies only
// to numeric gauges; state and counter_raw members keep their latest typed value.
// A cumulative profile is served on its own terms (see cumulativeWideShapeSupported);
// plain storage only, a storage mode stays refused. The instantaneous shapes share
// one judgement with the long layout, so both layouts accept or refuse the same
// member shapes.
func WideProfileShapeSupported(group *WriteGroup, tagTypes map[string]schema.DataType) bool {
	if !group.isWideProfile() || group.RecordingProfile == nil || group.StorageMode != "" {
		return false
	}
	switch group.RecordingProfile.Kind {
	case WriteGroupProfileInstantaneous:
		return recordingShapeSupported(group, tagTypes, true)
	case WriteGroupProfileCumulative:
		return cumulativeWideShapeSupported(group, tagTypes)
	}
	return false
}

// cumulativeWideShapeSupported judges a cumulative wide profile: no reduction and
// no role (a gauge mean can never be a counter baseline), and every member
// names its raw numeric Tag type as both its source and its result, plus the
// saved measurement that is to confirm its counter semantics. Only the structure
// is judged here; that the measurement really confirms a counter is proved by the
// readiness service and the runtime counter provider, never by the id alone.
func cumulativeWideShapeSupported(group *WriteGroup, tagTypes map[string]schema.DataType) bool {
	if group.RecordingProfile.EffectiveSampleReduction() != WriteGroupSampleReductionLatest {
		return false
	}
	for i := range group.Members {
		member := &group.Members[i]
		tagType, ok := tagTypes[member.TagID]
		if !ok || member.ValueRole != "" || member.SourceValueType == "" || member.SourceValueType != tagType {
			return false
		}
		exact, ok := measurement.ExactTypeForTag(tagType)
		if !ok || member.ResultValueType != exact {
			return false
		}
		switch exact {
		case measurement.ExactInt64, measurement.ExactUint64, measurement.ExactFloat64:
		default:
			return false
		}
		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			return false
		}
	}
	return len(group.Members) > 0
}

// memberExactType is the exact value type a member's destination column holds.
// A wide profile column is typed by the frozen result type, never by the tag, so
// a mean stored as a decimal or a float does not change the source Tag type.
// Every other layout keeps the tag's own exact type.
func memberExactType(group *WriteGroup, member WriteGroupMember, tagTypes map[string]schema.DataType) (measurement.ExactType, bool) {
	if group.isWideProfile() || isMixedWriteGroup(group) {
		return member.ResultValueType, member.ResultValueType != ""
	}
	return measurement.ExactTypeForTag(tagTypes[member.TagID])
}
