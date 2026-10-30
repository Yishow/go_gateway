package sourcerule

import (
	"context"
	"go-gateway/internal/datalink/runtimeview"
)

// RuntimeSourceObservationResult describes only a finalized reconcile attempt.
type RuntimeSourceObservationResult string

const (
	SourceObservationApplied     RuntimeSourceObservationResult = "applied"
	SourceObservationRejected    RuntimeSourceObservationResult = "rejected"
	SourceObservationUnavailable RuntimeSourceObservationResult = "unavailable"
)

// RuntimeSourceObservation is ephemeral execution evidence, never a saved DTO.
// ProjectionVersion corroborates the actual input only; it is NOT WorkspaceRevision.
// Applied is intermediate until the outer mutation finalizes all authorities.
type RuntimeSourceObservation struct {
	Scope                                                runtimeview.Scope
	RuleID, RuleRevision, ReferenceID, ProjectionVersion string
	Attempted, Applied, Verified                         bool
	Result                                               RuntimeSourceObservationResult
}

// RuntimeSourceObservationFinalizer observes the final result after compensation.
// Implementations must not change the mutation result or perform configuration writes.
type RuntimeSourceObservationFinalizer interface {
	FinalizeSourceRuleObservation(context.Context, RuntimeReconcileOutcome, error)
}

func (s *Service) finalizeSourceRuleObservation(ctx context.Context, outcome RuntimeReconcileOutcome, err error) {
	defer func() { _ = recover() }() // Diagnostics must never replace a successful or failed mutation.
	if observer, ok := s.runtimeSync.(RuntimeSourceObservationFinalizer); ok {
		observer.FinalizeSourceRuleObservation(ctx, outcome, err)
	}
}
