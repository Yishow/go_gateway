package workspace

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedSchemaCounterProofUsesExistingScopedFacts(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		t.Run(map[bool]string{false: "unscaled", true: "scaled"}[scaled], func(t *testing.T) {
			identity, member := exactCounterBindingFixture()
			member.Calculation = WriteGroupMemberCalculationDelta
			group, _ := mixedSchemaFixture()
			member.Cadence = group.Members[0].Cadence
			group.WorkspaceID = identity.WorkspaceID
			group.ID = identity.GroupID
			group.Revision = identity.AppliedRevision
			var proof CounterBindingProof
			var err error
			if scaled {
				proof, err = CounterBindingProofFromExactSource(identity, member, schema.DataTypeUint64)
			} else {
				member.ExactRecordingSource = nil
				facts := counterProofFacts()
				facts.ScaleMode = CounterScaleUnscaled
				facts.ExactScale = "1"
				facts.ResultValueType = measurement.ExactUint64
				proof, err = NewCounterBindingProof(facts)
			}
			require.NoError(t, err)
			group.Members = []WriteGroupMember{member}
			types := map[string]schema.DataType{member.TagID: schema.DataTypeUint64}
			proofs := CounterBindingProofs{identity.Key(): proof}
			digest, err := managedSourceDigest(group, types, proofs)
			require.NoError(t, err)
			require.NotEmpty(t, digest)
			_, err = managedSourceDigest(group, types)
			require.ErrorIs(t, err, ErrCounterBindingProofMissing)
			tampered := proof
			tampered.DefinitionRevision = "changed-without-seal"
			_, err = managedSourceDigest(group, types, CounterBindingProofs{identity.Key(): tampered})
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			foreign := proof.CounterBindingFacts
			foreign.Member.GroupID = "other-record"
			foreignProof, err := NewCounterBindingProof(foreign)
			require.NoError(t, err)
			_, err = managedSourceDigest(group, types, CounterBindingProofs{identity.Key(): foreignProof})
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			facts := proof.CounterBindingFacts
			facts.CounterPolicy.MaxRatePerSecond = "0.25"
			facts.DefinitionRevision = "def-2"
			refreshed, err := NewCounterBindingProof(facts)
			require.NoError(t, err)
			next, err := managedSourceDigest(group, types, CounterBindingProofs{identity.Key(): refreshed})
			if scaled {
				require.ErrorIs(t, err, ErrCounterBindingProofConflict, "scaled persisted source must be refreshed, never silently relabelled")
			} else {
				require.NoError(t, err)
				require.NotEqual(t, digest, next)
			}
		})
	}
}
func TestMixedLatestAndMeanDoNotRequireOrBorrowCounterProof(t *testing.T) {
	for _, method := range []WriteGroupMemberCalculation{WriteGroupMemberCalculationLatest, WriteGroupMemberCalculationMean} {
		group, types := mixedSchemaFixture()
		group.Members[0].Calculation = method
		_, err := managedSourceDigest(group, types)
		require.NoError(t, err)
		proof, err := NewCounterBindingProof(counterProofFacts())
		require.NoError(t, err)
		_, err = managedSourceDigest(group, types, CounterBindingProofs{proof.Member.Key(): proof})
		require.ErrorIs(t, err, ErrCounterBindingProofConflict)
	}
}
