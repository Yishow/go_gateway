package runtime

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type frozenFixtureLedger struct {
	*groupdelivery.Ledger
	store *groupdelivery.Store
}

// freezeDurableFixture supplies the descriptor that real Apply persists. It is
// only used by pre-existing algorithm fixtures; missing-proof tests never call
// it. Existing immutable descriptors are not replaced on restart.
func freezeDurableFixture(t *testing.T, f *longFixture) {
	t.Helper()
	if f.config.Ledger == nil || f.config.CounterBindings != nil {
		return
	}
	var ledger *frozenFixtureLedger
	switch l := f.config.Ledger.(type) {
	case *frozenFixtureLedger:
		ledger = l
	case duplicatingLedger:
		ledger, _ = l.inner.(*frozenFixtureLedger)
	case ledgerWithoutBaselines:
		ledger, _ = l.inner.(*frozenFixtureLedger)
	}
	if ledger == nil {
		return
	}
	proofs := workspace.CounterBindingProofs{}
	for _, member := range f.config.Group.Members {
		profile := f.config.Group.RecordingProfile
		if profile == nil || profile.Kind != workspace.WriteGroupProfileCumulative && member.ValueRole != workspace.WriteGroupValueRoleCounterRaw {
			continue
		}
		identity := workspace.CounterBindingMemberIdentity{WorkspaceID: f.config.Group.WorkspaceID, GroupID: f.config.Group.ID, AppliedRevision: f.config.Group.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
		if member.ExactRecordingSource != nil {
			proof, err := workspace.CounterBindingProofFromExactSource(identity, member, f.config.TagTypes[member.TagID])
			require.NoError(t, err)
			proofs[identity.Key()] = proof
			continue
		}
		source := f.config.Counters[member.TagID]
		kind, ok := measurement.ExactTypeForTag(f.config.TagTypes[member.TagID])
		require.True(t, ok)
		policy := measurement.CounterPolicy{AllowNegative: source.Kind == aggregation.CounterKindSigned}
		if source.RolloverModulus != nil {
			policy.RolloverModulus, _ = source.RolloverModulus.CanonicalDigits()
		}
		if source.MaxRatePerSecond != nil {
			policy.MaxRatePerSecond, _ = source.MaxRatePerSecond.CanonicalDigits()
		}
		proof, err := workspace.NewCounterBindingProof(workspace.CounterBindingFacts{
			Member: identity, MeasurementID: source.MeasurementID, DefinitionRevision: source.DefinitionRevision, SourceBindingRevision: source.SourceBindingRevision, SeriesEpoch: source.SeriesEpoch, Unit: source.Unit, SemanticKind: measurement.SemanticKind(source.Kind),
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, SourceValueType: f.config.TagTypes[member.TagID], ResultValueType: kind, ScaleMode: workspace.CounterScaleUnscaled, ExactScale: "1", CounterPolicy: policy,
		})
		require.NoError(t, err)
		proofs[identity.Key()] = proof
	}
	if len(proofs) == 0 {
		return
	}
	f.config.CounterBindings = proofs
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{WorkspaceID: f.config.Group.WorkspaceID, GroupID: f.config.Group.ID, AppliedRevision: f.config.Group.AppliedRevision, Group: f.config.Group, CounterBindings: proofs, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: f.config.TagTypes}})
	require.NoError(t, err)
	require.NoError(t, ledger.store.SaveRuntimeVersion(t.Context(), groupdelivery.GroupKey{WorkspaceID: f.config.Group.WorkspaceID, GroupID: f.config.Group.ID, GroupRevision: f.config.Group.AppliedRevision}, payload))
}
