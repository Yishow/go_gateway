package handlers

import (
	"math"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/require"
)

func TestModbusWriteRejectsTruncatedRegisterValuesBeforeSending(t *testing.T) {
	handler := NewTestHandler(nil, nil, nil)
	for _, value := range []any{-1, 65536, int64(-1), int64(65536), -1.0, 65536.0, 1.5, math.NaN(), math.Inf(1)} {
		transport := &modbusOverrideTransport{}
		client := modbus.NewClient(transport, 1)
		for _, operation := range []string{"write_single_register", "write_multiple_registers"} {
			input := value
			if operation == "write_multiple_registers" {
				input = []interface{}{value}
			}
			err := handler.executeWrite(client, "modbus_rtu", WriteRequest{Operation: operation, Values: input})
			require.Error(t, err, "value=%v operation=%s", value, operation)
			require.Empty(t, transport.requests)
		}
	}
	for _, value := range []any{0, 65535, int64(65535), 65535.0} {
		actual, err := testRegisterValue(value)
		require.NoError(t, err)
		if value == 0 {
			require.Zero(t, actual)
		} else {
			require.Equal(t, uint16(65535), actual)
		}
	}
}

func TestClientFactoryValidatesUnitIDBeforeCreatingTransport(t *testing.T) {
	handler := NewTestHandler(nil, nil, nil)
	for _, protocol := range []string{"modbus_tcp", "modbus_udp", "modbus_rtu"} {
		for _, value := range []any{-1, 256, 1.5, math.NaN(), "7"} {
			client, err := handler.createClientWithDebug(protocol, map[string]interface{}{"unitID": value}, "unit-validation")
			require.Error(t, err)
			require.Nil(t, client)
		}
	}
	for _, tc := range []struct {
		config map[string]interface{}
		want   byte
	}{
		{nil, 1}, {map[string]interface{}{"unitID": 0}, 0}, {map[string]interface{}{"unitID": 255.0}, 255},
	} {
		actual, err := testModbusUnitID(tc.config)
		require.NoError(t, err)
		require.Equal(t, tc.want, actual)
	}
}

func TestMalformedStoredConnectorConfigDoesNotExposePartialData(t *testing.T) {
	response := toDatabaseConnectorResponse(&schema.DatabaseConnector{
		ConnectionConfig: `{"host":"partial-host","password":"test-only", "broken": }`,
	})
	require.Empty(t, response.ConnectionConfig)
}

func TestMalformedSkippedAddressesDoNotExposePartialList(t *testing.T) {
	response := mapSourceRuleResponse(&schema.SourceRule{SkippedAddresses: `["D1", 2]`})
	require.Empty(t, response.SkippedAddresses)
}
