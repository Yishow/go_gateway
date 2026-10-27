package counterproof

import (
	"encoding/json"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// SnapshotScope selects the original durable applied revision, not a draft.
type SnapshotScope struct{ WorkspaceID, GroupID, AppliedRevision string }

type snapshotMember struct {
	DeviceID        string                `json:"device_id"`
	PointID         string                `json:"point_id"`
	TagID           string                `json:"tag_id"`
	MeasurementID   *string               `json:"measurement_id"`
	SourceRevision  string                `json:"source_revision"`
	MappingRevision string                `json:"mapping_revision"`
	SourceValueType schema.DataType       `json:"source_value_type"`
	ResultValueType measurement.ExactType `json:"result_value_type"`
	ValueRole       string                `json:"value_role"`
	ExactSource     *ExactRecordingSource `json:"exact_recording_source"`
}
type snapshotGroup struct {
	ID              string `json:"id"`
	WorkspaceID     string `json:"workspace_id"`
	Revision        string `json:"revision"`
	AppliedRevision string `json:"applied_revision"`
	LayoutVersion   string `json:"layout_version"`
	Destination     struct {
		StorageStrategy string `json:"storage_strategy"`
	} `json:"destination"`
	Profile *struct {
		Kind string `json:"kind"`
	} `json:"recording_profile"`
	Members []snapshotMember `json:"members"`
}
type snapshotDescriptor struct {
	WorkspaceID     string         `json:"workspace_id"`
	GroupID         string         `json:"group_id"`
	AppliedRevision string         `json:"applied_revision"`
	Group           *snapshotGroup `json:"group"`
	Layout          *struct {
		TagTypes map[string]schema.DataType `json:"tag_types"`
	} `json:"runtime_layout"`
	Bindings CounterBindingProofs `json:"counter_bindings"`
}

// ResolveSnapshot resolves only immutable stored evidence. Legacy exact facts
// are adapted in memory; no original bytes, journal digest or audit is rewritten.
// Missing unscaled facts are never reconstructed from an editable definition.
func ResolveSnapshot(data []byte, scope SnapshotScope) (CounterBindingProofs, error) {
	fail := func(reason error) (CounterBindingProofs, error) {
		return nil, &CounterBindingProofError{Reason: reason}
	}
	if !UniqueJSONObject(data, false) || !snapshotFieldNames(data) {
		return fail(ErrCounterBindingProofConflict)
	}
	var snap snapshotDescriptor
	if err := json.Unmarshal(data, &snap); err != nil {
		if _, ok := err.(*CounterBindingProofError); ok {
			return nil, err
		}
		return fail(ErrCounterBindingProofConflict)
	}
	if snap.WorkspaceID != scope.WorkspaceID || snap.GroupID != scope.GroupID || snap.AppliedRevision != scope.AppliedRevision {
		return fail(ErrCounterBindingProofConflict)
	}
	group := snap.Group
	if group == nil {
		if len(snap.Bindings) > 0 {
			return fail(ErrCounterBindingProofConflict)
		}
		return nil, nil
	}
	if group.ID != scope.GroupID || group.WorkspaceID != scope.WorkspaceID || group.Revision != scope.AppliedRevision || group.AppliedRevision != scope.AppliedRevision {
		return fail(ErrCounterBindingProofConflict)
	}
	resolved := CounterBindingProofs{}
	seen := map[string]bool{}
	for _, member := range group.Members {
		identity := CounterBindingMemberIdentity{WorkspaceID: scope.WorkspaceID, GroupID: scope.GroupID, AppliedRevision: scope.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
		key := identity.Key()
		if seen[key] {
			return fail(ErrCounterBindingProofConflict)
		}
		seen[key] = true
		counter := group.Profile != nil && (group.Profile.Kind == "cumulative" || group.Profile.Kind == "instantaneous" && member.ValueRole == "counter_raw")
		proof, exists := snap.Bindings[key]
		if !counter {
			if exists || member.ExactSource != nil {
				return fail(ErrCounterBindingProofConflict)
			}
			continue
		}
		if snap.Layout == nil {
			return fail(ErrCounterBindingProofMissing)
		}
		tagType := snap.Layout.TagTypes[member.TagID]
		if !exists {
			if snap.Bindings != nil || member.ExactSource == nil {
				return fail(ErrCounterBindingProofMissing)
			}
			var err error
			proof, err = legacyExactProof(group, member, identity, tagType)
			if err != nil {
				return nil, err
			}
		}
		if proof.Member != identity || member.MeasurementID == nil || *member.MeasurementID != proof.MeasurementID || proof.SourceRevision != member.SourceRevision || proof.MappingRevision != member.MappingRevision || proof.SourceValueType != tagType {
			return fail(ErrCounterBindingProofConflict)
		}
		if member.SourceValueType != "" && member.SourceValueType != proof.SourceValueType || member.ResultValueType != "" && member.ResultValueType != proof.ResultValueType {
			return fail(ErrCounterBindingProofConflict)
		}
		if proof.ScaleMode == CounterScaleExact || member.ExactSource != nil {
			exact, err := legacyExactProof(group, member, identity, tagType)
			if err != nil {
				return nil, err
			}
			if exact.Digest != proof.Digest {
				return fail(ErrCounterBindingProofConflict)
			}
		}
		resolved[key] = proof
	}
	for key := range snap.Bindings {
		if _, ok := resolved[key]; !ok {
			return fail(ErrCounterBindingProofConflict)
		}
	}
	return resolved, nil
}

func legacyExactProof(group *snapshotGroup, member snapshotMember, identity CounterBindingMemberIdentity, tagType schema.DataType) (CounterBindingProof, error) {
	fail := func() (CounterBindingProof, error) {
		return CounterBindingProof{}, &CounterBindingProofError{Reason: ErrCounterBindingProofConflict}
	}
	s := member.ExactSource
	if s == nil || group.LayoutVersion != "managed_long_v1" || group.Destination.StorageStrategy != "managed" || group.Profile == nil ||
		(group.Profile.Kind == "cumulative" && member.ValueRole != "") || (group.Profile.Kind != "cumulative" && (group.Profile.Kind != "instantaneous" || member.ValueRole != "counter_raw")) ||
		member.MeasurementID == nil || *member.MeasurementID != s.MeasurementID || s.OriginalDataType != tagType || member.SourceValueType != tagType || member.ResultValueType != measurement.ExactDecimal {
		return fail()
	}
	facts := CounterBindingFacts{Member: identity, MeasurementID: s.MeasurementID, DefinitionRevision: s.DefinitionRevision, SourceBindingRevision: s.SourceBindingRevision, SeriesEpoch: s.SeriesEpoch, SemanticKind: s.SemanticKind, Unit: s.Unit, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ScaleMode: CounterScaleExact, ExactScale: s.Scale, SourceValueType: tagType, ResultValueType: measurement.ExactDecimal}
	if s.CounterPolicy != nil {
		facts.CounterPolicy = *s.CounterPolicy
	}
	return NewCounterBindingProof(facts)
}

// Structural fields are canonical ASCII; map keys (including Unicode source
// identities and tag IDs) are data and deliberately not subjected to this rule.
func snapshotFieldNames(data []byte) bool {
	var object func(json.RawMessage, string) bool
	object = func(raw json.RawMessage, kind string) bool {
		if string(raw) == "null" || len(raw) == 0 {
			return true
		}
		var fields map[string]json.RawMessage
		if json.Unmarshal(raw, &fields) != nil {
			return false
		}
		for name, value := range fields {
			if name == "" || strings.IndexFunc(name, func(r rune) bool { return r != '_' && (r < 'a' || r > 'z') && (r < '0' || r > '9') }) >= 0 {
				return false
			}
			child := ""
			switch kind + "/" + name {
			case "snapshot/group":
				child = "group"
			case "snapshot/runtime_layout":
				child = "layout"
			case "group/recording_profile", "group/destination", "member/counter_policy", "exact/counter_policy":
				child = "fields"
			case "member/exact_recording_source":
				child = "exact"
			case "group/members":
				var members []json.RawMessage
				if json.Unmarshal(value, &members) != nil {
					return false
				}
				for _, m := range members {
					if !object(m, "member") {
						return false
					}
				}
			}
			if child != "" && !object(value, child) {
				return false
			}
		}
		return true
	}
	return object(data, "snapshot")
}
