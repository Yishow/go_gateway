package workspace

import "go-gateway/internal/datalink/measurement"

// mixedSchemaCounterProofs selects existing sealed counter facts only for delta
// members. The proof's result type is the counter input (scaled or native), not
// the selected delta output; member metadata separately binds that output type.
func mixedSchemaCounterProofs(group *WriteGroup, proofs CounterBindingProofs) (CounterBindingProofs, error) {
	selected := CounterBindingProofs{}
	for _, member := range group.Members {
		if member.Calculation != WriteGroupMemberCalculationDelta {
			continue
		}
		identity := CounterBindingMemberIdentity{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.Revision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
		proof, ok := proofs[identity.Key()]
		if !ok {
			return nil, ErrCounterBindingProofMissing
		}
		if proof.Validate() != nil || proof.Member != identity || member.MeasurementID == nil || proof.MeasurementID != *member.MeasurementID || proof.SourceRevision != member.SourceRevision || proof.MappingRevision != member.MappingRevision || proof.SourceValueType != member.SourceValueType {
			return nil, ErrCounterBindingProofConflict
		}
		if member.ExactRecordingSource != nil {
			expected, err := CounterBindingProofFromExactSource(identity, member, member.SourceValueType)
			if err != nil || expected.Digest != proof.Digest {
				return nil, ErrCounterBindingProofConflict
			}
		} else {
			source, ok := measurement.ExactTypeForTag(member.SourceValueType)
			if !ok || proof.ScaleMode != CounterScaleUnscaled || proof.ExactScale != "1" || proof.ResultValueType != source {
				return nil, ErrCounterBindingProofConflict
			}
		}
		selected[identity.Key()] = proof
	}
	if len(selected) != len(proofs) {
		return nil, ErrCounterBindingProofConflict
	}
	if len(selected) == 0 {
		return nil, nil
	}
	return selected, nil
}
