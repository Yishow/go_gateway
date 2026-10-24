package workspace

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// EnsureCircuitProfile creates or reuses the canonical managed group of one
// circuit profile. The create-once key and the group save share one local
// transaction, so a reload, a resend or a rename resolves to the saved group
// instead of adding a second writer. Display names are not part of intent.
func (s *WriteGroupService) EnsureCircuitProfile(ctx context.Context, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	// Create-once ensures share one lock; their keys are independent.
	s.basicManagedMu.Lock()
	defer s.basicManagedMu.Unlock()
	return retryWhileSQLiteBusy(ctx, func() (*WriteGroupSaveResult, error) {
		return s.ensureCircuitProfile(ctx, mutation)
	})
}

func (s *WriteGroupService) ensureCircuitProfile(ctx context.Context, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	mutation, err := prepareWriteGroupMutation(mutation, false, true)
	if err != nil {
		return nil, err
	}
	candidate, key, err := circuitProfileCandidate(mutation)
	if err != nil {
		return nil, err
	}
	if result, found, err := s.replayCircuitProfile(ctx, key, candidate); err != nil {
		return nil, err
	} else if found {
		return result, nil
	}

	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("ensure circuit profile")
		}
		created := cloneWriteGroup(candidate)
		created.WorkspaceID = record.ID
		if err := s.repo.CreateInTx(ctx, tx, created); err != nil {
			return err
		}
		if err := projectWriteGroup(record, created); err != nil {
			return err
		}
		saved = cloneWriteGroup(created)
		return nil
	})
	if err == nil {
		return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
	}

	// A competing save may have won the key after the optimistic workspace
	// check. Resolve the committed identity before reporting a stale revision
	// so a lost response stays idempotent.
	if result, found, replayErr := s.replayCircuitProfile(ctx, key, candidate); replayErr != nil {
		return nil, replayErr
	} else if found {
		return result, nil
	}
	return nil, normalizeWriteGroupServiceError("ensure circuit profile", err)
}

// circuitProfileCandidate resets every server-owned field of the caller's
// payload and derives the canonical key it would occupy.
func circuitProfileCandidate(mutation WriteGroupMutation) (*WriteGroup, circuitProfileKey, error) {
	candidate := cloneWriteGroup(mutation.Group)
	candidate.ID = ""
	candidate.Revision = ""
	candidate.AppliedRevision = ""
	candidate.BasicManagedDeviceID = ""
	candidate.Status = WriteGroupStatusDraft
	candidate.Migration = WriteGroupMigration{}
	candidate.WorkspaceID = mutation.WorkspaceID
	candidate.Name = strings.TrimSpace(candidate.Name)
	candidate.Destination.ConnectorID = strings.TrimSpace(candidate.Destination.ConnectorID)
	candidate.Destination.ConnectorRevision = mutation.ExpectedConnectorRevision
	candidate.Destination.TableSchema = strings.TrimSpace(candidate.Destination.TableSchema)
	candidate.Destination.TableName = strings.TrimSpace(candidate.Destination.TableName)
	candidate.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	candidate.Destination.SchemaRevision = ""
	candidate.Destination.SchemaDigest = ""
	candidate.RowPolicy.GroupKeyColumns = normalizeStringSet(candidate.RowPolicy.GroupKeyColumns)
	candidate.RowPolicy.UniqueKeyColumns = normalizeStringSet(candidate.RowPolicy.UniqueKeyColumns)
	for i := range candidate.Members {
		member := &candidate.Members[i]
		member.PointID = strings.TrimSpace(member.PointID)
		member.TagID = strings.TrimSpace(member.TagID)
		member.EntityKey = strings.TrimSpace(member.EntityKey)
	}
	if err := normalizeWriteGroupProfile(candidate); err != nil {
		return nil, circuitProfileKey{}, err
	}
	key, err := circuitProfileKeyOf(candidate)
	if err != nil {
		return nil, circuitProfileKey{}, err
	}
	if err := prepareManagedWriteGroup(candidate); err != nil {
		return nil, circuitProfileKey{}, err
	}
	return candidate, key, nil
}

func (s *WriteGroupService) replayCircuitProfile(ctx context.Context, key circuitProfileKey, candidate *WriteGroup) (*WriteGroupSaveResult, bool, error) {
	group, err := s.repo.findProfileGroup(ctx, key)
	if errors.Is(err, ErrWriteGroupNotFound) || errors.Is(err, ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, normalizeWriteGroupServiceError("find circuit profile group", err)
	}
	if group.Status == WriteGroupStatusDeleted {
		return nil, false, ErrWriteGroupLifecycleBlocked
	}
	if !sameCircuitProfileIntent(candidate, group) {
		return nil, false, ErrWriteGroupCircuitProfileConflict
	}
	result, err := s.currentGroupResult(ctx, group, "read circuit profile workspace")
	if err != nil {
		return nil, false, err
	}
	return result, true, nil
}

// sameCircuitProfileIntent compares destination, sources, layout and profile.
// Group and circuit display names are deliberately ignored: renaming never
// changes identity.
func sameCircuitProfileIntent(candidate, existing *WriteGroup) bool {
	if candidate.isLongLayout() && existing.isLongLayout() && candidate.Destination.TableName == "" {
		// A long group's table is server-derived from its scope, not from its ID.
		candidate = cloneWriteGroup(candidate)
		candidate.Destination.TableName = existing.Destination.TableName
	}
	return sameBasicManagedIntent(candidate, existing) &&
		candidate.EffectiveLayoutVersion() == existing.EffectiveLayoutVersion() &&
		candidate.StorageMode == existing.StorageMode &&
		sameRecordingProfile(candidate.RecordingProfile, existing.RecordingProfile)
}
