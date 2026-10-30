package adapters

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"go-gateway/lib/hsllogic"
	"math"
	"testing"
)

func TestMC3EResolvedProfileNumericMatrix(t *testing.T) {
	fixtures := []struct {
		dtype schema.DataType
		bits  uint64
		want  any
		words int
	}{
		{schema.DataTypeInt32, 0xfedcba98, int32(-19088744), 2},
		{schema.DataTypeUint32, 0x12345678, uint32(0x12345678), 2},
		{schema.DataTypeFloat32, uint64(math.Float32bits(123.25)), float32(123.25), 2},
		{schema.DataTypeInt64, 0xfedcba9876543210, int64(-81985529216486896), 4},
		{schema.DataTypeUint64, 0x1122334455667788, uint64(0x1122334455667788), 4},
		{schema.DataTypeFloat64, math.Float64bits(-123456.25), float64(-123456.25), 4},
	}
	for _, f := range fixtures {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			t.Run(string(f.dtype)+"/"+format, func(t *testing.T) {
				tr := &mcProfileTransport{words: profileWireWords(f.bits, f.words, format)}
				c := mcProfileConnector(tr)
				req := mcProfileRequest(t, f.dtype, format)
				result, err := c.Read(t.Context(), req)
				require.NoError(t, err)
				require.Equal(t, schema.QualityGood, result.Quality)
				require.Equal(t, f.want, result.Value)
				require.Equal(t, []int{f.words}, tr.requests)
				// MC Count continues to mean number of typed values, not native word count.
				req.Count = 2
				tr.words = append(tr.words, tr.words...)
				result, err = c.Read(t.Context(), req)
				require.NoError(t, err)
				require.Equal(t, []any{f.want, f.want}, result.Value)
				require.Equal(t, 2*f.words, tr.requests[1])
				tr.words = tr.words[:f.words-1]
				req.Count = 1
				result, err = c.Read(t.Context(), req)
				require.Error(t, err)
				require.Equal(t, schema.QualityBad, result.Quality)
				require.Nil(t, result.Value)
			})
		}
	}
}
func TestMC3EResolvedProfileInterleavedAndLegacyConnection(t *testing.T) {
	tr := &mcProfileTransport{words: []uint16{0x1122, 0x3344, 0x5566, 0x7788}}
	c := mcProfileConnector(tr)
	a := mcProfileRequest(t, schema.DataTypeUint64, "ABCD")
	b := mcProfileRequest(t, schema.DataTypeUint64, "CDAB")
	for _, tc := range []struct {
		req  connector.ReadRequest
		want uint64
	}{{a, 0x1122334455667788}, {b, 0x7788556633441122}, {a, 0x1122334455667788}} {
		got, err := c.Read(t.Context(), tc.req)
		require.NoError(t, err)
		require.Equal(t, tc.want, got.Value)
	}
	tr.words = []uint16{0x1234, 0x5678}
	// Legacy DataFormat never overrode the MC connection converter.
	legacy := connector.ReadRequest{Address: "D100", DataType: schema.DataTypeUint32, Count: 1, DataFormat: "ABCD"}
	got, err := c.Read(t.Context(), legacy)
	require.NoError(t, err)
	require.Equal(t, uint32(0x56781234), got.Value)
	tr.extra = []byte{0, 0}
	got, err = c.Read(t.Context(), legacy)
	require.NoError(t, err)
	require.Equal(t, uint32(0x56781234), got.Value, "legacy extra-payload policy remains unchanged")
}
func TestMC3EResolvedProfileInherit16KeepsNativeConverter(t *testing.T) {
	for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
		tr := &mcProfileTransport{words: []uint16{0x1234}}
		c := mcProfileConnector(tr)
		c.dataConverter = hsllogic.NewDataConverter(hsllogic.DataFormat(format))
		p := schema.Point{ID: "P", DeviceID: "D", Address: "D100", DataType: schema.DataTypeUint16, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolMC3E, ConnectionFormat: format, RuleOwned: true, RuleFormat: "ABCD"})
		require.NoError(t, err)
		req := connector.ReadRequest{Address: p.Address, DataType: p.DataType, Count: 1, Profile: profile}
		got, err := c.Read(t.Context(), req)
		require.NoError(t, err)
		require.Equal(t, c.dataConverter.RegistersToValue(tr.words, hsllogic.DataTypeUint16), got.Value)
	}
}
func TestMC3EResolvedProfileBadIdentityAndCountDoNotSend(t *testing.T) {
	tr := &mcProfileTransport{}
	c := mcProfileConnector(tr)
	req := mcProfileRequest(t, schema.DataTypeUint64, "ABCD")
	for _, count := range []int{-1, 32, int(^uint(0) >> 1)} {
		req.Count = count
		got, err := c.Read(t.Context(), req)
		require.Error(t, err)
		require.Equal(t, schema.QualityBad, got.Quality)
	}
	req.Count = 1
	req.Address = "D101"
	_, err := c.Read(t.Context(), req)
	require.ErrorIs(t, err, connector.ErrReadProfile)
	require.Empty(t, tr.requests)
}

func TestMC3EResolvedProfileLegacyTypedCountAndConcurrentIsolation(t *testing.T) {
	tr := &mcProfileTransport{words: []uint16{0x1122, 0x3344, 0x5566, 0x7788}}
	c := mcProfileConnector(tr)
	result, err := c.Read(t.Context(), connector.ReadRequest{Address: "D100", DataType: schema.DataTypeUint32, DataFormat: "ABCD", Count: 2})
	require.NoError(t, err)
	require.Equal(t, []any{uint32(0x33441122), uint32(0x77885566)}, result.Value)
	require.Equal(t, []int{4}, tr.requests)
	a, b := mcProfileRequest(t, schema.DataTypeUint64, "ABCD"), mcProfileRequest(t, schema.DataTypeUint64, "CDAB")
	type answer struct {
		value any
		err   error
		want  uint64
	}
	answers := make(chan answer, 20)
	for i := 0; i < 10; i++ {
		go func() { r, e := c.Read(t.Context(), a); answers <- answer{r.Value, e, 0x1122334455667788} }()
		go func() { r, e := c.Read(t.Context(), b); answers <- answer{r.Value, e, 0x7788556633441122} }()
	}
	for i := 0; i < 20; i++ {
		r := <-answers
		require.NoError(t, r.err)
		require.Equal(t, r.want, r.value)
	}
}
