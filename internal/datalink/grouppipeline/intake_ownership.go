package grouppipeline

import (
	"context"
	"errors"
	"maps"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
)

// Routing avoids asking known former owners to accept a source; it is advisory.
// Each journal transaction independently resolves the durable cutover again.
// A required but absent/quarantined/not-ready owner is retryable NOT-accepted,
// never a memory-only handoff or an ACK via the old worker.
func (p *Pipeline) acceptWithOwnership(ctx context.Context, envelope measurement.SampleEnvelope, current []*managed) error {
	if p.deps.Store == nil {
		return groupdelivery.ErrOwnershipUnavailable
	}
	source := groupdelivery.OwnershipSource{WorkspaceID: envelope.WorkspaceID, DeviceID: envelope.DeviceID, PointID: envelope.PointID, TagID: envelope.TagID}
	route, routeErr := p.deps.Store.ResolveCutoverOwnership(ctx, source, envelope.ObservedAt)
	errs := []error{routeErr}
	if routeErr == nil && len(route.Affected) > 0 && len(route.Owners) == 0 {
		errs = append(errs, groupdelivery.ErrOwnershipConflict)
	}
	attempted := map[groupdelivery.GroupKey]bool{}
	for _, m := range current {
		key := groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}
		if m.drainOnly || !m.boundary.Wants(envelope.DeviceID, envelope.PointID, envelope.TagID) {
			continue
		}
		if routeErr == nil && route.Affected[key] && !route.Owners[key] {
			continue
		}
		invoked := false
		err := p.observedProofGate(ctx, m, func() error {
			invoked = true
			if route.Owners[key] {
				return m.boundary.AcceptOwnedSample(ctx, envelope)
			}
			return m.boundary.AcceptSample(ctx, envelope)
		})
		if err != nil {
			errs = append(errs, err)
		}
		if invoked {
			attempted[key] = true
		}
	}
	for key := range route.Owners {
		if !attempted[key] {
			errs = append(errs, groupdelivery.ErrOwnershipUnavailable)
			p.recordOwnerObservation(ctx, ownerDescriptor{key: key, kind: OwnerRecordingRefusal, reason: groupdelivery.ErrOwnershipUnavailable.Error()})
		}
	}
	// A cutover can commit after advisory routing while a formerly irrelevant
	// cached boundary returns early without attempting a journal transaction.
	// Recheck the owner set before ACK, so such ignored work must be retried.
	// This does not authorize writes: AppendSample is still the linearization.
	latest, err := p.deps.Store.ResolveCutoverOwnership(ctx, source, envelope.ObservedAt)
	if err != nil {
		errs = append(errs, err)
	} else if !maps.Equal(route.Affected, latest.Affected) || !maps.Equal(route.Owners, latest.Owners) {
		errs = append(errs, groupdelivery.ErrOwnershipConflict)
	}
	return errors.Join(errs...)
}

// wantsCutoverSample runs only on a cache miss and never with p.mu held. The
// legacy interest interface has no context, so the existing Owns timeout bound
// also bounds this canonical-scope lookup on the ingestion path.
func (p *Pipeline) wantsCutoverSample(deviceID, pointID, tagID string) bool {
	if p.deps.Store == nil || p.deps.Groups == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	listed, err := p.deps.Groups.List(ctx)
	if err != nil || listed == nil {
		return true
	}
	workspaceID := listed.WorkspaceID
	if workspaceID == "" {
		// Lightweight test/providers may omit the redundant list header. Only one
		// explicit workspace present in their canonical group list is admissible.
		for _, group := range listed.Groups {
			if group == nil || group.WorkspaceID == "" {
				return true
			}
			if workspaceID != "" && workspaceID != group.WorkspaceID {
				return true
			}
			workspaceID = group.WorkspaceID
		}
		if workspaceID == "" {
			return false
		}
	}
	wanted, err := p.deps.Store.HasCutoverSource(ctx, groupdelivery.OwnershipSource{WorkspaceID: workspaceID, DeviceID: deviceID, PointID: pointID, TagID: tagID})
	// An unavailable hint must not silently discard input before the decisive
	// scoped acceptance transaction can report its typed refusal.
	return wanted || err != nil
}
