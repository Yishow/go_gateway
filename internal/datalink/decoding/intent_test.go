package decoding

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestPointIntent_EqualFormatDistinctIntentAndReset(t *testing.T) {
	defaults := FormatContext{Protocol: schema.ProtocolModbusTCP, RuleOwned: true, RuleFormat: "ABCD", ConnectionFormat: "DCBA"}
	a := schema.Point{DataType: schema.DataTypeUint32, DataFormat: "ABCD", FormatMode: schema.PointFormatOverride, FormatOverride: "ABCD", FormatVersion: 1, RevisionID: "A-original"}
	b := schema.Point{DataType: schema.DataTypeUint32, DataFormat: "ABCD", FormatMode: schema.PointFormatInherit, FormatVersion: 1, RevisionID: "B-original"}
	ar, err := ResolvePoint(a, defaults)
	require.NoError(t, err)
	br, err := ResolvePoint(b, defaults)
	require.NoError(t, err)
	require.Equal(t, "ABCD", ar.EffectiveFormat)
	require.Equal(t, "ABCD", br.EffectiveFormat)
	require.NotEqual(t, ar.Mode, br.Mode)
	defaults.RuleFormat = "CDAB"
	ar, err = ResolvePoint(a, defaults)
	require.NoError(t, err)
	br, err = ResolvePoint(b, defaults)
	require.NoError(t, err)
	require.Equal(t, "ABCD", ar.EffectiveFormat)
	require.Equal(t, "CDAB", br.EffectiveFormat)
	inherit := schema.PointFormatInherit
	reset, err := WithIntent(a, IntentEdit{Mode: &inherit})
	require.NoError(t, err)
	require.Empty(t, reset.FormatOverride)
	require.NotEqual(t, a.RevisionID, reset.RevisionID)
	ar, err = ResolvePoint(reset, defaults)
	require.NoError(t, err)
	require.Equal(t, "CDAB", ar.EffectiveFormat)
	require.Equal(t, schema.PointFormatOverride, a.FormatMode, "candidate edits must not mutate their input")
	require.Equal(t, "B-original", b.RevisionID)
}

func TestPointIntent_InheritPrecedenceAndReadback(t *testing.T) {
	point := schema.Point{DataType: schema.DataTypeUint64, FormatMode: schema.PointFormatInherit, FormatVersion: 1, RevisionID: "point-1", DataFormat: "DCBA"}
	defaults := FormatContext{Protocol: schema.ProtocolMC3E, RuleOwned: true, RuleFormat: "ABCD", ConnectionFormat: "BADC", WorkspaceRevision: "workspace-1", RuleRevision: "rule-1", ConnectionRevision: "connection-1"}
	for _, tc := range []struct {
		owned                  bool
		rule, connection, want string
		source                 Source
	}{
		{true, "ABCD", "BADC", "ABCD", "rule"},
		{true, "", "BADC", "BADC", SourceConnection},
		{false, "ABCD", "BADC", "BADC", SourceConnection},
		{false, "", "", "CDAB", SourceNative},
	} {
		defaults.RuleOwned, defaults.RuleFormat, defaults.ConnectionFormat = tc.owned, tc.rule, tc.connection
		got, err := ResolvePoint(point, defaults)
		require.NoError(t, err)
		require.True(t, got.Resolved)
		require.Equal(t, tc.want, got.EffectiveFormat)
		require.Equal(t, tc.source, got.Source)
		require.Empty(t, got.ConfiguredFormat)
		require.True(t, got.Applicable)
		require.Equal(t, "point-1", got.PointRevision)
		require.Equal(t, "workspace-1", got.WorkspaceRevision)
		require.Equal(t, "connection-1", got.ConnectionRevision)
		if tc.owned {
			require.Equal(t, "rule-1", got.RuleRevision)
		} else {
			require.Empty(t, got.RuleRevision)
		}
	}
}

func TestPointIntent_LegacyPolicyAndNativeRemainDistinct(t *testing.T) {
	for _, mode := range []schema.PointFormatMode{"", schema.PointFormatLegacy, schema.PointFormatInherit} {
		point := schema.Point{DataType: schema.DataTypeUint32, DataFormat: "ABCD", FormatMode: mode}
		defaults := FormatContext{Protocol: schema.ProtocolFatekFBs}
		got, err := ResolvePoint(point, defaults)
		require.NoError(t, err)
		require.Empty(t, got.EffectiveFormat)
		require.True(t, got.NativeCompatibility)
		defaults = FormatContext{Protocol: schema.ProtocolMC3E, ConnectionFormat: " cdab "}
		got, err = ResolvePoint(point, defaults)
		require.NoError(t, err)
		require.Equal(t, "ABCD", got.EffectiveFormat)
	}
}

func TestPointIntent_OmissionAndExplicitEqualOverride(t *testing.T) {
	point := schema.Point{DataType: schema.DataTypeUint32, DataFormat: "ABCD", RevisionID: "old"}
	unchanged, err := WithIntent(point, IntentEdit{})
	require.NoError(t, err)
	require.Equal(t, point, unchanged)
	override, format := schema.PointFormatOverride, "ABCD"
	edited, err := WithIntent(point, IntentEdit{Mode: &override, ConfiguredFormat: &format})
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatOverride, edited.FormatMode)
	require.Equal(t, "ABCD", edited.FormatOverride)
	require.Equal(t, schema.PointFormatVersion, edited.FormatVersion)
	require.NotEmpty(t, edited.RevisionID)
	require.NotEqual(t, "old", edited.RevisionID)
	same, err := WithIntent(edited, IntentEdit{Mode: &override, ConfiguredFormat: &format})
	require.NoError(t, err)
	require.Equal(t, edited, same, "identical candidate must not churn revision")
}

func TestPointIntent_InvalidMetadataFailsWithoutCandidateMutation(t *testing.T) {
	point := schema.Point{DataType: schema.DataTypeUint32, RevisionID: "old"}
	override, invalid, legacy := schema.PointFormatOverride, schema.PointFormatMode("secret://bad"), schema.PointFormatLegacy
	bad, canonical := "secret://bad", "ABCD"
	for _, edit := range []IntentEdit{
		{Mode: &invalid}, {Mode: &legacy}, {Mode: &override}, {Mode: &override, ConfiguredFormat: &bad}, {ConfiguredFormat: &canonical},
	} {
		got, err := WithIntent(point, edit)
		require.ErrorIs(t, err, ErrInvalidIntent)
		require.Equal(t, point, got)
		require.NotContains(t, err.Error(), "secret")
	}
	got, err := ResolvePoint(schema.Point{FormatMode: invalid}, FormatContext{Protocol: schema.ProtocolMC3E})
	require.ErrorIs(t, err, ErrInvalidIntent)
	require.Zero(t, got)
}

func TestPointIntent_ReadErrorsAndCanonicalFormats(t *testing.T) {
	for _, point := range []schema.Point{
		{FormatMode: schema.PointFormatOverride},
		{FormatMode: schema.PointFormatOverride, FormatOverride: "secret://bad"},
		{FormatMode: schema.PointFormatInherit, FormatOverride: "ABCD"},
		{FormatOverride: "ABCD"},
		{FormatVersion: 2},
		{FormatVersion: -1},
	} {
		got, err := ResolvePoint(point, FormatContext{Protocol: schema.ProtocolModbusTCP})
		require.ErrorIs(t, err, ErrInvalidIntent)
		require.Zero(t, got)
		require.NotContains(t, err.Error(), "secret")
	}
	for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
		got, err := ResolvePoint(schema.Point{DataType: schema.DataTypeFloat64, FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1},
			FormatContext{Protocol: schema.ProtocolModbusTCP, ConnectionFormat: "DCBA"})
		require.NoError(t, err)
		require.Equal(t, format, got.ConfiguredFormat)
		require.Equal(t, format, got.EffectiveFormat)
		require.Equal(t, SourceOverride, got.Source)
		require.Equal(t, PolicyOverride, got.Policy)
	}
	got, err := ResolvePoint(schema.Point{FormatMode: schema.PointFormatInherit}, FormatContext{Protocol: schema.ProtocolMQTT})
	require.ErrorIs(t, err, ErrUnsupportedProtocol)
	require.Zero(t, got)
}

func TestPointIntent_ResetAndRetuneCandidates(t *testing.T) {
	inherit, override := schema.PointFormatInherit, schema.PointFormatOverride
	empty, format := "", "CDAB"
	original := schema.Point{FormatMode: override, FormatOverride: "ABCD", FormatVersion: 1, RevisionID: "old"}
	reset, err := WithIntent(original, IntentEdit{Mode: &inherit, ConfiguredFormat: &empty})
	require.NoError(t, err)
	require.Empty(t, reset.FormatOverride)
	require.Equal(t, inherit, reset.FormatMode)
	retuned, err := WithIntent(original, IntentEdit{ConfiguredFormat: &format})
	require.NoError(t, err)
	require.Equal(t, "CDAB", retuned.FormatOverride)
	require.NotEqual(t, original.RevisionID, retuned.RevisionID)
	rejected, err := WithIntent(original, IntentEdit{Mode: &inherit, ConfiguredFormat: &format})
	require.ErrorIs(t, err, ErrInvalidIntent)
	require.Equal(t, original, rejected)
	future := original
	future.FormatVersion = 2
	rejected, err = WithIntent(future, IntentEdit{Mode: &inherit})
	require.ErrorIs(t, err, ErrInvalidIntent)
	require.Equal(t, future, rejected)
}

func TestPointIntent_InheritedSingleWordDriverPolicies(t *testing.T) {
	for _, tc := range []struct {
		protocol schema.ProtocolType
		want     string
		source   Source
		native   bool
	}{
		{schema.ProtocolModbusTCP, "BADC", SourceRule, false},
		{schema.ProtocolMC3E, "DCBA", SourceConnection, false},
		{schema.ProtocolFatekFBs, "", SourceNative, true},
	} {
		got, err := ResolvePoint(schema.Point{DataType: schema.DataTypeUint16, FormatMode: schema.PointFormatInherit, FormatVersion: 1},
			FormatContext{Protocol: tc.protocol, RuleOwned: true, RuleFormat: "BADC", ConnectionFormat: "DCBA"})
		require.NoError(t, err)
		require.False(t, got.Applicable)
		require.True(t, got.Resolved)
		require.Equal(t, tc.want, got.EffectiveFormat)
		require.Equal(t, tc.source, got.Source)
		require.Equal(t, tc.native, got.NativeCompatibility)
	}
}

func TestPointIntent_InheritedFatekExplicitDefaults(t *testing.T) {
	point := schema.Point{DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	for _, tc := range []struct {
		rule, connection, want string
		source                 Source
	}{
		{"BADC", "DCBA", "BADC", SourceRule}, {"", "CDAB", "CDAB", SourceConnection},
	} {
		got, err := ResolvePoint(point, FormatContext{Protocol: schema.ProtocolFatekFBs, RuleOwned: true, RuleFormat: tc.rule, ConnectionFormat: tc.connection})
		require.NoError(t, err)
		require.True(t, got.Resolved)
		require.Equal(t, tc.want, got.EffectiveFormat)
		require.Equal(t, tc.source, got.Source)
		require.False(t, got.NativeCompatibility)
	}
	for _, defaults := range []FormatContext{
		{Protocol: schema.ProtocolMC3E, RuleOwned: true, RuleFormat: "secret://bad"},
		{Protocol: schema.ProtocolFatekFBs, ConnectionFormat: "secret://bad"},
	} {
		got, err := ResolvePoint(point, defaults)
		require.ErrorIs(t, err, ErrInvalidIntent)
		require.Zero(t, got)
	}
}
