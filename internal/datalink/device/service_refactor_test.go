package device

import (
	"context"
	"fmt"
	"testing"

	"go-gateway/internal/datalink/connector"
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

func TestCheckReadinessReturnsErrorWhenUpdateFails(t *testing.T) {
	repo := &stubDeviceRepository{
		device: &schema.Device{
			ID:               "dev-1",
			Status:           schema.DeviceStatusActive,
			Protocol:         schema.ProtocolModbusTCP,
			ConnectionConfig: `{"host":"127.0.0.1"}`,
		},
		updateErr: fmt.Errorf("db update failed"),
	}
	svc := NewService(repo, nil)

	_, err := svc.CheckReadiness(context.Background(), "dev-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "更新設備就緒狀態失敗")
}

func TestTestConnectionReturnsErrorWhenResultUpdateFails(t *testing.T) {
	repo := &stubDeviceRepository{
		device: &schema.Device{
			ID:               "dev-1",
			Status:           schema.DeviceStatusDraft,
			Protocol:         schema.ProtocolType("unsupported"),
			ConnectionConfig: "{}",
		},
		updateTestResultErr: fmt.Errorf("update test result failed"),
	}
	svc := NewService(repo, connector.GetConnectionManager())

	err := svc.TestConnection(context.Background(), "dev-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "記錄測試結果失敗")
	assert.Equal(t, 1, repo.updateTestResultCalls)
}

type stubDeviceRepository struct {
	device                *schema.Device
	updateErr             error
	updateTestResultErr   error
	updateTestResultCalls int
}

func (r *stubDeviceRepository) Create(ctx context.Context, device *schema.Device) error {
	return nil
}

func (r *stubDeviceRepository) Update(ctx context.Context, device *schema.Device) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	deviceCopy := *device
	r.device = &deviceCopy
	return nil
}

func (r *stubDeviceRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *stubDeviceRepository) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	if r.device == nil {
		return nil, fmt.Errorf("device not found")
	}
	deviceCopy := *r.device
	return &deviceCopy, nil
}

func (r *stubDeviceRepository) List(ctx context.Context, filter ListFilter) ([]*schema.Device, error) {
	return nil, nil
}

func (r *stubDeviceRepository) UpdateTestResult(ctx context.Context, id string, success bool, errMsg string) error {
	r.updateTestResultCalls++
	if r.updateTestResultErr != nil {
		return r.updateTestResultErr
	}
	return nil
}

func (r *stubDeviceRepository) UpdateStatus(ctx context.Context, id string, status schema.DeviceStatus) error {
	return nil
}
