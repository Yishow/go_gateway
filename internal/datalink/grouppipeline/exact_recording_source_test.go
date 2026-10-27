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
	fixture := newFixture(t)
	pipeline := fixture.pipeline()
	snap := mysqlRuntimeSnapshot(t)
	member := &snap.Group.Members[0]
	measurementID := "scaled-measurement"
	member.MeasurementID = &measurementID
	member.SourceRevision, member.MappingRevision = "source-1", "mapping-1"
	member.SourceValueType, member.ResultValueType = schema.DataTypeInt64, measurement.ExactDecimal
	member.ValueRole = workspace.WriteGroupValueRoleCounterRaw
	member.ExactRecordingSource = &workspace.ExactRecordingSource{
		MeasurementID: measurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", OriginalDataType: schema.DataTypeInt64, Scale: "0.1", Unit: "kWh",
		SemanticKind: measurement.SemanticKindCounter,
	}
	snap.RuntimeLayout.TagTypes[member.TagID] = schema.DataTypeInt64
	reader := &unavailableExactCounterReader{}
	pipeline.deps.Counters = reader
	require.NoError(t, pipeline.saveRuntimeSnapshot(t.Context(), snap))
	frozen, err := pipeline.loadRuntimeSnapshot(t.Context(), groupKey(snap))
	require.NoError(t, err)
	built, reason := pipeline.buildFrozen(t.Context(), frozen, at(0), time.Time{})
	require.Empty(t, reason)
	require.NotNil(t, built)
	require.Zero(t, reader.calls, "accepted source facts survive a definition outage or edit")
	require.NoError(t, built.boundary.AcceptSample(t.Context(), envelope(frozen.Group, "scaled", at(1), int64(120004))))
	var payload string
	require.NoError(t, fixture.store.DB().QueryRowContext(t.Context(),
		`SELECT payload FROM wg_delivery_samples`).Scan(&payload))
	var sample struct {
		Value measurement.ExactValue `json:"value"`
	}
	require.NoError(t, json.Unmarshal([]byte(payload), &sample))
	digits, ok := sample.Value.DecimalDigits()
	require.True(t, ok)
	require.Equal(t, "12000.4", digits)
}
