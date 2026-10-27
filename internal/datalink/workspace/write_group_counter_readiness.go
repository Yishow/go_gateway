package workspace

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// Stable readiness codes of a cumulative source whose saved measurement does not
// confirm counter semantics. They are separate on purpose: a gauge, a foreign
// definition, an unreadable one and an unconfirmed one each have a different
// repair, and none of them is ever read as ready.
const (
	counterMeasurementNotCounter  = "counter-measurement-not-counter"
	counterMeasurementForeign     = "counter-measurement-foreign"
	counterMeasurementUnavailable = "counter-measurement-unavailable"
	counterMeasurementUnconfirmed = "counter-measurement-unconfirmed"
	counterTypeUnsupported        = "counter-type-unsupported"
)

var counterMeasurementMessages = map[string]string{
	counterMeasurementNotCounter:  "a counter source names a saved measurement that is not a counter or signed_counter",
	counterMeasurementForeign:     "a counter source names a saved measurement that belongs to another workspace, device, point or tag",
	counterMeasurementUnavailable: "a counter source names a saved measurement that cannot be read",
	counterMeasurementUnconfirmed: "a counter source names a measurement whose unit, revisions, series epoch or policy are not confirmed",
	counterTypeUnsupported:        "counter recording requires a numeric source tag; string and bool tags cannot be counted",
}

// counterMeasurementIssues proves, for a supported cumulative group, that every
// source's saved measurement really is the confirmed counter the runtime counts
// from. A source without a measurement id is already reported by the counter
// gate; here a non-empty id must also name a definition of the same workspace,
// device, point and tag whose own evidence confirms counter semantics. The
// result is the same rule the runtime applies when it builds the group, so a
// group is never ready here and blocked there.
func (s *WriteGroupService) counterMeasurementIssues(ctx context.Context, group *WriteGroup, tagTypes map[string]schema.DataType) []ReadinessIssue {
	if (!longRuntimeSupported(group) && !group.isWideProfile()) || group.RecordingProfile == nil {
		return nil
	}
	cumulative := group.RecordingProfile.Kind == WriteGroupProfileCumulative
	if !cumulative && group.RecordingProfile.Kind != WriteGroupProfileInstantaneous {
		return nil
	}
	reasons := map[string]bool{}
	repo := measurement.NewSQLRepository(s.repo.db)
	for _, member := range group.Members {
		if !cumulative && member.ValueRole != WriteGroupValueRoleCounterRaw {
			continue // gauges and states never need a counter definition
		}
		if kind, ok := measurement.ExactTypeForTag(tagTypes[member.TagID]); ok {
			switch kind {
			case measurement.ExactInt64, measurement.ExactUint64, measurement.ExactFloat64:
			default:
				reasons[counterTypeUnsupported] = true
				continue
			}
		}
		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			continue // reported by the counter gate
		}
		def, err := repo.GetByID(ctx, *member.MeasurementID)
		switch {
		case err != nil || def == nil:
			reasons[counterMeasurementUnavailable] = true
		case def.WorkspaceID != group.WorkspaceID || def.DeviceID != member.DeviceID ||
			def.PointID != member.PointID || def.TagID == nil || *def.TagID != member.TagID:
			reasons[counterMeasurementForeign] = true
		default:
			switch measurement.ConfirmCounter(*def) {
			case measurement.CounterNotCounter:
				reasons[counterMeasurementNotCounter] = true
			case measurement.CounterUnconfirmed:
				reasons[counterMeasurementUnconfirmed] = true
			}
			if member.ExactRecordingSource != nil {
				proof := member.ExactRecordingSource
				if def.ExactRecordingScale != proof.Scale ||
					def.DefinitionRevision != proof.DefinitionRevision ||
					def.SourceBindingRevision != proof.SourceBindingRevision ||
					def.SeriesEpoch != proof.SeriesEpoch ||
					strings.TrimSpace(def.Unit) != proof.Unit {
					reasons[counterMeasurementUnconfirmed] = true
				}
			} else if strings.TrimSpace(def.ExactRecordingScale) != "" && group.isLongLayout() && group.Destination.StorageStrategy == WriteGroupStorageStrategyManaged {
				reasons[counterMeasurementUnconfirmed] = true
			}
		}
	}
	var issues []ReadinessIssue
	for _, code := range []string{counterMeasurementNotCounter, counterMeasurementForeign, counterMeasurementUnavailable, counterMeasurementUnconfirmed, counterTypeUnsupported} {
		if reasons[code] {
			issues = append(issues, writeGroupReadinessIssue(code, counterMeasurementMessages[code], group.ID))
		}
	}
	return issues
}
