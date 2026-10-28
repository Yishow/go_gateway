package workspace

import (
	"strings"

	"go-gateway/internal/datalink/schema"
)

// WriteGroupTransitionError is a safe, typed structural rejection. ReasonCode
// identifies the unsupported condition without exposing source or destination data.
type WriteGroupTransitionError struct {
	ReasonCode string
}

func (e *WriteGroupTransitionError) Error() string { return "recording transition structure rejected" }

// Unwrap retains the existing validation error classification.
func (e *WriteGroupTransitionError) Unwrap() error { return ErrWriteGroupValidation }

// ValidateWriteGroupTransitionStructure checks only the successor matrix.
// A nil error is NOT permission to prepare or activate a successor. Review must
// independently prove current target connector, capability, schema and ownership
// (including Timescale/native_monthly), and confirm frozen counter meaning using
// the existing proof engine. No predecessor destination proof may be reused.
// A nil error does not validate eligibility, lateness, completeness or storage-
// column policies; existing policy validation, exact review/digest and boundary
// checks still apply before any live operation.
//
// Engines are backend-resolved identities: sqlite, postgres, mysql, mariadb or
// sqlserver. MySQL driver identity alone does not identify the server family.
// tagTypes must be resolved type evidence; this function does not establish its
// provenance. Review must bind it to the immutable source and reviewed target.
// This function never probes, normalizes, saves, or mutates either group.
func ValidateWriteGroupTransitionStructure(source, target *WriteGroup, sourceEngine, targetEngine string, tagTypes map[string]schema.DataType) error {
	if source == nil || target == nil || strings.TrimSpace(source.WorkspaceID) == "" || source.WorkspaceID != target.WorkspaceID {
		return transitionRejection("invalid-transition")
	}
	if !transitionEngineSupported(sourceEngine) || !transitionEngineSupported(targetEngine) {
		return transitionRejection("unsupported-engine")
	}
	if sourceEngine != targetEngine {
		return transitionRejection("cross-engine")
	}
	if !transitionLayoutSupported(source, sourceEngine) || !transitionLayoutSupported(target, targetEngine) {
		return transitionRejection("unsupported-layout")
	}
	sourceMode, sourceModeOK := transitionStorageMode(source, sourceEngine)
	targetMode, targetModeOK := transitionStorageMode(target, targetEngine)
	if !sourceModeOK || !targetModeOK {
		return transitionRejection("unsupported-storage-mode")
	}
	if sourceMode != targetMode {
		return transitionRejection("storage-mode-change")
	}
	if source.RecordingProfile == nil || target.RecordingProfile == nil {
		return transitionRejection("unsupported-shape")
	}
	if source.RecordingProfile.Kind != target.RecordingProfile.Kind {
		return transitionRejection("profile-kind-change")
	}
	if !sameRecordingProfile(source.RecordingProfile, target.RecordingProfile) {
		return transitionRejection("recording-meaning-change")
	}
	sourceMembers, sourceOK := transitionMembers(source)
	targetMembers, targetOK := transitionMembers(target)
	if !sourceOK || !targetOK || len(sourceMembers) != len(targetMembers) {
		return transitionRejection("partial-source-set")
	}
	for identity, before := range sourceMembers {
		after, exists := targetMembers[identity]
		if !exists {
			return transitionRejection("partial-source-set")
		}
		before = transitionComparableMember(source, before, tagTypes)
		after = transitionComparableMember(target, after, tagTypes)
		if !sameTransitionMemberMeaning(before, after) {
			return transitionRejection("recording-meaning-change")
		}
	}
	if !transitionShapeSupported(source, tagTypes) || !transitionShapeSupported(target, tagTypes) {
		return transitionRejection("unsupported-shape")
	}
	return nil
}

func transitionRejection(reason string) error {
	return &WriteGroupTransitionError{ReasonCode: reason}
}

func transitionEngineSupported(engine string) bool {
	switch engine {
	case writeGroupConnectorKindSQLite, writeGroupConnectorKindPostgres, writeGroupConnectorKindMySQL, "mariadb", writeGroupConnectorKindSQLServer:
		return true
	}
	return false
}

func transitionLayoutSupported(group *WriteGroup, engine string) bool {
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return false
	}
	switch group.EffectiveLayoutVersion() {
	case WriteGroupLayoutManagedLongV1:
		return true
	case WriteGroupLayoutManagedWideProfileV1:
		return engine == writeGroupConnectorKindSQLite || engine == writeGroupConnectorKindPostgres
	}
	// Legacy wide/custom stays on the existing cutover contract and gates;
	// it is never admitted by this new managed-profile successor validator.
	return false
}

func transitionStorageMode(group *WriteGroup, engine string) (WriteGroupStorageMode, bool) {
	if group.isWideProfile() {
		// Wide profile stores no mode; its conceptual plain mode is not an
		// instruction to mutate the saved group by writing "plain" into it.
		return WriteGroupStorageModePlain, group.StorageMode == ""
	}
	mode := group.StorageMode
	if mode == "" {
		mode = WriteGroupStorageModePlain
	}
	switch mode {
	case WriteGroupStorageModePlain:
		return mode, true
	case WriteGroupStorageModeTimescale:
		return mode, engine == writeGroupConnectorKindPostgres
	case WriteGroupStorageModeNativeMonthly:
		return mode, engine == writeGroupConnectorKindMySQL || engine == "mariadb" || engine == writeGroupConnectorKindSQLServer
	}
	return mode, false
}

func transitionMembers(group *WriteGroup) (map[CutoverMemberRef]WriteGroupMember, bool) {
	members := make(map[CutoverMemberRef]WriteGroupMember, len(group.Members))
	for _, member := range group.Members {
		if strings.TrimSpace(member.DeviceID) == "" || strings.TrimSpace(member.PointID) == "" || strings.TrimSpace(member.TagID) == "" {
			return nil, false
		}
		identity := memberRefOf(member)
		if _, duplicate := members[identity]; duplicate {
			return nil, false
		}
		members[identity] = member
	}
	return members, len(members) > 0
}

func sameTransitionMemberMeaning(before, after WriteGroupMember) bool {
	return before.CircuitID == after.CircuitID && before.EntityKey == after.EntityKey &&
		before.SourceRevision == after.SourceRevision && before.MappingRevision == after.MappingRevision &&
		before.ValueRole == after.ValueRole && sameOptionalString(before.MeasurementID, after.MeasurementID) &&
		before.SourceValueType == after.SourceValueType && before.ResultValueType == after.ResultValueType &&
		sameExactRecordingSource(before.ExactRecordingSource, after.ExactRecordingSource)
}
