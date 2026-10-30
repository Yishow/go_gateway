package workspace

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemberCadenceStrictJSONAndBounds(t *testing.T) {
	for _, wire := range []string{
		`{"interval":1.5,"unit":"second","alignment":"utc_epoch"}`,
		`{"interval":1e1,"unit":"second","alignment":"utc_epoch"}`,
		`{"interval":null,"unit":"second","alignment":"utc_epoch"}`,
		`{"interval":0,"unit":"second","alignment":"utc_epoch"}`,
		`{"interval":1,"unit":"hour","alignment":"utc_epoch"}`,
		`{"interval":1,"unit":"second","alignment":"activation"}`,
		`{"interval":9223372037,"unit":"second","alignment":"utc_epoch"}`,
		`{"interval":153722868,"unit":"minute","alignment":"utc_epoch"}`,
	} {
		var cadence WriteGroupMemberCadence
		require.Error(t, json.Unmarshal([]byte(wire), &cadence), wire)
	}
	for _, seconds := range []int64{5, 15, 60, 600, 9223372036} {
		cadence := WriteGroupMemberCadence{Interval: seconds, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}
		duration, err := cadence.Duration()
		require.NoError(t, err)
		require.Equal(t, time.Duration(seconds)*time.Second, duration)
		wire, err := json.Marshal(cadence)
		require.NoError(t, err)
		var again WriteGroupMemberCadence
		require.NoError(t, json.Unmarshal(wire, &again))
		require.Equal(t, cadence, again)
	}
}

func TestMemberCadenceAbsenceAndClone(t *testing.T) {
	var legacy WriteGroupMember
	require.NoError(t, json.Unmarshal([]byte(`{"device_id":"d"}`), &legacy))
	require.Nil(t, legacy.Cadence)
	wire, err := json.Marshal(legacy)
	require.NoError(t, err)
	require.NotContains(t, string(wire), "cadence")
	for _, raw := range []string{`{"cadence":null}`, `{"cadence":{}}`} {
		require.Error(t, json.Unmarshal([]byte(raw), &legacy))
	}
	original := &WriteGroup{Members: []WriteGroupMember{{Cadence: &WriteGroupMemberCadence{Interval: 15, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}}}}
	clone := cloneWriteGroup(original)
	clone.Members[0].Cadence.Interval = 60
	require.Equal(t, int64(15), original.Members[0].Cadence.Interval)
}

func TestMixedCadenceNormalizationHasNoGroupWindow(t *testing.T) {
	group := &WriteGroup{
		RecordID:         "44f0caba-3624-4d44-8f79-48ccae515bfb",
		RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1},
		Destination:      WriteGroupDestination{StorageStrategy: WriteGroupStorageStrategyManaged},
		Members:          []WriteGroupMember{{Calculation: WriteGroupMemberCalculationMean, Cadence: &WriteGroupMemberCadence{Interval: 15, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}}},
	}
	require.NoError(t, normalizeWriteGroupProfile(group))
	require.Zero(t, group.RowPolicy.IntervalSeconds)
	for _, change := range []func(*WriteGroup){
		func(g *WriteGroup) { g.RowPolicy.IntervalSeconds = 60 },
		func(g *WriteGroup) { g.Members[0].Cadence = nil },
		func(g *WriteGroup) { g.RecordID = "" },
		func(g *WriteGroup) { g.RecordingProfile.SampleReduction = WriteGroupSampleReductionMean },
	} {
		bad := cloneWriteGroup(group)
		change(bad)
		require.Error(t, normalizeWriteGroupProfile(bad))
	}
}

func TestMemberCadenceWindowMetadataDoesNotDrift(t *testing.T) {
	activation := time.Date(2026, 10, 30, 12, 3, 17, 0, time.UTC)
	for _, seconds := range []int64{5, 15, 60, 600} {
		cadence := WriteGroupMemberCadence{Interval: seconds, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}
		start, end, err := cadence.WindowAt(activation)
		require.NoError(t, err)
		require.False(t, start.After(activation))
		require.True(t, end.After(activation))
		require.Zero(t, start.Unix()%seconds)
		raw, err := json.Marshal(cadence)
		require.NoError(t, err)
		var restart WriteGroupMemberCadence
		require.NoError(t, json.Unmarshal(raw, &restart))
		afterStart, afterEnd, err := restart.WindowAt(activation.In(time.FixedZone("UTC+8", 8*3600)))
		require.NoError(t, err)
		require.Equal(t, start, afterStart)
		require.Equal(t, end, afterEnd)
		exact, _, err := restart.WindowAt(end)
		require.NoError(t, err)
		require.Equal(t, end, exact)
	}
}
