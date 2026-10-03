package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	basicManagedCanonicalRole = "basic"
	basicManagedRetryAttempts = 20
)

// ErrWriteGroupBasicManagedConflict reports a basic managed key whose saved
// group does not match the requested destination or source intent.
var ErrWriteGroupBasicManagedConflict = errors.New("basic managed write-group intent conflict")

// EnsureBasicManaged creates or reuses the canonical managed group for one
// device. The persisted create-once key and the group save share one local
// transaction.
func (s *WriteGroupService) EnsureBasicManaged(ctx context.Context, deviceID string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	s.basicManagedMu.Lock()
	defer s.basicManagedMu.Unlock()
	for attempt := 0; ; attempt++ {
		result, err := s.ensureBasicManaged(ctx, deviceID, mutation)
		if err == nil || !isSQLiteBusyError(err) || attempt == basicManagedRetryAttempts-1 {
			return result, err
		}
		timer := time.NewTimer(10 * time.Millisecond)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func (s *WriteGroupService) ensureBasicManaged(ctx context.Context, deviceID string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	deviceID = strings.TrimSpace(deviceID)
	if deviceID == "" {
		return nil, fmt.Errorf("ensure basic managed device: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, false, true)
	if err != nil {
		return nil, err
	}
	candidate, err := basicManagedCandidate(mutation, deviceID)
	if err != nil {
		return nil, err
	}

	if result, found, err := s.replayBasicManaged(ctx, mutation.WorkspaceID, deviceID, candidate); err != nil {
		return nil, err
	} else if found {
		return result, nil
	}

	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("ensure basic managed")
		}
		candidate := cloneWriteGroup(candidate)
		candidate.ID = ""
		candidate.Revision = ""
		candidate.AppliedRevision = ""
		candidate.Status = WriteGroupStatusDraft
		candidate.Migration = WriteGroupMigration{}
		candidate.Destination.SchemaRevision = ""
		candidate.Destination.SchemaDigest = ""
		candidate.WorkspaceID = record.ID
		candidate.Destination.ConnectorRevision = mutation.ExpectedConnectorRevision
		if err := s.repo.CreateInTx(ctx, tx, candidate); err != nil {
			return err
		}
		if err := s.repo.insertBasicManagedKeyInTx(ctx, tx, record.ID, deviceID, candidate.ID); err != nil {
			return err
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return err
		}
		saved = cloneWriteGroup(candidate)
		saved.BasicManagedDeviceID = deviceID
		return nil
	})
	if err == nil {
		return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
	}

	// A competing save may have won the persistent key after the optimistic
	// workspace check. Resolve the already-committed identity before returning
	// a stale-revision error so a lost response remains idempotent.
	if result, found, replayErr := s.replayBasicManaged(ctx, mutation.WorkspaceID, deviceID, candidate); replayErr != nil {
		return nil, replayErr
	} else if found {
		return result, nil
	}
	return nil, normalizeWriteGroupServiceError("ensure basic managed", err)
}

func isSQLiteBusyError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "database is locked") || strings.Contains(message, "sqlite_busy")
}

func basicManagedCandidate(mutation WriteGroupMutation, deviceID string) (*WriteGroup, error) {
	if mutation.Group == nil {
		return nil, fmt.Errorf("ensure basic managed payload: %w", ErrWriteGroupValidation)
	}
	candidate := cloneWriteGroup(mutation.Group)
	candidate.ID = ""
	candidate.Revision = ""
	candidate.AppliedRevision = ""
	candidate.Status = WriteGroupStatusDraft
	candidate.Migration = WriteGroupMigration{}
	candidate.WorkspaceID = mutation.WorkspaceID
	candidate.Name = strings.TrimSpace(candidate.Name)
	candidate.Destination.ConnectorID = strings.TrimSpace(candidate.Destination.ConnectorID)
	candidate.Destination.TableSchema = strings.TrimSpace(candidate.Destination.TableSchema)
	candidate.Destination.TableName = strings.TrimSpace(candidate.Destination.TableName)
	candidate.Destination.ConnectorRevision = mutation.ExpectedConnectorRevision
	candidate.Destination.StorageStrategy = WriteGroupStorageStrategyManaged
	candidate.Destination.SchemaRevision = ""
	candidate.Destination.SchemaDigest = ""
	candidate.RowPolicy.GroupKeyColumns = normalizeStringSet(candidate.RowPolicy.GroupKeyColumns)
	candidate.RowPolicy.UniqueKeyColumns = normalizeStringSet(candidate.RowPolicy.UniqueKeyColumns)
	if len(candidate.Members) == 0 || !slices.ContainsFunc(candidate.Members, func(member WriteGroupMember) bool {
		return strings.TrimSpace(member.DeviceID) == deviceID
	}) {
		return nil, fmt.Errorf("ensure basic managed member device: %w", ErrWriteGroupValidation)
	}
	for _, member := range candidate.Members {
		if strings.TrimSpace(member.DeviceID) != deviceID {
			return nil, fmt.Errorf("ensure basic managed requires one device: %w", ErrWriteGroupValidation)
		}
	}
	for i := range candidate.Members {
		candidate.Members[i].DeviceID = strings.TrimSpace(candidate.Members[i].DeviceID)
		candidate.Members[i].PointID = strings.TrimSpace(candidate.Members[i].PointID)
		candidate.Members[i].TagID = strings.TrimSpace(candidate.Members[i].TagID)
		candidate.Members[i].EntityKey = strings.TrimSpace(candidate.Members[i].EntityKey)
	}
	if err := prepareManagedWriteGroup(candidate); err != nil {
		return nil, err
	}
	return candidate, nil
}

func (s *WriteGroupService) replayBasicManaged(ctx context.Context, workspaceID, deviceID string, candidate *WriteGroup) (*WriteGroupSaveResult, bool, error) {
	group, err := s.repo.findBasicManaged(ctx, workspaceID, deviceID)
	if errors.Is(err, ErrWriteGroupNotFound) || errors.Is(err, ErrNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, normalizeWriteGroupServiceError("find basic managed", err)
	}
	if group.Status == WriteGroupStatusDeleted {
		return nil, false, ErrWriteGroupLifecycleBlocked
	}
	if !sameBasicManagedIntent(candidate, group) {
		return nil, false, ErrWriteGroupBasicManagedConflict
	}
	workspace, err := s.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return nil, false, normalizeWriteGroupServiceError("read basic managed workspace", err)
	}
	if workspace.ID != workspaceID {
		return nil, false, writeGroupNotFound("read basic managed workspace")
	}
	return &WriteGroupSaveResult{WorkspaceRevision: workspace.DatabaseSetupRevision, Group: cloneWriteGroup(group)}, true, nil
}

func sameBasicManagedIntent(candidate, existing *WriteGroup) bool {
	if candidate == nil || existing == nil || existing.Destination.StorageStrategy != WriteGroupStorageStrategyManaged {
		return false
	}
	if candidate.Destination.ConnectorID != existing.Destination.ConnectorID ||
		candidate.Destination.ConnectorRevision != existing.Destination.ConnectorRevision ||
		candidate.Destination.StorageStrategy != existing.Destination.StorageStrategy {
		return false
	}
	if candidate.Destination.TableName == "" {
		if existing.Destination.TableName != managedGroupTableName(existing.ID) {
			return false
		}
	} else if candidate.Destination.TableName != existing.Destination.TableName {
		return false
	}
	if candidate.Destination.Database != "" && candidate.Destination.Database != existing.Destination.Database {
		return false
	}
	if candidate.Destination.TableSchema != "" && candidate.Destination.TableSchema != existing.Destination.TableSchema {
		return false
	}
	comparableCandidate := cloneWriteGroup(candidate)
	if comparableCandidate.Destination.Database == "" {
		comparableCandidate.Destination.Database = existing.Destination.Database
	}
	if comparableCandidate.Destination.TableSchema == "" {
		comparableCandidate.Destination.TableSchema = existing.Destination.TableSchema
	}
	if comparableCandidate.Destination.TableName == "" {
		comparableCandidate.Destination.TableName = managedGroupTableName(existing.ID)
	}
	if !sameWriteGroupSchema(comparableCandidate, existing) {
		return false
	}
	if len(candidate.Members) != len(existing.Members) {
		return false
	}
	for _, member := range candidate.Members {
		if !slices.ContainsFunc(existing.Members, func(existingMember WriteGroupMember) bool {
			return member.DeviceID == existingMember.DeviceID &&
				member.PointID == existingMember.PointID &&
				member.TagID == existingMember.TagID &&
				member.EntityKey == existingMember.EntityKey &&
				member.TargetColumn == existingMember.TargetColumn &&
				member.Required == existingMember.Required &&
				sameOptionalString(member.MeasurementID, existingMember.MeasurementID) &&
				((member.MaxAgeSeconds == nil && existingMember.MaxAgeSeconds == nil) ||
					(member.MaxAgeSeconds != nil && existingMember.MaxAgeSeconds != nil && *member.MaxAgeSeconds == *existingMember.MaxAgeSeconds))
		}) {
			return false
		}
	}
	return true
}

func sameOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func (r *SQLWriteGroupRepository) findBasicManaged(ctx context.Context, workspaceID, deviceID string) (*WriteGroup, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("find basic managed group: %w", ErrWriteGroupNotFound)
	}
	var groupID string
	err := r.db.QueryRowContext(ctx, r.query(`
		SELECT group_id
		FROM write_group_basic_keys
		WHERE workspace_id = $1 AND device_id = $2 AND canonical_role = $3
	`), workspaceID, deviceID, basicManagedCanonicalRole).Scan(&groupID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("basic managed key not found: %w", ErrWriteGroupNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("read basic managed key: %w", err)
	}
	return r.Get(ctx, workspaceID, groupID)
}

func (r *SQLWriteGroupRepository) insertBasicManagedKeyInTx(ctx context.Context, tx *sql.Tx, workspaceID, deviceID, groupID string) error {
	if tx == nil || strings.TrimSpace(workspaceID) == "" || strings.TrimSpace(deviceID) == "" || strings.TrimSpace(groupID) == "" {
		return fmt.Errorf("insert basic managed key: %w", ErrWriteGroupValidation)
	}
	_, err := tx.ExecContext(ctx, r.query(`
		INSERT INTO write_group_basic_keys (workspace_id, device_id, canonical_role, group_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (workspace_id, device_id, canonical_role) DO NOTHING
	`), workspaceID, deviceID, basicManagedCanonicalRole, groupID)
	if err != nil {
		return fmt.Errorf("insert basic managed key: %w", err)
	}
	var storedGroupID string
	if err := tx.QueryRowContext(ctx, r.query(`
		SELECT group_id
		FROM write_group_basic_keys
		WHERE workspace_id = $1 AND device_id = $2 AND canonical_role = $3
	`), workspaceID, deviceID, basicManagedCanonicalRole).Scan(&storedGroupID); err != nil {
		return fmt.Errorf("confirm basic managed key: %w", err)
	}
	if storedGroupID != groupID {
		return fmt.Errorf("basic managed key already belongs to another group: %w", ErrWriteGroupBasicManagedConflict)
	}
	return nil
}
