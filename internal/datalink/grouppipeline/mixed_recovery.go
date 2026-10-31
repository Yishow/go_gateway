package grouppipeline

import (
	"context"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

// replaceMixedCache never clears uncertainty or a local drain on the old
// assembler. It restores a fresh worker from the original durable descriptor,
// then fences captured old workers before atomically installing the replacement.
func (p *Pipeline) replaceMixedCache(ctx context.Context, old *managed, snap *workspace.WriteGroupAppliedSnapshot) (*managed, string) {
	old.proofGate.Lock()
	defer old.proofGate.Unlock()
	if old.superseded || old.proofQuarantined != nil {
		return nil, reasonSnapshotUnavailable
	}
	fresh, reason := p.build(ctx, snap)
	if reason != "" {
		return nil, reason
	}
	key := boundaryKey(old.groupID, old.revision)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.boundaries[key] != old {
		return nil, reasonSnapshotUnavailable
	}
	old.superseded = true
	p.boundaries[key] = fresh
	delete(p.quarantines, groupdelivery.GroupKey{WorkspaceID: old.workspaceID, GroupID: old.groupID, GroupRevision: old.revision})
	return fresh, ""
}

// restoreMixedDrainLocked is only for an already-owned cache worker whose
// acknowledgment was uncertain. Caller holds old.proofGate. The old assembler
// is never retried or unquarantined; only original durable inputs are restored.
func (p *Pipeline) restoreMixedDrainLocked(ctx context.Context, old *managed) error {
	if old.proofQuarantined != nil {
		return groupdelivery.ErrOwnershipUnavailable
	}
	key := groupdelivery.GroupKey{WorkspaceID: old.workspaceID, GroupID: old.groupID, GroupRevision: old.revision}
	frozen, err := p.loadRuntimeSnapshot(ctx, key)
	if err != nil {
		return err
	}
	if frozen == nil || !mixedRecordingGroup(frozen.Group) {
		return errRuntimeSnapshotInvalid
	}
	// Keep an already established end unless canonical revision-end authority
	// supplies a newer one. The mixed builder independently validates ownership.
	frozen.EndsAt = old.boundary.ReplayState().Until
	if source, ok := p.deps.Groups.(interface {
		RevisionEnd(context.Context, string, string) (time.Time, error)
	}); ok {
		frozen.EndsAt, err = source.RevisionEnd(ctx, key.GroupID, key.GroupRevision)
		if err != nil {
			return err
		}
	} else if p.deps.Groups != nil {
		canonical, err := p.deps.Groups.ResolveAppliedAt(ctx, key.GroupID, frozen.EffectiveAt)
		if err != nil || canonical == nil || !sameAppliedSnapshot(frozen, canonical) {
			return groupdelivery.ErrOwnershipUnavailable
		}
		frozen.EndsAt = canonical.EndsAt
	}
	fresh, reason, err := p.buildFrozenWithCause(ctx, frozen, time.Time{}, time.Time{})
	if err != nil {
		return err
	}
	if reason != "" || fresh == nil {
		return errRuntimeSnapshotInvalid
	}
	if err = p.beginMixedDrain(ctx, fresh); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	cacheKey := boundaryKey(key.GroupID, key.GroupRevision)
	if p.boundaries[cacheKey] != old {
		return groupdelivery.ErrOwnershipConflict
	}
	old.superseded = true
	p.boundaries[cacheKey] = fresh
	return nil
}
