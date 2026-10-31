package runtime

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

type handoffRecognizer interface {
	RecognizeTransferredSample(context.Context, snapshot.Sample) (groupdelivery.HandoffRecognition, error)
}

// AcceptSample routes a typed sample to its bucket. Samples that do not
// belong to this group are ignored, not errors. Refusals (foreign workspace,
// stale revision, late, identity conflict) return a GroupSampleError.
func (b *GroupBoundary) AcceptSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	return b.acceptSample(ctx, envelope, false)
}

// AcceptOwnedSample requires a durable acceptance by this selected owner. A
// missing member, mismatched cached window or absent ledger cannot become an
// ignored-but-successful ACK. The ledger revalidates ownership in its write tx.
func (b *GroupBoundary) AcceptOwnedSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	return b.acceptSample(ctx, envelope, true)
}

func (b *GroupBoundary) acceptSample(ctx context.Context, envelope measurement.SampleEnvelope, owned bool) error {
	if owned && b.ledger == nil {
		return groupdelivery.ErrOwnershipUnavailable
	}
	member, ok := b.members[memberKey(envelope.DeviceID, envelope.PointID, envelope.TagID)]
	if !ok {
		if owned {
			return groupdelivery.ErrOwnershipConflict
		}
		return nil
	}
	if !b.notBefore.IsZero() && envelope.ObservedAt.Before(b.notBefore) {
		if owned {
			return groupdelivery.ErrOwnershipConflict
		}
		return nil // an earlier applied revision owns this sample
	}
	sample, err := sampleFrom(member, envelope)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.uncertainCommit != nil {
		return groupSamplePersistenceError(b.uncertainCommit)
	}
	if !b.until.IsZero() && !envelope.ObservedAt.Before(b.until) {
		if owned {
			return groupdelivery.ErrOwnershipConflict
		}
		return nil // a newer applied revision owns this sample
	}
	if envelope.WorkspaceID != b.workspaceID {
		return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: reasonWorkspaceMismatch}
	}
	if recognizer, ok := b.ledger.(handoffRecognizer); ok {
		recognized, err := recognizer.RecognizeTransferredSample(ctx, sample)
		if err != nil {
			return groupSamplePersistenceError(err)
		}
		if recognized.Found {
			if recognized.Finalized {
				if recognized.Sample != nil {
					return groupSamplePersistenceError(groupdelivery.ErrOwnershipUnavailable)
				}
				return nil // durable receipt, not reconstructed reclaimed bytes
			}
			if recognized.Sample == nil {
				return groupSamplePersistenceError(groupdelivery.ErrOwnershipUnavailable)
			}
			if err := b.validateDestinationSample(*recognized.Sample); err != nil {
				return err
			}
			// This is already accepted input, not a new late arrival. Do not
			// replace its proof or advance the assembler's captured handoff view.
			result := b.assembler.Replay(*recognized.Sample)
			switch result.Outcome {
			case snapshot.OfferSelected, snapshot.OfferNotSelected, snapshot.OfferDuplicate, snapshot.OfferLate:
				return nil
			default:
				return &GroupSampleError{Outcome: result.Outcome, Reason: result.Reason}
			}
		}
	}
	now := b.clock()
	result := b.assembler.Check(sample, now)
	switch result.Outcome {
	case snapshot.OfferDuplicate:
		// Even a cached duplicate must prove its durable authority: cutover
		// may have committed after this boundary last reconciled.
		if b.ledger == nil {
			return nil
		}
		if !sample.CounterProof.IsZero() {
			return b.verifyCounterDuplicate(ctx, sample)
		}
	case snapshot.OfferSelected, snapshot.OfferNotSelected:
	default:
		if result.Outcome == snapshot.OfferLate {
			b.assembler.Offer(sample, now) // counts the late sample; no state changes
		}
		return &GroupSampleError{Outcome: result.Outcome, Reason: result.Reason}
	}
	if err := b.validateDestinationSample(sample); err != nil {
		return err
	}
	if b.ledger != nil {
		// The ACK is the commit: refuse, and leave memory untouched, if it fails.
		var start time.Time
		if member.interval > 0 {
			var err error
			start, _, err = snapshot.CheckedWindowBounds(sample.ObservedAt, member.interval)
			if err != nil {
				return groupSamplePersistenceError(err)
			}
		} else {
			start = snapshot.BucketStart(sample.ObservedAt, b.interval)
		}
		if err := b.ledger.AppendSample(ctx, start, sample); err != nil {
			b.quarantineUncertainCommit(err)
			return groupSamplePersistenceError(err)
		}
	}
	b.assembler.Offer(sample, now)
	return nil
}

func groupSamplePersistenceError(err error) error {
	reason := reasonJournalFailed
	outcome := snapshot.OfferRejected
	var quotaErr *groupdelivery.QuotaError
	var ownerErr *groupdelivery.OwnershipError
	switch {
	case errors.Is(err, groupdelivery.ErrSampleConflict):
		reason, outcome = snapshot.ReasonIdentityConflict, snapshot.OfferConflict
	case errors.As(err, &ownerErr):
		reason = ownerErr.Reason
	case errors.As(err, &quotaErr):
		reason = quotaErr.Reason
	}
	return &GroupSampleError{Outcome: outcome, Reason: reason, Cause: err}
}
