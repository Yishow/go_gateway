package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-gateway/internal/datalink/measurement"
)

// RevisionEnd returns the immutable instant a cutover took the sources away
// from an applied revision; zero while the revision still owns them.
func (s *WriteGroupService) RevisionEnd(ctx context.Context, groupID, revision string) (time.Time, error) {
	if err := s.validate(); err != nil {
		return time.Time{}, err
	}
	return s.repo.revisionEndInTx(ctx, s.repo.db, groupID, revision)
}

// ApplyCutover carries out a reviewed cutover as one transaction. Every
// revision it reads is compared again inside the transaction, and the long
// groups are applied by the same code Apply runs, so a failure anywhere leaves
// the group that held the sources exactly as it was: nothing is disabled first
// and nothing is written twice.
func (s *WriteGroupService) ApplyCutover(ctx context.Context, req CutoverRequest) (*WriteGroupCutover, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	req, err := normalizeCutoverRequest(req)
	if err != nil {
		return nil, err
	}
	if req.Boundary.IsZero() || req.ReviewDigest == "" {
		return nil, fmt.Errorf("apply cutover needs the reviewed boundary and digest: %w", ErrWriteGroupValidation)
	}
	// An apply that committed but was never answered is asked again with the
	// same request: it gets the recorded cutover back and nothing is written.
	if recorded, err := s.recordedCutover(ctx, req); err != nil || recorded != nil {
		return recorded, err
	}
	plan, err := s.planCutover(ctx, req)
	if err != nil {
		return nil, err
	}
	if !plan.review.Applicable {
		return nil, fmt.Errorf("%w: %v", ErrCutoverNotApplicable, plan.review.Issues)
	}
	if plan.review.Digest != req.ReviewDigest {
		return nil, writeGroupRevisionConflict("cutover review is stale")
	}
	var result *WriteGroupCutover
	_, err = s.workspaceSvc.UpdateDatabaseSetup(ctx, req.ExpectedWorkspaceRevision, func(ctx context.Context, tx *sql.Tx, record *Record) error {
		if record.ID != req.WorkspaceID {
			return writeGroupNotFound("apply cutover")
		}
		cutover, err := s.applyCutoverInTx(ctx, tx, record, plan)
		if err != nil {
			return err
		}
		result = cutover
		return nil
	})
	if err != nil {
		return nil, normalizeWriteGroupServiceError("apply cutover", err)
	}
	return result, nil
}

// recordedCutover returns the cutover an identical earlier request already
// made. The review digest finds it, but only the whole intent decides: the same
// digest carried by any other request is a conflict, never that earlier success.
func (s *WriteGroupService) recordedCutover(ctx context.Context, req CutoverRequest) (*WriteGroupCutover, error) {
	recorded, intent, err := s.repo.cutoverByDigestInTx(ctx, s.repo.db, req.WorkspaceID, req.ReviewDigest)
	if err != nil || recorded == nil {
		return nil, err
	}
	want, err := cutoverIntent(req)
	if err != nil {
		return nil, err
	}
	if intent != want {
		return nil, writeGroupRevisionConflict("cutover review belongs to another request")
	}
	return recorded, nil
}

// GetCutover reads one recorded cutover of a source group, so an operator whose
// answer was lost can see whether it happened.
func (s *WriteGroupService) GetCutover(ctx context.Context, sourceGroupID, id string) (*WriteGroupCutover, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	saved, err := s.Get(ctx, sourceGroupID)
	if err != nil {
		return nil, err
	}
	cutover, err := s.repo.cutoverInTx(ctx, s.repo.db, saved.Group.WorkspaceID, id)
	if err != nil {
		return nil, err
	}
	if cutover.SourceGroupID != saved.Group.ID {
		return nil, writeGroupNotFound("read cutover")
	}
	return cutover, nil
}

// ListCutovers lists the recorded cutovers of a source group, oldest first.
func (s *WriteGroupService) ListCutovers(ctx context.Context, sourceGroupID string) ([]*WriteGroupCutover, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	saved, err := s.Get(ctx, sourceGroupID)
	if err != nil {
		return nil, err
	}
	return s.repo.listCutoversInTx(ctx, s.repo.db, saved.Group.WorkspaceID, saved.Group.ID)
}

func (s *WriteGroupService) applyCutoverInTx(ctx context.Context, tx *sql.Tx, record *Record, plan *cutoverPlan) (*WriteGroupCutover, error) {
	req := plan.req
	now := s.clockNow()
	if !req.Boundary.After(now) {
		return nil, fmt.Errorf("cutover boundary has already passed: %w", ErrWriteGroupLifecycleBlocked)
	}
	source, err := s.repo.GetInTx(ctx, tx, record.ID, plan.source.ID)
	if err != nil {
		return nil, err
	}
	if source.Revision != plan.source.Revision || source.AppliedRevision != plan.source.AppliedRevision {
		return nil, writeGroupRevisionConflict("write group revision is stale")
	}
	if err := s.checkLegacyCutoverInTx(ctx, tx, record.ID, req); err != nil {
		return nil, err
	}
	if err := requireNoFinalizedOverlap(ctx, tx, plan); err != nil {
		return nil, err
	}
	id := s.repo.newID()
	cutover := &WriteGroupCutover{
		ID: id, OwnershipFenceRevision: id, WorkspaceID: record.ID, Kind: req.Kind, SourceGroupID: source.ID, Boundary: req.Boundary,
		ReviewDigest: req.ReviewDigest, Moved: plan.moved, CreatedAt: now,
	}
	if plan.prior != nil {
		cutover.PriorCutoverID = plan.prior.ID
	}
	if cutover.intent, err = cutoverIntent(req); err != nil {
		return nil, err
	}
	if req.Kind == WriteGroupCutoverForward {
		err = s.startTargetsAt(ctx, tx, record, plan, cutover)
	} else {
		err = s.restoreSourceAt(ctx, tx, record, plan, source, cutover)
	}
	if err != nil {
		return nil, err
	}
	for _, end := range cutover.Ended {
		if err := s.repo.insertRevisionEndInTx(ctx, tx, record.ID, end.GroupID, end.Revision, cutover.ID, cutover.Boundary); err != nil {
			return nil, err
		}
	}
	if err := s.repo.insertCutoverInTx(ctx, tx, cutover); err != nil {
		return nil, err
	}
	if err := s.repo.insertOwnershipFenceInTx(ctx, tx, cutover); err != nil {
		return nil, err
	}
	return cutover, nil
}

// startTargetsAt applies each long group at the boundary through the ordinary
// apply path, and ends the source's revision at the same instant.
func (s *WriteGroupService) startTargetsAt(ctx context.Context, tx *sql.Tx, record *Record, plan *cutoverPlan, cutover *WriteGroupCutover) error {
	boundary := plan.req.Boundary
	for _, target := range plan.targets {
		mutation := WriteGroupMutation{
			WorkspaceID: record.ID, ExpectedWorkspaceRevision: plan.req.ExpectedWorkspaceRevision,
			ExpectedGroupRevision: target.Revision, ExpectedConnectorRevision: target.Destination.ConnectorRevision,
		}
		applied, err := s.applyGroupInTx(ctx, tx, record, target.ID, mutation, plan.readiness[target.ID],
			func(*WriteGroup, *WriteGroupAppliedSnapshot) (time.Time, error) { return boundary, nil }, true)
		if err != nil {
			return err
		}
		cutover.Started = append(cutover.Started, CutoverStart{GroupID: applied.ID, Revision: applied.AppliedRevision})
	}
	if plan.successor != nil {
		// The sources that were not chosen go on in the same wide group, as the
		// remaining revision saved beforehand, taking over at the same instant.
		mutation := WriteGroupMutation{
			WorkspaceID: record.ID, ExpectedWorkspaceRevision: plan.req.ExpectedWorkspaceRevision,
			ExpectedGroupRevision: plan.successor.Revision, ExpectedConnectorRevision: plan.successor.Destination.ConnectorRevision,
		}
		applied, err := s.applyGroupInTx(ctx, tx, record, plan.successor.ID, mutation, plan.readiness[plan.successor.ID],
			func(*WriteGroup, *WriteGroupAppliedSnapshot) (time.Time, error) { return boundary, nil }, true)
		if err != nil {
			return err
		}
		cutover.Started = append(cutover.Started, CutoverStart{GroupID: applied.ID, Revision: applied.AppliedRevision})
	}
	cutover.Ended = plan.ended
	return nil
}

// restoreSourceAt gives the sources back: the source group gets a new revision
// with the exact content of the one that ended, effective at the new boundary,
// and the long revisions end there. No earlier version, runtime descriptor,
// frozen row or table is rewritten or removed.
func (s *WriteGroupService) restoreSourceAt(ctx context.Context, tx *sql.Tx, record *Record, plan *cutoverPlan, source *WriteGroup, cutover *WriteGroupCutover) error {
	// The full revision that ended is what comes back, whatever the cutover left in its place.
	ended, err := s.repo.versionByRevisionInTx(ctx, tx, record.ID, source.ID, plan.prior.Ended[0].Revision)
	if err != nil {
		return err
	}
	for _, end := range plan.ended {
		target, err := s.repo.GetInTx(ctx, tx, record.ID, end.GroupID)
		if err != nil {
			return err
		}
		if target.AppliedRevision != end.Revision {
			return writeGroupRevisionConflict("write group revision is stale")
		}
	}
	restored := cloneWriteGroup(ended.Group)
	definitions := make(map[string]*measurement.MeasurementDefinition)
	if err := s.repo.validateWithMeasurements(ctx, tx, restored, definitions); err != nil {
		return err
	}
	proofs, err := s.verifyCounterBindingsInTx(ctx, tx, restored, plan.readiness[source.ID], definitions)
	if err != nil {
		return err
	}
	restored.Revision = s.repo.newID()
	restored.AppliedRevision = restored.Revision
	proofs, err = rekeyCounterBindings(proofs, restored.AppliedRevision)
	if err != nil {
		return err
	}
	restored.Status = WriteGroupStatusReady
	restored.UpdatedAt = s.clockNow()
	if err := s.repo.updateGroup(ctx, tx, restored); err != nil {
		return err
	}
	boundary := plan.req.Boundary
	if err := s.repo.insertVersionInTx(ctx, tx, restored, boundary); err != nil {
		return err
	}
	if err := s.saveRuntimeVersionInTx(ctx, tx, &WriteGroupAppliedSnapshot{
		WorkspaceID: record.ID, GroupID: restored.ID, GroupRevision: restored.Revision, AppliedRevision: restored.AppliedRevision,
		EffectiveAt: boundary.UTC(), Group: restored, RuntimeLayout: plan.readiness[source.ID].runtimeLayout, CounterBindings: proofs,
	}); err != nil {
		return err
	}
	if err := s.repo.setAppliedInTx(ctx, tx, restored); err != nil {
		return err
	}
	if err := projectWriteGroup(record, restored); err != nil {
		return err
	}
	if err := s.activateWriterOwnershipInTx(ctx, tx, WriterOwnershipActivationRequest{
		WorkspaceID: record.ID, GroupID: restored.ID, ExpectedDatabaseSetupRevision: plan.req.ExpectedWorkspaceRevision,
		ExpectedGroupRevision: source.Revision, ExpectedConnectorRevision: plan.req.ExpectedConnectorRevision,
		PreviousAppliedRevision: source.AppliedRevision, AppliedRevision: restored.AppliedRevision, EffectiveAt: boundary.UTC(),
	}); err != nil {
		return err
	}
	cutover.Ended = plan.ended
	cutover.Started = []CutoverStart{{GroupID: restored.ID, Revision: restored.Revision}}
	return nil
}

// requireNoActiveOwner refuses an ordinary Apply of a long or wide profile group
// whose sources an active wide or custom group still owns: two writers would record
// the same observation. Only a reviewed cutover may hand them over.
func (s *WriteGroupService) requireNoActiveOwner(ctx context.Context, tx *sql.Tx, candidate *WriteGroup) error {
	if !candidate.isLongLayout() && !candidate.isWideProfile() {
		return nil
	}
	overlap, err := s.repo.activeOverlapInTx(ctx, tx, candidate)
	if err != nil {
		return err
	}
	if len(overlap) > 0 {
		return fmt.Errorf("%w: %d source(s)", ErrWriteGroupCutoverRequired, len(overlap))
	}
	return nil
}
