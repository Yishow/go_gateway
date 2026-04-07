package adapters

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModbusConnectors_Read_UsesConnectionConfigDataFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		newConnector func(client *modbus.ModbusClient) connector.Protocol
	}{
		{
			name: "tcp",
			newConnector: func(client *modbus.ModbusClient) connector.Protocol {
				return &ModbusTCPConnector{client: client, connected: true, persistentMode: true}
			},
		},
		{
			name: "rtu",
			newConnector: func(client *modbus.ModbusClient) connector.Protocol {
				return &ModbusRTUConnector{client: client, connected: true, persistentMode: true}
			},
		},
		{
			name: "udp",
			newConnector: func(client *modbus.ModbusClient) connector.Protocol {
				return &ModbusUDPConnector{client: client, connected: true, persistentMode: true}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTrans := NewMockTransport()
			mockTrans.HoldingRegs[0] = 0x0000
			mockTrans.HoldingRegs[1] = 0x3f80
			mockTrans.HoldingRegs[2] = 0x0064
			mockTrans.HoldingRegs[3] = 0x0000

			conn := tt.newConnector(modbus.NewClient(mockTrans, 1))
			loadConnectionConfigDataFormat(t, conn, "CDAB")

			floatResult, err := conn.Read(context.Background(), connector.ReadRequest{
				Address:  "40001",
				DataType: schema.DataTypeFloat32,
			})
			require.NoError(t, err)
			assert.Equal(t, float32(1), floatResult.Value)

			intResult, err := conn.Read(context.Background(), connector.ReadRequest{
				Address:  "40003",
				DataType: schema.DataTypeInt32,
			})
			require.NoError(t, err)
			assert.Equal(t, int32(100), intResult.Value)
		})
	}
}

func loadConnectionConfigDataFormat(t *testing.T, protocol connector.Protocol, dataFormat string) {
	t.Helper()

	rawConfig, err := json.Marshal(map[string]string{"data_format": dataFormat})
	require.NoError(t, err)

	var configValue reflect.Value
	switch typed := protocol.(type) {
	case *ModbusTCPConnector:
		require.NoError(t, json.Unmarshal(rawConfig, &typed.config))
		configValue = reflect.ValueOf(typed.config)
	case *ModbusRTUConnector:
		require.NoError(t, json.Unmarshal(rawConfig, &typed.config))
		configValue = reflect.ValueOf(typed.config)
	case *ModbusUDPConnector:
		require.NoError(t, json.Unmarshal(rawConfig, &typed.config))
		configValue = reflect.ValueOf(typed.config)
	default:
		t.Fatalf("unsupported protocol type %T", protocol)
	}

	dataFormatField := configValue.FieldByName("DataFormat")
	require.True(t, dataFormatField.IsValid(), "connection config should expose DataFormat")
	assert.Equal(t, dataFormat, dataFormatField.String())
}
