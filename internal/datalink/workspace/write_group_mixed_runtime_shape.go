package workspace

import "go-gateway/internal/datalink/schema"

// mixedRuntimeShapeSupported is the typed calculation/encoder support check,
// not readiness. Actual schema, ownership and coherent sealed counter proofs
// must still pass the existing readiness and Apply transaction authorities.
func mixedRuntimeShapeSupported(group *WriteGroup, types map[string]schema.DataType) bool {
	if !mixedSchemaShapeSupported(group, types) {
		return false
	}
	for _, member := range group.Members {
		if member.ExactRecordingSource != nil && (!group.isLongLayout() || member.Calculation != WriteGroupMemberCalculationDelta || ValidateExactRecordingSource(member, types[member.TagID]) != nil) {
			return false
		}
	}
	return true
}
