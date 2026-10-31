package grouppipeline

import (
	"context"
	"errors"
	"sort"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Receipt participants select cache work only. Each replacement is constructed
// from its original frozen descriptor and Store.Restore's transactional view.
func (p *Pipeline) reconcileHandoffs(ctx context.Context, list *workspace.WriteGroupListResult) error {
	if p.deps.Store == nil {
		return nil
	}
	workspaces := map[string]bool{}
	ready := map[groupdelivery.GroupKey]bool{}
	if list.WorkspaceID != "" {
		workspaces[list.WorkspaceID] = true
	}
	for _, group := range list.Groups {
		if group == nil || group.WorkspaceID == "" {
			continue
		}
		if list.WorkspaceID != "" && group.WorkspaceID != list.WorkspaceID {
			continue
		}
		workspaces[group.WorkspaceID] = true
		if group.Status != workspace.WriteGroupStatusDisabled && group.Status != workspace.WriteGroupStatusDeleted && group.AppliedRevision != "" {
			ready[groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}] = true
		}
	}
	var failures []error
	for id := range workspaces {
		groups, err := p.prepareWorkspaceHandoffs(ctx, id)
		failures = append(failures, err)
		// An empty or wholly pre-boundary journal has no receipt participants,
		// but its cached assembler still predates the committed owner epoch.
		// Discovery is only a cache hint; rebuild captures inputs and authority
		// together again, under the same proof gate used for receipt participants.
		p.mu.RLock()
		var candidates []*managed
		for _, m := range p.boundaries {
			if m.workspaceID == id {
				candidates = append(candidates, m)
			}
		}
		p.mu.RUnlock()
		for _, m := range candidates {
			key := groupdelivery.GroupKey{WorkspaceID: id, GroupID: m.groupID, GroupRevision: m.revision}
			view, err := p.deps.Store.OwnershipView(ctx, key)
			if err != nil {
				failures = append(failures, groupdelivery.ErrOwnershipUnavailable)
				continue
			}
			if m.boundary.ReplayState().OwnershipView != view {
				groups = append(groups, []groupdelivery.GroupKey{key})
			}
		}
		for _, keys := range handoffComponents(groups) {
			if err := ctx.Err(); err != nil {
				return errors.Join(append(failures, err)...)
			}
			failures = append(failures, p.rebuildHandoffComponent(ctx, keys, ready))
		}
	}
	return errors.Join(failures...)
}

func handoffComponents(groups [][]groupdelivery.GroupKey) [][]groupdelivery.GroupKey {
	parent := map[groupdelivery.GroupKey]groupdelivery.GroupKey{}
	var root func(groupdelivery.GroupKey) groupdelivery.GroupKey
	root = func(key groupdelivery.GroupKey) groupdelivery.GroupKey {
		if parent[key] != key {
			parent[key] = root(parent[key])
		}
		return parent[key]
	}
	for _, group := range groups {
		for _, key := range group {
			if _, ok := parent[key]; !ok {
				parent[key] = key
			}
		}
		for _, key := range group[1:] {
			parent[root(key)] = root(group[0])
		}
	}
	sets := map[groupdelivery.GroupKey][]groupdelivery.GroupKey{}
	for key := range parent {
		r := root(key)
		sets[r] = append(sets[r], key)
	}
	var components [][]groupdelivery.GroupKey
	for _, keys := range sets {
		sortHandoffKeys(keys)
		components = append(components, keys)
	}
	sort.Slice(components, func(i, j int) bool { return handoffKeyLess(components[i][0], components[j][0]) })
	return components
}

func handoffKeyLess(a, b groupdelivery.GroupKey) bool {
	if a.WorkspaceID != b.WorkspaceID {
		return a.WorkspaceID < b.WorkspaceID
	}
	if a.GroupID != b.GroupID {
		return a.GroupID < b.GroupID
	}
	return a.GroupRevision < b.GroupRevision
}
func sortHandoffKeys(keys []groupdelivery.GroupKey) {
	sort.Slice(keys, func(i, j int) bool { return handoffKeyLess(keys[i], keys[j]) })
}

func (p *Pipeline) rebuildHandoffComponent(ctx context.Context, keys []groupdelivery.GroupKey, ready map[groupdelivery.GroupKey]bool) error {
	old := make(map[groupdelivery.GroupKey]*managed, len(keys))
	p.mu.RLock()
	for _, key := range keys {
		old[key] = p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)]
	}
	p.mu.RUnlock()
	var locked []*managed
	var blocked *workspace.WriteGroup
	var blockedReason string
	// proofGate precedes Pipeline.mu. Keep all related workers gated through
	// construction/publication, but never gate an unrelated receipt component.
	defer func() {
		for i := len(locked) - 1; i >= 0; i-- {
			locked[i].proofGate.Unlock()
		}
		if blocked != nil {
			p.setBlocked(ctx, blocked, blockedReason)
		}
	}()
	for _, key := range keys {
		if m := old[key]; m != nil {
			m.proofGate.Lock()
			locked = append(locked, m)
			if m.workspaceID != key.WorkspaceID || m.superseded || m.proofQuarantined != nil {
				return groupdelivery.ErrOwnershipUnavailable
			}
		}
	}
	p.mu.RLock()
	unchanged := true
	for _, key := range keys {
		if p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)] != old[key] {
			unchanged = false
		}
	}
	p.mu.RUnlock()
	if !unchanged {
		return groupdelivery.ErrOwnershipConflict
	}
	replacements := map[groupdelivery.GroupKey]*managed{}
	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return err
		}
		frozen, err := p.loadRuntimeSnapshot(ctx, key)
		if err != nil || frozen == nil {
			return groupdelivery.ErrOwnershipUnavailable
		}
		canonical, err := p.deps.Groups.ResolveAppliedAt(ctx, key.GroupID, frozen.EffectiveAt)
		if err != nil || canonical == nil || !sameAppliedSnapshot(frozen, canonical) {
			return groupdelivery.ErrOwnershipUnavailable
		}
		frozen.EndsAt = canonical.EndsAt
		var first, until time.Time
		drainOnly := !ready[key]
		mixed := mixedRecordingGroup(frozen.Group)
		if previous := old[key]; previous != nil {
			if mixed {
				// Member cursors and the canonical end are restored by the
				// existing mixed builder, never through a zero common interval.
				drainOnly = drainOnly || previous.localDrain.Load() || previous.drainOnly
			} else {
				state := previous.boundary.ReplayState()
				first = snapshot.BucketStart(state.NotBefore, previous.interval)
				if !state.Until.IsZero() {
					until = alignUp(state.Until, previous.interval)
				}
				drainOnly = previous.drainOnly
			}
		} else if drainOnly && !mixed {
			// Historical/disabled participants may drain accepted input, never become
			// new intake owners merely because a receipt mentions them.
			restored, err := p.deps.Store.Restore(ctx, key)
			if err != nil {
				return groupdelivery.ErrOwnershipUnavailable
			}
			interval := time.Duration(frozen.Group.RowPolicy.IntervalSeconds) * time.Second
			first = restored.NextClose
			for _, sample := range restored.Samples {
				start := snapshot.BucketStart(sample.ObservedAt, interval)
				if first.IsZero() || start.Before(first) {
					first = start
				}
				if end := start.Add(interval); end.After(until) {
					until = end
				}
			}
			if until.IsZero() {
				continue
			}
		}
		fresh, reason := p.buildAt(ctx, frozen, first, until)
		if reason != "" {
			blocked, blockedReason = frozen.Group, reason
			return groupdelivery.ErrOwnershipUnavailable
		}
		if mixed {
			if drainOnly {
				if err := p.beginMixedDrain(ctx, fresh); err != nil {
					return err
				}
			}
		} else {
			fresh.drainOnly = drainOnly
		}
		if previous := old[key]; previous != nil && !previous.boundary.NeedsDurableRestore() && previous.boundary.ReplayState() == fresh.boundary.ReplayState() && previous.localDrain.Load() == fresh.localDrain.Load() {
			continue
		}
		replacements[key] = fresh
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, key := range keys {
		if p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)] != old[key] {
			return groupdelivery.ErrOwnershipConflict
		}
	}
	for key, fresh := range replacements {
		if previous := old[key]; previous != nil {
			previous.superseded = true
		}
		p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)] = fresh
		delete(p.quarantines, key)
	}
	return nil
}
