package workspace

import (
	"encoding/json"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemberCalculationWireRejectsInvalidIntent(t *testing.T) {
	for _, value := range []string{`"unknown"`, `""`, `null`, `12`, `true`, `{}`} {
		t.Run(value, func(t *testing.T) {
			var member WriteGroupMember
			err := json.Unmarshal([]byte(`{"calculation":`+value+`}`), &member)
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		})
	}
}

func TestMemberCalculationWireRetainsExplicitChoice(t *testing.T) {
	for _, method := range []string{"latest", "mean", "delta"} {
		var member WriteGroupMember
		require.NoError(t, json.Unmarshal([]byte(`{"calculation":"`+method+`"}`), &member))
		wire, err := json.Marshal(member)
		require.NoError(t, err)
		require.Contains(t, string(wire), `"calculation":"`+method+`"`)
	}
}

// These assertions cross the canonical JSON and profile-normalization boundary.
func TestMemberCalculationShapeBoundary(t *testing.T) {
	for _, profile := range []string{
		`{"kind":"mixed","result_shape":"member-calculation-v1"}`,
	} {
		var p WriteGroupRecordingProfile
		require.NoError(t, json.Unmarshal([]byte(profile), &p))
		wire, err := json.Marshal(p)
		require.NoError(t, err)
		require.Contains(t, string(wire), `"result_shape":"member-calculation-v1"`)
		require.ErrorIs(t, normalizeWriteGroupProfile(&WriteGroup{RecordingProfile: &p}), ErrWriteGroupValidation, "runtime support is not advertised by a typed model")
	}
	for _, profile := range []string{
		`{"kind":"instantaneous","interval":15,"unit":"second","result_shape":"member-calculation-v1"}`,
		`{"kind":"mixed"}`,
		`{"kind":"mixed","result_shape":"unknown"}`,
		`{"kind":"mixed","result_shape":null}`,
		`{"kind":"mixed","result_shape":"member-calculation-v1","sample_reduction":"latest"}`,
	} {
		var p WriteGroupRecordingProfile
		require.ErrorIs(t, json.Unmarshal([]byte(profile), &p), ErrWriteGroupValidation)
	}
}

func TestMemberCalculationCannotAttachToLegacy(t *testing.T) {
	for _, kind := range []WriteGroupProfileKind{WriteGroupProfileInstantaneous, WriteGroupProfileCumulative} {
		p, err := DefaultWriteGroupRecordingProfile(kind)
		require.NoError(t, err)
		g := &WriteGroup{RecordingProfile: &p, Members: []WriteGroupMember{{Calculation: WriteGroupMemberCalculationLatest}}}
		require.ErrorIs(t, validateMemberMetadataShape(g), ErrWriteGroupValidation)
	}
}

func TestMemberCalculationInitializationIsNotSubmissionDefault(t *testing.T) {
	source := WriteGroupMember{DeviceID: "device", PointID: "source-X", TagID: "tag-X"}
	a, err := InitializeMemberCalculation(source, "")
	require.NoError(t, err)
	require.Equal(t, WriteGroupMemberCalculationLatest, a.Calculation)
	explicit := source
	explicit.Calculation = WriteGroupMemberCalculationMean
	for _, defaultMethod := range []WriteGroupMemberCalculation{"", WriteGroupMemberCalculationDelta, WriteGroupMemberCalculationLatest} {
		next, err := InitializeMemberCalculation(explicit, defaultMethod)
		require.NoError(t, err)
		require.Equal(t, WriteGroupMemberCalculationMean, next.Calculation)
	}
	b, err := InitializeMemberCalculation(source, WriteGroupMemberCalculationDelta)
	require.NoError(t, err)
	require.Equal(t, WriteGroupMemberCalculationDelta, b.Calculation)
	a.Calculation = WriteGroupMemberCalculationMean
	require.Equal(t, WriteGroupMemberCalculationDelta, b.Calculation)
	require.Empty(t, source.Calculation)
	_, err = InitializeMemberCalculation(source, "unknown")
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	explicit.Calculation = "unknown"
	_, err = InitializeMemberCalculation(explicit, "")
	require.ErrorIs(t, err, ErrWriteGroupValidation)
}

func TestMemberCalculationIntentRequiresExplicitMethod(t *testing.T) {
	for _, method := range []WriteGroupMemberCalculation{"", "unknown", "latest", "mean", "delta"} {
		g := &WriteGroup{RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}, Members: []WriteGroupMember{{Calculation: method}}}
		err := ValidateMemberCalculationIntent(g)
		if method == "" || method == "unknown" {
			require.ErrorIs(t, err, ErrWriteGroupValidation)
		} else {
			require.NoError(t, err)
		}
	}
	require.ErrorIs(t, ValidateMemberCalculationIntent(nil), ErrWriteGroupValidation)
	require.ErrorIs(t, ValidateMemberCalculationIntent(&WriteGroup{Members: []WriteGroupMember{{Calculation: WriteGroupMemberCalculationLatest}}}), ErrWriteGroupValidation)
	require.ErrorIs(t, ValidateMemberCalculationIntent(&WriteGroup{RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed}}), ErrWriteGroupValidation)
	require.ErrorIs(t, ValidateMemberCalculationIntent(&WriteGroup{RecordingProfile: &WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}}), ErrWriteGroupValidation)
}

func TestMemberCalculationRemainsRuntimeUnsupported(t *testing.T) {
	require.True(t, memberMetadataUnsupported(&WriteGroupMember{Calculation: WriteGroupMemberCalculationLatest}))
}

func TestMemberCalculationLegacySemanticsRegression(t *testing.T) {
	for _, kind := range []WriteGroupProfileKind{WriteGroupProfileInstantaneous, WriteGroupProfileCumulative} {
		p, err := DefaultWriteGroupRecordingProfile(kind)
		require.NoError(t, err)
		g := &WriteGroup{Name: "before", RecordingProfile: &p, Members: []WriteGroupMember{{PointID: "source-X"}}}
		g.Name = "unrelated edit"
		require.NoError(t, ValidateMemberCalculationIntent(g))
		wire, err := json.Marshal(g)
		require.NoError(t, err)
		require.NotContains(t, string(wire), `"calculation"`)
		require.NotContains(t, string(wire), `"result_shape"`)
		if kind == WriteGroupProfileCumulative {
			require.Equal(t, WriteGroupSampleReductionLatest, g.RecordingProfile.EffectiveSampleReduction())
			require.Equal(t, 600, g.RecordingProfile.IntervalSeconds())
		}
	}
	for _, tc := range []struct {
		role WriteGroupValueRole
		want measurement.ExactType
	}{
		{WriteGroupValueRoleGauge, measurement.ExactDecimal},
		{WriteGroupValueRoleCounterRaw, measurement.ExactUint64},
		{WriteGroupValueRoleState, measurement.ExactUint64},
	} {
		actual, err := deriveMemberResultType(WriteGroupSampleReductionMean, tc.role, schema.DataTypeUint64)
		require.NoError(t, err)
		require.Equal(t, tc.want, actual)
	}
}

func TestMemberCalculationProfileRoundTrip(t *testing.T) {
	p := WriteGroupRecordingProfile{Kind: WriteGroupProfileMixed, ResultShape: WriteGroupResultShapeMemberCalculationV1}
	wire, err := json.Marshal(p)
	require.NoError(t, err)
	var decoded WriteGroupRecordingProfile
	require.NoError(t, json.Unmarshal(wire, &decoded))
	require.Equal(t, p, decoded)
	for _, change := range []func(*WriteGroupRecordingProfile){
		func(p *WriteGroupRecordingProfile) { p.Interval = 15 },
		func(p *WriteGroupRecordingProfile) { p.Unit = WriteGroupIntervalUnitSecond },
		func(p *WriteGroupRecordingProfile) { p.SampleReduction = WriteGroupSampleReductionLatest },
	} {
		invalid := p
		change(&invalid)
		require.ErrorIs(t, ValidateMemberCalculationIntent(&WriteGroup{RecordingProfile: &invalid, Members: []WriteGroupMember{{Calculation: WriteGroupMemberCalculationLatest}}}), ErrWriteGroupValidation)
	}
}
