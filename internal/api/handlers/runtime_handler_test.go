package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func TestRuntimeHandler_WorkspaceContextReturnsOrderedDevicesAndDefaultDeviceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now().UTC()

	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "device-B",
		Name:             "Unavailable Filler",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ReadinessStatus:  `{"availability_status":"unavailable","availability_reason":"invalid Step 1 configuration"}`,
		ConnectionConfig: `{"host":"192.168.1.11","port":502,"slave_id":2,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device-B failed: %v", err)
	}
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "device-A",
		Name:             "Healthy Mixer",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.1.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device-A failed: %v", err)
	}

	if _, err := workspaceSvc.AttachDevice(context.Background(), "device-B"); err != nil {
		t.Fatalf("attach device-B failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "device-A"); err != nil {
		t.Fatalf("attach device-A failed: %v", err)
	}

	handler := NewRuntimeHandler(deviceSvc, nil, nil, nil, nil, workspaceSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.WorkspaceContext(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	var body struct {
		Success bool `json:"success"`
		Data    struct {
			WorkspaceID     string  `json:"workspace_id"`
			DefaultDeviceID *string `json:"default_device_id"`
			Devices         []struct {
				DeviceID           string  `json:"device_id"`
				Name               string  `json:"name"`
				Running            bool    `json:"running"`
				AvailabilityStatus string  `json:"availability_status"`
				AvailabilityReason *string `json:"availability_reason"`
			} `json:"devices"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}

	if !body.Success {
		t.Fatalf("expected success response, got %s", resp.Body.String())
	}
	if body.Data.WorkspaceID == "" {
		t.Fatalf("expected workspace id, got %s", resp.Body.String())
	}
	if len(body.Data.Devices) != 2 {
		t.Fatalf("expected two devices, got %+v", body.Data.Devices)
	}
	if body.Data.Devices[0].DeviceID != "device-B" || body.Data.Devices[1].DeviceID != "device-A" {
		t.Fatalf("expected ordered devices [device-B device-A], got %+v", body.Data.Devices)
	}
	if body.Data.Devices[0].AvailabilityStatus != "unavailable" {
		t.Fatalf("expected first device unavailable, got %+v", body.Data.Devices[0])
	}
	if body.Data.Devices[0].AvailabilityReason == nil || *body.Data.Devices[0].AvailabilityReason != "invalid Step 1 configuration" {
		t.Fatalf("expected unavailable reason, got %+v", body.Data.Devices[0])
	}
	if body.Data.DefaultDeviceID == nil || *body.Data.DefaultDeviceID != "device-A" {
		t.Fatalf("expected default device-A, got %+v", body.Data.DefaultDeviceID)
	}
}
