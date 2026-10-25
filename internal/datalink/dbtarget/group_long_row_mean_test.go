package dbtarget

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var meanWindowStart = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func meanWindowOutcome(t *testing.T, member LongRowMember, result snapshot.MemberResult) snapshot.Outcome {
	t.Helper()
	windowEnd := meanWindowStart.Add(15 * time.Second)
	if result.Sample != nil {
		result.Sample.MemberKey = member.MemberKey
	}
	return snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: member.MemberKey, RecordID: "rec-mean", EffectKey: "eff-mean",
		BucketStart: meanWindowStart, BucketEnd: windowEnd, Members: []snapshot.MemberResult{result},
	}
}

func meanOutcome(t *testing.T, member LongRowMember) snapshot.Outcome {
	t.Helper()
	decimal, err := measurement.NewDecimal("20")
	require.NoError(t, err)
	sum, err := measurement.NewDecimal("60")
	require.NoError(t, err)
	sampleTime := meanWindowStart.Add(9 * time.Second)
	return meanWindowOutcome(t, member, snapshot.MemberResult{
		MemberKey: member.MemberKey, Status: snapshot.MemberOK,
		Sample: &snapshot.Sample{SampleID: "sample-9", ObservedAt: sampleTime, Quality: schema.QualityGood, Value: measurement.NewInt64(30)},
		Reduction: &measurement.MeanResult{
			Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod,
			Value: decimal, Sum: sum, Count: 3, Basis: measurement.MeanBasisExact,
			Scale: 0, Rounding: measurement.MeanRoundingHalfEven,
		},
		Evidence: &snapshot.MeanEvidence{
			Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod,
			Basis: measurement.MeanBasisExact, SourceValueType: schema.DataTypeInt64, ResultValueType: measurement.ExactDecimal,
			ValidCount: 3, RejectedCount: 1,
			FirstObservedAt: meanWindowStart.Add(3 * time.Second), LastObservedAt: meanWindowStart.Add(12 * time.Second),
			Coverage: snapshot.CoveragePartial, Reason: "bad", RejectionReasons: []string{"bad"},
		},
	})
}

// TestMeanLongRowEncoding pins the typed mean encoder: the reduction value and
// its result type fill value_decimal/value_type, the raw sample value never
// becomes the window value, and provenance carries the frozen source/result
// contract and the window evidence.
func TestMeanLongRowEncoding(t *testing.T) {
	member := LongRowMember{
		MemberKey: "device-1/pt-1/tag-1", DeviceID: "device-1", PointID: "pt-1", TagID: "tag-1",
		Type: measurement.ExactDecimal, Reduction: "mean", ValueRole: "gauge", SourceValueType: schema.DataTypeInt64,
	}
	spec := LongRowSpec{
		Dialect: SQLDialectSQLite, Columns: longColumnsFor(t, "sqlite"), WorkspaceID: "ws", GroupID: "g1", GroupRevision: "r1",
		ProfileKind: "instantaneous", ProfileInterval: 15, ProfileUnit: "second", Members: []LongRowMember{member},
	}
	layout, issues := NewLongRowLayout(spec)
	require.Empty(t, issues)

	row, err := layout.EncodeRow(meanOutcome(t, member))
	require.NoError(t, err)
	cells := cellMap(row)
	require.Equal(t, "20", cells["value_decimal"])
	require.Nil(t, cells["value_int"], "the raw int sample is not the window value")
	require.Nil(t, cells["value_float"])
	require.Equal(t, "decimal", cells["value_type"])
	require.Equal(t, longStatusOK, cells["status"])

	var provenance map[string]any
	require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &provenance))
	mean, ok := provenance["mean"].(map[string]any)
	require.True(t, ok, "a mean record carries its reduction evidence: %v", provenance)
	require.Equal(t, "mean", mean["reduction"])
	require.Equal(t, "gauge", mean["role"])
	require.Equal(t, "int64", mean["source_type"])
	require.Equal(t, "decimal", mean["result_type"])
	require.Equal(t, measurement.MeanAggregationVersion, mean["version"])
	require.Equal(t, measurement.MeanAggregationMethod, mean["method"])
	require.Equal(t, "exact", mean["basis"])
	require.EqualValues(t, 3, mean["valid_count"])
	require.EqualValues(t, 1, mean["rejected_count"])
	require.Equal(t, "partial", mean["coverage"])
	require.Equal(t, "bad", mean["reason"])
	require.Equal(t, "60", mean["sum"])
	require.EqualValues(t, 0, mean["scale"])
	require.Equal(t, "half-even", mean["rounding"])
	require.Equal(t, meanWindowStart.Add(3*time.Second).Format(longTimeText), mean["first_observed_at"])
	require.Equal(t, meanWindowStart.Add(12*time.Second).Format(longTimeText), mean["last_observed_at"])

	// A window with no usable mean stays NULL and keeps the frozen typed
	// contract plus the missing evidence; no raw sample is substituted.
	member.Reduction, member.ValueRole, member.SourceValueType = "mean", "gauge", schema.DataTypeUint64
	member.Type = measurement.ExactDecimal
	layout, issues = NewLongRowLayout(spec)
	require.Empty(t, issues)
	silent := meanWindowOutcome(t, member, snapshot.MemberResult{
		MemberKey: member.MemberKey, Status: snapshot.MemberMissing, Reason: snapshot.ReasonNoSamples,
		Evidence: &snapshot.MeanEvidence{
			Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod,
			Basis: measurement.MeanBasisExact, SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactDecimal,
			Coverage: snapshot.CoverageMissing, Reason: snapshot.ReasonNoSamples,
		},
	})
	row, err = layout.EncodeRow(silent)
	require.NoError(t, err)
	cells = cellMap(row)
	require.Nil(t, cells["value_decimal"])
	require.Equal(t, "decimal", cells["value_type"], "the frozen result type survives a NULL window")
	require.Equal(t, "missing", cells["status"])
	require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &provenance))
	mean, ok = provenance["mean"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "missing", mean["coverage"])
	require.Equal(t, snapshot.ReasonNoSamples, mean["reason"])
	require.EqualValues(t, 0, mean["valid_count"])

	// A latest state/counter_raw member states its frozen source/result metadata
	// without borrowing a mean block.
	state := member
	state.Reduction, state.Type, state.ValueRole, state.SourceValueType = "", measurement.ExactText, "state", schema.DataTypeString
	layout, issues = NewLongRowLayout(LongRowSpec{
		Dialect: SQLDialectSQLite, Columns: longColumnsFor(t, "sqlite"), WorkspaceID: "ws", GroupID: "g1", GroupRevision: "r1",
		ProfileKind: "instantaneous", ProfileInterval: 15, ProfileUnit: "second", Members: []LongRowMember{state},
	})
	require.Empty(t, issues)
	latest := meanWindowOutcome(t, state, snapshot.MemberResult{
		MemberKey: state.MemberKey, Status: snapshot.MemberOK,
		Sample: &snapshot.Sample{SampleID: "state-1", ObservedAt: meanWindowStart.Add(3 * time.Second), Quality: schema.QualityGood, Value: measurement.NewText("on")},
	})
	row, err = layout.EncodeRow(latest)
	require.NoError(t, err)
	cells = cellMap(row)
	require.Equal(t, "on", cells["value_text"])
	require.Equal(t, "text", cells["value_type"])
	provenance = nil
	require.NoError(t, json.Unmarshal([]byte(cells["provenance"].(string)), &provenance))
	require.Nil(t, provenance["mean"])
	require.Equal(t, "string", provenance["source_type"], "the frozen Tag type is kept, not the exact carrier type")
	require.Equal(t, "text", provenance["result_type"])
	require.Equal(t, "state", provenance["role"])

	// A mean member without its reduction result must never be encoded from the
	// raw sample: fail closed instead.
	member.Reduction = "mean"
	layout, issues = NewLongRowLayout(spec)
	require.Empty(t, issues)
	broken := meanOutcome(t, member)
	broken.Members[0].Reduction = nil
	_, err = layout.EncodeRow(broken)
	require.Error(t, err, "a mean window without a reduction value is refused, never written as the raw sample")
}
