package workspace

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mixedScheduleGroup(at time.Time, seconds ...int64) *WriteGroup {
	g := &WriteGroup{RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed}, UpdatedAt: at}
	for _, n := range seconds {
		g.Members = append(g.Members, WriteGroupMember{Cadence: &WriteGroupMemberCadence{Interval: n, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}})
	}
	return g
}
func TestMixedOrdinaryApplySchedulePreservesIndependentPhase(t *testing.T) {
	now := time.Date(2026, 10, 31, 12, 0, 7, 123, time.UTC)
	candidate := mixedScheduleGroup(now, 15, 60)
	at, err := nextAlignedCutover(candidate, nil)
	require.NoError(t, err)
	require.Equal(t, now, at)
	active := &WriteGroupAppliedSnapshot{Group: mixedScheduleGroup(now.Add(-time.Hour), 15, 60)}
	at, err = nextAlignedCutover(candidate, active)
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 10, 31, 12, 1, 0, 0, time.UTC), at)
	candidate = mixedScheduleGroup(time.Date(2026, 10, 31, 12, 1, 0, 0, time.UTC), 5, 15)
	at, err = nextAlignedCutover(candidate, active)
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 10, 31, 12, 2, 0, 0, time.UTC), at, "replacement always chooses next common boundary")
	require.Zero(t, candidate.RowPolicy.IntervalSeconds)
	require.Zero(t, active.Group.RowPolicy.IntervalSeconds)
}
func TestMixedOrdinaryApplyScheduleRejectsUnrepresentableBoundaries(t *testing.T) {
	now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
	active := &WriteGroupAppliedSnapshot{Group: mixedScheduleGroup(now, 9223372036)}
	_, err := nextAlignedCutover(mixedScheduleGroup(now, 9223372035), active)
	require.ErrorIs(t, err, ErrWriteGroupIntervalOverflow)
	_, err = nextAlignedCutover(mixedScheduleGroup(time.Date(2300, 1, 1, 0, 0, 0, 0, time.UTC), 15), nil)
	require.ErrorIs(t, err, ErrWriteGroupIntervalOverflow)
}
func TestLegacyOrdinaryApplyScheduleUnchanged(t *testing.T) {
	now := time.Date(2026, 10, 31, 12, 0, 7, 0, time.UTC)
	candidate := &WriteGroup{UpdatedAt: now, RowPolicy: WriteGroupRowPolicy{IntervalSeconds: 15}}
	active := &WriteGroupAppliedSnapshot{Group: &WriteGroup{RowPolicy: WriteGroupRowPolicy{IntervalSeconds: 60}}}
	expected, err := nextWriteGroupCutover(now, 60, 15)
	require.NoError(t, err)
	got, err := nextAlignedCutover(candidate, active)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
