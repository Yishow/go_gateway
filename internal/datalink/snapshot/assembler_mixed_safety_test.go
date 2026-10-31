package snapshot

import (
	"go-gateway/internal/datalink/measurement"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMixedLatestFrozenTypeCannotBecomeCounterInput(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("b", "good", 10, 100))
	wrong := mixedSample("b", "wrong", 59, 0)
	wrong.Value = measurement.NewText("not a number")
	mixedOffer(t, a, wrong)
	p := a.Plan(mixedTime(60))
	r := p.Outcomes[3].Members[1]
	require.NotNil(t, r.Sample)
	require.Equal(t, "good", r.Sample.SampleID, "latest must rank frozen-type eligible reading")
	c = mixedCounterConfig(ModeMixedWide)
	a, err = NewAssembler(c)
	require.NoError(t, err)
	wrong = mixedIntegerSample("c", "wrong counter", 59, 1000)
	wrong.Value = measurement.NewInt64(1000)
	mixedOffer(t, a, wrong)
	p = a.Plan(mixedTime(60))
	require.Empty(t, p.CounterCheckpoints(), "wrong frozen exact type must not create baseline")
	require.Equal(t, MemberInvalid, p.Outcomes[3].Members[2].Status)
}

func TestMixedCutoverKeepsIndependentTruncatedWindows(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	mixedOffer(t, a, mixedSample("a", "before", 20, 22))
	mixedOffer(t, a, mixedSample("a", "after", 28, 999))
	mixedOffer(t, a, mixedSample("b", "slow-before", 20, 100))
	mixedOffer(t, a, mixedSample("b", "slow-after", 28, 999))
	require.NoError(t, a.SetMixedEndsAt(mixedTime(23)))
	require.Equal(t, OfferRejected, a.Check(mixedSample("b", "exact", 23, 100), mixedTime(23)).Outcome)
	p := a.Plan(mixedTime(60))
	require.NoError(t, p.Err)
	require.Len(t, p.Outcomes, 3)
	require.Equal(t, mixedTime(15), p.Outcomes[0].DueAt)
	require.Equal(t, mixedTime(30), p.Outcomes[1].DueAt)
	require.Equal(t, mixedTime(60), p.Outcomes[2].DueAt)
	require.Equal(t, "before", p.Outcomes[1].Members[0].Sample.SampleID)
	require.Equal(t, mixedTime(23), p.Outcomes[1].Members[0].Window.ActualEnd)
	require.True(t, p.Outcomes[2].Members[0].Window.NextDue.IsZero(), "ended member has no future due")
	require.True(t, p.Outcomes[2].Members[1].Window.NextDue.IsZero(), "last due result must not invent another window")
	require.Equal(t, "slow-before", p.Outcomes[2].Members[1].Sample.SampleID)
	require.Equal(t, mixedTime(23), p.Outcomes[2].Members[1].Window.ActualEnd)
	a.Commit(p)
	require.True(t, a.Done())
	require.True(t, a.Plan(mixedTime(120)).Empty())
	checkpoint, _ := a.MixedProgress()
	fresh, err := NewAssembler(a.config)
	require.NoError(t, err)
	require.NoError(t, fresh.RestoreMixedProgress(checkpoint))
	require.True(t, fresh.Done())
}

func TestMixedTypedProgressRejectsForeignAndStaleCommit(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	old := a.Plan(mixedTime(15))
	newer := a.Plan(mixedTime(60))
	a.Commit(newer)
	a.Commit(old)
	progress, _ := a.MixedProgress()
	require.Equal(t, mixedTime(60), progress.ClosedThrough, "stale closure must not rewind progress")
	c.GroupID = "foreign"
	b, err := NewAssembler(c)
	require.NoError(t, err)
	b.Commit(newer)
	progress, _ = b.MixedProgress()
	require.True(t, progress.ClosedThrough.IsZero())
	fresh, err := NewAssembler(mixedConfig(ModeMixedWide))
	require.NoError(t, err)
	require.Error(t, fresh.RestoreMixedProgress(DispatchWatermark{"w", "g", "r", mixedTime(16)}))
	require.Error(t, fresh.RestoreMixedProgress(DispatchWatermark{"w", "g", "r", mixedTime(0)}))
	require.NoError(t, fresh.RestoreMixedProgress(DispatchWatermark{"w", "g", "r", mixedTime(15)}))
	require.Error(t, fresh.SetMixedEndsAt(mixedTime(14)))
	require.True(t, fresh.NextClose().IsZero())
	require.Equal(t, time.Time{}, fresh.config.Until)
}

func TestMixedRestoreRejectsBoundaryOfAlreadyEndedMember(t *testing.T) {
	c := mixedConfig(ModeMixedWide)
	c.EndsAt = mixedTime(23)
	a, err := NewAssembler(c)
	require.NoError(t, err)
	// :45 is a 15s arithmetic boundary but not a dispatch after A ended at :30.
	require.Error(t, a.RestoreMixedProgress(DispatchWatermark{"w", "g", "r", mixedTime(45)}))
	require.NoError(t, a.RestoreMixedProgress(DispatchWatermark{"w", "g", "r", mixedTime(60)}))
	progress, err := a.MemberProgress()
	require.NoError(t, err)
	require.True(t, progress["a"].End.IsZero(), "ended member must not promise another due window")
	require.True(t, progress["b"].End.IsZero())
}
