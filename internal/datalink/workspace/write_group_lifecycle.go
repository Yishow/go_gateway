package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	// ErrWriteGroupAppliedRevisionUnavailable means that a group has no
	// immutable snapshot that can be resolved at the requested time.
	ErrWriteGroupAppliedRevisionUnavailable = errors.New("write group applied revision unavailable")
	// ErrWriteGroupApplyNotReady means that a draft has not passed the current
	// configuration and read-only schema evidence gate.
	ErrWriteGroupApplyNotReady = errors.New("write group apply is not ready")
	// ErrWriteGroupIntervalOverflow means an interval cannot be represented by
	// time.Duration or by a common old/new bucket boundary.
	ErrWriteGroupIntervalOverflow = errors.New("write group interval overflows time duration")
)

// WriteGroupAppliedSnapshot is an immutable applied group snapshot selected
// for one UTC bucket boundary. Group is always returned as a defensive copy.
type WriteGroupAppliedSnapshot struct {
	WorkspaceID     string                   `json:"workspace_id"`
	GroupID         string                   `json:"group_id"`
	GroupRevision   string                   `json:"group_revision"`
	AppliedRevision string                   `json:"applied_revision"`
	EffectiveAt     time.Time                `json:"effective_at"`
	Group           *WriteGroup              `json:"group"`
	RuntimeLayout   *WriteGroupRuntimeLayout `json:"runtime_layout,omitempty"`
}

// WriteGroupIntakeEligibility describes whether new samples may enter the
// currently resolved applied snapshot. It does not describe delivery status.
type WriteGroupIntakeEligibility struct {
	WorkspaceID     string           `json:"workspace_id"`
	GroupID         string           `json:"group_id"`
	GroupRevision   string           `json:"group_revision"`
	AppliedRevision string           `json:"applied_revision"`
	Status          WriteGroupStatus `json:"status"`
	Accepting       bool             `json:"accepting"`
	Reason          string           `json:"reason,omitempty"`
}

// WriteGroupBacklogOwnershipGuard proves that accepted records still retain
// an immutable group/revision/destination owner before an applied group is
// tombstoned. Implementations must only read through tx.
type WriteGroupBacklogOwnershipGuard interface {
	CheckWriteGroupBacklog(ctx context.Context, tx *sql.Tx, group *WriteGroup) error
}

// WithBacklogOwnershipGuard installs the transaction-scoped delete guard.
// Without one, deleting a group with an applied revision remains blocked.
func (s *WriteGroupService) WithBacklogOwnershipGuard(guard WriteGroupBacklogOwnershipGuard) *WriteGroupService {
	if s != nil {
		s.backlogGuard = guard
	}
	return s
}

// WithClock installs a deterministic clock for lifecycle bucket boundaries.
// A nil clock restores the service's UTC wall clock.
func (s *WriteGroupService) WithClock(clock func() time.Time) *WriteGroupService {
	if s == nil {
		return nil
	}
	if clock == nil {
		s.now = func() time.Time { return time.Now().UTC() }
	} else {
		s.now = clock
	}
	return s
}

// Disable stops new intake while retaining the applied revision and all
// accepted backlog identities. Group is intentionally optional in mutation.
func (s *WriteGroupService) Disable(ctx context.Context, id string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("disable write group: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, true, false)
	if err != nil {
		return nil, err
	}
	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("disable write group")
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
		if existing.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
			return writeGroupRevisionConflict("write group connector revision is stale")
		}
		// Disabling only stops new intake and never touches the destination, so it
		// must keep working after the endpoint was edited or removed; the group's
		// own saved destination revision above remains the compare-and-swap value.
		candidate := cloneWriteGroup(existing)
		candidate.Status = WriteGroupStatusDisabled
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
		return nil, normalizeWriteGroupServiceError("disable write group", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// Apply validates the saved draft and schedules it for the next UTC bucket.
// No external DDL or runtime activation is performed here.
func (s *WriteGroupService) Apply(ctx context.Context, id string, mutation WriteGroupMutation) (*WriteGroupSaveResult, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("apply write group: %w", ErrWriteGroupValidation)
	}
	mutation, err := prepareWriteGroupMutation(mutation, true, false)
	if err != nil {
		return nil, err
	}
	if mutation.Group != nil {
		return nil, fmt.Errorf("apply write group payload must be omitted: %w", ErrWriteGroupValidation)
	}
	readiness, err := s.Readiness(ctx, id)
	if err != nil {
		return nil, err
	}
	if !readiness.Ready {
		return nil, fmt.Errorf("apply write group: %w", ErrWriteGroupApplyNotReady)
	}

	var saved *WriteGroup
	updated, err := s.workspaceSvc.UpdateDatabaseSetup(ctx, mutation.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != mutation.WorkspaceID {
			return writeGroupNotFound("apply write group")
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
		if existing.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
			return writeGroupRevisionConflict("write group connector revision is stale")
		}
		if err := validateLiveWriteGroupConnector(ctx, tx, s.repo, existing.Destination.ConnectorID, mutation.ExpectedConnectorRevision); err != nil {
			return err
		}
		candidate := cloneWriteGroup(existing)
		if err := s.repo.validate(ctx, tx, candidate); err != nil {
			return err
		}
		if readiness.runtimeLayout != nil {
			currentTypes, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, candidate)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(currentTypes, readiness.runtimeLayout.TagTypes) {
				return writeGroupRevisionConflict("write group tag types changed during schema verification")
			}
		}
		candidate.Status = WriteGroupStatusReady
		candidate.UpdatedAt = s.clockNow()
		activeSnapshot, snapshotErr := s.repo.versionByRevisionInTx(ctx, tx, record.ID, existing.ID, existing.AppliedRevision)
		if snapshotErr != nil && !errors.Is(snapshotErr, ErrWriteGroupAppliedRevisionUnavailable) {
			return snapshotErr
		}
		if activeSnapshot != nil && sameWriteGroupAppliedSemantics(candidate, activeSnapshot.Group) {
			// A display-only rename must not restart the applied writer or create
			// a second bucket snapshot. Keep the prior active revision while
			// explicitly moving the saved draft back to ready.
			candidate.AppliedRevision = existing.AppliedRevision
			activeSnapshot.RuntimeLayout = readiness.runtimeLayout
			if err := s.saveRuntimeVersionInTx(ctx, tx, activeSnapshot); err != nil {
				return err
			}
			if err := s.repo.setAppliedInTx(ctx, tx, candidate); err != nil {
				return err
			}
			if err := projectWriteGroup(record, candidate); err != nil {
				return err
			}
			saved = cloneWriteGroup(candidate)
			return nil
		}
		if activeSnapshot != nil && activeSnapshot.EffectiveAt.After(candidate.UpdatedAt) {
			// Keep one pending semantic cutover so an older scheduled snapshot
			// cannot become effective again after a newer policy was applied.
			return fmt.Errorf("write-group cutover is still pending: %w", ErrWriteGroupLifecycleBlocked)
		}
		candidate.AppliedRevision = candidate.Revision
		oldIntervalSeconds := 0
		if activeSnapshot != nil {
			oldIntervalSeconds = activeSnapshot.Group.RowPolicy.IntervalSeconds
		}
		effectiveAt, err := nextWriteGroupCutover(
			candidate.UpdatedAt,
			oldIntervalSeconds,
			candidate.RowPolicy.IntervalSeconds,
		)
		if err != nil {
			return fmt.Errorf("schedule write-group cutover: %w", err)
		}
		if err := s.repo.insertVersionInTx(ctx, tx, candidate, effectiveAt); err != nil {
			return err
		}
		if err := s.saveRuntimeVersionInTx(ctx, tx, &WriteGroupAppliedSnapshot{
			WorkspaceID: record.ID, GroupID: candidate.ID, GroupRevision: candidate.Revision,
			AppliedRevision: candidate.AppliedRevision, EffectiveAt: effectiveAt.UTC(),
			Group: candidate, RuntimeLayout: readiness.runtimeLayout,
		}); err != nil {
			return err
		}
		if err := s.repo.setAppliedInTx(ctx, tx, candidate); err != nil {
			return err
		}
		if err := projectWriteGroup(record, candidate); err != nil {
			return err
		}
		if err := s.activateWriterOwnershipInTx(ctx, tx, WriterOwnershipActivationRequest{
			WorkspaceID:                   record.ID,
			GroupID:                       candidate.ID,
			ExpectedDatabaseSetupRevision: mutation.ExpectedWorkspaceRevision,
			ExpectedGroupRevision:         mutation.ExpectedGroupRevision,
			ExpectedConnectorRevision:     mutation.ExpectedConnectorRevision,
			PreviousAppliedRevision:       existing.AppliedRevision,
			AppliedRevision:               candidate.AppliedRevision,
			EffectiveAt:                   effectiveAt.UTC(),
		}); err != nil {
			return err
		}
		saved = cloneWriteGroup(candidate)
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("apply write group", err)
	}
	return &WriteGroupSaveResult{WorkspaceRevision: updated.DatabaseSetupRevision, Group: saved}, nil
}

// ResolveAppliedAt returns the immutable snapshot effective at at. Tombstones
// remain resolvable so accepted backlog can retain its original owner.
func (s *WriteGroupService) ResolveAppliedAt(ctx context.Context, id string, at time.Time) (result *WriteGroupAppliedSnapshot, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("resolve write-group snapshot")
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("resolve write-group snapshot", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback write-group snapshot read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("resolve write-group snapshot")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("resolve write-group snapshot", err)
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("resolve write-group snapshot", err)
	}
	if at.IsZero() {
		at = s.clockNow()
	}
	result, err = s.repo.resolveVersionInTx(ctx, setupTx.SQLTx(), record.ID, group.ID, at)
	if errors.Is(err, ErrWriteGroupAppliedRevisionUnavailable) {
		return nil, fmt.Errorf("resolve write-group snapshot: %w", ErrWriteGroupLifecycleBlocked)
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("resolve write-group snapshot", err)
	}
	result.WorkspaceID = record.ID
	result.GroupID = group.ID
	result.Group = cloneWriteGroup(result.Group)
	return result, nil
}

// IntakeEligibility reports whether a new sample may be accepted now. Draft
// edits keep an existing applied snapshot eligible; disabled/deleted groups
// always stop new intake while retaining their snapshots for delivery.
func (s *WriteGroupService) IntakeEligibility(ctx context.Context, id string) (result *WriteGroupIntakeEligibility, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("read write-group intake eligibility")
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read write-group intake eligibility", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback write-group intake read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("read write-group intake eligibility")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read write-group intake eligibility", err)
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read write-group intake eligibility", err)
	}
	result = &WriteGroupIntakeEligibility{
		WorkspaceID: record.ID, GroupID: group.ID, GroupRevision: group.Revision,
		AppliedRevision: group.AppliedRevision, Status: group.Status,
	}
	if group.Status == WriteGroupStatusDisabled {
		result.Reason = "group-disabled"
		return result, nil
	}
	if group.Status == WriteGroupStatusDeleted {
		result.Reason = "group-deleted"
		return result, nil
	}
	snapshot, resolveErr := s.repo.resolveVersionInTx(ctx, setupTx.SQLTx(), record.ID, group.ID, s.clockNow())
	if errors.Is(resolveErr, ErrWriteGroupAppliedRevisionUnavailable) {
		if group.AppliedRevision == "" {
			result.Reason = "no-applied-revision"
		} else {
			result.Reason = "pending-activation"
		}
		return result, nil
	}
	if resolveErr != nil {
		return nil, normalizeWriteGroupServiceError("read write-group intake eligibility", resolveErr)
	}
	result.AppliedRevision = snapshot.GroupRevision
	result.Accepting = true
	result.Reason = "applied-snapshot-active"
	return result, nil
}

const maxWriteGroupDurationSeconds = int64(1<<63-1) / int64(time.Second)

func nextWriteGroupCutover(now time.Time, oldIntervalSeconds, newIntervalSeconds int) (time.Time, error) {
	newInterval, err := writeGroupIntervalDuration(newIntervalSeconds)
	if err != nil {
		return time.Time{}, err
	}
	if oldIntervalSeconds == 0 {
		return nextWriteGroupBucketDuration(now, newInterval), nil
	}
	oldInterval, err := writeGroupIntervalDuration(oldIntervalSeconds)
	if err != nil {
		return time.Time{}, err
	}
	commonInterval, err := writeGroupCommonInterval(oldInterval, newInterval)
	if err != nil {
		return time.Time{}, err
	}
	return nextWriteGroupBucketDuration(now, commonInterval), nil
}

func writeGroupIntervalDuration(intervalSeconds int) (time.Duration, error) {
	if intervalSeconds <= 0 {
		return 0, fmt.Errorf("interval must be positive: %w", ErrWriteGroupValidation)
	}
	seconds := int64(intervalSeconds)
	if seconds > maxWriteGroupDurationSeconds {
		return 0, fmt.Errorf("interval seconds exceed time.Duration: %w", ErrWriteGroupIntervalOverflow)
	}
	return time.Duration(seconds) * time.Second, nil
}

func nextWriteGroupBucketDuration(now time.Time, interval time.Duration) time.Time {
	// Truncate(interval) uses Go's year-one zero time, which does not align
	// arbitrary integer-second intervals with Unix-epoch snapshot buckets.
	seconds := int64(interval / time.Second)
	remainder := now.Unix() % seconds
	if remainder < 0 {
		remainder += seconds
	}
	return now.UTC().Truncate(time.Second).Add(time.Duration(seconds-remainder) * time.Second)
}

func writeGroupCommonInterval(left, right time.Duration) (time.Duration, error) {
	leftSeconds := int64(left / time.Second)
	rightSeconds := int64(right / time.Second)
	divisor := greatestCommonDivisor(leftSeconds, rightSeconds)
	quotient := leftSeconds / divisor
	if quotient > maxWriteGroupDurationSeconds/rightSeconds {
		return 0, fmt.Errorf("old and new bucket boundary exceeds time.Duration: %w", ErrWriteGroupIntervalOverflow)
	}
	return time.Duration(quotient*rightSeconds) * time.Second, nil
}

func greatestCommonDivisor(left, right int64) int64 {
	for right != 0 {
		left, right = right, left%right
	}
	return left
}

func (s *WriteGroupService) clockNow() time.Time {
	if s != nil && s.now != nil {
		return s.now().UTC()
	}
	return time.Now().UTC()
}

func sameWriteGroupAppliedSemantics(candidate, applied *WriteGroup) bool {
	if candidate == nil || applied == nil {
		return false
	}
	left := cloneWriteGroup(candidate)
	right := cloneWriteGroup(applied)
	for _, group := range []*WriteGroup{left, right} {
		group.ID = ""
		group.WorkspaceID = ""
		group.Revision = ""
		group.AppliedRevision = ""
		group.Name = ""
		group.Status = ""
		group.Migration = WriteGroupMigration{}
		group.CreatedAt = time.Time{}
		group.UpdatedAt = time.Time{}
	}
	return reflect.DeepEqual(left, right)
}
