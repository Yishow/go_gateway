package handlers

import (
	"testing"

	"go-gateway/internal/protocol/modbus"
)

type modbusOverrideTransport struct {
	requests [][]byte
}

func (t *modbusOverrideTransport) Connect() error {
	return nil
}

func (t *modbusOverrideTransport) Close() error {
	return nil
}

func (t *modbusOverrideTransport) SendReceive(data []byte) ([]byte, error) {
	t.requests = append(t.requests, append([]byte(nil), data...))
	return modbus.BuildRTUFrame(data[0], modbus.FuncReadHoldingRegisters, []byte{2, 0x00, 0x01}), nil
}

func (t *modbusOverrideTransport) GetOriginalTransport() interface{} {
	return &modbus.RTUTransport{}
}

func TestPrepareOverrideClient_ReusesModbusTransportForUnitID(t *testing.T) {
	handler := NewTestHandler(nil, nil, nil)
	transport := &modbusOverrideTransport{}
	state := &ConnectionState{
		ID:       "conn-1",
		Protocol: "modbus_rtu",
		Config: map[string]interface{}{
			"port": "COM1",
		},
		Client: modbus.NewClient(transport, 1),
	}

	unitID := byte(7)
	clientToUse, tempClient, err := handler.prepareOverrideClient(state, &unitID, nil, "_temp")
	if err != nil {
		t.Fatalf("expected override client without error, got %v", err)
	}
	if tempClient != nil {
		t.Fatalf("expected no temporary client for modbus unit override, got %#v", tempClient)
	}

	result, err := handler.executeRead(clientToUse, state.Protocol, ReadRequest{
		ConnectionID: state.ID,
		Operation:    "read_holding_registers",
		Address:      0,
		Count:        1,
	})
	if err != nil {
		t.Fatalf("expected read to succeed via shared transport, got %v", err)
	}

	registers, ok := result.([]uint16)
	if !ok || len(registers) != 1 || registers[0] != 1 {
		t.Fatalf("expected one holding register with value 1, got %#v", result)
	}

	if len(transport.requests) != 1 {
		t.Fatalf("expected exactly one request on shared transport, got %d", len(transport.requests))
	}

	if gotUnitID := transport.requests[0][0]; gotUnitID != unitID {
		t.Fatalf("expected request unit id %d, got %d", unitID, gotUnitID)
	}
}
