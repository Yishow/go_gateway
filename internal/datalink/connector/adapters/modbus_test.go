package adapters

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/modbus"

	"github.com/stretchr/testify/assert"
)

func TestModbusTCPConnector_Read(t *testing.T) {
	mockTrans := NewMockTransport()
	// Pre-fill data
	mockTrans.HoldingRegs[100] = 12345
	mockTrans.HoldingRegs[101] = 6789

	// Manually construct connector
	client := modbus.NewClient(mockTrans, 1)
	conn := &ModbusTCPConnector{
		client:    client,
		connected: true,
	}

	ctx := context.Background()

	// Test Read Holding Register (Uint16)
	req := connector.ReadRequest{
		Address:  "40101", // HR 101 -> modbus addr 100
		DataType: schema.DataTypeUint16,
		Count:    1,
	}
	res, err := conn.Read(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, schema.QualityGood, res.Quality)
	assert.Equal(t, uint16(12345), res.Value)

	// Test Read Int16
	mockTrans.HoldingRegs[102] = uint16(0xFFFF) // -1
	req = connector.ReadRequest{
		Address:  "40103", // HR 103 -> modbus addr 102
		DataType: schema.DataTypeInt16,
		Count:    1,
	}
	res, err = conn.Read(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, int16(-1), res.Value)
}

func TestModbusTCPConnector_Write(t *testing.T) {
	mockTrans := NewMockTransport()
	client := modbus.NewClient(mockTrans, 1)
	conn := &ModbusTCPConnector{
		client:    client,
		connected: true,
	}

	ctx := context.Background()

	// Test Write Single Register
	req := connector.WriteRequest{
		Address: "40201", // HR 201 -> modbus addr 200
		Value:   123,
	}
	err := conn.Write(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, uint16(123), mockTrans.HoldingRegs[200])
}
