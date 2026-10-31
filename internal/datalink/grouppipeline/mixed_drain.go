package grouppipeline

import (
	"context"
	"encoding/json"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

// mixedDrain is a local dispatch bound derived only from accepted inputs.
// It is NOT an actual window end, canonical revision end, or durable cursor.
type mixedDrain struct {
	clockLimit time.Time
	members    map[string]time.Time
}

func (p *Pipeline) retireMixedGroup(ctx context.Context, id string) {
	p.clearBlocked(id)
	p.mu.RLock()
	var current []*managed
	for _, m := range p.boundaries {
		if m.groupID == id {
			current = append(current, m)
		}
	}
	p.mu.RUnlock()
	for _, m := range current {
		p.report(p.beginMixedDrain(ctx, m))
	}
}

// beginMixedDrain fences captured intake before reading durable pending inputs.
// A failed read remains fenced, and a later reconcile can retry preparation.
func (p *Pipeline) beginMixedDrain(ctx context.Context, m *managed) error {
	m.proofGate.Lock()
	defer m.proofGate.Unlock()
	if m.superseded {
		return nil
	}
	m.localDrain.Store(true)
	if m.boundary.NeedsDurableRestore() {
		return p.restoreMixedDrainLocked(ctx, m)
	}
	if m.drain != nil || m.drainFinished.Load() {
		return nil
	}
	key := groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}
	restored, err := p.deps.Store.Restore(ctx, key)
	if err != nil {
		return err
	}
	if len(restored.Samples) == 0 {
		m.drainFinished.Store(true)
		return nil
	}
	frozen, err := p.loadRuntimeSnapshot(ctx, key)
	if err != nil {
		return err
	}
	if frozen == nil || !mixedRecordingGroup(frozen.Group) {
		return errRuntimeSnapshotInvalid
	}
	intervals := map[string]time.Duration{}
	for _, member := range frozen.Group.Members {
		if member.Cadence == nil {
			return errRuntimeSnapshotInvalid
		}
		interval, err := member.Cadence.Duration()
		if err != nil {
			return err
		}
		raw, err := json.Marshal([]string{member.DeviceID, member.PointID, member.TagID})
		if err != nil {
			return err
		}
		intervals[string(raw)] = interval
	}
	drain := &mixedDrain{members: map[string]time.Time{}}
	var through time.Time
	for _, sample := range restored.Samples {
		interval, ok := intervals[sample.MemberKey]
		if !ok {
			return errRuntimeSnapshotInvalid
		}
		_, end, err := snapshot.CheckedWindowBounds(sample.ObservedAt, interval)
		if err != nil {
			return err
		}
		if end.After(drain.members[sample.MemberKey]) {
			drain.members[sample.MemberKey] = end
		}
		if end.After(through) {
			through = end
		}
	}
	var lateness time.Duration
	if seconds := frozen.Group.RowPolicy.AllowedLatenessSeconds; seconds != 0 {
		lateness, err = snapshot.CheckedCadenceDuration(int64(seconds), "second")
		if err != nil {
			return err
		}
	}
	drain.clockLimit, err = snapshot.CheckedWindowAdd(through, lateness, "local-drain")
	if err != nil {
		return err
	}
	m.drain = drain
	return nil
}

// tickManaged runs under proofGate, including captured workers after Disable.
func (p *Pipeline) tickManaged(ctx context.Context, m *managed, now time.Time) error {
	if !m.localDrain.Load() {
		return m.boundary.Tick(ctx, now)
	}
	if m.drainFinished.Load() || m.drain == nil {
		return nil
	}
	if now.After(m.drain.clockLimit) {
		now = m.drain.clockLimit
	}
	if err := m.boundary.Tick(ctx, now); err != nil {
		return err
	}
	key := groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}
	restored, err := p.deps.Store.Restore(ctx, key)
	if err != nil {
		return err
	}
	if restored.Mixed == nil {
		return nil
	}
	for member, end := range m.drain.members {
		progress, ok := restored.Mixed.Members[member]
		if !ok {
			return errRuntimeSnapshotInvalid
		}
		if !progress.End.IsZero() && progress.Start.Before(end) {
			return nil
		}
	}
	if len(restored.Samples) != 0 {
		return nil
	}
	m.drainFinished.Store(true)
	return nil
}
