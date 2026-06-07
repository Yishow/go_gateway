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

func TestStudioV2WorkspaceDevicesHandler_CreateReturnsValidationErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewStudioV2WorkspaceDevicesHandler(
		workspace.NewService(workspace.NewMemoryRepository()),
		device.NewService(device.NewMemoryRepository(), nil),
	)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", strings.NewReader(`{
		"id": "dev-A",
		"name": "",
		"protocol": "modbus_tcp",
		"connection_config": {
			"host": "",
			"port": 502,
			"slave_id": 1,
			"timeout": 5
		}
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.Create(c)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d body=%s", resp.Code, resp.Body.String())
	}
	if body := resp.Body.String(); !strings.Contains(body, "設備名稱不能為空") {
		t.Fatalf("expected actionable validation message, got %s", body)
	}
}

func TestStudioV2WorkspaceDevicesHandler_UpdateReturnsAppliedForRunningDevice(t *testing.T) {
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
		"name": "Line A Saved"
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
	if data["runtime_apply_status"] != "applied" {
		t.Fatalf("expected runtime_apply_status=applied, got %#v", data["runtime_apply_status"])
	}
	if syncer.upserted == nil || syncer.upserted.ID != "dev-active" {
		t.Fatalf("expected runtime upsert for dev-active, got %+v", syncer.upserted)
	}
}

func TestStudioV2WorkspaceDevicesHandler_UpdateReturnsApplyFailedWhenRuntimeApplyFails(t *testing.T) {
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

	syncer := &stubDeviceRuntimeSyncer{upsertErr: context.DeadlineExceeded}
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc, syncer)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-active", strings.NewReader(`{
		"name": "Line A Saved"
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
	if data["runtime_apply_status"] != "apply_failed" {
		t.Fatalf("expected runtime_apply_status=apply_failed, got %#v", data["runtime_apply_status"])
	}
	if data["runtime_apply_message"] == "" {
		t.Fatalf("expected runtime_apply_message, got %#v", data["runtime_apply_message"])
	}
}

func TestStudioV2WorkspaceDevicesHandler_UpdateReturnsNotRunningForDormantDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	now := time.Now()
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-draft",
		Name:             "Line A PLC",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}); err != nil {
		t.Fatalf("create device failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-draft"); err != nil {
		t.Fatalf("attach device failed: %v", err)
	}

	syncer := &stubDeviceRuntimeSyncer{}
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc, syncer)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-draft", strings.NewReader(`{
		"name": "Line A Saved",
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
	c.Params = gin.Params{{Key: "id", Value: "dev-draft"}}
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
	if data["runtime_apply_status"] != "not_running" {
		t.Fatalf("expected runtime_apply_status=not_running, got %#v", data["runtime_apply_status"])
	}
	if syncer.upserted != nil {
		t.Fatalf("expected no runtime upsert for draft device, got %+v", syncer.upserted)
	}
}

func TestStudioV2WorkspaceDevicesHandler_UpdateAvailabilityStopsRunningDevice(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodPut, "/api/v1/datalink/studio-v2/workspace/devices/dev-active/availability", strings.NewReader(`{
		"availability_status": "unavailable",
		"availability_reason": "device form is invalid"
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{{Key: "id", Value: "dev-active"}}
	c.Request = req

	handler.UpdateAvailability(c)

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
	if data["availability_status"] != "unavailable" {
		t.Fatalf("expected availability_status=unavailable, got %#v", data["availability_status"])
	}
	if data["availability_reason"] != "device form is invalid" {
		t.Fatalf("expected availability_reason to be preserved, got %#v", data["availability_reason"])
	}
	if data["running"] != false {
		t.Fatalf("expected running=false, got %#v", data["running"])
	}
	if syncer.removed != "dev-active" {
		t.Fatalf("expected runtime remove for dev-active, got %s", syncer.removed)
	}
}
