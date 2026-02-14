package device

import (
	"context"
	"testing"

	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCreateAndList(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil)
	ctx := context.Background()

	req := CreateDeviceRequest{
		Name:        "PLC-1",
		Protocol:    schema.ProtocolModbusTCP,
		Description: "test device",
		ConnectionConfig: map[string]interface{}{
			"host":     "127.0.0.1",
			"port":     float64(502),
			"slave_id": float64(1),
			"timeout":  float64(5),
		},
	}

	device, err := svc.Create(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, schema.DeviceStatusDraft, device.Status)
	assert.NotEmpty(t, device.ConnectionConfig)

	devices, err := svc.List(ctx, ListFilter{})
	require.NoError(t, err)
	assert.Len(t, devices, 1)
	assert.Equal(t, device.ID, devices[0].ID)
}

func TestValidateConnectionConfigUnknownProtocol(t *testing.T) {
	err := validateConnectionConfig(schema.ProtocolType("custom"), map[string]interface{}{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "未知的協議類型")
}

func TestParseConnectionConfig(t *testing.T) {
	device := &schema.Device{
		Protocol: schema.ProtocolModbusRTU,
		ConnectionConfig: `{
			"serial_port": "COM1",
			"baud_rate": 9600,
			"data_bits": 8,
			"stop_bits": 1,
			"parity": "none",
			"slave_id": 1,
			"timeout": 5
		}`,
	}

	cfg, err := ParseConnectionConfig(device)
	require.NoError(t, err)
	parsed, ok := cfg.(schema.ConnectionConfigModbusRTU)
	require.True(t, ok)
	assert.Equal(t, "COM1", parsed.SerialPort)
	assert.Equal(t, 9600, parsed.BaudRate)
}
