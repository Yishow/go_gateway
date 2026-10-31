package adapters

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"testing"
)

func TestModbusRTURawFactoryKeepsProductionConnectDefaults(t *testing.T) {
	var captured schema.ConnectionConfigModbusRTU
	calls := 0
	c := NewModbusRTUConnectorWithTransportFactory(func(config schema.ConnectionConfigModbusRTU) modbus.Transport {
		calls++
		captured = config
		return &profileWordTransport{words: []uint16{0}}
	})
	require.NoError(t, c.Connect(context.Background(), `{"serial_port":"raw-fixture"}`))
	t.Cleanup(func() { require.NoError(t, c.Close()) })
	require.Equal(t, 1, calls)
	require.Equal(t, "raw-fixture", captured.SerialPort)
	require.Equal(t, 9600, captured.BaudRate)
	require.Equal(t, 8, captured.DataBits)
	require.Equal(t, 1, captured.StopBits)
	require.Equal(t, "none", captured.Parity)
	require.Equal(t, byte(1), captured.SlaveID)
	require.True(t, c.IsConnected())
}

func TestModbusRTURawFactoryDefaultAndMissingTransport(t *testing.T) {
	original := NewModbusRTUConnector().(*ModbusRTUConnector)
	original.config = schema.ConnectionConfigModbusRTU{SerialPort: "never-opened", BaudRate: 19200, SlaveID: 7}
	transport, ok := original.newTransport().(*modbus.RTUTransport)
	require.True(t, ok, "default registry must retain original physical transport")
	require.Equal(t, "never-opened", transport.Port)
	require.Equal(t, 19200, transport.BaudRate)
	missing := NewModbusRTUConnectorWithTransportFactory(func(schema.ConnectionConfigModbusRTU) modbus.Transport { return nil })
	require.Error(t, missing.Connect(context.Background(), `{}`))
	require.False(t, missing.IsConnected())
	nilFactory := NewModbusRTUConnectorWithTransportFactory(nil)
	require.Nil(t, nilFactory.(*ModbusRTUConnector).newTransport(), "explicit missing factory must not open a physical port")
	require.Error(t, nilFactory.Connect(context.Background(), `{}`))
}
