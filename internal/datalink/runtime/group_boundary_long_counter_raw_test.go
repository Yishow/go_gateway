package runtime

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// counterRawFixture is a mean long group of one numeric gauge and one counter_raw
// member whose saved measurement is confirmed.
func counterRawFixture(t *testing.T) *longFixture {
	t.Helper()
	f := newLongFixture(t)
	f.members = []workspace.WriteGroupMember{
		{
			DeviceID: "device-B", PointID: "point-B", TagID: "tag-B", SourceRevision: "src-B", MappingRevision: "map-B", TargetColumn: "v_b", Required: true,
			ValueRole: workspace.WriteGroupValueRoleGauge, SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactDecimal,
		},
		{
			DeviceID: "device-B", PointID: "point-B3", TagID: "tag-count", SourceRevision: "src-count", MappingRevision: "map-count", TargetColumn: "v_count", Required: true,
			MeasurementID: ptr("measurement-count"),
			ValueRole:     workspace.WriteGroupValueRoleCounterRaw, SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactUint64,
		},
	}
	f.config.Group.Members = f.members
	f.config.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
	f.config.Counters = map[string]CounterSource{"tag-count": {
		MeasurementID: "measurement-count", DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
		SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
	}}
	return f
}

func TestMeanLongCounterRawKeepsTheLatestConfirmedRawValue(t *testing.T) {
	f := counterRawFixture(t)
	boundary := f.build(t)
	claimed := f.envelope(1, "raw-claim", at(6*time.Second), uint64(31))
	claimed.MeasurementID, claimed.SeriesEpoch, claimed.DefinitionRevision = "measurement-count", "epoch-1", "def-1"
	f.accept(t, boundary,
		f.envelope(0, "b-1", at(time.Second), int64(1)), f.envelope(0, "b-2", at(2*time.Second), int64(2)),
		f.envelope(1, "raw-1", at(time.Second), uint64(10)), f.envelope(1, "raw-2", at(5*time.Second), uint64(30)), claimed,
	)
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	rows := longRows(t, f.sink)
	require.Len(t, rows, 2)

	gauge := rows[longRowKey{"tag-B", bucketText(0)}]
	require.Equal(t, [2]any{"1.5", "decimal"}, [2]any{gauge["value_decimal"], gauge["value_type"]})
	raw := rows[longRowKey{"tag-count", bucketText(0)}]
	require.Equal(t, "31", raw["value_decimal"], "the newest confirmed reading is kept, never a mean of the three")
	require.Equal(t, "measurement-count", raw["measurement_id"])
	require.Equal(t, []string{"value_decimal"}, valueSlots(raw))
	require.Nil(t, raw["delta_value"])
	require.Nil(t, raw["delta_from"])
	require.Nil(t, raw["delta_to"])
	require.Nil(t, raw["series_epoch"], "no counter window is produced for a latest raw value")
	var provenance map[string]any
	require.NoError(t, json.Unmarshal([]byte(provenanceText(t, raw)), &provenance))
	require.Equal(t, [3]any{"counter_raw", "uint64", "uint64"}, [3]any{provenance["role"], provenance["source_type"], provenance["result_type"]})
	require.NotContains(t, provenance, "counter", "no baseline or delta document")
	require.NotContains(t, provenance, "mean", "a latest member carries no mean block")
}

func TestMeanLongCounterRawRefusesAnotherCounterSeries(t *testing.T) {
	f := counterRawFixture(t)
	boundary := f.build(t)
	for name, claim := range map[string]func(*measurement.SampleEnvelope){
		"another measurement":         func(e *measurement.SampleEnvelope) { e.MeasurementID = "measurement-other" },
		"another series epoch":        func(e *measurement.SampleEnvelope) { e.SeriesEpoch = "epoch-2" },
		"another definition revision": func(e *measurement.SampleEnvelope) { e.DefinitionRevision = "def-9" },
		"another source binding":      func(e *measurement.SampleEnvelope) { e.SourceBindingRevision = "bind-9" },
	} {
		envelope := f.envelope(1, "raw-"+name, at(3*time.Second), uint64(99))
		claim(&envelope)
		var refused *GroupSampleError
		require.ErrorAs(t, boundary.AcceptSample(t.Context(), envelope), &refused, name)
		require.Equal(t, [2]string{string(snapshot.OfferRejected), snapshot.ReasonRevisionMismatch}, [2]string{string(refused.Outcome), refused.Reason}, name)
	}
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	raw := longRows(t, f.sink)[longRowKey{"tag-count", bucketText(0)}]
	require.Equal(t, [2]any{"missing", nil}, [2]any{raw["status"], raw["value_decimal"]}, "a refused reading never becomes the record")
}

// TestMeanLongCounterRawNeedsAConfirmedBinding pins the activation refusals: a
// counter_raw member never runs on an id alone, a missing provider, a source of
// another measurement, a non-numeric tag or partial frozen facts.
func TestMeanLongCounterRawNeedsAConfirmedBinding(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(*longFixture)
		code   string
	}{
		"no counter provider output": {func(f *longFixture) { f.config.Counters = nil }, boundaryCounterSemantics},
		"the member is unconfirmed":  {func(f *longFixture) { delete(f.config.Counters, "tag-count") }, boundaryCounterSemantics},
		"an invalid counter kind": {func(f *longFixture) {
			source := f.config.Counters["tag-count"]
			source.Kind = ""
			f.config.Counters["tag-count"] = source
		}, boundaryCounterSemantics},
		"a source of another measurement": {func(f *longFixture) {
			source := f.config.Counters["tag-count"]
			source.MeasurementID = "measurement-other"
			f.config.Counters["tag-count"] = source
		}, boundaryCounterConflict},
		"a text tag": {func(f *longFixture) {
			f.config.TagTypes["tag-count"] = schema.DataTypeString
			f.members[1].SourceValueType, f.members[1].ResultValueType = schema.DataTypeString, measurement.ExactText
		}, boundaryCounterType},
		"no measurement id":    {func(f *longFixture) { f.members[1].MeasurementID = nil }, boundaryMemberMetadata},
		"an empty id":          {func(f *longFixture) { f.members[1].MeasurementID = ptr("") }, boundaryMemberMetadata},
		"an omitted source":    {func(f *longFixture) { f.members[1].SourceValueType = "" }, boundaryMemberMetadata},
		"a contradicting type": {func(f *longFixture) { f.members[1].ResultValueType = measurement.ExactDecimal }, boundaryMemberMetadata},
		"an unknown role":      {func(f *longFixture) { f.members[1].ValueRole = "average" }, boundaryMemberMetadata},
	} {
		t.Run(name, func(t *testing.T) {
			f := counterRawFixture(t)
			tc.mutate(f)
			f.config.Group.Members = f.members
			_, err := NewGroupBoundary(f.config)
			var blocked *GroupBoundaryError
			require.True(t, errors.As(err, &blocked), "%v", err)
			require.Equal(t, tc.code, blocked.Code)
		})
	}
}
