package adapters

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"
	"testing"
)

func TestModbusPlannedMalformedWholeBatchIsBadWithoutFanout(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		for _, bit := range []bool{false, true} {
			for _, fault := range []string{"empty", "short", "extra", "count", "wrong-unit", "wrong-function", "exception", "exception-extra", "exception-foreign", "envelope"} {
				t.Run(string(kind)+map[bool]string{false: "word", true: "bit"}[bit]+fault, func(t *testing.T) {
					area := "HR"
					dtype := schema.DataTypeUint16
					function := byte(3)
					if bit {
						area = "C"
						dtype = schema.DataTypeBool
						function = 1
					}
					req := modbusPlannedRequest(t, kind, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", area+"0", dtype, ""), modbusPlannedPoint("b", area+"1", dtype, ""))
					peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{function, 0}: 1, {function, 1}: 0}}
					peer.mutate = func(response []byte) []byte {
						if fault == "empty" {
							return nil
						}
						offset, end := 7, len(response)
						if kind == schema.ProtocolModbusRTU {
							offset = 1
							end -= 2
						}
						unit := byte(1)
						pdu := append([]byte(nil), response[offset:end]...)
						switch fault {
						case "short":
							pdu = pdu[:len(pdu)-1]
						case "extra":
							pdu = append(pdu, 0)
						case "count":
							pdu[1]++
						case "wrong-unit":
							unit = 2
						case "wrong-function":
							pdu[0]++
						case "exception":
							pdu = []byte{function | 0x80, 1}
						case "exception-extra":
							pdu = []byte{function | 0x80, 1, 0}
						case "exception-foreign":
							pdu = []byte{(function + 1) | 0x80, 1}
						}
						var out []byte
						if kind == schema.ProtocolModbusRTU {
							out = modbus.BuildRTUFrame(unit, pdu[0], pdu[1:])
							if fault == "envelope" {
								out[len(out)-1] ^= 1
							}
						} else {
							out = modbus.BuildTCPFrame(7, unit, pdu[0], pdu[1:])
							if fault == "envelope" {
								binary.BigEndian.PutUint16(out[4:6], 1)
							}
						}
						return out
					}
					ctx := wirebudget.With(t.Context(), 4, 0)
					got, err := requireModbusPlannedReader(t, modbusPlannedProtocol(kind, peer)).ReadPlanned(ctx, req)
					require.Error(t, err)
					require.Equal(t, fault == "exception", modbus.IsCommandRejected(err))
					require.Len(t, peer.frames, 1)
					require.Equal(t, 1, wirebudget.From(ctx).Used())
					require.Len(t, got, 2)
					for i, r := range got {
						require.Equal(t, schema.QualityBad, r.Quality)
						require.Nil(t, r.Value)
						require.Nil(t, r.RawBytes)
						require.NotEmpty(t, r.Error)
						require.Equal(t, got[0].Timestamp, r.Timestamp)
						require.Equal(t, req.Plan.Requests()[0].Points[i].Point, r.Binding.Point)
						require.Equal(t, uint64(19), r.Binding.InstallationLease)
					}
				})
			}
		}
	}
}
func TestModbusPlannedFailureDoesNotAlterAnotherRequest(t *testing.T) {
	req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""), modbusPlannedPoint("b", "HR10", schema.DataTypeUint16, ""))
	peer := &modbusPlannedPeer{memory: map[modbusPlannedAddress]uint16{{3, 0}: 7, {3, 10}: 8}}
	reader := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer))
	first := modbusPlannedAt(req, 0)
	got, err := reader.ReadPlanned(t.Context(), first)
	require.NoError(t, err)
	assertModbusPlannedResults(t, first, got, map[string]any{"a": uint16(7)})
	peer.mutate = func([]byte) []byte { return nil }
	bad, err := reader.ReadPlanned(t.Context(), modbusPlannedAt(req, 1))
	require.Error(t, err)
	require.Equal(t, schema.QualityBad, bad[0].Quality)
	assertModbusPlannedResults(t, first, got, map[string]any{"a": uint16(7)})
}
func TestModbusPlannedWrongBindingRejectedBeforeWire(t *testing.T) {
	for _, fault := range []string{"profile", "revision", "format", "bindings", "request-id", "request-index", "foreign-protocol"} {
		t.Run(fault, func(t *testing.T) {
			req := modbusPlannedRequest(t, schema.ProtocolModbusTCP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""))
			switch fault {
			case "profile":
				req.Bindings[0].Profile = connector.ReadProfile{}
			case "revision":
				req.Bindings[0].PointRevision = "new"
			case "format":
				req.Bindings[0].Point.Format = "BADC"
			case "bindings":
				req.Bindings = nil
			case "request-id":
				req.RequestID = ""
			case "request-index":
				req.RequestIndex = 3
			case "foreign-protocol":
				req = modbusPlannedRequest(t, schema.ProtocolModbusUDP, optimizer.ReadableCoverage{}, modbusPlannedPoint("a", "HR0", schema.DataTypeUint16, ""))
			}
			peer := &modbusPlannedPeer{}
			_, err := requireModbusPlannedReader(t, modbusPlannedProtocol(schema.ProtocolModbusTCP, peer)).ReadPlanned(t.Context(), req)
			require.Error(t, err)
			require.Empty(t, peer.frames)
		})
	}
}
