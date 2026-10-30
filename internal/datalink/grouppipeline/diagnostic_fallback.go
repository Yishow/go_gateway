package grouppipeline

import (
	"context"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
)

func (c OwnerObservationConfig) enabled() bool { return c.ObserveCaptured != nil || c.Observe != nil }

func reportGenericDiagnostic(config Config) {
	if config.OnDiagnosticError != nil {
		defer func() { _ = recover() }()
		config.OnDiagnosticError()
	}
}

func reportPipelineError(config Config, err error) {
	if err == nil {
		return
	}
	if config.OnError != nil {
		config.OnError(err)
	}
	reportGenericDiagnostic(config)
}

func (p *Pipeline) report(err error) { reportPipelineError(p.config, err) }

// reportTick isolates the descriptor interval of this exact proof-gated Tick.
// Native reporting stays here, before pruning/flush, with the identical error.
// Only the fixed generic diagnostic may wait for post-unlock capture acceptance.
func (p *Pipeline) reportTick(ctx context.Context, m *managed, operation func() error) {
	b, _ := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	start := 0
	if b != nil {
		start = len(b.descriptors)
	}
	err := p.observedProofGate(ctx, m, operation)
	end := 0
	if b != nil {
		end = len(b.descriptors)
	}
	if err == nil {
		return
	}
	if p.config.OnError != nil {
		p.config.OnError(err)
	}
	key := groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}
	if b != nil && b.dropped == 0 && end == start+1 && tickDescriptorCovers(err, key, b.descriptors[start]) {
		b.descriptors[start].pendingGeneric = true
		return
	}
	reportGenericDiagnostic(p.config)
}

// tickDescriptorCovers accepts only a closed native refusal, never errors.Is
// claims, error strings, or any joined error (even a join of one typed leaf).
// The bounded walk also fails closed on cyclic/hostile wrappers without keeping
// an error, another collection, or any cross-operation diagnostic authority.
func tickDescriptorCovers(err error, key groupdelivery.GroupKey, d ownerDescriptor) (covered bool) {
	defer func() {
		if recover() != nil {
			covered = false
		}
	}()
	for range 64 {
		if err == nil {
			return false
		}
		if _, joined := err.(interface{ Unwrap() []error }); joined {
			return false
		}
		switch e := err.(type) {
		case *counterproof.CounterBindingProofError:
			if e == nil {
				return false
			}
			// Identity is intentional: an Is implementation or matching text is not a
			// closed proof-reason leaf and must not cover an unobserved native error.
			reason := ""
			switch e.Reason {
			case counterproof.ErrCounterBindingProofMissing:
				reason = counterproof.ErrCounterBindingProofMissing.Error()
			case counterproof.ErrCounterBindingProofConflict:
				reason = counterproof.ErrCounterBindingProofConflict.Error()
			case counterproof.ErrCounterBindingProofUnsupported:
				reason = counterproof.ErrCounterBindingProofUnsupported.Error()
			default:
				return false
			}
			return d.intake == nil && d.key == key && d.kind == OwnerCounterQuarantine && d.reason == reason
		case *groupdelivery.OwnershipError:
			if e == nil {
				return false
			}
			switch e.Reason {
			case groupdelivery.ErrOwnershipUnavailable.Reason, groupdelivery.ErrOwnershipConflict.Reason, groupdelivery.ErrHandoffProofConflict.Reason:
				return d.intake == nil && d.key == key && d.kind == OwnerRecordingRefusal && d.reason == e.Reason
			}
			return false
		}
		wrapped, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = wrapped.Unwrap()
	}
	return false
}

func (p *Pipeline) flushOwnerObservation(ctx context.Context, b *ownerObservationBatch, d ownerDescriptor) (accepted bool) {
	observation := p.verifyOwnerObservation(ctx, d)
	defer func() {
		if recover() != nil {
			b.drop()
			accepted = false
		}
	}()
	cfg := p.deps.OwnerObservation
	if cfg.ObserveCaptured != nil {
		return cfg.ObserveCaptured(observation) && observation.Available
	}
	if cfg.Observe != nil {
		cfg.Observe(observation)
	}
	return false
}
