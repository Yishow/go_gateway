package runtime

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func counterClaimsFixture(t *testing.T, cumulative, scaled bool) (*longFixture, *groupdelivery.Store, groupdelivery.GroupKey, *GroupBoundary, time.Duration, workspace.CounterBindingProof) {
	t.Helper()
	f, store, key, _ := counterAcceptanceFixture(t)
	interval := counterWindow
	member := &f.config.Group.Members[0]
	if !cumulative {
		interval = 15 * time.Second
		f.config.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond}
		f.config.Group.RowPolicy.IntervalSeconds = 15
		member.ValueRole = workspace.WriteGroupValueRoleCounterRaw
		member.SourceValueType, member.ResultValueType = schema.DataTypeInt64, measurement.ExactInt64
	}
	var proof workspace.CounterBindingProof
	for id, original := range f.config.CounterBindings {
		proof = original
		if scaled {
			member.SourceValueType, member.ResultValueType = schema.DataTypeInt64, measurement.ExactDecimal
			member.ExactRecordingSource = &workspace.ExactRecordingSource{
				MeasurementID: "meas-int", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1",
				Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, OriginalDataType: schema.DataTypeInt64, Scale: "0.1",
			}
			var err error
			proof, err = workspace.CounterBindingProofFromExactSource(original.Member, *member, schema.DataTypeInt64)
			require.NoError(t, err)
			f.config.CounterBindings[id] = proof
		}
	}
	payload, err := json.Marshal(workspace.WriteGroupAppliedSnapshot{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, AppliedRevision: key.GroupRevision, Group: f.config.Group, RuntimeLayout: &workspace.WriteGroupRuntimeLayout{TagTypes: f.config.TagTypes}, CounterBindings: f.config.CounterBindings})
	require.NoError(t, err)
	require.NoError(t, store.SaveRuntimeVersion(t.Context(), key, payload))
	boundary := f.build(t)
	require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(0, "baseline-1", at(time.Second), int64(120000))))
	require.NoError(t, boundary.Tick(t.Context(), at(interval)))
	f.clock = at(2 * interval)
	return f, store, key, boundary, interval, proof
}

func matchingCounterClaims(e *measurement.SampleEnvelope) {
	e.MeasurementID, e.DefinitionRevision, e.SourceBindingRevision, e.SeriesEpoch = "meas-int", "def-1", "bind-1", "epoch-1"
}

func changeCounterClaim(e *measurement.SampleEnvelope, field, value string) {
	switch field {
	case "measurement_id":
		e.MeasurementID = value
	case "definition_revision":
		e.DefinitionRevision = value
	case "source_binding_revision":
		e.SourceBindingRevision = value
	case "series_epoch":
		e.SeriesEpoch = value
	}
}
