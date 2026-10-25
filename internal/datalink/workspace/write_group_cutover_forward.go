package workspace

import (
	"context"
	"reflect"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
)

func refsOf(group *WriteGroup) []CutoverMemberRef {
	refs := make([]CutoverMemberRef, 0, len(group.Members))
	for _, m := range group.Members {
		refs = append(refs, memberRefOf(m))
	}
	sortMemberRefs(refs)
	return refs
}

func profileKindOf(group *WriteGroup) string {
	if group.RecordingProfile == nil {
		return ""
	}
	return string(group.RecordingProfile.Kind)
}

// onlyMembersRemoved reports whether draft is the applied group with some of its
// sources taken away and nothing else changed: the sources it keeps are exactly
// the applied ones, and layout, policies and destination are the applied ones.
func onlyMembersRemoved(applied, draft *WriteGroup) bool {
	byRef := make(map[CutoverMemberRef]WriteGroupMember, len(applied.Members))
	for _, m := range applied.Members {
		byRef[memberRefOf(m)] = m
	}
	for _, m := range draft.Members {
		if kept, ok := byRef[memberRefOf(m)]; !ok || !reflect.DeepEqual(kept, m) {
			return false
		}
	}
	left, right := cloneWriteGroup(applied), cloneWriteGroup(draft)
	for _, group := range []*WriteGroup{left, right} {
		group.Members, group.ID, group.Revision, group.AppliedRevision, group.Status, group.Name = nil, "", "", "", "", ""
		group.Migration, group.CreatedAt, group.UpdatedAt = WriteGroupMigration{}, time.Time{}, time.Time{}
		// The schema confirmation is per revision; the layout it confirms is what must match.
		group.Destination.SchemaRevision, group.Destination.SchemaDigest = "", ""
	}
	return reflect.DeepEqual(left, right)
}

// omitsRequiredColumn reports whether the remaining revision leaves out a column
// of a moved source that the real destination table cannot do without: NOT NULL
// and no default. The catalog the readiness already read decides; nothing is
// altered and no value is invented. A layout that cannot be read is refused.
func omitsRequiredColumn(applied, remaining *WriteGroup, readiness *WriteGroupReadiness) bool {
	if readiness == nil || readiness.runtimeLayout == nil {
		return true
	}
	kept := make(map[string]struct{}, len(remaining.Members))
	for _, m := range remaining.Members {
		kept[strings.ToLower(m.TargetColumn)] = struct{}{}
	}
	columns := make(map[string]dbtarget.ColumnInfo, len(readiness.runtimeLayout.Columns))
	for _, column := range readiness.runtimeLayout.Columns {
		columns[strings.ToLower(column.Name)] = column
	}
	for _, m := range applied.Members {
		name := strings.ToLower(m.TargetColumn)
		if _, stays := kept[name]; stays {
			continue
		}
		if column, ok := columns[name]; ok && !column.Nullable && !column.HasDefault {
			return true
		}
	}
	return false
}

func sameRefs(a, b []CutoverMemberRef) bool {
	return len(a) == len(b) && reflect.DeepEqual(append([]CutoverMemberRef{}, a...), append([]CutoverMemberRef{}, b...))
}

// reviewForwardSource reads the group whose sources move: the applied revision
// that ends and, when the group was saved with only the sources that stay, the
// remaining revision that takes over at the boundary.
func (s *WriteGroupService) reviewForwardSource(ctx context.Context, plan *cutoverPlan) (issues []string, activeFrom time.Time, err error) {
	source := plan.source
	switch {
	case source.Status == WriteGroupStatusDeleted || source.Status == WriteGroupStatusDisabled:
		return []string{CutoverIssueSourceUnavailable}, time.Time{}, nil
	case source.AppliedRevision == "":
		return []string{CutoverIssueSourceNotApplied}, time.Time{}, nil
	}
	// The applied revision may not be effective yet; its own version and end are
	// what matter, not what a clock currently resolves.
	version, err := s.repo.versionByRevisionInTx(ctx, s.repo.db, source.WorkspaceID, source.ID, source.AppliedRevision)
	if err != nil {
		return nil, time.Time{}, err
	}
	plan.full, activeFrom = version.Group, version.EffectiveAt
	ends, err := s.repo.revisionEndInTx(ctx, s.repo.db, source.ID, source.AppliedRevision)
	if err != nil {
		return nil, time.Time{}, err
	}
	plan.ended = []CutoverEnd{{GroupID: source.ID, Revision: source.AppliedRevision}}
	plan.facts = append(plan.facts, groupFact("source", source))
	if plan.full.isLongLayout() {
		issues = append(issues, CutoverIssueSourceIsLong)
	}
	if !ends.IsZero() {
		issues = append(issues, CutoverIssueSourceEnded)
	}
	if source.Revision == source.AppliedRevision {
		return issues, activeFrom, nil
	}
	plan.successor = source
	ready, err := s.Readiness(ctx, source.ID)
	if err != nil {
		return nil, time.Time{}, err
	}
	plan.readiness[source.ID] = ready
	if !ready.Ready {
		issues = append(issues, CutoverIssueRemainingNotReady)
	}
	if !onlyMembersRemoved(plan.full, source) {
		issues = append(issues, CutoverIssueRemainingChangesPolicy)
	}
	if omitsRequiredColumn(plan.full, source, ready) {
		issues = append(issues, CutoverIssueRemainingOmitsColumn)
	}
	return issues, activeFrom, nil
}

func (s *WriteGroupService) planForward(ctx context.Context, req CutoverRequest) (*cutoverPlan, error) {
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
	plan := &cutoverPlan{req: req, source: source, readiness: map[string]*WriteGroupReadiness{}}
	issues, activeFrom, err := s.reviewForwardSource(ctx, plan)
	if err != nil {
		return nil, err
	}
	owned := source
	if plan.full != nil {
		owned = plan.full
	}
	ownedRefs := make(map[CutoverMemberRef]struct{}, len(owned.Members))
	for _, m := range owned.Members {
		ownedRefs[memberRefOf(m)] = struct{}{}
	}
	if len(req.Targets) == 0 {
		issues = append(issues, CutoverIssueNoTargets)
	}
	// One observation may feed several profiles of a circuit, but never two groups
	// of the same kind: that would record it twice.
	taken := map[CutoverMemberRef]map[string]bool{}
	seen := map[string]bool{}
	views := make([]CutoverGroupView, 0, len(req.Targets))
	for _, ref := range req.Targets {
		saved, err := s.Get(ctx, ref.GroupID)
		if err != nil {
			return nil, err
		}
		target := saved.Group
		if target.WorkspaceID != req.WorkspaceID {
			return nil, writeGroupNotFound("review cutover")
		}
		if target.Revision != ref.ExpectedRevision {
			return nil, writeGroupRevisionConflict("write group revision is stale")
		}
		ready, err := s.Readiness(ctx, target.ID)
		if err != nil {
			return nil, err
		}
		plan.targets = append(plan.targets, target)
		plan.readiness[target.ID] = ready
		plan.started = append(plan.started, CutoverStart{GroupID: target.ID, Revision: target.Revision})
		plan.facts = append(plan.facts, groupFact("target", target))
		views = append(views, cutoverGroupView(target, target.Destination.SchemaRevision != ""))
		switch {
		case seen[target.ID]:
			issues = append(issues, CutoverIssueTargetsOverlap)
		case !managedCutoverTarget(target):
			issues = append(issues, CutoverIssueTargetNotLong)
		case target.AppliedRevision == target.Revision:
			issues = append(issues, CutoverIssueTargetActive)
		case !ready.Ready:
			issues = append(issues, CutoverIssueTargetNotReady)
		}
		seen[target.ID] = true
		kind := profileKindOf(target)
		for _, m := range target.Members {
			ref := memberRefOf(m)
			if _, ok := ownedRefs[ref]; !ok {
				issues = append(issues, CutoverIssueTargetOutsideSource)
			}
			if taken[ref] == nil {
				taken[ref] = map[string]bool{}
			}
			if taken[ref][kind] {
				issues = append(issues, CutoverIssueTargetsOverlap)
			}
			taken[ref][kind] = true
		}
	}
	var movedRefs, unselected []CutoverMemberRef
	for ref := range ownedRefs {
		if _, ok := taken[ref]; ok {
			movedRefs = append(movedRefs, ref)
		} else {
			unselected = append(unselected, ref)
		}
	}
	sortMemberRefs(movedRefs)
	sortMemberRefs(unselected)
	plan.moved = movedRefs
	switch {
	case plan.successor != nil:
		plan.remaining = refsOf(plan.successor)
		if !sameRefs(plan.remaining, unselected) {
			// What the draft keeps must be exactly what no target takes: nothing is
			// moved by omission and nothing moved is left behind.
			issues = append(issues, CutoverIssueRemainingMismatch)
		}
	case len(unselected) > 0:
		// A source nobody chose is never moved silently. Save the group with only
		// the sources that stay, then review again.
		issues = append(issues, CutoverIssueUnselectedSources)
	}
	boundary, boundaryIssues := s.planBoundary(req, activeFrom)
	issues = append(issues, boundaryIssues...)
	return plan.finish(boundary, uniqueSorted(issues), cutoverGroupView(owned, owned.Destination.SchemaRevision != ""), views, unselected)
}
