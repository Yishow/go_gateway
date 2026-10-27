package runtime

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"
)

// frozenExactCounterSource reads only the immutable source proof carried by
// this applied long member. Editable definitions never replace accepted facts.
func frozenExactCounterSource(group *workspace.WriteGroup, member workspace.WriteGroupMember, kind measurement.ExactType) (CounterSource, bool) {
	if group == nil || group.EffectiveLayoutVersion() != workspace.WriteGroupLayoutManagedLongV1 ||
		group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged || group.RecordingProfile == nil ||
		member.ExactRecordingSource == nil || member.SourceRevision == "" || member.MappingRevision == "" {
		return CounterSource{}, false
	}
	if group.RecordingProfile.Kind == workspace.WriteGroupProfileCumulative {
		if member.ValueRole != "" {
			return CounterSource{}, false
		}
	} else if group.RecordingProfile.Kind != workspace.WriteGroupProfileInstantaneous || member.ValueRole != workspace.WriteGroupValueRoleCounterRaw {
		return CounterSource{}, false
	}
	proof := member.ExactRecordingSource
	if workspace.ValidateExactRecordingSource(member, proof.OriginalDataType) != nil {
		return CounterSource{}, false
	}
	if expected, ok := measurement.ExactTypeForTag(proof.OriginalDataType); !ok || expected != kind {
		return CounterSource{}, false
	}
	definition := proof.Definition()
	// Validation above confirmed the scale. The legacy mapper supplies only the
	// unchanged counter semantics; the boundary applies the frozen scale once.
	definition.ExactRecordingScale = ""
	return CounterSourceFromDefinition(definition)
}
