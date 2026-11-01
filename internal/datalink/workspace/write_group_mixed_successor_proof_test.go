package workspace

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMixedSuccessorStatusBindsReservedCounterProof(t *testing.T) {
	identity, member := exactCounterBindingFixture()
	group, _ := mixedSchemaFixture()
	member.Calculation = WriteGroupMemberCalculationDelta
	member.Cadence = group.Members[0].Cadence
	member.ExactRecordingSource = nil
	group.WorkspaceID, group.ID, group.Revision = identity.WorkspaceID, identity.GroupID, identity.AppliedRevision
	group.Members = []WriteGroupMember{member}
	facts := counterProofFacts()
	facts.ScaleMode, facts.ExactScale, facts.ResultValueType = CounterScaleUnscaled, "1", measurement.ExactUint64
	proof, err := NewCounterBindingProof(facts)
	require.NoError(t, err)
	proofs := CounterBindingProofs{identity.Key(): proof}
	types := map[string]schema.DataType{member.TagID: schema.DataTypeUint64}
	managed, err := managedSourceDigest(group, types, proofs)
	require.NoError(t, err)
	result := &WriteGroupSuccessorPreparation{Group: group, LineageID: "lineage", PredecessorGroupID: "previous", OperationID: "operation", IntentDigest: "intent"}
	b := successorReserveSchemaBinding{Version: successorReserveSchemaBindingVersion, WorkspaceID: group.WorkspaceID, LineageID: result.LineageID, PredecessorGroupID: result.PredecessorGroupID, GroupID: group.ID, OperationID: result.OperationID, IntentDigest: result.IntentDigest, ManagedSourceDigest: managed, Dialect: "sqlite", SourceFamily: "sqlite", TagTypes: types}
	b.SourceDigest, err = b.sourceDigest()
	require.NoError(t, err)
	raw, err := json.Marshal(b)
	require.NoError(t, err)
	var fields map[string]any
	require.NoError(t, json.Unmarshal(raw, &fields))
	fields["counter_binding_proofs"] = proofs
	raw, err = json.Marshal(fields)
	require.NoError(t, err)
	_, err = decodeSuccessorReserveSchemaBinding(raw, result)
	require.NoError(t, err, "historical status must validate the original reserved-group proof without consulting mutable source facts")
	delete(fields, "counter_binding_proofs")
	raw, err = json.Marshal(fields)
	require.NoError(t, err)
	_, err = decodeSuccessorReserveSchemaBinding(raw, result)
	require.Error(t, err, "missing original counter proof cannot replay")
	facts.Member.GroupID = "previous"
	old, err := NewCounterBindingProof(facts)
	require.NoError(t, err)
	fields["counter_binding_proofs"] = CounterBindingProofs{old.Member.Key(): old}
	raw, err = json.Marshal(fields)
	require.NoError(t, err)
	_, err = decodeSuccessorReserveSchemaBinding(raw, result)
	require.Error(t, err, "predecessor proof must not become successor baseline")
}
