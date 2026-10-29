package runtime

import (
	"context"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

type counterDuplicateVerifier interface {
	VerifyCounterDuplicate(context.Context, snapshot.Sample) error
}

func (b *GroupBoundary) verifyCounterDuplicate(ctx context.Context, sample snapshot.Sample) error {
	if b.ledger == nil || sample.CounterProof.IsZero() {
		return nil
	}
	verifier, ok := b.ledger.(counterDuplicateVerifier)
	if !ok {
		return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: reasonJournalFailed, Cause: &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofMissing}}
	}
	if err := verifier.VerifyCounterDuplicate(ctx, sample); err != nil {
		return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: reasonJournalFailed, Cause: err}
	}
	return nil
}
