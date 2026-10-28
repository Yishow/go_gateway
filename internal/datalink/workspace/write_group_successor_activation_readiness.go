package workspace

import (
	"context"
)

// SuccessorActivationReadiness explicitly probes the prepared target. It never
// writes schema, grants ordinary reserved access, or schedules the boundary.
func (s *WriteGroupService) SuccessorActivationReadiness(ctx context.Context, id string, req WriteGroupSuccessorActivationConfirmation) (*WriteGroupSuccessorActivationReadiness, error) {
	envelope, digest, err := normalizeSuccessorActivationEnvelope(req.Envelope)
	if err != nil {
		return nil, err
	}
	if digest != req.ReviewDigest {
		return nil, recreateReservationConflict("activation_review_changed")
	}
	review, snap, err := s.readSuccessorActivation(ctx, id, envelope.Request)
	if err != nil {
		return nil, err
	}
	if review.ReviewDigest != digest {
		return nil, recreateReservationConflict("activation_review_changed")
	}
	result := &WriteGroupSuccessorActivationReadiness{Review: review, snapshot: snap}
	if len(review.Blockers) > 0 {
		result.Readiness = newWriteGroupReadiness(snap.ready.record, snap.ready.group)
		for _, b := range review.Blockers {
			result.Readiness.Issues = append(result.Readiness.Issues, writeGroupReadinessIssue(b.ReasonCode, "successor activation requires reviewed local evidence", snap.reserved.Group.ID))
		}
		return result, nil
	}
	reload := func(ctx context.Context) (*writeGroupReadinessSnapshot, error) {
		current, latest, err := s.readSuccessorActivation(ctx, id, envelope.Request)
		if err != nil {
			return nil, err
		}
		if current.ReviewDigest != digest || len(current.Blockers) > 0 {
			return nil, recreateReservationConflict("activation_review_changed")
		}
		result.snapshot = latest
		return latest.ready, nil
	}
	result.Readiness, err = s.evaluateReadinessSnapshot(ctx, snap.reserved.Group.ID, nil, snap.ready, reload)
	if err != nil {
		return nil, err
	}
	// Even an early blocked result cannot retain an unverified current snapshot.
	if _, err = reload(ctx); err != nil {
		return nil, err
	}
	if !result.Readiness.Ready {
		return result, nil
	}
	result.ReadinessDigest, err = successorDomainDigest("go-gateway/successor-activation-readiness/v1", struct {
		Review    string
		Readiness *WriteGroupReadiness
		Layout    *WriteGroupRuntimeLayout
		Counters  CounterBindingProofs
	}{digest, result.Readiness, result.Readiness.runtimeLayout, result.Readiness.counterBindings})
	if err != nil {
		return nil, err
	}
	return result, nil
}
