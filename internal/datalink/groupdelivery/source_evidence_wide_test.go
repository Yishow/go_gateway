package groupdelivery

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var wideColumns = []dbtarget.ColumnInfo{
	{Name: "record_id", DataType: "TEXT", PrimaryKey: true},
	{Name: "group_id", DataType: "TEXT"},
	{Name: "bucket_start", DataType: "TEXT"},
	{Name: "provenance", DataType: "TEXT"},
	{Name: "v_mean", DataType: "TEXT", Nullable: true},
}

var wideScope = SourceScope{
	WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, ProfileKind: "instantaneous",
	ConnectorID: testDestination.ConnectorID, ConnectorRevision: testDestination.ConnectorRevision,
	Wide: true, IntervalSeconds: 15, RecordKeyColumn: "record_id", GroupIDColumn: "group_id", BucketStartColumn: "bucket_start", ProvenanceColumn: "provenance",
}

var wideRef = SourceRef{
	DeviceID: "dev-1", PointID: "pt-a", TagID: "tag-a", Column: "v_mean",
	Reduction: "mean", Role: "gauge", SourceType: "uint64", ResultType: "decimal",
}

// wideMeanOutcome is a closed wide window of one mean source. A cut window names its real bounds exactly as the
// assembler does when a forward cutover shortens it; the nominal bucket stays the aligned 15 seconds.
func wideMeanOutcome(t *testing.T, cut string, from, to time.Time) snapshot.Outcome {
	t.Helper()
	start := at(0)
	// The observations of a window lie inside the window that was really closed, cut or not.
	observed := start
	if cut != "" {
		observed = from
	}
	mean, err := measurement.NewDecimal("20")
	require.NoError(t, err)
	sum, err := measurement.NewDecimal("40")
	require.NoError(t, err)
	member := SourceEntityKey(wideRef.DeviceID, wideRef.PointID, wideRef.TagID)
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, RecordID: "rec-wide", EffectKey: "eff-wide", BucketStart: start, BucketEnd: start.Add(15 * time.Second), Partial: true,
		Members: []snapshot.MemberResult{{
			MemberKey: member, Status: snapshot.MemberOK,
			Sample:    &snapshot.Sample{SampleID: "s-2", MemberKey: member, ObservedAt: observed.Add(3 * time.Second), Quality: schema.QualityGood, Value: measurement.NewUint64(30)},
			Reduction: &measurement.MeanResult{Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod, Value: mean, Sum: sum, Count: 2, Basis: measurement.MeanBasisExact, Rounding: measurement.MeanRoundingHalfEven},
			Evidence: &snapshot.MeanEvidence{
				Version: measurement.MeanAggregationVersion, Method: measurement.MeanAggregationMethod, Basis: measurement.MeanBasisExact,
				SourceValueType: schema.DataTypeUint64, ResultValueType: measurement.ExactDecimal, ValidCount: 2, RejectedCount: 0,
				FirstObservedAt: observed.Add(time.Second), LastObservedAt: observed.Add(3 * time.Second), Coverage: snapshot.CoveragePartial, Reason: snapshot.ReasonBoundaryTruncated,
			},
		}},
	}
	if cut != "" {
		outcome.WindowStart, outcome.WindowEnd, outcome.Truncated = from, to, cut
	}
	return outcome
}

// enqueueWide encodes the outcome with the real wide encoder and closes it into the outbox.
func enqueueWide(t *testing.T, store *Store, outcome snapshot.Outcome) string {
	t.Helper()
	layout, issues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: dbtarget.SQLDialectSQLite, Columns: wideColumns, Partial: true,
		RecordKeyColumn: "record_id", BucketStartColumn: "bucket_start", ProvenanceColumn: "provenance", GroupIDColumn: "group_id", GroupID: testKey.GroupID,
		Members: []dbtarget.GroupRowMember{{
			MemberKey: SourceEntityKey(wideRef.DeviceID, wideRef.PointID, wideRef.TagID), Column: "v_mean", Type: measurement.ExactDecimal,
			Reduction: "mean", ValueRole: "gauge", SourceValueType: schema.DataTypeUint64,
		}},
	})
	require.Empty(t, issues)
	row, err := layout.EncodeRow(outcome)
	require.NoError(t, err)
	require.NoError(t, store.CommitClosure(t.Context(), testKey, Closure{
		Destination: testDestination, NextClose: outcome.BucketEnd, Buckets: []ClosedBucket{{Outcome: outcome, Row: &row}},
	}))
	return outcome.EffectKey
}

// TestAverageDeliveryEvidenceWideWindows pins the real bounds of a wide window: a whole window is the aligned
// interval, and a window a cutover cut short is reported exactly as cut, never as the whole interval.
func TestAverageDeliveryEvidenceWideWindows(t *testing.T) {
	for _, tc := range []struct {
		name, cut  string
		from, to   time.Time
		start, end time.Time
	}{
		{"a whole window is the aligned interval", "", time.Time{}, time.Time{}, at(0), at(15)},
		{"a window cut at its end is reported as cut", "end", at(0), at(7), at(0), at(7)},
		{"a window cut at its start is reported as cut", "start", at(5), at(15), at(5), at(15)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store, _, _ := newTestStore(t)
			enqueueWide(t, store, wideMeanOutcome(t, tc.cut, tc.from, tc.to))

			got, err := store.SourceRecords(t.Context(), wideScope, wideRef)
			require.NoError(t, err)
			require.NotNil(t, got.Latest)
			require.True(t, got.Latest.Available, "a window the real encoder wrote is evidence: %+v", got.Latest)
			require.True(t, tc.start.Equal(got.Latest.WindowStart), "start %s", got.Latest.WindowStart)
			require.True(t, tc.end.Equal(got.Latest.WindowEnd), "end %s", got.Latest.WindowEnd)
			require.NotNil(t, got.Latest.Aggregation)
			require.Equal(t, "partial", got.Latest.Aggregation.Coverage)
		})
	}
}
