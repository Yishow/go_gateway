package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// ErrWriteGroupLegacyWriteConflict reports a legacy target write whose output
// scope is already owned by a canonical write group.
var ErrWriteGroupLegacyWriteConflict = errors.New("legacy write is owned by a canonical write group")

// PreflightLegacyTargetWrite checks whether a legacy target mapping may be
// written without opening the external destination or changing local state.
func (s *WriteGroupService) PreflightLegacyTargetWrite(ctx context.Context, mappingID, tagID, connectorID string) (err error) {
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	mappingID = strings.TrimSpace(mappingID)
	tagID = strings.TrimSpace(tagID)
	connectorID = strings.TrimSpace(connectorID)
	if mappingID == "" && (tagID == "" || connectorID == "") {
		return fmt.Errorf("legacy target write identity: %w", ErrWriteGroupValidation)
	}

	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	tx, err := s.beginLegacyWriteTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback legacy target preflight: %w", rollbackErr))
		}
	}()
	record, err := readLegacyWriteWorkspace(ctx, tx)
	if err != nil {
		return err
	}
	return s.checkLegacyTargetWriteWithRecord(ctx, tx, record, mappingID, tagID, connectorID)
}

// RunLegacyTargetWrite runs one legacy mapping persistence callback in the
// workspace-owned local transaction after rechecking canonical ownership.
func (s *WriteGroupService) RunLegacyTargetWrite(
	ctx context.Context,
	mappingID string,
	candidate *schema.DatabaseTargetMapping,
	persist func(context.Context, *sql.Tx) error,
) (err error) {
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	if persist == nil {
		return fmt.Errorf("legacy target write persistence: %w", ErrWriteGroupValidation)
	}
	mappingID = strings.TrimSpace(mappingID)
	tagID, connectorID := legacyTargetIdentity(candidate)
	if mappingID == "" && (tagID == "" || connectorID == "") {
		return fmt.Errorf("legacy target write identity: %w", ErrWriteGroupValidation)
	}

	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	tx, err := s.beginLegacyWriteTx(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback legacy target write: %w", rollbackErr))
		}
	}()
	record, err := readLegacyWriteWorkspace(ctx, tx)
	if err != nil {
		return err
	}
	if err := s.checkLegacyTargetWriteWithRecord(ctx, tx, record, mappingID, tagID, connectorID); err != nil {
		return err
	}
	if err := persist(ctx, tx); err != nil {
		return fmt.Errorf("persist legacy target write: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit legacy target write: %w", err)
	}
	committed = true
	return nil
}

// CheckLegacyTargetWriteInTx rechecks legacy mapping ownership in tx before a
// caller writes a target mapping.
func (s *WriteGroupService) CheckLegacyTargetWriteInTx(
	ctx context.Context,
	tx *sql.Tx,
	mappingID, tagID, connectorID string,
) error {
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	if tx == nil {
		return fmt.Errorf("legacy target write transaction: %w", ErrWriteGroupValidation)
	}
	mappingID = strings.TrimSpace(mappingID)
	tagID = strings.TrimSpace(tagID)
	connectorID = strings.TrimSpace(connectorID)
	if mappingID == "" && (tagID == "" || connectorID == "") {
		return fmt.Errorf("legacy target write identity: %w", ErrWriteGroupValidation)
	}
	record, err := readLegacyWriteWorkspace(ctx, tx)
	if err != nil {
		return err
	}
	return s.checkLegacyTargetWriteWithRecord(ctx, tx, record, mappingID, tagID, connectorID)
}

// PreflightLegacyRowGroupReplacement checks whether a row-group replacement
// may proceed without changing local state.
func (s *WriteGroupService) PreflightLegacyRowGroupReplacement(
	ctx context.Context,
	connectorID string,
	groups []DatabaseRowGroup,
) (err error) {
	if groups == nil {
		return nil
	}
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	connectorID = strings.TrimSpace(connectorID)
	if connectorID == "" {
		if len(groups) == 0 {
			// The first save of a new destination replaces nothing: there is no
			// connector yet, so no canonical group can own any of its scope.
			return nil
		}
		return fmt.Errorf("legacy row-group replacement connector: %w", ErrWriteGroupValidation)
	}

	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	tx, err := s.beginLegacyWriteTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback legacy row-group preflight: %w", rollbackErr))
		}
	}()
	record, err := readLegacyWriteWorkspace(ctx, tx)
	if err != nil {
		return err
	}
	return s.checkLegacyRowGroupReplacementWithRecord(ctx, tx, record, connectorID, groups)
}

// CheckLegacyRowGroupReplacementInTx rechecks row-group ownership in tx before
// a caller persists a connector or workspace replacement.
func (s *WriteGroupService) CheckLegacyRowGroupReplacementInTx(
	ctx context.Context,
	tx *sql.Tx,
	current *Record,
	connectorID string,
	groups []DatabaseRowGroup,
) error {
	if groups == nil {
		return nil
	}
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	if tx == nil || current == nil {
		return fmt.Errorf("legacy row-group replacement transaction: %w", ErrWriteGroupValidation)
	}
	connectorID = strings.TrimSpace(connectorID)
	if connectorID == "" {
		if len(groups) == 0 {
			// The first save of a new destination replaces nothing: there is no
			// connector yet, so no canonical group can own any of its scope.
			return nil
		}
		return fmt.Errorf("legacy row-group replacement connector: %w", ErrWriteGroupValidation)
	}
	return s.checkLegacyRowGroupReplacementWithRecord(ctx, tx, current, connectorID, groups)
}

func (s *WriteGroupService) validateLegacyWriteService() error {
	if err := s.validate(); err != nil {
		return err
	}
	if s.repo.db == nil {
		return ErrWriteGroupServiceUnavailable
	}
	return nil
}

func (s *WriteGroupService) beginLegacyWriteTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := s.repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin legacy write transaction: %w", err)
	}
	return tx, nil
}

func readLegacyWriteWorkspace(ctx context.Context, tx *sql.Tx) (*Record, error) {
	record, err := readWorkspaceRecord(ctx, tx)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read legacy write workspace: %w", err)
	}
	return record, nil
}

func legacyTargetIdentity(candidate *schema.DatabaseTargetMapping) (tagID, connectorID string) {
	if candidate == nil {
		return "", ""
	}
	return strings.TrimSpace(candidate.TagID), strings.TrimSpace(candidate.ConnectorID)
}

// PreflightConnectorDelete refuses to delete a connector that any canonical
// group (including a tombstone) still uses as its destination, because that
// would orphan the group's frozen delivery identity and make every later
// migrated read fail.
func (s *WriteGroupService) PreflightConnectorDelete(ctx context.Context, connectorID string) (err error) {
	if err := s.validateLegacyWriteService(); err != nil {
		return err
	}
	connectorID = strings.TrimSpace(connectorID)
	if connectorID == "" {
		return fmt.Errorf("connector delete identity: %w", ErrWriteGroupValidation)
	}
	s.workspaceSvc.mu.Lock()
	defer s.workspaceSvc.mu.Unlock()
	tx, err := s.beginLegacyWriteTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback connector delete preflight: %w", rollbackErr))
		}
	}()
	record, err := readLegacyWriteWorkspace(ctx, tx)
	if err != nil {
		return err
	}
	groups, err := s.readLegacyOwnershipGroups(ctx, tx, record.ID)
	if err != nil {
		return err
	}
	for _, group := range groups {
		if group != nil && group.Destination.ConnectorID == connectorID {
			return legacyWriteConflict()
		}
	}
	return nil
}
