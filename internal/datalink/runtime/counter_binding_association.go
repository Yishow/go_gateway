package runtime

import (
	"context"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
)

type counterReplayPreparer interface {
	PrepareCounterReplay(context.Context, []counterproof.CounterBindingReference) error
}

// prepareCounterAssociations runs before journal memory is rebuilt, so an old
// duplicate cannot stand in for a durable association of its original meaning.
func prepareCounterAssociations(ctx context.Context, cfg GroupBoundaryConfig, restored groupdelivery.Restored) (bool, error) {
	needsLegacy := false
	for _, sample := range restored.Samples {
		var ref counterproof.CounterBindingReference
		for _, proof := range cfg.CounterBindings {
			if proof.Member.SampleMemberKey() == sample.MemberKey {
				var err error
				ref, err = proof.Reference()
				if err != nil {
					return false, err
				}
				break
			}
		}
		// A global legacy association result cannot authorize a modified or
		// erased per-sample transfer origin, including non-counter inputs.
		if !restored.ValidateHandoffOrigin(sample, ref) {
			return false, &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofConflict}
		}
		if sample.CounterProof.IsZero() && !ref.IsZero() {
			// Incoming legacy bytes retain their original source association.
			needsLegacy = needsLegacy || !restored.AuthorizesHandoff(sample, ref)
		}
	}
	if !needsLegacy {
		return false, nil
	}
	preparer, ok := cfg.Ledger.(counterReplayPreparer)
	if !ok || len(cfg.CounterBindings) == 0 {
		return false, nil
	}
	refs := make([]counterproof.CounterBindingReference, 0, len(cfg.CounterBindings))
	for _, proof := range cfg.CounterBindings {
		ref, err := proof.Reference()
		if err != nil {
			return false, err
		}
		refs = append(refs, ref)
	}
	if err := preparer.PrepareCounterReplay(ctx, refs); err != nil {
		return false, err
	}
	return true, nil
}
