package runtime

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/workspace"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/sqltx"
)

func (b *GroupBoundary) attachMixedClosure(closure *groupdelivery.Closure, plan snapshot.Plan) error {
	if _, mixed := b.assembler.MixedProgress(); !mixed {
		return nil
	}
	if b.descriptorDigest == "" {
		return &GroupBoundaryError{Code: boundaryRestoreFailed}
	}
	before, err := b.assembler.MixedCheckpoint()
	if err != nil {
		return err
	}
	after, err := b.assembler.MixedCheckpointForPlan(plan)
	if err != nil {
		return err
	}
	closure.Mixed = &groupdelivery.MixedClosureProgress{Before: before, After: after, DescriptorDigest: b.descriptorDigest}
	return nil
}

func counterCheckpointAfterCursor(cp aggregation.CounterCheckpoint, restored groupdelivery.Restored) bool {
	if restored.Mixed == nil {
		return cp.WindowEnd.After(restored.NextClose)
	}
	member, ok := restored.Mixed.Members[cp.Binding.MemberKey]
	return !ok || restored.Mixed.Progress.ClosedThrough.IsZero() || cp.WindowEnd.After(member.Start) || cp.WindowEnd.After(restored.Mixed.Progress.ClosedThrough)
}

// Unknown transaction outcomes are neither retryable failures nor successes.
// Preserve this boundary's pre-commit memory, but block all further writes and
// planning until a newly constructed boundary restores durable authority.
func (b *GroupBoundary) quarantineUncertainCommit(err error) {
	if errors.Is(err, sqltx.ErrCommitOutcomeUnknown) || errors.Is(err, sqltx.ErrConnectionUnusable) {
		b.uncertainCommit = &GroupBoundaryError{Code: "durable-commit-outcome-unknown", Cause: err}
	}
}

// validateMixedRuntimeDescriptor binds the builder inputs to the same immutable
// bytes Restore authenticated. A digest alone is not permission to reinterpret
// those bytes with a caller's changed method, requiredness, type or layout.
func validateMixedRuntimeDescriptor(ctx context.Context, cfg GroupBoundaryConfig, digest string) error {
	reader, ok := cfg.Ledger.(interface {
		RuntimeDescriptor(context.Context) ([]byte, error)
	})
	if !ok {
		return &GroupBoundaryError{Code: "mixed-durable-descriptor-unavailable"}
	}
	raw, err := reader.RuntimeDescriptor(ctx)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != digest {
		return &GroupBoundaryError{Code: "mixed-durable-descriptor-conflict"}
	}
	var frozen workspace.WriteGroupAppliedSnapshot
	if err = json.Unmarshal(raw, &frozen); err != nil {
		return err
	}
	if frozen.Group == nil || frozen.RuntimeLayout == nil || frozen.WorkspaceID != cfg.Group.WorkspaceID || frozen.GroupID != cfg.Group.ID || frozen.AppliedRevision != cfg.Group.AppliedRevision || !frozen.EffectiveAt.Equal(cfg.StartsAt) {
		return &GroupBoundaryError{Code: "mixed-durable-descriptor-conflict"}
	}
	same := func(a, b any) bool {
		left, e1 := json.Marshal(a)
		right, e2 := json.Marshal(b)
		return e1 == nil && e2 == nil && bytes.Equal(left, right)
	}
	proofs, err := counterproof.ResolveSnapshot(raw, counterproof.SnapshotScope{WorkspaceID: frozen.WorkspaceID, GroupID: frozen.GroupID, AppliedRevision: frozen.AppliedRevision})
	if err != nil {
		return err
	}
	layout := frozen.RuntimeLayout
	if !same(frozen.Group, cfg.Group) || !same(layout.TagTypes, cfg.TagTypes) || !same(proofs, cfg.CounterBindings) || layout.Dialect != cfg.Dialect || !same(layout.Columns, cfg.Columns) || layout.ReceiptTableReady != cfg.ReceiptTableReady {
		return &GroupBoundaryError{Code: "mixed-durable-descriptor-conflict"}
	}
	return nil
}

// mixedRestoreCoverage separates historical checkpoint bounds from a newly
// committed ownership cutoff. Only the exact durable end may narrow history;
// normal checkpoint validation remains unchanged.
func mixedRestoreCoverage(cfg GroupBoundaryConfig, restored groupdelivery.Restored) (initial, narrowed time.Time, err error) {
	initial = cfg.EndsAt
	if !mixedBoundaryGroup(cfg.Group) || restored.Mixed == nil {
		return
	}
	previous := restored.Mixed.EndsAt
	if cfg.EndsAt.IsZero() || (!previous.IsZero() && !cfg.EndsAt.Before(previous)) {
		return
	}
	key := groupdelivery.GroupKey{WorkspaceID: cfg.Group.WorkspaceID, GroupID: cfg.Group.ID, GroupRevision: cfg.Group.AppliedRevision}
	if err = restored.OwnershipView.AuthorizeMixedEnd(key, cfg.StartsAt, cfg.EndsAt); err != nil {
		return
	}
	return previous, cfg.EndsAt, nil
}
