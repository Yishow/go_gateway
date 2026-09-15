package modbus

import (
	"errors"
	"testing"
)

func TestClient_WriteMultipleQuantityBounds(t *testing.T) {
	transport := NewMockTransport()
	client := NewClient(transport, 1)

	if err := client.WriteMultipleCoils(0, make([]bool, 1968)); err != nil {
		t.Fatalf("1968 coils should be accepted: %v", err)
	}
	before := len(transport.sentData)
	if err := client.WriteMultipleCoils(0, make([]bool, 1969)); !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("1969 coils should return ErrInvalidQuantity, got %v", err)
	}
	if got := len(transport.sentData); got != before {
		t.Fatalf("1969 coils sent %d requests, want %d", got-before, 0)
	}

	if err := client.WriteMultipleRegisters(0, make([]uint16, 123)); err != nil {
		t.Fatalf("123 registers should be accepted: %v", err)
	}
	before = len(transport.sentData)
	if err := client.WriteMultipleRegisters(0, make([]uint16, 124)); !errors.Is(err, ErrInvalidQuantity) {
		t.Fatalf("124 registers should return ErrInvalidQuantity, got %v", err)
	}
	if got := len(transport.sentData); got != before {
		t.Fatalf("124 registers sent %d requests, want %d", got-before, 0)
	}
}
