package api

import (
	"context"
	"net/http"
	"strings"
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

func newStudioV2WorkspaceSourceRuleRouter(t *testing.T) *gin.Engine {
	t.Helper()

	deviceRepo := device.NewMemoryRepository()
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-A",
		Name:             "dev-A",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("seed dev-A failed: %v", err)
	}
	if err := deviceRepo.Create(context.Background(), &schema.Device{
		ID:               "dev-B",
		Name:             "dev-B",
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusDraft,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":2,"timeout":5}`,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}); err != nil {
		t.Fatalf("seed dev-B failed: %v", err)
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
	sourceRuleSvc.SetTagMappingServices(tagSvc, mappingSvc)

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

func TestNewRouter_StudioV2WorkspaceSourceRuleAutosaveEndpoints(t *testing.T) {
	router := newStudioV2WorkspaceSourceRuleRouter(t)

	createA := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", map[string]any{
		"id":            "rule-A",
		"device_id":     "dev-A",
		"start_address": "40001",
		"count":         1,
		"data_type":     "int16",
		"naming_prefix": "A_",
		"enabled":       true,
	})
	if createA.Code != http.StatusCreated {
		t.Fatalf("expected createA 201, got %d body=%s", createA.Code, createA.Body.String())
	}
	if runtimeStatus := decodeJSONBody(t, createA)["data"].(map[string]any)["runtime_apply_status"]; runtimeStatus != "not_running" {
		t.Fatalf("expected createA runtime_apply_status=not_running, got %#v", runtimeStatus)
	}

	createB := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", map[string]any{
		"id":            "rule-B",
		"device_id":     "dev-B",
		"start_address": "40011",
		"count":         1,
		"data_type":     "int16",
		"naming_prefix": "B_",
		"enabled":       true,
	})
	if createB.Code != http.StatusCreated {
		t.Fatalf("expected createB 201, got %d body=%s", createB.Code, createB.Body.String())
	}
	if runtimeStatus := decodeJSONBody(t, createB)["data"].(map[string]any)["runtime_apply_status"]; runtimeStatus != "not_running" {
		t.Fatalf("expected createB runtime_apply_status=not_running, got %#v", runtimeStatus)
	}

	listResp := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/source-rules", nil)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d body=%s", listResp.Code, listResp.Body.String())
	}

	listBody := decodeJSONBody(t, listResp)
	items, ok := listBody["data"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("expected 2 rules, got %#v", listBody["data"])
	}

	firstRule := items[0].(map[string]any)
	secondRule := items[1].(map[string]any)
	if firstRule["device_id"] != "dev-B" || secondRule["device_id"] != "dev-A" {
		t.Fatalf("expected workspace device order [dev-B dev-A], got %#v", listBody["data"])
	}
	if firstRule["workspace_id"] == "" || firstRule["save_state"] != "saved" {
		t.Fatalf("expected workspace-scoped response fields, got %#v", firstRule)
	}

	invalidUpdate := performJSONRequest(t, router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", map[string]any{
		"device_id":     "dev-A",
		"start_address": "",
		"count":         1,
		"data_type":     "int16",
		"naming_prefix": "A_",
	})
	if invalidUpdate.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid update 400, got %d body=%s", invalidUpdate.Code, invalidUpdate.Body.String())
	}

	getListAfterInvalid := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace/source-rules", nil)
	items = decodeJSONBody(t, getListAfterInvalid)["data"].([]any)
	ruleA := items[1].(map[string]any)
	if ruleA["start_address"] != "40001" {
		t.Fatalf("expected persisted start_address to stay unchanged, got %#v", ruleA["start_address"])
	}

	deleteResp := performJSONRequest(t, router, http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", nil)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected delete 200, got %d body=%s", deleteResp.Code, deleteResp.Body.String())
	}
}

func TestNewRouter_StudioV2WorkspaceSourceRuleUpdateReturnsValidationWhenDeviceIsNotReady(t *testing.T) {
	router := newStudioV2WorkspaceSourceRuleRouter(t)

	createRule := performJSONRequest(t, router, http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", map[string]any{
		"id":            "rule-A",
		"device_id":     "dev-A",
		"start_address": "40001",
		"count":         1,
		"data_type":     "int16",
		"naming_prefix": "A_",
		"enabled":       true,
	})
	if createRule.Code != http.StatusCreated {
		t.Fatalf("expected createRule 201, got %d body=%s", createRule.Code, createRule.Body.String())
	}

	updateRule := performJSONRequest(t, router, http.MethodPut, "/api/v1/datalink/studio-v2/workspace/source-rules/rule-A", map[string]any{
		"device_id":     "dev-A",
		"start_address": "40001",
		"count":         1,
		"data_type":     "int32",
		"naming_prefix": "A_",
		"enabled":       true,
	})
	if updateRule.Code != http.StatusBadRequest {
		t.Fatalf("expected updateRule 400, got %d body=%s", updateRule.Code, updateRule.Body.String())
	}
	if !strings.Contains(updateRule.Body.String(), "activation readiness") {
		t.Fatalf("expected validation message to mention activation readiness, got %s", updateRule.Body.String())
	}
}
