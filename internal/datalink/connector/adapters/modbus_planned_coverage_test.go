package adapters

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModbusPlannedSafeCoverageAndNoSourceRuleGapGrant(t *testing.T) {
	for _, tc := range []struct {
		name     string
		coverage optimizer.ReadableCoverage
		spans    [][2]int
	}{
		{"same-rule-no-grant", optimizer.ReadableCoverage{}, [][2]int{{0, 1}, {10, 1}, {20, 1}}},
		{"explicit", optimizer.ReadableCoverage{Revision: "coverage-v1", Readable: []optimizer.NativeRange{{Area: "HR", Start: 0, End: 21}}}, [][2]int{{0, 21}}},
		{"forbidden", optimizer.ReadableCoverage{Revision: "coverage-v1", Readable: []optimizer.NativeRange{{Area: "HR", Start: 0, End: 21}}, Forbidden: []optimizer.NativeRange{{Area: "HR", Start: 5, End: 6}}}, [][2]int{{0, 1}, {10, 11}}},
		{"excluded", optimizer.ReadableCoverage{Revision: "coverage-v1", Readable: []optimizer.NativeRange{{Area: "HR", Start: 0, End: 21}}, Excluded: []optimizer.NativeRange{{Area: "HR", Start: 5, End: 6}}}, [][2]int{{0, 1}, {10, 11}}},
		{"unversioned", optimizer.ReadableCoverage{Readable: []optimizer.NativeRange{{Area: "HR", Start: 0, End: 21}}}, [][2]int{{0, 1}, {10, 1}, {20, 1}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, tc.coverage, modbusPlannedPoint("a", "HR0", schema.DataTypeInt16, ""), modbusPlannedPoint("b", "HR10", schema.DataTypeInt16, ""), modbusPlannedPoint("c", "HR20", schema.DataTypeInt16, ""))
			require.Len(t, req.Plan.Requests(), len(tc.spans))
			memory := map[modbusPlannedAddress]uint16{}
			for _, span := range tc.spans {
				for i := span[0]; i < span[0]+span[1]; i++ {
					memory[modbusPlannedAddress{3, uint16(i)}] = uint16(i + 11)
				}
			}
			peer := &modbusPlannedPeer{memory: memory}
			reader := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer))
			seen := map[string]any{}
			for i, r := range req.Plan.Requests() {
				require.Equal(t, tc.spans[i], [2]int{r.Start, r.Length})
				sub := modbusPlannedAt(req, i)
				got, err := reader.ReadPlanned(t.Context(), sub)
				require.NoError(t, err)
				for _, p := range got {
					require.NotContains(t, seen, p.Binding.Point.ID)
					require.Equal(t, schema.QualityGood, p.Quality)
					seen[p.Binding.Point.ID] = p.Value
				}
			}
			require.Equal(t, map[string]any{"a": int16(11), "b": int16(21), "c": int16(31)}, seen)
			require.Len(t, peer.frames, len(tc.spans))
		})
	}
}
func TestModbusPlannedSeparatesFunctionsAndUsesNormalizedOffsets(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			points := []schema.Point{modbusPlannedPoint("coil", "00001", schema.DataTypeBool, ""), modbusPlannedPoint("di", "10001", schema.DataTypeBool, ""), modbusPlannedPoint("holding", "40001", schema.DataTypeUint16, ""), modbusPlannedPoint("input", "30001", schema.DataTypeUint16, "")}
			// Alias address wins conflicting Function just as the existing scalar path.
			points[2].Function = "04"
			req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, points...)
			require.Len(t, req.Plan.Requests(), 4)
			peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{1, 0}: 1, {2, 0}: 0, {3, 0}: 0x1234, {4, 0}: 0x5678}}
			reader := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer))
			want := map[string]any{"coil": true, "di": false, "holding": uint16(0x1234), "input": uint16(0x5678)}
			for i := range req.Plan.Requests() {
				sub := modbusPlannedAt(req, i)
				got, err := reader.ReadPlanned(t.Context(), sub)
				require.NoError(t, err)
				assertModbusPlannedResults(t, sub, got, map[string]any{sub.Bindings[0].Point.ID: want[sub.Bindings[0].Point.ID]})
			}
			require.Len(t, peer.frames, 4)
		})
	}
}
func TestModbusPlannedLimitsKeepWideScalarsWhole(t *testing.T) {
	for _, wide := range []bool{false, true} {
		t.Run(fmt.Sprint(wide), func(t *testing.T) {
			var points []schema.Point
			memory := map[modbusPlannedAddress]uint16{}
			count := 126
			for i := 0; i < count; i++ {
				dtype := schema.DataTypeUint16
				format := ""
				if wide && i == 124 {
					dtype = schema.DataTypeUint64
					format = "ABCD"
				}
				points = append(points, modbusPlannedPoint(fmt.Sprint(i), fmt.Sprint("HR", i), dtype, format))
				memory[modbusPlannedAddress{3, uint16(i)}] = 1
			}
			if wide {
				memory[modbusPlannedAddress{3, 126}] = 1
				memory[modbusPlannedAddress{3, 127}] = 1
			}
			req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, points...)
			require.Len(t, req.Plan.Requests(), 2)
			peer := &modbusPlannedPeer{memory: memory}
			reader := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer))
			seen := map[string]bool{}
			for i, r := range req.Plan.Requests() {
				require.LessOrEqual(t, r.Length, 125)
				sub := modbusPlannedAt(req, i)
				got, err := reader.ReadPlanned(t.Context(), sub)
				require.NoError(t, err)
				for _, p := range got {
					id := p.Binding.Point.ID
					require.False(t, seen[id])
					seen[id] = true
					want := any(uint16(1))
					if wide && id == "124" {
						want = uint64(0x0001000100010001)
					}
					require.Equal(t, want, p.Value)
					require.Equal(t, schema.QualityGood, p.Quality)
				}
			}
			require.Len(t, seen, count)
			require.Len(t, peer.frames, 2)
		})
	}
}
