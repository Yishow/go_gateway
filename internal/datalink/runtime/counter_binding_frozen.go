package runtime

import (
	"encoding/json"
	"errors"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"
)

// resolveBoundaryCounters keeps the documented no-ledger algorithm utility,
// but durable activation only uses complete immutable descriptor evidence.
func resolveBoundaryCounters(cfg GroupBoundaryConfig) (GroupBoundaryConfig, error) {
	if cfg.Ledger == nil && cfg.CounterBindings == nil {
		return cfg, nil
	}
	group := cfg.Group
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.AppliedRevision, Group: group,
		RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: cfg.TagTypes}, CounterBindings: cfg.CounterBindings,
	})
	if err != nil {
		return cfg, err
	}
	proofs, err := counterproof.ResolveSnapshot(payload, counterproof.SnapshotScope{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.AppliedRevision})
	if err != nil {
		return cfg, err
	}
	cfg.CounterBindings = proofs
	if cfg.Counters == nil {
		cfg.Counters = map[string]CounterSource{}
		for _, proof := range proofs {
			source, ok := CounterSourceFromDefinition(measurement.MeasurementDefinition{
				ID: proof.MeasurementID, DefinitionRevision: proof.DefinitionRevision, SourceBindingRevision: proof.SourceBindingRevision,
				SeriesEpoch: proof.SeriesEpoch, Unit: proof.Unit, SemanticKind: proof.SemanticKind, NumericEncoding: string(proof.SourceValueType), CounterPolicy: &proof.CounterPolicy,
			})
			if !ok {
				return cfg, &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofConflict}
			}
			cfg.Counters[proof.Member.TagID] = source
		}
	}
	return cfg, nil
}

func counterProofBoundaryError(err error) *GroupBoundaryError {
	var typed *counterproof.CounterBindingProofError
	if errors.As(err, &typed) {
		return &GroupBoundaryError{Code: typed.Reason.Error(), Cause: err}
	}
	return &GroupBoundaryError{Code: boundaryRestoreFailed, Cause: err}
}

// validateCounterReplay checks every accepted sample before rebuilding memory.
// Legacy payloads stay byte-identical and require successful durable association
// preparation before they can enter memory.
func validateCounterReplay(members map[string]boundaryMember, restored groupdelivery.Restored, legacyAssociated bool) error {
	conflict := func() error {
		return &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofConflict}
	}
	for _, sample := range restored.Samples {
		member, ok := members[sample.MemberKey]
		if !ok {
			return conflict()
		}
		if !restored.ValidateHandoffOrigin(sample, member.counterProof) {
			return conflict()
		}
		if member.counter == nil {
			if !sample.CounterProof.IsZero() {
				return conflict()
			}
			continue
		}
		authorized := !member.counterProof.IsZero() && restored.AuthorizesHandoff(sample, member.counterProof)
		if member.counterProof.IsZero() || (sample.CounterProof.IsZero() && !legacyAssociated && !authorized) {
			return &counterproof.CounterBindingProofError{Reason: counterproof.ErrCounterBindingProofMissing}
		}
		if !sample.CounterProof.IsZero() && sample.CounterProof != member.counterProof && !authorized {
			return conflict()
		}
		if sample.SourceRevision != member.sourceRevision || sample.MappingRevision != member.mappingRevision {
			return conflict()
		}
		kind := member.kind
		if member.exactScale != "" {
			kind = measurement.ExactDecimal
		}
		if sample.Value.Type() != "" && sample.Value.Type() != kind {
			return conflict()
		}
	}
	for _, baseline := range restored.Baselines {
		checkpoint, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if err != nil {
			continue
		} // existing baseline codec supplies its precise error
		member, ok := members[baseline.MemberKey]
		if !ok || member.counter == nil {
			return conflict()
		}
		got, err := checkpoint.Binding.Key()
		want, wantErr := member.counter.Key()
		if err == nil && wantErr == nil && got != want {
			// Foreign scope is diagnosed by the baseline validator, not reclassified.
			if checkpoint.Binding.WorkspaceID == member.counter.WorkspaceID && checkpoint.Binding.GroupID == member.counter.GroupID && checkpoint.Binding.ProfileRevision == member.counter.ProfileRevision {
				return conflict()
			}
		}
	}
	return nil
}
