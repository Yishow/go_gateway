package grouppipeline

import (
	"cmp"
	"context"
	"errors"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"slices"
)

// CounterReplayQuarantine describes accepted inputs, never SQL output rows.
// Its fields contain only scope, a classified reason and safe repair guidance.
type CounterReplayQuarantine struct {
	WorkspaceID    string `json:"workspace_id"`
	GroupID        string `json:"group_id"`
	GroupRevision  string `json:"group_revision"`
	PendingSamples int    `json:"pending_samples"`
	Reason         string `json:"reason"`
	Action         string `json:"action"`
}

// CounterReplayRepairAction never asks an operator to substitute current
// settings, reset a baseline or delete evidence.
const CounterReplayRepairAction = "preserve-inputs-and-restore-original-revision-evidence-or-contact-maintenance"

func isCounterProofReason(reason string) bool {
	return reason == counterproof.ErrCounterBindingProofMissing.Error() || reason == counterproof.ErrCounterBindingProofConflict.Error() || reason == counterproof.ErrCounterBindingProofUnsupported.Error()
}

// setCounterQuarantine records only a classified proof fault. State is derived
// from durable evidence on restart and cleared only after this original scope
// successfully validates and rebuilds; a draft or another revision cannot clear it.
func (p *Pipeline) setCounterQuarantine(key groupdelivery.GroupKey, reason string) bool {
	if !key.Valid() || !isCounterProofReason(reason) {
		return false
	}
	boundaryID := boundaryKey(key.GroupID, key.GroupRevision)
	for {
		p.mu.RLock()
		active := p.boundaries[boundaryID]
		p.mu.RUnlock()
		matches := active != nil && active.workspaceID == key.WorkspaceID
		// Wait only on this revision, never under Pipeline.mu. An operation which
		// already holds the gate may finish; a copied but unstarted operation cannot.
		if matches {
			active.proofGate.Lock()
		}
		p.mu.Lock()
		if p.boundaries[boundaryID] != active {
			p.mu.Unlock()
			if matches {
				active.proofGate.Unlock()
			}
			continue
		}
		if matches {
			active.proofQuarantined = true
		}
		p.recordCounterQuarantineLocked(key, reason)
		p.mu.Unlock()
		if matches {
			active.proofGate.Unlock()
		}
		return true
	}
}

// recordCounterQuarantineLocked requires Pipeline.mu; active callers also hold
// that managed revision's proofGate, always in that order before Pipeline.mu.
func (p *Pipeline) recordCounterQuarantineLocked(key groupdelivery.GroupKey, reason string) {
	p.quarantines[key] = CounterReplayQuarantine{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, GroupRevision: key.GroupRevision, Reason: reason, Action: CounterReplayRepairAction}
	if old, ok := p.blocked[key.GroupID]; ok && old.Revision == key.GroupRevision {
		delete(p.blocked, key.GroupID)
	}
	if active := p.boundaries[boundaryKey(key.GroupID, key.GroupRevision)]; active != nil && active.workspaceID == key.WorkspaceID {
		delete(p.boundaries, boundaryKey(key.GroupID, key.GroupRevision))
	}
}

func (p *Pipeline) withCounterProofGate(m *managed, operation func() error) error {
	m.proofGate.Lock()
	defer m.proofGate.Unlock()
	if m.superseded {
		return groupdelivery.ErrOwnershipUnavailable
	}
	if m.proofQuarantined {
		return nil
	}
	err := operation()
	var proof *counterproof.CounterBindingProofError
	if errors.As(err, &proof) && proof.Reason != nil && isCounterProofReason(proof.Reason.Error()) {
		m.proofQuarantined = true
		p.mu.Lock()
		// A stale managed instance cannot re-quarantine a newly validated instance.
		if p.boundaries[boundaryKey(m.groupID, m.revision)] == m {
			p.recordCounterQuarantineLocked(groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}, proof.Reason.Error())
		}
		p.mu.Unlock()
	}
	return err
}

func (p *Pipeline) installManaged(m *managed) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.boundaries[boundaryKey(m.groupID, m.revision)] = m
	delete(p.quarantines, groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision})
}

// CounterReplayQuarantines reads exact pending-input counts independently of
// output row counters and bounded attention lists. A read failure is an error,
// never a fabricated zero or a successful clearance.
func (p *Pipeline) CounterReplayQuarantines(ctx context.Context, groupID string) ([]CounterReplayQuarantine, error) {
	p.mu.RLock()
	entries := make([]CounterReplayQuarantine, 0)
	for _, entry := range p.quarantines {
		if entry.GroupID == groupID {
			entries = append(entries, entry)
		}
	}
	p.mu.RUnlock()
	for i := range entries {
		entry := &entries[i]
		count, err := p.deps.Store.PendingSampleCount(ctx, groupdelivery.GroupKey{WorkspaceID: entry.WorkspaceID, GroupID: entry.GroupID, GroupRevision: entry.GroupRevision})
		if err != nil {
			return nil, err
		}
		entry.PendingSamples = count
	}
	slices.SortFunc(entries, func(a, b CounterReplayQuarantine) int {
		if order := cmp.Compare(a.WorkspaceID, b.WorkspaceID); order != 0 {
			return order
		}
		return cmp.Compare(a.GroupRevision, b.GroupRevision)
	})
	return entries, nil
}
