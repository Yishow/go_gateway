package runtime

import (
	"errors"
	"maps"
	"slices"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Activation codes specific to cumulative profiles.
const (
	boundaryCounterSemantics = "counter-semantics-required"
	boundaryCounterConflict  = "counter-semantics-conflict"
	boundaryCounterType      = "counter-type-unsupported"
	boundaryCounterPolicy    = "counter-policy-invalid"
	boundaryCounterDurable   = "counter-baseline-not-durable"
	// boundaryCounterBaselineInvalid and ...Version diagnose persisted baseline
	// state that cannot be trusted. They differ from a legitimately revised
	// binding, which restores fine and simply starts a new baseline.
	boundaryCounterBaselineInvalid = "counter-baseline-invalid"
	boundaryCounterBaselineVersion = "counter-baseline-version-unsupported"
)

// BaselineLedger is implemented by a ledger that stores counter baselines in the
// transaction that closes a window and returns them on restore. A cumulative
// profile refuses a ledger that does not say so, because its closures would
// restart every counter from nothing.
type BaselineLedger interface{ PersistsCounterBaselines() bool }

func ledgerPersistsBaselines(ledger GroupLedger) bool {
	baselines, ok := ledger.(BaselineLedger)
	return ok && baselines.PersistsCounterBaselines()
}

// restoreCounterBaselines decodes the durable baselines strictly and installs
// them in the assembler. Damaged, foreign or version-mismatched state is an
// activation failure with its own code; it is never dropped or reinterpreted.
func restoreCounterBaselines(assembler *snapshot.Assembler, restored groupdelivery.Restored) string {
	if len(restored.Baselines) == 0 {
		return ""
	}
	checkpoints := make(map[string]aggregation.CounterCheckpoint, len(restored.Baselines))
	for _, baseline := range restored.Baselines {
		checkpoint, err := aggregation.DecodeCounterCheckpoint(baseline.Payload)
		if errors.Is(err, aggregation.ErrCounterCheckpointVersion) {
			return boundaryCounterBaselineVersion
		}
		if err != nil || checkpoint.Binding.MemberKey != baseline.MemberKey {
			return boundaryCounterBaselineInvalid
		}
		if _, ambiguous := checkpoints[baseline.MemberKey]; ambiguous || checkpoint.WindowEnd.After(restored.NextClose) {
			// Two baselines for one member, or one past the windows the checkpoint
			// closed, cannot both be true of a state this process committed.
			return boundaryCounterBaselineInvalid
		}
		checkpoints[baseline.MemberKey] = checkpoint
	}
	if err := assembler.RestoreCounterCheckpoints(checkpoints); err != nil {
		return boundaryCounterBaselineInvalid
	}
	return ""
}

// closureBaselines serializes the baselines a plan advances, in a stable order,
// to be committed with the closure that produced them.
func closureBaselines(plan snapshot.Plan) ([]groupdelivery.CounterBaseline, error) {
	checkpoints := plan.CounterCheckpoints()
	baselines := make([]groupdelivery.CounterBaseline, 0, len(checkpoints))
	for _, member := range slices.Sorted(maps.Keys(checkpoints)) {
		payload, err := aggregation.EncodeCounterCheckpoint(checkpoints[member])
		if err != nil {
			return nil, err
		}
		baselines = append(baselines, groupdelivery.CounterBaseline{MemberKey: member, Payload: payload})
	}
	return baselines, nil
}

// CounterSource is the confirmed measurement semantics of one cumulative
// member. It comes from the saved measurement definition and binding, never
// from a tag name or a unit.
type CounterSource struct {
	MeasurementID         string
	DefinitionRevision    string
	SourceBindingRevision string
	SeriesEpoch           string
	Unit                  string
	Kind                  aggregation.CounterKind
	// RolloverModulus and MaxRatePerSecond are the confirmed rollover evidence of
	// a monotonic counter; both are exact values, never floats.
	RolloverModulus  *measurement.ExactValue
	MaxRatePerSecond *measurement.ExactValue
}

// counterBindingFor binds one member to its confirmed counter semantics and
// returns the activation code that refuses it when it cannot be bound. A source
// that is not numeric, or whose measurement disagrees with the member's saved
// one, is never bridged by a guess.
func counterBindingFor(group *workspace.WriteGroup, member workspace.WriteGroupMember, key string, kind measurement.ExactType, sources map[string]CounterSource) (bound aggregation.CounterBinding, code string) {
	source, ok := sources[member.TagID]
	if !ok || !source.Kind.Valid() {
		return aggregation.CounterBinding{}, boundaryCounterSemantics
	}
	switch kind {
	case measurement.ExactInt64, measurement.ExactUint64, measurement.ExactFloat64:
	default:
		return aggregation.CounterBinding{}, boundaryCounterType
	}
	if member.MeasurementID != nil && *member.MeasurementID != source.MeasurementID {
		return aggregation.CounterBinding{}, boundaryCounterConflict
	}
	binding := aggregation.CounterBinding{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, ProfileRevision: group.Revision, MemberKey: key,
		MeasurementID: source.MeasurementID, DefinitionRevision: source.DefinitionRevision,
		SourceBindingRevision: source.SourceBindingRevision, SeriesEpoch: source.SeriesEpoch,
		Kind: source.Kind, Unit: source.Unit, RolloverModulus: source.RolloverModulus, MaxRatePerSecond: source.MaxRatePerSecond,
	}
	if err := binding.Validate(); err != nil {
		if errors.Is(err, aggregation.ErrCounterPolicyInvalid) {
			return aggregation.CounterBinding{}, boundaryCounterPolicy
		}
		return aggregation.CounterBinding{}, boundaryCounterSemantics
	}
	return binding, ""
}

// counterIdentityMatches reports whether an envelope that names a measurement
// series is the series the member was confirmed against. Fields the envelope
// leaves empty are not claims. A reading of another epoch, definition or source
// binding is not comparable and must not enter this baseline.
func counterIdentityMatches(counter *aggregation.CounterBinding, envelope measurement.SampleEnvelope) bool {
	for _, claim := range []struct{ carried, confirmed string }{
		{envelope.MeasurementID, counter.MeasurementID},
		{envelope.SeriesEpoch, counter.SeriesEpoch},
		{envelope.DefinitionRevision, counter.DefinitionRevision},
		{envelope.SourceBindingRevision, counter.SourceBindingRevision},
	} {
		if claim.carried != "" && claim.carried != claim.confirmed {
			return false
		}
	}
	return true
}
