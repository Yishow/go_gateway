package workspace

import (
	"context"
	"database/sql"
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

// counterBindingEvidenceInTx proves, for a supported cumulative group, that every
// source's saved measurement really is the confirmed counter the runtime counts
// from. A source without a measurement id is already reported by the counter
// gate; here a non-empty id must also name a definition of the same workspace,
// device, point and tag whose own evidence confirms counter semantics. The
// result is the same rule the runtime applies when it builds the group, so a
// group is never ready here and blocked there.
func (s *WriteGroupService) counterBindingEvidenceInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, tagTypes map[string]schema.DataType, definitions map[string]*measurement.MeasurementDefinition) (CounterBindingProofs, []ReadinessIssue) {
	mixed := isMixedWriteGroup(group)
	if (!mixed && !longRuntimeSupported(group) && !group.isWideProfile()) || group.RecordingProfile == nil {
		return nil, nil
	}
	cumulative := group.RecordingProfile.Kind == WriteGroupProfileCumulative
	if !mixed && !cumulative && group.RecordingProfile.Kind != WriteGroupProfileInstantaneous {
		return nil, nil
	}
	var proofs CounterBindingProofs
	reasons := map[string]bool{}
	repo := measurement.NewSQLRepository(s.repo.db)
	for _, member := range group.Members {
		if mixed && member.Calculation != WriteGroupMemberCalculationDelta {
			continue
		}
		if !mixed && !cumulative && member.ValueRole != WriteGroupValueRoleCounterRaw {
			continue // gauges and states never need a counter definition
		}
		kind, knownType := measurement.ExactTypeForTag(tagTypes[member.TagID])
		supportedType := knownType && (kind == measurement.ExactInt64 || kind == measurement.ExactUint64 || kind == measurement.ExactFloat64)
		if member.MeasurementID == nil || strings.TrimSpace(*member.MeasurementID) == "" {
			if mixed {
				reasons[counterMeasurementUnconfirmed] = true
			}
			if !supportedType {
				reasons[counterTypeUnsupported] = true
			}
			continue // missing measurement is reported by the counter gate
		}
		def, cached := definitions[*member.MeasurementID]
		var err error
		if !cached {
			def, err = repo.GetByIDInTx(ctx, tx, *member.MeasurementID)
		}
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
				continue
			case measurement.CounterUnconfirmed:
				reasons[counterMeasurementUnconfirmed] = true
				continue
			}
			// Validation can withhold tag types after rejecting a definition.
			// Diagnose that original evidence first; an unproven source type
			// still cannot seal a proof even when the definition is confirmed.
			if !supportedType {
				reasons[counterTypeUnsupported] = true
				continue
			}
			var proof CounterBindingProof
			if member.ExactRecordingSource != nil {
				source := member.ExactRecordingSource
				if def.ExactRecordingScale != source.Scale ||
					def.DefinitionRevision != source.DefinitionRevision ||
					def.SourceBindingRevision != source.SourceBindingRevision ||
					def.SeriesEpoch != source.SeriesEpoch ||
					strings.TrimSpace(def.Unit) != source.Unit {
					reasons[counterMeasurementUnconfirmed] = true
					continue
				}
				identity := CounterBindingMemberIdentity{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.Revision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
				proof, err = CounterBindingProofFromExactSource(identity, member, tagTypes[member.TagID])
			} else if strings.TrimSpace(def.ExactRecordingScale) != "" && group.isLongLayout() && group.Destination.StorageStrategy == WriteGroupStorageStrategyManaged {
				reasons[counterMeasurementUnconfirmed] = true
				continue
			} else {
				proof, err = unscaledCounterBindingProof(group, member, tagTypes[member.TagID], def)
			}
			if err != nil {
				reasons[counterMeasurementUnconfirmed] = true
				continue
			}
			if proofs == nil {
				proofs = make(CounterBindingProofs)
			}
			proofs[proof.Member.Key()] = proof
		}
	}
	var issues []ReadinessIssue
	for _, code := range []string{counterMeasurementNotCounter, counterMeasurementForeign, counterMeasurementUnavailable, counterMeasurementUnconfirmed, counterTypeUnsupported} {
		if reasons[code] {
			issues = append(issues, writeGroupReadinessIssue(code, counterMeasurementMessages[code], group.ID))
		}
	}
	return proofs, issues
}

func unscaledCounterBindingProof(group *WriteGroup, member WriteGroupMember, tagType schema.DataType, def *measurement.MeasurementDefinition) (CounterBindingProof, error) {
	resultType, _ := measurement.ExactTypeForTag(tagType)
	facts := CounterBindingFacts{
		Member:        CounterBindingMemberIdentity{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.Revision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID},
		MeasurementID: def.ID, DefinitionRevision: def.DefinitionRevision, SourceBindingRevision: def.SourceBindingRevision,
		SeriesEpoch: def.SeriesEpoch, SemanticKind: def.SemanticKind, Unit: def.Unit,
		SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ScaleMode: CounterScaleUnscaled, ExactScale: "1", SourceValueType: tagType, ResultValueType: resultType,
	}
	if def.CounterPolicy != nil {
		facts.CounterPolicy = *def.CounterPolicy
	}
	return NewCounterBindingProof(facts)
}
