package adapters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
)

func TestFatekPlannedMixedNativeWidthsAndPerPointFormats(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("abcd", "D0", schema.DataTypeUint32, "ABCD"), fatekPlannedPoint("cdab", "D0", schema.DataTypeUint32, "CDAB"), fatekPlannedPoint("dr", "DR10", schema.DataTypeFloat32, "ABCD"), fatekPlannedPoint("bit", "X20", schema.DataTypeBool, ""))
	require.Len(t, req.Plan.Requests(), 1)
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 0x0102, "D1": 0x0304, "DR10": 0xc1200000, "X20": true}}
	got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
	require.NoError(t, err)
	assertFatekPlannedResults(t, req, got, map[string]any{"abcd": uint32(0x01020304), "cdab": uint32(0x03040102), "dr": float32(-10), "bit": true})
	require.Len(t, peer.frames, 1)
	require.Equal(t, "04D00000D00001DR00010X0020", string(peer.frames[0][5:len(peer.frames[0])-3]))
	require.Equal(t, optimizer.ApplicationFrameBudget{RequestBytes: 34, ResponseBytes: 26}, req.Plan.Requests()[0].FrameBytes)
	require.Len(t, peer.frames[0], 34)
	require.Len(t, peer.responses[0], 26)
}
func TestFatekPlannedDenseOverlapAndAllBits(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("wide", "D0", schema.DataTypeUint64, "ABCD"), fatekPlannedPoint("small", "D1", schema.DataTypeUint16, ""))
	require.Equal(t, "contiguous", req.Plan.Requests()[0].Kind)
	peer := &fatekPlannedPeer{memory: map[string]any{"D0": 0x0102, "D1": 0x0304, "D2": 0x0506, "D3": 0x0708}}
	got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
	require.NoError(t, err)
	assertFatekPlannedResults(t, req, got, map[string]any{"wide": uint64(0x0102030405060708), "small": uint16(0x0304)})
	require.Equal(t, "46", string(peer.frames[0][3:5]))
	require.Equal(t, "04D00000", string(peer.frames[0][5:len(peer.frames[0])-3]))
	for _, random := range []bool{false, true} {
		t.Run(fmt.Sprint(random), func(t *testing.T) {
			addresses := []string{"X0", "X1", "X2"}
			if random {
				addresses = []string{"X0", "X10", "X20"}
			}
			req := fatekPlannedRequest(t, fatekPlannedPoint("a", addresses[0], schema.DataTypeBool, ""), fatekPlannedPoint("b", addresses[1], schema.DataTypeBool, ""), fatekPlannedPoint("c", addresses[2], schema.DataTypeBool, ""))
			peer := &fatekPlannedPeer{memory: map[string]any{addresses[0]: true, addresses[1]: false, addresses[2]: true}}
			got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
			require.NoError(t, err)
			assertFatekPlannedResults(t, req, got, map[string]any{"a": true, "b": false, "c": true})
			command := "44"
			if random {
				command = "48"
			}
			require.Equal(t, command, string(peer.frames[0][3:5]))
		})
	}
}
func TestFatekPlannedDScalarNumericMatrix(t *testing.T) {
	for _, tc := range []struct {
		dtype schema.DataType
		words []int
		want  any
	}{
		{schema.DataTypeInt32, []int{0x8123, 0x4567}, int32(-2128394905)},
		{schema.DataTypeUint32, []int{0x0102, 0x0304}, uint32(0x01020304)},
		{schema.DataTypeFloat32, []int{0xc120, 0}, float32(-10)},
		{schema.DataTypeInt64, []int{0x8123, 0x4567, 0x89ab, 0xcdef}, int64(-9141386507638288913)},
		{schema.DataTypeUint64, []int{0x0102, 0x0304, 0x0506, 0x0708}, uint64(0x0102030405060708)},
		{schema.DataTypeFloat64, []int{0xc024, 0, 0, 0}, float64(-10)},
	} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			t.Run(fmt.Sprint(tc.dtype, "/", format), func(t *testing.T) {
				words := append([]int(nil), tc.words...)
				if format == "CDAB" || format == "DCBA" {
					for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
						words[i], words[j] = words[j], words[i]
					}
				}
				if format == "BADC" || format == "DCBA" {
					for i, w := range words {
						words[i] = (w >> 8) | ((w & 255) << 8)
					}
				}
				req := fatekPlannedRequest(t, fatekPlannedPoint("scalar", "D0", tc.dtype, format), fatekPlannedPoint("far", "D20", schema.DataTypeInt16, ""))
				memory := map[string]any{"D20": 0xfffe}
				for i, w := range words {
					memory[fmt.Sprint("D", i)] = w
				}
				peer := &fatekPlannedPeer{memory: memory}
				got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
				require.NoError(t, err)
				assertFatekPlannedResults(t, req, got, map[string]any{"scalar": tc.want, "far": int16(-2)})
			})
		}
	}
}
func TestFatekPlannedDRScalarNumericMatrix(t *testing.T) {
	for _, tc := range []struct {
		dtype schema.DataType
		words []int
		want  any
	}{
		{schema.DataTypeInt32, []int{0x8123, 0x4567}, int32(-2128394905)},
		{schema.DataTypeUint32, []int{0x0102, 0x0304}, uint32(0x01020304)},
		{schema.DataTypeFloat32, []int{0xc120, 0}, float32(-10)},
		{schema.DataTypeInt64, []int{0x8123, 0x4567, 0x89ab, 0xcdef}, int64(-9141386507638288913)},
		{schema.DataTypeUint64, []int{0x0102, 0x0304, 0x0506, 0x0708}, uint64(0x0102030405060708)},
		{schema.DataTypeFloat64, []int{0xc024, 0, 0, 0}, float64(-10)},
	} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			t.Run(fmt.Sprint(tc.dtype, "/", format), func(t *testing.T) {
				words := append([]int(nil), tc.words...)
				if format == "CDAB" || format == "DCBA" {
					for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
						words[i], words[j] = words[j], words[i]
					}
				}
				if format == "BADC" || format == "DCBA" {
					for i, w := range words {
						words[i] = (w >> 8) | ((w & 255) << 8)
					}
				}
				req := fatekPlannedRequest(t, fatekPlannedPoint("scalar", "DR0", tc.dtype, format), fatekPlannedPoint("far", "D20", schema.DataTypeInt16, ""))
				memory := map[string]any{"D20": 0xfffe}
				for i := 0; i < len(words); i += 2 {
					memory[fmt.Sprint("DR", i)] = words[i]<<16 | words[i+1]
				}
				peer := &fatekPlannedPeer{memory: memory}
				got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
				require.NoError(t, err)
				assertFatekPlannedResults(t, req, got, map[string]any{"scalar": tc.want, "far": int16(-2)})
			})
		}
	}
}
func TestFatekPlannedLegacyKeepsNativeComponents(t *testing.T) {
	for _, address := range []string{"D0", "DR0"} {
		p := fatekPlannedPoint("legacy", address, schema.DataTypeUint32, "")
		p.FormatMode = schema.PointFormatLegacy
		p.DataFormat = "DCBA"
		req := fatekPlannedRequest(t, p)
		peer := &fatekPlannedPeer{memory: map[string]any{"D0": 0x0102, "D1": 0x0304, "DR0": 0x01020304}}
		got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
		require.NoError(t, err)
		assertFatekPlannedResults(t, req, got, map[string]any{"legacy": uint32(0x01020304)})
	}
}
func fatekSelectRequest(req connector.PlannedReadRequest, index int) connector.PlannedReadRequest {
	all := req.Bindings
	req.Bindings = nil
	req.RequestIndex = index
	req.RequestID = fmt.Sprint("request-", index)
	request, _ := req.Plan.Request(index)
	for _, p := range request.Points {
		for _, b := range all {
			if b.Point.ID == p.Point.ID {
				req.Bindings = append(req.Bindings, b)
			}
		}
	}
	return req
}
func TestFatekPlanned65DUnitsSplitWithoutLostPoints(t *testing.T) {
	var points []schema.Point
	memory := map[string]any{}
	for i := range 65 {
		address := fmt.Sprint("D", i*10)
		points = append(points, fatekPlannedPoint(fmt.Sprint(i), address, schema.DataTypeUint16, ""))
		memory[address] = i + 1
	}
	req := fatekPlannedRequest(t, points...)
	require.Len(t, req.Plan.Requests(), 2)
	peer := &fatekPlannedPeer{memory: memory}
	c := fatekPlannedConnector(peer)
	seen := map[string]bool{}
	for i := range req.Plan.Requests() {
		selected := fatekSelectRequest(req, i)
		got, err := readFatekPlanned(t, t.Context(), c, selected)
		require.NoError(t, err)
		for _, value := range got {
			require.Equal(t, schema.QualityGood, value.Quality)
			require.Equal(t, uint16(memory[value.Binding.Point.Address].(int)), value.Value)
			require.False(t, seen[value.Binding.Point.ID])
			seen[value.Binding.Point.ID] = true
		}
	}
	require.Len(t, seen, 65)
	require.Len(t, peer.frames, 2)
}
