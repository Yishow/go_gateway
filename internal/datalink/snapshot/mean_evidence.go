package snapshot

import (
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

// Coverage is how much of a window the evidence accounts for. Complete needs a
// saved, verifiable whole-window proof; without one a usable mean stays unknown,
// never complete.
type Coverage string

const (
	// CoverageComplete means the whole window's observations are accounted for
	// by a verifiable proof. No such proof is wired by this change.
	CoverageComplete Coverage = "complete"
	// CoveragePartial means usable values exist while a known failure, gap or
	// truncated window is on record.
	CoveragePartial Coverage = "partial"
	// CoverageUnknown means usable values exist without whole-window proof.
	CoverageUnknown Coverage = "unknown"
	// CoverageMissing means no value can be carried.
	CoverageMissing Coverage = "missing"
)

// ReasonCoverageUnverified marks a usable mean whose window coverage has no
// saved proof. It is never reported as complete.
const ReasonCoverageUnverified = "coverage-unverified"

// Stable public reasons of the mean evidence projection. A raw quality reason
// is never copied into operator-visible evidence; an unknown raw value folds
// into the bounded category of its quality or status.
const (
	MeanReasonBad        = "bad"
	MeanReasonInvalid    = "invalid"
	MeanReasonStale      = "stale"
	MeanReasonNonNumeric = "non-numeric"
	reasonMissingValue   = "missing-value"
)

// MeanEvidence is the value-free evidence of one mean window: the aggregation
// metadata, the actual counts and observed bounds of the accepted set, the
// coverage and stable reasons. It never holds an expected count, a percentage
// or a guessed cause.
type MeanEvidence struct {
	Version         string
	Method          string
	Basis           measurement.MeanBasis
	SourceValueType schema.DataType
	ResultValueType measurement.ExactType
	ValidCount      int64
	RejectedCount   int64
	FirstObservedAt time.Time
	LastObservedAt  time.Time
	Coverage        Coverage
	Reason          string
	// RejectionReasons lists the stable distinct reasons of the rejected
	// observations, sorted. It is empty when nothing was rejected.
	RejectionReasons []string
}
