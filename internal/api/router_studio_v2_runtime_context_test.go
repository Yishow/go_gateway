package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

func TestNewRouter_StudioV2WorkspaceRuntimeContextEndpointReturnsOrderedPayload(t *testing.T) {
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	services := &DatalinkServices{
		Device:       deviceSvc,
		Point:        point.NewService(point.NewMemoryRepository(), nil),
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		Workspace:    workspace.NewService(workspace.NewMemoryRepository()),
	}
	router := NewRouter(services)

	createB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "device-B",
		"name":     "Unavailable Filler",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.11",
			"port":     502,
			"slave_id": 2,
			"timeout":  5,
		},
	})
	if createB.Code != http.StatusCreated {
		t.Fatalf("expected createB 201, got %d body=%s", createB.Code, createB.Body.String())
	}

	createA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "device-A",
		"name":     "Healthy Mixer",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.10",
			"port":     502,
			"slave_id": 1,
			"timeout":  5,
		},
	})
	if createA.Code != http.StatusCreated {
		t.Fatalf("expected createA 201, got %d body=%s", createA.Code, createA.Body.String())
	}

	savedB, err := services.Device.GetByID(context.Background(), "device-B")
	if err != nil {
		t.Fatalf("load device-B failed: %v", err)
	}
	savedB.ReadinessStatus = `{"availability_status":"unavailable","availability_reason":"invalid Step 1 configuration"}`
	savedB.UpdatedAt = time.Now().UTC()
	if err := deviceRepo.Update(context.Background(), savedB); err != nil {
		t.Fatalf("mark device-B unavailable failed: %v", err)
	}

	resp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected runtime context 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	body := decodeJSONBody(t, resp)
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
	devices, ok := data["devices"].([]any)
	if !ok || len(devices) != 2 {
		t.Fatalf("expected two runtime context devices, got %#v", data["devices"])
	}
	first, _ := devices[0].(map[string]any)
	second, _ := devices[1].(map[string]any)
	if first["device_id"] != "device-B" || second["device_id"] != "device-A" {
		t.Fatalf("expected ordered runtime devices [device-B device-A], got %#v", data["devices"])
	}
	if first["availability_status"] != "unavailable" {
		t.Fatalf("expected first device unavailable, got %#v", first)
	}
	if data["default_device_id"] != "device-A" {
		t.Fatalf("expected default_device_id=device-A, got %#v", data["default_device_id"])
	}
}

func TestNewRouter_StudioV2WorkspaceRuntimeContextEndpointReturnsEmptyDefaultWhenNoAvailableDevice(t *testing.T) {
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	services := &DatalinkServices{
		Device:       deviceSvc,
		Point:        point.NewService(point.NewMemoryRepository(), nil),
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		Workspace:    workspace.NewService(workspace.NewMemoryRepository()),
	}
	router := NewRouter(services)

	createOnly := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "device-B",
		"name":     "Unavailable Filler",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "192.168.10.11",
			"port":     502,
			"slave_id": 2,
			"timeout":  5,
		},
	})
	if createOnly.Code != http.StatusCreated {
		t.Fatalf("expected createOnly 201, got %d body=%s", createOnly.Code, createOnly.Body.String())
	}

	savedDevice, err := services.Device.GetByID(context.Background(), "device-B")
	if err != nil {
		t.Fatalf("load device-B failed: %v", err)
	}
	savedDevice.Status = schema.DeviceStatusActive
	savedDevice.ReadinessStatus = `{"availability_status":"unavailable","availability_reason":"invalid Step 1 configuration"}`
	savedDevice.UpdatedAt = time.Now().UTC()
	if err := deviceRepo.Update(context.Background(), savedDevice); err != nil {
		t.Fatalf("mark device-B unavailable failed: %v", err)
	}

	resp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/runtime-context", nil)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected runtime context 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	body := decodeJSONBody(t, resp)
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
	if _, exists := data["default_device_id"]; !exists {
		t.Fatalf("expected default_device_id field, got %#v", data)
	}
	if data["default_device_id"] != nil {
		t.Fatalf("expected default_device_id=nil, got %#v", data["default_device_id"])
	}
}
