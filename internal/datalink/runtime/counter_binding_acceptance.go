package runtime

import (
	"strings"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
)

// installCounterProofReferences installs resolved immutable facts. A nil map
// is only for the documented non-durable algorithm utility.
func installCounterProofReferences(cfg GroupBoundaryConfig, members map[string]boundaryMember) error {
	if cfg.CounterBindings == nil {
		return nil
	}
	fail := func(reason error) error { return &counterproof.CounterBindingProofError{Reason: reason} }
	for key, proof := range cfg.CounterBindings {
		if key != proof.Member.Key() || proof.Member.WorkspaceID != cfg.Group.WorkspaceID || proof.Member.GroupID != cfg.Group.ID || proof.Member.AppliedRevision != cfg.Group.AppliedRevision {
			return fail(counterproof.ErrCounterBindingProofConflict)
		}
		memberKey := proof.Member.SampleMemberKey()
		member, exists := members[memberKey]
		if !exists || member.counter == nil {
			return fail(counterproof.ErrCounterBindingProofConflict)
		}
		if cfg.Ledger == nil && !mixedBoundaryGroup(cfg.Group) {
			return fail(counterproof.ErrCounterBindingProofMissing)
		}
		ref, err := proof.Reference()
		if err != nil {
			return err
		}
		if !counterProofMatchesBinding(proof, *member.counter) {
			return fail(counterproof.ErrCounterBindingProofConflict)
		}
		result := member.kind
		if member.exactScale != "" {
			result = measurement.ExactDecimal
		}
		if proof.SourceValueType != cfg.TagTypes[proof.Member.TagID] || proof.ResultValueType != result ||
			(proof.ScaleMode == counterproof.CounterScaleExact && proof.ExactScale != member.exactScale) ||
			(proof.ScaleMode == counterproof.CounterScaleUnscaled && member.exactScale != "") {
			return fail(counterproof.ErrCounterBindingProofConflict)
		}
		member.sourceRevision, member.mappingRevision = proof.SourceRevision, proof.MappingRevision
		member.counterProof = ref
		member.counterClaims = counterClaimIdentity{
			measurementID: proof.MeasurementID, definitionRevision: proof.DefinitionRevision,
			sourceBindingRevision: proof.SourceBindingRevision, seriesEpoch: proof.SeriesEpoch,
		}
		members[memberKey] = member
	}
	for _, member := range members {
		if member.counter != nil && member.counterProof.IsZero() {
			return fail(counterproof.ErrCounterBindingProofMissing)
		}
	}
	return nil
}

func counterProofMatchesBinding(proof counterproof.CounterBindingProof, actual aggregation.CounterBinding) bool {
	modulus, rate, err := proof.CounterPolicy.ExactRollover()
	if err != nil {
		return false
	}
	expected := aggregation.CounterBinding{
		WorkspaceID: proof.Member.WorkspaceID, GroupID: proof.Member.GroupID, ProfileRevision: proof.Member.AppliedRevision, MemberKey: proof.Member.SampleMemberKey(),
		MeasurementID: proof.MeasurementID, DefinitionRevision: proof.DefinitionRevision, SourceBindingRevision: proof.SourceBindingRevision, SeriesEpoch: proof.SeriesEpoch,
		Kind: aggregation.CounterKind(proof.SemanticKind), Unit: strings.TrimSpace(proof.Unit), RolloverModulus: modulus, MaxRatePerSecond: rate,
	}
	if proof.ScaleMode == counterproof.CounterScaleExact {
		expected.ExactRecordingScale = proof.ExactScale
	}
	want, expectedErr := expected.Key()
	got, actualErr := actual.Key()
	return expectedErr == nil && actualErr == nil && want == got
}
