package adapters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestFatekPlannedSliceRejectsWrongPhysicalUnit(t *testing.T) {
	req := fatekPlannedRequest(t, fatekPlannedPoint("a", "DR0", schema.DataTypeUint64, "ABCD"), fatekPlannedPoint("far", "DR10", schema.DataTypeUint32, "ABCD"))
	r, bindings, err := req.Snapshot()
	require.NoError(t, err)
	p := r.Points[0]
	r.Units[p.Offset+1].Address = 1 // DR1 is valid, but is not this scalar's second component.
	value, _, err := decodeFatekPlannedPoint(r, p, bindings[0].Profile, fatekPlannedRaw{random: map[string]any{"DR0": 0x01020304, "DR1": 0x05060708, "DR10": 7}})
	require.Error(t, err)
	require.Nil(t, value)
}
func TestFatekPlannedPhysicalStrideAndPhase(t *testing.T) {
	for _, tc := range []struct {
		name   string
		points []schema.Point
		memory map[string]any
		want   map[string]any
		body   string
	}{
		{"DR64", []schema.Point{fatekPlannedPoint("a", "DR0", schema.DataTypeUint64, "ABCD"), fatekPlannedPoint("far", "DR10", schema.DataTypeUint64, "ABCD")}, map[string]any{"DR0": 0x01020304, "DR2": 0x05060708, "DR10": 0x11121314, "DR12": 0x15161718}, map[string]any{"a": uint64(0x0102030405060708), "far": uint64(0x1112131415161718)}, "04DR00000DR00002DR00010DR00012"},
		{"cross-phase", []schema.Point{fatekPlannedPoint("a", "DR0", schema.DataTypeUint64, "ABCD"), fatekPlannedPoint("odd", "DR1", schema.DataTypeUint32, "ABCD"), fatekPlannedPoint("b", "DR2", schema.DataTypeUint64, "ABCD")}, map[string]any{"DR0": 0x03040102, "DR1": 0x05060304, "DR2": 0x07080506, "DR4": 0x0b0c090a}, map[string]any{"a": uint64(0x0304010207080506), "odd": uint32(0x05060304), "b": uint64(0x070805060b0c090a)}, "04DR00000DR00002DR00004DR00001"},
		{"W16", []schema.Point{fatekPlannedPoint("a", "WY8", schema.DataTypeUint32, ""), fatekPlannedPoint("far", "WY64", schema.DataTypeUint16, "")}, map[string]any{"WY8": 0x0102, "WY24": 0x0304, "WY64": 7}, map[string]any{"a": uint32(0x01020304), "far": uint16(7)}, "03WY0064WY0008WY0024"},
		{"DW32", []schema.Point{fatekPlannedPoint("a", "DWM8", schema.DataTypeUint64, "")}, map[string]any{"DWM8": 0x01020304, "DWM40": 0x05060708}, map[string]any{"a": uint64(0x0102030405060708)}, "02DWM0008"},
		{"DR-contiguous", []schema.Point{fatekPlannedPoint("a", "DR0", schema.DataTypeUint32, "ABCD"), fatekPlannedPoint("b", "DR2", schema.DataTypeUint32, "ABCD")}, map[string]any{"DR0": 0x01020304, "DR2": 0x05060708}, map[string]any{"a": uint32(0x01020304), "b": uint32(0x05060708)}, "02DR00000"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := fatekPlannedRequest(t, tc.points...)
			require.Len(t, req.Plan.Requests(), 1)
			peer := &fatekPlannedPeer{memory: tc.memory}
			got, err := readFatekPlanned(t, t.Context(), fatekPlannedConnector(peer), req)
			require.NoError(t, err)
			assertFatekPlannedResults(t, req, got, tc.want)
			require.Len(t, peer.frames, 1)
			require.Equal(t, tc.body, string(peer.frames[0][5:len(peer.frames[0])-3]))
			r := req.Plan.Requests()[0]
			require.Equal(t, r.FrameBytes.RequestBytes, len(peer.frames[0]))
			require.Equal(t, r.FrameBytes.ResponseBytes, len(peer.responses[0]))
		})
	}
}
func TestFatekPlannedWordCostWholeScalarLimits(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		dtype                 schema.DataType
		count, bits, requests int
	}{{"DR32", schema.DataTypeUint32, 32, 0, 1}, {"DR33", schema.DataTypeUint32, 33, 0, 2}, {"DR33-wide", schema.DataTypeUint64, 33, 0, 3}, {"mixed44", schema.DataTypeUint32, 20, 24, 1}, {"mixed45", schema.DataTypeUint32, 20, 25, 2}} {
		t.Run(tc.name, func(t *testing.T) {
			var points []schema.Point
			memory := map[string]any{}
			want := map[string]any{}
			for i := range tc.count {
				id, address := fmt.Sprint("dr", i), fmt.Sprint("DR", i*10)
				points = append(points, fatekPlannedPoint(id, address, tc.dtype, "ABCD"))
				memory[address] = 0x01020304
				want[id] = uint32(0x01020304)
				if tc.dtype == schema.DataTypeUint64 {
					memory[fmt.Sprint("DR", i*10+2)] = 0x05060708
					want[id] = uint64(0x0102030405060708)
				}
			}
			for i := range tc.bits {
				id, address := fmt.Sprint("bit", i), fmt.Sprint("X", i*10)
				points = append(points, fatekPlannedPoint(id, address, schema.DataTypeBool, ""))
				memory[address] = i%2 == 0
				want[id] = i%2 == 0
			}
			req := fatekPlannedRequest(t, points...)
			require.Len(t, req.Plan.Requests(), tc.requests)
			peer := &fatekPlannedPeer{memory: memory}
			c := fatekPlannedConnector(peer)
			seen := map[string]bool{}
			for i := range req.Plan.Requests() {
				selected := fatekSelectRequest(req, i)
				got, err := readFatekPlanned(t, t.Context(), c, selected)
				require.NoError(t, err)
				for _, r := range got {
					require.Equal(t, schema.QualityGood, r.Quality)
					require.Equal(t, want[r.Binding.Point.ID], r.Value)
					require.False(t, seen[r.Binding.Point.ID])
					seen[r.Binding.Point.ID] = true
				}
			}
			require.Len(t, seen, tc.count+tc.bits)
			require.Len(t, peer.frames, tc.requests)
			if tc.name == "mixed44" {
				require.Equal(t, "2C", string(peer.frames[0][5:7]))
			}
		})
	}
}
