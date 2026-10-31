package adapters

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestMCPlannedCompleteScalarMatrix(t *testing.T) {
	for _, tc := range []struct {
		dtype schema.DataType
		words []uint16
		want  any
	}{
		{schema.DataTypeInt32, []uint16{0x8123, 0x4567}, int32(-2128394905)},
		{schema.DataTypeUint32, []uint16{0x0102, 0x0304}, uint32(0x01020304)},
		{schema.DataTypeFloat32, []uint16{0xc120, 0}, float32(-10)},
		{schema.DataTypeInt64, []uint16{0x8123, 0x4567, 0x89ab, 0xcdef}, int64(-9141386507638288913)},
		{schema.DataTypeUint64, []uint16{0x0102, 0x0304, 0x0506, 0x0708}, uint64(0x0102030405060708)},
		{schema.DataTypeFloat64, []uint16{0xc024, 0, 0, 0}, float64(-10)},
	} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			for _, random := range []bool{false, true} {
				t.Run(fmt.Sprint(tc.dtype, "/", format, "/", random), func(t *testing.T) {
					words := append([]uint16(nil), tc.words...)
					if format == "CDAB" || format == "DCBA" {
						for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
							words[i], words[j] = words[j], words[i]
						}
					}
					if format == "BADC" || format == "DCBA" {
						for i, w := range words {
							words[i] = w>>8 | w<<8
						}
					}
					points := []schema.Point{mcPlannedPoint("scalar", "D0", tc.dtype, format)}
					want := map[string]any{"scalar": tc.want}
					if random {
						points = append(points, mcPlannedPoint("far", "D20", schema.DataTypeInt16, ""))
						words = append(words, 0xfffe)
						want["far"] = int16(-2)
					}
					req := mcPlannedRequest(t, points...)
					require.Len(t, req.Plan.Requests(), 1)
					tr := &mcPlannedTransport{response: mcPlannedWords(words...)}
					results, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
					require.NoError(t, err)
					assertMCPlannedResults(t, req, results, want)
					require.Len(t, tr.frames, 1)
				})
			}
		}
	}
}
func TestMCPlannedLegacyAndSingleWordKeepNativeSemantics(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		points := []schema.Point{mcPlannedPoint("signed", "D0", schema.DataTypeInt16, ""), mcPlannedPoint("unsigned", "D0", schema.DataTypeUint16, ""), mcPlannedPoint("bool", "D0", schema.DataTypeBool, "")}
		if legacy {
			for i := range points {
				points[i].FormatMode = schema.PointFormatLegacy
				points[i].DataFormat = "BADC"
			}
		}
		req := mcPlannedRequest(t, points...)
		tr := &mcPlannedTransport{response: mcPlannedWords(0xfffe)}
		got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
		require.NoError(t, err)
		assertMCPlannedResults(t, req, got, map[string]any{"signed": int16(-2), "unsigned": uint16(0xfffe), "bool": true})
	}
	p := mcPlannedPoint("legacy", "D0", schema.DataTypeUint64, "")
	p.FormatMode = schema.PointFormatLegacy
	p.DataFormat = "ABCD"
	req := mcPlannedRequest(t, p)
	tr := &mcPlannedTransport{response: mcPlannedWords(0x0708, 0x0506, 0x0304, 0x0102)}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"legacy": uint64(0x0102030405060708)})
}
func TestMCPlannedRadixAndAllBitPositions(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("a", "X10", schema.DataTypeBool, ""), mcPlannedPoint("b", "X11", schema.DataTypeBool, ""), mcPlannedPoint("c", "X12", schema.DataTypeBool, ""))
	tr := &mcPlannedTransport{response: []byte{0, 0, 0x10, 0x10}}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"a": true, "b": false, "c": true})
	require.Len(t, tr.frames, 1)
	require.Equal(t, []byte{1, 4, 1, 0, 8, 0, 0, 0x9c, 3, 0}, tr.frames[0][11:])
	req = mcPlannedRequest(t, mcPlannedPoint("word", "W10", schema.DataTypeUint16, ""))
	tr = &mcPlannedTransport{response: mcPlannedWords(0x1234)}
	got, err = mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"word": uint16(0x1234)})
	require.Equal(t, []byte{0x10, 0, 0, 0xb4, 1, 0}, tr.frames[0][15:])
}
func TestMCPlannedRandomNestedOverlap(t *testing.T) {
	req := mcPlannedRequest(t, mcPlannedPoint("wide", "D0", schema.DataTypeInt64, "ABCD"), mcPlannedPoint("small", "D1", schema.DataTypeInt16, ""), mcPlannedPoint("far", "D10", schema.DataTypeInt16, ""))
	require.Len(t, req.Plan.Requests(), 1)
	require.Equal(t, "random", req.Plan.Requests()[0].Kind)
	tr := &mcPlannedTransport{response: mcPlannedWords(0x0102, 0x0304, 0x0506, 0x0708, 9)}
	got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
	require.NoError(t, err)
	assertMCPlannedResults(t, req, got, map[string]any{"wide": int64(0x0102030405060708), "small": int16(0x0304), "far": int16(9)})
	require.Len(t, tr.frames, 1)
	require.Equal(t, byte(5), tr.frames[0][15])
	require.Len(t, tr.frames[0], 37)
}
func TestMCPlannedLimitsExecuteCompleteScalarRequests(t *testing.T) {
	for _, random := range []bool{false, true} {
		points := []schema.Point{}
		count := 126
		if random {
			count = 64
		}
		for i := 0; i < count; i++ {
			dtype := schema.DataTypeInt16
			address := fmt.Sprint("D", i)
			format := ""
			if random {
				dtype = schema.DataTypeInt64
				address = fmt.Sprint("D", i*10)
				format = "ABCD"
			}
			points = append(points, mcPlannedPoint(fmt.Sprint(i), address, dtype, format))
		}
		support := "unknown"
		if random {
			support = "supported"
		}
		req := mcPlannedRequestCapability(t, support, points...)
		require.Len(t, req.Plan.Requests(), 2)
		allBindings := req.Bindings
		seen := map[string]bool{}
		sends := 0
		for index, r := range req.Plan.Requests() {
			req.RequestIndex = index
			req.RequestID = fmt.Sprint("request-", index)
			req.Bindings = nil
			for _, p := range r.Points {
				for _, b := range allBindings {
					if b.Point.ID == p.Point.ID {
						req.Bindings = append(req.Bindings, b)
					}
				}
			}
			words := make([]uint16, r.Length)
			for i := range words {
				words[i] = 1
			}
			tr := &mcPlannedTransport{response: mcPlannedWords(words...)}
			got, err := mcPlannedConnector(tr).ReadPlanned(t.Context(), req)
			require.NoError(t, err)
			require.Len(t, tr.frames, 1)
			sends++
			for _, value := range got {
				require.False(t, seen[value.Binding.Point.ID])
				seen[value.Binding.Point.ID] = true
				if random {
					require.Equal(t, int64(0x0001000100010001), value.Value)
					require.Equal(t, "random", r.Kind)
					require.LessOrEqual(t, int(tr.frames[0][15]), 255)
				} else {
					require.Equal(t, int16(1), value.Value)
					require.LessOrEqual(t, int(binary.LittleEndian.Uint16(tr.frames[0][19:])), 125)
				}
				require.Equal(t, schema.QualityGood, value.Quality)
			}
		}
		require.Len(t, seen, count)
		require.Equal(t, 2, sends)
	}
}
