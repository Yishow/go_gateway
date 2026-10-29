package grouppipeline

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-gateway/internal/datalink/workspace"
)

// ErrAppliedConfigurationUnavailable means immutable configuration proof is absent or invalid.
var ErrAppliedConfigurationUnavailable = errors.New("applied configuration unavailable")

// ErrAppliedConfigurationStale means the read's immutable identities do not agree.
var ErrAppliedConfigurationStale = errors.New("applied configuration stale reference")

// ResolveAppliedConfigurationAt reads and verifies the immutable configuration
// at one explicit server observation. It never probes a destination, repairs a
// descriptor, or creates a runtime boundary. The returned retained snapshot is
// configuration evidence only: neither its presence nor its status grants intake.
func (p *Pipeline) ResolveAppliedConfigurationAt(ctx context.Context, groupID string, asOf time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p == nil || p.deps.Groups == nil || p.deps.Store == nil || p.deps.Store.DB() == nil || strings.TrimSpace(groupID) == "" || asOf.IsZero() {
		return nil, ErrAppliedConfigurationUnavailable
	}
	resolved, err := p.deps.Groups.ResolveAppliedAt(ctx, groupID, asOf)
	if err != nil {
		if errors.Is(err, workspace.ErrWriteGroupNotFound) {
			return nil, workspace.ErrWriteGroupNotFound
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, ErrAppliedConfigurationUnavailable
	}
	if resolved == nil || resolved.Group == nil {
		return nil, ErrAppliedConfigurationUnavailable
	}
	key := groupKey(resolved)
	if !key.Valid() || key.GroupID != groupID || resolved.GroupID != resolved.Group.ID || resolved.WorkspaceID != resolved.Group.WorkspaceID ||
		resolved.GroupRevision != key.GroupRevision || resolved.Group.Revision != key.GroupRevision || resolved.Group.AppliedRevision != key.GroupRevision ||
		resolved.EffectiveAt.IsZero() || resolved.EffectiveAt.After(asOf) {
		return nil, ErrAppliedConfigurationStale
	}
	frozen, err := p.loadRuntimeSnapshot(ctx, key)
	if err != nil || frozen == nil || frozen.RuntimeLayout == nil || frozen.GroupRevision != key.GroupRevision {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, ErrAppliedConfigurationUnavailable
	}
	if !sameAppliedSnapshot(resolved, frozen) {
		return nil, ErrAppliedConfigurationStale
	}
	// The recovery decoder also supports legacy custom descriptors. This public
	// configuration reader claims only layouts whose managed proof is verified.
	if !canonicalManagedGroup(frozen.Group) ||
		(longGroup(frozen.Group) && !validLongRuntimeProof(frozen.Group, frozen.RuntimeLayout)) ||
		(!longGroup(frozen.Group) && !validCanonicalManagedRuntimeProof(frozen.Group, frozen.RuntimeLayout)) {
		return nil, ErrAppliedConfigurationUnavailable
	}
	// JSON decoding made a detached descriptor; EndsAt instead belongs to the
	// canonical ledger read and must never be replaced by absent/stale frozen data.
	frozen.EndsAt = resolved.EndsAt
	return frozen, nil
}
