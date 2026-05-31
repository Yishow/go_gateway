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
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
)

func newStudioV2WorkspaceMappingRouter(t *testing.T) *gin.Engine {
	t.Helper()

	deviceRepo := device.NewMemoryRepository()
	for _, seed := range []struct {
		id       string
		slaveID  int
	}{
		{id: "dev-A", slaveID: 1},
		{id: "dev-B", slaveID: 2},
	} {
		if err := deviceRepo.Create(context.Background(), &schema.Device{
			ID:               seed.id,
			Name:             seed.id,
			Protocol:         schema.ProtocolModbusTCP,
			Status:           schema.DeviceStatusDraft,
			ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}); err != nil {
			t.Fatalf("seed %s failed: %v", seed.id, err)
		}
	}

	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	deviceSvc := device.NewService(deviceRepo, nil)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository())
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-B"); err != nil {
		t.Fatalf("attach dev-B failed: %v", err)
	}
	if _, err := workspaceSvc.AttachDevice(context.Background(), "dev-A"); err != nil {
		t.Fatalf("attach dev-A failed: %v", err)
	}

	sourceRuleSvc := sourcerule.NewService(sourcerule.NewMemoryRepository(), deviceSvc, pointSvc, nil)
	for _, createReq := range []sourcerule.CreateRuleRequest{
		{
			ID:           "rule-A",
			DeviceID:     "dev-A",
			StartAddress: "40001",
			Count:        1,
			DataType:     schema.DataTypeInt16,
			NamingPrefix: "A_",
			Enabled:      true,
		},
		{
			ID:           "rule-B",
			DeviceID:     "dev-B",
			StartAddress: "40011",
			Count:        1,
			DataType:     schema.DataTypeInt16,
			NamingPrefix: "B_",
			Enabled:      true,
		},
	} {
		if _, err := sourceRuleSvc.Create(context.Background(), createReq); err != nil {
			t.Fatalf("seed rule %s failed: %v", createReq.ID, err)
		}
	}

	return NewRouter(&DatalinkServices{
		Device:       deviceSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pollinggroup.NewService(pollinggroup.NewMemoryRepository()),
		Settings:     settings.NewService(settings.NewMemoryRepository()),
		SourceRule:   sourceRuleSvc,
		Workspace:    workspaceSvc,
	})
}

func TestNewRouter_StudioV2WorkspaceMappingAutosaveEndpoints(t *testing.T) {
	router := newStudioV2WorkspaceMappingRouter(t)

	createB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", map[string]any{
		"rule_id":      "rule-B",
		"address":      "40011",
		"tag_key":      "line.b.temp",
		"display_name": "Line B Temp",
		"unit":         "C",
		"target_type":  "float64",
		"scale":        1,
		"offset":       0,
		"enabled":      true,
	})
	if createB.Code != http.StatusCreated {
		t.Fatalf("expected createB 201, got %d body=%s", createB.Code, createB.Body.String())
	}
	if runtimeStatus := decodeJSONBody(t, createB)["data"].(map[string]any)["runtime_apply_status"]; runtimeStatus != "not_running" {
		t.Fatalf("expected createB runtime_apply_status=not_running, got %#v", runtimeStatus)
	}

	createA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/mappings", map[string]any{
		"rule_id":      "rule-A",
		"address":      "40001",
		"tag_key":      "line.a.temp",
		"display_name": "Line A Temp",
		"unit":         "C",
		"target_type":  "float64",
		"scale":        2,
		"offset":       5,
		"enabled":      true,
	})
	if createA.Code != http.StatusCreated {
		t.Fatalf("expected createA 201, got %d body=%s", createA.Code, createA.Body.String())
	}
	if runtimeStatus := decodeJSONBody(t, createA)["data"].(map[string]any)["runtime_apply_status"]; runtimeStatus != "not_running" {
		t.Fatalf("expected createA runtime_apply_status=not_running, got %#v", runtimeStatus)
	}

	listResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/mappings", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d body=%s", listResp.Code, listResp.Body.String())
	}

	items, ok := decodeJSONBody(t, listResp)["data"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("expected 2 mappings, got %#v", decodeJSONBody(t, listResp)["data"])
	}
	firstRow := items[0].(map[string]any)
	secondRow := items[1].(map[string]any)
	if firstRow["device_id"] != "dev-B" || secondRow["device_id"] != "dev-A" {
		t.Fatalf("expected workspace device order [dev-B dev-A], got %#v", items)
	}
	if firstRow["workspace_id"] == "" || firstRow["save_state"] != "saved" {
		t.Fatalf("expected workspace-scoped response fields, got %#v", firstRow)
	}

	mappingAID := decodeJSONBody(t, createA)["data"].(map[string]any)["id"].(string)
	invalidUpdate := performJSONRequest(t, router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/mappings/"+mappingAID, map[string]any{
		"tag_key":      "",
		"display_name": "Line A Temp",
		"unit":         "C",
		"target_type":  "float64",
		"scale":        2,
		"offset":       5,
		"enabled":      true,
	})
	if invalidUpdate.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid update 400, got %d body=%s", invalidUpdate.Code, invalidUpdate.Body.String())
	}

	listAfterInvalid := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/mappings", nil)
	items = decodeJSONBody(t, listAfterInvalid)["data"].([]any)
	rowA := items[1].(map[string]any)
	if rowA["tag_key"] != "line.a.temp" {
		t.Fatalf("expected persisted tag_key to stay unchanged, got %#v", rowA["tag_key"])
	}

	deleteResp := performJSONRequest(t, router, http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/mappings/"+mappingAID, nil)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected delete 200, got %d body=%s", deleteResp.Code, deleteResp.Body.String())
	}
}
