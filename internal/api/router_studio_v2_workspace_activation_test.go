package api

import (
	"net/http"
	"testing"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

func TestNewRouter_StudioV2WorkspaceActivateEndpointReturnsPartialResults(t *testing.T) {
	bank := memory.NewMemoryBank(64)
	if err := bank.WriteWord(0, 123); err != nil {
		t.Fatalf("seed modbus bank failed: %v", err)
	}

	server := virtualmodbus.NewServer(bank)
	if err := server.Start(0); err != nil {
		t.Fatalf("start modbus server failed: %v", err)
	}
	defer server.Stop()

	deviceSvc := device.NewService(
		device.NewMemoryRepository(),
		connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()),
	)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, nil, nil, nil)

	router := NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        point.NewService(point.NewMemoryRepository(), nil),
		Tag:          tag.NewService(tag.NewMemoryRepository()),
		Mapping:      mapping.NewService(mapping.NewMemoryRepository()),
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		Workspace:    workspaceSvc,
	})

	createA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "dev-A",
		"name":     "Line A PLC",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "127.0.0.1",
			"port":     server.Port(),
			"slave_id": 1,
			"timeout":  2,
		},
	})
	if createA.Code != http.StatusCreated {
		t.Fatalf("expected createA 201, got %d body=%s", createA.Code, createA.Body.String())
	}
	testA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/devices/dev-A/test", nil)
	if testA.Code != http.StatusOK {
		t.Fatalf("expected testA 200, got %d body=%s", testA.Code, testA.Body.String())
	}

	createB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/devices", map[string]any{
		"id":       "dev-B",
		"name":     "Line B PLC",
		"protocol": "modbus_tcp",
		"connection_config": map[string]any{
			"host":     "127.0.0.1",
			"port":     server.Port(),
			"slave_id": 1,
			"timeout":  2,
		},
	})
	if createB.Code != http.StatusCreated {
		t.Fatalf("expected createB 201, got %d body=%s", createB.Code, createB.Body.String())
	}
	testB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/devices/dev-B/test", nil)
	if testB.Code != http.StatusOK {
		t.Fatalf("expected testB 200, got %d body=%s", testB.Code, testB.Body.String())
	}

	activateResp := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/activate", nil)
	if activateResp.Code != http.StatusOK {
		t.Fatalf("expected activate 200, got %d body=%s", activateResp.Code, activateResp.Body.String())
	}

	body := decodeJSONBody(t, activateResp)
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
	results, ok := data["results"].([]any)
	if !ok || len(results) != 2 {
		t.Fatalf("expected two per-device results, got %#v", data["results"])
	}
}
