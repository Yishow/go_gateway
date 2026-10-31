package adapters

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"
)

type modbusPlannedAddress struct {
	function byte
	address  uint16
}
type modbusPlannedPeer struct {
	mu                sync.Mutex
	kind              schema.ProtocolType
	memory            map[modbusPlannedAddress]uint16
	frames, responses [][]byte
	failure           error
	mutate            func([]byte) []byte
	hook              func()
	wait              bool
}

func (*modbusPlannedPeer) Connect() error               { return nil }
func (*modbusPlannedPeer) Close() error                 { return nil }
func (*modbusPlannedPeer) IsConnected() bool            { return true }
func (*modbusPlannedPeer) GetNextTransactionID() uint16 { return 7 }
func (p *modbusPlannedPeer) GetOriginalTransport() any {
	if p.kind == schema.ProtocolModbusRTU {
		return &modbus.RTUTransport{}
	}
	if p.kind == schema.ProtocolModbusUDP {
		return &modbus.UDPTransport{}
	}
	return &modbus.TCPTransport{}
}
func (p *modbusPlannedPeer) SendReceive(frame []byte) ([]byte, error) { return p.respond(frame) }
func (p *modbusPlannedPeer) SendReceiveContext(ctx context.Context, frame []byte) ([]byte, error) {
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	if p.wait {
		p.mu.Lock()
		p.frames = append(p.frames, append([]byte(nil), frame...))
		p.mu.Unlock()
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return p.respond(frame)
}
func (p *modbusPlannedPeer) respond(frame []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.frames = append(p.frames, append([]byte(nil), frame...))
	var tx uint16
	var unit, function byte
	var data []byte
	var err error
	if p.kind == schema.ProtocolModbusRTU {
		unit, function, data, err = modbus.ParseRTUFrame(frame)
	} else {
		tx, unit, function, data, err = modbus.ParseTCPFrame(frame)
	}
	if err != nil || len(data) != 4 || unit != 1 || function < 1 || function > 4 {
		return nil, fmt.Errorf("invalid peer request %x: %v", frame, err)
	}
	address, count := binary.BigEndian.Uint16(data), binary.BigEndian.Uint16(data[2:])
	if count == 0 || int(address)+int(count) > 65536 {
		return nil, fmt.Errorf("invalid peer span")
	}
	if p.hook != nil {
		p.hook()
	}
	if p.failure != nil {
		return nil, p.failure
	}
	size := 2 * int(count)
	if function < 3 {
		size = (int(count) + 7) / 8
	}
	payload := make([]byte, size+1)
	payload[0] = byte(size)
	for i := 0; i < int(count); i++ {
		value, ok := p.memory[modbusPlannedAddress{function, address + uint16(i)}]
		if !ok {
			return nil, fmt.Errorf("peer unreadable FC%02d:%d", function, int(address)+i)
		}
		if function < 3 {
			if value > 1 {
				return nil, fmt.Errorf("invalid fixture bit")
			}
			if value == 1 {
				payload[1+i/8] |= 1 << uint(i%8)
			}
		} else {
			binary.BigEndian.PutUint16(payload[1+2*i:], value)
		}
	}
	var response []byte
	if p.kind == schema.ProtocolModbusRTU {
		response = modbus.BuildRTUFrame(unit, function, payload)
	} else {
		response = modbus.BuildTCPFrame(tx, unit, function, payload)
	}
	if p.mutate != nil {
		response = p.mutate(response)
	}
	p.responses = append(p.responses, append([]byte(nil), response...))
	return response, nil
}
func modbusPlannedProtocol(kind schema.ProtocolType, peer *modbusPlannedPeer) connector.Protocol {
	peer.kind = kind
	client := modbus.NewClient(peer, 1)
	switch kind {
	case schema.ProtocolModbusRTU:
		return &ModbusRTUConnector{client: client, transport: peer, connected: true, persistentMode: true}
	case schema.ProtocolModbusUDP:
		return &ModbusUDPConnector{client: client, transport: peer, connected: true, persistentMode: true}
	default:
		return &ModbusTCPConnector{client: client, transport: peer, connected: true, persistentMode: true}
	}
}
func modbusPlannedPoint(id, address string, dtype schema.DataType, format string) schema.Point {
	p := schema.Point{ID: id, DeviceID: "device", RevisionID: "revision-" + id, Address: address, DataType: dtype, FormatMode: schema.PointFormatInherit, FormatVersion: 1}
	if format != "" {
		p.FormatMode = schema.PointFormatOverride
		p.FormatOverride = format
	}
	return p
}
func modbusPlannedRequest(t *testing.T, kind schema.ProtocolType, coverage optimizer.ReadableCoverage, points ...schema.Point) connector.PlannedReadRequest {
	t.Helper()
	frame := "tcp"
	if kind == schema.ProtocolModbusRTU {
		frame = "rtu"
	}
	source := optimizer.AcquisitionSource{Identity: optimizer.AcquisitionIdentity{DeviceID: "device", DeviceRevision: "device-v1", Station: "1", ConnectionRevision: connector.ReadConnectionRevision("{}"), Protocol: "modbus", Frame: frame, ProtocolRevision: "modbus-v1", PollingGroup: "fast", Cadence: time.Second}, Capability: optimizer.AcquisitionCapability{Revision: "cap-v1", ProtocolLimit: 2000}, Coverage: coverage}
	req := connector.PlannedReadRequest{AcquisitionID: "scan-1", RequestID: "request-1"}
	for _, p := range points {
		profile, err := connector.ResolveReadProfile(p, decoding.FormatContext{Protocol: kind, RuleOwned: true, RuleID: "rule", RuleRevision: "source-v1", ConnectionRevision: connector.ReadConnectionRevision("{}")})
		require.NoError(t, err)
		resolution := profile.Resolution()
		point := optimizer.EffectivePoint{ID: p.ID, SourceIdentity: "rule", SourceRevision: "source-v1", Address: p.Address, Function: p.Function, DataType: p.DataType, NativeWidth: schema.RegisterCountForDataType(p.DataType), Format: resolution.EffectiveFormat, FormatSource: string(resolution.Source), FormatRevision: p.RevisionID}
		source.Points = append(source.Points, point)
		req.Bindings = append(req.Bindings, connector.PlannedPointBinding{Point: point, Profile: profile, PointRevision: p.RevisionID, ConfigFingerprint: "fingerprint-" + p.ID, InstallationScope: "installation", InstallationLease: 19})
	}
	req.Plan = optimizer.BuildAcquisitionPlan(optimizer.AcquisitionInput{ConfigurationRevision: "config-v1", Sources: []optimizer.AcquisitionSource{source}})
	require.Empty(t, req.Plan.Errors())
	return req
}
func modbusPlannedAt(req connector.PlannedReadRequest, index int) connector.PlannedReadRequest {
	req.RequestIndex = index
	req.RequestID = fmt.Sprintf("request-%d", index+1)
	all := req.Bindings
	req.Bindings = nil
	for _, p := range req.Plan.Requests()[index].Points {
		for _, b := range all {
			if p.Point.ID == b.Point.ID {
				req.Bindings = append(req.Bindings, b)
			}
		}
	}
	return req
}
func requireModbusPlannedReader(t *testing.T, c connector.Protocol) connector.PlannedReader {
	t.Helper()
	reader, ok := c.(connector.PlannedReader)
	require.True(t, ok, "production Modbus adapter must expose PlannedReader")
	return reader
}
func assertModbusPlannedResults(t *testing.T, req connector.PlannedReadRequest, got []connector.PlannedPointResult, want map[string]any) {
	t.Helper()
	require.Len(t, got, len(want))
	seen := map[string]bool{}
	for _, r := range got {
		id := r.Binding.Point.ID
		require.False(t, seen[id])
		seen[id] = true
		require.Equal(t, want[id], r.Value, id)
		require.Equal(t, schema.QualityGood, r.Quality, id)
		require.Equal(t, req.AcquisitionID, r.AcquisitionID)
		require.Equal(t, req.RequestID, r.RequestID)
		require.Equal(t, req.Plan.Fingerprint(), r.PlanFingerprint)
		require.Equal(t, req.Plan.Requests()[req.RequestIndex].Group, r.Group)
		for _, b := range req.Bindings {
			if b.Point.ID == id {
				require.Equal(t, b, r.Binding)
			}
		}
		require.False(t, r.Timestamp.IsZero())
		require.Equal(t, got[0].Timestamp, r.Timestamp)
		require.Equal(t, r.Timestamp, r.ObservedAt)
		require.Equal(t, r.Timestamp, r.ReceivedAt)
		require.Equal(t, "gateway", r.TimeOrigin)
	}
}

// The memory peer is an explicit whole-frame transport, unlike a serial stream.
func (p *modbusPlannedPeer) SendReceiveExactContext(ctx context.Context, frame []byte) ([]byte, error) {
	return p.SendReceiveContext(ctx, frame)
}
