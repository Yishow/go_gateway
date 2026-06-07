package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func TestStudioV2WorkspaceDevicesHandler_UpdateConnectionConfigReturnsRestartRequiredForRunningDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now()
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-active",
		Name:             "Line A PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-active"); err != nil {
		t.Fatalf("attach device failed: %v", err)
	}

	syncer := &stubDeviceRuntimeSyncer{}
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc, syncer)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-active", strings.NewReader(`{
		"connection_config": {
			"host": "192.168.10.20",
			"port": 502,
			"slave_id": 1,
			"timeout": 5
		}
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-active"}}
	c.Request = req

	handler.Update(c)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", resp.Code, resp.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response failed: %v body=%s", err, resp.Body.String())
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
	if data["runtime_apply_status"] != "restart-required" {
		t.Fatalf("expected runtime_apply_status=restart-required, got %#v", data["runtime_apply_status"])
	}
	if syncer.upserted != nil {
		t.Fatalf("expected no hot runtime upsert for connection change, got %+v", syncer.upserted)
	}
}
