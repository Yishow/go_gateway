package runtime

import (
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestLongBoundaryRefusesWhatItCannotWriteExactly(t *testing.T) {
	blockedCode := func(t *testing.T, f *longFixture) string {
		t.Helper()
		_, err := NewGroupBoundary(f.config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked)
		return blocked.Code
	}

	t.Run("a cumulative profile without confirmed counter semantics is never written as an instantaneous or wide row", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileCumulative, Interval: 10, Unit: workspace.WriteGroupIntervalUnitMinute}
		f.config.Group.RowPolicy.IntervalSeconds = 600
		require.Equal(t, "counter-semantics-required", blockedCode(t, f))
	})
	t.Run("an unknown profile kind is not supported", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: "average", Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond}
		require.Equal(t, "profile-unsupported", blockedCode(t, f))
	})
	t.Run("a long group without a profile is refused", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Group.RecordingProfile = nil
		require.Equal(t, "profile-unsupported", blockedCode(t, f))
	})
	t.Run("the row interval must be the one the profile derives", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Group.RowPolicy.IntervalSeconds = 60
		require.Equal(t, "profile-interval-mismatch", blockedCode(t, f))
	})
	t.Run("a table that cannot hold the fixed layout is refused with safe codes", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.Columns = f.config.Columns[1:] // bucket_start is gone
		_, err := NewGroupBoundary(f.config)
		var blocked *GroupBoundaryError
		require.ErrorAs(t, err, &blocked)
		require.Equal(t, "layout-blocked", blocked.Code)
		require.NotEmpty(t, blocked.Issues)
	})
	t.Run("a tag type without an exact SQL form is refused", func(t *testing.T) {
		f := newLongFixture(t)
		f.config.TagTypes["tag-A"] = schema.DataType("blob")
		require.Equal(t, "tag-type-unsupported", blockedCode(t, f))
	})
}

func TestProfileWideGroupsKeepSkipRowSemantics(t *testing.T) {
	f := newBoundaryFixture(t)
	boundary := f.build(t)
	values := []any{21.5, int64(101), true, "batch-001"}
	for i, value := range values { // the required counter member never reports
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(i, "s"+f.members[i].TagID, 4, value)))
	}
	require.NoError(t, boundary.Tick(t.Context(), boundaryAt(10)))
	require.Empty(t, f.sink.rows, "a wide group still skips a row whose required member is missing")
	require.Len(t, f.sink.outcomes, 1)
	require.Equal(t, snapshot.OutcomeSkipped, f.sink.outcomes[0].Kind)
	require.Equal(t, snapshot.ReasonIncompleteRequired, f.sink.outcomes[0].Reason)
}
