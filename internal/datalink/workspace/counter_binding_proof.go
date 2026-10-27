package workspace

import "go-gateway/internal/datalink/counterproof"

// CounterBindingMemberIdentity identifies one member of one immutable revision.
type CounterBindingMemberIdentity = counterproof.CounterBindingMemberIdentity

// CounterScaleMode distinguishes native values from exact integer scaling.
type CounterScaleMode = counterproof.CounterScaleMode

// CounterBindingFacts contains the complete frozen counter meaning.
type CounterBindingFacts = counterproof.CounterBindingFacts

// CounterBindingProof is the shared versioned, digest-verified contract.
type CounterBindingProof = counterproof.CounterBindingProof

// CounterBindingProofs is keyed by complete member identity.
type CounterBindingProofs = counterproof.CounterBindingProofs

// CounterBindingProofError classifies unusable evidence without exposing facts.
type CounterBindingProofError = counterproof.CounterBindingProofError

const (
	CounterScaleUnscaled       = counterproof.CounterScaleUnscaled
	CounterScaleExact          = counterproof.CounterScaleExact
	CounterBindingProofVersion = counterproof.CounterBindingProofVersion
)

var (
	ErrCounterBindingProofConflict    = counterproof.ErrCounterBindingProofConflict
	ErrCounterBindingProofUnsupported = counterproof.ErrCounterBindingProofUnsupported
	ErrCounterBindingProofMissing     = counterproof.ErrCounterBindingProofMissing
)

// NewCounterBindingProof seals independently supplied immutable facts.
func NewCounterBindingProof(facts CounterBindingFacts) (CounterBindingProof, error) {
	return counterproof.NewCounterBindingProof(facts)
}
