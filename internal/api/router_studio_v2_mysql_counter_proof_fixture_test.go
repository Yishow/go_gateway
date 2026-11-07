package api

import (
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// The hand-built fixture snapshots must carry the same immutable counter
// evidence as Apply. Missing proof must keep failing closed in production.
func mysqlCircuitCounterProofs(t *testing.T, group *workspace.WriteGroup, definition *measurement.MeasurementDefinition, types map[string]schema.DataType) workspace.CounterBindingProofs {
	t.Helper()
	proofs := workspace.CounterBindingProofs{}
	for _, member := range group.Members {
		identity := workspace.CounterBindingMemberIdentity{WorkspaceID: group.WorkspaceID, GroupID: group.ID, AppliedRevision: group.Revision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}
		var proof workspace.CounterBindingProof
		var err error
		if member.ExactRecordingSource != nil {
			proof, err = workspace.CounterBindingProofFromExactSource(identity, member, types[member.TagID])
		} else {
			result, supported := measurement.ExactTypeForTag(types[member.TagID])
			require.True(t, supported)
			facts := workspace.CounterBindingFacts{
				Member: identity, MeasurementID: definition.ID, DefinitionRevision: definition.DefinitionRevision,
				SourceBindingRevision: definition.SourceBindingRevision, SeriesEpoch: definition.SeriesEpoch,
				SemanticKind: definition.SemanticKind, Unit: definition.Unit, SourceRevision: member.SourceRevision,
				MappingRevision: member.MappingRevision, ScaleMode: workspace.CounterScaleUnscaled, ExactScale: "1",
				SourceValueType: types[member.TagID], ResultValueType: result,
			}
			if definition.CounterPolicy != nil {
				facts.CounterPolicy = *definition.CounterPolicy
			}
			proof, err = workspace.NewCounterBindingProof(facts)
		}
		require.NoError(t, err)
		proofs[identity.Key()] = proof
	}
	return proofs
}
