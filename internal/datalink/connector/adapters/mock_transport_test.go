package adapters

import (
	"encoding/binary"
	"fmt"
	"sync"

	"go-gateway/internal/protocol/modbus"
)

// MockTransport simulates a Modbus device
type MockTransport struct {
	mu            sync.Mutex
	HoldingRegs   map[uint16]uint16
	Coils         map[uint16]bool
	TransactionID uint16
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		HoldingRegs: make(map[uint16]uint16),
		Coils:       make(map[uint16]bool),
	}
}

func (m *MockTransport) Connect() error {
	return nil
}

func (m *MockTransport) Close() error {
	return nil
}

func (m *MockTransport) IsConnected() bool {
	return true
}

func (m *MockTransport) GetNextTransactionID() uint16 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TransactionID++
	return m.TransactionID
}

func (m *MockTransport) SendReceive(data []byte) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Parse request (Simulate TCP behavior)
	txID, unitID, funcCode, pduData, err := modbus.ParseTCPFrame(data)
	if err != nil {
		return nil, err
	}

	var responsePDU []byte

	switch funcCode {
	case modbus.FuncReadHoldingRegisters:
		if len(pduData) < 4 {
			return nil, modbus.ErrResponseTooShort
		}
		startAddr := binary.BigEndian.Uint16(pduData[0:])
		count := binary.BigEndian.Uint16(pduData[2:])

		// Byte count = count * 2
		byteCount := count * 2
		responseData := make([]byte, 1+byteCount)
		responseData[0] = byte(byteCount)

		for i := uint16(0); i < count; i++ {
			val := m.HoldingRegs[startAddr+i]
			binary.BigEndian.PutUint16(responseData[1+i*2:], val)
		}
		responsePDU = modbus.BuildPDU(funcCode, responseData)

	case modbus.FuncWriteSingleRegister:
		if len(pduData) < 4 {
			return nil, modbus.ErrResponseTooShort
		}
		addr := binary.BigEndian.Uint16(pduData[0:])
		val := binary.BigEndian.Uint16(pduData[2:])
		fmt.Printf("Mock Write: Addr=%d Val=%d\n", addr, val)
		m.HoldingRegs[addr] = val

		// Echo request as response
		responsePDU = modbus.BuildPDU(funcCode, pduData)

	default:
		return nil, fmt.Errorf("MockTransport: unsupported function code %d", funcCode)
	}

	// Build response frame using same transaction ID
	return modbus.BuildTCPFrame(txID, unitID, funcCode, responsePDU[1:]), nil
}
