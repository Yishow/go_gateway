package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrDeviceScopeRefresh means the device was deleted but Share needs recovery.
var ErrDeviceScopeRefresh = errors.New("device deleted but share scope refresh failed")

// DeviceRecordingReference identifies a recording that prevents device deletion.
type DeviceRecordingReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DeviceRecordingsReferencedError lists up to 20 recordings retaining a device.
type DeviceRecordingsReferencedError struct {
	Recordings []DeviceRecordingReference `json:"recordings"`
	Total      int                        `json:"total"`
}

func (e *DeviceRecordingsReferencedError) Error() string { return "DEVICE_REFERENCED_BY_RECORDINGS" }

// DeviceDeleteResult describes a committed device deletion.
type DeviceDeleteResult struct {
	WorkspaceID string
	DeviceID    string
	DeviceName  string
	Protocol    string
	Groups      []*WriteGroupSaveResult
}

// DeleteWorkspaceDevice checks references, purges eligible tombstones, and
// deletes the device and workspace membership in one configuration transaction.
func (s *WriteGroupService) DeleteWorkspaceDevice(ctx context.Context, deviceID string, deleteDevice func(context.Context, *sql.Tx, *sql.DB, string) error) (result *DeviceDeleteResult, err error) {
	if err := s.validate(); err != nil || deleteDevice == nil || strings.TrimSpace(deviceID) == "" {
		return nil, ErrWriteGroupServiceUnavailable
	}
	workspaceRepo, ok := s.workspaceSvc.repo.(*SQLRepository)
	if !ok || workspaceRepo.db != s.repo.db || s.repo.db == nil {
		return nil, ErrWriteGroupServiceUnavailable
	}
	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	var failCloser interface{ FailClosedMembership() }
	if s.workspaceSvc.shareScope != nil {
		var ok bool
		failCloser, ok = s.workspaceSvc.shareScope.(interface{ FailClosedMembership() })
		if !ok {
			return nil, ErrWriteGroupServiceUnavailable
		}
	}
	setup, err := workspaceRepo.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin device deletion: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			if rollbackErr := setup.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
				err = errors.Join(err, fmt.Errorf("rollback device deletion: %w", rollbackErr))
			}
		}
	}()
	record, err := setup.Get(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(record.OrderedDeviceIDs, deviceID) {
		return nil, ErrNotFound
	}
	tx := setup.SQLTx()
	active, tombstones, err := s.deviceRecordingReferencesInTx(ctx, tx, record.ID, deviceID)
	if err != nil {
		return nil, err
	}
	if len(active) > 0 {
		return nil, recordingReferencesError(active)
	}
	for _, ref := range tombstones {
		group, err := s.repo.GetInTx(ctx, tx, record.ID, ref.ID)
		if err != nil {
			return nil, err
		}
		ids, err := s.permanentDeleteGroupIDs(ctx, tx, group)
		if err != nil {
			return nil, err
		}
		for id := range ids {
			var otherDevices int
			err := tx.QueryRowContext(ctx, s.repo.query(`SELECT
				(SELECT COUNT(*) FROM write_group_members WHERE group_id=$1 AND device_id<>$2) +
				(SELECT COUNT(*) FROM write_group_basic_keys WHERE group_id=$3 AND device_id<>$4)`), id, deviceID, id, deviceID).Scan(&otherDevices)
			if err != nil {
				return nil, fmt.Errorf("check tombstone device ownership: %w", err)
			}
			if otherDevices > 0 {
				return nil, recordingReferencesError([]DeviceRecordingReference{ref})
			}
		}
	}
	result = &DeviceDeleteResult{WorkspaceID: record.ID, DeviceID: deviceID}
	if err := tx.QueryRowContext(ctx, s.repo.query(`SELECT name,protocol FROM devices WHERE id=$1`), deviceID).Scan(&result.DeviceName, &result.Protocol); err != nil {
		return nil, fmt.Errorf("read deleted device identity: %w", err)
	}
	for _, ref := range tombstones {
		group, getErr := s.repo.GetInTx(ctx, tx, record.ID, ref.ID)
		if errors.Is(getErr, ErrWriteGroupNotFound) {
			continue // A later generation already purged the complete lineage.
		}
		if getErr != nil {
			return nil, getErr
		}
		deleted, deleteErr := s.permanentDeleteInTx(ctx, tx, record, group)
		if deleteErr != nil {
			return nil, deleteErr
		}
		result.Groups = append(result.Groups, deleted)
	}
	if err := deleteDevice(ctx, tx, s.repo.db, deviceID); err != nil {
		return nil, err
	}
	record.OrderedDeviceIDs = slices.DeleteFunc(record.OrderedDeviceIDs, func(id string) bool { return id == deviceID })
	record.Status = workspaceStatusForDevices(record.OrderedDeviceIDs)
	if len(result.Groups) > 0 {
		record.DatabaseSetupRevision = s.workspaceSvc.newID()
	}
	record.UpdatedAt = s.workspaceSvc.now()
	if err := setup.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("save device deletion: %w", err)
	}
	if err := setup.Commit(); err != nil {
		return nil, fmt.Errorf("commit device deletion: %w", err)
	}
	committed = true
	for _, deleted := range result.Groups {
		s.recordPermanentDeleteAudit(ctx, deleted)
	}
	if s.workspaceSvc.shareScope != nil {
		if err := s.workspaceSvc.shareScope.InvalidateWorkspaceMembership(ctx, record.ID); err != nil {
			failCloser.FailClosedMembership()
			return result, fmt.Errorf("%w: %v", ErrDeviceScopeRefresh, err)
		}
	}
	return result, nil
}

func recordingReferencesError(refs []DeviceRecordingReference) *DeviceRecordingsReferencedError {
	total := len(refs)
	if len(refs) > 20 {
		refs = refs[:20]
	}
	return &DeviceRecordingsReferencedError{Recordings: refs, Total: total}
}

func (s *WriteGroupService) deviceRecordingReferencesInTx(ctx context.Context, tx *sql.Tx, workspaceID, deviceID string) (active, tombstones []DeviceRecordingReference, err error) {
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT w.id,w.name,w.status FROM write_groups w
		LEFT JOIN write_group_generations g ON g.group_id=w.id
		WHERE w.workspace_id=$1 AND (EXISTS (SELECT 1 FROM write_group_members m WHERE m.group_id=w.id AND m.device_id=$2)
		OR EXISTS (SELECT 1 FROM write_group_basic_keys b WHERE b.device_id=$3 AND (b.group_id=w.id OR b.lineage_id=g.lineage_id)))
		ORDER BY COALESCE(g.generation,0) DESC,w.id`), workspaceID, deviceID, deviceID)
	if err != nil {
		return nil, nil, fmt.Errorf("read device recording references: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var ref DeviceRecordingReference
		var status string
		if err := rows.Scan(&ref.ID, &ref.Name, &status); err != nil {
			return nil, nil, fmt.Errorf("scan device recording reference: %w", err)
		}
		if status == string(WriteGroupStatusDeleted) {
			tombstones = append(tombstones, ref)
		} else {
			active = append(active, ref)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate device recording references: %w", err)
	}
	return active, tombstones, nil
}
