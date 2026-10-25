package runtime

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// wideCounterFixture is an applied cumulative managed_wide_profile_v1 group over
// the same real SQLite table: two exact counters, each confirmed by its own saved
// measurement, one row per 15s window.
func wideCounterFixture(t *testing.T) *longFixture {
	t.Helper()
	f := wideProfileFixture(t)
	f.members = []workspace.WriteGroupMember{
		{
			DeviceID: "device-A", PointID: "point-G", TagID: "tag-G", SourceRevision: "src-G", MappingRevision: "map-G", TargetColumn: "v_gauge",
			MeasurementID: ptr("measurement-g"), SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactUint64,
		},
		{
			DeviceID: "device-A", PointID: "point-C", TagID: "tag-C", SourceRevision: "src-C", MappingRevision: "map-C", TargetColumn: "v_raw",
			MeasurementID: ptr("measurement-count"), SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactUint64,
		},
	}
	f.config.Group.Members = f.members
	f.config.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
	f.config.Group.RecordingProfile.SampleReduction = ""
	f.config.TagTypes["tag-G"] = schema.DataTypeUint64
	source := func(id string) CounterSource {
		return CounterSource{
			MeasurementID: id, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1",
			SeriesEpoch: "epoch-1", Unit: "kWh", Kind: aggregation.CounterKindMonotonic,
		}
	}
	f.config.Counters = map[string]CounterSource{"tag-G": source("measurement-g"), "tag-C": source("measurement-count")}
	return f
}

type wideCounterView struct {
	Status  string `json:"status"`
	Counter *struct {
		Status    string `json:"status"`
		Delta     string `json:"delta"`
		RawText   string `json:"raw_text"`
		DeltaFrom string `json:"delta_from"`
	} `json:"counter"`
}

func provenanceViews(t *testing.T, cells map[string]any) []wideCounterView {
	t.Helper()
	var views []wideCounterView
	text, ok := cells["provenance"].(string)
	require.True(t, ok, "provenance is stored as JSON text on SQLite")
	require.NoError(t, json.Unmarshal([]byte(text), &views), text)
	return views
}

func TestWideCumulativeBoundaryWritesOneRowWithEveryCounter(t *testing.T) {
	f := wideCounterFixture(t)
	boundary := f.build(t)
	f.accept(t, boundary, f.envelope(0, "g0", at(time.Second), uint64(100)), f.envelope(1, "c0", at(2*time.Second), uint64(10)))
	require.NoError(t, boundary.Tick(t.Context(), at(15*time.Second)))
	f.accept(t, boundary, f.envelope(0, "g1", at(16*time.Second), uint64(101)), f.envelope(1, "c1", at(17*time.Second), uint64(14)))
	require.NoError(t, boundary.Tick(t.Context(), at(30*time.Second)))

	require.Len(t, f.sink.rows, 2, "one record per window, not one per member")
	first, second := cellsOf(f.sink.rows[0]), cellsOf(f.sink.rows[1])
	require.Equal(t, [2]any{"100", "10"}, [2]any{first["v_gauge"], first["v_raw"]}, "the columns keep the raw readings")
	require.Equal(t, [2]any{"101", "14"}, [2]any{second["v_gauge"], second["v_raw"]})
	baseline := provenanceViews(t, first)
	require.Equal(t, [2]string{"baseline", "baseline"}, [2]string{baseline[0].Counter.Status, baseline[1].Counter.Status})
	require.Empty(t, baseline[0].Counter.Delta+baseline[1].Counter.Delta, "a baseline has no difference")
	next := provenanceViews(t, second)
	require.Equal(t, [2]string{"1", "4"}, [2]string{next[0].Counter.Delta, next[1].Counter.Delta}, "each member's own exact difference")
}

// TestWideCumulativeBoundaryFailsClosed pins every activation refusal of the
// cumulative wide shape: nothing becomes a writer on an id alone, a mean, a role,
// a contradicting type or a counter the provider does not confirm.
func TestWideCumulativeBoundaryFailsClosed(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(*longFixture)
		code   string
	}{
		"no counter provider output": {func(f *longFixture) { f.config.Counters = nil }, boundaryCounterSemantics},
		"a member the provider does not confirm": {func(f *longFixture) {
			delete(f.config.Counters, "tag-C")
		}, boundaryCounterSemantics},
		"a source of another measurement": {func(f *longFixture) {
			source := f.config.Counters["tag-G"]
			source.MeasurementID = "measurement-other"
			f.config.Counters["tag-G"] = source
		}, boundaryCounterConflict},
		"a mean as a counter baseline": {func(f *longFixture) {
			f.config.Group.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
		}, boundaryMemberMetadata},
		"an instantaneous role on a cumulative member": {func(f *longFixture) {
			f.members[0].ValueRole = workspace.WriteGroupValueRoleGauge
		}, boundaryMemberMetadata},
		"a result type that is not the raw type": {func(f *longFixture) {
			f.members[0].ResultValueType = measurement.ExactDecimal
		}, boundaryMemberMetadata},
		"an omitted source type": {func(f *longFixture) { f.members[0].SourceValueType = "" }, boundaryMemberMetadata},
		"a member without a saved measurement": {func(f *longFixture) {
			f.members[0].MeasurementID = nil
		}, boundaryProfileUnsupported},
		"a text tag": {func(f *longFixture) {
			f.config.TagTypes["tag-G"] = schema.DataTypeString
			f.members[0].SourceValueType, f.members[0].ResultValueType = schema.DataTypeString, measurement.ExactText
		}, boundaryCounterType},
		"a missing tag type": {func(f *longFixture) { delete(f.config.TagTypes, "tag-G") }, boundaryTagTypeUnsupport},
	} {
		t.Run(name, func(t *testing.T) {
			f := wideCounterFixture(t)
			tc.mutate(f)
			f.config.Group.Members = f.members
			require.Equal(t, tc.code, blockedCode(t, f.config))
		})
	}
}
