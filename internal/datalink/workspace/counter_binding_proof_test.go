package workspace

import (
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// The v1 bytes are independent contract fixtures, not encoder-generated expectations.
const counterProofBody = `{"version":1,"member":{"workspace_id":"workspace-1","group_id":"group-1","applied_revision":"applied-1","device_id":"device-1","point_id":"point-1","tag_id":"tag-1"},"measurement_id":"measurement-1","definition_revision":"def-1","source_binding_revision":"binding-1","series_epoch":"epoch-1","semantic_kind":"counter","unit":"kWh","source_revision":"source-1","mapping_revision":"mapping-1","scale_mode":"exact","exact_scale":"0.1","source_value_type":"uint64","result_value_type":"decimal","counter_policy":{"allow_negative":false,"rollover_modulus":"9007199254740993","max_rate_per_second":"0.125"}}`
const counterProofDigest = "901aaf8cbc6dfe0778514cadf5a648961d0d324362b6266bc704123d802eeee4"

var counterProofJSON = counterProofBody[:len(counterProofBody)-1] + `,"digest":"` + counterProofDigest + `"}`

const counterMemberKey = `["workspace-1","group-1","applied-1","device-1","point-1","tag-1"]`
const legacyNonCounterSnapshot = `{"workspace_id":"workspace-1","group_id":"group-1","group_revision":"revision-1","applied_revision":"applied-1","effective_at":"2026-10-01T00:00:00Z","group":null}`

func counterProofFacts() CounterBindingFacts {
	return CounterBindingFacts{
		Member:        CounterBindingMemberIdentity{WorkspaceID: "workspace-1", GroupID: "group-1", AppliedRevision: "applied-1", DeviceID: "device-1", PointID: "point-1", TagID: "tag-1"},
		MeasurementID: "measurement-1", DefinitionRevision: "def-1", SourceBindingRevision: "binding-1",
		SeriesEpoch: "epoch-1", SemanticKind: measurement.SemanticKindCounter, Unit: "kWh",
		SourceRevision: "source-1", MappingRevision: "mapping-1",
		ScaleMode: CounterScaleExact, ExactScale: "0.1", SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactDecimal,
		CounterPolicy: measurement.CounterPolicy{RolloverModulus: "9007199254740993", MaxRatePerSecond: "0.125"},
	}
}

func TestCounterBindingProofCanonicalRoundTrip(t *testing.T) {
	proof, err := NewCounterBindingProof(counterProofFacts())
	require.NoError(t, err)
	require.Equal(t, 1, proof.Version)
	require.Equal(t, counterProofDigest, proof.Digest)
	wire, err := json.Marshal(proof)
	require.NoError(t, err)
	require.Equal(t, counterProofJSON, string(wire))
	var decoded CounterBindingProof
	require.NoError(t, json.Unmarshal([]byte(counterProofJSON), &decoded))
	require.Equal(t, counterProofFacts(), decoded.CounterBindingFacts)
	require.Equal(t, counterProofDigest, decoded.Digest)
	// Whitespace and object order are not semantic changes.
	reordered := `{"digest":"` + counterProofDigest + `",` + counterProofBody[1:]
	require.NoError(t, json.Unmarshal([]byte("\n"+reordered+"\n"), &decoded))
	canonical, err := json.Marshal(decoded)
	require.NoError(t, err)
	require.Equal(t, counterProofJSON, string(canonical))
}

func TestCounterBindingProofRejectsTamperWithoutChangingReceiver(t *testing.T) {
	for _, field := range []string{"workspace-1", "group-1", "applied-1", "device-1", "point-1", "tag-1",
		"measurement-1", "def-1", "binding-1", "epoch-1", "counter", "kWh", "source-1", "mapping-1",
		"exact", "0.1", "uint64", "decimal", "9007199254740993", "0.125", counterProofDigest} {
		t.Run(field, func(t *testing.T) {
			var decoded CounterBindingProof
			require.NoError(t, json.Unmarshal([]byte(counterProofJSON), &decoded))
			before := decoded
			tampered := strings.Replace(counterProofJSON, `"`+field+`"`, `"`+field+`-changed"`, 1)
			if field == "kWh" {
				tampered = strings.Replace(counterProofJSON, `"kWh"`, `"MWh"`, 1)
			}
			err := json.Unmarshal([]byte(tampered), &decoded)
			require.ErrorIs(t, err, ErrCounterBindingProofConflict)
			var typed *CounterBindingProofError
			require.ErrorAs(t, err, &typed)
			require.Equal(t, before, decoded, "a rejected decode cannot partially replace trusted facts")
		})
	}
	proof, err := NewCounterBindingProof(counterProofFacts())
	require.NoError(t, err)
	proof.Unit = "MWh"
	_, err = json.Marshal(proof)
	require.ErrorIs(t, err, ErrCounterBindingProofConflict, "encoding never silently reseals edited facts")
}

func TestCounterBindingProofUnsupportedVersionIsTyped(t *testing.T) {
	for _, wire := range []string{
		strings.Replace(counterProofJSON, `"version":1`, `"version":2`, 1),
		`{"version":99,"future_semantics":true}`,
		`{"version":99,"future_semantics":[null,{"future":true}]}`,
		strings.Replace(counterProofJSON, `"version":1,`, "", 1),
	} {
		var decoded CounterBindingProof
		err := json.Unmarshal([]byte(wire), &decoded)
		require.ErrorIs(t, err, ErrCounterBindingProofUnsupported)
		var typed *CounterBindingProofError
		require.ErrorAs(t, err, &typed)
		require.Zero(t, decoded)
	}
}

func TestCounterBindingProofFullMemberIdentityNeverUsesTagAlone(t *testing.T) {
	base := counterProofFacts().Member
	require.Equal(t, counterMemberKey, base.Key())
	keys := map[string]bool{base.Key(): true}
	for i := range 6 {
		member := base
		fields := []*string{&member.WorkspaceID, &member.GroupID, &member.AppliedRevision, &member.DeviceID, &member.PointID, &member.TagID}
		*fields[i] += "-other"
		require.False(t, keys[member.Key()], "each identity component partitions the proof")
		keys[member.Key()] = true
	}
	left, right := base, base
	left.DeviceID, left.PointID = "device|point", `tag\"`
	right.DeviceID, right.PointID = "device", `point|tag\"`
	require.NotEqual(t, left.Key(), right.Key(), "opaque IDs cannot collide at a delimiter")
	var parts []string
	require.NoError(t, json.Unmarshal([]byte(left.Key()), &parts))
	require.Equal(t, []string{"workspace-1", "group-1", "applied-1", "device|point", `tag\"`, "tag-1"}, parts)
}

func TestCounterBindingProofSnapshotMapRoundTripAndKeyConflict(t *testing.T) {
	proof, err := NewCounterBindingProof(counterProofFacts())
	require.NoError(t, err)
	var snapshot WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(legacyNonCounterSnapshot), &snapshot))
	snapshot.CounterBindings = CounterBindingProofs{counterMemberKey: proof}
	wire, err := json.Marshal(snapshot)
	require.NoError(t, err)
	var decoded WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(wire, &decoded))
	require.Equal(t, snapshot, decoded)
	require.Contains(t, string(wire), `"counter_bindings":`)
	wrongKey := CounterBindingProofs{"tag-1": proof}
	_, err = json.Marshal(wrongKey)
	require.ErrorIs(t, err, ErrCounterBindingProofConflict)
	var proofs CounterBindingProofs
	err = json.Unmarshal([]byte(`{"tag-1":`+counterProofJSON+`}`), &proofs)
	require.ErrorIs(t, err, ErrCounterBindingProofConflict)
	require.Nil(t, proofs)
	// The map must preserve an unsupported version's classification, even if
	// that future format has a different shape from v1.
	err = json.Unmarshal([]byte(`{"future":{"version":99,"future_semantics":[null]}}`), &proofs)
	require.ErrorIs(t, err, ErrCounterBindingProofUnsupported)
	require.Nil(t, proofs)
}

func TestCounterBindingProofLegacyNonCounterSnapshotBytesUnchanged(t *testing.T) {
	var snapshot WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(legacyNonCounterSnapshot), &snapshot))
	require.Nil(t, snapshot.CounterBindings)
	wire, err := json.Marshal(snapshot)
	require.NoError(t, err)
	require.Equal(t, legacyNonCounterSnapshot, string(wire), "absence remains absence, not a fabricated counter proof")
}

func TestCounterBindingProofConstructorCopiesPolicyPointers(t *testing.T) {
	facts := counterProofFacts()
	reset, legacyRollover := 0.5, 65536.0
	facts.CounterPolicy.ResetThreshold, facts.CounterPolicy.MaxRollOver = &reset, &legacyRollover
	proof, err := NewCounterBindingProof(facts)
	require.NoError(t, err)
	reset, legacyRollover = 99, 99
	require.Equal(t, 0.5, *proof.CounterPolicy.ResetThreshold)
	require.Equal(t, 65536.0, *proof.CounterPolicy.MaxRollOver)
	wire, err := json.Marshal(proof)
	require.NoError(t, err)
	var decoded CounterBindingProof
	require.NoError(t, json.Unmarshal(wire, &decoded))
	require.Equal(t, proof, decoded)
	*decoded.CounterPolicy.ResetThreshold = 1
	require.Equal(t, 0.5, *proof.CounterPolicy.ResetThreshold)
	_, err = json.Marshal(decoded)
	require.ErrorIs(t, err, ErrCounterBindingProofConflict)
}

func TestCounterBindingProofPreservesUnicodeValues(t *testing.T) {
	facts := counterProofFacts()
	facts.Member.DeviceID = "設備 ſ"
	facts.Unit = "度"
	proof, err := NewCounterBindingProof(facts)
	require.NoError(t, err)
	wire, err := json.Marshal(proof)
	require.NoError(t, err)
	var decoded CounterBindingProof
	require.NoError(t, json.Unmarshal(wire, &decoded))
	require.Equal(t, "設備 ſ", decoded.Member.DeviceID)
	require.Equal(t, "度", decoded.Unit)
	require.Equal(t, `["workspace-1","group-1","applied-1","設備 ſ","point-1","tag-1"]`, decoded.Member.Key())
}
