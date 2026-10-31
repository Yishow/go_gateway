package runtime

import (
	"go-gateway/internal/datalink/snapshot"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedBoundaryCutoffExcludesAcceptedFutureAndKeepsSlowPartial(t *testing.T) {
	for _, kind := range []string{"until", "ends"} {
		t.Run(kind, func(t *testing.T) {
			f := mixedBoundaryFixture(t)
			b, err := NewGroupBoundary(f.config)
			require.NoError(t, err)
			f.accept(t, b, f.envelope(0, "within", boundaryAt(14), int64(20)), f.envelope(0, "future", boundaryAt(29), int64(999)), f.envelope(1, "slow-valid", boundaryAt(14), "ON"), f.envelope(1, "slow-future", boundaryAt(29), "OFF"))
			if kind == "until" {
				require.NoError(t, b.SetUntil(boundaryAt(15)))
			} else {
				require.NoError(t, b.SetEndsAt(boundaryAt(15)))
			}
			require.NoError(t, b.Tick(t.Context(), boundaryAt(60)))
			slowPartial := false
			for _, row := range f.sink.rows {
				for _, member := range row.Outcome.Members {
					if member.Status == snapshot.MemberNotDue {
						continue
					}
					require.NotNil(t, member.Window)
					require.False(t, member.Window.ActualEnd.After(boundaryAt(15)), "data beyond stopped coverage")
					if member.Sample != nil {
						require.True(t, member.Sample.ObservedAt.Before(boundaryAt(15)))
						require.NotContains(t, member.Sample.SampleID, "future")
					}
					if member.MemberKey == memberKey("device-A", "point-S", "tag-S") {
						require.Equal(t, boundaryAt(60), row.Outcome.DueAt, "slow member keeps its nominal dispatch boundary")
						require.Equal(t, boundaryAt(15), member.Window.ActualEnd)
						require.Equal(t, "end", member.Window.Truncated)
						require.Equal(t, "ON", member.Sample.Value.Value())
						slowPartial = true
					}
				}
			}
			require.True(t, slowPartial, "valid slow member must not be dropped for having DueAt after cutoff")
		})
	}
}

// Existing source API permits ignored returns; the new error-aware tests below
// additionally require callers to observe invalid cutoffs.
func TestMixedBoundaryInvalidCutoffDoesNotMutateIntake(t *testing.T) {
	for _, kind := range []string{"until", "ends"} {
		t.Run(kind, func(t *testing.T) {
			f := mixedBoundaryFixture(t)
			b, err := NewGroupBoundary(f.config)
			require.NoError(t, err)
			if kind == "until" {
				err = b.SetUntil(time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC))
			} else {
				err = b.SetEndsAt(time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC))
			}
			var arithmetic *snapshot.WindowArithmeticError
			require.ErrorAs(t, err, &arithmetic)
			require.False(t, b.Retiring(), "invalid cutoff cannot become an intake-only fake stop")
		})
	}
}

func TestMixedBoundaryCutoffReplayDoesNotExtendOrRejectCommittedPartial(t *testing.T) {
	f := mixedBoundaryFixture(t)
	b, err := NewGroupBoundary(f.config)
	require.NoError(t, err)
	f.accept(t, b, f.envelope(1, "slow", boundaryAt(14), "ON"))
	require.NoError(t, b.SetUntil(boundaryAt(15)))
	require.NoError(t, b.Tick(t.Context(), boundaryAt(60)))
	require.NoError(t, b.SetUntil(boundaryAt(15)), "same actual cutoff remains valid after later nominal dispatch")
	require.NoError(t, b.SetEndsAt(boundaryAt(30)), "a later request cannot extend earlier coverage")
	require.Equal(t, boundaryAt(15), b.ReplayState().Until)
}

func TestMixedBoundaryCutoffCannotRewriteClosedCoverage(t *testing.T) {
	for _, kind := range []string{"until", "ends"} {
		t.Run(kind, func(t *testing.T) {
			f := mixedBoundaryFixture(t)
			b, err := NewGroupBoundary(f.config)
			require.NoError(t, err)
			require.NoError(t, b.Tick(t.Context(), boundaryAt(15)))
			before, _ := b.assembler.MemberProgress()
			if kind == "until" {
				err = b.SetUntil(boundaryAt(14))
			} else {
				err = b.SetEndsAt(boundaryAt(14))
			}
			require.ErrorIs(t, err, snapshot.ErrInvalidConfig)
			require.False(t, b.Retiring())
			after, _ := b.assembler.MemberProgress()
			require.Equal(t, before, after)
		})
	}
}
