package adapters

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestModbusPlannedBitCoverageEmitsOnlyOriginalPoints(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, grant := range []bool{false, true} {
			for _, forbid := range []bool{false, true} {
				t.Run(fmt.Sprint(kind, "/", grant, "/", forbid), func(t *testing.T) {
					coverage := optimizer.ReadableCoverage{}
					if grant {
						coverage = optimizer.ReadableCoverage{Revision: "bit-coverage", Readable: []optimizer.NativeRange{{Area: "COIL", Start: 0, End: 21}}}
					}
					if forbid {
						coverage.Forbidden = []optimizer.NativeRange{{Area: "COIL", Start: 5, End: 6}}
					}
					req := modbusPlannedRequest(t, kind, coverage, modbusPlannedPoint("a", "C0", schema.DataTypeBool, ""), modbusPlannedPoint("b", "C10", schema.DataTypeBool, ""), modbusPlannedPoint("c", "C20", schema.DataTypeBool, ""))
					memory := map[modbusPlannedAddress]uint16{}
					for i := 0; i < 21; i++ {
						if forbid && i == 5 {
							continue
						}
						memory[modbusPlannedAddress{1, uint16(i)}] = 0
					}
					memory[modbusPlannedAddress{1, 0}] = 1
					memory[modbusPlannedAddress{1, 20}] = 1
					peer := &modbusPlannedPeer{memory: memory}
					reader := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer))
					seen := map[string]any{}
					for i := range req.Plan.Requests() {
						got, err := reader.ReadPlanned(t.Context(), modbusPlannedAt(req, i))
						require.NoError(t, err)
						for _, r := range got {
							require.NotContains(t, seen, r.Binding.Point.ID)
							require.Equal(t, schema.QualityGood, r.Quality)
							seen[r.Binding.Point.ID] = r.Value
						}
					}
					count := 3
					if grant {
						count = 1
						if forbid {
							count = 2
						}
					}
					require.Len(t, peer.frames, count)
					require.Equal(t, map[string]any{"a": true, "b": false, "c": true}, seen)
				})
			}
		}
	}
}
func TestModbusPlannedBitMaximumAndEndAddress(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			coverage := optimizer.ReadableCoverage{Revision: "last-bits", Readable: []optimizer.NativeRange{{Area: "DI", Start: 63536, End: 65536}}}
			req := modbusPlannedRequest(t, kind, coverage, modbusPlannedPoint("first", "DI63536", schema.DataTypeBool, ""), modbusPlannedPoint("last", "DI65535", schema.DataTypeBool, ""))
			require.Len(t, req.Plan.Requests(), 1)
			require.Equal(t, 2000, req.Plan.Requests()[0].Length)
			memory := map[modbusPlannedAddress]uint16{}
			for i := 63536; i <= 65535; i++ {
				memory[modbusPlannedAddress{2, uint16(i)}] = 0
			}
			memory[modbusPlannedAddress{2, 65535}] = 1
			peer := &modbusPlannedPeer{memory: memory}
			got, err := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer)).ReadPlanned(t.Context(), req)
			require.NoError(t, err)
			assertModbusPlannedResults(t, req, got, map[string]any{"first": false, "last": true})
			require.Equal(t, req.Plan.Requests()[0].FrameBytes.ResponseBytes, len(peer.responses[0]))
		})
	}
}
