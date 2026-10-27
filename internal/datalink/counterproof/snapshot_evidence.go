package counterproof

import "encoding/json"

// SnapshotEvidence is a read-only projection of a validated original descriptor.
// Members includes non-counter members so a mixed accepted set can be checked
// without treating an unknown member as an ordinary gauge.
type SnapshotEvidence struct {
	Bindings    CounterBindingProofs
	Members     map[string]bool
	LegacyExact bool
}

// ResolveSnapshotEvidence uses the same strict descriptor and proof decoder as
// activation and acceptance; it does not introduce another interpretation.
func ResolveSnapshotEvidence(data []byte, scope SnapshotScope) (SnapshotEvidence, error) {
	proofs, err := ResolveSnapshot(data, scope)
	if err != nil {
		return SnapshotEvidence{}, err
	}
	var descriptor snapshotDescriptor
	if err := json.Unmarshal(data, &descriptor); err != nil {
		return SnapshotEvidence{}, err
	}
	evidence := SnapshotEvidence{Bindings: proofs, Members: map[string]bool{}, LegacyExact: descriptor.Bindings == nil}
	if descriptor.Group != nil {
		for _, member := range descriptor.Group.Members {
			identity := CounterBindingMemberIdentity{DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
			evidence.Members[identity.SampleMemberKey()] = true
		}
	}
	return evidence, nil
}
