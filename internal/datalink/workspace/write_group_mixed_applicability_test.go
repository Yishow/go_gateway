package workspace

import (
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedCalculationApplicability(t *testing.T) {
	group := &WriteGroup{RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}}
	for _, kind := range []schema.DataType{schema.DataTypeBool, schema.DataTypeString} {
		member := WriteGroupMember{Calculation: WriteGroupMemberCalculationMean}
		require.Error(t, applyRecordingMetadata(group, &member, 0, kind))
		member.Calculation = WriteGroupMemberCalculationDelta
		require.Error(t, applyRecordingMetadata(group, &member, 0, kind))
		member.Calculation = WriteGroupMemberCalculationLatest
		require.NoError(t, applyRecordingMetadata(group, &member, 0, kind))
	}
	member := WriteGroupMember{Calculation: WriteGroupMemberCalculationMean}
	require.NoError(t, applyRecordingMetadata(group, &member, 0, schema.DataTypeInt64))
	require.Equal(t, measurement.ExactDecimal, member.ResultValueType)
	member = WriteGroupMember{Calculation: WriteGroupMemberCalculationDelta}
	require.Error(t, applyRecordingMetadata(group, &member, 0, schema.DataTypeInt64), "delta needs explicit confirmed binding")
}

func TestMixedDraftStillCannotActivateRuntime(t *testing.T) {
	group := longGroupFixture(WriteGroupProfileInstantaneous, WriteGroupStorageModePlain, false)
	group.RecordingProfile = &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
	group.RowPolicy.IntervalSeconds = 0
	for i := range group.Members {
		group.Members[i].Calculation = WriteGroupMemberCalculationLatest
		group.Members[i].Cadence = &WriteGroupMemberCadence{Interval: 15, Unit: WriteGroupIntervalUnitSecond, Alignment: WriteGroupAlignmentUTCEpoch}
	}
	result := finalizedLong(group)
	require.False(t, result.Ready)
	require.Contains(t, issueCodesOf(result), "recording-reduction-unavailable")
}
