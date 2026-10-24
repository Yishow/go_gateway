package workspace

import "context"

// planRollback reviews returning the sources of a forward cutover to the group
// that held them. The source gets a new revision with exactly the content of
// the full revision that ended, applied at a new boundary; the long revisions
// and the remaining wide revision (if the cutover left one) end there too.
func (s *WriteGroupService) planRollback(ctx context.Context, req CutoverRequest) (*cutoverPlan, error) {
	saved, err := s.Get(ctx, req.SourceGroupID)
	if err != nil {
		return nil, err
	}
	source := saved.Group
	if source.WorkspaceID != req.WorkspaceID {
		return nil, writeGroupNotFound("review cutover")
	}
	if source.Revision != req.ExpectedSourceRevision {
		return nil, writeGroupRevisionConflict("write group revision is stale")
	}
	if source.Destination.ConnectorRevision != req.ExpectedConnectorRevision {
		return nil, writeGroupRevisionConflict("write group connector revision is stale")
	}
	prior, rolledBack, err := s.loadCutover(ctx, req.WorkspaceID, req.CutoverID)
	if err != nil {
		return nil, err
	}
	plan := &cutoverPlan{req: req, source: source, prior: prior, readiness: map[string]*WriteGroupReadiness{}, moved: prior.Moved}
	var issues []string
	if prior.Kind != WriteGroupCutoverForward || prior.SourceGroupID != source.ID || rolledBack || len(prior.Ended) != 1 {
		issues = append(issues, CutoverIssueNotRollbackable)
	}
	views := make([]CutoverGroupView, 0, len(prior.Started))
	for _, start := range prior.Started {
		holder, err := s.Get(ctx, start.GroupID)
		if err != nil {
			return nil, err
		}
		ended, err := s.repo.revisionEndInTx(ctx, s.repo.db, start.GroupID, start.Revision)
		if err != nil {
			return nil, err
		}
		// Each revision the cutover started must still be the one applied, and not ended.
		if holder.Group.AppliedRevision != start.Revision || !ended.IsZero() {
			issues = append(issues, CutoverIssueNotRollbackable)
		}
		plan.facts = append(plan.facts, groupFact("started", holder.Group))
		plan.ended = append(plan.ended, CutoverEnd(start))
		if start.GroupID != source.ID {
			views = append(views, cutoverGroupView(holder.Group, holder.Group.Destination.SchemaRevision != ""))
		} else {
			plan.successor = holder.Group
			plan.remaining = refsOf(holder.Group)
		}
	}
	// The source must still be exactly the revision the cutover left in charge. A
	// newer saved draft is somebody's work, and anything else is a decision a
	// rollback must not overwrite.
	switch {
	case source.AppliedRevision != source.Revision:
		issues = append(issues, CutoverIssueSourceDraftPending)
	case plan.successor == nil && len(prior.Ended) == 1 && prior.Ended[0].Revision != source.Revision:
		issues = append(issues, CutoverIssueNotRollbackable)
	}
	// What comes back is the full group that ended, not the remaining one that
	// is ready now: it is proven as it would be restored, against the sources,
	// tag types and destination as they are today.
	if len(prior.Ended) == 1 {
		full, err := s.repo.versionByRevisionInTx(ctx, s.repo.db, source.WorkspaceID, source.ID, prior.Ended[0].Revision)
		if err != nil {
			return nil, err
		}
		restored := cloneWriteGroup(full.Group)
		ready, err := s.readinessOf(ctx, source.ID, restored)
		if err != nil {
			return nil, err
		}
		plan.readiness[source.ID] = ready
		if !ready.Ready {
			issues = append(issues, CutoverIssueRestoredNotReady)
		}
	}
	plan.facts = append(plan.facts, groupFact("source", source))
	boundary, boundaryIssues := s.planBoundary(req, prior.Boundary)
	issues = append(issues, boundaryIssues...)
	return plan.finish(boundary, uniqueSorted(issues), cutoverGroupView(source, source.Destination.SchemaRevision != ""), views, nil)
}

func (s *WriteGroupService) loadCutover(ctx context.Context, workspaceID, id string) (*WriteGroupCutover, bool, error) {
	cutover, err := s.repo.cutoverInTx(ctx, s.repo.db, workspaceID, id)
	if err != nil {
		return nil, false, err
	}
	rolledBack, err := s.repo.rolledBackInTx(ctx, s.repo.db, id)
	return cutover, rolledBack, err
}
