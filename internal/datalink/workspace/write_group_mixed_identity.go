package workspace

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const circuitProfileRoleRecord = "record"

func mixedRecordProfileKey(group *WriteGroup) (circuitProfileKey, error) {
	id, err := uuid.Parse(group.RecordID)
	if err != nil || id == uuid.Nil || id.String() != group.RecordID {
		return circuitProfileKey{}, fmt.Errorf("%w: mixed record_id must be a canonical nonzero UUID", ErrWriteGroupValidation)
	}
	return circuitProfileKey{WorkspaceID: group.WorkspaceID, Kind: WriteGroupProfileMixed, Role: circuitProfileRoleRecord, SelectionScope: group.RecordID}, nil
}

func (r *SQLWriteGroupRepository) loadMixedRecordID(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	group.RecordID = ""
	if !isMixedWriteGroup(group) {
		return nil
	}
	key, found, err := r.loadProfileKey(ctx, runner, group.WorkspaceID, group.ID)
	if err != nil {
		return err
	}
	// Storage-only 1.2 snapshots have no canonical key; they remain inspectable
	// but cannot enter normal draft mutation, which requires complete identity.
	if !found {
		return nil
	}
	if key.Kind != WriteGroupProfileMixed || key.Role != circuitProfileRoleRecord || key.CircuitID != "" {
		return fmt.Errorf("%w: invalid stored mixed record identity", ErrWriteGroupValidation)
	}
	if group.RowPolicy.IntervalSeconds != 0 {
		return fmt.Errorf("%w: stored mixed record carries legacy group interval", ErrWriteGroupValidation)
	}
	group.RecordID = key.SelectionScope
	if _, err = mixedRecordProfileKey(group); err != nil {
		return err
	}
	return validateMixedMemberCadences(group)
}

func sameMixedMemberRevisions(candidate, stored *WriteGroup) bool {
	if candidate.RecordID != stored.RecordID {
		return false
	}
	for _, member := range candidate.Members {
		for _, saved := range stored.Members {
			if member.DeviceID != saved.DeviceID || member.PointID != saved.PointID || member.TagID != saved.TagID {
				continue
			}
			if member.SourceRevision != "" && member.SourceRevision != saved.SourceRevision {
				return false
			}
			if member.MappingRevision != "" && member.MappingRevision != saved.MappingRevision {
				return false
			}
		}
	}
	return true
}
