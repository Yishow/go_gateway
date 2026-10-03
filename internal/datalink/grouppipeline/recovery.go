package grouppipeline

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
)

func groupKey(snap *workspace.WriteGroupAppliedSnapshot) groupdelivery.GroupKey {
	return groupdelivery.GroupKey{WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, GroupRevision: snap.AppliedRevision}
}

func sameAppliedSnapshot(a, b *workspace.WriteGroupAppliedSnapshot) bool {
	return a.WorkspaceID == b.WorkspaceID && a.GroupID == b.GroupID &&
		a.AppliedRevision == b.AppliedRevision && a.EffectiveAt.Equal(b.EffectiveAt) && reflect.DeepEqual(a.Group, b.Group)
}

func (p *Pipeline) saveRuntimeSnapshot(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot) error {
	payload, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	return p.deps.Store.SaveRuntimeVersion(ctx, groupKey(snap), payload)
}

func (p *Pipeline) loadRuntimeSnapshot(ctx context.Context, key groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error) {
	payload, err := p.deps.Store.RuntimeVersion(ctx, key)
	if err != nil || payload == nil {
		return nil, err
	}
	var snap workspace.WriteGroupAppliedSnapshot
	if err := json.Unmarshal(payload, &snap); err != nil {
		return nil, err
	}
	if groupKey(&snap) != key || snap.Group == nil || snap.RuntimeLayout == nil ||
		snap.Group.ID != key.GroupID || snap.Group.WorkspaceID != key.WorkspaceID ||
		snap.Group.Revision != key.GroupRevision || snap.Group.AppliedRevision != key.GroupRevision || snap.EffectiveAt.IsZero() {
		return nil, errors.New("invalid runtime version")
	}
	return &snap, nil
}

func (p *Pipeline) restoreOpenJournals(ctx context.Context) error {
	journals, err := p.deps.Store.OpenJournals(ctx)
	if err != nil {
		return err
	}
	for _, journal := range journals {
		key := boundaryKey(journal.Key.GroupID, journal.Key.GroupRevision)
		p.mu.RLock()
		_, exists := p.boundaries[key]
		p.mu.RUnlock()
		if exists {
			continue
		}
		snap, err := p.loadRuntimeSnapshot(ctx, journal.Key)
		if err != nil {
			p.setBlocked(&workspace.WriteGroup{ID: journal.Key.GroupID, AppliedRevision: journal.Key.GroupRevision}, reasonSnapshotUnavailable)
			continue
		}
		if snap == nil {
			// Legacy journals may be recovered only through their original immutable
			// revision and verified destination, never through the current draft.
			snap, err = p.deps.Groups.ResolveAppliedAt(ctx, journal.Key.GroupID, journal.FirstBucket)
			if err != nil || snap == nil || groupKey(snap) != journal.Key {
				p.setBlocked(&workspace.WriteGroup{ID: journal.Key.GroupID, AppliedRevision: journal.Key.GroupRevision}, reasonSnapshotUnavailable)
				continue
			}
			_, reason := p.build(ctx, snap)
			if reason != "" {
				p.setBlocked(snap.Group, reason)
				continue
			}
			frozen, loadErr := p.loadRuntimeSnapshot(ctx, journal.Key)
			if loadErr != nil || frozen == nil {
				p.setBlocked(snap.Group, reasonSnapshotUnavailable)
				continue
			}
			snap = frozen
		}
		until := journal.LastBucket.Add(time.Duration(snap.Group.RowPolicy.IntervalSeconds) * time.Second)
		m, reason := p.buildFrozen(ctx, snap, journal.FirstBucket, until)
		if reason != "" {
			p.setBlocked(snap.Group, reason)
			continue
		}
		m.drainOnly = true
		p.mu.Lock()
		p.boundaries[key] = m
		p.mu.Unlock()
	}
	return nil
}
