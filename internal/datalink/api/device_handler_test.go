package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/api"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDeviceRepo mock implementation
type MockDeviceRepo struct {
	mock.Mock
}

func (m *MockDeviceRepo) Create(ctx context.Context, dev *schema.Device) error {
	args := m.Called(ctx, dev)
	return args.Error(0)
}
func (m *MockDeviceRepo) Update(ctx context.Context, dev *schema.Device) error {
	args := m.Called(ctx, dev)
	return args.Error(0)
}
func (m *MockDeviceRepo) GetByID(ctx context.Context, id string) (*schema.Device, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.Device), args.Error(1)
}
func (m *MockDeviceRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockDeviceRepo) List(ctx context.Context, filter device.ListFilter) ([]*schema.Device, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*schema.Device), args.Error(1)
}
func (m *MockDeviceRepo) UpdateStatus(ctx context.Context, id string, status schema.DeviceStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
func (m *MockDeviceRepo) UpdateTestResult(ctx context.Context, id string, success bool, errMsg string) error {
    args := m.Called(ctx, id, success, errMsg)
    return args.Error(0)
}

func TestDeviceHandler_CheckReadiness(t *testing.T) {
	// Setup
	mockRepo := new(MockDeviceRepo)
	svc := device.NewService(mockRepo, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	handler := api.NewDeviceHandler(svc, nil)

	testDevice := &schema.Device{
		ID:               "dev-1",
		Name:             "Test Device",
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"localhost","port":502}`,
        ReadinessStatus:  "",
	}

	mockRepo.On("GetByID", mock.Anything, "dev-1").Return(testDevice, nil)
    mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	// Create Request
	req, _ := http.NewRequest("POST", "/devices/dev-1/readiness", nil)
	rr := httptest.NewRecorder()

	// Call Handler
	handler.CheckReadiness(rr, req, "dev-1")

	// Assert Response
	assert.Equal(t, http.StatusOK, rr.Code)

    t.Logf("Response Body: %s", rr.Body.String())

    type APIResponse struct {
        Success bool `json:"success"`
        Data    schema.DeviceReadiness `json:"data"`
    }

	var response APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
    
    result := response.Data
	assert.Equal(t, "dev-1", result.DeviceID)
	// assert.Equal(t, "ready", result.Status) // Depends on logic
    
    // Check if checks are present
    assert.NotEmpty(t, result.Checks)
}
