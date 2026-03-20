package modbus

import "testing"

type wrappedRTUTransport struct {
	sent [][]byte
}

func (w *wrappedRTUTransport) Connect() error {
	return nil
}

func (w *wrappedRTUTransport) Close() error {
	return nil
}

func (w *wrappedRTUTransport) SendReceive(data []byte) ([]byte, error) {
	w.sent = append(w.sent, append([]byte(nil), data...))
	return BuildRTUFrame(1, FuncReadHoldingRegisters, []byte{2, 0x30, 0x39}), nil
}

func (w *wrappedRTUTransport) GetNextTransactionID() uint16 {
	return 0
}

func (w *wrappedRTUTransport) GetOriginalTransport() interface{} {
	return &RTUTransport{}
}

func TestClient_ReadHoldingRegisters_WithWrappedRTUTransport(t *testing.T) {
	transport := &wrappedRTUTransport{}
	client := NewClient(transport, 1)

	registers, err := client.ReadHoldingRegisters(0, 1)
	if err != nil {
		t.Fatalf("expected wrapped RTU read to succeed, got error: %v", err)
	}

	if len(registers) != 1 || registers[0] != 12345 {
		t.Fatalf("expected one holding register with value 12345, got %#v", registers)
	}

	if len(transport.sent) != 1 {
		t.Fatalf("expected exactly one RTU request, got %d", len(transport.sent))
	}

	expectedFrame := BuildRTUFrame(1, FuncReadHoldingRegisters, BuildReadRequest(0, 1))
	if got := transport.sent[0]; string(got) != string(expectedFrame) {
		t.Fatalf("expected RTU frame %v, got %v", expectedFrame, got)
	}
}
