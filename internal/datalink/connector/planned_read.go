package connector

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/datalink/optimizer"
)

// PlannedPointBinding carries a captured scalar policy and observation identity.
// These opaque proofs are returned unchanged, never resolved from current state.
// The caller owns eligibility checks and must not mutate inputs during a call.
type PlannedPointBinding struct {
	Point                                               optimizer.EffectivePoint
	Profile                                             ReadProfile
	PointRevision, ConfigFingerprint, InstallationScope string
	InstallationLease                                   uint64
}

// PlannedReadRequest selects one immutable optimizer request, not arbitrary raw
// addresses. The adapter executes only this request; fallback belongs to the
// bounded scan owner and cannot be silently invented by an adapter.
type PlannedReadRequest struct {
	Plan                     optimizer.AcquisitionPlan
	RequestIndex             int
	AcquisitionID, RequestID string
	Bindings                 []PlannedPointBinding
}

// PlannedPointResult preserves the original point and all captured metadata even
// on failure. Timestamp is final wire completion, not scheduled scan start.
type PlannedPointResult struct {
	Binding                                   PlannedPointBinding
	Group                                     optimizer.AcquisitionGroup
	PlanFingerprint, AcquisitionID, RequestID string
	ReadResult
	ObservedAt, ReceivedAt time.Time
}

// PlannedReader is an optional adapter capability, not a second collector.
type PlannedReader interface {
	ReadPlanned(context.Context, PlannedReadRequest) ([]PlannedPointResult, error)
}

var ErrPlannedRead = errors.New("invalid or mismatched planned read")
