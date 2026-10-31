package recordingplan

// ValidateManagedGroupColumns checks only the existing canonical column/owner
// contract. Frozen runtime descriptors do not carry index catalog metadata;
// callers must not present this as fresh full-table or index verification.
func ValidateManagedGroupColumns(layout GroupSchemaLayout, state TargetTableInspection) error {
	normalized, err := ValidateGroupSchemaLayout(layout)
	if err != nil {
		return err
	}
	missing, err := inspectGroupExtension(state, normalized)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return ErrIncompatibleExistingTable
	}
	return nil
}
