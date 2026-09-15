package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func newStudioV2WorkspaceMeasurementRouter(t *testing.T) (
	routerResult *gin.Engine,
) {
	t.Helper()

	deviceRepo := device.NewMemoryRepository()
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-1",
		Name:             "Main Power Meter",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("seed dev-1 failed: %v", err)
	}

	deviceSvc := device.NewService(deviceRepo, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-1"); err != nil {
		t.Fatalf("attach dev-1 failed: %v", err)
	}

	measRepo := measurement.NewMemoryRepository()
	measSvc := measurement.NewService(measRepo)

	router := NewRouter(&DatalinkServices{
		Device:      deviceSvc,
		Workspace:   workspaceSvc,
		Measurement: measSvc,
	})

	return router
}

func TestStudioV2WorkspaceMeasurements_CRUD(t *testing.T) {
	router := newStudioV2WorkspaceMeasurementRouter(t)

	// 1. Create measurement
	createBody := `{
		"id": "meas-v-a",
		"device_id": "dev-1",
		"point_id": "pt-1",
		"name": "Phase A Voltage",
		"quantity": "voltage",
		"unit": "V",
		"semantic_kind": "gauge"
	}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/measurements", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create failed with status %d: %s", w.Code, w.Body.String())
	}

	// 2. List measurements
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/measurements", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list failed with status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "meas-v-a") {
		t.Fatalf("expected meas-v-a in list response: %s", w.Body.String())
	}

	// 3. Update measurement (rename)
	updateBody := `{
		"device_id": "dev-1",
		"point_id": "pt-1",
		"name": "Line Voltage Va",
		"quantity": "voltage",
		"unit": "V",
		"semantic_kind": "gauge"
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/api/v1/datalink/studio-v2/workspace/measurements/meas-v-a", strings.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("update failed with status %d: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Line Voltage Va") {
		t.Fatalf("expected updated name: %s", w.Body.String())
	}

	// 4. Delete measurement
	req = httptest.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/measurements/meas-v-a", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("delete failed with status %d: %s", w.Code, w.Body.String())
	}
}

func TestStudioV2WorkspaceMeasurements_TemplateFlow(t *testing.T) {
	router := newStudioV2WorkspaceMeasurementRouter(t)

	// 1. List templates
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/datalink/studio-v2/workspace/measurements/templates", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("list templates failed: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "three_phase_power_meter_v1") {
		t.Fatalf("expected 3-phase meter template in response: %s", w.Body.String())
	}

	// 2. Preview template
	previewBody := `{
		"template_id": "three_phase_power_meter_v1",
		"device_id": "dev-1",
		"base_address": "40001"
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/measurements/templates/preview", strings.NewReader(previewBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("preview template failed: %d %s", w.Code, w.Body.String())
	}

	// 3. Apply template with confirmed=true
	applyBody := `{
		"template_id": "three_phase_power_meter_v1",
		"device_id": "dev-1",
		"confirmed": true,
		"definitions": [
			{
				"id": "meas-dev-1-kwh",
				"device_id": "dev-1",
				"point_id": "pt-kwh",
				"name": "Active Energy",
				"quantity": "energy",
				"unit": "kWh",
				"semantic_kind": "counter"
			}
		]
	}`
	req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/measurements/templates/apply", strings.NewReader(applyBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("apply template failed: %d %s", w.Code, w.Body.String())
	}
}

func TestStudioV2WorkspaceMeasurements_DeviceNotBelongToWorkspace(t *testing.T) {
	router := newStudioV2WorkspaceMeasurementRouter(t)

	// dev-alien does not belong to workspace
	createBody := `{
		"id": "meas-alien",
		"device_id": "dev-alien",
		"point_id": "pt-1",
		"name": "Alien Sensor",
		"quantity": "temperature",
		"semantic_kind": "gauge"
	}`
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/measurements", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected validation rejection for unattached device, got %d: %s", w.Code, w.Body.String())
	}
}
