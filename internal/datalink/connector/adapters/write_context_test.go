package adapters

import (
	"context"
	"io"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
)

// Only in-memory wire calls are counted. No endpoint is opened by this fixture.
type canceledWriteProbe struct {
	kind   schema.ProtocolType
	sends  int
	cancel context.CancelFunc
	ack    bool
}

func (*canceledWriteProbe) Connect() error               { return nil }
func (*canceledWriteProbe) Close() error                 { return nil }
func (*canceledWriteProbe) IsConnected() bool            { return true }
func (*canceledWriteProbe) GetNextTransactionID() uint16 { return 1 }
func (p *canceledWriteProbe) GetOriginalTransport() any {
	if p.kind == schema.ProtocolModbusRTU {
		return &modbus.RTUTransport{}
	}
	return p
}
func (p *canceledWriteProbe) SendReceive(frame []byte) ([]byte, error) {
	p.sends++
	if p.ack {
		switch p.kind {
		case schema.ProtocolFatekFBs:
			return fatekPeerFrame(string(frame[1:3]), string(frame[3:5]), "0"), nil
		case schema.ProtocolMC3E:
			return []byte{0, 0}, nil
		default:
			return append([]byte(nil), frame...), nil
		}
	}
	if p.cancel != nil {
		p.cancel()
	}
	// The request crossed the wire boundary, but no ACK was observed.
	return nil, io.EOF
}
func canceledWriteAdapter(kind schema.ProtocolType, p *canceledWriteProbe) (connector.Protocol, connector.WriteRequest) {
	p.kind = kind
	req := connector.WriteRequest{Address: "D0", Value: 1}
	switch kind {
	case schema.ProtocolFatekFBs:
		return &FatekConnector{client: fatek.NewClient(p, 1), transport: p, connected: true, persistentMode: true}, req
	case schema.ProtocolMC3E:
		return &MC3EConnector{client: mcprotocol.NewClientWithTransport(p), connected: true, persistentMode: true}, req
	case schema.ProtocolModbusRTU:
		req.Address = "HR0"
		return &ModbusRTUConnector{client: modbus.NewClient(p, 1), transport: p, connected: true, persistentMode: true}, req
	case schema.ProtocolModbusUDP:
		req.Address = "HR0"
		return &ModbusUDPConnector{client: modbus.NewClient(p, 1), transport: p, connected: true, persistentMode: true}, req
	default:
		req.Address = "HR0"
		return &ModbusTCPConnector{client: modbus.NewClient(p, 1), transport: p, connected: true, persistentMode: true}, req
	}
}
func TestAdapterWriteCanceledBeforeWireDoesNotSend(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolFatekFBs, schema.ProtocolMC3E, schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			p := &canceledWriteProbe{}
			adapter, req := canceledWriteAdapter(kind, p)
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			require.ErrorIs(t, adapter.Write(ctx, req), context.Canceled)
			require.Zero(t, p.sends)
		})
	}
}
func TestAdapterWriteCanceledAfterUnknownEffectDoesNotRetry(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolFatekFBs, schema.ProtocolMC3E, schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			p := &canceledWriteProbe{cancel: cancel}
			adapter, req := canceledWriteAdapter(kind, p)
			require.ErrorIs(t, adapter.Write(ctx, req), io.EOF, "unknown wire outcome must not become successful")
			require.Equal(t, 1, p.sends, "caller cancellation after the first unknown effect must prevent another write")
		})
	}
}

func TestAdapterWriteAcknowledgedSuccessRemainsSingleSend(t *testing.T) {
	for _, kind := range []schema.ProtocolType{schema.ProtocolFatekFBs, schema.ProtocolMC3E, schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP} {
		t.Run(string(kind), func(t *testing.T) {
			p := &canceledWriteProbe{ack: true}
			adapter, req := canceledWriteAdapter(kind, p)
			require.NoError(t, adapter.Write(t.Context(), req))
			require.Equal(t, 1, p.sends, "known ACK must remain successful without retry")
		})
	}
}
