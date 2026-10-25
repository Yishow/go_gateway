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

func wideProfileColumns() []ColumnInfo {
	return []ColumnInfo{
		{Name: "record_id", DataType: "TEXT", PrimaryKey: true},
		{Name: "group_id", DataType: "TEXT"},
		{Name: "bucket_start", DataType: "TEXT"},
		{Name: "provenance", DataType: "TEXT"},
		{Name: "v_mean", DataType: "TEXT", Nullable: true},
		{Name: "v_state", DataType: "INTEGER", Nullable: true},
		{Name: "v_legacy", DataType: "TEXT", Nullable: true},
	}
}

func wideProfileSpec(mutate func(*GroupRowSpec)) GroupRowSpec {
	spec := GroupRowSpec{
		Dialect: SQLDialectSQLite, Columns: wideProfileColumns(), Partial: true,
		RecordKeyColumn: "record_id", BucketStartColumn: "bucket_start", ProvenanceColumn: "provenance",
		GroupIDColumn: "group_id", GroupID: "g1",
		Members: []GroupRowMember{
			{MemberKey: "mean", Column: "v_mean", Type: measurement.ExactDecimal, Reduction: "mean", ValueRole: "gauge", SourceValueType: schema.DataTypeInt16},
			{MemberKey: "state", Column: "v_state", Type: measurement.ExactInt64, ValueRole: "state", SourceValueType: schema.DataTypeInt16},
			{MemberKey: "legacy", Column: "v_legacy", Type: measurement.ExactText},
		},
	}
	if mutate != nil {
		mutate(&spec)
	}
	return spec
}

func wideProfileOutcome(t *testing.T) snapshot.Outcome {
	t.Helper()
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	twenty, err := measurement.NewDecimal("20")
	require.NoError(t, err)
	sixty, err := measurement.NewDecimal("60")
	require.NoError(t, err)
	return snapshot.Outcome{
		Kind: snapshot.OutcomeRow, RecordID: "rec-wide", EffectKey: "eff-wide", BucketStart: start, BucketEnd: start.Add(15 * time.Second),
		Partial: true,
		Members: []snapshot.MemberResult{
			{
				MemberKey: "mean", Status: snapshot.MemberOK,
				Sample: &snapshot.Sample{SampleID: "m-3", ObservedAt: start.Add(9 * time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(30)},
				Reduction: &measurement.MeanResult{
					Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod,
					Value: twenty, Sum: sixty, Count: 3, Basis: measurement.MeanBasisExact, Rounding: measurement.MeanRoundingHalfEven,
				},
				Evidence: &snapshot.MeanEvidence{
					Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod, Basis: measurement.MeanBasisExact,
					SourceValueType: schema.DataTypeInt16, ResultValueType: measurement.ExactDecimal, ValidCount: 3,
					FirstObservedAt: start.Add(3 * time.Second), LastObservedAt: start.Add(9 * time.Second),
					Coverage: snapshot.CoverageUnknown, Reason: snapshot.ReasonCoverageUnverified,
				},
			},
			{
				MemberKey: "state", Status: snapshot.MemberOK,
				Sample: &snapshot.Sample{SampleID: "s-2", ObservedAt: start.Add(5 * time.Second), Quality: schema.QualityGood, Value: measurement.NewInt64(30)},
			},
			{MemberKey: "legacy", Status: snapshot.MemberMissing, Reason: snapshot.ReasonNoSamples},
		},
	}
}

// TestWideProfileRowStoresTheReductionAndKeepsLegacyEvidence pins the typed wide
// encoder: a mean member stores its exact reduction value, never a raw sample; a
// latest member stores its newest reading; an unusable optional member is NULL;
// provenance names the frozen types and the window evidence of the new members and
// keeps the exact JSON a legacy member always had.
func TestWideProfileRowStoresTheReductionAndKeepsLegacyEvidence(t *testing.T) {
	layout, issues := NewGroupRowLayout(wideProfileSpec(nil))
	require.Empty(t, issues)
	row, err := layout.EncodeRow(wideProfileOutcome(t))
	require.NoError(t, err)
	cells := cellMap(row)
	require.Equal(t, "20", cells["v_mean"], "the exact mean, as exact text")
	require.Equal(t, int64(30), cells["v_state"], "a state keeps its newest reading")
	require.Nil(t, cells["v_legacy"], "an unusable optional member is NULL, never zero")
	require.True(t, row.Partial)

	provenance := cells["provenance"].(string)
	var entries []provenanceEntry
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	require.Len(t, entries, 3)
	require.NotNil(t, entries[0].Mean)
	require.Equal(t, [4]string{"mean", "gauge", "int16", "decimal"}, [4]string{entries[0].Mean.Reduction, entries[0].Mean.Role, entries[0].Mean.SourceType, entries[0].Mean.ResultType})
	require.Equal(t, [3]any{"60", int64(3), "unknown"}, [3]any{entries[0].Mean.Sum, entries[0].Mean.ValidCount, entries[0].Mean.Coverage})
	require.Equal(t, [3]string{"state", "int16", "int64"}, [3]string{entries[1].Role, entries[1].SourceType, entries[1].ResultType})
	require.Nil(t, entries[1].Mean)
	require.Contains(t, provenance, `{"member":"legacy","status":"missing","reason":"no-samples"}`,
		"a member without metadata keeps exactly its legacy JSON")
}

func TestWideProfileRowRefusesWhatItCannotStoreExactly(t *testing.T) {
	layout, issues := NewGroupRowLayout(wideProfileSpec(nil))
	require.Empty(t, issues)

	noReduction := wideProfileOutcome(t)
	noReduction.Members[0].Reduction = nil
	_, err := layout.EncodeRow(noReduction)
	require.Equal(t, rowErrTypeMismatch, rowCode(t, err), "a mean member never falls back to a raw sample")

	rawStandsIn := wideProfileOutcome(t)
	rawStandsIn.Members[1].Reduction = wideProfileOutcome(t).Members[0].Reduction
	_, err = layout.EncodeRow(rawStandsIn)
	require.Equal(t, rowErrTypeMismatch, rowCode(t, err), "a latest member never carries a reduction")

	float, convErr := measurement.NewFloat64(20)
	require.NoError(t, convErr)
	wrongType := wideProfileOutcome(t)
	wrongType.Members[0].Reduction.Value = float
	_, err = layout.EncodeRow(wrongType)
	require.Equal(t, rowErrTypeMismatch, rowCode(t, err), "a float result never fills a decimal column")
}

// TestWideProfileLayoutNeedsNullableColumnsAndProvenance pins the storage gate:
// a table whose member column cannot hold NULL, or that has no provenance
// column, is refused when the layout is built; nothing is altered to fit.
func TestWideProfileLayoutNeedsNullableColumnsAndProvenance(t *testing.T) {
	_, issues := NewGroupRowLayout(wideProfileSpec(func(s *GroupRowSpec) { s.Columns[4].Nullable = false }))
	require.Contains(t, codesOfIssues(issues), issueColumnNotNullable)

	_, issues = NewGroupRowLayout(wideProfileSpec(func(s *GroupRowSpec) { s.ProvenanceColumn = "" }))
	require.Contains(t, codesOfIssues(issues), issuePartialNeedsProvenanc)

	_, issues = NewGroupRowLayout(wideProfileSpec(func(s *GroupRowSpec) { s.Columns[4].DataType = "REAL" }))
	require.Contains(t, codesOfIssues(issues), issueUnsupportedSQLType, "a float column cannot hold an exact decimal")
}

func codesOfIssues(issues []LayoutIssue) []string {
	codes := make([]string, 0, len(issues))
	for _, issue := range issues {
		codes = append(codes, issue.Code)
	}
	return codes
}
