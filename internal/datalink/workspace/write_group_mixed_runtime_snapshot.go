package workspace

import "go-gateway/internal/datalink/schema"

// ValidateMixedRecordingSnapshot checks a frozen descriptor using the same
// member/schema and sealed proof contracts as confirmation. It does not grant
// activation, inspect a target, or consult mutable measurement definitions.
func ValidateMixedRecordingSnapshot(group *WriteGroup, types map[string]schema.DataType, proofs CounterBindingProofs) error {
	if !isMixedWriteGroup(group) || group.AppliedRevision == "" || group.AppliedRevision != group.Revision || group.RowPolicy.IntervalSeconds != 0 {
		return ErrWriteGroupValidation
	}
	if !mixedSchemaShapeSupported(group, types) {
		return ErrWriteGroupValidation
	}
	_, err := mixedSchemaCounterProofs(group, proofs)
	return err
}
