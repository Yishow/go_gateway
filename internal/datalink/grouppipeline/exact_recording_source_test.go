package grouppipeline

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type unavailableExactCounterReader struct{ calls int }

func (r *unavailableExactCounterReader) CounterSources(context.Context, *workspace.WriteGroup) (map[string]runtime.CounterSource, error) {
	r.calls++
	return nil, errors.New("editable definition repository unavailable")
}

func TestFrozenExactRecordingSourceDoesNotRefetchDefinition(t *testing.T) {
	for _, tc := range []struct {
		name        string
		withBinding bool
		tagType     schema.DataType
	}{
		{"legacy_exact_source", false, schema.DataTypeInt64},
		{"versioned_int64_binding", true, schema.DataTypeInt64},
		{"versioned_uint64_binding", true, schema.DataTypeUint64},
	} {
		t.Run(tc.name, func(t *testing.T) { testFrozenExactRecordingSource(t, tc.withBinding, tc.tagType) })
	}
}

func testFrozenExactRecordingSource(t *testing.T, withBinding bool, tagType schema.DataType) {
	t.Helper()
	fixture := newFixture(t)
	pipeline := fixture.pipeline()
	snap := mysqlRuntimeSnapshot(t)
	member := &snap.Group.Members[0]
	measurementID := "scaled-measurement"
	member.MeasurementID = &measurementID
	member.SourceRevision, member.MappingRevision = "source-1", "mapping-1"
	member.SourceValueType, member.ResultValueType = tagType, measurement.ExactDecimal
	member.ValueRole = workspace.WriteGroupValueRoleCounterRaw
	member.ExactRecordingSource = &workspace.ExactRecordingSource{
		MeasurementID: measurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", OriginalDataType: tagType, Scale: "0.1", Unit: "kWh",
		SemanticKind: measurement.SemanticKindCounter,
	}
	snap.RuntimeLayout.TagTypes[member.TagID] = tagType
	if withBinding {
		proof, err := workspace.CounterBindingProofFromExactSource(workspace.CounterBindingMemberIdentity{
			WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, AppliedRevision: snap.AppliedRevision,
			DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
		}, *member, tagType)
		require.NoError(t, err)
		snap.CounterBindings = workspace.CounterBindingProofs{proof.Member.Key(): proof}
	}
	reader := &unavailableExactCounterReader{}
	pipeline.deps.Counters = reader
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	frozen, err := pipeline.loadRuntimeSnapshot(t.Context(), groupKey(snap))
	require.NoError(t, err)
	built, reason := pipeline.buildFrozen(t.Context(), frozen, at(0), time.Time{})
	require.Empty(t, reason)
	require.NotNil(t, built)
	require.Zero(t, reader.calls, "accepted source facts survive a definition outage or edit")
	for i, input := range []struct {
		raw  uint64
		want string
	}{{120000, "12000.0"}, {120004, "12000.4"}, {9007199254740993, "900719925474099.3"}} {
		var raw any = input.raw
		if tagType == schema.DataTypeInt64 {
			raw = int64(input.raw)
		}
		require.NoError(t, built.boundary.AcceptSample(t.Context(), envelope(frozen.Group, input.want, at(i+1), raw)))
		var payload string
		require.NoError(t, fixture.store.DB().QueryRowContext(t.Context(),
			`SELECT payload FROM wg_delivery_samples ORDER BY seq DESC LIMIT 1`).Scan(&payload))
		var sample struct {
			Value measurement.ExactValue `json:"value"`
		}
		require.NoError(t, json.Unmarshal([]byte(payload), &sample))
		digits, ok := sample.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, input.want, digits, "exact integer scale runs once, with no float64 loss")
	}
	require.Zero(t, reader.calls)
}
