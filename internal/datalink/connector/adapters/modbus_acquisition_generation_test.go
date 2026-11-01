package adapters

import (
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"net"
	"testing"

	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
)

type modbusGenerationPeer struct{ acquisitionGenerationTransport }

func (p *modbusGenerationPeer) IsConnected() bool { return p.live }

func TestModbusAcquisitionGenerationReconnectAndDelegation(t *testing.T) {
	for _, kind := range []string{"tcp", "udp", "rtu"} {
		t.Run(kind, func(t *testing.T) {
			peer := &modbusGenerationPeer{}
			client := modbus.NewClient(peer, 1)
			var adapter any
			var reconnect func(context.Context) error
			var ensure func() error
			var proxy func() *ModbusTCPConnector
			switch kind {
			case "tcp":
				c := &ModbusTCPConnector{client: client, transport: peer}
				adapter = c
				reconnect = c.Reconnect
				ensure = c.ensureConnection
				proxy = func() *ModbusTCPConnector { return c }
			case "udp":
				c := &ModbusUDPConnector{client: client, transport: peer}
				adapter = c
				reconnect = c.Reconnect
				ensure = c.ensureConnection
				proxy = c.readProxy
			case "rtu":
				c := &ModbusRTUConnector{client: client, transport: peer}
				adapter = c
				reconnect = c.Reconnect
				ensure = c.ensureConnection
				proxy = c.readProxy
			}
			provider, ok := adapter.(interface{ AcquisitionTransportGeneration() uint64 })
			require.True(t, ok, "real Modbus adapter must expose current proof without I/O")
			require.Zero(t, provider.AcquisitionTransportGeneration())
			require.NoError(t, ensure())
			first := provider.AcquisitionTransportGeneration()
			require.NotZero(t, first)
			require.NoError(t, reconnect(t.Context()))
			second := provider.AcquisitionTransportGeneration()
			require.Greater(t, second, first)
			failure := errors.New("dial failed")
			peer.connect = func(context.Context) error { return failure }
			require.ErrorIs(t, reconnect(t.Context()), failure)
			require.Equal(t, second, provider.AcquisitionTransportGeneration())
			ctx, cancel := context.WithCancel(t.Context())
			peer.connect = func(context.Context) error { cancel(); return nil }
			require.ErrorIs(t, reconnect(ctx), context.Canceled)
			require.Equal(t, second, provider.AcquisitionTransportGeneration())
			peer.connect = nil
			require.NoError(t, proxy().ensureConnectionContext(t.Context()))
			third := provider.AcquisitionTransportGeneration()
			require.Greater(t, third, second, "proxy reconnect must advance owner, not temporary proxy")
			peer.live = true
			require.NoError(t, proxy().ensureConnectionContext(t.Context()))
			require.Equal(t, third, provider.AcquisitionTransportGeneration(), "liveness check must not mint a generation")
		})
	}
}

func TestModbusAcquisitionGenerationInitialConnect(t *testing.T) {
	var listen net.ListenConfig
	listener, err := listen.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	udp, err := listen.ListenPacket(t.Context(), "udp", "127.0.0.1:0")
	require.NoError(t, err)
	defer udp.Close()
	peer := &modbusGenerationPeer{}
	for _, tc := range []struct {
		name, config string
		adapter      connector.Protocol
	}{
		{"tcp", fmt.Sprintf(`{"host":"127.0.0.1","port":%d,"timeout":1}`, listener.Addr().(*net.TCPAddr).Port), &ModbusTCPConnector{}},
		{"udp", fmt.Sprintf(`{"host":"127.0.0.1","port":%d,"timeout":1}`, udp.LocalAddr().(*net.UDPAddr).Port), &ModbusUDPConnector{}},
		{"rtu", `{"port":"fixture"}`, &ModbusRTUConnector{transportFactory: func(schema.ConnectionConfigModbusRTU) modbus.Transport { return peer }}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			provider, ok := tc.adapter.(interface{ AcquisitionTransportGeneration() uint64 })
			require.True(t, ok)
			require.Zero(t, provider.AcquisitionTransportGeneration())
			require.NoError(t, tc.adapter.Connect(t.Context(), tc.config))
			defer tc.adapter.Close()
			require.NotZero(t, provider.AcquisitionTransportGeneration())
		})
	}
}
