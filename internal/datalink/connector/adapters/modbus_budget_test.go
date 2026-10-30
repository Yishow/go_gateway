package adapters

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"testing"
)

type modbusBudgetTransport struct{ sends int }

func (t *modbusBudgetTransport) Connect() error               { return nil }
func (t *modbusBudgetTransport) Close() error                 { return nil }
func (t *modbusBudgetTransport) IsConnected() bool            { return true }
func (t *modbusBudgetTransport) GetNextTransactionID() uint16 { return 1 }
func (t *modbusBudgetTransport) SendReceive([]byte) ([]byte, error) {
	t.sends++
	return []byte{0, 1, 0, 0, 0, 5, 1, 3, 2, 0, 7}, nil
}
func TestModbusCanceledReadDoesNotSend(t *testing.T) {
	tr := &modbusBudgetTransport{}
	c := &ModbusTCPConnector{client: modbus.NewClient(tr, 1), connected: true, persistentMode: true}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := c.Read(ctx, connector.ReadRequest{Address: "0", DataType: schema.DataTypeUint16})
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, tr.sends)
}

type modbusFailingBudgetTransport struct {
	modbusBudgetTransport
	failure error
}

func (t *modbusFailingBudgetTransport) SendReceive([]byte) ([]byte, error) {
	t.sends++
	return nil, t.failure
}
func TestModbusDelegatedReadsShareCallerBudget(t *testing.T) {
	for _, kind := range []string{"tcp", "rtu", "udp"} {
		t.Run(kind, func(t *testing.T) {
			tr := &modbusFailingBudgetTransport{failure: io.EOF}
			client := modbus.NewClient(tr, 1)
			var protocol connector.Protocol
			switch kind {
			case "tcp":
				protocol = &ModbusTCPConnector{client: client, connected: true, persistentMode: true}
			case "rtu":
				protocol = &ModbusRTUConnector{client: client, connected: true, persistentMode: true}
			case "udp":
				protocol = &ModbusUDPConnector{client: client, connected: true, persistentMode: true}
			}
			ctx := wirebudget.With(t.Context(), 3, 0)
			_, err := protocol.Read(ctx, connector.ReadRequest{Address: "0", DataType: schema.DataTypeUint16})
			require.Error(t, err)
			require.Equal(t, 3, tr.sends)
			require.Equal(t, 3, wirebudget.From(ctx).Used())
		})
	}
}
func TestModbusReadNeverRetriesTextualEOF(t *testing.T) {
	tr := &modbusFailingBudgetTransport{failure: errors.New("address eof is unsupported")}
	c := &ModbusTCPConnector{client: modbus.NewClient(tr, 1), connected: true, persistentMode: true}
	_, err := c.Read(wirebudget.With(t.Context(), 4, 0), connector.ReadRequest{Address: "0", DataType: schema.DataTypeUint16})
	require.Error(t, err)
	require.Equal(t, 1, tr.sends)
}
