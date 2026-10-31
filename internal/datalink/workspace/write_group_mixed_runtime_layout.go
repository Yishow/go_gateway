package workspace

import (
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

// MixedManagedSchemaLayout exposes the existing confirmed schema builder for
// frozen mixed storage validation. It grants neither readiness nor activation.
func MixedManagedSchemaLayout(group *WriteGroup, types map[string]schema.DataType, dialect string) (*recordingplan.GroupSchemaLayout, error) {
	if !isMixedWriteGroup(group) {
		return nil, ErrWriteGroupValidation
	}
	return managedSchemaLayout(group, types, dialect)
}
