package runtime

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func boundaryInstant(seconds float64) time.Time {
	return boundaryStart.Add(time.Duration(seconds * float64(time.Second)))
}

// acceptAt offers one complete set of readings, one per member, observed at at.
func (f *boundaryFixture) acceptAt(t *testing.T, boundary *GroupBoundary, id string, at time.Time) {
	t.Helper()
	for i, value := range []any{21.5, int64(101), true, "batch-001", uint64(7)} {
		envelope := f.envelope(i, id+f.members[i].TagID, 0, value)
		envelope.ObservedAt, envelope.ReceivedAt = at, at
		require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
	}
}

func (f *boundaryFixture) open(boundary *GroupBoundary) int {
	boundary.mu.Lock()
	defer boundary.mu.Unlock()
	return boundary.assembler.OpenBuckets()
}

// A cutover hands a source over at an exact instant that falls inside a window.
// Ownership is half-open: an observation before the instant stays with the
// revision that ends, one at or after it belongs to the revision that starts,
// and the window the instant cuts is closed with its real bounds and kept as a
// row that says it was cut: the accepted readings are not dropped, and a cut
// window is not passed off as a complete one.
func TestGroupBoundaryCutoverInstantsAreExactAndHalfOpen(t *testing.T) {
	instant := boundaryInstant(14.5) // inside the 10 s window [10,20)

	t.Run("the revision that ends keeps what was observed before the instant", func(t *testing.T) {
		f := newBoundaryFixture(t)
		f.clock = boundaryInstant(14.6)
		f.config.EndsAt, f.config.Until = instant, boundaryInstant(20)
		boundary := f.build(t)

		f.acceptAt(t, boundary, "at-the-instant", instant)
		require.Zero(t, f.open(boundary), "an observation at the instant belongs to the next owner")
		f.acceptAt(t, boundary, "just-before", instant.Add(-time.Microsecond))
		require.Equal(t, 1, f.open(boundary), "one window holds everything observed before the instant")

		f.clock = boundaryInstant(30)
		outcomes := boundary.assembler.Tick(f.clock)
		var cut snapshot.Outcome
		for _, outcome := range outcomes {
			if outcome.BucketStart.Equal(boundaryInstant(10)) {
				cut = outcome
			}
		}
		require.Equal(t, snapshot.OutcomeRow, cut.Kind, "what the window would have written is still written")
		require.Equal(t, "end", cut.Truncated)
		require.NotEmpty(t, cut.EffectKey, "with a stable identity")
		require.True(t, cut.WindowEnd.Equal(instant), "the window ends at the real instant: %s", cut.WindowEnd)
		require.True(t, cut.BucketEnd.Equal(boundaryInstant(20)), "its identity keeps the aligned end")

		encoded, err := boundary.layout.EncodeRow(cut)
		require.NoError(t, err)
		require.True(t, encoded.Partial, "the encoded row says it is cut")
	})

	t.Run("the revision that starts takes what is observed at the instant", func(t *testing.T) {
		f := newBoundaryFixture(t)
		f.clock = boundaryInstant(14.6)
		f.config.FirstBucket, f.config.StartsAt = boundaryInstant(10), instant
		boundary := f.build(t)

		f.acceptAt(t, boundary, "just-before", instant.Add(-time.Microsecond))
		require.Zero(t, f.open(boundary), "an observation before the instant belongs to the previous owner")
		f.acceptAt(t, boundary, "at-the-instant", instant)
		require.Equal(t, 1, f.open(boundary))

		outcomes := boundary.assembler.Tick(boundaryInstant(30))
		require.NotEmpty(t, outcomes)
		first := outcomes[0]
		require.True(t, first.BucketStart.Equal(boundaryInstant(10)))
		require.Equal(t, "start", first.Truncated)
		require.True(t, first.WindowStart.Equal(instant), "the window starts at the real instant: %s", first.WindowStart)
		require.True(t, first.WindowEnd.Equal(boundaryInstant(20)))
		require.Equal(t, snapshot.OutcomeRow, first.Kind)
		for _, later := range outcomes[1:] {
			require.Empty(t, later.Truncated, "only the window containing the instant is cut")
		}
	})

	t.Run("a cut window an unusable required member already skips stays skipped", func(t *testing.T) {
		f := newBoundaryFixture(t)
		f.clock = boundaryInstant(14.6)
		f.config.EndsAt, f.config.Until = instant, boundaryInstant(20)
		boundary := f.build(t)
		envelope := f.envelope(0, "only-one", 0, 21.5)
		envelope.ObservedAt, envelope.ReceivedAt = instant.Add(-time.Second), instant.Add(-time.Second)
		require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
		for _, outcome := range boundary.assembler.Tick(boundaryInstant(30)) {
			if outcome.BucketStart.Equal(boundaryInstant(10)) {
				require.Equal(t, snapshot.OutcomeSkipped, outcome.Kind, "the missing required members still decide")
				require.Equal(t, snapshot.ReasonIncompleteRequired, outcome.Reason)
			}
		}
	})

	t.Run("an instant that is a window boundary cuts nothing", func(t *testing.T) {
		f := newBoundaryFixture(t)
		f.clock = boundaryInstant(20)
		f.config.EndsAt, f.config.Until = boundaryInstant(20), boundaryInstant(20)
		boundary := f.build(t)
		f.acceptAt(t, boundary, "inside", boundaryInstant(19.999))
		for _, outcome := range boundary.assembler.Tick(boundaryInstant(30)) {
			require.Empty(t, outcome.Truncated)
		}
	})
}
