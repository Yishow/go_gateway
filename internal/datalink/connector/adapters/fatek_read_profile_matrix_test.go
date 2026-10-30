package adapters

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
	"math"
	"sync"
	"testing"
)

func TestFatekResolvedProfileNumericMatrix(t *testing.T) {
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
	for _, address := range []string{"D100", "DR100"} {
		for _, f := range fixtures {
			for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
				t.Run(address+"/"+string(f.dtype)+"/"+format, func(t *testing.T) {
					payload := ""
					for _, w := range profileWireWords(f.bits, f.words, format) {
						payload += fmt.Sprintf("%04X", w)
					}
					tr := &fatekProfileTransport{payload: payload}
					req := fatekProfileRequest(t, address, f.dtype, format)
					got, err := fatekProfileConnector(tr).Read(t.Context(), req)
					require.NoError(t, err)
					require.Equal(t, schema.QualityGood, got.Quality)
					require.Equal(t, f.want, got.Value)
					expected := f.words
					if address == "DR100" {
						expected /= 2
					}
					require.Equal(t, []int{expected}, tr.counts)
				})
			}
		}
	}
}

func TestFatekResolvedProfileInterleavingAndNativeCompatibility(t *testing.T) {
	tr := &fatekProfileTransport{payload: "11223344"}
	c := fatekProfileConnector(tr)
	a := fatekProfileRequest(t, "D100", schema.DataTypeUint32, "ABCD")
	b := fatekProfileRequest(t, "DR100", schema.DataTypeUint32, "CDAB")
	for _, req := range []connector.ReadRequest{a, b, a, b} {
		got, err := c.Read(t.Context(), req)
		require.NoError(t, err)
		want := uint32(0x11223344)
		if req.Address == "DR100" {
			want = 0x33441122
		}
		require.Equal(t, want, got.Value)
	}
	require.Equal(t, []int{2, 1, 2, 1}, tr.counts)
	for _, mode := range []schema.PointFormatMode{schema.PointFormatLegacy, schema.PointFormatInherit} {
		p := schema.Point{ID: "P", DeviceID: "F", Address: "D100", DataType: schema.DataTypeUint32, DataFormat: "DCBA", FormatMode: mode, FormatVersion: 1}
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs})
		require.NoError(t, err)
		req := a
		req.Profile = profile
		req.DataFormat = "DCBA"
		got, err := c.Read(t.Context(), req)
		require.NoError(t, err)
		require.Equal(t, uint32(0x11223344), got.Value)
	}
	// Legacy count remains a native-component count; the scalar assembler uses
	// only the components required by the type, not a vector of typed values.
	tr.payload = "112233445566"
	a.Profile = connector.ReadProfile{}
	a.Count = 3
	got, err := c.Read(t.Context(), a)
	require.NoError(t, err)
	require.Equal(t, uint32(0x11223344), got.Value)
}

func TestFatekResolvedProfileMalformedResponsesAreNotGood(t *testing.T) {
	for _, tc := range []struct {
		name, payload, status string
		badChecksum           bool
	}{
		{"short", "1234", "", false}, {"extra", "123456781234", "", false}, {"odd", "1234567", "", false},
		{"invalidhex", "1234ZZZZ", "", false}, {"checksum", "12345678", "", true}, {"status", "12345678", "1", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tr := &fatekProfileTransport{payload: tc.payload, status: tc.status, badChecksum: tc.badChecksum}
			got, err := fatekProfileConnector(tr).Read(t.Context(), fatekProfileRequest(t, "D100", schema.DataTypeUint32, "CDAB"))
			require.Error(t, err)
			require.NotEqual(t, schema.QualityGood, got.Quality)
			require.Nil(t, got.Value)
		})
	}
}

func TestFatekResolvedProfileRejectsMismatchedRequestBeforeWire(t *testing.T) {
	tr := &fatekProfileTransport{payload: "12345678"}
	req := fatekProfileRequest(t, "D100", schema.DataTypeUint32, "ABCD")
	req.Address = "D101"
	got, err := fatekProfileConnector(tr).Read(t.Context(), req)
	require.ErrorIs(t, err, connector.ErrReadProfile)
	require.Equal(t, schema.QualityBad, got.Quality)
	require.Empty(t, tr.counts)
}

func TestFatekResolvedProfileConcurrentRequestsKeepOwnFormat(t *testing.T) {
	tr := &fatekProfileTransport{payload: "11223344"}
	c := fatekProfileConnector(tr)
	a := fatekProfileRequest(t, "D100", schema.DataTypeUint32, "ABCD")
	b := fatekProfileRequest(t, "DR100", schema.DataTypeUint32, "CDAB")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := a
			want := uint32(0x11223344)
			if i%2 == 1 {
				req = b
				want = 0x33441122
			}
			got, err := c.Read(t.Context(), req)
			assert.NoError(t, err)
			assert.Equal(t, want, got.Value)
		}(i)
	}
	wg.Wait()
}

func TestFatekModernUnsupportedTypeDoesNotFallBack(t *testing.T) {
	for _, address := range []string{"D100", "X0"} {
		t.Run(address, func(t *testing.T) {
			p := schema.Point{ID: "P", DeviceID: "F", Address: address, DataType: schema.DataType("unknown"), FormatMode: schema.PointFormatInherit, FormatVersion: 1}
			profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs})
			require.NoError(t, err)
			tr := &fatekProfileTransport{payload: "1"}
			if address == "D100" {
				tr.payload = "1234"
			}
			got, err := fatekProfileConnector(tr).Read(t.Context(), connector.ReadRequest{Address: p.Address, DataType: p.DataType, Count: 1, Profile: profile})
			require.Error(t, err)
			require.NotEqual(t, schema.QualityGood, got.Quality)
			require.Empty(t, tr.counts)
		})
	}
}

func TestFatekInherit16AndExplicitNativeCountCompatibility(t *testing.T) {
	for _, dtype := range []schema.DataType{schema.DataTypeInt16, schema.DataTypeUint16} {
		p := schema.Point{ID: "P", DeviceID: "F", Address: "D100", DataType: dtype, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: schema.ProtocolFatekFBs})
		require.NoError(t, err)
		tr := &fatekProfileTransport{payload: "FEDC"}
		got, err := fatekProfileConnector(tr).Read(t.Context(), connector.ReadRequest{Address: p.Address, DataType: dtype, Count: 1, Profile: profile})
		require.NoError(t, err)
		if dtype == schema.DataTypeInt16 {
			require.Equal(t, int16(-292), got.Value)
		} else {
			require.Equal(t, uint16(0xfedc), got.Value)
		}
	}
	// Modern numeric Count also remains native components, not typed values.
	// Sufficient extra requested components retain the scalar-first contract.
	tr := &fatekProfileTransport{payload: "112233445566"}
	req := fatekProfileRequest(t, "D100", schema.DataTypeUint32, "CDAB")
	req.Count = 3
	got, err := fatekProfileConnector(tr).Read(t.Context(), req)
	require.NoError(t, err)
	require.Equal(t, uint32(0x33441122), got.Value)
	require.Equal(t, []int{3}, tr.counts)
	req.Count = 1
	tr.payload = "1122"
	got, err = fatekProfileConnector(tr).Read(t.Context(), req)
	require.Error(t, err)
	require.NotEqual(t, schema.QualityGood, got.Quality)
}
