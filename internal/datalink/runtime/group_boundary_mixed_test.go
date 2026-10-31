package runtime

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
	"testing"
	"time"
)

func mixedBoundaryFixture(t *testing.T) *longFixture {
	t.Helper()
	f := wideProfileFixture(t)
	_, err := f.db.Exec(`ALTER TABLE ` + wideFixtureTable + ` RENAME COLUMN bucket_start TO due_at`)
	require.NoError(t, err)
	f.config.Columns = sqliteColumns(t, f.db, wideFixtureTable)
	g := f.config.Group
	g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileMixed, ResultShape: workspace.WriteGroupResultShapeMemberCalculationV1}
	g.RecordID = "84ec2c79-2fd1-4f94-a514-5078b8b6d2af"
	g.RowPolicy.IntervalSeconds = 0
	g.RowPolicy.BucketStartColumn = ""
	g.RowPolicy.DueAtColumn = "due_at"
	g.Members = g.Members[:2]
	f.members = g.Members
	for i := range g.Members {
		g.Members[i].Cadence = &workspace.WriteGroupMemberCadence{Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond, Alignment: workspace.WriteGroupAlignmentUTCEpoch}
		g.Members[i].Calculation = workspace.WriteGroupMemberCalculationMean
	}
	g.Members[1].Calculation = workspace.WriteGroupMemberCalculationLatest
	g.Members[1].Cadence.Interval = 60
	f.config.FirstBucket = time.Time{}
	f.config.StartsAt = boundaryStart
	f.config.Counters = nil
	return f
}

func TestMixedProductionBoundaryUsesMemberWindowsAndFrozenIntent(t *testing.T) {
	f := mixedBoundaryFixture(t)
	b, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "mean-1", boundaryAt(1), int64(20)), f.envelope(0, "mean-2", boundaryAt(2), int64(24)), f.envelope(1, "state", boundaryAt(3), "ON"))
	// Mutating the caller's draft does not reinterpret accepted observations.
	f.config.Group.Members[0].Calculation = workspace.WriteGroupMemberCalculationLatest
	f.config.Group.Members[0].Cadence.Interval = 60
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	require.Len(t, f.sink.rows, 1)
	row := f.sink.rows[0]
	require.Equal(t, boundaryAt(15), row.Outcome.DueAt)
	require.Equal(t, snapshot.CalculationMean, row.Outcome.Members[0].Window.Calculation)
	digits, ok := row.Outcome.Members[0].Reduction.Value.DecimalDigits()
	require.True(t, ok)
	require.Equal(t, "22", digits)
	require.NoError(t, b.Tick(t.Context(), boundaryAt(60)))
	require.NotEmpty(t, f.sink.rows)
}

func TestMixedBoundaryTickRejectsUnrepresentableTime(t *testing.T) {
	f := mixedBoundaryFixture(t)
	b, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	err = b.Tick(t.Context(), time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC))
	var arithmetic *snapshot.WindowArithmeticError
	require.ErrorAs(t, err, &arithmetic)
	require.Empty(t, f.sink.rows)
}

func TestMixedBoundaryRefusesLegacyLedgerAndCommonCursor(t *testing.T) {
	f := mixedBoundaryFixture(t)
	f.config.Ledger = &boundsReplayLedger{}
	_, err := NewGroupBoundary(f.config)
	var blocked *GroupBoundaryError
	require.ErrorAs(t, err, &blocked)
	require.Equal(t, "mixed-durable-progress-unsupported", blocked.Code)
	f.config.Ledger = nil
	f.config.FirstBucket = boundaryStart
	_, err = NewGroupBoundary(f.config)
	require.Error(t, err)
}

func TestMixedBoundaryRejectsDurationOverflowBeforeMultiplication(t *testing.T) {
	f := mixedBoundaryFixture(t)
	f.config.Group.RowPolicy.AllowedLatenessSeconds = 1 << 55
	_, err := NewGroupBoundary(f.config)
	var arithmetic *snapshot.WindowArithmeticError
	require.ErrorAs(t, err, &arithmetic)
}

func TestMixedBoundaryFirstPartialAndExactEdge(t *testing.T) {
	f := mixedBoundaryFixture(t)
	f.config.StartsAt = boundaryAt(7)
	b, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(0, "before", boundaryAt(6), int64(999)), f.envelope(0, "first", boundaryAt(7), int64(20)), f.envelope(0, "edge", boundaryAt(15), int64(24)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
	require.Len(t, f.sink.rows, 1)
	first := f.sink.rows[0].Outcome.Members[0]
	require.Equal(t, boundaryAt(7), first.Window.ActualStart)
	require.Equal(t, "start", first.Window.Truncated)
	digits, _ := first.Reduction.Value.DecimalDigits()
	require.Equal(t, "20", digits)
	require.NoError(t, b.Tick(t.Context(), boundaryAt(30)))
	require.Len(t, f.sink.rows, 2)
	digits, _ = f.sink.rows[1].Outcome.Members[0].Reduction.Value.DecimalDigits()
	require.Equal(t, "24", digits)
}
