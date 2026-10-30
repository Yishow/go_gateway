package connector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestReadProfileFrozenAndValidated(t *testing.T) {
	p := schema.Point{ID: "P", RevisionID: "point-1", Address: "HR100", Function: "03", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	defaults := decoding.FormatContext{Protocol: schema.ProtocolModbusTCP, RuleOwned: true, RuleFormat: "CDAB", RuleRevision: "rule-2", ConnectionFormat: "BADC", ConnectionRevision: "conn-3"}
	profile, err := ResolveReadProfile(p, defaults)
	require.NoError(t, err)
	req := ReadRequest{Address: p.Address, Function: p.Function, DataType: p.DataType, Profile: profile}
	require.NoError(t, ValidateReadProfile(req))
	p.Address = "IR1"
	defaults.RuleFormat = "DCBA"
	detached := profile.Resolution()
	detached.EffectiveFormat = "ABCD"
	require.Equal(t, "CDAB", profile.Resolution().EffectiveFormat)
	require.Equal(t, "rule-2", profile.Resolution().RuleRevision)
	for _, mutate := range []func(*ReadRequest){
		func(r *ReadRequest) { r.Address = "HR101" },
		func(r *ReadRequest) { r.Function = "04" },
		func(r *ReadRequest) { r.DataType = schema.DataTypeUint64 },
		func(r *ReadRequest) { r.Profile.version = 2 },
		func(r *ReadRequest) { r.Profile.resolution.Resolved = false },
	} {
		changed := req
		mutate(&changed)
		require.ErrorIs(t, ValidateReadProfile(changed), ErrReadProfile)
	}
	require.NoError(t, ValidateReadProfile(ReadRequest{})) // Old API is legacy only.
}
func TestReadProfileRejectsInvalidModernIntent(t *testing.T) {
	for _, format := range []string{"UNKNOWN", ""} {
		profile, err := ResolveReadProfile(schema.Point{ID: "P", Address: "HR0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: format}, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP})
		require.Error(t, err)
		require.False(t, profile.Present())
	}
}

func TestReadProfileManagedConnectionRejectsChangedDefaultsBeforeWire(t *testing.T) {
	point := schema.Point{ID: "P", Address: "HR0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	profile, err := ResolveReadProfile(point, decoding.FormatContext{Protocol: schema.ProtocolModbusTCP, ConnectionFormat: "CDAB", ConnectionRevision: ReadConnectionRevision(`{"data_format":"CDAB"}`)})
	require.NoError(t, err)
	protocol := &stubProtocol{}
	conn := &ManagedConnection{Protocol: protocol, ProtocolType: schema.ProtocolModbusTCP, Config: `{"data_format":"ABCD"}`}
	result, err := conn.Read(t.Context(), ReadRequest{Address: point.Address, DataType: point.DataType, Profile: profile})
	require.ErrorIs(t, err, ErrReadProfile)
	require.Equal(t, schema.QualityBad, result.Quality)
}
