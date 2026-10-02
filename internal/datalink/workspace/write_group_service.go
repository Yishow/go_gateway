package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrWriteGroupRevisionConflict reports a mutation made against a stale
	// workspace, group, or connector revision.
	ErrWriteGroupRevisionConflict = errors.New("write group revision conflict")
	// ErrWriteGroupLifecycleBlocked reports a lifecycle operation that cannot
	// be applied while the group still has an active applied revision.
	ErrWriteGroupLifecycleBlocked = errors.New("write group lifecycle is blocked")
	// ErrWriteGroupServiceUnavailable reports that the service cannot share a
	// local workspace transaction with the canonical repository.
	ErrWriteGroupServiceUnavailable = errors.New("write group service unavailable")
)

// WriteGroupMutation is the expected-revision envelope for one full group
// replacement. The service owns all server-generated and readonly fields.
type WriteGroupMutation struct {
	WorkspaceID               string      `json:"workspace_id"`
	ExpectedWorkspaceRevision string      `json:"expected_workspace_revision"`
	ExpectedGroupRevision     string      `json:"expected_group_revision"`
	ExpectedConnectorRevision string      `json:"expected_connector_revision"`
	Group                     *WriteGroup `json:"group"`
}

// WriteGroupSaveResult contains the coherent workspace and canonical group
// snapshot returned after a successful create, update, or delete.
type WriteGroupSaveResult struct {
	WorkspaceRevision string      `json:"workspace_revision"`
	Group             *WriteGroup `json:"group"`
}

// WriteGroupListResult contains one workspace revision and its group snapshot.
type WriteGroupListResult struct {
	WorkspaceID       string        `json:"workspace_id"`
	WorkspaceRevision string        `json:"workspace_revision"`
	Groups            []*WriteGroup `json:"groups"`
}

// WriteGroupService coordinates canonical WriteGroup persistence with the
// workspace compatibility projection in one local transaction.
type WriteGroupService struct {
	workspaceSvc              *Service
	repo                      *SQLWriteGroupRepository
	tableInspector            WriteGroupTableInspector
	backlogGuard              WriteGroupBacklogOwnershipGuard
	writerOwnershipActivation WriterOwnershipActivationBarrier
	now                       func() time.Time
}

// NewWriteGroupService creates the service used by the Studio V2 write-group
// API. The repository must use the same database as workspaceSvc.
func NewWriteGroupService(workspaceSvc *Service, repo *SQLWriteGroupRepository) *WriteGroupService {
	return &WriteGroupService{
		workspaceSvc: workspaceSvc,
		repo:         repo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

// Create persists a server-owned draft write group.
func (s *WriteGroupService) Create(ctx context.Context, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	mutation, err := prepareWriteGroupMutation(mutation, false, true)
	if err != nil {
		return nil, err
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("create write group")
		}
		candidate := cloneWriteGroup(mutation.Group)
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
		if err := projectWriteGroup(record, candidate); err != nil {
			return err
		}
		saved = cloneWriteGroup(candidate)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("create write group", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// Update replaces one draft write group subject to group, workspace, and
// connector compare-and-swap revisions.
func (s *WriteGroupService) Update(ctx context.Context, id string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("update write group: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, true, true)
	if err != nil {
		return nil, err
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("update write group")
		}
		existing, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return err
		}
		if existing.Revision != mutation.ExpectedGroupRevision {
			return writeGroupRevisionConflict("write group revision is stale")
		}
		if existing.Status == WriteGroupStatusDeleted {
			return ErrWriteGroupLifecycleBlocked
		}
		candidate := cloneWriteGroup(mutation.Group)
		candidate.ID = id
		candidate.WorkspaceID = record.ID
		candidate.Status = WriteGroupStatusDraft
		if existing.Status == WriteGroupStatusDisabled {
			candidate.Status = WriteGroupStatusDisabled
		}
		candidate.Destination.ConnectorRevision = mutation.ExpectedConnectorRevision
		if err := s.repo.updateInTx(ctx, tx, candidate, true); err != nil {
			return err
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return err
		}
		saved = cloneWriteGroup(candidate)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("update write group", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// Delete tombstones one un-applied write group while retaining its identity.
func (s *WriteGroupService) Delete(ctx context.Context, id string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("delete write group: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, true, false)
	if err != nil {
		return nil, err
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("delete write group")
		}
		existing, err := s.repo.GetInTx(ctx, tx, record.ID, id)
		if err != nil {
			return err
		}
		if existing.Revision != mutation.ExpectedGroupRevision {
			return writeGroupRevisionConflict("write group revision is stale")
		}
		if existing.Status == WriteGroupStatusDeleted {
			return ErrWriteGroupLifecycleBlocked
		}
		if existing.AppliedRevision != "" {
			if s.backlogGuard == nil {
				return ErrWriteGroupLifecycleBlocked
			}
			if err := s.backlogGuard.CheckWriteGroupBacklog(ctx, tx, existing); err != nil {
				return errors.Join(ErrWriteGroupLifecycleBlocked, err)
			}
		}
		if existing.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
			return writeGroupRevisionConflict("write group connector revision is stale")
		}
		if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, existing.Destination.ConnectorID, mutation.ExpectedConnectorRevision); err != nil {
			return err
		}
		candidate := cloneWriteGroup(existing)
		candidate.Status = WriteGroupStatusDeleted
		if err := s.repo.updateInTx(ctx, tx, candidate, false); err != nil {
			return err
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return err
		}
		saved = cloneWriteGroup(candidate)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("delete write group", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// List returns a coherent workspace and group snapshot.
func (s *WriteGroupService) List(ctx context.Context) (result *WriteGroupListResult, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("list write groups", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback list write-groups read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("list write groups")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("list write groups", err)
	}
	groups, err := s.repo.ListInTx(ctx, setupTx.SQLTx(), record.ID)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("list write groups", err)
	}
	result = &WriteGroupListResult{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		Groups:            cloneWriteGroups(groups),
	}
	return result, nil
}

// Get returns a coherent workspace and canonical group snapshot.
func (s *WriteGroupService) Get(ctx context.Context, id string) (result *WriteGroupSaveResult, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("get write group")
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("get write group", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback write-group read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("get write group")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("get write group", err)
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("get write group", err)
	}
	result = &WriteGroupSaveResult{
		WorkspaceRevision: record.DatabaseSetupRevision,
		Group:             cloneWriteGroup(group),
	}
	return result, nil
}

func (s *WriteGroupService) validate() error {
	if s == nil || s.workspaceSvc == nil || s.repo == nil || s.workspaceSvc.repo == nil {
		return ErrWriteGroupServiceUnavailable
	}
	return nil
}

func prepareWriteGroupMutation(mutation WriteGroupMutation, requireGroupRevision, requireGroup bool) (WriteGroupMutation, error) {
	mutation.WorkspaceID = strings.TrimSpace(mutation.WorkspaceID)
	mutation.ExpectedWorkspaceRevision = strings.TrimSpace(mutation.ExpectedWorkspaceRevision)
	mutation.ExpectedGroupRevision = strings.TrimSpace(mutation.ExpectedGroupRevision)
	mutation.ExpectedConnectorRevision = strings.TrimSpace(mutation.ExpectedConnectorRevision)
	if mutation.WorkspaceID == "" {
		return WriteGroupMutation{}, fmt.Errorf("write group workspace id: %w", ErrWriteGroupValidation)
	}
	if requireGroupRevision && mutation.ExpectedGroupRevision == "" {
		return WriteGroupMutation{}, fmt.Errorf("write group expected revision: %w", ErrWriteGroupValidation)
	}
	if mutation.ExpectedConnectorRevision == "" {
		return WriteGroupMutation{}, fmt.Errorf("write group connector revision: %w", ErrWriteGroupValidation)
	}
	if mutation.Group == nil {
		if requireGroup {
			return WriteGroupMutation{}, fmt.Errorf("write group payload: %w", ErrWriteGroupValidation)
		}
		return mutation, nil
	}
	mutation.Group = cloneWriteGroup(mutation.Group)
	mutation.Group.WorkspaceID = strings.TrimSpace(mutation.Group.WorkspaceID)
	if mutation.Group.WorkspaceID == "" {
		return WriteGroupMutation{}, fmt.Errorf("write group payload workspace id: %w", ErrWriteGroupValidation)
	}
	if mutation.Group.WorkspaceID != mutation.WorkspaceID {
		return WriteGroupMutation{}, writeGroupNotFound("write group payload")
	}
	if provided := strings.TrimSpace(mutation.Group.Destination.ConnectorRevision); provided != "" &&
		provided != mutation.ExpectedConnectorRevision {
		return WriteGroupMutation{}, writeGroupRevisionConflict("write group connector revision is stale")
	}
	mutation.Group.Destination.ConnectorRevision = mutation.ExpectedConnectorRevision
	return mutation, nil
}

func writeGroupNotFound(operation string) error {
	return fmt.Errorf("%s: %w: %w", operation, ErrWriteGroupNotFound, ErrNotFound)
}

func writeGroupRevisionConflict(reason string) error {
	return fmt.Errorf("%s: %w: %w", reason, ErrWriteGroupRevisionConflict, ErrSetupRevisionConflict)
}

func normalizeWriteGroupServiceError(operation string, err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrDatabaseSetupUnavailable):
		return fmt.Errorf("%s: %w", operation, errors.Join(ErrWriteGroupServiceUnavailable, err))
	case errors.Is(err, ErrSetupRevisionConflict):
		return fmt.Errorf("%s: %w: %w", operation, ErrWriteGroupRevisionConflict, err)
	case errors.Is(err, ErrWriteGroupNotFound), errors.Is(err, ErrNotFound),
		errors.Is(err, ErrWriteGroupValidation),
		errors.Is(err, ErrWriteGroupConnectorNotFound),
		errors.Is(err, ErrWriteGroupConnectorRevisionConflict),
		errors.Is(err, ErrWriteGroupSourceRevisionConflict),
		errors.Is(err, ErrWriteGroupUnsupportedConnectorKind),
		errors.Is(err, ErrWriteGroupRevisionConflict),
		errors.Is(err, ErrWriteGroupLifecycleBlocked):
		if errors.Is(err, ErrSetupRevisionConflict) && !errors.Is(err, ErrWriteGroupRevisionConflict) {
			return fmt.Errorf("%s: %w: %w", operation, ErrWriteGroupRevisionConflict, err)
		}
		return err
	default:
		return fmt.Errorf("%s: %w", operation, err)
	}
}

func cloneWriteGroups(groups []*WriteGroup) []*WriteGroup {
	if groups == nil {
		return nil
	}
	cloned := make([]*WriteGroup, len(groups))
	for i, group := range groups {
		cloned[i] = cloneWriteGroup(group)
	}
	return cloned
}

func validateLiveWriteGroupConnector(ctx context.Context, tx *sql.Tx, repo *SQLWriteGroupRepository, connectorID, expectedRevision string) error {
	connectorID = strings.TrimSpace(connectorID)
	expectedRevision = strings.TrimSpace(expectedRevision)
	if connectorID == "" || expectedRevision == "" {
		return fmt.Errorf("write group connector identity: %w", ErrWriteGroupValidation)
	}
	var revision string
	err := tx.QueryRowContext(ctx, repo.query(
		"SELECT identity_revision FROM database_connectors WHERE id = $1",
	), connectorID).Scan(&revision)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %w", ErrWriteGroupConnectorNotFound, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("read write-group connector identity: %w", err)
	}
	if revision != expectedRevision {
		return fmt.Errorf("%w: %w", ErrWriteGroupConnectorRevisionConflict, ErrSetupRevisionConflict)
	}
	return nil
}
