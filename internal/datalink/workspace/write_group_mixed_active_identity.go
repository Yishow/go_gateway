package workspace

import "context"

// mixedActiveRecordIdentity is read-only evidence from the existing canonical
// key and lineage. It creates no source owner or second identity authority.
type mixedActiveRecordIdentity struct {
	groupID, recordID, lineageID string
}

func (r *SQLWriteGroupRepository) mixedRecordIdentityInTx(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) (*mixedActiveRecordIdentity, error) {
	if !isMixedWriteGroup(group) {
		return nil, nil
	}
	if group.ID == "" || group.WorkspaceID == "" || group.RecordingProfile.ResultShape != WriteGroupResultShapeMemberCalculationV1 ||
		(!group.isLongLayout() && !group.isWideProfile()) || group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return nil, ErrWriteGroupValidation
	}
	expected, err := mixedRecordProfileKey(group)
	if err != nil {
		return nil, err
	}
	meta, found, err := r.loadWriteGroupLineage(ctx, runner, group.WorkspaceID, group.ID)
	if err != nil {
		return nil, err
	}
	if !found || meta.CurrentGroupID != group.ID {
		return nil, lineageCorrupt("mixed source reference requires current canonical record")
	}
	key, found, err := r.loadProfileKey(ctx, runner, group.WorkspaceID, group.ID)
	if err != nil {
		return nil, err
	}
	if !found || key != expected {
		return nil, lineageCorrupt("mixed applied record identity disagrees with canonical key")
	}
	return &mixedActiveRecordIdentity{group.ID, group.RecordID, meta.LineageID}, nil
}

func independentMixedRecordIdentities(a, b *mixedActiveRecordIdentity) bool {
	return a != nil && b != nil && a.groupID != b.groupID && a.recordID != b.recordID && a.lineageID != b.lineageID
}
